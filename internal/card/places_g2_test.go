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

func TestUnmarshalPlacesG2(t *testing.T) {
	input, err := os.ReadFile("testdata/places_g2.b64")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	inputBytes, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		t.Fatalf("Failed to decode input: %v", err)
	}
	places, err := (UnmarshalOptions{}).unmarshalPlacesG2(inputBytes)
	if err != nil {
		t.Fatalf("UnmarshalPlacesG2 failed: %v", err)
	}
	// Convert to JSON with stable formatting
	// First convert to JSON using protojson
	jsonBytes, err := protojson.Marshal(places)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
	}
	// Then reformat with json.Indent for stable output
	var stableJSON bytes.Buffer
	if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}
	if err := os.WriteFile("testdata/places_g2.golden.json", stableJSON.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write golden JSON: %v", err)
	}
}

// TestPlacesG2RoundTrip verifies that PlacesG2 can be marshalled and unmarshalled
// with perfect binary fidelity (unmarshal → marshal → unmarshal produces identical results).
func TestPlacesG2RoundTrip(t *testing.T) {
	// Read test data
	input, err := os.ReadFile("testdata/places_g2.b64")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	originalBytes, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		t.Fatalf("Failed to decode input: %v", err)
	}

	t.Logf("Original data: %d bytes", len(originalBytes))

	// First unmarshal
	places1, err := (UnmarshalOptions{}).unmarshalPlacesG2(originalBytes)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	t.Logf("First unmarshal: newestRecordIndex=%d, numRecords=%d",
		places1.GetNewestRecordIndex(), len(places1.GetRecords()))

	// Marshal
	marshalledBytes, err := appendPlacesG2(nil, places1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	t.Logf("Marshalled data: %d bytes", len(marshalledBytes))

	// Verify binary fidelity
	if !bytes.Equal(originalBytes, marshalledBytes) {
		t.Errorf("Binary round-trip failed: lengths differ")
		t.Logf("Original length: %d", len(originalBytes))
		t.Logf("Marshalled length: %d", len(marshalledBytes))

		// Find first difference
		minLen := len(originalBytes)
		if len(marshalledBytes) < minLen {
			minLen = len(marshalledBytes)
		}
		for i := 0; i < minLen; i++ {
			if originalBytes[i] != marshalledBytes[i] {
				t.Logf("First difference at byte %d: original=0x%02x, marshalled=0x%02x",
					i, originalBytes[i], marshalledBytes[i])
				break
			}
		}
	}

	// Second unmarshal to verify structural fidelity
	places2, err := (UnmarshalOptions{}).unmarshalPlacesG2(marshalledBytes)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	t.Logf("Second unmarshal: newestRecordIndex=%d, numRecords=%d",
		places2.GetNewestRecordIndex(), len(places2.GetRecords()))

	// Compare structures
	if diff := cmp.Diff(places1, places2, protocmp.Transform()); diff != "" {
		t.Errorf("Structure round-trip failed (-first +second):\n%s", diff)
	}

	// Additional structural checks
	if places1.GetNewestRecordIndex() != places2.GetNewestRecordIndex() {
		t.Errorf("newestRecordIndex mismatch: first=%d, second=%d",
			places1.GetNewestRecordIndex(), places2.GetNewestRecordIndex())
	}

	if len(places1.GetRecords()) != len(places2.GetRecords()) {
		t.Errorf("record count mismatch: first=%d, second=%d",
			len(places1.GetRecords()), len(places2.GetRecords()))
	}
}

// TestPlacesG2Anonymization verifies that anonymization is deterministic and stable.
// By default, this test asserts that re-anonymizing the test data produces
// identical results (no changes). When run with -update flag, it regenerates
// the anonymized test data:
//
//	go test -run TestPlacesG2Anonymization -update -v
//
// Since anonymization is deterministic, this acts as a golden file test that
// catches unintended changes while allowing intentional updates.
func TestPlacesG2Anonymization(t *testing.T) {
	// Read current test data
	input, err := os.ReadFile("testdata/places_g2.b64")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	currentBytes, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		t.Fatalf("Failed to decode input: %v", err)
	}

	// Unmarshal
	places, err := (UnmarshalOptions{}).unmarshalPlacesG2(currentBytes)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Re-anonymize (should be idempotent since data is already anonymized)
	anonymized := AnonymizePlacesG2(places)

	// Marshal back
	anonymizedBytes, err := appendPlacesG2(nil, anonymized)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify round-trip works
	places2, err := (UnmarshalOptions{}).unmarshalPlacesG2(anonymizedBytes)
	if err != nil {
		t.Fatalf("Round-trip unmarshal failed: %v", err)
	}

	// Generate golden JSON with stable formatting
	// First convert to JSON using protojson
	jsonBytes, err := protojson.Marshal(places2)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
	}
	// Then reformat with json.Indent for stable output
	var stableJSON bytes.Buffer
	if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}
	jsonData := stableJSON.Bytes()

	if *update {
		// Update mode: write new files
		anonymizedBase64 := base64.StdEncoding.EncodeToString(anonymizedBytes)
		if err := os.WriteFile("testdata/places_g2.b64", []byte(anonymizedBase64), 0o644); err != nil {
			t.Fatalf("Failed to write places_g2.b64: %v", err)
		}
		t.Logf("✅ Updated: testdata/places_g2.b64 (%d bytes)", len(anonymizedBytes))

		if err := os.WriteFile("testdata/places_g2.golden.json", jsonData, 0o644); err != nil {
			t.Fatalf("Failed to write golden JSON: %v", err)
		}
		t.Logf("✅ Updated: testdata/places_g2.golden.json")
	} else {
		// Assert mode: verify files haven't changed
		if !bytes.Equal(currentBytes, anonymizedBytes) {
			t.Errorf("Re-anonymizing places_g2.b64 produced different output.\n" +
				"This means anonymization is not deterministic or the test data is stale.\n" +
				"Run 'go test -update' to regenerate the golden files.")
		}

		currentJSON, err := os.ReadFile("testdata/places_g2.golden.json")
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
	const testEpoch = int64(1577836800) // 2020-01-01 00:00:00 UTC
	var earliestSeconds int64 = -1
	for i, record := range places2.GetRecords() {
		if record.GetDailyWorkPeriodCountry() != ddv1.NationNumeric_FINLAND {
			t.Errorf("Record %d has wrong country: got %v, want FINLAND", i, record.GetDailyWorkPeriodCountry())
		}
		if gnss := record.GetEntryGnssPlaceRecord(); gnss != nil {
			if gnss.GetGeoCoordinates().GetLatitude() != 60100 || gnss.GetGeoCoordinates().GetLongitude() != 24560 {
				t.Errorf("Record %d has wrong GNSS coordinates: got lat=%d, lon=%d, want lat=60100, lon=24560",
					i, gnss.GetGeoCoordinates().GetLatitude(), gnss.GetGeoCoordinates().GetLongitude())
			}
		}
		if ts := record.GetEntryTime(); ts != nil {
			if earliestSeconds == -1 || ts.Seconds < earliestSeconds {
				earliestSeconds = ts.Seconds
			}
		}
		if gnss := record.GetEntryGnssPlaceRecord(); gnss != nil {
			if ts := gnss.GetTimestamp(); ts != nil {
				if earliestSeconds == -1 || ts.Seconds < earliestSeconds {
					earliestSeconds = ts.Seconds
				}
			}
		}
	}
	if earliestSeconds != testEpoch {
		t.Errorf("Earliest timestamp should be test epoch: got %d, want %d", earliestSeconds, testEpoch)
	}
}
