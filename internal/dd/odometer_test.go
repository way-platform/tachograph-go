package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalOdometer(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		want       uint32
		wantErr    bool
		errMessage string
	}{
		{
			name:  "maximum value 999999",
			input: []byte{0x0F, 0x42, 0x3F},
			want:  999999,
		},
		{
			name:  "zero value",
			input: []byte{0x00, 0x00, 0x00},
			want:  0,
		},
		{
			name:  "middle value 123456",
			input: []byte{0x01, 0xE2, 0x40},
			want:  123456,
		},
		{
			name:  "value 1",
			input: []byte{0x00, 0x00, 0x01},
			want:  1,
		},
		{
			name:  "value 255",
			input: []byte{0x00, 0x00, 0xFF},
			want:  255,
		},
		{
			name:  "value 65535",
			input: []byte{0x00, 0xFF, 0xFF},
			want:  65535,
		},
		{
			name:       "buffer larger than needed - exact length required",
			input:      []byte{0x01, 0xE2, 0x40, 0xFF, 0xFF},
			wantErr:    true,
			errMessage: "invalid data length for OdometerShort",
		},
		{
			name:       "insufficient data - 2 bytes",
			input:      []byte{0x01, 0x02},
			wantErr:    true,
			errMessage: "invalid data length for OdometerShort",
		},
		{
			name:       "insufficient data - 1 byte",
			input:      []byte{0x01},
			wantErr:    true,
			errMessage: "invalid data length for OdometerShort",
		},
		{
			name:       "insufficient data - 0 bytes",
			input:      []byte{},
			wantErr:    true,
			errMessage: "invalid data length for OdometerShort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts UnmarshalOptions
			got, err := opts.UnmarshalOdometer(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalOdometer() expected error containing %q, got nil", tt.errMessage)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalOdometer() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("UnmarshalOdometer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendOdometer(t *testing.T) {
	tests := []struct {
		name  string
		input uint32
		want  []byte
	}{
		{
			name:  "maximum value 999999",
			input: 999999,
			want:  []byte{0x0F, 0x42, 0x3F},
		},
		{
			name:  "zero value",
			input: 0,
			want:  []byte{0x00, 0x00, 0x00},
		},
		{
			name:  "middle value 123456",
			input: 123456,
			want:  []byte{0x01, 0xE2, 0x40},
		},
		{
			name:  "value 1",
			input: 1,
			want:  []byte{0x00, 0x00, 0x01},
		},
		{
			name:  "value 255",
			input: 255,
			want:  []byte{0x00, 0x00, 0xFF},
		},
		{
			name:  "value 65535",
			input: 65535,
			want:  []byte{0x00, 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte{}
			got := AppendOdometer(dst, tt.input)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendOdometer() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOdometerRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "maximum value",
			input: []byte{0x0F, 0x42, 0x3F},
		},
		{
			name:  "zero value",
			input: []byte{0x00, 0x00, 0x00},
		},
		{
			name:  "middle value",
			input: []byte{0x01, 0xE2, 0x40},
		},
		{
			name:  "small value",
			input: []byte{0x00, 0x00, 0x01},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			odometer, err := opts.UnmarshalOdometer(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalOdometer() unexpected error: %v", err)
			}

			// Marshal
			got := AppendOdometer(nil, odometer)

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
