package tachograph

import (
	"encoding/binary"
	"io"
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

// readOdometerFromBytes reads an odometer value (3 bytes) from a byte slice at the given offset.
//
// The data type `OdometerShort` is specified in the Data Dictionary, Section 2.113.
//
// ASN.1 Definition:
//
//     OdometerShort ::= INTEGER(0..999999)
//
// Binary Layout (3 bytes):
//   - Odometer Value (3 bytes): Big-endian unsigned integer
func readOdometerFromBytes(data []byte, offset int) (uint32, int, error) {
	if offset+3 > len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	// Convert 3-byte big-endian to uint32
	value := uint32(data[offset])<<16 | uint32(data[offset+1])<<8 | uint32(data[offset+2])
	return value, offset + 3, nil
}

