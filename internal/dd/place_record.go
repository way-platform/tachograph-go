package dd

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPlaceRecord unmarshals a PlaceRecord from binary data.
//
// The data type `PlaceRecord` is specified in the Data Dictionary, Section 2.117.
//
// ASN.1 Definition (Gen1):
//
//	PlaceRecord ::= SEQUENCE {
//	    entryTime TimeReal,
//	    entryTypeDailyWorkPeriod EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry NationNumeric,
//	    dailyWorkPeriodRegion RegionNumeric,
//	    vehicleOdometerValue OdometerShort
//	}
//
// For Gen2, the following component is added:
//
//	entryGNSSPlaceRecord GNSSPlaceRecord
func UnmarshalPlaceRecord(data []byte, generation ddv1.Generation) (*ddv1.PlaceRecord, error) {
	const (
		lenPlaceRecordGen1 = 10
		lenPlaceRecordGen2 = 22 // 10 (base) + 12 (GNSS: 8 bytes coords + 4 bytes timestamp)
		idxEntryTime       = 0
		idxEntryType       = 4
		idxCountry         = 5
		idxRegion          = 6
		idxOdometer        = 7
		idxGNSS            = 10
	)

	expectedLen := lenPlaceRecordGen1
	if generation == ddv1.Generation_GENERATION_2 {
		expectedLen = lenPlaceRecordGen2
	}

	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid data length for PlaceRecord (gen %d): got %d, want %d", generation, len(data), expectedLen)
	}

	record := &ddv1.PlaceRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, err := UnmarshalTimeReal(data[idxEntryTime : idxEntryTime+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// Parse entryTypeDailyWorkPeriod (1 byte)
	SetEnumFromProtocolValueGeneric(
		ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED.Descriptor(),
		int32(data[idxEntryType]),
		func(enumNum protoreflect.EnumNumber) {
			record.SetEntryTypeDailyWorkPeriod(ddv1.EntryTypeDailyWorkPeriod(enumNum))
		},
		func(rawValue int32) {
			record.SetUnrecognizedEntryTypeDailyWorkPeriod(rawValue)
		},
	)

	// Parse dailyWorkPeriodCountry (NationNumeric - 1 byte)
	SetEnumFromProtocolValueGeneric(
		ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(),
		int32(data[idxCountry]),
		func(enumNum protoreflect.EnumNumber) {
			record.SetDailyWorkPeriodCountry(ddv1.NationNumeric(enumNum))
		},
		func(rawValue int32) {
			record.SetUnrecognizedDailyWorkPeriodCountry(rawValue)
		},
	)

	// Parse dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	record.SetDailyWorkPeriodRegion([]byte{data[idxRegion]})

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, err := UnmarshalOdometer(data[idxOdometer : idxOdometer+3])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal odometer: %w", err)
	}
	record.SetVehicleOdometerKm(int32(odometerValue))

	// For Gen2, parse the GNSS place record (12 bytes)
	if generation == ddv1.Generation_GENERATION_2 {
		gnssData := data[idxGNSS : idxGNSS+12]
		gnssRecord, err := UnmarshalGNSSPlaceRecordCardVariant(gnssData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal GNSS place record: %w", err)
		}
		record.SetEntryGnssPlaceRecord(gnssRecord)
	}

	return record, nil
}

// AppendPlaceRecord appends a PlaceRecord to dst.
//
// The data type `PlaceRecord` is specified in the Data Dictionary, Section 2.117.
//
// ASN.1 Definition (Gen1):
//
//	PlaceRecord ::= SEQUENCE {
//	    entryTime TimeReal,
//	    entryTypeDailyWorkPeriod EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry NationNumeric,
//	    dailyWorkPeriodRegion RegionNumeric,
//	    vehicleOdometerValue OdometerShort
//	}
//
// For Gen2, the following component is added:
//
//	entryGNSSPlaceRecord GNSSPlaceRecord
func AppendPlaceRecord(dst []byte, record *ddv1.PlaceRecord, generation ddv1.Generation) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("place record cannot be nil")
	}

	// Append entryTime (TimeReal - 4 bytes)
	var err error
	entryTime := record.GetEntryTime()
	if entryTime == nil {
		// Append zero timestamp if nil
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	} else {
		dst, err = AppendTimeReal(dst, entryTime)
		if err != nil {
			return nil, fmt.Errorf("failed to append entry time: %w", err)
		}
	}

	// Append entryTypeDailyWorkPeriod (1 byte)
	entryTypeValue, _ := GetProtocolValueForEnum(record.GetEntryTypeDailyWorkPeriod())
	dst = append(dst, byte(entryTypeValue))

	// Append dailyWorkPeriodCountry (NationNumeric - 1 byte)
	countryValue, _ := GetProtocolValueForEnum(record.GetDailyWorkPeriodCountry())
	dst = append(dst, byte(countryValue))

	// Append dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	region := record.GetDailyWorkPeriodRegion()
	if len(region) > 0 {
		dst = append(dst, region[0])
	} else {
		dst = append(dst, 0)
	}

	// Append vehicleOdometerValue (OdometerShort - 3 bytes)
	dst = AppendOdometer(dst, uint32(record.GetVehicleOdometerKm()))

	// For Gen2, append the GNSS place record (12 bytes)
	if generation == ddv1.Generation_GENERATION_2 {
		gnssRecord := record.GetEntryGnssPlaceRecord()
		dst, err = AppendGNSSPlaceRecordCardVariant(dst, gnssRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to append GNSS place record: %w", err)
		}
	}

	return dst, nil
}
