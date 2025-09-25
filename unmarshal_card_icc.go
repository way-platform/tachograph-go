package tachograph

import (
	"bytes"
	"errors"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalIcc parses the binary data for an EF_ICC record.
func unmarshalIcc(data []byte) (*cardv1.IccIdentification, error) {
	var icc cardv1.IccIdentification
	if len(data) < 25 {
		return nil, errors.New("not enough data for IccIdentification")
	}
	r := bytes.NewReader(data)
	clockStop, _ := r.ReadByte()
	icc.SetClockStop(int32(clockStop))
	icc.SetCardExtendedSerialNumber(readString(r, 8))
	icc.SetCardApprovalNumber(readString(r, 8))
	personaliser, _ := r.ReadByte()
	icc.SetCardPersonaliserId(int32(personaliser))
	embedder := make([]byte, 5)
	r.Read(embedder)
	icc.SetEmbedderIcAssemblerId(embedder)
	icIdentifier := make([]byte, 2)
	r.Read(icIdentifier)
	icc.SetIcIdentifier(icIdentifier)
	return &icc, nil
}
