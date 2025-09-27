package tachograph

import (
	"bytes"
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuEventsAndFaults unmarshals VU events and faults data from a VU transfer.
func UnmarshalVuEventsAndFaults(r *bytes.Reader, target *vuv1.EventsAndFaults, generation int) (int, error) {
	initialLen := r.Len()

	// Set generation
	if generation == 1 {
		target.SetGeneration(vuv1.Generation_GENERATION_1)
	} else {
		target.SetGeneration(vuv1.Generation_GENERATION_2)
	}

	// For now, implement a simplified version that just reads the data
	// without fully parsing all the complex structures
	// This ensures the interface is complete while allowing for future enhancement

	// Read all remaining data as signature for now
	remainingData := make([]byte, r.Len())
	if _, err := r.Read(remainingData); err != nil {
		return 0, fmt.Errorf("failed to read events and faults data: %w", err)
	}

	// Set as signature based on generation
	if generation == 1 {
		target.SetSignatureGen1(remainingData)
	} else {
		target.SetSignatureGen2(remainingData)
	}

	return initialLen - r.Len(), nil
}

