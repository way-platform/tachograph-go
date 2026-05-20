package card

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalPlacesG2 unmarshals the EF_Places data (Gen2 format).
//
// Gen2 Structure (TCS_152):
// - placePointerNewestRecord: 2 bytes (not 1!)
// - placeRecords: N × 21 bytes (112 records for driver cards)
func (opts UnmarshalOptions) unmarshalPlacesG2(data []byte) (*cardv1.PlacesG2, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least 2", len(data))
	}

	target := &cardv1.PlacesG2{}

	// Save complete raw data for painting
	target.SetRawData(data)

	// Read the newest record index (2 bytes for Gen2)
	newestRecordIndex := binary.BigEndian.Uint16(data[0:2])
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Remaining data contains the circular buffer of place records
	remainingData := data[2:]

	// Parse Gen2 records (21 bytes each)
	records, _ := opts.unmarshalCircularPlaceRecordsG2(remainingData, int(newestRecordIndex))
	target.SetRecords(records)

	return target, nil
}

// unmarshalCircularPlaceRecordsG2 parses place records from a circular buffer (Gen2: 21 bytes each).
func (opts UnmarshalOptions) unmarshalCircularPlaceRecordsG2(data []byte, newestIndex int) ([]*ddv1.PlaceRecordG2, []byte) {
	const recordSize = 21
	numFullRecords := len(data) / recordSize
	trailingBytes := data[numFullRecords*recordSize:]

	records := make([]*ddv1.PlaceRecordG2, 0, numFullRecords)

	for i := 0; i < numFullRecords; i++ {
		start := i * recordSize
		end := start + recordSize
		recordData := data[start:end]

		record, err := opts.UnmarshalPlaceRecordG2(recordData)
		if err != nil {
			// Mark record as invalid on parse error
			record = &ddv1.PlaceRecordG2{}
			record.SetValid(false)
			record.SetRawData(recordData)
		}

		records = append(records, record)
	}

	return records, trailingBytes
}

// MarshalPlacesG2 marshals the EF_Places data (Gen2 format).
//
// Gen2 Structure (TCS_152):
// - placePointerNewestRecord: 2 bytes (not 1!)
// - placeRecords: N × 21 bytes
func (opts MarshalOptions) MarshalPlacesG2(p *cardv1.PlacesG2) ([]byte, error) {
	if p == nil {
		return nil, nil
	}

	// Calculate expected size: 2 bytes (pointer) + N records × 21 bytes
	const recordSize = 21
	numRecords := len(p.GetRecords())
	expectedSize := 2 + (numRecords * recordSize)

	// Use raw_data as canvas if available and correct size
	if rawData := p.GetRawData(); len(rawData) == expectedSize {
		// Make a copy to use as canvas
		canvas := make([]byte, expectedSize)
		copy(canvas, rawData)

		// Paint newest record index over canvas (2 bytes for Gen2)
		binary.BigEndian.PutUint16(canvas[0:2], uint16(p.GetNewestRecordIndex()))

		// Paint each record over canvas

		offset := 2
		for _, record := range p.GetRecords() {
			recordBytes, err := opts.MarshalPlaceRecordG2(record)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Gen2 place record: %w", err)
			}
			if len(recordBytes) != recordSize {
				return nil, fmt.Errorf("invalid Gen2 place record size: got %d, want %d", len(recordBytes), recordSize)
			}
			copy(canvas[offset:offset+recordSize], recordBytes)
			offset += recordSize
		}

		return canvas, nil
	}

	// Fall back to building from scratch
	var dst []byte
	newestRecordIndex := uint16(p.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	for _, record := range p.GetRecords() {
		recordBytes, err := opts.MarshalPlaceRecordG2(record)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Gen2 place record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// AnonymizePlacesG2 creates an anonymized copy of PlacesG2 (Gen2), replacing potentially
// sensitive location data (including GNSS coordinates) while preserving the structure
// for testing.
//
// Timestamp anonymization strategy:
// - Replaces all timestamps (entry time + GNSS timestamp) with static test values
// - Base: 2020-01-01 00:00:00, incremented by 1 hour per record
// - Deterministic: same input always produces same output
// - Maintains record ordering for parsing tests
func (opts AnonymizeOptions) anonymizePlacesG2(p *cardv1.PlacesG2) *cardv1.PlacesG2 {
	if p == nil {
		return nil
	}

	result := &cardv1.PlacesG2{}

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Preserve structural metadata
	result.SetNewestRecordIndex(p.GetNewestRecordIndex())

	// Anonymize each record (timestamps anonymized below)
	var anonymizedRecords []*ddv1.PlaceRecordG2
	for _, record := range p.GetRecords() {
		anonymizedRecords = append(anonymizedRecords, ddOpts.AnonymizePlaceRecordG2(record))
	}
	result.SetRecords(anonymizedRecords)

	// Apply dataset-specific timestamp normalization (unless timestamps are preserved)
	// This makes the anonymization non-reversible
	if !opts.PreserveTimestamps {
		anonymizeTimestampsInPlaceG2(anonymizedRecords)
	}

	// Signature and raw_data fields left unset (nil) - TLV marshaller will omit these blocks

	return result
}

// anonymizeTimestampsInPlaceG2 replaces all timestamps with static test values.
// Uses a fixed base timestamp (2020-01-01 00:00:00) and increments by 1 hour per record
// to maintain ordering while providing deterministic, anonymized test data.
func anonymizeTimestampsInPlaceG2(records []*ddv1.PlaceRecordG2) {
	if len(records) == 0 {
		return
	}

	// Test epoch: 2020-01-01 00:00:00 UTC
	const testEpoch = int64(1577836800)
	const oneHour = int64(3600)

	// Replace all timestamps with static incremented values
	for i, record := range records {
		// Set main entry time: base + (i * 1 hour)
		staticTimestamp := testEpoch + (int64(i) * oneHour)
		if record.GetEntryTime() != nil || i == 0 {
			record.SetEntryTime(&timestamppb.Timestamp{
				Seconds: staticTimestamp,
				Nanos:   0,
			})
		}

		// Set GNSS timestamp if present (same as main timestamp for simplicity)
		if gnss := record.GetEntryGnssPlaceRecord(); gnss != nil {
			if gnss.GetTimestamp() != nil {
				gnss.SetTimestamp(&timestamppb.Timestamp{
					Seconds: staticTimestamp,
					Nanos:   0,
				})
			}
		}
	}
}
