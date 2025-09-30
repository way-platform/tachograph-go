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
// Unknown position marker: 0x7FFFFF (8388607 decimal)
func (opts UnmarshalOptions) UnmarshalGeoCoordinates(data []byte) (*ddv1.GeoCoordinates, error) {
	const (
		lenGeoCoordinates = 6 // 3 bytes latitude + 3 bytes longitude
	)
	if len(data) != lenGeoCoordinates {
		return nil, fmt.Errorf("invalid data length for GeoCoordinates: got %d, want %d", len(data), lenGeoCoordinates)
	}
	var output ddv1.GeoCoordinates
	output.SetLatitude(readInt24(data[0:3]))
	output.SetLongitude(readInt24(data[3:6]))
	return &output, nil
}

// readInt24 reads a 3-byte big-endian signed integer.
// The value is sign-extended to 32 bits.
func readInt24(data []byte) int32 {
	// Read as unsigned 24-bit value
	val := uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	// Sign extend from 24 bits to 32 bits
	// If bit 23 is set (negative number), extend with 1s
	if val&0x800000 != 0 {
		val |= 0xFF000000
	}
	return int32(val)
}

// AppendGeoCoordinates appends geo coordinates data to dst.
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
func AppendGeoCoordinates(dst []byte, geoCoords *ddv1.GeoCoordinates) ([]byte, error) {
	// Append latitude (3-byte signed integer)
	dst = appendInt24(dst, geoCoords.GetLatitude())
	// Append longitude (3-byte signed integer)
	dst = appendInt24(dst, geoCoords.GetLongitude())
	return dst, nil
}

// appendInt24 appends a 3-byte big-endian signed integer to dst.
// Only the lower 24 bits of the value are written.
func appendInt24(dst []byte, val int32) []byte {
	// Write the lower 24 bits in big-endian order
	dst = append(dst, byte(val>>16), byte(val>>8), byte(val))
	return dst
}
