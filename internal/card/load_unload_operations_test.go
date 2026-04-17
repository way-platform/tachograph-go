package card

import (
	"encoding/binary"
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// makeLoadUnloadRecord builds a 20-byte CardLoadUnloadRecord payload.
//
// Layout: 4B timestamp | 1B operationType | 12B gnssPlaceAuthRecord | 3B odometer
func makeLoadUnloadRecord(timestampSecs uint32, opType byte, odometer [3]byte) []byte {
	r := make([]byte, 20)
	binary.BigEndian.PutUint32(r[0:4], timestampSecs)
	r[4] = opType
	// gnssPlaceAuthRecord (12 bytes): timestamp(4) accuracy(1) lat(3) lon(3) authStatus(1)
	binary.BigEndian.PutUint32(r[5:9], 0)
	r[9] = 0x00
	copy(r[10:13], []byte{0x00, 0x00, 0x00})
	copy(r[13:16], []byte{0x00, 0x00, 0x00})
	r[16] = 0x00
	copy(r[17:20], odometer[:])
	return r
}

func TestUnmarshalLoadUnloadOperations(t *testing.T) {
	t.Run("empty input error", func(t *testing.T) {
		_, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations([]byte{})
		if err == nil {
			t.Fatal("expected error for empty input, got nil")
		}
	})

	t.Run("pointer only, zero records", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		got, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations(data)
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
		rec := makeLoadUnloadRecord(1000, 0x01, [3]byte{0x00, 0x00, 0x64}) // op=LOAD, 100 km
		data := append([]byte{0x00, 0x00}, rec...)

		got, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.GetRecords()) != 1 {
			t.Fatalf("Records count: got %d, want 1", len(got.GetRecords()))
		}
		r := got.GetRecords()[0]
		if r.GetTimestamp().GetSeconds() != 1000 {
			t.Errorf("Timestamp: got %d, want 1000", r.GetTimestamp().GetSeconds())
		}
		if r.GetOperationType() != ddv1.OperationType_LOAD_OPERATION {
			t.Errorf("OperationType: got %v, want LOAD_OPERATION", r.GetOperationType())
		}
		if r.GetVehicleOdometerKm() != 100 {
			t.Errorf("VehicleOdometerKm: got %d, want 100", r.GetVehicleOdometerKm())
		}
	})

	t.Run("two records", func(t *testing.T) {
		rec1 := makeLoadUnloadRecord(500, 0x01, [3]byte{0x00, 0x00, 0x0A})
		rec2 := makeLoadUnloadRecord(600, 0x02, [3]byte{0x00, 0x00, 0x14})
		data := append([]byte{0x00, 0x01}, rec1...)
		data = append(data, rec2...)

		got, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations(data)
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
		// 2-byte pointer + 15 bytes (not a multiple of 20)
		data := make([]byte, 17)
		_, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations(data)
		if err == nil {
			t.Fatal("expected error for truncated record, got nil")
		}
	})

	t.Run("sentinel odometer 0xFFFFFF leaves field unset", func(t *testing.T) {
		rec := makeLoadUnloadRecord(1000, 0x01, [3]byte{0xFF, 0xFF, 0xFF})
		data := append([]byte{0x00, 0x00}, rec...)

		got, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		r := got.GetRecords()[0]
		if r.HasVehicleOdometerKm() {
			t.Errorf("VehicleOdometerKm should be unset for sentinel 0xFFFFFF, got %d", r.GetVehicleOdometerKm())
		}
	})
}

func TestLoadUnloadOperations_ParseRawDriverCardFile(t *testing.T) {
	rec := makeLoadUnloadRecord(2000, 0x01, [3]byte{0x00, 0x03, 0xE8}) // 1000 km
	efData := append([]byte{0x00, 0x00}, rec...)

	rawRec := &cardv1.RawCardFile_Record{}
	rawRec.SetFile(cardv1.ElementaryFileType_EF_LOAD_UNLOAD_OPERATIONS)
	rawRec.SetGeneration(ddv1.Generation_GENERATION_2)
	rawRec.SetContentType(cardv1.ContentType_DATA)
	rawRec.SetValue(efData)

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords([]*cardv1.RawCardFile_Record{rawRec})

	cardFile, err := ParseOptions{}.ParseRawDriverCardFile(rawFile)
	if err != nil {
		t.Fatalf("ParseRawDriverCardFile: %v", err)
	}
	lu := cardFile.GetTachographG2().GetLoadUnloadOperations()
	if lu == nil {
		t.Fatal("LoadUnloadOperations is nil — EF_LOAD_UNLOAD_OPERATIONS not wired into switch")
	}
	if len(lu.GetRecords()) != 1 {
		t.Fatalf("Records count: got %d, want 1", len(lu.GetRecords()))
	}
	if lu.GetRecords()[0].GetVehicleOdometerKm() != 1000 {
		t.Errorf("VehicleOdometerKm: got %d, want 1000", lu.GetRecords()[0].GetVehicleOdometerKm())
	}
}

func TestLoadUnloadOperations_Fixtures(t *testing.T) {
	hexdumpFiles, err := findHexdumpFiles(
		cardv1.ElementaryFileType_EF_LOAD_UNLOAD_OPERATIONS,
		ddv1.Generation_GENERATION_2,
		cardv1.ContentType_DATA,
	)
	if err != nil {
		t.Fatalf("failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("no hexdump fixtures found for EF_LOAD_UNLOAD_OPERATIONS GENERATION_2")
	}
	for _, hexdumpPath := range hexdumpFiles {
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		t.Run(strings.TrimSuffix(relPath, ".hexdump"), func(t *testing.T) {
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("failed to read hexdump: %v", err)
			}
			got, err := UnmarshalOptions{}.unmarshalLoadUnloadOperations(data)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			loadOrCreateGolden(t, got, goldenJSONPath(hexdumpPath))
		})
	}
}
