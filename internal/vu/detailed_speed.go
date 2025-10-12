package vu

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfDetailedSpeed dispatches to generation-specific size calculation.
func sizeOfDetailedSpeed(data []byte, transferType vuv1.TransferType) (int, error) {
	switch transferType {
	case vuv1.TransferType_DETAILED_SPEED_GEN1:
		return sizeOfDetailedSpeedGen1(data)
	case vuv1.TransferType_DETAILED_SPEED_GEN2:
		return sizeOfDetailedSpeedGen2(data)
	default:
		return 0, fmt.Errorf("unsupported transfer type for DetailedSpeed: %v", transferType)
	}
}

// sizeOfDetailedSpeedGen1 calculates total size for Gen1 Detailed Speed including signature.
//
// Detailed Speed Gen1 structure (from Appendix 7, Section 2.2.6.5):
// - VuDetailedSpeedData (Data Dictionary 2.192): 2 bytes + (noOfSpeedBlocks * 64 bytes)
//   - noOfSpeedBlocks: INTEGER(0..2^16-1) = 2 bytes
//   - vuDetailedSpeedBlocks: SET SIZE(noOfSpeedBlocks) OF VuDetailedSpeedBlock
//   - VuDetailedSpeedBlock (Data Dictionary 2.190): 64 bytes total
//   - speedBlockBeginDate: TimeReal = 4 bytes
//   - speedsPerSecond: 60 bytes (60 Speed values, one per second)
//
// - Signature: 128 bytes (RSA)
func sizeOfDetailedSpeedGen1(data []byte) (int, error) {
	offset := 0

	// VuDetailedSpeedData: 2 bytes count + variable speed blocks
	if len(data[offset:]) < 2 {
		return 0, fmt.Errorf("insufficient data for noOfSpeedBlocks")
	}
	noOfSpeedBlocks := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each VuDetailedSpeedBlock: 64 bytes (4 TimeReal + 60 Speed bytes)
	// Per Data Dictionary 2.190
	const vuDetailedSpeedBlockSize = 64
	offset += int(noOfSpeedBlocks) * vuDetailedSpeedBlockSize

	// Signature: 128 bytes for Gen1 RSA
	offset += 128

	return offset, nil
}

// sizeOfDetailedSpeedGen2 calculates size by parsing Gen2 RecordArrays.
func sizeOfDetailedSpeedGen2(data []byte) (int, error) {
	offset := 0

	// VuDetailedSpeedBlockRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuDetailedSpeedBlockRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// ===== Unmarshal Functions =====

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
func UnmarshalVuDetailedSpeed(data []byte, offset int, target *vuv1.DetailedSpeed, generation int) (int, error) {
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
