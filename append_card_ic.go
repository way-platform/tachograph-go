package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardIc appends IC identification data to a byte slice.
//
// ASN.1 Specification (Data Dictionary 2.13):
//
//	CardChipIdentification ::= SEQUENCE {
//	    icSerialNumber              OCTET STRING (SIZE(4)),
//	    icManufacturingReferences   OCTET STRING (SIZE(4))
//	}
func AppendCardIc(data []byte, ic *cardv1.Ic) ([]byte, error) {
	const (
		// CardChipIdentification layout constants
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
