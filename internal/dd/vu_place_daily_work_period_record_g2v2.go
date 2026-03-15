package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuPlaceDailyWorkPeriodRecordG2V2 parses a Generation 2 version 2 VuPlaceDailyWorkPeriodRecord (41 bytes).
//
// The data type `VuPlaceDailyWorkPeriodRecord` is specified in the Data Dictionary, Section 2.219.
//
// ASN.1 Definition (Gen2 V2):
//
//	VuPlaceDailyWorkPeriodRecord ::= SEQUENCE {
//	    fullCardNumberAndGeneration FullCardNumberAndGeneration,
//	    placeAuthRecord             PlaceAuthRecord
//	}
//
// Binary Layout (fixed length, 41 bytes):
//   - Bytes 0-18: fullCardNumberAndGeneration (FullCardNumberAndGeneration)
//   - Bytes 19-40: placeAuthRecord (PlaceAuthRecord)
func (opts UnmarshalOptions) UnmarshalVuPlaceDailyWorkPeriodRecordG2V2(data []byte) (*ddv1.VuPlaceDailyWorkPeriodRecordG2V2, error) {
	const (
		idxFullCardNumber                   = 0
		idxPlaceAuthRecord                  = 19
		lenVuPlaceDailyWorkPeriodRecordG2V2 = 41

		lenFullCardNumberAndGeneration = 19
		lenPlaceAuthRecord             = 22
	)

	if len(data) != lenVuPlaceDailyWorkPeriodRecordG2V2 {
		return nil, fmt.Errorf(
			"invalid data length for VuPlaceDailyWorkPeriodRecordG2V2: got %d, want %d",
			len(data), lenVuPlaceDailyWorkPeriodRecordG2V2,
		)
	}

	record := &ddv1.VuPlaceDailyWorkPeriodRecordG2V2{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// fullCardNumberAndGeneration (19 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxFullCardNumber : idxFullCardNumber+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number and generation: %w", err)
	}
	record.SetFullCardNumber(fullCardNumber)

	// placeAuthRecord (22 bytes)
	placeAuthRecord, err := opts.UnmarshalPlaceAuthRecord(data[idxPlaceAuthRecord : idxPlaceAuthRecord+lenPlaceAuthRecord])
	if err != nil {
		return nil, fmt.Errorf("unmarshal place auth record: %w", err)
	}
	record.SetPlaceAuthRecord(placeAuthRecord)

	return record, nil
}

// MarshalVuPlaceDailyWorkPeriodRecordG2V2 marshals a VuPlaceDailyWorkPeriodRecordG2V2 (41 bytes) to bytes.
func (opts MarshalOptions) MarshalVuPlaceDailyWorkPeriodRecordG2V2(record *ddv1.VuPlaceDailyWorkPeriodRecordG2V2) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuPlaceDailyWorkPeriodRecordG2V2 = 41

	// Use raw data painting strategy if available
	var canvas [lenVuPlaceDailyWorkPeriodRecordG2V2]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenVuPlaceDailyWorkPeriodRecordG2V2 {
			return nil, fmt.Errorf(
				"invalid raw_data length for VuPlaceDailyWorkPeriodRecordG2V2: got %d, want %d",
				len(rawData), lenVuPlaceDailyWorkPeriodRecordG2V2,
			)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// fullCardNumberAndGeneration (19 bytes)
	fullCardNumberBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetFullCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal full card number and generation: %w", err)
	}
	copy(canvas[offset:offset+19], fullCardNumberBytes)
	offset += 19

	// placeAuthRecord (22 bytes)
	placeAuthRecordBytes, err := opts.MarshalPlaceAuthRecord(record.GetPlaceAuthRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal place auth record: %w", err)
	}
	copy(canvas[offset:offset+22], placeAuthRecordBytes)

	return canvas[:], nil
}
