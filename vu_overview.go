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
	// Gen2 Overview structure - uses record arrays instead of fixed fields
	// VuOverviewSecondGen ::= SEQUENCE {
	//     memberStateCertificateRecordArray    MemberStateCertificateRecordArray,
	//     vuCertificateRecordArray             VuCertificateRecordArray,
	//     vehicleIdentificationNumberRecordArray VehicleIdentificationNumberRecordArray,
	//     vehicleRegistrationIdentificationRecordArray VehicleRegistrationIdentificationRecordArray,
	//     currentDateTimeRecordArray           CurrentDateTimeRecordArray,
	//     vuDownloadablePeriodRecordArray      VuDownloadablePeriodRecordArray,
	//     cardSlotsStatusRecordArray           CardSlotsStatusRecordArray,
	//     vuDownloadActivityDataRecordArray    VuDownloadActivityDataRecordArray,
	//     vuCompanyLocksRecordArray            VuCompanyLocksRecordArray,
	//     vuControlActivityRecordArray         VuControlActivityRecordArray,
	//     signatureRecordArray                 SignatureRecordArray
	// }

	// Parse MemberStateCertificateRecordArray
	memberStateCerts, offset, err := parseMemberStateCertificateRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse member state certificate record array: %w", err)
	}
	if len(memberStateCerts) > 0 {
		overview.SetMemberStateCertificate(memberStateCerts[0]) // Use first certificate
	}

	// Parse VuCertificateRecordArray
	vuCerts, offset, err := parseVuCertificateRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse VU certificate record array: %w", err)
	}
	if len(vuCerts) > 0 {
		overview.SetVuCertificate(vuCerts[0]) // Use first certificate
	}

	// Parse VehicleIdentificationNumberRecordArray
	vins, offset, err := parseVehicleIdentificationNumberRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse vehicle identification number record array: %w", err)
	}
	if len(vins) > 0 {
		overview.SetVehicleIdentificationNumber(vins[0]) // Use first VIN
	}

	// Parse VehicleRegistrationIdentificationRecordArray
	vehicleRegs, offset, err := parseVehicleRegistrationIdentificationRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse vehicle registration identification record array: %w", err)
	}
	if len(vehicleRegs) > 0 {
		overview.SetVehicleRegistrationWithNation(vehicleRegs[0]) // Use first registration
	}

	// Parse CurrentDateTimeRecordArray
	currentDateTimes, offset, err := parseCurrentDateTimeRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse current date time record array: %w", err)
	}
	if len(currentDateTimes) > 0 {
		overview.SetCurrentDateTime(currentDateTimes[0]) // Use first date time
	}

	// Parse VuDownloadablePeriodRecordArray
	downloadablePeriods, offset, err := parseVuDownloadablePeriodRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse VU downloadable period record array: %w", err)
	}
	if len(downloadablePeriods) > 0 {
		overview.SetDownloadablePeriod(downloadablePeriods[0]) // Use first period
	}

	// Parse CardSlotsStatusRecordArray
	cardSlotsStatuses, offset, err := parseCardSlotsStatusRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card slots status record array: %w", err)
	}
	if len(cardSlotsStatuses) > 0 {
		overview.SetDriverSlotCard(cardSlotsStatuses[0].DriverSlotCardStatus)
		overview.SetCoDriverSlotCard(cardSlotsStatuses[0].CodriverSlotCardStatus)
	}

	// Parse VuDownloadActivityDataRecordArray
	downloadActivityData, offset, err := parseVuDownloadActivityDataRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse VU download activity data record array: %w", err)
	}
	overview.SetDownloadActivities(downloadActivityData)

	// Parse VuCompanyLocksRecordArray
	companyLocks, offset, err := parseVuCompanyLocksRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse VU company locks record array: %w", err)
	}
	overview.SetCompanyLocks(companyLocks)

	// Parse VuControlActivityRecordArray
	controlActivities, offset, err := parseVuControlActivityRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse VU control activity record array: %w", err)
	}
	overview.SetControlActivities(controlActivities)

	// Parse SignatureRecordArray
	signature, offset, err := parseSignatureRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signature record array: %w", err)
	}
	overview.SetSignatureGen2(signature)

	return offset - startOffset, nil
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
//
//nolint:unused
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

//nolint:unused
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
		companyCardNumberAndGeneration := lock.GetCompanyCardNumberAndGeneration()
		var companyCardNumber *ddv1.FullCardNumber
		if companyCardNumberAndGeneration != nil {
			companyCardNumber = companyCardNumberAndGeneration.GetFullCardNumber()
		}
		buf.Write(appendVuFullCardNumber(nil, companyCardNumber, 16)) // Card number field
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
		controlCardNumberAndGeneration := control.GetControlCardNumberAndGeneration()
		var controlCardNumber *ddv1.FullCardNumber
		if controlCardNumberAndGeneration != nil {
			controlCardNumber = controlCardNumberAndGeneration.GetFullCardNumber()
		}
		buf.Write(appendVuFullCardNumber(nil, controlCardNumber, 16))
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

//nolint:unused
func appendOverviewGen2(buf *bytes.Buffer, overview *vuv1.Overview) {
	// Gen2 Overview structure - uses record arrays instead of fixed fields
	// VuOverviewSecondGen ::= SEQUENCE {
	//     memberStateCertificateRecordArray    MemberStateCertificateRecordArray,
	//     vuCertificateRecordArray             VuCertificateRecordArray,
	//     vehicleIdentificationNumberRecordArray VehicleIdentificationNumberRecordArray,
	//     vehicleRegistrationIdentificationRecordArray VehicleRegistrationIdentificationRecordArray,
	//     currentDateTimeRecordArray           CurrentDateTimeRecordArray,
	//     vuDownloadablePeriodRecordArray      VuDownloadablePeriodRecordArray,
	//     cardSlotsStatusRecordArray           CardSlotsStatusRecordArray,
	//     vuDownloadActivityDataRecordArray    VuDownloadActivityDataRecordArray,
	//     vuCompanyLocksRecordArray            VuCompanyLocksRecordArray,
	//     vuControlActivityRecordArray         VuControlActivityRecordArray,
	//     signatureRecordArray                 SignatureRecordArray
	// }

	// Append MemberStateCertificateRecordArray
	memberStateCert := overview.GetMemberStateCertificate()
	if len(memberStateCert) > 0 {
		appendMemberStateCertificateRecordArray(buf, [][]byte{memberStateCert})
	}

	// Append VuCertificateRecordArray
	vuCert := overview.GetVuCertificate()
	if len(vuCert) > 0 {
		appendVuCertificateRecordArray(buf, [][]byte{vuCert})
	}

	// Append VehicleIdentificationNumberRecordArray
	vin := overview.GetVehicleIdentificationNumber()
	if vin != nil {
		appendVehicleIdentificationNumberRecordArray(buf, []*ddv1.StringValue{vin})
	}

	// Append VehicleRegistrationIdentificationRecordArray
	vehicleReg := overview.GetVehicleRegistrationWithNation()
	if vehicleReg != nil {
		appendVehicleRegistrationIdentificationRecordArray(buf, []*ddv1.VehicleRegistrationIdentification{vehicleReg})
	}

	// Append CurrentDateTimeRecordArray
	currentDateTime := overview.GetCurrentDateTime()
	if currentDateTime != nil {
		appendCurrentDateTimeRecordArray(buf, []*timestamppb.Timestamp{currentDateTime})
	}

	// TODO: Implement record array marshalling for:
	// - VuDownloadablePeriodRecordArray
	// - CardSlotsStatusRecordArray
	// - VuDownloadActivityDataRecordArray
	// - VuCompanyLocksRecordArray
	// - VuControlActivityRecordArray

	// Append SignatureRecordArray
	signature := overview.GetSignatureGen2()
	if len(signature) > 0 {
		buf.Write(signature)
	}
}

//nolint:unused
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
//
//nolint:unused
func appendVuBytes(dst []byte, data []byte) []byte {
	return append(dst, data...)
}

// appendVuString appends a string to dst with a fixed length, padding with null bytes
//
//nolint:unused
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
//
//nolint:unused
func appendVuFullCardNumber(dst []byte, cardNumber *ddv1.FullCardNumber, length int) []byte {
	if cardNumber == nil {
		return append(dst, make([]byte, length)...)
	}
	// TODO: Implement proper FullCardNumber serialization
	return append(dst, make([]byte, length)...)
}

// Gen2 Overview Record Array Helper Functions

// parseMemberStateCertificateRecordArray parses MemberStateCertificateRecordArray
func parseMemberStateCertificateRecordArray(data []byte, offset int) ([][]byte, int, error) {
	// For now, implement a simplified version that reads a single certificate
	// In a full implementation, this would parse the record array header and multiple certificates
	if offset+194 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for member state certificate")
	}
	cert := make([]byte, 194)
	copy(cert, data[offset:offset+194])
	return [][]byte{cert}, offset + 194, nil
}

// appendMemberStateCertificateRecordArray appends MemberStateCertificateRecordArray
//
//nolint:unused
func appendMemberStateCertificateRecordArray(buf *bytes.Buffer, certs [][]byte) {
	// For now, implement a simplified version that writes a single certificate
	// In a full implementation, this would write the record array header and multiple certificates
	if len(certs) > 0 {
		buf.Write(certs[0])
	}
}

// parseVuCertificateRecordArray parses VuCertificateRecordArray
func parseVuCertificateRecordArray(data []byte, offset int) ([][]byte, int, error) {
	// For now, implement a simplified version that reads a single certificate
	if offset+194 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for VU certificate")
	}
	cert := make([]byte, 194)
	copy(cert, data[offset:offset+194])
	return [][]byte{cert}, offset + 194, nil
}

// appendVuCertificateRecordArray appends VuCertificateRecordArray
//
//nolint:unused
func appendVuCertificateRecordArray(buf *bytes.Buffer, certs [][]byte) {
	// For now, implement a simplified version that writes a single certificate
	if len(certs) > 0 {
		buf.Write(certs[0])
	}
}

// parseVehicleIdentificationNumberRecordArray parses VehicleIdentificationNumberRecordArray
func parseVehicleIdentificationNumberRecordArray(data []byte, offset int) ([]*ddv1.StringValue, int, error) {
	// For now, implement a simplified version that reads a single VIN
	if offset+17 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for vehicle identification number")
	}
	vinBytes := data[offset : offset+17]
	vin, err := unmarshalIA5StringValue(vinBytes)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse VIN: %w", err)
	}
	return []*ddv1.StringValue{vin}, offset + 17, nil
}

// appendVehicleIdentificationNumberRecordArray appends VehicleIdentificationNumberRecordArray
//
//nolint:unused
func appendVehicleIdentificationNumberRecordArray(buf *bytes.Buffer, vins []*ddv1.StringValue) {
	// For now, implement a simplified version that writes a single VIN
	if len(vins) > 0 && vins[0] != nil {
		vinBytes := make([]byte, 17)
		copy(vinBytes, []byte(vins[0].GetDecoded()))
		buf.Write(vinBytes)
	}
}

// parseVehicleRegistrationIdentificationRecordArray parses VehicleRegistrationIdentificationRecordArray
func parseVehicleRegistrationIdentificationRecordArray(data []byte, offset int) ([]*ddv1.VehicleRegistrationIdentification, int, error) {
	// For now, implement a simplified version that reads a single registration
	if offset+15 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for vehicle registration identification")
	}
	nation := ddv1.NationNumeric(data[offset])
	regNumBytes := data[offset+1 : offset+15]
	regNum, err := unmarshalIA5StringValue(regNumBytes[1:]) // Skip codepage byte
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse registration number: %w", err)
	}
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)
	vehicleReg.SetNumber(regNum)
	return []*ddv1.VehicleRegistrationIdentification{vehicleReg}, offset + 15, nil
}

// appendVehicleRegistrationIdentificationRecordArray appends VehicleRegistrationIdentificationRecordArray
//
//nolint:unused
func appendVehicleRegistrationIdentificationRecordArray(buf *bytes.Buffer, regs []*ddv1.VehicleRegistrationIdentification) {
	// For now, implement a simplified version that writes a single registration
	if len(regs) > 0 && regs[0] != nil {
		reg := regs[0]
		buf.WriteByte(byte(reg.GetNation()))
		regBytes := make([]byte, 14)
		regBytes[0] = 0 // Codepage
		copy(regBytes[1:], []byte(reg.GetNumber().GetDecoded()))
		buf.Write(regBytes)
	}
}

// parseCurrentDateTimeRecordArray parses CurrentDateTimeRecordArray
func parseCurrentDateTimeRecordArray(data []byte, offset int) ([]*timestamppb.Timestamp, int, error) {
	// For now, implement a simplified version that reads a single timestamp
	if offset+4 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for current date time")
	}
	timeValue, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read current date time: %w", err)
	}
	timestamp := timestamppb.New(time.Unix(timeValue, 0))
	return []*timestamppb.Timestamp{timestamp}, offset, nil
}

// appendCurrentDateTimeRecordArray appends CurrentDateTimeRecordArray
//
//nolint:unused
func appendCurrentDateTimeRecordArray(buf *bytes.Buffer, timestamps []*timestamppb.Timestamp) {
	// For now, implement a simplified version that writes a single timestamp
	if len(timestamps) > 0 && timestamps[0] != nil {
		timeBytes := appendVuTimeReal(nil, timestamps[0])
		buf.Write(timeBytes)
	}
}

// parseVuDownloadablePeriodRecordArray parses VuDownloadablePeriodRecordArray
func parseVuDownloadablePeriodRecordArray(data []byte, offset int) ([]*ddv1.DownloadablePeriod, int, error) {
	// For now, implement a simplified version that reads a single period
	if offset+8 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for VU downloadable period")
	}
	// Parse downloadable period (8 bytes: 4 bytes start + 4 bytes end)
	startTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read start time: %w", err)
	}
	endTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read end time: %w", err)
	}
	period := &ddv1.DownloadablePeriod{}
	period.SetMinTime(timestamppb.New(time.Unix(startTime, 0)))
	period.SetMaxTime(timestamppb.New(time.Unix(endTime, 0)))
	return []*ddv1.DownloadablePeriod{period}, offset, nil
}

// appendVuDownloadablePeriodRecordArray appends VuDownloadablePeriodRecordArray

// CardSlotsStatus represents card slot status information
type CardSlotsStatus struct {
	DriverSlotCardStatus   ddv1.SlotCardType
	CodriverSlotCardStatus ddv1.SlotCardType
}

// parseCardSlotsStatusRecordArray parses CardSlotsStatusRecordArray
func parseCardSlotsStatusRecordArray(data []byte, offset int) ([]*CardSlotsStatus, int, error) {
	// For now, implement a simplified version that reads a single status
	if offset+2 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for card slots status")
	}
	// Parse card slots status (2 bytes)
	status := &CardSlotsStatus{
		DriverSlotCardStatus:   ddv1.SlotCardType(data[offset]),
		CodriverSlotCardStatus: ddv1.SlotCardType(data[offset+1]),
	}
	return []*CardSlotsStatus{status}, offset + 2, nil
}

// appendCardSlotsStatusRecordArray appends CardSlotsStatusRecordArray

// parseVuDownloadActivityDataRecordArray parses VuDownloadActivityDataRecordArray
func parseVuDownloadActivityDataRecordArray(data []byte, offset int) ([]*vuv1.Overview_DownloadActivity, int, error) {
	// For now, implement a simplified version that reads a single data record
	// This is a placeholder - the actual structure would need to be defined
	return []*vuv1.Overview_DownloadActivity{}, offset, nil
}

// appendVuDownloadActivityDataRecordArray appends VuDownloadActivityDataRecordArray

// parseVuCompanyLocksRecordArray parses VuCompanyLocksRecordArray
func parseVuCompanyLocksRecordArray(data []byte, offset int) ([]*vuv1.Overview_CompanyLock, int, error) {
	// For now, implement a simplified version that reads a single locks record
	// This is a placeholder - the actual structure would need to be defined
	return []*vuv1.Overview_CompanyLock{}, offset, nil
}

// appendVuCompanyLocksRecordArray appends VuCompanyLocksRecordArray

// parseVuControlActivityRecordArray parses VuControlActivityRecordArray
func parseVuControlActivityRecordArray(data []byte, offset int) ([]*vuv1.Overview_ControlActivity, int, error) {
	// For now, implement a simplified version that reads a single control activity record
	// This is a placeholder - the actual structure would need to be defined
	return []*vuv1.Overview_ControlActivity{}, offset, nil
}

// appendVuControlActivityRecordArray appends VuControlActivityRecordArray
