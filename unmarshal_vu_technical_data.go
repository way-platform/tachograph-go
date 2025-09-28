package tachograph

import (
	"fmt"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuTechnicalData unmarshals VU technical data from a VU transfer.
//
// ASN.1 Specification (Data Dictionary 2.2.6.5):
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
//
// Binary Layout (Gen1):
//
//	Variable size structure containing:
//	- vuIdentification (variable size)
//	- vuCalibrationData (variable size)
//	- vuCardData (variable size)
//	- signature (128 bytes)
//

func UnmarshalVuTechnicalData(data []byte, offset int, target *vuv1.TechnicalData, generation int) (int, error) {
	startOffset := offset

	// Set generation
	if generation == 1 {
		target.SetGeneration(datadictionaryv1.Generation_GENERATION_1)
	} else {
		target.SetGeneration(datadictionaryv1.Generation_GENERATION_2)
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
