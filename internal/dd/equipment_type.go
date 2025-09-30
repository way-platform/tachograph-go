package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalEquipmentType parses equipment type from a single byte.
//
// The data type `EquipmentType` is specified in the Data Dictionary, Section 2.67.
//
// ASN.1 Definition:
//
//	EquipmentType ::= INTEGER (0..255)
//
// Binary Layout (1 byte):
//   - Equipment Type (1 byte): Raw integer value
//
//nolint:unused
func UnmarshalEquipmentType(data []byte) (ddv1.EquipmentType, error) {
	if len(data) < 1 {
		return ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED, fmt.Errorf("insufficient data for EquipmentType: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	equipmentType := ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED
	if enumNumber, found := SetEnumFromProtocolValue(ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED.Descriptor(), rawValue); found {
		equipmentType = ddv1.EquipmentType(enumNumber)
	} else {
		equipmentType = ddv1.EquipmentType_EQUIPMENT_TYPE_UNRECOGNIZED
	}

	return equipmentType, nil
}

// appendEquipmentType appends equipment type as a single byte.
//
// The data type `EquipmentType` is specified in the Data Dictionary, Section 2.67.
//
// ASN.1 Definition:
//
//	EquipmentType ::= INTEGER (0..255)
//
// Binary Layout (1 byte):
//   - Equipment Type (1 byte): Raw integer value
//
//nolint:unused
func AppendEquipmentType(dst []byte, equipmentType ddv1.EquipmentType) []byte {
	// Get the protocol value for the enum
	if protocolValue, ok := GetProtocolValueFromEnum(equipmentType); ok {
		return append(dst, byte(protocolValue))
	}
	return append(dst, 0) // Default to 0 if not found
}
