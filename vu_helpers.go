package tachograph

import (
	"encoding/binary"
	"io"
	"time"
)

// VU-specific helper functions for parsing byte slice data

// readUint8FromBytes reads a single byte from a byte slice at the given offset
func readUint8FromBytes(data []byte, offset int) (uint8, int, error) {
	if offset >= len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	return data[offset], offset + 1, nil
}

// readUint16FromBytes reads a 16-bit unsigned integer from a byte slice at the given offset
func readUint16FromBytes(data []byte, offset int) (uint16, int, error) {
	if offset+2 > len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	value := binary.BigEndian.Uint16(data[offset:])
	return value, offset + 2, nil
}

// readUint32FromBytes reads a 32-bit unsigned integer from a byte slice at the given offset
func readUint32FromBytes(data []byte, offset int) (uint32, int, error) {
	if offset+4 > len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	value := binary.BigEndian.Uint32(data[offset:])
	return value, offset + 4, nil
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

// readVuStringFromBytes reads a string with codepage handling from a byte slice
func readVuStringFromBytes(data []byte, offset int, length int) (string, int, error) {
	if offset+length > len(data) {
		return "", offset, io.ErrUnexpectedEOF
	}
	stringData := data[offset : offset+length]
	result := readStringFromBytes(stringData, length)
	return result, offset + length, nil
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

// readVuDatefFromBytes reads a Datef value (4 bytes: year(2), month(1), day(1))
func readVuDatefFromBytes(data []byte, offset int) (string, int, error) {
	if offset+4 > len(data) {
		return "", offset, io.ErrUnexpectedEOF
	}

	year := binary.BigEndian.Uint16(data[offset:])
	month := data[offset+2]
	day := data[offset+3]

	// Convert BCD to actual values - for VU, these might be plain binary
	yearActual := year
	monthActual := uint16(month)
	dayActual := uint16(day)

	// Format as ISO date
	date := time.Date(int(yearActual), time.Month(monthActual), int(dayActual), 0, 0, 0, 0, time.UTC)
	return date.Format("2006-01-02"), offset + 4, nil
}

// readVuOdometerFromBytes reads an odometer value (3 bytes)
func readVuOdometerFromBytes(data []byte, offset int) (uint32, int, error) {
	if offset+3 > len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	// Convert 3-byte big-endian to uint32
	value := uint32(data[offset])<<16 | uint32(data[offset+1])<<8 | uint32(data[offset+2])
	return value, offset + 3, nil
}
