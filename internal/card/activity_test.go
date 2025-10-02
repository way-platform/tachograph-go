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

// TestDriverActivityDataRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestDriverActivityDataRoundTrip(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/activity.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	opts := UnmarshalOptions{}
	activity1, err := opts.unmarshalDriverActivityData(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	marshaled, err := appendDriverActivity(nil, activity1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	activity2, err := opts.unmarshalDriverActivityData(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	if diff := cmp.Diff(activity1, activity2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestDriverActivityDataAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestDriverActivityDataAnonymization -update -v  # regenerate
func TestDriverActivityDataAnonymization(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/activity.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	opts := UnmarshalOptions{}
	activity, err := opts.unmarshalDriverActivityData(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	anonymized := AnonymizeDriverActivityData(activity)

	anonymizedData, err := appendDriverActivity(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/activity.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write activity.b64: %v", err)
		}

		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/activity.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write activity.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		expectedB64, err := os.ReadFile("testdata/activity.b64")
		if err != nil {
			t.Fatalf("Failed to read expected activity.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		expectedJSON, err := os.ReadFile("testdata/activity.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.DriverActivityData
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	if anonymized == nil {
		t.Fatal("Anonymized DriverActivityData is nil")
	}

	// Verify record count is preserved
	if len(anonymized.GetDailyRecords()) != len(activity.GetDailyRecords()) {
		t.Errorf("Daily record count changed: got %d, want %d",
			len(anonymized.GetDailyRecords()), len(activity.GetDailyRecords()))
	}
}

// AnonymizeDriverActivityData creates an anonymized copy of DriverActivityData.
// Due to the complexity of the cyclic buffer structure with byte offsets and pointers,
// we preserve the raw data as-is but anonymize dates in the parsed records.
// Note: This approach preserves binary fidelity while providing anonymized semantic data.
func AnonymizeDriverActivityData(activity *cardv1.DriverActivityData) *cardv1.DriverActivityData {
	if activity == nil {
		return nil
	}

	anonymized := &cardv1.DriverActivityData{}

	// Preserve cyclic buffer structure (pointers and raw data)
	anonymized.SetOldestDayRecordIndex(activity.GetOldestDayRecordIndex())
	anonymized.SetNewestDayRecordIndex(activity.GetNewestDayRecordIndex())
	anonymized.SetRawData(activity.GetRawData())

	// Base timestamp for anonymization: 2020-01-01 00:00:00 UTC
	baseEpoch := int64(1577836800)

	// Anonymize only the parsed record dates
	var anonymizedRecords []*cardv1.DriverActivityData_DailyRecord
	for i, record := range activity.GetDailyRecords() {
		anonymizedRecord := &cardv1.DriverActivityData_DailyRecord{}
		anonymizedRecord.SetValid(record.GetValid())
		anonymizedRecord.SetRawData(record.GetRawData())

		if record.GetValid() {
			// Anonymize only the date field
			recordDate := &timestamppb.Timestamp{Seconds: baseEpoch + int64(i)*86400}
			anonymizedRecord.SetActivityRecordDate(recordDate)

			// Preserve all other fields
			anonymizedRecord.SetActivityPreviousRecordLength(record.GetActivityPreviousRecordLength())
			anonymizedRecord.SetActivityRecordLength(record.GetActivityRecordLength())
			anonymizedRecord.SetActivityDailyPresenceCounter(record.GetActivityDailyPresenceCounter())
			anonymizedRecord.SetActivityDayDistance(record.GetActivityDayDistance())
			if changes := record.GetActivityChangeInfo(); changes != nil {
				anonymizedRecord.SetActivityChangeInfo(changes)
			}
		}

		anonymizedRecords = append(anonymizedRecords, anonymizedRecord)
	}

	anonymized.SetDailyRecords(anonymizedRecords)

	return anonymized
}
