package dd

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuFaultRecord parses a VuFaultRecord (Generation 1).
//
// The data type `VuFaultRecord` is specified in the Data Dictionary, Section 2.201.
//
// ASN.1 Specification (Gen1):
//
//	VuFaultRecord ::= SEQUENCE {
//	    faultType                       EventFaultType,               -- 1 byte
//	    faultRecordPurpose              EventFaultRecordPurpose,      -- 1 byte
//	    faultBeginTime                  TimeReal,                     -- 4 bytes
//	    faultEndTime                    TimeReal,                     -- 4 bytes
//	    cardNumberDriverSlotBegin       FullCardNumber,               -- 18 bytes
//	    cardNumberCodriverSlotBegin     FullCardNumber,               -- 18 bytes
//	    cardNumberDriverSlotEnd         FullCardNumber,               -- 18 bytes
//	    cardNumberCodriverSlotEnd       FullCardNumber                -- 18 bytes
//	}
func (opts UnmarshalOptions) UnmarshalVuFaultRecord(data []byte) (*ddv1.VuFaultRecord, error) {
	const (
		idxFaultType                   = 0
		idxFaultRecordPurpose          = 1
		idxFaultBeginTime              = 2
		lenFaultBeginTime              = 4
		idxFaultEndTime                = 6
		lenFaultEndTime                = 4
		idxCardNumberDriverSlotBegin   = 10
		lenCardNumber                  = 18
		idxCardNumberCodriverSlotBegin = 28
		idxCardNumberDriverSlotEnd     = 46
		idxCardNumberCodriverSlotEnd   = 64
		lenVuFaultRecord               = 82
	)

	if len(data) != lenVuFaultRecord {
		return nil, fmt.Errorf("invalid length for VuFaultRecord: got %d, want %d", len(data), lenVuFaultRecord)
	}

	record := &ddv1.VuFaultRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// Parse faultType (1 byte)
	faultType, unrecognizedFaultType := opts.parseEventFaultType(data[idxFaultType])
	if faultType != ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
		record.SetFaultType(faultType)
	}
	if unrecognizedFaultType != 0 {
		record.SetUnrecognizedFaultType(unrecognizedFaultType)
	}

	// Parse faultRecordPurpose (1 byte)
	recordPurpose, unrecognizedPurpose := opts.parseEventFaultRecordPurpose(data[idxFaultRecordPurpose])
	if recordPurpose != ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
		record.SetRecordPurpose(recordPurpose)
	}
	if unrecognizedPurpose != 0 {
		record.SetUnrecognizedRecordPurpose(unrecognizedPurpose)
	}

	// Parse faultBeginTime (4 bytes)
	beginTime, err := opts.UnmarshalTimeReal(data[idxFaultBeginTime : idxFaultBeginTime+lenFaultBeginTime])
	if err != nil {
		return nil, fmt.Errorf("failed to parse fault begin time: %w", err)
	}
	record.SetBeginTime(beginTime)

	// Parse faultEndTime (4 bytes)
	endTime, err := opts.UnmarshalTimeReal(data[idxFaultEndTime : idxFaultEndTime+lenFaultEndTime])
	if err != nil {
		return nil, fmt.Errorf("failed to parse fault end time: %w", err)
	}
	record.SetEndTime(endTime)

	// Parse cardNumberDriverSlotBegin (18 bytes)
	cardDriverBegin, err := opts.UnmarshalFullCardNumber(data[idxCardNumberDriverSlotBegin : idxCardNumberDriverSlotBegin+lenCardNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to parse driver card begin: %w", err)
	}
	record.SetCardNumberDriverSlotBegin(cardDriverBegin)

	// Parse cardNumberCodriverSlotBegin (18 bytes)
	cardCodriverBegin, err := opts.UnmarshalFullCardNumber(data[idxCardNumberCodriverSlotBegin : idxCardNumberCodriverSlotBegin+lenCardNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to parse codriver card begin: %w", err)
	}
	record.SetCardNumberCodriverSlotBegin(cardCodriverBegin)

	// Parse cardNumberDriverSlotEnd (18 bytes)
	cardDriverEnd, err := opts.UnmarshalFullCardNumber(data[idxCardNumberDriverSlotEnd : idxCardNumberDriverSlotEnd+lenCardNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to parse driver card end: %w", err)
	}
	record.SetCardNumberDriverSlotEnd(cardDriverEnd)

	// Parse cardNumberCodriverSlotEnd (18 bytes)
	cardCodriverEnd, err := opts.UnmarshalFullCardNumber(data[idxCardNumberCodriverSlotEnd : idxCardNumberCodriverSlotEnd+lenCardNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to parse codriver card end: %w", err)
	}
	record.SetCardNumberCodriverSlotEnd(cardCodriverEnd)

	return record, nil
}

// MarshalVuFaultRecord marshals a VuFaultRecord to binary format (Generation 1).
func (opts MarshalOptions) MarshalVuFaultRecord(record *ddv1.VuFaultRecord) ([]byte, error) {
	const lenVuFaultRecord = 82

	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	// Use raw data painting strategy if available
	var canvas [lenVuFaultRecord]byte
	if record.HasRawData() {
		if len(record.GetRawData()) != lenVuFaultRecord {
			return nil, fmt.Errorf(
				"invalid raw_data length for VuFaultRecord: got %d, want %d",
				len(record.GetRawData()), lenVuFaultRecord,
			)
		}
		copy(canvas[:], record.GetRawData())
	}

	// Paint semantic values over the canvas
	const (
		idxFaultType                   = 0
		idxFaultRecordPurpose          = 1
		idxFaultBeginTime              = 2
		lenFaultBeginTime              = 4
		idxFaultEndTime                = 6
		lenFaultEndTime                = 4
		idxCardNumberDriverSlotBegin   = 10
		lenCardNumber                  = 18
		idxCardNumberCodriverSlotBegin = 28
		idxCardNumberDriverSlotEnd     = 46
		idxCardNumberCodriverSlotEnd   = 64
	)

	// Marshal faultType (1 byte)
	canvas[idxFaultType] = opts.marshalEventFaultType(record.GetFaultType(), record.GetUnrecognizedFaultType())

	// Marshal faultRecordPurpose (1 byte)
	canvas[idxFaultRecordPurpose] = opts.marshalEventFaultRecordPurpose(record.GetRecordPurpose(), record.GetUnrecognizedRecordPurpose())

	// Marshal faultBeginTime (4 bytes)
	beginTime, err := opts.MarshalTimeReal(record.GetBeginTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal begin time: %w", err)
	}
	copy(canvas[idxFaultBeginTime:idxFaultBeginTime+lenFaultBeginTime], beginTime)

	// Marshal faultEndTime (4 bytes)
	endTime, err := opts.MarshalTimeReal(record.GetEndTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal end time: %w", err)
	}
	copy(canvas[idxFaultEndTime:idxFaultEndTime+lenFaultEndTime], endTime)

	// Marshal cardNumberDriverSlotBegin (18 bytes)
	cardDriverBegin, err := opts.MarshalFullCardNumber(record.GetCardNumberDriverSlotBegin())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal driver card begin: %w", err)
	}
	copy(canvas[idxCardNumberDriverSlotBegin:idxCardNumberDriverSlotBegin+lenCardNumber], cardDriverBegin)

	// Marshal cardNumberCodriverSlotBegin (18 bytes)
	cardCodriverBegin, err := opts.MarshalFullCardNumber(record.GetCardNumberCodriverSlotBegin())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal codriver card begin: %w", err)
	}
	copy(canvas[idxCardNumberCodriverSlotBegin:idxCardNumberCodriverSlotBegin+lenCardNumber], cardCodriverBegin)

	// Marshal cardNumberDriverSlotEnd (18 bytes)
	cardDriverEnd, err := opts.MarshalFullCardNumber(record.GetCardNumberDriverSlotEnd())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal driver card end: %w", err)
	}
	copy(canvas[idxCardNumberDriverSlotEnd:idxCardNumberDriverSlotEnd+lenCardNumber], cardDriverEnd)

	// Marshal cardNumberCodriverSlotEnd (18 bytes)
	cardCodriverEnd, err := opts.MarshalFullCardNumber(record.GetCardNumberCodriverSlotEnd())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal codriver card end: %w", err)
	}
	copy(canvas[idxCardNumberCodriverSlotEnd:idxCardNumberCodriverSlotEnd+lenCardNumber], cardCodriverEnd)

	return canvas[:], nil
}

// parseEventFaultType parses an EventFaultType byte value.
//
// The protocol value is not the enum number: the mapping is declared in the
// protocol_enum_value annotation of each enum value, as in the rest of the
// package. Casting the raw byte to the enum shifted every value by the two
// sentinels (UNSPECIFIED, UNRECOGNIZED), so protocol '07'H (over speeding)
// was reported as GENERAL_CARD_INSERTION_WHILE_DRIVING.
func (opts UnmarshalOptions) parseEventFaultType(b byte) (ddv1.EventFaultType, int32) {
	eventFaultType, err := UnmarshalEnum[ddv1.EventFaultType](b)
	if err != nil {
		return ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED, int32(b)
	}
	return eventFaultType, 0
}

// marshalEventFaultType marshals an EventFaultType to a byte.
func (opts MarshalOptions) marshalEventFaultType(eventFaultType ddv1.EventFaultType, unrecognized int32) byte {
	if unrecognized != 0 {
		return byte(unrecognized)
	}
	b, err := MarshalEnum(eventFaultType)
	if err != nil {
		return 0
	}
	return b
}

// parseEventFaultRecordPurpose parses an EventFaultRecordPurpose byte value.
//
// Same annotation-driven mapping as parseEventFaultType.
func (opts UnmarshalOptions) parseEventFaultRecordPurpose(b byte) (ddv1.EventFaultRecordPurpose, int32) {
	purpose, err := UnmarshalEnum[ddv1.EventFaultRecordPurpose](b)
	if err != nil {
		return ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED, int32(b)
	}
	return purpose, 0
}

// marshalEventFaultRecordPurpose marshals an EventFaultRecordPurpose to a byte.
func (opts MarshalOptions) marshalEventFaultRecordPurpose(purpose ddv1.EventFaultRecordPurpose, unrecognized int32) byte {
	if unrecognized != 0 {
		return byte(unrecognized)
	}
	b, err := MarshalEnum(purpose)
	if err != nil {
		return 0
	}
	return b
}

// AnonymizeVuFaultRecord anonymizes a VU fault record.
func (opts AnonymizeOptions) AnonymizeVuFaultRecord(rec *ddv1.VuFaultRecord) *ddv1.VuFaultRecord {
	if rec == nil {
		return nil
	}

	result := proto.Clone(rec).(*ddv1.VuFaultRecord)

	// Anonymize timestamps
	result.SetBeginTime(opts.AnonymizeTimestamp(rec.GetBeginTime()))
	result.SetEndTime(opts.AnonymizeTimestamp(rec.GetEndTime()))

	// Anonymize card numbers
	result.SetCardNumberDriverSlotBegin(opts.AnonymizeFullCardNumber(rec.GetCardNumberDriverSlotBegin()))
	result.SetCardNumberCodriverSlotBegin(opts.AnonymizeFullCardNumber(rec.GetCardNumberCodriverSlotBegin()))
	result.SetCardNumberDriverSlotEnd(opts.AnonymizeFullCardNumber(rec.GetCardNumberDriverSlotEnd()))
	result.SetCardNumberCodriverSlotEnd(opts.AnonymizeFullCardNumber(rec.GetCardNumberCodriverSlotEnd()))

	// Clear raw_data
	result.ClearRawData()

	return result
}
