package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendFaultRecord appends the binary representation of a single fault record to dst.
func AppendFaultRecord(dst []byte, rec *cardv1.FaultData_Record) ([]byte, error) {
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
	faultTypeProtocol := GetEventFaultTypeProtocolValue(rec.GetFaultType(), rec.GetUnrecognizedFaultType())
	dst = append(dst, byte(faultTypeProtocol))
	dst = appendTimeReal(dst, rec.GetFaultBeginTime())
	dst = appendTimeReal(dst, rec.GetFaultEndTime())
	dst = append(dst, byte(0)) // Placeholder for vehicleRegistrationNation, assuming BCD
	dst = appendString(dst, rec.GetVehicleRegistrationNumber(), 14)
	return dst, nil
}
