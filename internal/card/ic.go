package card

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardIc unmarshals IC identification data from EF_IC.
//
// The data type `CardChipIdentification` is specified in the Data Dictionary, Section 2.13.
//
// ASN.1 Definition:
//
//	CardChipIdentification ::= SEQUENCE {
//	    icSerialNumber              OCTET STRING (SIZE(4)),
//	    icManufacturingReferences   OCTET STRING (SIZE(4))
//	}
func (opts UnmarshalOptions) unmarshalIc(data []byte) (*cardv1.Ic, error) {
	const (
		lenIcSerialNumber            = 4
		lenIcManufacturingReferences = 4
		lenCardChipIdentification    = lenIcSerialNumber + lenIcManufacturingReferences
	)

	if len(data) < lenCardChipIdentification {
		return nil, fmt.Errorf("insufficient data for IC identification: got %d bytes, need %d", len(data), lenCardChipIdentification)
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

// AppendCardIc appends IC identification data to a byte slice.
//
// The data type `CardChipIdentification` is specified in the Data Dictionary, Section 2.13.
//
// ASN.1 Definition:
//
//	CardChipIdentification ::= SEQUENCE {
//	    icSerialNumber              OCTET STRING (SIZE(4)),
//	    icManufacturingReferences   OCTET STRING (SIZE(4))
//	}
func appendCardIc(data []byte, ic *cardv1.Ic) ([]byte, error) {
	const (
		lenIcSerialNumber            = 4
		lenIcManufacturingReferences = 4
	)

	if ic == nil {
		return data, nil
	}

	// Append IC Serial Number (4 bytes)
	serialBytes := ic.GetIcSerialNumber()
	if len(serialBytes) >= lenIcSerialNumber {
		data = append(data, serialBytes[:lenIcSerialNumber]...)
	} else {
		// Pad with zeros
		padded := make([]byte, lenIcSerialNumber)
		copy(padded, serialBytes)
		data = append(data, padded...)
	}

	// Append IC Manufacturing References (4 bytes)
	mfgBytes := ic.GetIcManufacturingReferences()
	if len(mfgBytes) >= lenIcManufacturingReferences {
		data = append(data, mfgBytes[:lenIcManufacturingReferences]...)
	} else {
		// Pad with zeros
		padded := make([]byte, lenIcManufacturingReferences)
		copy(padded, mfgBytes)
		data = append(data, padded...)
	}

	return data, nil
}
