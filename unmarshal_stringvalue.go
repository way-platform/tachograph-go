package tachograph

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalStringValue unmarshals a code-paged string value from binary data.
// The input should contain a code page byte followed by the encoded string data.
func unmarshalStringValue(input []byte) (*datadictionaryv1.StringValue, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("insufficient data for string value: %w", io.ErrUnexpectedEOF)
	}

	codePage := input[0]
	data := input[1:]

	var output datadictionaryv1.StringValue
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

// unmarshalIA5StringValue unmarshals an IA5 (ASCII) string value from binary data.
// IA5 strings have a fixed encoding and may include padding.
func unmarshalIA5StringValue(input []byte) (*datadictionaryv1.StringValue, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("insufficient data for string value: %w", io.ErrUnexpectedEOF)
	}

	var output datadictionaryv1.StringValue
	output.SetEncoding(datadictionaryv1.Encoding_IA5)
	output.SetEncoded(input)

	// Decode and trim the input bytes
	decoded := trimSpaceAndZero(string(input))
	output.SetDecoded(decoded)

	return &output, nil
}

// getEncodingFromCodePage maps a code page byte to the corresponding Encoding enum.
func getEncodingFromCodePage(codePage byte) datadictionaryv1.Encoding {
	switch codePage {
	case 0:
		return datadictionaryv1.Encoding_ENCODING_DEFAULT
	case 1:
		return datadictionaryv1.Encoding_ISO_8859_1
	case 2:
		return datadictionaryv1.Encoding_ISO_8859_2
	case 3:
		return datadictionaryv1.Encoding_ISO_8859_3
	case 5:
		return datadictionaryv1.Encoding_ISO_8859_5
	case 7:
		return datadictionaryv1.Encoding_ISO_8859_7
	case 9:
		return datadictionaryv1.Encoding_ISO_8859_9
	case 13:
		return datadictionaryv1.Encoding_ISO_8859_13
	case 15:
		return datadictionaryv1.Encoding_ISO_8859_15
	case 16:
		return datadictionaryv1.Encoding_ISO_8859_16
	case 80:
		return datadictionaryv1.Encoding_KOI8_R
	case 85:
		return datadictionaryv1.Encoding_KOI8_U
	case 255:
		return datadictionaryv1.Encoding_ENCODING_EMPTY
	default:
		return datadictionaryv1.Encoding_ENCODING_UNRECOGNIZED
	}
}

// trimSpaceAndZero trims spaces, 0x00 and 0xff values off a string
func trimSpaceAndZero(s string) string {
	w := "\t\n\v\f\r \x85\xA0\x00\xFF"
	return strings.Trim(s, w)
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
	return trimSpaceAndZero(res), nil
}

// Protocol-compliant helper functions

// unmarshalStringValueFromReader reads a code-paged string from a reader
func unmarshalStringValueFromReader(r io.Reader, length int) (*datadictionaryv1.StringValue, error) {
	data := make([]byte, length)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}
	return unmarshalStringValue(data)
}

// unmarshalIA5StringValueFromReader reads an IA5 string from a reader
func unmarshalIA5StringValueFromReader(r io.Reader, length int) (*datadictionaryv1.StringValue, error) {
	data := make([]byte, length)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}
	return unmarshalIA5StringValue(data)
}

// unmarshalNationNumericFromReader reads a nation code from a reader
func unmarshalNationNumericFromReader(r io.Reader) (datadictionaryv1.NationNumeric, error) {
	var nation byte
	if err := binary.Read(r, binary.BigEndian, &nation); err != nil {
		return datadictionaryv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED, err
	}
	return datadictionaryv1.NationNumeric(nation), nil
}

// unmarshalControlTypeFromReader reads a control type byte from a reader
func unmarshalControlTypeFromReader(r io.Reader) (*datadictionaryv1.ControlType, error) {
	var b byte
	if err := binary.Read(r, binary.BigEndian, &b); err != nil {
		return nil, err
	}

	ct := &datadictionaryv1.ControlType{}
	ct.SetCardDownloading((b & 0x80) != 0)
	ct.SetVuDownloading((b & 0x40) != 0)
	ct.SetPrinting((b & 0x20) != 0)
	ct.SetDisplay((b & 0x10) != 0)
	ct.SetCalibrationChecking((b & 0x08) != 0)
	return ct, nil
}

// unmarshalDateFromReader reads a BCD-encoded date from a reader
func unmarshalDateFromReader(r io.Reader) (*datadictionaryv1.Date, error) {
	data := make([]byte, 4)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}

	// Parse BCD format: YYYYMMDD
	year := int32(int32((data[0]&0xF0)>>4)*1000 + int32(data[0]&0x0F)*100 + int32((data[1]&0xF0)>>4)*10 + int32(data[1]&0x0F))
	month := int32(int32((data[2]&0xF0)>>4)*10 + int32(data[2]&0x0F))
	day := int32(int32((data[3]&0xF0)>>4)*10 + int32(data[3]&0x0F))

	date := &datadictionaryv1.Date{}
	date.SetYear(year)
	date.SetMonth(month)
	date.SetDay(day)
	return date, nil
}
