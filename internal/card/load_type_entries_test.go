package card

import (
	"encoding/binary"
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// makeLoadTypeRecord builds a 5-byte CardLoadTypeEntryRecord payload.
//
// Layout: 4B timestamp | 1B loadType
func makeLoadTypeRecord(timestampSecs uint32, loadType byte) []byte {
	r := make([]byte, 5)
	binary.BigEndian.PutUint32(r[0:4], timestampSecs)
	r[4] = loadType
	return r
}

func TestUnmarshalLoadTypeEntries(t *testing.T) {
	t.Run("empty input error", func(t *testing.T) {
		_, err := UnmarshalOptions{}.unmarshalLoadTypeEntries([]byte{})
		if err == nil {
			t.Fatal("expected error for empty input, got nil")
		}
	})

	t.Run("pointer only, zero records", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		got, err := UnmarshalOptions{}.unmarshalLoadTypeEntries(data)
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

	t.Run("one record", func(t *testing.T) {
		// loadType byte 0x01 → GOODS (protocol_enum_value=1)
		rec := makeLoadTypeRecord(3000, 0x01)
		data := append([]byte{0x00, 0x00}, rec...)

		got, err := UnmarshalOptions{}.unmarshalLoadTypeEntries(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.GetRecords()) != 1 {
			t.Fatalf("Records count: got %d, want 1", len(got.GetRecords()))
		}
		r := got.GetRecords()[0]
		if r.GetTimestamp().GetSeconds() != 3000 {
			t.Errorf("Timestamp: got %d, want 3000", r.GetTimestamp().GetSeconds())
		}
		if r.GetLoadTypeEntered() != ddv1.LoadType_GOODS {
			t.Errorf("LoadTypeEntered: got %v, want GOODS", r.GetLoadTypeEntered())
		}
	})

	t.Run("multiple records", func(t *testing.T) {
		rec1 := makeLoadTypeRecord(1000, 0x01) // GOODS
		rec2 := makeLoadTypeRecord(2000, 0x02) // PASSENGERS
		rec3 := makeLoadTypeRecord(3000, 0x00) // NOT_DEFINED
		data := append([]byte{0x00, 0x02}, rec1...)
		data = append(data, rec2...)
		data = append(data, rec3...)

		got, err := UnmarshalOptions{}.unmarshalLoadTypeEntries(data)
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
		_, err := UnmarshalOptions{}.unmarshalLoadTypeEntries(data)
		if err == nil {
			t.Fatal("expected error for truncated record, got nil")
		}
	})
}

func TestLoadTypeEntries_ParseRawDriverCardFile(t *testing.T) {
	rec := makeLoadTypeRecord(5000, 0x02) // PASSENGERS
	efData := append([]byte{0x00, 0x00}, rec...)

	rawRec := &cardv1.RawCardFile_Record{}
	rawRec.SetFile(cardv1.ElementaryFileType_EF_LOAD_TYPE_ENTRIES)
	rawRec.SetGeneration(ddv1.Generation_GENERATION_2)
	rawRec.SetContentType(cardv1.ContentType_DATA)
	rawRec.SetValue(efData)

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords([]*cardv1.RawCardFile_Record{rawRec})

	cardFile, err := ParseOptions{}.ParseRawDriverCardFile(rawFile)
	if err != nil {
		t.Fatalf("ParseRawDriverCardFile: %v", err)
	}
	lt := cardFile.GetTachographG2().GetLoadTypeEntries()
	if lt == nil {
		t.Fatal("LoadTypeEntries is nil — EF_LOAD_TYPE_ENTRIES not wired into switch")
	}
	if len(lt.GetRecords()) != 1 {
		t.Fatalf("Records count: got %d, want 1", len(lt.GetRecords()))
	}
	if lt.GetRecords()[0].GetLoadTypeEntered() != ddv1.LoadType_PASSENGERS {
		t.Errorf("LoadTypeEntered: got %v, want PASSENGERS", lt.GetRecords()[0].GetLoadTypeEntered())
	}
}

func TestLoadTypeEntries_Fixtures(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_LOAD_TYPE_ENTRIES,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("no hexdump fixtures found for EF_LOAD_TYPE_ENTRIES GENERATION_2")
	}
	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		t.Run(strings.TrimSuffix(relPath, ".hexdump"), func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("failed to read hexdump: %v", err)
			}
			got, err := UnmarshalOptions{}.unmarshalLoadTypeEntries(data)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			loadOrCreateGolden(t, got, goldenJSONPath(hexdumpPath))
		})
	}
}
