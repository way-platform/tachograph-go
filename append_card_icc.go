package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendIcc appends the binary representation of an EF_ICC message to dst.
func AppendIcc(dst []byte, icc *cardv1.IccIdentification) ([]byte, error) {
	dst = append(dst, byte(icc.GetClockStop()))
	dst = appendString(dst, icc.GetCardExtendedSerialNumber(), 8)
	dst = appendString(dst, icc.GetCardApprovalNumber(), 8)
	dst = append(dst, byte(icc.GetCardPersonaliserId()))
	dst = append(dst, icc.GetEmbedderIcAssemblerId()...)
	dst = append(dst, icc.GetIcIdentifier()...)
	return dst, nil
}
