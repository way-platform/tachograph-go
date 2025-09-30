package card

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// Test_roundTrip_rawCardFile tests that RawCardFile → Binary → RawCardFile conversion is 100% perfect
func Test_roundTrip_rawCardFile(t *testing.T) {
	// Dynamically discover test files
	files, err := os.ReadDir("../../testdata/card/driver")
	if err != nil {
		t.Fatalf("Failed to read testdata/card/driver directory: %v", err)
	}

	var testFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".DDD" {
			testFiles = append(testFiles, filepath.Join("../../testdata/card/driver", file.Name()))
		}
	}

	for _, filePath := range testFiles {
		t.Run(filePath, func(t *testing.T) {
			// Read original binary data
			originalData, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			// Step 1: Binary → RawCardFile
			originalRawFile, err := UnmarshalRawCardFile(originalData)
			if err != nil {
				t.Fatalf("Failed to unmarshal to RawCardFile: %v", err)
			}

			// Step 2: RawCardFile → Binary
			marshalledData, err := MarshalRawCardFile(originalRawFile)
			if err != nil {
				t.Fatalf("Failed to marshal RawCardFile to binary: %v", err)
			}

			// Step 3: Binary → RawCardFile (roundtrip)
			roundtripRawFile, err := UnmarshalRawCardFile(marshalledData)
			if err != nil {
				t.Fatalf("Failed to unmarshal roundtrip binary: %v", err)
			}

			// Validate binary data is identical
			if len(originalData) != len(marshalledData) {
				t.Errorf("Binary length mismatch: original=%d, marshalled=%d", len(originalData), len(marshalledData))
			}

			// Find first binary difference
			firstDiff := -1
			minLen := len(originalData)
			if len(marshalledData) < minLen {
				minLen = len(marshalledData)
			}

			for i := 0; i < minLen; i++ {
				if originalData[i] != marshalledData[i] {
					firstDiff = i
					break
				}
			}

			if firstDiff != -1 {
				t.Errorf("Binary data differs at byte %d: original=0x%02X, marshalled=0x%02X",
					firstDiff, originalData[firstDiff], marshalledData[firstDiff])

				// Show context around the difference
				start := firstDiff - 10
				if start < 0 {
					start = 0
				}
				end := firstDiff + 10
				if end > minLen {
					end = minLen
				}

				t.Logf("Context around difference:")
				t.Logf("  Original:   %X", originalData[start:end])
				t.Logf("  Marshalled: %X", marshalledData[start:end])
			} else if len(originalData) != len(marshalledData) {
				// Lengths differ but no byte difference found within common length
				t.Errorf("Binary lengths differ: original=%d, marshalled=%d", len(originalData), len(marshalledData))
			}

			// Validate RawCardFile structures are identical
			if diff := cmp.Diff(originalRawFile, roundtripRawFile, protocmp.Transform()); diff != "" {
				t.Errorf("RawCardFile structures differ (-original +roundtrip):\n%s", diff)
			}

			// Success metrics
			if firstDiff == -1 && len(originalData) == len(marshalledData) {
				t.Logf("✅ Perfect binary roundtrip: %d bytes", len(originalData))
				t.Logf("✅ Perfect structure roundtrip: %d records", len(originalRawFile.GetRecords()))
			}
		})
	}
}
