package tachograph

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

func bcdBytesToInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	s := hex.EncodeToString(b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid BCD value: %s", s)
	}
	return int(i), nil
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

// readTimeReal reads a TimeReal value (4 bytes) and converts to Timestamp
// TimeReal is INTEGER (0..2^32-1) representing seconds since epoch
func readTimeReal(r *bytes.Reader) *timestamppb.Timestamp {
	var timeVal uint32
	binary.Read(r, binary.BigEndian, &timeVal)
	if timeVal == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(int64(timeVal), 0))
}

// readDatef reads a Datef value (4 bytes BCD) and converts to Date
// Datef is OCTET STRING (SIZE(4)) with BCD-encoded YYYYMMDD format
func readDatef(r *bytes.Reader) (*datadictionaryv1.Date, error) {
	b := make([]byte, 4)
	r.Read(b)

	// Parse BCD format: YYYYMMDD
	year := int32(int32((b[0]&0xF0)>>4)*1000 + int32(b[0]&0x0F)*100 + int32((b[1]&0xF0)>>4)*10 + int32(b[1]&0x0F))
	month := int32(int32((b[2]&0xF0)>>4)*10 + int32(b[2]&0x0F))
	day := int32(int32((b[3]&0xF0)>>4)*10 + int32(b[3]&0x0F))

	// Validate the date
	if year < 1900 || year > 9999 || month < 1 || month > 12 || day < 1 || day > 31 {
		return nil, nil // Return nil for invalid or zero dates
	}

	date := &datadictionaryv1.Date{}
	date.SetYear(year)
	date.SetMonth(month)
	date.SetDay(day)
	return date, nil
}
