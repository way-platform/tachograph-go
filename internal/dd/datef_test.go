package dd

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestReadDatef(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantYear  int32
		wantMonth int32
		wantDay   int32
		wantNil   bool
	}{
		{
			name:      "valid date 2025-09-30",
			input:     []byte{0x20, 0x25, 0x09, 0x30},
			wantYear:  2025,
			wantMonth: 9,
			wantDay:   30,
		},
		{
			name:      "valid date 2024-01-01",
			input:     []byte{0x20, 0x24, 0x01, 0x01},
			wantYear:  2024,
			wantMonth: 1,
			wantDay:   1,
		},
		{
			name:      "valid date 2023-12-31",
			input:     []byte{0x20, 0x23, 0x12, 0x31},
			wantYear:  2023,
			wantMonth: 12,
			wantDay:   31,
		},
		{
			name:    "zero date",
			input:   []byte{0x00, 0x00, 0x00, 0x00},
			wantNil: true,
		},
		{
			name:    "invalid month 13",
			input:   []byte{0x20, 0x25, 0x13, 0x01},
			wantNil: true,
		},
		{
			name:    "invalid day 32",
			input:   []byte{0x20, 0x25, 0x01, 0x32},
			wantNil: true,
		},
		{
			name:    "invalid month 00",
			input:   []byte{0x20, 0x25, 0x00, 0x15},
			wantNil: true,
		},
		{
			name:    "year too old",
			input:   []byte{0x18, 0x99, 0x12, 0x31},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			got, err := ReadDatef(r)
			if err != nil {
				t.Fatalf("ReadDatef() unexpected error: %v", err)
			}

			if tt.wantNil {
				if got != nil {
					t.Errorf("ReadDatef() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("ReadDatef() returned nil")
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
		})
	}
}

func TestAppendDatef(t *testing.T) {
	tests := []struct {
		name      string
		timestamp *timestamppb.Timestamp
		want      []byte
	}{
		{
			name:      "valid date 2025-09-30",
			timestamp: timestamppb.New(time.Date(2025, 9, 30, 10, 0, 0, 0, time.UTC)),
			want:      []byte{0x20, 0x25, 0x09, 0x30},
		},
		{
			name:      "valid date 2024-01-01",
			timestamp: timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			want:      []byte{0x20, 0x24, 0x01, 0x01},
		},
		{
			name:      "valid date 2023-12-31",
			timestamp: timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
			want:      []byte{0x20, 0x23, 0x12, 0x31},
		},
		{
			name:      "valid date 1999-05-15",
			timestamp: timestamppb.New(time.Date(1999, 5, 15, 12, 30, 0, 0, time.UTC)),
			want:      []byte{0x19, 0x99, 0x05, 0x15},
		},
		{
			name:      "nil timestamp",
			timestamp: nil,
			want:      []byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte{}
			got := AppendDatef(dst, tt.timestamp)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendDatef() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDatefRoundTrip(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			r := bytes.NewReader(tt.input)
			date, err := ReadDatef(r)
			if err != nil {
				t.Fatalf("ReadDatef() unexpected error: %v", err)
			}
			if date == nil {
				t.Fatal("ReadDatef() returned nil")
			}

			// Convert Date to Timestamp for marshalling
			ts := timestamppb.New(time.Date(
				int(date.GetYear()),
				time.Month(date.GetMonth()),
				int(date.GetDay()),
				0, 0, 0, 0, time.UTC,
			))

			// Marshal
			dst := []byte{}
			got := AppendDatef(dst, ts)

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
