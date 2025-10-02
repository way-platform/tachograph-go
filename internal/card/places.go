package card

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalPlaces unmarshals the EF_Places data (Gen1 format).
//
// Gen1 Structure (TCS_150):
// - placePointerNewestRecord: 1 byte (not 2!)
// - placeRecords: N × 10 bytes (84-112 records for driver cards)
func (opts UnmarshalOptions) unmarshalPlaces(data []byte) (*cardv1.Places, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least 1", len(data))
	}

	target := &cardv1.Places{}

	// Save complete raw data for painting
	target.SetRawData(data)

	// Read the newest record index (1 byte for Gen1)
	newestRecordIndex := data[0]
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Remaining data contains the circular buffer of place records
	remainingData := data[1:]

	// Create dd.UnmarshalOptions from card-level UnmarshalOptions
	ddOpts := dd.UnmarshalOptions{
		Generation: opts.Generation,
		Version:    opts.Version,
	}

	// Parse Gen1 records (10 bytes each)
	records, _ := parseCircularPlaceRecordsGen1(remainingData, int(newestRecordIndex), ddOpts)
	target.SetRecords(records)

	return target, nil
}

// parseCircularPlaceRecordsGen1 parses place records from a circular buffer (Gen1: 10 bytes each).
func parseCircularPlaceRecordsGen1(data []byte, newestIndex int, opts dd.UnmarshalOptions) ([]*ddv1.PlaceRecord, []byte) {
	const recordSize = 10
	numFullRecords := len(data) / recordSize
	trailingBytes := data[numFullRecords*recordSize:]

	records := make([]*ddv1.PlaceRecord, 0, numFullRecords)

	for i := 0; i < numFullRecords; i++ {
		start := i * recordSize
		end := start + recordSize
		recordData := data[start:end]

		record, err := opts.UnmarshalPlaceRecord(recordData)
		if err != nil {
			// Mark record as invalid on parse error
			record = &ddv1.PlaceRecord{}
			record.SetValid(false)
			record.SetRawData(recordData)
		}

		records = append(records, record)
	}

	return records, trailingBytes
}

// appendPlaces marshals the EF_Places data (Gen1 format).
//
// Gen1 Structure (TCS_150):
// - placePointerNewestRecord: 1 byte (not 2!)
// - placeRecords: N × 10 bytes
func appendPlaces(dst []byte, p *cardv1.Places) ([]byte, error) {
	if p == nil {
		return dst, nil
	}

	// Calculate expected size: 1 byte (pointer) + N records × 10 bytes
	const recordSize = 10
	numRecords := len(p.GetRecords())
	expectedSize := 1 + (numRecords * recordSize)

	// Use raw_data as canvas if available and correct size
	if rawData := p.GetRawData(); len(rawData) == expectedSize {
		// Make a copy to use as canvas
		canvas := make([]byte, expectedSize)
		copy(canvas, rawData)

		// Paint newest record index over canvas (1 byte for Gen1)
		canvas[0] = byte(p.GetNewestRecordIndex())

		// Paint each record over canvas
		offset := 1
		for _, record := range p.GetRecords() {
			recordBytes, err := dd.AppendPlaceRecord(nil, record)
			if err != nil {
				return nil, fmt.Errorf("failed to append Gen1 place record: %w", err)
			}
			if len(recordBytes) != recordSize {
				return nil, fmt.Errorf("invalid Gen1 place record size: got %d, want %d", len(recordBytes), recordSize)
			}
			copy(canvas[offset:offset+recordSize], recordBytes)
			offset += recordSize
		}

		return append(dst, canvas...), nil
	}

	// Fall back to building from scratch
	newestRecordIndex := byte(p.GetNewestRecordIndex())
	dst = append(dst, newestRecordIndex)

	for _, record := range p.GetRecords() {
		recordBytes, err := dd.AppendPlaceRecord(nil, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append Gen1 place record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// AnonymizePlaces creates an anonymized copy of Places (Gen1), replacing potentially
// sensitive location data while preserving the structure for testing.
//
// Timestamp anonymization strategy:
// - Replaces all timestamps with static test values
// - Base: 2020-01-01 00:00:00, incremented by 1 hour per record
// - Deterministic: same input always produces same output
// - Maintains record ordering for parsing tests
func AnonymizePlaces(p *cardv1.Places) *cardv1.Places {
	if p == nil {
		return nil
	}

	result := &cardv1.Places{}

	// Preserve structural metadata
	result.SetNewestRecordIndex(p.GetNewestRecordIndex())

	// Anonymize each record (timestamps anonymized below)
	var anonymizedRecords []*ddv1.PlaceRecord
	for _, record := range p.GetRecords() {
		anonymizedRecords = append(anonymizedRecords, dd.AnonymizePlaceRecord(record))
	}
	result.SetRecords(anonymizedRecords)

	// Apply dataset-specific timestamp normalization
	// This makes the anonymization non-reversible
	anonymizeTimestampsInPlace(anonymizedRecords)

	// Regenerate raw_data for each record after timestamp modification
	for _, record := range anonymizedRecords {
		recordBytes, err := dd.AppendPlaceRecord(nil, record)
		if err == nil {
			record.SetRawData(recordBytes)
		}
	}

	// Regenerate raw_data to match anonymized content
	// This ensures round-trip fidelity after anonymization
	anonymizedBytes, err := appendPlaces(nil, result)
	if err == nil {
		result.SetRawData(anonymizedBytes)
	}
	// If marshalling fails, we'll have no raw_data, which is acceptable

	// Don't preserve signature - it will be invalid

	return result
}

// anonymizeTimestampsInPlace replaces all timestamps with static test values.
// Uses a fixed base timestamp (2020-01-01 00:00:00) and increments by 1 hour per record
// to maintain ordering while providing deterministic, anonymized test data.
func anonymizeTimestampsInPlace(records []*ddv1.PlaceRecord) {
	if len(records) == 0 {
		return
	}

	// Test epoch: 2020-01-01 00:00:00 UTC
	const testEpoch = int64(1577836800)
	const oneHour = int64(3600)

	// Replace all timestamps with static incremented values
	for i, record := range records {
		// Set entry time: base + (i * 1 hour)
		staticTimestamp := testEpoch + (int64(i) * oneHour)
		if record.GetEntryTime() != nil || i == 0 {
			record.SetEntryTime(&timestamppb.Timestamp{
				Seconds: staticTimestamp,
				Nanos:   0,
			})
		}
	}
}
