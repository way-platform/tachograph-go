package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestAnonymizeEventsAndFaults_Gen2V2(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for EVENTS_AND_FAULTS_GEN2_V2")
	}

	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			ef, err := unmarshalEventsAndFaultsGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			anon := AnonymizeOptions{}.anonymizeEventsAndFaultsGen2V2(ef)
			if anon == nil {
				t.Fatal("anonymize returned nil")
			}

			// Faults: FaultType and timestamps preserved, card numbers changed
			for i, origFault := range ef.GetFaults() {
				anonFault := anon.GetFaults()[i]
				if diff := cmp.Diff(origFault.GetFaultType(), anonFault.GetFaultType()); diff != "" {
					t.Errorf("Faults[%d].FaultType changed (-want +got):\n%s", i, diff)
				}
				if diff := cmp.Diff(origFault.GetBeginTime().GetSeconds(), anonFault.GetBeginTime().GetSeconds()); diff != "" {
					t.Errorf("Faults[%d].BeginTime changed (-want +got):\n%s", i, diff)
				}
				if diff := cmp.Diff(origFault.GetEndTime().GetSeconds(), anonFault.GetEndTime().GetSeconds()); diff != "" {
					t.Errorf("Faults[%d].EndTime changed (-want +got):\n%s", i, diff)
				}
			}

			// Events: EventType and timestamps preserved
			for i, origEvent := range ef.GetEvents() {
				anonEvent := anon.GetEvents()[i]
				if diff := cmp.Diff(origEvent.GetEventType(), anonEvent.GetEventType()); diff != "" {
					t.Errorf("Events[%d].EventType changed (-want +got):\n%s", i, diff)
				}
				if diff := cmp.Diff(origEvent.GetBeginTime().GetSeconds(), anonEvent.GetBeginTime().GetSeconds()); diff != "" {
					t.Errorf("Events[%d].BeginTime changed (-want +got):\n%s", i, diff)
				}
			}
		})
	}
}

func TestEventsAndFaults_Gen2V2(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for EVENTS_AND_FAULTS_GEN2_V2")
	}

	// Run subtest for each discovered file
	for _, hexdumpPath := range hexdumpFiles {
		// Use relative path from testdata as subtest name
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".data.hexdump")

		t.Run(testName, func(t *testing.T) {
			// Read hexdump
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			// Unmarshal
			eventsAndFaults, err := unmarshalEventsAndFaultsGen2V2(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if eventsAndFaults == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, eventsAndFaults, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalEventsAndFaultsGen2V2(eventsAndFaults)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
