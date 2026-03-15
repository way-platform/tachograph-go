package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestAnonymizeOverview_Gen2V2(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_OVERVIEW_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for OVERVIEW_GEN2_V2")
	}

	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			overview, err := unmarshalOverviewGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			anon := AnonymizeOptions{}.anonymizeOverviewGen2V2(overview)
			if anon == nil {
				t.Fatal("anonymize returned nil")
			}

			// Certificates cleared (set to empty bytes)
			if len(anon.GetMemberStateCertificate()) != 0 {
				t.Error("expected MemberStateCertificate to be empty after anonymization")
			}
			if len(anon.GetVuCertificate()) != 0 {
				t.Error("expected VuCertificate to be empty after anonymization")
			}

			// VIN and VRN are anonymized (all asterisks)
			if vinVal := anon.GetVehicleIdentificationNumber().GetValue(); vinVal != "" && strings.Trim(vinVal, "*") != "" {
				t.Errorf("expected VehicleIdentificationNumber to be all asterisks, got %q", vinVal)
			}
			if vrnVal := anon.GetVehicleRegistrationNumber().GetValue(); vrnVal != "" && strings.Trim(vrnVal, "*") != "" {
				t.Errorf("expected VehicleRegistrationNumber to be all asterisks, got %q", vrnVal)
			}

			// CurrentDateTime preserved
			if diff := cmp.Diff(
				overview.GetCurrentDateTime().GetSeconds(),
				anon.GetCurrentDateTime().GetSeconds(),
			); diff != "" {
				t.Errorf("CurrentDateTime changed (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOverview_Gen2V2(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_OVERVIEW_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for OVERVIEW_GEN2_V2")
	}

	// Run subtest for each discovered file
	for _, hexdumpPath := range hexdumpFiles {
		// Use relative path from testdata as subtest name
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".data.hexdump")

		t.Run(testName, func(t *testing.T) {
			// Read hexdump
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			// Unmarshal
			overview, err := unmarshalOverviewGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if overview == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, overview, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalOverviewGen2V2(overview)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
