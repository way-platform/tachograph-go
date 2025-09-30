package dd

import (
	"testing"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalControlType(t *testing.T) {
	tests := []struct {
		name                    string
		input                   []byte
		wantCardDownloading     bool
		wantVuDownloading       bool
		wantPrinting            bool
		wantDisplay             bool
		wantCalibrationChecking bool
		wantErr                 bool
	}{
		{
			name:                "card downloading only",
			input:               []byte{0x80},
			wantCardDownloading: true,
		},
		{
			name:              "VU downloading only",
			input:             []byte{0x40},
			wantVuDownloading: true,
		},
		{
			name:         "printing only",
			input:        []byte{0x20},
			wantPrinting: true,
		},
		{
			name:        "display only",
			input:       []byte{0x10},
			wantDisplay: true,
		},
		{
			name:                    "calibration checking only",
			input:                   []byte{0x08},
			wantCalibrationChecking: true,
		},
		{
			name:                    "all flags set",
			input:                   []byte{0xF8},
			wantCardDownloading:     true,
			wantVuDownloading:       true,
			wantPrinting:            true,
			wantDisplay:             true,
			wantCalibrationChecking: true,
		},
		{
			name:  "no flags set",
			input: []byte{0x00},
		},
		{
			name:                "card and VU downloading",
			input:               []byte{0xC0},
			wantCardDownloading: true,
			wantVuDownloading:   true,
		},
		{
			name:                    "display and calibration",
			input:                   []byte{0x18},
			wantDisplay:             true,
			wantCalibrationChecking: true,
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalControlType(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalControlType() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalControlType() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("UnmarshalControlType() returned nil")
			}
			if got.GetCardDownloading() != tt.wantCardDownloading {
				t.Errorf("GetCardDownloading() = %v, want %v", got.GetCardDownloading(), tt.wantCardDownloading)
			}
			if got.GetVuDownloading() != tt.wantVuDownloading {
				t.Errorf("GetVuDownloading() = %v, want %v", got.GetVuDownloading(), tt.wantVuDownloading)
			}
			if got.GetPrinting() != tt.wantPrinting {
				t.Errorf("GetPrinting() = %v, want %v", got.GetPrinting(), tt.wantPrinting)
			}
			if got.GetDisplay() != tt.wantDisplay {
				t.Errorf("GetDisplay() = %v, want %v", got.GetDisplay(), tt.wantDisplay)
			}
			if got.GetCalibrationChecking() != tt.wantCalibrationChecking {
				t.Errorf("GetCalibrationChecking() = %v, want %v", got.GetCalibrationChecking(), tt.wantCalibrationChecking)
			}
		})
	}
}

func TestControlTypeRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "all flags set",
			input: []byte{0xF8},
		},
		{
			name:  "no flags set",
			input: []byte{0x00},
		},
		{
			name:  "card downloading only",
			input: []byte{0x80},
		},
		{
			name:  "mixed flags",
			input: []byte{0xA8},
		},
		{
			name:  "reserved bits set",
			input: []byte{0xF8 | 0x07}, // All flags + all reserved bits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			ct, err := UnmarshalControlType(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalControlType() unexpected error: %v", err)
			}
			if ct == nil {
				t.Fatal("UnmarshalControlType() returned nil")
			}

			// Marshal back
			var dst []byte
			dst, err = AppendControlType(dst, ct)
			if err != nil {
				t.Fatalf("AppendControlType() unexpected error: %v", err)
			}

			// Compare
			if len(dst) != 1 {
				t.Fatalf("AppendControlType() returned %d bytes, want 1", len(dst))
			}
			if dst[0] != tt.input[0] {
				t.Errorf("Round-trip failed: got 0x%02X, want 0x%02X", dst[0], tt.input[0])
			}
		})
	}
}

func TestAppendControlType_PreservesReservedBits(t *testing.T) {
	tests := []struct {
		name         string
		rawData      []byte
		setFlags     func(*ddv1.ControlType)
		wantReserved byte // Expected reserved bits (0-2)
	}{
		{
			name:    "preserves reserved bits 0-2 when set",
			rawData: []byte{0xF8 | 0x07}, // All flags + reserved bits 111
			setFlags: func(ct *ddv1.ControlType) {
				// Don't change any flags
			},
			wantReserved: 0x07,
		},
		{
			name:    "preserves reserved bit 0",
			rawData: []byte{0x80 | 0x01}, // Card downloading + reserved bit 0
			setFlags: func(ct *ddv1.ControlType) {
				ct.SetPrinting(true) // Add printing flag
			},
			wantReserved: 0x01,
		},
		{
			name:    "preserves mixed reserved bits",
			rawData: []byte{0x40 | 0x05}, // VU downloading + reserved bits 101
			setFlags: func(ct *ddv1.ControlType) {
				ct.SetDisplay(true) // Add display flag
			},
			wantReserved: 0x05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal (this sets raw_data)
			ct, err := UnmarshalControlType(tt.rawData)
			if err != nil {
				t.Fatalf("UnmarshalControlType() unexpected error: %v", err)
			}

			// Apply flag changes
			tt.setFlags(ct)

			// Marshal back
			var dst []byte
			dst, err = AppendControlType(dst, ct)
			if err != nil {
				t.Fatalf("AppendControlType() unexpected error: %v", err)
			}

			// Check that reserved bits are preserved
			gotReserved := dst[0] & 0x07
			if gotReserved != tt.wantReserved {
				t.Errorf("Reserved bits not preserved: got 0x%02X, want 0x%02X", gotReserved, tt.wantReserved)
			}
		})
	}
}
