package card

import (
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalLoadUnloadOperations parses the EF_Load_Unload_Operations file.
//
// The data type `CardLoadUnloadOperations` is specified in the Data Dictionary, Section 2.24c.
//
// ASN.1 Definition:
//
//	CardLoadUnloadOperations ::= SEQUENCE {
//	    loadUnloadPointerNewestRecord INTEGER(0..NoOfLoadUnloadRecords -1),
//	    cardLoadUnloadRecords SET SIZE (NoOfLoadUnloadRecords) OF CardLoadUnloadRecord
//	}
//
// Binary Layout:
//   - Bytes 0-1: loadUnloadPointerNewestRecord (2 bytes, big-endian)
//   - Bytes 2..N: Array of CardLoadUnloadRecord (20 bytes each)
func (opts UnmarshalOptions) unmarshalLoadUnloadOperations(data []byte) (*cardv1.LoadUnloadOperations, error) {
	const (
		lenPointer = 2
		lenRecord  = 20
	)

	if len(data) < lenPointer {
		return nil, fmt.Errorf("data too short for CardLoadUnloadOperations: got %d, want at least %d", len(data), lenPointer)
	}

	recordsData := data[lenPointer:]
	if len(recordsData)%lenRecord != 0 {
		return nil, fmt.Errorf("invalid data length for CardLoadUnloadOperations records: %d is not a multiple of %d", len(recordsData), lenRecord)
	}

	newestRecordIndex := int32(binary.BigEndian.Uint16(data[0:lenPointer]))

	target := &cardv1.LoadUnloadOperations{}
	target.SetNewestRecordIndex(newestRecordIndex)

	count := len(recordsData) / lenRecord
	records := make([]*cardv1.LoadUnloadOperations_Record, count)

	for i := 0; i < count; i++ {
		start := i * lenRecord
		ddRecord, err := opts.UnmarshalCardLoadUnloadRecord(recordsData[start : start+lenRecord])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal load/unload record %d: %w", i, err)
		}

		r := &cardv1.LoadUnloadOperations_Record{}
		r.SetTimestamp(ddRecord.GetTimeStamp())
		r.SetOperationType(ddRecord.GetOperationType())
		r.SetGnssPlaceAuthRecord(ddRecord.GetGnssPlaceAuthRecord())
		if ddRecord.HasVehicleOdometerKm() {
			r.SetVehicleOdometerKm(ddRecord.GetVehicleOdometerKm())
		}
		records[i] = r
	}
	target.SetRecords(records)

	return target, nil
}

// MarshalLoadUnloadOperations marshals the EF_Load_Unload_Operations file.
func (opts MarshalOptions) MarshalLoadUnloadOperations(msg *cardv1.LoadUnloadOperations) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}

	var dst []byte
	dst = binary.BigEndian.AppendUint16(dst, uint16(msg.GetNewestRecordIndex()))

	for i, record := range msg.GetRecords() {
		ddRecord := &ddv1.CardLoadUnloadRecord{}
		ddRecord.SetTimeStamp(record.GetTimestamp())
		ddRecord.SetOperationType(record.GetOperationType())
		ddRecord.SetGnssPlaceAuthRecord(record.GetGnssPlaceAuthRecord())
		if record.HasVehicleOdometerKm() {
			ddRecord.SetVehicleOdometerKm(record.GetVehicleOdometerKm())
		}

		b, err := opts.MarshalCardLoadUnloadRecord(ddRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal load/unload record %d: %w", i, err)
		}
		dst = append(dst, b...)
	}

	return dst, nil
}
