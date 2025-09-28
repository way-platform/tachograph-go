package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalEventFaultType parses event fault type from raw data.
//
// The data type `EventFaultType` is specified in the Data Dictionary, Section 2.70.
//
// ASN.1 Definition:
//
//	EventFaultType ::= INTEGER (0..255)
//
// Binary Layout (1 byte):
//   - Event Fault Type (1 byte): Raw integer value (0-255)
func unmarshalEventFaultType(data []byte) (ddv1.EventFaultType, error) {
	if len(data) < 1 {
		return ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED, fmt.Errorf("insufficient data for EventFaultType: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	eventFaultType := ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED
	SetEnumFromProtocolValue(ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			eventFaultType = ddv1.EventFaultType(enumNum)
		}, func(unrecognized int32) {
			eventFaultType = ddv1.EventFaultType_EVENT_FAULT_TYPE_UNRECOGNIZED
		})

	return eventFaultType, nil
}

// appendEventFaultType appends event fault type as a single byte.
//
// The data type `EventFaultType` is specified in the Data Dictionary, Section 2.70.
//
// ASN.1 Definition:
//
//	EventFaultType ::= INTEGER (0..255)
//
// Binary Layout (1 byte):
//   - Event Fault Type (1 byte): Raw integer value (0-255)
func appendEventFaultType(dst []byte, eventFaultType ddv1.EventFaultType) []byte {
	// Get the protocol value for the enum
	protocolValue := GetEventFaultTypeProtocolValue(eventFaultType, 0)
	return append(dst, byte(protocolValue))
}
