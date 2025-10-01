package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPreviousVehicleInfoG2 unmarshals a Generation 2 PreviousVehicleInfo (20 bytes).
//
// The data type `PreviousVehicleInfo` (Gen2 variant) is specified in the Data Dictionary, Section 2.118.
//
// ASN.1 Definition (Gen2):
//
//	PreviousVehicleInfo ::= SEQUENCE {
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,
//	    cardWithdrawalTime                TimeReal,
//	    vuGeneration                      Generation
//	}
func (opts UnmarshalOptions) UnmarshalPreviousVehicleInfoG2(data []byte) (*ddv1.PreviousVehicleInfoG2, error) {
	const (
		idxVehicleReg          = 0
		idxCardWithdrawalTime  = 15
		idxVuGeneration        = 19
		lenPreviousVehicleInfo = 20 // Fixed size for Gen2
	)

	if len(data) != lenPreviousVehicleInfo {
		return nil, fmt.Errorf("invalid data length for Gen2 PreviousVehicleInfo: got %d, want %d", len(data), lenPreviousVehicleInfo)
	}

	info := &ddv1.PreviousVehicleInfoG2{}
	info.SetRawData(data)

	// Parse vehicleRegistrationIdentification (15 bytes)
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[idxVehicleReg : idxVehicleReg+15])
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

	// Parse vuGeneration (1 byte)
	vuGen, err := UnmarshalEnum[ddv1.Generation](data[idxVuGeneration])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vu generation: %w", err)
	}
	info.SetVuGeneration(vuGen)

	return info, nil
}

// AppendPreviousVehicleInfoG2 appends a Generation 2 PreviousVehicleInfo (20 bytes).
func AppendPreviousVehicleInfoG2(dst []byte, info *ddv1.PreviousVehicleInfoG2) ([]byte, error) {
	const lenPreviousVehicleInfo = 20 // Fixed size for Gen2

	// Use raw data painting strategy if available
	var canvas [lenPreviousVehicleInfo]byte
	if rawData := info.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenPreviousVehicleInfo {
			return nil, fmt.Errorf("invalid raw_data length for PreviousVehicleInfoG2: got %d, want %d", len(rawData), lenPreviousVehicleInfo)
		}
		copy(canvas[:], rawData)
	}

	// Paint semantic values over the canvas
	var err error

	// Vehicle registration (15 bytes)
	vehicleRegBytes, err := AppendVehicleRegistration(nil, info.GetVehicleRegistration())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle registration: %w", err)
	}
	copy(canvas[0:15], vehicleRegBytes)

	// Card withdrawal time (4 bytes)
	timeBytes, err := AppendTimeReal(nil, info.GetCardWithdrawalTime())
	if err != nil {
		return nil, fmt.Errorf("failed to append card withdrawal time: %w", err)
	}
	copy(canvas[15:19], timeBytes)

	// VU generation (1 byte)
	vuGenByte, err := MarshalEnum(info.GetVuGeneration())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vu generation: %w", err)
	}
	canvas[19] = vuGenByte

	return append(dst, canvas[:]...), nil
}
