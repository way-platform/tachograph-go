package tachograph

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDate unmarshals a BCD-encoded date from a byte slice
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func unmarshalDate(data []byte) (*ddv1.Date, error) {
	const (
		lenDatef = 4
	)

	if len(data) < lenDatef {
		return nil, fmt.Errorf("insufficient data for date: got %d, want %d", len(data), lenDatef)
	}

	// Extract BCD-encoded date components
	year := bcdToInt(data[0])*100 + bcdToInt(data[1])
	month := bcdToInt(data[2])
	day := bcdToInt(data[3])

	date := &ddv1.Date{}
	date.SetYear(int32(year))
	date.SetMonth(int32(month))
	date.SetDay(int32(day))

	return date, nil
}

// appendDate appends a 4-byte BCD-encoded date from the new Date type.
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func appendDate(dst []byte, date *ddv1.Date) []byte {
	if date == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := int(date.GetYear())
	month := int(date.GetMonth())
	day := int(date.GetDay())

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
}

// bcdToInt converts a BCD-encoded byte to an integer
func bcdToInt(b byte) int {
	return int((b>>4)*10 + (b & 0x0F))
}
