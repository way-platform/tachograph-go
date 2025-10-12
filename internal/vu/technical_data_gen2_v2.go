package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalTechnicalDataGen2V2 parses Gen2 V2 Technical Data from the complete transfer value.
//
// Gen2 V2 Technical Data structure is identical to Gen2 V1.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalTechnicalDataGen2V2(value []byte) (*vuv1.TechnicalDataGen2V2, error) {
	technicalData := &vuv1.TechnicalDataGen2V2{}
	technicalData.SetRawData(value)

	// Validate structure by skipping through all record arrays
	offset := 0
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(value, offset)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		offset += size
		return nil
	}

	// Skip all record arrays
	if err := skipRecordArray("VuApprovalNumber"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuSoftwareIdentification"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuManufacturerName"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuManufacturerAddress"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuPartNumber"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuSerialNumber"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("SensorPaired"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("Signature"); err != nil {
		return nil, err
	}

	if offset != len(value) {
		return nil, fmt.Errorf("Technical Data Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	return technicalData, nil
}

// appendTechnicalDataGen2V2 marshals Gen2 V2 Technical Data using raw data painting.
func appendTechnicalDataGen2V2(dst []byte, technicalData *vuv1.TechnicalDataGen2V2) ([]byte, error) {
	if technicalData == nil {
		return nil, fmt.Errorf("technicalData cannot be nil")
	}

	raw := technicalData.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	return nil, fmt.Errorf("cannot marshal Technical Data Gen2 V2 without raw_data (semantic marshalling not yet implemented)")
}
