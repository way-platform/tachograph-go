package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalCardPlaces unmarshals places data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.4):
//
//	CardPlaceDailyWorkPeriod ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort
//	}
//
// Binary Layout (variable size):
//
//	0-1:   newestRecordIndex (2 bytes, big-endian)
//	2+:    place records (12 bytes each)
//	  - 0-3:   entryTime (4 bytes)
//	  - 4-4:   entryTypeDailyWorkPeriod (1 byte)
//	  - 5-5:   dailyWorkPeriodCountry (1 byte)
//	  - 6-7:   dailyWorkPeriodRegion (2 bytes, big-endian)
//	  - 8-10:  vehicleOdometerValue (3 bytes)
//	  - 11-11: reserved (1 byte)
func unmarshalCardPlaces(data []byte) (*cardv1.Places, error) {
	const (
		// Minimum EF_Places record size
		MIN_EF_PLACES_SIZE = 2

		// Field offsets within place record
		ENTRY_TIME_OFFSET                = 0
		ENTRY_TYPE_OFFSET                = 4
		DAILY_WORK_PERIOD_COUNTRY_OFFSET = 5
		DAILY_WORK_PERIOD_REGION_OFFSET  = 6
		VEHICLE_ODOMETER_VALUE_OFFSET    = 8
		RESERVED_OFFSET                  = 11

		// Field sizes
		ENTRY_TIME_SIZE                = 4
		ENTRY_TYPE_SIZE                = 1
		DAILY_WORK_PERIOD_COUNTRY_SIZE = 1
		DAILY_WORK_PERIOD_REGION_SIZE  = 2
		VEHICLE_ODOMETER_VALUE_SIZE    = 3
		RESERVED_SIZE                  = 1
	)

	if len(data) < MIN_EF_PLACES_SIZE {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least %d", len(data), MIN_EF_PLACES_SIZE)
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
	SetEnumFromProtocolValue(datadictionaryv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED.Descriptor(),
		int32(entryType),
		func(enumNum protoreflect.EnumNumber) {
			record.SetEntryType(datadictionaryv1.EntryTypeDailyWorkPeriod(enumNum))
		}, nil)

	// Read daily work period country (1 byte)
	var country byte
	if err := binary.Read(r, binary.BigEndian, &country); err != nil {
		return nil, fmt.Errorf("failed to read country: %w", err)
	}
	// Convert raw country to enum using protocol annotations
	SetEnumFromProtocolValue(datadictionaryv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(),
		int32(country),
		func(enumNum protoreflect.EnumNumber) {
			record.SetDailyWorkPeriodCountry(datadictionaryv1.NationNumeric(enumNum))
		}, nil)

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
