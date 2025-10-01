package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalRawCardFile parses raw card data.
func UnmarshalRawCardFile(input []byte) (*cardv1.RawCardFile, error) {
	var output cardv1.RawCardFile
	sc := bufio.NewScanner(bytes.NewReader(input))
	sc.Split(scanCardFile)
	for sc.Scan() {
		record, err := unmarshalRawCardFileRecord(sc.Bytes())
		if err != nil {
			return nil, err
		}
		output.SetRecords(append(output.GetRecords(), record))
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &output, nil
}

// MarshalRawCardFile serializes a RawCardFile into binary format.
func MarshalRawCardFile(file *cardv1.RawCardFile) ([]byte, error) {
	var result []byte
	for _, record := range file.GetRecords() {
		// Write tag (FID + appendix)
		result = binary.BigEndian.AppendUint16(result, uint16(record.GetTag()>>8))
		result = append(result, byte(record.GetTag()&0xFF))
		// Write length
		result = binary.BigEndian.AppendUint16(result, uint16(record.GetLength()))
		// Write value
		result = append(result, record.GetValue()...)
	}
	return result, nil
}

// scanCardFile is a [bufio.SplitFunc] that splits a card file into separate TLV records.
func scanCardFile(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Need at least 5 bytes for TLV header (3 bytes tag + 2 bytes length)
	if len(data) < 5 {
		if atEOF {
			if len(data) == 0 {
				// No more data - this is normal EOF
				return 0, nil, nil
			}
			// We have some data but not enough for a complete header
			return 0, nil, io.ErrUnexpectedEOF
		}
		// Request more data
		return 0, nil, nil
	}
	// Read the length field (bytes 3-4, big-endian)
	length := binary.BigEndian.Uint16(data[3:5])
	// Calculate total record size: 5 bytes header + length bytes value
	totalSize := 5 + int(length)
	// Check if we have enough data for the complete record
	if len(data) < totalSize {
		if atEOF {
			// We're at EOF but don't have enough data - this is an error condition
			return 0, nil, io.ErrUnexpectedEOF
		}
		// Request more data
		return 0, nil, nil
	}
	// We have a complete TLV record.
	return totalSize, data[:totalSize], nil
}

// unmarshalRawCardFileRecord unmarshals a single raw card file record
func unmarshalRawCardFileRecord(input []byte) (*cardv1.RawCardFile_Record, error) {
	var output cardv1.RawCardFile_Record
	// Parse tag: FID (2 bytes) + appendix (1 byte)
	fid := binary.BigEndian.Uint16(input[0:2])
	appendix := input[2]
	output.SetTag((int32(fid) << 8) | int32(appendix))
	// Parse length (2 bytes)
	length := binary.BigEndian.Uint16(input[3:5])
	output.SetLength(int32(length))
	// Parse value - make a copy since input slice may be reused by scanner
	value := make([]byte, length)
	copy(value, input[5:5+length])
	output.SetValue(value)
	// Determine content type and generation based on appendix byte
	// Per Chapter 12: Appendix encodes both content type and generation in bit pattern
	// Bit 0 (LSB): 0 = DATA, 1 = SIGNATURE
	// Bit 1: 0 = Tachograph DF (Gen1), 1 = Tachograph_G2 DF (Gen2)
	//
	// Examples:
	// 0x00 (0b00) = DATA in Tachograph DF (Gen1)
	// 0x01 (0b01) = SIGNATURE in Tachograph DF (Gen1)
	// 0x02 (0b10) = DATA in Tachograph_G2 DF (Gen2)
	// 0x03 (0b11) = SIGNATURE in Tachograph_G2 DF (Gen2)
	// Extract content type from bit 0
	if appendix&0x01 != 0 {
		output.SetContentType(cardv1.ContentType_SIGNATURE)
	} else {
		output.SetContentType(cardv1.ContentType_DATA)
	}
	// Extract generation from bit 1
	if appendix&0x02 != 0 {
		output.SetGeneration(ddv1.Generation_GENERATION_2)
	} else {
		output.SetGeneration(ddv1.Generation_GENERATION_1)
	}
	// Map FID to elementary file type
	fileType, _ := mapFidToElementaryFileType(fid)
	output.SetFile(fileType)
	return &output, nil
}
