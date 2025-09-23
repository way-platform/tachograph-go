package tachograph

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// TestRawCardFileRoundtrip tests that RawCardFile → Binary → RawCardFile conversion is 100% perfect
func TestRawCardFileRoundtrip(t *testing.T) {
	// Dynamically discover test files
	files, err := os.ReadDir("testdata/card")
	if err != nil {
		t.Fatalf("Failed to read testdata/card directory: %v", err)
	}

	var testFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".DDD" {
			testFiles = append(testFiles, filepath.Join("testdata/card", file.Name()))
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

// TestRawCardFileStructureConsistency tests that all card files have consistent RawCardFile structure
func TestRawCardFileStructureConsistency(t *testing.T) {
	testFiles := []string{
		"testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD",
		"testdata/card/proprietary-Omar_Khyam_Khawaja_2025-09-12_12-02-20.DDD",
		"testdata/card/proprietary-Teemu_Samuli_Hyvärinen_2025-09-12_12-03-47.DDD",
		"testdata/card/proprietary-Ville_Petteri_Kalske_2025-09-12_11-41-51.DDD",
	}

	type RecordSummary struct {
		Tag         int32
		File        string
		ContentType string
		Length      int32
	}

	var fileStructures []struct {
		path    string
		records []*RecordSummary
	}

	// Analyze each file structure
	for _, filePath := range testFiles {
		originalData, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", filePath, err)
		}

		rawFile, err := UnmarshalRawCardFile(originalData)
		if err != nil {
			t.Fatalf("Failed to unmarshal file %s: %v", filePath, err)
		}

		var records []*RecordSummary
		for _, record := range rawFile.GetRecords() {
			records = append(records, &RecordSummary{
				Tag:         record.GetTag(),
				File:        record.GetFile().String(),
				ContentType: record.GetContentType().String(),
				Length:      record.GetLength(),
			})
		}

		fileStructures = append(fileStructures, struct {
			path    string
			records []*RecordSummary
		}{filePath, records})

		t.Logf("File: %s", filePath)
		t.Logf("  Records: %d", len(records))
		for i, record := range records {
			t.Logf("  [%2d] Tag=0x%06X, File=%s, Type=%s, Length=%d",
				i, record.Tag, record.File, record.ContentType, record.Length)
		}
		t.Logf("")
	}

	// Compare structures (they should be very similar for same card type)
	if len(fileStructures) > 1 {
		baseStructure := fileStructures[0]
		for i := 1; i < len(fileStructures); i++ {
			compareStructure := fileStructures[i]

			if len(baseStructure.records) != len(compareStructure.records) {
				t.Logf("Structure difference: %s has %d records, %s has %d records",
					baseStructure.path, len(baseStructure.records),
					compareStructure.path, len(compareStructure.records))
			}

			minLen := len(baseStructure.records)
			if len(compareStructure.records) < minLen {
				minLen = len(compareStructure.records)
			}

			for j := 0; j < minLen; j++ {
				base := baseStructure.records[j]
				comp := compareStructure.records[j]

				if base.Tag != comp.Tag || base.File != comp.File || base.ContentType != comp.ContentType {
					t.Logf("Structure difference at record %d:", j)
					t.Logf("  %s: Tag=0x%06X, File=%s, Type=%s", baseStructure.path, base.Tag, base.File, base.ContentType)
					t.Logf("  %s: Tag=0x%06X, File=%s, Type=%s", compareStructure.path, comp.Tag, comp.File, comp.ContentType)
				}
			}
		}
	}
}
