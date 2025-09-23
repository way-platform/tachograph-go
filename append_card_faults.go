package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendFaultRecord appends the binary representation of a single fault record to dst.
func AppendFaultRecord(dst []byte, rec *cardv1.FaultData_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, 24)...), nil
	}
	dst = append(dst, byte(rec.GetFaultType()))
	dst = appendTimeReal(dst, rec.GetFaultBeginTime())
	dst = appendTimeReal(dst, rec.GetFaultEndTime())
	dst = append(dst, byte(0)) // Placeholder for vehicleRegistrationNation, assuming BCD
	dst = appendString(dst, rec.GetVehicleRegistrationNumber(), 14)
	return dst, nil
}
