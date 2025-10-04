package dd

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalIa5StringValue unmarshals an IA5 (ASCII) string value from binary data.
// IA5 strings have a fixed encoding (ASCII) and may include padding.
//
// The data type `IA5String` is specified in the Data Dictionary, Section 2.89.
//
// ASN.1 Definition:
//
//	IA5String ::= OCTET STRING (SIZE(0..255))
func (opts UnmarshalOptions) UnmarshalIa5StringValue(input []byte) (*ddv1.Ia5StringValue, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("insufficient data for IA5 string value: %w", io.ErrUnexpectedEOF)
	}

	var output ddv1.Ia5StringValue
	output.SetRawData(input)
	output.SetLength(int32(len(input)))

	// Decode and trim the input bytes
	trimmed := trimSpaceAndZeroBytes(input)
	decoded := string(trimmed)

	// Ensure the result is valid UTF-8
	if !utf8.ValidString(decoded) {
		// Convert invalid UTF-8 sequences to replacement characters
		decoded = strings.ToValidUTF8(decoded, string(utf8.RuneError))
	}

	output.SetValue(decoded)

	return &output, nil
}

// AppendIa5StringValue appends an Ia5StringValue to dst.
//
// This function handles IA5String format (fixed-length ASCII strings) defined as:
//
//	IA5String ::= OCTET STRING (SIZE(N))
//
// Binary Layout: stringData (N bytes, space-padded)
//
// The length is taken from the Ia5StringValue's 'length' field.
//
// If 'raw_data' is available, it is used directly (for round-trip fidelity).
// Otherwise, the 'value' string is encoded as ASCII and padded with spaces to the specified length.
func AppendIa5StringValue(dst []byte, sv *ddv1.Ia5StringValue) ([]byte, error) {
	// Handle nil - return empty bytes (no code page for IA5)
	if sv == nil {
		if dst == nil {
			return []byte{}, nil
		}
		return dst, nil
	}

	// Validate that raw_data and length agree if both are provided
	if sv.HasRawData() && sv.HasLength() {
		rawData := sv.GetRawData()
		length := int(sv.GetLength())
		if len(rawData) != length {
			return nil, fmt.Errorf("raw_data length (%d) does not match length field (%d)", len(rawData), length)
		}
	}

	// Length is required for IA5String
	if !sv.HasLength() {
		return nil, fmt.Errorf("IA5String requires length field to be set")
	}

	length := int(sv.GetLength())

	// Prefer raw bytes if available and of correct length
	if sv.HasRawData() && len(sv.GetRawData()) == length {
		return append(dst, sv.GetRawData()...), nil
	}

	// Fallback: use value string and pad with spaces
	value := sv.GetValue()
	if len(value) > length {
		return nil, fmt.Errorf("string value '%s' is longer than the allowed length %d", value, length)
	}
	result := make([]byte, length)
	copy(result, []byte(value))
	for i := len(value); i < length; i++ {
		result[i] = ' '
	}
	return append(dst, result...), nil
}
