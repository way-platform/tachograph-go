package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendEventsData appends the binary representation of EventData to dst.
//
// ASN.1 Specification (Data Dictionary 2.19):
//
//	CardEventData ::= SEQUENCE OF CardEventRecord
//
//	CardEventRecord ::= SEQUENCE {
//	    eventType                        EventFaultType,
//	    eventBeginTime                   TimeReal,
//	    eventEndTime                     TimeReal,
//	    eventVehicleRegistration         VehicleRegistrationIdentification,
//	    eventTypeSpecificData            OCTET STRING (SIZE (2))
//	}
//
// Binary Layout (variable size):
//
//	Each record: 35 bytes
//	- 0-0:   eventType (1 byte)
//	- 1-4:   eventBeginTime (4 bytes, TimeReal)
//	- 5-8:   eventEndTime (4 bytes, TimeReal)
//	- 9-23:  eventVehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	- 24-25: eventTypeSpecificData (2 bytes)
//

func AppendEventsData(dst []byte, data *cardv1.EventsData) ([]byte, error) {
	var err error
	for _, r := range data.GetRecords() {
		dst, err = AppendEventRecord(dst, r)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendEventRecord appends a single event record to dst.
//
// ASN.1 Specification (Data Dictionary 2.20):
//
//	CardEventRecord ::= SEQUENCE {
//	    eventType                        EventFaultType,
//	    eventBeginTime                   TimeReal,
//	    eventEndTime                     TimeReal,
//	    eventVehicleRegistration         VehicleRegistrationIdentification,
//	    eventTypeSpecificData            OCTET STRING (SIZE (2))
//	}
//
// Binary Layout (35 bytes):
//
//	0-0:   eventType (1 byte)
//	1-4:   eventBeginTime (4 bytes, TimeReal)
//	5-8:   eventEndTime (4 bytes, TimeReal)
//	9-23:  eventVehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	24-25: eventTypeSpecificData (2 bytes)
func AppendEventRecord(dst []byte, record *cardv1.EventsData_Record) ([]byte, error) {
	if !record.GetValid() {
		return append(dst, record.GetRawData()...), nil
	}

	protocolValue := GetEventFaultTypeProtocolValue(record.GetEventType(), 0)
	dst = append(dst, byte(protocolValue))
	dst = appendTimeReal(dst, record.GetEventBeginTime())
	dst = appendTimeReal(dst, record.GetEventEndTime())
	dst, err := appendVehicleRegistration(dst, record.GetEventVehicleRegistration())
	if err != nil {
		return nil, err
	}

	return dst, nil
}
