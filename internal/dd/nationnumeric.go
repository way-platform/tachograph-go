package dd

import (
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalNationNumeric unmarshals a nation code from a byte slice
//
// The data type `NationNumeric` is specified in the Data Dictionary, Section 2.118.
//
// ASN.1 Definition:
//
//     NationNumeric ::= OCTET STRING (SIZE(1))
func UnmarshalNationNumeric(data []byte) (ddv1.NationNumeric, error) {
	if len(data) == 0 {
		return ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED, nil
	}
	return ddv1.NationNumeric(data[0]), nil
}
