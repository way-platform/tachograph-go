package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalStringValue(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		wantEncoding ddv1.Encoding
		wantEncoded  []byte
		wantDecoded  string
		wantErr      bool
		errMessage   string
	}{
		{
			name:         "ISO-8859-1 with name 'John'",
			input:        []byte{0x01, 0x4A, 0x6F, 0x68, 0x6E}, // Code page 1 + "John"
			wantEncoding: ddv1.Encoding_ISO_8859_1,
			wantEncoded:  []byte{0x4A, 0x6F, 0x68, 0x6E},
			wantDecoded:  "John",
		},
		{
			name:         "Default code page with text",
			input:        []byte{0x00, 0x48, 0x65, 0x6C, 0x6C, 0x6F}, // Code page 0 + "Hello"
			wantEncoding: ddv1.Encoding_ENCODING_DEFAULT,
			wantEncoded:  []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F},
			wantDecoded:  "Hello",
		},
		{
			name:         "ISO-8859-2 Polish text",
			input:        []byte{0x02, 0xB1, 0xE6, 0xEA, 0xB3, 0xF1}, // Code page 2 + ISO-8859-2 bytes
			wantEncoding: ddv1.Encoding_ISO_8859_2,
			wantEncoded:  []byte{0xB1, 0xE6, 0xEA, 0xB3, 0xF1},
			wantDecoded:  "ąćęłń",
		},
		{
			name:         "Empty code page (255) with padding",
			input:        []byte{0xFF, 0x00, 0x00, 0x00},
			wantEncoding: ddv1.Encoding_ENCODING_EMPTY,
			wantEncoded:  []byte{0x00, 0x00, 0x00},
			wantDecoded:  "",
		},
		{
			name:         "ISO-8859-15 with Euro sign area",
			input:        []byte{0x0F, 0x45, 0x75, 0x72, 0x6F}, // Code page 15 + "Euro"
			wantEncoding: ddv1.Encoding_ISO_8859_15,
			wantEncoded:  []byte{0x45, 0x75, 0x72, 0x6F},
			wantDecoded:  "Euro",
		},
		{
			name:       "insufficient data - only code page",
			input:      []byte{0x01},
			wantErr:    true,
			errMessage: "insufficient data",
		},
		{
			name:       "insufficient data - empty",
			input:      []byte{},
			wantErr:    true,
			errMessage: "insufficient data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts UnmarshalOptions
			got, err := opts.UnmarshalStringValue(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalStringValue() expected error containing %q, got nil", tt.errMessage)
				} else if tt.errMessage != "" && !contains(err.Error(), tt.errMessage) {
					t.Errorf("UnmarshalStringValue() error = %v, want error containing %q", err, tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("UnmarshalStringValue() unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("UnmarshalStringValue() returned nil")
			}

			if got.GetEncoding() != tt.wantEncoding {
				t.Errorf("UnmarshalStringValue().GetEncoding() = %v, want %v", got.GetEncoding(), tt.wantEncoding)
			}

			if diff := cmp.Diff(tt.wantEncoded, got.GetRawData()); diff != "" {
				t.Errorf("UnmarshalStringValue().GetRawData() mismatch (-want +got):\n%s", diff)
			}

			if got.GetValue() != tt.wantDecoded {
				t.Errorf("UnmarshalStringValue().GetValue() = %q, want %q", got.GetValue(), tt.wantDecoded)
			}
		})
	}
}

func TestUnmarshalIA5StringValue(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantEncoded []byte
		wantDecoded string
		wantErr     bool
		errMessage  string
	}{
		{
			name:        "simple ASCII string",
			input:       []byte("Hello"),
			wantEncoded: []byte("Hello"),
			wantDecoded: "Hello",
		},
		{
			name:        "ASCII with trailing spaces",
			input:       []byte("Test    "),
			wantEncoded: []byte("Test    "),
			wantDecoded: "Test",
		},
		{
			name:        "ASCII with trailing zeros",
			input:       []byte{0x41, 0x42, 0x43, 0x00, 0x00},
			wantEncoded: []byte{0x41, 0x42, 0x43, 0x00, 0x00},
			wantDecoded: "ABC",
		},
		{
			name:        "ASCII with mixed padding",
			input:       []byte{0x54, 0x65, 0x73, 0x74, 0x20, 0x00, 0xFF},
			wantEncoded: []byte{0x54, 0x65, 0x73, 0x74, 0x20, 0x00, 0xFF},
			wantDecoded: "Test",
		},
		{
			name:        "14-char registration number",
			input:       []byte("FPA-829       "),
			wantEncoded: []byte("FPA-829       "),
			wantDecoded: "FPA-829",
		},
		{
			name:        "all spaces",
			input:       []byte("      "),
			wantEncoded: []byte("      "),
			wantDecoded: "",
		},
		{
			name:       "empty input",
			input:      []byte{},
			wantErr:    true,
			errMessage: "insufficient data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts UnmarshalOptions
			got, err := opts.UnmarshalIa5StringValue(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalIA5StringValue() expected error containing %q, got nil", tt.errMessage)
				} else if tt.errMessage != "" && !contains(err.Error(), tt.errMessage) {
					t.Errorf("UnmarshalIA5StringValue() error = %v, want error containing %q", err, tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("UnmarshalIA5StringValue() unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("UnmarshalIA5StringValue() returned nil")
			}

			// Ia5StringValue no longer has encoding field - it's always IA5/ASCII

			if diff := cmp.Diff(tt.wantEncoded, got.GetRawData()); diff != "" {
				t.Errorf("UnmarshalIA5StringValue().GetRawData() mismatch (-want +got):\n%s", diff)
			}

			if got.GetValue() != tt.wantDecoded {
				t.Errorf("UnmarshalIA5StringValue().GetValue() = %q, want %q", got.GetValue(), tt.wantDecoded)
			}
		})
	}
}

func TestAppendStringValue_CodePaged(t *testing.T) {
	tests := []struct {
		name        string
		stringValue *ddv1.StringValue
		want        []byte
	}{
		{
			name: "ISO-8859-1 with encoded bytes",
			stringValue: func() *ddv1.StringValue {
				sv := &ddv1.StringValue{}
				sv.SetEncoding(ddv1.Encoding_ISO_8859_1)
				sv.SetRawData([]byte{0x4A, 0x6F, 0x68, 0x6E})
				sv.SetValue("John")
				return sv
			}(),
			want: []byte{0x01, 0x4A, 0x6F, 0x68, 0x6E},
		},
		{
			name: "Default encoding from decoded string",
			stringValue: func() *ddv1.StringValue {
				sv := &ddv1.StringValue{}
				sv.SetEncoding(ddv1.Encoding_ENCODING_DEFAULT)
				sv.SetValue("Hello")
				return sv
			}(),
			want: []byte{0x00, 0x48, 0x65, 0x6C, 0x6C, 0x6F},
		},
		{
			name: "Empty string with default encoding",
			stringValue: func() *ddv1.StringValue {
				sv := &ddv1.StringValue{}
				sv.SetEncoding(ddv1.Encoding_ENCODING_DEFAULT)
				sv.SetValue("")
				return sv
			}(),
			want: []byte{0x00},
		},
		{
			name: "Empty encoding (255)",
			stringValue: func() *ddv1.StringValue {
				sv := &ddv1.StringValue{}
				sv.SetEncoding(ddv1.Encoding_ENCODING_EMPTY)
				sv.SetValue("")
				return sv
			}(),
			want: []byte{0xFF},
		},
		{
			name:        "nil string value for code-paged",
			stringValue: nil,
			want:        []byte{0xFF},
		},
		{
			name: "ISO-8859-15 with text",
			stringValue: func() *ddv1.StringValue {
				sv := &ddv1.StringValue{}
				sv.SetEncoding(ddv1.Encoding_ISO_8859_15)
				sv.SetValue("Test")
				return sv
			}(),
			want: []byte{0x0F, 0x54, 0x65, 0x73, 0x74},
		},
		{
			name: "Prefer encoded bytes over decoded",
			stringValue: func() *ddv1.StringValue {
				sv := &ddv1.StringValue{}
				sv.SetEncoding(ddv1.Encoding_ISO_8859_1)
				sv.SetRawData([]byte{0x41, 0x42, 0x43})
				sv.SetValue("XYZ") // Should be ignored
				return sv
			}(),
			want: []byte{0x01, 0x41, 0x42, 0x43},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendStringValue(nil, tt.stringValue)
			if err != nil {
				t.Fatalf("AppendStringValue() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendStringValue() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendIa5StringValue(t *testing.T) {
	tests := []struct {
		name        string
		stringValue *ddv1.Ia5StringValue
		length      int
		want        []byte
		wantErr     bool
	}{
		{
			name: "simple string with padding",
			stringValue: func() *ddv1.Ia5StringValue {
				sv := &ddv1.Ia5StringValue{}
				sv.SetLength(10)
				sv.SetValue("ABC")
				return sv
			}(),
			length: 10,
			want:   []byte{0x41, 0x42, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "prefer raw_data bytes if correct length",
			stringValue: func() *ddv1.Ia5StringValue {
				sv := &ddv1.Ia5StringValue{}
				sv.SetLength(5)
				sv.SetRawData([]byte{0x41, 0x42, 0x43, 0x20, 0x20})
				sv.SetValue("XYZ") // Should be ignored
				return sv
			}(),
			length: 5,
			want:   []byte{0x41, 0x42, 0x43, 0x20, 0x20},
		},
		{
			name: "error if raw_data length disagrees with length field",
			stringValue: func() *ddv1.Ia5StringValue {
				sv := &ddv1.Ia5StringValue{}
				sv.SetLength(10)
				sv.SetRawData([]byte{0x41, 0x42}) // Wrong length - should error
				sv.SetValue("Test")
				return sv
			}(),
			length:  10,
			wantErr: true,
		},
		{
			name:        "nil string value",
			stringValue: nil,
			length:      5,
			want:        []byte{}, // Nil returns empty (no code page for IA5)
		},
		{
			name: "full length string",
			stringValue: func() *ddv1.Ia5StringValue {
				sv := &ddv1.Ia5StringValue{}
				sv.SetLength(10)
				sv.SetValue("1234567890")
				return sv
			}(),
			length: 10,
			want:   []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
		},
		{
			name: "empty string with length",
			stringValue: func() *ddv1.Ia5StringValue {
				sv := &ddv1.Ia5StringValue{}
				sv.SetLength(3)
				sv.SetValue("")
				return sv
			}(),
			length: 3,
			want:   []byte{0x20, 0x20, 0x20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendIa5StringValue(nil, tt.stringValue)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("AppendStringValue() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendStringValue() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendStringValue() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStringValueRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "ISO-8859-1 with name",
			input: []byte{0x01, 0x4A, 0x6F, 0x68, 0x6E, 0x20, 0x44, 0x6F, 0x65},
		},
		{
			name:  "Default encoding with text",
			input: []byte{0x00, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64},
		},
		{
			name:  "ISO-8859-2 Polish text",
			input: []byte{0x02, 0xB1, 0xE6, 0xEA, 0xB3, 0xF1},
		},
		{
			name:  "Empty code page with padding",
			input: []byte{0xFF, 0x00, 0x00},
		},
		{
			name:  "ISO-8859-15",
			input: []byte{0x0F, 0x54, 0x65, 0x73, 0x74},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			sv, err := opts.UnmarshalStringValue(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalStringValue() error: %v", err)
			}

			if sv == nil {
				t.Fatal("UnmarshalStringValue() returned nil")
			}

			// Marshal back (code-paged format)
			got, err := AppendStringValue(nil, sv)
			if err != nil {
				t.Fatalf("AppendStringValue() error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-original +got):\n%s", diff)
			}
		})
	}
}

func TestIA5StringValueRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		length int
	}{
		{
			name:   "simple ASCII",
			input:  []byte("Hello     "),
			length: 10,
		},
		{
			name:   "registration number",
			input:  []byte("FPA-829       "),
			length: 14,
		},
		{
			name:   "full length",
			input:  []byte("1234567890"),
			length: 10,
		},
		{
			name:   "all spaces",
			input:  []byte("     "),
			length: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			sv, err := opts.UnmarshalIa5StringValue(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalIA5StringValue() error: %v", err)
			}

			if sv == nil {
				t.Fatal("UnmarshalIA5StringValue() returned nil")
			}

			// Marshal back
			got, err := AppendIa5StringValue(nil, sv)
			if err != nil {
				t.Fatalf("AppendIa5StringValue() error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-original +got):\n%s", diff)
			}
		})
	}
}

func TestAppendIa5StringValue_WithExistingData(t *testing.T) {
	existing := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	sv := &ddv1.Ia5StringValue{}
	// IA5 encoding is implicit for Ia5StringValue
	sv.SetLength(10)
	sv.SetValue("Test")

	got, err := AppendIa5StringValue(existing, sv)
	if err != nil {
		t.Fatalf("AppendIa5StringValue() unexpected error: %v", err)
	}

	want := []byte{
		0xDE, 0xAD, 0xBE, 0xEF, // existing data
		0x54, 0x65, 0x73, 0x74, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, // "Test" + padding
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AppendIa5StringValue() with existing data mismatch (-want +got):\n%s", diff)
	}
}
