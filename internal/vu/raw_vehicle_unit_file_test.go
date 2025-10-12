package vu

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
)

var update = flag.Bool("update", false, "update golden files")

// TestRawVehicleUnitFileGolden is a golden file test for RawVehicleUnitFile unmarshaling.
//
// Run with -update to regenerate golden files:
//
//	go test -run TestRawVehicleUnitFileGolden -update -v
//
// This test verifies that we correctly identify transfer record boundaries and preserve
// the exact binary values for later parsing and signature verification.
func TestRawVehicleUnitFileGolden(t *testing.T) {
	// Find all VU test files
	testFiles, err := filepath.Glob("../../testdata/vu/*.DDD")
	if err != nil {
		t.Fatalf("failed to glob test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("no VU test files found")
	}

	for _, testFile := range testFiles {
		baseName := strings.TrimSuffix(filepath.Base(testFile), ".DDD")
		t.Run(baseName, func(t *testing.T) {
			// Read VU file
			data, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			// Unmarshal to RawVehicleUnitFile
			rawFile, err := unmarshalRawVehicleUnitFile(data)
			if err != nil {
				t.Fatalf("unmarshalRawVehicleUnitFile failed: %v", err)
			}

			if rawFile == nil {
				t.Fatal("unmarshalRawVehicleUnitFile returned nil")
			}

			// Verify we got some records
			if len(rawFile.GetRecords()) == 0 {
				t.Errorf("expected at least one record, got 0")
			}

			// Convert to JSON (with pretty printing for readability)
			jsonMarshaler := protojson.MarshalOptions{
				Multiline:       true,
				Indent:          "  ",
				EmitUnpopulated: false,
			}
			gotJSON, err := jsonMarshaler.Marshal(rawFile)
			if err != nil {
				t.Fatalf("failed to marshal to JSON: %v", err)
			}

			// Golden file path in the testdata directory
			goldenPath := filepath.Join("testdata", baseName+".raw.golden.json")

			if *update {
				// Update mode: write the golden file
				if err := os.MkdirAll("testdata", 0o755); err != nil {
					t.Fatalf("failed to create testdata directory: %v", err)
				}
				if err := os.WriteFile(goldenPath, gotJSON, 0o644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenPath)
			} else {
				// Comparison mode: read and compare with golden file
				wantJSON, err := os.ReadFile(goldenPath)
				if err != nil {
					t.Fatalf("failed to read golden file (run with -update to create): %v", err)
				}

				// Unmarshal both to compare structures (ignores whitespace differences)
				var got, want map[string]interface{}
				if err := json.Unmarshal(gotJSON, &got); err != nil {
					t.Fatalf("failed to unmarshal got JSON: %v", err)
				}
				if err := json.Unmarshal(wantJSON, &want); err != nil {
					t.Fatalf("failed to unmarshal want JSON: %v", err)
				}

				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("RawVehicleUnitFile mismatch (-want +got):\n%s", diff)
				}
			}

			// Log summary
			t.Logf("Successfully parsed %d transfer records", len(rawFile.GetRecords()))
			for _, record := range rawFile.GetRecords() {
				t.Logf("  - %s (tag=0x%04X, gen=%s, value=%d bytes)",
					record.GetType().String(),
					record.GetTag(),
					record.GetGeneration().String(),
					len(record.GetValue()))
			}
		})
	}
}

// TestRawVehicleUnitFileRoundTrip verifies binary fidelity by checking that
// the concatenation of all record values (with tags) reconstructs the original file.
func TestRawVehicleUnitFileRoundTrip(t *testing.T) {
	// Find all VU test files
	testFiles, err := filepath.Glob("../../testdata/vu/*.DDD")
	if err != nil {
		t.Fatalf("failed to glob test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("no VU test files found")
	}

	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			// Read VU file
			originalData, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			// Unmarshal to RawVehicleUnitFile
			rawFile, err := unmarshalRawVehicleUnitFile(originalData)
			if err != nil {
				t.Fatalf("unmarshalRawVehicleUnitFile failed: %v", err)
			}

			// Reconstruct binary by concatenating tags and values
			var reconstructed []byte
			for _, record := range rawFile.GetRecords() {
				// Append 2-byte tag (big-endian)
				tag := uint16(record.GetTag())
				reconstructed = append(reconstructed, byte(tag>>8), byte(tag))
				// Append value
				reconstructed = append(reconstructed, record.GetValue()...)
			}

			// Compare binary data
			if diff := cmp.Diff(originalData, reconstructed); diff != "" {
				t.Errorf("Binary round-trip mismatch (-original +reconstructed):\n%s", diff)
				t.Errorf("Original length: %d, Reconstructed length: %d", len(originalData), len(reconstructed))
			}
		})
	}
}
