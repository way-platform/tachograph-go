package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDriverActivityValue parses driver activity value from raw data.
//
// The data type `DriverActivityValue` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	DriverActivityValue ::= INTEGER (0..3)
//
// Binary Layout (2 bits):
//   - Activity Value (2 bits): Raw integer value (0-3)
func unmarshalDriverActivityValue(data []byte) (ddv1.DriverActivityValue, error) {
	if len(data) < 1 {
		return ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED, fmt.Errorf("insufficient data for DriverActivityValue: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	activityValue := ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED
	SetDriverActivityValue(ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED.Descriptor(), rawValue, func(en protoreflect.EnumNumber) {
		activityValue = ddv1.DriverActivityValue(en)
	}, func(unrecognized int32) {
		activityValue = ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNRECOGNIZED
	})

	return activityValue, nil
}

// appendDriverActivityValue appends driver activity value as a single byte.
//
// The data type `DriverActivityValue` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	DriverActivityValue ::= INTEGER (0..3)
//
// Binary Layout (2 bits):
//   - Activity Value (2 bits): Raw integer value (0-3)
func appendDriverActivityValue(dst []byte, activityValue ddv1.DriverActivityValue) []byte {
	// Get the protocol value for the enum
	protocolValue := GetDriverActivityValue(activityValue, 0)
	return append(dst, byte(protocolValue))
}
