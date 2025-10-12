package vu

import (
	"encoding/binary"
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfCardDownload calculates the size of a card download transfer.
// Card downloads contain TLV-formatted card data (Tag-Length-Value format).
//
// The format is defined in Chapter 12, Section 3.4.2:
// - Tag: 3 bytes (FID + '00' for data, FID + '01' for signature)
// - Length: 2 bytes (big-endian, number of bytes in value field)
// - Value: Variable length
//
// This function parses through the TLV records until all data is consumed.
func sizeOfCardDownload(data []byte, transferType vuv1.TransferType) (int, error) {
	offset := 0

	// Parse TLV records until we've consumed all data
	for offset < len(data) {
		// Need at least 5 bytes for TLV header (3-byte tag + 2-byte length)
		const tlvHeaderSize = 5
		if len(data[offset:]) < tlvHeaderSize {
			// If we have less than a full header, we've reached the end
			break
		}

		// Read length (bytes 3-4 of TLV header, big-endian)
		length := binary.BigEndian.Uint16(data[offset+3:])

		// TLV record size = header + value
		recordSize := tlvHeaderSize + int(length)

		// Check if we have enough data for this record
		if offset+recordSize > len(data) {
			return 0, fmt.Errorf("incomplete TLV record at offset %d: need %d bytes, have %d", offset, recordSize, len(data)-offset)
		}

		offset += recordSize
	}

	return offset, nil
}

// ===== Unmarshal Functions =====

// unmarshalCardDownload parses a card download transfer.
// The card download payload is raw TLV-formatted card data.
