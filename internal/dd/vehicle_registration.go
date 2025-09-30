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
func (opts UnmarshalOptions) UnmarshalVehicleRegistration(data []byte) (*ddv1.VehicleRegistrationIdentification, error) {
	const lenVehicleRegistration = 15

	if len(data) != lenVehicleRegistration {
		return nil, fmt.Errorf("invalid data length for VehicleRegistrationIdentification: got %d, want %d", len(data), lenVehicleRegistration)
	}

	vehicleReg := &ddv1.VehicleRegistrationIdentification{}

	// Read nation code (1 byte) and convert using protocol annotations
	if nation, err := UnmarshalEnum[ddv1.NationNumeric](data[0]); err == nil {
		vehicleReg.SetNation(nation)
	} else {
		// Value not recognized - set UNRECOGNIZED (no unrecognized field for this type)
		vehicleReg.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}

	// Read registration number (14 bytes)
	regNumber, err := opts.UnmarshalIA5StringValue(data[1:lenVehicleRegistration])
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
	var nationByte byte
	if nation == ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED {
		// UNRECOGNIZED values should not occur during marshalling
		return nil, fmt.Errorf("cannot marshal UNRECOGNIZED nation (no unrecognized field)")
	} else {
		var err error
		nationByte, err = MarshalEnum(nation)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal nation: %w", err)
		}
	}
	dst = append(dst, nationByte)

	// Append registration number (14 bytes, padded with spaces)
	number := vehicleReg.GetNumber()
	if number == nil {
		// Create empty StringValue with correct length for VehicleRegistrationNumber (SIZE(14))
		number = &ddv1.StringValue{}
		number.SetValue("")
		number.SetLength(14)
		number.SetEncoding(ddv1.Encoding_IA5)
	}
	return AppendStringValue(dst, number)
}
