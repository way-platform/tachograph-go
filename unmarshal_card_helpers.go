package tachograph

import (
	"bytes"
	"encoding/binary"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromBCD(b byte) int {
	return (int(b>>4) * 10) + int(b&0x0F)
}

func readString(r *bytes.Reader, len int) string {
	b := make([]byte, len)
	r.Read(b)
	// Trim trailing spaces and null bytes
	b = bytes.TrimRight(b, " \x00")
	
	// Check if the result is valid UTF-8, if not convert to hex representation
	if !isValidUTF8(b) {
		return bytesToHexString(b)
	}
	
	return string(b)
}

// isValidUTF8 checks if the byte slice contains valid UTF-8
func isValidUTF8(b []byte) bool {
	// Check if all bytes are printable ASCII or valid UTF-8
	for _, byte := range b {
		if byte < 0x20 || byte > 0x7E {
			// Contains non-printable characters, treat as binary
			return false
		}
	}
	return true
}

// bytesToHexString converts binary data to a hex string representation
func bytesToHexString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	result := make([]byte, len(b)*2)
	const hexDigits = "0123456789ABCDEF"
	for i, byte := range b {
		result[i*2] = hexDigits[byte>>4]
		result[i*2+1] = hexDigits[byte&0x0F]
	}
	return string(result)
}

func readTimeReal(r *bytes.Reader) *timestamppb.Timestamp {
	var timeVal uint32
	binary.Read(r, binary.BigEndian, &timeVal)
	if timeVal == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(int64(timeVal), 0))
}

func readDatef(r *bytes.Reader) *timestamppb.Timestamp {
	b := make([]byte, 4)
	r.Read(b)
	year := fromBCD(b[0])*100 + fromBCD(b[1])
	month := fromBCD(b[2])
	day := fromBCD(b[3])
	if year == 0 || month == 0 || day == 0 {
		return nil
	}
	return timestamppb.New(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
}
