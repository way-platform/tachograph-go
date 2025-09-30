package dd

import (
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalDate unmarshals a BCD-encoded date from a byte slice
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func UnmarshalDate(data []byte) (*ddv1.Date, error) {
	// TODO: Move implementation from unmarshalDate in top-level package
	panic("not implemented yet - will be moved during migration")
}

// AppendDate appends a 4-byte BCD-encoded date from the new Date type.
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func AppendDate(dst []byte, date *ddv1.Date) []byte {
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




