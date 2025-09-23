package tachograph

import (
	"os"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func TestDebugMarshalOrder(t *testing.T) {
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

	t.Logf("=== MARSHALLING STEP BY STEP ===")

	// Start with empty buffer
	dst := make([]byte, 0)

	// 1. EF_ICC (no signature)
	t.Logf("1. Adding EF_ICC...")
	if driverCard.GetIcc() != nil {
		beforeLen := len(dst)
		dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_ICC, driverCard.GetIcc(), AppendIcc)
		if err != nil {
			t.Fatalf("Failed to append EF_ICC: %v", err)
		}
		afterLen := len(dst)
		t.Logf("   EF_ICC added %d bytes: %X", afterLen-beforeLen, dst[beforeLen:afterLen])
	} else {
		t.Logf("   EF_ICC is nil!")
	}

	// 2. EF_IC (no signature)
	t.Logf("2. Adding EF_IC...")
	if driverCard.GetIc() != nil {
		beforeLen := len(dst)
		dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_IC, driverCard.GetIc(), AppendCardIc)
		if err != nil {
			t.Fatalf("Failed to append EF_IC: %v", err)
		}
		afterLen := len(dst)
		t.Logf("   EF_IC added %d bytes: %X", afterLen-beforeLen, dst[beforeLen:afterLen])
	} else {
		t.Logf("   EF_IC is nil!")
	}

	// Compare first 50 bytes with original
	t.Logf("\n=== COMPARISON (first 50 bytes) ===")
	t.Logf("Original:   %X", originalData[:50])
	t.Logf("Marshalled: %X", dst[:minInt(50, len(dst))])

	// Find first difference
	for i := 0; i < minInt(len(originalData), len(dst)); i++ {
		if originalData[i] != dst[i] {
			t.Logf("First difference at byte %d: original=0x%02X, marshalled=0x%02X", i, originalData[i], dst[i])
			break
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
