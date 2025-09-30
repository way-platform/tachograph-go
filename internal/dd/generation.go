package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGeneration parses generation from raw data.
//
// The data type `Generation` is specified in the Data Dictionary.
//
// ASN.1 Definition:
//
//	Generation ::= INTEGER (1..2)
//
// Binary Layout (1 byte):
//   - Generation (1 byte): Raw integer value (1-2)
func UnmarshalGeneration(data []byte) (ddv1.Generation, error) {
	if len(data) != 1 {
		return ddv1.Generation_GENERATION_UNSPECIFIED, fmt.Errorf("invalid data length for Generation: got %d, want 1", len(data))
	}

	protocolValue := int32(data[0])

	// Use reflection to find matching enum value by protocol_enum_value annotation
	enumDesc := ddv1.Generation_GENERATION_1.Descriptor()
	if enumNumber, found := GetEnumForProtocolValue(enumDesc, protocolValue); found {
		return ddv1.Generation(enumNumber), nil
	}

	return ddv1.Generation_GENERATION_UNSPECIFIED, fmt.Errorf("invalid generation value: %d", protocolValue)
}

// AppendGeneration appends generation as a single byte.
//
// The data type `Generation` is specified in the Data Dictionary.
//
// ASN.1 Definition:
//
//	Generation ::= INTEGER (1..2)
//
// Binary Layout (1 byte):
//   - Generation (1 byte): Raw integer value (1-2)
func AppendGeneration(dst []byte, generation ddv1.Generation) []byte {
	if protocolValue, found := GetProtocolValueForEnum(generation); found {
		return append(dst, byte(protocolValue))
	}
	// Default to 0 for unspecified/unrecognized
	return append(dst, 0)
}
