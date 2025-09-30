package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalVehicleRegistration(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		wantNation    ddv1.NationNumeric
		wantRegNumber string
		wantErr       bool
		errMessage    string
	}{
		{
			name:          "Finland with registration FPA-829",
			input:         []byte{0x12, 0x46, 0x50, 0x41, 0x2D, 0x38, 0x32, 0x39, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_FINLAND,
			wantRegNumber: "FPA-829",
		},
		{
			name:          "Spain with registration ABC-1234",
			input:         []byte{0x0F, 0x41, 0x42, 0x43, 0x2D, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_SPAIN,
			wantRegNumber: "ABC-1234",
		},
		{
			name:          "Germany with registration B-MW-1234",
			input:         []byte{0x0D, 0x42, 0x2D, 0x4D, 0x57, 0x2D, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_GERMANY,
			wantRegNumber: "B-MW-1234",
		},
		{
			name:          "France with registration 75ABC12",
			input:         []byte{0x11, 0x37, 0x35, 0x41, 0x42, 0x43, 0x31, 0x32, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_FRANCE,
			wantRegNumber: "75ABC12",
		},
		{
			name:          "Italy with registration AB123CD",
			input:         []byte{0x1A, 0x41, 0x42, 0x31, 0x32, 0x33, 0x43, 0x44, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_ITALY,
			wantRegNumber: "AB123CD",
		},
		{
			name:          "Sweden with full 14-char registration",
			input:         []byte{0x2C, 0x41, 0x42, 0x43, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x58},
			wantNation:    ddv1.NationNumeric_SWEDEN,
			wantRegNumber: "ABC1234567890X",
		},
		{
			name:          "Empty nation (0xFF) with spaces",
			input:         []byte{0xFF, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_NATION_NUMERIC_EMPTY,
			wantRegNumber: "",
		},
		{
			name:          "Default nation (0x00) with short registration",
			input:         []byte{0x00, 0x41, 0x42, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_NATION_NUMERIC_DEFAULT,
			wantRegNumber: "ABC",
		},
		{
			name:          "Unrecognized nation code (100)",
			input:         []byte{0x64, 0x54, 0x45, 0x53, 0x54, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantNation:    ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED,
			wantRegNumber: "TEST",
		},
		{
			name:       "insufficient data - 14 bytes",
			input:      []byte{0x12, 0x46, 0x50, 0x41, 0x2D, 0x38, 0x32, 0x39, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
			wantErr:    true,
			errMessage: "invalid data length for VehicleRegistrationIdentification",
		},
		{
			name:       "insufficient data - 10 bytes",
			input:      []byte{0x12, 0x46, 0x50, 0x41, 0x2D, 0x38, 0x32, 0x39, 0x20, 0x20},
			wantErr:    true,
			errMessage: "invalid data length for VehicleRegistrationIdentification",
		},
		{
			name:       "insufficient data - 0 bytes",
			input:      []byte{},
			wantErr:    true,
			errMessage: "invalid data length for VehicleRegistrationIdentification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalVehicleRegistration(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalVehicleRegistration() expected error containing %q, got nil", tt.errMessage)
				} else if tt.errMessage != "" && !contains(err.Error(), tt.errMessage) {
					t.Errorf("UnmarshalVehicleRegistration() error = %v, want error containing %q", err, tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("UnmarshalVehicleRegistration() unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("UnmarshalVehicleRegistration() returned nil")
			}

			if got.GetNation() != tt.wantNation {
				t.Errorf("UnmarshalVehicleRegistration().GetNation() = %v, want %v", got.GetNation(), tt.wantNation)
			}

			gotRegNumber := ""
			if got.GetNumber() != nil {
				gotRegNumber = got.GetNumber().GetValue()
			}
			if gotRegNumber != tt.wantRegNumber {
				t.Errorf("UnmarshalVehicleRegistration().GetNumber().GetValue() = %q, want %q", gotRegNumber, tt.wantRegNumber)
			}
		})
	}
}

func TestAppendVehicleRegistration(t *testing.T) {
	tests := []struct {
		name       string
		vehicleReg *ddv1.VehicleRegistrationIdentification
		want       []byte
		wantErr    bool
	}{
		{
			name: "Finland with registration FPA-829",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_FINLAND)
				num := &ddv1.StringValue{}
				num.SetValue("FPA-829")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0x12, 0x46, 0x50, 0x41, 0x2D, 0x38, 0x32, 0x39, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "Spain with registration ABC-1234",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_SPAIN)
				num := &ddv1.StringValue{}
				num.SetValue("ABC-1234")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0x0F, 0x41, 0x42, 0x43, 0x2D, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "Germany with short registration ABC",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_GERMANY)
				num := &ddv1.StringValue{}
				num.SetValue("ABC")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0x0D, 0x41, 0x42, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "Italy with empty registration",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_ITALY)
				num := &ddv1.StringValue{}
				num.SetValue("")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0x1A, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "France with nil number",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_FRANCE)
				return vr
			}(),
			want: []byte{0x11, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "Empty nation with registration",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_NATION_NUMERIC_EMPTY)
				num := &ddv1.StringValue{}
				num.SetValue("TEST")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0xFF, 0x54, 0x45, 0x53, 0x54, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "Unspecified nation",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED)
				num := &ddv1.StringValue{}
				num.SetValue("ABC")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0xFF, 0x41, 0x42, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name: "Unrecognized nation",
			vehicleReg: func() *ddv1.VehicleRegistrationIdentification {
				vr := &ddv1.VehicleRegistrationIdentification{}
				vr.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
				num := &ddv1.StringValue{}
				num.SetValue("TEST")
				vr.SetNumber(num)
				return vr
			}(),
			want: []byte{0xFF, 0x54, 0x45, 0x53, 0x54, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:       "nil vehicle registration",
			vehicleReg: nil,
			wantErr:    true, // Errors because it calls nested AppendStringValue
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppendVehicleRegistration(nil, tt.vehicleReg)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AppendVehicleRegistration() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendVehicleRegistration() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendVehicleRegistration() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendVehicleRegistration_WithExistingData(t *testing.T) {
	existing := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	vr := &ddv1.VehicleRegistrationIdentification{}
	vr.SetNation(ddv1.NationNumeric_SWEDEN)
	num := &ddv1.StringValue{}
	num.SetValue("ABC123")
	vr.SetNumber(num)

	got, err := AppendVehicleRegistration(existing, vr)
	if err != nil {
		t.Fatalf("AppendVehicleRegistration() unexpected error: %v", err)
	}

	want := []byte{
		0xDE, 0xAD, 0xBE, 0xEF, // existing data
		0x2C,                               // Sweden (protocol value 44)
		0x41, 0x42, 0x43, 0x31, 0x32, 0x33, // "ABC123"
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, // padding
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AppendVehicleRegistration() with existing data mismatch (-want +got):\n%s", diff)
	}
}

func TestVehicleRegistrationRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "Finland with registration FPA-829",
			input: []byte{0x12, 0x46, 0x50, 0x41, 0x2D, 0x38, 0x32, 0x39, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "Spain with registration ABC-1234",
			input: []byte{0x0F, 0x41, 0x42, 0x43, 0x2D, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "Germany with registration B-MW-1234",
			input: []byte{0x0D, 0x42, 0x2D, 0x4D, 0x57, 0x2D, 0x31, 0x32, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "France with registration 75ABC12",
			input: []byte{0x11, 0x37, 0x35, 0x41, 0x42, 0x43, 0x31, 0x32, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "Italy with full 14-char registration",
			input: []byte{0x1A, 0x41, 0x42, 0x43, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x58},
		},
		{
			name:  "Sweden with short registration",
			input: []byte{0x2C, 0x41, 0x42, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "Empty nation with spaces",
			input: []byte{0xFF, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
		{
			name:  "Default nation with registration",
			input: []byte{0x00, 0x54, 0x45, 0x53, 0x54, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			vehicleReg, err := UnmarshalVehicleRegistration(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalVehicleRegistration() error: %v", err)
			}

			if vehicleReg == nil {
				t.Fatal("UnmarshalVehicleRegistration() returned nil")
			}

			// Marshal back
			got, err := AppendVehicleRegistration(nil, vehicleReg)
			if err != nil {
				t.Fatalf("AppendVehicleRegistration() error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-original +got):\n%s", diff)
			}
		})
	}
}
