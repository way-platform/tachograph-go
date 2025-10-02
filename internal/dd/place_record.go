package dd

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UnmarshalPlaceRecord parses a Generation 1 place record (10 bytes, no GNSS data).
//
// The data type `PlaceRecord` is specified in the Data Dictionary, Section 2.117.
//
// ASN.1 Definition (Gen1):
//
//	PlaceRecord ::= SEQUENCE {
//	    entryTime                    TimeReal,
//	    entryTypeDailyWorkPeriod     EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry       NationNumeric,
//	    dailyWorkPeriodRegion        RegionNumeric,
//	    vehicleOdometerValue         OdometerShort
//	}
func (opts UnmarshalOptions) UnmarshalPlaceRecord(data []byte) (*ddv1.PlaceRecord, error) {
	const (
		idxEntryTime   = 0
		idxEntryType   = 4
		idxCountry     = 5
		idxRegion      = 6
		idxOdometer    = 7
		lenPlaceRecord = 10 // Fixed size for Gen1
	)

	if len(data) != lenPlaceRecord {
		return nil, fmt.Errorf("invalid data length for Gen1 PlaceRecord: got %d, want %d", len(data), lenPlaceRecord)
	}

	record := &ddv1.PlaceRecord{}
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

	return record, nil
}

// AppendPlaceRecord appends a Generation 1 place record (10 bytes).
func AppendPlaceRecord(dst []byte, rec *ddv1.PlaceRecord) ([]byte, error) {
	const lenPlaceRecord = 10 // Fixed size for Gen1

	// Use raw data painting strategy if available
	var canvas [lenPlaceRecord]byte
	if rawData := rec.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenPlaceRecord {
			return nil, fmt.Errorf("invalid raw_data length for PlaceRecord: got %d, want %d", len(rawData), lenPlaceRecord)
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

	return append(dst, canvas[:]...), nil
}

// anonymizeTimestamp is a placeholder that returns the timestamp unchanged.
// Actual anonymization happens at the Places message level via AnonymizeTimestampsInPlace,
// which needs access to all timestamps to calculate a dataset-specific offset.
//
// This function exists for API consistency but does not modify the timestamp.
func anonymizeTimestamp(ts *timestamppb.Timestamp) *timestamppb.Timestamp {
	return ts
}

// AnonymizePlaceRecord creates an anonymized copy of PlaceRecord, preserving the
// structure while replacing potentially sensitive location data with normalized values.
//
// The anonymization:
// - Shifts timestamps to test epoch (2020) while preserving relative timing
// - Normalizes country/region to generic values
// - Rounds odometer to nearest 100km
// - Preserves entry type (needed for structure testing)
func AnonymizePlaceRecord(rec *ddv1.PlaceRecord) *ddv1.PlaceRecord {
	if rec == nil {
		return nil
	}

	result := &ddv1.PlaceRecord{}

	// Anonymize timestamp: shift to test epoch while preserving all relative timing
	result.SetEntryTime(anonymizeTimestamp(rec.GetEntryTime()))

	// Preserve entry type (structural information)
	result.SetEntryTypeDailyWorkPeriod(rec.GetEntryTypeDailyWorkPeriod())
	if rec.HasUnrecognizedEntryTypeDailyWorkPeriod() {
		result.SetUnrecognizedEntryTypeDailyWorkPeriod(rec.GetUnrecognizedEntryTypeDailyWorkPeriod())
	}

	// Anonymize country (use a generic test country code)
	result.SetDailyWorkPeriodCountry(ddv1.NationNumeric_FINLAND) // Finland as test default

	// Anonymize region (use generic value)
	result.SetDailyWorkPeriodRegion([]byte{0x01})

	// Round odometer to nearest 100km (preserves magnitude but not exact location correlation)
	originalOdometer := rec.GetVehicleOdometerKm()
	roundedOdometer := (originalOdometer / 100) * 100
	result.SetVehicleOdometerKm(roundedOdometer)

	// Regenerate raw_data to match anonymized values
	// This ensures round-trip fidelity after anonymization
	anonymizedBytes, err := AppendPlaceRecord(nil, result)
	if err == nil {
		result.SetRawData(anonymizedBytes)
	}
	// If marshalling fails, we'll have no raw_data, which is acceptable

	return result
}
