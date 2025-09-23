package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardApplicationIdentification appends application identification data to a byte slice.
func AppendCardApplicationIdentification(data []byte, appId *cardv1.DriverCardApplicationIdentification) ([]byte, error) {
	if appId == nil {
		return data, nil
	}

	// Type of tachograph card ID (1 byte)
	if appId.HasTypeOfTachographCardId() {
		protocolValue := GetEquipmentTypeProtocolValue(appId.GetTypeOfTachographCardId(), appId.GetUnrecognizedTypeOfTachographCardId())
		data = append(data, byte(protocolValue))
	} else {
		data = append(data, 0x00)
	}

	// Card structure version (2 bytes)
	structureVersion := appId.GetCardStructureVersion()
	if len(structureVersion) >= 2 {
		data = append(data, structureVersion[:2]...)
	} else {
		data = append(data, 0x00, 0x01) // Default version
	}

	// Events per type count (1 byte)
	if appId.HasEventsPerTypeCount() {
		data = append(data, byte(appId.GetEventsPerTypeCount()))
	} else {
		data = append(data, 0x00)
	}

	// Faults per type count (1 byte)
	if appId.HasFaultsPerTypeCount() {
		data = append(data, byte(appId.GetFaultsPerTypeCount()))
	} else {
		data = append(data, 0x00)
	}

	// Activity structure length (2 bytes)
	if appId.HasActivityStructureLength() {
		activityLength := make([]byte, 2)
		binary.BigEndian.PutUint16(activityLength, uint16(appId.GetActivityStructureLength()))
		data = append(data, activityLength...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Card vehicle records count (1 byte)
	if appId.HasCardVehicleRecordsCount() {
		data = append(data, byte(appId.GetCardVehicleRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// Card place records count (1 byte)
	if appId.HasCardPlaceRecordsCount() {
		data = append(data, byte(appId.GetCardPlaceRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// Gen2 fields
	if appId.HasGnssAdRecordsCount() {
		data = append(data, byte(appId.GetGnssAdRecordsCount()))
	}

	if appId.HasSpecificConditionRecordsCount() {
		data = append(data, byte(appId.GetSpecificConditionRecordsCount()))
	}

	// Gen2v2 fields
	if appId.HasCardVehicleUnitRecordsCount() {
		data = append(data, byte(appId.GetCardVehicleUnitRecordsCount()))
	}

	return data, nil
}
