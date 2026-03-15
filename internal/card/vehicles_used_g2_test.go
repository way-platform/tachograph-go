package card

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/way-platform/tachograph-go/internal/dd"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestVehiclesUsedG2_Generation2(t *testing.T) {
	// Discover all matching hexdump files using type-safe enums
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_VEHICLES_USED,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("no hexdump files found for EF_VEHICLES_USED GENERATION_2 (run extract-testdata-records to regenerate)")
	}

	// Run subtest for each discovered file
	for _, hexdumpPath := range hexdumpFiles {
		// Use relative path from testdata as subtest name
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			// Read hexdump
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			// Unmarshal for golden comparison (no raw_data for readable JSON)
			opts := UnmarshalOptions{}
			vehicles, err := opts.unmarshalVehiclesUsedG2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, vehicles, goldenPath)

			// Round-trip test: unmarshal with PreserveRawData for binary fidelity.
			// Without raw_data, null-padded strings lose their padding and get
			// re-encoded with space padding, breaking the round-trip.
			rtOpts := UnmarshalOptions{UnmarshalOptions: dd.UnmarshalOptions{PreserveRawData: true}}
			rtVehicles, err := rtOpts.unmarshalVehiclesUsedG2(data)
			if err != nil {
				t.Fatalf("Round-trip unmarshal failed: %v", err)
			}
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalVehiclesUsedG2(rtVehicles)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
