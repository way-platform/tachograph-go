package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// Record sizes are fixed per generation - this is one of the benefits of type splitting!
const (
	placeRecordSizeGen1 = 10 // Fixed size: no conditionals needed
	placeRecordSizeGen2 = 21 // Fixed size: no conditionals needed
	lenMinEfPlaces      = 2  // Minimum EF_Places file size (for the pointer)
)

// unmarshalPlaces unmarshals places data from a card EF.
//
// The generation determines the record format:
// - Gen1: 10-byte records (PlaceRecord)
// - Gen2: 21-byte records (PlaceRecordG2 with GNSS data)
//
// The data type `CardPlaceDailyWorkPeriod` is specified in the Data Dictionary, Section 2.27.
//
// ASN.1 Definition:
//
//	CardPlaceDailyWorkPeriod ::= SEQUENCE {
//	    placePointerNewestRecord INTEGER(0..NoOfCardPlaceRecords-1),
//	    placeRecords SET SIZE(NoOfCardPlaceRecords) OF PlaceRecord
//	}
func (opts UnmarshalOptions) unmarshalPlaces(data []byte) (*cardv1.Places, error) {
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

	// Parse place records from circular buffer
	remainingData := make([]byte, r.Len())
	if _, err := r.Read(remainingData); err != nil {
		return nil, fmt.Errorf("failed to read remaining data: %w", err)
	}

	// Create dd.UnmarshalOptions from card-level UnmarshalOptions
	ddOpts := dd.UnmarshalOptions{
		Generation: opts.Generation,
		Version:    opts.Version,
	}

	// Parse records based on generation - note how clean this is with type splitting!
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		records, trailingBytes := parseCircularPlaceRecordsG2(remainingData, int(newestRecordIndex), ddOpts)
		target.SetRecordsG2(records)
		target.SetTrailingBytes(trailingBytes)
	} else {
		records, trailingBytes := parseCircularPlaceRecordsGen1(remainingData, int(newestRecordIndex), ddOpts)
		target.SetRecords(records)
		target.SetTrailingBytes(trailingBytes)
	}

	return &target, nil
}

// parseCircularPlaceRecordsGen1 parses Gen1 place records (10 bytes each) from a circular buffer.
func parseCircularPlaceRecordsGen1(data []byte, newestRecordIndex int, opts dd.UnmarshalOptions) ([]*ddv1.PlaceRecord, []byte) {
	const recordSize = placeRecordSizeGen1 // Fixed size - no conditionals!

	if len(data) < recordSize {
		return nil, data // Not enough data for even one record
	}

	totalRecords := len(data) / recordSize
	if totalRecords == 0 {
		return nil, data
	}

	var validRecords []*ddv1.PlaceRecord

	// Start from the record after the newest (which should be the oldest)
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
			break
		}

		recordData := data[recordOffset : recordOffset+recordSize]
		record, valid := unmarshalPlaceRecordGen1WithValidation(recordData, opts)

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

// parseCircularPlaceRecordsG2 parses Gen2 place records (21 bytes each) from a circular buffer.
func parseCircularPlaceRecordsG2(data []byte, newestRecordIndex int, opts dd.UnmarshalOptions) ([]*ddv1.PlaceRecordG2, []byte) {
	const recordSize = placeRecordSizeGen2 // Fixed size - no conditionals!

	if len(data) < recordSize {
		return nil, data // Not enough data for even one record
	}

	totalRecords := len(data) / recordSize
	if totalRecords == 0 {
		return nil, data
	}

	var validRecords []*ddv1.PlaceRecordG2

	// Start from the record after the newest (which should be the oldest)
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
			break
		}

		recordData := data[recordOffset : recordOffset+recordSize]
		record, valid := unmarshalPlaceRecordG2WithValidation(recordData, opts)

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

// unmarshalPlaceRecordGen1WithValidation parses and validates a Gen1 place record.
func unmarshalPlaceRecordGen1WithValidation(data []byte, opts dd.UnmarshalOptions) (*ddv1.PlaceRecord, bool) {
	record, err := opts.UnmarshalPlaceRecord(data)
	if err != nil {
		// If parsing fails, treat as invalid but keep raw data
		invalidRecord := &ddv1.PlaceRecord{}
		invalidRecord.SetValid(false)
		invalidRecord.SetRawData(data)
		return invalidRecord, false
	}

	// Validate the record
	if !isValidPlaceRecordGen1(record) {
		record.SetValid(false)
		return record, false
	}

	record.SetValid(true)
	return record, true
}

// unmarshalPlaceRecordG2WithValidation parses and validates a Gen2 place record.
func unmarshalPlaceRecordG2WithValidation(data []byte, opts dd.UnmarshalOptions) (*ddv1.PlaceRecordG2, bool) {
	record, err := opts.UnmarshalPlaceRecordG2(data)
	if err != nil {
		// If parsing fails, treat as invalid but keep raw data
		invalidRecord := &ddv1.PlaceRecordG2{}
		invalidRecord.SetValid(false)
		invalidRecord.SetRawData(data)
		return invalidRecord, false
	}

	// Validate the record
	if !isValidPlaceRecordG2(record) {
		record.SetValid(false)
		return record, false
	}

	record.SetValid(true)
	return record, true
}

// isValidPlaceRecordGen1 validates a Gen1 place record for reasonable values.
func isValidPlaceRecordGen1(record *ddv1.PlaceRecord) bool {
	// Check timestamp validity
	entryTime := record.GetEntryTime()
	if entryTime != nil {
		year := entryTime.AsTime().Year()
		if year < 1980 || year > 2050 {
			return false
		}
	}

	// Check odometer value
	odometer := record.GetVehicleOdometerKm()
	if odometer < 0 || odometer > 10000000 {
		return false
	}

	return true
}

// isValidPlaceRecordG2 validates a Gen2 place record for reasonable values.
func isValidPlaceRecordG2(record *ddv1.PlaceRecordG2) bool {
	// Check timestamp validity
	entryTime := record.GetEntryTime()
	if entryTime != nil {
		year := entryTime.AsTime().Year()
		if year < 1980 || year > 2050 {
			return false
		}
	}

	// Check odometer value
	odometer := record.GetVehicleOdometerKm()
	if odometer < 0 || odometer > 10000000 {
		return false
	}

	return true
}

// appendPlaces appends the binary representation of Places to dst.
func appendPlaces(dst []byte, p *cardv1.Places, generation ddv1.Generation) ([]byte, error) {
	if p == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(p.GetNewestRecordIndex()))

	var err error

	// Append Gen1 or Gen2 records based on which is populated
	if generation == ddv1.Generation_GENERATION_2 {
		for _, rec := range p.GetRecordsG2() {
			dst, err = appendPlaceRecordG2(dst, rec)
			if err != nil {
				return nil, err
			}
		}
	} else {
		for _, rec := range p.GetRecords() {
			dst, err = appendPlaceRecordGen1(dst, rec)
			if err != nil {
				return nil, err
			}
		}
	}

	// Append trailing bytes for roundtrip accuracy
	if trailingBytes := p.GetTrailingBytes(); len(trailingBytes) > 0 {
		dst = append(dst, trailingBytes...)
	}

	return dst, nil
}

// appendPlaceRecordGen1 appends a Gen1 place record (10 bytes).
func appendPlaceRecordGen1(dst []byte, rec *ddv1.PlaceRecord) ([]byte, error) {
	const recordSize = placeRecordSizeGen1 // Fixed size - no conditionals!

	if rec == nil || !rec.GetValid() {
		if raw := rec.GetRawData(); len(raw) > 0 {
			return append(dst, raw...), nil
		}
		return append(dst, make([]byte, recordSize)...), nil
	}

	return dd.AppendPlaceRecord(dst, rec)
}

// appendPlaceRecordG2 appends a Gen2 place record (21 bytes).
func appendPlaceRecordG2(dst []byte, rec *ddv1.PlaceRecordG2) ([]byte, error) {
	const recordSize = placeRecordSizeGen2 // Fixed size - no conditionals!

	if rec == nil || !rec.GetValid() {
		if raw := rec.GetRawData(); len(raw) > 0 {
			return append(dst, raw...), nil
		}
		return append(dst, make([]byte, recordSize)...), nil
	}

	return dd.AppendPlaceRecordG2(dst, rec)
}
