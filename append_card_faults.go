package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendFaultsData appends the binary representation of FaultData to dst.
func AppendFaultsData(dst []byte, data *cardv1.FaultsData) ([]byte, error) {
	var err error
	for _, r := range data.GetRecords() {
		dst, err = AppendFaultRecord(dst, r)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendFaultRecord appends a single fault record.
func AppendFaultRecord(dst []byte, record *cardv1.FaultsData_Record) ([]byte, error) {
	if !record.GetValid() {
		return append(dst, record.GetRawData()...), nil
	}

	protocolValue := GetEventFaultTypeProtocolValue(record.GetFaultType(), 0)
	dst = append(dst, byte(protocolValue))
	dst = appendTimeReal(dst, record.GetFaultBeginTime())
	dst = appendTimeReal(dst, record.GetFaultEndTime())
	dst, err := appendVehicleRegistration(dst, record.GetFaultVehicleRegistration())
	if err != nil {
		return nil, err
	}

	return dst, nil
}
