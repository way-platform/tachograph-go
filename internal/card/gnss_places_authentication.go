package card

import (
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalGnssPlacesAuthentication parses the EF_GNSS_Places_Authentication file.
//
// The data type `GNSSAuthAccumulatedDriving` is specified in the Data Dictionary, Section 2.79a.
//
// ASN.1 Definition:
//
//	GNSSAuthAccumulatedDriving ::= SEQUENCE {
//	    gnssAuthADPointerNewestRecord INTEGER(0..NoOfGNSSADRecords -1),
//	    gnssAuthStatusADRecords SET SIZE(NoOfGNSSADRecords) OF GNSSAuthStatusADRecord
//	}
//
// Binary Layout:
//   - Bytes 0-1: gnssAuthADPointerNewestRecord (2 bytes, big-endian)
//   - Bytes 2..N: Array of GNSSAuthStatusADRecord (5 bytes each)
//
// GNSSAuthStatusADRecord has the same binary layout as PlaceAuthStatusRecord:
// 4B TimeReal + 1B PositionAuthenticationStatus.
func (opts UnmarshalOptions) unmarshalGnssPlacesAuthentication(data []byte) (*cardv1.GnssPlacesAuthentication, error) {
	const (
		lenPointer = 2
		lenRecord  = 5
	)

	if len(data) < lenPointer {
		return nil, fmt.Errorf("data too short for GNSSAuthAccumulatedDriving: got %d, want at least %d", len(data), lenPointer)
	}

	recordsData := data[lenPointer:]
	if len(recordsData)%lenRecord != 0 {
		return nil, fmt.Errorf("invalid data length for GNSSAuthAccumulatedDriving records: %d is not a multiple of %d", len(recordsData), lenRecord)
	}

	newestRecordIndex := int32(binary.BigEndian.Uint16(data[0:lenPointer]))

	target := &cardv1.GnssPlacesAuthentication{}
	target.SetNewestRecordIndex(newestRecordIndex)

	count := len(recordsData) / lenRecord
	records := make([]*cardv1.GnssPlacesAuthentication_Record, count)

	for i := 0; i < count; i++ {
		start := i * lenRecord
		// GNSSAuthStatusADRecord has the same layout as PlaceAuthStatusRecord.
		ddRecord, err := opts.UnmarshalPlaceAuthStatusRecord(recordsData[start : start+lenRecord])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal GNSS auth status record %d: %w", i, err)
		}

		r := &cardv1.GnssPlacesAuthentication_Record{}
		r.SetTimestamp(ddRecord.GetEntryTime())
		r.SetAuthenticationStatus(ddRecord.GetAuthenticationStatus())
		records[i] = r
	}
	target.SetRecords(records)

	return target, nil
}
