package vu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalTechnicalDataGen2V2 parses Gen2 V2 Technical Data from the complete transfer value.
//
// Gen2 V2 adds SensorExternalGNSSCoupledRecordArray, VuITSConsentRecordArray,
// and VuPowerSupplyInterruptionRecordArray to the V1 structure.
//
// Structure:
//
//	VuTechnicalDataSecondGenV2 ::= SEQUENCE {
//	    vuIdentificationRecordArray              VuIdentificationRecordArray,
//	    vuSensorPairedRecordArray                VuSensorPairedRecordArray,
//	    vuSensorExternalGNSSCoupledRecordArray   VuSensorExternalGNSSCoupledRecordArray,
//	    vuCalibrationRecordArray                 VuCalibrationRecordArray,
//	    vuITSConsentRecordArray                  VuITSConsentRecordArray,
//	    vuPowerSupplyInterruptionRecordArray      VuPowerSupplyInterruptionRecordArray,
//	    signatureRecordArray                     SignatureRecordArray
//	}
func unmarshalTechnicalDataGen2V2(value []byte) (*vuv1.TechnicalDataGen2V2, error) {
	totalSize, signatureSize, err := sizeOfTechnicalDataGen2V2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	td := &vuv1.TechnicalDataGen2V2{}
	td.SetRawData(value)
	offset := 0

	// VuIdentificationRecordArray
	ddIdent, bytesRead, err := parseVuIdentificationRecordArrayGen2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuIdentificationRecordArray: %w", err)
	}
	td.SetVuIdentification(vuIdentToGen2V2(ddIdent))
	offset += bytesRead

	// SensorPairedRecordArray
	ddSensors, bytesRead, err := parseSensorPairedRecordArrayGen2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse SensorPairedRecordArray: %w", err)
	}
	pairedSensors := make([]*vuv1.TechnicalDataGen2V2_PairedSensor, len(ddSensors))
	for i, s := range ddSensors {
		pairedSensors[i] = sensorPairedToGen2V2(s)
	}
	td.SetPairedSensors(pairedSensors)
	offset += bytesRead

	// SensorExternalGNSSCoupledRecordArray
	ddGnss, bytesRead, err := parseSensorExternalGNSSCoupledRecordArrayGen2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse SensorExternalGNSSCoupledRecordArray: %w", err)
	}
	td.SetCoupledGnssFacilities(ddGnss)
	offset += bytesRead

	// VuCalibrationRecordArray
	calRecords, bytesRead, err := parseCalibrationRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuCalibrationRecordArray: %w", err)
	}
	td.SetCalibrationRecords(calRecords)
	offset += bytesRead

	// VuITSConsentRecordArray
	itsConsents, bytesRead, err := parseItsConsentRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuITSConsentRecordArray: %w", err)
	}
	td.SetItsConsentRecords(itsConsents)
	offset += bytesRead

	// VuPowerSupplyInterruptionRecordArray
	powerInterruptions, bytesRead, err := parsePowerSupplyInterruptionRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuPowerSupplyInterruptionRecordArray: %w", err)
	}
	td.SetPowerSupplyInterruptions(powerInterruptions)
	offset += bytesRead

	td.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Technical Data Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return td, nil
}

// MarshalTechnicalDataGen2V2 marshals Gen2 V2 Technical Data.
func (opts MarshalOptions) MarshalTechnicalDataGen2V2(td *vuv1.TechnicalDataGen2V2) ([]byte, error) {
	if td == nil {
		return nil, fmt.Errorf("technicalData cannot be nil")
	}

	raw := td.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	var result []byte
	marshalOpts := dd.MarshalOptions{}

	// VuIdentificationRecordArray (1 record × 124 bytes)
	vuIdentData, err := marshalVuIdentificationGen2V2(marshalOpts, td.GetVuIdentification())
	if err != nil {
		return nil, fmt.Errorf("marshal VuIdentificationRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x01, 124, 1)
	result = append(result, vuIdentData...)

	// SensorPairedRecordArray (N records × 28 bytes)
	sensors := td.GetPairedSensors()
	sensorData, err := marshalSensorPairedRecordsGen2V2(marshalOpts, sensors)
	if err != nil {
		return nil, fmt.Errorf("marshal SensorPairedRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x02, 28, uint16(len(sensors)))
	result = append(result, sensorData...)

	// SensorExternalGNSSCoupledRecordArray (N records × 28 bytes)
	gnss := td.GetCoupledGnssFacilities()
	gnssData, err := marshalCoupledGnssRecordsGen2V2(marshalOpts, gnss)
	if err != nil {
		return nil, fmt.Errorf("marshal SensorExternalGNSSCoupledRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, 28, uint16(len(gnss)))
	result = append(result, gnssData...)

	// VuCalibrationRecordArray (N records × 252 bytes)
	calRecords := td.GetCalibrationRecords()
	calData, err := marshalCalibrationRecordsGen2V2(marshalOpts, calRecords)
	if err != nil {
		return nil, fmt.Errorf("marshal VuCalibrationRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x04, 252, uint16(len(calRecords)))
	result = append(result, calData...)

	// VuITSConsentRecordArray (N records × 20 bytes)
	itsConsents := td.GetItsConsentRecords()
	itsData, err := marshalItsConsentRecordsGen2V2(marshalOpts, itsConsents)
	if err != nil {
		return nil, fmt.Errorf("marshal VuITSConsentRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x05, 20, uint16(len(itsConsents)))
	result = append(result, itsData...)

	// VuPowerSupplyInterruptionRecordArray (N records × 5 bytes)
	powerInterruptions := td.GetPowerSupplyInterruptions()
	powerData, err := marshalPowerSupplyInterruptionRecordsGen2V2(marshalOpts, powerInterruptions)
	if err != nil {
		return nil, fmt.Errorf("marshal VuPowerSupplyInterruptionRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x06, 5, uint16(len(powerInterruptions)))
	result = append(result, powerData...)

	// Signature: stored as complete SignatureRecordArray bytes (header + sig bytes).
	// When empty (anonymized data), include a placeholder header so sizeOf can parse the output.
	if sig := td.GetSignature(); len(sig) > 0 {
		result = append(result, sig...)
	} else {
		result = appendRecordArrayHeader(result, 0x07, 0, 0)
	}
	return result, nil
}

// anonymizeTechnicalDataGen2V2 anonymizes Gen2 V2 Technical Data.
func (opts AnonymizeOptions) anonymizeTechnicalDataGen2V2(td *vuv1.TechnicalDataGen2V2) *vuv1.TechnicalDataGen2V2 {
	if td == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &vuv1.TechnicalDataGen2V2{}

	// Anonymize VU identification
	if vuIdent := td.GetVuIdentification(); vuIdent != nil {
		anon := &vuv1.TechnicalDataGen2V2_VuIdentification{}
		anon.SetManufacturerName(ddOpts.AnonymizeStringValue(vuIdent.GetManufacturerName()))
		anon.SetManufacturerAddress(ddOpts.AnonymizeStringValue(vuIdent.GetManufacturerAddress()))
		anon.SetPartNumber(ddOpts.AnonymizeIa5StringValue(vuIdent.GetPartNumber()))
		if sn := vuIdent.GetSerialNumber(); sn != nil {
			anonSn := &ddv1.ExtendedSerialNumber{}
			anonSn.SetType(sn.GetType())
			anonSn.SetManufacturerCode(sn.GetManufacturerCode())
			anonSn.SetSerialNumber(0)
			anon.SetSerialNumber(anonSn)
		}
		anon.SetSoftwareIdentification(vuIdent.GetSoftwareIdentification())
		anon.SetManufacturingDate(vuIdent.GetManufacturingDate())
		anon.SetApprovalNumber(dd.NewIa5StringValue(16, "TEST0001"))
		result.SetVuIdentification(anon)
	}

	// Anonymize paired sensors
	anonSensors := make([]*vuv1.TechnicalDataGen2V2_PairedSensor, len(td.GetPairedSensors()))
	for i, sensor := range td.GetPairedSensors() {
		anon := &vuv1.TechnicalDataGen2V2_PairedSensor{}
		if sn := sensor.GetSerialNumber(); sn != nil {
			anonSn := &ddv1.ExtendedSerialNumber{}
			anonSn.SetType(sn.GetType())
			anonSn.SetManufacturerCode(sn.GetManufacturerCode())
			anonSn.SetSerialNumber(0)
			anon.SetSerialNumber(anonSn)
		}
		anon.SetApprovalNumber(dd.NewIa5StringValue(16, "SENSOR01"))
		anon.SetPairingDate(sensor.GetPairingDate())
		anonSensors[i] = anon
	}
	result.SetPairedSensors(anonSensors)

	// Anonymize coupled GNSS facilities
	anonGnss := make([]*vuv1.TechnicalDataGen2V2_CoupledGnss, len(td.GetCoupledGnssFacilities()))
	for i, gnss := range td.GetCoupledGnssFacilities() {
		anon := &vuv1.TechnicalDataGen2V2_CoupledGnss{}
		if sn := gnss.GetSerialNumber(); sn != nil {
			anonSn := &ddv1.ExtendedSerialNumber{}
			anonSn.SetType(sn.GetType())
			anonSn.SetManufacturerCode(sn.GetManufacturerCode())
			anonSn.SetSerialNumber(0)
			anon.SetSerialNumber(anonSn)
		}
		anon.SetApprovalNumber(dd.NewIa5StringValue(16, "GNSS0001"))
		anon.SetCouplingDate(gnss.GetCouplingDate())
		anonGnss[i] = anon
	}
	result.SetCoupledGnssFacilities(anonGnss)

	// Anonymize calibration records
	anonCals := make([]*vuv1.TechnicalDataGen2V2_CalibrationRecord, len(td.GetCalibrationRecords()))
	for i, cal := range td.GetCalibrationRecords() {
		anon := &vuv1.TechnicalDataGen2V2_CalibrationRecord{}
		anon.SetPurpose(cal.GetPurpose())
		anon.SetUnrecognizedPurpose(cal.GetUnrecognizedPurpose())
		anon.SetWorkshopName(ddOpts.AnonymizeStringValue(cal.GetWorkshopName()))
		anon.SetWorkshopAddress(ddOpts.AnonymizeStringValue(cal.GetWorkshopAddress()))
		anon.SetWorkshopCardNumber(ddOpts.AnonymizeFullCardNumber(cal.GetWorkshopCardNumber()))
		anon.SetWorkshopCardExpiryDate(cal.GetWorkshopCardExpiryDate())
		anon.SetVin(ddOpts.AnonymizeIa5StringValue(cal.GetVin()))
		anon.SetVehicleRegistration(ddOpts.AnonymizeVehicleRegistrationIdentification(cal.GetVehicleRegistration()))
		anon.SetWVehicleCharacteristicConstant(cal.GetWVehicleCharacteristicConstant())
		anon.SetKConstantOfRecordingEquipment(cal.GetKConstantOfRecordingEquipment())
		anon.SetLTyreCircumferenceEighthsMm(cal.GetLTyreCircumferenceEighthsMm())
		anon.SetTyreSize(ddOpts.AnonymizeIa5StringValue(cal.GetTyreSize()))
		anon.SetAuthorisedSpeedKmh(cal.GetAuthorisedSpeedKmh())
		anon.SetOldOdometerValueKm(ddOpts.AnonymizeOdometerValue(cal.GetOldOdometerValueKm()))
		anon.SetNewOdometerValueKm(ddOpts.AnonymizeOdometerValue(cal.GetNewOdometerValueKm()))
		anon.SetOldTimeValue(ddOpts.AnonymizeTimestamp(cal.GetOldTimeValue()))
		anon.SetNewTimeValue(ddOpts.AnonymizeTimestamp(cal.GetNewTimeValue()))
		anon.SetNextCalibrationDate(ddOpts.AnonymizeTimestamp(cal.GetNextCalibrationDate()))
		// V2 extension fields
		anon.SetSensorSerialNumber(anonymizeExtendedSerialNumber(cal.GetSensorSerialNumber()))
		anon.SetSensorGnssSerialNumber(anonymizeExtendedSerialNumber(cal.GetSensorGnssSerialNumber()))
		anon.SetRcmSerialNumber(anonymizeExtendedSerialNumber(cal.GetRcmSerialNumber()))
		anon.SetSealRecords(cal.GetSealRecords()) // Seal records have no PII
		anon.SetLoadType(cal.GetLoadType())
		anon.SetCalibrationCountry(cal.GetCalibrationCountry())
		anon.SetCalibrationCountryTimestamp(cal.GetCalibrationCountryTimestamp())
		anonCals[i] = anon
	}
	result.SetCalibrationRecords(anonCals)

	// Anonymize ITS consent records (clear card numbers)
	anonIts := make([]*vuv1.TechnicalDataGen2V2_ItsConsentRecord, len(td.GetItsConsentRecords()))
	for i, its := range td.GetItsConsentRecords() {
		anon := &vuv1.TechnicalDataGen2V2_ItsConsentRecord{}
		anon.SetFullCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(its.GetFullCardNumberAndGeneration()))
		anon.SetConsentStatus(its.GetConsentStatus())
		anonIts[i] = anon
	}
	result.SetItsConsentRecords(anonIts)

	// Preserve power supply interruption records (no PII — just timestamps and slot numbers)
	result.SetPowerSupplyInterruptions(td.GetPowerSupplyInterruptions())

	result.SetSignature([]byte{})
	return result
}

// ===== V2-specific parse helpers =====

// parseSensorExternalGNSSCoupledRecordArrayGen2 parses a SensorExternalGNSSCoupledRecordArray.
//
// Same binary layout as SensorPaired (28 bytes for Gen2):
//
//	sensorSerialNumber SensorSerialNumber            -- 8 bytes
//	sensorApprovalNumber SensorApprovalNumber        -- 16 bytes (Gen2 IA5String)
//	sensorCouplingDate SensorGNSSCouplingDate        -- 4 bytes
func parseSensorExternalGNSSCoupledRecordArrayGen2(data []byte, offset int) ([]*vuv1.TechnicalDataGen2V2_CoupledGnss, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 28
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected CoupledGnss record size %d, got %d", expectedRecordSize, recordSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.TechnicalDataGen2V2_CoupledGnss, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for CoupledGnss record %d", i)
		}
		// Reuse SensorPaired parser — same binary layout
		sensor, err := unmarshalOpts.UnmarshalSensorPaired(data[recStart:recEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("CoupledGnss record %d: %w", i, err)
		}
		gnss := &vuv1.TechnicalDataGen2V2_CoupledGnss{}
		gnss.SetSerialNumber(sensor.GetSerialNumber())
		gnss.SetApprovalNumber(sensor.GetApprovalNumber())
		gnss.SetCouplingDate(sensor.GetPairingDate())
		records = append(records, gnss)
		recStart = recEnd
	}

	return records, headerSize + int(recordSize)*int(noOfRecords), nil
}

// parseCalibrationRecordArrayGen2V2 parses a VuCalibrationRecordArray for Gen2 V2.
//
// Dispatches based on record size from the RecordArray header:
//   - 168 bytes: Gen2 V1 layout (FullCardNumberAndGeneration(19) + Datef expiry)
//   - 252 bytes: Gen2 V2 layout (FullCardNumber(18) + TimeReal expiry + V2 extensions)
//
// Some Gen2 V2 VUs send 168-byte calibration records (V1 layout), so we must handle both.
func parseCalibrationRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.TechnicalDataGen2V2_CalibrationRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.TechnicalDataGen2V2_CalibrationRecord, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for CalibrationRecord %d", i)
		}

		var parsed *vuv1.TechnicalDataGen2V2_CalibrationRecord
		switch {
		case recordSize >= 252:
			// Full V2 layout with extension fields
			parsed, err = parseOneCalibrationRecordGen2V2(unmarshalOpts, data[recStart:recEnd])
		case recordSize >= 168:
			// V1 layout (168 bytes) — convert to V2 proto type
			v1Rec, parseErr := parseOneCalibrationRecordGen2(unmarshalOpts, data[recStart:recEnd])
			if parseErr != nil {
				err = parseErr
			} else {
				parsed = gen2V1CalibrationToV2(v1Rec)
			}
		default:
			err = fmt.Errorf("unexpected CalibrationRecord size %d (want >= 168)", recordSize)
		}
		if err != nil {
			return nil, 0, fmt.Errorf("CalibrationRecord %d: %w", i, err)
		}
		records = append(records, parsed)
		recStart = recEnd
	}

	return records, headerSize + int(recordSize)*int(noOfRecords), nil
}

// parseOneCalibrationRecordGen2V2 parses a single Gen2 V2 VuCalibrationRecord.
//
// Gen2 V2 layout (252 bytes):
//
//	calibrationPurpose CalibrationPurpose                            -- 1 byte (offset 0)
//	workshopName Name                                                -- 36 bytes (offset 1)
//	workshopAddress Address                                          -- 36 bytes (offset 37)
//	workshopCardNumber FullCardNumber                                -- 18 bytes (offset 73)
//	workshopCardExpiryDate TimeReal                                  -- 4 bytes (offset 91)
//	vehicleIdentificationNumber VehicleIdentificationNumber          -- 17 bytes (offset 95)
//	vehicleRegistrationIdentification VehicleRegistrationIdentification -- 15 bytes (offset 112)
//	wVehicleCharacteristicConstant W-VehicleCharacteristicConstant   -- 2 bytes (offset 127)
//	kConstantOfRecordingEquipment K-ConstantOfRecordingEquipment     -- 2 bytes (offset 129)
//	lTyreCircumference L-TyreCircumference                           -- 2 bytes (offset 131)
//	tyreSize TyreSize                                                -- 15 bytes (offset 133)
//	authorisedSpeed SpeedAuthorised                                  -- 1 byte (offset 148)
//	oldOdometerValue OdometerShort                                   -- 3 bytes (offset 149)
//	newOdometerValue OdometerShort                                   -- 3 bytes (offset 152)
//	oldTimeValue TimeReal                                            -- 4 bytes (offset 155)
//	newTimeValue TimeReal                                            -- 4 bytes (offset 159)
//	nextCalibrationDate TimeReal                                     -- 4 bytes (offset 163)
//	sensorSerialNumber ExtendedSerialNumber                          -- 8 bytes (offset 167)
//	sensorGNSSSerialNumber ExtendedSerialNumber                      -- 8 bytes (offset 175)
//	rcmSerialNumber ExtendedSerialNumber                             -- 8 bytes (offset 183)
//	sealDataVu SealDataVu (5 × SealRecord(11))                      -- 55 bytes (offset 191)
//	byDefaultLoadType LoadType                                       -- 1 byte (offset 246)
//	calibrationCountry NationNumeric                                 -- 1 byte (offset 247)
//	calibrationCountryTimestamp TimeReal                              -- 4 bytes (offset 248)
func parseOneCalibrationRecordGen2V2(opts dd.UnmarshalOptions, data []byte) (*vuv1.TechnicalDataGen2V2_CalibrationRecord, error) {
	const (
		idxPurpose          = 0
		idxWorkshopName     = 1
		lenWorkshopName     = 36
		idxWorkshopAddress  = 37
		lenWorkshopAddress  = 36
		idxWorkshopCard     = 73
		lenWorkshopCard     = 18
		idxCardExpiry       = 91
		lenCardExpiry       = 4
		idxVIN              = 95
		lenVIN              = 17
		idxVehicleReg       = 112
		lenVehicleReg       = 15
		idxWVehicleChar     = 127
		idxKConstant        = 129
		idxLTyreCirc        = 131
		idxTyreSize         = 133
		lenTyreSize         = 15
		idxAuthorisedSpeed  = 148
		idxOldOdometer      = 149
		idxNewOdometer      = 152
		idxOldTimeValue     = 155
		idxNewTimeValue     = 159
		idxNextCalDate      = 163
		idxSensorSerial     = 167
		idxGNSSSerial       = 175
		idxRCMSerial        = 183
		idxSealData         = 191
		lenSealRecord       = 11
		numSealRecords      = 5
		lenSealData         = numSealRecords * lenSealRecord // 55
		idxLoadType         = idxSealData + lenSealData      // 246
		idxCalCountry       = 247
		idxCalCountryTs     = 248
		lenRecord           = 252
	)

	if len(data) < lenRecord {
		return nil, fmt.Errorf("invalid Gen2V2 CalibrationRecord size: got %d, want >= %d", len(data), lenRecord)
	}

	rec := &vuv1.TechnicalDataGen2V2_CalibrationRecord{}

	// calibrationPurpose (1 byte)
	purposeValue := int32(data[idxPurpose])
	purpose, err := dd.UnmarshalEnum[ddv1.CalibrationPurpose](byte(purposeValue))
	if err != nil {
		rec.SetUnrecognizedPurpose(purposeValue)
		rec.SetPurpose(ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED)
	} else {
		rec.SetPurpose(purpose)
	}

	// workshopName (36 bytes)
	workshopName, err := opts.UnmarshalStringValue(data[idxWorkshopName : idxWorkshopName+lenWorkshopName])
	if err != nil {
		return nil, fmt.Errorf("workshop name: %w", err)
	}
	rec.SetWorkshopName(workshopName)

	// workshopAddress (36 bytes)
	workshopAddr, err := opts.UnmarshalStringValue(data[idxWorkshopAddress : idxWorkshopAddress+lenWorkshopAddress])
	if err != nil {
		return nil, fmt.Errorf("workshop address: %w", err)
	}
	rec.SetWorkshopAddress(workshopAddr)

	// workshopCardNumber (18 bytes, FullCardNumber)
	workshopCard, err := opts.UnmarshalFullCardNumber(data[idxWorkshopCard : idxWorkshopCard+lenWorkshopCard])
	if err != nil {
		return nil, fmt.Errorf("workshop card number: %w", err)
	}
	rec.SetWorkshopCardNumber(workshopCard)

	// workshopCardExpiryDate (4 bytes, TimeReal)
	expiryDate, err := opts.UnmarshalTimeReal(data[idxCardExpiry : idxCardExpiry+lenCardExpiry])
	if err != nil {
		return nil, fmt.Errorf("workshop card expiry date: %w", err)
	}
	rec.SetWorkshopCardExpiryDate(expiryDate)

	// VIN (17 bytes)
	vin, err := opts.UnmarshalIa5StringValue(data[idxVIN : idxVIN+lenVIN])
	if err != nil {
		return nil, fmt.Errorf("VIN: %w", err)
	}
	rec.SetVin(vin)

	// vehicleRegistration (15 bytes)
	vehicleReg, err := opts.UnmarshalVehicleRegistrationIdentification(data[idxVehicleReg : idxVehicleReg+lenVehicleReg])
	if err != nil {
		return nil, fmt.Errorf("vehicle registration: %w", err)
	}
	rec.SetVehicleRegistration(vehicleReg)

	// W-Vehicle Characteristic Constant (2 bytes)
	rec.SetWVehicleCharacteristicConstant(int32(binary.BigEndian.Uint16(data[idxWVehicleChar : idxWVehicleChar+2])))

	// K-Constant (2 bytes)
	rec.SetKConstantOfRecordingEquipment(int32(binary.BigEndian.Uint16(data[idxKConstant : idxKConstant+2])))

	// L-Tyre Circumference (2 bytes)
	rec.SetLTyreCircumferenceEighthsMm(int32(binary.BigEndian.Uint16(data[idxLTyreCirc : idxLTyreCirc+2])))

	// tyreSize (15 bytes)
	tyreSize, err := opts.UnmarshalIa5StringValue(data[idxTyreSize : idxTyreSize+lenTyreSize])
	if err != nil {
		return nil, fmt.Errorf("tyre size: %w", err)
	}
	rec.SetTyreSize(tyreSize)

	// authorisedSpeed (1 byte)
	rec.SetAuthorisedSpeedKmh(int32(data[idxAuthorisedSpeed]))

	// oldOdometerValue (3 bytes, 24-bit big-endian)
	rec.SetOldOdometerValueKm(int32(data[idxOldOdometer])<<16 |
		int32(data[idxOldOdometer+1])<<8 |
		int32(data[idxOldOdometer+2]))

	// newOdometerValue (3 bytes)
	rec.SetNewOdometerValueKm(int32(data[idxNewOdometer])<<16 |
		int32(data[idxNewOdometer+1])<<8 |
		int32(data[idxNewOdometer+2]))

	// oldTimeValue (4 bytes)
	oldTime, err := opts.UnmarshalTimeReal(data[idxOldTimeValue : idxOldTimeValue+4])
	if err != nil {
		return nil, fmt.Errorf("old time value: %w", err)
	}
	rec.SetOldTimeValue(oldTime)

	// newTimeValue (4 bytes)
	newTime, err := opts.UnmarshalTimeReal(data[idxNewTimeValue : idxNewTimeValue+4])
	if err != nil {
		return nil, fmt.Errorf("new time value: %w", err)
	}
	rec.SetNewTimeValue(newTime)

	// nextCalibrationDate (4 bytes)
	nextCal, err := opts.UnmarshalTimeReal(data[idxNextCalDate : idxNextCalDate+4])
	if err != nil {
		return nil, fmt.Errorf("next calibration date: %w", err)
	}
	rec.SetNextCalibrationDate(nextCal)

	// === V2 extension ===

	// sensorSerialNumber (8 bytes)
	sensorSN, err := opts.UnmarshalExtendedSerialNumber(data[idxSensorSerial : idxSensorSerial+8])
	if err != nil {
		return nil, fmt.Errorf("sensor serial number: %w", err)
	}
	rec.SetSensorSerialNumber(sensorSN)

	// sensorGNSSSerialNumber (8 bytes)
	gnssSN, err := opts.UnmarshalExtendedSerialNumber(data[idxGNSSSerial : idxGNSSSerial+8])
	if err != nil {
		return nil, fmt.Errorf("GNSS serial number: %w", err)
	}
	rec.SetSensorGnssSerialNumber(gnssSN)

	// rcmSerialNumber (8 bytes)
	rcmSN, err := opts.UnmarshalExtendedSerialNumber(data[idxRCMSerial : idxRCMSerial+8])
	if err != nil {
		return nil, fmt.Errorf("RCM serial number: %w", err)
	}
	rec.SetRcmSerialNumber(rcmSN)

	// sealDataVu (55 bytes = 5 × SealRecord(11))
	sealRecords := make([]*vuv1.TechnicalDataGen2V2_SealRecord, numSealRecords)
	for s := range numSealRecords {
		sealStart := idxSealData + s*lenSealRecord
		sealRec := &vuv1.TechnicalDataGen2V2_SealRecord{}

		equipType, eqErr := dd.UnmarshalEnum[ddv1.EquipmentType](data[sealStart])
		if eqErr != nil {
			sealRec.SetUnrecognizedEquipmentType(int32(data[sealStart]))
		} else {
			sealRec.SetEquipmentType(equipType)
		}
		sealRec.SetManufacturerCode(string(bytes.Trim(data[sealStart+1:sealStart+3], "\x00\xff ")))
		sealRec.SetSealIdentifier(string(bytes.Trim(data[sealStart+3:sealStart+11], "\x00\xff ")))
		sealRecords[s] = sealRec
	}
	rec.SetSealRecords(sealRecords)

	// byDefaultLoadType (1 byte)
	rec.SetLoadType(int32(data[idxLoadType]))

	// calibrationCountry (1 byte)
	rec.SetCalibrationCountry(ddv1.NationNumeric(int32(data[idxCalCountry])))

	// calibrationCountryTimestamp (4 bytes)
	calCountryTs, err := opts.UnmarshalTimeReal(data[idxCalCountryTs : idxCalCountryTs+4])
	if err != nil {
		return nil, fmt.Errorf("calibration country timestamp: %w", err)
	}
	rec.SetCalibrationCountryTimestamp(calCountryTs)

	return rec, nil
}

// parseItsConsentRecordArrayGen2V2 parses a VuITSConsentRecordArray.
//
// VuITSConsentRecord (20 bytes):
//
//	fullCardNumberAndGeneration FullCardNumberAndGeneration  -- 19 bytes
//	vuITSConsentGranted bool                                 -- 1 byte
func parseItsConsentRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.TechnicalDataGen2V2_ItsConsentRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	// Gen2V2 devices may send larger ItsConsentRecords (45 bytes observed vs spec 20).
	// Parse the shared 20-byte prefix; extra bytes are preserved in the parent TV raw_data.
	const minRecordSize = 20
	if int(recordSize) < minRecordSize {
		return nil, 0, fmt.Errorf("expected ItsConsentRecord size >= %d, got %d", minRecordSize, recordSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.TechnicalDataGen2V2_ItsConsentRecord, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for ItsConsentRecord %d", i)
		}
		rec := data[recStart : recStart+minRecordSize] // parse first 20 bytes only

		const lenFullCardNumberAndGen = 19
		cardNumber, err := unmarshalOpts.UnmarshalFullCardNumberAndGeneration(rec[:lenFullCardNumberAndGen])
		if err != nil {
			return nil, 0, fmt.Errorf("ItsConsentRecord %d card number: %w", i, err)
		}

		consent := &vuv1.TechnicalDataGen2V2_ItsConsentRecord{}
		consent.SetFullCardNumberAndGeneration(cardNumber)
		consent.SetConsentStatus(rec[lenFullCardNumberAndGen] != 0)
		records = append(records, consent)
		recStart = recEnd
	}

	return records, headerSize + int(recordSize)*int(noOfRecords), nil
}

// parsePowerSupplyInterruptionRecordArrayGen2V2 parses a VuPowerSupplyInterruptionRecordArray.
//
// VuPowerSupplyInterruptionRecord (5 bytes):
//
//	eventTimestamp TimeReal     -- 4 bytes
//	cardSlotNumber CardSlotNumber -- 1 byte
func parsePowerSupplyInterruptionRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.TechnicalDataGen2V2_PowerSupplyInterruptionRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	// Gen2V2 devices may send larger PowerSupplyInterruptionRecords (20 bytes observed vs spec 5).
	// Parse the shared 5-byte prefix; extra bytes are preserved in the parent TV raw_data.
	const minRecordSize = 5
	if int(recordSize) < minRecordSize {
		return nil, 0, fmt.Errorf("expected PowerSupplyInterruptionRecord size >= %d, got %d", minRecordSize, recordSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.TechnicalDataGen2V2_PowerSupplyInterruptionRecord, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for PowerSupplyInterruptionRecord %d", i)
		}
		rec := data[recStart : recStart+minRecordSize] // parse first 5 bytes only

		timestamp, err := unmarshalOpts.UnmarshalTimeReal(rec[:4])
		if err != nil {
			return nil, 0, fmt.Errorf("PowerSupplyInterruptionRecord %d timestamp: %w", i, err)
		}

		psi := &vuv1.TechnicalDataGen2V2_PowerSupplyInterruptionRecord{}
		psi.SetTimestamp(timestamp)
		slotNumber, err := dd.UnmarshalEnum[ddv1.CardSlotNumber](rec[4])
		if err != nil {
			psi.SetUnrecognizedCardSlotNumber(int32(rec[4]))
		} else {
			psi.SetCardSlotNumber(slotNumber)
		}
		records = append(records, psi)
		recStart = recEnd
	}

	return records, headerSize + int(recordSize)*int(noOfRecords), nil
}

// ===== V2-specific marshal helpers =====

// marshalVuIdentificationGen2V2 marshals a V2 VuIdentification to binary (124 bytes).
func marshalVuIdentificationGen2V2(opts dd.MarshalOptions, ident *vuv1.TechnicalDataGen2V2_VuIdentification) ([]byte, error) {
	if ident == nil {
		return make([]byte, 124), nil
	}
	ddIdent := &ddv1.VuIdentification{}
	ddIdent.SetManufacturerName(ident.GetManufacturerName())
	ddIdent.SetManufacturerAddress(ident.GetManufacturerAddress())
	ddIdent.SetPartNumber(ident.GetPartNumber())
	ddIdent.SetSerialNumber(ident.GetSerialNumber())
	ddIdent.SetSoftwareIdentification(ident.GetSoftwareIdentification())
	ddIdent.SetManufacturingDate(ident.GetManufacturingDate())
	ddIdent.SetApprovalNumber(ident.GetApprovalNumber())
	return opts.MarshalVuIdentification(ddIdent)
}

// marshalSensorPairedRecordsGen2V2 marshals V2 paired sensor records to binary.
func marshalSensorPairedRecordsGen2V2(opts dd.MarshalOptions, sensors []*vuv1.TechnicalDataGen2V2_PairedSensor) ([]byte, error) {
	result := make([]byte, 0, len(sensors)*28)
	for i, sensor := range sensors {
		ddSensor := &ddv1.SensorPaired{}
		ddSensor.SetSerialNumber(sensor.GetSerialNumber())
		ddSensor.SetApprovalNumber(sensor.GetApprovalNumber())
		ddSensor.SetPairingDate(sensor.GetPairingDate())
		b, err := opts.MarshalSensorPaired(ddSensor)
		if err != nil {
			return nil, fmt.Errorf("sensor %d: %w", i, err)
		}
		result = append(result, b...)
	}
	return result, nil
}

// marshalCoupledGnssRecordsGen2V2 marshals V2 coupled GNSS records to binary.
// Uses SensorPaired binary layout (same 28-byte structure).
func marshalCoupledGnssRecordsGen2V2(opts dd.MarshalOptions, gnssRecords []*vuv1.TechnicalDataGen2V2_CoupledGnss) ([]byte, error) {
	result := make([]byte, 0, len(gnssRecords)*28)
	for i, gnss := range gnssRecords {
		ddSensor := &ddv1.SensorPaired{}
		ddSensor.SetSerialNumber(gnss.GetSerialNumber())
		ddSensor.SetApprovalNumber(gnss.GetApprovalNumber())
		ddSensor.SetPairingDate(gnss.GetCouplingDate())
		b, err := opts.MarshalSensorPaired(ddSensor)
		if err != nil {
			return nil, fmt.Errorf("coupled GNSS record %d: %w", i, err)
		}
		result = append(result, b...)
	}
	return result, nil
}

// marshalCalibrationRecordsGen2V2 marshals V2 calibration records to binary (252 bytes each).
func marshalCalibrationRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.TechnicalDataGen2V2_CalibrationRecord) ([]byte, error) {
	const lenRecord = 252
	result := make([]byte, 0, len(records)*lenRecord)
	for i, rec := range records {
		b, err := marshalOneCalibrationRecordGen2V2(opts, rec)
		if err != nil {
			return nil, fmt.Errorf("calibration record %d: %w", i, err)
		}
		result = append(result, b...)
	}
	return result, nil
}

// marshalOneCalibrationRecordGen2V2 marshals a single Gen2 V2 calibration record to 252 bytes.
func marshalOneCalibrationRecordGen2V2(opts dd.MarshalOptions, rec *vuv1.TechnicalDataGen2V2_CalibrationRecord) ([]byte, error) {
	const lenRecord = 252
	buf := make([]byte, lenRecord)

	// calibrationPurpose (1 byte)
	if rec.GetUnrecognizedPurpose() != 0 {
		buf[0] = byte(rec.GetUnrecognizedPurpose())
	} else {
		buf[0] = byte(rec.GetPurpose())
	}

	// workshopName (36 bytes at offset 1)
	nameBytes, err := opts.MarshalStringValue(rec.GetWorkshopName())
	if err != nil {
		return nil, fmt.Errorf("workshop name: %w", err)
	}
	copy(buf[1:37], nameBytes)

	// workshopAddress (36 bytes at offset 37)
	addrBytes, err := opts.MarshalStringValue(rec.GetWorkshopAddress())
	if err != nil {
		return nil, fmt.Errorf("workshop address: %w", err)
	}
	copy(buf[37:73], addrBytes)

	// workshopCardNumber (18 bytes at offset 73)
	cardBytes, err := opts.MarshalFullCardNumber(rec.GetWorkshopCardNumber())
	if err != nil {
		return nil, fmt.Errorf("workshop card number: %w", err)
	}
	copy(buf[73:91], cardBytes)

	// workshopCardExpiryDate (4 bytes at offset 91)
	expiryBytes, err := opts.MarshalTimeReal(rec.GetWorkshopCardExpiryDate())
	if err != nil {
		return nil, fmt.Errorf("workshop card expiry: %w", err)
	}
	copy(buf[91:95], expiryBytes)

	// VIN (17 bytes at offset 95)
	vinBytes, err := opts.MarshalIa5StringValue(rec.GetVin())
	if err != nil {
		return nil, fmt.Errorf("VIN: %w", err)
	}
	copy(buf[95:112], vinBytes)

	// vehicleRegistration (15 bytes at offset 112)
	vregBytes, err := opts.MarshalVehicleRegistrationIdentification(rec.GetVehicleRegistration())
	if err != nil {
		return nil, fmt.Errorf("vehicle registration: %w", err)
	}
	copy(buf[112:127], vregBytes)

	// W-Vehicle Characteristic Constant (2 bytes at offset 127)
	binary.BigEndian.PutUint16(buf[127:129], uint16(rec.GetWVehicleCharacteristicConstant()))

	// K-Constant (2 bytes at offset 129)
	binary.BigEndian.PutUint16(buf[129:131], uint16(rec.GetKConstantOfRecordingEquipment()))

	// L-Tyre Circumference (2 bytes at offset 131)
	binary.BigEndian.PutUint16(buf[131:133], uint16(rec.GetLTyreCircumferenceEighthsMm()))

	// tyreSize (15 bytes at offset 133)
	tyreSizeBytes, err := opts.MarshalIa5StringValue(rec.GetTyreSize())
	if err != nil {
		return nil, fmt.Errorf("tyre size: %w", err)
	}
	copy(buf[133:148], tyreSizeBytes)

	// authorisedSpeed (1 byte at offset 148)
	buf[148] = byte(rec.GetAuthorisedSpeedKmh())

	// oldOdometerValue (3 bytes at offset 149)
	v := rec.GetOldOdometerValueKm()
	buf[149] = byte(v >> 16)
	buf[150] = byte(v >> 8)
	buf[151] = byte(v)

	// newOdometerValue (3 bytes at offset 152)
	v = rec.GetNewOdometerValueKm()
	buf[152] = byte(v >> 16)
	buf[153] = byte(v >> 8)
	buf[154] = byte(v)

	// oldTimeValue (4 bytes at offset 155)
	oldTimeBytes, err := opts.MarshalTimeReal(rec.GetOldTimeValue())
	if err != nil {
		return nil, fmt.Errorf("old time value: %w", err)
	}
	copy(buf[155:159], oldTimeBytes)

	// newTimeValue (4 bytes at offset 159)
	newTimeBytes, err := opts.MarshalTimeReal(rec.GetNewTimeValue())
	if err != nil {
		return nil, fmt.Errorf("new time value: %w", err)
	}
	copy(buf[159:163], newTimeBytes)

	// nextCalibrationDate (4 bytes at offset 163)
	nextCalBytes, err := opts.MarshalTimeReal(rec.GetNextCalibrationDate())
	if err != nil {
		return nil, fmt.Errorf("next calibration date: %w", err)
	}
	copy(buf[163:167], nextCalBytes)

	// === V2 extension ===

	// sensorSerialNumber (8 bytes at offset 167)
	sensorSNBytes, err := opts.MarshalExtendedSerialNumber(rec.GetSensorSerialNumber())
	if err != nil {
		return nil, fmt.Errorf("sensor serial number: %w", err)
	}
	copy(buf[167:175], sensorSNBytes)

	// sensorGNSSSerialNumber (8 bytes at offset 175)
	gnssSNBytes, err := opts.MarshalExtendedSerialNumber(rec.GetSensorGnssSerialNumber())
	if err != nil {
		return nil, fmt.Errorf("GNSS serial number: %w", err)
	}
	copy(buf[175:183], gnssSNBytes)

	// rcmSerialNumber (8 bytes at offset 183)
	rcmSNBytes, err := opts.MarshalExtendedSerialNumber(rec.GetRcmSerialNumber())
	if err != nil {
		return nil, fmt.Errorf("RCM serial number: %w", err)
	}
	copy(buf[183:191], rcmSNBytes)

	// sealDataVu (55 bytes at offset 191 = 5 × SealRecord(11))
	for s, seal := range rec.GetSealRecords() {
		if s >= 5 {
			break
		}
		off := 191 + s*11
		if seal.GetUnrecognizedEquipmentType() != 0 {
			buf[off] = byte(seal.GetUnrecognizedEquipmentType())
		} else {
			buf[off] = byte(seal.GetEquipmentType())
		}
		copy(buf[off+1:off+3], []byte(seal.GetManufacturerCode()))
		copy(buf[off+3:off+11], []byte(seal.GetSealIdentifier()))
	}

	// byDefaultLoadType (1 byte at offset 246)
	buf[246] = byte(rec.GetLoadType())

	// calibrationCountry (1 byte at offset 247)
	buf[247] = byte(rec.GetCalibrationCountry())

	// calibrationCountryTimestamp (4 bytes at offset 248)
	calCountryTsBytes, err := opts.MarshalTimeReal(rec.GetCalibrationCountryTimestamp())
	if err != nil {
		return nil, fmt.Errorf("calibration country timestamp: %w", err)
	}
	copy(buf[248:252], calCountryTsBytes)

	return buf, nil
}

// marshalItsConsentRecordsGen2V2 marshals V2 ITS consent records to binary.
func marshalItsConsentRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.TechnicalDataGen2V2_ItsConsentRecord) ([]byte, error) {
	result := make([]byte, 0, len(records)*20)
	for i, rec := range records {
		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(rec.GetFullCardNumberAndGeneration())
		if err != nil {
			return nil, fmt.Errorf("ITS consent record %d card number: %w", i, err)
		}
		if len(cardBytes) != 19 {
			return nil, fmt.Errorf("ITS consent record %d: expected 19 bytes for card number, got %d", i, len(cardBytes))
		}
		result = append(result, cardBytes...)
		if rec.GetConsentStatus() {
			result = append(result, 1)
		} else {
			result = append(result, 0)
		}
	}
	return result, nil
}

// marshalPowerSupplyInterruptionRecordsGen2V2 marshals V2 power supply interruption records to binary.
func marshalPowerSupplyInterruptionRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.TechnicalDataGen2V2_PowerSupplyInterruptionRecord) ([]byte, error) {
	result := make([]byte, 0, len(records)*5)
	for i, rec := range records {
		tsBytes, err := opts.MarshalTimeReal(rec.GetTimestamp())
		if err != nil {
			return nil, fmt.Errorf("power supply interruption record %d timestamp: %w", i, err)
		}
		result = append(result, tsBytes...)
		if rec.GetUnrecognizedCardSlotNumber() != 0 {
			result = append(result, byte(rec.GetUnrecognizedCardSlotNumber()))
		} else {
			result = append(result, byte(rec.GetCardSlotNumber()))
		}
	}
	return result, nil
}

// ===== V2-specific conversion helpers =====

// vuIdentToGen2V2 converts a dd VuIdentification to the V2 proto nested type.
func vuIdentToGen2V2(ddIdent *ddv1.VuIdentification) *vuv1.TechnicalDataGen2V2_VuIdentification {
	if ddIdent == nil {
		return nil
	}
	ident := &vuv1.TechnicalDataGen2V2_VuIdentification{}
	ident.SetManufacturerName(ddIdent.GetManufacturerName())
	ident.SetManufacturerAddress(ddIdent.GetManufacturerAddress())
	ident.SetPartNumber(ddIdent.GetPartNumber())
	ident.SetSerialNumber(ddIdent.GetSerialNumber())
	ident.SetSoftwareIdentification(ddIdent.GetSoftwareIdentification())
	ident.SetManufacturingDate(ddIdent.GetManufacturingDate())
	ident.SetApprovalNumber(ddIdent.GetApprovalNumber())
	return ident
}

// sensorPairedToGen2V2 converts a dd SensorPaired to the V2 proto nested type.
func sensorPairedToGen2V2(s *ddv1.SensorPaired) *vuv1.TechnicalDataGen2V2_PairedSensor {
	if s == nil {
		return nil
	}
	sensor := &vuv1.TechnicalDataGen2V2_PairedSensor{}
	sensor.SetSerialNumber(s.GetSerialNumber())
	sensor.SetApprovalNumber(s.GetApprovalNumber())
	sensor.SetPairingDate(s.GetPairingDate())
	return sensor
}

// gen2V1CalibrationToV2 converts a V1-parsed gen2CalibrationRecord to the V2 proto type.
//
// Used when Gen2 V2 VUs send 168-byte (V1 layout) calibration records.
// The V2 extension fields are left as zero values.
func gen2V1CalibrationToV2(r gen2CalibrationRecord) *vuv1.TechnicalDataGen2V2_CalibrationRecord {
	rec := &vuv1.TechnicalDataGen2V2_CalibrationRecord{}
	rec.SetPurpose(r.purpose)
	rec.SetUnrecognizedPurpose(r.unrecognizedPurpose)
	rec.SetWorkshopName(r.workshopName)
	rec.SetWorkshopAddress(r.workshopAddress)
	// V1 has FullCardNumberAndGeneration(19); V2 proto has FullCardNumber(18).
	// Extract just the FullCardNumber from FullCardNumberAndGeneration.
	if fcng := r.workshopCardNumberAndGen; fcng != nil {
		rec.SetWorkshopCardNumber(fcng.GetFullCardNumber())
	}
	// V1 has Datef expiry; V2 proto has Timestamp.
	if d := r.workshopCardExpiryDate; d != nil {
		t := time.Date(int(d.GetYear()), time.Month(d.GetMonth()), int(d.GetDay()), 0, 0, 0, 0, time.UTC)
		rec.SetWorkshopCardExpiryDate(timestamppb.New(t))
	}
	rec.SetVin(r.vin)
	rec.SetVehicleRegistration(r.vehicleRegistration)
	rec.SetWVehicleCharacteristicConstant(r.wVehicleCharConst)
	rec.SetKConstantOfRecordingEquipment(r.kConstantRecordEquip)
	rec.SetLTyreCircumferenceEighthsMm(r.lTyreCircumference)
	rec.SetTyreSize(r.tyreSize)
	rec.SetAuthorisedSpeedKmh(r.authorisedSpeedKmh)
	rec.SetOldOdometerValueKm(r.oldOdometerValueKm)
	rec.SetNewOdometerValueKm(r.newOdometerValueKm)
	rec.SetOldTimeValue(r.oldTimeValue)
	rec.SetNewTimeValue(r.newTimeValue)
	rec.SetNextCalibrationDate(r.nextCalibrationDate)
	return rec
}

// anonymizeExtendedSerialNumber anonymizes an ExtendedSerialNumber by zeroing the serial number.
func anonymizeExtendedSerialNumber(esn *ddv1.ExtendedSerialNumber) *ddv1.ExtendedSerialNumber {
	if esn == nil {
		return nil
	}
	anon := &ddv1.ExtendedSerialNumber{}
	anon.SetType(esn.GetType())
	anon.SetManufacturerCode(esn.GetManufacturerCode())
	anon.SetSerialNumber(0)
	return anon
}
