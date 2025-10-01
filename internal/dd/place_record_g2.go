package dd

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPlaceRecordG2 parses a Generation 2 place record (21 bytes, includes GNSS data).
//
// The data type `PlaceRecord` (Gen2 variant) is specified in the Data Dictionary, Section 2.117.
//
// ASN.1 Definition (Gen2):
//
//	PlaceRecord ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort,
//	    entryGNSSPlaceRecord         GNSSPlaceRecord
//	}
func (opts UnmarshalOptions) UnmarshalPlaceRecordG2(data []byte) (*ddv1.PlaceRecordG2, error) {
	const (
		idxEntryTime   = 0
		idxEntryType   = 4
		idxCountry     = 5
		idxRegion      = 6
		idxOdometer    = 7
		idxGNSS        = 10
		lenPlaceRecord = 21 // Fixed size for Gen2
	)

	if len(data) != lenPlaceRecord {
		return nil, fmt.Errorf("invalid data length for Gen2 PlaceRecord: got %d, want %d", len(data), lenPlaceRecord)
	}

	record := &ddv1.PlaceRecordG2{}
	record.SetRawData(data)

	// Parse entry time (4 bytes)
	entryTime, err := opts.UnmarshalTimeReal(data[idxEntryTime : idxEntryTime+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// Parse entry type (1 byte)
	entryTypeByte := data[idxEntryType]
	entryType, err := UnmarshalEnum[ddv1.EntryTypeDailyWorkPeriod](entryTypeByte)
	if err != nil {
		record.SetEntryTypeDailyWorkPeriod(ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNRECOGNIZED)
		record.SetUnrecognizedEntryTypeDailyWorkPeriod(int32(entryTypeByte))
	} else {
		record.SetEntryTypeDailyWorkPeriod(entryType)
	}

	// Parse country (1 byte)
	countryByte := data[idxCountry]
	country, err := UnmarshalEnum[ddv1.NationNumeric](countryByte)
	if err != nil {
		record.SetDailyWorkPeriodCountry(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
		record.SetUnrecognizedDailyWorkPeriodCountry(int32(countryByte))
	} else {
		record.SetDailyWorkPeriodCountry(country)
	}

	// Parse region (1 byte)
	record.SetDailyWorkPeriodRegion([]byte{data[idxRegion]})

	// Parse odometer (3 bytes)
	odometerBytes := data[idxOdometer : idxOdometer+3]
	record.SetVehicleOdometerKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBytes...))))

	// Parse GNSS place record (11 bytes)
	gnssRecord, err := opts.UnmarshalGNSSPlaceRecord(data[idxGNSS : idxGNSS+11])
	if err != nil {
		return nil, fmt.Errorf("failed to parse GNSS place record: %w", err)
	}
	record.SetEntryGnssPlaceRecord(gnssRecord)

	return record, nil
}

// AppendPlaceRecordG2 appends a Generation 2 place record (21 bytes).
func AppendPlaceRecordG2(dst []byte, rec *ddv1.PlaceRecordG2) ([]byte, error) {
	const lenPlaceRecord = 21 // Fixed size for Gen2

	// Use raw data painting strategy if available
	var canvas [lenPlaceRecord]byte
	if rawData := rec.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenPlaceRecord {
			return nil, fmt.Errorf("invalid raw_data length for PlaceRecordG2: got %d, want %d", len(rawData), lenPlaceRecord)
		}
		copy(canvas[:], rawData)
	}

	// Paint semantic values over the canvas
	var err error
	timeBytes, err := AppendTimeReal(nil, rec.GetEntryTime())
	if err != nil {
		return nil, fmt.Errorf("failed to append entry time: %w", err)
	}
	copy(canvas[0:4], timeBytes)

	// Entry type (1 byte)
	var entryTypeProtocol byte
	if rec.GetEntryTypeDailyWorkPeriod() == ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNRECOGNIZED {
		entryTypeProtocol = byte(rec.GetUnrecognizedEntryTypeDailyWorkPeriod())
	} else {
		entryTypeProtocol, _ = MarshalEnum(rec.GetEntryTypeDailyWorkPeriod())
	}
	canvas[4] = entryTypeProtocol

	// Country (1 byte)
	var countryProtocol byte
	if rec.GetDailyWorkPeriodCountry() == ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED {
		countryProtocol = byte(rec.GetUnrecognizedDailyWorkPeriodCountry())
	} else {
		countryProtocol, _ = MarshalEnum(rec.GetDailyWorkPeriodCountry())
	}
	canvas[5] = countryProtocol

	// Region (1 byte)
	regionBytes := rec.GetDailyWorkPeriodRegion()
	if len(regionBytes) > 0 {
		canvas[6] = regionBytes[0]
	}
	// Otherwise leave as zero (or preserved from raw_data)

	// Odometer (3 bytes)
	odometerBytes := AppendOdometer(nil, uint32(rec.GetVehicleOdometerKm()))
	copy(canvas[7:10], odometerBytes)

	// GNSS place record (11 bytes)
	gnssBytes, err := AppendGNSSPlaceRecord(nil, rec.GetEntryGnssPlaceRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to append GNSS place record: %w", err)
	}
	copy(canvas[10:21], gnssBytes)

	return append(dst, canvas[:]...), nil
}
