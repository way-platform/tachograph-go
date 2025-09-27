package tachograph

import (
	"bytes"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendVuDetailedSpeed appends VU detailed speed data to a buffer.
func AppendVuDetailedSpeed(buf *bytes.Buffer, detailedSpeed *vuv1.DetailedSpeed) error {
	if detailedSpeed == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if detailedSpeed.GetGeneration() == datadictionaryv1.Generation_GENERATION_1 {
		signature := detailedSpeed.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := detailedSpeed.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}
