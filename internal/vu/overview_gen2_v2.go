package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalOverviewGen2V2 parses Gen2 V2 Overview data from the complete transfer value.
//
// Gen2 V2 adds VehicleRegistrationNumberRecordArray at position 4 (after VIN, before CurrentDateTime).
//
// Structure:
//
//	VuOverviewSecondGenV2 ::= SEQUENCE {
//	    memberStateCertificateRecordArray                MemberStateCertificateRecordArray,
//	    vuCertificateRecordArray                         VuCertificateRecordArray,
//	    vehicleIdentificationNumberRecordArray           VehicleIdentificationNumberRecordArray,
//	    vehicleRegistrationNumberRecordArray             VehicleRegistrationNumberRecordArray,
//	    currentDateTimeRecordArray                       CurrentDateTimeRecordArray,
//	    vuDownloadablePeriodRecordArray                  VuDownloadablePeriodRecordArray,
//	    cardSlotsStatusRecordArray                       CardSlotsStatusRecordArray,
//	    vuDownloadActivityDataRecordArray                VuDownloadActivityDataRecordArray,
//	    vuCompanyLocksRecordArray                        VuCompanyLocksRecordArray,
//	    vuControlActivityRecordArray                     VuControlActivityRecordArray,
//	    signatureRecordArray                             SignatureRecordArray
//	}
func unmarshalOverviewGen2V2(value []byte) (*vuv1.OverviewGen2V2, error) {
	totalSize, signatureSize, err := sizeOfOverviewGen2V2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	overview := &vuv1.OverviewGen2V2{}
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

	// VehicleRegistrationNumberRecordArray (Gen2 V2 addition: 13-byte IA5String)
	vrn, bytesRead, err := parseVehicleRegistrationNumberRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VehicleRegistrationNumberRecordArray: %w", err)
	}
	overview.SetVehicleRegistrationNumber(vrn)
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
	activities, bytesRead, err := parseDownloadActivityDataRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuDownloadActivityDataRecordArray: %w", err)
	}
	overview.SetDownloadActivities(activities)
	offset += bytesRead

	// VuCompanyLocksRecordArray
	companyLocks, bytesRead, err := parseCompanyLocksRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuCompanyLocksRecordArray: %w", err)
	}
	overview.SetCompanyLocks(companyLocks)
	offset += bytesRead

	// VuControlActivityRecordArray
	controlActivities, bytesRead, err := parseControlActivityRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuControlActivityRecordArray: %w", err)
	}
	overview.SetControlActivities(controlActivities)
	offset += bytesRead

	overview.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Overview Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return overview, nil
}

// MarshalOverviewGen2V2 marshals Gen2 V2 Overview data.
func (opts MarshalOptions) MarshalOverviewGen2V2(overview *vuv1.OverviewGen2V2) ([]byte, error) {
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
	result = appendCertificateRecordArray(result, 0x01, overview.GetMemberStateCertificate())

	// VuCertificateRecordArray
	result = appendCertificateRecordArray(result, 0x02, overview.GetVuCertificate())

	// VehicleIdentificationNumberRecordArray (17 bytes)
	vinData, err := marshalOpts.MarshalIa5StringValue(overview.GetVehicleIdentificationNumber())
	if err != nil {
		return nil, fmt.Errorf("marshal VIN: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, uint16(len(vinData)), 1)
	result = append(result, vinData...)

	// VehicleRegistrationNumberRecordArray (13 bytes)
	vrnData, err := marshalOpts.MarshalIa5StringValue(overview.GetVehicleRegistrationNumber())
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

	// VuDownloadablePeriodRecordArray (8 bytes)
	downloadablePeriodData, err := marshalDownloadablePeriodData(marshalOpts, overview.GetDownloadablePeriod())
	if err != nil {
		return nil, fmt.Errorf("marshal DownloadablePeriod: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x06, 8, 1)
	result = append(result, downloadablePeriodData...)

	// CardSlotsStatusRecordArray (1 byte)
	cardSlotsData, err := marshalCardSlotsStatusByte(overview.GetDriverSlotCard(), overview.GetCoDriverSlotCard())
	if err != nil {
		return nil, fmt.Errorf("marshal CardSlotsStatus: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x07, 1, 1)
	result = append(result, cardSlotsData...)

	// VuDownloadActivityDataRecordArray (59 bytes per record)
	activitiesData, err := marshalDownloadActivitiesGen2V2(marshalOpts, overview.GetDownloadActivities())
	if err != nil {
		return nil, fmt.Errorf("marshal VuDownloadActivityDataRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x08, 59, uint16(len(overview.GetDownloadActivities())))
	result = append(result, activitiesData...)

	// VuCompanyLocksRecordArray (99 bytes per record)
	locksData, err := marshalCompanyLocksGen2V2(marshalOpts, overview.GetCompanyLocks())
	if err != nil {
		return nil, fmt.Errorf("marshal VuCompanyLocksRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x09, 99, uint16(len(overview.GetCompanyLocks())))
	result = append(result, locksData...)

	// VuControlActivityRecordArray (32 bytes per record)
	controlsData, err := marshalControlActivitiesGen2V2(marshalOpts, overview.GetControlActivities())
	if err != nil {
		return nil, fmt.Errorf("marshal VuControlActivityRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x0A, 32, uint16(len(overview.GetControlActivities())))
	result = append(result, controlsData...)

	result = append(result, overview.GetSignature()...)
	return result, nil
}

// anonymizeOverviewGen2V2 anonymizes Gen2 V2 Overview data.
func (opts AnonymizeOptions) anonymizeOverviewGen2V2(overview *vuv1.OverviewGen2V2) *vuv1.OverviewGen2V2 {
	if overview == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &vuv1.OverviewGen2V2{}

	// Certificates are cleared during anonymization
	result.SetMemberStateCertificate(nil)
	result.SetVuCertificate(nil)

	// Anonymize VIN and VRN
	result.SetVehicleIdentificationNumber(ddOpts.AnonymizeIa5StringValue(overview.GetVehicleIdentificationNumber()))
	result.SetVehicleRegistrationNumber(ddOpts.AnonymizeIa5StringValue(overview.GetVehicleRegistrationNumber()))

	// Preserve timestamps and period (no PII)
	result.SetCurrentDateTime(overview.GetCurrentDateTime())
	result.SetDownloadablePeriod(overview.GetDownloadablePeriod())

	// Preserve card slot status (no PII)
	result.SetDriverSlotCard(overview.GetDriverSlotCard())
	result.SetCoDriverSlotCard(overview.GetCoDriverSlotCard())

	// Anonymize download activities
	anonActivities := make([]*vuv1.OverviewGen2V2_DownloadActivity, len(overview.GetDownloadActivities()))
	for i, a := range overview.GetDownloadActivities() {
		anon := &vuv1.OverviewGen2V2_DownloadActivity{}
		anon.SetDownloadingTime(a.GetDownloadingTime())
		anon.SetFullCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(a.GetFullCardNumberAndGeneration()))
		anon.SetCompanyOrWorkshopName(ddOpts.AnonymizeStringValue(a.GetCompanyOrWorkshopName()))
		anonActivities[i] = anon
	}
	result.SetDownloadActivities(anonActivities)

	// Anonymize company locks
	anonLocks := make([]*vuv1.OverviewGen2V2_CompanyLock, len(overview.GetCompanyLocks()))
	for i, lock := range overview.GetCompanyLocks() {
		anon := &vuv1.OverviewGen2V2_CompanyLock{}
		anon.SetLockInTime(lock.GetLockInTime())
		anon.SetLockOutTime(lock.GetLockOutTime())
		anon.SetCompanyName(ddOpts.AnonymizeStringValue(lock.GetCompanyName()))
		anon.SetCompanyAddress(ddOpts.AnonymizeStringValue(lock.GetCompanyAddress()))
		anon.SetCompanyCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(lock.GetCompanyCardNumberAndGeneration()))
		anonLocks[i] = anon
	}
	result.SetCompanyLocks(anonLocks)

	// Anonymize control activities
	anonControls := make([]*vuv1.OverviewGen2V2_ControlActivity, len(overview.GetControlActivities()))
	for i, ctrl := range overview.GetControlActivities() {
		anon := &vuv1.OverviewGen2V2_ControlActivity{}
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

// ===== V2-specific parse helpers =====

// parseVehicleRegistrationNumberRecordArray parses a VehicleRegistrationNumberRecordArray.
// Expected: 1 record × 13 bytes (IA5String, Gen2 V2 only).
func parseVehicleRegistrationNumberRecordArray(data []byte, offset int) (*ddv1.Ia5StringValue, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}
	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 VRN record, got %d", noOfRecords)
	}
	recordStart := offset + headerSize
	recordEnd := recordStart + int(recordSize)
	if recordEnd > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for VRN record")
	}
	var unmarshalOpts dd.UnmarshalOptions
	vrn, err := unmarshalOpts.UnmarshalIa5StringValue(data[recordStart:recordEnd])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal VRN: %w", err)
	}
	totalSize := headerSize + int(recordSize)
	return vrn, totalSize, nil
}

// parseDownloadActivityDataRecordArrayGen2V2 parses a VuDownloadActivityDataRecordArray for V2.
func parseDownloadActivityDataRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.OverviewGen2V2_DownloadActivity, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedSize = 4 + 19 + 36 // 59
	if recordSize != 0 && int(recordSize) < expectedSize {
		return nil, 0, fmt.Errorf("VuDownloadActivityData size %d too small (need at least %d)", recordSize, expectedSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.OverviewGen2V2_DownloadActivity, 0, noOfRecords)
	recordStart := offset + headerSize

	for range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuDownloadActivityData")
		}
		rec := data[recordStart:recordEnd]

		activity := &vuv1.OverviewGen2V2_DownloadActivity{}

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

// parseCompanyLocksRecordArrayGen2V2 parses a VuCompanyLocksRecordArray for V2.
func parseCompanyLocksRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.OverviewGen2V2_CompanyLock, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedSize = 4 + 4 + 36 + 36 + 19 // 99
	if recordSize != 0 && int(recordSize) < expectedSize {
		return nil, 0, fmt.Errorf("VuCompanyLocksRecord size %d too small (need at least %d)", recordSize, expectedSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.OverviewGen2V2_CompanyLock, 0, noOfRecords)
	recordStart := offset + headerSize

	for range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuCompanyLocksRecord")
		}
		rec := data[recordStart:recordEnd]

		lock := &vuv1.OverviewGen2V2_CompanyLock{}

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

// parseControlActivityRecordArrayGen2V2 parses a VuControlActivityRecordArray for V2.
func parseControlActivityRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.OverviewGen2V2_ControlActivity, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedSize = 1 + 4 + 19 + 4 + 4 // 32
	if recordSize != 0 && int(recordSize) < expectedSize {
		return nil, 0, fmt.Errorf("VuControlActivityRecord size %d too small (need at least %d)", recordSize, expectedSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.OverviewGen2V2_ControlActivity, 0, noOfRecords)
	recordStart := offset + headerSize

	for range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuControlActivityRecord")
		}
		rec := data[recordStart:recordEnd]

		ctrl := &vuv1.OverviewGen2V2_ControlActivity{}

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

// ===== V2-specific marshal helpers =====

func marshalDownloadActivitiesGen2V2(opts dd.MarshalOptions, activities []*vuv1.OverviewGen2V2_DownloadActivity) ([]byte, error) {
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

func marshalCompanyLocksGen2V2(opts dd.MarshalOptions, locks []*vuv1.OverviewGen2V2_CompanyLock) ([]byte, error) {
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

func marshalControlActivitiesGen2V2(opts dd.MarshalOptions, controls []*vuv1.OverviewGen2V2_ControlActivity) ([]byte, error) {
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
