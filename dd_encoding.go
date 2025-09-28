package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalEncoding parses encoding from raw data.
//
// The data type `Encoding` is specified in the Data Dictionary, used within StringValue.
//
// ASN.1 Definition:
//
//	Encoding ::= INTEGER (0..255)
//	-- Code page values for character encodings
//
// Binary Layout (1 byte):
//   - Encoding (1 byte): Raw integer value (0-255)
func unmarshalEncoding(data []byte) (ddv1.Encoding, error) {
	if len(data) < 1 {
		return ddv1.Encoding_ENCODING_UNSPECIFIED, fmt.Errorf("insufficient data for Encoding: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	encoding := ddv1.Encoding_ENCODING_UNSPECIFIED
	SetEncoding(ddv1.Encoding_ENCODING_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			encoding = ddv1.Encoding(enumNum)
		}, func(unrecognized int32) {
			encoding = ddv1.Encoding_ENCODING_UNRECOGNIZED
		})

	return encoding, nil
}

// appendEncoding appends encoding as a single byte.
//
// The data type `Encoding` is specified in the Data Dictionary, used within StringValue.
//
// ASN.1 Definition:
//
//	Encoding ::= INTEGER (0..255)
//	-- Code page values for character encodings
//
// Binary Layout (1 byte):
//   - Encoding (1 byte): Raw integer value (0-255)
func appendEncoding(dst []byte, encoding ddv1.Encoding) []byte {
	// Get the protocol value for the enum
	protocolValue := GetEncodingProtocolValue(encoding, 0)
	return append(dst, byte(protocolValue))
}
