package vu

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalTechnicalDataGen2V1 parses Gen2 V1 Technical Data from the complete transfer value.
//
// Structure:
//
//	VuTechnicalDataSecondGenV1 ::= SEQUENCE {
//	    vuIdentificationRecordArray  VuIdentificationRecordArray,
//	    vuSensorPairedRecordArray    VuSensorPairedRecordArray,
//	    vuCalibrationRecordArray     VuCalibrationRecordArray,
//	    signatureRecordArray         SignatureRecordArray
//	}
func unmarshalTechnicalDataGen2V1(value []byte) (*vuv1.TechnicalDataGen2V1, error) {
	totalSize, signatureSize, err := sizeOfTechnicalDataGen2V1(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	td := &vuv1.TechnicalDataGen2V1{}
	td.SetRawData(value)
	offset := 0

	// VuIdentificationRecordArray
	ddIdent, bytesRead, err := parseVuIdentificationRecordArrayGen2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuIdentificationRecordArray: %w", err)
	}
	td.SetVuIdentification(vuIdentToGen2V1(ddIdent))
	offset += bytesRead

	// SensorPairedRecordArray
	ddSensors, bytesRead, err := parseSensorPairedRecordArrayGen2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse SensorPairedRecordArray: %w", err)
	}
	pairedSensors := make([]*vuv1.TechnicalDataGen2V1_PairedSensor, len(ddSensors))
	for i, s := range ddSensors {
		pairedSensors[i] = sensorPairedToGen2V1(s)
	}
	td.SetPairedSensors(pairedSensors)
	offset += bytesRead

	// VuCalibrationRecordArray
	calRecords, bytesRead, err := parseCalibrationRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuCalibrationRecordArray: %w", err)
	}
	td.SetCalibrationRecords(calRecords)
	offset += bytesRead

	td.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Technical Data Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return td, nil
}

// MarshalTechnicalDataGen2V1 marshals Gen2 V1 Technical Data.
func (opts MarshalOptions) MarshalTechnicalDataGen2V1(td *vuv1.TechnicalDataGen2V1) ([]byte, error) {
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
	vuIdentData, err := marshalVuIdentificationGen2(marshalOpts, td.GetVuIdentification())
	if err != nil {
		return nil, fmt.Errorf("marshal VuIdentificationRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x01, 124, 1)
	result = append(result, vuIdentData...)

	// SensorPairedRecordArray (N records × 28 bytes)
	sensors := td.GetPairedSensors()
	sensorData, err := marshalSensorPairedRecordsGen2V1(marshalOpts, sensors)
	if err != nil {
		return nil, fmt.Errorf("marshal SensorPairedRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x02, 28, uint16(len(sensors)))
	result = append(result, sensorData...)

	// VuCalibrationRecordArray (N records × 168 bytes)
	calRecords := td.GetCalibrationRecords()
	calData, err := marshalCalibrationRecordsGen2V1(marshalOpts, calRecords)
	if err != nil {
		return nil, fmt.Errorf("marshal VuCalibrationRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, 168, uint16(len(calRecords)))
	result = append(result, calData...)

	// Signature: stored as complete SignatureRecordArray bytes (header + sig bytes).
	// When empty (anonymized data), include a placeholder header so sizeOf can parse the output.
	if sig := td.GetSignature(); len(sig) > 0 {
		result = append(result, sig...)
	} else {
		result = appendRecordArrayHeader(result, 0x04, 0, 0)
	}
	return result, nil
}

// anonymizeTechnicalDataGen2V1 anonymizes Gen2 V1 Technical Data.
func (opts AnonymizeOptions) anonymizeTechnicalDataGen2V1(td *vuv1.TechnicalDataGen2V1) *vuv1.TechnicalDataGen2V1 {
	if td == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &vuv1.TechnicalDataGen2V1{}

	// Anonymize VU identification
	if vuIdent := td.GetVuIdentification(); vuIdent != nil {
		anon := &vuv1.TechnicalDataGen2V1_VuIdentification{}
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

	// Anonymize paired sensors (zero serial numbers, anonymize approval numbers)
	anonSensors := make([]*vuv1.TechnicalDataGen2V1_PairedSensor, len(td.GetPairedSensors()))
	for i, sensor := range td.GetPairedSensors() {
		anon := &vuv1.TechnicalDataGen2V1_PairedSensor{}
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

	// Anonymize calibration records
	anonCals := make([]*vuv1.TechnicalDataGen2V1_CalibrationRecord, len(td.GetCalibrationRecords()))
	for i, cal := range td.GetCalibrationRecords() {
		anon := &vuv1.TechnicalDataGen2V1_CalibrationRecord{}
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

	result.SetSignature([]byte{})
	return result
}

// ===== Shared Gen2 parse helpers =====

// parseVuIdentificationRecordArrayGen2 parses a VuIdentificationRecordArray.
//
// Gen2V1 sends 124 bytes; Gen2V2 sends 138 bytes (additional certification fields).
// UnmarshalVuIdentification handles variable-length records; extra bytes are preserved in raw_data.
func parseVuIdentificationRecordArrayGen2(data []byte, offset int) (*ddv1.VuIdentification, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const minRecordSize = 124 // Gen2V1; Gen2V2 is 138
	if int(recordSize) < minRecordSize {
		return nil, 0, fmt.Errorf("expected VuIdentification record size >= %d, got %d", minRecordSize, recordSize)
	}
	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 VuIdentification record, got %d", noOfRecords)
	}

	recStart := offset + headerSize
	recEnd := recStart + int(recordSize)
	if recEnd > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for VuIdentification record")
	}

	var unmarshalOpts dd.UnmarshalOptions
	ident, err := unmarshalOpts.UnmarshalVuIdentification(data[recStart:recEnd])
	if err != nil {
		return nil, 0, fmt.Errorf("unmarshal VuIdentification: %w", err)
	}

	return ident, headerSize + int(recordSize), nil
}

// parseSensorPairedRecordArrayGen2 parses a SensorPairedRecordArray.
//
// Expects N records × 28 bytes (Gen2).
func parseSensorPairedRecordArrayGen2(data []byte, offset int) ([]*ddv1.SensorPaired, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 28
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected SensorPaired record size %d, got %d", expectedRecordSize, recordSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	sensors := make([]*ddv1.SensorPaired, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for SensorPaired record %d", i)
		}
		sensor, err := unmarshalOpts.UnmarshalSensorPaired(data[recStart:recEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("SensorPaired record %d: %w", i, err)
		}
		sensors = append(sensors, sensor)
		recStart = recEnd
	}

	return sensors, headerSize + int(recordSize)*int(noOfRecords), nil
}

// parseCalibrationRecordArrayGen2V1 parses a VuCalibrationRecordArray for Gen2 V1.
func parseCalibrationRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.TechnicalDataGen2V1_CalibrationRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 168
	if recordSize != expectedRecordSize {
		// Some Gen2 V1 VU files produced by certain firmware versions have
		// non-standard calibration record sizes (e.g. 28 bytes), likely a
		// GNSS-coupled layout or a national implementation deviation. Skip
		// the records rather than failing the entire parse and advance the
		// cursor by the declared record count × record size so the array
		// framing stays honoured. This is a workaround: a proper fix would
		// decode the alternate layout and emit best-effort records.
		totalConsumed := headerSize + int(recordSize)*int(noOfRecords)
		return nil, totalConsumed, nil
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.TechnicalDataGen2V1_CalibrationRecord, 0, noOfRecords)
	recStart := offset + headerSize

	for i := range noOfRecords {
		recEnd := recStart + int(recordSize)
		if recEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for CalibrationRecord %d", i)
		}
		parsed, err := parseOneCalibrationRecordGen2(unmarshalOpts, data[recStart:recEnd])
		if err != nil {
			return nil, 0, fmt.Errorf("CalibrationRecord %d: %w", i, err)
		}
		records = append(records, gen2CalibrationToV1(parsed))
		recStart = recEnd
	}

	return records, headerSize + int(recordSize)*int(noOfRecords), nil
}

// ===== Gen2 calibration record (shared by V1 and V2) =====

// gen2CalibrationRecord holds the parsed fields of a Gen2 VuCalibrationRecord.
type gen2CalibrationRecord struct {
	purpose                  ddv1.CalibrationPurpose
	unrecognizedPurpose      int32
	workshopName             *ddv1.StringValue
	workshopAddress          *ddv1.StringValue
	workshopCardNumberAndGen *ddv1.FullCardNumberAndGeneration
	workshopCardExpiryDate   *ddv1.Date
	vin                      *ddv1.Ia5StringValue
	vehicleRegistration      *ddv1.VehicleRegistrationIdentification
	wVehicleCharConst        int32
	kConstantRecordEquip     int32
	lTyreCircumference       int32
	tyreSize                 *ddv1.Ia5StringValue
	authorisedSpeedKmh       int32
	oldOdometerValueKm       int32
	newOdometerValueKm       int32
	oldTimeValue             *timestamppb.Timestamp
	newTimeValue             *timestamppb.Timestamp
	nextCalibrationDate      *timestamppb.Timestamp
}

// parseOneCalibrationRecordGen2 parses a single 168-byte Gen2 VuCalibrationRecord.
//
// Gen2 layout (168 bytes) — identical to Gen1 except workshopCardNumber
// is FullCardNumberAndGeneration (19 bytes) instead of FullCardNumber (18 bytes):
//
//	calibrationPurpose CalibrationPurpose                          -- 1 byte (offset 0)
//	workshopName Name                                              -- 36 bytes (offset 1)
//	workshopAddress Address                                        -- 36 bytes (offset 37)
//	workshopCardNumberAndGeneration FullCardNumberAndGeneration    -- 19 bytes (offset 73)
//	workshopCardExpiryDate Datef                                   -- 4 bytes (offset 92)
//	vehicleIdentificationNumber VehicleIdentificationNumber        -- 17 bytes (offset 96)
//	vehicleRegistrationIdentification VehicleRegistrationIdentification -- 15 bytes (offset 113)
//	wVehicleCharacteristicConstant W-VehicleCharacteristicConstant -- 2 bytes (offset 128)
//	kConstantOfRecordingEquipment K-ConstantOfRecordingEquipment   -- 2 bytes (offset 130)
//	lTyreCircumference L-TyreCircumference                         -- 2 bytes (offset 132)
//	tyreSize TyreSize                                              -- 15 bytes (offset 134)
//	authorisedSpeed SpeedAuthorised                                -- 1 byte (offset 149)
//	oldOdometerValue OdometerShort                                 -- 3 bytes (offset 150)
//	newOdometerValue OdometerShort                                 -- 3 bytes (offset 153)
//	oldTimeValue TimeReal                                          -- 4 bytes (offset 156)
//	newTimeValue TimeReal                                          -- 4 bytes (offset 160)
//	nextCalibrationDate TimeReal                                   -- 4 bytes (offset 164)
func parseOneCalibrationRecordGen2(opts dd.UnmarshalOptions, data []byte) (gen2CalibrationRecord, error) {
	const lenRecord = 168
	if len(data) != lenRecord {
		return gen2CalibrationRecord{}, fmt.Errorf("invalid Gen2 CalibrationRecord size: got %d, want %d", len(data), lenRecord)
	}

	const (
		idxCalibrationPurpose = 0
		idxWorkshopName       = 1
		lenWorkshopName       = 36
		idxWorkshopAddress    = 37
		lenWorkshopAddress    = 36
		idxWorkshopCardAndGen = 73
		lenWorkshopCardAndGen = 19
		idxWorkshopCardExpiry = 92
		lenWorkshopCardExpiry = 4
		idxVIN                = 96
		lenVIN                = 17
		idxVehicleReg         = 113
		lenVehicleReg         = 15
		idxWVehicleChar       = 128
		idxKConstant          = 130
		idxLTyreCirc          = 132
		idxTyreSize           = 134
		lenTyreSize           = 15
		idxAuthorisedSpeed    = 149
		idxOldOdometer        = 150
		idxNewOdometer        = 153
		idxOldTimeValue       = 156
		lenTimeValue          = 4
		idxNewTimeValue       = 160
		idxNextCalDate        = 164
	)

	var r gen2CalibrationRecord

	// calibrationPurpose (1 byte)
	purposeValue := int32(data[idxCalibrationPurpose])
	purpose := ddv1.CalibrationPurpose(purposeValue)
	valueDesc := purpose.Descriptor().Values().ByNumber(protoreflect.EnumNumber(purposeValue))
	if valueDesc == nil {
		r.unrecognizedPurpose = purposeValue
		r.purpose = ddv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED
	} else {
		r.purpose = purpose
	}

	var err error

	// workshopName (36 bytes, StringValue)
	r.workshopName, err = opts.UnmarshalStringValue(data[idxWorkshopName : idxWorkshopName+lenWorkshopName])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("workshop name: %w", err)
	}

	// workshopAddress (36 bytes, StringValue)
	r.workshopAddress, err = opts.UnmarshalStringValue(data[idxWorkshopAddress : idxWorkshopAddress+lenWorkshopAddress])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("workshop address: %w", err)
	}

	// workshopCardNumberAndGeneration (19 bytes)
	r.workshopCardNumberAndGen, err = opts.UnmarshalFullCardNumberAndGeneration(data[idxWorkshopCardAndGen : idxWorkshopCardAndGen+lenWorkshopCardAndGen])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("workshop card number and generation: %w", err)
	}

	// workshopCardExpiryDate (4 bytes, Datef)
	r.workshopCardExpiryDate, err = opts.UnmarshalDate(data[idxWorkshopCardExpiry : idxWorkshopCardExpiry+lenWorkshopCardExpiry])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("workshop card expiry date: %w", err)
	}

	// VIN (17 bytes, IA5String)
	r.vin, err = opts.UnmarshalIa5StringValue(data[idxVIN : idxVIN+lenVIN])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("VIN: %w", err)
	}

	// vehicleRegistration (15 bytes)
	r.vehicleRegistration, err = opts.UnmarshalVehicleRegistrationIdentification(data[idxVehicleReg : idxVehicleReg+lenVehicleReg])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("vehicle registration: %w", err)
	}

	// W-Vehicle Characteristic Constant (2 bytes, big-endian)
	r.wVehicleCharConst = int32(binary.BigEndian.Uint16(data[idxWVehicleChar : idxWVehicleChar+2]))

	// K-Constant of Recording Equipment (2 bytes, big-endian)
	r.kConstantRecordEquip = int32(binary.BigEndian.Uint16(data[idxKConstant : idxKConstant+2]))

	// L-Tyre Circumference (2 bytes, big-endian)
	r.lTyreCircumference = int32(binary.BigEndian.Uint16(data[idxLTyreCirc : idxLTyreCirc+2]))

	// tyreSize (15 bytes, IA5String)
	r.tyreSize, err = opts.UnmarshalIa5StringValue(data[idxTyreSize : idxTyreSize+lenTyreSize])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("tyre size: %w", err)
	}

	// authorisedSpeed (1 byte)
	r.authorisedSpeedKmh = int32(data[idxAuthorisedSpeed])

	// oldOdometerValue (3 bytes, 24-bit big-endian)
	r.oldOdometerValueKm = int32(data[idxOldOdometer])<<16 |
		int32(data[idxOldOdometer+1])<<8 |
		int32(data[idxOldOdometer+2])

	// newOdometerValue (3 bytes, 24-bit big-endian)
	r.newOdometerValueKm = int32(data[idxNewOdometer])<<16 |
		int32(data[idxNewOdometer+1])<<8 |
		int32(data[idxNewOdometer+2])

	// oldTimeValue (4 bytes, TimeReal)
	r.oldTimeValue, err = opts.UnmarshalTimeReal(data[idxOldTimeValue : idxOldTimeValue+lenTimeValue])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("old time value: %w", err)
	}

	// newTimeValue (4 bytes, TimeReal)
	r.newTimeValue, err = opts.UnmarshalTimeReal(data[idxNewTimeValue : idxNewTimeValue+lenTimeValue])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("new time value: %w", err)
	}

	// nextCalibrationDate (4 bytes, TimeReal)
	r.nextCalibrationDate, err = opts.UnmarshalTimeReal(data[idxNextCalDate : idxNextCalDate+lenTimeValue])
	if err != nil {
		return gen2CalibrationRecord{}, fmt.Errorf("next calibration date: %w", err)
	}

	return r, nil
}

// marshalOneCalibrationRecordGen2 marshals a gen2CalibrationRecord to 168 bytes.
func marshalOneCalibrationRecordGen2(opts dd.MarshalOptions, r gen2CalibrationRecord) ([]byte, error) {
	const size = 168
	var canvas [size]byte

	const (
		idxCalibrationPurpose = 0
		idxWorkshopName       = 1
		idxWorkshopAddress    = 37
		idxWorkshopCardAndGen = 73
		idxWorkshopCardExpiry = 92
		idxVIN                = 96
		idxVehicleReg         = 113
		idxWVehicleChar       = 128
		idxKConstant          = 130
		idxLTyreCirc          = 132
		idxTyreSize           = 134
		idxAuthorisedSpeed    = 149
		idxOldOdometer        = 150
		idxNewOdometer        = 153
		idxOldTimeValue       = 156
		idxNewTimeValue       = 160
		idxNextCalDate        = 164
	)

	// calibrationPurpose (1 byte)
	if r.unrecognizedPurpose != 0 {
		canvas[idxCalibrationPurpose] = byte(r.unrecognizedPurpose)
	} else {
		canvas[idxCalibrationPurpose] = byte(r.purpose)
	}

	// workshopName (36 bytes)
	workshopNameBytes, err := opts.MarshalStringValue(r.workshopName)
	if err != nil {
		return nil, fmt.Errorf("workshop name: %w", err)
	}
	if len(workshopNameBytes) != 36 {
		return nil, fmt.Errorf("workshop name: expected 36 bytes, got %d", len(workshopNameBytes))
	}
	copy(canvas[idxWorkshopName:idxWorkshopName+36], workshopNameBytes)

	// workshopAddress (36 bytes)
	workshopAddressBytes, err := opts.MarshalStringValue(r.workshopAddress)
	if err != nil {
		return nil, fmt.Errorf("workshop address: %w", err)
	}
	if len(workshopAddressBytes) != 36 {
		return nil, fmt.Errorf("workshop address: expected 36 bytes, got %d", len(workshopAddressBytes))
	}
	copy(canvas[idxWorkshopAddress:idxWorkshopAddress+36], workshopAddressBytes)

	// workshopCardNumberAndGeneration (19 bytes)
	workshopCardBytes, err := opts.MarshalFullCardNumberAndGeneration(r.workshopCardNumberAndGen)
	if err != nil {
		return nil, fmt.Errorf("workshop card number and generation: %w", err)
	}
	if len(workshopCardBytes) != 19 {
		return nil, fmt.Errorf("workshop card number and generation: expected 19 bytes, got %d", len(workshopCardBytes))
	}
	copy(canvas[idxWorkshopCardAndGen:idxWorkshopCardAndGen+19], workshopCardBytes)

	// workshopCardExpiryDate (4 bytes)
	expiryDateBytes, err := opts.MarshalDate(r.workshopCardExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("workshop card expiry date: %w", err)
	}
	if len(expiryDateBytes) != 4 {
		return nil, fmt.Errorf("workshop card expiry date: expected 4 bytes, got %d", len(expiryDateBytes))
	}
	copy(canvas[idxWorkshopCardExpiry:idxWorkshopCardExpiry+4], expiryDateBytes)

	// VIN (17 bytes)
	vinBytes, err := opts.MarshalIa5StringValue(r.vin)
	if err != nil {
		return nil, fmt.Errorf("VIN: %w", err)
	}
	if len(vinBytes) != 17 {
		return nil, fmt.Errorf("VIN: expected 17 bytes, got %d", len(vinBytes))
	}
	copy(canvas[idxVIN:idxVIN+17], vinBytes)

	// vehicleRegistration (15 bytes)
	vehicleRegBytes, err := opts.MarshalVehicleRegistrationIdentification(r.vehicleRegistration)
	if err != nil {
		return nil, fmt.Errorf("vehicle registration: %w", err)
	}
	if len(vehicleRegBytes) != 15 {
		return nil, fmt.Errorf("vehicle registration: expected 15 bytes, got %d", len(vehicleRegBytes))
	}
	copy(canvas[idxVehicleReg:idxVehicleReg+15], vehicleRegBytes)

	// W-Vehicle Characteristic Constant (2 bytes, big-endian)
	binary.BigEndian.PutUint16(canvas[idxWVehicleChar:idxWVehicleChar+2], uint16(r.wVehicleCharConst))

	// K-Constant of Recording Equipment (2 bytes, big-endian)
	binary.BigEndian.PutUint16(canvas[idxKConstant:idxKConstant+2], uint16(r.kConstantRecordEquip))

	// L-Tyre Circumference (2 bytes, big-endian)
	binary.BigEndian.PutUint16(canvas[idxLTyreCirc:idxLTyreCirc+2], uint16(r.lTyreCircumference))

	// tyreSize (15 bytes)
	tyreSizeBytes, err := opts.MarshalIa5StringValue(r.tyreSize)
	if err != nil {
		return nil, fmt.Errorf("tyre size: %w", err)
	}
	if len(tyreSizeBytes) != 15 {
		return nil, fmt.Errorf("tyre size: expected 15 bytes, got %d", len(tyreSizeBytes))
	}
	copy(canvas[idxTyreSize:idxTyreSize+15], tyreSizeBytes)

	// authorisedSpeed (1 byte)
	canvas[idxAuthorisedSpeed] = byte(r.authorisedSpeedKmh)

	// oldOdometerValue (3 bytes, 24-bit big-endian)
	canvas[idxOldOdometer] = byte((r.oldOdometerValueKm >> 16) & 0xFF)
	canvas[idxOldOdometer+1] = byte((r.oldOdometerValueKm >> 8) & 0xFF)
	canvas[idxOldOdometer+2] = byte(r.oldOdometerValueKm & 0xFF)

	// newOdometerValue (3 bytes, 24-bit big-endian)
	canvas[idxNewOdometer] = byte((r.newOdometerValueKm >> 16) & 0xFF)
	canvas[idxNewOdometer+1] = byte((r.newOdometerValueKm >> 8) & 0xFF)
	canvas[idxNewOdometer+2] = byte(r.newOdometerValueKm & 0xFF)

	// oldTimeValue (4 bytes)
	oldTimeBytes, err := opts.MarshalTimeReal(r.oldTimeValue)
	if err != nil {
		return nil, fmt.Errorf("old time value: %w", err)
	}
	copy(canvas[idxOldTimeValue:idxOldTimeValue+4], oldTimeBytes)

	// newTimeValue (4 bytes)
	newTimeBytes, err := opts.MarshalTimeReal(r.newTimeValue)
	if err != nil {
		return nil, fmt.Errorf("new time value: %w", err)
	}
	copy(canvas[idxNewTimeValue:idxNewTimeValue+4], newTimeBytes)

	// nextCalibrationDate (4 bytes)
	nextCalBytes, err := opts.MarshalTimeReal(r.nextCalibrationDate)
	if err != nil {
		return nil, fmt.Errorf("next calibration date: %w", err)
	}
	copy(canvas[idxNextCalDate:idxNextCalDate+4], nextCalBytes)

	return canvas[:], nil
}

// gen2CalibrationToV1 converts a parsed gen2CalibrationRecord to the V1 proto type.
func gen2CalibrationToV1(r gen2CalibrationRecord) *vuv1.TechnicalDataGen2V1_CalibrationRecord {
	rec := &vuv1.TechnicalDataGen2V1_CalibrationRecord{}
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

// gen2CalibrationFromV1 converts a V1 CalibrationRecord to the shared gen2CalibrationRecord.
func gen2CalibrationFromV1(rec *vuv1.TechnicalDataGen2V1_CalibrationRecord) gen2CalibrationRecord {
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

// ===== V1-specific marshal helpers =====

// marshalVuIdentificationGen2 marshals a V1 VuIdentification to binary (124 bytes).
func marshalVuIdentificationGen2(opts dd.MarshalOptions, ident *vuv1.TechnicalDataGen2V1_VuIdentification) ([]byte, error) {
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

// marshalSensorPairedRecordsGen2V1 marshals V1 paired sensor records to binary.
func marshalSensorPairedRecordsGen2V1(opts dd.MarshalOptions, sensors []*vuv1.TechnicalDataGen2V1_PairedSensor) ([]byte, error) {
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

// marshalCalibrationRecordsGen2V1 marshals V1 calibration records to binary.
func marshalCalibrationRecordsGen2V1(opts dd.MarshalOptions, records []*vuv1.TechnicalDataGen2V1_CalibrationRecord) ([]byte, error) {
	result := make([]byte, 0, len(records)*168)
	for i, rec := range records {
		b, err := marshalOneCalibrationRecordGen2(opts, gen2CalibrationFromV1(rec))
		if err != nil {
			return nil, fmt.Errorf("calibration record %d: %w", i, err)
		}
		result = append(result, b...)
	}
	return result, nil
}

// ===== V1-specific conversion helpers =====

// vuIdentToGen2V1 converts a dd VuIdentification to the V1 proto nested type.
func vuIdentToGen2V1(ddIdent *ddv1.VuIdentification) *vuv1.TechnicalDataGen2V1_VuIdentification {
	if ddIdent == nil {
		return nil
	}
	ident := &vuv1.TechnicalDataGen2V1_VuIdentification{}
	ident.SetManufacturerName(ddIdent.GetManufacturerName())
	ident.SetManufacturerAddress(ddIdent.GetManufacturerAddress())
	ident.SetPartNumber(ddIdent.GetPartNumber())
	ident.SetSerialNumber(ddIdent.GetSerialNumber())
	ident.SetSoftwareIdentification(ddIdent.GetSoftwareIdentification())
	ident.SetManufacturingDate(ddIdent.GetManufacturingDate())
	ident.SetApprovalNumber(ddIdent.GetApprovalNumber())
	return ident
}

// sensorPairedToGen2V1 converts a dd SensorPaired to the V1 proto nested type.
func sensorPairedToGen2V1(s *ddv1.SensorPaired) *vuv1.TechnicalDataGen2V1_PairedSensor {
	if s == nil {
		return nil
	}
	sensor := &vuv1.TechnicalDataGen2V1_PairedSensor{}
	sensor.SetSerialNumber(s.GetSerialNumber())
	sensor.SetApprovalNumber(s.GetApprovalNumber())
	sensor.SetPairingDate(s.GetPairingDate())
	return sensor
}
