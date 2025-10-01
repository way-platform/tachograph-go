package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalApplicationIdentificationG2 parses the binary data for an EF_ApplicationIdentification record (Gen2 format).
//
// The data type `ApplicationIdentification` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition (Gen2):
//
//	ApplicationIdentification ::= SEQUENCE {
//	    typeOfTachographCardId    EquipmentType,
//	    cardStructureVersion      CardStructureVersion,
//	    noOfEventsPerType         INTEGER(0..255),
//	    noOfFaultsPerType         INTEGER(0..255),
//	    activityStructureLength   INTEGER(0..65535),
//	    noOfCardVehicleRecords    INTEGER(0..255),
//	    noOfCardPlaceRecords      INTEGER(0..255),
//	    noOfGNSSADRecords         INTEGER(0..255),
//	    noOfSpecificConditionRecords INTEGER(0..255),
//	    noOfCardVehicleUnitRecords   INTEGER(0..255)
//	}
func (opts UnmarshalOptions) unmarshalApplicationIdentificationG2(data []byte) (*cardv1.ApplicationIdentificationG2, error) {
	const (
		lenEfApplicationIdentificationG2 = 17 // Gen2: 1 + 2 + 1 + 1 + 2 + 2 + 2 + 2 + 2 + 2 = 17 bytes
	)

	if len(data) != lenEfApplicationIdentificationG2 {
		return nil, fmt.Errorf("invalid data length for Gen2 application identification: got %d bytes, want %d", len(data), lenEfApplicationIdentificationG2)
	}

	target := &cardv1.ApplicationIdentificationG2{}
	r := bytes.NewReader(data)

	// Read type of tachograph card ID (1 byte)
	var cardType byte
	if err := binary.Read(r, binary.BigEndian, &cardType); err != nil {
		return nil, fmt.Errorf("failed to read card type: %w", err)
	}
	// Convert raw card type to enum using protocol annotations
	if equipmentType, err := dd.UnmarshalEnum[ddv1.EquipmentType](cardType); err == nil {
		target.SetTypeOfTachographCardId(equipmentType)
	} else {
		return nil, fmt.Errorf("invalid equipment type: %w", err)
	}

	// Read card structure version (2 bytes)
	structureVersionBytes := make([]byte, 2)
	if _, err := r.Read(structureVersionBytes); err != nil {
		return nil, fmt.Errorf("failed to read card structure version: %w", err)
	}
	// Parse BCD structure version using centralized helper
	cardStructureVersion, err := opts.UnmarshalCardStructureVersion(structureVersionBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal card structure version: %w", err)
	}
	target.SetCardStructureVersion(cardStructureVersion)

	// For now, assume this is a driver card and create the driver data
	driver := &cardv1.ApplicationIdentificationG2_Driver{}

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

	// Read card vehicle records count (2 bytes in Gen2)
	var vehicleRecords uint16
	if err := binary.Read(r, binary.BigEndian, &vehicleRecords); err != nil {
		return nil, fmt.Errorf("failed to read vehicle records count: %w", err)
	}
	driver.SetCardVehicleRecordsCount(int32(vehicleRecords))

	// Read card place records count (2 bytes in Gen2)
	var placeRecords uint16
	if err := binary.Read(r, binary.BigEndian, &placeRecords); err != nil {
		return nil, fmt.Errorf("failed to read place records count: %w", err)
	}
	driver.SetCardPlaceRecordsCount(int32(placeRecords))

	// Gen2-specific fields:

	// Read GNSS AD records count (2 bytes)
	var gnssAdRecords uint16
	if err := binary.Read(r, binary.BigEndian, &gnssAdRecords); err != nil {
		return nil, fmt.Errorf("failed to read GNSS AD records count: %w", err)
	}
	driver.SetGnssAdRecordsCount(int32(gnssAdRecords))

	// Read specific condition records count (2 bytes)
	var specificConditionRecords uint16
	if err := binary.Read(r, binary.BigEndian, &specificConditionRecords); err != nil {
		return nil, fmt.Errorf("failed to read specific condition records count: %w", err)
	}
	driver.SetSpecificConditionRecordsCount(int32(specificConditionRecords))

	// Read card vehicle unit records count (2 bytes)
	var vehicleUnitRecords uint16
	if err := binary.Read(r, binary.BigEndian, &vehicleUnitRecords); err != nil {
		return nil, fmt.Errorf("failed to read vehicle unit records count: %w", err)
	}
	driver.SetCardVehicleUnitRecordsCount(int32(vehicleUnitRecords))

	// Set the driver data and card type
	target.SetDriver(driver)
	target.SetCardType(cardv1.CardType_DRIVER_CARD)

	return target, nil
}

// appendCardApplicationIdentificationG2 appends Gen2 application identification data to a byte slice.
//
// The data type `ApplicationIdentification` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition (Gen2):
//
//	ApplicationIdentification ::= SEQUENCE {
//	    typeOfTachographCardId    EquipmentType,
//	    cardStructureVersion      CardStructureVersion,
//	    noOfEventsPerType         INTEGER(0..255),
//	    noOfFaultsPerType         INTEGER(0..255),
//	    activityStructureLength   INTEGER(0..65535),
//	    noOfCardVehicleRecords    INTEGER(0..255),
//	    noOfCardPlaceRecords      INTEGER(0..255),
//	    noOfGNSSADRecords         INTEGER(0..255),
//	    noOfSpecificConditionRecords INTEGER(0..255),
//	    noOfCardVehicleUnitRecords   INTEGER(0..255)
//	}
func appendCardApplicationIdentificationG2(data []byte, appId *cardv1.ApplicationIdentificationG2) ([]byte, error) {
	if appId == nil {
		return data, nil
	}

	// Type of tachograph card ID (1 byte)
	if appId.HasTypeOfTachographCardId() {
		protocolValue, _ := dd.MarshalEnum(appId.GetTypeOfTachographCardId())
		data = append(data, protocolValue)
	} else {
		data = append(data, 0x00)
	}

	// Card structure version (2 bytes)
	structureVersion := appId.GetCardStructureVersion()
	if structureVersion != nil {
		// Append using centralized helper
		var err error
		data, err = dd.AppendCardStructureVersion(data, structureVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to append card structure version: %w", err)
		}
	} else {
		data = append(data, 0x00, 0x01) // Default version
	}

	// Get driver data for the specific fields
	var driver *cardv1.ApplicationIdentificationG2_Driver
	switch appId.GetCardType() {
	case cardv1.CardType_DRIVER_CARD:
		driver = appId.GetDriver()
	}

	if driver == nil {
		// If no driver data, append zeros for all driver-specific fields (14 bytes)
		data = append(data, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
		return data, nil
	}

	// Events per type count (1 byte)
	if driver.HasEventsPerTypeCount() {
		data = append(data, byte(driver.GetEventsPerTypeCount()))
	} else {
		data = append(data, 0x00)
	}

	// Faults per type count (1 byte)
	if driver.HasFaultsPerTypeCount() {
		data = append(data, byte(driver.GetFaultsPerTypeCount()))
	} else {
		data = append(data, 0x00)
	}

	// Activity structure length (2 bytes)
	if driver.HasActivityStructureLength() {
		activityLength := make([]byte, 2)
		binary.BigEndian.PutUint16(activityLength, uint16(driver.GetActivityStructureLength()))
		data = append(data, activityLength...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Card vehicle records count (2 bytes in Gen2)
	if driver.HasCardVehicleRecordsCount() {
		vehicleRecords := make([]byte, 2)
		binary.BigEndian.PutUint16(vehicleRecords, uint16(driver.GetCardVehicleRecordsCount()))
		data = append(data, vehicleRecords...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Card place records count (2 bytes in Gen2)
	if driver.HasCardPlaceRecordsCount() {
		placeRecords := make([]byte, 2)
		binary.BigEndian.PutUint16(placeRecords, uint16(driver.GetCardPlaceRecordsCount()))
		data = append(data, placeRecords...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Gen2-specific fields (always present in Gen2):

	// GNSS AD records count (2 bytes)
	if driver.HasGnssAdRecordsCount() {
		gnssAdRecords := make([]byte, 2)
		binary.BigEndian.PutUint16(gnssAdRecords, uint16(driver.GetGnssAdRecordsCount()))
		data = append(data, gnssAdRecords...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Specific condition records count (2 bytes)
	if driver.HasSpecificConditionRecordsCount() {
		specificConditionRecords := make([]byte, 2)
		binary.BigEndian.PutUint16(specificConditionRecords, uint16(driver.GetSpecificConditionRecordsCount()))
		data = append(data, specificConditionRecords...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Card vehicle unit records count (2 bytes)
	if driver.HasCardVehicleUnitRecordsCount() {
		vehicleUnitRecords := make([]byte, 2)
		binary.BigEndian.PutUint16(vehicleUnitRecords, uint16(driver.GetCardVehicleUnitRecordsCount()))
		data = append(data, vehicleUnitRecords...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	return data, nil
}
