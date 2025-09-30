package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalBcdString parses BCD string data.
//
// The data type `BcdString` is specified in the Data Dictionary, Section 2.7.
//
// ASN.1 Definition:
//
//	BCDString ::= CharacterStringType
//
// Binary Layout (variable length):
//   - BCD String (variable): BCD-encoded bytes
func UnmarshalBcdString(data []byte) (*ddv1.BcdString, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for BcdString: got %d, want at least 1", len(data))
	}
	decoded, err := decodeBCD(data)
	if err != nil {
		return nil, err
	}
	var output ddv1.BcdString
	output.SetRawData(data)
	output.SetValue(int32(decoded))
	return &output, nil
}

// AppendBcdString appends BCD string data to dst.
//
// The data type `BcdString` is specified in the Data Dictionary, Section 2.7.
//
// ASN.1 Definition:
//
//	BCDString ::= CharacterStringType
//
// Binary Layout (variable length):
//   - BCD String (variable): BCD-encoded bytes
func AppendBcdString(dst []byte, bcdString *ddv1.BcdString) ([]byte, error) {
	// Handle nil BcdString - return empty (nothing to append)
	if bcdString == nil {
		return dst, nil
	}

	// Get the semantic value
	value := bcdString.GetValue()
	intValue := int(value)
	if intValue < 0 {
		return nil, fmt.Errorf("cannot encode negative BCD value: %d", value)
	}

	// Use raw_data as canvas if available (raw data painting approach)
	rawData := bcdString.GetRawData()
	if len(rawData) > 0 {
		// Paint semantic value over raw_data canvas
		// Use raw_data length to preserve the exact byte count (may include padding)
		encodedBytes, err := encodeBCD(intValue, len(rawData))
		if err != nil {
			return nil, fmt.Errorf("failed to encode BCD value %d into %d bytes: %w", intValue, len(rawData), err)
		}
		return append(dst, encodedBytes...), nil
	}

	// Fall back to encoding from value with minimal byte count
	// Count digits (each byte holds 2 digits)
	// For zero, we need 1 digit
	digitCount := 1
	if intValue > 0 {
		digitCount = 0
		temp := intValue
		for temp > 0 {
			digitCount++
			temp /= 10
		}
	}

	// Calculate required bytes (round up to nearest byte)
	byteCount := (digitCount + 1) / 2

	// Encode the value
	encodedBytes, err := encodeBCD(intValue, byteCount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode BCD string: %w", err)
	}

	return append(dst, encodedBytes...), nil
}
