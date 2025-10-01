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
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Create dd.UnmarshalOptions from card-level UnmarshalOptions
	ddOpts := dd.UnmarshalOptions{
		Generation: opts.Generation,
		Version:    opts.Version,
	}

	// Parse Gen1 vehicle records (31 bytes each)
	records, err := parseVehicleRecordsGen1(r, ddOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gen1 vehicle records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// unmarshalVehiclesUsedG2 unmarshals vehicles used data from a Gen2 card EF.
func (opts UnmarshalOptions) unmarshalVehiclesUsedG2(data []byte) (*cardv1.VehiclesUsedG2, error) {
	const (
		lenMinEfVehiclesUsed = 2 // Minimum EF_Vehicles_Used record size
	)

	if len(data) < lenMinEfVehiclesUsed {
		return nil, fmt.Errorf("insufficient data for vehicles used: got %d bytes, need at least %d", len(data), lenMinEfVehiclesUsed)
	}

	var target cardv1.VehiclesUsedG2
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Create dd.UnmarshalOptions from card-level UnmarshalOptions
	ddOpts := dd.UnmarshalOptions{
		Generation: opts.Generation,
		Version:    opts.Version,
	}

	// Parse Gen2 vehicle records (48 bytes each)
	records, err := parseVehicleRecordsGen2(r, ddOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gen2 vehicle records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// parseVehicleRecordsGen1 parses Gen1 vehicle records (31 bytes each).
func parseVehicleRecordsGen1(r *bytes.Reader, opts dd.UnmarshalOptions) ([]*ddv1.CardVehicleRecord, error) {
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

// parseVehicleRecordsGen2 parses Gen2 vehicle records (48 bytes each).
func parseVehicleRecordsGen2(r *bytes.Reader, opts dd.UnmarshalOptions) ([]*ddv1.CardVehicleRecordG2, error) {
	const lenCardVehicleRecord = 48

	var records []*ddv1.CardVehicleRecordG2
	for r.Len() >= lenCardVehicleRecord {
		recordBytes := make([]byte, lenCardVehicleRecord)
		if _, err := r.Read(recordBytes); err != nil {
			break // Stop parsing on error, but return what we have
		}

		record, err := opts.UnmarshalCardVehicleRecordG2(recordBytes)
		if err != nil {
			return records, fmt.Errorf("failed to parse Gen2 vehicle record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// appendVehiclesUsed appends Gen1 vehicles used data for marshalling.
func appendVehiclesUsed(dst []byte, vehiclesUsed *cardv1.VehiclesUsed) ([]byte, error) {
	if vehiclesUsed == nil {
		return dst, nil
	}

	// Write newest record index (2 bytes)
	newestRecordIndex := uint16(vehiclesUsed.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	// Write Gen1 records (31 bytes each)
	for _, record := range vehiclesUsed.GetRecords() {
		recordBytes, err := dd.AppendCardVehicleRecord(nil, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append Gen1 vehicle record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// appendVehiclesUsedG2 appends Gen2 vehicles used data for marshalling.
func appendVehiclesUsedG2(dst []byte, vehiclesUsed *cardv1.VehiclesUsedG2) ([]byte, error) {
	if vehiclesUsed == nil {
		return dst, nil
	}

	// Write newest record index (2 bytes)
	newestRecordIndex := uint16(vehiclesUsed.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	// Write Gen2 records (48 bytes each)
	for _, record := range vehiclesUsed.GetRecords() {
		recordBytes, err := dd.AppendCardVehicleRecordG2(nil, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append Gen2 vehicle record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}
