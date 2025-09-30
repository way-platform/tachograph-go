package dd

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalEnum converts a raw protocol byte value to a typed enum.
// Returns an error if no enum value has a matching protocol_enum_value annotation.
//
// The type parameter T must be a protobuf enum type (underlying type is int32).
//
// The caller decides how to handle unrecognized values - either save to an
// unrecognized_ field or return an error.
func UnmarshalEnum[T interface {
	~int32
	protoreflect.Enum
}](rawValue byte) (T, error) {
	var zero T
	enumDesc := zero.Descriptor()
	values := enumDesc.Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, ddv1.E_ProtocolEnumValue) {
			protocolValue := proto.GetExtension(opts, ddv1.E_ProtocolEnumValue).(int32)
			if protocolValue == int32(rawValue) {
				return T(valueDesc.Number()), nil
			}
		}
	}
	return zero, fmt.Errorf(
		"no enum value in %s has protocol_enum_value=%d",
		enumDesc.FullName(), rawValue,
	)
}

// MarshalEnum converts a typed enum to a raw protocol byte value.
// Returns an error if the enum value doesn't have a protocol_enum_value annotation.
//
// The type parameter T must be a protobuf enum type (underlying type is int32).
//
// This function should only be called after explicitly handling UNRECOGNIZED values.
// If this returns an error, it indicates an invalid enum value (e.g., UNSPECIFIED).
func MarshalEnum[T interface {
	~int32
	protoreflect.Enum
}](value T) (byte, error) {
	enumDesc := value.Descriptor()
	valueDesc := enumDesc.Values().ByNumber(value.Number())
	if valueDesc != nil {
		opts := valueDesc.Options()
		if proto.HasExtension(opts, ddv1.E_ProtocolEnumValue) {
			protocolValue := proto.GetExtension(opts, ddv1.E_ProtocolEnumValue).(int32)
			return byte(protocolValue), nil
		}
	}
	return 0, fmt.Errorf(
		"enum %s value %d has no protocol_enum_value annotation",
		enumDesc.FullName(), value.Number(),
	)
}
