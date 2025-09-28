package tachograph

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
//     ControlType ::= OCTET STRING (SIZE(1))
func unmarshalControlType(data []byte) (*ddv1.ControlType, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for control type")
	}
	b := data[0]
	ct := &ddv1.ControlType{}
	ct.SetCardDownloading((b & 0x80) != 0)
	ct.SetVuDownloading((b & 0x40) != 0)
	ct.SetPrinting((b & 0x20) != 0)
	ct.SetDisplay((b & 0x10) != 0)
	ct.SetCalibrationChecking((b & 0x08) != 0)
	return ct, nil
}
