package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalFaultsData parses the binary data for an EF_Faults_Data record.
func unmarshalFaultsData(data []byte) (*cardv1.FaultsData, error) {
	r := bytes.NewReader(data)
	var records []*cardv1.FaultsData_Record

	for r.Len() >= cardEventFaultRecordSize {
		recordData := make([]byte, cardEventFaultRecordSize)
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
//
// ASN.1 Specification (Data Dictionary 2.22):
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                   EventFaultType,                     -- 1 byte
//	    faultBeginTime              TimeReal,                         -- 4 bytes
//	    faultEndTime                TimeReal,                         -- 4 bytes
//	    faultVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
//	}
func UnmarshalFaultRecord(data []byte, rec *cardv1.FaultsData_Record) error {
	const (
		// CardFaultRecord layout constants
		lenFaultType                = 1
		lenFaultBeginTime           = 4
		lenFaultEndTime             = 4
		lenFaultVehicleRegistration = 15
		totalLength                 = lenFaultType + lenFaultBeginTime + lenFaultEndTime + lenFaultVehicleRegistration
	)

	if len(data) < totalLength {
		return fmt.Errorf("insufficient data for fault record: got %d bytes, need %d", len(data), totalLength)
	}

	offset := 0

	// Read fault type (1 byte) and convert using generic enum helper
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for fault type")
	}
	faultType := data[offset]
	enumDesc := datadictionaryv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED.Descriptor()
	SetEnumFromProtocolValue(enumDesc, int32(faultType),
		func(enumNum protoreflect.EnumNumber) {
			rec.SetFaultType(datadictionaryv1.EventFaultType(enumNum))
		}, nil)
	offset++

	// Read fault begin time (4 bytes)
	if offset+4 > len(data) {
		return fmt.Errorf("insufficient data for fault begin time")
	}
	rec.SetFaultBeginTime(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read fault end time (4 bytes)
	if offset+4 > len(data) {
		return fmt.Errorf("insufficient data for fault end time")
	}
	rec.SetFaultEndTime(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read vehicle registration nation (1 byte)
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for vehicle registration nation")
	}
	nation, err := unmarshalNationNumeric(data[offset : offset+1])
	if err != nil {
		return fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	offset++

	// Create VehicleRegistrationIdentification structure
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)

	// Read vehicle registration number (14 bytes)
	if offset+14 > len(data) {
		return fmt.Errorf("insufficient data for vehicle registration number")
	}
	regNumber, err := unmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	offset += 14
	vehicleReg.SetNumber(regNumber)
	rec.SetFaultVehicleRegistration(vehicleReg)
	return nil
}
