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
//   - Registration number (14 bytes): codePage (1 byte) + vehicleRegNumber (13 bytes)
//
// Some Gen2 V1 vehicle units emit 14-byte records (VehicleRegistrationNumber only,
// without the nation byte). When 14 bytes are provided, the nation defaults to
// NATION_NUMERIC_UNSPECIFIED. Note: the marshal function always produces the canonical
// 15-byte format, so 14-byte inputs are normalized to 15 bytes on round-trip.
func (opts UnmarshalOptions) UnmarshalVehicleRegistration(data []byte) (*ddv1.VehicleRegistrationIdentification, error) {
	const (
		lenWithNation    = 15 // VehicleRegistrationIdentification: nation (1) + VRN (14)
		lenWithoutNation = 14 // VehicleRegistrationNumber only: codePage (1) + regNumber (13)
	)

	vehicleReg := &ddv1.VehicleRegistrationIdentification{}

	switch len(data) {
	case lenWithNation:
		// Standard: 1 byte nation + 14 bytes VehicleRegistrationNumber
		if nation, err := UnmarshalEnum[ddv1.NationNumeric](data[0]); err == nil {
			vehicleReg.SetNation(nation)
		} else {
			vehicleReg.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
		}

		regNumber, err := opts.UnmarshalStringValue(data[1:lenWithNation])
		if err != nil {
			return nil, fmt.Errorf("failed to parse registration number: %w", err)
		}
		vehicleReg.SetNumber(regNumber)

	case lenWithoutNation:
		// Some Gen2 V1 VUs omit the nation byte, emitting only VehicleRegistrationNumber (14 bytes)
		vehicleReg.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED)

		regNumber, err := opts.UnmarshalStringValue(data[0:lenWithoutNation])
		if err != nil {
			return nil, fmt.Errorf("failed to parse registration number: %w", err)
		}
		vehicleReg.SetNumber(regNumber)

	default:
		return nil, fmt.Errorf("invalid data length for VehicleRegistrationIdentification: got %d, want %d or %d",
			len(data), lenWithNation, lenWithoutNation)
	}

	return vehicleReg, nil
}

// MarshalVehicleRegistration marshals a VehicleRegistrationIdentification to bytes.
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
//   - Registration number (14 bytes): codePage (1 byte) + vehicleRegNumber (13 bytes)
func (opts MarshalOptions) MarshalVehicleRegistration(vehicleReg *ddv1.VehicleRegistrationIdentification) ([]byte, error) {
	if vehicleReg == nil {
		return nil, fmt.Errorf("vehicleRegistration cannot be nil")
	}

	var dst []byte

	// Marshal nation (1 byte) - get protocol value from enum
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

	// Marshal registration number (14 bytes: 1 byte code page + 13 bytes string data)
	number := vehicleReg.GetNumber()
	if number == nil {
		// Create empty StringValue with correct length for VehicleRegistrationNumber
		// Default to code page 0 (no national characters) with 13 bytes of data
		number = &ddv1.StringValue{}
		number.SetValue("")
		number.SetLength(13) // Length of the string data, not including code page byte
		number.SetEncoding(ddv1.Encoding_ENCODING_DEFAULT)
	}
	numberBytes, err := opts.MarshalStringValue(number)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle registration number: %w", err)
	}
	return append(dst, numberBytes...), nil
}
