package tachograph

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardIc unmarshals IC identification data from EF_IC.
func UnmarshalCardIc(data []byte, target *cardv1.ChipIdentification) error {
	if len(data) < 8 {
		return fmt.Errorf("insufficient data for IC identification: got %d bytes, need 8", len(data))
	}

	r := bytes.NewReader(data)

	// According to Data Dictionary Section 2.13, EF_IC contains:
	// - IC Serial Number (4 bytes)
	// - IC Manufacturing References (4 bytes)

	// Read IC Serial Number (4 bytes)
	serialBytes := make([]byte, 4)
	if _, err := r.Read(serialBytes); err != nil {
		return fmt.Errorf("failed to read IC serial number: %w", err)
	}
	target.SetIcSerialNumber(fmt.Sprintf("%08X", serialBytes))

	// Read IC Manufacturing References (4 bytes)
	mfgBytes := make([]byte, 4)
	if _, err := r.Read(mfgBytes); err != nil {
		return fmt.Errorf("failed to read IC manufacturing references: %w", err)
	}
	target.SetIcManufacturingReferences(fmt.Sprintf("%08X", mfgBytes))

	return nil
}
