package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalEventFaultRecordPurpose parses event fault record purpose from raw data.
//
// The data type `EventFaultRecordPurpose` is specified in the Data Dictionary, Section 2.69.
//
// ASN.1 Definition:
//
//	EventFaultRecordPurpose ::= INTEGER {
//	    tenMostRecent(0), longestInLast10Days(1), fiveLongestInLast365Days(2),
//	    lastInLast10Days(3), mostSeriousInLast10Days(4), fiveMostSeriousInLast365Days(5),
//	    firstAfterLastCalibration(6), activeOrOngoing(7)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Event Fault Record Purpose (1 byte): Raw integer value (0-7)
func unmarshalEventFaultRecordPurpose(data []byte) (ddv1.EventFaultRecordPurpose, error) {
	if len(data) < 1 {
		return ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED, fmt.Errorf("insufficient data for EventFaultRecordPurpose: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	eventFaultRecordPurpose := ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED
	setEnumFromProtocolValue(ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			eventFaultRecordPurpose = ddv1.EventFaultRecordPurpose(enumNum)
		}, func(unrecognized int32) {
			eventFaultRecordPurpose = ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNRECOGNIZED
		})

	return eventFaultRecordPurpose, nil
}

// appendEventFaultRecordPurpose appends event fault record purpose as a single byte.
//
// The data type `EventFaultRecordPurpose` is specified in the Data Dictionary, Section 2.69.
//
// ASN.1 Definition:
//
//	EventFaultRecordPurpose ::= INTEGER {
//	    tenMostRecent(0), longestInLast10Days(1), fiveLongestInLast365Days(2),
//	    lastInLast10Days(3), mostSeriousInLast10Days(4), fiveMostSeriousInLast365Days(5),
//	    firstAfterLastCalibration(6), activeOrOngoing(7)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Event Fault Record Purpose (1 byte): Raw integer value (0-7)
func appendEventFaultRecordPurpose(dst []byte, eventFaultRecordPurpose ddv1.EventFaultRecordPurpose) []byte {
	// Get the protocol value for the enum
	protocolValue := getProtocolValueFromEnum(eventFaultRecordPurpose, 0)
	return append(dst, byte(protocolValue))
}
