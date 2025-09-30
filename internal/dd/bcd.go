package dd

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// decodeBCD converts BCD-encoded bytes to an integer
func decodeBCD(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	s := hex.EncodeToString(b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid BCD value: %s", s)
	}
	return int(i), nil
}
