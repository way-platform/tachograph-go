package tachograph

import (
	"encoding/binary"
	"io"
)

// VU-specific binary parsing functions for reading structured data from byte slices
// These functions are used across multiple VU-related files for consistent data parsing

// readUint8FromBytes reads a single byte from a byte slice at the given offset
func readUint8FromBytes(data []byte, offset int) (uint8, int, error) {
	if offset >= len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	return data[offset], offset + 1, nil
}

// readBytesFromBytes reads n bytes from a byte slice at the given offset
func readBytesFromBytes(data []byte, offset int, n int) ([]byte, int, error) {
	if offset+n > len(data) {
		return nil, offset, io.ErrUnexpectedEOF
	}
	result := make([]byte, n)
	copy(result, data[offset:offset+n])
	return result, offset + n, nil
}

// readVuTimeRealFromBytes reads a TimeReal value (4 bytes) and converts to Unix timestamp
func readVuTimeRealFromBytes(data []byte, offset int) (int64, int, error) {
	if offset+4 > len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	value := binary.BigEndian.Uint32(data[offset:])
	// TimeReal is seconds since 00:00:00 UTC, 1 January 1970
	return int64(value), offset + 4, nil
}

