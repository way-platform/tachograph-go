package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalTechnicalDataGen1 parses Gen1 Technical Data from the complete transfer value.
//
// Gen1 Technical Data structure (from Data Dictionary and Appendix 7, Section 2.2.6.7):
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuApprovalNumber                VuApprovalNumber,
//	    vuSoftwareIdentification        VuSoftwareIdentification,
//	    vuManufacturerName              VuManufacturerName,
//	    vuManufacturerAddress           VuManufacturerAddress,
//	    vuPartNumber                    VuPartNumber,
//	    vuSerialNumber                  ExtendedSerialNumber,
//	    sensorPaired                    SensorPaired,
//	    signature                       SignatureFirstGen
//	}
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalTechnicalDataGen1(value []byte) (*vuv1.TechnicalDataGen1, error) {
	technicalData := &vuv1.TechnicalDataGen1{}
	technicalData.SetRawData(value)

	// TODO: Implement full semantic parsing
	// For now, validate that we have enough data for the structure
	if len(value) < 128 { // At minimum, signature is 128 bytes
		return nil, fmt.Errorf("insufficient data for Technical Data Gen1")
	}

	// Store the signature (last 128 bytes)
	signatureStart := len(value) - 128
	technicalData.SetSignature(value[signatureStart:])

	return technicalData, nil
}

// appendTechnicalDataGen1 marshals Gen1 Technical Data using raw data painting.
func appendTechnicalDataGen1(dst []byte, technicalData *vuv1.TechnicalDataGen1) ([]byte, error) {
	if technicalData == nil {
		return nil, fmt.Errorf("technicalData cannot be nil")
	}

	raw := technicalData.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	return nil, fmt.Errorf("cannot marshal Technical Data Gen1 without raw_data (semantic marshalling not yet implemented)")
}
