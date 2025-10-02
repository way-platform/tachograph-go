package card

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// TestVehicleUnitsUsedRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestVehicleUnitsUsedRoundTrip(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/vehicle_units_used.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	opts := UnmarshalOptions{}
	vu1, err := opts.unmarshalVehicleUnitsUsed(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	marshaled, err := appendCardVehicleUnitsUsed(nil, vu1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	vu2, err := opts.unmarshalVehicleUnitsUsed(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	if diff := cmp.Diff(vu1, vu2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestVehicleUnitsUsedAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestVehicleUnitsUsedAnonymization -update -v  # regenerate
func TestVehicleUnitsUsedAnonymization(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/vehicle_units_used.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	opts := UnmarshalOptions{}
	vu, err := opts.unmarshalVehicleUnitsUsed(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	anonymized := AnonymizeVehicleUnitsUsed(vu)

	anonymizedData, err := appendCardVehicleUnitsUsed(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/vehicle_units_used.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write vehicle_units_used.b64: %v", err)
		}

		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/vehicle_units_used.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write vehicle_units_used.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		expectedB64, err := os.ReadFile("testdata/vehicle_units_used.b64")
		if err != nil {
			t.Fatalf("Failed to read expected vehicle_units_used.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		expectedJSON, err := os.ReadFile("testdata/vehicle_units_used.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.VehicleUnitsUsed
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	if anonymized == nil {
		t.Fatal("Anonymized VehicleUnitsUsed is nil")
	}

	// Verify that all records have been anonymized
	records := anonymized.GetRecords()
	for i, record := range records {
		if record == nil {
			continue
		}

		// Check timestamp is in the anonymization sequence
		ts := record.GetTimestamp()
		if ts != nil && ts.Seconds != 0 {
			// Verify timestamp is 2020-01-01 00:00:00 UTC + (i * 1 hour)
			expectedSeconds := int64(1577836800 + i*3600)
			if ts.Seconds != expectedSeconds {
				t.Errorf("Record %d: timestamp = %d, want %d", i, ts.Seconds, expectedSeconds)
			}
		}

		// Manufacturer code should be 0x40 (test value) for non-zero records
		mc := record.GetManufacturerCode()
		if mc != 0x40 && mc != 0 {
			t.Errorf("Record %d: manufacturer_code = 0x%02x, want 0x40 or 0x00", i, mc)
		}

		// Software version should be "0000" or zeros for non-zero records
		swVersion := record.GetVuSoftwareVersion()
		if len(swVersion) == 4 && !bytes.Equal(swVersion, []byte("0000")) && !bytes.Equal(swVersion, []byte{0, 0, 0, 0}) {
			t.Errorf("Record %d: vu_software_version = %q, want \"0000\" or zeros", i, swVersion)
		}
	}
}

// AnonymizeVehicleUnitsUsed creates an anonymized copy of VehicleUnitsUsed,
// replacing sensitive data with static, deterministic test values.
func AnonymizeVehicleUnitsUsed(vu *cardv1.VehicleUnitsUsed) *cardv1.VehicleUnitsUsed {
	if vu == nil {
		return nil
	}

	result := &cardv1.VehicleUnitsUsed{}

	// Preserve the pointer to newest record
	result.SetVehicleUnitPointerNewestRecord(vu.GetVehicleUnitPointerNewestRecord())

	// Anonymize each record
	originalRecords := vu.GetRecords()
	anonymizedRecords := make([]*cardv1.VehicleUnitsUsed_Record, len(originalRecords))
	for i, record := range originalRecords {
		anonymizedRecords[i] = anonymizeVehicleUnitRecord(record, i)
	}
	result.SetRecords(anonymizedRecords)

	// Preserve signature if present
	if vu.HasSignature() {
		result.SetSignature(vu.GetSignature())
	}

	return result
}

// anonymizeVehicleUnitRecord anonymizes a single vehicle unit record.
// Uses index to create sequential timestamps.
func anonymizeVehicleUnitRecord(record *cardv1.VehicleUnitsUsed_Record, index int) *cardv1.VehicleUnitsUsed_Record {
	if record == nil {
		// Return a zero-filled record for nil entries
		zeroRecord := &cardv1.VehicleUnitsUsed_Record{}
		zeroRecord.SetManufacturerCode(0)
		zeroRecord.SetDeviceId([]byte{0})
		zeroRecord.SetVuSoftwareVersion([]byte{0, 0, 0, 0})
		return zeroRecord
	}

	result := &cardv1.VehicleUnitsUsed_Record{}

	// Replace timestamp with sequential test timestamps
	// Base: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	// Increment by 1 hour per record
	baseEpoch := int64(1577836800)
	ts := record.GetTimestamp()
	if ts != nil && ts.Seconds != 0 {
		// Non-zero timestamp - replace with sequential test timestamp
		result.SetTimestamp(&timestamppb.Timestamp{
			Seconds: baseEpoch + int64(index)*3600,
			Nanos:   0,
		})
	}
	// else: Zero timestamp - leave unset (nil)

	// Replace manufacturer code with test value 0x40 for non-zero values
	if record.GetManufacturerCode() != 0 {
		result.SetManufacturerCode(0x40)
	} else {
		result.SetManufacturerCode(0)
	}

	// Replace device ID with test value 0x00
	if len(record.GetDeviceId()) > 0 {
		result.SetDeviceId([]byte{0x00})
	}

	// Replace VU software version with "0000" for non-zero values
	swVersion := record.GetVuSoftwareVersion()
	if len(swVersion) > 0 && !bytes.Equal(swVersion, []byte{0, 0, 0, 0}) {
		result.SetVuSoftwareVersion([]byte("0000"))
	} else {
		result.SetVuSoftwareVersion([]byte{0, 0, 0, 0})
	}

	return result
}
