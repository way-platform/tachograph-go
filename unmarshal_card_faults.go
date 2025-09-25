package tachograph

import (
	"bytes"
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalFaultsData parses the binary data for an EF_Faults_Data record.
func unmarshalFaultsData(data []byte) (*cardv1.FaultData, error) {
	const recordSize = 24
	r := bytes.NewReader(data)
	var records []*cardv1.FaultData_Record

	for r.Len() >= recordSize {
		recordData := make([]byte, recordSize)
		r.Read(recordData)

		// Check if this is a valid record by examining the fault begin time (first 4 bytes after fault type)
		// Fault type is 1 byte, so fault begin time starts at byte 1
		faultBeginTime := binary.BigEndian.Uint32(recordData[1:5])

		rec := &cardv1.FaultData_Record{}

		if faultBeginTime == 0 {
			// Non-valid record: preserve original bytes
			rec.SetValid(false)
			rec.SetRawData(recordData)
		} else {
			// Valid record: parse semantic data
			rec.SetValid(true)
			if err := UnmarshalFaultRecord(recordData, rec); err != nil {
				return nil, err
			}
		}

		records = append(records, rec)
	}

	var fd cardv1.FaultData
	fd.SetRecords(records)
	return &fd, nil
}

// UnmarshalFaultsData parses the binary data for an EF_Faults_Data record (legacy function).
// Deprecated: Use unmarshalFaultsData instead.
func UnmarshalFaultsData(data []byte, fd *cardv1.FaultData) error {
	result, err := unmarshalFaultsData(data)
	if err != nil {
		return err
	}
	*fd = *result
	return nil
}

// UnmarshalFaultRecord parses a single 24-byte fault record.
func UnmarshalFaultRecord(data []byte, rec *cardv1.FaultData_Record) error {
	r := bytes.NewReader(data)
	faultType, _ := r.ReadByte()

	// Convert raw fault type to enum using protocol annotations
	SetEventFaultType(int32(faultType), rec.SetFaultType, rec.SetUnrecognizedFaultType)

	rec.SetFaultBeginTime(readTimeReal(r))
	rec.SetFaultEndTime(readTimeReal(r))
	// TODO: Read BCD nation code
	r.ReadByte()
	rec.SetVehicleRegistrationNumber(readString(r, 14))
	return nil
}
