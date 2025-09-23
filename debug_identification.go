package tachograph

import (
	"os"
	"testing"
)

func TestDebugIdentification(t *testing.T) {
	filePath := "testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"

	originalData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal to get the identification data
	file, err := UnmarshalFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	driverCard := file.GetDriverCard()
	if driverCard == nil {
		t.Fatalf("Not a driver card")
	}

	if id := driverCard.GetIdentification(); id != nil {
		t.Logf("CardIssuingMemberState: %q (len=%d)", id.GetCardIssuingMemberState(), len(id.GetCardIssuingMemberState()))
		for i, b := range []byte(id.GetCardIssuingMemberState()) {
			t.Logf("  byte[%d] = 0x%02X ('%c')", i, b, b)
		}
		t.Logf("CardNumber: %q", id.GetCardNumber())
		t.Logf("CardIssuingAuthorityName: %q", id.GetCardIssuingAuthorityName())
	}
}
