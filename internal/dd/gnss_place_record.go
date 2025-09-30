package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGNSSPlaceRecord unmarshals a GNSSPlaceRecord from binary data using the
// STANDARD DD binary format (13 bytes).
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
// Binary Layout (13 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (8 bytes)
//
// Note: For the card-specific PlaceRecord Gen2 variant (12 bytes, different field order,
// no accuracy), use UnmarshalGNSSPlaceRecordCardVariant instead.
func UnmarshalGNSSPlaceRecord(data []byte) (*ddv1.GNSSPlaceRecord, error) {
	const (
		lenGNSSPlaceRecord = 13
		idxTimestamp       = 0
		idxAccuracy        = 4
		idxGeoCoords       = 5
	)

	if len(data) != lenGNSSPlaceRecord {
		return nil, fmt.Errorf("invalid data length for GNSSPlaceRecord: got %d, want %d", len(data), lenGNSSPlaceRecord)
	}

	record := &ddv1.GNSSPlaceRecord{}

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := UnmarshalTimeReal(data[idxTimestamp : idxTimestamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse gnssAccuracy (1 byte)
	accuracy := int32(data[idxAccuracy])
	record.SetGnssAccuracy(accuracy)

	// Parse geoCoordinates (8 bytes)
	geoCoords, err := UnmarshalGeoCoordinates(data[idxGeoCoords : idxGeoCoords+8])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo coordinates: %w", err)
	}
	record.SetGeoCoordinates(geoCoords)

	return record, nil
}

// AppendGNSSPlaceRecord appends a GNSSPlaceRecord to dst using the
// STANDARD DD binary format (13 bytes).
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
// Binary Layout (13 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (8 bytes)
//
// Note: For the card-specific PlaceRecord Gen2 variant (12 bytes, different field order,
// no accuracy), use AppendGNSSPlaceRecordCardVariant instead.
func AppendGNSSPlaceRecord(dst []byte, record *ddv1.GNSSPlaceRecord) ([]byte, error) {
	if record == nil {
		// Append 13 zero bytes if no GNSS data
		return append(dst, make([]byte, 13)...), nil
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

	// Append geoCoordinates (8 bytes)
	dst, err = AppendGeoCoordinates(dst, record.GetGeoCoordinates())
	if err != nil {
		return nil, fmt.Errorf("failed to append geo coordinates: %w", err)
	}

	return dst, nil
}

// UnmarshalGNSSPlaceRecordCardVariant unmarshals a GNSSPlaceRecord from the
// CARD-SPECIFIC binary format used in PlaceRecord Gen2 (12 bytes).
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
// IMPORTANT: Card Binary Layout (12 bytes, DIFFERENT from standard DD format):
//
//	Offset 0: geoCoordinates (8 bytes)  ← DIFFERENT ORDER than standard
//	Offset 8: timeStamp (4 bytes)        ← DIFFERENT ORDER than standard
//	gnssAccuracy field is OMITTED        ← NOT PRESENT in card format
//
// This card-specific format is embedded within PlaceRecord_G2 (Gen2 card places).
// For the standard DD format (13 bytes), use UnmarshalGNSSPlaceRecord instead.
func UnmarshalGNSSPlaceRecordCardVariant(data []byte) (*ddv1.GNSSPlaceRecord, error) {
	const (
		lenGNSSPlaceRecordCard = 12
		idxGeoCoords           = 0
		idxTimestamp           = 8
	)

	if len(data) != lenGNSSPlaceRecordCard {
		return nil, fmt.Errorf("invalid data length for GNSSPlaceRecord (card variant): got %d, want %d", len(data), lenGNSSPlaceRecordCard)
	}

	record := &ddv1.GNSSPlaceRecord{}

	// Parse geoCoordinates (8 bytes)
	geoCoords, err := UnmarshalGeoCoordinates(data[idxGeoCoords : idxGeoCoords+8])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo coordinates: %w", err)
	}
	record.SetGeoCoordinates(geoCoords)

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := UnmarshalTimeReal(data[idxTimestamp : idxTimestamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Note: gnssAccuracy is not present in the card file format

	return record, nil
}

// AppendGNSSPlaceRecordCardVariant appends a GNSSPlaceRecord to dst using the
// CARD-SPECIFIC binary format used in PlaceRecord Gen2 (12 bytes).
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
// IMPORTANT: Card Binary Layout (12 bytes, DIFFERENT from standard DD format):
//
//	Offset 0: geoCoordinates (8 bytes)  ← DIFFERENT ORDER than standard
//	Offset 8: timeStamp (4 bytes)        ← DIFFERENT ORDER than standard
//	gnssAccuracy field is OMITTED        ← NOT PRESENT in card format
//
// This card-specific format is embedded within PlaceRecord_G2 (Gen2 card places).
// For the standard DD format (13 bytes), use AppendGNSSPlaceRecord instead.
func AppendGNSSPlaceRecordCardVariant(dst []byte, record *ddv1.GNSSPlaceRecord) ([]byte, error) {
	if record == nil {
		// Append 12 zero bytes if no GNSS data
		return append(dst, make([]byte, 12)...), nil
	}

	// Append geoCoordinates (8 bytes)
	var err error
	dst, err = AppendGeoCoordinates(dst, record.GetGeoCoordinates())
	if err != nil {
		return nil, fmt.Errorf("failed to append geo coordinates: %w", err)
	}

	// Append timestamp (TimeReal - 4 bytes)
	timestamp := record.GetTimestamp()
	if timestamp == nil {
		// Append zero timestamp if nil
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	} else {
		dst, err = AppendTimeReal(dst, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to append timestamp: %w", err)
		}
	}

	// Note: gnssAccuracy is not written in the card file format

	return dst, nil
}
