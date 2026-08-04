package dd

import (
	"encoding/binary"
	"testing"
)

// encodeActivityChangeInfo builds the 2-byte big-endian representation of an
// ActivityChangeInfo value from its constituent bit fields.
//
// Bit layout (16 bits): s(1) | c(1) | p(1) | aa(2) | ttttttttttt(11)
//   - s:   slot        (0 = DRIVER, 1 = CO_DRIVER)
//   - c:   crew        (0 = SINGLE, 1 = CREW)
//   - p:   card status (0 = INSERTED, 1 = NOT_INSERTED)
//   - aa:  activity    (0 = BREAK_REST, 1 = AVAILABILITY, 2 = WORK, 3 = DRIVING)
//   - t:   time of change in minutes since 00:00
func encodeActivityChangeInfo(slot, crew, card, activity uint16, timeMinutes uint16) []byte {
	var v uint16
	v |= (slot & 0x1) << 15
	v |= (crew & 0x1) << 14
	v |= (card & 0x1) << 13
	v |= (activity & 0x3) << 11
	v |= timeMinutes & 0x7FF
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf
}

// ---------------------------------------------------------------------------
// H-1: TimeOfChangeMinutes must be in the range 0..1439 (DD 2.1).
// The 11-bit field can hold 0..2047; values 1440..2047 are invalid.
// ---------------------------------------------------------------------------

func TestActivityChangeInfo_TimeZero(t *testing.T) {
	// time=0 (00:00) — valid lower bound
	input := encodeActivityChangeInfo(0, 0, 0, 0, 0) // 0x0000
	opts := UnmarshalOptions{}
	aci, err := opts.UnmarshalActivityChangeInfo(input)
	if err != nil {
		t.Fatalf("expected no error for time=0, got: %v", err)
	}
	if got := aci.GetTimeOfChangeMinutes(); got != 0 {
		t.Errorf("time = %d, want 0", got)
	}
}

func TestActivityChangeInfo_TimeMidDay(t *testing.T) {
	// time=720 (12:00) — valid middle value
	input := encodeActivityChangeInfo(0, 0, 0, 0, 720)
	opts := UnmarshalOptions{}
	aci, err := opts.UnmarshalActivityChangeInfo(input)
	if err != nil {
		t.Fatalf("expected no error for time=720, got: %v", err)
	}
	if got := aci.GetTimeOfChangeMinutes(); got != 720 {
		t.Errorf("time = %d, want 720", got)
	}
}

func TestActivityChangeInfo_TimeEndOfDay(t *testing.T) {
	// time=1439 (23:59) — valid upper bound
	input := encodeActivityChangeInfo(0, 0, 0, 0, 1439)
	opts := UnmarshalOptions{}
	aci, err := opts.UnmarshalActivityChangeInfo(input)
	if err != nil {
		t.Fatalf("expected no error for time=1439, got: %v", err)
	}
	if got := aci.GetTimeOfChangeMinutes(); got != 1439 {
		t.Errorf("time = %d, want 1439", got)
	}
}

func TestActivityChangeInfo_TimeOutOfRange_1440(t *testing.T) {
	// time=1440 — first invalid value (one past max)
	input := encodeActivityChangeInfo(0, 0, 0, 0, 1440)
	opts := UnmarshalOptions{}
	_, err := opts.UnmarshalActivityChangeInfo(input)
	if err == nil {
		t.Fatal("expected error for time=1440, got nil")
	}
}

func TestActivityChangeInfo_TimeOutOfRange_2047(t *testing.T) {
	// time=2047 — maximum 11-bit value, far out of range
	input := encodeActivityChangeInfo(0, 0, 0, 0, 2047)
	opts := UnmarshalOptions{}
	_, err := opts.UnmarshalActivityChangeInfo(input)
	if err == nil {
		t.Fatal("expected error for time=2047, got nil")
	}
}
