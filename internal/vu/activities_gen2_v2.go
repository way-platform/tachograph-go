package vu

import (
	"fmt"
	"time"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalActivitiesGen2V2 parses Gen2 V2 Activities data from the complete transfer value.
//
// Gen2 V2 extends Gen2 V1 with border crossing and load/unload records:
//
// ASN.1 Definition:
//
//	VuActivitiesSecondGenV2 ::= SEQUENCE {
//	    timeRealRecordArray                   TimeRealRecordArray,
//	    odometerValueMidnightRecordArray      OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                   VuCardIWRecordArray,
//	    vuActivityDailyRecordArray            VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray     VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                   VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray        VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray           VuBorderCrossingRecordArray,
//	    vuLoadUnloadRecordArray               VuLoadUnloadRecordArray,
//	    signatureRecordArray                  SignatureRecordArray
//	}
//
// Each RecordArray has a 5-byte header:
//
//	recordType (1 byte) + recordSize (2 bytes, big-endian) + noOfRecords (2 bytes, big-endian)
func unmarshalActivitiesGen2V2(value []byte) (*vuv1.ActivitiesGen2V2, error) {
	// Split transfer value into data and signature
	// Gen2 uses variable-length ECDSA signatures stored as SignatureRecordArray
	// We use the sizeOf function to determine where to split
	totalSize, signatureSize, err := sizeOfActivitiesGen2V2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	activities := &vuv1.ActivitiesGen2V2{}
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

	// VuCardIWRecordArray (Gen2 - 132 bytes per record, same as V1)
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

	// VuPlaceDailyWorkPeriodRecordArray (Gen2v2 - 41 bytes per record: 19 FullCardNumber + 22 PlaceAuthRecord)
	// PlaceAuthRecord[0:21] == PlaceRecordG2 (same layout; auth byte at index 21 is dropped).
	vuPlaceRecords, bytesRead, err := parseVuPlaceDailyWorkPeriodRecordArrayG2V2(data, offset)
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

	// VuGNSSADRecordArray (Gen2v2 - 59 bytes per record with authentication)
	gnssADRecords, bytesRead, err := parseVuGNSSADRecordArrayG2(data, offset)
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

	// VuBorderCrossingRecordArray (Gen2v2 - 57 bytes per record)
	borderCrossings, bytesRead, err := parseVuBorderCrossingRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuBorderCrossingRecordArray: %w", err)
	}
	activities.SetBorderCrossings(borderCrossings)
	offset += bytesRead

	// VuLoadUnloadRecordArray (Gen2v2 - 60 bytes per record)
	loadUnloadRecs, bytesRead, err := parseVuLoadUnloadRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuLoadUnloadRecordArray: %w", err)
	}
	activities.SetLoadUnloadOperations(loadUnloadRecs)
	offset += bytesRead

	// Store signature (extracted at the beginning)
	activities.SetSignature(signature)

	// Verify we consumed exactly the right amount of data
	if offset != len(data) {
		return nil, fmt.Errorf("Activities Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return activities, nil
}

// MarshalActivitiesGen2V2 marshals Gen2 V2 Activities data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and matches the structure, it uses it as the output.
func (opts MarshalOptions) MarshalActivitiesGen2V2(activities *vuv1.ActivitiesGen2V2) ([]byte, error) {
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

	// VuCardIWRecordArray (Gen2 - 131 bytes per record)
	cardIWData, err := marshalCardIWRecordsG2V2(activities.GetCardIwData())
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

	// VuPlaceDailyWorkPeriodRecordArray (Gen2v2 - 40 bytes per record)
	placeData, err := marshalPlaceRecordsG2V2(activities.GetPlaces())
	if err != nil {
		return nil, fmt.Errorf("marshal VuPlaceDailyWorkPeriodRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x05, 40, uint16(len(activities.GetPlaces())))
	result = append(result, placeData...)

	// VuGNSSADRecordArray (Gen2v2 - 57 bytes per record with authentication)
	gnssData, err := marshalGnssAccumulatedDrivingRecordsV2(activities.GetGnssAccumulatedDriving())
	if err != nil {
		return nil, fmt.Errorf("marshal VuGNSSADRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x06, 57, uint16(len(activities.GetGnssAccumulatedDriving())))
	result = append(result, gnssData...)

	// VuSpecificConditionRecordArray (5 bytes per record)
	specificCondData, err := marshalSpecificConditionRecords(activities.GetSpecificConditions())
	if err != nil {
		return nil, fmt.Errorf("marshal VuSpecificConditionRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x07, 5, uint16(len(activities.GetSpecificConditions())))
	result = append(result, specificCondData...)

	// VuBorderCrossingRecordArray (Gen2v2 - 55 bytes per record)
	borderCrossingData, err := marshalBorderCrossingRecords(activities.GetBorderCrossings())
	if err != nil {
		return nil, fmt.Errorf("marshal VuBorderCrossingRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x08, 55, uint16(len(activities.GetBorderCrossings())))
	result = append(result, borderCrossingData...)

	// VuLoadUnloadRecordArray (Gen2v2 - 58 bytes per record)
	loadUnloadData, err := marshalLoadUnloadRecords(activities.GetLoadUnloadOperations())
	if err != nil {
		return nil, fmt.Errorf("marshal VuLoadUnloadRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x09, 58, uint16(len(activities.GetLoadUnloadOperations())))
	result = append(result, loadUnloadData...)

	// Append signature at the end (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result = appendSignatureRecordArray(result, activities.GetSignature())

	return result, nil
}

// Helper functions for parsing Gen2 V2 RecordArrays

// parseVuPlaceDailyWorkPeriodRecordArrayG2V2 parses a VuPlaceDailyWorkPeriodRecordArray for Gen2v2.
//
// Gen2v2 place records are 41 bytes: 19 (FullCardNumberAndGeneration) + 22 (PlaceAuthRecord).
// PlaceAuthRecord[0:21] has the same binary layout as PlaceRecordG2[0:21]; the final byte
// of PlaceAuthRecord is an authentication status tag that we drop when converting to PlaceRecordG2.
func parseVuPlaceDailyWorkPeriodRecordArrayG2V2(data []byte, offset int) ([]*ddv1.VuPlaceDailyWorkPeriodRecordG2, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	// Accept both 40 bytes (Gen2v1 layout, from our own marshal) and 41 bytes
	// (Gen2v2 device format: 19 FullCardNumber + 22 PlaceAuthRecord). The
	// PlaceAuthRecord auth-status byte at index 40 is dropped in both cases.
	const (
		v1RecordSize = 40
		v2RecordSize = 41
	)
	if recordSize != v1RecordSize && recordSize != v2RecordSize {
		return nil, 0, fmt.Errorf("expected VuPlaceDailyWorkPeriodRecord size %d or %d, got %d", v1RecordSize, v2RecordSize, recordSize)
	}

	opts := dd.UnmarshalOptions{PreserveRawData: true}

	records := make([]*ddv1.VuPlaceDailyWorkPeriodRecordG2, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuPlaceDailyWorkPeriodRecord %d", i)
		}
		// Parse exactly 40 bytes; the auth-status byte (when recordSize=41) is dropped
		record, err := opts.UnmarshalVuPlaceDailyWorkPeriodRecordG2(data[recStart : recStart+v1RecordSize])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuPlaceDailyWorkPeriodRecord %d: %w", i, err)
		}
		records = append(records, record)
		recStart = recEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuGNSSADRecordArrayG2 parses a VuGNSSADRecordArray (Gen2v2 - 57 bytes per record with authentication).
func parseVuGNSSADRecordArrayG2(data []byte, offset int) ([]*ddv1.VuGNSSADRecordG2, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 57 // Gen2v2
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuGNSSADRecordG2 size %d, got %d", expectedRecordSize, recordSize)
	}

	opts := dd.UnmarshalOptions{PreserveRawData: true}

	records := make([]*ddv1.VuGNSSADRecordG2, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuGNSSADRecordG2 %d", i)
		}

		record, err := opts.UnmarshalVuGNSSADRecordG2(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuGNSSADRecordG2 %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuBorderCrossingRecordArray parses a VuBorderCrossingRecordArray (Gen2v2 - 57 bytes per record).
func parseVuBorderCrossingRecordArray(data []byte, offset int) ([]*ddv1.VuBorderCrossingRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 55
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuBorderCrossingRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	opts := dd.UnmarshalOptions{PreserveRawData: true}

	records := make([]*ddv1.VuBorderCrossingRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuBorderCrossingRecord %d", i)
		}

		record, err := opts.UnmarshalVuBorderCrossingRecord(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuBorderCrossingRecord %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuLoadUnloadRecordArray parses a VuLoadUnloadRecordArray (Gen2v2 - 60 bytes per record).
func parseVuLoadUnloadRecordArray(data []byte, offset int) ([]*ddv1.VuLoadUnloadRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 58
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuLoadUnloadRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	opts := dd.UnmarshalOptions{PreserveRawData: true}

	records := make([]*ddv1.VuLoadUnloadRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuLoadUnloadRecord %d", i)
		}

		record, err := opts.UnmarshalVuLoadUnloadRecord(data[recordStart:recordEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("unmarshal VuLoadUnloadRecord %d: %w", i, err)
		}

		records = append(records, record)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// Helper functions for marshalling Gen2 V2 RecordArrays

// marshalCardIWRecordsG2V2 marshals CardIWRecords for Gen2v2 (same as V1).
func marshalCardIWRecordsG2V2(records []*ddv1.VuCardIWRecordG2) ([]byte, error) {
	// Gen2v2 uses same format as V1
	return marshalCardIWRecordsG2(records)
}

// marshalPlaceRecordsG2V2 marshals PlaceRecords for Gen2v2 (same format as V1).
func marshalPlaceRecordsG2V2(records []*ddv1.PlaceRecordG2) ([]byte, error) {
	// Gen2v2 uses same format as V1
	return marshalPlaceRecordsG2V1(records)
}

// marshalGnssAccumulatedDrivingRecordsV2 marshals GnssAccumulatedDrivingRecords for Gen2v2.
func marshalGnssAccumulatedDrivingRecordsV2(records []*ddv1.VuGNSSADRecordG2) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalVuGNSSADRecordG2(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal GnssAccumulatedDrivingRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// marshalBorderCrossingRecords marshals BorderCrossingRecords for Gen2v2.
func marshalBorderCrossingRecords(records []*ddv1.VuBorderCrossingRecord) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalVuBorderCrossingRecord(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal BorderCrossingRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// marshalLoadUnloadRecords marshals LoadUnloadRecords for Gen2v2.
func marshalLoadUnloadRecords(records []*ddv1.VuLoadUnloadRecord) ([]byte, error) {
	var result []byte
	var opts dd.MarshalOptions

	for i, rec := range records {
		recordData, err := opts.MarshalVuLoadUnloadRecord(rec)
		if err != nil {
			return nil, fmt.Errorf("marshal LoadUnloadRecord %d: %w", i, err)
		}
		result = append(result, recordData...)
	}
	return result, nil
}

// anonymizeActivitiesGen2V2 anonymizes Gen2 V2 Activities data.
//
// Anonymization strategy (same as V1 plus border crossings and load/unload):
// - Replaces timestamps with deterministic sequential values
// - Replaces card numbers and holder names with generic test data
// - Normalizes locations to generic values (Finland/Helsinki)
// - Rounds odometer values to nearest 100km
// - Clears signatures and raw_data
func (opts AnonymizeOptions) anonymizeActivitiesGen2V2(activities *vuv1.ActivitiesGen2V2) *vuv1.ActivitiesGen2V2 {
	if activities == nil {
		return nil
	}

	result := &vuv1.ActivitiesGen2V2{}

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

	// Anonymize card_iw_data (same as V1)
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

	// Anonymize places (same as V1)
	anonPlaces := make([]*ddv1.PlaceRecordG2, len(activities.GetPlaces()))
	for i, place := range activities.GetPlaces() {
		anonPlaces[i] = ddOpts.AnonymizePlaceRecordG2(place)
	}
	result.SetPlaces(anonPlaces)

	// Anonymize gnss_accumulated_driving
	anonGnss := make([]*ddv1.VuGNSSADRecordG2, len(activities.GetGnssAccumulatedDriving()))
	for i, gnss := range activities.GetGnssAccumulatedDriving() {
		anonGnss[i] = &ddv1.VuGNSSADRecordG2{}
		anonGnss[i].SetTimeStamp(timestamppb.New(baseTime.Add(time.Duration(i*3) * time.Hour)))
		anonGnss[i].SetCardNumberDriverSlot(&ddv1.FullCardNumberAndGeneration{})
		anonGnss[i].SetCardNumberCodriverSlot(&ddv1.FullCardNumberAndGeneration{})

		// Create anonymized GNSS place auth record
		gnssAuthRec := &ddv1.GNSSPlaceAuthRecord{}
		gnssAuthRec.SetTimestamp(timestamppb.New(baseTime.Add(time.Duration(i*3) * time.Hour)))
		gnssAuthRec.SetGnssAccuracy(gnss.GetGnssPlaceAuthRecord().GetGnssAccuracy())
		testCoords := &ddv1.GeoCoordinates{}
		testCoords.SetLatitude(60170000)  // 60.17°N (Helsinki)
		testCoords.SetLongitude(24940000) // 24.94°E (Helsinki)
		gnssAuthRec.SetGeoCoordinates(testCoords)
		gnssAuthRec.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus_AUTHENTICATED)
		anonGnss[i].SetGnssPlaceAuthRecord(gnssAuthRec)
		anonGnss[i].SetVehicleOdometerKm((gnss.GetVehicleOdometerKm() / 100) * 100)
	}
	result.SetGnssAccumulatedDriving(anonGnss)

	// Preserve specific_conditions (no PII)
	result.SetSpecificConditions(activities.GetSpecificConditions())

	// Anonymize border_crossings (Gen2v2 specific)
	anonBorderCrossings := make([]*ddv1.VuBorderCrossingRecord, len(activities.GetBorderCrossings()))
	for i, bc := range activities.GetBorderCrossings() {
		anonBorderCrossings[i] = &ddv1.VuBorderCrossingRecord{}
		anonBorderCrossings[i].SetCardNumberDriverSlot(&ddv1.FullCardNumberAndGeneration{})
		anonBorderCrossings[i].SetCardNumberCodriverSlot(&ddv1.FullCardNumberAndGeneration{})
		anonBorderCrossings[i].SetCountryLeft(ddv1.NationNumeric_FINLAND)
		anonBorderCrossings[i].SetCountryEntered(ddv1.NationNumeric_SWEDEN)
		anonBorderCrossings[i].SetVehicleOdometerKm((bc.GetVehicleOdometerKm() / 100) * 100)

		// Anonymize GNSS auth record
		anonGnssAuth := &ddv1.GNSSPlaceAuthRecord{}
		anonGnssAuth.SetTimestamp(timestamppb.New(baseTime.Add(time.Duration(i*4) * time.Hour)))
		anonGnssAuth.SetGnssAccuracy(10)
		testCoordsBc := &ddv1.GeoCoordinates{}
		testCoordsBc.SetLatitude(60170000)
		testCoordsBc.SetLongitude(24940000)
		anonGnssAuth.SetGeoCoordinates(testCoordsBc)
		anonGnssAuth.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus_AUTHENTICATED)
		anonBorderCrossings[i].SetGnssPlaceAuthRecord(anonGnssAuth)
	}
	result.SetBorderCrossings(anonBorderCrossings)

	// Anonymize load_unload_operations (Gen2v2 specific)
	anonLoadUnload := make([]*ddv1.VuLoadUnloadRecord, len(activities.GetLoadUnloadOperations()))
	for i, lu := range activities.GetLoadUnloadOperations() {
		anonLoadUnload[i] = &ddv1.VuLoadUnloadRecord{}
		anonLoadUnload[i].SetTimeStamp(timestamppb.New(baseTime.Add(time.Duration(i*5) * time.Hour)))
		anonLoadUnload[i].SetOperationType(lu.GetOperationType())
		anonLoadUnload[i].SetCardNumberDriverSlot(&ddv1.FullCardNumberAndGeneration{})
		anonLoadUnload[i].SetCardNumberCodriverSlot(&ddv1.FullCardNumberAndGeneration{})
		anonLoadUnload[i].SetVehicleOdometerKm((lu.GetVehicleOdometerKm() / 100) * 100)

		// Anonymize GNSS auth record
		anonGnssAuthLu := &ddv1.GNSSPlaceAuthRecord{}
		anonGnssAuthLu.SetTimestamp(timestamppb.New(baseTime.Add(time.Duration(i*5) * time.Hour)))
		anonGnssAuthLu.SetGnssAccuracy(10)
		testCoordsLu := &ddv1.GeoCoordinates{}
		testCoordsLu.SetLatitude(60170000)
		testCoordsLu.SetLongitude(24940000)
		anonGnssAuthLu.SetGeoCoordinates(testCoordsLu)
		anonGnssAuthLu.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus_AUTHENTICATED)
		anonLoadUnload[i].SetGnssPlaceAuthRecord(anonGnssAuthLu)
	}
	result.SetLoadUnloadOperations(anonLoadUnload)

	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})
	result.ClearRawData()

	return result
}
