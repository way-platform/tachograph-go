package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendEventRecord appends the binary representation of a single event record to dst.
func AppendEventRecord(dst []byte, rec *cardv1.EventData_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, 24)...), nil
	}

	if !rec.GetValid() {
		// Non-valid record: use preserved raw data
		rawData := rec.GetRawData()
		if len(rawData) != 24 {
			// Fallback to zeros if raw data is invalid
			return append(dst, make([]byte, 24)...), nil
		}
		return append(dst, rawData...), nil
	}

	// Valid record: serialize semantic data
	eventTypeProtocol := GetEventFaultTypeProtocolValue(rec.GetEventType(), rec.GetUnrecognizedEventType())
	dst = append(dst, byte(eventTypeProtocol))
	dst = appendTimeReal(dst, rec.GetEventBeginTime())
	dst = appendTimeReal(dst, rec.GetEventEndTime())

	// Convert hex string nation back to byte
	nationByte := byte(0) // Default fallback
	if nationStr := rec.GetVehicleRegistrationNation(); len(nationStr) >= 2 {
		if parsedNation, err := parseHexByte(nationStr); err == nil {
			nationByte = parsedNation
		}
	}
	dst = append(dst, nationByte)

	dst = appendString(dst, rec.GetVehicleRegistrationNumber(), 14)
	return dst, nil
}
