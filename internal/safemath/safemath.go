// Package safemath provides overflow-safe integer arithmetic for parsing
// untrusted binary data. These functions detect integer overflow before it
// occurs, preventing silent wraparound that could lead to incorrect buffer
// sizes, out-of-bounds reads, or memory exhaustion.
package safemath

import "errors"

// ErrOverflow is returned when an integer arithmetic operation would overflow.
var ErrOverflow = errors.New("integer overflow")

// MulInt returns a * b, or ErrOverflow if the result would overflow int.
// This is safe for both 32-bit and 64-bit platforms.
func MulInt(a, b int) (int, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	result := a * b
	if result/a != b {
		return 0, ErrOverflow
	}
	return result, nil
}
