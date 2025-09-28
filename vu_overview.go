package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalOverview parses VU overview data for different generations
//
// The data type `VuOverview` is specified in the Data Dictionary, Section 2.2.6.
//
// ASN.1 Definition:
//
//	VuOverviewFirstGen ::= SEQUENCE {
//	    memberStateCertificate            MemberStateCertificateFirstGen,
//	    vuCertificate                     VuCertificateFirstGen,
//	    vehicleIdentificationNumber       VehicleIdentificationNumber,
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,
//	    currentDateTime                   CurrentDateTime,
//	    vuDownloadablePeriod              VuDownloadablePeriod,
//	    cardSlotsStatus                   CardSlotsStatus,
//	    vuDownloadActivityData            VuDownloadActivityDataFirstGen,
//	    vuCompanyLocksData                VuCompanyLocksDataFirstGen,
//	    vuControlActivityData             VuControlActivityDataFirstGen,
//	    signature                         SignatureFirstGen
//	}
func unmarshalOverview(data []byte, offset int, overview *vuv1.Overview, generation int) (int, error) {
	startOffset := offset

	switch generation {
	case 1:
		overview.SetGeneration(ddv1.Generation_GENERATION_1)
		return unmarshalOverviewGen1(data, offset, overview, startOffset)
	case 2:
		overview.SetGeneration(ddv1.Generation_GENERATION_2)
		// For now, assume version 1 - we can add version detection later
		overview.SetVersion(vuv1.Version_VERSION_1)
		return unmarshalOverviewGen2(data, offset, overview, startOffset)
	default:
		return 0, nil
	}
}

func unmarshalOverviewGen1(data []byte, offset int, overview *vuv1.Overview, startOffset int) (int, error) {
	// Gen1 Overview structure based on benchmark definitions
	// See VuOverviewFirstGen in benchmark/tachoparser/pkg/decoder/definitions.go

	// MemberStateCertificate (194 bytes)
	memberStateCert, offset, err := readBytesFromBytes(data, offset, 194)
	if err != nil {
		return 0, err
	}
	overview.SetMemberStateCertificate(memberStateCert)

	// VuCertificate (194 bytes)
	vuCert, offset, err := readBytesFromBytes(data, offset, 194)
	if err != nil {
		return 0, err
	}
	overview.SetVuCertificate(vuCert)

	// VehicleIdentificationNumber (17 bytes)
	vinBytes, offset, err := readBytesFromBytes(data, offset, 17)
	if err != nil {
		return 0, err
	}
	vinStrValue, err := unmarshalIA5StringValue(vinBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to read VIN: %w", err)
	}
	overview.SetVehicleIdentificationNumber(vinStrValue)

	// VehicleRegistrationIdentification (15 bytes: nation(1) + regnum(14))
	nation, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return 0, err
	}
	regNumBytes, offset, err := readBytesFromBytes(data, offset, 14)
	if err != nil {
		return 0, err
	}

	// First byte is codepage, rest is registration number
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(ddv1.NationNumeric(nation))

	regNumStrValue, err := unmarshalIA5StringValue(regNumBytes[1:])
	if err != nil {
		return 0, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	vehicleReg.SetNumber(regNumStrValue)
	overview.SetVehicleRegistrationWithNation(vehicleReg)

	// CurrentDateTime (4 bytes)
	currentTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, err
	}
	overview.SetCurrentDateTime(timestamppb.New(time.Unix(currentTime, 0)))

	// VuDownloadablePeriod (8 bytes: 4 bytes min + 4 bytes max)
	minTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, err
	}
	maxTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, err
	}

	downloadablePeriod := &ddv1.DownloadablePeriod{}
	downloadablePeriod.SetMinTime(timestamppb.New(time.Unix(minTime, 0)))
	downloadablePeriod.SetMaxTime(timestamppb.New(time.Unix(maxTime, 0)))
	overview.SetDownloadablePeriod(downloadablePeriod)

	// CardSlotsStatus (1 byte - driver and co-driver slots)
	slotsStatus, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return 0, err
	}
	// Extract driver and co-driver slot info from the byte
	driverSlot := (slotsStatus >> 4) & 0x0F
	coDriverSlot := slotsStatus & 0x0F

	overview.SetDriverSlotCard(mapSlotCardType(driverSlot))
	overview.SetCoDriverSlotCard(mapSlotCardType(coDriverSlot))

	// VuDownloadActivityData (4 bytes - last download time)
	lastDownloadTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, err
	}
	// Create download activity record
	downloadActivity := &vuv1.Overview_DownloadActivity{}
	downloadActivity.SetDownloadingTime(timestamppb.New(time.Unix(lastDownloadTime, 0)))
	overview.SetDownloadActivities([]*vuv1.Overview_DownloadActivity{downloadActivity})

	// VuCompanyLocksData - variable length, need to determine size
	// For now, we'll skip this complex structure

	// VuControlActivityData - variable length
	// For now, we'll skip this complex structure

	// Signature (128 bytes for Gen1)
	signature, offset, err := readBytesFromBytes(data, offset, 128)
	if err != nil {
		return 0, err
	}
	overview.SetSignatureGen1(signature)

	bytesRead := offset - startOffset
	return bytesRead, nil
}

func unmarshalOverviewGen2(data []byte, offset int, overview *vuv1.Overview, startOffset int) (int, error) {
	// Gen2 Overview structure - more complex, implement basic version for now
	// This would need detailed implementation based on the specific Gen2 structure

	// For now, just read a minimal amount to avoid errors
	// In a full implementation, this would parse the complete Gen2 structure

	bytesRead := offset - startOffset
	return bytesRead, nil
}

// AppendOverview marshals VU overview data for different generations
//
// The data type `VuOverview` is specified in the Data Dictionary, Section 2.2.6.
//
// ASN.1 Definition:
//
//	VuOverviewFirstGen ::= SEQUENCE {
//	    memberStateCertificate            MemberStateCertificateFirstGen,
//	    vuCertificate                     VuCertificateFirstGen,
//	    vehicleIdentificationNumber       VehicleIdentificationNumber,
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,
//	    currentDateTime                   CurrentDateTime,
//	    vuDownloadablePeriod              VuDownloadablePeriod,
//	    cardSlotsStatus                   CardSlotsStatus,
//	    vuDownloadActivityData            VuDownloadActivityDataFirstGen,
//	    vuCompanyLocksData                VuCompanyLocksDataFirstGen,
//	    vuControlActivityData             VuControlActivityDataFirstGen,
//	    signature                         SignatureFirstGen
//	}
func appendOverview(buf *bytes.Buffer, overview *vuv1.Overview) {
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
		buf.WriteByte(uint8(vehicleReg.GetNation()))
		// First byte of registration is codepage (assume codepage 1 = ISO-8859-1)
		buf.WriteByte(1)
		// Registration number (13 bytes)
		number := vehicleReg.GetNumber()
		if number != nil {
			buf.Write(appendVuString(nil, number.GetDecoded(), 13))
		} else {
			buf.Write(appendVuString(nil, "", 13))
		}
	} else {
		// Default values
		buf.WriteByte(0)                       // nation
		buf.WriteByte(1)                       // codepage
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
		buf.Write(make([]byte, 4)) // 4 zero bytes
		buf.Write(make([]byte, 4)) // 4 zero bytes
	}

	// CardSlotsStatus (1 byte - driver and co-driver slots)
	driverSlot := mapSlotCardTypeToUint8(overview.GetDriverSlotCard())
	coDriverSlot := mapSlotCardTypeToUint8(overview.GetCoDriverSlotCard())
	slotsStatus := (driverSlot << 4) | (coDriverSlot & 0x0F)
	buf.WriteByte(slotsStatus)

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
			buf.WriteByte(b)
		} else {
			buf.WriteByte(0)
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

// VU-specific helper functions for binary operations

// appendVuBytes appends a byte slice to dst
func appendVuBytes(dst []byte, data []byte) []byte {
	return append(dst, data...)
}

// appendVuString appends a string to dst with a fixed length, padding with null bytes
func appendVuString(dst []byte, s string, length int) []byte {
	result := make([]byte, length)
	copy(result, []byte(s))
	// Pad with null bytes
	for i := len(s); i < length; i++ {
		result[i] = 0
	}
	return append(dst, result...)
}

// appendVuTimeReal appends a TimeReal value (4 bytes) to dst
func appendVuTimeReal(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	return binary.BigEndian.AppendUint32(dst, uint32(ts.GetSeconds()))
}

// appendVuFullCardNumber appends a FullCardNumber to dst with a fixed length
func appendVuFullCardNumber(dst []byte, cardNumber *ddv1.FullCardNumber, length int) []byte {
	if cardNumber == nil {
		return append(dst, make([]byte, length)...)
	}
	// TODO: Implement proper FullCardNumber serialization
	return append(dst, make([]byte, length)...)
}
