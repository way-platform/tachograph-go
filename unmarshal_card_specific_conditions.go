package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalCardSpecificConditions unmarshals specific conditions data from a card EF.
//
// TODO: Add ASN.1 Specification - SpecificConditionRecord definition not found in data dictionary
//
// Binary Layout (variable size):
//
//	Each record: 5 bytes (4 bytes time + 1 byte condition type)
//
// Constants:
const (
	// SpecificConditionRecord size
	specificConditionRecordSize = 5
)

func unmarshalCardSpecificConditions(data []byte) (*cardv1.SpecificConditions, error) {
	if len(data) == 0 {
		// Empty data is valid - no specific conditions
		var target cardv1.SpecificConditions
		target.SetRecords([]*cardv1.SpecificConditions_Record{})
		return &target, nil
	}

	var target cardv1.SpecificConditions
	r := bytes.NewReader(data)
	var records []*cardv1.SpecificConditions_Record

	// Each record is 5 bytes: 4 bytes time + 1 byte condition type
	recordSize := specificConditionRecordSize

	for r.Len() >= recordSize {
		record, err := parseSpecificConditionRecord(r)
		if err != nil {
			break // Stop parsing on error, but return what we have
		}
		records = append(records, record)
	}

	target.SetRecords(records)
	return &target, nil
}

// UnmarshalCardSpecificConditions unmarshals specific conditions data from a card EF (legacy function).
// Deprecated: Use unmarshalCardSpecificConditions instead.
func UnmarshalCardSpecificConditions(data []byte, target *cardv1.SpecificConditions) error {
	result, err := unmarshalCardSpecificConditions(data)
	if err != nil {
		return err
	}
	*target = *result
	return nil
}

// parseSpecificConditionRecord parses a single specific condition record
func parseSpecificConditionRecord(r *bytes.Reader) (*cardv1.SpecificConditions_Record, error) {
	if r.Len() < 5 {
		return nil, fmt.Errorf("insufficient data for specific condition record")
	}

	record := &cardv1.SpecificConditions_Record{}

	// Read entry time (4 bytes)
	record.SetEntryTime(readTimeReal(r))

	// Read specific condition type (1 byte)
	var conditionType byte
	if err := binary.Read(r, binary.BigEndian, &conditionType); err != nil {
		return nil, fmt.Errorf("failed to read condition type: %w", err)
	}
	// Convert raw condition type to enum using protocol annotations
	SetSpecificConditionType(datadictionaryv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED.Descriptor(), int32(conditionType), func(st protoreflect.EnumNumber) {
		record.SetSpecificConditionType(datadictionaryv1.SpecificConditionType(st))
	}, nil)

	return record, nil
}
