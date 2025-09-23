package tachograph

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestDebugRawComparison(t *testing.T) {
	filePath := "testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"

	originalData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// 1. Parse original file into RawCardFile
	originalRaw, err := UnmarshalRawCardFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal original to raw: %v", err)
	}

	// 2. Parse original file into DriverCardFile
	file, err := UnmarshalFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal file: %v", err)
	}

	driverCard := file.GetDriverCard()
	if driverCard == nil {
		t.Fatalf("No driver card")
	}

	// 3. Convert DriverCardFile to RawCardFile
	marshalledRaw, err := DriverCardFileToRaw(driverCard)
	if err != nil {
		t.Fatalf("Failed to convert driver card to raw: %v", err)
	}

	// 4. Compare the two RawCardFile structures
	t.Logf("=== RAW CARD FILE COMPARISON ===")
	t.Logf("Original records: %d", len(originalRaw.GetRecords()))
	t.Logf("Marshalled records: %d", len(marshalledRaw.GetRecords()))

	// Compare each record
	minRecords := len(originalRaw.GetRecords())
	if len(marshalledRaw.GetRecords()) < minRecords {
		minRecords = len(marshalledRaw.GetRecords())
	}

	for i := 0; i < minRecords; i++ {
		origRecord := originalRaw.GetRecords()[i]
		marshRecord := marshalledRaw.GetRecords()[i]

		t.Logf("\n--- Record %d ---", i)
		t.Logf("Original:   Tag=0x%06X, File=%s, ContentType=%s, Length=%d",
			origRecord.GetTag(), origRecord.GetFile(), origRecord.GetContentType(), origRecord.GetLength())
		t.Logf("Marshalled: Tag=0x%06X, File=%s, ContentType=%s, Length=%d",
			marshRecord.GetTag(), marshRecord.GetFile(), marshRecord.GetContentType(), marshRecord.GetLength())

		if origRecord.GetTag() != marshRecord.GetTag() {
			t.Errorf("  ❌ Tag mismatch: 0x%06X vs 0x%06X", origRecord.GetTag(), marshRecord.GetTag())
		}

		if origRecord.GetLength() != marshRecord.GetLength() {
			t.Errorf("  ❌ Length mismatch: %d vs %d", origRecord.GetLength(), marshRecord.GetLength())
		}

		// Compare values as hex strings for better readability
		origHex := hex.EncodeToString(origRecord.GetValue())
		marshHex := hex.EncodeToString(marshRecord.GetValue())

		if origHex != marshHex {
			t.Logf("  ❌ Value mismatch:")
			t.Logf("    Original:   %s", origHex)
			t.Logf("    Marshalled: %s", marshHex)

			// Use cmp.Diff for detailed comparison
			diff := cmp.Diff(origHex, marshHex)
			t.Logf("    Diff: %s", diff)
		} else {
			t.Logf("  ✅ Values match")
		}
	}

	// Report missing records
	if len(originalRaw.GetRecords()) > len(marshalledRaw.GetRecords()) {
		t.Errorf("Missing %d records in marshalled data", len(originalRaw.GetRecords())-len(marshalledRaw.GetRecords()))
		for i := minRecords; i < len(originalRaw.GetRecords()); i++ {
			record := originalRaw.GetRecords()[i]
			t.Logf("Missing record %d: Tag=0x%06X, File=%s, Length=%d",
				i, record.GetTag(), record.GetFile(), record.GetLength())
		}
	}

	// 5. Also compare the full protobuf structures
	t.Logf("\n=== PROTOBUF STRUCTURE COMPARISON ===")
	if diff := cmp.Diff(originalRaw, marshalledRaw, protocmp.Transform()); diff != "" {
		t.Logf("RawCardFile structures differ:\n%s", diff)
	} else {
		t.Logf("✅ RawCardFile structures are identical!")
	}
}
