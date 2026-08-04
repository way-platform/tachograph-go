package vu

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/safemath"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfDetailedSpeed dispatches to generation-specific size calculation.
func sizeOfDetailedSpeed(data []byte, transferType vuv1.TransferType) (totalSize, signatureSize int, err error) {
	switch transferType {
	case vuv1.TransferType_DETAILED_SPEED_GEN1:
		return sizeOfDetailedSpeedGen1(data)
	case vuv1.TransferType_DETAILED_SPEED_GEN2:
		return sizeOfDetailedSpeedGen2(data)
	default:
		return 0, 0, fmt.Errorf("unsupported transfer type for DetailedSpeed: %v", transferType)
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
func sizeOfDetailedSpeedGen1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// VuDetailedSpeedData: 2 bytes count + variable speed blocks
	if len(data[offset:]) < 2 {
		return 0, 0, fmt.Errorf("insufficient data for noOfSpeedBlocks")
	}
	noOfSpeedBlocks := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Each VuDetailedSpeedBlock: 64 bytes (4 TimeReal + 60 Speed bytes)
	// Per Data Dictionary 2.190
	const vuDetailedSpeedBlockSize = 64
	speedBlocksSize, mulErr := safemath.MulInt(int(noOfSpeedBlocks), vuDetailedSpeedBlockSize)
	if mulErr != nil {
		return 0, 0, fmt.Errorf("speed blocks overflow: noOfSpeedBlocks=%d blockSize=%d: %w", noOfSpeedBlocks, vuDetailedSpeedBlockSize, mulErr)
	}
	offset += speedBlocksSize

	// Signature: 128 bytes for Gen1 RSA
	const gen1SignatureSize = 128
	offset += gen1SignatureSize

	return offset, gen1SignatureSize, nil
}

// sizeOfDetailedSpeedGen2 calculates size by parsing Gen2 RecordArrays.
func sizeOfDetailedSpeedGen2(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// VuDetailedSpeedBlockRecordArray
	size, sizeErr := sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuDetailedSpeedBlockRecordArray: %w", sizeErr)
	}
	offset += size

	// SignatureRecordArray (last)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("SignatureRecordArray: %w", sizeErr)
	}
	signatureSizeGen2 := size
	offset += size

	return offset, signatureSizeGen2, nil
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
