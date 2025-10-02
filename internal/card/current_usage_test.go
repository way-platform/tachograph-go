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

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestCurrentUsageRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestCurrentUsageRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/current_usage.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	opts := UnmarshalOptions{}
	cu1, err := opts.unmarshalCurrentUsage(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := appendCurrentUsage(nil, cu1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	cu2, err := opts.unmarshalCurrentUsage(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(cu1, cu2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestCurrentUsageAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestCurrentUsageAnonymization -update -v  # regenerate
func TestCurrentUsageAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/current_usage.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	opts := UnmarshalOptions{}
	cu, err := opts.unmarshalCurrentUsage(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeCurrentUsage(cu)

	// Marshal anonymized data
	anonymizedData, err := appendCurrentUsage(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/current_usage.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write current_usage.b64: %v", err)
		}

		// Write golden JSON with stable formatting
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
		if err := os.WriteFile("testdata/current_usage.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write current_usage.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/current_usage.b64")
		if err != nil {
			t.Fatalf("Failed to read expected current_usage.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/current_usage.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.CurrentUsage
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized CurrentUsage is nil")
	}

	// Verify session open time is set to static test value
	if anonymized.GetSessionOpenTime() == nil {
		t.Error("Session open time should not be nil")
	}

	// Verify vehicle registration is anonymized
	vehicleReg := anonymized.GetSessionOpenVehicle()
	if vehicleReg == nil {
		t.Error("Vehicle registration should not be nil")
	} else {
		// Should be FINLAND
		if vehicleReg.GetNation() != ddv1.NationNumeric_FINLAND {
			t.Errorf("Nation = %v, want FINLAND", vehicleReg.GetNation())
		}
		// Should have a registration number
		if vehicleReg.GetNumber() == nil || vehicleReg.GetNumber().GetValue() == "" {
			t.Error("Registration number should not be empty")
		}
	}
}

// AnonymizeCurrentUsage creates an anonymized copy of CurrentUsage data,
// replacing sensitive information with static, deterministic test values.
func AnonymizeCurrentUsage(cu *cardv1.CurrentUsage) *cardv1.CurrentUsage {
	if cu == nil {
		return nil
	}

	anonymized := &cardv1.CurrentUsage{}

	// Use static test timestamp: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	anonymized.SetSessionOpenTime(&timestamppb.Timestamp{Seconds: 1577836800})

	// Anonymize vehicle registration
	if vehicleReg := cu.GetSessionOpenVehicle(); vehicleReg != nil {
		anonymizedReg := &ddv1.VehicleRegistrationIdentification{}

		// Country → FINLAND (always)
		anonymizedReg.SetNation(ddv1.NationNumeric_FINLAND)

		// Registration number → static test value
		if regNum := vehicleReg.GetNumber(); regNum != nil {
			anonymizedReg.SetNumber(dd.AnonymizeStringValue(regNum, "TEST-123"))
		}

		anonymized.SetSessionOpenVehicle(anonymizedReg)
	}

	return anonymized
}
