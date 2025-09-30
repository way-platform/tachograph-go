package card

import (
	"github.com/way-platform/tachograph-go/internal/dd"
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

	// Write vehicle unit records
	records := vehicleUnits.GetRecords()
	if len(records) > 0 {
		// Write number of records (1 byte)
		if len(records) > 255 {
			return nil, fmt.Errorf("too many vehicle unit records: %d", len(records))
		}
		data = append(data, byte(len(records)))

		// Write each record
		for _, record := range records {
			var err error
			data, err = appendCardVehicleUnitRecord(data, record)
			if err != nil {
				return nil, fmt.Errorf("failed to append vehicle unit record: %w", err)
			}
		}
	} else {
		// Write 0 records
		data = append(data, 0x00)
	}

	return data, nil
}

// appendCardVehicleUnitRecord appends a single vehicle unit record to dst
func appendCardVehicleUnitRecord(dst []byte, record *cardv1.VehicleUnitsUsed_Record) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// Timestamp (TimeReal - 4 bytes)
	dst = dd.AppendTimeReal(dst, record.GetTimestamp())

	// Manufacturer code (1 byte)
	manufacturerCode := record.GetManufacturerCode()
	if manufacturerCode < 0 || manufacturerCode > 255 {
		return nil, fmt.Errorf("invalid manufacturer code: %d", manufacturerCode)
	}
	dst = append(dst, byte(manufacturerCode))

	// Device ID (1 byte)
	deviceID := record.GetDeviceId()
	if len(deviceID) > 1 {
		return nil, fmt.Errorf("device ID too long: %d bytes", len(deviceID))
	}
	if len(deviceID) == 1 {
		dst = append(dst, deviceID[0])
	} else {
		dst = append(dst, 0x00)
	}

	// VU software version (4 bytes)
	vuSoftwareVersion := record.GetVuSoftwareVersion()
	if len(vuSoftwareVersion) > 4 {
		return nil, fmt.Errorf("VU software version too long: %d bytes", len(vuSoftwareVersion))
	}
	if len(vuSoftwareVersion) == 4 {
		dst = append(dst, vuSoftwareVersion...)
	} else {
		// Pad with zeros if shorter than 4 bytes
		padded := make([]byte, 4)
		copy(padded, vuSoftwareVersion)
		dst = append(dst, padded...)
	}

	return dst, nil
}
