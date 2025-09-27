package tachograph

import (
	"bytes"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalOverview parses VU overview data for different generations
func UnmarshalOverview(r *bytes.Reader, overview *vuv1.Overview, generation int) (int, error) {
	startPos := int64(r.Len())

	switch generation {
	case 1:
		overview.SetGeneration(datadictionaryv1.Generation_GENERATION_1)
		return unmarshalOverviewGen1(r, overview, startPos)
	case 2:
		overview.SetGeneration(datadictionaryv1.Generation_GENERATION_2)
		// For now, assume version 1 - we can add version detection later
		overview.SetVersion(vuv1.Version_VERSION_1)
		return unmarshalOverviewGen2(r, overview, startPos)
	default:
		return 0, nil
	}
}

func unmarshalOverviewGen1(r *bytes.Reader, overview *vuv1.Overview, startPos int64) (int, error) {
	// Gen1 Overview structure based on benchmark definitions
	// See VuOverviewFirstGen in benchmark/tachoparser/pkg/decoder/definitions.go

	// MemberStateCertificate (194 bytes)
	memberStateCert, err := readBytes(r, 194)
	if err != nil {
		return 0, err
	}
	overview.SetMemberStateCertificate(memberStateCert)

	// VuCertificate (194 bytes)
	vuCert, err := readBytes(r, 194)
	if err != nil {
		return 0, err
	}
	overview.SetVuCertificate(vuCert)

	// VehicleIdentificationNumber (17 bytes)
	vinBytes, err := readBytes(r, 17)
	if err != nil {
		return 0, err
	}
	vin := readString(bytes.NewReader(vinBytes), 17)
	overview.SetVehicleIdentificationNumber(vin)

	// VehicleRegistrationIdentification (15 bytes: nation(1) + regnum(14))
	nation, err := readUint8(r)
	if err != nil {
		return 0, err
	}
	regNumBytes, err := readBytes(r, 14)
	if err != nil {
		return 0, err
	}

	// First byte is codepage, rest is registration number
	regNum := ""
	if len(regNumBytes) > 1 {
		regNum = readString(bytes.NewReader(regNumBytes[1:]), 13)
	}

	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(int32(nation))
	vehicleReg.SetNumber(regNum)
	overview.SetVehicleRegistrationWithNation(vehicleReg)

	// CurrentDateTime (4 bytes)
	currentTime, err := readVuTimeReal(r)
	if err != nil {
		return 0, err
	}
	overview.SetCurrentDateTime(timestamppb.New(time.Unix(currentTime, 0)))

	// VuDownloadablePeriod (8 bytes: 4 bytes min + 4 bytes max)
	minTime, err := readVuTimeReal(r)
	if err != nil {
		return 0, err
	}
	maxTime, err := readVuTimeReal(r)
	if err != nil {
		return 0, err
	}

	downloadablePeriod := &vuv1.Overview_DownloadablePeriod{}
	downloadablePeriod.SetMinTime(timestamppb.New(time.Unix(minTime, 0)))
	downloadablePeriod.SetMaxTime(timestamppb.New(time.Unix(maxTime, 0)))
	overview.SetDownloadablePeriod(downloadablePeriod)

	// CardSlotsStatus (1 byte - driver and co-driver slots)
	slotsStatus, err := readUint8(r)
	if err != nil {
		return 0, err
	}
	// Extract driver and co-driver slot info from the byte
	driverSlot := (slotsStatus >> 4) & 0x0F
	coDriverSlot := slotsStatus & 0x0F

	overview.SetDriverSlotCard(mapSlotCardType(driverSlot))
	overview.SetCoDriverSlotCard(mapSlotCardType(coDriverSlot))

	// VuDownloadActivityData (4 bytes - last download time)
	lastDownloadTime, err := readVuTimeReal(r)
	if err != nil {
		return 0, err
	}
	overview.SetLastDownloadTime(timestamppb.New(time.Unix(lastDownloadTime, 0)))

	// VuCompanyLocksData - variable length, need to determine size
	// For now, we'll skip this complex structure

	// VuControlActivityData - variable length
	// For now, we'll skip this complex structure

	// Signature (128 bytes for Gen1)
	signature, err := readBytes(r, 128)
	if err != nil {
		return 0, err
	}
	overview.SetSignatureGen1(signature)

	bytesRead := int(startPos - int64(r.Len()))
	return bytesRead, nil
}

func unmarshalOverviewGen2(r *bytes.Reader, overview *vuv1.Overview, startPos int64) (int, error) {
	// Gen2 Overview structure - more complex, implement basic version for now
	// This would need detailed implementation based on the specific Gen2 structure

	// For now, just read a minimal amount to avoid errors
	// In a full implementation, this would parse the complete Gen2 structure

	bytesRead := int(startPos - int64(r.Len()))
	return bytesRead, nil
}

func mapSlotCardType(slotValue uint8) vuv1.Overview_SlotCardType {
	switch slotValue {
	case 0:
		return vuv1.Overview_NO_CARD
	case 1:
		return vuv1.Overview_DRIVER_CARD
	case 2:
		return vuv1.Overview_WORKSHOP_CARD
	case 3:
		return vuv1.Overview_CONTROL_CARD
	case 4:
		return vuv1.Overview_COMPANY_CARD
	default:
		return vuv1.Overview_SLOT_CARD_TYPE_UNSPECIFIED
	}
}
