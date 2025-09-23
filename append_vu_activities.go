package tachograph

import (
	"bytes"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendVuActivities appends VU activities data to a buffer.
func AppendVuActivities(buf *bytes.Buffer, activities *vuv1.Activities) error {
	if activities == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if activities.GetGeneration() == vuv1.Generation_GENERATION_1 {
		signature := activities.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := activities.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}
