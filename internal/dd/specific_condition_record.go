package dd

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalSpecificConditionRecord unmarshals a SpecificConditionRecord from binary data.
//
// The data type `SpecificConditionRecord` is specified in the Data Dictionary, Section 2.152.
//
// ASN.1 Definition:
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime TimeReal,
//	    specificConditionType SpecificConditionType
//	}
func (opts UnmarshalOptions) UnmarshalSpecificConditionRecord(data []byte) (*ddv1.SpecificConditionRecord, error) {
	const (
		lenSpecificConditionRecord = 5
		idxEntryTime               = 0
		idxSpecificConditionType   = 4
	)

	if len(data) != lenSpecificConditionRecord {
		return nil, fmt.Errorf("invalid data length for SpecificConditionRecord: got %d, want %d", len(data), lenSpecificConditionRecord)
	}

	record := &ddv1.SpecificConditionRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, err := opts.UnmarshalTimeReal(data[idxEntryTime : idxEntryTime+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// Parse specificConditionType (1 byte)
	SetEnumFromProtocolValueGeneric(
		ddv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED.Descriptor(),
		int32(data[idxSpecificConditionType]),
		func(enumNum protoreflect.EnumNumber) {
			record.SetSpecificConditionType(ddv1.SpecificConditionType(enumNum))
		},
		func(rawValue int32) {
			record.SetUnrecognizedSpecificConditionType(rawValue)
		},
	)

	return record, nil
}

// AppendSpecificConditionRecord appends a SpecificConditionRecord to dst.
//
// The data type `SpecificConditionRecord` is specified in the Data Dictionary, Section 2.152.
//
// ASN.1 Definition:
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime TimeReal,
//	    specificConditionType SpecificConditionType
//	}
func AppendSpecificConditionRecord(dst []byte, record *ddv1.SpecificConditionRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("specific condition record cannot be nil")
	}

	// Append entryTime (TimeReal - 4 bytes)
	entryTime := record.GetEntryTime()
	if entryTime == nil {
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	} else {
		var err error
		dst, err = AppendTimeReal(dst, entryTime)
		if err != nil {
			return nil, fmt.Errorf("failed to append entry time: %w", err)
		}
	}

	// Append specificConditionType (1 byte)
	conditionType, _ := GetProtocolValueForEnum(record.GetSpecificConditionType())
	dst = append(dst, byte(conditionType))

	return dst, nil
}
