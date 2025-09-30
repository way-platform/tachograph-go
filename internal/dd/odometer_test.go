package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadOdometerFromBytes(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		offset     int
		want       uint32
		wantOffset int
		wantErr    bool
	}{
		{
			name:       "maximum value 999999",
			input:      []byte{0x0F, 0x42, 0x3F},
			offset:     0,
			want:       999999,
			wantOffset: 3,
		},
		{
			name:       "zero value",
			input:      []byte{0x00, 0x00, 0x00},
			offset:     0,
			want:       0,
			wantOffset: 3,
		},
		{
			name:       "middle value 123456",
			input:      []byte{0x01, 0xE2, 0x40},
			offset:     0,
			want:       123456,
			wantOffset: 3,
		},
		{
			name:       "value 1",
			input:      []byte{0x00, 0x00, 0x01},
			offset:     0,
			want:       1,
			wantOffset: 3,
		},
		{
			name:       "with offset",
			input:      []byte{0xFF, 0xFF, 0x01, 0xE2, 0x40, 0xFF},
			offset:     2,
			want:       123456,
			wantOffset: 5,
		},
		{
			name:    "insufficient data at offset",
			input:   []byte{0x01, 0x02},
			offset:  0,
			wantErr: true,
		},
		{
			name:    "offset beyond buffer",
			input:   []byte{0x01, 0x02, 0x03},
			offset:  2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOffset, err := ReadOdometerFromBytes(tt.input, tt.offset)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ReadOdometerFromBytes() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadOdometerFromBytes() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ReadOdometerFromBytes() = %v, want %v", got, tt.want)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("ReadOdometerFromBytes() offset = %v, want %v", gotOffset, tt.wantOffset)
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
			odometer, offset, err := ReadOdometerFromBytes(tt.input, 0)
			if err != nil {
				t.Fatalf("ReadOdometerFromBytes() unexpected error: %v", err)
			}
			if offset != 3 {
				t.Errorf("ReadOdometerFromBytes() offset = %v, want 3", offset)
			}

			// Marshal
			dst := []byte{}
			got := AppendOdometer(dst, odometer)

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
