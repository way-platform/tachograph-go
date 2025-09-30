package dd

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnmarshalTimeReal(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantUnix  int64
		wantIsNil bool
		wantErr   bool
	}{
		{
			name:     "2025-09-30 10:00:00 UTC",
			input:    []byte{0x68, 0xDB, 0xAA, 0x20},
			wantUnix: 1759226400,
		},
		{
			name:     "2024-01-01 00:00:00 UTC",
			input:    []byte{0x65, 0x92, 0x00, 0x80},
			wantUnix: 1704067200,
		},
		{
			name:     "Unix epoch + 1",
			input:    []byte{0x00, 0x00, 0x00, 0x01},
			wantUnix: 1,
		},
		{
			name:      "zero value",
			input:     []byte{0x00, 0x00, 0x00, 0x00},
			wantIsNil: true,
		},
		{
			name:     "2038-01-19 03:14:07 UTC (max int32)",
			input:    []byte{0x7F, 0xFF, 0xFF, 0xFF},
			wantUnix: 2147483647,
		},
		{
			name:     "value beyond int32 max",
			input:    []byte{0x80, 0x00, 0x00, 0x00},
			wantUnix: 2147483648,
		},
		{
			name:    "insufficient data (3 bytes)",
			input:   []byte{0x68, 0xDB, 0xAA},
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
			got, err := opts.UnmarshalTimeReal(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalTimeReal() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalTimeReal() unexpected error: %v", err)
			}

			if tt.wantIsNil {
				if got != nil {
					t.Errorf("UnmarshalTimeReal() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("UnmarshalTimeReal() returned nil")
			}
			if got.GetSeconds() != tt.wantUnix {
				t.Errorf("UnmarshalTimeReal().GetSeconds() = %v, want %v", got.GetSeconds(), tt.wantUnix)
			}
		})
	}
}

func TestAppendTimeReal(t *testing.T) {
	tests := []struct {
		name      string
		timestamp *timestamppb.Timestamp
		want      []byte
	}{
		{
			name:      "2025-09-30 10:00:00 UTC",
			timestamp: timestamppb.New(time.Unix(1759226400, 0)),
			want:      []byte{0x68, 0xDB, 0xAA, 0x20},
		},
		{
			name:      "2024-01-01 00:00:00 UTC",
			timestamp: timestamppb.New(time.Unix(1704067200, 0)),
			want:      []byte{0x65, 0x92, 0x00, 0x80},
		},
		{
			name:      "Unix epoch + 1",
			timestamp: timestamppb.New(time.Unix(1, 0)),
			want:      []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			name:      "nil timestamp",
			timestamp: nil,
			want:      []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:      "2038-01-19 03:14:07 UTC",
			timestamp: timestamppb.New(time.Unix(2147483647, 0)),
			want:      []byte{0x7F, 0xFF, 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte{}
			got, err := AppendTimeReal(dst, tt.timestamp)
			if err != nil {
				t.Fatalf("AppendTimeReal() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendTimeReal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTimeRealRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "2025-09-30",
			input: []byte{0x68, 0xDB, 0xAA, 0x20},
		},
		{
			name:  "2024-01-01",
			input: []byte{0x65, 0x92, 0x00, 0x80},
		},
		{
			name:  "epoch + 1",
			input: []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			name:  "max int32",
			input: []byte{0x7F, 0xFF, 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			ts, err := opts.UnmarshalTimeReal(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalTimeReal() unexpected error: %v", err)
			}
			if ts == nil {
				t.Fatal("UnmarshalTimeReal() returned nil")
			}

			// Marshal
			dst := []byte{}
			got, err := AppendTimeReal(dst, ts)
			if err != nil {
				t.Fatalf("AppendTimeReal() error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
