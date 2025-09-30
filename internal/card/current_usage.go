package card

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
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
func unmarshalCardCurrentUsage(data []byte) (*cardv1.CurrentUsage, error) {
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
	sessionOpenTime, err := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse session open time: %w", err)
	}
	target.SetSessionOpenTime(sessionOpenTime)
	offset += 4

	// Read session open vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration nation")
	}
	nationByte := data[offset]
	offset++

	// Create VehicleRegistrationIdentification structure
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	if enumNum, found := dd.GetEnumForProtocolValue(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(), int32(nationByte)); found {
		vehicleReg.SetNation(ddv1.NationNumeric(enumNum))
	} else {
		vehicleReg.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}

	if offset+14 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration number")
	}
	regNumber, err := dd.UnmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	// offset += 14 // Not needed as this is the last field
	vehicleReg.SetNumber(regNumber)
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
	data = dd.AppendTimeReal(data, currentUsage.GetSessionOpenTime())

	// Session open vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	vehicleReg := currentUsage.GetSessionOpenVehicle()
	if vehicleReg != nil {
		var err error
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
