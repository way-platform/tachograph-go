package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalEventsData parses the binary data for an EF_Events_Data record.
func UnmarshalEventsData(data []byte, ed *cardv1.EventData) error {
	const recordSize = 24
	r := bytes.NewReader(data)
	var records []*cardv1.EventData_Record

	for r.Len() >= recordSize {
		recordData := make([]byte, recordSize)
		r.Read(recordData)

		// Check if this is a valid record by examining the event begin time (first 4 bytes after event type)
		// Event type is 1 byte, so event begin time starts at byte 1
		eventBeginTime := binary.BigEndian.Uint32(recordData[1:5])

		rec := &cardv1.EventData_Record{}

		if eventBeginTime == 0 {
			// Non-valid record: preserve original bytes
			rec.SetValid(false)
			rec.SetRawData(recordData)
		} else {
			// Valid record: parse semantic data
			rec.SetValid(true)
			if err := UnmarshalEventRecord(recordData, rec); err != nil {
				return err
			}
		}

		records = append(records, rec)
	}

	ed.SetRecords(records)
	return nil
}

// UnmarshalEventRecord parses a single 24-byte event record.
func UnmarshalEventRecord(data []byte, rec *cardv1.EventData_Record) error {
	r := bytes.NewReader(data)
	eventType, _ := r.ReadByte()

	// Convert raw event type to enum using protocol annotations
	SetEventFaultType(int32(eventType), rec.SetEventType, rec.SetUnrecognizedEventType)

	rec.SetEventBeginTime(readTimeReal(r))
	rec.SetEventEndTime(readTimeReal(r))

	// Read vehicle registration nation (1 byte)
	var nation byte
	binary.Read(r, binary.BigEndian, &nation)
	rec.SetVehicleRegistrationNation(fmt.Sprintf("%02X", nation))

	rec.SetVehicleRegistrationNumber(readString(r, 14))
	return nil
}
