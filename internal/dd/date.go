package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalDate unmarshals a BCD-encoded date from a byte slice.
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
//
// Binary Layout (4 bytes):
//   - Year (2 bytes): BCD-encoded YYYY
//   - Month (1 byte): BCD-encoded MM
//   - Day (1 byte): BCD-encoded DD
func UnmarshalDate(data []byte) (*ddv1.Date, error) {
	const lenDatef = 4

	if len(data) < lenDatef {
		return nil, fmt.Errorf("insufficient data for Date: got %d, want %d", len(data), lenDatef)
	}

	date := &ddv1.Date{}

	// Store the original encoded bytes for round-trip fidelity
	date.SetEncoded(data[:lenDatef])

	// Parse BCD format: YYYYMMDD
	year := int32(int32((data[0]&0xF0)>>4)*1000 + int32(data[0]&0x0F)*100 + int32((data[1]&0xF0)>>4)*10 + int32(data[1]&0x0F))
	month := int32(int32((data[2]&0xF0)>>4)*10 + int32(data[2]&0x0F))
	day := int32(int32((data[3]&0xF0)>>4)*10 + int32(data[3]&0x0F))

	// Validate the date
	if year < 1900 || year > 9999 {
		return nil, fmt.Errorf("invalid year in Date: %d", year)
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month in Date: %d", month)
	}
	if day < 1 || day > 31 {
		return nil, fmt.Errorf("invalid day in Date: %d", day)
	}

	date.SetYear(year)
	date.SetMonth(month)
	date.SetDay(day)

	return date, nil
}

// AppendDate appends a 4-byte BCD-encoded date from the Date type.
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
//
// Binary Layout (4 bytes):
//   - Year (2 bytes): BCD-encoded YYYY
//   - Month (1 byte): BCD-encoded MM
//   - Day (1 byte): BCD-encoded DD
func AppendDate(dst []byte, date *ddv1.Date) []byte {
	const lenDatef = 4

	if date == nil {
		return append(dst, 0, 0, 0, 0)
	}

	// Prefer the original encoded bytes for perfect round-trip fidelity
	if encoded := date.GetEncoded(); len(encoded) >= lenDatef {
		return append(dst, encoded[:lenDatef]...)
	}

	// Fall back to encoding from decoded values
	year := int(date.GetYear())
	month := int(date.GetMonth())
	day := int(date.GetDay())

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
}
