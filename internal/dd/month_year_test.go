package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalMonthYear(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantMonth  int32
		wantYear   int32
		wantErr    bool
		errMessage string
	}{
		{
			name:      "January 2025 (0125)",
			input:     []byte{0x01, 0x25},
			wantMonth: 1,
			wantYear:  2025,
		},
		{
			name:      "December 2024 (1224)",
			input:     []byte{0x12, 0x24},
			wantMonth: 12,
			wantYear:  2024,
		},
		{
			name:      "September 2023 (0923)",
			input:     []byte{0x09, 0x23},
			wantMonth: 9,
			wantYear:  2023,
		},
		{
			name:      "June 2000 (0600)",
			input:     []byte{0x06, 0x00},
			wantMonth: 6,
			wantYear:  2000,
		},
		{
			name:      "December 2099 (1299)",
			input:     []byte{0x12, 0x99},
			wantMonth: 12,
			wantYear:  2099,
		},
		{
			name:      "February 2020 (0220)",
			input:     []byte{0x02, 0x20},
			wantMonth: 2,
			wantYear:  2020,
		},
		{
			name:      "October 2015 (1015)",
			input:     []byte{0x10, 0x15},
			wantMonth: 10,
			wantYear:  2015,
		},
		{
			name:      "zero value (0000) - should parse but with zero values",
			input:     []byte{0x00, 0x00},
			wantMonth: 0,
			wantYear:  0,
		},
		{
			name:       "insufficient data - 1 byte",
			input:      []byte{0x01},
			wantErr:    true,
			errMessage: "invalid data length for MonthYear",
		},
		{
			name:       "insufficient data - 0 bytes",
			input:      []byte{},
			wantErr:    true,
			errMessage: "invalid data length for MonthYear",
		},
		{
			name:      "invalid BCD (0xAB, 0xCD) - should still preserve encoded bytes",
			input:     []byte{0xAB, 0xCD},
			wantMonth: 0, // Will fail to decode but still preserve encoded
			wantYear:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts UnmarshalOptions
			got, err := opts.UnmarshalMonthYear(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalMonthYear() expected error containing %q, got nil", tt.errMessage)
				} else if tt.errMessage != "" && !contains(err.Error(), tt.errMessage) {
					t.Errorf("UnmarshalMonthYear() error = %v, want error containing %q", err, tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("UnmarshalMonthYear() unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("UnmarshalMonthYear() returned nil")
			}

			if got.GetMonth() != tt.wantMonth {
				t.Errorf("UnmarshalMonthYear().GetMonth() = %v, want %v", got.GetMonth(), tt.wantMonth)
			}

			if got.GetYear() != tt.wantYear {
				t.Errorf("UnmarshalMonthYear().GetYear() = %v, want %v", got.GetYear(), tt.wantYear)
			}

			// Verify encoded bytes are preserved
			if !tt.wantErr && len(tt.input) >= 2 {
				if diff := cmp.Diff(tt.input[:2], got.GetRawData()); diff != "" {
					t.Errorf("UnmarshalMonthYear().GetRawData() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestAppendMonthYear(t *testing.T) {
	tests := []struct {
		name      string
		monthYear *ddv1.MonthYear
		want      []byte
	}{
		{
			name: "January 2025",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(1)
				my.SetYear(2025)
				return my
			}(),
			want: []byte{0x01, 0x25},
		},
		{
			name: "December 2024",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(12)
				my.SetYear(2024)
				return my
			}(),
			want: []byte{0x12, 0x24},
		},
		{
			name: "September 2023",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(9)
				my.SetYear(2023)
				return my
			}(),
			want: []byte{0x09, 0x23},
		},
		{
			name: "June 2000",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(6)
				my.SetYear(2000)
				return my
			}(),
			want: []byte{0x06, 0x00},
		},
		{
			name: "December 2099",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(12)
				my.SetYear(2099)
				return my
			}(),
			want: []byte{0x12, 0x99},
		},
		{
			name: "October 2015",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(10)
				my.SetYear(2015)
				return my
			}(),
			want: []byte{0x10, 0x15},
		},
		{
			name: "with raw_data - should paint semantic values over canvas",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetRawData([]byte{0x03, 0x22})
				my.SetMonth(3)   // Matches raw_data (0x03)
				my.SetYear(2022) // Matches raw_data (0x22 = 22 + 2000)
				return my
			}(),
			want: []byte{0x03, 0x22}, // Semantic values painted over raw_data canvas
		},
		{
			name:      "nil monthYear",
			monthYear: nil,
			want:      []byte{0x00, 0x00},
		},
		{
			name: "zero values",
			monthYear: func() *ddv1.MonthYear {
				my := &ddv1.MonthYear{}
				my.SetMonth(0)
				my.SetYear(0)
				return my
			}(),
			want: []byte{0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendMonthYear(nil, tt.monthYear)
			if err != nil {
				t.Fatalf("AppendMonthYear() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendMonthYear() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendMonthYear_WithExistingData(t *testing.T) {
	existing := []byte{0xFF, 0xEE}

	my := &ddv1.MonthYear{}
	my.SetMonth(6)
	my.SetYear(2023)

	got, err := AppendMonthYear(existing, my)
	if err != nil {
		t.Fatalf("AppendMonthYear() unexpected error: %v", err)
	}

	want := []byte{0xFF, 0xEE, 0x06, 0x23}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AppendMonthYear() with existing data mismatch (-want +got):\n%s", diff)
	}
}

func TestMonthYearRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "January 2025",
			input: []byte{0x01, 0x25},
		},
		{
			name:  "December 2024",
			input: []byte{0x12, 0x24},
		},
		{
			name:  "September 2023",
			input: []byte{0x09, 0x23},
		},
		{
			name:  "June 2000",
			input: []byte{0x06, 0x00},
		},
		{
			name:  "December 2099",
			input: []byte{0x12, 0x99},
		},
		{
			name:  "zero value",
			input: []byte{0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			monthYear, err := opts.UnmarshalMonthYear(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalMonthYear() error: %v", err)
			}

			if monthYear == nil {
				t.Fatal("UnmarshalMonthYear() returned nil")
			}

			// Marshal back
			got, err := AppendMonthYear(nil, monthYear)
			if err != nil {
				t.Fatalf("AppendMonthYear() error: %v", err)
			}

			// Should match original input
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-original +got):\n%s", diff)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
