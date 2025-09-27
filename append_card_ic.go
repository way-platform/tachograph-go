package tachograph

import (
	"encoding/hex"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardIc appends IC identification data to a byte slice.
func AppendCardIc(data []byte, ic *cardv1.Ic) ([]byte, error) {
	if ic == nil {
		return data, nil
	}

	// IC Serial Number (4 bytes)
	serialStr := ic.GetIcSerialNumber()
	if len(serialStr) > 0 {
		// Remove any spaces and decode hex
		serialStr = removeSpaces(serialStr)
		if len(serialStr) >= 8 {
			serialBytes, err := hex.DecodeString(serialStr[:8])
			if err == nil && len(serialBytes) == 4 {
				data = append(data, serialBytes...)
			} else {
				// Fallback: pad or truncate to 4 bytes
				serialBytes := make([]byte, 4)
				copy(serialBytes, serialStr)
				data = append(data, serialBytes...)
			}
		} else {
			// Pad with zeros
			serialBytes := make([]byte, 4)
			copy(serialBytes, serialStr)
			data = append(data, serialBytes...)
		}
	} else {
		data = append(data, 0x00, 0x00, 0x00, 0x00)
	}

	// IC Manufacturing References (4 bytes)
	mfgStr := ic.GetIcManufacturingReferences()
	if len(mfgStr) > 0 {
		// Remove any spaces and decode hex
		mfgStr = removeSpaces(mfgStr)
		if len(mfgStr) >= 8 {
			mfgBytes, err := hex.DecodeString(mfgStr[:8])
			if err == nil && len(mfgBytes) == 4 {
				data = append(data, mfgBytes...)
			} else {
				// Fallback: pad or truncate to 4 bytes
				mfgBytes := make([]byte, 4)
				copy(mfgBytes, mfgStr)
				data = append(data, mfgBytes...)
			}
		} else {
			// Pad with zeros
			mfgBytes := make([]byte, 4)
			copy(mfgBytes, mfgStr)
			data = append(data, mfgBytes...)
		}
	} else {
		data = append(data, 0x00, 0x00, 0x00, 0x00)
	}

	return data, nil
}

// removeSpaces removes spaces from a string
func removeSpaces(s string) string {
	result := ""
	for _, c := range s {
		if c != ' ' {
			result += string(c)
		}
	}
	return result
}
