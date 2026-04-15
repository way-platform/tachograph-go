package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGeoCoordinates parses geo coordinates data.
//
// The data type `GeoCoordinates` is specified in the Data Dictionary, Section 2.76.
//
// ASN.1 Definition:
//
//	GeoCoordinates ::= SEQUENCE {
//	    latitude INTEGER(-90000..90001),
//	    longitude INTEGER(-180000..180001)
//	}
//
// Binary Layout (6 bytes total):
//   - Latitude (3 bytes): Signed 24-bit integer in ±DDMM.M × 10 format
//   - Longitude (3 bytes): Signed 24-bit integer in ±DDDMM.M × 10 format
//
// Returns nil when position is unavailable. The canonical sentinel is both
// latitude and longitude equal to 0x7FFFFF (8388607). If only one field
// carries the sentinel the data is malformed; the position is treated as
// unavailable in that case too.
func (opts UnmarshalOptions) UnmarshalGeoCoordinates(data []byte) (*ddv1.GeoCoordinates, error) {
	const (
		lenGeoCoordinates  = 6 // 3 bytes latitude + 3 bytes longitude
		geoSentinel       int32 = 0x7FFFFF
	)
	if len(data) != lenGeoCoordinates {
		return nil, fmt.Errorf("invalid data length for GeoCoordinates: got %d, want %d", len(data), lenGeoCoordinates)
	}
	lat := readInt24(data[0:3])
	lon := readInt24(data[3:6])
	if lat == geoSentinel || lon == geoSentinel {
		return nil, nil
	}
	var output ddv1.GeoCoordinates
	output.SetLatitude(lat)
	output.SetLongitude(lon)
	return &output, nil
}

// MarshalGeoCoordinates marshals geo coordinates data to bytes.
//
// The data type `GeoCoordinates` is specified in the Data Dictionary, Section 2.76.
//
// ASN.1 Definition:
//
//	GeoCoordinates ::= SEQUENCE {
//	    latitude INTEGER(-90000..90001),
//	    longitude INTEGER(-180000..180001)
//	}
//
// Binary Layout (6 bytes total):
//   - Latitude (3 bytes): Signed 24-bit integer in ±DDMM.M × 10 format
//   - Longitude (3 bytes): Signed 24-bit integer in ±DDDMM.M × 10 format
//
// Unknown position marker: 0x7FFFFF (8388607 decimal)
func (opts MarshalOptions) MarshalGeoCoordinates(geoCoords *ddv1.GeoCoordinates) ([]byte, error) {
	const lenGeoCoordinates = 6
	var canvas [lenGeoCoordinates]byte
	// Marshal latitude (3-byte signed integer)
	latBytes := marshalInt24(geoCoords.GetLatitude())
	copy(canvas[0:3], latBytes)
	// Marshal longitude (3-byte signed integer)
	longBytes := marshalInt24(geoCoords.GetLongitude())
	copy(canvas[3:6], longBytes)
	return canvas[:], nil
}
