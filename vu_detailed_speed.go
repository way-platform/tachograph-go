package tachograph

import (
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

// appendVuDetailedSpeedBytes appends VU detailed speed data to a byte slice
func appendVuDetailedSpeedBytes(dst []byte, detailedSpeed *vuv1.DetailedSpeed) ([]byte, error) {
	if detailedSpeed == nil {
		return dst, nil
	}

	if detailedSpeed.GetGeneration() == ddv1.Generation_GENERATION_1 {
		return appendVuDetailedSpeedGen1Bytes(dst, detailedSpeed)
	} else {
		return appendVuDetailedSpeedGen2Bytes(dst, detailedSpeed)
	}
}

// appendVuDetailedSpeedGen1Bytes appends Generation 1 VU detailed speed data
func appendVuDetailedSpeedGen1Bytes(dst []byte, detailedSpeed *vuv1.DetailedSpeed) ([]byte, error) {
	// For now, implement a simplified version that just writes signature data
	// This matches the current unmarshal behavior which reads all data as signature
	// This ensures the interface is complete while allowing for future enhancement

	signature := detailedSpeed.GetSignatureGen1()
	if len(signature) > 0 {
		dst = append(dst, signature...)
	}

	return dst, nil
}

// appendVuDetailedSpeedGen2Bytes appends Generation 2 VU detailed speed data
func appendVuDetailedSpeedGen2Bytes(dst []byte, detailedSpeed *vuv1.DetailedSpeed) ([]byte, error) {
	// For now, implement a simplified version that just writes signature data
	// This matches the current unmarshal behavior which reads all data as signature
	// This ensures the interface is complete while allowing for future enhancement

	signature := detailedSpeed.GetSignatureGen2()
	if len(signature) > 0 {
		dst = append(dst, signature...)
	}

	return dst, nil
}
