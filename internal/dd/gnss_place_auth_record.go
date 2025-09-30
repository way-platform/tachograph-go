package dd

import (
	"fmt"

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
//
// Binary Layout (12 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (6 bytes)
//	Offset 11: authenticationStatus (1 byte)
func (opts UnmarshalOptions) UnmarshalGNSSPlaceAuthRecord(data []byte) (*ddv1.GNSSPlaceAuthRecord, error) {
	const (
		lenGNSSPlaceAuthRecord = 12
		idxTimestamp           = 0
		idxAccuracy            = 4
		idxGeoCoords           = 5
		idxAuthStatus          = 11
	)

	if len(data) != lenGNSSPlaceAuthRecord {
		return nil, fmt.Errorf("invalid data length for GNSSPlaceAuthRecord: got %d, want %d", len(data), lenGNSSPlaceAuthRecord)
	}

	record := &ddv1.GNSSPlaceAuthRecord{}

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[idxTimestamp : idxTimestamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse gnssAccuracy (1 byte)
	record.SetGnssAccuracy(int32(data[idxAccuracy]))

	// Parse geoCoordinates (6 bytes)
	geoCoords, err := opts.UnmarshalGeoCoordinates(data[idxGeoCoords : idxGeoCoords+6])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo coordinates: %w", err)
	}
	record.SetGeoCoordinates(geoCoords)

	// Parse authenticationStatus (1 byte)
	if authStatus, err := UnmarshalEnum[ddv1.PositionAuthenticationStatus](data[idxAuthStatus]); err == nil {
		record.SetAuthenticationStatus(authStatus)
	} else {
		record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNRECOGNIZED)
		record.SetUnrecognizedAuthenticationStatus(int32(data[idxAuthStatus]))
	}

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
//
// Binary Layout (12 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (6 bytes)
//	Offset 11: authenticationStatus (1 byte)
func AppendGNSSPlaceAuthRecord(dst []byte, record *ddv1.GNSSPlaceAuthRecord) ([]byte, error) {
	if record == nil {
		// Append 12 zero bytes if no data
		return append(dst, make([]byte, 12)...), nil
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

	// Append geoCoordinates (6 bytes)
	var err error
	dst, err = AppendGeoCoordinates(dst, record.GetGeoCoordinates())
	if err != nil {
		return nil, fmt.Errorf("failed to append geo coordinates: %w", err)
	}

	// Append authenticationStatus (1 byte)
	var authStatus byte
	if record.GetAuthenticationStatus() == ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNRECOGNIZED {
		authStatus = byte(record.GetUnrecognizedAuthenticationStatus())
	} else {
		var err error
		authStatus, err = MarshalEnum(record.GetAuthenticationStatus())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal authentication status: %w", err)
		}
	}
	dst = append(dst, authStatus)

	return dst, nil
}
