package card

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalSpecificConditionsG2 unmarshals specific conditions data from a card EF (Gen2 format with circular buffer).
//
// The data type `SpecificConditions` is specified in the Data Dictionary, Section 2.153.
//
// ASN.1 Definition (Gen2):
//
//	SpecificConditions ::= SEQUENCE {
//	    conditionPointerNewestRecord NoOfSpecificConditionRecords,
//	    specificConditionRecords SET SIZE(NoOfSpecificConditionRecords) OF SpecificConditionRecord
//	}
func (opts UnmarshalOptions) unmarshalSpecificConditionsG2(data []byte) (*cardv1.SpecificConditionsG2, error) {
	const (
		lenPointer                 = 2 // Gen2 uses 2-byte pointer (INT(0..65535))
		lenSpecificConditionRecord = 5
	)

	if len(data) < lenPointer {
		return nil, fmt.Errorf("insufficient data for Gen2 specific conditions: got %d bytes, need at least %d", len(data), lenPointer)
	}

	target := &cardv1.SpecificConditionsG2{}

	// Save complete raw data for painting (includes pointer + records + trailing bytes)
	target.SetRawData(data)

	// Read newest record pointer (2 bytes)
	newestRecordPointer := binary.BigEndian.Uint16(data[0:lenPointer])
	target.SetNewestRecordIndex(int32(newestRecordPointer))

	// Parse records
	recordsData := data[lenPointer:]
	var records []*ddv1.SpecificConditionRecord

	offset := 0
	for offset+lenSpecificConditionRecord <= len(recordsData) {
		record, err := opts.UnmarshalSpecificConditionRecord(recordsData[offset : offset+lenSpecificConditionRecord])
		if err != nil {
			break // Stop parsing on error, but return what we have
		}
		records = append(records, record)
		offset += lenSpecificConditionRecord
	}

	target.SetRecords(records)
	return target, nil
}

// appendCardSpecificConditionsG2 appends Gen2 specific conditions data to a byte slice.
//
// The data type `SpecificConditions` is specified in the Data Dictionary, Section 2.153.
//
// ASN.1 Definition (Gen2):
//
//	SpecificConditions ::= SEQUENCE {
//	    conditionPointerNewestRecord NoOfSpecificConditionRecords,
//	    specificConditionRecords SET SIZE(NoOfSpecificConditionRecords) OF SpecificConditionRecord
//	}
func appendCardSpecificConditionsG2(data []byte, conditions *cardv1.SpecificConditionsG2) ([]byte, error) {
	if conditions == nil {
		return data, nil
	}

	// Use raw_data as canvas if available (includes pointer + records + trailing bytes)
	if rawData := conditions.GetRawData(); len(rawData) > 0 {
		const (
			lenPointer = 2
			recordSize = 5
		)

		// Make a copy to use as canvas
		canvas := make([]byte, len(rawData))
		copy(canvas, rawData)

		// Paint newest record pointer over canvas
		binary.BigEndian.PutUint16(canvas[0:lenPointer], uint16(conditions.GetNewestRecordIndex()))

		// Paint each record over canvas
		offset := lenPointer
		for _, record := range conditions.GetRecords() {
			if offset+recordSize > len(canvas) {
				// Canvas too small for all records, fall back to building from scratch
				goto fallback
			}
			recordBytes, err := dd.AppendSpecificConditionRecord(nil, record)
			if err != nil {
				return nil, fmt.Errorf("failed to append specific condition record: %w", err)
			}
			if len(recordBytes) != recordSize {
				return nil, fmt.Errorf("invalid specific condition record size: got %d, want %d", len(recordBytes), recordSize)
			}
			copy(canvas[offset:offset+recordSize], recordBytes)
			offset += recordSize
		}

		// Canvas preserves trailing bytes automatically
		return append(data, canvas...), nil
	}

fallback:
	// Fall back to building from scratch
	pointerBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(pointerBytes, uint16(conditions.GetNewestRecordIndex()))
	data = append(data, pointerBytes...)

	for _, record := range conditions.GetRecords() {
		var err error
		data, err = dd.AppendSpecificConditionRecord(data, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append specific condition record: %w", err)
		}
	}

	return data, nil
}
