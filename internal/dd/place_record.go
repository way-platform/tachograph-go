package dd

import (
	"fmt"

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
func (opts UnmarshalOptions) UnmarshalPlaceRecord(data []byte) (*ddv1.PlaceRecord, error) {
	const (
		lenPlaceRecordGen1 = 10
		lenPlaceRecordGen2 = 21 // 10 (base) + 11 (GNSS: 4 timestamp + 1 accuracy + 6 coords)
		idxEntryTime       = 0
		idxEntryType       = 4
		idxCountry         = 5
		idxRegion          = 6
		idxOdometer        = 7
		idxGNSS            = 10
	)

	// Check for Gen2; otherwise assume Gen1 (including zero value)
	expectedLen := lenPlaceRecordGen1
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		expectedLen = lenPlaceRecordGen2
	}

	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid data length for PlaceRecord: got %d, want %d (Gen1) or %d (Gen2)", len(data), lenPlaceRecordGen1, lenPlaceRecordGen2)
	}

	record := &ddv1.PlaceRecord{}
	// Populate generation from unmarshal context
	record.SetGeneration(opts.Generation)
	// Store raw data for round-trip fidelity
	record.SetRawData(data)
	// Mark as valid (parsing succeeded)
	record.SetValid(true)

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, err := opts.UnmarshalTimeReal(data[idxEntryTime : idxEntryTime+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// Parse entryTypeDailyWorkPeriod (1 byte)
	if entryType, err := UnmarshalEnum[ddv1.EntryTypeDailyWorkPeriod](data[idxEntryType]); err == nil {
		record.SetEntryTypeDailyWorkPeriod(entryType)
	} else {
		record.SetEntryTypeDailyWorkPeriod(ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNRECOGNIZED)
		record.SetUnrecognizedEntryTypeDailyWorkPeriod(int32(data[idxEntryType]))
	}

	// Parse dailyWorkPeriodCountry (NationNumeric - 1 byte)
	if country, err := UnmarshalEnum[ddv1.NationNumeric](data[idxCountry]); err == nil {
		record.SetDailyWorkPeriodCountry(country)
	} else {
		record.SetDailyWorkPeriodCountry(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
		record.SetUnrecognizedDailyWorkPeriodCountry(int32(data[idxCountry]))
	}

	// Parse dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	record.SetDailyWorkPeriodRegion([]byte{data[idxRegion]})

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, err := opts.UnmarshalOdometer(data[idxOdometer : idxOdometer+3])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal odometer: %w", err)
	}
	record.SetVehicleOdometerKm(int32(odometerValue))

	// For Gen2, parse the GNSS place record (11 bytes)
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		gnssData := data[idxGNSS : idxGNSS+11]
		gnssRecord, err := opts.UnmarshalGNSSPlaceRecord(gnssData)
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
func AppendPlaceRecord(dst []byte, record *ddv1.PlaceRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("place record cannot be nil")
	}

	// Get generation from the self-describing record
	generation := record.GetGeneration()
	if generation == ddv1.Generation_GENERATION_UNSPECIFIED {
		return nil, fmt.Errorf("PlaceRecord.generation must be specified (got GENERATION_UNSPECIFIED)")
	}

	const (
		lenPlaceRecordGen1 = 10
		lenPlaceRecordGen2 = 21
	)

	recordSize := lenPlaceRecordGen1
	if generation == ddv1.Generation_GENERATION_2 {
		recordSize = lenPlaceRecordGen2
	}

	// Raw data painting: Use raw_data as canvas if available, otherwise zero-filled buffer
	var canvas []byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != recordSize {
			return nil, fmt.Errorf("invalid raw_data length for PlaceRecord: got %d, want %d", len(rawData), recordSize)
		}
		canvas = make([]byte, recordSize)
		copy(canvas, rawData)
	} else {
		canvas = make([]byte, recordSize)
	}

	// Paint semantic values over the canvas
	const (
		idxEntryTime = 0
		idxEntryType = 4
		idxCountry   = 5
		idxRegion    = 6
		idxOdometer  = 7
		idxGNSS      = 10
	)

	// Paint entryTime (4 bytes)
	var err error
	timeBytes, err := AppendTimeReal(nil, record.GetEntryTime())
	if err != nil {
		return nil, fmt.Errorf("failed to encode entry time: %w", err)
	}
	copy(canvas[idxEntryTime:idxEntryTime+4], timeBytes)

	// Paint entryTypeDailyWorkPeriod (1 byte)
	var entryTypeValue byte
	if record.GetEntryTypeDailyWorkPeriod() == ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNRECOGNIZED {
		entryTypeValue = byte(record.GetUnrecognizedEntryTypeDailyWorkPeriod())
	} else {
		var err error
		entryTypeValue, err = MarshalEnum(record.GetEntryTypeDailyWorkPeriod())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal entry type: %w", err)
		}
	}
	canvas[idxEntryType] = entryTypeValue

	// Paint dailyWorkPeriodCountry (1 byte)
	var countryValue byte
	if record.GetDailyWorkPeriodCountry() == ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED {
		countryValue = byte(record.GetUnrecognizedDailyWorkPeriodCountry())
	} else {
		var err error
		countryValue, err = MarshalEnum(record.GetDailyWorkPeriodCountry())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal country: %w", err)
		}
	}
	canvas[idxCountry] = countryValue

	// Paint dailyWorkPeriodRegion (1 byte)
	region := record.GetDailyWorkPeriodRegion()
	if len(region) > 0 {
		canvas[idxRegion] = region[0]
	} else {
		canvas[idxRegion] = 0
	}

	// Paint vehicleOdometerValue (3 bytes)
	odometerBytes := AppendOdometer(nil, uint32(record.GetVehicleOdometerKm()))
	copy(canvas[idxOdometer:idxOdometer+3], odometerBytes)

	// For Gen2, paint the GNSS place record (11 bytes)
	if generation == ddv1.Generation_GENERATION_2 {
		gnssBytes, err := AppendGNSSPlaceRecord(nil, record.GetEntryGnssPlaceRecord())
		if err != nil {
			return nil, fmt.Errorf("failed to encode GNSS place record: %w", err)
		}
		copy(canvas[idxGNSS:idxGNSS+11], gnssBytes)
	}

	return append(dst, canvas...), nil
}
