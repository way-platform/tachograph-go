package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalPositionAuthenticationStatus parses position authentication status from raw data.
//
// The data type `PositionAuthenticationStatus` is specified in the Data Dictionary, Section 2.117a.
//
// ASN.1 Definition:
//
//	PositionAuthenticationStatus ::= INTEGER {
//	    notAvailable(0), authenticated(1), notAuthenticated(2),
//	    authenticationCorrupted(3)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Position Authentication Status (1 byte): Raw integer value (0-3)
func unmarshalPositionAuthenticationStatus(data []byte) (ddv1.PositionAuthenticationStatus, error) {
	if len(data) < 1 {
		return ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNSPECIFIED, fmt.Errorf("insufficient data for PositionAuthenticationStatus: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	positionAuthStatus := ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNSPECIFIED
	SetPositionAuthenticationStatus(ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			positionAuthStatus = ddv1.PositionAuthenticationStatus(enumNum)
		}, func(unrecognized int32) {
			positionAuthStatus = ddv1.PositionAuthenticationStatus_POSITION_AUTHENTICATION_STATUS_UNRECOGNIZED
		})

	return positionAuthStatus, nil
}

// appendPositionAuthenticationStatus appends position authentication status as a single byte.
//
// The data type `PositionAuthenticationStatus` is specified in the Data Dictionary, Section 2.117a.
//
// ASN.1 Definition:
//
//	PositionAuthenticationStatus ::= INTEGER {
//	    notAvailable(0), authenticated(1), notAuthenticated(2),
//	    authenticationCorrupted(3)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Position Authentication Status (1 byte): Raw integer value (0-3)
func appendPositionAuthenticationStatus(dst []byte, positionAuthStatus ddv1.PositionAuthenticationStatus) []byte {
	// Get the protocol value for the enum
	protocolValue := GetPositionAuthenticationStatusProtocolValue(positionAuthStatus, 0)
	return append(dst, byte(protocolValue))
}
