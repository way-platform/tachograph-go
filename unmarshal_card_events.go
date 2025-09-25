package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func unmarshalEventsData(data []byte) (*cardv1.EventData, error) {
	const recordSize = 24
	r := bytes.NewReader(data)
	var records []*cardv1.EventData_Record
	for r.Len() >= recordSize {
		recordData := make([]byte, recordSize)
		r.Read(recordData)
		// Check if this is a valid record by examining the event begin time (first 4 bytes after event type)
		// Event type is 1 byte, so event begin time starts at byte 1
		eventBeginTime := binary.BigEndian.Uint32(recordData[1:5])
		if eventBeginTime == 0 {
			// Non-valid record: preserve original bytes
			rec := &cardv1.EventData_Record{}
			rec.SetValid(false)
			rec.SetRawData(recordData)
			records = append(records, rec)
		} else {
			// Valid record: parse semantic data
			rec, err := unmarshalEventRecord(recordData)
			if err != nil {
				return nil, err
			}
			rec.SetValid(true)
			records = append(records, rec)
		}
	}
	var ed cardv1.EventData
	ed.SetRecords(records)
	return &ed, nil
}

// unmarshalEventRecord parses a single 24-byte event record.
func unmarshalEventRecord(data []byte) (*cardv1.EventData_Record, error) {
	var rec cardv1.EventData_Record
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
	return &rec, nil
}
