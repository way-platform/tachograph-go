package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendPlaces appends the binary representation of Places to dst.
//
// ASN.1 Specification (Data Dictionary 2.4):
//
//	CardPlaceDailyWorkPeriod ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort
//	}
//
// Binary Layout (variable size):
//
//	0-1:   newestRecordIndex (2 bytes, big-endian)
//	2+:    place records (12 bytes each)
//	  - 0-3:   entryTime (4 bytes)
//	  - 4-4:   entryTypeDailyWorkPeriod (1 byte)
//	  - 5-5:   dailyWorkPeriodCountry (1 byte)
//	  - 6-7:   dailyWorkPeriodRegion (2 bytes, big-endian)
//	  - 8-10:  vehicleOdometerValue (3 bytes)
//	  - 11-11: reserved (1 byte)
func AppendPlaces(dst []byte, p *cardv1.Places) ([]byte, error) {
	if p == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(p.GetNewestRecordIndex()))

	var err error
	for _, rec := range p.GetRecords() {
		dst, err = AppendPlaceRecord(dst, rec)
		if err != nil {
			return nil, err
		}
	}

	// Append trailing bytes for roundtrip accuracy
	if trailingBytes := p.GetTrailingBytes(); len(trailingBytes) > 0 {
		dst = append(dst, trailingBytes...)
	}

	return dst, nil
}

// AppendPlaceRecord appends a single 12-byte place record.
func AppendPlaceRecord(dst []byte, rec *cardv1.Places_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, placeRecordSize)...), nil
	}
	dst = appendTimeReal(dst, rec.GetEntryTime()) // 4 bytes

	// Entry type with protocol value conversion using generic helper
	entryTypeProtocol := GetProtocolValueFromEnum(rec.GetEntryType(), 0)
	dst = append(dst, byte(entryTypeProtocol)) // 1 byte

	// Country with protocol value conversion using generic helper
	countryProtocol := GetProtocolValueFromEnum(rec.GetDailyWorkPeriodCountry(), 0)
	dst = append(dst, byte(countryProtocol)) // 1 byte

	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetDailyWorkPeriodRegion())) // 2 bytes
	dst = appendOdometer(dst, uint32(rec.GetVehicleOdometerKm()))                    // 3 bytes
	dst = append(dst, byte(rec.GetReservedByte()))                                   // 1 byte reserved (preserved)
	return dst, nil
}
