package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalExtendedSerialNumber(t *testing.T) {
	tests := []struct {
		name                 string
		input                []byte
		wantSerialNumber     int64
		wantMonth            int32
		wantYear             int32
		wantEquipmentType    ddv1.EquipmentType
		wantManufacturerCode int32
		wantErr              bool
	}{
		{
			name:                 "valid extended serial number",
			input:                []byte{0x00, 0xBC, 0x61, 0x4E, 0x09, 0x25, 0x06, 0x0A},
			wantSerialNumber:     12345678,
			wantMonth:            9,
			wantYear:             2025,
			wantEquipmentType:    ddv1.EquipmentType_VEHICLE_UNIT,
			wantManufacturerCode: 10,
		},
		{
			name:                 "serial number with month/year 01/24",
			input:                []byte{0x00, 0x00, 0x00, 0x01, 0x01, 0x24, 0x01, 0x05},
			wantSerialNumber:     1,
			wantMonth:            1,
			wantYear:             2024,
			wantEquipmentType:    ddv1.EquipmentType_DRIVER_CARD,
			wantManufacturerCode: 5,
		},
		{
			name:                 "zero serial number",
			input:                []byte{0x00, 0x00, 0x00, 0x00, 0x12, 0x23, 0x07, 0x0F},
			wantSerialNumber:     0,
			wantMonth:            12,
			wantYear:             2023,
			wantEquipmentType:    ddv1.EquipmentType_MOTION_SENSOR,
			wantManufacturerCode: 15,
		},
		{
			name:                 "maximum serial number",
			input:                []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x12, 0x99, 0x02, 0xFF},
			wantSerialNumber:     4294967295,
			wantMonth:            12,
			wantYear:             2099,
			wantEquipmentType:    ddv1.EquipmentType_WORKSHOP_CARD,
			wantManufacturerCode: 255,
		},
		{
			name:    "insufficient data",
			input:   []byte{0x00, 0x00, 0x00},
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
			got, err := UnmarshalExtendedSerialNumber(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalExtendedSerialNumber() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalExtendedSerialNumber() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("UnmarshalExtendedSerialNumber() returned nil")
			}
			if got.GetSerialNumber() != tt.wantSerialNumber {
				t.Errorf("GetSerialNumber() = %v, want %v", got.GetSerialNumber(), tt.wantSerialNumber)
			}
			if got.GetMonthYear() == nil {
				t.Fatal("GetMonthYear() returned nil")
			}
			if got.GetMonthYear().GetMonth() != tt.wantMonth {
				t.Errorf("GetMonthYear().GetMonth() = %v, want %v", got.GetMonthYear().GetMonth(), tt.wantMonth)
			}
			if got.GetMonthYear().GetYear() != tt.wantYear {
				t.Errorf("GetMonthYear().GetYear() = %v, want %v", got.GetMonthYear().GetYear(), tt.wantYear)
			}
			if got.GetType() != tt.wantEquipmentType {
				t.Errorf("GetType() = %v, want %v", got.GetType(), tt.wantEquipmentType)
			}
			if got.GetManufacturerCode() != tt.wantManufacturerCode {
				t.Errorf("GetManufacturerCode() = %v, want %v", got.GetManufacturerCode(), tt.wantManufacturerCode)
			}
		})
	}
}

func TestAppendExtendedSerialNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   *ddv1.ExtendedSerialNumber
		want    []byte
		wantErr bool
	}{
		{
			name: "valid extended serial number",
			input: func() *ddv1.ExtendedSerialNumber {
				esn := &ddv1.ExtendedSerialNumber{}
				esn.SetSerialNumber(12345678)
				my := &ddv1.MonthYear{}
				my.SetEncoded([]byte{0x09, 0x25})
				my.SetMonth(9)
				my.SetYear(2025)
				esn.SetMonthYear(my)
				esn.SetType(ddv1.EquipmentType_VEHICLE_UNIT)
				esn.SetManufacturerCode(10)
				return esn
			}(),
			want: []byte{0x00, 0xBC, 0x61, 0x4E, 0x09, 0x25, 0x06, 0x0A},
		},
		{
			name: "with decoded month/year only",
			input: func() *ddv1.ExtendedSerialNumber {
				esn := &ddv1.ExtendedSerialNumber{}
				esn.SetSerialNumber(1)
				my := &ddv1.MonthYear{}
				my.SetMonth(1)
				my.SetYear(2024)
				esn.SetMonthYear(my)
				esn.SetType(ddv1.EquipmentType_DRIVER_CARD)
				esn.SetManufacturerCode(5)
				return esn
			}(),
			want: []byte{0x00, 0x00, 0x00, 0x01, 0x01, 0x24, 0x01, 0x05},
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true, // Errors because it calls nested AppendMonthYear
		},
		{
			name: "zero values",
			input: func() *ddv1.ExtendedSerialNumber {
				esn := &ddv1.ExtendedSerialNumber{}
				esn.SetSerialNumber(0)
				esn.SetType(ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED)
				esn.SetManufacturerCode(0)
				return esn
			}(),
			want: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte{}
			got, err := AppendExtendedSerialNumber(dst, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AppendExtendedSerialNumber() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendExtendedSerialNumber() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AppendExtendedSerialNumber() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExtendedSerialNumberRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "full data",
			input: []byte{0x00, 0xBC, 0x61, 0x4E, 0x09, 0x25, 0x06, 0x0A},
		},
		{
			name:  "different values",
			input: []byte{0x00, 0x00, 0x00, 0x01, 0x01, 0x24, 0x01, 0x05},
		},
		{
			name:  "zero values",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "maximum values",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x12, 0x99, 0x07, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			esn, err := UnmarshalExtendedSerialNumber(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalExtendedSerialNumber() unexpected error: %v", err)
			}
			if esn == nil {
				t.Fatal("UnmarshalExtendedSerialNumber() returned nil")
			}

			// Marshal
			dst := []byte{}
			got, err := AppendExtendedSerialNumber(dst, esn)
			if err != nil {
				t.Fatalf("AppendExtendedSerialNumber() unexpected error: %v", err)
			}

			// Verify round-trip
			if diff := cmp.Diff(tt.input, got); diff != "" {
				t.Errorf("Round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
