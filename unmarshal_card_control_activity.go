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
	var controlType byte
	if err := binary.Read(r, binary.BigEndian, &controlType); err != nil {
		return nil, fmt.Errorf("failed to read control type: %w", err)
	}
	target.SetControlType([]byte{controlType})
	// Read control time (4 bytes)
	target.SetControlTime(readTimeReal(r))
	// Read control card number (18 bytes)
	controlCardBytes := make([]byte, 18)
	if _, err := r.Read(controlCardBytes); err != nil {
		return nil, fmt.Errorf("failed to read control card number: %w", err)
	}
	// Create FullCardNumber structure
	fullCardNumber := &datadictionaryv1.FullCardNumber{}
	fullCardNumber.SetCardNumber(readString(bytes.NewReader(controlCardBytes), 18))
	target.SetControlCardNumber(fullCardNumber)
	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	var nation byte
	if err := binary.Read(r, binary.BigEndian, &nation); err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	// Create VehicleRegistrationIdentification structure
	vehicleReg := &datadictionaryv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(int32(nation))
	vehicleReg.SetNumber(readString(r, 14))
	target.SetControlVehicleRegistration(vehicleReg)
	// Read control download period begin (4 bytes)
	target.SetControlDownloadPeriodBegin(readTimeReal(r))
	// Read control download period end (4 bytes)
	target.SetControlDownloadPeriodEnd(readTimeReal(r))
	return &target, nil
}
