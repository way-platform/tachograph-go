package vu

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/safemath"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfActivities dispatches to generation-specific size calculation.
func sizeOfActivities(data []byte, transferType vuv1.TransferType) (totalSize, signatureSize int, err error) {
	switch transferType {
	case vuv1.TransferType_ACTIVITIES_GEN1:
		return sizeOfActivitiesGen1(data)
	case vuv1.TransferType_ACTIVITIES_GEN2_V1:
		return sizeOfActivitiesGen2V1(data)
	case vuv1.TransferType_ACTIVITIES_GEN2_V2:
		return sizeOfActivitiesGen2V2(data)
	default:
		return 0, 0, fmt.Errorf("unsupported transfer type for Activities: %v", transferType)
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
func sizeOfActivitiesGen1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// Fixed-size header sections (7 bytes total)
	offset += 4 // TimeReal (date of day downloaded)
	offset += 3 // OdometerValueMidnight

	// VuCardIWData: 2 bytes count + variable records
	if len(data[offset:]) < 2 {
		return 0, 0, fmt.Errorf("insufficient data for noOfIWRecords")
	}
	noOfIWRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each VuCardIWRecordFirstGen: 129 bytes
	// (cardHolderName 72, fullCardNumber 18, cardExpiryDate 4, cardInsertionTime 4,
	//  vehicleOdometerValueAtInsertion 3, cardSlotNumber 1, cardWithdrawalTime 4,
	//  vehicleOdometerValueAtWithdrawal 3, previousVehicleInfo 19, manualInputFlag 1)
	const vuCardIWRecordSize = 129
	iwSize, mulErr := safemath.MulInt(int(noOfIWRecords), vuCardIWRecordSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("IW records overflow: count=%d recordSize=%d: %w", noOfIWRecords, vuCardIWRecordSize, mulErr)
	}
	offset += iwSize

	// VuActivityDailyData: 2 bytes count + variable activity changes
	if len(data[offset:]) < 2 {
		return 0, 0, fmt.Errorf("insufficient data for noOfActivityChanges")
	}
	noOfActivityChanges := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each ActivityChangeInfo: 2 bytes
	const activityChangeInfoSize = 2
	actSize, mulErr := safemath.MulInt(int(noOfActivityChanges), activityChangeInfoSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("activity changes overflow: count=%d recordSize=%d: %w", noOfActivityChanges, activityChangeInfoSize, mulErr)
	}
	offset += actSize

	// VuPlaceDailyWorkPeriodData: 1 byte count + variable place records
	if len(data[offset:]) < 1 {
		return 0, 0, fmt.Errorf("insufficient data for noOfPlaceRecords")
	}
	noOfPlaceRecords := data[offset]
	offset += 1

	// Each VuPlaceDailyWorkPeriodRecordFirstGen: 28 bytes (18 FullCardNumber + 10 PlaceRecordFirstGen)
	const vuPlaceDailyWorkPeriodRecordSize = 28
	placeSize, mulErr := safemath.MulInt(int(noOfPlaceRecords), vuPlaceDailyWorkPeriodRecordSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("place records overflow: count=%d recordSize=%d: %w", noOfPlaceRecords, vuPlaceDailyWorkPeriodRecordSize, mulErr)
	}
	offset += placeSize

	// VuSpecificConditionData: 2 bytes count + variable condition records
	if len(data[offset:]) < 2 {
		return 0, 0, fmt.Errorf("insufficient data for noOfSpecificConditionRecords")
	}
	noOfSpecificConditionRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each SpecificConditionRecord: 5 bytes (4 TimeReal + 1 SpecificConditionType)
	const specificConditionRecordSize = 5
	condSize, mulErr := safemath.MulInt(int(noOfSpecificConditionRecords), specificConditionRecordSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("specific condition records overflow: count=%d recordSize=%d: %w", noOfSpecificConditionRecords, specificConditionRecordSize, mulErr)
	}
	offset += condSize

	// Signature: 128 bytes for Gen1 RSA
	const gen1SignatureSize = 128
	offset += gen1SignatureSize

	return offset, gen1SignatureSize, nil
}

// sizeOfActivitiesGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
func sizeOfActivitiesGen2V1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// DateOfDayDownloadedRecordArray
	size, sizeErr := sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("DateOfDayDownloadedRecordArray: %w", sizeErr)
	}
	offset += size

	// OdometerValueMidnightRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("OdometerValueMidnightRecordArray: %w", sizeErr)
	}
	offset += size

	// VuCardIWRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuCardIWRecordArray: %w", sizeErr)
	}
	offset += size

	// VuActivityDailyRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuActivityDailyRecordArray: %w", sizeErr)
	}
	offset += size

	// VuPlaceDailyWorkPeriodRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuPlaceDailyWorkPeriodRecordArray: %w", sizeErr)
	}
	offset += size

	// VuSpecificConditionRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuSpecificConditionRecordArray: %w", sizeErr)
	}
	offset += size

	// VuGNSSADRecordArray (Gen2+)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuGNSSADRecordArray: %w", sizeErr)
	}
	offset += size

	// SignatureRecordArray (last)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("SignatureRecordArray: %w", sizeErr)
	}
	signatureSizeGen2 := size
	offset += size

	return offset, signatureSizeGen2, nil
}

// sizeOfActivitiesGen2V2 calculates size by parsing all Gen2 V2 RecordArrays.
// Must handle VuBorderCrossingRecordArray and VuLoadUnloadRecordArray.
func sizeOfActivitiesGen2V2(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// DateOfDayDownloadedRecordArray
	size, sizeErr := sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("DateOfDayDownloadedRecordArray: %w", sizeErr)
	}
	offset += size

	// OdometerValueMidnightRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("OdometerValueMidnightRecordArray: %w", sizeErr)
	}
	offset += size

	// VuCardIWRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuCardIWRecordArray: %w", sizeErr)
	}
	offset += size

	// VuActivityDailyRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuActivityDailyRecordArray: %w", sizeErr)
	}
	offset += size

	// VuPlaceDailyWorkPeriodRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuPlaceDailyWorkPeriodRecordArray: %w", sizeErr)
	}
	offset += size

	// VuSpecificConditionRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuSpecificConditionRecordArray: %w", sizeErr)
	}
	offset += size

	// VuGNSSADRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuGNSSADRecordArray: %w", sizeErr)
	}
	offset += size

	// VuBorderCrossingRecordArray (Gen2 V2+)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuBorderCrossingRecordArray: %w", sizeErr)
	}
	offset += size

	// VuLoadUnloadRecordArray (Gen2 V2+)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuLoadUnloadRecordArray: %w", sizeErr)
	}
	offset += size

	// SignatureRecordArray (last)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("SignatureRecordArray: %w", sizeErr)
	}
	signatureSizeGen2 := size
	offset += size

	return offset, signatureSizeGen2, nil
}
