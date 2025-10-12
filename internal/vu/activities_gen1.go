package vu

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalActivitiesGen1 parses Gen1 Activities data from the complete transfer value.
//
// Gen1 Activities structure (from Data Dictionary and Appendix 7, Section 2.2.6.3):
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    timeReal                      TimeReal,                              -- 4 bytes
//	    odometerValueMidnight         OdometerShort,                         -- 3 bytes
//	    vuCardIWData                  VuCardIWDataFirstGen,                  -- 2 + (N * 129) bytes
//	    vuActivityDailyData           VuActivityDailyDataFirstGen,           -- 2 + (M * 2) bytes
//	    vuPlaceDailyWorkPeriodData    VuPlaceDailyWorkPeriodDataFirstGen,    -- 1 + (P * 28) bytes
//	    vuSpecificConditionData       VuSpecificConditionDataFirstGen,       -- 2 + (Q * 5) bytes
//	    signature                     SignatureFirstGen                      -- 128 bytes (RSA)
//	}
//
// Binary Layout:
// - TimeReal: 4 bytes (date of day downloaded)
// - OdometerValueMidnight: 3 bytes (OdometerShort)
// - VuCardIWData: 2 bytes (noOfIWRecords) + (noOfIWRecords * 129 bytes)
//   - Each VuCardIWRecordFirstGen: 129 bytes
//   - FullCardNumber: 18 bytes
//   - ManufacturerCode: 1 byte
//   - DownloadTime: 4 bytes
//   - ... (rest of record)
//
// - VuActivityDailyData: 2 bytes (noOfActivityChanges) + (noOfActivityChanges * 2 bytes)
//   - Each ActivityChangeInfo: 2 bytes
//
// - VuPlaceDailyWorkPeriodData: 1 byte (noOfPlaceRecords) + (noOfPlaceRecords * 28 bytes)
//   - Each VuPlaceDailyWorkPeriodRecordFirstGen: 28 bytes
//   - FullCardNumber: 18 bytes
//   - PlaceRecord: 10 bytes
//
// - VuSpecificConditionData: 2 bytes (noOfSpecificConditionRecords) + (noOfSpecificConditionRecords * 5 bytes)
//   - Each SpecificConditionRecord: 5 bytes
//   - TimeReal: 4 bytes
//   - SpecificConditionType: 1 byte
//
// - Signature: 128 bytes (RSA)
//
// Note: This is a minimal implementation that validates the binary structure and stores raw_data.
// Full semantic parsing of all nested records is TODO.
func unmarshalActivitiesGen1(value []byte) (*vuv1.ActivitiesGen1, error) {
	activities := &vuv1.ActivitiesGen1{}
	activities.SetRawData(value)

	offset := 0
	var opts dd.UnmarshalOptions

	// TimeReal (4 bytes) - date of day downloaded
	if offset+4 > len(value) {
		return nil, fmt.Errorf("insufficient data for TimeReal")
	}
	timeReal, err := opts.UnmarshalTimeReal(value[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal TimeReal: %w", err)
	}
	activities.SetDateOfDay(timeReal)
	offset += 4

	// OdometerValueMidnight (3 bytes - OdometerShort)
	if offset+3 > len(value) {
		return nil, fmt.Errorf("insufficient data for OdometerValueMidnight")
	}
	odometer, err := opts.UnmarshalOdometer(value[offset : offset+3])
	if err != nil {
		return nil, fmt.Errorf("unmarshal OdometerValueMidnight: %w", err)
	}
	activities.SetOdometerMidnightKm(int32(odometer))
	offset += 3

	// VuCardIWData: 2 bytes (noOfIWRecords) + (noOfIWRecords * 129 bytes)
	if offset+2 > len(value) {
		return nil, fmt.Errorf("insufficient data for noOfIWRecords")
	}
	noOfIWRecords := binary.BigEndian.Uint16(value[offset : offset+2])
	offset += 2

	// Parse each CardIWRecord (129 bytes each for Gen1)
	cardIWRecords := make([]*ddv1.VuCardIWRecord, noOfIWRecords)
	for i := uint16(0); i < noOfIWRecords; i++ {
		const cardIWRecordSize = 129
		if offset+cardIWRecordSize > len(value) {
			return nil, fmt.Errorf("insufficient data for CardIWRecord %d", i)
		}

		record, err := opts.UnmarshalVuCardIWRecord(value[offset : offset+cardIWRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal CardIWRecord %d: %w", i, err)
		}

		cardIWRecords[i] = record
		offset += cardIWRecordSize
	}
	activities.SetCardIwData(cardIWRecords)

	// VuActivityDailyData: 2 bytes (noOfActivityChanges) + (noOfActivityChanges * 2 bytes)
	if offset+2 > len(value) {
		return nil, fmt.Errorf("insufficient data for noOfActivityChanges")
	}
	noOfActivityChanges := binary.BigEndian.Uint16(value[offset : offset+2])
	offset += 2

	// Parse each ActivityChangeInfo (2 bytes each)
	activityChanges := make([]*ddv1.ActivityChangeInfo, noOfActivityChanges)
	for i := uint16(0); i < noOfActivityChanges; i++ {
		const activityChangeSize = 2
		if offset+activityChangeSize > len(value) {
			return nil, fmt.Errorf("insufficient data for ActivityChangeInfo %d", i)
		}

		activityChange, err := opts.UnmarshalActivityChangeInfo(value[offset : offset+activityChangeSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal activity change %d: %w", i, err)
		}
		activityChanges[i] = activityChange
		offset += activityChangeSize
	}
	activities.SetActivityChanges(activityChanges)

	// VuPlaceDailyWorkPeriodData: 1 byte (noOfPlaceRecords) + (noOfPlaceRecords * 28 bytes)
	// Note: Each record is 28 bytes (18 FullCardNumber + 10 PlaceRecord)
	if offset+1 > len(value) {
		return nil, fmt.Errorf("insufficient data for noOfPlaceRecords")
	}
	noOfPlaceRecords := value[offset]
	offset += 1

	// Parse each VuPlaceDailyWorkPeriodRecord (28 bytes each)
	placeRecords := make([]*vuv1.ActivitiesGen1_PlaceRecord, noOfPlaceRecords)
	for i := uint8(0); i < noOfPlaceRecords; i++ {
		const placeRecordSize = 28 // 18 bytes FullCardNumber + 10 bytes PlaceRecord
		if offset+placeRecordSize > len(value) {
			return nil, fmt.Errorf("insufficient data for PlaceRecord %d", i)
		}

		record := &vuv1.ActivitiesGen1_PlaceRecord{}
		recordOffset := 0

		// Skip FullCardNumber (18 bytes) - not exposed in the proto for Gen1 PlaceRecord
		// This is the card number associated with this place entry
		recordOffset += 18

		// PlaceRecord (10 bytes)
		// entryTime (4 bytes)
		entryTime, err := opts.UnmarshalTimeReal(value[offset+recordOffset : offset+recordOffset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal place entry time: %w", err)
		}
		record.SetEntryTime(entryTime)
		recordOffset += 4

		// entryTypeDailyWorkPeriod (1 byte)
		entryType, err := dd.UnmarshalEnum[ddv1.EntryTypeDailyWorkPeriod](value[offset+recordOffset])
		if err != nil {
			return nil, fmt.Errorf("unmarshal entry type: %w", err)
		}
		record.SetEntryType(entryType)
		recordOffset += 1

		// dailyWorkPeriodCountry (1 byte)
		country, err := dd.UnmarshalEnum[ddv1.NationNumeric](value[offset+recordOffset])
		if err != nil {
			return nil, fmt.Errorf("unmarshal country: %w", err)
		}
		record.SetCountry(country)
		recordOffset += 1

		// dailyWorkPeriodRegion (1 byte)
		region := value[offset+recordOffset : offset+recordOffset+1]
		record.SetRegion(region)
		recordOffset += 1

		// vehicleOdometerValue (3 bytes)
		odometerValue, err := opts.UnmarshalOdometer(value[offset+recordOffset : offset+recordOffset+3])
		if err != nil {
			return nil, fmt.Errorf("unmarshal place odometer: %w", err)
		}
		record.SetOdometerKm(int32(odometerValue))
		recordOffset += 3

		placeRecords[i] = record
		offset += placeRecordSize
	}
	activities.SetPlaces(placeRecords)

	// VuSpecificConditionData: 2 bytes (noOfSpecificConditionRecords) + (noOfSpecificConditionRecords * 5 bytes)
	if offset+2 > len(value) {
		return nil, fmt.Errorf("insufficient data for noOfSpecificConditionRecords")
	}
	noOfSpecificConditionRecords := binary.BigEndian.Uint16(value[offset : offset+2])
	offset += 2

	// Parse each SpecificConditionRecord (5 bytes each)
	specificConditions := make([]*ddv1.SpecificConditionRecord, noOfSpecificConditionRecords)
	for i := uint16(0); i < noOfSpecificConditionRecords; i++ {
		const specificConditionSize = 5
		if offset+specificConditionSize > len(value) {
			return nil, fmt.Errorf("insufficient data for SpecificConditionRecord %d", i)
		}

		specificCondition, err := opts.UnmarshalSpecificConditionRecord(value[offset : offset+specificConditionSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal specific condition %d: %w", i, err)
		}
		specificConditions[i] = specificCondition
		offset += specificConditionSize
	}
	activities.SetSpecificConditions(specificConditions)

	// Signature (128 bytes - RSA for Gen1)
	if offset+128 > len(value) {
		return nil, fmt.Errorf("insufficient data for Signature")
	}
	activities.SetSignature(value[offset : offset+128])
	offset += 128

	// Verify we consumed exactly the right amount of data
	if offset != len(value) {
		return nil, fmt.Errorf("Activities Gen1 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	return activities, nil
}

// appendActivitiesGen1 marshals Gen1 Activities data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and has the correct length, it uses it as a canvas and paints semantic values over it.
// Otherwise, it creates a zero-filled canvas and encodes from semantic fields.
func appendActivitiesGen1(dst []byte, activities *vuv1.ActivitiesGen1) ([]byte, error) {
	if activities == nil {
		return nil, fmt.Errorf("activities cannot be nil")
	}

	// For now, use raw_data directly if available
	// Full semantic marshalling requires implementing all record types
	raw := activities.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	// TODO: Implement marshalling from semantic fields
	// This would require:
	// 1. Writing TimeReal
	// 2. Writing OdometerValueMidnight
	// 3. Writing VuCardIWData (count + records)
	// 4. Writing VuActivityDailyData (count + records)
	// 5. Writing VuPlaceDailyWorkPeriodData (count + records)
	// 6. Writing VuSpecificConditionData (count + records)
	// 7. Writing Signature
	return nil, fmt.Errorf("cannot marshal Activities Gen1 without raw_data (semantic marshalling not yet implemented)")
}
