package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendEventRecord appends the binary representation of a single event record to dst.
func AppendEventRecord(dst []byte, rec *cardv1.EventData_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, 24)...), nil
	}
	dst = append(dst, byte(rec.GetEventType()))
	dst = appendTimeReal(dst, rec.GetEventBeginTime())
	dst = appendTimeReal(dst, rec.GetEventEndTime())
	dst = append(dst, byte(0)) // Placeholder for vehicleRegistrationNation, assuming BCD
	dst = appendString(dst, rec.GetVehicleRegistrationNumber(), 14)
	return dst, nil
}
