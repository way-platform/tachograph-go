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

	// Capture trailing bytes for round-trip fidelity
	if offset < len(recordsData) {
		target.SetTrailingBytes(recordsData[offset:])
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

	// Write newest record pointer (2 bytes)
	pointerBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(pointerBytes, uint16(conditions.GetNewestRecordIndex()))
	data = append(data, pointerBytes...)

	// Write each specific condition record using the DD package
	for _, record := range conditions.GetRecords() {
		var err error
		data, err = dd.AppendSpecificConditionRecord(data, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append specific condition record: %w", err)
		}
	}

	// Append trailing bytes for round-trip fidelity
	if trailingBytes := conditions.GetTrailingBytes(); len(trailingBytes) > 0 {
		data = append(data, trailingBytes...)
	}

	return data, nil
}
