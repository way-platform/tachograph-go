package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestAnonymizeTechnicalData_Gen2V2(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_TECHNICAL_DATA_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for TECHNICAL_DATA_GEN2_V2")
	}

	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			td, err := unmarshalTechnicalDataGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			anon := AnonymizeOptions{}.anonymizeTechnicalDataGen2V2(td)
			if anon == nil {
				t.Fatal("anonymize returned nil")
			}

			// ManufacturerName anonymized (all asterisks)
			if mfgName := anon.GetVuIdentification().GetManufacturerName().GetValue(); mfgName != "" && strings.Trim(mfgName, "*") != "" {
				t.Errorf("expected ManufacturerName to be all asterisks, got %q", mfgName)
			}

			// VU serial number zeroed
			if anon.GetVuIdentification().GetSerialNumber() != nil {
				if diff := cmp.Diff(int64(0), anon.GetVuIdentification().GetSerialNumber().GetSerialNumber()); diff != "" {
					t.Errorf("VuIdentification.SerialNumber not zeroed (-want +got):\n%s", diff)
				}
			}

			// Sensor serial numbers zeroed
			for i, sensor := range anon.GetPairedSensors() {
				if sensor.GetSerialNumber() != nil {
					if diff := cmp.Diff(int64(0), sensor.GetSerialNumber().GetSerialNumber()); diff != "" {
						t.Errorf("PairedSensors[%d].SerialNumber not zeroed (-want +got):\n%s", i, diff)
					}
				}
			}

			// CoupledGnssFacilities serial numbers zeroed
			for i, gnss := range anon.GetCoupledGnssFacilities() {
				if gnss.GetSerialNumber() != nil {
					if diff := cmp.Diff(int64(0), gnss.GetSerialNumber().GetSerialNumber()); diff != "" {
						t.Errorf("CoupledGnssFacilities[%d].SerialNumber not zeroed (-want +got):\n%s", i, diff)
					}
				}
			}

			// Signature empty
			if len(anon.GetSignature()) != 0 {
				t.Error("expected Signature to be empty after anonymization")
			}
		})
	}
}

func TestTechnicalData_Gen2V2(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_TECHNICAL_DATA_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for TECHNICAL_DATA_GEN2_V2")
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
			technicalData, err := unmarshalTechnicalDataGen2V2(data)
			if err != nil {
				t.Fatalf("Failed to unmarshal TechnicalData Gen2V2: %v", err)
			}
			if technicalData == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, technicalData, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalTechnicalDataGen2V2(technicalData)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
