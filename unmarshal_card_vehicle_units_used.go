package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalCardVehicleUnitsUsed unmarshals vehicle units used data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.40):
//
//	CardVehicleUnitsUsed ::= SEQUENCE {
//	    vehicleUnitPointerNewestRecord     INTEGER(0..NoOfCardVehicleUnitRecords-1),
//	    cardVehicleUnitRecords             SET SIZE(NoOfCardVehicleUnitRecords) OF CardVehicleUnitRecord
//	}
//
//	CardVehicleUnitRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,
//	    manufacturerCode                   ManufacturerCode,
//	    deviceID                           DeviceID,
//	    vuSoftwareVersion                  VuSoftwareVersion
//	}
//
// TODO: Add detailed ASN.1 Specification - CardVehicleUnitRecord definition not found in data dictionary
//
// Binary Layout (variable size):
//
//	0-1:   vehicleUnitPointerNewestRecord (2 bytes, big-endian)
//	2+:    cardVehicleUnitRecords (variable size each)
//
// Constants:
const (
	// CardVehicleUnitsUsed header size
	cardVehicleUnitsUsedHeaderSize = 2
)

func unmarshalCardVehicleUnitsUsed(data []byte) (*cardv1.VehicleUnitsUsed, error) {
	if len(data) < cardVehicleUnitsUsedHeaderSize {
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
	proto.Merge(target, result)
	return nil
}
