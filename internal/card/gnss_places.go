package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
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
	records, err := parseGNSSAccumulatedDrivingRecords(data[lenNewestRecordIndex:], opts)
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

// parseGNSSAccumulatedDrivingRecords parses the fixed-size array of GNSS accumulated driving records.
func parseGNSSAccumulatedDrivingRecords(data []byte, opts UnmarshalOptions) ([]*cardv1.GnssPlaces_Record, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(splitGNSSAccumulatedDrivingRecord)

	var records []*cardv1.GnssPlaces_Record
	for scanner.Scan() {
		record, err := unmarshalGNSSAccumulatedDrivingRecord(scanner.Bytes(), opts)
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
func unmarshalGNSSAccumulatedDrivingRecord(data []byte, opts UnmarshalOptions) (*cardv1.GnssPlaces_Record, error) {
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
	record.SetVehicleOdometerKm(int32(odometer))

	return &record, nil
}

// appendCardGnssPlaces appends GNSS places data to a byte slice.
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
func appendCardGnssPlaces(dst []byte, gnssPlaces *cardv1.GnssPlaces) ([]byte, error) {
	if gnssPlaces == nil {
		return dst, nil
	}

	// Append newest record index (2 bytes)
	newestRecordIndex := gnssPlaces.GetNewestRecordIndex()
	dst = binary.BigEndian.AppendUint16(dst, uint16(newestRecordIndex))

	// Append all GNSS accumulated driving records
	// The binary format is a fixed-size array, so we write exactly what we have
	records := gnssPlaces.GetRecords()
	for _, record := range records {
		var err error
		dst, err = appendGNSSAccumulatedDrivingRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append GNSS accumulated driving record: %w", err)
		}
	}

	return dst, nil
}

// appendGNSSAccumulatedDrivingRecord appends a single GNSS accumulated driving record to dst.
//
// Binary structure (18 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 11 bytes: GNSSPlaceRecord
//   - 3 bytes: OdometerShort vehicleOdometerValue
func appendGNSSAccumulatedDrivingRecord(dst []byte, record *cardv1.GnssPlaces_Record) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// Append timestamp (TimeReal - 4 bytes)
	var err error
	dst, err = dd.AppendTimeReal(dst, record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to append timestamp: %w", err)
	}

	// Append GNSS place record (11 bytes)
	dst, err = dd.AppendGNSSPlaceRecord(dst, record.GetGnssPlaceRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to append GNSS place record: %w", err)
	}

	// Append vehicle odometer (OdometerShort - 3 bytes)
	odometer := record.GetVehicleOdometerKm()
	if odometer < 0 || odometer > 999999 {
		return nil, fmt.Errorf("invalid vehicle odometer value: %d", odometer)
	}
	dst = dd.AppendOdometer(dst, uint32(odometer))

	return dst, nil
}
