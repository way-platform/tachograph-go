package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendEventsData appends the binary representation of EventData to dst.
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

// AppendEventRecord appends a single event record.
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
