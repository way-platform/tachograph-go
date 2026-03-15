package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestAnonymizeActivities_Gen2V2(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_ACTIVITIES_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for ACTIVITIES_GEN2_V2")
	}

	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			activities, err := unmarshalActivitiesGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			anon := AnonymizeOptions{}.anonymizeActivitiesGen2V2(activities)
			if anon == nil {
				t.Fatal("anonymize returned nil")
			}

			// CardIW holder names replaced with TEST/DRIVER
			for i, rec := range anon.GetCardIwData() {
				if got := rec.GetCardHolderName().GetHolderSurname().GetValue(); got != "TEST" {
					t.Errorf("CardIwData[%d].HolderSurname = %q, want TEST", i, got)
				}
				if got := rec.GetCardHolderName().GetHolderFirstNames().GetValue(); got != "DRIVER" {
					t.Errorf("CardIwData[%d].HolderFirstNames = %q, want DRIVER", i, got)
				}
			}

			// Signature empty, RawData empty
			if len(anon.GetSignature()) != 0 {
				t.Error("expected Signature to be empty after anonymization")
			}
			if len(anon.GetRawData()) != 0 {
				t.Error("expected RawData to be empty after anonymization")
			}
		})
	}
}

func TestActivities_Gen2V2(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_ACTIVITIES_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for ACTIVITIES_GEN2_V2")
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
			activities, err := unmarshalActivitiesGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if activities == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, activities, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalActivitiesGen2V2(activities)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
