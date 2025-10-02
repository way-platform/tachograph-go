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
)

// TestVehiclesRoundTrip verifies that VehiclesUsed can be marshalled and unmarshalled
// with perfect binary fidelity (unmarshal → marshal → unmarshal produces identical results).
func TestVehiclesRoundTrip(t *testing.T) {
	// Read the base64-encoded test data
	b64Data, err := os.ReadFile("testdata/vehicles.b64")
	if err != nil {
		t.Fatalf("Failed to read testdata/vehicles.b64: %v", err)
	}

	// Decode from base64
	originalBytes, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	vehicles1, err := (UnmarshalOptions{}).unmarshalVehiclesUsed(originalBytes)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}
	if vehicles1 == nil {
		t.Fatal("First unmarshal returned nil")
	}

	// Marshal back to binary
	marshalledBytes, err := appendVehiclesUsed(nil, vehicles1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if len(marshalledBytes) != len(originalBytes) {
		t.Errorf("Binary round-trip failed: length mismatch: got %d, want %d", len(marshalledBytes), len(originalBytes))
	}
	if diff := cmp.Diff(originalBytes, marshalledBytes); diff != "" {
		t.Errorf("Binary round-trip failed: original and marshalled bytes differ (-want +got):\n%s", diff)
		t.Fatalf("Original length: %d, Marshalled length: %d", len(originalBytes), len(marshalledBytes))
	}

	// Second unmarshal (from marshalled bytes)
	vehicles2, err := (UnmarshalOptions{}).unmarshalVehiclesUsed(marshalledBytes)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}
	if vehicles2 == nil {
		t.Fatal("Second unmarshal returned nil")
	}

	// Verify structural equality
	if diff := cmp.Diff(vehicles1, vehicles2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural round-trip failed: parsed messages differ (-want +got):\n%s", diff)
	}
}

// TestVehiclesAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestVehiclesAnonymization -update -v  # regenerate
func TestVehiclesAnonymization(t *testing.T) {
	// Read the original base64-encoded test data
	b64Data, err := os.ReadFile("testdata/vehicles.b64")
	if err != nil {
		t.Fatalf("Failed to read testdata/vehicles.b64: %v", err)
	}

	// Decode from base64
	originalBytes, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal the original data
	vehicles, err := (UnmarshalOptions{}).unmarshalVehiclesUsed(originalBytes)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeVehiclesUsed(vehicles)

	// Marshal the anonymized data
	anonymizedBytes, err := appendVehiclesUsed(nil, anonymized)
	if err != nil {
		t.Fatalf("Marshal anonymized failed: %v", err)
	}

	// Encode to base64
	anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedBytes)

	// Convert to JSON with stable formatting
	// First convert to JSON using protojson
	jsonBytes, err := protojson.Marshal(anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
	}
	// Then reformat with json.Indent for stable output
	var stableJSON bytes.Buffer
	if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}

	if *update {
		// Update mode: write new golden files
		if err := os.WriteFile("testdata/vehicles.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write testdata/vehicles.b64: %v", err)
		}
		t.Logf("✅ Updated: testdata/vehicles.b64 (%d bytes)", len(anonymizedBytes))

		if err := os.WriteFile("testdata/vehicles.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write testdata/vehicles.golden.json: %v", err)
		}
		t.Logf("✅ Updated: testdata/vehicles.golden.json")
	} else {
		// Assert mode: verify matches golden files

		// Re-anonymizing should produce identical output (determinism check)
		reanonymizedBytes, err := appendVehiclesUsed(nil, AnonymizeVehiclesUsed(vehicles))
		if err != nil {
			t.Fatalf("Re-anonymize marshal failed: %v", err)
		}
		if diff := cmp.Diff(anonymizedBytes, reanonymizedBytes); diff != "" {
			t.Errorf("Re-anonymizing produced different output (-want +got):\n%s", diff)
		}

		// Verify binary matches testdata/vehicles.b64
		goldenB64, err := os.ReadFile("testdata/vehicles.b64")
		if err != nil {
			t.Fatalf("Failed to read testdata/vehicles.b64: %v", err)
		}
		if diff := cmp.Diff(string(goldenB64), anonymizedB64); diff != "" {
			t.Errorf("Binary mismatch with testdata/vehicles.b64 (-want +got):\n%s", diff)
		}

		// Verify JSON matches testdata/vehicles.golden.json
		goldenJSON, err := os.ReadFile("testdata/vehicles.golden.json")
		if err != nil {
			t.Fatalf("Failed to read testdata/vehicles.golden.json: %v", err)
		}
		if diff := cmp.Diff(string(goldenJSON), stableJSON.String()); diff != "" {
			t.Errorf("Golden JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized.GetNewestRecordIndex() != vehicles.GetNewestRecordIndex() {
		t.Errorf("Anonymization changed newest record index: got %d, want %d",
			anonymized.GetNewestRecordIndex(), vehicles.GetNewestRecordIndex())
	}

	if len(anonymized.GetRecords()) != len(vehicles.GetRecords()) {
		t.Errorf("Anonymization changed record count: got %d, want %d",
			len(anonymized.GetRecords()), len(vehicles.GetRecords()))
	}

	// Verify vehicle registrations are anonymized
	for i, record := range anonymized.GetRecords() {
		vrn := record.GetVehicleRegistration().GetNumber().GetValue()
		if vrn != "TEST-VRN" {
			t.Errorf("Record %d: vehicle registration not anonymized: got %q, want %q", i, vrn, "TEST-VRN")
		}

		// Verify country is preserved (structural info)
		// We don't test specific country value as it depends on the original data
	}
}
