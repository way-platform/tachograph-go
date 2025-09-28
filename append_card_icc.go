package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendIcc appends the binary representation of an EF_ICC message to dst.
//
// ASN.1 Specification (Data Dictionary 2.23):
// CardIccIdentification ::= SEQUENCE {
//     clockStop                   OCTET STRING (SIZE(1)),
//     cardExtendedSerialNumber    ExtendedSerialNumber,    -- 8 bytes
//     cardApprovalNumber          CardApprovalNumber,      -- 8 bytes
//     cardPersonaliserID          ManufacturerCode,        -- 1 byte
//     embedderIcAssemblerId       EmbedderIcAssemblerId,   -- 5 bytes
//     icIdentifier                OCTET STRING (SIZE(2))
// }
func AppendIcc(dst []byte, icc *cardv1.Icc) ([]byte, error) {
	const (
		// CardIccIdentification layout constants
		lenClockStop                = 1
		lenCardExtendedSerialNumber = 8
		lenCardApprovalNumber       = 8
		lenCardPersonaliserId       = 1
		lenEmbedderIcAssemblerId    = 5
		lenIcIdentifier             = 2
	)

	var err error
	// Append clock stop (1 byte)
	dst = append(dst, byte(icc.GetClockStop()))
	
	// Append extended serial number (8 bytes)
	dst, err = appendExtendedSerialNumber(dst, icc.GetCardExtendedSerialNumber(), lenCardExtendedSerialNumber)
	if err != nil {
		return nil, err
	}
	
	// Append card approval number (8 bytes)
	dst, err = appendStringValue(dst, icc.GetCardApprovalNumber(), lenCardApprovalNumber)
	if err != nil {
		return nil, err
	}
	
	// Append card personaliser ID (1 byte)
	dst = append(dst, byte(icc.GetCardPersonaliserId()))
	
	// Append embedder IC assembler ID (5 bytes)
	dst, err = appendEmbedderIcAssemblerId(dst, icc.GetEmbedderIcAssemblerId())
	if err != nil {
		return nil, err
	}
	
	// Append IC identifier (2 bytes)
	dst = append(dst, icc.GetIcIdentifier()...)
	return dst, nil
}
