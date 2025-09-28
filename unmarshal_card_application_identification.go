package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalCardApplicationIdentification parses the binary data for an EF_ApplicationIdentification record.
//
// ASN.1 Specification (Data Dictionary 2.2):
//
//	ApplicationIdentification ::= SEQUENCE {
//	    typeOfTachographCardId    EquipmentType,
//	    cardStructureVersion      CardStructureVersion,
//	    noOfEventsPerType         INTEGER(0..255),
//	    noOfFaultsPerType         INTEGER(0..255),
//	    activityStructureLength   INTEGER(0..65535),
//	    noOfCardVehicleRecords    INTEGER(0..255),
//	    noOfCardPlaceRecords      INTEGER(0..255),
//	    noOfGNSSADRecords         INTEGER(0..255) OPTIONAL,
//	    noOfSpecificConditionRecords INTEGER(0..255) OPTIONAL,
//	    noOfCardVehicleUnitRecords   INTEGER(0..255) OPTIONAL
//	}
//
// Binary Layout (7-10 bytes depending on version):
//
//	0-0:   typeOfTachographCardId (1 byte)
//	1-2:   cardStructureVersion (2 bytes, BCD format)
//	3-3:   noOfEventsPerType (1 byte)
//	4-4:   noOfFaultsPerType (1 byte)
//	5-6:   activityStructureLength (2 bytes, big-endian)
//	7-7:   noOfCardVehicleRecords (1 byte)
//	8-8:   noOfCardPlaceRecords (1 byte)
//	9-9:   noOfGNSSADRecords (1 byte, Gen2+)
//	10-10: noOfSpecificConditionRecords (1 byte, Gen2+)
//	11-11: noOfCardVehicleUnitRecords (1 byte, Gen2v2+)
func unmarshalCardApplicationIdentification(data []byte) (*cardv1.ApplicationIdentification, error) {
	const (
		// Minimum EF_ApplicationIdentification record size
		MIN_EF_APPLICATION_IDENTIFICATION_SIZE = 7
	)

	if len(data) < MIN_EF_APPLICATION_IDENTIFICATION_SIZE {
		return nil, fmt.Errorf("insufficient data for application identification: got %d bytes, need at least %d", len(data), MIN_EF_APPLICATION_IDENTIFICATION_SIZE)
	}

	target := &cardv1.ApplicationIdentification{}
	r := bytes.NewReader(data)

	// Read type of tachograph card ID (1 byte)
	var cardType byte
	if err := binary.Read(r, binary.BigEndian, &cardType); err != nil {
		return nil, fmt.Errorf("failed to read card type: %w", err)
	}
	// Convert raw card type to enum using protocol annotations
	SetEquipmentType(ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED.Descriptor(), int32(cardType), func(et protoreflect.EnumNumber) {
		target.SetTypeOfTachographCardId(ddv1.EquipmentType(et))
	}, nil)

	// Read card structure version (2 bytes)
	structureVersionBytes := make([]byte, 2)
	if _, err := r.Read(structureVersionBytes); err != nil {
		return nil, fmt.Errorf("failed to read card structure version: %w", err)
	}
	// Parse BCD structure version
	major := int32((structureVersionBytes[0]&0xF0)>>4)*10 + int32(structureVersionBytes[0]&0x0F)
	minor := int32((structureVersionBytes[1]&0xF0)>>4)*10 + int32(structureVersionBytes[1]&0x0F)
	cardStructureVersion := &ddv1.CardStructureVersion{}
	cardStructureVersion.SetMajor(major)
	cardStructureVersion.SetMinor(minor)
	target.SetCardStructureVersion(cardStructureVersion)

	// For now, assume this is a driver card and create the driver data
	driver := &cardv1.ApplicationIdentification_Driver{}

	// Read events per type count (1 byte)
	var eventsPerType byte
	if err := binary.Read(r, binary.BigEndian, &eventsPerType); err != nil {
		return nil, fmt.Errorf("failed to read events per type count: %w", err)
	}
	driver.SetEventsPerTypeCount(int32(eventsPerType))

	// Read faults per type count (1 byte)
	var faultsPerType byte
	if err := binary.Read(r, binary.BigEndian, &faultsPerType); err != nil {
		return nil, fmt.Errorf("failed to read faults per type count: %w", err)
	}
	driver.SetFaultsPerTypeCount(int32(faultsPerType))

	// Read activity structure length (2 bytes)
	var activityLength uint16
	if err := binary.Read(r, binary.BigEndian, &activityLength); err != nil {
		return nil, fmt.Errorf("failed to read activity structure length: %w", err)
	}
	driver.SetActivityStructureLength(int32(activityLength))

	// Read card vehicle records count (1 byte)
	var vehicleRecords byte
	if err := binary.Read(r, binary.BigEndian, &vehicleRecords); err != nil {
		return nil, fmt.Errorf("failed to read vehicle records count: %w", err)
	}
	driver.SetCardVehicleRecordsCount(int32(vehicleRecords))

	// Read card place records count (1 byte)
	var placeRecords byte
	if err := binary.Read(r, binary.BigEndian, &placeRecords); err != nil {
		return nil, fmt.Errorf("failed to read place records count: %w", err)
	}
	driver.SetCardPlaceRecordsCount(int32(placeRecords))

	// Gen2 fields (if more data available)
	if r.Len() >= 2 {
		// Read GNSS AD records count (1 byte)
		var gnssAdRecords byte
		if err := binary.Read(r, binary.BigEndian, &gnssAdRecords); err == nil {
			driver.SetGnssAdRecordsCount(int32(gnssAdRecords))
		}
		// Read specific condition records count (1 byte)
		var specificConditionRecords byte
		if err := binary.Read(r, binary.BigEndian, &specificConditionRecords); err == nil {
			driver.SetSpecificConditionRecordsCount(int32(specificConditionRecords))
		}
	}

	// Gen2v2 fields (if more data available)
	if r.Len() >= 1 {
		// Read card vehicle unit records count (1 byte)
		var vehicleUnitRecords byte
		if err := binary.Read(r, binary.BigEndian, &vehicleUnitRecords); err == nil {
			driver.SetCardVehicleUnitRecordsCount(int32(vehicleUnitRecords))
		}
	}

	// Set the driver data and card type
	target.SetDriver(driver)
	target.SetCardType(cardv1.CardType_DRIVER_CARD)

	return target, nil
}
