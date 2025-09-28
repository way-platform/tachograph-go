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
// The data type `ApplicationIdentification` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition:
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
func unmarshalCardApplicationIdentification(data []byte) (*cardv1.ApplicationIdentification, error) {
	const (
		lenMinEfApplicationIdentification = 7 // Minimum EF_ApplicationIdentification record size
	)

	if len(data) < lenMinEfApplicationIdentification {
		return nil, fmt.Errorf("insufficient data for application identification: got %d bytes, need at least %d", len(data), lenMinEfApplicationIdentification)
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

// AppendCardApplicationIdentification appends application identification data to a byte slice.
//
// The data type `ApplicationIdentification` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition:
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
func appendCardApplicationIdentification(data []byte, appId *cardv1.ApplicationIdentification) ([]byte, error) {
	if appId == nil {
		return data, nil
	}

	// Type of tachograph card ID (1 byte)
	if appId.HasTypeOfTachographCardId() {
		protocolValue := GetProtocolValueFromEnum(appId.GetTypeOfTachographCardId(), 0)
		data = append(data, byte(protocolValue))
	} else {
		data = append(data, 0x00)
	}

	// Card structure version (2 bytes)
	structureVersion := appId.GetCardStructureVersion()
	if structureVersion != nil {
		// Convert major and minor to BCD format
		major := structureVersion.GetMajor()
		minor := structureVersion.GetMinor()
		majorBCD := ((major / 10) << 4) | (major % 10)
		minorBCD := ((minor / 10) << 4) | (minor % 10)
		data = append(data, byte(majorBCD), byte(minorBCD))
	} else {
		data = append(data, 0x00, 0x01) // Default version
	}

	// Get driver data for the specific fields
	var driver *cardv1.ApplicationIdentification_Driver
	switch appId.GetCardType() {
	case cardv1.CardType_DRIVER_CARD:
		driver = appId.GetDriver()
	}

	if driver == nil {
		// If no driver data, append zeros for all driver-specific fields
		data = append(data, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // events, faults, activity length (2 bytes), vehicle records, place records
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

	// Card vehicle records count (1 byte)
	if driver.HasCardVehicleRecordsCount() {
		data = append(data, byte(driver.GetCardVehicleRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// Card place records count (1 byte)
	if driver.HasCardPlaceRecordsCount() {
		data = append(data, byte(driver.GetCardPlaceRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// Gen2 fields
	if driver.HasGnssAdRecordsCount() {
		data = append(data, byte(driver.GetGnssAdRecordsCount()))
	}

	if driver.HasSpecificConditionRecordsCount() {
		data = append(data, byte(driver.GetSpecificConditionRecordsCount()))
	}

	// Gen2v2 fields
	if driver.HasCardVehicleUnitRecordsCount() {
		data = append(data, byte(driver.GetCardVehicleUnitRecordsCount()))
	}

	return data, nil
}
