package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalCalibrationPurpose parses calibration purpose from raw data.
//
// The data type `CalibrationPurpose` is specified in the Data Dictionary, Section 2.8.
//
// ASN.1 Definition:
//
//	CalibrationPurpose ::= INTEGER {
//	    reserved(0), activation(1), firstInstallation(2), installation(3),
//	    periodicInspection(4), vrnEntryByCompany(5), timeAdjustment(6)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Calibration Purpose (1 byte): Raw integer value (0-6)
func unmarshalCalibrationPurpose(data []byte) (ddv1.CalibrationPurpose, error) {
	if len(data) < 1 {
		return ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED, fmt.Errorf("insufficient data for CalibrationPurpose: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	calibrationPurpose := ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED
	SetCalibrationPurpose(ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			calibrationPurpose = ddv1.CalibrationPurpose(enumNum)
		}, func(unrecognized int32) {
			calibrationPurpose = ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNRECOGNIZED
		})

	return calibrationPurpose, nil
}

// appendCalibrationPurpose appends calibration purpose as a single byte.
//
// The data type `CalibrationPurpose` is specified in the Data Dictionary, Section 2.8.
//
// ASN.1 Definition:
//
//	CalibrationPurpose ::= INTEGER {
//	    reserved(0), activation(1), firstInstallation(2), installation(3),
//	    periodicInspection(4), vrnEntryByCompany(5), timeAdjustment(6)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Calibration Purpose (1 byte): Raw integer value (0-6)
func appendCalibrationPurpose(dst []byte, calibrationPurpose ddv1.CalibrationPurpose) []byte {
	// Get the protocol value for the enum
	protocolValue := GetCalibrationPurposeProtocolValue(calibrationPurpose, 0)
	return append(dst, byte(protocolValue))
}
