package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfOverview dispatches to generation-specific size calculation.
func sizeOfOverview(data []byte, transferType vuv1.TransferType) (int, error) {
	switch transferType {
	case vuv1.TransferType_OVERVIEW_GEN1:
		return sizeOfOverviewGen1(data)
	case vuv1.TransferType_OVERVIEW_GEN2_V1:
		return sizeOfOverviewGen2V1(data)
	case vuv1.TransferType_OVERVIEW_GEN2_V2:
		return sizeOfOverviewGen2V2(data)
	default:
		return 0, fmt.Errorf("unsupported transfer type for Overview: %v", transferType)
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
func sizeOfOverviewGen1(data []byte) (int, error) {
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
		return 0, fmt.Errorf("insufficient data for noOfLocks")
	}
	noOfLocks := data[offset]
	offset += 1

	// Each VuCompanyLocksRecordFirstGen: 4 + 4 + 36 + 36 + 18 = 98 bytes
	// (lockInTime, lockOutTime, companyName, companyAddress, companyCardNumber)
	const vuCompanyLocksRecordSize = 98
	offset += int(noOfLocks) * vuCompanyLocksRecordSize

	// VuControlActivityData: 1 byte count + variable records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfControls")
	}
	noOfControls := data[offset]
	offset += 1

	// Each VuControlActivityRecordFirstGen: 1 + 4 + 18 + 4 + 4 = 31 bytes
	// (controlType, controlTime, controlCardNumber, downloadPeriodBeginTime, downloadPeriodEndTime)
	const vuControlActivityRecordSize = 31
	offset += int(noOfControls) * vuControlActivityRecordSize

	// Signature: 128 bytes for Gen1 RSA
	offset += 128

	return offset, nil
}

// sizeOfOverviewGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
//
// Gen2 uses RecordArray structures with 5-byte headers that include the size.
// We parse each RecordArray header sequentially to determine the total size.
func sizeOfOverviewGen2V1(data []byte) (int, error) {
	offset := 0

	// MemberStateCertificateRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("MemberStateCertificateRecordArray: %w", err)
	}
	offset += size

	// VUCertificateRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VUCertificateRecordArray: %w", err)
	}
	offset += size

	// VehicleIdentificationNumberRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VehicleIdentificationNumberRecordArray: %w", err)
	}
	offset += size

	// VehicleRegistrationIdentificationRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VehicleRegistrationIdentificationRecordArray: %w", err)
	}
	offset += size

	// CurrentDateTimeRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("CurrentDateTimeRecordArray: %w", err)
	}
	offset += size

	// VuDownloadablePeriodRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuDownloadablePeriodRecordArray: %w", err)
	}
	offset += size

	// CardSlotsStatusRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("CardSlotsStatusRecordArray: %w", err)
	}
	offset += size

	// VuDownloadActivityDataRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuDownloadActivityDataRecordArray: %w", err)
	}
	offset += size

	// VuCompanyLocksRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuCompanyLocksRecordArray: %w", err)
	}
	offset += size

	// VuControlActivityRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuControlActivityRecordArray: %w", err)
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

// sizeOfOverviewGen2V2 calculates size by parsing all Gen2 V2 RecordArrays.
//
// Gen2 V2 has an additional VehicleRegistrationNumberRecordArray between
// VehicleIdentificationNumberRecordArray and CurrentDateTimeRecordArray.
func sizeOfOverviewGen2V2(data []byte) (int, error) {
	offset := 0

	// MemberStateCertificateRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("MemberStateCertificateRecordArray: %w", err)
	}
	offset += size

	// VUCertificateRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VUCertificateRecordArray: %w", err)
	}
	offset += size

	// VehicleIdentificationNumberRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VehicleIdentificationNumberRecordArray: %w", err)
	}
	offset += size

	// VehicleRegistrationNumberRecordArray (Gen2 V2 addition)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VehicleRegistrationNumberRecordArray: %w", err)
	}
	offset += size

	// CurrentDateTimeRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("CurrentDateTimeRecordArray: %w", err)
	}
	offset += size

	// VuDownloadablePeriodRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuDownloadablePeriodRecordArray: %w", err)
	}
	offset += size

	// CardSlotsStatusRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("CardSlotsStatusRecordArray: %w", err)
	}
	offset += size

	// VuDownloadActivityDataRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuDownloadActivityDataRecordArray: %w", err)
	}
	offset += size

	// VuCompanyLocksRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuCompanyLocksRecordArray: %w", err)
	}
	offset += size

	// VuControlActivityRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuControlActivityRecordArray: %w", err)
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
