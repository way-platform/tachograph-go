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
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestApplicationIdentificationRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestApplicationIdentificationRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/application_identification.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	opts := UnmarshalOptions{}
	appId1, err := opts.unmarshalApplicationIdentification(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := appendCardApplicationIdentification(nil, appId1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	appId2, err := opts.unmarshalApplicationIdentification(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(appId1, appId2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestApplicationIdentificationAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestApplicationIdentificationAnonymization -update -v  # regenerate
func TestApplicationIdentificationAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/application_identification.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	opts := UnmarshalOptions{}
	appId, err := opts.unmarshalApplicationIdentification(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeApplicationIdentification(appId)

	// Marshal anonymized data
	anonymizedData, err := appendCardApplicationIdentification(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/application_identification.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write application_identification.b64: %v", err)
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
		if err := os.WriteFile("testdata/application_identification.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write application_identification.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/application_identification.b64")
		if err != nil {
			t.Fatalf("Failed to read expected application_identification.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/application_identification.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.ApplicationIdentification
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized ApplicationIdentification is nil")
	}

	// Verify card type is set
	if anonymized.GetCardType() != cardv1.CardType_DRIVER_CARD {
		t.Errorf("Card type = %v, want DRIVER_CARD", anonymized.GetCardType())
	}

	// Verify type of tachograph card ID is set
	if anonymized.GetTypeOfTachographCardId() == ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED {
		t.Error("Type of tachograph card ID should be set")
	}

	// Verify card structure version is set
	if anonymized.GetCardStructureVersion() == nil {
		t.Error("Card structure version should not be nil")
	}

	// Verify driver data is present
	driver := anonymized.GetDriver()
	if driver == nil {
		t.Fatal("Driver data should not be nil for driver card")
	}
}

// AnonymizeApplicationIdentification creates an anonymized copy of ApplicationIdentification data,
// preserving structural information while using static test values for version information.
func AnonymizeApplicationIdentification(appId *cardv1.ApplicationIdentification) *cardv1.ApplicationIdentification {
	if appId == nil {
		return nil
	}

	anonymized := &cardv1.ApplicationIdentification{}

	// Preserve type of tachograph card ID (not sensitive, categorical)
	anonymized.SetTypeOfTachographCardId(appId.GetTypeOfTachographCardId())

	// Preserve card structure version (not sensitive, technical metadata)
	anonymized.SetCardStructureVersion(appId.GetCardStructureVersion())

	// Preserve card type (not sensitive, categorical)
	anonymized.SetCardType(appId.GetCardType())

	// Preserve driver data structure (counts are not sensitive)
	if driver := appId.GetDriver(); driver != nil {
		anonymizedDriver := &cardv1.ApplicationIdentification_Driver{}
		anonymizedDriver.SetEventsPerTypeCount(driver.GetEventsPerTypeCount())
		anonymizedDriver.SetFaultsPerTypeCount(driver.GetFaultsPerTypeCount())
		anonymizedDriver.SetActivityStructureLength(driver.GetActivityStructureLength())
		anonymizedDriver.SetCardVehicleRecordsCount(driver.GetCardVehicleRecordsCount())
		anonymizedDriver.SetCardPlaceRecordsCount(driver.GetCardPlaceRecordsCount())
		anonymized.SetDriver(anonymizedDriver)
	}

	return anonymized
}

