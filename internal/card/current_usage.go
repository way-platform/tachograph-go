package card

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardCurrentUsage unmarshals current usage data from a card EF.
//
// The data type `CardCurrentUse` is specified in the Data Dictionary, Section 2.16.
//
// ASN.1 Definition:
//
//	CardCurrentUse ::= SEQUENCE {
//	    sessionOpenTime                   TimeReal,
//	    sessionOpenVehicle                VehicleRegistrationIdentification
//	}
func (opts UnmarshalOptions) unmarshalCurrentUsage(data []byte) (*cardv1.CurrentUsage, error) {
	const (
		lenCardCurrentUse = 19 // 4 bytes time + 15 bytes vehicle registration
	)

	if len(data) < lenCardCurrentUse {
		return nil, fmt.Errorf("insufficient data for current usage")
	}
	var target cardv1.CurrentUsage
	offset := 0

	// Read session open time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for session open time")
	}
	sessionOpenTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse session open time: %w", err)
	}
	target.SetSessionOpenTime(sessionOpenTime)
	offset += 4

	// Read session open vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	if offset+15 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration")
	}
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[offset : offset+15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration: %w", err)
	}
	// offset += 15 // Not needed as this is the last field
	target.SetSessionOpenVehicle(vehicleReg)
	return &target, nil
}

// AppendCurrentUsage appends current usage data to a byte slice.
//
// The data type `CardCurrentUse` is specified in the Data Dictionary, Section 2.16.
//
// ASN.1 Definition:
//
//	CardCurrentUse ::= SEQUENCE {
//	    sessionOpenTime                   TimeReal,
//	    sessionOpenVehicle                VehicleRegistrationIdentification
//	}
func appendCurrentUsage(data []byte, currentUsage *cardv1.CurrentUsage) ([]byte, error) {
	if currentUsage == nil {
		return data, nil
	}

	// Session open time (4 bytes)
	var err error
	data, err = dd.AppendTimeReal(data, currentUsage.GetSessionOpenTime())
	if err != nil {
		return data, fmt.Errorf("failed to append session open time: %w", err)
	}

	// Session open vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	vehicleReg := currentUsage.GetSessionOpenVehicle()
	if vehicleReg != nil {
		data, err = dd.AppendVehicleRegistration(data, vehicleReg)
		if err != nil {
			return nil, err
		}
	} else {
		// No vehicle registration - pad with zeros
		data = append(data, make([]byte, 15)...)
	}

	return data, nil
}
