package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalCompanyActivityType parses company activity type from raw data.
//
// The data type `CompanyActivityType` is specified in the Data Dictionary, Section 2.47.
//
// ASN.1 Definition:
//
//	CompanyActivityType ::= INTEGER (1..4)
//
// Binary Layout (1 byte):
//   - Company Activity Type (1 byte): Raw integer value (1-4)
func unmarshalCompanyActivityType(data []byte) (ddv1.CompanyActivityType, error) {
	if len(data) < 1 {
		return ddv1.CompanyActivityType_COMPANY_ACTIVITY_TYPE_UNSPECIFIED, fmt.Errorf("insufficient data for CompanyActivityType: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	companyActivityType := ddv1.CompanyActivityType_COMPANY_ACTIVITY_TYPE_UNSPECIFIED
	SetCompanyActivityType(ddv1.CompanyActivityType_COMPANY_ACTIVITY_TYPE_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			companyActivityType = ddv1.CompanyActivityType(enumNum)
		}, func(unrecognized int32) {
			companyActivityType = ddv1.CompanyActivityType_COMPANY_ACTIVITY_TYPE_UNRECOGNIZED
		})

	return companyActivityType, nil
}

// appendCompanyActivityType appends company activity type as a single byte.
//
// The data type `CompanyActivityType` is specified in the Data Dictionary, Section 2.47.
//
// ASN.1 Definition:
//
//	CompanyActivityType ::= INTEGER (1..4)
//
// Binary Layout (1 byte):
//   - Company Activity Type (1 byte): Raw integer value (1-4)
func appendCompanyActivityType(dst []byte, companyActivityType ddv1.CompanyActivityType) []byte {
	// Get the protocol value for the enum
	protocolValue := GetCompanyActivityTypeProtocolValue(companyActivityType, 0)
	return append(dst, byte(protocolValue))
}
