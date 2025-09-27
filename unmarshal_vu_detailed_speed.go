package tachograph

import (
	"bytes"
	"fmt"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuDetailedSpeed unmarshals VU detailed speed data from a VU transfer.
func UnmarshalVuDetailedSpeed(r *bytes.Reader, target *vuv1.DetailedSpeed, generation int) (int, error) {
	initialLen := r.Len()

	// Set generation
	if generation == 1 {
		target.SetGeneration(datadictionaryv1.Generation_GENERATION_1)
	} else {
		target.SetGeneration(datadictionaryv1.Generation_GENERATION_2)
	}

	// For now, implement a simplified version that just reads the data
	// This ensures the interface is complete while allowing for future enhancement

	// Read all remaining data
	remainingData := make([]byte, r.Len())
	if _, err := r.Read(remainingData); err != nil {
		return 0, fmt.Errorf("failed to read detailed speed data: %w", err)
	}

	// Set as signature based on generation
	if generation == 1 {
		target.SetSignatureGen1(remainingData)
	} else {
		target.SetSignatureGen2(remainingData)
	}

	return initialLen - r.Len(), nil
}
