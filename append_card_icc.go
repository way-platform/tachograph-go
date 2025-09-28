package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendIcc appends the binary representation of an EF_ICC message to dst.
func AppendIcc(dst []byte, icc *cardv1.Icc) ([]byte, error) {
	var err error
	dst = append(dst, byte(icc.GetClockStop()))
	dst, err = appendExtendedSerialNumber(dst, icc.GetCardExtendedSerialNumber(), 8)
	if err != nil {
		return nil, err
	}
	dst, err = appendStringValue(dst, icc.GetCardApprovalNumber(), 8)
	if err != nil {
		return nil, err
	}
	dst = append(dst, byte(icc.GetCardPersonaliserId()))
	dst, err = appendEmbedderIcAssemblerId(dst, icc.GetEmbedderIcAssemblerId())
	if err != nil {
		return nil, err
	}
	dst = append(dst, icc.GetIcIdentifier()...)
	return dst, nil
}
