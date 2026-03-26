package dd

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuTimeAdjustmentRecord parses a VuTimeAdjustmentRecord (Generation 1).
//
// The data type `VuTimeAdjustmentRecord` is specified in the Data Dictionary, Section 2.232.
//
// ASN.1 Specification (Gen1):
//
//	VuTimeAdjustmentRecord ::= SEQUENCE {
//	    oldTimeValue            TimeReal,           -- 4 bytes
//	    newTimeValue            TimeReal,           -- 4 bytes
//	    workshopName            Name,               -- 36 bytes (1 + 35)
//	    workshopAddress         Address,            -- 36 bytes (1 + 35)
//	    workshopCardNumber      FullCardNumber      -- 18 bytes
//	}
func (opts UnmarshalOptions) UnmarshalVuTimeAdjustmentRecord(data []byte) (*ddv1.VuTimeAdjustmentRecord, error) {
	const (
		idxOldTimeValue           = 0
		lenOldTimeValue           = 4
		idxNewTimeValue           = 4
		lenNewTimeValue           = 4
		idxWorkshopName           = 8
		lenWorkshopName           = 36
		idxWorkshopAddress        = 44
		lenWorkshopAddress        = 36
		idxWorkshopCardNumber     = 80
		lenWorkshopCardNumber     = 18
		lenVuTimeAdjustmentRecord = 98
	)

	if len(data) != lenVuTimeAdjustmentRecord {
		return nil, fmt.Errorf("invalid length for VuTimeAdjustmentRecord: got %d, want %d", len(data), lenVuTimeAdjustmentRecord)
	}

	record := &ddv1.VuTimeAdjustmentRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// Parse oldTimeValue (4 bytes)
	oldTime, err := opts.UnmarshalTimeReal(data[idxOldTimeValue : idxOldTimeValue+lenOldTimeValue])
	if err != nil {
		return nil, fmt.Errorf("failed to parse old time value: %w", err)
	}
	record.SetOldTimeValue(oldTime)

	// Parse newTimeValue (4 bytes)
	newTime, err := opts.UnmarshalTimeReal(data[idxNewTimeValue : idxNewTimeValue+lenNewTimeValue])
	if err != nil {
		return nil, fmt.Errorf("failed to parse new time value: %w", err)
	}
	record.SetNewTimeValue(newTime)

	// Parse workshopName (36 bytes: 1 byte code page + 35 bytes name)
	workshopName, err := opts.UnmarshalStringValue(data[idxWorkshopName : idxWorkshopName+lenWorkshopName])
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop name: %w", err)
	}
	record.SetWorkshopName(workshopName)

	// Parse workshopAddress (36 bytes: 1 byte code page + 35 bytes address)
	workshopAddress, err := opts.UnmarshalStringValue(data[idxWorkshopAddress : idxWorkshopAddress+lenWorkshopAddress])
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop address: %w", err)
	}
	record.SetWorkshopAddress(workshopAddress)

	// Parse workshopCardNumber (18 bytes)
	workshopCardNumber, err := opts.UnmarshalFullCardNumber(data[idxWorkshopCardNumber : idxWorkshopCardNumber+lenWorkshopCardNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to parse workshop card number: %w", err)
	}
	record.SetWorkshopCardNumber(workshopCardNumber)

	return record, nil
}

// MarshalVuTimeAdjustmentRecord marshals a VuTimeAdjustmentRecord to binary format (Generation 1).
func (opts MarshalOptions) MarshalVuTimeAdjustmentRecord(record *ddv1.VuTimeAdjustmentRecord) ([]byte, error) {
	const lenVuTimeAdjustmentRecord = 98

	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	// Use raw data painting strategy if available
	var canvas [lenVuTimeAdjustmentRecord]byte
	if record.HasRawData() {
		if len(record.GetRawData()) != lenVuTimeAdjustmentRecord {
			return nil, fmt.Errorf(
				"invalid raw_data length for VuTimeAdjustmentRecord: got %d, want %d",
				len(record.GetRawData()), lenVuTimeAdjustmentRecord,
			)
		}
		copy(canvas[:], record.GetRawData())
	}

	// Paint semantic values over the canvas
	const (
		idxOldTimeValue       = 0
		lenOldTimeValue       = 4
		idxNewTimeValue       = 4
		lenNewTimeValue       = 4
		idxWorkshopName       = 8
		lenWorkshopName       = 36
		idxWorkshopAddress    = 44
		lenWorkshopAddress    = 36
		idxWorkshopCardNumber = 80
		lenWorkshopCardNumber = 18
	)

	// Marshal oldTimeValue (4 bytes)
	oldTime, err := opts.MarshalTimeReal(record.GetOldTimeValue())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal old time value: %w", err)
	}
	if record.GetOldTimeValue() != nil {
		copy(canvas[idxOldTimeValue:idxOldTimeValue+lenOldTimeValue], oldTime)
	}

	// Marshal newTimeValue (4 bytes)
	newTime, err := opts.MarshalTimeReal(record.GetNewTimeValue())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new time value: %w", err)
	}
	if record.GetNewTimeValue() != nil {
		copy(canvas[idxNewTimeValue:idxNewTimeValue+lenNewTimeValue], newTime)
	}

	// Marshal workshopName (36 bytes: 1 byte code page + 35 bytes name)
	workshopName, err := opts.MarshalStringValue(record.GetWorkshopName())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop name: %w", err)
	}
	if len(workshopName) != lenWorkshopName {
		return nil, fmt.Errorf("invalid workshop name length: got %d, want %d", len(workshopName), lenWorkshopName)
	}
	copy(canvas[idxWorkshopName:idxWorkshopName+lenWorkshopName], workshopName)

	// Marshal workshopAddress (36 bytes: 1 byte code page + 35 bytes address)
	workshopAddress, err := opts.MarshalStringValue(record.GetWorkshopAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop address: %w", err)
	}
	if len(workshopAddress) != lenWorkshopAddress {
		return nil, fmt.Errorf("invalid workshop address length: got %d, want %d", len(workshopAddress), lenWorkshopAddress)
	}
	copy(canvas[idxWorkshopAddress:idxWorkshopAddress+lenWorkshopAddress], workshopAddress)

	// Marshal workshopCardNumber (18 bytes)
	workshopCardNumber, err := opts.MarshalFullCardNumber(record.GetWorkshopCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workshop card number: %w", err)
	}
	if len(workshopCardNumber) != lenWorkshopCardNumber {
		return nil, fmt.Errorf("invalid workshop card number length: got %d, want %d", len(workshopCardNumber), lenWorkshopCardNumber)
	}
	copy(canvas[idxWorkshopCardNumber:idxWorkshopCardNumber+lenWorkshopCardNumber], workshopCardNumber)

	return canvas[:], nil
}

// AnonymizeVuTimeAdjustmentRecord anonymizes a VU time adjustment record.
func (opts AnonymizeOptions) AnonymizeVuTimeAdjustmentRecord(rec *ddv1.VuTimeAdjustmentRecord) *ddv1.VuTimeAdjustmentRecord {
	if rec == nil {
		return nil
	}

	result := proto.Clone(rec).(*ddv1.VuTimeAdjustmentRecord)

	// Anonymize timestamps
	result.SetOldTimeValue(opts.AnonymizeTimestamp(rec.GetOldTimeValue()))
	result.SetNewTimeValue(opts.AnonymizeTimestamp(rec.GetNewTimeValue()))

	// Anonymize workshop name (35 bytes)
	result.SetWorkshopName(NewStringValue(ddv1.Encoding_ISO_8859_1, 35, "TEST WORKSHOP"))

	// Anonymize workshop address (35 bytes)
	result.SetWorkshopAddress(NewStringValue(ddv1.Encoding_ISO_8859_1, 35, "TEST ADDRESS, 00000 TEST CITY"))

	// Anonymize workshop card number
	result.SetWorkshopCardNumber(opts.AnonymizeFullCardNumber(rec.GetWorkshopCardNumber()))

	// Clear raw_data
	result.ClearRawData()

	return result
}
