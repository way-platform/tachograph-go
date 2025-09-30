package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalEquipmentType(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    ddv1.EquipmentType
		wantErr bool
	}{
		{
			name:  "driver card",
			input: []byte{0x01},
			want:  ddv1.EquipmentType_DRIVER_CARD,
		},
		{
			name:  "workshop card",
			input: []byte{0x02},
			want:  ddv1.EquipmentType_WORKSHOP_CARD,
		},
		{
			name:  "control card",
			input: []byte{0x03},
			want:  ddv1.EquipmentType_CONTROL_CARD,
		},
		{
			name:  "company card",
			input: []byte{0x04},
			want:  ddv1.EquipmentType_COMPANY_CARD,
		},
		{
			name:  "vehicle unit",
			input: []byte{0x06},
			want:  ddv1.EquipmentType_VEHICLE_UNIT,
		},
		{
			name:  "motion sensor",
			input: []byte{0x07},
			want:  ddv1.EquipmentType_MOTION_SENSOR,
		},
		{
			name:  "unassigned value",
			input: []byte{0xFF},
			want:  ddv1.EquipmentType_EQUIPMENT_TYPE_UNRECOGNIZED,
		},
		{
			name:  "zero value - reserved",
			input: []byte{0x00},
			want:  ddv1.EquipmentType_RESERVED_MEMBER_STATE_OR_EUROPE,
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
			got, err := opts.UnmarshalEquipmentType(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalEquipmentType() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalEquipmentType() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("UnmarshalEquipmentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendEquipmentType(t *testing.T) {
	tests := []struct {
		name  string
		input ddv1.EquipmentType
		want  []byte
	}{
		{
			name:  "driver card",
			input: ddv1.EquipmentType_DRIVER_CARD,
			want:  []byte{0x01},
		},
		{
			name:  "workshop card",
			input: ddv1.EquipmentType_WORKSHOP_CARD,
			want:  []byte{0x02},
		},
		{
			name:  "control card",
			input: ddv1.EquipmentType_CONTROL_CARD,
			want:  []byte{0x03},
		},
		{
			name:  "company card",
			input: ddv1.EquipmentType_COMPANY_CARD,
			want:  []byte{0x04},
		},
		{
			name:  "vehicle unit",
			input: ddv1.EquipmentType_VEHICLE_UNIT,
			want:  []byte{0x06},
		},
		{
			name:  "motion sensor",
			input: ddv1.EquipmentType_MOTION_SENSOR,
			want:  []byte{0x07},
		},
		{
			name:  "reserved member state or europe",
			input: ddv1.EquipmentType_RESERVED_MEMBER_STATE_OR_EUROPE,
			want:  []byte{0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte{}
			got := AppendEquipmentType(dst, tt.input)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendEquipmentType() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEquipmentTypeRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "driver card",
			input: []byte{0x01},
		},
		{
			name:  "workshop card",
			input: []byte{0x02},
		},
		{
			name:  "control card",
			input: []byte{0x03},
		},
		{
			name:  "company card",
			input: []byte{0x04},
		},
		{
			name:  "vehicle unit",
			input: []byte{0x06},
		},
		{
			name:  "motion sensor",
			input: []byte{0x07},
		},
		{
			name:  "reserved member state or europe",
			input: []byte{0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var opts UnmarshalOptions
			equipmentType, err := opts.UnmarshalEquipmentType(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalEquipmentType() unexpected error: %v", err)
			}

			// Marshal
			dst := []byte{}
			got := AppendEquipmentType(dst, equipmentType)

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
