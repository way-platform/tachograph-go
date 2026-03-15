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

func TestAnonymizeOverview_Gen2V1(t *testing.T) {
	// No real Gen2V1 VU hexdumps available; use synthetic construction.
	vin := &ddv1.Ia5StringValue{}
	vin.SetValue("WDB9634031L123456")

	overview := &vuv1.OverviewGen2V1{}
	overview.SetVehicleIdentificationNumber(vin)
	ts := timestamppb.New(time.Date(2025, 9, 12, 10, 0, 0, 0, time.UTC))
	overview.SetCurrentDateTime(ts)
	period := &ddv1.DownloadablePeriod{}
	period.SetMinTime(timestamppb.New(time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)))
	period.SetMaxTime(timestamppb.New(time.Date(2025, 9, 12, 0, 0, 0, 0, time.UTC)))
	overview.SetDownloadablePeriod(period)

	anon := AnonymizeOptions{}.anonymizeOverviewGen2V1(overview)
	if anon == nil {
		t.Fatal("anonymize returned nil")
	}

	// Certificates cleared (set to empty bytes)
	if len(anon.GetMemberStateCertificate()) != 0 {
		t.Error("expected MemberStateCertificate to be empty after anonymization")
	}
	if len(anon.GetVuCertificate()) != 0 {
		t.Error("expected VuCertificate to be empty after anonymization")
	}

	// VIN anonymized (all asterisks)
	if vinVal := anon.GetVehicleIdentificationNumber().GetValue(); strings.Trim(vinVal, "*") != "" {
		t.Errorf("expected VehicleIdentificationNumber to be all asterisks, got %q", vinVal)
	}

	// CurrentDateTime and DownloadablePeriod preserved
	if diff := cmp.Diff(overview.GetCurrentDateTime().GetSeconds(), anon.GetCurrentDateTime().GetSeconds()); diff != "" {
		t.Errorf("CurrentDateTime changed (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(
		overview.GetDownloadablePeriod().GetMinTime().GetSeconds(),
		anon.GetDownloadablePeriod().GetMinTime().GetSeconds(),
	); diff != "" {
		t.Errorf("DownloadablePeriod.MinTime changed (-want +got):\n%s", diff)
	}
}

func TestOverview_Gen2V1(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_OVERVIEW_GEN2_V1)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for OVERVIEW_GEN2_V1")
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
			overview, err := unmarshalOverviewGen2V1(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if overview == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, overview, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalOverviewGen2V1(overview)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
