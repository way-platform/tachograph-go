package tachograph

import (
	"bytes"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalEventsData parses the binary data for an EF_Events_Data record.
func UnmarshalEventsData(data []byte, ed *cardv1.EventData) error {
	const recordSize = 24
	r := bytes.NewReader(data)
	for r.Len() >= recordSize {
		recordData := make([]byte, recordSize)
		r.Read(recordData)

		// Check if the record is empty (padded)
		isEmpty := true
		for _, b := range recordData {
			if b != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		rec := &cardv1.EventData_Record{}
		if err := UnmarshalEventRecord(recordData, rec); err != nil {
			return err
		}
		records := ed.GetRecords()
		records = append(records, rec)
		ed.SetRecords(records)
	}
	return nil
}

// UnmarshalEventRecord parses a single 24-byte event record.
func UnmarshalEventRecord(data []byte, rec *cardv1.EventData_Record) error {
	r := bytes.NewReader(data)
	eventType, _ := r.ReadByte()
	rec.SetEventType(int32(eventType))
	rec.SetEventBeginTime(readTimeReal(r))
	rec.SetEventEndTime(readTimeReal(r))
	// TODO: Read BCD nation code
	r.ReadByte()
	rec.SetVehicleRegistrationNumber(readString(r, 14))
	return nil
}
