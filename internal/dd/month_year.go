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
func UnmarshalMonthYear(data []byte) (*ddv1.MonthYear, error) {
	const lenMonthYear = 2

	if len(data) < lenMonthYear {
		return nil, fmt.Errorf("insufficient data for MonthYear: got %d, want %d", len(data), lenMonthYear)
	}

	monthYear := &ddv1.MonthYear{}
	monthYear.SetEncoded(data[:lenMonthYear])

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

	if monthYear == nil {
		// Append zeros
		return append(dst, 0, 0), nil
	}

	// Prefer the original encoded bytes for perfect round-trip fidelity
	if encoded := monthYear.GetEncoded(); len(encoded) >= lenMonthYear {
		return append(dst, encoded[:lenMonthYear]...), nil
	}

	// Fall back to encoding from decoded values
	month := monthYear.GetMonth()
	year := monthYear.GetYear()

	// Convert 4-digit year to 2-digit for BCD encoding
	year2Digit := year % 100

	// Encode month as BCD
	monthBCD := byte(((month / 10) << 4) | (month % 10))
	dst = append(dst, monthBCD)

	// Encode year as BCD
	yearBCD := byte(((year2Digit / 10) << 4) | (year2Digit % 10))
	dst = append(dst, yearBCD)

	return dst, nil
}
