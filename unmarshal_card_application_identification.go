package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func unmarshalCardApplicationIdentification(data []byte) (*cardv1.DriverCardApplicationIdentification, error) {
	if len(data) < 7 {
		return nil, fmt.Errorf("insufficient data for application identification: got %d bytes, need at least 7", len(data))
	}
	var target cardv1.DriverCardApplicationIdentification
	r := bytes.NewReader(data)
	// Read type of tachograph card ID (1 byte)
	var cardType byte
	if err := binary.Read(r, binary.BigEndian, &cardType); err != nil {
		return nil, fmt.Errorf("failed to read card type: %w", err)
	}
	// Convert raw card type to enum using protocol annotations
	SetEquipmentType(int32(cardType), target.SetTypeOfTachographCardId, target.SetUnrecognizedTypeOfTachographCardId)
	// Read card structure version (2 bytes)
	structureVersionBytes := make([]byte, 2)
	if _, err := r.Read(structureVersionBytes); err != nil {
		return nil, fmt.Errorf("failed to read card structure version: %w", err)
	}
	target.SetCardStructureVersion(structureVersionBytes)
	// Read events per type count (1 byte)
	var eventsPerType byte
	if err := binary.Read(r, binary.BigEndian, &eventsPerType); err != nil {
		return nil, fmt.Errorf("failed to read events per type count: %w", err)
	}
	target.SetEventsPerTypeCount(int32(eventsPerType))
	// Read faults per type count (1 byte)
	var faultsPerType byte
	if err := binary.Read(r, binary.BigEndian, &faultsPerType); err != nil {
		return nil, fmt.Errorf("failed to read faults per type count: %w", err)
	}
	target.SetFaultsPerTypeCount(int32(faultsPerType))
	// Read activity structure length (2 bytes)
	var activityLength uint16
	if err := binary.Read(r, binary.BigEndian, &activityLength); err != nil {
		return nil, fmt.Errorf("failed to read activity structure length: %w", err)
	}
	target.SetActivityStructureLength(int32(activityLength))
	// Read card vehicle records count (1 byte)
	var vehicleRecords byte
	if err := binary.Read(r, binary.BigEndian, &vehicleRecords); err != nil {
		return nil, fmt.Errorf("failed to read vehicle records count: %w", err)
	}
	target.SetCardVehicleRecordsCount(int32(vehicleRecords))
	// Read card place records count (1 byte)
	var placeRecords byte
	if err := binary.Read(r, binary.BigEndian, &placeRecords); err != nil {
		return nil, fmt.Errorf("failed to read place records count: %w", err)
	}
	target.SetCardPlaceRecordsCount(int32(placeRecords))
	// Gen2 fields (if more data available)
	if r.Len() >= 2 {
		// Read GNSS AD records count (1 byte)
		var gnssAdRecords byte
		if err := binary.Read(r, binary.BigEndian, &gnssAdRecords); err == nil {
			target.SetGnssAdRecordsCount(int32(gnssAdRecords))
		}
		// Read specific condition records count (1 byte)
		var specificConditionRecords byte
		if err := binary.Read(r, binary.BigEndian, &specificConditionRecords); err == nil {
			target.SetSpecificConditionRecordsCount(int32(specificConditionRecords))
		}
	}
	// Gen2v2 fields (if more data available)
	if r.Len() >= 1 {
		// Read card vehicle unit records count (1 byte)
		var vehicleUnitRecords byte
		if err := binary.Read(r, binary.BigEndian, &vehicleUnitRecords); err == nil {
			target.SetCardVehicleUnitRecordsCount(int32(vehicleUnitRecords))
		}
	}
	return &target, nil
}
