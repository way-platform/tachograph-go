package dd

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// getEncodingFromCodePage maps a code page byte to the corresponding Encoding enum.
//
// This function provides the mapping between tachograph protocol code page values
// and the protobuf Encoding enum values. Code pages are used in the binary format
// to specify the character encoding of string data.
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

// getCodePageFromEncoding maps an Encoding enum to a code page byte.
//
// This function provides the reverse mapping from protobuf Encoding enum values
// to tachograph protocol code page bytes. Used when marshalling string data
// to ensure the correct code page is written to the binary format.
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
	case ddv1.Encoding_ENCODING_EMPTY, ddv1.Encoding_ENCODING_UNSPECIFIED, ddv1.Encoding_ENCODING_UNRECOGNIZED:
		return 255
	default:
		return 255
	}
}

// trimSpaceAndZeroBytes trims spaces, control characters, and padding bytes off a byte slice.
//
// This function removes common padding and control characters that are often
// used in tachograph string data to pad fixed-length fields. The trimmed
// characters include:
// - ASCII control characters: 0x00-0x1F (null, SOH, STX, ..., US)
// - Whitespace: space (0x20)
// - Control characters: 0x85 (NEL), 0xA0 (non-breaking space)
// - Padding bytes: 0xFF (often used as padding in binary protocols)
func trimSpaceAndZeroBytes(b []byte) []byte {
	// Define cutset as string - bytes.Trim handles this properly
	cutset := "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f \x85\xA0\xFF"
	return bytes.Trim(b, cutset)
}

// decodeWithCodePage decodes a byte slice with the given code page, returns the trimmed decoded string.
//
// This function handles the conversion from tachograph protocol string data (encoded
// with various character sets) to UTF-8 strings. It supports multiple code pages
// including ISO-8859 variants and KOI8 variants.
//
// Parameters:
//   - codePage: The code page byte from the protocol (0-255)
//   - data: The raw string data bytes to decode
//
// Returns:
//   - The decoded UTF-8 string with padding characters trimmed
//   - An error if decoding fails
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

// encodeWithCodePage encodes a string to bytes using the specified code page.
//
// This function handles the conversion from UTF-8 strings to tachograph protocol
// string data using various character encodings.
//
// Parameters:
//   - codePage: The code page byte to use for encoding (0-255)
//   - s: The UTF-8 string to encode
//
// Returns:
//   - The encoded bytes in the specified character set
//   - An error if encoding fails
func encodeWithCodePage(codePage byte, s string) ([]byte, error) {
	// For code page 255 (empty), return empty bytes
	if codePage == 255 {
		return []byte{}, nil
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

	enc := cmap.NewEncoder()
	result, err := enc.String(s)
	if err != nil {
		return nil, fmt.Errorf("could not encode string to code page %d: %w", codePage, err)
	}

	return []byte(result), nil
}
