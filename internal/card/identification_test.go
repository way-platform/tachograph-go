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

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestIdentificationRoundTrip verifies that Identification can be marshalled and unmarshalled
// with perfect binary fidelity (unmarshal → marshal → unmarshal produces identical results).
func TestIdentificationRoundTrip(t *testing.T) {
	input, err := os.ReadFile("testdata/identification.b64")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	originalBytes, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		t.Fatalf("Failed to decode input: %v", err)
	}

	t.Logf("Original data: %d bytes", len(originalBytes))

	// First unmarshal
	identification1, err := (UnmarshalOptions{}).unmarshalIdentification(originalBytes)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}
	cardNumber := ""
	if card := identification1.GetCard(); card != nil {
		if driverID := card.GetDriverIdentification(); driverID != nil {
			cardNumber = driverID.GetDriverIdentificationNumber().GetValue()
		}
	}
	holderSurname := ""
	if holder := identification1.GetDriverCardHolder(); holder != nil {
		holderSurname = holder.GetCardHolderSurname().GetValue()
	}
	t.Logf("First unmarshal: card_number=%s, holder_surname=%s", cardNumber, holderSurname)

	// Marshal both parts
	marshalledBytes, err := appendCardIdentification(nil, identification1.GetCard())
	if err != nil {
		t.Fatalf("Marshal Card failed: %v", err)
	}
	marshalledBytes, err = appendDriverCardHolderIdentification(marshalledBytes, identification1.GetDriverCardHolder())
	if err != nil {
		t.Fatalf("Marshal DriverCardHolder failed: %v", err)
	}
	t.Logf("Marshalled data: %d bytes", len(marshalledBytes))

	// Verify binary equality
	if !bytes.Equal(originalBytes, marshalledBytes) {
		t.Errorf("Binary round-trip failed: original and marshalled bytes differ")
		t.Logf("Original length: %d, Marshalled length: %d", len(originalBytes), len(marshalledBytes))
	}

	// Second unmarshal
	identification2, err := (UnmarshalOptions{}).unmarshalIdentification(marshalledBytes)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}
	cardNumber2 := ""
	if card := identification2.GetCard(); card != nil {
		if driverID := card.GetDriverIdentification(); driverID != nil {
			cardNumber2 = driverID.GetDriverIdentificationNumber().GetValue()
		}
	}
	holderSurname2 := ""
	if holder := identification2.GetDriverCardHolder(); holder != nil {
		holderSurname2 = holder.GetCardHolderSurname().GetValue()
	}
	t.Logf("Second unmarshal: card_number=%s, holder_surname=%s", cardNumber2, holderSurname2)

	// Verify structural equality
	if diff := cmp.Diff(identification1, identification2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural round-trip failed (-first +second):\n%s", diff)
	}
}

// TestIdentificationAnonymization verifies that anonymization is deterministic and stable.
// By default, this test asserts that re-anonymizing the test data produces
// identical results (no changes). When run with -update flag, it regenerates
// the anonymized test data:
//
//	go test -run TestIdentificationAnonymization -update -v
//
// Since anonymization is deterministic, this acts as a golden file test that
// catches unintended changes while allowing intentional updates.
func TestIdentificationAnonymization(t *testing.T) {
	// Read current test data
	input, err := os.ReadFile("testdata/identification.b64")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	currentBytes, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		t.Fatalf("Failed to decode input: %v", err)
	}

	// Unmarshal
	identification, err := (UnmarshalOptions{}).unmarshalIdentification(currentBytes)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Re-anonymize (should be idempotent since data is already anonymized)
	anonymized := AnonymizeIdentification(identification)

	// Marshal both parts
	anonymizedBytes, err := appendCardIdentification(nil, anonymized.GetCard())
	if err != nil {
		t.Fatalf("Marshal Card failed: %v", err)
	}
	anonymizedBytes, err = appendDriverCardHolderIdentification(anonymizedBytes, anonymized.GetDriverCardHolder())
	if err != nil {
		t.Fatalf("Marshal DriverCardHolder failed: %v", err)
	}

	// Verify round-trip works
	identification2, err := (UnmarshalOptions{}).unmarshalIdentification(anonymizedBytes)
	if err != nil {
		t.Fatalf("Round-trip unmarshal failed: %v", err)
	}

	// Generate golden JSON
	// Convert to JSON with stable formatting
	jsonBytes, err := protojson.Marshal(identification2)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
	}
	// Then reformat with json.Indent for stable output
	var stableJSON bytes.Buffer
	if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}
	jsonData := stableJSON.Bytes()
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	if *update {
		// Update mode: write new files
		anonymizedBase64 := base64.StdEncoding.EncodeToString(anonymizedBytes)
		if err := os.WriteFile("testdata/identification.b64", []byte(anonymizedBase64), 0o644); err != nil {
			t.Fatalf("Failed to write identification.b64: %v", err)
		}
		t.Logf("✅ Updated: testdata/identification.b64 (%d bytes)", len(anonymizedBytes))

		if err := os.WriteFile("testdata/identification.golden.json", jsonData, 0o644); err != nil {
			t.Fatalf("Failed to write golden JSON: %v", err)
		}
		t.Logf("✅ Updated: testdata/identification.golden.json")
	} else {
		// Assert mode: verify files haven't changed
		if !bytes.Equal(currentBytes, anonymizedBytes) {
			t.Errorf("Re-anonymizing identification.b64 produced different output.\n" +
				"This means anonymization is not deterministic or the test data is stale.\n" +
				"Run 'go test -update' to regenerate the golden files.")
		}

		currentJSON, err := os.ReadFile("testdata/identification.golden.json")
		if err != nil {
			t.Fatalf("Failed to read golden JSON: %v", err)
		}
		if !bytes.Equal(currentJSON, jsonData) {
			t.Errorf("Golden JSON mismatch.\n"+
				"Run 'go test -update' to regenerate the golden files.\n"+
				"Diff:\n%s", cmp.Diff(string(currentJSON), string(jsonData)))
		}
	}

	// Additional structural assertions (always run)
	card := identification2.GetCard()
	if card == nil {
		t.Fatal("Card is nil")
	}

	// Verify anonymized card number pattern (DRIVER00000001)
	driverID := card.GetDriverIdentification()
	if driverID == nil {
		t.Fatal("DriverIdentification is nil")
	}
	cardNumberStr := driverID.GetDriverIdentificationNumber().GetValue()
	if len(cardNumberStr) < 6 || cardNumberStr[0:6] != "DRIVER" {
		t.Errorf("Card number should start with 'DRIVER': got %s", cardNumberStr)
	}

	// Verify anonymized names
	holder := identification2.GetDriverCardHolder()
	if holder == nil {
		t.Fatal("DriverCardHolder is nil")
	}
	if holder.GetCardHolderSurname().GetValue() != "TEST_SURNAME" {
		t.Errorf("Expected anonymized surname 'TEST_SURNAME', got %s", holder.GetCardHolderSurname().GetValue())
	}
	if holder.GetCardHolderFirstNames().GetValue() != "TEST_FIRSTNAME" {
		t.Errorf("Expected anonymized first names 'TEST_FIRSTNAME', got %s", holder.GetCardHolderFirstNames().GetValue())
	}

	// Verify country is Finland (our test default)
	if card.GetCardIssuingMemberState() != ddv1.NationNumeric_FINLAND {
		t.Errorf("Expected issuing country FINLAND, got %v", card.GetCardIssuingMemberState())
	}
}
