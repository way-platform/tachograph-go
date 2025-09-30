package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPreviousVehicleInfo unmarshals a PreviousVehicleInfo from binary data.
//
// The data type `PreviousVehicleInfo` is specified in the Data Dictionary, Section 2.118.
//
// ASN.1 Definition (Gen1):
//
//	PreviousVehicleInfo ::= SEQUENCE {
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,
//	    cardWithdrawalTime                TimeReal
//	}
//
// For Gen2, the following component is added:
//
//	vuGeneration Generation
func (opts UnmarshalOptions) UnmarshalPreviousVehicleInfo(data []byte) (*ddv1.PreviousVehicleInfo, error) {
	const (
		lenPreviousVehicleInfoGen1 = 19 // 15 bytes vehicle reg + 4 bytes time
		lenPreviousVehicleInfoGen2 = 20 // Gen1 + 1 byte generation
		lenVehicleReg              = 15
		idxVehicleReg              = 0
		idxCardWithdrawalTime      = 15
		idxVuGeneration            = 19
	)

	// Check for Gen2; otherwise assume Gen1 (including zero value)
	expectedLen := lenPreviousVehicleInfoGen1
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		expectedLen = lenPreviousVehicleInfoGen2
	}

	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid data length for PreviousVehicleInfo: got %d, want %d (Gen1) or %d (Gen2)", len(data), lenPreviousVehicleInfoGen1, lenPreviousVehicleInfoGen2)
	}

	info := &ddv1.PreviousVehicleInfo{}
	// Populate generation from unmarshal context
	info.SetGeneration(opts.Generation)

	// Parse vehicleRegistrationIdentification (15 bytes)
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[idxVehicleReg : idxVehicleReg+lenVehicleReg])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle registration: %w", err)
	}
	info.SetVehicleRegistration(vehicleReg)

	// Parse cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime, err := opts.UnmarshalTimeReal(data[idxCardWithdrawalTime : idxCardWithdrawalTime+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal card withdrawal time: %w", err)
	}
	info.SetCardWithdrawalTime(withdrawalTime)

	// For Gen2, parse vuGeneration (1 byte)
	if opts.Generation == ddv1.Generation_GENERATION_2 {
		vuGen := ddv1.Generation(data[idxVuGeneration])
		info.SetVuGeneration(vuGen)
	}

	return info, nil
}

// AppendPreviousVehicleInfo appends a PreviousVehicleInfo to dst.
//
// The data type `PreviousVehicleInfo` is specified in the Data Dictionary, Section 2.118.
//
// ASN.1 Definition (Gen1):
//
//	PreviousVehicleInfo ::= SEQUENCE {
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,
//	    cardWithdrawalTime                TimeReal
//	}
//
// For Gen2, the following component is added:
//
//	vuGeneration Generation
func AppendPreviousVehicleInfo(dst []byte, info *ddv1.PreviousVehicleInfo) ([]byte, error) {
	if info == nil {
		return nil, fmt.Errorf("previous vehicle info cannot be nil")
	}

	// Get generation from the self-describing record
	generation := info.GetGeneration()
	if generation == ddv1.Generation_GENERATION_UNSPECIFIED {
		return nil, fmt.Errorf("PreviousVehicleInfo.generation must be specified (got GENERATION_UNSPECIFIED)")
	}

	// Append vehicleRegistrationIdentification (15 bytes)
	var err error
	dst, err = AppendVehicleRegistration(dst, info.GetVehicleRegistration())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle registration: %w", err)
	}

	// Append cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime := info.GetCardWithdrawalTime()
	if withdrawalTime == nil {
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	} else {
		dst, err = AppendTimeReal(dst, withdrawalTime)
		if err != nil {
			return nil, fmt.Errorf("failed to append card withdrawal time: %w", err)
		}
	}

	// For Gen2, append vuGeneration (1 byte)
	if generation == ddv1.Generation_GENERATION_2 {
		dst = append(dst, byte(info.GetVuGeneration()))
	}

	return dst, nil
}
