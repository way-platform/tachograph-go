package vu

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalActivitiesGen1 parses Gen1 Activities data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
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
// - Signature: 128 bytes (RSA-1024)
//
func unmarshalActivitiesGen1(value []byte) (*vuv1.ActivitiesGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	activities := &vuv1.ActivitiesGen1{}
	activities.SetRawData(value) // Store complete transfer value for painting

	offset := 0
	opts := dd.UnmarshalOptions{PreserveRawData: true}

	// TimeReal (4 bytes) - date of day downloaded
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for TimeReal")
	}
	timeReal, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal TimeReal: %w", err)
	}
	activities.SetDateOfDay(timeReal)
	offset += 4

	// OdometerValueMidnight (3 bytes - OdometerShort)
	if offset+3 > len(data) {
		return nil, fmt.Errorf("insufficient data for OdometerValueMidnight")
	}
	odometer, err := opts.UnmarshalOdometer(data[offset : offset+3])
	if err != nil {
		return nil, fmt.Errorf("unmarshal OdometerValueMidnight: %w", err)
	}
	activities.SetOdometerMidnightKm(int32(odometer))
	offset += 3

	// VuCardIWData: 2 bytes (noOfIWRecords) + (noOfIWRecords * 129 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfIWRecords")
	}
	noOfIWRecords := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	// Parse each CardIWRecord (129 bytes each for Gen1)
	cardIWRecords := make([]*ddv1.VuCardIWRecord, noOfIWRecords)
	for i := uint16(0); i < noOfIWRecords; i++ {
		const cardIWRecordSize = 129
		if offset+cardIWRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for CardIWRecord %d", i)
		}

		record, err := opts.UnmarshalVuCardIWRecord(data[offset : offset+cardIWRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal CardIWRecord %d: %w", i, err)
		}

		cardIWRecords[i] = record
		offset += cardIWRecordSize
	}
	activities.SetCardIwData(cardIWRecords)

	// VuActivityDailyData: 2 bytes (noOfActivityChanges) + (noOfActivityChanges * 2 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfActivityChanges")
	}
	noOfActivityChanges := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	// Parse each ActivityChangeInfo (2 bytes each)
	activityChanges := make([]*ddv1.ActivityChangeInfo, noOfActivityChanges)
	for i := uint16(0); i < noOfActivityChanges; i++ {
		const activityChangeSize = 2
		if offset+activityChangeSize > len(data) {
			return nil, fmt.Errorf("insufficient data for ActivityChangeInfo %d", i)
		}

		activityChange, err := opts.UnmarshalActivityChangeInfo(data[offset : offset+activityChangeSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal activity change %d: %w", i, err)
		}
		activityChanges[i] = activityChange
		offset += activityChangeSize
	}
	activities.SetActivityChanges(activityChanges)

	// VuPlaceDailyWorkPeriodData: 1 byte (noOfPlaceRecords) + (noOfPlaceRecords * 28 bytes)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfPlaceRecords")
	}
	noOfPlaceRecords := data[offset]
	offset += 1

	// Parse each VuPlaceDailyWorkPeriodRecord (28 bytes each)
	placeRecords := make([]*ddv1.VuPlaceDailyWorkPeriodRecord, noOfPlaceRecords)
	for i := uint8(0); i < noOfPlaceRecords; i++ {
		const placeRecordSize = 28 // 18 bytes FullCardNumber + 10 bytes PlaceRecord
		if offset+placeRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for VuPlaceDailyWorkPeriodRecord %d", i)
		}

		vuPlaceRecord, err := opts.UnmarshalVuPlaceDailyWorkPeriodRecord(data[offset : offset+placeRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuPlaceDailyWorkPeriodRecord %d: %w", i, err)
		}

		placeRecords[i] = vuPlaceRecord
		offset += placeRecordSize
	}
	activities.SetPlaceRecords(placeRecords)

	// VuSpecificConditionData: 2 bytes (noOfSpecificConditionRecords) + (noOfSpecificConditionRecords * 5 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfSpecificConditionRecords")
	}
	noOfSpecificConditionRecords := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	// Parse each SpecificConditionRecord (5 bytes each)
	specificConditions := make([]*ddv1.SpecificConditionRecord, noOfSpecificConditionRecords)
	for i := uint16(0); i < noOfSpecificConditionRecords; i++ {
		const specificConditionSize = 5
		if offset+specificConditionSize > len(data) {
			return nil, fmt.Errorf("insufficient data for SpecificConditionRecord %d", i)
		}

		specificCondition, err := opts.UnmarshalSpecificConditionRecord(data[offset : offset+specificConditionSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal specific condition %d: %w", i, err)
		}
		specificConditions[i] = specificCondition
		offset += specificConditionSize
	}
	activities.SetSpecificConditions(specificConditions)

	// Store signature (extracted at the beginning)
	activities.SetSignature(signature)

	// Verify we consumed exactly the right amount of data
	if offset != len(data) {
		return nil, fmt.Errorf("Activities Gen1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return activities, nil
}

// MarshalActivitiesGen1 marshals Gen1 Activities data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and has the correct length, it uses it as a canvas and paints semantic values over it.
// Otherwise, it creates a zero-filled canvas and encodes from semantic fields.
func (opts MarshalOptions) MarshalActivitiesGen1(activities *vuv1.ActivitiesGen1) ([]byte, error) {
	if activities == nil {
		return nil, fmt.Errorf("activities cannot be nil")
	}

	// Calculate expected size (signature is stored separately and appended at the end)
	noOfIWRecords := len(activities.GetCardIwData())
	noOfActivityChanges := len(activities.GetActivityChanges())
	noOfPlaceRecords := len(activities.GetPlaceRecords())
	noOfSpecificConditions := len(activities.GetSpecificConditions())

	// Fixed header: 4 (TimeReal) + 3 (OdometerShort)
	// VuCardIWData: 2 (count) + N*129 (records)
	// VuActivityDailyData: 2 (count) + M*2 (records)
	// VuPlaceDailyWorkPeriodData: 1 (count) + P*28 (records)
	// VuSpecificConditionData: 2 (count) + Q*5 (records)
	const headerSize = 4 + 3
	dataSize := headerSize + 2 + (noOfIWRecords * 129) + 2 + (noOfActivityChanges * 2) + 1 + (noOfPlaceRecords * 28) + 2 + (noOfSpecificConditions * 5)

	// Use raw_data as canvas if available
	var canvas []byte
	raw := activities.GetRawData()
	if len(raw) == dataSize+128 { // dataSize + 128-byte signature
		// Extract data portion (without signature) to use as canvas
		canvas = make([]byte, dataSize)
		copy(canvas, raw[:dataSize])
	} else if len(raw) == dataSize {
		// raw_data is just the data portion without signature
		canvas = make([]byte, dataSize)
		copy(canvas, raw)
	} else {
		// Create zero-filled canvas
		canvas = make([]byte, dataSize)
	}

	offset := 0

	// TimeReal (4 bytes)
	timeRealBytes, err := opts.MarshalTimeReal(activities.GetDateOfDay())
	if err != nil {
		return nil, fmt.Errorf("marshal TimeReal: %w", err)
	}
	copy(canvas[offset:offset+4], timeRealBytes)
	offset += 4

	// OdometerValueMidnight (3 bytes - OdometerShort)
	odometerBytes, err := opts.MarshalOdometer(activities.GetOdometerMidnightKm())
	if err != nil {
		return nil, fmt.Errorf("marshal OdometerValueMidnight: %w", err)
	}
	copy(canvas[offset:offset+3], odometerBytes)
	offset += 3

	// VuCardIWData: 2 bytes (count) + records
	binary.BigEndian.PutUint16(canvas[offset:offset+2], uint16(noOfIWRecords))
	offset += 2

	for i, record := range activities.GetCardIwData() {
		recordBytes, err := opts.MarshalVuCardIWRecord(record)
		if err != nil {
			return nil, fmt.Errorf("marshal CardIWRecord %d: %w", i, err)
		}
		if len(recordBytes) != 129 {
			return nil, fmt.Errorf("CardIWRecord %d has invalid length: got %d, want 129", i, len(recordBytes))
		}
		copy(canvas[offset:offset+129], recordBytes)
		offset += 129
	}

	// VuActivityDailyData: 2 bytes (count) + records
	binary.BigEndian.PutUint16(canvas[offset:offset+2], uint16(noOfActivityChanges))
	offset += 2

	for i, activityChange := range activities.GetActivityChanges() {
		activityBytes, err := opts.MarshalActivityChangeInfo(activityChange)
		if err != nil {
			return nil, fmt.Errorf("marshal ActivityChangeInfo %d: %w", i, err)
		}
		if len(activityBytes) != 2 {
			return nil, fmt.Errorf("ActivityChangeInfo %d has invalid length: got %d, want 2", i, len(activityBytes))
		}
		copy(canvas[offset:offset+2], activityBytes)
		offset += 2
	}

	// VuPlaceDailyWorkPeriodData: 1 byte (count) + records
	canvas[offset] = byte(noOfPlaceRecords)
	offset += 1

	// Marshal each VuPlaceDailyWorkPeriodRecord (28 bytes each)
	for i, placeRecord := range activities.GetPlaceRecords() {
		recordBytes, err := opts.MarshalVuPlaceDailyWorkPeriodRecord(placeRecord)
		if err != nil {
			return nil, fmt.Errorf("marshal VuPlaceDailyWorkPeriodRecord %d: %w", i, err)
		}
		if len(recordBytes) != 28 {
			return nil, fmt.Errorf("VuPlaceDailyWorkPeriodRecord %d has invalid length: got %d, want 28", i, len(recordBytes))
		}
		copy(canvas[offset:offset+28], recordBytes)
		offset += 28
	}

	// VuSpecificConditionData: 2 bytes (count) + records
	binary.BigEndian.PutUint16(canvas[offset:offset+2], uint16(noOfSpecificConditions))
	offset += 2

	for i, specificCondition := range activities.GetSpecificConditions() {
		specificBytes, err := opts.MarshalSpecificConditionRecord(specificCondition)
		if err != nil {
			return nil, fmt.Errorf("marshal SpecificConditionRecord %d: %w", i, err)
		}
		if len(specificBytes) != 5 {
			return nil, fmt.Errorf("SpecificConditionRecord %d has invalid length: got %d, want 5", i, len(specificBytes))
		}
		copy(canvas[offset:offset+5], specificBytes)
		offset += 5
	}

	// Verify we used exactly the expected amount of data
	if offset != dataSize {
		return nil, fmt.Errorf("Activities Gen1 marshalling mismatch: wrote %d bytes, expected %d", offset, dataSize)
	}

	// Append signature to create complete transfer value
	signature := activities.GetSignature()
	if len(signature) == 0 {
		// Gen1 uses fixed 128-byte RSA-1024 signatures
		signature = make([]byte, 128)
	}
	if len(signature) != 128 {
		return nil, fmt.Errorf("invalid signature length: got %d, want 128", len(signature))
	}

	transferValue := append(canvas, signature...)
	return transferValue, nil
}

// anonymizeActivitiesGen1 anonymizes Gen1 Activities data.
func (opts AnonymizeOptions) anonymizeActivitiesGen1(activities *vuv1.ActivitiesGen1) *vuv1.ActivitiesGen1 {
	if activities == nil {
		return nil
	}
	result := proto.Clone(activities).(*vuv1.ActivitiesGen1)

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize card IW data (cards inserted/withdrawn)
	var anonymizedCardIWRecords []*ddv1.VuCardIWRecord
	for _, record := range result.GetCardIwData() {
		if record == nil {
			continue
		}
		anonRecord := proto.Clone(record).(*ddv1.VuCardIWRecord)

		// Anonymize holder name
		if origHolderName := record.GetCardHolderName(); origHolderName != nil {
			holderName := &ddv1.HolderName{}
			holderName.SetHolderSurname(ddOpts.AnonymizeStringValue(origHolderName.GetHolderSurname()))
			holderName.SetHolderFirstNames(ddOpts.AnonymizeStringValue(origHolderName.GetHolderFirstNames()))
			anonRecord.SetCardHolderName(holderName)
		}

		// Anonymize card number
		anonRecord.SetFullCardNumber(ddOpts.AnonymizeFullCardNumber(record.GetFullCardNumber()))

		// Anonymize card expiry date (preserve or anonymize timestamp)
		if expiryDate := record.GetCardExpiryDate(); expiryDate != nil && !opts.PreserveTimestamps {
			anonRecord.SetCardExpiryDate(dd.NewDate(2025, 12, 31))
		}

		// Anonymize card insertion/withdrawal timestamps
		anonRecord.SetCardInsertionTime(ddOpts.AnonymizeTimestamp(record.GetCardInsertionTime()))
		anonRecord.SetCardSlotNumber(record.GetCardSlotNumber()) // Not PII

		// Anonymize card withdrawal time if present
		if record.GetCardWithdrawalTime() != nil {
			anonRecord.SetCardWithdrawalTime(ddOpts.AnonymizeTimestamp(record.GetCardWithdrawalTime()))
		}

		// Anonymize odometer values
		anonRecord.SetOdometerAtInsertionKm(ddOpts.AnonymizeOdometerValue(record.GetOdometerAtInsertionKm()))
		anonRecord.SetOdometerAtWithdrawalKm(ddOpts.AnonymizeOdometerValue(record.GetOdometerAtWithdrawalKm()))

		// Anonymize previous vehicle info (contains vehicle registration which is PII)
		if prevVehicle := record.GetPreviousVehicleInfo(); prevVehicle != nil {
			anonRecord.SetPreviousVehicleInfo(ddOpts.AnonymizePreviousVehicleInfo(prevVehicle))
		}

		// Manual input flag is not PII - keep as-is
		anonRecord.SetManualInputFlag(record.GetManualInputFlag())

		// Clear raw_data
		anonRecord.ClearRawData()
		anonymizedCardIWRecords = append(anonymizedCardIWRecords, anonRecord)
	}
	result.SetCardIwData(anonymizedCardIWRecords)

	// Anonymize activity changes
	var anonymizedActivityChanges []*ddv1.ActivityChangeInfo
	for _, activity := range result.GetActivityChanges() {
		if activity == nil {
			continue
		}
		anonActivity := proto.Clone(activity).(*ddv1.ActivityChangeInfo)

		// Activity changes don't contain PII - just clear raw_data
		// (slot, crew mode, inserted status, activity type, time of change are not personally identifiable)

		// Clear raw_data
		anonActivity.ClearRawData()
		anonymizedActivityChanges = append(anonymizedActivityChanges, anonActivity)
	}
	result.SetActivityChanges(anonymizedActivityChanges)

	// Anonymize place records (VU place daily work period records)
	var anonymizedPlaceRecords []*ddv1.VuPlaceDailyWorkPeriodRecord
	for _, placeRecord := range result.GetPlaceRecords() {
		anonymizedPlaceRecords = append(anonymizedPlaceRecords, ddOpts.AnonymizeVuPlaceDailyWorkPeriodRecord(placeRecord))
	}
	result.SetPlaceRecords(anonymizedPlaceRecords)

	// Anonymize specific condition records (timestamps are PII)
	var anonymizedSpecificConditions []*ddv1.SpecificConditionRecord
	for _, sc := range result.GetSpecificConditions() {
		if sc == nil {
			continue
		}
		anonSC := proto.Clone(sc).(*ddv1.SpecificConditionRecord)
		anonSC.SetEntryTime(ddOpts.AnonymizeTimestamp(sc.GetEntryTime()))
		anonymizedSpecificConditions = append(anonymizedSpecificConditions, anonSC)
	}
	result.SetSpecificConditions(anonymizedSpecificConditions)

	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Clear raw_data to force semantic marshalling
	result.ClearRawData()

	return result
}
