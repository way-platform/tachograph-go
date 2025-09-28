package tachograph

import (
	"fmt"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalDate unmarshals a BCD-encoded date from a byte slice
func unmarshalDate(data []byte) (*datadictionaryv1.Date, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for date")
	}
	// Parse BCD format: YYYYMMDD
	year := int32(int32((data[0]&0xF0)>>4)*1000 + int32(data[0]&0x0F)*100 + int32((data[1]&0xF0)>>4)*10 + int32(data[1]&0x0F))
	month := int32(int32((data[2]&0xF0)>>4)*10 + int32(data[2]&0x0F))
	day := int32(int32((data[3]&0xF0)>>4)*10 + int32(data[3]&0x0F))
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
