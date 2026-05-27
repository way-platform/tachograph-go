package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestUnmarshalVehicleRegistrationIdentification_Nation guards against a
// regression where the nation byte was cast directly to the NationNumeric enum
// (data[0] interpreted as the proto enum value) instead of being mapped through
// the protocol_enum_value annotation. With the raw cast, a Norwegian record
// (protocol value 37) decoded as MONACO (proto enum value 37). See tacho#1090.
func TestUnmarshalVehicleRegistrationIdentification_Nation(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		wantNation    ddv1.NationNumeric
		wantRegNumber string
	}{
		{
			// Norway protocol value = 37 (0x25). The raw-cast bug decoded this
			// as MONACO (proto enum value 37). Real plate from tacho#1090.
			name:          "Norway with registration JD92527",
			input:         []byte{0x25, 0x00, 0x4A, 0x44, 0x39, 0x32, 0x35, 0x32, 0x37, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_NORWAY,
			wantRegNumber: "JD92527",
		},
		{
			// Monaco protocol value = 34 (0x22), distinct from Norway.
			name:          "Monaco with registration MC1234",
			input:         []byte{0x22, 0x00, 0x4D, 0x43, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_MONACO,
			wantRegNumber: "MC1234",
		},
		{
			name:          "Unrecognized nation code (100)",
			input:         []byte{0x64, 0x00, 0x54, 0x45, 0x53, 0x54, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED,
			wantRegNumber: "TEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalOpts UnmarshalOptions
			got, err := unmarshalOpts.UnmarshalVehicleRegistrationIdentification(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalVehicleRegistrationIdentification() unexpected error: %v", err)
			}
			if got.GetNation() != tt.wantNation {
				t.Errorf("GetNation() = %v, want %v", got.GetNation(), tt.wantNation)
			}
			gotRegNumber := ""
			if got.GetNumber() != nil {
				gotRegNumber = got.GetNumber().GetValue()
			}
			if gotRegNumber != tt.wantRegNumber {
				t.Errorf("GetNumber().GetValue() = %q, want %q", gotRegNumber, tt.wantRegNumber)
			}
		})
	}
}

// TestVehicleRegistrationIdentificationRoundTrip verifies that recognised
// nation codes survive an unmarshal → marshal round-trip byte-for-byte.
func TestVehicleRegistrationIdentificationRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "Norway with registration JD92527",
			input: []byte{0x25, 0x00, 0x4A, 0x44, 0x39, 0x32, 0x35, 0x32, 0x37, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "Monaco with registration MC1234",
			input: []byte{0x22, 0x00, 0x4D, 0x43, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalOpts UnmarshalOptions
			var marshalOpts MarshalOptions
			vri, err := unmarshalOpts.UnmarshalVehicleRegistrationIdentification(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalVehicleRegistrationIdentification() error: %v", err)
			}
			got, err := marshalOpts.MarshalVehicleRegistrationIdentification(vri)
			if err != nil {
				t.Fatalf("MarshalVehicleRegistrationIdentification() error: %v", err)
			}
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-original +got):\n%s", diff)
			}
		})
	}
}
