package tachograph

import (
	"encoding/binary"
	"os"
	"testing"
)

func TestAnalyzeTLVStructure(t *testing.T) {
	filePath := "testdata/card/driver/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	t.Logf("Analyzing TLV structure of file: %s", filePath)
	t.Logf("Total file size: %d bytes", len(data))

	offset := 0
	for offset < len(data) {
		if offset+5 > len(data) {
			break
		}

		// Read tag (3 bytes: FID + appendix) and length (2 bytes)
		fid := binary.BigEndian.Uint16(data[offset : offset+2])
		appendix := data[offset+2]
		length := binary.BigEndian.Uint16(data[offset+3 : offset+5])

		t.Logf("Offset %d: FID=0x%04X, Appendix=0x%02X, Length=%d", offset, fid, appendix, length)

		// Move to next TLV
		offset += 5 + int(length)

		// Stop if we've read enough or if length seems invalid
		if length > 20000 || offset >= len(data) {
			break
		}
	}
}
