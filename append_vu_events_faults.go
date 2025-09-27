package tachograph

import (
	"bytes"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendVuEventsAndFaults appends VU events and faults data to a buffer.
func AppendVuEventsAndFaults(buf *bytes.Buffer, eventsAndFaults *vuv1.EventsAndFaults) error {
	if eventsAndFaults == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if eventsAndFaults.GetGeneration() == vuv1.Generation_GENERATION_1 {
		signature := eventsAndFaults.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := eventsAndFaults.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}

