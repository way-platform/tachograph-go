package tachograph

import (
	"bytes"
	"encoding/binary"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// appendTimeReal appends a 4-byte TimeReal value.
func appendTimeReal(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	return binary.BigEndian.AppendUint32(dst, uint32(ts.GetSeconds()))
}

// appendDatef appends a 4-byte BCD-encoded date.
func appendDatef(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := ts.AsTime().Year()
	month := int(ts.AsTime().Month())
	day := ts.AsTime().Day()

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
}

// appendDate appends a 4-byte BCD-encoded date from the new Date type.
func appendDate(dst []byte, date *datadictionaryv1.Date) []byte {
	if date == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := int(date.GetYear())
	month := int(date.GetMonth())
	day := int(date.GetDay())

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
}

// readTimeReal reads a TimeReal value (4 bytes) from a bytes.Reader and converts to Timestamp
func readTimeReal(r *bytes.Reader) *timestamppb.Timestamp {
	var timeVal uint32
	_ = binary.Read(r, binary.BigEndian, &timeVal) // ignore error as we're reading from in-memory buffer
	if timeVal == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(int64(timeVal), 0))
}

// readDatef reads a Datef value (4 bytes BCD) from a bytes.Reader and converts to Date
func readDatef(r *bytes.Reader) (*datadictionaryv1.Date, error) {
	b := make([]byte, 4)
	_, _ = r.Read(b) // ignore error as we're reading from in-memory buffer

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
