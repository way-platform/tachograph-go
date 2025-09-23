package tachograph

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestDebugUnmarshalLoop(t *testing.T) {
	filePath := "testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	t.Logf("File size: %d bytes", len(data))

	r := bytes.NewReader(data)
	iteration := 0

	for r.Len() > 0 {
		iteration++
		remainingBytes := r.Len()

		if remainingBytes < 5 {
			t.Logf("Iteration %d: Only %d bytes remaining, stopping", iteration, remainingBytes)
			break
		}

		// Read Tag - 3 bytes (FID + appendix)
		tagBytes := make([]byte, 3)
		if _, err := r.Read(tagBytes); err != nil {
			t.Logf("Iteration %d: Error reading tag: %v", iteration, err)
			break
		}

		fid := binary.BigEndian.Uint16(tagBytes[0:2])
		appendix := tagBytes[2]

		// Read Length - 2 bytes
		var length uint16
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			t.Logf("Iteration %d: Error reading length: %v", iteration, err)
			break
		}

		t.Logf("Iteration %d: FID=0x%04X, Appendix=0x%02X, Length=%d, Remaining after header=%d",
			iteration, fid, appendix, length, r.Len())

		// Skip the value data
		if int(length) > r.Len() {
			t.Logf("Iteration %d: Length %d exceeds remaining bytes %d, stopping", iteration, length, r.Len())
			break
		}

		value := make([]byte, length)
		if _, err := r.Read(value); err != nil {
			t.Logf("Iteration %d: Error reading value: %v", iteration, err)
			break
		}

		// Check if this is a proprietary EF
		if fid >= 0xC000 {
			t.Logf("  -> PROPRIETARY EF detected! FID=0x%04X", fid)
		}

		if iteration > 30 { // Safety limit
			t.Logf("Stopping after %d iterations for safety", iteration)
			break
		}
	}

	t.Logf("Loop completed after %d iterations, %d bytes remaining", iteration, r.Len())
}
