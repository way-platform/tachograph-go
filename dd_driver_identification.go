package tachograph

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
func unmarshalDriverIdentification(data []byte) (*ddv1.DriverIdentification, error) {
	const (
		lenDriverIdentification = 14
	)

	if len(data) < lenDriverIdentification {
		return nil, fmt.Errorf("insufficient data for DriverIdentification: got %d, want %d", len(data), lenDriverIdentification)
	}

	driverID := &ddv1.DriverIdentification{}

	// Parse driver identification number (14 bytes)
	identificationNumber, err := unmarshalIA5StringValue(data[0:14])
	if err != nil {
		return nil, fmt.Errorf("failed to parse driver identification number: %w", err)
	}
	driverID.SetIdentificationNumber(identificationNumber)

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
func appendDriverIdentification(dst []byte, driverID *ddv1.DriverIdentification) ([]byte, error) {
	if driverID == nil {
		// Append default values (14 zero bytes)
		return append(dst, make([]byte, 14)...), nil
	}

	// Append driver identification number (14 bytes)
	identificationNumber := driverID.GetIdentificationNumber()
	if identificationNumber != nil {
		return appendString(dst, identificationNumber.GetDecoded(), 14)
	}
	return appendString(dst, "", 14)
}
