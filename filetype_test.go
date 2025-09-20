package tacho

import "testing"

func TestInferFileType(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected FileType
	}{
		{
			name:     "Valid Card File",
			data:     []byte{0x00, 0x02, 0x00, 0x0A, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A},
			expected: CardFileType,
		},
		{
			name:     "Valid Unit File (TV format - VuOverviewFirstGen)",
			data:     []byte{0x76, 0x01, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Valid Unit File (TV format - VuActivitiesFirstGen)",
			data:     []byte{0x76, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Valid Unit File (TV format - VuEventsAndFaultsFirstGen)",
			data:     []byte{0x76, 0x03, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Valid Unit File (TV format - VuOverviewSecondGen)",
			data:     []byte{0x76, 0x21, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Valid Unit File (TV format - VuOverviewSecondGenV2)",
			data:     []byte{0x76, 0x31, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Valid Unit File (TV format - VuDownloadInterfaceVersion)",
			data:     []byte{0x76, 0x00, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Valid Unit File (TRTP format)",
			data:     []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Another Valid Unit File (TRTP format)",
			data:     []byte{0x1F, 0x02, 0x03, 0x04, 0x05},
			expected: UnitFileType,
		},
		{
			name:     "Unknown File",
			data:     []byte{0xFF, 0xFE, 0xFD, 0xFC},
			expected: UnknownFileType,
		},
		{
			name:     "Empty File",
			data:     []byte{},
			expected: UnknownFileType,
		},
		{
			name:     "Too Short File",
			data:     []byte{0x00},
			expected: UnknownFileType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InferFileType(tt.data); got != tt.expected {
				t.Errorf("InferFileType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
