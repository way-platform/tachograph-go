package card

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalPlaces unmarshals the EF_Places data (Gen1 format).
func (opts UnmarshalOptions) unmarshalPlaces(data []byte) (*cardv1.Places, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least 2", len(data))
	}

	target := &cardv1.Places{}

	// Read the newest record index (2 bytes)
	newestRecordIndex := binary.BigEndian.Uint16(data[0:2])
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Remaining data contains the circular buffer of place records
	remainingData := data[2:]

	// Create dd.UnmarshalOptions from card-level UnmarshalOptions
	ddOpts := dd.UnmarshalOptions{
		Generation: opts.Generation,
		Version:    opts.Version,
	}

	// Parse Gen1 records (10 bytes each)
	records, trailingBytes := parseCircularPlaceRecordsGen1(remainingData, int(newestRecordIndex), ddOpts)
	target.SetRecords(records)

	// Store trailing bytes for round-trip fidelity
	if len(trailingBytes) > 0 {
		// Note: Gen1 Places doesn't have trailing_bytes field, we just discard them
		// This is okay since Gen1 cards have fixed record sizes
	}

	return target, nil
}

// unmarshalPlacesG2 unmarshals the EF_Places data (Gen2 format).
func (opts UnmarshalOptions) unmarshalPlacesG2(data []byte) (*cardv1.PlacesG2, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for places: got %d bytes, need at least 2", len(data))
	}

	target := &cardv1.PlacesG2{}

	// Read the newest record index (2 bytes)
	newestRecordIndex := binary.BigEndian.Uint16(data[0:2])
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Remaining data contains the circular buffer of place records
	remainingData := data[2:]

	// Create dd.UnmarshalOptions from card-level UnmarshalOptions
	ddOpts := dd.UnmarshalOptions{
		Generation: opts.Generation,
		Version:    opts.Version,
	}

	// Parse Gen2 records (21 bytes each)
	records, _ := parseCircularPlaceRecordsG2(remainingData, int(newestRecordIndex), ddOpts)
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

// parseCircularPlaceRecordsG2 parses place records from a circular buffer (Gen2: 21 bytes each).
func parseCircularPlaceRecordsG2(data []byte, newestIndex int, opts dd.UnmarshalOptions) ([]*ddv1.PlaceRecordG2, []byte) {
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

// appendPlaces marshals the EF_Places data (Gen1 format).
func appendPlaces(dst []byte, p *cardv1.Places) ([]byte, error) {
	if p == nil {
		return dst, nil
	}

	// Write newest record index (2 bytes)
	newestRecordIndex := uint16(p.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	// Write Gen1 records (10 bytes each)
	for _, record := range p.GetRecords() {
		recordBytes, err := dd.AppendPlaceRecord(nil, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append Gen1 place record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// appendPlacesG2 marshals the EF_Places data (Gen2 format).
func appendPlacesG2(dst []byte, p *cardv1.PlacesG2) ([]byte, error) {
	if p == nil {
		return dst, nil
	}

	// Write newest record index (2 bytes)
	newestRecordIndex := uint16(p.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	// Write Gen2 records (21 bytes each)
	for _, record := range p.GetRecords() {
		recordBytes, err := dd.AppendPlaceRecordG2(nil, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append Gen2 place record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}
