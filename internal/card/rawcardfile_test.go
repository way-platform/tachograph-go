package card

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

// Test_roundTrip_rawCardFile tests that RawCardFile → Binary → RawCardFile conversion is 100% perfect
func Test_roundTrip_rawCardFile(t *testing.T) {
	// Check if testdata directory exists
	if _, err := os.Stat("../../testdata/card/driver"); err != nil {
		if os.IsNotExist(err) {
			t.Skip("testdata/card/driver directory not present (proprietary test files not available)")
		}
		t.Fatalf("Failed to access testdata/card/driver directory: %v", err)
	}

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

func TestUnmarshalRawCardFile_golden(t *testing.T) {
	// Check if testdata/card directory exists
	if _, err := os.Stat("../../testdata/card"); err != nil {
		if os.IsNotExist(err) {
			t.Skip("testdata/card directory not present (proprietary test files not available)")
		}
		t.Fatalf("Failed to access testdata/card directory: %v", err)
	}

	if err := filepath.WalkDir("../../testdata/card", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".DDD") {
			return nil
		}

		// Create relative path from testdata/card root (e.g., "driver/file.DDD")
		relPath, err := filepath.Rel("../../testdata/card", path)
		if err != nil {
			t.Fatalf("Failed to compute relative path: %v", err)
		}

		// Create golden file path in internal/card/testdata/raw/card with same structure
		goldenFile := filepath.Join("testdata/raw/card", strings.TrimSuffix(relPath, ".DDD")+".json")

		t.Run(path, func(t *testing.T) {
			// Read and parse the DDD file
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read DDD file %s: %v", path, err)
			}

			var actual string
			rawFile, err := UnmarshalRawCardFile(data)
			if err != nil {
				// If parsing fails, use the error message as the golden content
				actual = `{"error":"` + err.Error() + `"}`
			} else {
				// Validate the parsed file against protovalidate annotations
				validator, err := protovalidate.New()
				if err != nil {
					t.Fatalf("Failed to create validator: %v", err)
				}
				if err := validator.Validate(rawFile); err != nil {
					t.Errorf("Validation failed for %s: %v", path, err)
				}

				// If parsing succeeds, use the JSON representation
				actualBytes, err := (protojson.MarshalOptions{}).Marshal(rawFile)
				if err != nil {
					t.Fatalf("Failed to marshal RawCardFile %s: %v", path, err)
				}
				var actualIndented bytes.Buffer
				_ = json.Indent(&actualIndented, actualBytes, "", "  ") // ignore error as JSON marshaling already succeeded
				actual = actualIndented.String()
			}

			if *update {
				// Ensure the directory exists
				goldenDir := filepath.Dir(goldenFile)
				if err := os.MkdirAll(goldenDir, 0o755); err != nil {
					t.Fatalf("Failed to create golden file directory %s: %v", goldenDir, err)
				}

				if err := os.WriteFile(goldenFile, []byte(actual), 0o644); err != nil {
					t.Fatalf("Failed to write golden file %s: %v", goldenFile, err)
				}
				t.Logf("Updated golden file: %s", goldenFile)
				return
			}

			expected, err := os.ReadFile(goldenFile)
			if err != nil {
				if os.IsNotExist(err) {
					t.Fatalf("Golden file %s does not exist. Run with -update to create it.", goldenFile)
				}
				t.Fatalf("Failed to read golden file %s: %v", goldenFile, err)
			}

			if diff := cmp.Diff(string(expected), actual); diff != "" {
				t.Errorf("Golden file mismatch for %s (-expected +actual):\n%s", path, diff)
				t.Logf("To update the golden file, run: go test -update")
			}
		})
		return nil
	}); err != nil {
		t.Fatalf("Failed to walk testdata/card directory: %v", err)
	}
}
