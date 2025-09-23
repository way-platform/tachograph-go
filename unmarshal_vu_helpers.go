package tachograph

import (
	"bytes"
	"encoding/binary"
	"time"
)

// VU-specific helper functions for parsing TV format data

// readUint8 reads a single byte
func readUint8(r *bytes.Reader) (uint8, error) {
	var value uint8
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

// readUint16 reads a 16-bit unsigned integer
func readUint16(r *bytes.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

// readUint32 reads a 32-bit unsigned integer
func readUint32(r *bytes.Reader) (uint32, error) {
	var value uint32
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

// readBytes reads n bytes
func readBytes(r *bytes.Reader, n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := r.Read(data)
	return data, err
}

// readVuString reads a string with codepage handling (similar to cards)
func readVuString(r *bytes.Reader, length int) (string, error) {
	data, err := readBytes(r, length)
	if err != nil {
		return "", err
	}
	return readString(bytes.NewReader(data), length), nil
}

// readVuTimeReal reads a TimeReal value (4 bytes) and converts to Unix timestamp
func readVuTimeReal(r *bytes.Reader) (int64, error) {
	value, err := readUint32(r)
	if err != nil {
		return 0, err
	}
	// TimeReal is seconds since 00:00:00 UTC, 1 January 1970
	return int64(value), nil
}

// readVuDatef reads a Datef value (4 bytes: year(2), month(1), day(1))
func readVuDatef(r *bytes.Reader) (string, error) {
	year, err := readUint16(r)
	if err != nil {
		return "", err
	}
	month, err := readUint8(r)
	if err != nil {
		return "", err
	}
	day, err := readUint8(r)
	if err != nil {
		return "", err
	}

	// Convert BCD to actual values - for VU, these might be plain binary
	yearActual := year
	monthActual := uint16(month)
	dayActual := uint16(day)

	// Format as ISO date
	date := time.Date(int(yearActual), time.Month(monthActual), int(dayActual), 0, 0, 0, 0, time.UTC)
	return date.Format("2006-01-02"), nil
}

// readVuOdometer reads an odometer value (3 bytes)
func readVuOdometer(r *bytes.Reader) (uint32, error) {
	data, err := readBytes(r, 3)
	if err != nil {
		return 0, err
	}
	// Convert 3-byte big-endian to uint32
	return uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2]), nil
}
