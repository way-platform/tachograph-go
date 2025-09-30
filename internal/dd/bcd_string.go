package dd

import (
	"encoding/hex"
	"fmt"
	"strconv"

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
//
//nolint:unused
func UnmarshalBcdString(data []byte) (*ddv1.BcdString, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for BcdString: got %d, want at least 1", len(data))
	}

	// Convert BCD-encoded bytes to BcdString
	return CreateBcdString(data)
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
//
//nolint:unused
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

// bcdBytesToInt converts BCD-encoded bytes to an integer
func BcdBytesToInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	s := hex.EncodeToString(b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid BCD value: %s", s)
	}
	return int(i), nil
}

// createBcdString creates a BcdString message from BCD-encoded bytes
func CreateBcdString(bcdBytes []byte) (*ddv1.BcdString, error) {
	decoded, err := BcdBytesToInt(bcdBytes)
	if err != nil {
		return nil, err
	}

	bcdString := &ddv1.BcdString{}
	bcdString.SetEncoded(bcdBytes)
	bcdString.SetDecoded(int32(decoded))
	return bcdString, nil
}
