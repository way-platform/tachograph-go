package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func ptr32(v int32) *int32 { return &v }

func TestUnmarshalOdometer(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		want       *int32
		wantErr    bool
		errMessage string
	}{
		// Sentinel: 0xFFFFFF → nil (odometer not available)
		{
			name:  "sentinel 0xFFFFFF returns nil",
			input: []byte{0xFF, 0xFF, 0xFF},
			want:  nil,
		},
		// Zero is a valid odometer value, not a sentinel
		{
			name:  "zero value is valid",
			input: []byte{0x00, 0x00, 0x00},
			want:  ptr32(0),
		},
		// Value just below the sentinel is valid
		{
			name:  "0xFFFFFE is valid (16777214)",
			input: []byte{0xFF, 0xFF, 0xFE},
			want:  ptr32(16777214),
		},
		// Normal values
		{
			name:  "maximum spec value 999999",
			input: []byte{0x0F, 0x42, 0x3F},
			want:  ptr32(999999),
		},
		{
			name:  "middle value 123456",
			input: []byte{0x01, 0xE2, 0x40},
			want:  ptr32(123456),
		},
		{
			name:  "value 375000 (0x05B8D8)",
			input: []byte{0x05, 0xB8, 0xD8},
			want:  ptr32(375000),
		},
		{
			name:  "value 1",
			input: []byte{0x00, 0x00, 0x01},
			want:  ptr32(1),
		},
		// Error cases
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
			name:       "insufficient data - 0 bytes",
			input:      []byte{},
			wantErr:    true,
			errMessage: "invalid data length for OdometerShort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unmarshalOpts := UnmarshalOptions{PreserveRawData: true}
			got, err := unmarshalOpts.UnmarshalOdometer(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalOdometer() expected error containing %q, got nil", tt.errMessage)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalOdometer() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("UnmarshalOdometer() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendOdometer(t *testing.T) {
	tests := []struct {
		name  string
		input int32
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
			opts := MarshalOptions{}
			got, err := opts.MarshalOdometer(tt.input)
			if err != nil {
				t.Fatalf("MarshalOdometer() unexpected error: %v", err)
			}
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
			unmarshalOpts := UnmarshalOptions{}
			marshalOpts := MarshalOptions{}
			odometer, err := unmarshalOpts.UnmarshalOdometer(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalOdometer() unexpected error: %v", err)
			}
			if odometer == nil {
				t.Fatal("UnmarshalOdometer() returned nil for non-sentinel value")
			}

			got, err := marshalOpts.MarshalOdometer(*odometer)
			if err != nil {
				t.Fatalf("MarshalOdometer() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
