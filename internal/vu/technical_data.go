package vu

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuTechnicalData unmarshals VU technical data from a VU transfer.
//
// The data type `VuTechnicalData` is specified in the Data Dictionary, Section 2.2.6.5.
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuIdentification                  VuIdentification,
//	    vuCalibrationData                 VuCalibrationData,
//	    vuCardData                        VuCardData,
//	    signature                         SignatureFirstGen
//	}
//
//	VuTechnicalDataSecondGen ::= SEQUENCE {
//	    vuIdentificationRecordArray       VuIdentificationRecordArray,
//	    vuCalibrationRecordArray          VuCalibrationRecordArray,
//	    vuCardRecordArray                 VuCardRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}
func UnmarshalVuTechnicalData(data []byte, offset int, target *vuv1.TechnicalData, generation int) (int, error) {
	startOffset := offset

	// Set generation
	if generation == 1 {
		target.SetGeneration(ddv1.Generation_GENERATION_1)
	} else {
		target.SetGeneration(ddv1.Generation_GENERATION_2)
	}

	// For now, implement a simplified version that just reads the data
	// This ensures the interface is complete while allowing for future enhancement

	// Read all remaining data
	remainingData, offset, err := readBytesFromBytes(data, offset, len(data)-offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read technical data: %w", err)
	}

	// Set as signature based on generation
	if generation == 1 {
		target.SetSignatureGen1(remainingData)
	} else {
		target.SetSignatureGen2(remainingData)
	}

	return offset - startOffset, nil
}

// AppendVuTechnicalData appends VU technical data to a buffer.
//
// The data type `VuTechnicalData` is specified in the Data Dictionary, Section 2.2.6.5.
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuIdentification                  VuIdentification,
//	    vuCalibrationData                 VuCalibrationData,
//	    vuCardData                        VuCardData,
//	    signature                         SignatureFirstGen
//	}
//
//	VuTechnicalDataSecondGen ::= SEQUENCE {
//	    vuIdentificationRecordArray       VuIdentificationRecordArray,
//	    vuCalibrationRecordArray          VuCalibrationRecordArray,
//	    vuCardRecordArray                 VuCardRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}

// appendVuTechnicalDataBytes appends VU technical data to a byte slice
func appendVuTechnicalDataBytes(dst []byte, technicalData *vuv1.TechnicalData) ([]byte, error) {
	if technicalData == nil {
		return dst, nil
	}

	if technicalData.GetGeneration() == ddv1.Generation_GENERATION_1 {
		return appendVuTechnicalDataGen1Bytes(dst, technicalData)
	} else {
		return appendVuTechnicalDataGen2Bytes(dst, technicalData)
	}
}

// appendVuTechnicalDataGen1Bytes appends Generation 1 VU technical data
func appendVuTechnicalDataGen1Bytes(dst []byte, technicalData *vuv1.TechnicalData) ([]byte, error) {
	// For now, implement a simplified version that just writes signature data
	// This matches the current unmarshal behavior which reads all data as signature
	// This ensures the interface is complete while allowing for future enhancement

	signature := technicalData.GetSignatureGen1()
	if len(signature) > 0 {
		dst = append(dst, signature...)
	}

	return dst, nil
}

// appendVuTechnicalDataGen2Bytes appends Generation 2 VU technical data
func appendVuTechnicalDataGen2Bytes(dst []byte, technicalData *vuv1.TechnicalData) ([]byte, error) {
	// For now, implement a simplified version that just writes signature data
	// This matches the current unmarshal behavior which reads all data as signature
	// This ensures the interface is complete while allowing for future enhancement

	signature := technicalData.GetSignatureGen2()
	if len(signature) > 0 {
		dst = append(dst, signature...)
	}

	return dst, nil
}
