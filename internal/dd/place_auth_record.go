package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPlaceAuthRecord parses a PlaceAuthRecord (22 bytes).
//
// The data type `PlaceAuthRecord` is specified in the Data Dictionary, Section 2.116a.
//
// ASN.1 Definition (Gen2 V2):
//
//	PlaceAuthRecord ::= SEQUENCE {
//	    entryTime                       TimeReal,
//	    entryTypeDailyWorkPeriod        EntryTypeDailyWorkPeriod,
//	    dailyWorkPeriodCountry          NationNumeric,
//	    dailyWorkPeriodRegion           RegionNumeric,
//	    vehicleOdometerValue            OdometerShort,
//	    entryGNSSPlaceAuthRecord        GNSSPlaceAuthRecord
//	}
//
// Binary Layout (fixed length, 22 bytes):
//   - Bytes 0-3: entryTime (TimeReal)
//   - Byte 4: entryTypeDailyWorkPeriod (EntryTypeDailyWorkPeriod)
//   - Byte 5: dailyWorkPeriodCountry (NationNumeric)
//   - Byte 6: dailyWorkPeriodRegion (RegionNumeric)
//   - Bytes 7-9: vehicleOdometerValue (OdometerShort)
//   - Bytes 10-21: entryGNSSPlaceAuthRecord (GNSSPlaceAuthRecord)
func (opts UnmarshalOptions) UnmarshalPlaceAuthRecord(data []byte) (*ddv1.PlaceAuthRecord, error) {
	const (
		idxEntryTime                = 0
		idxEntryTypeDailyWorkPeriod = 4
		idxDailyWorkPeriodCountry   = 5
		idxDailyWorkPeriodRegion    = 6
		idxVehicleOdometerValue     = 7
		idxEntryGNSSPlaceAuthRecord = 10
		lenPlaceAuthRecord          = 22

		lenTimeReal                 = 4
		lenEntryTypeDailyWorkPeriod = 1
		lenNationNumeric            = 1
		lenRegionNumeric            = 1
		lenOdometerShort            = 3
		lenGNSSPlaceAuthRecord      = 12
	)

	if len(data) != lenPlaceAuthRecord {
		return nil, fmt.Errorf("invalid data length for PlaceAuthRecord: got %d, want %d", len(data), lenPlaceAuthRecord)
	}

	record := &ddv1.PlaceAuthRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// entryTime (4 bytes)
	entryTime, err := opts.UnmarshalTimeReal(data[idxEntryTime : idxEntryTime+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// entryTypeDailyWorkPeriod (1 byte)
	entryTypeDailyWorkPeriod, err := UnmarshalEnum[ddv1.EntryTypeDailyWorkPeriod](data[idxEntryTypeDailyWorkPeriod])
	if err != nil {
		return nil, fmt.Errorf("unmarshal entry type daily work period: %w", err)
	}
	record.SetEntryTypeDailyWorkPeriod(entryTypeDailyWorkPeriod)

	// dailyWorkPeriodCountry (1 byte)
	dailyWorkPeriodCountry, err := UnmarshalEnum[ddv1.NationNumeric](data[idxDailyWorkPeriodCountry])
	if err != nil {
		return nil, fmt.Errorf("unmarshal daily work period country: %w", err)
	}
	record.SetDailyWorkPeriodCountry(dailyWorkPeriodCountry)

	// dailyWorkPeriodRegion (1 byte)
	dailyWorkPeriodRegion, err := UnmarshalEnum[ddv1.RegionNumeric](data[idxDailyWorkPeriodRegion])
	if err != nil {
		return nil, fmt.Errorf("unmarshal daily work period region: %w", err)
	}
	record.SetDailyWorkPeriodRegion(dailyWorkPeriodRegion)

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerValue, err := opts.UnmarshalOdometer(data[idxVehicleOdometerValue : idxVehicleOdometerValue+lenOdometerShort])
	if err != nil {
		return nil, fmt.Errorf("unmarshal vehicle odometer value: %w", err)
	}
	if vehicleOdometerValue != nil {
		record.SetVehicleOdometerKm(*vehicleOdometerValue)
	}

	// entryGNSSPlaceAuthRecord (12 bytes)
	entryGNSSPlaceAuthRecord, err := opts.UnmarshalGNSSPlaceAuthRecord(data[idxEntryGNSSPlaceAuthRecord : idxEntryGNSSPlaceAuthRecord+lenGNSSPlaceAuthRecord])
	if err != nil {
		return nil, fmt.Errorf("unmarshal entry GNSS place auth record: %w", err)
	}
	record.SetEntryGnssPlaceAuthRecord(entryGNSSPlaceAuthRecord)

	return record, nil
}

// MarshalPlaceAuthRecord marshals a PlaceAuthRecord (22 bytes) to bytes.
func (opts MarshalOptions) MarshalPlaceAuthRecord(record *ddv1.PlaceAuthRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenPlaceAuthRecord = 22

	// Use raw data painting strategy if available
	var canvas [lenPlaceAuthRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenPlaceAuthRecord {
			return nil, fmt.Errorf("invalid raw_data length for PlaceAuthRecord: got %d, want %d", len(rawData), lenPlaceAuthRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// entryTime (4 bytes)
	entryTimeBytes, err := opts.MarshalTimeReal(record.GetEntryTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry time: %w", err)
	}
	copy(canvas[offset:offset+4], entryTimeBytes)
	offset += 4

	// entryTypeDailyWorkPeriod (1 byte)
	entryTypeDailyWorkPeriodByte, _ := MarshalEnum(record.GetEntryTypeDailyWorkPeriod())
	canvas[offset] = entryTypeDailyWorkPeriodByte
	offset += 1

	// dailyWorkPeriodCountry (1 byte)
	dailyWorkPeriodCountryByte, _ := MarshalEnum(record.GetDailyWorkPeriodCountry())
	canvas[offset] = dailyWorkPeriodCountryByte
	offset += 1

	// dailyWorkPeriodRegion (1 byte)
	dailyWorkPeriodRegionByte, _ := MarshalEnum(record.GetDailyWorkPeriodRegion())
	canvas[offset] = dailyWorkPeriodRegionByte
	offset += 1

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerBytes, err := opts.MarshalOdometer(record.GetVehicleOdometerKm())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle odometer value: %w", err)
	}
	copy(canvas[offset:offset+3], vehicleOdometerBytes)
	offset += 3

	// entryGNSSPlaceAuthRecord (12 bytes)
	entryGNSSPlaceAuthRecordBytes, err := opts.MarshalGNSSPlaceAuthRecord(record.GetEntryGnssPlaceAuthRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry GNSS place auth record: %w", err)
	}
	copy(canvas[offset:offset+12], entryGNSSPlaceAuthRecordBytes)

	return canvas[:], nil
}
