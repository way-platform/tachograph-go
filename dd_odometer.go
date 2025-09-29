package tachograph

import (
	"encoding/binary"
)

// appendOdometer appends a 3-byte odometer value.
//
// The data type `OdometerShort` is specified in the Data Dictionary, Section 2.113.
//
// ASN.1 Definition:
//
//     OdometerShort ::= INTEGER(0..999999)
//
// Binary Layout (3 bytes):
//   - Odometer Value (3 bytes): Big-endian unsigned integer
func appendOdometer(dst []byte, odometer uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, odometer)
	return append(dst, b[1:]...)
}

