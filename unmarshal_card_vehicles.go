package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalCardVehiclesUsed unmarshals vehicles used data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.6):
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBeginKm       OdometerShort,
//	    vehicleOdometerEndKm         OdometerShort,
//	    vehicleFirstUse              TimeReal,
//	    vehicleLastUse               TimeReal,
//	    vehicleRegistration          VehicleRegistrationIdentification,
//	    vehicleIdentificationNumber  VehicleIdentificationNumber OPTIONAL,
//	    vehicleRegistrationNation    NationNumeric OPTIONAL,
//	    vehicleRegistrationNumber    RegistrationNumber OPTIONAL
//	}
//
// Binary Layout (variable size):
//
//	0-1:   newestRecordIndex (2 bytes, big-endian)
//	2+:    vehicle records (31 bytes Gen1, 48 bytes Gen2 each)
//	  - 0-2:   vehicleOdometerBeginKm (3 bytes)
//	  - 3-5:   vehicleOdometerEndKm (3 bytes)
//	  - 6-9:   vehicleFirstUse (4 bytes)
//	  - 10-13: vehicleLastUse (4 bytes)
//	  - 14-28: vehicleRegistration (15 bytes: 1 byte nation + 14 bytes registration)
//	  - 29-30: vehicleIdentificationNumber (2 bytes, Gen2+)
//	  - 31-31: vehicleRegistrationNation (1 byte, Gen2+)
//	  - 32-47: vehicleRegistrationNumber (16 bytes, Gen2+)
func unmarshalCardVehiclesUsed(data []byte) (*cardv1.VehiclesUsed, error) {
	const (
		// Minimum EF_VehiclesUsed record size
		MIN_EF_VEHICLES_USED_SIZE = 2
	)

	if len(data) < MIN_EF_VEHICLES_USED_SIZE {
		return nil, fmt.Errorf("insufficient data for vehicles used: got %d bytes, need at least %d", len(data), MIN_EF_VEHICLES_USED_SIZE)
	}

	var target cardv1.VehiclesUsed
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse vehicle records
	var records []*cardv1.VehiclesUsed_Record
	recordSize := determineVehicleRecordSize(data)

	for r.Len() >= recordSize {
		record, err := parseVehicleRecord(r, recordSize)
		if err != nil {
			break // Stop parsing on error, but return what we have
		}
		records = append(records, record)
	}

	target.SetRecords(records)
	return &target, nil
}

// UnmarshalCardVehiclesUsed unmarshals vehicles used data from a card EF (legacy function).
// Deprecated: Use unmarshalCardVehiclesUsed instead.
func UnmarshalCardVehiclesUsed(data []byte, target *cardv1.VehiclesUsed) error {
	result, err := unmarshalCardVehiclesUsed(data)
	if err != nil {
		return err
	}
	proto.Merge(target, result)
	return nil
}

// determineVehicleRecordSize determines the size of vehicle records based on data length
// Gen1: 31 bytes per record, Gen2: 48 bytes per record
func determineVehicleRecordSize(data []byte) int {
	remainingData := len(data) - 2 // Subtract pointer size

	// Try to determine generation based on typical record counts and sizes
	if remainingData%48 == 0 {
		return 48 // Gen2 record size
	}
	return 31 // Gen1 record size (default)
}

// parseVehicleRecord parses a single vehicle record
func parseVehicleRecord(r *bytes.Reader, recordSize int) (*cardv1.VehiclesUsed_Record, error) {
	if r.Len() < recordSize {
		return nil, fmt.Errorf("insufficient data for vehicle record")
	}

	record := &cardv1.VehiclesUsed_Record{}

	// Read odometer begin (3 bytes)
	odometerBeginBytes := make([]byte, 3)
	if _, err := r.Read(odometerBeginBytes); err != nil {
		return nil, fmt.Errorf("failed to read odometer begin: %w", err)
	}
	record.SetVehicleOdometerBeginKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBeginBytes...))))

	// Read odometer end (3 bytes)
	odometerEndBytes := make([]byte, 3)
	if _, err := r.Read(odometerEndBytes); err != nil {
		return nil, fmt.Errorf("failed to read odometer end: %w", err)
	}
	record.SetVehicleOdometerEndKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerEndBytes...))))

	// Read vehicle first use (4 bytes)
	record.SetVehicleFirstUse(readTimeReal(r))

	// Read vehicle last use (4 bytes)
	record.SetVehicleLastUse(readTimeReal(r))

	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes registration number)
	var nation byte
	if err := binary.Read(r, binary.BigEndian, &nation); err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	// Create VehicleRegistrationIdentification structure
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(datadictionaryv1.NationNumeric(nation))

	// Read vehicle registration number (14 bytes)
	regNumberBytes := make([]byte, 14)
	if _, err := r.Read(regNumberBytes); err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	regNumber, err := unmarshalIA5StringValue(regNumberBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration number: %w", err)
	}
	vehicleReg.SetNumber(regNumber)
	record.SetVehicleRegistration(vehicleReg)

	// Read VU data block counter (2 bytes)
	var vuDataBlockCounter uint16
	if err := binary.Read(r, binary.BigEndian, &vuDataBlockCounter); err != nil {
		return nil, fmt.Errorf("failed to read VU data block counter: %w", err)
	}
	record.SetVuDataBlockCounter(int32(vuDataBlockCounter))

	// For Gen2 records, read VIN (17 bytes)
	if recordSize == 48 {
		vinBytes := make([]byte, 17)
		if _, err := r.Read(vinBytes); err != nil {
			return nil, fmt.Errorf("failed to read vehicle identification number: %w", err)
		}
		record.SetVehicleIdentificationNumber(readString(bytes.NewReader(vinBytes), 17))
	}

	return record, nil
}
