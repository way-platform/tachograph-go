package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalCardPlaces unmarshals places data from a card EF.
//
// The data type `CardPlaceDailyWorkPeriod` is specified in the Data Dictionary, Section 2.4.
//
// ASN.1 Definition:
//
//	CardPlaceDailyWorkPeriod ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort
//	}
const (
	// PlaceRecord size (12 bytes total)
	placeRecordSize = 12
)

func unmarshalCardPlaces(data []byte) (*cardv1.Places, error) {
	const (
		lenMinEfPlaces = 2 // Minimum EF_Places record size
	)

	if len(data) < lenMinEfPlaces {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least %d", len(data), lenMinEfPlaces)
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

// parsePlaceRecord parses a single place record
func parsePlaceRecord(r *bytes.Reader) (*cardv1.Places_Record, error) {
	const (
		lenPlaceRecord = 12
	)

	if r.Len() < lenPlaceRecord {
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
	setEnumFromProtocolValue(ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED.Descriptor(),
		int32(entryType),
		func(enumNum protoreflect.EnumNumber) {
			record.SetEntryType(ddv1.EntryTypeDailyWorkPeriod(enumNum))
		}, nil)

	// Read daily work period country (1 byte)
	var country byte
	if err := binary.Read(r, binary.BigEndian, &country); err != nil {
		return nil, fmt.Errorf("failed to read country: %w", err)
	}
	// Convert raw country to enum using protocol annotations
	setEnumFromProtocolValue(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(),
		int32(country),
		func(enumNum protoreflect.EnumNumber) {
			record.SetDailyWorkPeriodCountry(ddv1.NationNumeric(enumNum))
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

// AppendPlaces appends the binary representation of Places to dst.
//
// The data type `CardPlaceDailyWorkPeriod` is specified in the Data Dictionary, Section 2.4.
//
// ASN.1 Definition:
//
//	CardPlaceDailyWorkPeriod ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort
//	}
func appendPlaces(dst []byte, p *cardv1.Places) ([]byte, error) {
	if p == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(p.GetNewestRecordIndex()))

	var err error
	for _, rec := range p.GetRecords() {
		dst, err = appendPlaceRecord(dst, rec)
		if err != nil {
			return nil, err
		}
	}

	// Append trailing bytes for roundtrip accuracy
	if trailingBytes := p.GetTrailingBytes(); len(trailingBytes) > 0 {
		dst = append(dst, trailingBytes...)
	}

	return dst, nil
}

// AppendPlaceRecord appends a single place record.
func appendPlaceRecord(dst []byte, rec *cardv1.Places_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, placeRecordSize)...), nil
	}
	dst = appendTimeReal(dst, rec.GetEntryTime()) // 4 bytes

	// Entry type with protocol value conversion using generic helper
	entryTypeProtocol := getProtocolValueFromEnum(rec.GetEntryType(), 0)
	dst = append(dst, byte(entryTypeProtocol)) // 1 byte

	// Country with protocol value conversion using generic helper
	countryProtocol := getProtocolValueFromEnum(rec.GetDailyWorkPeriodCountry(), 0)
	dst = append(dst, byte(countryProtocol)) // 1 byte

	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetDailyWorkPeriodRegion())) // 2 bytes
	dst = appendOdometer(dst, uint32(rec.GetVehicleOdometerKm()))                    // 3 bytes
	dst = append(dst, byte(rec.GetReservedByte()))                                   // 1 byte reserved (preserved)
	return dst, nil
}
