package tachograph

import (
	"bytes"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendOverview marshals VU overview data for different generations
func AppendOverview(buf *bytes.Buffer, overview *vuv1.Overview) {
	if overview == nil {
		return
	}

	switch overview.GetGeneration() {
	case datadictionaryv1.Generation_GENERATION_1:
		appendOverviewGen1(buf, overview)
	case datadictionaryv1.Generation_GENERATION_2:
		appendOverviewGen2(buf, overview)
	}
}

func appendOverviewGen1(buf *bytes.Buffer, overview *vuv1.Overview) {
	// Gen1 Overview structure based on benchmark definitions
	// See VuOverviewFirstGen in benchmark/tachoparser/pkg/decoder/definitions.go

	// MemberStateCertificate (194 bytes)
	memberStateCert := overview.GetMemberStateCertificate()
	if len(memberStateCert) >= 194 {
		appendVuBytes(buf, memberStateCert[:194])
	} else {
		// Pad to 194 bytes
		padded := make([]byte, 194)
		copy(padded, memberStateCert)
		appendVuBytes(buf, padded)
	}

	// VuCertificate (194 bytes)
	vuCert := overview.GetVuCertificate()
	if len(vuCert) >= 194 {
		appendVuBytes(buf, vuCert[:194])
	} else {
		// Pad to 194 bytes
		padded := make([]byte, 194)
		copy(padded, vuCert)
		appendVuBytes(buf, padded)
	}

	// VehicleIdentificationNumber (17 bytes)
	vin := overview.GetVehicleIdentificationNumber()
	appendVuString(buf, vin, 17)

	// VehicleRegistrationIdentification (15 bytes: nation(1) + regnum(14))
	vehicleReg := overview.GetVehicleRegistrationWithNation()
	if vehicleReg != nil {
		appendUint8(buf, uint8(vehicleReg.GetNation()))
		// First byte of registration is codepage (assume codepage 1 = ISO-8859-1)
		appendUint8(buf, 1)
		// Registration number (13 bytes)
		appendVuString(buf, vehicleReg.GetNumber(), 13)
	} else {
		// Default values
		appendUint8(buf, 0)         // nation
		appendUint8(buf, 1)         // codepage
		appendVuString(buf, "", 13) // empty registration
	}

	// CurrentDateTime (4 bytes)
	appendVuTimeReal(buf, overview.GetCurrentDateTime())

	// VuDownloadablePeriod (8 bytes: 4 bytes min + 4 bytes max)
	downloadablePeriod := overview.GetDownloadablePeriod()
	if downloadablePeriod != nil {
		appendVuTimeReal(buf, downloadablePeriod.GetMinTime())
		appendVuTimeReal(buf, downloadablePeriod.GetMaxTime())
	} else {
		appendUint32(buf, 0)
		appendUint32(buf, 0)
	}

	// CardSlotsStatus (1 byte - driver and co-driver slots)
	driverSlot := mapSlotCardTypeToUint8(overview.GetDriverSlotCard())
	coDriverSlot := mapSlotCardTypeToUint8(overview.GetCoDriverSlotCard())
	slotsStatus := (driverSlot << 4) | (coDriverSlot & 0x0F)
	appendUint8(buf, slotsStatus)

	// VuDownloadActivityData (4 bytes - last download time)
	appendVuTimeReal(buf, overview.GetLastDownloadTime())

	// VuCompanyLocksData - variable length
	// For now, we'll append the company locks in a simplified format
	companyLocks := overview.GetCompanyLocks()
	for _, lock := range companyLocks {
		appendVuTimeReal(buf, lock.GetLockInTime())
		appendVuTimeReal(buf, lock.GetLockOutTime())
		appendVuString(buf, lock.GetCompanyName(), 36)               // Name field is typically 36 bytes
		appendVuString(buf, lock.GetCompanyAddress(), 36)            // Address field
		appendVuFullCardNumber(buf, lock.GetCompanyCardNumber(), 16) // Card number field
	}

	// VuControlActivityData - variable length
	controlActivities := overview.GetControlActivities()
	for _, control := range controlActivities {
		controlType := control.GetControlType()
		if len(controlType) > 0 {
			appendUint8(buf, controlType[0])
		} else {
			appendUint8(buf, 0)
		}
		appendVuTimeReal(buf, control.GetControlTime())
		appendVuFullCardNumber(buf, control.GetControlCardNumber(), 16)
		appendVuTimeReal(buf, control.GetDownloadPeriodBeginTime())
		appendVuTimeReal(buf, control.GetDownloadPeriodEndTime())
	}

	// Signature (128 bytes for Gen1)
	signature := overview.GetSignatureGen1()
	if len(signature) >= 128 {
		appendVuBytes(buf, signature[:128])
	} else {
		// Pad to 128 bytes
		padded := make([]byte, 128)
		copy(padded, signature)
		appendVuBytes(buf, padded)
	}
}

func appendOverviewGen2(buf *bytes.Buffer, overview *vuv1.Overview) {
	// Gen2 Overview structure - more complex
	// For now, implement a basic version
	// In a full implementation, this would marshal the complete Gen2 structure

	// Add basic Gen2 fields as they become available
	// This is a placeholder for future Gen2 implementation
}

func mapSlotCardTypeToUint8(cardType vuv1.Overview_SlotCardType) uint8 {
	switch cardType {
	case vuv1.Overview_NO_CARD:
		return 0
	case vuv1.Overview_DRIVER_CARD:
		return 1
	case vuv1.Overview_WORKSHOP_CARD:
		return 2
	case vuv1.Overview_CONTROL_CARD:
		return 3
	case vuv1.Overview_COMPANY_CARD:
		return 4
	default:
		return 0
	}
}
