package safemath

import (
	"errors"
	"math"
	"testing"
)

func TestMulInt(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr error
	}{
		{name: "zero_a", a: 0, b: 100, want: 0},
		{name: "zero_b", a: 100, b: 0, want: 0},
		{name: "both_zero", a: 0, b: 0, want: 0},
		{name: "one_times_one", a: 1, b: 1, want: 1},
		{name: "small_positive", a: 7, b: 8, want: 56},
		{name: "large_valid", a: 65535, b: 65535, want: 65535 * 65535},
		{name: "negative_a", a: -3, b: 5, want: -15},
		{name: "both_negative", a: -4, b: -6, want: 24},
		{name: "max_int_overflow", a: math.MaxInt, b: 2, wantErr: ErrOverflow},
		{name: "near_max_overflow", a: math.MaxInt/2 + 1, b: 2, wantErr: ErrOverflow},
		{name: "min_int_overflow", a: math.MinInt, b: 2, wantErr: ErrOverflow},
		{name: "uint16_max_product", a: 0xFFFF, b: 0xFFFF, want: 0xFFFF * 0xFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulInt(tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("MulInt(%d, %d) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("MulInt(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				return
			}
			if got != tt.want {
				t.Errorf("MulInt(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMulInt_Commutative(t *testing.T) {
	// Overflow detection must work regardless of argument order
	got1, err1 := MulInt(math.MaxInt, 2)
	got2, err2 := MulInt(2, math.MaxInt)

	if !errors.Is(err1, ErrOverflow) || !errors.Is(err2, ErrOverflow) {
		t.Errorf("expected overflow in both directions, got err1=%v err2=%v", err1, err2)
	}
	_ = got1
	_ = got2
}
