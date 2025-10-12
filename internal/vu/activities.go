package vu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/way-platform/tachograph-go/internal/dd"

	"google.golang.org/protobuf/types/known/timestamppb"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfActivities dispatches to generation-specific size calculation.
func sizeOfActivities(data []byte, transferType vuv1.TransferType) (int, error) {
	switch transferType {
	case vuv1.TransferType_ACTIVITIES_GEN1:
		return sizeOfActivitiesGen1(data)
	case vuv1.TransferType_ACTIVITIES_GEN2_V1:
		return sizeOfActivitiesGen2V1(data)
	case vuv1.TransferType_ACTIVITIES_GEN2_V2:
		return sizeOfActivitiesGen2V2(data)
	default:
		return 0, fmt.Errorf("unsupported transfer type for Activities: %v", transferType)
	}
}

// sizeOfActivitiesGen1 calculates total size for Gen1 Activities including signature.
//
// Activities Gen1 structure (from Appendix 7, Section 2.2.6.3):
// - TimeReal: 4 bytes (date of day downloaded)
// - OdometerValueMidnight: 3 bytes (OdometerShort)
// - VuCardIWData: 2 bytes (noOfIWRecords) + (noOfIWRecords * 129 bytes)
//   - VuCardIWRecordFirstGen: 129 bytes (72+18+4+4+3+1+4+3+19+1)
//
// - VuActivityDailyData: 2 bytes (noOfActivityChanges) + (noOfActivityChanges * 2 bytes)
//   - ActivityChangeInfo: 2 bytes
//
// - VuPlaceDailyWorkPeriodData: 1 byte (noOfPlaceRecords) + (noOfPlaceRecords * 28 bytes)
//   - VuPlaceDailyWorkPeriodRecordFirstGen: 28 bytes (18 + 10)
//
// - VuSpecificConditionData: 2 bytes (noOfSpecificConditionRecords) + (noOfSpecificConditionRecords * 5 bytes)
//   - SpecificConditionRecord: 5 bytes (4 + 1)
//
// - Signature: 128 bytes (RSA)
func sizeOfActivitiesGen1(data []byte) (int, error) {
	offset := 0

	// Fixed-size header sections (7 bytes total)
	offset += 4 // TimeReal (date of day downloaded)
	offset += 3 // OdometerValueMidnight

	// VuCardIWData: 2 bytes count + variable records
	if len(data[offset:]) < 2 {
		return 0, fmt.Errorf("insufficient data for noOfIWRecords")
	}
	noOfIWRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each VuCardIWRecordFirstGen: 129 bytes
	// (cardHolderName 72, fullCardNumber 18, cardExpiryDate 4, cardInsertionTime 4,
	//  vehicleOdometerValueAtInsertion 3, cardSlotNumber 1, cardWithdrawalTime 4,
	//  vehicleOdometerValueAtWithdrawal 3, previousVehicleInfo 19, manualInputFlag 1)
	const vuCardIWRecordSize = 129
	offset += int(noOfIWRecords) * vuCardIWRecordSize

	// VuActivityDailyData: 2 bytes count + variable activity changes
	if len(data[offset:]) < 2 {
		return 0, fmt.Errorf("insufficient data for noOfActivityChanges")
	}
	noOfActivityChanges := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each ActivityChangeInfo: 2 bytes
	const activityChangeInfoSize = 2
	offset += int(noOfActivityChanges) * activityChangeInfoSize

	// VuPlaceDailyWorkPeriodData: 1 byte count + variable place records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfPlaceRecords")
	}
	noOfPlaceRecords := data[offset]
	offset += 1

	// Each VuPlaceDailyWorkPeriodRecordFirstGen: 28 bytes (18 FullCardNumber + 10 PlaceRecordFirstGen)
	const vuPlaceDailyWorkPeriodRecordSize = 28
	offset += int(noOfPlaceRecords) * vuPlaceDailyWorkPeriodRecordSize

	// VuSpecificConditionData: 2 bytes count + variable condition records
	if len(data[offset:]) < 2 {
		return 0, fmt.Errorf("insufficient data for noOfSpecificConditionRecords")
	}
	noOfSpecificConditionRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each SpecificConditionRecord: 5 bytes (4 TimeReal + 1 SpecificConditionType)
	const specificConditionRecordSize = 5
	offset += int(noOfSpecificConditionRecords) * specificConditionRecordSize

	// Signature: 128 bytes for Gen1 RSA
	offset += 128

	return offset, nil
}

// sizeOfActivitiesGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
func sizeOfActivitiesGen2V1(data []byte) (int, error) {
	offset := 0

	// DateOfDayDownloadedRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("DateOfDayDownloadedRecordArray: %w", err)
	}
	offset += size

	// OdometerValueMidnightRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("OdometerValueMidnightRecordArray: %w", err)
	}
	offset += size

	// VuCardIWRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuCardIWRecordArray: %w", err)
	}
	offset += size

	// VuActivityDailyRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuActivityDailyRecordArray: %w", err)
	}
	offset += size

	// VuPlaceDailyWorkPeriodRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuPlaceDailyWorkPeriodRecordArray: %w", err)
	}
	offset += size

	// VuSpecificConditionRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuSpecificConditionRecordArray: %w", err)
	}
	offset += size

	// VuGNSSADRecordArray (Gen2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuGNSSADRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// sizeOfActivitiesGen2V2 calculates size by parsing all Gen2 V2 RecordArrays.
// Must handle VuBorderCrossingRecordArray and VuLoadUnloadRecordArray.
func sizeOfActivitiesGen2V2(data []byte) (int, error) {
	offset := 0

	// DateOfDayDownloadedRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("DateOfDayDownloadedRecordArray: %w", err)
	}
	offset += size

	// OdometerValueMidnightRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("OdometerValueMidnightRecordArray: %w", err)
	}
	offset += size

	// VuCardIWRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuCardIWRecordArray: %w", err)
	}
	offset += size

	// VuActivityDailyRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuActivityDailyRecordArray: %w", err)
	}
	offset += size

	// VuPlaceDailyWorkPeriodRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuPlaceDailyWorkPeriodRecordArray: %w", err)
	}
	offset += size

	// VuSpecificConditionRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuSpecificConditionRecordArray: %w", err)
	}
	offset += size

	// VuGNSSADRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuGNSSADRecordArray: %w", err)
	}
	offset += size

	// VuBorderCrossingRecordArray (Gen2 V2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuBorderCrossingRecordArray: %w", err)
	}
	offset += size

	// VuLoadUnloadRecordArray (Gen2 V2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuLoadUnloadRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// UnmarshalVuActivities unmarshals VU activities data from a VU transfer.
//
// The data type `VuActivities` is specified in the Data Dictionary, Section 2.2.6.2.
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    dateOfDay                        TimeReal,
//	    odometerValueMidnight            OdometerValueMidnight,
//	    vuCardIWData                     VuCardIWData,
//	    vuActivityDailyData              VuActivityDailyData,
//	    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
//	    vuSpecificConditionData          VuSpecificConditionData,
//	    signature                        SignatureFirstGen
//	}
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
//	    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                      VuCardIWRecordArray,
//	    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                      VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
//	    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
//	    signatureRecordArray                     SignatureRecordArray
//	}
func UnmarshalVuActivities(data []byte, offset int, target *vuv1.Activities, generation int) (int, error) {
	switch generation {
	case 1:
		return unmarshalVuActivitiesGen1(data, offset, target)
	case 2:
		return unmarshalVuActivitiesGen2(data, offset, target)
	default:
		return 0, fmt.Errorf("unsupported generation: %d", generation)
	}
}

// unmarshalVuActivitiesGen1 unmarshals Generation 1 VU activities
func unmarshalVuActivitiesGen1(data []byte, offset int, target *vuv1.Activities) (int, error) {
	var opts dd.UnmarshalOptions
	startOffset := offset
	target.SetGeneration(ddv1.Generation_GENERATION_1)

	// Read TimeReal (4 bytes) - this is the date of the day
	timeReal, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read date of day: %w", err)
	}
	target.SetDateOfDay(timestamppb.New(time.Unix(timeReal, 0)))

	// Read OdometerValueMidnight (3 bytes)
	odometerValue, err := opts.UnmarshalOdometer(data[offset : offset+3])
	if err != nil {
		return 0, fmt.Errorf("failed to read odometer value midnight: %w", err)
	}
	offset += 3
	target.SetOdometerMidnightKm(int32(odometerValue))

	// Parse VuCardIWData
	cardIWData, offset, err := parseVuCardIWData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW data: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyData
	activityChanges, offset, err := parseVuActivityDailyData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily data: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodData
	places, offset, err := parseVuPlaceDailyWorkPeriodData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period data: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuSpecificConditionData
	specificConditions, offset, err := parseVuSpecificConditionData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition data: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Read signature (128 bytes for Gen1)
	signatureBytes, offset, err := readBytesFromBytes(data, offset, 128)
	if err != nil {
		return 0, fmt.Errorf("failed to read signature: %w", err)
	}
	target.SetSignatureGen1(signatureBytes)

	return offset - startOffset, nil
}

// unmarshalVuActivitiesGen2 unmarshals Generation 2 VU activities
func unmarshalVuActivitiesGen2(data []byte, offset int, target *vuv1.Activities) (int, error) {
	startOffset := offset
	target.SetGeneration(ddv1.Generation_GENERATION_2)

	// Gen2 format uses record arrays, each with a header
	// Parse DateOfDayDownloadedRecordArray
	dates, offset, err := parseDateOfDayDownloadedRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse date of day downloaded record array: %w", err)
	}
	if len(dates) > 0 {
		target.SetDateOfDay(dates[0]) // Use first date
	}

	// Parse OdometerValueMidnightRecordArray
	odometerValues, offset, err := parseOdometerValueMidnightRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse odometer value midnight record array: %w", err)
	}
	if len(odometerValues) > 0 {
		target.SetOdometerMidnightKm(odometerValues[0])
	}

	// Parse VuCardIWRecordArray
	cardIWData, offset, err := parseVuCardIWRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW record array: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyRecordArray
	activityChanges, offset, err := parseVuActivityDailyRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily record array: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodRecordArray
	places, offset, err := parseVuPlaceDailyWorkPeriodRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period record array: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuGNSSADRecordArray (Gen2+)
	gnssRecords, offset, err := parseVuGNSSADRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse GNSS AD record array: %w", err)
	}
	target.SetGnssAccumulatedDriving(gnssRecords)

	// Parse VuSpecificConditionRecordArray
	specificConditions, offset, err := parseVuSpecificConditionRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition record array: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Try to parse Gen2v2 specific arrays if there's more data
	if offset+10 <= len(data) { // Need some minimum data for arrays
		// Parse VuBorderCrossingRecordArray (Gen2v2+)
		borderCrossings, newOffset, err := parseVuBorderCrossingRecordArray(data, offset)
		if err == nil {
			target.SetBorderCrossings(borderCrossings)
			target.SetVersion(ddv1.Version_VERSION_2)
			offset = newOffset
		}

		// Parse VuLoadUnloadRecordArray (Gen2v2+)
		if offset+5 <= len(data) {
			loadUnloadRecords, newOffset, err := parseVuLoadUnloadRecordArray(data, offset)
			if err == nil {
				target.SetLoadUnloadOperations(loadUnloadRecords)
				offset = newOffset
			}
		}
	}

	// Parse SignatureRecordArray
	signatureBytes, offset, err := parseSignatureRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signature record array: %w", err)
	}
	target.SetSignatureGen2(signatureBytes)

	return offset - startOffset, nil
}

// Helper functions for parsing different record types
// These are simplified implementations - in a full implementation,
// each would need to properly handle the record array format

// splitVuCardIWRecord splits data into 126-byte VuCardIWRecord records
func splitVuCardIWRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const cardIWRecordSize = 126

	if len(data) < cardIWRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return cardIWRecordSize, data[:cardIWRecordSize], nil
}

func parseVuCardIWData(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWData ::= SEQUENCE {
	//     noOfIWRecords INTEGER(0..255),
	//     vuCardIWRecords SET SIZE(noOfIWRecords) OF VuCardIWRecord -- 126 bytes each
	// }

	// Read number of records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of IW records: %w", err)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuCardIWRecord)

	var records []*vuv1.Activities_CardIWRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuCardIWRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse IW record %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 126

	return records, offset, nil
}

// unmarshalVuCardIWRecord parses a single VuCardIWRecord from a byte slice
func unmarshalVuCardIWRecord(data []byte) (*vuv1.Activities_CardIWRecord, error) {
	// VuCardIWRecord ::= SEQUENCE {
	//     cardHolderName HolderName,                    -- 72 bytes
	//     fullCardNumber FullCardNumber,                -- 19 bytes
	//     cardExpiryDate Datef,                         -- 4 bytes
	//     cardInsertionTime TimeReal,                   -- 4 bytes
	//     vehicleOdometerValueAtInsertion OdometerShort, -- 3 bytes
	//     cardSlotNumber CardSlotNumber,                -- 1 byte
	//     cardWithdrawalTime TimeReal,                  -- 4 bytes
	//     vehicleOdometerValueAtWithdrawal OdometerShort, -- 3 bytes
	//     previousVehicleInfo PreviousVehicleInfo,      -- 19 bytes
	//     manualInputFlag ManualInputFlag               -- 1 byte
	// }
	// Total: 130 bytes

	if len(data) < 130 {
		return nil, fmt.Errorf("insufficient data for card IW record: got %d, need 130", len(data))
	}

	var opts dd.UnmarshalOptions
	record := &vuv1.Activities_CardIWRecord{}

	// Parse cardHolderName (HolderName - 72 bytes)
	holderNameData := data[0:72]
	holderName, err := opts.UnmarshalHolderName(holderNameData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// Parse fullCardNumber (FullCardNumber - 19 bytes)
	fullCardNumberData := data[72:91]
	fullCardNumber, err := opts.UnmarshalFullCardNumber(fullCardNumberData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal full card number: %w", err)
	}
	// Create FullCardNumberAndGeneration wrapper
	fullCardNumberAndGeneration := &ddv1.FullCardNumberAndGeneration{}
	fullCardNumberAndGeneration.SetFullCardNumber(fullCardNumber)
	fullCardNumberAndGeneration.SetGeneration(ddv1.Generation_GENERATION_1) // Default to Gen1
	record.SetFullCardNumberAndGeneration(fullCardNumberAndGeneration)

	// Parse cardExpiryDate (Datef - 4 bytes)
	cardExpiryDate, err := opts.UnmarshalDate(data[91:95])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card expiry date: %w", err)
	}
	record.SetCardExpiryDate(cardExpiryDate)

	// Parse cardInsertionTime (TimeReal - 4 bytes)
	insertionTime, _, err := readVuTimeRealFromBytes(data, 95)
	if err != nil {
		return nil, fmt.Errorf("failed to read card insertion time: %w", err)
	}
	record.SetCardInsertionTime(timestamppb.New(time.Unix(insertionTime, 0)))

	// Parse vehicleOdometerValueAtInsertion (OdometerShort - 3 bytes)
	odometerAtInsertion, err := opts.UnmarshalOdometer(data[99 : 99+3])
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer at insertion: %w", err)
	}
	record.SetOdometerAtInsertionKm(int32(odometerAtInsertion))

	// Parse cardSlotNumber (CardSlotNumber - 1 byte)
	slotNumber := data[102]
	record.SetCardSlotNumber(ddv1.CardSlotNumber(slotNumber))

	// Parse cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime, _, err := readVuTimeRealFromBytes(data, 103)
	if err != nil {
		return nil, fmt.Errorf("failed to read card withdrawal time: %w", err)
	}
	record.SetCardWithdrawalTime(timestamppb.New(time.Unix(withdrawalTime, 0)))

	// Parse vehicleOdometerValueAtWithdrawal (OdometerShort - 3 bytes)
	odometerAtWithdrawal, err := opts.UnmarshalOdometer(data[107 : 107+3])
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer at withdrawal: %w", err)
	}
	record.SetOdometerAtWithdrawalKm(int32(odometerAtWithdrawal))

	// Parse previousVehicleInfo (PreviousVehicleInfo - 19 bytes: 15 vehicle reg + 4 cardWithdrawalTime)
	previousVehicleData := data[110:129]
	opts.Generation = ddv1.Generation_GENERATION_1
	previousVehicleInfo, err := opts.UnmarshalPreviousVehicleInfo(previousVehicleData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal previous vehicle info: %w", err)
	}
	record.SetPreviousVehicleInfo(previousVehicleInfo)

	// Parse manualInputFlag (ManualInputFlag - 1 byte)
	manualInputFlag := data[125]
	record.SetManualInputFlag(manualInputFlag != 0)

	return record, nil
}

func parseVuActivityDailyData(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	// VuActivityDailyData ::= SEQUENCE {
	//     noOfActivityChanges INTEGER(0..255),
	//     activityChanges SET SIZE(noOfActivityChanges) OF ActivityChangeInfo
	// }

	// Read number of activity changes (1 byte)
	noOfChanges, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of activity changes: %w", err)
	}

	var changes []*ddv1.ActivityChangeInfo

	// Parse each ActivityChangeInfo
	for i := 0; i < int(noOfChanges); i++ {
		change, newOffset, err := parseActivityChangeInfo(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse activity change %d: %w", i, err)
		}
		changes = append(changes, change)
		offset = newOffset
	}

	return changes, offset, nil
}

// parseActivityChangeInfo parses a single ActivityChangeInfo record (legacy function for Gen1)
func parseActivityChangeInfo(data []byte, offset int) (*ddv1.ActivityChangeInfo, int, error) {
	var opts dd.UnmarshalOptions
	change, err := opts.UnmarshalActivityChangeInfo(data[offset : offset+2])
	if err != nil {
		return nil, offset, err
	}
	return change, offset + 2, nil
}

func parseVuPlaceDailyWorkPeriodData(data []byte, offset int) ([]*vuv1.Activities_PlaceRecord, int, error) {
	// VuPlaceDailyWorkPeriodData ::= SEQUENCE {
	//     noOfPlaceRecords INTEGER(0..255),
	//     placeRecords SET SIZE(noOfPlaceRecords) OF PlaceRecord
	// }

	// Read number of place records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of place records: %w", err)
	}

	var records []*vuv1.Activities_PlaceRecord

	// Parse each PlaceRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuPlaceRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse place record %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// unmarshalVuPlaceRecord parses a single VuPlaceRecord from a byte slice
func unmarshalVuPlaceRecord(data []byte) (*vuv1.Activities_PlaceRecord, error) {
	// PlaceRecord ::= SEQUENCE {
	//     entryTime TimeReal,                           -- 4 bytes
	//     entryTypeDailyWorkPeriod EntryTypeDailyWorkPeriod, -- 1 byte
	//     dailyWorkPeriodCountry NationNumeric,         -- 1 byte
	//     dailyWorkPeriodRegion RegionNumeric,          -- 1 byte
	//     vehicleOdometerValue OdometerShort            -- 3 bytes
	// }

	if len(data) < 10 {
		return nil, fmt.Errorf("insufficient data for place record: got %d, need 10", len(data))
	}

	var opts dd.UnmarshalOptions
	record := &vuv1.Activities_PlaceRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, _, err := readVuTimeRealFromBytes(data, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read entry time: %w", err)
	}
	record.SetEntryTime(timestamppb.New(time.Unix(entryTime, 0)))

	// Parse entryTypeDailyWorkPeriod (1 byte)
	entryType := data[4]
	record.SetEntryType(ddv1.EntryTypeDailyWorkPeriod(entryType))

	// Parse dailyWorkPeriodCountry (NationNumeric - 1 byte)
	country := data[5]
	record.SetCountry(ddv1.NationNumeric(country))

	// Parse dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	region := data[6]
	record.SetRegion([]byte{region})

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, err := opts.UnmarshalOdometer(data[7 : 7+3])
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, nil
}

// parseVuPlaceRecord parses a single PlaceRecord (legacy function for Gen1)
func parseVuPlaceRecord(data []byte, offset int) (*vuv1.Activities_PlaceRecord, int, error) {
	record, err := unmarshalVuPlaceRecord(data[offset : offset+10])
	if err != nil {
		return nil, offset, err
	}
	return record, offset + 10, nil
}

func parseVuSpecificConditionData(data []byte, offset int) ([]*ddv1.SpecificConditionRecord, int, error) {
	// VuSpecificConditionData ::= SEQUENCE {
	//     noOfSpecificConditionRecords INTEGER(0..255),
	//     specificConditionRecords SET SIZE(noOfSpecificConditionRecords) OF SpecificConditionRecord
	// }

	// Read number of specific condition records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of specific condition records: %w", err)
	}

	var records []*ddv1.SpecificConditionRecord

	// Parse each SpecificConditionRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuSpecificConditionRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse specific condition record %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuSpecificConditionRecord parses a single SpecificConditionRecord (legacy function for Gen1)
func parseVuSpecificConditionRecord(data []byte, offset int) (*ddv1.SpecificConditionRecord, int, error) {
	var opts dd.UnmarshalOptions
	record, err := opts.UnmarshalSpecificConditionRecord(data[offset : offset+5])
	if err != nil {
		return nil, offset, err
	}
	return record, offset + 5, nil
}

// Gen2 record array parsers
func parseDateOfDayDownloadedRecordArray(data []byte, offset int) ([]*timestamppb.Timestamp, int, error) {
	// DateOfDayDownloadedRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF TimeReal -- 4 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x01 for TimeReal)
	if recordType != 0x01 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 1 (TimeReal)", recordType)
	}

	// Validate record size (should be 4 bytes for TimeReal)
	if recordSize != 4 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 4", recordSize)
	}

	var timestamps []*timestamppb.Timestamp

	// Parse each TimeReal record
	for i := 0; i < int(noOfRecords); i++ {
		if offset+4 > len(data) {
			return nil, offset, fmt.Errorf("insufficient data for TimeReal record %d: got %d, need 4", i, len(data)-offset)
		}

		timeValue, newOffset, err := readVuTimeRealFromBytes(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to read TimeReal record %d: %w", i, err)
		}

		timestamps = append(timestamps, timestamppb.New(time.Unix(timeValue, 0)))
		offset = newOffset
	}

	return timestamps, offset, nil
}

func parseOdometerValueMidnightRecordArray(data []byte, offset int) ([]int32, int, error) {
	// OdometerValueMidnightRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF OdometerShort -- 3 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x02 for OdometerShort)
	if recordType != 0x02 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 2 (OdometerShort)", recordType)
	}

	// Validate record size (should be 3 bytes for OdometerShort)
	if recordSize != 3 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 3", recordSize)
	}

	var opts dd.UnmarshalOptions
	var odometerValues []int32

	// Parse each OdometerShort record
	for i := 0; i < int(noOfRecords); i++ {
		odometerValue, err := opts.UnmarshalOdometer(data[offset : offset+3])
		if err != nil {
			return nil, offset, fmt.Errorf("failed to read OdometerShort record %d: %w", i, err)
		}

		odometerValues = append(odometerValues, int32(odometerValue))
		offset += 3
	}

	return odometerValues, offset, nil
}

func parseVuCardIWRecordArray(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuCardIWRecord -- Variable size each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes) - not used for variable-length records
	_ = binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x0D for VuCardIWRecord)
	if recordType != 0x0D {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 13 (VuCardIWRecord)", recordType)
	}

	var records []*vuv1.Activities_CardIWRecord

	// Parse each VuCardIWRecord
	for i := 0; i < int(noOfRecords); i++ {
		// For Gen2, we need to determine the record size dynamically
		// VuCardIWRecord has variable length due to HolderName and PreviousVehicleInfo
		// We'll use a conservative approach and try to parse with known minimum size
		record, newOffset, err := parseVuCardIWRecordGen2(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuCardIWRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuCardIWRecordGen2 parses a single VuCardIWRecord for Gen2
func parseVuCardIWRecordGen2(data []byte, offset int) (*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWRecord (Gen2) ::= SEQUENCE {
	//     cardHolderName HolderName,                    -- 72 bytes
	//     fullCardNumberAndGeneration FullCardNumberAndGeneration, -- 20 bytes
	//     cardExpiryDate Datef,                         -- 4 bytes
	//     cardInsertionTime TimeReal,                   -- 4 bytes
	//     vehicleOdometerValueAtInsertion OdometerShort, -- 3 bytes
	//     cardSlotNumber CardSlotNumber,                -- 1 byte
	//     cardWithdrawalTime TimeReal,                  -- 4 bytes
	//     vehicleOdometerValueAtWithdrawal OdometerShort, -- 3 bytes
	//     previousVehicleInfo PreviousVehicleInfo,      -- 20 bytes
	//     manualInputFlag ManualInputFlag               -- 1 byte
	// }
	// Total: 132 bytes

	var opts dd.UnmarshalOptions
	record := &vuv1.Activities_CardIWRecord{}

	// Parse cardHolderName (HolderName - 72 bytes)
	holderNameData, offset, err := readBytesFromBytes(data, offset, 72)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card holder name: %w", err)
	}
	holderName, err := opts.UnmarshalHolderName(holderNameData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// Parse fullCardNumberAndGeneration (FullCardNumberAndGeneration - 20 bytes)
	fullCardNumberAndGenData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read full card number and generation: %w", err)
	}
	_, err = opts.UnmarshalFullCardNumberAndGeneration(fullCardNumberAndGenData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal full card number and generation: %w", err)
	}
	// Note: The protobuf might not have this field yet, so we'll set the regular fullCardNumber
	// This is a limitation of the current schema that should be addressed
	// Create placeholder FullCardNumberAndGeneration
	placeholder := &ddv1.FullCardNumberAndGeneration{}
	placeholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	placeholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetFullCardNumberAndGeneration(placeholder)

	// Parse cardExpiryDate (Datef - 4 bytes)
	if offset+4 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for card expiry date")
	}
	cardExpiryDate, err := opts.UnmarshalDate(data[offset : offset+4])
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse card expiry date: %w", err)
	}
	offset += 4
	record.SetCardExpiryDate(cardExpiryDate)

	// Parse cardInsertionTime (TimeReal - 4 bytes)
	insertionTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card insertion time: %w", err)
	}
	record.SetCardInsertionTime(timestamppb.New(time.Unix(insertionTime, 0)))

	// Parse vehicleOdometerValueAtInsertion (OdometerShort - 3 bytes)
	odometerAtInsertion, err := opts.UnmarshalOdometer(data[offset : offset+3])
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer at insertion: %w", err)
	}
	offset += 3
	record.SetOdometerAtInsertionKm(int32(odometerAtInsertion))

	// Parse cardSlotNumber (CardSlotNumber - 1 byte)
	slotNumber, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card slot number: %w", err)
	}
	record.SetCardSlotNumber(ddv1.CardSlotNumber(slotNumber))

	// Parse cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card withdrawal time: %w", err)
	}
	record.SetCardWithdrawalTime(timestamppb.New(time.Unix(withdrawalTime, 0)))

	// Parse vehicleOdometerValueAtWithdrawal (OdometerShort - 3 bytes)
	odometerAtWithdrawal, err := opts.UnmarshalOdometer(data[offset : offset+3])
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer at withdrawal: %w", err)
	}
	offset += 3
	record.SetOdometerAtWithdrawalKm(int32(odometerAtWithdrawal))

	// Parse previousVehicleInfo (PreviousVehicleInfo - 20 bytes: 15 vehicle reg + 4 cardWithdrawalTime + 1 vuGeneration)
	previousVehicleData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read previous vehicle info: %w", err)
	}
	opts.Generation = ddv1.Generation_GENERATION_2
	previousVehicleInfo, err := opts.UnmarshalPreviousVehicleInfo(previousVehicleData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal previous vehicle info: %w", err)
	}
	record.SetPreviousVehicleInfo(previousVehicleInfo)

	// Parse manualInputFlag (ManualInputFlag - 1 byte)
	manualInputFlag, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read manual input flag: %w", err)
	}
	record.SetManualInputFlag(manualInputFlag != 0)

	return record, offset, nil
}

// splitActivityChangeInfoRecord splits data into 2-byte ActivityChangeInfo records
func splitActivityChangeInfoRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const activityChangeInfoRecordSize = 2

	if len(data) < activityChangeInfoRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return activityChangeInfoRecordSize, data[:activityChangeInfoRecordSize], nil
}

func parseVuActivityDailyRecordArray(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	// VuActivityDailyRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF ActivityChangeInfo -- 2 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x03 for ActivityChangeInfo)
	if recordType != 0x03 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 3 (ActivityChangeInfo)", recordType)
	}

	// Validate record size (should be 2 bytes for ActivityChangeInfo)
	if recordSize != 2 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 2", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitActivityChangeInfoRecord)

	var opts dd.UnmarshalOptions
	var changes []*ddv1.ActivityChangeInfo
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		change, err := opts.UnmarshalActivityChangeInfo(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse ActivityChangeInfo record %d: %w", recordCount, err)
		}
		changes = append(changes, change)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 2

	return changes, offset, nil
}

// splitVuPlaceRecord splits data into 10-byte VuPlaceRecord records
func splitVuPlaceRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const placeRecordSize = 10

	if len(data) < placeRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return placeRecordSize, data[:placeRecordSize], nil
}

func parseVuPlaceDailyWorkPeriodRecordArray(data []byte, offset int) ([]*vuv1.Activities_PlaceRecord, int, error) {
	// VuPlaceDailyWorkPeriodRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF PlaceRecord -- 10 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x04 for PlaceRecord)
	if recordType != 0x04 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 4 (PlaceRecord)", recordType)
	}

	// Validate record size (should be 10 bytes for PlaceRecord)
	if recordSize != 10 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 10", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuPlaceRecord)

	var records []*vuv1.Activities_PlaceRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuPlaceRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse PlaceRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 10

	return records, offset, nil
}

// splitVuGNSSADRecord splits data into 14-byte VuGNSSADRecord records
func splitVuGNSSADRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const gnssRecordSize = 14

	if len(data) < gnssRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return gnssRecordSize, data[:gnssRecordSize], nil
}

func parseVuGNSSADRecordArray(data []byte, offset int) ([]*vuv1.Activities_GnssRecord, int, error) {
	// VuGNSSADRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuGNSSADRecord -- 14 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for GNSS record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x16 for VuGNSSADRecord)
	if recordType != 0x16 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 22 (VuGNSSADRecord)", recordType)
	}

	// Validate record size (should be 14 bytes for VuGNSSADRecord)
	if recordSize != 14 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 14", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuGNSSADRecord)

	var records []*vuv1.Activities_GnssRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuGNSSADRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuGNSSADRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 14

	return records, offset, nil
}

// unmarshalVuGNSSADRecord parses a single VuGNSSADRecord from a byte slice
func unmarshalVuGNSSADRecord(data []byte) (*vuv1.Activities_GnssRecord, error) {
	// VuGNSSADRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                    -- 4 bytes
	//     gnssAccuracy GNSSAccuracy,            -- 1 byte
	//     geoCoordinates GeoCoordinates,        -- 8 bytes (latitude + longitude)
	//     positionAuthenticationStatus PositionAuthenticationStatus -- 1 byte
	// }

	if len(data) < 14 {
		return nil, fmt.Errorf("insufficient data for GNSS record: got %d, need 14", len(data))
	}

	record := &vuv1.Activities_GnssRecord{}

	// Parse timeStamp (TimeReal - 4 bytes)
	timestamp, _, err := readVuTimeRealFromBytes(data, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read timestamp: %w", err)
	}
	record.SetTimestamp(timestamppb.New(time.Unix(timestamp, 0)))

	// Parse gnssAccuracy (GNSSAccuracy - 1 byte)
	accuracy := data[4]
	record.SetGnssAccuracy(int32(accuracy))

	// Parse geoCoordinates (GeoCoordinates - 8 bytes: 4 bytes latitude + 4 bytes longitude)
	// Latitude (4 bytes, signed integer)
	latBytes := data[5:9]
	latitude := int32(binary.BigEndian.Uint32(latBytes))

	// Longitude (4 bytes, signed integer)
	lonBytes := data[9:13]
	longitude := int32(binary.BigEndian.Uint32(lonBytes))

	// Create GeoCoordinates
	geoCoords := &ddv1.GeoCoordinates{}
	geoCoords.SetLatitude(latitude)
	geoCoords.SetLongitude(longitude)
	record.SetGeoCoordinates(geoCoords)

	// Parse positionAuthenticationStatus (PositionAuthenticationStatus - 1 byte)
	authStatus := data[13]
	record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus(authStatus))

	return record, nil
}

// splitSpecificConditionRecord splits data into 5-byte SpecificConditionRecord records
func splitSpecificConditionRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const specificConditionRecordSize = 5

	if len(data) < specificConditionRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return specificConditionRecordSize, data[:specificConditionRecordSize], nil
}

func parseVuSpecificConditionRecordArray(data []byte, offset int) ([]*ddv1.SpecificConditionRecord, int, error) {
	// VuSpecificConditionRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF SpecificConditionRecord -- 5 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x09 for SpecificConditionRecord)
	if recordType != 0x09 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 9 (SpecificConditionRecord)", recordType)
	}

	// Validate record size (should be 5 bytes for SpecificConditionRecord)
	if recordSize != 5 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 5", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitSpecificConditionRecord)

	var opts dd.UnmarshalOptions
	var records []*ddv1.SpecificConditionRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := opts.UnmarshalSpecificConditionRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse SpecificConditionRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 5

	return records, offset, nil
}

// splitVuBorderCrossingRecord splits data into 57-byte VuBorderCrossingRecord records
func splitVuBorderCrossingRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const borderCrossingRecordSize = 57

	if len(data) < borderCrossingRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return borderCrossingRecordSize, data[:borderCrossingRecordSize], nil
}

func parseVuBorderCrossingRecordArray(data []byte, offset int) ([]*vuv1.Activities_BorderCrossingRecord, int, error) {
	// VuBorderCrossingRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuBorderCrossingRecord -- 59 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for border crossing record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x17 for VuBorderCrossingRecord)
	if recordType != 0x17 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 23 (VuBorderCrossingRecord)", recordType)
	}

	// Validate record size (should be 59 bytes for VuBorderCrossingRecord)
	if recordSize != 59 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 59", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuBorderCrossingRecord)

	var records []*vuv1.Activities_BorderCrossingRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuBorderCrossingRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuBorderCrossingRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 57

	return records, offset, nil
}

// unmarshalVuBorderCrossingRecord parses a single VuBorderCrossingRecord from a byte slice
func unmarshalVuBorderCrossingRecord(data []byte) (*vuv1.Activities_BorderCrossingRecord, error) {
	// VuBorderCrossingRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration,   -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     countryLeft NationNumeric,                                -- 1 byte
	//     countryEntered NationNumeric,                             -- 1 byte
	//     gnssPlaceAuthRecord GNSSPlaceAuthRecord,                  -- 12 bytes
	//     vehicleOdometerValue OdometerShort                        -- 3 bytes
	// }
	// Total: 57 bytes

	const vuBorderCrossingRecordSize = 57
	if len(data) < vuBorderCrossingRecordSize {
		return nil, fmt.Errorf("insufficient data for border crossing record: got %d, need %d", len(data), vuBorderCrossingRecordSize)
	}

	var opts dd.UnmarshalOptions
	record := &vuv1.Activities_BorderCrossingRecord{}

	// Parse cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	driverCardData := data[0:20]
	_, err := opts.UnmarshalFullCardNumberAndGeneration(driverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal driver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	driverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	driverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	driverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberDriverSlot(driverPlaceholder)

	// Parse cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	codriverCardData := data[20:40]
	_, err = opts.UnmarshalFullCardNumberAndGeneration(codriverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal codriver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	codriverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	codriverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	codriverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberCodriverSlot(codriverPlaceholder)

	// Parse countryLeft (NationNumeric - 1 byte)
	countryLeft := data[40]
	record.SetCountryLeft(ddv1.NationNumeric(countryLeft))

	// Parse countryEntered (NationNumeric - 1 byte)
	countryEntered := data[41]
	record.SetCountryEntered(ddv1.NationNumeric(countryEntered))

	// Parse gnssPlaceAuthRecord (GNSSPlaceAuthRecord - 12 bytes)
	placeData := data[42:54]
	placeRecord, err := unmarshalGNSSPlaceAuthRecord(placeData, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GNSS place auth record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, err := opts.UnmarshalOdometer(data[54:57])
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, nil
}

// unmarshalGNSSPlaceAuthRecord parses a GNSSPlaceAuthRecord and converts it to VU GnssRecord format.
//
// This uses the dd package implementation for parsing, then converts to the VU-specific message type.
func unmarshalGNSSPlaceAuthRecord(data []byte, opts dd.UnmarshalOptions) (*vuv1.Activities_GnssRecord, error) {
	// Use dd package to parse GNSSPlaceAuthRecord (12 bytes)
	ddRecord, err := opts.UnmarshalGNSSPlaceAuthRecord(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GNSSPlaceAuthRecord: %w", err)
	}

	// Convert dd.GNSSPlaceAuthRecord to vu.Activities_GnssRecord
	vuRecord := &vuv1.Activities_GnssRecord{}
	vuRecord.SetTimestamp(ddRecord.GetTimestamp())
	vuRecord.SetGnssAccuracy(ddRecord.GetGnssAccuracy())
	vuRecord.SetGeoCoordinates(ddRecord.GetGeoCoordinates())
	vuRecord.SetAuthenticationStatus(ddRecord.GetAuthenticationStatus())

	return vuRecord, nil
}

// splitVuLoadUnloadRecord splits data into 60-byte VuLoadUnloadRecord records
func splitVuLoadUnloadRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const loadUnloadRecordSize = 60

	if len(data) < loadUnloadRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return loadUnloadRecordSize, data[:loadUnloadRecordSize], nil
}

func parseVuLoadUnloadRecordArray(data []byte, offset int) ([]*vuv1.Activities_LoadUnloadRecord, int, error) {
	// VuLoadUnloadRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuLoadUnloadRecord -- 58 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for load/unload record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x18 for VuLoadUnloadRecord)
	if recordType != 0x18 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 24 (VuLoadUnloadRecord)", recordType)
	}

	// Validate record size (should be 58 bytes for VuLoadUnloadRecord)
	if recordSize != 58 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 58", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuLoadUnloadRecord)

	var records []*vuv1.Activities_LoadUnloadRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuLoadUnloadRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuLoadUnloadRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 60

	return records, offset, nil
}

// unmarshalVuLoadUnloadRecord parses a single VuLoadUnloadRecord from a byte slice
func unmarshalVuLoadUnloadRecord(data []byte) (*vuv1.Activities_LoadUnloadRecord, error) {
	// VuLoadUnloadRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                                     -- 4 bytes
	//     operationType OperationType,                            -- 1 byte
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     gnssPlaceAuthRecord GNSSPlaceAuthRecord,                -- 12 bytes
	//     vehicleOdometerValue OdometerShort                      -- 3 bytes
	// }
	// Total: 60 bytes

	const vuLoadUnloadRecordSize = 60
	if len(data) < vuLoadUnloadRecordSize {
		return nil, fmt.Errorf("insufficient data for load/unload record: got %d, need %d", len(data), vuLoadUnloadRecordSize)
	}

	var opts dd.UnmarshalOptions
	record := &vuv1.Activities_LoadUnloadRecord{}

	// Parse timeStamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[0:4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse operationType (OperationType - 1 byte)
	operationType := data[4]
	record.SetOperationType(ddv1.OperationType(operationType))

	// Parse cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	driverCardData := data[5:25]
	_, err = opts.UnmarshalFullCardNumberAndGeneration(driverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal driver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	driverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	driverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	driverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberDriverSlot(driverPlaceholder)

	// Parse cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	codriverCardData := data[25:45]
	_, err = opts.UnmarshalFullCardNumberAndGeneration(codriverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal codriver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	codriverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	codriverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	codriverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberCodriverSlot(codriverPlaceholder)

	// Parse gnssPlaceAuthRecord (GNSSPlaceAuthRecord - 12 bytes)
	placeData := data[45:57]
	placeRecord, err := unmarshalGNSSPlaceAuthRecord(placeData, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GNSS place auth record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, err := opts.UnmarshalOdometer(data[57:60])
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, nil
}

// splitSignatureRecord splits data into variable-length Signature records
func splitSignatureRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// We need at least 6 bytes for the record array header
	if len(data) < 6 {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	// Read recordSize from the header (bytes 2-3)
	recordSize := binary.BigEndian.Uint16(data[2:4])

	// Validate record size
	if recordSize == 0 {
		return 0, nil, fmt.Errorf("invalid record size: got 0")
	}

	// Total record size = header (6 bytes) + signature data (recordSize bytes)
	totalSize := 6 + int(recordSize)

	if len(data) < totalSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return totalSize, data[:totalSize], nil
}

func parseSignatureRecordArray(data []byte, offset int) ([]byte, int, error) {
	// SignatureRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF Signature -- Variable size each
	// }

	// Use bufio.Scanner to parse the signature record
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitSignatureRecord)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, offset, fmt.Errorf("scanner error: %w", err)
		}
		return nil, offset, fmt.Errorf("no signature record found")
	}

	recordData := scanner.Bytes()

	// Validate record type (should be 0x08 for Signature)
	recordType := binary.BigEndian.Uint16(recordData[0:2])
	if recordType != 0x08 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 8 (Signature)", recordType)
	}

	// Extract signature data (skip 6-byte header)
	signature := make([]byte, len(recordData)-6)
	copy(signature, recordData[6:])

	// Update offset to reflect consumed data
	offset += len(recordData)

	return signature, offset, nil
}

// AppendVuActivities appends VU activities data to a buffer.
//
// The data type `VuActivities` is specified in the Data Dictionary, Section 2.2.6.2.
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    dateOfDay                        TimeReal,
//	    odometerValueMidnight            OdometerValueMidnight,
//	    vuCardIWData                     VuCardIWData,
//	    vuActivityDailyData              VuActivityDailyData,
//	    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
//	    vuSpecificConditionData          VuSpecificConditionData,
//	    signature                        SignatureFirstGen
//	}
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
//	    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                      VuCardIWRecordArray,
//	    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                      VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
//	    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
//	    signatureRecordArray                     SignatureRecordArray
//	}
