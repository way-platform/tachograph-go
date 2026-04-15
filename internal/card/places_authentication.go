package card

import (
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalPlacesAuthentication parses the EF_Places_Authentication file.
//
// The data type `CardPlaceAuthDailyWorkPeriod` is specified in the Data Dictionary, Section 2.26a.
//
// ASN.1 Definition:
//
//	CardPlaceAuthDailyWorkPeriod ::= SEQUENCE {
//	    placeAuthPointerNewestRecord INTEGER(0..NoOfCardPlaceRecords -1),
//	    placeAuthStatusRecords SET SIZE(NoOfCardPlaceRecords) OF PlaceAuthStatusRecord
//	}
//
// Binary Layout:
//   - Bytes 0-1: placeAuthPointerNewestRecord (2 bytes, big-endian)
//   - Bytes 2..N: Array of PlaceAuthStatusRecord (5 bytes each)
func (opts UnmarshalOptions) unmarshalPlacesAuthentication(data []byte) (*cardv1.PlacesAuthentication, error) {
	const (
		lenPointer = 2
		lenRecord  = 5
	)

	if len(data) < lenPointer {
		return nil, fmt.Errorf("data too short for CardPlaceAuthDailyWorkPeriod: got %d, want at least %d", len(data), lenPointer)
	}

	recordsData := data[lenPointer:]
	if len(recordsData)%lenRecord != 0 {
		return nil, fmt.Errorf("invalid data length for CardPlaceAuthDailyWorkPeriod records: %d is not a multiple of %d", len(recordsData), lenRecord)
	}

	newestRecordIndex := int32(binary.BigEndian.Uint16(data[0:lenPointer]))

	target := &cardv1.PlacesAuthentication{}
	target.SetNewestRecordIndex(newestRecordIndex)

	count := len(recordsData) / lenRecord
	records := make([]*cardv1.PlacesAuthentication_Record, count)

	for i := 0; i < count; i++ {
		start := i * lenRecord
		ddRecord, err := opts.UnmarshalPlaceAuthStatusRecord(recordsData[start : start+lenRecord])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal place auth status record %d: %w", i, err)
		}

		r := &cardv1.PlacesAuthentication_Record{}
		r.SetEntryTime(ddRecord.GetEntryTime())
		r.SetAuthenticationStatus(ddRecord.GetAuthenticationStatus())
		records[i] = r
	}
	target.SetRecords(records)

	return target, nil
}
