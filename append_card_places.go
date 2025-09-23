package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendPlaces appends the binary representation of Places to dst.
func AppendPlaces(dst []byte, p *cardv1.Places) ([]byte, error) {
	if p == nil {
		return dst, nil
	}
	dst = append(dst, byte(p.GetNewestRecordIndex()))

	var err error
	for _, rec := range p.GetRecords() {
		dst, err = AppendPlaceRecord(dst, rec)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendPlaceRecord appends a single 10-byte place record.
func AppendPlaceRecord(dst []byte, rec *cardv1.Places_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, 10)...), nil
	}
	dst = appendTimeReal(dst, rec.GetEntryTime())
	dst = append(dst, byte(rec.GetEntryType()))
	dst = append(dst, byte(rec.GetDailyWorkPeriodCountry()))
	dst = append(dst, byte(rec.GetDailyWorkPeriodRegion()))
	dst = appendOdometer(dst, rec.GetVehicleOdometerKm())
	return dst, nil
}
