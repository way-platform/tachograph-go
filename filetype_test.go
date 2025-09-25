package tachograph

import (
	"testing"

	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

func TestInferFileType(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected tachographv1.File_Type
	}{
		{
			name:     "Valid Card File",
			data:     []byte{0x00, 0x02, 0x00, 0x0A, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A},
			expected: tachographv1.File_DRIVER_CARD,
		},
		{
			name:     "Valid Unit File (TV format - VuOverviewFirstGen)",
			data:     []byte{0x76, 0x01, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Valid Unit File (TV format - VuActivitiesFirstGen)",
			data:     []byte{0x76, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Valid Unit File (TV format - VuEventsAndFaultsFirstGen)",
			data:     []byte{0x76, 0x03, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Valid Unit File (TV format - VuOverviewSecondGen)",
			data:     []byte{0x76, 0x21, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Valid Unit File (TV format - VuOverviewSecondGenV2)",
			data:     []byte{0x76, 0x31, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Valid Unit File (TV format - VuDownloadInterfaceVersion)",
			data:     []byte{0x76, 0x00, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Valid Unit File (TRTP format)",
			data:     []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Another Valid Unit File (TRTP format)",
			data:     []byte{0x1F, 0x02, 0x03, 0x04, 0x05},
			expected: tachographv1.File_VEHICLE_UNIT,
		},
		{
			name:     "Unknown File",
			data:     []byte{0xFF, 0xFE, 0xFD, 0xFC},
			expected: tachographv1.File_TYPE_UNSPECIFIED,
		},
		{
			name:     "Empty File",
			data:     []byte{},
			expected: tachographv1.File_TYPE_UNSPECIFIED,
		},
		{
			name:     "Too Short File",
			data:     []byte{0x00},
			expected: tachographv1.File_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferFileType(tt.data); got != tt.expected {
				t.Errorf("InferFileType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
