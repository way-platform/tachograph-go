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
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestFaultsDataRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestFaultsDataRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/faults.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	opts := UnmarshalOptions{}
	faults1, err := opts.unmarshalFaultsData(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := appendFaultsData(nil, faults1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	faults2, err := opts.unmarshalFaultsData(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(faults1, faults2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestFaultsDataAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestFaultsDataAnonymization -update -v  # regenerate
func TestFaultsDataAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/faults.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	opts := UnmarshalOptions{}
	faults, err := opts.unmarshalFaultsData(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeFaultsData(faults)

	// Marshal anonymized data
	anonymizedData, err := appendFaultsData(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/faults.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write faults.b64: %v", err)
		}

		// Write golden JSON with stable formatting
		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/faults.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write faults.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/faults.b64")
		if err != nil {
			t.Fatalf("Failed to read expected faults.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/faults.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.FaultsData
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized FaultsData is nil")
	}

	// Verify fault count is preserved
	if len(anonymized.GetFaults()) != len(faults.GetFaults()) {
		t.Errorf("Fault count changed: got %d, want %d", len(anonymized.GetFaults()), len(faults.GetFaults()))
	}

	// Verify valid/invalid status is preserved for all faults
	for i, fault := range anonymized.GetFaults() {
		origFault := faults.GetFaults()[i]
		if fault.GetValid() != origFault.GetValid() {
			t.Errorf("Fault %d valid status changed: got %v, want %v", i, fault.GetValid(), origFault.GetValid())
		}

		// For valid faults, check vehicle registration is anonymized to FINLAND
		if fault.GetValid() {
			vehicleReg := fault.GetFaultVehicleRegistration()
			if vehicleReg != nil && vehicleReg.GetNation() != ddv1.NationNumeric_FINLAND {
				t.Errorf("Fault %d vehicle nation = %v, want FINLAND", i, vehicleReg.GetNation())
			}
		}
	}
}

// AnonymizeFaultsData creates an anonymized copy of FaultsData,
// replacing sensitive information with static, deterministic test values.
func AnonymizeFaultsData(faults *cardv1.FaultsData) *cardv1.FaultsData {
	if faults == nil {
		return nil
	}

	anonymized := &cardv1.FaultsData{}

	// Base timestamp for anonymization: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	baseEpoch := int64(1577836800)

	var anonymizedFaults []*cardv1.FaultsData_Record
	for i, fault := range faults.GetFaults() {
		anonymizedFault := &cardv1.FaultsData_Record{}

		// Preserve valid flag
		anonymizedFault.SetValid(fault.GetValid())

		if fault.GetValid() {
			// Preserve fault type (not sensitive, categorical)
			anonymizedFault.SetFaultType(fault.GetFaultType())

			// Use incrementing timestamps based on index (1 hour apart)
			beginTime := &timestamppb.Timestamp{Seconds: baseEpoch + int64(i)*3600}
			endTime := &timestamppb.Timestamp{Seconds: baseEpoch + int64(i)*3600 + 1800} // 30 mins later
			anonymizedFault.SetFaultBeginTime(beginTime)
			anonymizedFault.SetFaultEndTime(endTime)

			// Anonymize vehicle registration
			if vehicleReg := fault.GetFaultVehicleRegistration(); vehicleReg != nil {
				anonymizedReg := &ddv1.VehicleRegistrationIdentification{}
				anonymizedReg.SetNation(ddv1.NationNumeric_FINLAND)

				// Use static test registration number
				// VehicleRegistrationNumber is: 1 byte code page + 13 bytes data
				testRegNum := &ddv1.StringValue{}
				testRegNum.SetValue("TEST-VRN")
				testRegNum.SetEncoding(ddv1.Encoding_ISO_8859_1) // Code page 1 (Latin-1)
				testRegNum.SetLength(13)                               // Length of data bytes (not including code page)
				anonymizedReg.SetNumber(testRegNum)

				anonymizedFault.SetFaultVehicleRegistration(anonymizedReg)
			}

			// Regenerate raw_data for binary fidelity
			rawData, err := appendFaultRecord(nil, anonymizedFault)
			if err == nil {
				anonymizedFault.SetRawData(rawData)
			}
		} else {
			// Preserve invalid records as-is
			anonymizedFault.SetRawData(fault.GetRawData())
		}

		anonymizedFaults = append(anonymizedFaults, anonymizedFault)
	}

	anonymized.SetFaults(anonymizedFaults)
	return anonymized
}
