package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardPlaces unmarshals places data from a card EF.
func UnmarshalCardPlaces(data []byte, target *cardv1.Places) error {
	if len(data) < 2 {
		return fmt.Errorf("insufficient data for places")
	}

	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return fmt.Errorf("failed to read newest record index: %w", err)
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
	record.SetEntryType(int32(entryType))

	// Read daily work period country (1 byte)
	var country byte
	if err := binary.Read(r, binary.BigEndian, &country); err != nil {
		return nil, fmt.Errorf("failed to read country: %w", err)
	}
	record.SetDailyWorkPeriodCountry(int32(country))

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

	// Skip reserved byte (1 byte)
	var reserved byte
	binary.Read(r, binary.BigEndian, &reserved)

	return record, nil
}
