package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/proto"
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
		_, _ = r.Read(trailingBytes) // ignore error as we're reading from in-memory buffer
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
	proto.Merge(target, result)
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
	_ = binary.Read(r, binary.BigEndian, &reserved) // ignore error as we're reading from in-memory buffer
	record.SetReservedByte(int32(reserved))

	return record, nil
}
