package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalCardPlaces unmarshals places data from a card EF, handling both Gen1 and Gen2 formats.
//
// The data type `CardPlaceDailyWorkPeriod` is specified in the Data Dictionary, Section 2.4.
//
// ASN.1 Definition (Gen1):
//
//	CardPlaceDailyWorkPeriod ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort
//	}
//
// ASN.1 Definition (Gen2):
//
//	PlaceRecord_G2 ::= SEQUENCE {
//		entryTime                    TimeReal,
//		entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//		dailyWorkPeriodCountry       NationNumeric,
//		dailyWorkPeriodRegion        RegionNumeric,
//		vehicleOdometerValue         OdometerShort,
//		gnssPlaceRecord              GNSSPlaceRecord
//	}
const (
	placeRecordSizeGen1 = 10
	placeRecordSizeGen2 = 21
	lenMinEfPlaces      = 2 // Minimum EF_Places file size (for the pointer)
)

func unmarshalCardPlaces(data []byte, generation ddv1.Generation) (*cardv1.Places, error) {
	if len(data) < lenMinEfPlaces {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least %d", len(data), lenMinEfPlaces)
	}

	// Create unmarshal options from generation
	opts := dd.UnmarshalOptions{
		Generation: generation,
	}

	var target cardv1.Places
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse place records from circular buffer
	remainingData := make([]byte, r.Len())
	if _, err := r.Read(remainingData); err != nil {
		return nil, fmt.Errorf("failed to read remaining data: %w", err)
	}

	records, trailingBytes := parseCircularPlaceRecords(remainingData, int(newestRecordIndex), opts)
	target.SetRecords(records)
	target.SetTrailingBytes(trailingBytes)

	return &target, nil
}

// parseCircularPlaceRecords parses place records from a circular buffer, starting from the oldest valid record.
// It returns the parsed records and any trailing bytes that do not form a complete record.
func parseCircularPlaceRecords(data []byte, newestRecordIndex int, opts dd.UnmarshalOptions) ([]*ddv1.PlaceRecord, []byte) {
	recordSize := placeRecordSizeGen1
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		recordSize = placeRecordSizeGen2
	}

	if len(data) < recordSize {
		return nil, data // Not enough data for even one record
	}

	totalRecords := len(data) / recordSize
	if totalRecords == 0 {
		return nil, data
	}

	var validRecords []*ddv1.PlaceRecord

	// Start from the record after the newest (which should be the oldest)
	// If newestRecordIndex is out of bounds, start from beginning
	startIndex := 0
	if newestRecordIndex >= 0 && newestRecordIndex < totalRecords {
		startIndex = (newestRecordIndex + 1) % totalRecords
	}

	// Read records in chronological order (oldest to newest)
	// Stop as soon as we encounter an invalid record
	for i := 0; i < totalRecords; i++ {
		recordIndex := (startIndex + i) % totalRecords
		recordOffset := recordIndex * recordSize

		if recordOffset+recordSize > len(data) {
			break // Should not happen with the totalRecords calculation, but as a safeguard
		}

		recordData := data[recordOffset : recordOffset+recordSize]
		record, valid := unmarshalPlaceRecordWithValidation(recordData, opts)

		if valid {
			validRecords = append(validRecords, record)
		} else {
			// First invalid record marks the end of valid data in the circular buffer
			break
		}
	}

	// Capture any remaining trailing bytes for roundtrip accuracy
	totalRecordsSize := totalRecords * recordSize
	var trailingBytes []byte
	if len(data) > totalRecordsSize {
		trailingBytes = data[totalRecordsSize:]
	}

	return validRecords, trailingBytes
}

// unmarshalPlaceRecordWithValidation parses a place record and validates it
func unmarshalPlaceRecordWithValidation(data []byte, opts dd.UnmarshalOptions) (*ddv1.PlaceRecord, bool) {
	record, err := unmarshalPlaceRecord(data, opts)
	if err != nil {
		// If parsing fails, treat as invalid but keep raw data
		invalidRecord := &ddv1.PlaceRecord{}
		invalidRecord.SetValid(false)
		invalidRecord.SetRawData(data)
		return invalidRecord, false
	}

	// Validate the record
	if !isValidPlaceRecord(record) {
		record.SetValid(false)
		record.SetRawData(data)
		return record, false
	}

	record.SetValid(true)
	return record, true
}

// isValidPlaceRecord validates a place record for reasonable values
func isValidPlaceRecord(record *ddv1.PlaceRecord) bool {
	// Check timestamp validity (reasonable range for tachograph data)
	// Tachographs were introduced around 1985, so anything before 1980 is suspicious
	// Anything after 2050 is also suspicious
	entryTime := record.GetEntryTime()
	if entryTime != nil {
		year := entryTime.AsTime().Year()
		if year < 1980 || year > 2050 {
			return false
		}
	}

	// Check odometer value (should be reasonable - not negative, not extremely large)
	odometer := record.GetVehicleOdometerKm()
	if odometer < 0 || odometer > 10000000 { // 10M km is unreasonably high
		return false
	}

	// All other fields can have various values including unrecognized ones,
	// so we don't validate them as strictly
	return true
}

// unmarshalPlaceRecord parses a single place record from a byte slice based on card generation.
func unmarshalPlaceRecord(data []byte, opts dd.UnmarshalOptions) (*ddv1.PlaceRecord, error) {
	recordSize := placeRecordSizeGen1
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		recordSize = placeRecordSizeGen2
	}

	if len(data) < recordSize {
		return nil, fmt.Errorf("insufficient data for place record: got %d, want %d", len(data), recordSize)
	}

	r := bytes.NewReader(data)
	record := &ddv1.PlaceRecord{}

	// Read entry time (4 bytes)
	entryTimeBytes := make([]byte, 4)
	if _, err := r.Read(entryTimeBytes); err != nil {
		return nil, fmt.Errorf("failed to read entry time: %w", err)
	}
	entryTime, err := opts.UnmarshalTimeReal(entryTimeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// Read entry type (1 byte)
	entryType, _ := r.ReadByte()
	dd.SetEnumFromProtocolValueGeneric(ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED.Descriptor(),
		int32(entryType),
		func(enumNum protoreflect.EnumNumber) {
			record.SetEntryTypeDailyWorkPeriod(ddv1.EntryTypeDailyWorkPeriod(enumNum))
		}, func(rawValue int32) {
			record.SetUnrecognizedEntryTypeDailyWorkPeriod(rawValue)
		})

	// Read daily work period country (1 byte)
	country, _ := r.ReadByte()
	dd.SetEnumFromProtocolValueGeneric(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(),
		int32(country),
		func(enumNum protoreflect.EnumNumber) {
			record.SetDailyWorkPeriodCountry(ddv1.NationNumeric(enumNum))
		}, func(rawValue int32) {
			record.SetUnrecognizedDailyWorkPeriodCountry(rawValue)
		})

	// Read daily work period region (1 byte for both Gen1 and Gen2 as per spec)
	region, _ := r.ReadByte()
	record.SetDailyWorkPeriodRegion([]byte{region})

	// Read vehicle odometer (3 bytes)
	odometerBytes := make([]byte, 3)
	if _, err := r.Read(odometerBytes); err != nil {
		return nil, fmt.Errorf("failed to read odometer: %w", err)
	}
	record.SetVehicleOdometerKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBytes...))))

	if opts.Generation == ddv1.Generation_GENERATION_2 {
		// Read GNSS place record (11 bytes: 4 timestamp + 1 accuracy + 6 coords)
		gnssBytes := make([]byte, 11)
		if _, err := r.Read(gnssBytes); err != nil {
			return nil, fmt.Errorf("failed to read GNSS place record: %w", err)
		}
		gnssRecord, err := opts.UnmarshalGNSSPlaceRecord(gnssBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GNSS place record: %w", err)
		}
		record.SetEntryGnssPlaceRecord(gnssRecord)
	}

	return record, nil
}

// appendPlaces appends the binary representation of Places to dst.
func appendPlaces(dst []byte, p *cardv1.Places, generation ddv1.Generation) ([]byte, error) {
	if p == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(p.GetNewestRecordIndex()))

	var err error
	for _, rec := range p.GetRecords() {
		dst, err = appendPlaceRecord(dst, rec, generation)
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

// appendPlaceRecord appends a single place record based on card generation.
func appendPlaceRecord(dst []byte, rec *ddv1.PlaceRecord, generation ddv1.Generation) ([]byte, error) {
	recordSize := placeRecordSizeGen1
	if generation == ddv1.Generation_GENERATION_2 {
		recordSize = placeRecordSizeGen2
	}

	if rec == nil || !rec.GetValid() {
		if raw := rec.GetRawData(); len(raw) > 0 {
			return append(dst, raw...), nil
		}
		return append(dst, make([]byte, recordSize)...), nil
	}

	var err error
	dst, err = dd.AppendTimeReal(dst, rec.GetEntryTime()) // 4 bytes
	if err != nil {
		return nil, fmt.Errorf("failed to append entry time: %w", err)
	}

	entryTypeProtocol, _ := dd.GetProtocolValueForEnum(rec.GetEntryTypeDailyWorkPeriod())
	dst = append(dst, byte(entryTypeProtocol)) // 1 byte

	countryProtocol, _ := dd.GetProtocolValueForEnum(rec.GetDailyWorkPeriodCountry())
	dst = append(dst, byte(countryProtocol)) // 1 byte

	// Append region byte (1 byte)
	regionBytes := rec.GetDailyWorkPeriodRegion()
	if len(regionBytes) > 0 {
		dst = append(dst, regionBytes[0])
	} else {
		dst = append(dst, 0)
	}

	dst = dd.AppendOdometer(dst, uint32(rec.GetVehicleOdometerKm())) // 3 bytes

	if generation == ddv1.Generation_GENERATION_2 {
		gnssRecord := rec.GetEntryGnssPlaceRecord()
		// Append GNSS place record (11 bytes: 4 timestamp + 1 accuracy + 6 coords)
		dst, err = dd.AppendGNSSPlaceRecord(dst, gnssRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to append GNSS place record: %w", err)
		}
	}

	return dst, nil
}
