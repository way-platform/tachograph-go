package tachograph

import (
	"bytes"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendOverview marshals VU overview data for different generations
func AppendOverview(buf *bytes.Buffer, overview *vuv1.Overview) {
	if overview == nil {
		return
	}

	switch overview.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		appendOverviewGen1(buf, overview)
	case ddv1.Generation_GENERATION_2:
		appendOverviewGen2(buf, overview)
	}
}

func appendOverviewGen1(buf *bytes.Buffer, overview *vuv1.Overview) {
	// Gen1 Overview structure based on benchmark definitions
	// See VuOverviewFirstGen in benchmark/tachoparser/pkg/decoder/definitions.go

	// MemberStateCertificate (194 bytes)
	memberStateCert := overview.GetMemberStateCertificate()
	if len(memberStateCert) >= 194 {
		buf.Write(appendVuBytes(nil, memberStateCert[:194]))
	} else {
		// Pad to 194 bytes
		padded := make([]byte, 194)
		copy(padded, memberStateCert)
		buf.Write(appendVuBytes(nil, padded))
	}

	// VuCertificate (194 bytes)
	vuCert := overview.GetVuCertificate()
	if len(vuCert) >= 194 {
		buf.Write(appendVuBytes(nil, vuCert[:194]))
	} else {
		// Pad to 194 bytes
		padded := make([]byte, 194)
		copy(padded, vuCert)
		buf.Write(appendVuBytes(nil, padded))
	}

	// VehicleIdentificationNumber (17 bytes)
	vin := overview.GetVehicleIdentificationNumber()
	if vin != nil {
		buf.Write(appendVuString(nil, vin.GetDecoded(), 17))
	} else {
		buf.Write(appendVuString(nil, "", 17))
	}

	// VehicleRegistrationIdentification (15 bytes: nation(1) + regnum(14))
	vehicleReg := overview.GetVehicleRegistrationWithNation()
	if vehicleReg != nil {
		buf.Write(appendUint8(nil, uint8(vehicleReg.GetNation())))
		// First byte of registration is codepage (assume codepage 1 = ISO-8859-1)
		buf.Write(appendUint8(nil, 1))
		// Registration number (13 bytes)
		number := vehicleReg.GetNumber()
		if number != nil {
			buf.Write(appendVuString(nil, number.GetDecoded(), 13))
		} else {
			buf.Write(appendVuString(nil, "", 13))
		}
	} else {
		// Default values
		buf.Write(appendUint8(nil, 0))         // nation
		buf.Write(appendUint8(nil, 1))         // codepage
		buf.Write(appendVuString(nil, "", 13)) // empty registration
	}

	// CurrentDateTime (4 bytes)
	buf.Write(appendVuTimeReal(nil, overview.GetCurrentDateTime()))

	// VuDownloadablePeriod (8 bytes: 4 bytes min + 4 bytes max)
	downloadablePeriod := overview.GetDownloadablePeriod()
	if downloadablePeriod != nil {
		buf.Write(appendVuTimeReal(nil, downloadablePeriod.GetMinTime()))
		buf.Write(appendVuTimeReal(nil, downloadablePeriod.GetMaxTime()))
	} else {
		buf.Write(appendUint32(nil, 0))
		buf.Write(appendUint32(nil, 0))
	}

	// CardSlotsStatus (1 byte - driver and co-driver slots)
	driverSlot := mapSlotCardTypeToUint8(overview.GetDriverSlotCard())
	coDriverSlot := mapSlotCardTypeToUint8(overview.GetCoDriverSlotCard())
	slotsStatus := (driverSlot << 4) | (coDriverSlot & 0x0F)
	buf.Write(appendUint8(nil, slotsStatus))

	// VuDownloadActivityData (4 bytes - last download time)
	downloadActivities := overview.GetDownloadActivities()
	if len(downloadActivities) > 0 {
		buf.Write(appendVuTimeReal(nil, downloadActivities[0].GetDownloadingTime()))
	} else {
		buf.Write(appendVuTimeReal(nil, nil))
	}

	// VuCompanyLocksData - variable length
	// For now, we'll append the company locks in a simplified format
	companyLocks := overview.GetCompanyLocks()
	for _, lock := range companyLocks {
		buf.Write(appendVuTimeReal(nil, lock.GetLockInTime()))
		buf.Write(appendVuTimeReal(nil, lock.GetLockOutTime()))
		companyName := lock.GetCompanyName()
		if companyName != nil {
			buf.Write(appendVuString(nil, companyName.GetDecoded(), 36))
		} else {
			buf.Write(appendVuString(nil, "", 36))
		}
		companyAddress := lock.GetCompanyAddress()
		if companyAddress != nil {
			buf.Write(appendVuString(nil, companyAddress.GetDecoded(), 36))
		} else {
			buf.Write(appendVuString(nil, "", 36))
		}
		buf.Write(appendVuFullCardNumber(nil, lock.GetCompanyCardNumber(), 16)) // Card number field
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
			buf.Write(appendUint8(nil, b))
		} else {
			buf.Write(appendUint8(nil, 0))
		}
		buf.Write(appendVuTimeReal(nil, control.GetControlTime()))
		buf.Write(appendVuFullCardNumber(nil, control.GetControlCardNumber(), 16))
		buf.Write(appendVuTimeReal(nil, control.GetDownloadPeriodBeginTime()))
		buf.Write(appendVuTimeReal(nil, control.GetDownloadPeriodEndTime()))
	}

	// Signature (128 bytes for Gen1)
	signature := overview.GetSignatureGen1()
	if len(signature) >= 128 {
		buf.Write(appendVuBytes(nil, signature[:128]))
	} else {
		// Pad to 128 bytes
		padded := make([]byte, 128)
		copy(padded, signature)
		buf.Write(appendVuBytes(nil, padded))
	}
}

func appendOverviewGen2(buf *bytes.Buffer, overview *vuv1.Overview) {
	// Gen2 Overview structure - more complex
	// For now, implement a basic version
	// In a full implementation, this would marshal the complete Gen2 structure

	// Add basic Gen2 fields as they become available
	// This is a placeholder for future Gen2 implementation
}

func mapSlotCardTypeToUint8(cardType ddv1.SlotCardType) uint8 {
	switch cardType {
	case ddv1.SlotCardType_NO_CARD:
		return 0
	case ddv1.SlotCardType_DRIVER_CARD_INSERTED:
		return 1
	case ddv1.SlotCardType_WORKSHOP_CARD_INSERTED:
		return 2
	case ddv1.SlotCardType_CONTROL_CARD_INSERTED:
		return 3
	case ddv1.SlotCardType_COMPANY_CARD_INSERTED:
		return 4
	default:
		return 0
	}
}
