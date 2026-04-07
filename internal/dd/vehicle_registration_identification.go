package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVehicleRegistrationIdentification parses the VehicleRegistrationIdentification structure.
//
// See Data Dictionary, Section 2.166, `VehicleRegistrationIdentification`.
//
// ASN.1 Specification:
//
//	VehicleRegistrationIdentification ::= SEQUENCE {
//	    vehicleRegistrationNation NationNumeric,     -- 1 byte
//	    vehicleRegistrationNumber VehicleRegistrationNumber  -- 14 bytes (1 code page + 13 data)
//	}
//
// Binary Layout (fixed length: 15 bytes):
//   - Vehicle Registration Nation (1 byte): NationNumeric
//   - Vehicle Registration Number (14 bytes): StringValue (1 byte code page + 13 bytes data)
//
// Some Gen2 V1 vehicle units emit 14-byte records (VehicleRegistrationNumber only,
// without the nation byte). When 14 bytes are provided, the nation defaults to
// NATION_NUMERIC_UNSPECIFIED. Note: the marshal function always produces the canonical
// 15-byte format, so 14-byte inputs are normalized to 15 bytes on round-trip.
func (opts UnmarshalOptions) UnmarshalVehicleRegistrationIdentification(data []byte) (*ddv1.VehicleRegistrationIdentification, error) {
	const (
		lenWithNation    = 15 // VehicleRegistrationIdentification: nation (1) + VRN (14)
		lenWithoutNation = 14 // VehicleRegistrationNumber only: codePage (1) + regNumber (13)
	)

	vrn := &ddv1.VehicleRegistrationIdentification{}

	switch len(data) {
	case lenWithNation:
		// Standard: 1 byte nation + 14 bytes VehicleRegistrationNumber
		nationValue := int32(data[0])
		nation := ddv1.NationNumeric(nationValue)
		vrn.SetNation(nation)

		number, err := opts.UnmarshalStringValue(data[1:lenWithNation])
		if err != nil {
			return nil, fmt.Errorf("failed to parse vehicle registration number: %w", err)
		}
		vrn.SetNumber(number)

	case lenWithoutNation:
		// Some Gen2 V1 VUs omit the nation byte, emitting only VehicleRegistrationNumber (14 bytes)
		vrn.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED)

		number, err := opts.UnmarshalStringValue(data[0:lenWithoutNation])
		if err != nil {
			return nil, fmt.Errorf("failed to parse vehicle registration number: %w", err)
		}
		vrn.SetNumber(number)

	default:
		return nil, fmt.Errorf(
			"invalid data length for VehicleRegistrationIdentification: got %d, want %d or %d",
			len(data), lenWithNation, lenWithoutNation,
		)
	}

	return vrn, nil
}

// MarshalVehicleRegistrationIdentification marshals the VehicleRegistrationIdentification structure.
//
// See Data Dictionary, Section 2.166, `VehicleRegistrationIdentification`.
func (opts MarshalOptions) MarshalVehicleRegistrationIdentification(vrn *ddv1.VehicleRegistrationIdentification) ([]byte, error) {
	if vrn == nil {
		return nil, fmt.Errorf("vrn cannot be nil")
	}

	const size = 15
	var canvas [size]byte

	offset := 0

	// Marshal nation (1 byte)
	canvas[offset] = byte(vrn.GetNation())
	offset += 1

	// Marshal registration number (14 bytes)
	numberBytes, err := opts.MarshalStringValue(vrn.GetNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle registration number: %w", err)
	}
	if len(numberBytes) != 14 {
		return nil, fmt.Errorf(
			"invalid vehicle registration number length: got %d, want 14",
			len(numberBytes),
		)
	}
	copy(canvas[offset:offset+14], numberBytes)
	offset += 14

	if offset != size {
		return nil, fmt.Errorf(
			"VehicleRegistrationIdentification marshalling size mismatch: wrote %d bytes, expected %d",
			offset, size,
		)
	}

	return canvas[:], nil
}

// AnonymizeVehicleRegistrationIdentification anonymizes vehicle registration data.
func (opts AnonymizeOptions) AnonymizeVehicleRegistrationIdentification(vreg *ddv1.VehicleRegistrationIdentification) *ddv1.VehicleRegistrationIdentification {
	if vreg == nil {
		return nil
	}

	result := &ddv1.VehicleRegistrationIdentification{}
	// Preserve country (structural info)
	result.SetNation(vreg.GetNation())
	// Anonymize the registration number
	result.SetNumber(opts.AnonymizeStringValue(vreg.GetNumber()))
	return result
}
