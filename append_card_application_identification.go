package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardApplicationIdentification appends application identification data to a byte slice.
func AppendCardApplicationIdentification(data []byte, appId *cardv1.ApplicationIdentification) ([]byte, error) {
	if appId == nil {
		return data, nil
	}

	// Type of tachograph card ID (1 byte)
	if appId.HasTypeOfTachographCardId() {
		protocolValue := GetEquipmentTypeProtocolValue(appId.GetTypeOfTachographCardId(), 0)
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
