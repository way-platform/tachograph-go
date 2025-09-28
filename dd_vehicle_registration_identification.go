package tachograph

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalVehicleRegistrationIdentification parses vehicle registration identification data.
//
// The data type `VehicleRegistrationIdentification` is specified in the Data Dictionary, Section 2.166.
//
// ASN.1 Definition:
//
//	VehicleRegistrationIdentification ::= SEQUENCE {
//	    vehicleRegistrationNation NationNumeric,
//	    vehicleRegistrationNumber VehicleRegistrationNumber
//	}
//
// Binary Layout (15 bytes total):
//   - Nation (1 byte): NationNumeric
//   - Registration Number (14 bytes): VehicleRegistrationNumber (IA5String)
func unmarshalVehicleRegistrationIdentification(data []byte) (*ddv1.VehicleRegistrationIdentification, error) {
	const (
		lenVehicleRegistrationIdentification = 15 // 1 byte nation + 14 bytes number
	)

	if len(data) < lenVehicleRegistrationIdentification {
		return nil, fmt.Errorf("insufficient data for VehicleRegistrationIdentification: got %d, want %d", len(data), lenVehicleRegistrationIdentification)
	}

	// Parse nation (1 byte)
	nation, err := unmarshalNationNumeric(data[0:1])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}

	// Parse registration number (14 bytes)
	regNumber, err := unmarshalIA5StringValue(data[1:15])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}

	// Create VehicleRegistrationIdentification structure
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)
	vehicleReg.SetNumber(regNumber)

	return vehicleReg, nil
}

// appendVehicleRegistrationIdentification appends vehicle registration identification data to dst.
//
// The data type `VehicleRegistrationIdentification` is specified in the Data Dictionary, Section 2.166.
//
// ASN.1 Definition:
//
//	VehicleRegistrationIdentification ::= SEQUENCE {
//	    vehicleRegistrationNation NationNumeric,
//	    vehicleRegistrationNumber VehicleRegistrationNumber
//	}
//
// Binary Layout (15 bytes total):
//   - Nation (1 byte): NationNumeric
//   - Registration Number (14 bytes): VehicleRegistrationNumber (IA5String)
func appendVehicleRegistrationIdentification(dst []byte, vehicleReg *ddv1.VehicleRegistrationIdentification) ([]byte, error) {
	if vehicleReg == nil {
		// Append default values: 1 byte nation (0xFF) + 14 bytes registration number (spaces)
		dst = append(dst, 0xFF)
		return appendString(dst, "", 14)
	}

	// Append nation (1 byte)
	dst = append(dst, byte(vehicleReg.GetNation()))

	// Append registration number (14 bytes, padded with spaces)
	number := vehicleReg.GetNumber()
	if number != nil {
		return appendString(dst, number.GetDecoded(), 14)
	}
	return appendString(dst, "", 14)
}
