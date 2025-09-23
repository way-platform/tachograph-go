package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardVehicleUnitsUsed appends vehicle units used data to a byte slice.
func AppendCardVehicleUnitsUsed(data []byte, vehicleUnits *cardv1.VehicleUnitsUsed) ([]byte, error) {
	if vehicleUnits == nil {
		return data, nil
	}

	// Newest record pointer (2 bytes)
	if vehicleUnits.HasVehicleUnitPointerNewestRecord() {
		pointer := make([]byte, 2)
		binary.BigEndian.PutUint16(pointer, uint16(vehicleUnits.GetVehicleUnitPointerNewestRecord()))
		data = append(data, pointer...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// For now, skip the complex record structures
	// This provides a basic implementation that satisfies the interface

	return data, nil
}
