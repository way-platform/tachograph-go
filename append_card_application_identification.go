package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardApplicationIdentification appends application identification data to a byte slice.
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
func AppendCardApplicationIdentification(data []byte, appId *cardv1.ApplicationIdentification) ([]byte, error) {
	const (
		// Minimum EF_ApplicationIdentification record size
		MIN_EF_APPLICATION_IDENTIFICATION_SIZE = 7

		// Field offsets
		TYPE_OF_TACHOGRAPH_CARD_ID_OFFSET       = 0
		CARD_STRUCTURE_VERSION_OFFSET           = 1
		NO_OF_EVENTS_PER_TYPE_OFFSET            = 3
		NO_OF_FAULTS_PER_TYPE_OFFSET            = 4
		ACTIVITY_STRUCTURE_LENGTH_OFFSET        = 5
		NO_OF_CARD_VEHICLE_RECORDS_OFFSET       = 7
		NO_OF_CARD_PLACE_RECORDS_OFFSET         = 8
		NO_OF_GNSS_AD_RECORDS_OFFSET            = 9
		NO_OF_SPECIFIC_CONDITION_RECORDS_OFFSET = 10
		NO_OF_CARD_VEHICLE_UNIT_RECORDS_OFFSET  = 11

		// Field sizes
		CARD_STRUCTURE_VERSION_SIZE    = 2
		ACTIVITY_STRUCTURE_LENGTH_SIZE = 2
	)

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
