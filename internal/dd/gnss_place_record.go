package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGNSSPlaceRecord unmarshals a GNSSPlaceRecord from binary data.
//
// The data type `GNSSPlaceRecord` is specified in the Data Dictionary, Section 2.80.
//
// ASN.1 Definition:
//
//	GNSSPlaceRecord ::= SEQUENCE {
//	    timeStamp TimeReal,
//	    gnssAccuracy GNSSAccuracy,
//	    geoCoordinates GeoCoordinates
//	}
//
// Binary Layout (11 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (6 bytes)
//
// This format is used consistently across all contexts (VU downloads, card files, etc.).
func (opts UnmarshalOptions) UnmarshalGNSSPlaceRecord(data []byte) (*ddv1.GNSSPlaceRecord, error) {
	const (
		lenGNSSPlaceRecord = 11
		idxTimestamp       = 0
		idxAccuracy        = 4
		idxGeoCoords       = 5
	)

	if len(data) != lenGNSSPlaceRecord {
		return nil, fmt.Errorf("invalid data length for GNSSPlaceRecord: got %d, want %d", len(data), lenGNSSPlaceRecord)
	}

	record := &ddv1.GNSSPlaceRecord{}

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[idxTimestamp : idxTimestamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse gnssAccuracy (1 byte)
	accuracy := int32(data[idxAccuracy])
	record.SetGnssAccuracy(accuracy)

	// Parse geoCoordinates (6 bytes)
	geoCoords, err := opts.UnmarshalGeoCoordinates(data[idxGeoCoords : idxGeoCoords+6])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo coordinates: %w", err)
	}
	record.SetGeoCoordinates(geoCoords)

	return record, nil
}

// AppendGNSSPlaceRecord appends a GNSSPlaceRecord to dst.
//
// The data type `GNSSPlaceRecord` is specified in the Data Dictionary, Section 2.80.
//
// ASN.1 Definition:
//
//	GNSSPlaceRecord ::= SEQUENCE {
//	    timeStamp TimeReal,
//	    gnssAccuracy GNSSAccuracy,
//	    geoCoordinates GeoCoordinates
//	}
//
// Binary Layout (11 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (6 bytes)
//
// This format is used consistently across all contexts (VU downloads, card files, etc.).
func AppendGNSSPlaceRecord(dst []byte, record *ddv1.GNSSPlaceRecord) ([]byte, error) {
	if record == nil {
		// Append 11 zero bytes if no GNSS data
		return append(dst, make([]byte, 11)...), nil
	}

	// Append timestamp (TimeReal - 4 bytes)
	var err error
	dst, err = AppendTimeReal(dst, record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to append timestamp: %w", err)
	}

	// Append gnssAccuracy (1 byte)
	accuracy := record.GetGnssAccuracy()
	if accuracy < 0 || accuracy > 255 {
		return nil, fmt.Errorf("invalid GNSS accuracy: %d (must be 0-255)", accuracy)
	}
	dst = append(dst, byte(accuracy))

	// Append geoCoordinates (6 bytes)
	dst, err = AppendGeoCoordinates(dst, record.GetGeoCoordinates())
	if err != nil {
		return nil, fmt.Errorf("failed to append geo coordinates: %w", err)
	}

	return dst, nil
}
