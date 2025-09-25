package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardVehicleUnitsUsed unmarshals vehicle units used data from a card EF.
func unmarshalCardVehicleUnitsUsed(data []byte) (*cardv1.VehicleUnitsUsed, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for vehicle units used")
	}

	var target cardv1.VehicleUnitsUsed
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordPointer uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordPointer); err != nil {
		return nil, fmt.Errorf("failed to read newest record pointer: %w", err)
	}
	target.SetVehicleUnitPointerNewestRecord(int32(newestRecordPointer))

	// For now, just set empty records to satisfy the interface
	// The actual vehicle units structure is complex and would need detailed parsing
	target.SetRecords([]*cardv1.VehicleUnitsUsed_Record{})

	return &target, nil
}

// UnmarshalCardVehicleUnitsUsed unmarshals vehicle units used data from a card EF (legacy function).
// Deprecated: Use unmarshalCardVehicleUnitsUsed instead.
func UnmarshalCardVehicleUnitsUsed(data []byte, target *cardv1.VehicleUnitsUsed) error {
	result, err := unmarshalCardVehicleUnitsUsed(data)
	if err != nil {
		return err
	}
	*target = *result
	return nil
}
