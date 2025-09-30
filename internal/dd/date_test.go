package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalDate(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantYear    int32
		wantMonth   int32
		wantDay     int32
		wantEncoded []byte
		wantErr     bool
	}{
		{
			name:        "valid date 2025-09-30",
			input:       []byte{0x20, 0x25, 0x09, 0x30},
			wantYear:    2025,
			wantMonth:   9,
			wantDay:     30,
			wantEncoded: []byte{0x20, 0x25, 0x09, 0x30},
		},
		{
			name:        "valid date 2024-01-01",
			input:       []byte{0x20, 0x24, 0x01, 0x01},
			wantYear:    2024,
			wantMonth:   1,
			wantDay:     1,
			wantEncoded: []byte{0x20, 0x24, 0x01, 0x01},
		},
		{
			name:        "valid date 2023-12-31",
			input:       []byte{0x20, 0x23, 0x12, 0x31},
			wantYear:    2023,
			wantMonth:   12,
			wantDay:     31,
			wantEncoded: []byte{0x20, 0x23, 0x12, 0x31},
		},
		{
			name:        "valid date 1999-05-15",
			input:       []byte{0x19, 0x99, 0x05, 0x15},
			wantYear:    1999,
			wantMonth:   5,
			wantDay:     15,
			wantEncoded: []byte{0x19, 0x99, 0x05, 0x15},
		},
		{
			name:        "zero date",
			input:       []byte{0x00, 0x00, 0x00, 0x00},
			wantYear:    0,
			wantMonth:   0,
			wantDay:     0,
			wantEncoded: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:        "month 13 (invalid but parsed)",
			input:       []byte{0x20, 0x25, 0x13, 0x01},
			wantYear:    2025,
			wantMonth:   13,
			wantDay:     1,
			wantEncoded: []byte{0x20, 0x25, 0x13, 0x01},
		},
		{
			name:        "day 32 (invalid but parsed)",
			input:       []byte{0x20, 0x25, 0x01, 0x32},
			wantYear:    2025,
			wantMonth:   1,
			wantDay:     32,
			wantEncoded: []byte{0x20, 0x25, 0x01, 0x32},
		},
		{
			name:        "month 00 (invalid but parsed)",
			input:       []byte{0x20, 0x25, 0x00, 0x15},
			wantYear:    2025,
			wantMonth:   0,
			wantDay:     15,
			wantEncoded: []byte{0x20, 0x25, 0x00, 0x15},
		},
		{
			name:        "day 00 (invalid but parsed)",
			input:       []byte{0x20, 0x25, 0x05, 0x00},
			wantYear:    2025,
			wantMonth:   5,
			wantDay:     0,
			wantEncoded: []byte{0x20, 0x25, 0x05, 0x00},
		},
		{
			name:        "year 1899",
			input:       []byte{0x18, 0x99, 0x12, 0x31},
			wantYear:    1899,
			wantMonth:   12,
			wantDay:     31,
			wantEncoded: []byte{0x18, 0x99, 0x12, 0x31},
		},
		{
			name:    "insufficient data",
			input:   []byte{0x20, 0x25},
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts UnmarshalOptions
			got, err := opts.UnmarshalDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalDate() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalDate() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("UnmarshalDate() returned nil")
			}
			if got.GetYear() != tt.wantYear {
				t.Errorf("GetYear() = %v, want %v", got.GetYear(), tt.wantYear)
			}
			if got.GetMonth() != tt.wantMonth {
				t.Errorf("GetMonth() = %v, want %v", got.GetMonth(), tt.wantMonth)
			}
			if got.GetDay() != tt.wantDay {
				t.Errorf("GetDay() = %v, want %v", got.GetDay(), tt.wantDay)
			}
			if diff := cmp.Diff(tt.wantEncoded, got.GetRawData()); diff != "" {
				t.Errorf("GetEncoded() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendDate(t *testing.T) {
	tests := []struct {
		name string
		date *ddv1.Date
		want []byte
	}{
		{
			name: "valid date with encoded bytes",
			date: func() *ddv1.Date {
				d := &ddv1.Date{}
				d.SetRawData([]byte{0x20, 0x25, 0x09, 0x30})
				d.SetYear(2025)
				d.SetMonth(9)
				d.SetDay(30)
				return d
			}(),
			want: []byte{0x20, 0x25, 0x09, 0x30},
		},
		{
			name: "valid date from decoded values only",
			date: func() *ddv1.Date {
				d := &ddv1.Date{}
				d.SetYear(2024)
				d.SetMonth(1)
				d.SetDay(1)
				return d
			}(),
			want: []byte{0x20, 0x24, 0x01, 0x01},
		},
		{
			name: "date 2023-12-31",
			date: func() *ddv1.Date {
				d := &ddv1.Date{}
				d.SetYear(2023)
				d.SetMonth(12)
				d.SetDay(31)
				return d
			}(),
			want: []byte{0x20, 0x23, 0x12, 0x31},
		},
		{
			name: "date 1999-05-15",
			date: func() *ddv1.Date {
				d := &ddv1.Date{}
				d.SetYear(1999)
				d.SetMonth(5)
				d.SetDay(15)
				return d
			}(),
			want: []byte{0x19, 0x99, 0x05, 0x15},
		},
		{
			name: "nil date",
			date: nil,
			want: []byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte{}
			got, err := AppendDate(dst, tt.date)
			if err != nil {
				t.Fatalf("AppendDate() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendDate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDateRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "date 2025-09-30",
			input: []byte{0x20, 0x25, 0x09, 0x30},
		},
		{
			name:  "date 2024-01-01",
			input: []byte{0x20, 0x24, 0x01, 0x01},
		},
		{
			name:  "date 2023-12-31",
			input: []byte{0x20, 0x23, 0x12, 0x31},
		},
		{
			name:  "date 1999-05-15",
			input: []byte{0x19, 0x99, 0x05, 0x15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			date, err := opts.UnmarshalDate(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalDate() unexpected error: %v", err)
			}
			if date == nil {
				t.Fatal("UnmarshalDate() returned nil")
			}

			// Marshal
			dst := []byte{}
			got, err := AppendDate(dst, date)
			if err != nil {
				t.Fatalf("AppendDate() error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
