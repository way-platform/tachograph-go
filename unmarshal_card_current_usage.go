package tachograph

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalCardCurrentUsage unmarshals current usage data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.16):
//
//	CardCurrentUse ::= SEQUENCE {
//	    sessionOpenTime                   TimeReal,
//	    sessionOpenVehicle                VehicleRegistrationIdentification
//	}
//
// Binary Layout (19 bytes):
//
//	0-3:   sessionOpenTime (4 bytes, TimeReal)
//	4-18:  sessionOpenVehicle (15 bytes: 1 byte nation + 14 bytes number)
//
// Constants:
const (
	// CardCurrentUse total size
	cardCurrentUseSize = 19
)

func unmarshalCardCurrentUsage(data []byte) (*cardv1.CurrentUsage, error) {
	if len(data) < cardCurrentUseSize {
		return nil, fmt.Errorf("insufficient data for current usage")
	}
	var target cardv1.CurrentUsage
	offset := 0

	// Read session open time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for session open time")
	}
	target.SetSessionOpenTime(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read session open vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
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
	// offset += 14 // Not needed as this is the last field
	vehicleReg.SetNumber(regNumber)
	target.SetSessionOpenVehicle(vehicleReg)
	return &target, nil
}
