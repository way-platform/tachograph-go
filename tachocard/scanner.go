package tachocard

import (
	"encoding/binary"
	"fmt"
	"io"
)

// SplitFunc is a bufio.Scanner split function for Tachograph Card TLV records.
func SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// The constant header size for a TLV record is 5 bytes (3 for Tag, 2 for Length).
	const headerSize = 5
	// 1. Check for the header.
	if len(data) < headerSize {
		if atEOF && len(data) > 0 {
			// Reached EOF but have a fragment that can't be a full header.
			return len(data), nil, fmt.Errorf("incomplete TLV header: got %d bytes", len(data))
		}
		// Not enough data for a header, request more.
		return 0, nil, nil
	}
	// 2. Read the length of the value payload.
	valueLength := binary.BigEndian.Uint16(data[3:5])
	// 3. Calculate the total length of the full TLV record.
	totalLength := headerSize + int(valueLength)
	// 4. Check if the complete record is in the buffer.
	if len(data) < totalLength {
		if atEOF {
			// The buffer doesn't contain the full record, and we're at the end.
			return len(data), nil, io.ErrUnexpectedEOF
		}
		// Not enough data for the full record, request more.
		return 0, nil, nil
	}
	// 5. We have a full record, so we can return it.
	advance = totalLength
	token = data[:totalLength]
	err = nil
	return
}
