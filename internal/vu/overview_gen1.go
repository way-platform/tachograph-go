package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalOverviewGen1 parses Gen1 Overview data from the complete transfer value.
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
// - Signature: 128 bytes (RSA)
func unmarshalOverviewGen1(value []byte) (*vuv1.OverviewGen1, error) {
	overview := &vuv1.OverviewGen1{}
	overview.SetRawData(value)

	offset := 0
	var opts dd.UnmarshalOptions

	// MemberStateCertificate (194 bytes)
	if offset+194 > len(value) {
		return nil, fmt.Errorf("insufficient data for MemberStateCertificate")
	}
	overview.SetMemberStateCertificate(value[offset : offset+194])
	offset += 194

	// VuCertificate (194 bytes)
	if offset+194 > len(value) {
		return nil, fmt.Errorf("insufficient data for VuCertificate")
	}
	overview.SetVuCertificate(value[offset : offset+194])
	offset += 194

	// VehicleIdentificationNumber (17 bytes)
	if offset+17 > len(value) {
		return nil, fmt.Errorf("insufficient data for VehicleIdentificationNumber")
	}
	vin, err := opts.UnmarshalIa5StringValue(value[offset : offset+17])
	if err != nil {
		return nil, fmt.Errorf("unmarshal VIN: %w", err)
	}
	overview.SetVehicleIdentificationNumber(vin)
	offset += 17

	// VehicleRegistrationIdentification (15 bytes)
	if offset+15 > len(value) {
		return nil, fmt.Errorf("insufficient data for VehicleRegistrationIdentification")
	}
	vrn, err := opts.UnmarshalVehicleRegistration(value[offset : offset+15])
	if err != nil {
		return nil, fmt.Errorf("unmarshal VehicleRegistrationIdentification: %w", err)
	}
	overview.SetVehicleRegistrationWithNation(vrn)
	offset += 15

	// CurrentDateTime (4 bytes)
	if offset+4 > len(value) {
		return nil, fmt.Errorf("insufficient data for CurrentDateTime")
	}
	currentTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal CurrentDateTime: %w", err)
	}
	overview.SetCurrentDateTime(currentTime)
	offset += 4

	// VuDownloadablePeriod (8 bytes: 2 x TimeReal)
	if offset+8 > len(value) {
		return nil, fmt.Errorf("insufficient data for VuDownloadablePeriod")
	}
	minTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal DownloadablePeriod minTime: %w", err)
	}
	maxTime, err := opts.UnmarshalTimeReal(value[offset+4 : offset+8])
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
	if offset+1 > len(value) {
		return nil, fmt.Errorf("insufficient data for CardSlotsStatus")
	}
	cardSlotsStatus := value[offset]
	driverSlot := ddv1.SlotCardType(cardSlotsStatus & 0x0F)
	coDriverSlot := ddv1.SlotCardType((cardSlotsStatus >> 4) & 0x0F)
	overview.SetDriverSlotCard(driverSlot)
	overview.SetCoDriverSlotCard(coDriverSlot)
	offset += 1

	// VuDownloadActivityData (58 bytes: 4 + 18 + 36)
	if offset+58 > len(value) {
		return nil, fmt.Errorf("insufficient data for VuDownloadActivityData")
	}

	downloadActivity := &vuv1.OverviewGen1_DownloadActivity{}

	// DownloadingTime (4 bytes)
	downloadingTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal downloading time: %w", err)
	}
	downloadActivity.SetDownloadingTime(downloadingTime)
	offset += 4

	// FullCardNumber (18 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumber(value[offset : offset+18])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number: %w", err)
	}
	downloadActivity.SetFullCardNumber(fullCardNumber)
	offset += 18

	// CompanyOrWorkshopName (36 bytes: 1 code page + 35 name)
	companyName, err := opts.UnmarshalStringValue(value[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("unmarshal company name: %w", err)
	}
	downloadActivity.SetCompanyOrWorkshopName(companyName)
	offset += 36

	overview.SetDownloadActivities([]*vuv1.OverviewGen1_DownloadActivity{downloadActivity})

	// VuCompanyLocksData: 1 byte (noOfLocks) + (noOfLocks * 98 bytes per record)
	if offset+1 > len(value) {
		return nil, fmt.Errorf("insufficient data for VuCompanyLocksData noOfLocks")
	}
	noOfLocks := value[offset]
	offset += 1

	const companyLockRecordSize = 98 // 4 + 4 + 36 + 36 + 18
	if offset+int(noOfLocks)*companyLockRecordSize > len(value) {
		return nil, fmt.Errorf("insufficient data for VuCompanyLocksData records")
	}

	companyLocks := make([]*vuv1.OverviewGen1_CompanyLock, noOfLocks)
	for i := 0; i < int(noOfLocks); i++ {
		lock := &vuv1.OverviewGen1_CompanyLock{}

		// LockInTime (4 bytes)
		lockInTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal lockInTime: %w", err)
		}
		lock.SetLockInTime(lockInTime)
		offset += 4

		// LockOutTime (4 bytes)
		lockOutTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal lockOutTime: %w", err)
		}
		lock.SetLockOutTime(lockOutTime)
		offset += 4

		// CompanyName (36 bytes)
		companyName, err := opts.UnmarshalStringValue(value[offset : offset+36])
		if err != nil {
			return nil, fmt.Errorf("unmarshal company name: %w", err)
		}
		lock.SetCompanyName(companyName)
		offset += 36

		// CompanyAddress (36 bytes)
		companyAddress, err := opts.UnmarshalStringValue(value[offset : offset+36])
		if err != nil {
			return nil, fmt.Errorf("unmarshal company address: %w", err)
		}
		lock.SetCompanyAddress(companyAddress)
		offset += 36

		// CompanyCardNumber (18 bytes)
		companyCardNumber, err := opts.UnmarshalFullCardNumber(value[offset : offset+18])
		if err != nil {
			return nil, fmt.Errorf("unmarshal company card number: %w", err)
		}
		lock.SetCompanyCardNumber(companyCardNumber)
		offset += 18

		companyLocks[i] = lock
	}
	overview.SetCompanyLocks(companyLocks)

	// VuControlActivityData: 1 byte (noOfControls) + (noOfControls * 31 bytes per record)
	if offset+1 > len(value) {
		return nil, fmt.Errorf("insufficient data for VuControlActivityData noOfControls")
	}
	noOfControls := value[offset]
	offset += 1

	const controlActivityRecordSize = 31 // 1 + 4 + 18 + 4 + 4
	if offset+int(noOfControls)*controlActivityRecordSize > len(value) {
		return nil, fmt.Errorf("insufficient data for VuControlActivityData records")
	}

	controlActivities := make([]*vuv1.OverviewGen1_ControlActivity, noOfControls)
	for i := 0; i < int(noOfControls); i++ {
		control := &vuv1.OverviewGen1_ControlActivity{}

		// ControlType (1 byte)
		controlType, err := opts.UnmarshalControlType(value[offset : offset+1])
		if err != nil {
			return nil, fmt.Errorf("unmarshal control type: %w", err)
		}
		control.SetControlType(controlType)
		offset += 1

		// ControlTime (4 bytes)
		controlTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal control time: %w", err)
		}
		control.SetControlTime(controlTime)
		offset += 4

		// ControlCardNumber (18 bytes)
		controlCardNumber, err := opts.UnmarshalFullCardNumber(value[offset : offset+18])
		if err != nil {
			return nil, fmt.Errorf("unmarshal control card number: %w", err)
		}
		control.SetControlCardNumber(controlCardNumber)
		offset += 18

		// DownloadPeriodBeginTime (4 bytes)
		downloadPeriodBeginTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal download period begin time: %w", err)
		}
		control.SetDownloadPeriodBeginTime(downloadPeriodBeginTime)
		offset += 4

		// DownloadPeriodEndTime (4 bytes)
		downloadPeriodEndTime, err := opts.UnmarshalTimeReal(value[offset : offset+4])
		if err != nil {
			return nil, fmt.Errorf("unmarshal download period end time: %w", err)
		}
		control.SetDownloadPeriodEndTime(downloadPeriodEndTime)
		offset += 4

		controlActivities[i] = control
	}
	overview.SetControlActivities(controlActivities)

	// Signature (128 bytes - RSA for Gen1)
	if offset+128 > len(value) {
		return nil, fmt.Errorf("insufficient data for Signature")
	}
	overview.SetSignature(value[offset : offset+128])
	offset += 128

	// Verify we consumed exactly the right amount of data
	if offset != len(value) {
		return nil, fmt.Errorf("Overview Gen1 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	return overview, nil
}

// appendOverviewGen1 marshals Gen1 Overview data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and has the correct length, it uses it as a canvas and paints semantic values over it.
// Otherwise, it creates a zero-filled canvas and encodes from semantic fields.
func appendOverviewGen1(dst []byte, overview *vuv1.OverviewGen1) ([]byte, error) {
	if overview == nil {
		return nil, fmt.Errorf("overview cannot be nil")
	}

	// Calculate expected size
	noOfLocks := len(overview.GetCompanyLocks())
	noOfControls := len(overview.GetControlActivities())
	expectedSize := 491 + 1 + (noOfLocks * 98) + 1 + (noOfControls * 31) + 128
	// 491 = 194 + 194 + 17 + 15 + 4 + 8 + 1 + 58

	// Use raw_data as canvas if available
	var canvas []byte
	if raw := overview.GetRawData(); len(raw) == expectedSize {
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
		vinBytes, err := dd.AppendIa5StringValue(nil, vin)
		if err != nil {
			return nil, fmt.Errorf("append VIN: %w", err)
		}
		copy(canvas[offset:offset+17], vinBytes)
	}
	offset += 17

	// VehicleRegistrationIdentification (15 bytes)
	vrn := overview.GetVehicleRegistrationWithNation()
	if vrn != nil {
		vrnBytes, err := dd.AppendVehicleRegistration(nil, vrn)
		if err != nil {
			return nil, fmt.Errorf("append VRN: %w", err)
		}
		copy(canvas[offset:offset+15], vrnBytes)
	}
	offset += 15

	// CurrentDateTime (4 bytes)
	currentTime := overview.GetCurrentDateTime()
	if currentTime != nil {
		timeBytes, err := dd.AppendTimeReal(nil, currentTime)
		if err != nil {
			return nil, fmt.Errorf("append current time: %w", err)
		}
		copy(canvas[offset:offset+4], timeBytes)
	}
	offset += 4

	// VuDownloadablePeriod (8 bytes)
	downloadablePeriod := overview.GetDownloadablePeriod()
	if downloadablePeriod != nil {
		minTimeBytes, err := dd.AppendTimeReal(nil, downloadablePeriod.GetMinTime())
		if err != nil {
			return nil, fmt.Errorf("append min time: %w", err)
		}
		copy(canvas[offset:offset+4], minTimeBytes)
		offset += 4

		maxTimeBytes, err := dd.AppendTimeReal(nil, downloadablePeriod.GetMaxTime())
		if err != nil {
			return nil, fmt.Errorf("append max time: %w", err)
		}
		copy(canvas[offset:offset+4], maxTimeBytes)
		offset += 4
	} else {
		offset += 8
	}

	// CardSlotsStatus (1 byte)
	driverSlot, err := dd.MarshalEnum(overview.GetDriverSlotCard())
	if err != nil {
		driverSlot = 0
	}
	coDriverSlot, err := dd.MarshalEnum(overview.GetCoDriverSlotCard())
	if err != nil {
		coDriverSlot = 0
	}
	canvas[offset] = (coDriverSlot << 4) | (driverSlot & 0x0F)
	offset += 1

	// VuDownloadActivityData (58 bytes)
	downloadActivities := overview.GetDownloadActivities()
	if len(downloadActivities) > 0 {
		activity := downloadActivities[0]

		// DownloadingTime (4 bytes)
		downloadingTimeBytes, err := dd.AppendTimeReal(nil, activity.GetDownloadingTime())
		if err != nil {
			return nil, fmt.Errorf("append downloading time: %w", err)
		}
		copy(canvas[offset:offset+4], downloadingTimeBytes)
		offset += 4

		// FullCardNumber (18 bytes)
		cardNumberBytes, err := dd.AppendFullCardNumber(nil, activity.GetFullCardNumber())
		if err != nil {
			return nil, fmt.Errorf("append full card number: %w", err)
		}
		copy(canvas[offset:offset+18], cardNumberBytes)
		offset += 18

		// CompanyOrWorkshopName (36 bytes)
		companyNameBytes, err := dd.AppendStringValue(nil, activity.GetCompanyOrWorkshopName())
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
		lockInTimeBytes, err := dd.AppendTimeReal(nil, lock.GetLockInTime())
		if err != nil {
			return nil, fmt.Errorf("append lock in time: %w", err)
		}
		copy(canvas[offset:offset+4], lockInTimeBytes)
		offset += 4

		// LockOutTime (4 bytes)
		lockOutTimeBytes, err := dd.AppendTimeReal(nil, lock.GetLockOutTime())
		if err != nil {
			return nil, fmt.Errorf("append lock out time: %w", err)
		}
		copy(canvas[offset:offset+4], lockOutTimeBytes)
		offset += 4

		// CompanyName (36 bytes)
		companyNameBytes, err := dd.AppendStringValue(nil, lock.GetCompanyName())
		if err != nil {
			return nil, fmt.Errorf("append company name: %w", err)
		}
		copy(canvas[offset:offset+36], companyNameBytes)
		offset += 36

		// CompanyAddress (36 bytes)
		companyAddressBytes, err := dd.AppendStringValue(nil, lock.GetCompanyAddress())
		if err != nil {
			return nil, fmt.Errorf("append company address: %w", err)
		}
		copy(canvas[offset:offset+36], companyAddressBytes)
		offset += 36

		// CompanyCardNumber (18 bytes)
		companyCardNumberBytes, err := dd.AppendFullCardNumber(nil, lock.GetCompanyCardNumber())
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
		controlTypeBytes, err := dd.AppendControlType(nil, control.GetControlType())
		if err != nil {
			return nil, fmt.Errorf("append control type: %w", err)
		}
		copy(canvas[offset:offset+1], controlTypeBytes)
		offset += 1

		// ControlTime (4 bytes)
		controlTimeBytes, err := dd.AppendTimeReal(nil, control.GetControlTime())
		if err != nil {
			return nil, fmt.Errorf("append control time: %w", err)
		}
		copy(canvas[offset:offset+4], controlTimeBytes)
		offset += 4

		// ControlCardNumber (18 bytes)
		controlCardNumberBytes, err := dd.AppendFullCardNumber(nil, control.GetControlCardNumber())
		if err != nil {
			return nil, fmt.Errorf("append control card number: %w", err)
		}
		copy(canvas[offset:offset+18], controlCardNumberBytes)
		offset += 18

		// DownloadPeriodBeginTime (4 bytes)
		downloadPeriodBeginTimeBytes, err := dd.AppendTimeReal(nil, control.GetDownloadPeriodBeginTime())
		if err != nil {
			return nil, fmt.Errorf("append download period begin time: %w", err)
		}
		copy(canvas[offset:offset+4], downloadPeriodBeginTimeBytes)
		offset += 4

		// DownloadPeriodEndTime (4 bytes)
		downloadPeriodEndTimeBytes, err := dd.AppendTimeReal(nil, control.GetDownloadPeriodEndTime())
		if err != nil {
			return nil, fmt.Errorf("append download period end time: %w", err)
		}
		copy(canvas[offset:offset+4], downloadPeriodEndTimeBytes)
		offset += 4
	}

	// Signature (128 bytes)
	copy(canvas[offset:offset+128], overview.GetSignature())

	return append(dst, canvas...), nil
}
