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
	if vin != nil {
		appendVuString(buf, vin.GetDecoded(), 17)
	} else {
		appendVuString(buf, "", 17)
	}

	// VehicleRegistrationIdentification (15 bytes: nation(1) + regnum(14))
	vehicleReg := overview.GetVehicleRegistrationWithNation()
	if vehicleReg != nil {
		appendUint8(buf, uint8(vehicleReg.GetNation()))
		// First byte of registration is codepage (assume codepage 1 = ISO-8859-1)
		appendUint8(buf, 1)
		// Registration number (13 bytes)
		number := vehicleReg.GetNumber()
		if number != nil {
			appendVuString(buf, number.GetDecoded(), 13)
		} else {
			appendVuString(buf, "", 13)
		}
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
		companyName := lock.GetCompanyName()
		if companyName != nil {
			appendVuString(buf, companyName.GetDecoded(), 36)
		} else {
			appendVuString(buf, "", 36)
		}
		companyAddress := lock.GetCompanyAddress()
		if companyAddress != nil {
			appendVuString(buf, companyAddress.GetDecoded(), 36)
		} else {
			appendVuString(buf, "", 36)
		}
		appendVuFullCardNumber(buf, lock.GetCompanyCardNumber(), 16) // Card number field
	}

	// VuControlActivityData - variable length
	controlActivities := overview.GetControlActivities()
	for _, control := range controlActivities {
		controlType := control.GetControlType()
		if controlType != nil {
			// Convert ControlType to byte bitmask
			var b byte
			if controlType.GetCardDownloading() {
				b |= 0x80 // bit 'c'
			}
			if controlType.GetVuDownloading() {
				b |= 0x40 // bit 'v'
			}
			if controlType.GetPrinting() {
				b |= 0x20 // bit 'p'
			}
			if controlType.GetDisplay() {
				b |= 0x10 // bit 'd'
			}
			if controlType.GetCalibrationChecking() {
				b |= 0x08 // bit 'e'
			}
			appendUint8(buf, b)
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

func mapSlotCardTypeToUint8(cardType datadictionaryv1.SlotCardType) uint8 {
	switch cardType {
	case datadictionaryv1.SlotCardType_NO_CARD:
		return 0
	case datadictionaryv1.SlotCardType_DRIVER_CARD_INSERTED:
		return 1
	case datadictionaryv1.SlotCardType_WORKSHOP_CARD_INSERTED:
		return 2
	case datadictionaryv1.SlotCardType_CONTROL_CARD_INSERTED:
		return 3
	case datadictionaryv1.SlotCardType_COMPANY_CARD_INSERTED:
		return 4
	default:
		return 0
	}
}
