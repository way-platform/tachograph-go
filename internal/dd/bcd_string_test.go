package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalBcdString(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantDecoded int32
		wantEncoded []byte
		wantErr     bool
	}{
		{
			name:        "two-byte BCD",
			input:       []byte{0x12, 0x34},
			wantDecoded: 1234,
			wantEncoded: []byte{0x12, 0x34},
		},
		{
			name:        "three-byte BCD",
			input:       []byte{0x12, 0x34, 0x56},
			wantDecoded: 123456,
			wantEncoded: []byte{0x12, 0x34, 0x56},
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
			name:        "leading zeros",
			input:       []byte{0x00, 0x42},
			wantDecoded: 42,
			wantEncoded: []byte{0x00, 0x42},
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "invalid BCD characters",
			input:   []byte{0xab, 0xcd},
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
			if got.GetValue() != tt.wantDecoded {
				t.Errorf("UnmarshalBcdString().GetValue() = %v, want %v", got.GetValue(), tt.wantDecoded)
			}
			if diff := cmp.Diff(tt.wantEncoded, got.GetRawData()); diff != "" {
				t.Errorf("UnmarshalBcdString().GetRawData() mismatch (-want +got):\n%s", diff)
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
				bcdString, err = UnmarshalBcdString(tt.input)
				if err != nil {
					t.Fatalf("UnmarshalBcdString() unexpected error: %v", err)
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

func TestAppendBcdString_SemanticFallback(t *testing.T) {
	tests := []struct {
		name        string
		decoded     int32
		wantBytes   []byte
		wantErr     bool
		description string
	}{
		{
			name:        "single digit",
			decoded:     5,
			wantBytes:   []byte{0x05},
			description: "5 should encode to 0x05 (1 byte, minimal)",
		},
		{
			name:        "two digits",
			decoded:     42,
			wantBytes:   []byte{0x42},
			description: "42 should encode to 0x42 (1 byte)",
		},
		{
			name:        "three digits",
			decoded:     123,
			wantBytes:   []byte{0x01, 0x23},
			description: "123 should encode to 0x01 0x23 (2 bytes, minimal)",
		},
		{
			name:        "four digits",
			decoded:     1234,
			wantBytes:   []byte{0x12, 0x34},
			description: "1234 should encode to 0x12 0x34 (2 bytes)",
		},
		{
			name:        "five digits",
			decoded:     12345,
			wantBytes:   []byte{0x01, 0x23, 0x45},
			description: "12345 should encode to 0x01 0x23 0x45 (3 bytes, minimal)",
		},
		{
			name:        "zero",
			decoded:     0,
			wantBytes:   []byte{0x00},
			description: "0 should encode to 0x00 (1 byte)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create BCD string with only value (no raw_data)
			bcdString := &ddv1.BcdString{}
			bcdString.SetValue(tt.decoded)
			// Explicitly ensure raw_data is empty to test fallback
			// (protobuf will return empty bytes by default, but let's be explicit)

			dst := []byte{}
			got, err := AppendBcdString(dst, bcdString)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AppendBcdString() expected error for %s, got nil", tt.description)
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendBcdString() unexpected error for %s: %v", tt.description, err)
			}

			if diff := cmp.Diff(tt.wantBytes, got); diff != "" {
				t.Errorf("AppendBcdString() for %s mismatch (-want +got):\n%s", tt.description, diff)
			}

			// Verify it can be decoded back
			decoded, err := decodeBCD(got)
			if err != nil {
				t.Fatalf("decodeBCD() failed: %v", err)
			}
			if int32(decoded) != tt.decoded {
				t.Errorf("Decoded value = %d, want %d", decoded, tt.decoded)
			}
		})
	}
}
