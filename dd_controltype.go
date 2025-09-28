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

// appendControlType appends a ControlType as a single byte bitmask.
//
// The data type `ControlType` is specified in the Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//     ControlType ::= OCTET STRING (SIZE(1))
func appendControlType(dst []byte, ct *ddv1.ControlType) []byte {
	if ct == nil {
		return append(dst, 0)
	}

	var b byte
	if ct.GetCardDownloading() {
		b |= 0x80 // bit 'c'
	}
	if ct.GetVuDownloading() {
		b |= 0x40 // bit 'v'
	}
	if ct.GetPrinting() {
		b |= 0x20 // bit 'p'
	}
	if ct.GetDisplay() {
		b |= 0x10 // bit 'd'
	}
	if ct.GetCalibrationChecking() {
		b |= 0x08 // bit 'e'
	}
	// bits 0-2 are RFU (Reserved for Future Use)

	return append(dst, b)
}
