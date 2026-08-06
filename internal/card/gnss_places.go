package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalGnssPlaces unmarshals GNSS places data from a card EF.
//
// The data type `GNSSAccumulatedDriving` is specified in the Data Dictionary, Section 2.78.
//
// ASN.1 Definition:
//
//	GNSSAccumulatedDriving ::= SEQUENCE {
//	    gnssADPointerNewestRecord          INTEGER(0..NoOfGNSSADRecords-1),
//	    gnssAccumulatedDrivingRecords      SET SIZE(NoOfGNSSADRecords) OF GNSSAccumulatedDrivingRecord
//	}
//
//	GNSSAccumulatedDrivingRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,             -- 4 bytes
//	    gnssPlaceRecord                    GNSSPlaceRecord,      -- 11 bytes
//	    vehicleOdometerValue               OdometerShort         -- 3 bytes
//	}
//
// Binary structure:
//   - 2 bytes: gnssADPointerNewestRecord
//   - N * 18 bytes: fixed-size array of GNSSAccumulatedDrivingRecord (N determined by data length)
//
// Typical sizes:
//   - Control/Company cards: 434 bytes (24 records)
func (opts UnmarshalOptions) unmarshalGnssPlaces(data []byte) (*cardv1.GnssPlaces, error) {
	const (
		idxNewestRecordIndex             = 0
		lenNewestRecordIndex             = 2
		lenGNSSAccumulatedDrivingRecord  = 18
		lenGNSSAccumulatedDrivingMinimum = lenNewestRecordIndex
	)

	if len(data) < lenGNSSAccumulatedDrivingMinimum {
		return nil, fmt.Errorf("invalid data length for GNSSAccumulatedDriving: got %d, want at least %d", len(data), lenGNSSAccumulatedDrivingMinimum)
	}

	// Validate that the records section is a multiple of record size
	recordsDataLen := len(data) - lenNewestRecordIndex
	if recordsDataLen%lenGNSSAccumulatedDrivingRecord != 0 {
		return nil, fmt.Errorf("invalid records data length for GNSSAccumulatedDriving: got %d bytes, not a multiple of %d", recordsDataLen, lenGNSSAccumulatedDrivingRecord)
	}

	var target cardv1.GnssPlaces

	// Parse newest record index
	newestRecordIndex := binary.BigEndian.Uint16(data[idxNewestRecordIndex:])
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse records using bufio.Scanner pattern
	records, err := opts.unmarshalGNSSAccumulatedDrivingRecords(data[lenNewestRecordIndex:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse GNSS accumulated driving records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// splitGNSSAccumulatedDrivingRecord is a bufio.SplitFunc for parsing GNSSAccumulatedDrivingRecord entries.
func splitGNSSAccumulatedDrivingRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const lenGNSSAccumulatedDrivingRecord = 18

	if len(data) < lenGNSSAccumulatedDrivingRecord {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return lenGNSSAccumulatedDrivingRecord, data[:lenGNSSAccumulatedDrivingRecord], nil
}

// unmarshalGNSSAccumulatedDrivingRecords parses the fixed-size array of GNSS accumulated driving records.
func (opts UnmarshalOptions) unmarshalGNSSAccumulatedDrivingRecords(data []byte) ([]*cardv1.GnssPlaces_Record, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(splitGNSSAccumulatedDrivingRecord)

	var records []*cardv1.GnssPlaces_Record
	for scanner.Scan() {
		record, err := opts.unmarshalGNSSAccumulatedDrivingRecord(scanner.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal GNSS accumulated driving record: %w", err)
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return records, nil
}

// unmarshalGNSSAccumulatedDrivingRecord unmarshals a single GNSSAccumulatedDrivingRecord.
//
// Binary structure (18 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 11 bytes: GNSSPlaceRecord
//   - 3 bytes: OdometerShort vehicleOdometerValue
func (opts UnmarshalOptions) unmarshalGNSSAccumulatedDrivingRecord(data []byte) (*cardv1.GnssPlaces_Record, error) {
	const (
		idxTimeStamp       = 0
		idxGnssPlaceRecord = 4
		idxVehicleOdometer = 15
		lenRecord          = 18
	)

	if len(data) != lenRecord {
		return nil, fmt.Errorf("invalid data length for GNSSAccumulatedDrivingRecord: got %d, want %d", len(data), lenRecord)
	}

	var record cardv1.GnssPlaces_Record

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse GNSS place record (11 bytes)
	gnssPlaceRecord, err := opts.UnmarshalGNSSPlaceRecord(data[idxGnssPlaceRecord : idxGnssPlaceRecord+11])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GNSS place record: %w", err)
	}
	record.SetGnssPlaceRecord(gnssPlaceRecord)

	// Parse vehicle odometer (OdometerShort - 3 bytes)
	odometer, err := opts.UnmarshalOdometer(data[idxVehicleOdometer : idxVehicleOdometer+3])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle odometer: %w", err)
	}
	if odometer != nil {
		record.SetVehicleOdometerKm(*odometer)
	}

	return &record, nil
}

// MarshalCardGnssPlaces marshals GNSS places data.
//
// The data type `GNSSAccumulatedDriving` is specified in the Data Dictionary, Section 2.78.
//
// ASN.1 Definition:
//
//	GNSSAccumulatedDriving ::= SEQUENCE {
//	    gnssADPointerNewestRecord          INTEGER(0..NoOfGNSSADRecords-1),
//	    gnssAccumulatedDrivingRecords      SET SIZE(NoOfGNSSADRecords) OF GNSSAccumulatedDrivingRecord
//	}
//
//	GNSSAccumulatedDrivingRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,             -- 4 bytes
//	    gnssPlaceRecord                    GNSSPlaceRecord,      -- 11 bytes
//	    vehicleOdometerValue               OdometerShort         -- 3 bytes
//	}
//
// Binary structure:
//   - 2 bytes: gnssADPointerNewestRecord
//   - N * 18 bytes: fixed-size array of GNSSAccumulatedDrivingRecord
//
// The number of records (N) is determined from the original data, not explicitly stored.
func (opts MarshalOptions) MarshalCardGnssPlaces(gnssPlaces *cardv1.GnssPlaces) ([]byte, error) {
	if gnssPlaces == nil {
		return nil, nil
	}

	var dst []byte

	// Append newest record index (2 bytes)
	newestRecordIndex := gnssPlaces.GetNewestRecordIndex()
	dst = binary.BigEndian.AppendUint16(dst, uint16(newestRecordIndex))

	// Append all GNSS accumulated driving records
	// The binary format is a fixed-size array, so we write exactly what we have
	records := gnssPlaces.GetRecords()
	for _, record := range records {
		recordBytes, err := opts.MarshalGNSSAccumulatedDrivingRecord(record)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal GNSS accumulated driving record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// MarshalGNSSAccumulatedDrivingRecord marshals a single GNSS accumulated driving record.
//
// Binary structure (18 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 11 bytes: GNSSPlaceRecord
//   - 3 bytes: OdometerShort vehicleOdometerValue
func (opts MarshalOptions) MarshalGNSSAccumulatedDrivingRecord(record *cardv1.GnssPlaces_Record) ([]byte, error) {
	if record == nil {
		return nil, nil
	}

	var dst []byte

	// Append timestamp (TimeReal - 4 bytes)

	timestampBytes, err := opts.MarshalTimeReal(record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal timestamp: %w", err)
	}
	dst = append(dst, timestampBytes...)

	// Append GNSS place record (11 bytes)
	gnssBytes, err := opts.MarshalGNSSPlaceRecord(record.GetGnssPlaceRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GNSS place record: %w", err)
	}
	dst = append(dst, gnssBytes...)

	// Append vehicle odometer (OdometerShort - 3 bytes)
	odometer := record.GetVehicleOdometerKm()
	if odometer < 0 || odometer > 999999 {
		return nil, fmt.Errorf("invalid vehicle odometer value: %d", odometer)
	}
	odometerBytes, err := opts.MarshalOdometer(odometer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle odometer: %w", err)
	}
	dst = append(dst, odometerBytes...)

	return dst, nil
}

// anonymizeGnssPlaces creates an anonymized copy of GnssPlaces,
// replacing sensitive information with static, deterministic test values.
func (opts AnonymizeOptions) anonymizeGnssPlaces(gnssPlaces *cardv1.GnssPlaces) *cardv1.GnssPlaces {
	if gnssPlaces == nil {
		return nil
	}

	result := &cardv1.GnssPlaces{}

	// Preserve the pointer to newest record
	result.SetNewestRecordIndex(gnssPlaces.GetNewestRecordIndex())

	// Anonymize each record
	originalRecords := gnssPlaces.GetRecords()
	anonymizedRecords := make([]*cardv1.GnssPlaces_Record, len(originalRecords))
	for i, record := range originalRecords {
		anonymizedRecords[i] = opts.anonymizeGNSSAccumulatedDrivingRecord(record, i)
	}
	result.SetRecords(anonymizedRecords)

	// Signature field left unset (nil) - TLV marshaller will omit the signature block

	return result
}

// anonymizeGNSSAccumulatedDrivingRecord anonymizes a single GNSS accumulated driving record.
// Uses index to create sequential timestamps.
func (opts AnonymizeOptions) anonymizeGNSSAccumulatedDrivingRecord(record *cardv1.GnssPlaces_Record, index int) *cardv1.GnssPlaces_Record {
	if record == nil {
		// Return a zero-filled record for nil entries
		zeroRecord := &cardv1.GnssPlaces_Record{}
		zeroRecord.SetVehicleOdometerKm(0)
		return zeroRecord
	}

	result := &cardv1.GnssPlaces_Record{}

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Replace outer timestamp with sequential test timestamps
	// Base: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	// Increment by 1 hour per record
	baseEpoch := int64(1577836800)
	ts := record.GetTimestamp()
	if ts != nil && ts.Seconds != 0 {
		// Non-zero timestamp - replace with sequential test timestamp
		result.SetTimestamp(&timestamppb.Timestamp{
			Seconds: baseEpoch + int64(index)*3600,
			Nanos:   0,
		})
	}
	// else: Zero timestamp - leave unset (nil)

	// Anonymize GNSS place record (replaces coordinates with Helsinki)
	gnssPlaceRecord := record.GetGnssPlaceRecord()
	if gnssPlaceRecord != nil {
		anonymizedGnssPlace := ddOpts.AnonymizeGNSSPlaceRecord(gnssPlaceRecord)
		// Update the inner timestamp to match the outer one (for consistency)
		if result.GetTimestamp() != nil {
			anonymizedGnssPlace.SetTimestamp(result.GetTimestamp())
		}
		result.SetGnssPlaceRecord(anonymizedGnssPlace)
	}

	// Round odometer to nearest 100km (preserves magnitude but not exact correlation)
	originalOdometer := record.GetVehicleOdometerKm()
	roundedOdometer := (originalOdometer / 100) * 100
	result.SetVehicleOdometerKm(roundedOdometer)

	return result
}
