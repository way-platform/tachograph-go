package dd

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
//
//nolint:unused
func UnmarshalGeoCoordinates(data []byte) (*ddv1.GeoCoordinates, error) {
	const (
		lenGeoCoordinates = 8 // 4 bytes latitude + 4 bytes longitude
	)
	if len(data) != lenGeoCoordinates {
		return nil, fmt.Errorf("invalid data length for GeoCoordinates: got %d, want %d", len(data), lenGeoCoordinates)
	}
	var output ddv1.GeoCoordinates
	output.SetLatitude(int32(binary.BigEndian.Uint32(data[0:4])))
	output.SetLongitude(int32(binary.BigEndian.Uint32(data[4:8])))
	return &output, nil
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
func AppendGeoCoordinates(dst []byte, geoCoords *ddv1.GeoCoordinates) ([]byte, error) {
	dst = binary.BigEndian.AppendUint32(dst, uint32(geoCoords.GetLatitude()))
	dst = binary.BigEndian.AppendUint32(dst, uint32(geoCoords.GetLongitude()))
	return dst, nil
}
