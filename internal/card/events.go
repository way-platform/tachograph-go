package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalEventsData unmarshals events data from a card EF.
//
// The data type `CardEventData` is specified in the Data Dictionary, Section 2.19.
//
// ASN.1 Definition:
//
//	CardEventData ::= SEQUENCE OF CardEventRecord
//
//	CardEventRecord ::= SEQUENCE {
//	    eventType                   EventFaultType,                     -- 1 byte
//	    eventBeginTime              TimeReal,                         -- 4 bytes
//	    eventEndTime                TimeReal,                         -- 4 bytes
//	    eventVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
//	}
const (
	// CardEventRecord size (24 bytes total)
	cardEventRecordSize = 24
)

// splitCardEventRecord returns a SplitFunc that splits data into 24-byte event records
func splitCardEventRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < cardEventRecordSize {
		if atEOF {
			return 0, nil, nil // No more complete records, but not an error
		}
		return 0, nil, nil // Need more data
	}

	return cardEventRecordSize, data[:cardEventRecordSize], nil
}

func unmarshalEventsData(data []byte) (*cardv1.EventsData, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(splitCardEventRecord)

	var records []*cardv1.EventsData_Record
	for scanner.Scan() {
		recordData := scanner.Bytes()
		// Check if this is a valid record by examining the event begin time (first 4 bytes after event type)
		// Event type is 1 byte, so event begin time starts at byte 1
		eventBeginTime := binary.BigEndian.Uint32(recordData[1:5])
		if eventBeginTime == 0 {
			// Non-valid record: preserve original bytes
			rec := &cardv1.EventsData_Record{}
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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Use simplified schema with single events array in chronological order
	var ed cardv1.EventsData
	ed.SetEvents(records)
	return &ed, nil
}

// unmarshalEventRecord parses a single event record.
//
// The data type `CardEventRecord` is specified in the Data Dictionary, Section 2.20.
//
// ASN.1 Definition:
//
//	CardEventRecord ::= SEQUENCE {
//	    eventType                   EventFaultType,                     -- 1 byte
//	    eventBeginTime              TimeReal,                         -- 4 bytes
//	    eventEndTime                TimeReal,                         -- 4 bytes
//	    eventVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
//	}
func unmarshalEventRecord(data []byte) (*cardv1.EventsData_Record, error) {
	const (
		lenEventType                = 1
		lenEventBeginTime           = 4
		lenEventEndTime             = 4
		lenEventVehicleRegistration = 15
		lenCardEventRecord          = lenEventType + lenEventBeginTime + lenEventEndTime + lenEventVehicleRegistration
	)

	if len(data) < lenCardEventRecord {
		return nil, fmt.Errorf("insufficient data for event record: got %d bytes, need %d", len(data), lenCardEventRecord)
	}

	var rec cardv1.EventsData_Record
	offset := 0

	// Read event type (1 byte) and convert using generic enum helper
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for event type")
	}
	eventType := data[offset]
	enumDesc := ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED.Descriptor()
	if enumNum, found := dd.SetEnumFromProtocolValue(enumDesc, int32(eventType)); found {
		rec.SetEventType(ddv1.EventFaultType(enumNum))
	}
	offset++

	// Read event begin time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for event begin time")
	}
	rec.SetEventBeginTime(dd.ReadTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read event end time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for event end time")
	}
	rec.SetEventEndTime(dd.ReadTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read vehicle registration nation (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration nation")
	}
	nation, err := dd.UnmarshalNationNumeric(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	offset++

	// Create VehicleRegistrationIdentification structure
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)

	// Read vehicle registration number (14 bytes)
	if offset+14 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration number")
	}
	regNumber, err := dd.UnmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	// offset += 14 // Not needed as this is the last field
	vehicleReg.SetNumber(regNumber)
	rec.SetEventVehicleRegistration(vehicleReg)
	return &rec, nil
}

// AppendEventsData appends the binary representation of EventData to dst.
//
// The data type `CardEventData` is specified in the Data Dictionary, Section 2.19.
//
// ASN.1 Definition:
//
//	CardEventData ::= SEQUENCE OF CardEventRecord
//
//	CardEventRecord ::= SEQUENCE {
//	    eventType                        EventFaultType,                     -- 1 byte
//	    eventBeginTime                   TimeReal,                         -- 4 bytes
//	    eventEndTime                     TimeReal,                         -- 4 bytes
//	    eventVehicleRegistration         VehicleRegistrationIdentification -- 15 bytes
//	}
func appendEventsData(dst []byte, data *cardv1.EventsData) ([]byte, error) {
	var err error

	// Process events in their chronological order
	for _, r := range data.GetEvents() {
		dst, err = appendEventRecord(dst, r)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendEventRecord appends a single event record to dst.
//
// The data type `CardEventRecord` is specified in the Data Dictionary, Section 2.20.
//
// ASN.1 Definition:
//
//	CardEventRecord ::= SEQUENCE {
//	    eventType                        EventFaultType,                     -- 1 byte
//	    eventBeginTime                   TimeReal,                         -- 4 bytes
//	    eventEndTime                     TimeReal,                         -- 4 bytes
//	    eventVehicleRegistration         VehicleRegistrationIdentification -- 15 bytes
//	}
func appendEventRecord(dst []byte, record *cardv1.EventsData_Record) ([]byte, error) {
	if !record.GetValid() {
		return append(dst, record.GetRawData()...), nil
	}

	protocolValue := dd.GetEventFaultTypeProtocolValue(record.GetEventType(), 0)
	dst = append(dst, byte(protocolValue))
	dst = dd.AppendTimeReal(dst, record.GetEventBeginTime())
	dst = dd.AppendTimeReal(dst, record.GetEventEndTime())
	dst, err := dd.AppendVehicleRegistration(dst, record.GetEventVehicleRegistration())
	if err != nil {
		return nil, err
	}

	return dst, nil
}
