package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalFaultsData parses the binary data for an EF_Faults_Data record.
func unmarshalFaultsData(data []byte) (*cardv1.FaultsData, error) {
	const recordSize = 24
	r := bytes.NewReader(data)
	var records []*cardv1.FaultsData_Record

	for r.Len() >= recordSize {
		recordData := make([]byte, recordSize)
		r.Read(recordData)

		// Check if this is a valid record by examining the fault begin time (first 4 bytes after fault type)
		// Fault type is 1 byte, so fault begin time starts at byte 1
		faultBeginTime := binary.BigEndian.Uint32(recordData[1:5])

		rec := &cardv1.FaultsData_Record{}

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

	var fd cardv1.FaultsData
	fd.SetRecords(records)
	return &fd, nil
}

// UnmarshalFaultsData parses the binary data for an EF_Faults_Data record (legacy function).
// Deprecated: Use unmarshalFaultsData instead.
func UnmarshalFaultsData(data []byte, fd *cardv1.FaultsData) error {
	result, err := unmarshalFaultsData(data)
	if err != nil {
		return err
	}
	*fd = *result
	return nil
}

// UnmarshalFaultRecord parses a single 24-byte fault record.
func UnmarshalFaultRecord(data []byte, rec *cardv1.FaultsData_Record) error {
	r := bytes.NewReader(data)
	faultType, _ := r.ReadByte()

	// Convert raw fault type to enum using protocol annotations
	SetEventFaultType(int32(faultType), rec.SetFaultType, nil)

	rec.SetFaultBeginTime(readTimeReal(r))
	rec.SetFaultEndTime(readTimeReal(r))
	// Read vehicle registration
	nation, err := unmarshalNationNumericFromReader(r)
	if err != nil {
		return fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	// Create VehicleRegistrationIdentification structure
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)

	regNumber, err := unmarshalIA5StringValueFromReader(r, 14)
	if err != nil {
		return fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	vehicleReg.SetNumber(regNumber)
	rec.SetFaultVehicleRegistration(vehicleReg)
	return nil
}
