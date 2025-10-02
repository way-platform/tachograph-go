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
)

// TestGnssPlacesRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestGnssPlacesRoundTrip(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/gnss_places.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	opts := UnmarshalOptions{}
	gnssPlaces1, err := opts.unmarshalGnssPlaces(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	marshaled, err := appendCardGnssPlaces(nil, gnssPlaces1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	gnssPlaces2, err := opts.unmarshalGnssPlaces(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	if diff := cmp.Diff(gnssPlaces1, gnssPlaces2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestGnssPlacesAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestGnssPlacesAnonymization -update -v  # regenerate
func TestGnssPlacesAnonymization(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/gnss_places.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	opts := UnmarshalOptions{}
	gnssPlaces, err := opts.unmarshalGnssPlaces(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	anonymized := AnonymizeGnssPlaces(gnssPlaces)

	anonymizedData, err := appendCardGnssPlaces(nil, anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/gnss_places.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write gnss_places.b64: %v", err)
		}

		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/gnss_places.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write gnss_places.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		expectedB64, err := os.ReadFile("testdata/gnss_places.b64")
		if err != nil {
			t.Fatalf("Failed to read expected gnss_places.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		expectedJSON, err := os.ReadFile("testdata/gnss_places.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.GnssPlaces
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	if anonymized == nil {
		t.Fatal("Anonymized GnssPlaces is nil")
	}

	// Verify that all records have been anonymized
	records := anonymized.GetRecords()
	for i, record := range records {
		if record == nil {
			continue
		}

		// Check outer timestamp is in the anonymization sequence
		ts := record.GetTimestamp()
		if ts != nil && ts.Seconds != 0 {
			// Verify timestamp is 2020-01-01 00:00:00 UTC + (i * 1 hour)
			expectedSeconds := int64(1577836800 + i*3600)
			if ts.Seconds != expectedSeconds {
				t.Errorf("Record %d: outer timestamp = %d, want %d", i, ts.Seconds, expectedSeconds)
			}
		}

		// Check GNSS place record has Helsinki coordinates
		gnssPlace := record.GetGnssPlaceRecord()
		if gnssPlace != nil && gnssPlace.GetGeoCoordinates() != nil {
			coords := gnssPlace.GetGeoCoordinates()
			// Helsinki: 60°10.0'N (60100), 24°56.0'E (24560)
			if coords.GetLatitude() != 60100 || coords.GetLongitude() != 24560 {
				t.Errorf("Record %d: coordinates = (%d, %d), want (60100, 24560)",
					i, coords.GetLatitude(), coords.GetLongitude())
			}
		}

		// Verify odometer is rounded to nearest 100km
		odometer := record.GetVehicleOdometerKm()
		if odometer%100 != 0 {
			t.Errorf("Record %d: odometer = %d, should be rounded to nearest 100km", i, odometer)
		}
	}
}

// AnonymizeGnssPlaces creates an anonymized copy of GnssPlaces,
// replacing sensitive GNSS data with static, deterministic test values.
func AnonymizeGnssPlaces(gnssPlaces *cardv1.GnssPlaces) *cardv1.GnssPlaces {
	if gnssPlaces == nil {
		return nil
	}

	result := &cardv1.GnssPlaces{}

	// Preserve the pointer to newest record
	result.SetNewestRecordIndex(gnssPlaces.GetNewestRecordIndex())

	// Anonymize each record
	originalRecords := gnssPlaces.GetRecords()
	anonymizedRecords := make([]*cardv1.GnssPlaces_Record, len(originalRecords))
	for i, record := range originalRecords {
		anonymizedRecords[i] = anonymizeGNSSAccumulatedDrivingRecord(record, i)
	}
	result.SetRecords(anonymizedRecords)

	// Preserve signature if present
	if gnssPlaces.HasSignature() {
		result.SetSignature(gnssPlaces.GetSignature())
	}

	return result
}

// anonymizeGNSSAccumulatedDrivingRecord anonymizes a single GNSS accumulated driving record.
// Uses index to create sequential timestamps.
func anonymizeGNSSAccumulatedDrivingRecord(record *cardv1.GnssPlaces_Record, index int) *cardv1.GnssPlaces_Record {
	if record == nil {
		// Return a zero-filled record for nil entries
		zeroRecord := &cardv1.GnssPlaces_Record{}
		zeroRecord.SetVehicleOdometerKm(0)
		return zeroRecord
	}

	result := &cardv1.GnssPlaces_Record{}

	// Replace outer timestamp with sequential test timestamps
	// Base: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	// Increment by 1 hour per record
	baseEpoch := int64(1577836800)
	ts := record.GetTimestamp()
	if ts != nil && ts.Seconds != 0 {
		// Non-zero timestamp - replace with sequential test timestamp
		result.SetTimestamp(&timestamppb.Timestamp{
			Seconds: baseEpoch + int64(index)*3600,
			Nanos:   0,
		})
	}
	// else: Zero timestamp - leave unset (nil)

	// Anonymize GNSS place record (replaces coordinates with Helsinki)
	gnssPlaceRecord := record.GetGnssPlaceRecord()
	if gnssPlaceRecord != nil {
		anonymizedGnssPlace := dd.AnonymizeGNSSPlaceRecord(gnssPlaceRecord)
		// Update the inner timestamp to match the outer one (for consistency)
		if result.GetTimestamp() != nil {
			anonymizedGnssPlace.SetTimestamp(result.GetTimestamp())
		}
		result.SetGnssPlaceRecord(anonymizedGnssPlace)
	}

	// Round odometer to nearest 100km (preserves magnitude but not exact correlation)
	originalOdometer := record.GetVehicleOdometerKm()
	roundedOdometer := (originalOdometer / 100) * 100
	result.SetVehicleOdometerKm(roundedOdometer)

	return result
}
