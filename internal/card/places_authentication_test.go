package card

import (
	"encoding/binary"
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// makePlaceAuthStatusRecord builds a 5-byte PlaceAuthStatusRecord payload.
//
// Layout: 4B entryTime | 1B authenticationStatus
func makePlaceAuthStatusRecord(entryTimeSecs uint32, authStatus byte) []byte {
	r := make([]byte, 5)
	binary.BigEndian.PutUint32(r[0:4], entryTimeSecs)
	r[4] = authStatus
	return r
}

func TestUnmarshalPlacesAuthentication(t *testing.T) {
	t.Run("empty input error", func(t *testing.T) {
		_, err := UnmarshalOptions{}.unmarshalPlacesAuthentication([]byte{})
		if err == nil {
			t.Fatal("expected error for empty input, got nil")
		}
	})

	t.Run("pointer only, zero records", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		got, err := UnmarshalOptions{}.unmarshalPlacesAuthentication(data)
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
		// authStatus byte 0x01 → AUTHENTICATED
		rec := makePlaceAuthStatusRecord(7200, 0x01)
		data := append([]byte{0x00, 0x00}, rec...)

		got, err := UnmarshalOptions{}.unmarshalPlacesAuthentication(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.GetRecords()) != 1 {
			t.Fatalf("Records count: got %d, want 1", len(got.GetRecords()))
		}
		r := got.GetRecords()[0]
		if r.GetEntryTime().GetSeconds() != 7200 {
			t.Errorf("EntryTime: got %d, want 7200", r.GetEntryTime().GetSeconds())
		}
		if r.GetAuthenticationStatus() != ddv1.PositionAuthenticationStatus_AUTHENTICATED {
			t.Errorf("AuthenticationStatus: got %v, want AUTHENTICATED", r.GetAuthenticationStatus())
		}
	})

	t.Run("multiple records", func(t *testing.T) {
		rec1 := makePlaceAuthStatusRecord(1000, 0x01) // AUTHENTICATED
		rec2 := makePlaceAuthStatusRecord(2000, 0x00) // NOT_AUTHENTICATED
		rec3 := makePlaceAuthStatusRecord(3000, 0x00) // NOT_AUTHENTICATED
		data := append([]byte{0x00, 0x02}, rec1...)
		data = append(data, rec2...)
		data = append(data, rec3...)

		got, err := UnmarshalOptions{}.unmarshalPlacesAuthentication(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetNewestRecordIndex() != 2 {
			t.Errorf("NewestRecordIndex: got %d, want 2", got.GetNewestRecordIndex())
		}
		if len(got.GetRecords()) != 3 {
			t.Fatalf("Records count: got %d, want 3", len(got.GetRecords()))
		}
	})

	t.Run("truncated record returns error", func(t *testing.T) {
		// 2-byte pointer + 3 bytes (not a multiple of 5)
		data := make([]byte, 5)
		_, err := UnmarshalOptions{}.unmarshalPlacesAuthentication(data)
		if err == nil {
			t.Fatal("expected error for truncated record, got nil")
		}
	})
}

func TestPlacesAuthentication_ParseRawDriverCardFile(t *testing.T) {
	rec := makePlaceAuthStatusRecord(3600, 0x01) // AUTHENTICATED
	efData := append([]byte{0x00, 0x00}, rec...)

	rawRec := &cardv1.RawCardFile_Record{}
	rawRec.SetFile(cardv1.ElementaryFileType_EF_PLACES_AUTHENTICATION)
	rawRec.SetGeneration(ddv1.Generation_GENERATION_2)
	rawRec.SetContentType(cardv1.ContentType_DATA)
	rawRec.SetValue(efData)

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords([]*cardv1.RawCardFile_Record{rawRec})

	cardFile, err := ParseOptions{}.ParseRawDriverCardFile(rawFile)
	if err != nil {
		t.Fatalf("ParseRawDriverCardFile: %v", err)
	}
	pa := cardFile.GetTachographG2().GetPlacesAuthentication()
	if pa == nil {
		t.Fatal("PlacesAuthentication is nil — EF_PLACES_AUTHENTICATION not wired into switch")
	}
	if len(pa.GetRecords()) != 1 {
		t.Fatalf("Records count: got %d, want 1", len(pa.GetRecords()))
	}
	if pa.GetRecords()[0].GetAuthenticationStatus() != ddv1.PositionAuthenticationStatus_AUTHENTICATED {
		t.Errorf("AuthenticationStatus: got %v, want AUTHENTICATED", pa.GetRecords()[0].GetAuthenticationStatus())
	}
}

func TestPlacesAuthentication_Fixtures(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_PLACES_AUTHENTICATION,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("no hexdump fixtures found for EF_PLACES_AUTHENTICATION GENERATION_2")
	}
	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		t.Run(strings.TrimSuffix(relPath, ".hexdump"), func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("failed to read hexdump: %v", err)
			}
			got, err := UnmarshalOptions{}.unmarshalPlacesAuthentication(data)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			loadOrCreateGolden(t, got, goldenJSONPath(hexdumpPath))
		})
	}
}
