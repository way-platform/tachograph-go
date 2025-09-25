package tachograph

import (
	"encoding/binary"
	"io"
)

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
