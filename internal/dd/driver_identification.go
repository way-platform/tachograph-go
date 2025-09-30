package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDriverIdentification parses driver identification data.
//
// The data type `DriverIdentification` is specified in the Data Dictionary, Section 2.26.
//
// ASN.1 Definition:
//
//	driverIdentification SEQUENCE {
//	    driverIdentificationNumber IA5String(SIZE(14))
//	}
//
// Binary Layout (14 bytes):
//   - Driver Identification Number (14 bytes): IA5String
func UnmarshalDriverIdentification(data []byte) (*ddv1.DriverIdentification, error) {
	const (
		lenDriverIdentification = 14
	)

	if len(data) != lenDriverIdentification {
		return nil, fmt.Errorf("invalid data length for DriverIdentification: got %d, want %d", len(data), lenDriverIdentification)
	}

	driverID := &ddv1.DriverIdentification{}

	// Parse driver identification number (14 bytes)
	identificationNumber, err := UnmarshalIA5StringValue(data[0:14])
	if err != nil {
		return nil, fmt.Errorf("failed to parse driver identification number: %w", err)
	}
	driverID.SetDriverIdentificationNumber(identificationNumber)

	return driverID, nil
}

// appendDriverIdentification appends driver identification data to dst.
//
// The data type `DriverIdentification` is specified in the Data Dictionary, Section 2.26.
//
// ASN.1 Definition:
//
//	driverIdentification SEQUENCE {
//	    driverIdentificationNumber IA5String(SIZE(14))
//	}
//
// Binary Layout (14 bytes):
//   - Driver Identification Number (14 bytes): IA5String
func AppendDriverIdentification(dst []byte, driverID *ddv1.DriverIdentification) ([]byte, error) {
	if driverID == nil {
		return nil, fmt.Errorf("driverID cannot be nil")
	}

	// Append driver identification number (14 bytes)
	identificationNumber := driverID.GetDriverIdentificationNumber()
	if identificationNumber != nil {
		return AppendString(dst, identificationNumber.GetValue(), 14)
	}
	return AppendString(dst, "", 14)
}
