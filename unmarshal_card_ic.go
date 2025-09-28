package tachograph

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalCardIc unmarshals IC identification data from EF_IC.
//
// ASN.1 Specification (Data Dictionary 2.13):
//
//	CardChipIdentification ::= SEQUENCE {
//	    icSerialNumber              OCTET STRING (SIZE(4)),
//	    icManufacturingReferences   OCTET STRING (SIZE(4))
//	}
func unmarshalCardIc(data []byte) (*cardv1.Ic, error) {
	const (
		// CardChipIdentification layout constants
		lenIcSerialNumber            = 4
		lenIcManufacturingReferences = 4
		totalLength                  = lenIcSerialNumber + lenIcManufacturingReferences
	)

	if len(data) < totalLength {
		return nil, fmt.Errorf("insufficient data for IC identification: got %d bytes, need %d", len(data), totalLength)
	}

	var target cardv1.Ic
	r := bytes.NewReader(data)

	// Read IC Serial Number (4 bytes)
	serialBytes := make([]byte, lenIcSerialNumber)
	if _, err := r.Read(serialBytes); err != nil {
		return nil, fmt.Errorf("failed to read IC serial number: %w", err)
	}
	target.SetIcSerialNumber(serialBytes)

	// Read IC Manufacturing References (4 bytes)
	mfgBytes := make([]byte, lenIcManufacturingReferences)
	if _, err := r.Read(mfgBytes); err != nil {
		return nil, fmt.Errorf("failed to read IC manufacturing references: %w", err)
	}
	target.SetIcManufacturingReferences(mfgBytes)

	return &target, nil
}

// UnmarshalCardIc unmarshals IC identification data from EF_IC (legacy function).
// Deprecated: Use unmarshalCardIc instead.
func UnmarshalCardIc(data []byte, target *cardv1.Ic) error {
	result, err := unmarshalCardIc(data)
	if err != nil {
		return err
	}
	proto.Merge(target, result)
	return nil
}
