package tachograph

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// MarshalFile serializes a protobuf File message into the binary DDD file format.
func MarshalFile(file *tachographv1.File) ([]byte, error) {
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		return card.MarshalDriverCardFile(file.GetDriverCard())
	case tachographv1.File_VEHICLE_UNIT:
		return vu.MarshalVehicleUnitFile(file.GetVehicleUnit())
	case tachographv1.File_RAW_CARD:
		return card.MarshalRawCardFile(file.GetRawCard())
	default:
		return nil, fmt.Errorf("unsupported file type for marshaling: %v", file.GetType())
	}
}
