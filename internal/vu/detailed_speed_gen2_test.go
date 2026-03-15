package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestAnonymizeDetailedSpeed_Gen2(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_DETAILED_SPEED_GEN2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for DETAILED_SPEED_GEN2")
	}

	// Use just first file to keep test fast (files are large)
	hexdumpPath := hexdumpFiles[0]
	relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
	testName := strings.TrimSuffix(relPath, ".hexdump")

	t.Run(testName, func(t *testing.T) {
		data, err := readHexdump(hexdumpPath)
		if err != nil {
			t.Fatalf("Failed to read hexdump: %v", err)
		}

		ds, err := unmarshalDetailedSpeedGen2(data)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		anon := AnonymizeOptions{}.anonymizeDetailedSpeedGen2(ds)
		if anon == nil {
			t.Fatal("anonymize returned nil")
		}

		// SpeedBlocks count preserved
		if diff := cmp.Diff(len(ds.GetSpeedBlocks()), len(anon.GetSpeedBlocks())); diff != "" {
			t.Errorf("SpeedBlocks count changed (-want +got):\n%s", diff)
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

func TestDetailedSpeed_Gen2(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_DETAILED_SPEED_GEN2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for DETAILED_SPEED_GEN2")
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
			detailedSpeed, err := unmarshalDetailedSpeedGen2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if detailedSpeed == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, detailedSpeed, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalDetailedSpeedGen2(detailedSpeed)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
