package tachograph

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// TestSemanticRoundtrip tests the DriverCardFile → RawCardFile conversion step specifically
// This isolates semantic conversion issues from binary serialization issues
func TestSemanticRoundtrip(t *testing.T) {
	// Dynamically discover test files
	files, err := os.ReadDir("testdata/card/driver")
	if err != nil {
		t.Fatalf("Failed to read testdata/card/driver directory: %v", err)
	}

	var testFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".DDD" {
			testFiles = append(testFiles, filepath.Join("testdata/card/driver", file.Name()))
		}
	}

	for _, filePath := range testFiles {
		t.Run(filePath, func(t *testing.T) {
			// Read original binary data
			originalData, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			// Step 1: Binary → RawCardFile (ground truth)
			originalRawFile, err := unmarshalRawCardFile(originalData)
			if err != nil {
				t.Fatalf("Failed to unmarshal original binary to RawCardFile: %v", err)
			}

			// Step 2: Binary → DriverCardFile (semantic unmarshalling)
			file, err := UnmarshalFile(originalData)
			if err != nil {
				t.Fatalf("Failed to unmarshal to DriverCardFile: %v", err)
			}

			driverCard := file.GetDriverCard()
			if driverCard == nil {
				t.Fatalf("File is not a driver card")
			}

			// Step 3: DriverCardFile → RawCardFile (semantic marshalling - this is what we're testing)
			// Pass the original RawCardFile to preserve signatures
			marshalledRawFile, err := DriverCardFileToRawWithSignatures(driverCard, originalRawFile)
			if err != nil {
				t.Fatalf("Failed to convert DriverCardFile to RawCardFile: %v", err)
			}

			// Compare the raw structures
			t.Logf("=== SEMANTIC CONVERSION ANALYSIS ===")
			t.Logf("Original records: %d", len(originalRawFile.GetRecords()))
			t.Logf("Marshalled records: %d", len(marshalledRawFile.GetRecords()))

			// Detailed record-by-record comparison
			originalRecords := originalRawFile.GetRecords()
			marshalledRecords := marshalledRawFile.GetRecords()

			maxRecords := len(originalRecords)
			if len(marshalledRecords) > maxRecords {
				maxRecords = len(marshalledRecords)
			}

			perfectMatches := 0
			recordIssues := 0

			for i := 0; i < maxRecords; i++ {
				t.Logf("\n--- Record %d ---", i)

				if i >= len(originalRecords) {
					t.Logf("❌ Missing in original: %s", formatRecord(marshalledRecords[i]))
					recordIssues++
					continue
				}

				if i >= len(marshalledRecords) {
					t.Logf("❌ Missing in marshalled: %s", formatRecord(originalRecords[i]))
					recordIssues++
					continue
				}

				orig := originalRecords[i]
				marsh := marshalledRecords[i]

				t.Logf("Original:   %s", formatRecord(orig))
				t.Logf("Marshalled: %s", formatRecord(marsh))

				// Check each field
				issues := []string{}

				if orig.GetTag() != marsh.GetTag() {
					issues = append(issues, "Tag")
				}
				if orig.GetFile() != marsh.GetFile() {
					issues = append(issues, "File")
				}
				if orig.GetContentType() != marsh.GetContentType() {
					issues = append(issues, "ContentType")
				}
				if orig.GetLength() != marsh.GetLength() {
					issues = append(issues, "Length")
				}

				// Compare values using bytes
				if !bytesEqual(orig.GetValue(), marsh.GetValue()) {
					issues = append(issues, "Value")
				}

				if len(issues) == 0 {
					t.Logf("✅ Perfect match")
					perfectMatches++
				} else {
					t.Logf("❌ Issues: %v", issues)
					recordIssues++

					// Show detailed value comparison for mismatches
					if contains(issues, "Value") {
						showValueDiff(t, orig.GetValue(), marsh.GetValue())
					}
				}
			}

			// Summary
			t.Logf("\n=== SUMMARY ===")
			t.Logf("Perfect matches: %d", perfectMatches)
			t.Logf("Records with issues: %d", recordIssues)

			if recordIssues == 0 {
				t.Logf("🎉 PERFECT SEMANTIC CONVERSION!")
			} else {
				// Use protocmp for detailed structural diff
				t.Logf("\n=== DETAILED PROTOBUF DIFF ===")
				if diff := cmp.Diff(originalRawFile, marshalledRawFile, protocmp.Transform()); diff != "" {
					t.Logf("RawCardFile structures differ (-original +marshalled):\n%s", diff)
				}
			}
		})
	}
}

// formatRecord creates a human-readable summary of a RawCardFile_Record
func formatRecord(record interface{}) string {
	if record == nil {
		return "nil"
	}

	// Import the actual protobuf types
	if r, ok := record.(*cardv1.RawCardFile_Record); ok {
		return fmt.Sprintf("Tag=0x%06X, File=%s, Type=%s, Length=%d",
			r.GetTag(), r.GetFile().String(), r.GetContentType().String(), r.GetLength())
	}

	return "unknown"
}

// bytesEqual compares two byte slices
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// showValueDiff shows a detailed comparison of byte values
func showValueDiff(t *testing.T, original, marshalled []byte) {
	t.Logf("  Value length: original=%d, marshalled=%d", len(original), len(marshalled))

	// Show first few bytes for context
	maxShow := 32
	if len(original) < maxShow {
		maxShow = len(original)
	}
	if len(marshalled) < maxShow && len(marshalled) > maxShow {
		maxShow = len(marshalled)
	}

	if maxShow > 0 {
		origHex := fmt.Sprintf("%X", original[:minValue(maxShow, len(original))])
		marshHex := fmt.Sprintf("%X", marshalled[:minValue(maxShow, len(marshalled))])

		t.Logf("  Original (first %d bytes):   %s", minValue(maxShow, len(original)), origHex)
		t.Logf("  Marshalled (first %d bytes): %s", minValue(maxShow, len(marshalled)), marshHex)

		// Find first difference
		minLen := len(original)
		if len(marshalled) < minLen {
			minLen = len(marshalled)
		}

		for i := 0; i < minLen; i++ {
			if original[i] != marshalled[i] {
				t.Logf("  First difference at byte %d: original=0x%02X, marshalled=0x%02X",
					i, original[i], marshalled[i])
				break
			}
		}
	}
}

// minValue returns the minimum of two integers
func minValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}
