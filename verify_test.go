package tachograph_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/way-platform/tachograph-go"
)

func TestVerifyFile_goldenFiles(t *testing.T) {
	testDataDir := "testdata/card/driver"
	entries, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Skipf("test data directory not available: %v", err)
	}
	// Find all .DDD files in the test data directory
	var testFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".DDD" {
			testFiles = append(testFiles, filepath.Join(testDataDir, entry.Name()))
		}
	}
	if len(testFiles) == 0 {
		t.Skip("no test files found")
	}
	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			// Read the test file
			data, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}
			// Parse the file
			file, err := tachograph.UnmarshalFile(data)
			if err != nil {
				t.Fatalf("failed to unmarshal file: %v", err)
			}
			// Verify the certificates using the default resolver
			// Golden files should be valid, so verification should succeed
			err = tachograph.VerifyFile(t.Context(), file)
			if err != nil {
				t.Errorf("certificate verification failed: %v", err)
			}
		})
	}
}
