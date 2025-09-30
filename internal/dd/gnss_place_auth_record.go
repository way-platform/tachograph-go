package dd

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGNSSPlaceAuthRecord unmarshals a GNSSPlaceAuthRecord from binary data.
//
// The data type `GNSSPlaceAuthRecord` is specified in the Data Dictionary, Section 2.79c.
//
// ASN.1 Definition:
//
//	GNSSPlaceAuthRecord ::= SEQUENCE {
//	    timeStamp TimeReal,
//	    gnssAccuracy GNSSAccuracy,
//	    geoCoordinates GeoCoordinates,
//	    authenticationStatus PositionAuthenticationStatus
//	}
func UnmarshalGNSSPlaceAuthRecord(data []byte) (*ddv1.GNSSPlaceAuthRecord, error) {
	const (
		lenGNSSPlaceAuthRecord = 14
		idxTimestamp           = 0
		idxAccuracy            = 4
		idxGeoCoords           = 5
		idxAuthStatus          = 13
	)

	if len(data) != lenGNSSPlaceAuthRecord {
		return nil, fmt.Errorf("invalid data length for GNSSPlaceAuthRecord: got %d, want %d", len(data), lenGNSSPlaceAuthRecord)
	}

	record := &ddv1.GNSSPlaceAuthRecord{}

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := UnmarshalTimeReal(data[idxTimestamp : idxTimestamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse gnssAccuracy (1 byte)
	record.SetGnssAccuracy(int32(data[idxAccuracy]))

	// Parse geoCoordinates (8 bytes)
	geoCoords, err := UnmarshalGeoCoordinates(data[idxGeoCoords : idxGeoCoords+8])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo coordinates: %w", err)
	}
	record.SetGeoCoordinates(geoCoords)

	// Parse authenticationStatus (1 byte)
	SetEnumFromProtocolValueGeneric(
		ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNSPECIFIED.Descriptor(),
		int32(data[idxAuthStatus]),
		func(enumNum protoreflect.EnumNumber) {
			record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus(enumNum))
		},
		nil, // No unrecognized field for this enum
	)

	return record, nil
}

// AppendGNSSPlaceAuthRecord appends a GNSSPlaceAuthRecord to dst.
//
// The data type `GNSSPlaceAuthRecord` is specified in the Data Dictionary, Section 2.79c.
//
// ASN.1 Definition:
//
//	GNSSPlaceAuthRecord ::= SEQUENCE {
//	    timeStamp TimeReal,
//	    gnssAccuracy GNSSAccuracy,
//	    geoCoordinates GeoCoordinates,
//	    authenticationStatus PositionAuthenticationStatus
//	}
func AppendGNSSPlaceAuthRecord(dst []byte, record *ddv1.GNSSPlaceAuthRecord) ([]byte, error) {
	if record == nil {
		// Append 14 zero bytes if no data
		return append(dst, make([]byte, 14)...), nil
	}

	// Append timestamp (TimeReal - 4 bytes)
	timestamp := record.GetTimestamp()
	if timestamp == nil {
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	} else {
		var err error
		dst, err = AppendTimeReal(dst, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to append timestamp: %w", err)
		}
	}

	// Append gnssAccuracy (1 byte)
	dst = append(dst, byte(record.GetGnssAccuracy()))

	// Append geoCoordinates (8 bytes)
	var err error
	dst, err = AppendGeoCoordinates(dst, record.GetGeoCoordinates())
	if err != nil {
		return nil, fmt.Errorf("failed to append geo coordinates: %w", err)
	}

	// Append authenticationStatus (1 byte)
	authStatus, _ := GetProtocolValueForEnum(record.GetAuthenticationStatus())
	dst = append(dst, byte(authStatus))

	return dst, nil
}
