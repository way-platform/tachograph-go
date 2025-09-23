package tachograph

import (
	"bytes"
	"encoding/binary"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// VU-specific helper functions for marshaling TV format data

// appendUint8 appends a single byte
func appendUint8(buf *bytes.Buffer, value uint8) {
	buf.WriteByte(value)
}

// appendUint16 appends a 16-bit unsigned integer in big-endian format
func appendUint16(buf *bytes.Buffer, value uint16) {
	binary.Write(buf, binary.BigEndian, value)
}

// appendUint32 appends a 32-bit unsigned integer in big-endian format
func appendUint32(buf *bytes.Buffer, value uint32) {
	binary.Write(buf, binary.BigEndian, value)
}

// appendVuBytes appends raw bytes
func appendVuBytes(buf *bytes.Buffer, data []byte) {
	buf.Write(data)
}

// appendVuString appends a string with proper padding to specified length
func appendVuString(buf *bytes.Buffer, str string, length int) {
	data := make([]byte, length)
	copy(data, []byte(str))
	// Pad with spaces if necessary
	for i := len(str); i < length; i++ {
		data[i] = ' '
	}
	buf.Write(data)
}

// appendVuTimeReal appends a TimeReal value (4 bytes) from a timestamp
func appendVuTimeReal(buf *bytes.Buffer, ts *timestamppb.Timestamp) {
	if ts == nil {
		appendUint32(buf, 0)
		return
	}
	// Convert timestamp to Unix seconds
	appendUint32(buf, uint32(ts.GetSeconds()))
}

// appendVuDatef appends a Datef value (4 bytes: year(2), month(1), day(1))
func appendVuDatef(buf *bytes.Buffer, dateStr string) {
	if dateStr == "" {
		appendUint16(buf, 0)
		appendUint8(buf, 0)
		appendUint8(buf, 0)
		return
	}

	// Parse ISO date format
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Default to zeros on parse error
		appendUint16(buf, 0)
		appendUint8(buf, 0)
		appendUint8(buf, 0)
		return
	}

	// Append as plain binary (not BCD for VU)
	appendUint16(buf, uint16(date.Year()))
	appendUint8(buf, uint8(date.Month()))
	appendUint8(buf, uint8(date.Day()))
}

// appendVuOdometer appends an odometer value (3 bytes)
func appendVuOdometer(buf *bytes.Buffer, value uint32) {
	// Convert uint32 to 3-byte big-endian
	buf.WriteByte(byte(value >> 16))
	buf.WriteByte(byte(value >> 8))
	buf.WriteByte(byte(value))
}

// appendVuTag appends a 2-byte VU tag
func appendVuTag(buf *bytes.Buffer, tag uint16) {
	appendUint16(buf, tag)
}
