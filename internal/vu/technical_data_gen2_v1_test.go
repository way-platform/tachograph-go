package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestAnonymizeTechnicalData_Gen2V1(t *testing.T) {
	// No real Gen2V1 VU hexdumps available; use synthetic construction.
	mfgName := &ddv1.StringValue{}
	mfgName.SetValue("Real Manufacturer GmbH")

	sn := &ddv1.ExtendedSerialNumber{}
	sn.SetSerialNumber(99999)
	sn.SetManufacturerCode(0xAB)

	vuIdent := &vuv1.TechnicalDataGen2V1_VuIdentification{}
	vuIdent.SetManufacturerName(mfgName)
	vuIdent.SetSerialNumber(sn)

	sensorSn := &ddv1.ExtendedSerialNumber{}
	sensorSn.SetSerialNumber(12345)
	sensor := &vuv1.TechnicalDataGen2V1_PairedSensor{}
	sensor.SetSerialNumber(sensorSn)

	calPurpose := ddv1.CalibrationPurpose_ACTIVATION
	cal := &vuv1.TechnicalDataGen2V1_CalibrationRecord{}
	cal.SetPurpose(calPurpose)
	cal.SetAuthorisedSpeedKmh(90)

	td := &vuv1.TechnicalDataGen2V1{}
	td.SetVuIdentification(vuIdent)
	td.SetPairedSensors([]*vuv1.TechnicalDataGen2V1_PairedSensor{sensor})
	td.SetCalibrationRecords([]*vuv1.TechnicalDataGen2V1_CalibrationRecord{cal})

	anon := AnonymizeOptions{}.anonymizeTechnicalDataGen2V1(td)
	if anon == nil {
		t.Fatal("anonymize returned nil")
	}

	// ManufacturerName anonymized
	if anon.GetVuIdentification().GetManufacturerName().GetValue() == "Real Manufacturer GmbH" {
		t.Error("expected ManufacturerName to be anonymized")
	}

	// SerialNumber zeroed
	if diff := cmp.Diff(int64(0), anon.GetVuIdentification().GetSerialNumber().GetSerialNumber()); diff != "" {
		t.Errorf("VuIdentification.SerialNumber not zeroed (-want +got):\n%s", diff)
	}
	// ManufacturerCode preserved
	if diff := cmp.Diff(int32(0xAB), int32(anon.GetVuIdentification().GetSerialNumber().GetManufacturerCode())); diff != "" {
		t.Errorf("VuIdentification.ManufacturerCode changed (-want +got):\n%s", diff)
	}

	// Sensor serial number zeroed
	if diff := cmp.Diff(int64(0), anon.GetPairedSensors()[0].GetSerialNumber().GetSerialNumber()); diff != "" {
		t.Errorf("PairedSensors[0].SerialNumber not zeroed (-want +got):\n%s", diff)
	}

	// Calibration purpose and speed preserved
	if diff := cmp.Diff(calPurpose, anon.GetCalibrationRecords()[0].GetPurpose()); diff != "" {
		t.Errorf("CalibrationRecords[0].Purpose changed (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(int32(90), anon.GetCalibrationRecords()[0].GetAuthorisedSpeedKmh()); diff != "" {
		t.Errorf("CalibrationRecords[0].AuthorisedSpeedKmh changed (-want +got):\n%s", diff)
	}

	// Signature empty
	if len(anon.GetSignature()) != 0 {
		t.Error("expected Signature to be empty after anonymization")
	}
}

func TestTechnicalData_Gen2V1(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_TECHNICAL_DATA_GEN2_V1)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for TECHNICAL_DATA_GEN2_V1")
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
			technicalData, err := unmarshalTechnicalDataGen2V1(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if technicalData == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, technicalData, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalTechnicalDataGen2V1(technicalData)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
