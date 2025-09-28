package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardGnssPlaces unmarshals GNSS places data from a card EF.
//
// The data type `CardGNSSPlaceRecord` is specified in the Data Dictionary, Section 2.5.
//
// ASN.1 Definition:
//
//	CardGNSSPlaceRecord ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    gnssPlaceAccuracy           GNSSPlaceAccuracy,
//	    longitude                    Longitude,
//	    latitude                     Latitude,
//	    vehicleOdometerValue         OdometerShort
//	}
func unmarshalCardGnssPlaces(data []byte) (*cardv1.GnssPlaces, error) {
	const (
		lenMinEfGnssPlaces = 2 // Minimum EF_GNSSPlaces record size
	)

	if len(data) < lenMinEfGnssPlaces {
		return nil, fmt.Errorf("insufficient data for GNSS places: got %d bytes, need at least %d", len(data), lenMinEfGnssPlaces)
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

// AppendCardGnssPlaces appends GNSS places data to a byte slice.
//
// The data type `CardGNSSPlaceRecord` is specified in the Data Dictionary, Section 2.5.
//
// ASN.1 Definition:
//
//	CardGNSSPlaceRecord ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    gnssPlaceAccuracy           GNSSPlaceAccuracy,
//	    longitude                    Longitude,
//	    latitude                     Latitude,
//	    vehicleOdometerValue         OdometerShort
//	}
func appendCardGnssPlaces(data []byte, gnssPlaces *cardv1.GnssPlaces) ([]byte, error) {
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
