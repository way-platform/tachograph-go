package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardIc appends IC identification data to a byte slice.
func AppendCardIc(data []byte, ic *cardv1.Ic) ([]byte, error) {
	if ic == nil {
		return data, nil
	}

	// IC Serial Number (4 bytes)
	serialBytes := ic.GetIcSerialNumber()
	if len(serialBytes) >= 4 {
		data = append(data, serialBytes[:4]...)
	} else {
		// Pad with zeros
		padded := make([]byte, 4)
		copy(padded, serialBytes)
		data = append(data, padded...)
	}

	// IC Manufacturing References (4 bytes)
	mfgBytes := ic.GetIcManufacturingReferences()
	if len(mfgBytes) >= 4 {
		data = append(data, mfgBytes[:4]...)
	} else {
		// Pad with zeros
		padded := make([]byte, 4)
		copy(padded, mfgBytes)
		data = append(data, padded...)
	}

	return data, nil
}
