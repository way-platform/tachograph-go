package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalCardGnssPlaces unmarshals GNSS places data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.5):
//
//	CardGNSSPlaceRecord ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    gnssPlaceAccuracy           GNSSPlaceAccuracy,
//	    longitude                    Longitude,
//	    latitude                     Latitude,
//	    vehicleOdometerValue         OdometerShort
//	}
//
// Binary Layout (variable size):
//
//	0-1:   newestRecordIndex (2 bytes, big-endian)
//	2+:    GNSS place records (variable size each)
//	  - 0-3:   entryTime (4 bytes)
//	  - 4-4:   gnssPlaceAccuracy (1 byte)
//	  - 5-8:   longitude (4 bytes, big-endian)
//	  - 9-12:  latitude (4 bytes, big-endian)
//	  - 13-15: vehicleOdometerValue (3 bytes)
func unmarshalCardGnssPlaces(data []byte) (*cardv1.GnssPlaces, error) {
	const (
		// Minimum EF_GNSSPlaces record size
		MIN_EF_GNSS_PLACES_SIZE = 2
	)

	if len(data) < MIN_EF_GNSS_PLACES_SIZE {
		return nil, fmt.Errorf("insufficient data for GNSS places: got %d bytes, need at least %d", len(data), MIN_EF_GNSS_PLACES_SIZE)
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
	proto.Merge(target, result)
	return nil
}
