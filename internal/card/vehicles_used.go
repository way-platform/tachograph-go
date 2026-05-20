package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalVehiclesUsed unmarshals vehicles used data from a Gen1 card EF.
func (opts UnmarshalOptions) unmarshalVehiclesUsed(data []byte) (*cardv1.VehiclesUsed, error) {
	const (
		lenMinEfVehiclesUsed = 2 // Minimum EF_Vehicles_Used record size
	)

	if len(data) < lenMinEfVehiclesUsed {
		return nil, fmt.Errorf("insufficient data for vehicles used: got %d bytes, need at least %d", len(data), lenMinEfVehiclesUsed)
	}

	var target cardv1.VehiclesUsed

	// Save complete raw data for painting
	target.SetRawData(data)

	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse Gen1 vehicle records (31 bytes each)
	records, err := opts.unmarshalVehicleRecordsGen1(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gen1 vehicle records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// unmarshalVehicleRecordsGen1 parses Gen1 vehicle records (31 bytes each).
func (opts UnmarshalOptions) unmarshalVehicleRecordsGen1(r *bytes.Reader) ([]*ddv1.CardVehicleRecord, error) {
	const lenCardVehicleRecord = 31

	var records []*ddv1.CardVehicleRecord
	for r.Len() >= lenCardVehicleRecord {
		recordBytes := make([]byte, lenCardVehicleRecord)
		if _, err := r.Read(recordBytes); err != nil {
			break // Stop parsing on error, but return what we have
		}

		record, err := opts.UnmarshalCardVehicleRecord(recordBytes)
		if err != nil {
			return records, fmt.Errorf("failed to parse Gen1 vehicle record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// MarshalVehiclesUsed marshals Gen1 vehicles used data.
func (opts MarshalOptions) MarshalVehiclesUsed(vehiclesUsed *cardv1.VehiclesUsed) ([]byte, error) {
	if vehiclesUsed == nil {
		return nil, nil
	}

	// Calculate expected size: 2 bytes (pointer) + N records × 31 bytes
	const recordSize = 31
	numRecords := len(vehiclesUsed.GetRecords())
	expectedSize := 2 + (numRecords * recordSize)

	// Use raw_data as canvas if available and correct size
	if rawData := vehiclesUsed.GetRawData(); len(rawData) == expectedSize {
		// Make a copy to use as canvas
		canvas := make([]byte, expectedSize)
		copy(canvas, rawData)

		// Paint newest record index over canvas
		binary.BigEndian.PutUint16(canvas[0:2], uint16(vehiclesUsed.GetNewestRecordIndex()))

		// Paint each record over canvas

		offset := 2
		for _, record := range vehiclesUsed.GetRecords() {
			recordBytes, err := opts.MarshalCardVehicleRecord(record)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Gen1 vehicle record: %w", err)
			}
			if len(recordBytes) != recordSize {
				return nil, fmt.Errorf("invalid Gen1 vehicle record size: got %d, want %d", len(recordBytes), recordSize)
			}
			copy(canvas[offset:offset+recordSize], recordBytes)
			offset += recordSize
		}

		return canvas, nil
	}

	// Fall back to building from scratch
	var dst []byte
	newestRecordIndex := uint16(vehiclesUsed.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	for _, record := range vehiclesUsed.GetRecords() {
		recordBytes, err := opts.MarshalCardVehicleRecord(record)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Gen1 vehicle record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// AnonymizeVehiclesUsed creates an anonymized copy, replacing sensitive data
// with static, deterministic test values while preserving structure.
//
// Anonymization strategy:
// - Vehicle registrations: Replaced with "TEST-VRN"
// - Timestamps: Static base (2020-01-01 00:00:00) + 1 day per record
// - Odometer readings: Rounded to nearest 1000km
// - Countries: Preserved (structural info)
// - Pointer: Preserved (structural info)
// - VU counters: Preserved (structural info)
func (opts AnonymizeOptions) anonymizeVehiclesUsed(v *cardv1.VehiclesUsed) *cardv1.VehiclesUsed {
	if v == nil {
		return nil
	}

	result := &cardv1.VehiclesUsed{}

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Preserve pointer (structural info)
	result.SetNewestRecordIndex(v.GetNewestRecordIndex())

	// Anonymize records
	var anonymizedRecords []*ddv1.CardVehicleRecord
	for i, record := range v.GetRecords() {
		anonymizedRecord := ddOpts.AnonymizeCardVehicleRecord(record, i)
		anonymizedRecords = append(anonymizedRecords, anonymizedRecord)
	}
	result.SetRecords(anonymizedRecords)

	// Signature and raw_data fields left unset (nil) - TLV marshaller will omit these blocks

	return result
}
