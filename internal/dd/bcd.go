package dd

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// decodeBCD converts BCD-encoded bytes to an integer
func decodeBCD(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	s := hex.EncodeToString(b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid BCD value: %s", s)
	}
	return int(i), nil
}

// encodeBCD converts an integer to BCD-encoded bytes with a specified length.
// The length parameter specifies the number of bytes to produce.
// For example, encodeBCD(123, 2) produces []byte{0x01, 0x23}.
func encodeBCD(value int, length int) ([]byte, error) {
	if value < 0 {
		return nil, fmt.Errorf("cannot encode negative value as BCD: %d", value)
	}
	
	// Convert to decimal string
	s := strconv.Itoa(value)
	
	// Check if value fits in the specified length
	maxDigits := length * 2
	if len(s) > maxDigits {
		return nil, fmt.Errorf("value %d requires more than %d bytes (has %d digits, max %d)", value, length, len(s), maxDigits)
	}
	
	// Pad with leading zeros to fill the length
	for len(s) < maxDigits {
		s = "0" + s
	}
	
	// Convert hex string to bytes
	result, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to encode BCD: %w", err)
	}
	
	return result, nil
}
