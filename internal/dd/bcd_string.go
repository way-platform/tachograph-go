package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalBcdString parses BCD string data.
//
// The data type `BcdString` is specified in the Data Dictionary, Section 2.7.
//
// ASN.1 Definition:
//
//	BCDString ::= CharacterStringType
//
// Binary Layout (variable length):
//   - BCD String (variable): BCD-encoded bytes
func UnmarshalBcdString(data []byte) (*ddv1.BcdString, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for BcdString: got %d, want at least 1", len(data))
	}
	decoded, err := decodeBCD(data)
	if err != nil {
		return nil, err
	}
	var output ddv1.BcdString
	output.SetEncoded(data)
	output.SetDecoded(int32(decoded))
	return &output, nil
}

// appendBcdString appends BCD string data to dst.
//
// The data type `BcdString` is specified in the Data Dictionary, Section 2.7.
//
// ASN.1 Definition:
//
//	BCDString ::= CharacterStringType
//
// Binary Layout (variable length):
//   - BCD String (variable): BCD-encoded bytes
func AppendBcdString(dst []byte, bcdString *ddv1.BcdString) ([]byte, error) {
	if bcdString == nil {
		// Append empty BCD string (0 bytes)
		return dst, nil
	}
	// Append the raw BCD-encoded bytes
	encoded := bcdString.GetEncoded()
	if encoded != nil {
		dst = append(dst, encoded...)
	}
	return dst, nil
}
