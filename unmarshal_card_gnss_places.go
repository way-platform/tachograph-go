package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardGnssPlaces unmarshals GNSS places data from a card EF.
func unmarshalCardGnssPlaces(data []byte) (*cardv1.GnssPlaces, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GNSS places")
	}

	var target cardv1.GnssPlaces
	r := bytes.NewReader(data)

	// Read newest record index (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}
	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// For now, just set empty records to satisfy the interface
	// The actual GNSS places structure is complex and would need detailed parsing
	target.SetRecords([]*cardv1.GnssPlaces_Record{})

	return &target, nil
}

// UnmarshalCardGnssPlaces unmarshals GNSS places data from a card EF (legacy function).
// Deprecated: Use unmarshalCardGnssPlaces instead.
func UnmarshalCardGnssPlaces(data []byte, target *cardv1.GnssPlaces) error {
	result, err := unmarshalCardGnssPlaces(data)
	if err != nil {
		return err
	}
	*target = *result
	return nil
}
