package dd

import (
	"encoding/binary"
	"fmt"
)

// appendOdometer appends a 3-byte odometer value.
//
// The data type `OdometerShort` is specified in the Data Dictionary, Section 2.113.
//
// ASN.1 Definition:
//
//	OdometerShort ::= INTEGER(0..999999)
//
// Binary Layout (3 bytes):
//   - Odometer Value (3 bytes): Big-endian unsigned integer
func AppendOdometer(dst []byte, odometer uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, odometer)
	return append(dst, b[1:]...)
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
func UnmarshalOdometer(data []byte) (uint32, error) {
	const lenOdometerShort = 3

	if len(data) != lenOdometerShort {
		return 0, fmt.Errorf("invalid data length for OdometerShort: got %d, want %d", len(data), lenOdometerShort)
	}

	// Convert 3-byte big-endian to uint32
	value := uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	return value, nil
}
