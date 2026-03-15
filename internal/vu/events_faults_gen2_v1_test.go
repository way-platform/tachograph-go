package vu

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestAnonymizeEventsAndFaults_Gen2V1(t *testing.T) {
	// No real Gen2V1 VU hexdumps available; use synthetic construction.
	beginTime := timestamppb.New(time.Date(2025, 9, 10, 8, 0, 0, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2025, 9, 10, 9, 0, 0, 0, time.UTC))

	driverIdNum := &ddv1.Ia5StringValue{}
	driverIdNum.SetValue("FIN12345678")
	driverID := &ddv1.DriverIdentification{}
	driverID.SetDriverIdentificationNumber(driverIdNum)
	innerCard := &ddv1.FullCardNumber{}
	innerCard.SetDriverIdentification(driverID)
	cardNum := &ddv1.FullCardNumberAndGeneration{}
	cardNum.SetFullCardNumber(innerCard)

	fault := &vuv1.EventsAndFaultsGen2V1_FaultRecord{}
	fault.SetFaultType(ddv1.EventFaultType_VU_SEC_MOTION_SENSOR_AUTH_FAILURE)
	fault.SetBeginTime(beginTime)
	fault.SetEndTime(endTime)
	fault.SetCardNumberAndGenDriverSlotBegin(cardNum)

	event := &vuv1.EventsAndFaultsGen2V1_EventRecord{}
	event.SetEventType(ddv1.EventFaultType_GENERAL_POWER_SUPPLY_INTERRUPTION)
	event.SetBeginTime(beginTime)
	event.SetEndTime(endTime)
	event.SetCardNumberAndGenDriverSlotBegin(cardNum)

	ef := &vuv1.EventsAndFaultsGen2V1{}
	ef.SetFaults([]*vuv1.EventsAndFaultsGen2V1_FaultRecord{fault})
	ef.SetEvents([]*vuv1.EventsAndFaultsGen2V1_EventRecord{event})

	anon := AnonymizeOptions{}.anonymizeEventsAndFaultsGen2V1(ef)
	if anon == nil {
		t.Fatal("anonymize returned nil")
	}

	// FaultType and timestamps preserved
	if anon.GetFaults()[0].GetFaultType() != ddv1.EventFaultType_VU_SEC_MOTION_SENSOR_AUTH_FAILURE {
		t.Error("expected FaultType to be preserved")
	}
	if diff := cmp.Diff(beginTime.GetSeconds(), anon.GetFaults()[0].GetBeginTime().GetSeconds()); diff != "" {
		t.Errorf("Faults[0].BeginTime changed (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(endTime.GetSeconds(), anon.GetFaults()[0].GetEndTime().GetSeconds()); diff != "" {
		t.Errorf("Faults[0].EndTime changed (-want +got):\n%s", diff)
	}

	// Card DriverIdentificationNumber anonymized (asterisks, not original)
	origDriverId := "FIN12345678"
	anonDriverId := anon.GetFaults()[0].GetCardNumberAndGenDriverSlotBegin().GetFullCardNumber().GetDriverIdentification().GetDriverIdentificationNumber().GetValue()
	if anonDriverId == origDriverId {
		t.Error("expected CardNumberAndGenDriverSlotBegin DriverIdentificationNumber to be anonymized")
	}

	// EventType preserved
	if anon.GetEvents()[0].GetEventType() != ddv1.EventFaultType_GENERAL_POWER_SUPPLY_INTERRUPTION {
		t.Error("expected EventType to be preserved")
	}

	// Signature empty
	if len(anon.GetSignature()) != 0 {
		t.Error("expected Signature to be empty after anonymization")
	}
}

func TestEventsAndFaults_Gen2V1(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for EVENTS_AND_FAULTS_GEN2_V1")
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
			eventsAndFaults, err := unmarshalEventsAndFaultsGen2V1(data)
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
			marshaled, err := marshalOpts.MarshalEventsAndFaultsGen2V1(eventsAndFaults)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
