package tachograph

import (
	"bytes"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// Data dictionary string parsing and serialization functions
// These functions are used across multiple files for consistent string handling

// readString reads a fixed-length string from a bytes.Reader
func readString(r *bytes.Reader, len int) string {
	b := make([]byte, len)
	_, _ = r.Read(b) // ignore error as we're reading from in-memory buffer
	// Trim trailing spaces and null bytes
	b = bytes.TrimRight(b, " \x00")

	// Check if the result is valid UTF-8, if not convert to hex representation
	if !isValidUTF8(b) {
		return bytesToHexString(b)
	}

	return string(b)
}

// isValidUTF8 checks if the byte slice contains valid UTF-8
func isValidUTF8(b []byte) bool {
	// Check if all bytes are printable ASCII or valid UTF-8
	for _, byte := range b {
		if byte < 0x20 || byte > 0x7E {
			// Contains non-printable characters, treat as binary
			return false
		}
	}
	return true
}

// bytesToHexString converts binary data to a hex string representation
func bytesToHexString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	result := make([]byte, len(b)*2)
	const hexDigits = "0123456789ABCDEF"
	for i, byte := range b {
		result[i*2] = hexDigits[byte>>4]
		result[i*2+1] = hexDigits[byte&0x0F]
	}
	return string(result)
}

// appendString appends a fixed-length string, padding with spaces.
func appendString(dst []byte, s string, length int) ([]byte, error) {
	if len(s) > length {
		return nil, fmt.Errorf("string '%s' is longer than the allowed length %d", s, length)
	}
	result := make([]byte, length)
	copy(result, []byte(s))
	for i := len(s); i < length; i++ {
		result[i] = ' '
	}
	return append(dst, result...), nil
}

// appendStringValue appends a fixed-length string from a StringValue, padding with spaces.
func appendStringValue(dst []byte, sv *ddv1.StringValue, length int) ([]byte, error) {
	if sv == nil {
		return appendString(dst, "", length)
	}
	// Use the decoded string if available, otherwise use the encoded bytes as string
	s := sv.GetDecoded()
	if s == "" {
		s = string(sv.GetEncoded())
	}
	return appendString(dst, s, length)
}

// appendBCDNation appends a BCD-encoded nation number.
func appendBCDNation(dst []byte, nation string) ([]byte, error) {
	// This is a placeholder. A real implementation would convert the nation string
	// to its numeric code and then to BCD.
	return append(dst, 0), nil // Append a single zero byte for now
}

// appendVehicleRegistration appends vehicle registration from VehicleRegistrationIdentification
func appendVehicleRegistration(dst []byte, vehicleReg *ddv1.VehicleRegistrationIdentification) ([]byte, error) {
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
