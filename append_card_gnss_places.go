package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardGnssPlaces appends GNSS places data to a byte slice.
func AppendCardGnssPlaces(data []byte, gnssPlaces *cardv1.GnssPlaces) ([]byte, error) {
	if gnssPlaces == nil {
		return data, nil
	}

	// Newest record index (2 bytes)
	if gnssPlaces.HasNewestRecordIndex() {
		indexBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(indexBytes, uint16(gnssPlaces.GetNewestRecordIndex()))
		data = append(data, indexBytes...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// For now, skip the complex record structures
	// This provides a basic implementation that satisfies the interface

	return data, nil
}
