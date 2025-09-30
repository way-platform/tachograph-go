package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

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
func (opts UnmarshalOptions) unmarshalGnssPlaces(data []byte) (*cardv1.GnssPlaces, error) {
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

	// Write GNSS place records
	records := gnssPlaces.GetRecords()
	if len(records) > 0 {
		// Write number of records (1 byte)
		if len(records) > 255 {
			return nil, fmt.Errorf("too many GNSS place records: %d", len(records))
		}
		data = append(data, byte(len(records)))

		// Write each record
		for _, record := range records {
			var err error
			data, err = appendCardGnssPlaceRecord(data, record)
			if err != nil {
				return nil, fmt.Errorf("failed to append GNSS place record: %w", err)
			}
		}
	} else {
		// Write 0 records
		data = append(data, 0x00)
	}

	return data, nil
}

// appendCardGnssPlaceRecord appends a single GNSS place record to dst
func appendCardGnssPlaceRecord(dst []byte, record *cardv1.GnssPlaces_Record) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// Entry time (TimeReal - 4 bytes)
	var err error
	dst, err = dd.AppendTimeReal(dst, record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to append timestamp: %w", err)
	}

	// GNSS place record (GNSSPlaceRecord from DD)
	gnssPlaceRecord := record.GetGnssPlaceRecord()
	if gnssPlaceRecord != nil {
		// Append the GNSS place record using DD function
		var err error
		dst, err = dd.AppendGNSSPlaceRecord(dst, gnssPlaceRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to append GNSS place record: %w", err)
		}
	} else {
		// Append default values (12 bytes: 4 timestamp + 1 accuracy + 8 geo coords - 1 = 13 bytes total)
		// Actually the full GNSSPlaceRecord is 13 bytes: 4 (timestamp) + 1 (accuracy) + 8 (coords)
		dst = append(dst, make([]byte, 13)...)
	}

	// Vehicle odometer value (OdometerShort - 3 bytes)
	odometer := record.GetVehicleOdometerKm()
	if odometer < 0 || odometer > 999999 {
		return nil, fmt.Errorf("invalid odometer value: %d", odometer)
	}
	dst = dd.AppendOdometer(dst, uint32(odometer))

	return dst, nil
}
