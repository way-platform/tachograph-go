package card

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

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
	controlType, err := dd.UnmarshalControlType(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read control type: %w", err)
	}
	target.SetControlType(controlType)
	offset++

	// Read control time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control time")
	}
	controlTimestamp, err2 := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err2 != nil {
		return nil, fmt.Errorf("failed to parse control time: %w", err2)
	}
	target.SetControlTime(controlTimestamp)
	offset += 4

	// Read control card number (18 bytes) - this should be parsed as a proper FullCardNumberAndGeneration
	// For now, create a basic structure - this needs proper protocol parsing
	fullCardNumberAndGeneration := &ddv1.FullCardNumberAndGeneration{}

	fullCardNumber := &ddv1.FullCardNumber{}
	fullCardNumber.SetCardType(ddv1.EquipmentType_DRIVER_CARD)

	// Read the card number as IA5 string
	if offset+18 > len(data) {
		return nil, fmt.Errorf("insufficient data for control card number")
	}
	cardNumberStr, err := dd.UnmarshalIA5StringValue(data[offset : offset+18])
	if err != nil {
		return nil, fmt.Errorf("failed to read control card number: %w", err)
	}
	offset += 18

	// Create driver identification with the card number
	driverID := &ddv1.DriverIdentification{}
	driverID.SetDriverIdentificationNumber(cardNumberStr)
	fullCardNumber.SetDriverIdentification(driverID)

	// Set the full card number in the generation wrapper
	fullCardNumberAndGeneration.SetFullCardNumber(fullCardNumber)
	// Default to Generation 1 for now - this should be determined from context
	fullCardNumberAndGeneration.SetGeneration(ddv1.Generation_GENERATION_1)

	target.SetControlCardNumber(fullCardNumberAndGeneration)

	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	if offset+15 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration")
	}
	vehicleReg, err := dd.UnmarshalVehicleRegistration(data[offset : offset+15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration: %w", err)
	}
	offset += 15
	target.SetControlVehicleRegistration(vehicleReg)

	// Read control download period begin (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control download period begin")
	}
	controlDownloadPeriodBegin, err3 := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err3 != nil {
		return nil, fmt.Errorf("failed to parse control download period begin: %w", err3)
	}
	target.SetControlDownloadPeriodBegin(controlDownloadPeriodBegin)
	offset += 4

	// Read control download period end (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control download period end")
	}
	controlDownloadPeriodEnd, err4 := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err4 != nil {
		return nil, fmt.Errorf("failed to parse control download period end: %w", err4)
	}
	target.SetControlDownloadPeriodEnd(controlDownloadPeriodEnd)
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
	data = dd.AppendTimeReal(data, controlData.GetControlTime())

	var err error
	// Control card number (18 bytes)
	data, err = dd.AppendFullCardNumberAsString(data, controlData.GetControlCardNumber().GetFullCardNumber(), 18)
	if err != nil {
		return nil, err
	}

	// Vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	data, err = dd.AppendVehicleRegistration(data, controlData.GetControlVehicleRegistration())
	if err != nil {
		return nil, err
	}

	// Control download period begin (4 bytes)
	data = dd.AppendTimeReal(data, controlData.GetControlDownloadPeriodBegin())

	// Control download period end (4 bytes)
	data = dd.AppendTimeReal(data, controlData.GetControlDownloadPeriodEnd())

	return data, nil
}
