package tachograph

import (
	"bytes"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendVuTechnicalData appends VU technical data to a buffer.
func AppendVuTechnicalData(buf *bytes.Buffer, technicalData *vuv1.TechnicalData) error {
	if technicalData == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if technicalData.GetGeneration() == vuv1.Generation_GENERATION_1 {
		signature := technicalData.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := technicalData.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}

