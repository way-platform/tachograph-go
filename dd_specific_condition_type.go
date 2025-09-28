package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalSpecificConditionType parses specific condition type from raw data.
//
// The data type `SpecificConditionType` is specified in the Data Dictionary, Section 2.154.
//
// ASN.1 Definition:
//
//	SpecificConditionType ::= INTEGER (0..3)
//
// Binary Layout (1 byte):
//   - Specific Condition Type (1 byte): Raw integer value (0-3)
func unmarshalSpecificConditionType(data []byte) (ddv1.SpecificConditionType, error) {
	if len(data) < 1 {
		return ddv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED, fmt.Errorf("insufficient data for SpecificConditionType: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	specificConditionType := ddv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED
	setEnumFromProtocolValue(ddv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			specificConditionType = ddv1.SpecificConditionType(enumNum)
		}, func(unrecognized int32) {
			specificConditionType = ddv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNRECOGNIZED
		})

	return specificConditionType, nil
}

// appendSpecificConditionType appends specific condition type as a single byte.
//
// The data type `SpecificConditionType` is specified in the Data Dictionary, Section 2.154.
//
// ASN.1 Definition:
//
//	SpecificConditionType ::= INTEGER (0..3)
//
// Binary Layout (1 byte):
//   - Specific Condition Type (1 byte): Raw integer value (0-3)
func appendSpecificConditionType(dst []byte, specificConditionType ddv1.SpecificConditionType) []byte {
	// Get the protocol value for the enum
	protocolValue := getProtocolValueFromEnum(specificConditionType, 0)
	return append(dst, byte(protocolValue))
}
