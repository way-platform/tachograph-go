package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalOverviewGen1 parses Gen1 Overview data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Overview structure (from Data Dictionary and Appendix 7, Section 2.2.6.2):
//
// ASN.1 Definition:
//
//	VuOverviewFirstGen ::= SEQUENCE {
//	    memberStateCertificate            MemberStateCertificateFirstGen,    -- 194 bytes
//	    vuCertificate                     VuCertificateFirstGen,              -- 194 bytes
//	    vehicleIdentificationNumber       VehicleIdentificationNumber,        -- 17 bytes
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,  -- 15 bytes
//	    currentDateTime                   CurrentDateTime,                    -- 4 bytes
//	    vuDownloadablePeriod              VuDownloadablePeriod,               -- 8 bytes
//	    cardSlotsStatus                   CardSlotsStatus,                    -- 1 byte
//	    vuDownloadActivityData            VuDownloadActivityDataFirstGen,     -- 58 bytes
//	    vuCompanyLocksData                VuCompanyLocksDataFirstGen,         -- 1 + (N * 98) bytes
//	    vuControlActivityData             VuControlActivityDataFirstGen,      -- 1 + (M * 31) bytes
//	    signature                         SignatureFirstGen                   -- 128 bytes (RSA)
//	}
//
// Binary Layout:
// - MemberStateCertificate: 194 bytes
// - VuCertificate: 194 bytes
// - VehicleIdentificationNumber: 17 bytes (IA5String)
// - VehicleRegistrationIdentification: 15 bytes (1 nation + 1 codePage + 13 vrn)
// - CurrentDateTime: 4 bytes (TimeReal)
// - VuDownloadablePeriod: 8 bytes (2 x TimeReal)
// - CardSlotsStatus: 1 byte (4-bit driver slot | 4-bit co-driver slot)
// - VuDownloadActivityData: 58 bytes
//   - DownloadingTime: 4 bytes (TimeReal)
//   - FullCardNumber: 18 bytes (1 EquipmentType + 1 NationNumeric + 16 CardNumber)
//   - CompanyOrWorkshopName: 36 bytes (1 CodePage + 35 Name bytes)
//
// - VuCompanyLocksData: 1 byte (noOfLocks) + (noOfLocks * 98 bytes per record)
//   - Each VuCompanyLocksRecordFirstGen: 98 bytes
//   - LockInTime: 4 bytes
//   - LockOutTime: 4 bytes
//   - CompanyName: 36 bytes
//   - CompanyAddress: 36 bytes
//   - CompanyCardNumber: 18 bytes
//
// - VuControlActivityData: 1 byte (noOfControls) + (noOfControls * 31 bytes per record)
//   - Each VuControlActivityRecordFirstGen: 31 bytes
//   - ControlType: 1 byte
//   - ControlTime: 4 bytes
//   - ControlCardNumber: 18 bytes
//   - DownloadPeriodBeginTime: 4 bytes
//   - DownloadPeriodEndTime: 4 bytes
//
// - Signature: 128 bytes (RSA-1024)
func unmarshalOverviewGen1(value []byte) (*vuv1.OverviewGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	overview := &vuv1.OverviewGen1{}
	overview.SetRawData(value) // Store complete transfer value for painting

	offset := 0
	opts := dd.UnmarshalOptions{PreserveRawData: true}

	// MemberStateCertificate (194 bytes)
	if offset+194 > len(data) {
		return nil, fmt.Errorf("insufficient data for MemberStateCertificate")
	}
	overview.SetMemberStateCertificate(data[offset : offset+194])
	offset += 194

	// VuCertificate (194 bytes)
	if offset+194 > len(data) {
		return nil, fmt.Errorf("insufficient data for VuCertificate")
	}
	overview.SetVuCertificate(data[offset : offset+194])
	offset += 194

	// VehicleIdentificationNumber (17 bytes)
	if offset+17 > len(data) {
		return nil, fmt.Errorf("insufficient data for VehicleIdentificationNumber")
	}
	vin, err := opts.UnmarshalIa5StringValue(data[offset : offset+17])
	if err != nil {
		return nil, fmt.Errorf("unmarshal VIN: %w", err)
	}
	overview.SetVehicleIdentificationNumber(vin)
	offset += 17

	// VehicleRegistrationIdentification (15 bytes)
	if offset+15 > len(data) {
		return nil, fmt.Errorf("insufficient data for VehicleRegistrationIdentification")
	}
	vrn, err := opts.UnmarshalVehicleRegistration(data[offset : offset+15])
	if err != nil {
		return nil, fmt.Errorf("unmarshal VehicleRegistrationIdentification: %w", err)
	}
	overview.SetVehicleRegistrationWithNation(vrn)
	offset += 15

	// CurrentDateTime (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for CurrentDateTime")
	}
	currentTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal CurrentDateTime: %w", err)
	}
	overview.SetCurrentDateTime(currentTime)
	offset += 4

	// VuDownloadablePeriod (8 bytes: 2 x TimeReal)
	if offset+8 > len(data) {
		return nil, fmt.Errorf("insufficient data for VuDownloadablePeriod")
	}
	minTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal DownloadablePeriod minTime: %w", err)
	}
	maxTime, err := opts.UnmarshalTimeReal(data[offset+4 : offset+8])
	if err != nil {
		return nil, fmt.Errorf("unmarshal DownloadablePeriod maxTime: %w", err)
	}
	downloadablePeriod := &ddv1.DownloadablePeriod{}
	downloadablePeriod.SetMinTime(minTime)
	downloadablePeriod.SetMaxTime(maxTime)
	overview.SetDownloadablePeriod(downloadablePeriod)
	offset += 8

	// CardSlotsStatus (1 byte)
	// Lower 4 bits (0-3): driver slot
	// Upper 4 bits (4-7): co-driver slot
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for CardSlotsStatus")
	}
	cardSlotsStatus := data[offset]
	driverSlotRaw := byte(cardSlotsStatus & 0x0F)
	coDriverSlotRaw := byte((cardSlotsStatus >> 4) & 0x0F)

	// Use UnmarshalEnum to properly map protocol values to enum values
	driverSlot, err := dd.UnmarshalEnum[ddv1.SlotCardType](driverSlotRaw)
	if err != nil {
		// Unrecognized value - set to UNRECOGNIZED
		driverSlot = ddv1.SlotCardType_SLOT_CARD_TYPE_UNRECOGNIZED
	}
	coDriverSlot, err := dd.UnmarshalEnum[ddv1.SlotCardType](coDriverSlotRaw)
	if err != nil {
		// Unrecognized value - set to UNRECOGNIZED
		coDriverSlot = ddv1.SlotCardType_SLOT_CARD_TYPE_UNRECOGNIZED
	}

	overview.SetDriverSlotCard(driverSlot)
	overview.SetCoDriverSlotCard(coDriverSlot)
	offset += 1

	// VuDownloadActivityData (58 bytes: 4 + 18 + 36)
	if offset+58 > len(data) {
		return nil, fmt.Errorf("insufficient data for VuDownloadActivityData")
	}

	downloadActivity := &vuv1.OverviewGen1_DownloadActivity{}

	// DownloadingTime (4 bytes)
	downloadingTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal downloading time: %w", err)
	}
	downloadActivity.SetDownloadingTime(downloadingTime)
	offset += 4

	// FullCardNumber (18 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumber(data[offset : offset+18])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number: %w", err)
	}
	downloadActivity.SetFullCardNumber(fullCardNumber)
	offset += 18

	// CompanyOrWorkshopName (36 bytes: 1 code page + 35 name)
	companyName, err := opts.UnmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("unmarshal company name: %w", err)
	}
	downloadActivity.SetCompanyOrWorkshopName(companyName)
	offset += 36

	overview.SetDownloadActivities([]*vuv1.OverviewGen1_DownloadActivity{downloadActivity})

	// VuCompanyLocksData: 1 byte (noOfLocks) + (noOfLocks * 98 bytes per record)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for VuCompanyLocksData noOfLocks")
	}
	noOfLocks := data[offset]
	offset += 1

	const companyLockRecordSize = 98 // 4 + 4 + 36 + 36 + 18
	if offset+int(noOfLocks)*companyLockRecordSize > len(data) {
		return nil, fmt.Errorf("insufficient data for VuCompanyLocksData records")
	}

	companyLocks := make([]*vuv1.OverviewGen1_CompanyLock, noOfLocks)
	for i := 0; i < int(noOfLocks); i++ {
		lock := &vuv1.OverviewGen1_CompanyLock{}

		// LockInTime (4 bytes)
		lockInTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal lockInTime: %w", err)
		}
		lock.SetLockInTime(lockInTime)
		offset += 4

		// LockOutTime (4 bytes)
		lockOutTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal lockOutTime: %w", err)
		}
		lock.SetLockOutTime(lockOutTime)
		offset += 4

		// CompanyName (36 bytes)
		companyName, err := opts.UnmarshalStringValue(data[offset : offset+36])
		if err != nil {
			return nil, fmt.Errorf("unmarshal company name: %w", err)
		}
		lock.SetCompanyName(companyName)
		offset += 36

		// CompanyAddress (36 bytes)
		companyAddress, err := opts.UnmarshalStringValue(data[offset : offset+36])
		if err != nil {
			return nil, fmt.Errorf("unmarshal company address: %w", err)
		}
		lock.SetCompanyAddress(companyAddress)
		offset += 36

		// CompanyCardNumber (18 bytes)
		companyCardNumber, err := opts.UnmarshalFullCardNumber(data[offset : offset+18])
		if err != nil {
			return nil, fmt.Errorf("unmarshal company card number: %w", err)
		}
		lock.SetCompanyCardNumber(companyCardNumber)
		offset += 18

		companyLocks[i] = lock
	}
	overview.SetCompanyLocks(companyLocks)

	// VuControlActivityData: 1 byte (noOfControls) + (noOfControls * 31 bytes per record)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for VuControlActivityData noOfControls")
	}
	noOfControls := data[offset]
	offset += 1

	const controlActivityRecordSize = 31 // 1 + 4 + 18 + 4 + 4
	if offset+int(noOfControls)*controlActivityRecordSize > len(data) {
		return nil, fmt.Errorf("insufficient data for VuControlActivityData records")
	}

	controlActivities := make([]*vuv1.OverviewGen1_ControlActivity, noOfControls)
	for i := 0; i < int(noOfControls); i++ {
		control := &vuv1.OverviewGen1_ControlActivity{}

		// ControlType (1 byte)
		controlType, err := opts.UnmarshalControlType(data[offset : offset+1])
		if err != nil {
			return nil, fmt.Errorf("unmarshal control type: %w", err)
		}
		control.SetControlType(controlType)
		offset += 1

		// ControlTime (4 bytes)
		controlTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal control time: %w", err)
		}
		control.SetControlTime(controlTime)
		offset += 4

		// ControlCardNumber (18 bytes)
		controlCardNumber, err := opts.UnmarshalFullCardNumber(data[offset : offset+18])
		if err != nil {
			return nil, fmt.Errorf("unmarshal control card number: %w", err)
		}
		control.SetControlCardNumber(controlCardNumber)
		offset += 18

		// DownloadPeriodBeginTime (4 bytes)
		downloadPeriodBeginTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal download period begin time: %w", err)
		}
		control.SetDownloadPeriodBeginTime(downloadPeriodBeginTime)
		offset += 4

		// DownloadPeriodEndTime (4 bytes)
		downloadPeriodEndTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal download period end time: %w", err)
		}
		control.SetDownloadPeriodEndTime(downloadPeriodEndTime)
		offset += 4

		controlActivities[i] = control
	}
	overview.SetControlActivities(controlActivities)

	// Store signature (extracted at the beginning)
	overview.SetSignature(signature)

	// Verify we consumed exactly the right amount of data
	if offset != len(data) {
		return nil, fmt.Errorf("Overview Gen1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return overview, nil
}

// MarshalOverviewGen1 marshals Gen1 Overview data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and has the correct length, it uses it as a canvas and paints semantic values over it.
// Otherwise, it creates a zero-filled canvas and encodes from semantic fields.
func (opts MarshalOptions) MarshalOverviewGen1(overview *vuv1.OverviewGen1) ([]byte, error) {
	if overview == nil {
		return nil, fmt.Errorf("overview cannot be nil")
	}

	// Calculate expected size (signature is stored separately in RawVehicleUnitFile_Record)
	noOfLocks := len(overview.GetCompanyLocks())
	noOfControls := len(overview.GetControlActivities())
	expectedSize := 491 + 1 + (noOfLocks * 98) + 1 + (noOfControls * 31)
	// 491 = 194 + 194 + 17 + 15 + 4 + 8 + 1 + 58

	// Use raw_data as canvas if available
	var canvas []byte
	if raw := overview.GetRawData(); len(raw) == expectedSize+128 {
		canvas = make([]byte, expectedSize)
		copy(canvas, raw[:expectedSize])
	} else if raw := overview.GetRawData(); len(raw) == expectedSize {
		canvas = make([]byte, len(raw))
		copy(canvas, raw)
	} else {
		canvas = make([]byte, expectedSize)
	}

	// Paint semantic values over canvas
	offset := 0

	// MemberStateCertificate (194 bytes)
	copy(canvas[offset:offset+194], overview.GetMemberStateCertificate())
	offset += 194

	// VuCertificate (194 bytes)
	copy(canvas[offset:offset+194], overview.GetVuCertificate())
	offset += 194

	// VehicleIdentificationNumber (17 bytes)
	vin := overview.GetVehicleIdentificationNumber()
	if vin != nil {
		vinBytes, err := opts.MarshalIa5StringValue(vin)
		if err != nil {
			return nil, fmt.Errorf("append VIN: %w", err)
		}
		copy(canvas[offset:offset+17], vinBytes)
	}
	offset += 17

	// VehicleRegistrationIdentification (15 bytes)
	vrn := overview.GetVehicleRegistrationWithNation()
	if vrn != nil {
		vrnBytes, err := opts.MarshalVehicleRegistration(vrn)
		if err != nil {
			return nil, fmt.Errorf("append VRN: %w", err)
		}
		copy(canvas[offset:offset+15], vrnBytes)
	}
	offset += 15

	// CurrentDateTime (4 bytes)
	currentTime := overview.GetCurrentDateTime()
	if currentTime != nil {
		timeBytes, err := opts.MarshalTimeReal(currentTime)
		if err != nil {
			return nil, fmt.Errorf("append current time: %w", err)
		}
		copy(canvas[offset:offset+4], timeBytes)
	}
	offset += 4

	// VuDownloadablePeriod (8 bytes)
	downloadablePeriod := overview.GetDownloadablePeriod()
	if downloadablePeriod != nil {
		minTimeBytes, err := opts.MarshalTimeReal(downloadablePeriod.GetMinTime())
		if err != nil {
			return nil, fmt.Errorf("append min time: %w", err)
		}
		copy(canvas[offset:offset+4], minTimeBytes)
		offset += 4

		maxTimeBytes, err := opts.MarshalTimeReal(downloadablePeriod.GetMaxTime())
		if err != nil {
			return nil, fmt.Errorf("append max time: %w", err)
		}
		copy(canvas[offset:offset+4], maxTimeBytes)
		offset += 4
	} else {
		offset += 8
	}

	// CardSlotsStatus (1 byte)
	// Only overwrite if both values can be marshalled (not UNRECOGNIZED)
	driverSlot, driverErr := dd.MarshalEnum(overview.GetDriverSlotCard())
	coDriverSlot, coDriverErr := dd.MarshalEnum(overview.GetCoDriverSlotCard())

	// If both values are recognized, write them; otherwise preserve canvas value (from raw_data)
	if driverErr == nil && coDriverErr == nil {
		canvas[offset] = (coDriverSlot << 4) | (driverSlot & 0x0F)
	}
	offset += 1

	// VuDownloadActivityData (58 bytes)
	downloadActivities := overview.GetDownloadActivities()
	if len(downloadActivities) > 0 {
		activity := downloadActivities[0]

		// DownloadingTime (4 bytes)
		downloadingTimeBytes, err := opts.MarshalTimeReal(activity.GetDownloadingTime())
		if err != nil {
			return nil, fmt.Errorf("append downloading time: %w", err)
		}
		copy(canvas[offset:offset+4], downloadingTimeBytes)
		offset += 4

		// FullCardNumber (18 bytes)
		cardNumberBytes, err := opts.MarshalFullCardNumber(activity.GetFullCardNumber())
		if err != nil {
			return nil, fmt.Errorf("append full card number: %w", err)
		}
		copy(canvas[offset:offset+18], cardNumberBytes)
		offset += 18

		// CompanyOrWorkshopName (36 bytes)
		companyNameBytes, err := opts.MarshalStringValue(activity.GetCompanyOrWorkshopName())
		if err != nil {
			return nil, fmt.Errorf("append company name: %w", err)
		}
		copy(canvas[offset:offset+36], companyNameBytes)
		offset += 36
	} else {
		offset += 58
	}

	// VuCompanyLocksData
	canvas[offset] = byte(noOfLocks)
	offset += 1

	for _, lock := range overview.GetCompanyLocks() {
		// LockInTime (4 bytes)
		lockInTimeBytes, err := opts.MarshalTimeReal(lock.GetLockInTime())
		if err != nil {
			return nil, fmt.Errorf("append lock in time: %w", err)
		}
		if lock.GetLockInTime() != nil {
			copy(canvas[offset:offset+4], lockInTimeBytes)
		}
		offset += 4

		// LockOutTime (4 bytes)
		lockOutTimeBytes, err := opts.MarshalTimeReal(lock.GetLockOutTime())
		if err != nil {
			return nil, fmt.Errorf("append lock out time: %w", err)
		}
		if lock.GetLockOutTime() != nil {
			copy(canvas[offset:offset+4], lockOutTimeBytes)
		}
		offset += 4

		// CompanyName (36 bytes)
		companyNameBytes, err := opts.MarshalStringValue(lock.GetCompanyName())
		if err != nil {
			return nil, fmt.Errorf("append company name: %w", err)
		}
		copy(canvas[offset:offset+36], companyNameBytes)
		offset += 36

		// CompanyAddress (36 bytes)
		companyAddressBytes, err := opts.MarshalStringValue(lock.GetCompanyAddress())
		if err != nil {
			return nil, fmt.Errorf("append company address: %w", err)
		}
		copy(canvas[offset:offset+36], companyAddressBytes)
		offset += 36

		// CompanyCardNumber (18 bytes)
		companyCardNumberBytes, err := opts.MarshalFullCardNumber(lock.GetCompanyCardNumber())
		if err != nil {
			return nil, fmt.Errorf("append company card number: %w", err)
		}
		copy(canvas[offset:offset+18], companyCardNumberBytes)
		offset += 18
	}

	// VuControlActivityData
	canvas[offset] = byte(noOfControls)
	offset += 1

	for _, control := range overview.GetControlActivities() {
		// ControlType (1 byte)
		controlTypeBytes, err := opts.MarshalControlType(control.GetControlType())
		if err != nil {
			return nil, fmt.Errorf("append control type: %w", err)
		}
		copy(canvas[offset:offset+1], controlTypeBytes)
		offset += 1

		// ControlTime (4 bytes)
		controlTimeBytes, err := opts.MarshalTimeReal(control.GetControlTime())
		if err != nil {
			return nil, fmt.Errorf("append control time: %w", err)
		}
		copy(canvas[offset:offset+4], controlTimeBytes)
		offset += 4

		// ControlCardNumber (18 bytes)
		controlCardNumberBytes, err := opts.MarshalFullCardNumber(control.GetControlCardNumber())
		if err != nil {
			return nil, fmt.Errorf("append control card number: %w", err)
		}
		copy(canvas[offset:offset+18], controlCardNumberBytes)
		offset += 18

		// DownloadPeriodBeginTime (4 bytes)
		downloadPeriodBeginTimeBytes, err := opts.MarshalTimeReal(control.GetDownloadPeriodBeginTime())
		if err != nil {
			return nil, fmt.Errorf("append download period begin time: %w", err)
		}
		copy(canvas[offset:offset+4], downloadPeriodBeginTimeBytes)
		offset += 4

		// DownloadPeriodEndTime (4 bytes)
		downloadPeriodEndTimeBytes, err := opts.MarshalTimeReal(control.GetDownloadPeriodEndTime())
		if err != nil {
			return nil, fmt.Errorf("append download period end time: %w", err)
		}
		copy(canvas[offset:offset+4], downloadPeriodEndTimeBytes)
		offset += 4
	}

	// Append signature to create complete transfer value
	signature := overview.GetSignature()
	if len(signature) == 0 {
		// Gen1 uses fixed 128-byte RSA-1024 signatures
		signature = make([]byte, 128)
	}
	transferValue := append(canvas, signature...)

	return transferValue, nil
}

// anonymizeOverviewGen1 anonymizes Gen1 Overview data.
func (opts AnonymizeOptions) anonymizeOverviewGen1(overview *vuv1.OverviewGen1) *vuv1.OverviewGen1 {
	if overview == nil {
		return nil
	}

	result := proto.Clone(overview).(*vuv1.OverviewGen1)

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize VIN
	if vin := result.GetVehicleIdentificationNumber(); vin != nil {
		result.SetVehicleIdentificationNumber(ddOpts.AnonymizeIa5StringValue(vin))
	}

	// Anonymize VRN
	if vrn := result.GetVehicleRegistrationWithNation(); vrn != nil {
		result.SetVehicleRegistrationWithNation(ddOpts.AnonymizeVehicleRegistrationIdentification(vrn))
	}

	// Clear certificates (will be invalid after anonymization anyway)
	result.SetMemberStateCertificate(nil)
	result.SetVuCertificate(nil)

	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Anonymize download activities
	var anonymizedDownloadActivities []*vuv1.OverviewGen1_DownloadActivity
	for _, activity := range result.GetDownloadActivities() {
		anonActivity := proto.Clone(activity).(*vuv1.OverviewGen1_DownloadActivity)
		// Anonymize card number
		if anonActivity.GetFullCardNumber() != nil {
			anonActivity.SetFullCardNumber(ddOpts.AnonymizeFullCardNumber(anonActivity.GetFullCardNumber()))
		}
		// Anonymize company/workshop name
		if anonActivity.GetCompanyOrWorkshopName() != nil {
			anonActivity.SetCompanyOrWorkshopName(ddOpts.AnonymizeStringValue(anonActivity.GetCompanyOrWorkshopName()))
		}
		anonymizedDownloadActivities = append(anonymizedDownloadActivities, anonActivity)
	}
	result.SetDownloadActivities(anonymizedDownloadActivities)

	// Anonymize company locks
	var anonymizedCompanyLocks []*vuv1.OverviewGen1_CompanyLock
	for _, lock := range result.GetCompanyLocks() {
		anonLock := proto.Clone(lock).(*vuv1.OverviewGen1_CompanyLock)
		// Anonymize company name
		if anonLock.GetCompanyName() != nil {
			anonLock.SetCompanyName(ddOpts.AnonymizeStringValue(anonLock.GetCompanyName()))
		}
		// Anonymize company address
		if anonLock.GetCompanyAddress() != nil {
			anonLock.SetCompanyAddress(ddOpts.AnonymizeStringValue(anonLock.GetCompanyAddress()))
		}
		// Anonymize company card number
		if anonLock.GetCompanyCardNumber() != nil {
			anonLock.SetCompanyCardNumber(ddOpts.AnonymizeFullCardNumber(anonLock.GetCompanyCardNumber()))
		}
		anonymizedCompanyLocks = append(anonymizedCompanyLocks, anonLock)
	}
	result.SetCompanyLocks(anonymizedCompanyLocks)

	// Anonymize control activities
	var anonymizedControlActivities []*vuv1.OverviewGen1_ControlActivity
	for _, activity := range result.GetControlActivities() {
		anonActivity := proto.Clone(activity).(*vuv1.OverviewGen1_ControlActivity)
		// Anonymize control card number
		if anonActivity.GetControlCardNumber() != nil {
			anonActivity.SetControlCardNumber(ddOpts.AnonymizeFullCardNumber(anonActivity.GetControlCardNumber()))
		}
		anonymizedControlActivities = append(anonymizedControlActivities, anonActivity)
	}
	result.SetControlActivities(anonymizedControlActivities)

	// Note: We intentionally keep raw_data here because MarshalOverviewGen1 uses
	// raw_data painting to serialize. The painting will apply the anonymized
	// semantic values (VIN, VRN, company data) over the raw_data canvas during marshalling.
	// This is the recommended approach per the raw data painting policy.

	return result
}
