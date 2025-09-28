package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalCardControlActivityData unmarshals control activity data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.15):
//
//	CardControlActivityDataRecord ::= SEQUENCE {
//	    controlType                        ControlType,
//	    controlTime                        TimeReal,
//	    controlCardNumber                  FullCardNumber,
//	    controlVehicleRegistration         VehicleRegistrationIdentification,
//	    controlDownloadPeriodBegin         TimeReal,
//	    controlDownloadPeriodEnd           TimeReal
//	}
//
// Binary Layout (46 bytes):
//
//	0-0:   controlType (1 byte)
//	1-4:   controlTime (4 bytes, TimeReal)
//	5-22:  controlCardNumber (18 bytes, FullCardNumber)
//	23-37: controlVehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	38-41: controlDownloadPeriodBegin (4 bytes, TimeReal)
//	42-45: controlDownloadPeriodEnd (4 bytes, TimeReal)
//
// Constants:
const (
	// CardControlActivityDataRecord total size
	cardControlActivityDataRecordSize = 46
)

func unmarshalCardControlActivityData(data []byte) (*cardv1.ControlActivityData, error) {
	if len(data) < cardControlActivityDataRecordSize {
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
	fullCardNumber := &datadictionaryv1.FullCardNumber{}
	fullCardNumber.SetCardType(datadictionaryv1.EquipmentType_DRIVER_CARD)

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
	driverID := &datadictionaryv1.FullCardNumber_DriverIdentification{}
	driverID.SetIdentification(cardNumberStr)
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
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
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
	offset += 4

	return &target, nil
}
