package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalOverviewGen2V1 parses Gen2 V1 Overview data from the complete transfer value.
//
// Gen2 V1 Overview structure uses RecordArray format:
//
//	VuOverviewSecondGen ::= SEQUENCE {
//	    memberStateCertificateRecordArray                MemberStateCertificateRecordArray,
//	    vuCertificateRecordArray                         VuCertificateRecordArray,
//	    vehicleIdentificationNumberRecordArray           VehicleIdentificationNumberRecordArray,
//	    vehicleRegistrationIdentificationRecordArray     VehicleRegistrationIdentificationRecordArray,
//	    currentDateTimeRecordArray                       CurrentDateTimeRecordArray,
//	    vuDownloadablePeriodRecordArray                  VuDownloadablePeriodRecordArray,
//	    cardSlotsStatusRecordArray                       CardSlotsStatusRecordArray,
//	    vuDownloadActivityDataRecordArray                VuDownloadActivityDataRecordArray,
//	    vuCompanyLocksRecordArray                        VuCompanyLocksRecordArray,
//	    vuControlActivityRecordArray                     VuControlActivityRecordArray,
//	    signatureRecordArray                             SignatureRecordArray
//	}
func unmarshalOverviewGen2V1(value []byte) (*vuv1.OverviewGen2V1, error) {
	totalSize, signatureSize, err := sizeOfOverviewGen2V1(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	overview := &vuv1.OverviewGen2V1{}
	overview.SetRawData(value)
	offset := 0

	// MemberStateCertificateRecordArray
	msc, bytesRead, err := parseCertificateRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse MemberStateCertificateRecordArray: %w", err)
	}
	overview.SetMemberStateCertificate(msc)
	offset += bytesRead

	// VuCertificateRecordArray
	vuCert, bytesRead, err := parseCertificateRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuCertificateRecordArray: %w", err)
	}
	overview.SetVuCertificate(vuCert)
	offset += bytesRead

	// VehicleIdentificationNumberRecordArray
	vin, bytesRead, err := parseVehicleIdentificationNumberRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VehicleIdentificationNumberRecordArray: %w", err)
	}
	overview.SetVehicleIdentificationNumber(vin)
	offset += bytesRead

	// VehicleRegistrationIdentificationRecordArray
	vrn, bytesRead, err := parseVehicleRegistrationIdentificationRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VehicleRegistrationIdentificationRecordArray: %w", err)
	}
	overview.SetVehicleRegistrationWithNation(vrn)
	offset += bytesRead

	// CurrentDateTimeRecordArray
	currentDateTime, bytesRead, err := parseCurrentDateTimeRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse CurrentDateTimeRecordArray: %w", err)
	}
	overview.SetCurrentDateTime(currentDateTime)
	offset += bytesRead

	// VuDownloadablePeriodRecordArray
	downloadablePeriod, bytesRead, err := parseDownloadablePeriodRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuDownloadablePeriodRecordArray: %w", err)
	}
	overview.SetDownloadablePeriod(downloadablePeriod)
	offset += bytesRead

	// CardSlotsStatusRecordArray
	driverSlot, coDriverSlot, bytesRead, err := parseCardSlotsStatusRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse CardSlotsStatusRecordArray: %w", err)
	}
	overview.SetDriverSlotCard(driverSlot)
	overview.SetCoDriverSlotCard(coDriverSlot)
	offset += bytesRead

	// VuDownloadActivityDataRecordArray
	activities, bytesRead, err := parseDownloadActivityDataRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuDownloadActivityDataRecordArray: %w", err)
	}
	overview.SetDownloadActivities(activities)
	offset += bytesRead

	// VuCompanyLocksRecordArray
	companyLocks, bytesRead, err := parseCompanyLocksRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuCompanyLocksRecordArray: %w", err)
	}
	overview.SetCompanyLocks(companyLocks)
	offset += bytesRead

	// VuControlActivityRecordArray
	controlActivities, bytesRead, err := parseControlActivityRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuControlActivityRecordArray: %w", err)
	}
	overview.SetControlActivities(controlActivities)
	offset += bytesRead

	overview.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Overview Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return overview, nil
}

// MarshalOverviewGen2V1 marshals Gen2 V1 Overview data.
func (opts MarshalOptions) MarshalOverviewGen2V1(overview *vuv1.OverviewGen2V1) ([]byte, error) {
	if overview == nil {
		return nil, fmt.Errorf("overview cannot be nil")
	}

	raw := overview.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	marshalOpts := dd.MarshalOptions{}
	var result []byte

	// MemberStateCertificateRecordArray
	msc := overview.GetMemberStateCertificate()
	result = appendCertificateRecordArray(result, 0x01, msc)

	// VuCertificateRecordArray
	vuCert := overview.GetVuCertificate()
	result = appendCertificateRecordArray(result, 0x02, vuCert)

	// VehicleIdentificationNumberRecordArray (17 bytes)
	vinData, err := marshalOpts.MarshalIa5StringValue(overview.GetVehicleIdentificationNumber())
	if err != nil {
		return nil, fmt.Errorf("marshal VIN: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, uint16(len(vinData)), 1)
	result = append(result, vinData...)

	// VehicleRegistrationIdentificationRecordArray (15 bytes)
	vrnData, err := marshalOpts.MarshalVehicleRegistrationIdentification(overview.GetVehicleRegistrationWithNation())
	if err != nil {
		return nil, fmt.Errorf("marshal VRN: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x04, uint16(len(vrnData)), 1)
	result = append(result, vrnData...)

	// CurrentDateTimeRecordArray (4 bytes)
	currentDateTimeData, err := marshalOpts.MarshalTimeReal(overview.GetCurrentDateTime())
	if err != nil {
		return nil, fmt.Errorf("marshal CurrentDateTime: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x05, 4, 1)
	result = append(result, currentDateTimeData...)

	// VuDownloadablePeriodRecordArray (8 bytes: 2 × TimeReal)
	downloadablePeriodData, err := marshalDownloadablePeriodData(marshalOpts, overview.GetDownloadablePeriod())
	if err != nil {
		return nil, fmt.Errorf("marshal DownloadablePeriod: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x06, 8, 1)
	result = append(result, downloadablePeriodData...)

	// CardSlotsStatusRecordArray (1 byte: packed nibble)
	cardSlotsData, err := marshalCardSlotsStatusByte(overview.GetDriverSlotCard(), overview.GetCoDriverSlotCard())
	if err != nil {
		return nil, fmt.Errorf("marshal CardSlotsStatus: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x07, 1, 1)
	result = append(result, cardSlotsData...)

	// VuDownloadActivityDataRecordArray (59 bytes per record: 4+19+36)
	activitiesData, err := marshalDownloadActivitiesGen2V1(marshalOpts, overview.GetDownloadActivities())
	if err != nil {
		return nil, fmt.Errorf("marshal VuDownloadActivityDataRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x08, 59, uint16(len(overview.GetDownloadActivities())))
	result = append(result, activitiesData...)

	// VuCompanyLocksRecordArray (99 bytes per record: 4+4+36+36+19)
	locksData, err := marshalCompanyLocksGen2V1(marshalOpts, overview.GetCompanyLocks())
	if err != nil {
		return nil, fmt.Errorf("marshal VuCompanyLocksRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x09, 99, uint16(len(overview.GetCompanyLocks())))
	result = append(result, locksData...)

	// VuControlActivityRecordArray (32 bytes per record: 1+4+19+4+4)
	controlsData, err := marshalControlActivitiesGen2V1(marshalOpts, overview.GetControlActivities())
	if err != nil {
		return nil, fmt.Errorf("marshal VuControlActivityRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x0A, 32, uint16(len(overview.GetControlActivities())))
	result = append(result, controlsData...)

	result = append(result, overview.GetSignature()...)
	return result, nil
}

// anonymizeOverviewGen2V1 anonymizes Gen2 V1 Overview data.
func (opts AnonymizeOptions) anonymizeOverviewGen2V1(overview *vuv1.OverviewGen2V1) *vuv1.OverviewGen2V1 {
	if overview == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &vuv1.OverviewGen2V1{}

	// Certificates are cleared during anonymization (contain device identity)
	result.SetMemberStateCertificate(nil)
	result.SetVuCertificate(nil)

	// Anonymize VIN
	result.SetVehicleIdentificationNumber(ddOpts.AnonymizeIa5StringValue(overview.GetVehicleIdentificationNumber()))

	// Anonymize VRN
	result.SetVehicleRegistrationWithNation(ddOpts.AnonymizeVehicleRegistrationIdentification(overview.GetVehicleRegistrationWithNation()))

	// Preserve timestamps and period (no PII)
	result.SetCurrentDateTime(overview.GetCurrentDateTime())
	result.SetDownloadablePeriod(overview.GetDownloadablePeriod())

	// Preserve card slot status (no PII)
	result.SetDriverSlotCard(overview.GetDriverSlotCard())
	result.SetCoDriverSlotCard(overview.GetCoDriverSlotCard())

	// Anonymize download activities
	anonActivities := make([]*vuv1.OverviewGen2V1_DownloadActivity, len(overview.GetDownloadActivities()))
	for i, a := range overview.GetDownloadActivities() {
		anon := &vuv1.OverviewGen2V1_DownloadActivity{}
		anon.SetDownloadingTime(a.GetDownloadingTime())
		anon.SetFullCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(a.GetFullCardNumberAndGeneration()))
		anon.SetCompanyOrWorkshopName(ddOpts.AnonymizeStringValue(a.GetCompanyOrWorkshopName()))
		anonActivities[i] = anon
	}
	result.SetDownloadActivities(anonActivities)

	// Anonymize company locks
	anonLocks := make([]*vuv1.OverviewGen2V1_CompanyLock, len(overview.GetCompanyLocks()))
	for i, lock := range overview.GetCompanyLocks() {
		anon := &vuv1.OverviewGen2V1_CompanyLock{}
		anon.SetLockInTime(lock.GetLockInTime())
		anon.SetLockOutTime(lock.GetLockOutTime())
		anon.SetCompanyName(ddOpts.AnonymizeStringValue(lock.GetCompanyName()))
		anon.SetCompanyAddress(ddOpts.AnonymizeStringValue(lock.GetCompanyAddress()))
		anon.SetCompanyCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(lock.GetCompanyCardNumberAndGeneration()))
		anonLocks[i] = anon
	}
	result.SetCompanyLocks(anonLocks)

	// Anonymize control activities
	anonControls := make([]*vuv1.OverviewGen2V1_ControlActivity, len(overview.GetControlActivities()))
	for i, ctrl := range overview.GetControlActivities() {
		anon := &vuv1.OverviewGen2V1_ControlActivity{}
		anon.SetControlType(ctrl.GetControlType())
		anon.SetControlTime(ctrl.GetControlTime())
		anon.SetControlCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(ctrl.GetControlCardNumberAndGeneration()))
		anon.SetDownloadPeriodBeginTime(ctrl.GetDownloadPeriodBeginTime())
		anon.SetDownloadPeriodEndTime(ctrl.GetDownloadPeriodEndTime())
		anonControls[i] = anon
	}
	result.SetControlActivities(anonControls)

	result.SetSignature([]byte{})
	return result
}

// ===== Shared parse helpers (reused by overview_gen2_v2.go) =====

// parseCertificateRecordArray parses a certificate RecordArray (MemberState or VuCertificate).
// Returns the raw certificate bytes.
func parseCertificateRecordArray(data []byte, offset int) ([]byte, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}
	if noOfRecords == 0 {
		return nil, headerSize, nil
	}
	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 certificate record, got %d", noOfRecords)
	}
	recordStart := offset + headerSize
	recordEnd := recordStart + int(recordSize)
	if recordEnd > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for certificate record")
	}
	cert := make([]byte, recordSize)
	copy(cert, data[recordStart:recordEnd])
	totalSize := headerSize + int(recordSize)
	return cert, totalSize, nil
}

// parseVehicleIdentificationNumberRecordArray parses a VehicleIdentificationNumberRecordArray.
// Expected: 1 record × 17 bytes (IA5String).
func parseVehicleIdentificationNumberRecordArray(data []byte, offset int) (*ddv1.Ia5StringValue, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}
	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 VIN record, got %d", noOfRecords)
	}
	recordStart := offset + headerSize
	recordEnd := recordStart + int(recordSize)
	if recordEnd > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for VIN record")
	}
	var unmarshalOpts dd.UnmarshalOptions
	vin, err := unmarshalOpts.UnmarshalIa5StringValue(data[recordStart:recordEnd])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal VIN: %w", err)
	}
	totalSize := headerSize + int(recordSize)
	return vin, totalSize, nil
}

// parseVehicleRegistrationIdentificationRecordArray parses a VehicleRegistrationIdentificationRecordArray.
// Expected: 1 record × 15 bytes (1 nation + 14 StringValue).
func parseVehicleRegistrationIdentificationRecordArray(data []byte, offset int) (*ddv1.VehicleRegistrationIdentification, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}
	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 VRI record, got %d", noOfRecords)
	}
	recordStart := offset + headerSize
	recordEnd := recordStart + int(recordSize)
	if recordEnd > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for VRI record")
	}
	var unmarshalOpts dd.UnmarshalOptions
	vrn, err := unmarshalOpts.UnmarshalVehicleRegistrationIdentification(data[recordStart:recordEnd])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal VRI: %w", err)
	}
	totalSize := headerSize + int(recordSize)
	return vrn, totalSize, nil
}

// parseCurrentDateTimeRecordArray parses a CurrentDateTimeRecordArray.
// Expected: 1 record × 4 bytes (TimeReal).
func parseCurrentDateTimeRecordArray(data []byte, offset int) (*timestamppb.Timestamp, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}
	if noOfRecords != 1 || recordSize != 4 {
		return nil, 0, fmt.Errorf("expected CurrentDateTime: 1 record × 4 bytes, got %d records × %d bytes", noOfRecords, recordSize)
	}
	recordStart := offset + headerSize
	if recordStart+4 > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for CurrentDateTime")
	}
	var unmarshalOpts dd.UnmarshalOptions
	ts, err := unmarshalOpts.UnmarshalTimeReal(data[recordStart : recordStart+4])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal CurrentDateTime: %w", err)
	}
	totalSize := headerSize + 4
	return ts, totalSize, nil
}

// parseDownloadablePeriodRecordArray parses a VuDownloadablePeriodRecordArray.
// Expected: 1 record × 8 bytes (2 × TimeReal).
func parseDownloadablePeriodRecordArray(data []byte, offset int) (*ddv1.DownloadablePeriod, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}
	if noOfRecords != 1 || recordSize != 8 {
		return nil, 0, fmt.Errorf("expected DownloadablePeriod: 1 record × 8 bytes, got %d records × %d bytes", noOfRecords, recordSize)
	}
	recordStart := offset + headerSize
	if recordStart+8 > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for DownloadablePeriod")
	}
	var unmarshalOpts dd.UnmarshalOptions
	minTime, err := unmarshalOpts.UnmarshalTimeReal(data[recordStart : recordStart+4])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal DownloadablePeriod minTime: %w", err)
	}
	maxTime, err := unmarshalOpts.UnmarshalTimeReal(data[recordStart+4 : recordStart+8])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal DownloadablePeriod maxTime: %w", err)
	}
	dp := &ddv1.DownloadablePeriod{}
	dp.SetMinTime(minTime)
	dp.SetMaxTime(maxTime)
	totalSize := headerSize + 8
	return dp, totalSize, nil
}

// parseCardSlotsStatusRecordArray parses a CardSlotsStatusRecordArray.
// Expected: 1 record × 1 byte (packed nibble: upper=coDriver, lower=driver).
func parseCardSlotsStatusRecordArray(data []byte, offset int) (driver, coDriver ddv1.SlotCardType, bytesRead int, err error) {
	_, recordSize, noOfRecords, headerSize, parseErr := parseRecordArrayHeader(data, offset)
	if parseErr != nil {
		return 0, 0, 0, parseErr
	}
	if noOfRecords != 1 || recordSize != 1 {
		return 0, 0, 0, fmt.Errorf("expected CardSlotsStatus: 1 record × 1 byte, got %d records × %d bytes", noOfRecords, recordSize)
	}
	recordStart := offset + headerSize
	if recordStart+1 > len(data) {
		return 0, 0, 0, fmt.Errorf("insufficient data for CardSlotsStatus")
	}
	b := data[recordStart]
	driverRaw := b & 0x0F
	coDriverRaw := (b >> 4) & 0x0F

	driverSlot, err := dd.UnmarshalEnum[ddv1.SlotCardType](driverRaw)
	if err != nil {
		driverSlot = ddv1.SlotCardType_SLOT_CARD_TYPE_UNRECOGNIZED
	}
	coDriverSlot, err := dd.UnmarshalEnum[ddv1.SlotCardType](coDriverRaw)
	if err != nil {
		coDriverSlot = ddv1.SlotCardType_SLOT_CARD_TYPE_UNRECOGNIZED
	}

	totalSize := headerSize + 1
	return driverSlot, coDriverSlot, totalSize, nil
}

// parseDownloadActivityDataRecordArrayGen2V1 parses a VuDownloadActivityDataRecordArray for V1.
//
// VuDownloadActivityData (Gen2) layout:
//
//	downloadingTime (4) + fullCardNumberAndGeneration (19) + companyOrWorkshopName (36) = 59 bytes
func parseDownloadActivityDataRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.OverviewGen2V1_DownloadActivity, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedSize = 4 + 19 + 36 // 59
	if recordSize != 0 && int(recordSize) < expectedSize {
		return nil, 0, fmt.Errorf("VuDownloadActivityData size %d too small (need at least %d)", recordSize, expectedSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.OverviewGen2V1_DownloadActivity, 0, noOfRecords)
	recordStart := offset + headerSize

	for range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuDownloadActivityData")
		}
		rec := data[recordStart:recordEnd]

		activity := &vuv1.OverviewGen2V1_DownloadActivity{}

		downloadingTime, err := unmarshalOpts.UnmarshalTimeReal(rec[0:4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuDownloadActivityData downloadingTime: %w", err)
		}
		activity.SetDownloadingTime(downloadingTime)

		cardNum, err := unmarshalOpts.UnmarshalFullCardNumberAndGeneration(rec[4:23])
		if err != nil {
			return nil, 0, fmt.Errorf("VuDownloadActivityData card number: %w", err)
		}
		activity.SetFullCardNumberAndGeneration(cardNum)

		name, err := unmarshalOpts.UnmarshalStringValue(rec[23:59])
		if err != nil {
			return nil, 0, fmt.Errorf("VuDownloadActivityData company name: %w", err)
		}
		activity.SetCompanyOrWorkshopName(name)

		records = append(records, activity)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseCompanyLocksRecordArrayGen2V1 parses a VuCompanyLocksRecordArray for V1.
//
// VuCompanyLocksRecord (Gen2) layout:
//
//	lockInTime (4) + lockOutTime (4) + companyName (36) + companyAddress (36) +
//	companyCardNumberAndGeneration (19) = 99 bytes
func parseCompanyLocksRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.OverviewGen2V1_CompanyLock, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedSize = 4 + 4 + 36 + 36 + 19 // 99
	if recordSize != 0 && int(recordSize) < expectedSize {
		return nil, 0, fmt.Errorf("VuCompanyLocksRecord size %d too small (need at least %d)", recordSize, expectedSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.OverviewGen2V1_CompanyLock, 0, noOfRecords)
	recordStart := offset + headerSize

	for range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuCompanyLocksRecord")
		}
		rec := data[recordStart:recordEnd]

		lock := &vuv1.OverviewGen2V1_CompanyLock{}

		lockInTime, err := unmarshalOpts.UnmarshalTimeReal(rec[0:4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuCompanyLocksRecord lockInTime: %w", err)
		}
		lock.SetLockInTime(lockInTime)

		lockOutTime, err := unmarshalOpts.UnmarshalTimeReal(rec[4:8])
		if err != nil {
			return nil, 0, fmt.Errorf("VuCompanyLocksRecord lockOutTime: %w", err)
		}
		lock.SetLockOutTime(lockOutTime)

		companyName, err := unmarshalOpts.UnmarshalStringValue(rec[8:44])
		if err != nil {
			return nil, 0, fmt.Errorf("VuCompanyLocksRecord companyName: %w", err)
		}
		lock.SetCompanyName(companyName)

		companyAddress, err := unmarshalOpts.UnmarshalStringValue(rec[44:80])
		if err != nil {
			return nil, 0, fmt.Errorf("VuCompanyLocksRecord companyAddress: %w", err)
		}
		lock.SetCompanyAddress(companyAddress)

		cardNum, err := unmarshalOpts.UnmarshalFullCardNumberAndGeneration(rec[80:99])
		if err != nil {
			return nil, 0, fmt.Errorf("VuCompanyLocksRecord cardNumber: %w", err)
		}
		lock.SetCompanyCardNumberAndGeneration(cardNum)

		records = append(records, lock)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseControlActivityRecordArrayGen2V1 parses a VuControlActivityRecordArray for V1.
//
// VuControlActivityRecord (Gen2) layout:
//
//	controlType (1) + controlTime (4) + controlCardNumberAndGeneration (19) +
//	downloadPeriodBeginTime (4) + downloadPeriodEndTime (4) = 32 bytes
func parseControlActivityRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.OverviewGen2V1_ControlActivity, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedSize = 1 + 4 + 19 + 4 + 4 // 32
	if recordSize != 0 && int(recordSize) < expectedSize {
		return nil, 0, fmt.Errorf("VuControlActivityRecord size %d too small (need at least %d)", recordSize, expectedSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.OverviewGen2V1_ControlActivity, 0, noOfRecords)
	recordStart := offset + headerSize

	for range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuControlActivityRecord")
		}
		rec := data[recordStart:recordEnd]

		ctrl := &vuv1.OverviewGen2V1_ControlActivity{}

		controlType, err := unmarshalOpts.UnmarshalControlType(rec[0:1])
		if err != nil {
			return nil, 0, fmt.Errorf("VuControlActivityRecord controlType: %w", err)
		}
		ctrl.SetControlType(controlType)

		controlTime, err := unmarshalOpts.UnmarshalTimeReal(rec[1:5])
		if err != nil {
			return nil, 0, fmt.Errorf("VuControlActivityRecord controlTime: %w", err)
		}
		ctrl.SetControlTime(controlTime)

		cardNum, err := unmarshalOpts.UnmarshalFullCardNumberAndGeneration(rec[5:24])
		if err != nil {
			return nil, 0, fmt.Errorf("VuControlActivityRecord card number: %w", err)
		}
		ctrl.SetControlCardNumberAndGeneration(cardNum)

		beginTime, err := unmarshalOpts.UnmarshalTimeReal(rec[24:28])
		if err != nil {
			return nil, 0, fmt.Errorf("VuControlActivityRecord downloadPeriodBeginTime: %w", err)
		}
		ctrl.SetDownloadPeriodBeginTime(beginTime)

		endTime, err := unmarshalOpts.UnmarshalTimeReal(rec[28:32])
		if err != nil {
			return nil, 0, fmt.Errorf("VuControlActivityRecord downloadPeriodEndTime: %w", err)
		}
		ctrl.SetDownloadPeriodEndTime(endTime)

		records = append(records, ctrl)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// ===== Shared marshal helpers (reused by overview_gen2_v2.go) =====

// appendCertificateRecordArray appends a certificate RecordArray to dst.
func appendCertificateRecordArray(dst []byte, recordType byte, cert []byte) []byte {
	if len(cert) == 0 {
		dst = appendRecordArrayHeader(dst, recordType, 0, 0)
		return dst
	}
	dst = appendRecordArrayHeader(dst, recordType, uint16(len(cert)), 1)
	dst = append(dst, cert...)
	return dst
}

// marshalDownloadablePeriodData marshals a DownloadablePeriod to 8 bytes.
func marshalDownloadablePeriodData(opts dd.MarshalOptions, dp *ddv1.DownloadablePeriod) ([]byte, error) {
	var minTime, maxTime *timestamppb.Timestamp
	if dp != nil {
		minTime = dp.GetMinTime()
		maxTime = dp.GetMaxTime()
	}
	minBytes, err := opts.MarshalTimeReal(minTime)
	if err != nil {
		return nil, fmt.Errorf("marshal minTime: %w", err)
	}
	maxBytes, err := opts.MarshalTimeReal(maxTime)
	if err != nil {
		return nil, fmt.Errorf("marshal maxTime: %w", err)
	}
	return append(minBytes, maxBytes...), nil
}

// marshalCardSlotsStatusByte marshals driver and co-driver slot types to 1 packed byte.
func marshalCardSlotsStatusByte(driver, coDriver ddv1.SlotCardType) ([]byte, error) {
	driverByte, err := dd.MarshalEnum(driver)
	if err != nil {
		driverByte = 0
	}
	coDriverByte, err := dd.MarshalEnum(coDriver)
	if err != nil {
		coDriverByte = 0
	}
	return []byte{(coDriverByte << 4) | (driverByte & 0x0F)}, nil
}

// marshalDownloadActivitiesGen2V1 marshals download activity records for V1.
func marshalDownloadActivitiesGen2V1(opts dd.MarshalOptions, activities []*vuv1.OverviewGen2V1_DownloadActivity) ([]byte, error) {
	result := make([]byte, 0, len(activities)*59)
	for i, a := range activities {
		timeBytes, err := opts.MarshalTimeReal(a.GetDownloadingTime())
		if err != nil {
			return nil, fmt.Errorf("DownloadActivity %d downloadingTime: %w", i, err)
		}
		result = append(result, timeBytes...)

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(a.GetFullCardNumberAndGeneration())
		if err != nil {
			return nil, fmt.Errorf("DownloadActivity %d card number: %w", i, err)
		}
		if len(cardBytes) != 19 {
			return nil, fmt.Errorf("DownloadActivity %d card number length: got %d, want 19", i, len(cardBytes))
		}
		result = append(result, cardBytes...)

		nameBytes, err := opts.MarshalStringValue(a.GetCompanyOrWorkshopName())
		if err != nil {
			return nil, fmt.Errorf("DownloadActivity %d company name: %w", i, err)
		}
		if len(nameBytes) != 36 {
			return nil, fmt.Errorf("DownloadActivity %d company name length: got %d, want 36", i, len(nameBytes))
		}
		result = append(result, nameBytes...)
	}
	return result, nil
}

// marshalCompanyLocksGen2V1 marshals company lock records for V1.
func marshalCompanyLocksGen2V1(opts dd.MarshalOptions, locks []*vuv1.OverviewGen2V1_CompanyLock) ([]byte, error) {
	result := make([]byte, 0, len(locks)*99)
	for i, lock := range locks {
		lockInBytes, err := opts.MarshalTimeReal(lock.GetLockInTime())
		if err != nil {
			return nil, fmt.Errorf("CompanyLock %d lockInTime: %w", i, err)
		}
		result = append(result, lockInBytes...)

		lockOutBytes, err := opts.MarshalTimeReal(lock.GetLockOutTime())
		if err != nil {
			return nil, fmt.Errorf("CompanyLock %d lockOutTime: %w", i, err)
		}
		result = append(result, lockOutBytes...)

		nameBytes, err := opts.MarshalStringValue(lock.GetCompanyName())
		if err != nil {
			return nil, fmt.Errorf("CompanyLock %d companyName: %w", i, err)
		}
		if len(nameBytes) != 36 {
			return nil, fmt.Errorf("CompanyLock %d company name length: got %d, want 36", i, len(nameBytes))
		}
		result = append(result, nameBytes...)

		addrBytes, err := opts.MarshalStringValue(lock.GetCompanyAddress())
		if err != nil {
			return nil, fmt.Errorf("CompanyLock %d companyAddress: %w", i, err)
		}
		if len(addrBytes) != 36 {
			return nil, fmt.Errorf("CompanyLock %d company address length: got %d, want 36", i, len(addrBytes))
		}
		result = append(result, addrBytes...)

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(lock.GetCompanyCardNumberAndGeneration())
		if err != nil {
			return nil, fmt.Errorf("CompanyLock %d card number: %w", i, err)
		}
		if len(cardBytes) != 19 {
			return nil, fmt.Errorf("CompanyLock %d card number length: got %d, want 19", i, len(cardBytes))
		}
		result = append(result, cardBytes...)
	}
	return result, nil
}

// marshalControlActivitiesGen2V1 marshals control activity records for V1.
func marshalControlActivitiesGen2V1(opts dd.MarshalOptions, controls []*vuv1.OverviewGen2V1_ControlActivity) ([]byte, error) {
	result := make([]byte, 0, len(controls)*32)
	for i, ctrl := range controls {
		ctrlTypeBytes, err := opts.MarshalControlType(ctrl.GetControlType())
		if err != nil {
			return nil, fmt.Errorf("ControlActivity %d controlType: %w", i, err)
		}
		if len(ctrlTypeBytes) != 1 {
			return nil, fmt.Errorf("ControlActivity %d control type length: got %d, want 1", i, len(ctrlTypeBytes))
		}
		result = append(result, ctrlTypeBytes...)

		ctrlTimeBytes, err := opts.MarshalTimeReal(ctrl.GetControlTime())
		if err != nil {
			return nil, fmt.Errorf("ControlActivity %d controlTime: %w", i, err)
		}
		result = append(result, ctrlTimeBytes...)

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(ctrl.GetControlCardNumberAndGeneration())
		if err != nil {
			return nil, fmt.Errorf("ControlActivity %d card number: %w", i, err)
		}
		if len(cardBytes) != 19 {
			return nil, fmt.Errorf("ControlActivity %d card number length: got %d, want 19", i, len(cardBytes))
		}
		result = append(result, cardBytes...)

		beginBytes, err := opts.MarshalTimeReal(ctrl.GetDownloadPeriodBeginTime())
		if err != nil {
			return nil, fmt.Errorf("ControlActivity %d downloadPeriodBeginTime: %w", i, err)
		}
		result = append(result, beginBytes...)

		endBytes, err := opts.MarshalTimeReal(ctrl.GetDownloadPeriodEndTime())
		if err != nil {
			return nil, fmt.Errorf("ControlActivity %d downloadPeriodEndTime: %w", i, err)
		}
		result = append(result, endBytes...)
	}
	return result, nil
}
