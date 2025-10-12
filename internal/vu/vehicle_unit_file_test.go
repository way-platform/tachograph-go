package vu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestUnmarshalVehicleUnitFile tests the full semantic parsing of VU files.
//
// This test verifies that we can parse VU files all the way to the
// semantic VehicleUnitFile message, including Overview transfers.
func TestUnmarshalVehicleUnitFile(t *testing.T) {
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
			data, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			// Unmarshal to VehicleUnitFile (full semantic parsing)
			vuFile, err := unmarshalVehicleUnitFile(data)
			if err != nil {
				t.Fatalf("unmarshalVehicleUnitFile failed: %v", err)
			}

			if vuFile == nil {
				t.Fatal("unmarshalVehicleUnitFile returned nil")
			}

		// Verify we got generation-specific data
		generation := vuFile.GetGeneration()
		t.Logf("Successfully parsed VU file: generation=%s", generation)

		// Check generation-specific fields
		switch generation {
		case ddv1.Generation_GENERATION_1:
			gen1File := vuFile.GetGen1()
			if gen1File == nil {
				t.Fatal("expected Gen1 file data, got nil")
			}

			// Check Overview
			if overview := gen1File.GetOverview(); overview != nil {
				t.Logf("  - Overview: VIN=%s, certificates present=%v",
					overview.GetVehicleIdentificationNumber().GetValue(),
					len(overview.GetMemberStateCertificate()) > 0)

				// Verify raw_data is preserved
				if len(overview.GetRawData()) == 0 {
					t.Errorf("    - Overview raw_data is empty")
				}
			}

			// Log other transfers
			t.Logf("  - Activities: %d records", len(gen1File.GetActivities()))
			t.Logf("  - EventsAndFaults: %d records", len(gen1File.GetEventsAndFaults()))
			t.Logf("  - DetailedSpeed: %d records", len(gen1File.GetDetailedSpeed()))
			t.Logf("  - TechnicalData: %d records", len(gen1File.GetTechnicalData()))

		case ddv1.Generation_GENERATION_2:
			version := vuFile.GetVersion()
			if version == ddv1.Version_VERSION_2 {
				gen2v2File := vuFile.GetGen2V2()
				if gen2v2File == nil {
					t.Fatal("expected Gen2 V2 file data, got nil")
				}
				t.Logf("  - Version: V2")

				// Check Overview
				if overview := gen2v2File.GetOverview(); overview != nil {
					t.Logf("  - Overview present, raw_data size: %d bytes", len(overview.GetRawData()))
					if len(overview.GetRawData()) == 0 {
						t.Errorf("    - Overview raw_data is empty")
					}
				}

				t.Logf("  - Activities: %d records", len(gen2v2File.GetActivities()))
				t.Logf("  - EventsAndFaults: %d records", len(gen2v2File.GetEventsAndFaults()))
				t.Logf("  - DetailedSpeed: %d records", len(gen2v2File.GetDetailedSpeed()))
				t.Logf("  - TechnicalData: %d records", len(gen2v2File.GetTechnicalData()))
			} else {
				gen2v1File := vuFile.GetGen2V1()
				if gen2v1File == nil {
					t.Fatal("expected Gen2 V1 file data, got nil")
				}
				t.Logf("  - Version: V1")

				// Check Overview
				if overview := gen2v1File.GetOverview(); overview != nil {
					t.Logf("  - Overview present, raw_data size: %d bytes", len(overview.GetRawData()))
					if len(overview.GetRawData()) == 0 {
						t.Errorf("    - Overview raw_data is empty")
					}
				}

				t.Logf("  - Activities: %d records", len(gen2v1File.GetActivities()))
				t.Logf("  - EventsAndFaults: %d records", len(gen2v1File.GetEventsAndFaults()))
				t.Logf("  - DetailedSpeed: %d records", len(gen2v1File.GetDetailedSpeed()))
				t.Logf("  - TechnicalData: %d records", len(gen2v1File.GetTechnicalData()))
			}

		default:
			t.Errorf("unexpected generation: %s", generation)
		}
		})
	}
}

// TestVehicleUnitFileGolden validates semantic VU file parsing against golden JSON files.
//
// Run with -update flag to regenerate golden files.
func TestVehicleUnitFileGolden(t *testing.T) {
	testFiles, err := filepath.Glob("../../testdata/vu/*.DDD")
	if err != nil {
		t.Fatalf("failed to glob test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("no VU test files found")
	}

	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			data, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			vuFile, err := unmarshalVehicleUnitFile(data)
			if err != nil {
				t.Fatalf("unmarshalVehicleUnitFile failed: %v", err)
			}

			goldenPath := filepath.Join("testdata", filepath.Base(testFile)+".golden.json")

			m := protojson.MarshalOptions{
				Indent:        "  ",
				UseProtoNames: true,
			}
			gotJSON, err := m.Marshal(vuFile)
			if err != nil {
				t.Fatalf("failed to marshal to JSON: %v", err)
			}

			if *update {
				// Update mode: write the golden file
				if err := os.MkdirAll("testdata", 0o755); err != nil {
					t.Fatalf("failed to create testdata directory: %v", err)
				}
				if err := os.WriteFile(goldenPath, gotJSON, 0o644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenPath)
				return
			}

			// Normal mode: compare with golden file
			wantJSON, err := os.ReadFile(goldenPath)
			if err != nil {
				if os.IsNotExist(err) {
					t.Skipf("Golden file does not exist: %s. Run with -update to create it.", goldenPath)
				}
				t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
			}

			if diff := cmp.Diff(string(wantJSON), string(gotJSON)); diff != "" {
				t.Errorf("Golden file mismatch (-want +got):\n%s", diff)
			}

			t.Logf("Successfully validated against golden file")
		})
	}
}
