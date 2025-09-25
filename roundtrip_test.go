package tachograph

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestRoundtripCard tests that card files can be unmarshalled and marshalled back to identical binary data.
func TestRoundtripCard(t *testing.T) {
	cardDir := "testdata/card/driver"

	entries, err := os.ReadDir(cardDir)
	if err != nil {
		t.Fatalf("Failed to read card directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".DDD" {
			t.Run(entry.Name(), func(t *testing.T) {
				testRoundtripFile(t, filepath.Join(cardDir, entry.Name()))
			})
		}
	}
}

// TestRoundtripVU tests that VU files can be unmarshalled and marshalled back to identical binary data.
func TestRoundtripVU(t *testing.T) {
	vuDir := "testdata/vu"

	entries, err := os.ReadDir(vuDir)
	if err != nil {
		t.Fatalf("Failed to read VU directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".DDD" {
			t.Run(entry.Name(), func(t *testing.T) {
				testRoundtripFile(t, filepath.Join(vuDir, entry.Name()))
			})
		}
	}
}

// testRoundtripFile performs a roundtrip test on a single DDD file.
func testRoundtripFile(t *testing.T, filePath string) {
	// Step 1: Read the original file
	originalData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	t.Logf("Testing file: %s (size: %d bytes)", filePath, len(originalData))

	// Step 2: Unmarshal the data
	file, err := UnmarshalFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal file %s: %v", filePath, err)
	}

	t.Logf("Successfully unmarshalled file type: %v", file.GetType())

	// Step 3: Marshal the data back with original signatures preserved
	marshalledData, err := MarshalWithOriginal(file, originalData)
	if err != nil {
		t.Fatalf("Failed to marshal file %s: %v", filePath, err)
	}

	t.Logf("Successfully marshalled to %d bytes", len(marshalledData))

	// Step 4: Compare the binary data
	if len(originalData) != len(marshalledData) {
		t.Errorf("Length mismatch for %s: original=%d, marshalled=%d",
			filePath, len(originalData), len(marshalledData))

		// Provide detailed analysis of the length difference
		if len(marshalledData) > len(originalData) {
			t.Errorf("Marshalled data is %d bytes longer than original",
				len(marshalledData)-len(originalData))
		} else {
			t.Errorf("Marshalled data is %d bytes shorter than original",
				len(originalData)-len(marshalledData))
		}
	}

	if !bytes.Equal(originalData, marshalledData) {
		t.Errorf("Binary data mismatch for %s", filePath)

		// Find the first differing byte for debugging
		minLen := len(originalData)
		if len(marshalledData) < minLen {
			minLen = len(marshalledData)
		}

		for i := 0; i < minLen; i++ {
			if originalData[i] != marshalledData[i] {
				t.Errorf("First difference at byte %d: original=0x%02X, marshalled=0x%02X",
					i, originalData[i], marshalledData[i])

				// Show some context around the difference
				start := max(0, i-8)
				end := min(minLen, i+8)
				t.Errorf("Original context [%d:%d]: %X", start, end, originalData[start:end])
				t.Errorf("Marshalled context [%d:%d]: %X", start, end, marshalledData[start:end])
				break
			}
		}

		// If lengths differ but all compared bytes are equal, note that
		if minLen < len(originalData) || minLen < len(marshalledData) {
			t.Errorf("Length difference detected after %d matching bytes", minLen)
		}
	} else {
		t.Logf("✅ Perfect roundtrip: binary data matches exactly")
	}
}

// TestRoundtripAllFiles is a convenience test that runs all roundtrip tests.
func TestRoundtripAllFiles(t *testing.T) {
	t.Run("Card Files", func(t *testing.T) {
		TestRoundtripCard(t)
	})

	t.Run("VU Files", func(t *testing.T) {
		TestRoundtripVU(t)
	})
}

// BenchmarkRoundtripCard benchmarks the roundtrip performance for card files.
func BenchmarkRoundtripCard(b *testing.B) {
	cardDir := "testdata/card/driver"

	entries, err := os.ReadDir(cardDir)
	if err != nil {
		b.Fatalf("Failed to read card directory: %v", err)
	}

	// Use the first card file for benchmarking
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".DDD" {
			filePath := filepath.Join(cardDir, entry.Name())
			originalData, err := os.ReadFile(filePath)
			if err != nil {
				b.Fatalf("Failed to read file: %v", err)
			}

			b.Run(entry.Name(), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					file, err := UnmarshalFile(originalData)
					if err != nil {
						b.Fatalf("Unmarshal failed: %v", err)
					}

					_, err = Marshal(file)
					if err != nil {
						b.Fatalf("Marshal failed: %v", err)
					}
				}
			})
			break // Only benchmark the first file
		}
	}
}

// BenchmarkRoundtripVU benchmarks the roundtrip performance for VU files.
func BenchmarkRoundtripVU(b *testing.B) {
	vuDir := "testdata/vu"

	entries, err := os.ReadDir(vuDir)
	if err != nil {
		b.Fatalf("Failed to read VU directory: %v", err)
	}

	// Use the first VU file for benchmarking
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".DDD" {
			filePath := filepath.Join(vuDir, entry.Name())
			originalData, err := os.ReadFile(filePath)
			if err != nil {
				b.Fatalf("Failed to read file: %v", err)
			}

			b.Run(entry.Name(), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					file, err := UnmarshalFile(originalData)
					if err != nil {
						b.Fatalf("Unmarshal failed: %v", err)
					}

					_, err = Marshal(file)
					if err != nil {
						b.Fatalf("Marshal failed: %v", err)
					}
				}
			})
			break // Only benchmark the first file
		}
	}
}

// Helper functions for Go versions that don't have min/max built-ins
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
