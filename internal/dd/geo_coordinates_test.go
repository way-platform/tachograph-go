package dd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func TestUnmarshalGeoCoordinates(t *testing.T) {
	// Stuttgart: ~48°46.8'N, 9°18.0'E
	// Encoded as DDMM.M×10: lat=48468 (0x00BD54), lon=9180 (0x0023DC)
	stuttgartLat := int32(48468)
	stuttgartLon := int32(9180)
	wantStuttgart := &ddv1.GeoCoordinates{}
	wantStuttgart.SetLatitude(stuttgartLat)
	wantStuttgart.SetLongitude(stuttgartLon)

	wantNullIsland := &ddv1.GeoCoordinates{}
	wantNullIsland.SetLatitude(0)
	wantNullIsland.SetLongitude(0)

	tests := []struct {
		name  string
		input []byte
		want  *ddv1.GeoCoordinates
	}{
		{
			// Sentinel: both 0x7FFFFF → position unavailable
			name:  "both fields 0x7FFFFF returns nil",
			input: []byte{0x7F, 0xFF, 0xFF, 0x7F, 0xFF, 0xFF},
			want:  nil,
		},
		{
			// Null Island (0°N, 0°E) is a valid coordinate, not a sentinel
			name:  "both fields 0x000000 is valid null island",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  wantNullIsland,
		},
		{
			// Realistic coordinates: Stuttgart ~48°46.8'N 9°18.0'E
			name:  "Stuttgart coordinates",
			input: []byte{0x00, 0xBD, 0x54, 0x00, 0x23, 0xDC},
			want:  wantStuttgart,
		},
		{
			// Mixed sentinel (lat=0x7FFFFF, lon=0): malformed → treat as unavailable
			name:  "lat sentinel only returns nil",
			input: []byte{0x7F, 0xFF, 0xFF, 0x00, 0x00, 0x00},
			want:  nil,
		},
		{
			// Mixed sentinel (lat=0, lon=0x7FFFFF): malformed → treat as unavailable
			name:  "lon sentinel only returns nil",
			input: []byte{0x00, 0x00, 0x00, 0x7F, 0xFF, 0xFF},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := UnmarshalOptions{}
			got, err := opts.UnmarshalGeoCoordinates(tt.input)
			if err != nil {
				t.Fatalf("UnmarshalGeoCoordinates() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("UnmarshalGeoCoordinates() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
