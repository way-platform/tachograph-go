package card

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestCardDownload_Generation1(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER,
		ddv1.Generation_GENERATION_1,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Fatal("No hexdump files found for EF_CARD_DOWNLOAD_DRIVER GENERATION_1")
	}

	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			opts := UnmarshalOptions{}
			result, err := opts.unmarshalCardDownload(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, result, goldenPath)

			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalCardDownload(result)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCardDownload_Generation2(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Fatal("No hexdump files found for EF_CARD_DOWNLOAD_DRIVER GENERATION_2")
	}

	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".hexdump")

		t.Run(testName, func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			opts := UnmarshalOptions{}
			result, err := opts.unmarshalCardDownload(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, result, goldenPath)

			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalCardDownload(result)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
