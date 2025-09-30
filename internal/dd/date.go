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
func (opts UnmarshalOptions) UnmarshalDate(input []byte) (*ddv1.Date, error) {
	const lenDatef = 4
	if len(input) != lenDatef {
		return nil, fmt.Errorf("invalid data length for Date: got %d, want %d", len(input), lenDatef)
	}
	var output ddv1.Date
	output.SetRawData(input[:lenDatef])
	year := int32(int32((input[0]&0xF0)>>4)*1000 + int32(input[0]&0x0F)*100 + int32((input[1]&0xF0)>>4)*10 + int32(input[1]&0x0F))
	month := int32(int32((input[2]&0xF0)>>4)*10 + int32(input[2]&0x0F))
	day := int32(int32((input[3]&0xF0)>>4)*10 + int32(input[3]&0x0F))
	output.SetYear(year)
	output.SetMonth(month)
	output.SetDay(day)
	return &output, nil
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
	var canvas [lenDatef]byte
	if date.HasRawData() {
		if len(date.GetRawData()) != lenDatef {
			return nil, fmt.Errorf(
				"invalid raw_data length for Date: got %d, want %d",
				len(date.GetRawData()), lenDatef,
			)
		}
		copy(canvas[:], date.GetRawData())
	}
	year := int(date.GetYear())
	month := int(date.GetMonth())
	day := int(date.GetDay())
	canvas[0] = byte((year/1000)%10<<4 | (year/100)%10)
	canvas[1] = byte((year/10)%10<<4 | year%10)
	canvas[2] = byte((month/10)%10<<4 | month%10)
	canvas[3] = byte((day/10)%10<<4 | day%10)
	return append(dst, canvas[:]...), nil
}
