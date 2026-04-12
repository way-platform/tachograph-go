package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/safemath"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfOverview dispatches to generation-specific size calculation.
func sizeOfOverview(data []byte, transferType vuv1.TransferType) (totalSize, signatureSize int, err error) {
	switch transferType {
	case vuv1.TransferType_OVERVIEW_GEN1:
		return sizeOfOverviewGen1(data)
	case vuv1.TransferType_OVERVIEW_GEN2_V1:
		return sizeOfOverviewGen2V1(data)
	case vuv1.TransferType_OVERVIEW_GEN2_V2:
		return sizeOfOverviewGen2V2(data)
	default:
		return 0, 0, fmt.Errorf("unsupported transfer type for Overview: %v", transferType)
	}
}

// sizeOfOverviewGen1 calculates total size for Gen1 Overview including signature.
//
// Overview Gen1 structure (from Appendix 7, Section 2.2.6.2):
// - MemberStateCertificate: 194 bytes (CertificateFirstGen)
// - VuCertificate: 194 bytes (CertificateFirstGen)
// - VehicleIdentificationNumber: 17 bytes
// - VehicleRegistrationIdentification: 15 bytes (1 nation + 1 codePage + 13 vrn)
// - CurrentDateTime: 4 bytes (TimeReal)
// - VuDownloadablePeriod: 8 bytes (2 x TimeReal)
// - CardSlotsStatus: 1 byte
// - VuDownloadActivityData: 58 bytes (4 + 18 + 36)
//   - DownloadingTime: 4 bytes (TimeReal)
//   - FullCardNumber: 18 bytes (1 EquipmentType + 1 NationNumeric + 16 CardNumber)
//   - CompanyOrWorkshopName: 36 bytes (1 CodePage + 35 Name bytes)
//
// - VuCompanyLocksData: 1 byte (noOfLocks) + (noOfLocks * 98 bytes per record)
//   - VuCompanyLocksRecordFirstGen: 98 bytes (4 LockInTime + 4 LockOutTime + 36 CompanyName + 36 CompanyAddress + 18 CompanyCardNumber)
//
// - VuControlActivityData: 1 byte (noOfControls) + (noOfControls * 31 bytes per record)
//   - VuControlActivityRecordFirstGen: 31 bytes (1 ControlType + 4 ControlTime + 18 ControlCardNumber + 4 DownloadPeriodBegin + 4 DownloadPeriodEnd)
//
// - Signature: 128 bytes (RSA)
func sizeOfOverviewGen1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// Fixed-size header sections (491 bytes total)
	offset += 194 // MemberStateCertificate
	offset += 194 // VuCertificate
	offset += 17  // VehicleIdentificationNumber
	offset += 15  // VehicleRegistrationIdentification (1 nation + 1 codePage + 13 vrn)
	offset += 4   // CurrentDateTime (TimeReal)
	offset += 8   // VuDownloadablePeriod (2 x TimeReal)
	offset += 1   // CardSlotsStatus
	offset += 58  // VuDownloadActivityData (4 + 18 + 36)

	// VuCompanyLocksData: 1 byte count + variable records
	if len(data[offset:]) < 1 {
		return 0, 0, fmt.Errorf("insufficient data for noOfLocks")
	}
	noOfLocks := data[offset]
	offset += 1

	// Each VuCompanyLocksRecordFirstGen: 4 + 4 + 36 + 36 + 18 = 98 bytes
	// (lockInTime, lockOutTime, companyName, companyAddress, companyCardNumber)
	const vuCompanyLocksRecordSize = 98
	locksSize, mulErr := safemath.MulInt(int(noOfLocks), vuCompanyLocksRecordSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("company locks overflow: count=%d recordSize=%d: %w", noOfLocks, vuCompanyLocksRecordSize, mulErr)
	}
	offset += locksSize

	// VuControlActivityData: 1 byte count + variable records
	if len(data[offset:]) < 1 {
		return 0, 0, fmt.Errorf("insufficient data for noOfControls")
	}
	noOfControls := data[offset]
	offset += 1

	// Each VuControlActivityRecordFirstGen: 1 + 4 + 18 + 4 + 4 = 31 bytes
	// (controlType, controlTime, controlCardNumber, downloadPeriodBeginTime, downloadPeriodEndTime)
	const vuControlActivityRecordSize = 31
	controlsSize, mulErr := safemath.MulInt(int(noOfControls), vuControlActivityRecordSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("control activity overflow: count=%d recordSize=%d: %w", noOfControls, vuControlActivityRecordSize, mulErr)
	}
	offset += controlsSize

	// Signature: 128 bytes for Gen1 RSA
	const gen1SignatureSize = 128
	offset += gen1SignatureSize

	return offset, gen1SignatureSize, nil
}

// sizeOfOverviewGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
//
// Gen2 uses RecordArray structures with 5-byte headers that include the size.
// We parse each RecordArray header sequentially to determine the total size.
func sizeOfOverviewGen2V1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// MemberStateCertificateRecordArray
	size, sizeErr := sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("MemberStateCertificateRecordArray: %w", sizeErr)
	}
	offset += size

	// VUCertificateRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VUCertificateRecordArray: %w", sizeErr)
	}
	offset += size

	// VehicleIdentificationNumberRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VehicleIdentificationNumberRecordArray: %w", sizeErr)
	}
	offset += size

	// VehicleRegistrationIdentificationRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VehicleRegistrationIdentificationRecordArray: %w", sizeErr)
	}
	offset += size

	// CurrentDateTimeRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("CurrentDateTimeRecordArray: %w", sizeErr)
	}
	offset += size

	// VuDownloadablePeriodRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuDownloadablePeriodRecordArray: %w", sizeErr)
	}
	offset += size

	// CardSlotsStatusRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("CardSlotsStatusRecordArray: %w", sizeErr)
	}
	offset += size

	// VuDownloadActivityDataRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuDownloadActivityDataRecordArray: %w", sizeErr)
	}
	offset += size

	// VuCompanyLocksRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuCompanyLocksRecordArray: %w", sizeErr)
	}
	offset += size

	// VuControlActivityRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuControlActivityRecordArray: %w", sizeErr)
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

// sizeOfOverviewGen2V2 calculates size by parsing all Gen2 V2 RecordArrays.
//
// Gen2 V2 has an additional VehicleRegistrationNumberRecordArray between
// VehicleIdentificationNumberRecordArray and CurrentDateTimeRecordArray.
func sizeOfOverviewGen2V2(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// MemberStateCertificateRecordArray
	size, sizeErr := sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("MemberStateCertificateRecordArray: %w", sizeErr)
	}
	offset += size

	// VUCertificateRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VUCertificateRecordArray: %w", sizeErr)
	}
	offset += size

	// VehicleIdentificationNumberRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VehicleIdentificationNumberRecordArray: %w", sizeErr)
	}
	offset += size

	// VehicleRegistrationNumberRecordArray (Gen2 V2 addition)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VehicleRegistrationNumberRecordArray: %w", sizeErr)
	}
	offset += size

	// CurrentDateTimeRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("CurrentDateTimeRecordArray: %w", sizeErr)
	}
	offset += size

	// VuDownloadablePeriodRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuDownloadablePeriodRecordArray: %w", sizeErr)
	}
	offset += size

	// CardSlotsStatusRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("CardSlotsStatusRecordArray: %w", sizeErr)
	}
	offset += size

	// VuDownloadActivityDataRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuDownloadActivityDataRecordArray: %w", sizeErr)
	}
	offset += size

	// VuCompanyLocksRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuCompanyLocksRecordArray: %w", sizeErr)
	}
	offset += size

	// VuControlActivityRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuControlActivityRecordArray: %w", sizeErr)
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
