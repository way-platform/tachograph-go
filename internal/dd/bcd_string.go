package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalBcdString parses BCD string data.
//
// The data type `BcdString` is specified in the Data Dictionary, Section 2.7.
//
// ASN.1 Definition:
//
//	BCDString ::= CharacterStringType
//
// Binary Layout (variable length):
//   - BCD String (variable): BCD-encoded bytes
func UnmarshalBcdString(input []byte) (*ddv1.BcdString, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("insufficient data for BcdString: got %d, want at least 1", len(input))
	}
	value, err := decodeBCD(input)
	if err != nil {
		return nil, err
	}
	var output ddv1.BcdString
	output.SetValue(int32(value))
	output.SetLength(int32(len(input)))
	return &output, nil
}

// AppendBcdString appends BCD string data to dst.
//
// The data type `BcdString` is specified in the Data Dictionary, Section 2.7.
//
// ASN.1 Definition:
//
//	BCDString ::= CharacterStringType
//
// Binary Layout (variable length):
//   - BCD String (variable): BCD-encoded bytes
//
// The regulation specifies canonical encoding: digits 0-9 only, with zero
// padding. This makes encoding deterministic from value + length.
func AppendBcdString(dst []byte, bcdString *ddv1.BcdString) ([]byte, error) {
	if bcdString == nil {
		return dst, nil
	}
	value := bcdString.GetValue()
	if value < 0 {
		return nil, fmt.Errorf("cannot encode negative BCD value: %d", value)
	}
	length := bcdString.GetLength()
	if length <= 0 {
		return nil, fmt.Errorf("invalid BCD string length: %d", length)
	}
	return appendBCD(dst, int(value), int(length))
}
