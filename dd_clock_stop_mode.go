package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalClockStopMode parses clock stop mode from raw data.
//
// The data type `ClockStopMode` is specified in the Data Dictionary, Section 2.23.
//
// ASN.1 Definition:
//
//	clockStop OCTET STRING (SIZE(1))
//	-- Bitmask defining clock stop behavior
//
// Binary Layout (3 bits):
//   - Clock Stop Mode (3 bits): Raw integer value (0-5)
func unmarshalClockStopMode(data []byte) (ddv1.ClockStopMode, error) {
	if len(data) < 1 {
		return ddv1.ClockStopMode_CLOCK_STOP_MODE_UNSPECIFIED, fmt.Errorf("insufficient data for ClockStopMode: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	clockStopMode := ddv1.ClockStopMode_CLOCK_STOP_MODE_UNSPECIFIED
	SetEnumFromProtocolValue(ddv1.ClockStopMode_CLOCK_STOP_MODE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			clockStopMode = ddv1.ClockStopMode(enumNum)
		}, func(unrecognized int32) {
			clockStopMode = ddv1.ClockStopMode_CLOCK_STOP_MODE_UNRECOGNIZED
		})

	return clockStopMode, nil
}

// appendClockStopMode appends clock stop mode as a single byte.
//
// The data type `ClockStopMode` is specified in the Data Dictionary, Section 2.23.
//
// ASN.1 Definition:
//
//	clockStop OCTET STRING (SIZE(1))
//	-- Bitmask defining clock stop behavior
//
// Binary Layout (3 bits):
//   - Clock Stop Mode (3 bits): Raw integer value (0-5)
func appendClockStopMode(dst []byte, clockStopMode ddv1.ClockStopMode) []byte {
	// Get the protocol value for the enum
	protocolValue := GetClockStopModeProtocolValue(clockStopMode, 0)
	return append(dst, byte(protocolValue))
}
