package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVehicleRegistration unmarshals a VehicleRegistrationIdentification from a byte slice.
//
// The data type `VehicleRegistrationIdentification` is specified in the Data Dictionary, Section 2.187.
//
// ASN.1 Definition:
//
//	VehicleRegistrationIdentification ::= SEQUENCE {
//	    vehicleRegistrationNation    NationNumeric,    -- 1 byte
//	    vehicleRegistrationNumber    VehicleRegistrationNumber  -- 14 bytes
//	}
//
// Binary Layout (15 bytes):
//   - Nation code (1 byte): NationNumeric
//   - Registration number (14 bytes): IA5String
func UnmarshalVehicleRegistration(data []byte) (*ddv1.VehicleRegistrationIdentification, error) {
	const lenVehicleRegistration = 15

	if len(data) != lenVehicleRegistration {
		return nil, fmt.Errorf("invalid data length for VehicleRegistrationIdentification: got %d, want %d", len(data), lenVehicleRegistration)
	}

	vehicleReg := &ddv1.VehicleRegistrationIdentification{}

	// Read nation code (1 byte) and convert using protocol annotations
	nationByte := data[0]
	if enumNum, found := GetEnumForProtocolValue(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(), int32(nationByte)); found {
		vehicleReg.SetNation(ddv1.NationNumeric(enumNum))
	} else {
		vehicleReg.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}

	// Read registration number (14 bytes)
	regNumber, err := UnmarshalIA5StringValue(data[1:lenVehicleRegistration])
	if err != nil {
		return nil, fmt.Errorf("failed to parse registration number: %w", err)
	}
	vehicleReg.SetNumber(regNumber)

	return vehicleReg, nil
}

// AppendVehicleRegistration appends a VehicleRegistrationIdentification to dst.
//
// The data type `VehicleRegistrationIdentification` is specified in the Data Dictionary, Section 2.187.
//
// ASN.1 Definition:
//
//	VehicleRegistrationIdentification ::= SEQUENCE {
//	    vehicleRegistrationNation    NationNumeric,    -- 1 byte
//	    vehicleRegistrationNumber    VehicleRegistrationNumber  -- 14 bytes
//	}
//
// Binary Layout (15 bytes):
//   - Nation code (1 byte): NationNumeric
//   - Registration number (14 bytes): IA5String
func AppendVehicleRegistration(dst []byte, vehicleReg *ddv1.VehicleRegistrationIdentification) ([]byte, error) {
	if vehicleReg == nil {
		return nil, fmt.Errorf("vehicleRegistration cannot be nil")
	}

	// Append nation (1 byte) - get protocol value from enum
	nation := vehicleReg.GetNation()
	if protocolValue, found := GetProtocolValueForEnum(nation); found {
		dst = append(dst, byte(protocolValue))
	} else {
		// Default to 0xFF (EMPTY) for unspecified/unrecognized
		dst = append(dst, 0xFF)
	}

	// Append registration number (14 bytes, padded with spaces)
	number := vehicleReg.GetNumber()
	if number != nil {
		return AppendString(dst, number.GetValue(), 14)
	}
	return AppendString(dst, "", 14)
}
