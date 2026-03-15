package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalDetailedSpeedGen2 parses Gen2 Detailed Speed data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen2 Detailed Speed structure uses RecordArray format.
//
// Note: Gen2 has no V2 variant - both V1 and V2 use the same structure.
func unmarshalDetailedSpeedGen2(value []byte) (*vuv1.DetailedSpeedGen2, error) {
	totalSize, signatureSize, err := sizeOfDetailedSpeedGen2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	detailedSpeed := &vuv1.DetailedSpeedGen2{}
	detailedSpeed.SetRawData(value)

	offset := 0

	// VuDetailedSpeedBlockRecordArray
	speedBlocks, bytesRead, err := parseVuDetailedSpeedBlockRecordArray(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuDetailedSpeedBlockRecordArray: %w", err)
	}
	detailedSpeed.SetSpeedBlocks(speedBlocks)
	offset += bytesRead

	detailedSpeed.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Detailed Speed Gen2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return detailedSpeed, nil
}

// MarshalDetailedSpeedGen2 marshals Gen2 Detailed Speed data.
func (opts MarshalOptions) MarshalDetailedSpeedGen2(detailedSpeed *vuv1.DetailedSpeedGen2) ([]byte, error) {
	if detailedSpeed == nil {
		return nil, fmt.Errorf("detailedSpeed cannot be nil")
	}

	raw := detailedSpeed.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	var result []byte
	marshalOpts := dd.MarshalOptions{}

	// VuDetailedSpeedBlockRecordArray (64 bytes per block)
	blocks := detailedSpeed.GetSpeedBlocks()
	speedBlockData, err := marshalVuDetailedSpeedBlocks(marshalOpts, blocks)
	if err != nil {
		return nil, fmt.Errorf("marshal VuDetailedSpeedBlockRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x01, 64, uint16(len(blocks)))
	result = append(result, speedBlockData...)

	result = appendSignature(result, detailedSpeed.GetSignature(), 0x02)
	return result, nil
}

// anonymizeDetailedSpeedGen2 anonymizes Gen2 Detailed Speed data.
// Speed data contains no PII; only signature and raw_data are cleared.
func (opts AnonymizeOptions) anonymizeDetailedSpeedGen2(ds *vuv1.DetailedSpeedGen2) *vuv1.DetailedSpeedGen2 {
	if ds == nil {
		return nil
	}

	result := &vuv1.DetailedSpeedGen2{}
	result.SetSpeedBlocks(ds.GetSpeedBlocks())
	result.SetSignature([]byte{})
	return result
}

// parseVuDetailedSpeedBlockRecordArray parses a VuDetailedSpeedBlockRecordArray.
//
// Each VuDetailedSpeedBlock is 64 bytes:
//   - speedBlockBeginDate: TimeReal = 4 bytes
//   - speedsPerSecond: 60 bytes (1 byte per km/h value)
func parseVuDetailedSpeedBlockRecordArray(data []byte, offset int) ([]*vuv1.DetailedSpeedGen2_DetailedSpeedBlock, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 64
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuDetailedSpeedBlock size %d, got %d", expectedRecordSize, recordSize)
	}

	var unmarshalOpts dd.UnmarshalOptions
	records := make([]*vuv1.DetailedSpeedGen2_DetailedSpeedBlock, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for DetailedSpeedBlock %d", i)
		}

		rec := data[recordStart:recordEnd]

		beginDate, err := unmarshalOpts.UnmarshalTimeReal(rec[:4])
		if err != nil {
			return nil, 0, fmt.Errorf("DetailedSpeedBlock %d begin date: %w", i, err)
		}

		speedsKmh := make([]int32, 60)
		for j := range 60 {
			speedsKmh[j] = int32(rec[4+j])
		}

		block := &vuv1.DetailedSpeedGen2_DetailedSpeedBlock{}
		block.SetBeginDate(beginDate)
		block.SetSpeedsKmh(speedsKmh)

		records = append(records, block)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// marshalVuDetailedSpeedBlocks marshals detailed speed blocks to binary.
func marshalVuDetailedSpeedBlocks(opts dd.MarshalOptions, blocks []*vuv1.DetailedSpeedGen2_DetailedSpeedBlock) ([]byte, error) {
	result := make([]byte, 0, len(blocks)*64)
	for i, block := range blocks {
		beginDateBytes, err := opts.MarshalTimeReal(block.GetBeginDate())
		if err != nil {
			return nil, fmt.Errorf("DetailedSpeedBlock %d begin date: %w", i, err)
		}
		result = append(result, beginDateBytes...)

		speeds := block.GetSpeedsKmh()
		for j := range 60 {
			if j < len(speeds) {
				result = append(result, byte(speeds[j]))
			} else {
				result = append(result, 0)
			}
		}
	}
	return result, nil
}
