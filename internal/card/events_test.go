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

// TestEventsDataRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestEventsDataRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/events.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	opts := UnmarshalOptions{}
	events1, err := opts.unmarshalEventsData(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := appendEventsData(nil, events1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	events2, err := opts.unmarshalEventsData(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(events1, events2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestEventsDataAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestEventsDataAnonymization -update -v  # regenerate
func TestEventsDataAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/events.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	opts := UnmarshalOptions{}
	events, err := opts.unmarshalEventsData(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymized := AnonymizeEventsData(events)

	// Marshal anonymized data
	anonymizedData, err := appendEventsData(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/events.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write events.b64: %v", err)
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
		if err := os.WriteFile("testdata/events.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write events.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/events.b64")
		if err != nil {
			t.Fatalf("Failed to read expected events.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/events.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.EventsData
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized EventsData is nil")
	}

	// Verify event count is preserved
	if len(anonymized.GetEvents()) != len(events.GetEvents()) {
		t.Errorf("Event count changed: got %d, want %d", len(anonymized.GetEvents()), len(events.GetEvents()))
	}

	// Verify valid/invalid status is preserved for all events
	for i, event := range anonymized.GetEvents() {
		origEvent := events.GetEvents()[i]
		if event.GetValid() != origEvent.GetValid() {
			t.Errorf("Event %d valid status changed: got %v, want %v", i, event.GetValid(), origEvent.GetValid())
		}

		// For valid events, check vehicle registration is anonymized to FINLAND
		if event.GetValid() {
			vehicleReg := event.GetEventVehicleRegistration()
			if vehicleReg != nil && vehicleReg.GetNation() != ddv1.NationNumeric_FINLAND {
				t.Errorf("Event %d vehicle nation = %v, want FINLAND", i, vehicleReg.GetNation())
			}
		}
	}
}

// AnonymizeEventsData creates an anonymized copy of EventsData,
// replacing sensitive information with static, deterministic test values.
func AnonymizeEventsData(events *cardv1.EventsData) *cardv1.EventsData {
	if events == nil {
		return nil
	}

	anonymized := &cardv1.EventsData{}

	// Base timestamp for anonymization: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	baseEpoch := int64(1577836800)

	var anonymizedEvents []*cardv1.EventsData_Record
	for i, event := range events.GetEvents() {
		anonymizedEvent := &cardv1.EventsData_Record{}

		// Preserve valid flag
		anonymizedEvent.SetValid(event.GetValid())

		if event.GetValid() {
			// Preserve event type (not sensitive, categorical)
			anonymizedEvent.SetEventType(event.GetEventType())

			// Use incrementing timestamps based on index (1 hour apart)
			beginTime := &timestamppb.Timestamp{Seconds: baseEpoch + int64(i)*3600}
			endTime := &timestamppb.Timestamp{Seconds: baseEpoch + int64(i)*3600 + 1800} // 30 mins later
			anonymizedEvent.SetEventBeginTime(beginTime)
			anonymizedEvent.SetEventEndTime(endTime)

			// Anonymize vehicle registration
			if vehicleReg := event.GetEventVehicleRegistration(); vehicleReg != nil {
				anonymizedReg := &ddv1.VehicleRegistrationIdentification{}
				anonymizedReg.SetNation(ddv1.NationNumeric_FINLAND)

				// Use static test registration number
				// VehicleRegistrationNumber is: 1 byte code page + 13 bytes data
				testRegNum := &ddv1.StringValue{}
				testRegNum.SetValue("TEST-VRN")
				testRegNum.SetEncoding(ddv1.Encoding_ISO_8859_1) // Code page 1 (Latin-1)
				testRegNum.SetLength(13)                               // Length of data bytes (not including code page)
				anonymizedReg.SetNumber(testRegNum)

				anonymizedEvent.SetEventVehicleRegistration(anonymizedReg)
			}

			// Regenerate raw_data for binary fidelity
			rawData, err := appendEventRecord(nil, anonymizedEvent)
			if err == nil {
				anonymizedEvent.SetRawData(rawData)
			}
		} else {
			// Preserve invalid records as-is
			anonymizedEvent.SetRawData(event.GetRawData())
		}

		anonymizedEvents = append(anonymizedEvents, anonymizedEvent)
	}

	anonymized.SetEvents(anonymizedEvents)
	return anonymized
}
