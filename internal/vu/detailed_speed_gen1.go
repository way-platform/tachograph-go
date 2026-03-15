package vu

import (
	"encoding/binary"
	"fmt"

	dd "github.com/way-platform/tachograph-go/internal/dd"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalDetailedSpeedGen1 parses Gen1 Detailed Speed data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Detailed Speed structure (from Data Dictionary and Appendix 7, Section 2.2.6.5):
//
// ASN.1 Definition:
//
//	VuDetailedSpeedFirstGen ::= SEQUENCE {
//	    vuDetailedSpeedData        VuDetailedSpeedData,
//	    signature                  SignatureFirstGen
//	}
//
// Binary Layout:
//   - 2 bytes: noOfSpeedBlocks (INTEGER 0..65535)
//   - N x 64 bytes: VuDetailedSpeedBlock records
//   - 128 bytes: RSA-1024 signature
func unmarshalDetailedSpeedGen1(value []byte) (*vuv1.DetailedSpeedGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	detailedSpeed := &vuv1.DetailedSpeedGen1{}
	detailedSpeed.SetRawData(value) // Store complete transfer value for painting

	offset := 0
	opts := dd.UnmarshalOptions{PreserveRawData: true}

	// VuDetailedSpeedData: 2 bytes (noOfSpeedBlocks) + (noOfSpeedBlocks * 64 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfSpeedBlocks")
	}
	noOfSpeedBlocks := int(data[offset])<<8 | int(data[offset+1])
	offset += 2

	// Parse each VuDetailedSpeedBlock (64 bytes each)
	speedBlocks := make([]*vuv1.DetailedSpeedGen1_DetailedSpeedBlock, noOfSpeedBlocks)
	for i := 0; i < noOfSpeedBlocks; i++ {
		const speedBlockSize = 64 // 4 bytes TimeReal + 60 bytes Speed
		if offset+speedBlockSize > len(data) {
			return nil, fmt.Errorf("insufficient data for VuDetailedSpeedBlock %d", i)
		}

		speedBlock, err := unmarshalDetailedSpeedBlock(opts, data[offset:offset+speedBlockSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuDetailedSpeedBlock %d: %w", i, err)
		}

		speedBlocks[i] = speedBlock
		offset += speedBlockSize
	}
	detailedSpeed.SetSpeedBlocks(speedBlocks)

	// Verify we consumed all data
	if offset != len(data) {
		return nil, fmt.Errorf("unexpected extra data: consumed %d bytes, total %d bytes", offset, len(data))
	}

	// Store signature
	detailedSpeed.SetSignature(signature)

	return detailedSpeed, nil
}

// unmarshalDetailedSpeedBlock parses a single VuDetailedSpeedBlock (64 bytes).
//
// Binary Layout:
//   - 4 bytes: speedBlockBeginDate (TimeReal)
//   - 60 bytes: speedsPerSecond (60 x Speed, 1 byte each)
func unmarshalDetailedSpeedBlock(opts dd.UnmarshalOptions, data []byte) (*vuv1.DetailedSpeedGen1_DetailedSpeedBlock, error) {
	const (
		idxBeginDate          = 0
		lenBeginDate          = 4
		idxSpeeds             = 4
		lenSpeeds             = 60
		lenDetailedSpeedBlock = 64
	)

	if len(data) != lenDetailedSpeedBlock {
		return nil, fmt.Errorf("invalid data length for VuDetailedSpeedBlock: got %d, want %d", len(data), lenDetailedSpeedBlock)
	}

	block := &vuv1.DetailedSpeedGen1_DetailedSpeedBlock{}
	if opts.PreserveRawData {
		block.SetRawData(data)
	}

	// Parse begin date (4 bytes)
	beginDate, err := opts.UnmarshalTimeReal(data[idxBeginDate : idxBeginDate+lenBeginDate])
	if err != nil {
		return nil, fmt.Errorf("unmarshal begin date: %w", err)
	}
	block.SetBeginDate(beginDate)

	// Parse 60 speed values (1 byte each, in km/h)
	speeds := make([]int32, lenSpeeds)
	for i := 0; i < lenSpeeds; i++ {
		speeds[i] = int32(data[idxSpeeds+i])
	}
	block.SetSpeedsKmh(speeds)

	return block, nil
}

// MarshalDetailedSpeedGen1 marshals Gen1 Detailed Speed data using raw data painting.
func (opts MarshalOptions) MarshalDetailedSpeedGen1(detailedSpeed *vuv1.DetailedSpeedGen1) ([]byte, error) {
	if detailedSpeed == nil {
		return nil, fmt.Errorf("detailedSpeed cannot be nil")
	}

	// Calculate expected size (signature is stored separately and appended at the end)
	noOfSpeedBlocks := len(detailedSpeed.GetSpeedBlocks())

	// Data portion: 2 (count) + N*64 (blocks)
	dataSize := 2 + (noOfSpeedBlocks * 64)

	// Use raw_data as canvas if available
	var canvas []byte
	raw := detailedSpeed.GetRawData()
	if len(raw) == dataSize+128 { // dataSize + 128-byte signature
		// Extract data portion (without signature) to use as canvas
		canvas = make([]byte, dataSize)
		copy(canvas, raw[:dataSize])
	} else if len(raw) == dataSize {
		// raw_data is just the data portion without signature
		canvas = make([]byte, dataSize)
		copy(canvas, raw)
	} else {
		// Create zero-filled canvas
		canvas = make([]byte, dataSize)
	}

	offset := 0

	// Write noOfSpeedBlocks (2 bytes, big-endian)
	binary.BigEndian.PutUint16(canvas[offset:offset+2], uint16(noOfSpeedBlocks))
	offset += 2

	// Marshal each VuDetailedSpeedBlock (64 bytes each)
	for i, speedBlock := range detailedSpeed.GetSpeedBlocks() {
		blockBytes, err := opts.marshalDetailedSpeedBlock(speedBlock)
		if err != nil {
			return nil, fmt.Errorf("marshal VuDetailedSpeedBlock %d: %w", i, err)
		}
		if len(blockBytes) != 64 {
			return nil, fmt.Errorf("VuDetailedSpeedBlock %d has invalid length: got %d, want 64", i, len(blockBytes))
		}
		copy(canvas[offset:offset+64], blockBytes)
		offset += 64
	}

	// Verify we used exactly the expected amount of data
	if offset != dataSize {
		return nil, fmt.Errorf("Detailed Speed Gen1 marshalling mismatch: wrote %d bytes, expected %d", offset, dataSize)
	}

	// Append signature to create complete transfer value
	signature := detailedSpeed.GetSignature()
	if len(signature) == 0 {
		// Gen1 uses fixed 128-byte RSA-1024 signatures
		signature = make([]byte, 128)
	}
	if len(signature) != 128 {
		return nil, fmt.Errorf("invalid signature length: got %d, want 128", len(signature))
	}

	transferValue := append(canvas, signature...)
	return transferValue, nil
}

// marshalDetailedSpeedBlock marshals a single VuDetailedSpeedBlock (64 bytes).
func (opts MarshalOptions) marshalDetailedSpeedBlock(block *vuv1.DetailedSpeedGen1_DetailedSpeedBlock) ([]byte, error) {
	if block == nil {
		return nil, fmt.Errorf("block cannot be nil")
	}

	const lenDetailedSpeedBlock = 64

	// Use raw data painting strategy if available
	var canvas [lenDetailedSpeedBlock]byte
	if rawData := block.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenDetailedSpeedBlock {
			return nil, fmt.Errorf("invalid raw_data length for VuDetailedSpeedBlock: got %d, want %d", len(rawData), lenDetailedSpeedBlock)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// Marshal begin date (4 bytes)
	beginDateBytes, err := opts.MarshalTimeReal(block.GetBeginDate())
	if err != nil {
		return nil, fmt.Errorf("marshal begin date: %w", err)
	}
	copy(canvas[offset:offset+4], beginDateBytes)
	offset += 4

	// Marshal 60 speed values (1 byte each)
	speeds := block.GetSpeedsKmh()
	if len(speeds) != 60 {
		return nil, fmt.Errorf("invalid number of speeds: got %d, want 60", len(speeds))
	}
	for i, speed := range speeds {
		canvas[offset+i] = byte(speed)
	}

	return canvas[:], nil
}

// anonymizeDetailedSpeedGen1 anonymizes Gen1 Detailed Speed data.
func (opts AnonymizeOptions) anonymizeDetailedSpeedGen1(ds *vuv1.DetailedSpeedGen1) *vuv1.DetailedSpeedGen1 {
	if ds == nil {
		return nil
	}
	result := proto.Clone(ds).(*vuv1.DetailedSpeedGen1)

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize blocks (timestamp only - speed values are not PII)
	var anonymizedBlocks []*vuv1.DetailedSpeedGen1_DetailedSpeedBlock
	for _, block := range result.GetSpeedBlocks() {
		if block == nil {
			continue
		}
		anonBlock := proto.Clone(block).(*vuv1.DetailedSpeedGen1_DetailedSpeedBlock)

		// Anonymize begin date
		anonBlock.SetBeginDate(ddOpts.AnonymizeTimestamp(block.GetBeginDate()))

		// Speed values are not PII - keep as-is
		// (speeds_kmh array contains actual speed measurements which are not personally identifiable)

		// Clear raw_data
		anonBlock.ClearRawData()

		anonymizedBlocks = append(anonymizedBlocks, anonBlock)
	}
	result.SetSpeedBlocks(anonymizedBlocks)

	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Clear raw_data to force semantic marshalling
	result.ClearRawData()

	return result
}
