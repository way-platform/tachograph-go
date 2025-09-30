package dd

import (
	"testing"
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

			// Note: The AppendControlType function is not complete in the source file.
			// We can only test unmarshalling for now.
			// When AppendControlType is implemented, add marshalling test here.
		})
	}
}
