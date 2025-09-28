package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalCardControlActivityData unmarshals control activity data from a card EF.
//
// The data type `CardControlActivityDataRecord` is specified in the Data Dictionary, Section 2.15.
//
// ASN.1 Definition:
//
//	CardControlActivityDataRecord ::= SEQUENCE {
//	    controlType                        ControlType,
//	    controlTime                        TimeReal,
//	    controlCardNumber                  FullCardNumber,
//	    controlVehicleRegistration         VehicleRegistrationIdentification,
//	    controlDownloadPeriodBegin         TimeReal,
//	    controlDownloadPeriodEnd           TimeReal
//	}
func unmarshalCardControlActivityData(data []byte) (*cardv1.ControlActivityData, error) {
	const (
		lenCardControlActivityDataRecord = 46 // CardControlActivityDataRecord total size
	)

	if len(data) < lenCardControlActivityDataRecord {
		return nil, fmt.Errorf("insufficient data for control activity data")
	}
	var target cardv1.ControlActivityData
	controlTime := binary.BigEndian.Uint32(data[1:5])
	if controlTime == 0 {
		target.SetValid(false)
		target.SetRawData(data)
		return &target, nil
	}
	target.SetValid(true)

	offset := 0

	// Read control type (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for control type")
	}
	controlType, err := unmarshalControlType(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read control type: %w", err)
	}
	target.SetControlType(controlType)
	offset++

	// Read control time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control time")
	}
	target.SetControlTime(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read control card number (18 bytes) - this should be parsed as a proper FullCardNumber
	// For now, create a basic structure - this needs proper protocol parsing
	fullCardNumber := &ddv1.FullCardNumber{}
	fullCardNumber.SetCardType(ddv1.EquipmentType_DRIVER_CARD)

	// Read the card number as IA5 string
	if offset+18 > len(data) {
		return nil, fmt.Errorf("insufficient data for control card number")
	}
	cardNumberStr, err := unmarshalIA5StringValue(data[offset : offset+18])
	if err != nil {
		return nil, fmt.Errorf("failed to read control card number: %w", err)
	}
	offset += 18

	// Create driver identification with the card number
	driverID := &ddv1.DriverIdentification{}
	driverID.SetIdentificationNumber(cardNumberStr)
	fullCardNumber.SetDriverIdentification(driverID)
	target.SetControlCardNumber(fullCardNumber)

	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration nation")
	}
	nation, err := unmarshalNationNumeric(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	offset++

	// Create VehicleRegistrationIdentification structure
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)

	if offset+14 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration number")
	}
	regNumber, err := unmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	offset += 14
	vehicleReg.SetNumber(regNumber)
	target.SetControlVehicleRegistration(vehicleReg)

	// Read control download period begin (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control download period begin")
	}
	target.SetControlDownloadPeriodBegin(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read control download period end (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control download period end")
	}
	target.SetControlDownloadPeriodEnd(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	// offset += 4 // Not needed as this is the last field

	return &target, nil
}

// AppendCardControlActivityData appends control activity data to a byte slice.
//
// The data type `CardControlActivityDataRecord` is specified in the Data Dictionary, Section 2.15.
//
// ASN.1 Definition:
//
//	CardControlActivityDataRecord ::= SEQUENCE {
//	    controlType                        ControlType,
//	    controlTime                        TimeReal,
//	    controlCardNumber                  FullCardNumber,
//	    controlVehicleRegistration         VehicleRegistrationIdentification,
//	    controlDownloadPeriodBegin         TimeReal,
//	    controlDownloadPeriodEnd           TimeReal
//	}
func appendCardControlActivityData(data []byte, controlData *cardv1.ControlActivityData) ([]byte, error) {
	if controlData == nil {
		return data, nil
	}

	if !controlData.GetValid() {
		// Non-valid record: use preserved raw data
		rawData := controlData.GetRawData()
		if len(rawData) != 46 {
			// Fallback to zeros if raw data is invalid
			return append(data, make([]byte, 46)...), nil
		}
		return append(data, rawData...), nil
	}

	// Valid record: serialize semantic data
	// Control type (1 byte)
	controlType := controlData.GetControlType()
	var controlTypeByte byte
	if controlType != nil {
		// Build bitmask from boolean fields
		// Structure: 'cvpdexxx'B
		// - 'c': card downloading
		// - 'v': VU downloading
		// - 'p': printing
		// - 'd': display
		// - 'e': calibration checking
		if controlType.GetCardDownloading() {
			controlTypeByte |= 0x80 // bit 7
		}
		if controlType.GetVuDownloading() {
			controlTypeByte |= 0x40 // bit 6
		}
		if controlType.GetPrinting() {
			controlTypeByte |= 0x20 // bit 5
		}
		if controlType.GetDisplay() {
			controlTypeByte |= 0x10 // bit 4
		}
		if controlType.GetCalibrationChecking() {
			controlTypeByte |= 0x08 // bit 3
		}
	}
	data = append(data, controlTypeByte)

	// Control time (4 bytes)
	data = appendTimeReal(data, controlData.GetControlTime())

	var err error
	// Control card number (18 bytes)
	data, err = appendFullCardNumber(data, controlData.GetControlCardNumber(), 18)
	if err != nil {
		return nil, err
	}

	// Vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	data, err = appendVehicleRegistration(data, controlData.GetControlVehicleRegistration())
	if err != nil {
		return nil, err
	}

	// Control download period begin (4 bytes)
	data = appendTimeReal(data, controlData.GetControlDownloadPeriodBegin())

	// Control download period end (4 bytes)
	data = appendTimeReal(data, controlData.GetControlDownloadPeriodEnd())

	return data, nil
}
