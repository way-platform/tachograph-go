package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalCardVehiclesUsed unmarshals vehicles used data from a card EF.
func unmarshalCardVehiclesUsed(data []byte) (*cardv1.VehiclesUsed, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for vehicles used")
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
	*target = *result
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
	vehicleReg.SetNation(int32(nation))
	vehicleReg.SetNumber(readString(r, 14))
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
