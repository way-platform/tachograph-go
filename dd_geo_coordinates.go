package tachograph

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalGeoCoordinates parses geo coordinates data.
//
// The data type `GeoCoordinates` is specified in the Data Dictionary, Section 2.76.
//
// ASN.1 Definition:
//
//	GeoCoordinates ::= SEQUENCE {
//	    latitude INTEGER(-90*3600*1000..90*3600*1000),
//	    longitude INTEGER(-180*3600*1000+1..180*3600*1000)
//	}
//
// Binary Layout (8 bytes total):
//   - Latitude (4 bytes): Signed integer in millionths of a degree
//   - Longitude (4 bytes): Signed integer in millionths of a degree
func unmarshalGeoCoordinates(data []byte) (*ddv1.GeoCoordinates, error) {
	const (
		lenGeoCoordinates = 8 // 4 bytes latitude + 4 bytes longitude
	)

	if len(data) < lenGeoCoordinates {
		return nil, fmt.Errorf("insufficient data for GeoCoordinates: got %d, want %d", len(data), lenGeoCoordinates)
	}

	geoCoords := &ddv1.GeoCoordinates{}

	// Parse latitude (4 bytes, signed big-endian)
	latitude := int32(binary.BigEndian.Uint32(data[0:4]))
	geoCoords.SetLatitude(latitude)

	// Parse longitude (4 bytes, signed big-endian)
	longitude := int32(binary.BigEndian.Uint32(data[4:8]))
	geoCoords.SetLongitude(longitude)

	return geoCoords, nil
}

// appendGeoCoordinates appends geo coordinates data to dst.
//
// The data type `GeoCoordinates` is specified in the Data Dictionary, Section 2.76.
//
// ASN.1 Definition:
//
//	GeoCoordinates ::= SEQUENCE {
//	    latitude INTEGER(-90*3600*1000..90*3600*1000),
//	    longitude INTEGER(-180*3600*1000+1..180*3600*1000)
//	}
//
// Binary Layout (8 bytes total):
//   - Latitude (4 bytes): Signed integer in millionths of a degree
//   - Longitude (4 bytes): Signed integer in millionths of a degree
func appendGeoCoordinates(dst []byte, geoCoords *ddv1.GeoCoordinates) ([]byte, error) {
	if geoCoords == nil {
		// Append default values (8 zero bytes)
		return append(dst, make([]byte, 8)...), nil
	}

	// Append latitude (4 bytes, signed big-endian)
	latitude := geoCoords.GetLatitude()
	dst = binary.BigEndian.AppendUint32(dst, uint32(latitude))

	// Append longitude (4 bytes, signed big-endian)
	longitude := geoCoords.GetLongitude()
	dst = binary.BigEndian.AppendUint32(dst, uint32(longitude))

	return dst, nil
}
