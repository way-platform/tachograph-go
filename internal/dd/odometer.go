package dd

import (
	"fmt"
)

// MarshalOdometer marshals a 3-byte odometer value.
//
// The data type `Odometer` is specified in the Data Dictionary, Section 2.99.
//
// ASN.1 Definition:
//
//	Odometer ::= INTEGER (0..2^24-1)
//
// Binary Layout (3 bytes):
//   - Odometer reading in km (3 bytes): Big-endian uint24
func (opts MarshalOptions) MarshalOdometer(odometer int32) ([]byte, error) {
	const lenOdometer = 3
	var canvas [lenOdometer]byte
	// Convert 32-bit int to 24-bit big-endian (use only lower 24 bits)
	canvas[0] = byte((odometer >> 16) & 0xFF)
	canvas[1] = byte((odometer >> 8) & 0xFF)
	canvas[2] = byte(odometer & 0xFF)
	return canvas[:], nil
}

// UnmarshalOdometer unmarshals a 3-byte odometer value.
//
// The data type `OdometerShort` is specified in the Data Dictionary, Section 2.113.
//
// ASN.1 Definition:
//
//	OdometerShort ::= INTEGER(0..999999)
//
// Binary Layout (3 bytes):
//   - Odometer Value (3 bytes): Big-endian unsigned integer
//
// Returns nil when the sentinel value 0xFFFFFF is present, indicating the odometer
// is not available. Callers must not set the proto field when nil is returned.
func (opts UnmarshalOptions) UnmarshalOdometer(data []byte) (*int32, error) {
	const (
		lenOdometerShort  = 3
		odometerSentinel  = 0xFFFFFF
	)

	if len(data) != lenOdometerShort {
		return nil, fmt.Errorf("invalid data length for OdometerShort: got %d, want %d", len(data), lenOdometerShort)
	}

	value := int32(data[0])<<16 | int32(data[1])<<8 | int32(data[2])
	if value == odometerSentinel {
		return nil, nil
	}
	return &value, nil
}

// AnonymizeOdometerValue anonymizes odometer values based on options.
// If PreserveDistanceAndTrips is false, rounds to nearest 1000km.
func (opts AnonymizeOptions) AnonymizeOdometerValue(km int32) int32 {
	if opts.PreserveDistanceAndTrips {
		return km
	}

	// Round to nearest 1000km
	return (km / 1000) * 1000
}
