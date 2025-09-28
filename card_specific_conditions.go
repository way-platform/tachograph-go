package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalCardSpecificConditions unmarshals specific conditions data from a card EF.
//
// The data type `SpecificConditionRecord` is specified in the Data Dictionary, Section 2.19.
//
// ASN.1 Definition:
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime                TimeReal,
//	    specificConditionType    SpecificConditionType
//	}
func unmarshalCardSpecificConditions(data []byte) (*cardv1.SpecificConditions, error) {
	const (
		lenSpecificConditionRecord = 5 // SpecificConditionRecord size
	)

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
	recordSize := lenSpecificConditionRecord

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

// parseSpecificConditionRecord parses a single specific condition record
func parseSpecificConditionRecord(r *bytes.Reader) (*cardv1.SpecificConditions_Record, error) {
	const (
		lenSpecificConditionRecord = 5
	)

	if r.Len() < lenSpecificConditionRecord {
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
	setEnumFromProtocolValue(ddv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED.Descriptor(), int32(conditionType), func(st protoreflect.EnumNumber) {
		record.SetSpecificConditionType(ddv1.SpecificConditionType(st))
	}, nil)

	return record, nil
}

// AppendCardSpecificConditions appends specific conditions data to a byte slice.
//
// The data type `SpecificConditionRecord` is specified in the Data Dictionary, Section 2.19.
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

	records := conditions.GetRecords()
	for _, record := range records {
		if record == nil {
			continue
		}

		// Entry time (4 bytes)
		data = appendTimeReal(data, record.GetEntryTime())

		// Specific condition type (1 byte) - convert enum to protocol value
		conditionTypeProtocol := getProtocolValueFromEnum(record.GetSpecificConditionType(), 0)
		data = append(data, byte(conditionTypeProtocol))
	}

	return data, nil
}
