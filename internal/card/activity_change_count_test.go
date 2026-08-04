package card

import (
	"encoding/binary"
	"testing"
)

// buildActivityDayRecord constructs the binary representation of a single
// CardActivityDailyRecord with the specified number of activity changes.
//
// Binary layout:
//
//	[prevRecordLength: 2] [currentRecordLength: 2]
//	[activityRecordDate (TimeReal): 4]
//	[activityDailyPresenceCounter (BCD): 2]
//	[activityDayDistance: 2]
//	[activityChangeInfo: numChanges * 2]
//
// Each activity change is encoded as slot=0, crew=0, card=0(inserted),
// activity=3(DRIVING), time=(i+1) to avoid the 0x0000 sentinel value.
func buildActivityDayRecord(numChanges int) []byte {
	const (
		headerLen  = 4 // prevRecordLength(2) + currentRecordLength(2)
		dateLen    = 4 // TimeReal
		counterLen = 2 // BCD DailyPresenceCounter
		distLen    = 2 // Distance
		changeLen  = 2 // per ActivityChangeInfo
	)
	fixedLen := headerLen + dateLen + counterLen + distLen
	totalLen := fixedLen + numChanges*changeLen

	buf := make([]byte, totalLen)
	offset := 0

	// Header
	binary.BigEndian.PutUint16(buf[offset:], 0)               // prevRecordLength = 0 (first record)
	binary.BigEndian.PutUint16(buf[offset+2:], uint16(totalLen)) // currentRecordLength
	offset += headerLen

	// TimeReal: 2024-01-01 00:00:00 UTC = 1704067200 = 0x65920080
	binary.BigEndian.PutUint32(buf[offset:], 0x65920080)
	offset += dateLen

	// BCD DailyPresenceCounter: "0001" = 0x00, 0x01
	buf[offset] = 0x00
	buf[offset+1] = 0x01
	offset += counterLen

	// Distance: 100 km
	binary.BigEndian.PutUint16(buf[offset:], 100)
	offset += distLen

	// Activity changes: each with a unique valid time (1..numChanges),
	// using activity=DRIVING(3) so the uint16 is never 0x0000 or 0xFFFF.
	for i := 0; i < numChanges; i++ {
		timeMinutes := uint16((i + 1) % 1440) // wrap to stay in valid range
		// bit layout: s(1)|c(1)|p(1)|aa(2)|t(11)
		var v uint16
		v |= (3 & 0x3) << 11       // activity = DRIVING
		v |= timeMinutes & 0x7FF
		binary.BigEndian.PutUint16(buf[offset:], v)
		offset += changeLen
	}

	return buf
}

// ---------------------------------------------------------------------------
// H-2: Activity changes per day must be in the range 1..1440 (DD 2.1).
// SET SIZE (1..1440) OF ActivityChangeInfo.
// ---------------------------------------------------------------------------

func TestActivityDay_CountWithinLimit(t *testing.T) {
	// 100 activity changes — well within the 1440 limit.
	data := buildActivityDayRecord(100)
	opts := UnmarshalOptions{}
	record, err := opts.parseSingleActivityDailyRecord(data)
	if err != nil {
		t.Fatalf("expected no error for 100 changes, got: %v", err)
	}
	got := len(record.GetActivityChangeInfo())
	if got != 100 {
		t.Errorf("change count = %d, want 100", got)
	}
}

func TestActivityDay_CountAtLimit(t *testing.T) {
	// 1440 activity changes — exactly at the boundary.
	data := buildActivityDayRecord(1440)
	opts := UnmarshalOptions{}
	record, err := opts.parseSingleActivityDailyRecord(data)
	if err != nil {
		t.Fatalf("expected no error for 1440 changes, got: %v", err)
	}
	got := len(record.GetActivityChangeInfo())
	if got != 1440 {
		t.Errorf("change count = %d, want 1440", got)
	}
}

func TestActivityDay_CountExceedsLimit(t *testing.T) {
	// 1441 activity changes — one past the maximum allowed by DD 2.1.
	data := buildActivityDayRecord(1441)
	opts := UnmarshalOptions{}
	_, err := opts.parseSingleActivityDailyRecord(data)
	if err == nil {
		t.Fatal("expected error for 1441 changes, got nil")
	}
}
