package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardPlaces unmarshals places data from a card EF.
func unmarshalCardPlaces(data []byte) (*cardv1.Places, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for places")
	}

	var target cardv1.Places
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse place records
	var records []*cardv1.Places_Record
	recordSize := 12 // Fixed size: 4 bytes time + 1 byte entry type + 1 byte country + 2 bytes region + 3 bytes odometer + 1 byte reserved

	for r.Len() >= recordSize {
		record, err := parsePlaceRecord(r)
		if err != nil {
			break // Stop parsing on error, but return what we have
		}
		records = append(records, record)
	}

	target.SetRecords(records)

	// Capture any remaining trailing bytes for roundtrip accuracy
	if r.Len() > 0 {
		trailingBytes := make([]byte, r.Len())
		r.Read(trailingBytes)
		target.SetTrailingBytes(trailingBytes)
	}

	return &target, nil
}

// UnmarshalCardPlaces unmarshals places data from a card EF (legacy function).
// Deprecated: Use unmarshalCardPlaces instead.
func UnmarshalCardPlaces(data []byte, target *cardv1.Places) error {
	result, err := unmarshalCardPlaces(data)
	if err != nil {
		return err
	}
	*target = *result
	return nil
}

// parsePlaceRecord parses a single place record
func parsePlaceRecord(r *bytes.Reader) (*cardv1.Places_Record, error) {
	if r.Len() < 12 {
		return nil, fmt.Errorf("insufficient data for place record")
	}

	record := &cardv1.Places_Record{}

	// Read entry time (4 bytes)
	record.SetEntryTime(readTimeReal(r))

	// Read entry type (1 byte)
	var entryType byte
	if err := binary.Read(r, binary.BigEndian, &entryType); err != nil {
		return nil, fmt.Errorf("failed to read entry type: %w", err)
	}
	// Convert raw entry type to enum using protocol annotations
	SetEntryTypeDailyWorkPeriod(int32(entryType), record.SetEntryType, record.SetUnrecognizedEntryType)

	// Read daily work period country (1 byte)
	var country byte
	if err := binary.Read(r, binary.BigEndian, &country); err != nil {
		return nil, fmt.Errorf("failed to read country: %w", err)
	}
	// Convert raw country to enum using protocol annotations
	SetNationNumeric(int32(country), record.SetDailyWorkPeriodCountry, record.SetUnrecognizedDailyWorkPeriodCountry)

	// Read daily work period region (2 bytes)
	var region uint16
	if err := binary.Read(r, binary.BigEndian, &region); err != nil {
		return nil, fmt.Errorf("failed to read region: %w", err)
	}
	record.SetDailyWorkPeriodRegion(int32(region))

	// Read vehicle odometer (3 bytes)
	odometerBytes := make([]byte, 3)
	if _, err := r.Read(odometerBytes); err != nil {
		return nil, fmt.Errorf("failed to read odometer: %w", err)
	}
	record.SetVehicleOdometerKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBytes...))))

	// Read reserved byte (1 byte) and store it for roundtrip accuracy
	var reserved byte
	binary.Read(r, binary.BigEndian, &reserved)
	record.SetReservedByte(int32(reserved))

	return record, nil
}
