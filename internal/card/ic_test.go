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

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// TestIcRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestIcRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/ic.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	opts := UnmarshalOptions{}
	ic1, err := opts.unmarshalIc(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := appendCardIc(nil, ic1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	ic2, err := opts.unmarshalIc(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(ic1, ic2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestIcAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestIcAnonymization -update -v  # regenerate
func TestIcAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/ic.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	opts := UnmarshalOptions{}
	ic, err := opts.unmarshalIc(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeIc(ic)

	// Marshal anonymized data
	anonymizedData, err := appendCardIc(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/ic.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write ic.b64: %v", err)
		}

		// Write golden JSON with stable formatting
		// First convert to JSON using protojson
		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		// Then reformat with json.Indent for stable, deterministic output
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/ic.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write ic.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/ic.b64")
		if err != nil {
			t.Fatalf("Failed to read expected ic.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/ic.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.Ic
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized IC is nil")
	}

	// IC contains mostly hardware identifiers - anonymized values should be static test data
	if len(anonymized.GetIcSerialNumber()) != 4 {
		t.Errorf("IC serial number length = %d, want 4", len(anonymized.GetIcSerialNumber()))
	}
	if len(anonymized.GetIcManufacturingReferences()) != 4 {
		t.Errorf("IC manufacturing references length = %d, want 4", len(anonymized.GetIcManufacturingReferences()))
	}
}

// AnonymizeIc creates an anonymized copy of IC data, replacing hardware identifiers
// with static, deterministic test values.
func AnonymizeIc(ic *cardv1.Ic) *cardv1.Ic {
	if ic == nil {
		return nil
	}

	anonymized := &cardv1.Ic{}

	// Replace IC serial number with static test value
	// IC serial number is a hardware identifier - use static placeholder
	anonymized.SetIcSerialNumber([]byte{0x00, 0x00, 0x00, 0x01})

	// Replace manufacturing references with static test value
	// Manufacturing references are hardware identifiers - use static placeholder
	anonymized.SetIcManufacturingReferences([]byte{0xAA, 0xBB, 0xCC, 0xDD})

	return anonymized
}
