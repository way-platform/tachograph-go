package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestBcdBytesToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int
		wantErr bool
	}{
		{
			name:  "valid multi-byte BCD",
			input: []byte{0x12, 0x34},
			want:  1234,
		},
		{
			name:  "valid single-byte BCD",
			input: []byte{0x56},
			want:  56,
		},
		{
			name:  "zero value",
			input: []byte{0x00, 0x00},
			want:  0,
		},
		{
			name:  "empty input",
			input: []byte{},
			want:  0,
		},
		{
			name:  "leading zeros",
			input: []byte{0x00, 0x42},
			want:  42,
		},
		{
			name:  "three-byte BCD",
			input: []byte{0x12, 0x34, 0x56},
			want:  123456,
		},
		{
			name:    "invalid BCD characters",
			input:   []byte{0xab},
			wantErr: true,
		},
		{
			name:    "partially invalid BCD",
			input:   []byte{0x12, 0xfa},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BcdBytesToInt(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("BcdBytesToInt() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("BcdBytesToInt() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("BcdBytesToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateBcdString(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantDecoded int32
		wantEncoded []byte
		wantErr     bool
	}{
		{
			name:        "valid BCD string",
			input:       []byte{0x12, 0x34},
			wantDecoded: 1234,
			wantEncoded: []byte{0x12, 0x34},
		},
		{
			name:        "single byte",
			input:       []byte{0x99},
			wantDecoded: 99,
			wantEncoded: []byte{0x99},
		},
		{
			name:        "zero value",
			input:       []byte{0x00},
			wantDecoded: 0,
			wantEncoded: []byte{0x00},
		},
		{
			name:    "invalid BCD",
			input:   []byte{0xab, 0xcd},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateBcdString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateBcdString() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("CreateBcdString() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("CreateBcdString() returned nil")
			}
			if got.GetDecoded() != tt.wantDecoded {
				t.Errorf("CreateBcdString().GetDecoded() = %v, want %v", got.GetDecoded(), tt.wantDecoded)
			}
			if diff := cmp.Diff(tt.wantEncoded, got.GetEncoded()); diff != "" {
				t.Errorf("CreateBcdString().GetEncoded() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUnmarshalBcdString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:  "valid BCD string",
			input: []byte{0x12, 0x34, 0x56},
		},
		{
			name:  "single byte",
			input: []byte{0x42},
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalBcdString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalBcdString() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalBcdString() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("UnmarshalBcdString() returned nil")
			}
			if diff := cmp.Diff(tt.input, got.GetEncoded()); diff != "" {
				t.Errorf("UnmarshalBcdString().GetEncoded() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendBcdString(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantBytes []byte
	}{
		{
			name:      "valid BCD string",
			input:     []byte{0x12, 0x34},
			wantBytes: []byte{0x12, 0x34},
		},
		{
			name:      "single byte",
			input:     []byte{0x99},
			wantBytes: []byte{0x99},
		},
		{
			name:      "empty BCD string",
			input:     []byte{},
			wantBytes: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create BCD string from input
			var bcdString *ddv1.BcdString
			if len(tt.input) > 0 {
				var err error
				bcdString, err = CreateBcdString(tt.input)
				if err != nil {
					t.Fatalf("CreateBcdString() unexpected error: %v", err)
				}
			}

			dst := []byte{}
			got, err := AppendBcdString(dst, bcdString)
			if err != nil {
				t.Fatalf("AppendBcdString() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.wantBytes, got); diff != "" {
				t.Errorf("AppendBcdString() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBcdStringRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "two bytes",
			input: []byte{0x12, 0x34},
		},
		{
			name:  "four bytes",
			input: []byte{0x20, 0x25, 0x09, 0x30},
		},
		{
			name:  "single byte",
			input: []byte{0x42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			bcdString, err := UnmarshalBcdString(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalBcdString() unexpected error: %v", err)
			}
			if bcdString == nil {
				t.Fatal("UnmarshalBcdString() returned nil")
			}

			// Marshal
			dst := []byte{}
			got, err := AppendBcdString(dst, bcdString)
			if err != nil {
				t.Fatalf("AppendBcdString() unexpected error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
