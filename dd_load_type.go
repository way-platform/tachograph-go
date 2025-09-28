package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalLoadType parses load type from raw data.
//
// The data type `LoadType` is specified in the Data Dictionary, Section 2.90a.
//
// ASN.1 Definition:
//
//	LoadType ::= INTEGER {
//	    not-defined(0), passengers(1), goods(2)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Load Type (1 byte): Raw integer value (0-2)
func unmarshalLoadType(data []byte) (ddv1.LoadType, error) {
	if len(data) < 1 {
		return ddv1.LoadType_LOAD_TYPE_UNSPECIFIED, fmt.Errorf("insufficient data for LoadType: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	loadType := ddv1.LoadType_LOAD_TYPE_UNSPECIFIED
	SetLoadType(ddv1.LoadType_LOAD_TYPE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			loadType = ddv1.LoadType(enumNum)
		}, func(unrecognized int32) {
			loadType = ddv1.LoadType_LOAD_TYPE_UNRECOGNIZED
		})

	return loadType, nil
}

// appendLoadType appends load type as a single byte.
//
// The data type `LoadType` is specified in the Data Dictionary, Section 2.90a.
//
// ASN.1 Definition:
//
//	LoadType ::= INTEGER {
//	    not-defined(0), passengers(1), goods(2)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Load Type (1 byte): Raw integer value (0-2)
func appendLoadType(dst []byte, loadType ddv1.LoadType) []byte {
	// Get the protocol value for the enum
	protocolValue := GetLoadTypeProtocolValue(loadType, 0)
	return append(dst, byte(protocolValue))
}
