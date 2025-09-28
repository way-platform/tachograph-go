package tachograph

import (
	"bytes"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuDetailedSpeed unmarshals VU detailed speed data from a VU transfer.
//
// The data type `VuDetailedSpeed` is specified in the Data Dictionary, Section 2.2.6.4.
//
// ASN.1 Definition:
//
//	VuDetailedSpeedFirstGen ::= SEQUENCE {
//	    vuDetailedSpeedBlock              VuDetailedSpeedBlock,
//	    signature                         SignatureFirstGen
//	}
//
//	VuDetailedSpeedSecondGen ::= SEQUENCE {
//	    vuDetailedSpeedBlockRecordArray   VuDetailedSpeedBlockRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}
func unmarshalVuDetailedSpeed(data []byte, offset int, target *vuv1.DetailedSpeed, generation int) (int, error) {
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
		return 0, fmt.Errorf("failed to read detailed speed data: %w", err)
	}

	// Set as signature based on generation
	if generation == 1 {
		target.SetSignatureGen1(remainingData)
	} else {
		target.SetSignatureGen2(remainingData)
	}

	return offset - startOffset, nil
}

// AppendVuDetailedSpeed appends VU detailed speed data to a buffer.
//
// The data type `VuDetailedSpeed` is specified in the Data Dictionary, Section 2.2.6.4.
//
// ASN.1 Definition:
//
//	VuDetailedSpeedFirstGen ::= SEQUENCE {
//	    vuDetailedSpeedBlock              VuDetailedSpeedBlock,
//	    signature                         SignatureFirstGen
//	}
//
//	VuDetailedSpeedSecondGen ::= SEQUENCE {
//	    vuDetailedSpeedBlockRecordArray   VuDetailedSpeedBlockRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}
func appendVuDetailedSpeed(buf *bytes.Buffer, detailedSpeed *vuv1.DetailedSpeed) error {
	if detailedSpeed == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if detailedSpeed.GetGeneration() == ddv1.Generation_GENERATION_1 {
		signature := detailedSpeed.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := detailedSpeed.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}
