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

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestDrivingLicenceInfoRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestDrivingLicenceInfoRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/driving_licence.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	opts := UnmarshalOptions{}
	dli1, err := opts.unmarshalDrivingLicenceInfo(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := appendDrivingLicenceInfo(nil, dli1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	dli2, err := opts.unmarshalDrivingLicenceInfo(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(dli1, dli2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestDrivingLicenceInfoAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestDrivingLicenceInfoAnonymization -update -v  # regenerate
func TestDrivingLicenceInfoAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/driving_licence.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	opts := UnmarshalOptions{}
	dli, err := opts.unmarshalDrivingLicenceInfo(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeDrivingLicenceInfo(dli)

	// Marshal anonymized data
	anonymizedData, err := appendDrivingLicenceInfo(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/driving_licence.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write driving_licence.b64: %v", err)
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
		if err := os.WriteFile("testdata/driving_licence.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write driving_licence.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/driving_licence.b64")
		if err != nil {
			t.Fatalf("Failed to read expected driving_licence.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/driving_licence.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.DrivingLicenceInfo
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized DrivingLicenceInfo is nil")
	}

	// Verify issuing authority is anonymized
	authority := anonymized.GetDrivingLicenceIssuingAuthority()
	if authority == nil || authority.GetValue() == "" {
		t.Error("Driving licence issuing authority should not be empty")
	}

	// Verify issuing nation is FINLAND
	if anonymized.GetDrivingLicenceIssuingNation() != ddv1.NationNumeric_FINLAND {
		t.Errorf("Issuing nation = %v, want FINLAND", anonymized.GetDrivingLicenceIssuingNation())
	}

	// Verify licence number is anonymized
	licenceNumber := anonymized.GetDrivingLicenceNumber()
	if licenceNumber == nil || licenceNumber.GetValue() == "" {
		t.Error("Driving licence number should not be empty")
	}
}

// AnonymizeDrivingLicenceInfo creates an anonymized copy of DrivingLicenceInfo data,
// replacing sensitive information with static, deterministic test values.
func AnonymizeDrivingLicenceInfo(dli *cardv1.DrivingLicenceInfo) *cardv1.DrivingLicenceInfo {
	if dli == nil {
		return nil
	}

	anonymized := &cardv1.DrivingLicenceInfo{}

	// Replace issuing authority with static test value
	if authority := dli.GetDrivingLicenceIssuingAuthority(); authority != nil {
		anonymized.SetDrivingLicenceIssuingAuthority(dd.AnonymizeStringValue(authority, "TEST AUTHORITY"))
	}

	// Country → FINLAND (always)
	anonymized.SetDrivingLicenceIssuingNation(ddv1.NationNumeric_FINLAND)

	// Replace licence number with static test value
	if licenceNumber := dli.GetDrivingLicenceNumber(); licenceNumber != nil {
		anonymized.SetDrivingLicenceNumber(dd.AnonymizeStringValue(licenceNumber, "TEST-DL-123"))
	}

	return anonymized
}
