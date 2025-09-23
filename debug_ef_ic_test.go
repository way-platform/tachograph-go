package tachograph

import (
	"os"
	"testing"
)

func TestDebugEFIC(t *testing.T) {
	filePath := "testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"
	
	originalData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	file, err := UnmarshalFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	driverCard := file.GetDriverCard()
	if driverCard == nil {
		t.Fatalf("No driver card")
	}

	// Check EF_IC specifically
	ic := driverCard.GetIc()
	if ic == nil {
		t.Fatalf("EF_IC is nil!")
	}

	t.Logf("EF_IC present:")
	t.Logf("  IC Serial Number: %q", ic.GetIcSerialNumber())
	t.Logf("  IC Manufacturing References: %q", ic.GetIcManufacturingReferences())

	// Test marshalling EF_IC specifically
	icData, err := AppendCardIc(nil, ic)
	if err != nil {
		t.Fatalf("Failed to append EF_IC: %v", err)
	}
	
	t.Logf("EF_IC marshalled data (%d bytes): %X", len(icData), icData)
	
	// EF_IC data starts at offset 35 (after 30 bytes ICC + 5 bytes IC tag)
	expectedICData := originalData[35:43] // 8 bytes of EF_IC data
	t.Logf("Expected EF_IC data (%d bytes): %X", len(expectedICData), expectedICData)
	
	if len(icData) != len(expectedICData) {
		t.Errorf("Length mismatch: got %d, expected %d", len(icData), len(expectedICData))
	}
	
	for i := 0; i < len(icData) && i < len(expectedICData); i++ {
		if icData[i] != expectedICData[i] {
			t.Errorf("Byte %d mismatch: got 0x%02X, expected 0x%02X", i, icData[i], expectedICData[i])
		}
	}
}
