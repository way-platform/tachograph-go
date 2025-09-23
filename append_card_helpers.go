package tachograph

import (
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toBCD(val int) byte {
	return byte(((val / 10) << 4) | (val % 10))
}

func appendDatef(dst []byte, t *timestamppb.Timestamp) []byte {
	if t == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := t.AsTime().Year()
	month := int(t.AsTime().Month())
	day := t.AsTime().Day()

	dst = append(dst, toBCD(year/100), toBCD(year%100))
	dst = append(dst, toBCD(month), toBCD(day))
	return dst
}

func appendString(dst []byte, s string, fixedLen int) []byte {
	// Check if this is a hex string (even length, all hex digits)
	if len(s)%2 == 0 && len(s) <= fixedLen*2 && isHexString(s) {
		// Convert hex string back to binary
		b := hexStringToBytes(s)
		if len(b) > fixedLen {
			b = b[:fixedLen]
		}
		dst = append(dst, b...)
		// Pad with spaces to match original string field padding style
		for i := len(b); i < fixedLen; i++ {
			dst = append(dst, ' ')
		}
		return dst
	}

	// Treat as regular ASCII string
	b := []byte(s)
	if len(b) > fixedLen {
		b = b[:fixedLen]
	}
	dst = append(dst, b...)
	for i := len(b); i < fixedLen; i++ {
		dst = append(dst, ' ') // Pad with spaces
	}
	return dst
}

// isHexString checks if a string contains only hex digits AND has hex characters (A-F)
// This ensures we only convert actual hex-encoded binary data, not numeric strings
func isHexString(s string) bool {
	if len(s) == 0 {
		return false
	}

	hasHexChars := false
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
		// Check if we have actual hex characters (A-F, a-f)
		if (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f') {
			hasHexChars = true
		}
	}

	// Only treat as hex if it contains actual hex characters (A-F, a-f)
	// Pure numeric strings like "1069650899000001" should be treated as ASCII
	return hasHexChars
}

// parseHexByte parses a 2-character hex string into a single byte
func parseHexByte(hexStr string) (byte, error) {
	if len(hexStr) < 2 {
		return 0, fmt.Errorf("hex string too short: %s", hexStr)
	}
	// Take first 2 characters and parse as hex
	bytes := hexStringToBytes(hexStr[:2])
	if len(bytes) == 0 {
		return 0, fmt.Errorf("failed to parse hex string: %s", hexStr[:2])
	}
	return bytes[0], nil
}

// hexStringToBytes converts a hex string to bytes
func hexStringToBytes(s string) []byte {
	result := make([]byte, len(s)/2)
	for i := 0; i < len(result); i++ {
		high := hexDigitToByte(s[i*2])
		low := hexDigitToByte(s[i*2+1])
		result[i] = (high << 4) | low
	}
	return result
}

// hexDigitToByte converts a single hex digit character to its byte value
func hexDigitToByte(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	if c >= 'A' && c <= 'F' {
		return c - 'A' + 10
	}
	if c >= 'a' && c <= 'f' {
		return c - 'a' + 10
	}
	return 0
}

func appendTimeReal(dst []byte, t *timestamppb.Timestamp) []byte {
	var timeVal uint32
	if t != nil {
		timeVal = uint32(t.GetSeconds())
	}
	return binary.BigEndian.AppendUint32(dst, timeVal)
}

func appendOdometer(dst []byte, km int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(km))
	return append(dst, b[1:]...) // Append 3 bytes
}

// appendBCDNation appends a BCD-encoded nation code (1 byte)
// The input string should be a hex representation like "12" -> 0x12
func appendBCDNation(dst []byte, nationStr string) []byte {
	if len(nationStr) == 0 {
		return append(dst, 0x00)
	}

	// If it's a 2-character hex string (even if just digits), convert it back to byte
	if len(nationStr) == 2 && isAllHexDigits(nationStr) {
		b := hexStringToBytes(nationStr)
		if len(b) > 0 {
			return append(dst, b[0])
		}
	}

	// Fallback: treat as regular string and take first byte
	if len(nationStr) > 0 {
		return append(dst, nationStr[0])
	}

	return append(dst, 0x00)
}

// isAllHexDigits checks if string contains only hex digits (0-9, A-F, a-f)
func isAllHexDigits(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return len(s) > 0
}
