package card

import (
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalGnssPlacesAuthentication(t *testing.T) {
	t.Run("empty input error", func(t *testing.T) {
		_, err := UnmarshalOptions{}.unmarshalGnssPlacesAuthentication([]byte{})
		if err == nil {
			t.Fatal("expected error for empty input, got nil")
		}
	})

	t.Run("pointer only, zero records", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		got, err := UnmarshalOptions{}.unmarshalGnssPlacesAuthentication(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetNewestRecordIndex() != 0 {
			t.Errorf("NewestRecordIndex: got %d, want 0", got.GetNewestRecordIndex())
		}
		if len(got.GetRecords()) != 0 {
			t.Errorf("Records count: got %d, want 0", len(got.GetRecords()))
		}
	})

	t.Run("one record authenticated", func(t *testing.T) {
		// makePlaceAuthStatusRecord reused: 4B timestamp + 1B authStatus
		rec := makePlaceAuthStatusRecord(9000, 0x01) // AUTHENTICATED
		data := append([]byte{0x00, 0x00}, rec...)

		got, err := UnmarshalOptions{}.unmarshalGnssPlacesAuthentication(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.GetRecords()) != 1 {
			t.Fatalf("Records count: got %d, want 1", len(got.GetRecords()))
		}
		r := got.GetRecords()[0]
		if r.GetTimestamp().GetSeconds() != 9000 {
			t.Errorf("Timestamp: got %d, want 9000", r.GetTimestamp().GetSeconds())
		}
		if r.GetAuthenticationStatus() != ddv1.PositionAuthenticationStatus_AUTHENTICATED {
			t.Errorf("AuthenticationStatus: got %v, want AUTHENTICATED", r.GetAuthenticationStatus())
		}
	})

	t.Run("multiple records", func(t *testing.T) {
		rec1 := makePlaceAuthStatusRecord(1000, 0x01) // AUTHENTICATED
		rec2 := makePlaceAuthStatusRecord(2000, 0x00) // NOT_AUTHENTICATED
		data := append([]byte{0x00, 0x01}, rec1...)
		data = append(data, rec2...)

		got, err := UnmarshalOptions{}.unmarshalGnssPlacesAuthentication(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetNewestRecordIndex() != 1 {
			t.Errorf("NewestRecordIndex: got %d, want 1", got.GetNewestRecordIndex())
		}
		if len(got.GetRecords()) != 2 {
			t.Fatalf("Records count: got %d, want 2", len(got.GetRecords()))
		}
	})

	t.Run("truncated record returns error", func(t *testing.T) {
		// 2-byte pointer + 3 bytes (not a multiple of 5)
		data := make([]byte, 5)
		_, err := UnmarshalOptions{}.unmarshalGnssPlacesAuthentication(data)
		if err == nil {
			t.Fatal("expected error for truncated record, got nil")
		}
	})
}

func TestGnssPlacesAuthentication_ParseRawDriverCardFile(t *testing.T) {
	rec := makePlaceAuthStatusRecord(4800, 0x01) // AUTHENTICATED
	efData := append([]byte{0x00, 0x00}, rec...)

	rawRec := &cardv1.RawCardFile_Record{}
	rawRec.SetFile(cardv1.ElementaryFileType_EF_GNSS_PLACES_AUTHENTICATION)
	rawRec.SetGeneration(ddv1.Generation_GENERATION_2)
	rawRec.SetContentType(cardv1.ContentType_DATA)
	rawRec.SetValue(efData)

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords([]*cardv1.RawCardFile_Record{rawRec})

	cardFile, err := ParseOptions{}.ParseRawDriverCardFile(rawFile)
	if err != nil {
		t.Fatalf("ParseRawDriverCardFile: %v", err)
	}
	gpa := cardFile.GetTachographG2().GetGnssPlacesAuthentication()
	if gpa == nil {
		t.Fatal("GnssPlacesAuthentication is nil — EF_GNSS_PLACES_AUTHENTICATION not wired into switch")
	}
	if len(gpa.GetRecords()) != 1 {
		t.Fatalf("Records count: got %d, want 1", len(gpa.GetRecords()))
	}
	if gpa.GetRecords()[0].GetAuthenticationStatus() != ddv1.PositionAuthenticationStatus_AUTHENTICATED {
		t.Errorf("AuthenticationStatus: got %v, want AUTHENTICATED", gpa.GetRecords()[0].GetAuthenticationStatus())
	}
}

func TestGnssPlacesAuthentication_Fixtures(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_GNSS_PLACES_AUTHENTICATION,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("no hexdump fixtures found for EF_GNSS_PLACES_AUTHENTICATION GENERATION_2")
	}
	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		t.Run(strings.TrimSuffix(relPath, ".hexdump"), func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("failed to read hexdump: %v", err)
			}
			got, err := UnmarshalOptions{}.unmarshalGnssPlacesAuthentication(data)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			loadOrCreateGolden(t, got, goldenJSONPath(hexdumpPath))
		})
	}
}
