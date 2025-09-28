package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDrivingStatus parses driving status from raw data.
//
// The data type `DrivingStatus` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	DrivingStatus ::= INTEGER (0..1)
//
// Binary Layout (1 bit):
//   - Driving Status (1 bit): Raw integer value (0-1)
func unmarshalDrivingStatus(data []byte) (ddv1.DrivingStatus, error) {
	if len(data) < 1 {
		return ddv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED, fmt.Errorf("insufficient data for DrivingStatus: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	drivingStatus := ddv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED
	SetDrivingStatus(ddv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED.Descriptor(), rawValue, func(en protoreflect.EnumNumber) {
		drivingStatus = ddv1.DrivingStatus(en)
	}, func(unrecognized int32) {
		drivingStatus = ddv1.DrivingStatus_DRIVING_STATUS_UNRECOGNIZED
	})

	return drivingStatus, nil
}

// appendDrivingStatus appends driving status as a single byte.
//
// The data type `DrivingStatus` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	DrivingStatus ::= INTEGER (0..1)
//
// Binary Layout (1 bit):
//   - Driving Status (1 bit): Raw integer value (0-1)
func appendDrivingStatus(dst []byte, drivingStatus ddv1.DrivingStatus) []byte {
	// Get the protocol value for the enum
	protocolValue := GetDrivingStatus(drivingStatus, 0)
	return append(dst, byte(protocolValue))
}
