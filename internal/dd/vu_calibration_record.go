package dd

import (
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// UnmarshalVuCalibrationRecord parses the VuCalibrationRecord structure.
//
// See Data Dictionary, Section 2.174, `VuCalibrationRecord`.
//
// ASN.1 Specification (Generation 1):
//
//	VuCalibrationRecord ::= SEQUENCE {
//	    calibrationPurpose CalibrationPurpose,                               -- 1 byte
//	    workshopName Name,                                                   -- 36 bytes
//	    workshopAddress Address,                                             -- 36 bytes
//	    workshopCardNumber FullCardNumber,                                   -- 18 bytes
//	    workshopCardExpiryDate Datef,                                        -- 4 bytes
//	    vehicleIdentificationNumber VehicleIdentificationNumber,             -- 17 bytes
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification, -- 15 bytes
//	    wVehicleCharacteristicConstant W-VehicleCharacteristicConstant,     -- 2 bytes
//	    kConstantOfRecordingEquipment K-ConstantOfRecordingEquipment,       -- 2 bytes
//	    lTyreCircumference L-TyreCircumference,                              -- 2 bytes
//	    tyreSize TyreSize,                                                   -- 15 bytes
//	    authorisedSpeed SpeedAuthorised,                                     -- 1 byte
//	    oldOdometerValue OdometerShort,                                      -- 3 bytes
//	    newOdometerValue OdometerShort,                                      -- 3 bytes
//	    oldTimeValue TimeReal,                                               -- 4 bytes
//	    newTimeValue TimeReal,                                               -- 4 bytes
//	    nextCalibrationDate TimeReal                                         -- 4 bytes
//	}
//
// Binary Layout (Generation 1, fixed length: 167 bytes)
func (opts UnmarshalOptions) UnmarshalVuCalibrationRecord(data []byte) (*ddv1.VuCalibrationRecord, error) {
	const lenVuCalibrationRecord = 167

	if len(data) != lenVuCalibrationRecord {
		return nil, fmt.Errorf(
			"invalid data length for VuCalibrationRecord: got %d, want %d",
			len(data), lenVuCalibrationRecord,
		)
	}

	record := &ddv1.VuCalibrationRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	const (
		idxCalibrationPurpose      = 0
		lenCalibrationPurpose      = 1
		idxWorkshopName            = 1
		lenWorkshopName            = 36
		idxWorkshopAddress         = 37
		lenWorkshopAddress         = 36
		idxWorkshopCardNumber      = 73
		lenWorkshopCardNumber      = 18
		idxWorkshopCardExpiryDate  = 91
		lenWorkshopCardExpiryDate  = 4
		idxVIN                     = 95
		lenVIN                     = 17
		idxVehicleRegistration     = 112
		lenVehicleRegistration     = 15
		idxWVehicleCharConstant    = 127
		lenWVehicleCharConstant    = 2
		idxKConstantRecordingEquip = 129
		lenKConstantRecordingEquip = 2
		idxLTyreCircumference      = 131
		lenLTyreCircumference      = 2
		idxTyreSize                = 133
		lenTyreSize                = 15
		idxAuthorisedSpeed         = 148
		lenAuthorisedSpeed         = 1
		idxOldOdometerValue        = 149
		lenOldOdometerValue        = 3
		idxNewOdometerValue        = 152
		lenNewOdometerValue        = 3
		idxOldTimeValue            = 155
		lenOldTimeValue            = 4
		idxNewTimeValue            = 159
		lenNewTimeValue            = 4
		idxNextCalibrationDate     = 163
		lenNextCalibrationDate     = 4
	)

	// Parse calibration purpose (1 byte)
	calibrationPurposeValue := int32(data[idxCalibrationPurpose])
	calibrationPurpose := ddv1.CalibrationPurpose(calibrationPurposeValue)
	valueDesc := calibrationPurpose.Descriptor().Values().ByNumber(protoreflect.EnumNumber(calibrationPurposeValue))
	if valueDesc == nil {
		record.SetUnrecognizedPurpose(calibrationPurposeValue)
		calibrationPurpose = ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED
	}
	record.SetPurpose(calibrationPurpose)

	// Parse workshop name (36 bytes)
	workshopName, err := opts.UnmarshalStringValue(
		data[idxWorkshopName : idxWorkshopName+lenWorkshopName],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop name: %w", err)
	}
	record.SetWorkshopName(workshopName)

	// Parse workshop address (36 bytes)
	workshopAddress, err := opts.UnmarshalStringValue(
		data[idxWorkshopAddress : idxWorkshopAddress+lenWorkshopAddress],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop address: %w", err)
	}
	record.SetWorkshopAddress(workshopAddress)

	// Parse workshop card number (18 bytes)
	workshopCardNumber, err := opts.UnmarshalFullCardNumber(
		data[idxWorkshopCardNumber : idxWorkshopCardNumber+lenWorkshopCardNumber],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop card number: %w", err)
	}
	record.SetWorkshopCardNumber(workshopCardNumber)

	// Parse workshop card expiry date (4 bytes)
	workshopCardExpiryDate, err := opts.UnmarshalDate(
		data[idxWorkshopCardExpiryDate : idxWorkshopCardExpiryDate+lenWorkshopCardExpiryDate],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop card expiry date: %w", err)
	}
	record.SetWorkshopCardExpiryDate(workshopCardExpiryDate)

	// Parse VIN (17 bytes)
	vin, err := opts.UnmarshalIa5StringValue(data[idxVIN : idxVIN+lenVIN])
	if err != nil {
		return nil, fmt.Errorf("failed to parse VIN: %w", err)
	}
	record.SetVin(vin)

	// Parse vehicle registration (15 bytes)
	vehicleRegistration, err := opts.UnmarshalVehicleRegistrationIdentification(
		data[idxVehicleRegistration : idxVehicleRegistration+lenVehicleRegistration],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration: %w", err)
	}
	record.SetVehicleRegistration(vehicleRegistration)

	// Parse W-Vehicle Characteristic Constant (2 bytes)
	wVehicleCharConstant := int32(binary.BigEndian.Uint16(
		data[idxWVehicleCharConstant : idxWVehicleCharConstant+lenWVehicleCharConstant],
	))
	record.SetWVehicleCharacteristicConstant(wVehicleCharConstant)

	// Parse K-Constant of Recording Equipment (2 bytes)
	kConstantRecordingEquip := int32(binary.BigEndian.Uint16(
		data[idxKConstantRecordingEquip : idxKConstantRecordingEquip+lenKConstantRecordingEquip],
	))
	record.SetKConstantOfRecordingEquipment(kConstantRecordingEquip)

	// Parse L-Tyre Circumference (2 bytes)
	lTyreCircumference := int32(binary.BigEndian.Uint16(
		data[idxLTyreCircumference : idxLTyreCircumference+lenLTyreCircumference],
	))
	record.SetLTyreCircumferenceEighthsMm(lTyreCircumference)

	// Parse tyre size (15 bytes)
	tyreSize, err := opts.UnmarshalIa5StringValue(data[idxTyreSize : idxTyreSize+lenTyreSize])
	if err != nil {
		return nil, fmt.Errorf("failed to parse tyre size: %w", err)
	}
	record.SetTyreSize(tyreSize)

	// Parse authorised speed (1 byte)
	authorisedSpeed := int32(data[idxAuthorisedSpeed])
	record.SetAuthorisedSpeedKmh(authorisedSpeed)

	// Parse old odometer value (3 bytes - 24-bit integer)
	oldOdometerValue := int32(data[idxOldOdometerValue])<<16 |
		int32(data[idxOldOdometerValue+1])<<8 |
		int32(data[idxOldOdometerValue+2])
	record.SetOldOdometerValueKm(oldOdometerValue)

	// Parse new odometer value (3 bytes - 24-bit integer)
	newOdometerValue := int32(data[idxNewOdometerValue])<<16 |
		int32(data[idxNewOdometerValue+1])<<8 |
		int32(data[idxNewOdometerValue+2])
	record.SetNewOdometerValueKm(newOdometerValue)

	// Parse old time value (4 bytes)
	oldTimeValue, err := opts.UnmarshalTimeReal(
		data[idxOldTimeValue : idxOldTimeValue+lenOldTimeValue],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse old time value: %w", err)
	}
	record.SetOldTimeValue(oldTimeValue)

	// Parse new time value (4 bytes)
	newTimeValue, err := opts.UnmarshalTimeReal(
		data[idxNewTimeValue : idxNewTimeValue+lenNewTimeValue],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new time value: %w", err)
	}
	record.SetNewTimeValue(newTimeValue)

	// Parse next calibration date (4 bytes)
	nextCalibrationDate, err := opts.UnmarshalTimeReal(
		data[idxNextCalibrationDate : idxNextCalibrationDate+lenNextCalibrationDate],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse next calibration date: %w", err)
	}
	record.SetNextCalibrationDate(nextCalibrationDate)

	return record, nil
}

// MarshalVuCalibrationRecord marshals the VuCalibrationRecord structure using raw data painting.
//
// See Data Dictionary, Section 2.174, `VuCalibrationRecord`.
func (opts MarshalOptions) MarshalVuCalibrationRecord(record *ddv1.VuCalibrationRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const size = 167

	// Use raw data painting strategy
	var canvas [size]byte
	if raw := record.GetRawData(); len(raw) > 0 {
		if len(raw) != size {
			return nil, fmt.Errorf(
				"invalid raw_data length for VuCalibrationRecord: got %d, want %d",
				len(raw), size,
			)
		}
		copy(canvas[:], raw)
	}

	offset := 0

	// Marshal calibration purpose (1 byte)
	var purposeValue int32
	if record.GetUnrecognizedPurpose() != 0 {
		purposeValue = record.GetUnrecognizedPurpose()
	} else {
		purposeValue = int32(record.GetPurpose())
	}
	canvas[offset] = byte(purposeValue)
	offset += 1

	// Marshal workshop name (36 bytes)
	workshopNameBytes, err := opts.MarshalStringValue(record.GetWorkshopName())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop name: %w", err)
	}
	if len(workshopNameBytes) != 36 {
		return nil, fmt.Errorf("invalid workshop name length: got %d, want 36", len(workshopNameBytes))
	}
	copy(canvas[offset:offset+36], workshopNameBytes)
	offset += 36

	// Marshal workshop address (36 bytes)
	workshopAddressBytes, err := opts.MarshalStringValue(record.GetWorkshopAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop address: %w", err)
	}
	if len(workshopAddressBytes) != 36 {
		return nil, fmt.Errorf("invalid workshop address length: got %d, want 36", len(workshopAddressBytes))
	}
	copy(canvas[offset:offset+36], workshopAddressBytes)
	offset += 36

	// Marshal workshop card number (18 bytes)
	workshopCardNumberBytes, err := opts.MarshalFullCardNumber(record.GetWorkshopCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop card number: %w", err)
	}
	if len(workshopCardNumberBytes) != 18 {
		return nil, fmt.Errorf("invalid workshop card number length: got %d, want 18", len(workshopCardNumberBytes))
	}
	copy(canvas[offset:offset+18], workshopCardNumberBytes)
	offset += 18

	// Marshal workshop card expiry date (4 bytes)
	workshopCardExpiryDateBytes, err := opts.MarshalDate(record.GetWorkshopCardExpiryDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop card expiry date: %w", err)
	}
	if len(workshopCardExpiryDateBytes) != 4 {
		return nil, fmt.Errorf("invalid workshop card expiry date length: got %d, want 4", len(workshopCardExpiryDateBytes))
	}
	copy(canvas[offset:offset+4], workshopCardExpiryDateBytes)
	offset += 4

	// Marshal VIN (17 bytes)
	vinBytes, err := opts.MarshalIa5StringValue(record.GetVin())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VIN: %w", err)
	}
	if len(vinBytes) != 17 {
		return nil, fmt.Errorf("invalid VIN length: got %d, want 17", len(vinBytes))
	}
	copy(canvas[offset:offset+17], vinBytes)
	offset += 17

	// Marshal vehicle registration (15 bytes)
	vehicleRegistrationBytes, err := opts.MarshalVehicleRegistrationIdentification(record.GetVehicleRegistration())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle registration: %w", err)
	}
	if len(vehicleRegistrationBytes) != 15 {
		return nil, fmt.Errorf("invalid vehicle registration length: got %d, want 15", len(vehicleRegistrationBytes))
	}
	copy(canvas[offset:offset+15], vehicleRegistrationBytes)
	offset += 15

	// Marshal W-Vehicle Characteristic Constant (2 bytes)
	binary.BigEndian.PutUint16(
		canvas[offset:offset+2],
		uint16(record.GetWVehicleCharacteristicConstant()),
	)
	offset += 2

	// Marshal K-Constant of Recording Equipment (2 bytes)
	binary.BigEndian.PutUint16(
		canvas[offset:offset+2],
		uint16(record.GetKConstantOfRecordingEquipment()),
	)
	offset += 2

	// Marshal L-Tyre Circumference (2 bytes)
	binary.BigEndian.PutUint16(
		canvas[offset:offset+2],
		uint16(record.GetLTyreCircumferenceEighthsMm()),
	)
	offset += 2

	// Marshal tyre size (15 bytes)
	tyreSizeBytes, err := opts.MarshalIa5StringValue(record.GetTyreSize())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tyre size: %w", err)
	}
	if len(tyreSizeBytes) != 15 {
		return nil, fmt.Errorf("invalid tyre size length: got %d, want 15", len(tyreSizeBytes))
	}
	copy(canvas[offset:offset+15], tyreSizeBytes)
	offset += 15

	// Marshal authorised speed (1 byte)
	canvas[offset] = byte(record.GetAuthorisedSpeedKmh())
	offset += 1

	// Marshal old odometer value (3 bytes - 24-bit integer)
	oldOdometer := record.GetOldOdometerValueKm()
	canvas[offset] = byte((oldOdometer >> 16) & 0xFF)
	canvas[offset+1] = byte((oldOdometer >> 8) & 0xFF)
	canvas[offset+2] = byte(oldOdometer & 0xFF)
	offset += 3

	// Marshal new odometer value (3 bytes - 24-bit integer)
	newOdometer := record.GetNewOdometerValueKm()
	canvas[offset] = byte((newOdometer >> 16) & 0xFF)
	canvas[offset+1] = byte((newOdometer >> 8) & 0xFF)
	canvas[offset+2] = byte(newOdometer & 0xFF)
	offset += 3

	// Marshal old time value (4 bytes)
	oldTimeValueBytes, err := opts.MarshalTimeReal(record.GetOldTimeValue())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal old time value: %w", err)
	}
	if len(oldTimeValueBytes) != 4 && record.GetOldTimeValue() != nil {
		return nil, fmt.Errorf("invalid old time value length: got %d, want 4", len(oldTimeValueBytes))
	}
	if record.GetOldTimeValue() != nil {
		copy(canvas[offset:offset+4], oldTimeValueBytes)
	}
	offset += 4

	// Marshal new time value (4 bytes)
	newTimeValueBytes, err := opts.MarshalTimeReal(record.GetNewTimeValue())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new time value: %w", err)
	}
	if len(newTimeValueBytes) != 4 && record.GetNewTimeValue() != nil {
		return nil, fmt.Errorf("invalid new time value length: got %d, want 4", len(newTimeValueBytes))
	}
	if record.GetNewTimeValue() != nil {
		copy(canvas[offset:offset+4], newTimeValueBytes)
	}
	offset += 4

	// Marshal next calibration date (4 bytes)
	nextCalibrationDateBytes, err := opts.MarshalTimeReal(record.GetNextCalibrationDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal next calibration date: %w", err)
	}
	if len(nextCalibrationDateBytes) != 4 && record.GetNextCalibrationDate() != nil {
		return nil, fmt.Errorf("invalid next calibration date length: got %d, want 4", len(nextCalibrationDateBytes))
	}
	if record.GetNextCalibrationDate() != nil {
		copy(canvas[offset:offset+4], nextCalibrationDateBytes)
	}
	offset += 4

	if offset != size {
		return nil, fmt.Errorf(
			"VuCalibrationRecord marshalling size mismatch: wrote %d bytes, expected %d",
			offset, size,
		)
	}

	return canvas[:], nil
}

// AnonymizeVuCalibrationRecord anonymizes a VU calibration record.
func (opts AnonymizeOptions) AnonymizeVuCalibrationRecord(rec *ddv1.VuCalibrationRecord) *ddv1.VuCalibrationRecord {
	if rec == nil {
		return nil
	}

	result := proto.Clone(rec).(*ddv1.VuCalibrationRecord)

	// Anonymize workshop name
	result.SetWorkshopName(opts.AnonymizeStringValue(rec.GetWorkshopName()))

	// Anonymize workshop address
	result.SetWorkshopAddress(opts.AnonymizeStringValue(rec.GetWorkshopAddress()))

	// Anonymize workshop card number
	result.SetWorkshopCardNumber(opts.AnonymizeFullCardNumber(rec.GetWorkshopCardNumber()))

	// Anonymize workshop card expiry date (preserve or anonymize timestamp)
	if expiryDate := rec.GetWorkshopCardExpiryDate(); expiryDate != nil && !opts.PreserveTimestamps {
		// Create test expiry date
		result.SetWorkshopCardExpiryDate(NewDate(2025, 12, 31))
	}

	// Anonymize VIN
	result.SetVin(opts.AnonymizeIa5StringValue(rec.GetVin()))

	// Anonymize vehicle registration
	result.SetVehicleRegistration(opts.AnonymizeVehicleRegistrationIdentification(rec.GetVehicleRegistration()))

	// Anonymize odometer values
	result.SetOldOdometerValueKm(opts.AnonymizeOdometerValue(rec.GetOldOdometerValueKm()))
	result.SetNewOdometerValueKm(opts.AnonymizeOdometerValue(rec.GetNewOdometerValueKm()))

	// Anonymize timestamps
	result.SetOldTimeValue(opts.AnonymizeTimestamp(rec.GetOldTimeValue()))
	result.SetNewTimeValue(opts.AnonymizeTimestamp(rec.GetNewTimeValue()))
	result.SetNextCalibrationDate(opts.AnonymizeTimestamp(rec.GetNextCalibrationDate()))

	// Keep technical values (not PII): w_vehicle_characteristic_constant,
	// k_constant_of_recording_equipment, l_tyre_circumference, tyre_size, authorised_speed

	// Clear raw_data
	result.ClearRawData()

	return result
}
