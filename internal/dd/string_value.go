package dd

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalStringValue unmarshals a code-paged string value from binary data.
// The input should contain a code page byte followed by the encoded string data.
//
// The data type `StringValue` is specified in the Data Dictionary, Section 2.158.
//
// ASN.1 Definition:
//
//	StringValue ::= SEQUENCE {
//	    codePage    OCTET STRING (SIZE(1)),
//	    stringData  OCTET STRING (SIZE(0..255))
//	}
func UnmarshalStringValue(input []byte) (*ddv1.StringValue, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("insufficient data for string value: %w", io.ErrUnexpectedEOF)
	}

	codePage := input[0]
	data := input[1:]

	var output ddv1.StringValue
	output.SetEncoding(getEncodingFromCodePage(codePage))
	output.SetRawData(data)
	output.SetLength(uint32(len(data)))

	// Decode the string based on the code page
	decoded, err := decodeWithCodePage(codePage, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode string with code page %d: %w", codePage, err)
	}
	output.SetValue(decoded)

	return &output, nil
}

// UnmarshalIA5StringValue unmarshals an IA5 (ASCII) string value from binary data.
// IA5 strings have a fixed encoding and may include padding.
//
// The data type `IA5String` is specified in the Data Dictionary, Section 2.89.
//
// ASN.1 Definition:
//
//	IA5String ::= OCTET STRING (SIZE(0..255))
func UnmarshalIA5StringValue(input []byte) (*ddv1.StringValue, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("insufficient data for string value: %w", io.ErrUnexpectedEOF)
	}

	var output ddv1.StringValue
	output.SetEncoding(ddv1.Encoding_IA5)
	output.SetRawData(input)
	output.SetLength(uint32(len(input))) // Store the length for self-describing marshalling

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

// AppendStringValue appends a StringValue to dst.
//
// This function handles two ASN.1 formats based on the encoding:
//
//  1. IA5String format (Encoding_IA5):
//     Fixed-length ASCII strings defined as:
//     IA5String ::= OCTET STRING (SIZE(N))
//     Binary Layout: stringData (N bytes, space-padded)
//     The length is taken from the StringValue's 'length' field.
//
//  2. Code-paged format (other encodings):
//     Variable-length strings with code page prefix, defined as:
//     StringValue ::= SEQUENCE {
//     codePage    OCTET STRING (SIZE(1)),
//     stringData  OCTET STRING (SIZE(0..255))
//     }
//     Binary Layout: codePage (1 byte) + stringData (variable bytes)
//
// The function prefers to use 'raw_data' bytes if available (for round-trip fidelity),
// otherwise it encodes the 'value' string using the specified encoding.
//
// If both raw_data and length are provided, they must agree or an error is returned.
func AppendStringValue(dst []byte, sv *ddv1.StringValue) ([]byte, error) {
	// Handle nil
	if sv == nil {
		// Empty string value: code page 255 (EMPTY) + no data
		return append(dst, 0xFF), nil
	}

	// Validate that raw_data and length agree if both are provided
	if sv.HasRawData() && sv.HasLength() {
		rawData := sv.GetRawData()
		length := int(sv.GetLength())
		if len(rawData) != length {
			return nil, fmt.Errorf("raw_data length (%d) does not match length field (%d)", len(rawData), length)
		}
	}

	// Handle IA5String format (fixed-length, no code page)
	if sv.GetEncoding() == ddv1.Encoding_IA5 {
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

	// Handle code-paged format (variable-length with code page byte)

	// Determine the code page byte
	codePage := getCodePageFromEncoding(sv.GetEncoding())

	// Prefer raw bytes if available (round-trip fidelity)
	if raw := sv.GetRawData(); len(raw) > 0 {
		dst = append(dst, codePage)
		return append(dst, raw...), nil
	}

	// Fallback: encode the value string
	value := sv.GetValue()
	if value == "" {
		// Empty string
		return append(dst, codePage), nil
	}

	// For now, we only support encoding to ISO-8859-1 and IA5 (ASCII)
	// More sophisticated encoding would require charmap encoding
	encoded, err := encodeWithCodePage(codePage, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode string with code page %d: %w", codePage, err)
	}

	dst = append(dst, codePage)
	return append(dst, encoded...), nil
}

// getEncodingFromCodePage maps a code page byte to the corresponding Encoding enum.
func getEncodingFromCodePage(codePage byte) ddv1.Encoding {
	switch codePage {
	case 0:
		return ddv1.Encoding_ENCODING_DEFAULT
	case 1:
		return ddv1.Encoding_ISO_8859_1
	case 2:
		return ddv1.Encoding_ISO_8859_2
	case 3:
		return ddv1.Encoding_ISO_8859_3
	case 5:
		return ddv1.Encoding_ISO_8859_5
	case 7:
		return ddv1.Encoding_ISO_8859_7
	case 9:
		return ddv1.Encoding_ISO_8859_9
	case 13:
		return ddv1.Encoding_ISO_8859_13
	case 15:
		return ddv1.Encoding_ISO_8859_15
	case 16:
		return ddv1.Encoding_ISO_8859_16
	case 80:
		return ddv1.Encoding_KOI8_R
	case 85:
		return ddv1.Encoding_KOI8_U
	case 255:
		return ddv1.Encoding_ENCODING_EMPTY
	default:
		return ddv1.Encoding_ENCODING_UNRECOGNIZED
	}
}

// trimSpaceAndZeroBytes trims spaces, 0x00 and 0xff values off a byte slice
func trimSpaceAndZeroBytes(b []byte) []byte {
	// Define cutset as string - bytes.Trim handles this properly
	cutset := "\t\n\v\f\r \x85\xA0\x00\xFF"
	return bytes.Trim(b, cutset)
}

// decodeWithCodePage decodes a byte slice with the given code page, returns the trimmed decoded string
func decodeWithCodePage(codePage byte, data []byte) (string, error) {
	if codePage == 255 {
		// codepage 255 means empty/unassigned string
		return "", nil
	}

	// Check if the data contains any valid characters
	ok := false
	for i := 0; i < len(data); i++ {
		if data[i] > 0 && data[i] < 255 {
			ok = true
			break
		}
	}
	if !ok {
		return "", nil
	}

	// Map code page to character map
	var cmap *charmap.Charmap
	switch codePage {
	case 0:
		// Default to ISO-8859-1 for code page 0 (ASCII-compatible)
		cmap = charmap.ISO8859_1
	case 1:
		cmap = charmap.ISO8859_1
	case 2:
		cmap = charmap.ISO8859_2
	case 3:
		cmap = charmap.ISO8859_3
	case 5:
		cmap = charmap.ISO8859_5
	case 7:
		cmap = charmap.ISO8859_7
	case 9:
		cmap = charmap.ISO8859_9
	case 13:
		cmap = charmap.ISO8859_13
	case 15:
		cmap = charmap.ISO8859_15
	case 16:
		cmap = charmap.ISO8859_16
	case 80:
		cmap = charmap.KOI8R
	case 85:
		cmap = charmap.KOI8U
	default:
		// For unrecognized code pages, fall back to ISO-8859-1
		cmap = charmap.ISO8859_1
	}

	dec := cmap.NewDecoder()
	res, err := dec.String(string(data))
	if err != nil {
		return "", fmt.Errorf("could not decode code page %d string: %w", codePage, err)
	}

	// The character map decoder should produce valid UTF-8, but let's be safe
	trimmed := string(trimSpaceAndZeroBytes([]byte(res)))

	// If the result is not valid UTF-8, convert it to valid UTF-8
	if !utf8.ValidString(trimmed) {
		// Convert invalid UTF-8 sequences to replacement characters
		trimmed = strings.ToValidUTF8(trimmed, string(utf8.RuneError))
	}

	return trimmed, nil
}

// getCodePageFromEncoding maps an Encoding enum to a code page byte.
func getCodePageFromEncoding(encoding ddv1.Encoding) byte {
	switch encoding {
	case ddv1.Encoding_ENCODING_DEFAULT:
		return 0
	case ddv1.Encoding_ISO_8859_1:
		return 1
	case ddv1.Encoding_ISO_8859_2:
		return 2
	case ddv1.Encoding_ISO_8859_3:
		return 3
	case ddv1.Encoding_ISO_8859_5:
		return 5
	case ddv1.Encoding_ISO_8859_7:
		return 7
	case ddv1.Encoding_ISO_8859_9:
		return 9
	case ddv1.Encoding_ISO_8859_13:
		return 13
	case ddv1.Encoding_ISO_8859_15:
		return 15
	case ddv1.Encoding_ISO_8859_16:
		return 16
	case ddv1.Encoding_KOI8_R:
		return 80
	case ddv1.Encoding_KOI8_U:
		return 85
	case ddv1.Encoding_IA5:
		return 0 // IA5 uses default encoding (ASCII-compatible)
	case ddv1.Encoding_ENCODING_EMPTY, ddv1.Encoding_ENCODING_UNSPECIFIED, ddv1.Encoding_ENCODING_UNRECOGNIZED:
		return 255
	default:
		return 255
	}
}

// encodeWithCodePage encodes a string to bytes using the specified code page.
// For now, this is a simple implementation that only handles ASCII-compatible encodings.
func encodeWithCodePage(codePage byte, s string) ([]byte, error) {
	// For code page 255 (empty), return empty bytes
	if codePage == 255 {
		return []byte{}, nil
	}

	// For ASCII-compatible code pages (0, 1), we can just convert to bytes
	// TODO: Implement proper encoding for other code pages using charmap.Encoder
	if codePage == 0 || codePage == 1 {
		return []byte(s), nil
	}

	// For now, fall back to ISO-8859-1 encoding for other code pages
	// This is a simplification and should be enhanced for full support
	return []byte(s), nil
}
