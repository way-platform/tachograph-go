package card

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalCardSpecificConditions unmarshals specific conditions data from a card EF.
//
// The data type `SpecificConditionRecord` is specified in the Data Dictionary, Section 2.152.
//
// ASN.1 Definition:
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime                TimeReal,
//	    specificConditionType    SpecificConditionType
//	}
func unmarshalCardSpecificConditions(data []byte) (*cardv1.SpecificConditions, error) {
	const lenSpecificConditionRecord = 5

	if len(data) == 0 {
		// Empty data is valid - no specific conditions
		var target cardv1.SpecificConditions
		target.SetRecords([]*ddv1.SpecificConditionRecord{})
		return &target, nil
	}

	var target cardv1.SpecificConditions
	var records []*ddv1.SpecificConditionRecord

	// Parse each 5-byte SpecificConditionRecord using the DD package
	offset := 0
	for offset+lenSpecificConditionRecord <= len(data) {
		record, err := dd.UnmarshalSpecificConditionRecord(data[offset : offset+lenSpecificConditionRecord])
		if err != nil {
			break // Stop parsing on error, but return what we have
		}
		records = append(records, record)
		offset += lenSpecificConditionRecord
	}

	target.SetRecords(records)
	return &target, nil
}

// AppendCardSpecificConditions appends specific conditions data to a byte slice.
//
// The data type `SpecificConditionRecord` is specified in the Data Dictionary, Section 2.152.
//
// ASN.1 Definition:
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime                TimeReal,
//	    specificConditionType    SpecificConditionType
//	}
func appendCardSpecificConditions(data []byte, conditions *cardv1.SpecificConditions) ([]byte, error) {
	if conditions == nil {
		return data, nil
	}

	// Write each specific condition record using the DD package
	for _, record := range conditions.GetRecords() {
		var err error
		data, err = dd.AppendSpecificConditionRecord(data, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append specific condition record: %w", err)
		}
	}

	return data, nil
}
