package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// AppendDriverActivity appends the binary representation of DriverActivityData to dst.
func AppendDriverActivity(dst []byte, activity *cardv1.DriverActivityData) ([]byte, error) {
	if activity == nil {
		return dst, nil
	}

	// Append header (pointers)
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetOldestDayRecordIndex()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetNewestDayRecordIndex()))

	// Append the records
	var err error
	for _, rec := range activity.GetDailyRecords() {
		if rec.GetValid() {
			// It's a parsed record, so we need to serialize it.
			dst, err = appendParsedDailyRecord(dst, rec)
			if err != nil {
				return nil, err
			}
		} else {
			// It's a raw record, just append the bytes.
			dst = append(dst, rec.GetRaw()...)
		}
	}
	return dst, nil
}

// appendParsedDailyRecord appends a single parsed daily activity record.
func appendParsedDailyRecord(dst []byte, rec *cardv1.DriverActivityData_DailyRecord) ([]byte, error) {
	// Check if this is an empty record (for roundtrip compatibility)
	isEmpty := rec.GetActivityDayDistance() == 0 &&
		len(rec.GetActivityChangeInfo()) == 0 &&
		(rec.GetActivityRecordDate() == nil || rec.GetActivityRecordDate().GetSeconds() == 0)

	if isEmpty && rec.GetActivityRecordLength() > 0 {
		// For empty records, use the original record length and write zeros
		originalLength := int(rec.GetActivityRecordLength())
		dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetActivityPreviousRecordLength()))
		dst = binary.BigEndian.AppendUint16(dst, uint16(originalLength))

		// Write the content as zeros (originalLength - 4 for the header we already wrote)
		contentLength := originalLength - 4
		if contentLength > 0 {
			dst = append(dst, make([]byte, contentLength)...)
		}
		return dst, nil
	}

	// Normal record processing - serialize content to temporary buffer first
	contentBuf := make([]byte, 0, 2048)

	// Activity record date (4 bytes BCD)
	contentBuf = appendDatef(contentBuf, rec.GetActivityRecordDate())

	// Activity daily presence counter (2 bytes BCD)
	contentBuf = binary.BigEndian.AppendUint16(contentBuf, uint16(rec.GetActivityDailyPresenceCounter()))

	// Activity day distance (2 bytes)
	contentBuf = binary.BigEndian.AppendUint16(contentBuf, uint16(rec.GetActivityDayDistance()))

	// Activity change info (2 bytes each)
	for _, ac := range rec.GetActivityChangeInfo() {
		var err error
		contentBuf, err = AppendActivityChange(contentBuf, ac)
		if err != nil {
			return nil, err
		}
	}

	// Calculate total record length (content + 4-byte header)
	recordLength := len(contentBuf) + 4

	// Append header with lengths
	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetActivityPreviousRecordLength()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(recordLength))

	// Append content
	dst = append(dst, contentBuf...)

	return dst, nil
}

// AppendActivityChange appends a single 2-byte activity change info.
func AppendActivityChange(dst []byte, ac *datadictionaryv1.ActivityChangeInfo) ([]byte, error) {
	var aci uint16

	// Reconstruct the bitfield from enum values
	slot := GetCardSlotNumber(ac.GetSlot(), 0)
	drivingStatus := GetDrivingStatus(ac.GetDrivingStatus(), 0)
	cardInserted := GetCardInserted(ac.GetInserted())
	activity := GetDriverActivityValue(ac.GetActivity(), 0)

	aci |= (uint16(slot) & 0x1) << 15
	aci |= (uint16(drivingStatus) & 0x1) << 14
	aci |= (uint16(cardInserted) & 0x1) << 13
	aci |= (uint16(activity) & 0x3) << 11
	aci |= (uint16(ac.GetTimeOfChangeMinutes()) & 0x7FF)

	return binary.BigEndian.AppendUint16(dst, aci), nil
}
