package vu

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalActivitiesGen2V1 parses Gen2 V1 Activities data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen2 V1 Activities structure uses RecordArray format (from Data Dictionary):
//
// ASN.1 Definition:
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    timeRealRecordArray                   TimeRealRecordArray,
//	    odometerValueMidnightRecordArray      OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                   VuCardIWRecordArray,
//	    vuActivityDailyRecordArray            VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray     VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                   VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray        VuSpecificConditionRecordArray,
//	    signatureRecordArray                  SignatureRecordArray
//	}
//
// Each RecordArray has a 5-byte header:
//
//	recordType (1 byte) + recordSize (2 bytes, big-endian) + noOfRecords (2 bytes, big-endian)
func unmarshalActivitiesGen2V1(value []byte) (*vuv1.ActivitiesGen2V1, error) {
	// Split transfer value into data and signature
	// Gen2 uses variable-length ECDSA signatures stored as SignatureRecordArray
	// We use the sizeOf function to determine where to split
	totalSize, signatureSize, err := sizeOfActivitiesGen2V1(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	activities := &vuv1.ActivitiesGen2V1{}
	activities.SetRawData(value) // Store complete transfer value for painting

	offset := 0

	// TimeRealRecordArray
	dateOfDay, bytesRead, err := parseTimeRealRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse TimeRealRecordArray: %w", err)
	}
	activities.SetDateOfDay(dateOfDay)
	offset += bytesRead

	// OdometerValueMidnightRecordArray
	odometerMidnightKm, bytesRead, err := parseOdometerValueMidnightRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse OdometerValueMidnightRecordArray: %w", err)
	}
	activities.SetOdometerMidnightKm(odometerMidnightKm)
	offset += bytesRead

	// VuCardIWRecordArray (Gen2 - 132 bytes per record)
	cardIWRecords, bytesRead, err := parseVuCardIWRecordArrayG2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuCardIWRecordArray: %w", err)
	}
	activities.SetCardIwData(cardIWRecords)
	offset += bytesRead

	// VuActivityDailyRecordArray
	activityChanges, bytesRead, err := parseVuActivityDailyRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuActivityDailyRecordArray: %w", err)
	}
	activities.SetActivityChanges(activityChanges)
	offset += bytesRead

	// VuPlaceDailyWorkPeriodRecordArray (Gen2v1 - 41 bytes per record)
	vuPlaceRecords, bytesRead, err := parseVuPlaceDailyWorkPeriodRecordArrayG2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuPlaceDailyWorkPeriodRecordArray: %w", err)
	}
	// Extract PlaceRecordG2 from VuPlaceDailyWorkPeriodRecordG2 wrapper
	placeRecords := make([]*ddv1.PlaceRecordG2, 0, len(vuPlaceRecords))
	for _, vuPlaceRec := range vuPlaceRecords {
		placeRecords = append(placeRecords, vuPlaceRec.GetPlaceRecord())
	}
	activities.SetPlaces(placeRecords)
	offset += bytesRead

	// VuGNSSADRecordArray (Gen2v1 - 58 bytes per record)
	gnssADRecords, bytesRead, err := parseVuGNSSADRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuGNSSADRecordArray: %w", err)
	}
	activities.SetGnssAccumulatedDriving(gnssADRecords)
	offset += bytesRead

	// VuSpecificConditionRecordArray
	specificConditions, bytesRead, err := parseVuSpecificConditionRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuSpecificConditionRecordArray: %w", err)
	}
	activities.SetSpecificConditions(specificConditions)
	offset += bytesRead

	// Store signature (extracted at the beginning)
	activities.SetSignature(signature)

	// Verify we consumed exactly the right amount of data
	if offset != len(data) {
		return nil, fmt.Errorf("Activities Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return activities, nil
}

// MarshalActivitiesGen2V1 marshals Gen2 V1 Activities data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and matches the structure, it uses it as the output.
func (opts MarshalOptions) MarshalActivitiesGen2V1(activities *vuv1.ActivitiesGen2V1) ([]byte, error) {
	if activities == nil {
		return nil, fmt.Errorf("activities cannot be nil")
	}

	// For Gen2 structures with RecordArrays, raw data painting is straightforward
	raw := activities.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	// Marshal from semantic fields
	var result []byte
	marshalOpts := dd.MarshalOptions{}

	// TimeRealRecordArray (1 record of 4 bytes)
	timeRealData, err := marshalOpts.MarshalTimeReal(activities.GetDateOfDay())
	if err != nil {
		return nil, fmt.Errorf("marshal TimeReal: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x01, 4, 1)
	result = append(result, timeRealData...)

	// OdometerValueMidnightRecordArray (1 record of 3 bytes)
	odometerData, err := marshalOpts.MarshalOdometer(activities.GetOdometerMidnightKm())
	if err != nil {
		return nil, fmt.Errorf("marshal OdometerValueMidnight: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x02, 3, 1)
	result = append(result, odometerData...)

	// VuCardIWRecordArray (Gen2 - 132 bytes per record)
	cardIWData, err := marshalCardIWRecordsG2(activities.GetCardIwData())
	if err != nil {
		return nil, fmt.Errorf("marshal VuCardIWRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, 131, uint16(len(activities.GetCardIwData())))
	result = append(result, cardIWData...)

	// VuActivityDailyRecordArray (2 bytes per record)
	activityData, err := marshalActivityChangeInfos(activities.GetActivityChanges())
	if err != nil {
		return nil, fmt.Errorf("marshal VuActivityDailyRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x04, 2, uint16(len(activities.GetActivityChanges())))
	result = append(result, activityData...)

	// VuPlaceDailyWorkPeriodRecordArray (Gen2v1 - 40 bytes per record)
	placeData, err := marshalPlaceRecordsG2V1(activities.GetPlaces())
	if err != nil {
		return nil, fmt.Errorf("marshal VuPlaceDailyWorkPeriodRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x05, 40, uint16(len(activities.GetPlaces())))
	result = append(result, placeData...)

	// VuGNSSADRecordArray (Gen2v1 - 56 bytes per record)
	gnssData, err := marshalGnssAccumulatedDrivingRecordsV1(activities.GetGnssAccumulatedDriving())
	if err != nil {
		return nil, fmt.Errorf("marshal VuGNSSADRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x06, 56, uint16(len(activities.GetGnssAccumulatedDriving())))
	result = append(result, gnssData...)

	// VuSpecificConditionRecordArray (5 bytes per record)
	specificCondData, err := marshalSpecificConditionRecords(activities.GetSpecificConditions())
	if err != nil {
		return nil, fmt.Errorf("marshal VuSpecificConditionRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x07, 5, uint16(len(activities.GetSpecificConditions())))
	result = append(result, specificCondData...)

	result = appendSignature(result, activities.GetSignature(), 0x09)

	return result, nil
}

// Helper functions for parsing Gen2 V1 RecordArrays

// parseRecordArrayHeader parses the 5-byte RecordArray header.
// Returns: recordType, recordSize, noOfRecords, bytesConsumed, error
func parseRecordArrayHeader(data []byte, offset int) (byte, uint16, uint16, int, error) {
	const headerSize = 5
	if offset+headerSize > len(data) {
		return 0, 0, 0, 0, fmt.Errorf("insufficient data for RecordArray header at offset %d", offset)
	}

	recordType := data[offset]
	recordSize := binary.BigEndian.Uint16(data[offset+1 : offset+3])
	noOfRecords := binary.BigEndian.Uint16(data[offset+3 : offset+5])

	return recordType, recordSize, noOfRecords, headerSize, nil
}

// parseTimeRealRecordArray parses a TimeRealRecordArray (should have 1 record of 4 bytes).
func parseTimeRealRecordArray(data []byte, offset int) (*timestamppb.Timestamp, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 TimeReal record, got %d", noOfRecords)
	}

	if recordSize != 4 {
		return nil, 0, fmt.Errorf("expected TimeReal record size 4, got %d", recordSize)
	}

	recordStart := offset + headerSize
	var opts dd.UnmarshalOptions
	timeReal, err := opts.UnmarshalTimeReal(data[recordStart : recordStart+int(recordSize)])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal TimeReal: %w", err)
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return timeReal, totalSize, nil
}

// parseOdometerValueMidnightRecordArray parses an OdometerValueMidnightRecordArray (should have 1 record of 3 bytes).
func parseOdometerValueMidnightRecordArray(data []byte, offset int) (int32, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return 0, 0, err
	}

	if noOfRecords != 1 {
		return 0, 0, fmt.Errorf("expected 1 OdometerValueMidnight record, got %d", noOfRecords)
	}

	if recordSize != 3 {
		return 0, 0, fmt.Errorf("expected OdometerValueMidnight record size 3, got %d", recordSize)
	}

	recordStart := offset + headerSize
	var opts dd.UnmarshalOptions
	odometer, err := opts.UnmarshalOdometer(data[recordStart : recordStart+int(recordSize)])
	if err != nil {
		return 0, 0, fmt.Errorf("unmarshal Odometer: %w", err)
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return int32(odometer), totalSize, nil
}

// parseVuCardIWRecordArrayG2 parses a VuCardIWRecordArray (Gen2 - 132 bytes per record).
func parseVuCardIWRecordArrayG2(data []byte, offset int) ([]*ddv1.VuCardIWRecordG2, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 131 // Gen2
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuCardIWRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	opts.PreserveRawData = true

	records := make([]*ddv1.VuCardIWRecordG2, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := uint16(0); i < noOfRecords; i++ {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuCardIWRecord %d", i)
		}

		record, err := opts.UnmarshalVuCardIWRecordG2(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuCardIWRecord %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuActivityDailyRecordArray parses a VuActivityDailyRecordArray (2 bytes per record).
func parseVuActivityDailyRecordArray(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 2
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected ActivityChangeInfo size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	opts.PreserveRawData = true

	records := make([]*ddv1.ActivityChangeInfo, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := uint16(0); i < noOfRecords; i++ {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for ActivityChangeInfo %d", i)
		}

		record, err := opts.UnmarshalActivityChangeInfo(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal ActivityChangeInfo %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuPlaceDailyWorkPeriodRecordArrayG2 parses a VuPlaceDailyWorkPeriodRecordArray (Gen2v1 - 40 bytes per record).
func parseVuPlaceDailyWorkPeriodRecordArrayG2(data []byte, offset int) ([]*ddv1.VuPlaceDailyWorkPeriodRecordG2, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 40 // Gen2v1
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuPlaceDailyWorkPeriodRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	opts.PreserveRawData = true

	records := make([]*ddv1.VuPlaceDailyWorkPeriodRecordG2, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := uint16(0); i < noOfRecords; i++ {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuPlaceDailyWorkPeriodRecord %d", i)
		}

		record, err := opts.UnmarshalVuPlaceDailyWorkPeriodRecordG2(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuPlaceDailyWorkPeriodRecord %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuGNSSADRecordArray parses a VuGNSSADRecordArray (Gen2v1 - 56 bytes per record).
func parseVuGNSSADRecordArray(data []byte, offset int) ([]*ddv1.VuGNSSADRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 56 // Gen2v1
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuGNSSADRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	opts := dd.UnmarshalOptions{PreserveRawData: true}

	records := make([]*ddv1.VuGNSSADRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := uint16(0); i < noOfRecords; i++ {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuGNSSADRecord %d", i)
		}

		record, err := opts.UnmarshalVuGNSSADRecord(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuGNSSADRecord %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuSpecificConditionRecordArray parses a VuSpecificConditionRecordArray (5 bytes per record).
func parseVuSpecificConditionRecordArray(data []byte, offset int) ([]*ddv1.SpecificConditionRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 5
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected SpecificConditionRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	opts.PreserveRawData = true

	records := make([]*ddv1.SpecificConditionRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := uint16(0); i < noOfRecords; i++ {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for SpecificConditionRecord %d", i)
		}

		record, err := opts.UnmarshalSpecificConditionRecord(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal SpecificConditionRecord %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// Helper functions for marshalling Gen2 V1 RecordArrays

// appendRecordArrayHeader appends a 5-byte RecordArray header.
// Header format: recordType (1 byte) + recordSize (2 bytes BE) + noOfRecords (2 bytes BE)
func appendRecordArrayHeader(dst []byte, recordType byte, recordSize uint16, noOfRecords uint16) []byte {
	dst = append(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)
	return dst
}

// appendSignature appends signature bytes to dst.
//
// The signature field stores the complete SignatureRecordArray (5-byte header + sig bytes).
// When the signature is empty (anonymized data has no real signature), a placeholder 5-byte
// empty RecordArray header is appended so that sizeOf* functions can still parse the output.
func appendSignature(dst []byte, signature []byte, recordType byte) []byte {
	if len(signature) > 0 {
		return append(dst, signature...)
	}
	return appendRecordArrayHeader(dst, recordType, 0, 0)
}

// marshalCardIWRecordsG2 marshals CardIWRecords for Gen2.
func marshalCardIWRecordsG2(records []*ddv1.VuCardIWRecordG2) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalVuCardIWRecordG2(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal CardIWRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// marshalActivityChangeInfos marshals ActivityChangeInfo records.
func marshalActivityChangeInfos(records []*ddv1.ActivityChangeInfo) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalActivityChangeInfo(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal ActivityChangeInfo %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// marshalPlaceRecordsG2V1 marshals PlaceRecords for Gen2v1.
func marshalPlaceRecordsG2V1(records []*ddv1.PlaceRecordG2) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, placeRec := range records {
		// Wrap in VuPlaceDailyWorkPeriodRecordG2 (40 bytes = 19 bytes FullCardNumberAndGeneration + 21 bytes PlaceRecordG2)
		ddRecord := &ddv1.VuPlaceDailyWorkPeriodRecordG2{}
		// Note: VU place records include a card number, but Gen2v1 proto doesn't expose it
		// Use empty/zero card number for now
		ddRecord.SetFullCardNumber(&ddv1.FullCardNumberAndGeneration{})
		ddRecord.SetPlaceRecord(placeRec)

		recordData, err := opts.MarshalVuPlaceDailyWorkPeriodRecordG2(ddRecord)
		if err != nil {
			return nil, fmt.Errorf("marshal PlaceRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// marshalGnssAccumulatedDrivingRecordsV1 marshals GnssAccumulatedDrivingRecords for Gen2v1.
func marshalGnssAccumulatedDrivingRecordsV1(records []*ddv1.VuGNSSADRecord) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalVuGNSSADRecord(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal GnssAccumulatedDrivingRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// marshalSpecificConditionRecords marshals SpecificConditionRecord records.
func marshalSpecificConditionRecords(records []*ddv1.SpecificConditionRecord) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalSpecificConditionRecord(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal SpecificConditionRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// anonymizeActivitiesGen2V1 anonymizes Gen2 V1 Activities data.
//
// Anonymization strategy:
// - Replaces timestamps with deterministic sequential values
// - Replaces card numbers and holder names with generic test data
// - Normalizes locations to generic values (Finland/Helsinki)
// - Rounds odometer values to nearest 100km
// - Clears signatures and raw_data
func (opts AnonymizeOptions) anonymizeActivitiesGen2V1(activities *vuv1.ActivitiesGen2V1) *vuv1.ActivitiesGen2V1 {
	if activities == nil {
		return nil
	}

	result := &vuv1.ActivitiesGen2V1{}

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize date_of_day - use a fixed date (2024-01-01 00:00:00 UTC)
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result.SetDateOfDay(timestamppb.New(baseTime))

	// Round odometer to nearest 100km
	originalOdometer := activities.GetOdometerMidnightKm()
	roundedOdometer := (originalOdometer / 100) * 100
	result.SetOdometerMidnightKm(roundedOdometer)

	// Anonymize card_iw_data
	anonCardIW := make([]*ddv1.VuCardIWRecordG2, len(activities.GetCardIwData()))
	for i, rec := range activities.GetCardIwData() {
		anonCardIW[i] = &ddv1.VuCardIWRecordG2{}

		// Generic test holder name
		testSurname := &ddv1.StringValue{}
		testSurname.SetValue("TEST")
		testFirstName := &ddv1.StringValue{}
		testFirstName.SetValue("DRIVER")
		testName := &ddv1.HolderName{}
		testName.SetHolderSurname(testSurname)
		testName.SetHolderFirstNames(testFirstName)
		anonCardIW[i].SetCardHolderName(testName)

		// Generic test card number - use empty FullCardNumberAndGeneration
		anonCardIW[i].SetFullCardNumber(&ddv1.FullCardNumberAndGeneration{})

		// Use fixed dates
		testDate := &ddv1.Date{}
		testDate.SetYear(2030)
		testDate.SetMonth(12)
		testDate.SetDay(31)
		anonCardIW[i].SetCardExpiryDate(testDate)
		anonCardIW[i].SetCardInsertionTime(timestamppb.New(baseTime.Add(time.Duration(i*2) * time.Hour)))
		anonCardIW[i].SetCardWithdrawalTime(timestamppb.New(baseTime.Add(time.Duration(i*2+1) * time.Hour)))

		// Round odometer values
		anonCardIW[i].SetOdometerAtInsertionKm((rec.GetOdometerAtInsertionKm() / 100) * 100)
		anonCardIW[i].SetOdometerAtWithdrawalKm((rec.GetOdometerAtWithdrawalKm() / 100) * 100)

		// Preserve slot and manual input flag
		anonCardIW[i].SetCardSlotNumber(rec.GetCardSlotNumber())
		anonCardIW[i].SetManualInputFlag(rec.GetManualInputFlag())

		// Anonymize previous vehicle info
		if prevVehicle := rec.GetPreviousVehicleInfo(); prevVehicle != nil {
			anonCardIW[i].SetPreviousVehicleInfo(ddOpts.AnonymizePreviousVehicleInfoG2(prevVehicle))
		}
	}
	result.SetCardIwData(anonCardIW)

	// Anonymize activity_changes
	anonActivityChanges := make([]*ddv1.ActivityChangeInfo, len(activities.GetActivityChanges()))
	for i, ac := range activities.GetActivityChanges() {
		anonActivityChanges[i] = ddOpts.AnonymizeActivityChangeInfo(ac, i)
	}
	result.SetActivityChanges(anonActivityChanges)

	// Anonymize places
	anonPlaces := make([]*ddv1.PlaceRecordG2, len(activities.GetPlaces()))
	for i, place := range activities.GetPlaces() {
		anonPlaces[i] = ddOpts.AnonymizePlaceRecordG2(place)
	}
	result.SetPlaces(anonPlaces)

	// Anonymize gnss_accumulated_driving
	anonGnss := make([]*ddv1.VuGNSSADRecord, len(activities.GetGnssAccumulatedDriving()))
	for i, gnss := range activities.GetGnssAccumulatedDriving() {
		anonGnss[i] = &ddv1.VuGNSSADRecord{}
		anonGnss[i].SetTimeStamp(timestamppb.New(baseTime.Add(time.Duration(i*3) * time.Hour)))
		anonGnss[i].SetCardNumberDriverSlot(&ddv1.FullCardNumberAndGeneration{})
		anonGnss[i].SetCardNumberCodriverSlot(&ddv1.FullCardNumberAndGeneration{})
		// Create anonymized GNSS place record
		gnssPlace := &ddv1.GNSSPlaceRecord{}
		gnssPlace.SetTimestamp(timestamppb.New(baseTime.Add(time.Duration(i*3) * time.Hour)))
		gnssPlace.SetGnssAccuracy(gnss.GetGnssPlaceRecord().GetGnssAccuracy())
		// Use generic Finland/Helsinki coordinates
		testCoordsV1 := &ddv1.GeoCoordinates{}
		testCoordsV1.SetLatitude(60170000)  // 60.17°N (Helsinki)
		testCoordsV1.SetLongitude(24940000) // 24.94°E (Helsinki)
		gnssPlace.SetGeoCoordinates(testCoordsV1)
		anonGnss[i].SetGnssPlaceRecord(gnssPlace)
		anonGnss[i].SetVehicleOdometerKm((gnss.GetVehicleOdometerKm() / 100) * 100)
	}
	result.SetGnssAccumulatedDriving(anonGnss)

	// Preserve specific_conditions (no PII)
	result.SetSpecificConditions(activities.GetSpecificConditions())

	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})
	result.ClearRawData()

	return result
}
