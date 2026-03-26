package security

import "testing"

func TestUnmarshalTimeRealSentinels(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "zero value",
			input: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "all ones sentinel",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalTimeReal(tt.input)
			if err != nil {
				t.Fatalf("unmarshalTimeReal() unexpected error: %v", err)
			}
			if got != nil {
				t.Fatalf("unmarshalTimeReal() = %v, want nil", got)
			}
		})
	}
}
