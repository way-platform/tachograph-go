package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

func unmarshalCardControlActivityData(data []byte) (*cardv1.ControlActivityData, error) {
	if len(data) < 46 {
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
	r := bytes.NewReader(data)
	// Read control type (1 byte)
	controlType, err := unmarshalControlTypeFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read control type: %w", err)
	}
	target.SetControlType(controlType)
	// Read control time (4 bytes)
	target.SetControlTime(readTimeReal(r))
	// Read control card number (18 bytes) - this should be parsed as a proper FullCardNumber
	// For now, create a basic structure - this needs proper protocol parsing
	fullCardNumber := &datadictionaryv1.FullCardNumber{}
	fullCardNumber.SetCardType(datadictionaryv1.EquipmentType_DRIVER_CARD)

	// Read the card number as IA5 string
	cardNumberStr, err := unmarshalIA5StringValueFromReader(r, 18)
	if err != nil {
		return nil, fmt.Errorf("failed to read control card number: %w", err)
	}

	// Create driver identification with the card number
	driverID := &datadictionaryv1.FullCardNumber_DriverIdentification{}
	driverID.SetIdentification(cardNumberStr)
	fullCardNumber.SetDriverIdentification(driverID)
	target.SetControlCardNumber(fullCardNumber)
	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	nation, err := unmarshalNationNumericFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	// Create VehicleRegistrationIdentification structure
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)

	regNumber, err := unmarshalIA5StringValueFromReader(r, 14)
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	vehicleReg.SetNumber(regNumber)
	target.SetControlVehicleRegistration(vehicleReg)
	// Read control download period begin (4 bytes)
	target.SetControlDownloadPeriodBegin(readTimeReal(r))
	// Read control download period end (4 bytes)
	target.SetControlDownloadPeriodEnd(readTimeReal(r))
	return &target, nil
}
