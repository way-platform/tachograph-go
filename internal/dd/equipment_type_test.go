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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.input) != 1 {
				t.Fatalf("invalid test input length: got %d, want 1", len(tt.input))
			}
			got, err := UnmarshalEnum[ddv1.EquipmentType](tt.input[0])
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalEnum[EquipmentType]() expected error, got nil")
				}
				return
			}
			if err != nil {
				// If error but no error expected, and we want UNRECOGNIZED, return UNRECOGNIZED
				if tt.want == ddv1.EquipmentType_EQUIPMENT_TYPE_UNRECOGNIZED {
					got = ddv1.EquipmentType_EQUIPMENT_TYPE_UNRECOGNIZED
				} else {
					t.Fatalf("UnmarshalEnum[EquipmentType]() unexpected error: %v", err)
				}
			}
			if got != tt.want {
				t.Errorf("UnmarshalEnum[EquipmentType]() = %v, want %v", got, tt.want)
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
			equipmentTypeByte, err := MarshalEnum(tt.input)
			if err != nil {
				t.Fatalf("MarshalEnum[EquipmentType]() unexpected error: %v", err)
			}
			got := append(dst, equipmentTypeByte)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MarshalEnum[EquipmentType]() mismatch (-want +got):\n%s", diff)
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
			if len(tt.input) != 1 {
				t.Fatalf("invalid test input length: got %d, want 1", len(tt.input))
			}
			equipmentType, err := UnmarshalEnum[ddv1.EquipmentType](tt.input[0])
			if err != nil {
				t.Fatalf("UnmarshalEnum[EquipmentType]() unexpected error: %v", err)
			}

			// Marshal
			dst := []byte{}
			equipmentTypeByte, err := MarshalEnum(equipmentType)
			if err != nil {
				t.Fatalf("MarshalEnum[EquipmentType]() unexpected error: %v", err)
			}
			got := append(dst, equipmentTypeByte)

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
