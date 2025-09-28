package tachograph

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// getProtocolValueFromEnum returns the protocol value for an enum using protobuf reflection.
// This replaces all the hardcoded Get*ProtocolValue functions in enum_helpers.go.
func getProtocolValueFromEnum(enumValue protoreflect.Enum) (int32, bool) {
	enumDesc := enumValue.Descriptor()
	valueDesc := enumDesc.Values().ByNumber(enumValue.Number())
	if valueDesc == nil {
		return 0, false
	}

	opts := valueDesc.Options()
	if !proto.HasExtension(opts, ddv1.E_ProtocolEnumValue) {
		return 0, false
	}

	protocolValue := proto.GetExtension(opts, ddv1.E_ProtocolEnumValue).(int32)
	return protocolValue, true
}

// setEnumFromProtocolValue sets an enum from a protocol value using protobuf reflection.
// This replaces all the hardcoded Set* functions in enum_helpers.go.
func setEnumFromProtocolValue(enumDesc protoreflect.EnumDescriptor, rawValue int32) (protoreflect.EnumNumber, bool) {
	values := enumDesc.Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()

		// Check if this value has the protocol_enum_value annotation
		if proto.HasExtension(opts, ddv1.E_ProtocolEnumValue) {
			protocolValue := proto.GetExtension(opts, ddv1.E_ProtocolEnumValue).(int32)
			if protocolValue == rawValue {
				return valueDesc.Number(), true
			}
		}
	}
	return 0, false
}

// Generic enum conversion helpers that can be used throughout the codebase

// SetEnumFromProtocolValue is a generic helper that converts a raw protocol value to any enum type.
// Usage: SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
func SetEnumFromProtocolValue(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	if enumNumber, found := setEnumFromProtocolValue(enumDesc, rawValue); found {
		setEnum(enumNumber)
	} else {
		// For unknown values, use the first available enum value as default
		values := enumDesc.Values()
		if values.Len() > 0 {
			setEnum(values.Get(0).Number())
		}
		// Only call setUnrecognized if it's not nil (for backwards compatibility)
		if setUnrecognized != nil {
			setUnrecognized(rawValue)
		}
	}
}

// GetProtocolValueFromEnum is a generic helper that converts any enum to its protocol value.
// Usage: GetProtocolValueFromEnum(enumValue, unrecognizedValue)
func GetProtocolValueFromEnum(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	// Check if this is an UNRECOGNIZED value by checking if the number is negative
	if enumValue.Number() < 0 {
		return unrecognizedValue
	}

	if protocolValue, ok := getProtocolValueFromEnum(enumValue); ok {
		return protocolValue
	}

	// Fallback - this shouldn't happen in well-formed data
	return int32(enumValue.Number())
}

// Specific enum conversion functions for commonly used enums
// These provide convenient wrappers around the generic functions

// GetCardSlotNumber returns the protocol value for a CardSlotNumber enum
func GetCardSlotNumber(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	return GetProtocolValueFromEnum(enumValue, unrecognizedValue)
}

// GetDrivingStatus returns the protocol value for a DrivingStatus enum
func GetDrivingStatus(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	return GetProtocolValueFromEnum(enumValue, unrecognizedValue)
}

// GetCardInserted returns the protocol value for a CardInserted enum
func GetCardInserted(enumValue protoreflect.Enum) int32 {
	return GetProtocolValueFromEnum(enumValue, 0)
}

// GetCardInsertedFromBool returns the protocol value for a CardInserted from a boolean
func GetCardInsertedFromBool(inserted bool) int32 {
	if inserted {
		return 0 // Card is inserted
	}
	return 1 // Card is not inserted
}

// GetDriverActivityValue returns the protocol value for a DriverActivity enum
func GetDriverActivityValue(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	return GetProtocolValueFromEnum(enumValue, unrecognizedValue)
}

// GetEventFaultTypeProtocolValue returns the protocol value for an EventFaultType enum
func GetEventFaultTypeProtocolValue(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	return GetProtocolValueFromEnum(enumValue, unrecognizedValue)
}

// GetSpecificConditionTypeProtocolValue returns the protocol value for a SpecificConditionType enum
func GetSpecificConditionTypeProtocolValue(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	return GetProtocolValueFromEnum(enumValue, unrecognizedValue)
}

// Set* functions for setting enum values from protocol values
// These provide convenient wrappers around the generic SetEnumFromProtocolValue function

// SetCardSlotNumber sets a CardSlotNumber enum from a protocol value
func SetCardSlotNumber(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
}

// SetDrivingStatus sets a DrivingStatus enum from a protocol value
func SetDrivingStatus(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
}

// SetCardInserted sets a CardInserted enum from a protocol value
func SetCardInserted(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
}

// SetDriverActivityValue sets a DriverActivity enum from a protocol value
func SetDriverActivityValue(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
}

// SetEquipmentType sets an EquipmentType enum from a protocol value
func SetEquipmentType(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
}

// SetSpecificConditionType sets a SpecificConditionType enum from a protocol value
func SetSpecificConditionType(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	SetEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
}
