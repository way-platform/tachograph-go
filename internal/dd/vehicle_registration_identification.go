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
func (opts UnmarshalOptions) UnmarshalVehicleRegistrationIdentification(data []byte) (*ddv1.VehicleRegistrationIdentification, error) {
	// Some Gen2 V1 Vehicle Unit firmwares emit a 14-byte VRI record
	// (1 nation + 13 raw registration bytes) without an explicit codepage
	// prefix, while the canonical layout documented above is 15 bytes
	// (1 nation + 1 codepage + 13 number). Accept both sizes.
	const (
		lenVRIWithCodepage    = 15 // 1 nation + 1 codepage + 13 number
		lenVRIWithoutCodepage = 14 // 1 nation + 13 number (compact, older format)
	)

	if len(data) != lenVRIWithCodepage && len(data) != lenVRIWithoutCodepage {
		return nil, fmt.Errorf(
			"invalid data length for VehicleRegistrationIdentification: got %d, want %d or %d",
			len(data), lenVRIWithoutCodepage, lenVRIWithCodepage,
		)
	}

	vrn := &ddv1.VehicleRegistrationIdentification{}

	// Parse nation (always 1 byte at offset 0)
	nationValue := int32(data[0])
	nation := ddv1.NationNumeric(nationValue)
	vrn.SetNation(nation)

	var numberBytes []byte
	if len(data) == lenVRIWithCodepage {
		// Standard format: pass the 14-byte StringValue (codepage + 13 data bytes) as-is.
		numberBytes = data[1:15]
	} else {
		// Compact format: no explicit codepage. Synthesize a StringValue by
		// prepending codepage 0 (default) to the 13 data bytes.
		synth := make([]byte, 14)
		synth[0] = 0
		copy(synth[1:], data[1:14])
		numberBytes = synth
	}

	number, err := opts.UnmarshalStringValue(numberBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration number: %w", err)
	}
	vrn.SetNumber(number)

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
