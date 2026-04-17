package card

import (
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalLoadTypeEntries parses the EF_Load_Type_Entries file.
//
// The data type `CardLoadTypeEntries` is specified in the Data Dictionary, Section 2.24a.
//
// ASN.1 Definition:
//
//	CardLoadTypeEntries ::= SEQUENCE {
//	    loadTypeEntryPointerNewestRecord INTEGER(0..NoOfLoadTypeEntryRecords-1),
//	    cardLoadTypeEntryRecords SET SIZE(NoOfLoadTypeEntryRecords) OF CardLoadTypeEntryRecord
//	}
//
// Binary Layout:
//   - Bytes 0-1: loadTypeEntryPointerNewestRecord (2 bytes, big-endian)
//   - Bytes 2..N: Array of CardLoadTypeEntryRecord (5 bytes each)
func (opts UnmarshalOptions) unmarshalLoadTypeEntries(data []byte) (*cardv1.LoadTypeEntries, error) {
	const (
		lenPointer = 2
		lenRecord  = 5
	)

	if len(data) < lenPointer {
		return nil, fmt.Errorf("data too short for CardLoadTypeEntries: got %d, want at least %d", len(data), lenPointer)
	}

	recordsData := data[lenPointer:]
	if len(recordsData)%lenRecord != 0 {
		return nil, fmt.Errorf("invalid data length for CardLoadTypeEntries records: %d is not a multiple of %d", len(recordsData), lenRecord)
	}

	newestRecordIndex := int32(binary.BigEndian.Uint16(data[0:lenPointer]))

	target := &cardv1.LoadTypeEntries{}
	target.SetNewestRecordIndex(newestRecordIndex)

	count := len(recordsData) / lenRecord
	records := make([]*cardv1.LoadTypeEntries_Record, count)

	for i := 0; i < count; i++ {
		start := i * lenRecord
		ddRecord, err := opts.UnmarshalCardLoadTypeEntryRecord(recordsData[start : start+lenRecord])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal load type entry record %d: %w", i, err)
		}

		r := &cardv1.LoadTypeEntries_Record{}
		r.SetTimestamp(ddRecord.GetTimeStamp())
		r.SetLoadTypeEntered(ddRecord.GetLoadTypeEntered())
		records[i] = r
	}
	target.SetRecords(records)

	return target, nil
}

// MarshalLoadTypeEntries marshals the EF_Load_Type_Entries file.
func (opts MarshalOptions) MarshalLoadTypeEntries(msg *cardv1.LoadTypeEntries) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}

	var dst []byte
	dst = binary.BigEndian.AppendUint16(dst, uint16(msg.GetNewestRecordIndex()))

	for i, record := range msg.GetRecords() {
		ddRecord := &ddv1.CardLoadTypeEntryRecord{}
		ddRecord.SetTimeStamp(record.GetTimestamp())
		ddRecord.SetLoadTypeEntered(record.GetLoadTypeEntered())

		b, err := opts.MarshalCardLoadTypeEntryRecord(ddRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal load type entry record %d: %w", i, err)
		}
		dst = append(dst, b...)
	}

	return dst, nil
}
