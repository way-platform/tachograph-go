package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalOperationType parses operation type from raw data.
//
// The data type `OperationType` is specified in the Data Dictionary, Section 2.114a.
//
// ASN.1 Definition:
//
//	OperationType ::= INTEGER {
//	    load(1), unload(2), simultaneous(3)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Operation Type (1 byte): Raw integer value (1-3)
func unmarshalOperationType(data []byte) (ddv1.OperationType, error) {
	if len(data) < 1 {
		return ddv1.OperationType_OPERATION_TYPE_UNSPECIFIED, fmt.Errorf("insufficient data for OperationType: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	operationType := ddv1.OperationType_OPERATION_TYPE_UNSPECIFIED
	setEnumFromProtocolValue(ddv1.OperationType_OPERATION_TYPE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			operationType = ddv1.OperationType(enumNum)
		}, func(unrecognized int32) {
			operationType = ddv1.OperationType_OPERATION_TYPE_UNRECOGNIZED
		})

	return operationType, nil
}

// appendOperationType appends operation type as a single byte.
//
// The data type `OperationType` is specified in the Data Dictionary, Section 2.114a.
//
// ASN.1 Definition:
//
//	OperationType ::= INTEGER {
//	    load(1), unload(2), simultaneous(3)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Operation Type (1 byte): Raw integer value (1-3)
func appendOperationType(dst []byte, operationType ddv1.OperationType) []byte {
	// Get the protocol value for the enum
	protocolValue := getProtocolValueFromEnum(operationType, 0)
	return append(dst, byte(protocolValue))
}
