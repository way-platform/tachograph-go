package tachograph_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/way-platform/tachograph-go"
)

// TestVerifyFile tests certificate verification for driver card files.
// This test demonstrates that the VerifyFile function can successfully verify
// certificates in real driver card files using the default certificate resolver.
//
// The default resolver fetches CA certificates from embedded stores and remote
// sources, providing complete verification chains.
func TestVerifyFile(t *testing.T) {
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

			// Check that signature_valid flag was set to true
			if file.GetDriverCard() != nil {
				if tach := file.GetDriverCard().GetTachograph(); tach != nil {
					if cert := tach.GetCardCertificate().GetRsaCertificate(); cert != nil {
						if !cert.GetSignatureValid() {
							t.Error("Gen1 card certificate signature_valid is false, expected true")
						}
					}
				}
				if tachG2 := file.GetDriverCard().GetTachographG2(); tachG2 != nil {
					if cert := tachG2.GetCardSignCertificate().GetEccCertificate(); cert != nil {
						if !cert.GetSignatureValid() {
							t.Error("Gen2 card sign certificate signature_valid is false, expected true")
						}
					}
				}
			}
		})
	}
}
