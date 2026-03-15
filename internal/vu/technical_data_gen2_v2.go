package vu

import (
	"fmt"

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

	// VuCalibrationRecordArray (N records × 168 bytes)
	calRecords := td.GetCalibrationRecords()
	calData, err := marshalCalibrationRecordsGen2V2(marshalOpts, calRecords)
	if err != nil {
		return nil, fmt.Errorf("marshal VuCalibrationRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x04, 168, uint16(len(calRecords)))
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
		anon.SetWorkshopCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(cal.GetWorkshopCardNumberAndGeneration()))
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
func parseCalibrationRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.TechnicalDataGen2V2_CalibrationRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	// Gen2V1 sends 168 bytes; Gen2V2 sends 252 bytes (additional certification/GNSS fields).
	// Parse the shared 168-byte prefix; extra bytes are preserved in the parent TV raw_data.
	const minRecordSize = 168
	if int(recordSize) < minRecordSize {
		return nil, 0, fmt.Errorf("expected Gen2 CalibrationRecord size >= %d, got %d", minRecordSize, recordSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.TechnicalDataGen2V2_CalibrationRecord, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for CalibrationRecord %d", i)
		}
		// Parse first 168 bytes; remaining Gen2V2-specific bytes are skipped semantically.
		parsed, err := parseOneCalibrationRecordGen2(unmarshalOpts, data[recStart:recStart+minRecordSize])
		if err != nil {
			return nil, 0, fmt.Errorf("CalibrationRecord %d: %w", i, err)
		}
		records = append(records, gen2CalibrationToV2(parsed))
		recStart = recEnd
	}

	return records, headerSize + int(recordSize)*int(noOfRecords), nil
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

// marshalCalibrationRecordsGen2V2 marshals V2 calibration records to binary.
func marshalCalibrationRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.TechnicalDataGen2V2_CalibrationRecord) ([]byte, error) {
	result := make([]byte, 0, len(records)*168)
	for i, rec := range records {
		b, err := marshalOneCalibrationRecordGen2(opts, gen2CalibrationFromV2(rec))
		if err != nil {
			return nil, fmt.Errorf("calibration record %d: %w", i, err)
		}
		result = append(result, b...)
	}
	return result, nil
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

// gen2CalibrationToV2 converts a parsed gen2CalibrationRecord to the V2 proto type.
func gen2CalibrationToV2(r gen2CalibrationRecord) *vuv1.TechnicalDataGen2V2_CalibrationRecord {
	rec := &vuv1.TechnicalDataGen2V2_CalibrationRecord{}
	rec.SetPurpose(r.purpose)
	rec.SetUnrecognizedPurpose(r.unrecognizedPurpose)
	rec.SetWorkshopName(r.workshopName)
	rec.SetWorkshopAddress(r.workshopAddress)
	rec.SetWorkshopCardNumberAndGeneration(r.workshopCardNumberAndGen)
	rec.SetWorkshopCardExpiryDate(r.workshopCardExpiryDate)
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

// gen2CalibrationFromV2 converts a V2 CalibrationRecord to the shared gen2CalibrationRecord.
func gen2CalibrationFromV2(rec *vuv1.TechnicalDataGen2V2_CalibrationRecord) gen2CalibrationRecord {
	if rec == nil {
		return gen2CalibrationRecord{}
	}
	return gen2CalibrationRecord{
		purpose:                  rec.GetPurpose(),
		unrecognizedPurpose:      rec.GetUnrecognizedPurpose(),
		workshopName:             rec.GetWorkshopName(),
		workshopAddress:          rec.GetWorkshopAddress(),
		workshopCardNumberAndGen: rec.GetWorkshopCardNumberAndGeneration(),
		workshopCardExpiryDate:   rec.GetWorkshopCardExpiryDate(),
		vin:                      rec.GetVin(),
		vehicleRegistration:      rec.GetVehicleRegistration(),
		wVehicleCharConst:        rec.GetWVehicleCharacteristicConstant(),
		kConstantRecordEquip:     rec.GetKConstantOfRecordingEquipment(),
		lTyreCircumference:       rec.GetLTyreCircumferenceEighthsMm(),
		tyreSize:                 rec.GetTyreSize(),
		authorisedSpeedKmh:       rec.GetAuthorisedSpeedKmh(),
		oldOdometerValueKm:       rec.GetOldOdometerValueKm(),
		newOdometerValueKm:       rec.GetNewOdometerValueKm(),
		oldTimeValue:             rec.GetOldTimeValue(),
		newTimeValue:             rec.GetNewTimeValue(),
		nextCalibrationDate:      rec.GetNextCalibrationDate(),
	}
}
