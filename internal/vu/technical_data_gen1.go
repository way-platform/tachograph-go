package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalTechnicalDataGen1 parses Gen1 Technical Data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Technical Data structure (from Data Dictionary and Appendix 7, Section 2.2.6.7):
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuIdentification VuIdentification,
//	    sensorPaired SensorPaired,
//	    vuCalibrationData VuCalibrationData,
//	    signature SignatureFirstGen
//	}
func unmarshalTechnicalDataGen1(value []byte) (*vuv1.TechnicalDataGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	technicalData := &vuv1.TechnicalDataGen1{}
	technicalData.SetRawData(value) // Store complete transfer value for painting
	offset := 0
	opts := dd.UnmarshalOptions{PreserveRawData: true}

	// Parse VuIdentification (116 bytes for Gen1: 36+36+16+8+8+4+8)
	const vuIdentificationSize = 116
	if offset+vuIdentificationSize > len(data) {
		return nil, fmt.Errorf("insufficient data for VuIdentification")
	}
	vuIdentification, err := opts.UnmarshalVuIdentification(data[offset : offset+vuIdentificationSize])
	if err != nil {
		return nil, fmt.Errorf("unmarshal VuIdentification: %w", err)
	}
	technicalData.SetVuIdentification(vuIdentification)
	offset += vuIdentificationSize

	// Parse SensorPaired (20 bytes for Gen1)
	const sensorPairedSize = 20
	if offset+sensorPairedSize > len(data) {
		return nil, fmt.Errorf("insufficient data for SensorPaired")
	}
	sensorPaired, err := opts.UnmarshalSensorPaired(data[offset : offset+sensorPairedSize])
	if err != nil {
		return nil, fmt.Errorf("unmarshal SensorPaired: %w", err)
	}
	technicalData.SetPairedSensor(sensorPaired)
	offset += sensorPairedSize

	// Parse VuCalibrationData (1 byte count + calibration records)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfVuCalibrationRecords")
	}
	noOfVuCalibrationRecords := data[offset]
	offset += 1

	// Parse calibration records - only parse as many as will fit in the available data
	// (The count byte may be incorrect or there may be trailing data)
	calibrationRecords := []*ddv1.VuCalibrationRecord{}
	for i := 0; i < int(noOfVuCalibrationRecords); i++ {
		const calibrationRecordSize = 167
		if offset+calibrationRecordSize > len(data) {
			// Not enough data for another complete record - stop parsing
			// This can happen if the count is wrong or there's trailing data
			break
		}
		calibrationRecord, err := opts.UnmarshalVuCalibrationRecord(data[offset : offset+calibrationRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuCalibrationRecord %d: %w", i, err)
		}
		calibrationRecords = append(calibrationRecords, calibrationRecord)
		offset += calibrationRecordSize
	}
	technicalData.SetCalibrationRecords(calibrationRecords)

	// Note: There may be trailing bytes after calibration records (typically 20 bytes or so)
	// This appears to be padding or reserved data, so we don't require exact consumption

	// Store signature (extracted at the beginning)
	technicalData.SetSignature(signature)

	return technicalData, nil
}

// MarshalTechnicalDataGen1 marshals Gen1 Technical Data using raw data painting.
func (opts MarshalOptions) MarshalTechnicalDataGen1(technicalData *vuv1.TechnicalDataGen1) ([]byte, error) {
	if technicalData == nil {
		return nil, fmt.Errorf("technicalData cannot be nil")
	}

	// Calculate data size
	// VuIdentification: 116 bytes (Gen1: 36+36+16+8+8+4+8)
	// SensorPaired: 20 bytes (Gen1)
	// VuCalibrationData: 1 byte count + (n * 167 bytes)
	noOfCalibrationRecords := len(technicalData.GetCalibrationRecords())
	dataSize := 116 + 20 + 1 + (noOfCalibrationRecords * 167)

	// Use raw data painting with canvas
	var canvas []byte
	raw := technicalData.GetRawData()
	if len(raw) == dataSize+128 {
		// raw_data includes signature
		canvas = make([]byte, dataSize)
		copy(canvas, raw[:dataSize])
	} else if len(raw) == dataSize {
		// raw_data is just data portion
		canvas = make([]byte, dataSize)
		copy(canvas, raw)
	} else {
		// No valid raw_data, start with zeros
		canvas = make([]byte, dataSize)
	}

	offset := 0
	marshalOpts := dd.MarshalOptions{}

	// Marshal VuIdentification (116 bytes for Gen1)
	vuIdentBytes, err := marshalOpts.MarshalVuIdentification(technicalData.GetVuIdentification())
	if err != nil {
		return nil, fmt.Errorf("marshal VuIdentification: %w", err)
	}
	if len(vuIdentBytes) != 116 {
		return nil, fmt.Errorf("VuIdentification has invalid length: got %d, want 116", len(vuIdentBytes))
	}
	copy(canvas[offset:offset+116], vuIdentBytes)
	offset += 116

	// Marshal SensorPaired (20 bytes)
	sensorPairedBytes, err := marshalOpts.MarshalSensorPaired(technicalData.GetPairedSensor())
	if err != nil {
		return nil, fmt.Errorf("marshal SensorPaired: %w", err)
	}
	if len(sensorPairedBytes) != 20 {
		return nil, fmt.Errorf("SensorPaired has invalid length: got %d, want 20", len(sensorPairedBytes))
	}
	copy(canvas[offset:offset+20], sensorPairedBytes)
	offset += 20

	// Marshal VuCalibrationData
	canvas[offset] = byte(noOfCalibrationRecords)
	offset += 1

	for i, calibrationRecord := range technicalData.GetCalibrationRecords() {
		calibrationRecordBytes, err := marshalOpts.MarshalVuCalibrationRecord(calibrationRecord)
		if err != nil {
			return nil, fmt.Errorf("marshal VuCalibrationRecord %d: %w", i, err)
		}
		if len(calibrationRecordBytes) != 167 {
			return nil, fmt.Errorf("VuCalibrationRecord %d has invalid length: got %d, want 167", i, len(calibrationRecordBytes))
		}
		copy(canvas[offset:offset+167], calibrationRecordBytes)
		offset += 167
	}

	// Verify we wrote all data
	if offset != dataSize {
		return nil, fmt.Errorf("Technical Data Gen1 marshalling mismatch: wrote %d bytes, expected %d", offset, dataSize)
	}

	// Append signature
	signature := technicalData.GetSignature()
	if len(signature) == 0 {
		signature = make([]byte, 128)
	}
	if len(signature) != 128 {
		return nil, fmt.Errorf("invalid signature length: got %d, want 128", len(signature))
	}
	transferValue := append(canvas, signature...)

	return transferValue, nil
}

// anonymizeTechnicalDataGen1 anonymizes Gen1 Technical Data.
func (opts AnonymizeOptions) anonymizeTechnicalDataGen1(td *vuv1.TechnicalDataGen1) *vuv1.TechnicalDataGen1 {
	if td == nil {
		return nil
	}
	result := proto.Clone(td).(*vuv1.TechnicalDataGen1)

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize VU identification
	if vuIdent := result.GetVuIdentification(); vuIdent != nil {
		result.SetVuIdentification(ddOpts.AnonymizeVuIdentification(vuIdent))
	}

	// Anonymize paired sensor
	if sensor := result.GetPairedSensor(); sensor != nil {
		result.SetPairedSensor(ddOpts.AnonymizeSensorPaired(sensor))
	}

	// Anonymize calibration records
	var anonymizedCalibrationRecords []*ddv1.VuCalibrationRecord
	for _, calibrationRecord := range result.GetCalibrationRecords() {
		anonymizedCalibrationRecords = append(anonymizedCalibrationRecords, ddOpts.AnonymizeVuCalibrationRecord(calibrationRecord))
	}
	result.SetCalibrationRecords(anonymizedCalibrationRecords)

	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Clear raw_data to force semantic marshalling
	result.ClearRawData()

	return result
}
