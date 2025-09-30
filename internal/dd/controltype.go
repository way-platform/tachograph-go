package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalControlType unmarshals a control type from a byte slice
//
// The data type `ControlType` is specified in the Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//	ControlType ::= OCTET STRING (SIZE(1))
func UnmarshalControlType(input []byte) (*ddv1.ControlType, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("insufficient data for control type")
	}
	b := input[0]
	var output ddv1.ControlType
	output.SetRawValue(input)
	output.SetCardDownloading((b & 0x80) != 0)
	output.SetVuDownloading((b & 0x40) != 0)
	output.SetPrinting((b & 0x20) != 0)
	output.SetDisplay((b & 0x10) != 0)
	output.SetCalibrationChecking((b & 0x08) != 0)
	return &output, nil
}

// appendControlType appends a ControlType as a single byte bitmask.
//
// The data type `ControlType` is specified in the Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//     ControlType ::= OCTET STRING (SIZE(1))
