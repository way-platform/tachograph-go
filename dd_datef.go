package tachograph

import (
	"bytes"

	"google.golang.org/protobuf/types/known/timestamppb"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// appendDatef appends a 4-byte BCD-encoded date.
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func appendDatef(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := ts.AsTime().Year()
	month := int(ts.AsTime().Month())
	day := ts.AsTime().Day()

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
}

// readDatef reads a Datef value (4 bytes BCD) from a bytes.Reader and converts to Date
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func readDatef(r *bytes.Reader) (*ddv1.Date, error) {
	b := make([]byte, 4)
	_, _ = r.Read(b) // ignore error as we're reading from in-memory buffer

	// Parse BCD format: YYYYMMDD
	year := int32(int32((b[0]&0xF0)>>4)*1000 + int32(b[0]&0x0F)*100 + int32((b[1]&0xF0)>>4)*10 + int32(b[1]&0x0F))
	month := int32(int32((b[2]&0xF0)>>4)*10 + int32(b[2]&0x0F))
	day := int32(int32((b[3]&0xF0)>>4)*10 + int32(b[3]&0x0F))

	// Validate the date
	if year < 1900 || year > 9999 || month < 1 || month > 12 || day < 1 || day > 31 {
		return nil, nil // Return nil for invalid or zero dates
	}

	date := &ddv1.Date{}
	date.SetYear(year)
	date.SetMonth(month)
	date.SetDay(day)
	return date, nil
}
