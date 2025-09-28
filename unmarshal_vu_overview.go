package tachograph

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalOverview parses VU overview data for different generations
//
// ASN.1 Specification (Data Dictionary 2.2.6):
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
// Binary Layout (Gen1):
//
//	0-193:    memberStateCertificate (194 bytes)
//	194-387:  vuCertificate (194 bytes)
//	388-404:  vehicleIdentificationNumber (17 bytes)
//	405-419:  vehicleRegistrationIdentification (15 bytes: 1 byte nation + 14 bytes number)
//	420-423:  currentDateTime (4 bytes, TimeReal)
//	424-431:  vuDownloadablePeriod (8 bytes: 4 bytes min + 4 bytes max)
//	432-432:  cardSlotsStatus (1 byte)
//	433-451:  vuDownloadActivityData (19 bytes: 4 bytes time + 18 bytes card + 35 bytes name)
//	452+:     vuCompanyLocksData (variable size)
//	...:      vuControlActivityData (variable size)
//	...:      signature (128 bytes)
//
// Constants:
const (
	// VuOverviewFirstGen fixed fields size
	vuOverviewFirstGenFixedSize = 2*194 + 17 + 15 + 4 + 8 + 1 + 19 // 2*194 + 17 + 15 + 4 + 8 + 1 + 19 = 456 bytes

	// Certificate sizes
	memberStateCertificateSize = 194
	vuCertificateSize          = 194

	// Vehicle identification
	vehicleIdentificationNumberSize       = 17
	vehicleRegistrationIdentificationSize = 15

	// Time fields
	vuDownloadablePeriodSize = 8

	// Status and activity
	cardSlotsStatusSize                = 1
	vuDownloadActivityDataFirstGenSize = 19
)

func UnmarshalOverview(data []byte, offset int, overview *vuv1.Overview, generation int) (int, error) {
	startOffset := offset

	switch generation {
	case 1:
		overview.SetGeneration(datadictionaryv1.Generation_GENERATION_1)
		return unmarshalOverviewGen1(data, offset, overview, startOffset)
	case 2:
		overview.SetGeneration(datadictionaryv1.Generation_GENERATION_2)
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
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(datadictionaryv1.NationNumeric(nation))

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

	downloadablePeriod := &datadictionaryv1.DownloadablePeriod{}
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

func mapSlotCardType(slotValue uint8) datadictionaryv1.SlotCardType {
	switch slotValue {
	case 0:
		return datadictionaryv1.SlotCardType_NO_CARD
	case 1:
		return datadictionaryv1.SlotCardType_DRIVER_CARD_INSERTED
	case 2:
		return datadictionaryv1.SlotCardType_WORKSHOP_CARD_INSERTED
	case 3:
		return datadictionaryv1.SlotCardType_CONTROL_CARD_INSERTED
	case 4:
		return datadictionaryv1.SlotCardType_COMPANY_CARD_INSERTED
	default:
		return datadictionaryv1.SlotCardType_SLOT_CARD_TYPE_UNSPECIFIED
	}
}
