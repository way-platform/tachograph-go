package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardVehicleUnitsUsed unmarshals vehicle units used data from a card EF.
func UnmarshalCardVehicleUnitsUsed(data []byte, target *cardv1.VehicleUnitsUsed) error {
	if len(data) < 2 {
		return fmt.Errorf("insufficient data for vehicle units used")
	}

	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordPointer uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordPointer); err != nil {
		return fmt.Errorf("failed to read newest record pointer: %w", err)
	}
	target.SetVehicleUnitPointerNewestRecord(int32(newestRecordPointer))

	// For now, just set empty records to satisfy the interface
	// The actual vehicle units structure is complex and would need detailed parsing
	target.SetRecords([]*cardv1.VehicleUnitsUsed_Record{})

	return nil
}
