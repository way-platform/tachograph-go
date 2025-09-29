package tachograph

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// getProtocolValueFromEnumInternal returns the protocol value for an enum using protobuf reflection.
// This replaces all the hardcoded Get*ProtocolValue functions in enum_helpers.go.
func getProtocolValueFromEnumInternal(enumValue protoreflect.Enum) (int32, bool) {
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

// setEnumFromProtocolValueInternal sets an enum from a protocol value using protobuf reflection.
// This replaces all the hardcoded Set* functions in enum_helpers.go.
func setEnumFromProtocolValueInternal(enumDesc protoreflect.EnumDescriptor, rawValue int32) (protoreflect.EnumNumber, bool) {
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

// setEnumFromProtocolValue is a generic helper that converts a raw protocol value to any enum type.
// Usage: setEnumFromProtocolValue(enumDesc, rawValue, setEnum, setUnrecognized)
func setEnumFromProtocolValue(enumDesc protoreflect.EnumDescriptor, rawValue int32, setEnum func(protoreflect.EnumNumber), setUnrecognized func(int32)) {
	if enumNumber, found := setEnumFromProtocolValueInternal(enumDesc, rawValue); found {
		setEnum(enumNumber)
	} else {
		// For unknown values, use UNRECOGNIZED instead of UNSPECIFIED
		// UNRECOGNIZED is typically the second enum value (index 1)
		values := enumDesc.Values()
		if values.Len() > 1 {
			// Use the UNRECOGNIZED value (typically index 1)
			setEnum(values.Get(1).Number())
		} else if values.Len() > 0 {
			// Fallback to first value if UNRECOGNIZED doesn't exist
			setEnum(values.Get(0).Number())
		}
		// Call setUnrecognized to preserve the raw value for data fidelity
		if setUnrecognized != nil {
			setUnrecognized(rawValue)
		}
	}
}

// getProtocolValueFromEnum is a generic helper that converts any enum to its protocol value.
// Usage: getProtocolValueFromEnum(enumValue, unrecognizedValue)
func getProtocolValueFromEnum(enumValue protoreflect.Enum, unrecognizedValue int32) int32 {
	// Check if this is an UNRECOGNIZED value by checking if the number is negative
	if enumValue.Number() < 0 {
		return unrecognizedValue
	}

	if protocolValue, ok := getProtocolValueFromEnumInternal(enumValue); ok {
		return protocolValue
	}

	// Fallback - this shouldn't happen in well-formed data
	return int32(enumValue.Number())
}

// getCardInsertedFromBool returns the protocol value for a CardInserted from a boolean
func getCardInsertedFromBool(inserted bool) int32 {
	if inserted {
		return 0 // Card is inserted
	}
	return 1 // Card is not inserted
}

// Helper functions for specific enum types that return protoreflect.Enum interfaces
func getCardSlotNumberProtocolValue(slot ddv1.CardSlotNumber, unrecognizedValue int32) int32 {
	return getProtocolValueFromEnum(slot, unrecognizedValue)
}

func getDrivingStatusProtocolValue(status ddv1.DrivingStatus, unrecognizedValue int32) int32 {
	return getProtocolValueFromEnum(status, unrecognizedValue)
}

func getDriverActivityValueProtocolValue(activity ddv1.DriverActivityValue, unrecognizedValue int32) int32 {
	return getProtocolValueFromEnum(activity, unrecognizedValue)
}

func getEventFaultTypeProtocolValue(eventType ddv1.EventFaultType, unrecognizedValue int32) int32 {
	return getProtocolValueFromEnum(eventType, unrecognizedValue)
}
