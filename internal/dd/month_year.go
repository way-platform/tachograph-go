package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalMonthYear parses BCD-encoded month and year data.
//
// The data type `monthYear` is specified in the Data Dictionary, Section 2.72
// as part of ExtendedSerialNumber.
//
// ASN.1 Definition:
//
//	monthYear BCDString(SIZE(2))
//
// Binary Layout (2 bytes):
//   - BCD-encoded MMYY format (2 bytes)
func (opts UnmarshalOptions) UnmarshalMonthYear(data []byte) (*ddv1.MonthYear, error) {
	const lenMonthYear = 2

	if len(data) != lenMonthYear {
		return nil, fmt.Errorf("invalid data length for MonthYear: got %d, want %d", len(data), lenMonthYear)
	}

	monthYear := &ddv1.MonthYear{}
	monthYear.SetRawData(data[:lenMonthYear])

	// Decode BCD month/year as 4-digit number MMYY
	monthYearInt, err := decodeBCD(data[:lenMonthYear])
	if err == nil && monthYearInt > 0 {
		month := int32(monthYearInt / 100)
		year := int32(monthYearInt % 100)

		// Convert 2-digit year to 4-digit (assuming 20xx for years 00-99)
		if year >= 0 && year <= 99 {
			year += 2000
		}

		monthYear.SetMonth(month)
		monthYear.SetYear(year)
	}

	return monthYear, nil
}

// AppendMonthYear appends BCD-encoded month and year data to dst.
//
// The data type `monthYear` is specified in the Data Dictionary, Section 2.72
// as part of ExtendedSerialNumber.
//
// ASN.1 Definition:
//
//	monthYear BCDString(SIZE(2))
//
// Binary Layout (2 bytes):
//   - BCD-encoded MMYY format (2 bytes)
func AppendMonthYear(dst []byte, monthYear *ddv1.MonthYear) ([]byte, error) {
	const lenMonthYear = 2

	// No nil check needed - protobuf returns zero values for nil, which are valid
	// This function only reads primitive int32 fields (month, year) and bytes

	// Use stack-allocated array for the canvas (fixed size, avoids heap allocation)
	var canvas [lenMonthYear]byte

	// Start with raw_data as canvas if available (raw data painting approach)
	if rawData := monthYear.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenMonthYear {
			return nil, fmt.Errorf("invalid raw_data length for MonthYear: got %d, want %d", len(rawData), lenMonthYear)
		}
		copy(canvas[:], rawData)
	}
	// Otherwise canvas is zero-initialized (Go default)

	// Paint semantic values over the canvas
	month := monthYear.GetMonth()
	year := monthYear.GetYear()

	// Convert 4-digit year to 2-digit for BCD encoding
	year2Digit := year % 100

	// Encode month as BCD
	canvas[0] = byte(((month / 10) << 4) | (month % 10))

	// Encode year as BCD
	canvas[1] = byte(((year2Digit / 10) << 4) | (year2Digit % 10))

	return append(dst, canvas[:]...), nil
}
