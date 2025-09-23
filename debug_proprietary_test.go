package tachograph

import (
	"os"
	"testing"
	"unsafe"
)

func TestDebugProprietaryEFs(t *testing.T) {
	filePath := "testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"

	originalData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	t.Logf("Starting unmarshalling...")
	file, err := UnmarshalFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check if proprietary EFs were stored
	cardPtr := uintptr(unsafe.Pointer(file.GetDriverCard()))
	proprietaryEFs := GetProprietaryEFs(cardPtr)

	if proprietaryEFs == nil {
		t.Logf("No proprietary EFs found")
	} else {
		t.Logf("Found %d proprietary EFs:", len(proprietaryEFs.EFs))
		for i, ef := range proprietaryEFs.EFs {
			t.Logf("  EF %d: FID=0x%04X, Data length=%d", i+1, ef.FID, len(ef.Data))
		}
	}
}
