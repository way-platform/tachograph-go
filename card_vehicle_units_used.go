package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardVehicleUnitsUsed unmarshals vehicle units used data from a card EF.
//
// The data type `CardVehicleUnitsUsed` is specified in the Data Dictionary, Section 2.40.
//
// ASN.1 Definition:
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
func unmarshalCardVehicleUnitsUsed(data []byte) (*cardv1.VehicleUnitsUsed, error) {
	const (
		lenCardVehicleUnitsUsedHeader = 2 // CardVehicleUnitsUsed header size
	)

	if len(data) < lenCardVehicleUnitsUsedHeader {
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

// AppendCardVehicleUnitsUsed appends vehicle units used data to a byte slice.
//
// The data type `CardVehicleUnitsUsed` is specified in the Data Dictionary, Section 2.40.
//
// ASN.1 Definition:
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
func appendCardVehicleUnitsUsed(data []byte, vehicleUnits *cardv1.VehicleUnitsUsed) ([]byte, error) {
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
