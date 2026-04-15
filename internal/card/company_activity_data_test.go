package card

import (
	"encoding/binary"
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// makeCompanyActivityRecord builds a 47-byte companyActivityRecord payload.
//
// Layout:
//   - Byte  0:    companyActivityType (1 byte)
//   - Bytes 1-4:  companyActivityTime (4 bytes, big-endian TimeReal)
//   - Bytes 5-23: cardNumberInformation (19 bytes, zeroed)
//   - Bytes 24-38: vehicleRegistrationInformation (15 bytes, zeroed)
//   - Bytes 39-42: downloadPeriodBegin (4 bytes, big-endian TimeReal)
//   - Bytes 43-46: downloadPeriodEnd (4 bytes, big-endian TimeReal)
func makeCompanyActivityRecord(activityType byte, activityTimeSecs, beginSecs, endSecs uint32) []byte {
	r := make([]byte, 47)
	r[0] = activityType
	binary.BigEndian.PutUint32(r[1:5], activityTimeSecs)
	// cardNumberInformation (19 bytes): all zeros = valid FullCardNumber(18)+Generation(1)
	// vehicleRegistrationInformation (15 bytes): all zeros
	binary.BigEndian.PutUint32(r[39:43], beginSecs)
	binary.BigEndian.PutUint32(r[43:47], endSecs)
	return r
}

func TestUnmarshalCompanyActivityData(t *testing.T) {
	t.Run("empty input error", func(t *testing.T) {
		_, err := UnmarshalOptions{}.unmarshalCompanyActivityData([]byte{})
		if err == nil {
			t.Fatal("expected error for empty input, got nil")
		}
	})

	t.Run("pointer only, zero records", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		got, err := UnmarshalOptions{}.unmarshalCompanyActivityData(data)
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

	t.Run("one record card-downloading", func(t *testing.T) {
		// protocol_enum_value=1 → CARD_DOWNLOADING
		rec := makeCompanyActivityRecord(0x01, 5000, 1000, 4000)
		data := append([]byte{0x00, 0x00}, rec...)

		got, err := UnmarshalOptions{}.unmarshalCompanyActivityData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.GetRecords()) != 1 {
			t.Fatalf("Records count: got %d, want 1", len(got.GetRecords()))
		}
		r := got.GetRecords()[0]
		if r.GetCompanyActivityType() != ddv1.CompanyActivityType_CARD_DOWNLOADING {
			t.Errorf("CompanyActivityType: got %v, want CARD_DOWNLOADING", r.GetCompanyActivityType())
		}
		if r.GetCompanyActivityTime().GetSeconds() != 5000 {
			t.Errorf("CompanyActivityTime: got %d, want 5000", r.GetCompanyActivityTime().GetSeconds())
		}
		if r.GetDownloadPeriodBegin().GetSeconds() != 1000 {
			t.Errorf("DownloadPeriodBegin: got %d, want 1000", r.GetDownloadPeriodBegin().GetSeconds())
		}
		if r.GetDownloadPeriodEnd().GetSeconds() != 4000 {
			t.Errorf("DownloadPeriodEnd: got %d, want 4000", r.GetDownloadPeriodEnd().GetSeconds())
		}
	})

	t.Run("two records", func(t *testing.T) {
		rec1 := makeCompanyActivityRecord(0x01, 1000, 500, 900)  // CARD_DOWNLOADING
		rec2 := makeCompanyActivityRecord(0x02, 2000, 1500, 1900) // VU_DOWNLOADING
		data := append([]byte{0x00, 0x01}, rec1...)
		data = append(data, rec2...)

		got, err := UnmarshalOptions{}.unmarshalCompanyActivityData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetNewestRecordIndex() != 1 {
			t.Errorf("NewestRecordIndex: got %d, want 1", got.GetNewestRecordIndex())
		}
		if len(got.GetRecords()) != 2 {
			t.Fatalf("Records count: got %d, want 2", len(got.GetRecords()))
		}
		if got.GetRecords()[1].GetCompanyActivityType() != ddv1.CompanyActivityType_VU_DOWNLOADING {
			t.Errorf("Record[1] CompanyActivityType: got %v, want VU_DOWNLOADING", got.GetRecords()[1].GetCompanyActivityType())
		}
	})

	t.Run("truncated record returns error", func(t *testing.T) {
		// 2-byte pointer + 30 bytes (not a multiple of 47)
		data := make([]byte, 32)
		_, err := UnmarshalOptions{}.unmarshalCompanyActivityData(data)
		if err == nil {
			t.Fatal("expected error for truncated record, got nil")
		}
	})
}

func TestCompanyActivityData_ParseRawDriverCardFile(t *testing.T) {
	rec := makeCompanyActivityRecord(0x01, 8000, 6000, 7500) // CARD_DOWNLOADING
	efData := append([]byte{0x00, 0x00}, rec...)

	rawRec := &cardv1.RawCardFile_Record{}
	rawRec.SetFile(cardv1.ElementaryFileType_EF_COMPANY_ACTIVITY_DATA)
	rawRec.SetGeneration(ddv1.Generation_GENERATION_2)
	rawRec.SetContentType(cardv1.ContentType_DATA)
	rawRec.SetValue(efData)

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords([]*cardv1.RawCardFile_Record{rawRec})

	cardFile, err := ParseOptions{}.ParseRawDriverCardFile(rawFile)
	if err != nil {
		t.Fatalf("ParseRawDriverCardFile: %v", err)
	}
	cad := cardFile.GetTachographG2().GetCompanyActivityData()
	if cad == nil {
		t.Fatal("CompanyActivityData is nil — EF_COMPANY_ACTIVITY_DATA not wired into switch")
	}
	if len(cad.GetRecords()) != 1 {
		t.Fatalf("Records count: got %d, want 1", len(cad.GetRecords()))
	}
	if cad.GetRecords()[0].GetCompanyActivityType() != ddv1.CompanyActivityType_CARD_DOWNLOADING {
		t.Errorf("CompanyActivityType: got %v, want CARD_DOWNLOADING", cad.GetRecords()[0].GetCompanyActivityType())
	}
}

func TestCompanyActivityData_Fixtures(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_COMPANY_ACTIVITY_DATA,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("no hexdump fixtures found for EF_COMPANY_ACTIVITY_DATA GENERATION_2")
	}
	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		t.Run(strings.TrimSuffix(relPath, ".hexdump"), func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("failed to read hexdump: %v", err)
			}
			got, err := UnmarshalOptions{}.unmarshalCompanyActivityData(data)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			loadOrCreateGolden(t, got, goldenJSONPath(hexdumpPath))
		})
	}
}
