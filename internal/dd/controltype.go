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
	output.SetRawData(input)
	output.SetCardDownloading((b & 0x80) != 0)
	output.SetVuDownloading((b & 0x40) != 0)
	output.SetPrinting((b & 0x20) != 0)
	output.SetDisplay((b & 0x10) != 0)
	output.SetCalibrationChecking((b & 0x08) != 0)
	return &output, nil
}

// AppendControlType appends a ControlType as a single byte bitmask.
//
// The data type `ControlType` is specified in the Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//	ControlType ::= OCTET STRING (SIZE(1))
//
// Binary Layout (1 byte):
//   - Bit 7: card downloading
//   - Bit 6: VU downloading
//   - Bit 5: printing
//   - Bit 4: display
//   - Bit 3: calibration checking (Gen2+)
//   - Bits 2-0: Reserved (RFU)
func AppendControlType(dst []byte, controlType *ddv1.ControlType) ([]byte, error) {
	const lenControlType = 1
	var canvas [lenControlType]byte
	if controlType.HasRawData() {
		if len(controlType.GetRawData()) != lenControlType {
			return nil, fmt.Errorf(
				"invalid raw_data length for ControlType: got %d, want %d",
				len(controlType.GetRawData()), lenControlType,
			)
		}
		copy(canvas[:], controlType.GetRawData())
	}
	if controlType.GetCardDownloading() {
		canvas[0] |= 0x80
	}
	if controlType.GetVuDownloading() {
		canvas[0] |= 0x40
	}
	if controlType.GetPrinting() {
		canvas[0] |= 0x20
	}
	if controlType.GetDisplay() {
		canvas[0] |= 0x10
	}
	if controlType.GetCalibrationChecking() {
		canvas[0] |= 0x08
	}
	return append(dst, canvas[:]...), nil
}
