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

	if len(data) != lenDatef {
		return nil, fmt.Errorf("invalid data length for Date: got %d, want %d", len(data), lenDatef)
	}
	date := &ddv1.Date{}
	// Store the original encoded bytes for round-trip fidelity
	date.SetRawData(data[:lenDatef])
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
func AppendDate(dst []byte, date *ddv1.Date) ([]byte, error) {
	const lenDatef = 4

	// Use stack-allocated array for the canvas (fixed size, avoids heap allocation)
	var canvas [lenDatef]byte

	// Start with raw_data as canvas if available (raw data painting approach)
	if rawData := date.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenDatef {
			return nil, fmt.Errorf("invalid raw_data length for Date: got %d, want %d", len(rawData), lenDatef)
		}
		copy(canvas[:], rawData)
	}
	// Otherwise canvas is zero-initialized (Go default)

	// Paint semantic values over the canvas
	year := int(date.GetYear())
	month := int(date.GetMonth())
	day := int(date.GetDay())
	canvas[0] = byte((year/1000)%10<<4 | (year/100)%10)
	canvas[1] = byte((year/10)%10<<4 | year%10)
	canvas[2] = byte((month/10)%10<<4 | month%10)
	canvas[3] = byte((day/10)%10<<4 | day%10)

	return append(dst, canvas[:]...), nil
}
