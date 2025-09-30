package dd

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AppendDatef appends a 4-byte BCD-encoded date from a Timestamp.
//
// DEPRECATED: This function exists for backwards compatibility with protobuf
// schemas that incorrectly use Timestamp for date-only fields (e.g., card_expiry_date
// in card identification). New code should use AppendDate with dd.v1.Date instead.
//
// The data type `Datef` is specified in the Data Dictionary, Section 2.57.
//
// ASN.1 Definition:
//
//	Datef ::= OCTET STRING (SIZE(4))
func AppendDatef(dst []byte, ts *timestamppb.Timestamp) []byte {
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
