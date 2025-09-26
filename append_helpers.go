package tachograph

import (
	"encoding/binary"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"
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

	dst = append(dst, byte((year/1000)%10<<4| (year/100)%10))
	dst = append(dst, byte((year/10)%10<<4| year%10))
	dst = append(dst, byte((month/10)%10<<4| month%10))
	dst = append(dst, byte((day/10)%10<<4| day%10))
	return dst
}

func appendOdometer(dst []byte, odometer uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, odometer)
	return append(dst, b[1:]...)
}

// appendString appends a fixed-length string, padding with spaces.
func appendString(dst []byte, s string, length int) ([]byte, error) {
	if len(s) > length {
		return nil, fmt.Errorf("string '%s' is longer than the allowed length %d", s, length)
	}
	result := make([]byte, length)
	copy(result, []byte(s))
	for i := len(s); i < length; i++ {
		result[i] = ' '
	}
	return append(dst, result...), nil
}

// appendBCDNation appends a BCD-encoded nation number.
func appendBCDNation(dst []byte, nation string) ([]byte, error) {
	// This is a placeholder. A real implementation would convert the nation string
	// to its numeric code and then to BCD.
	return append(dst, 0), nil // Append a single zero byte for now
}

// AppendVehicleRegistration appends a VehicleRegistrationIdentification structure.
func AppendVehicleRegistration(dst []byte, nation string, number string) ([]byte, error) {
	// This is also a placeholder.
	dst = append(dst, 0) // Nation
	dst = append(dst, []byte(strings.Repeat(" ", 14))...)
	copy(dst[1:], []byte(number))
	return dst, nil
}
