package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendDriverActivity appends the binary representation of DriverActivity to dst.
func AppendDriverActivity(dst []byte, activity *cardv1.DriverActivity) ([]byte, error) {
	if activity == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetOldestDayRecordIndex()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetNewestDayRecordIndex()))

	var err error
	for _, rec := range activity.GetDailyRecords() {
		dst, err = AppendActivityDailyRecord(dst, rec)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendActivityDailyRecord appends a single daily activity record.
func AppendActivityDailyRecord(dst []byte, rec *cardv1.DriverActivity_DailyRecord) ([]byte, error) {
	// Marshal the record content to a temporary buffer to calculate its length.
	contentBuf := make([]byte, 0, 2048)
	contentBuf = appendTimeReal(contentBuf, rec.GetActivityRecordDate())
	// TODO: Append BCD daily presence counter
	contentBuf = append(contentBuf, 0, 0)
	contentBuf = binary.BigEndian.AppendUint16(contentBuf, uint16(rec.GetActivityDayDistance()))

	for _, ac := range rec.GetActivityChangeInfo() {
		var err error
		contentBuf, err = AppendActivityChange(contentBuf, ac)
		if err != nil {
			return nil, err
		}
	}

	recordLength := len(contentBuf) + 4 // +4 for the two length fields
	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetActivityPreviousRecordLength()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(recordLength))
	dst = append(dst, contentBuf...)

	return dst, nil
}

// AppendActivityChange appends a single 2-byte activity change info.
func AppendActivityChange(dst []byte, ac *cardv1.DriverActivity_DailyRecord_ActivityChange) ([]byte, error) {
	var aci uint16
	aci |= (uint16(ac.GetSlot()) & 0x1) << 15
	aci |= (uint16(ac.GetDrivingStatus()) & 0x1) << 14
	aci |= (uint16(ac.GetCardStatus()) & 0x1) << 13
	aci |= (uint16(ac.GetActivity()) & 0x3) << 11
	aci |= (uint16(ac.GetTimeOfChangeMinutes()) & 0x7FF)

	return binary.BigEndian.AppendUint16(dst, aci), nil
}
