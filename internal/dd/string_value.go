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
	output.SetEncoded(data)

	// Decode the string based on the code page
	decoded, err := decodeWithCodePage(codePage, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode string with code page %d: %w", codePage, err)
	}
	output.SetDecoded(decoded)

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
	output.SetEncoded(input)

	// Decode and trim the input bytes
	trimmed := trimSpaceAndZeroBytes(input)
	decoded := string(trimmed)

	// Ensure the result is valid UTF-8
	if !utf8.ValidString(decoded) {
		// Convert invalid UTF-8 sequences to replacement characters
		decoded = strings.ToValidUTF8(decoded, string(utf8.RuneError))
	}

	output.SetDecoded(decoded)

	return &output, nil
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

// CreateStringValue creates a StringValue message from a string
func CreateStringValue(s string) *ddv1.StringValue {
	stringValue := &ddv1.StringValue{}
	stringValue.SetDecoded(s)
	return stringValue
}



