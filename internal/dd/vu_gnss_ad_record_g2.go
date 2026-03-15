package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuGNSSADRecordG2 parses a VuGNSSADRecord (Generation 2, version 2 - 57 bytes).
//
// The data type `VuGNSSADRecord` is specified in the Data Dictionary, Section 2.203.
//
// ASN.1 Definition (Gen2 V2):
//
//	VuGNSSADRecord ::= SEQUENCE {
//	    timeStamp                       TimeReal,
//	    cardNumberAndGenDriverSlot      FullCardNumberAndGeneration,
//	    cardNumberAndGenCodriverSlot    FullCardNumberAndGeneration,
//	    gnssPlaceAuthRecord             GNSSPlaceAuthRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 57 bytes):
//   - Bytes 0-3: timeStamp (TimeReal)
//   - Bytes 4-22: cardNumberAndGenDriverSlot (FullCardNumberAndGeneration, 19 bytes)
//   - Bytes 23-41: cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration, 19 bytes)
//   - Bytes 42-53: gnssPlaceAuthRecord (GNSSPlaceAuthRecord)
//   - Bytes 54-56: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalVuGNSSADRecordG2(data []byte) (*ddv1.VuGNSSADRecordG2, error) {
	const (
		idxTimeStamp              = 0
		idxCardNumberDriverSlot   = 4
		idxCardNumberCodriverSlot = 23
		idxGnssPlaceAuthRecord    = 42
		idxVehicleOdometerValue   = 54
		lenVuGNSSADRecordG2       = 57

		lenTimeReal                    = 4
		lenFullCardNumberAndGeneration = 19
		lenGNSSPlaceAuthRecord         = 12
		lenOdometerShort               = 3
	)

	if len(data) != lenVuGNSSADRecordG2 {
		return nil, fmt.Errorf("invalid data length for VuGNSSADRecordG2: got %d, want %d", len(data), lenVuGNSSADRecordG2)
	}

	record := &ddv1.VuGNSSADRecordG2{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// timeStamp (4 bytes)
	timeStamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal time stamp: %w", err)
	}
	record.SetTimeStamp(timeStamp)

	// cardNumberAndGenDriverSlot (20 bytes)
	cardNumberDriverSlot, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxCardNumberDriverSlot : idxCardNumberDriverSlot+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card number driver slot: %w", err)
	}
	record.SetCardNumberDriverSlot(cardNumberDriverSlot)

	// cardNumberAndGenCodriverSlot (20 bytes)
	cardNumberCodriverSlot, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxCardNumberCodriverSlot : idxCardNumberCodriverSlot+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card number codriver slot: %w", err)
	}
	record.SetCardNumberCodriverSlot(cardNumberCodriverSlot)

	// gnssPlaceAuthRecord (12 bytes)
	gnssPlaceAuthRecord, err := opts.UnmarshalGNSSPlaceAuthRecord(data[idxGnssPlaceAuthRecord : idxGnssPlaceAuthRecord+lenGNSSPlaceAuthRecord])
	if err != nil {
		return nil, fmt.Errorf("unmarshal GNSS place auth record: %w", err)
	}
	record.SetGnssPlaceAuthRecord(gnssPlaceAuthRecord)

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerValue, err := opts.UnmarshalOdometer(data[idxVehicleOdometerValue : idxVehicleOdometerValue+lenOdometerShort])
	if err != nil {
		return nil, fmt.Errorf("unmarshal vehicle odometer value: %w", err)
	}
	record.SetVehicleOdometerKm(int32(vehicleOdometerValue))

	return record, nil
}

// MarshalVuGNSSADRecordG2 marshals a VuGNSSADRecordG2 (59 bytes) to bytes.
func (opts MarshalOptions) MarshalVuGNSSADRecordG2(record *ddv1.VuGNSSADRecordG2) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuGNSSADRecordG2 = 57

	// Use raw data painting strategy if available
	var canvas [lenVuGNSSADRecordG2]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenVuGNSSADRecordG2 {
			return nil, fmt.Errorf("invalid raw_data length for VuGNSSADRecordG2: got %d, want %d", len(rawData), lenVuGNSSADRecordG2)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// timeStamp (4 bytes)
	timeStampBytes, err := opts.MarshalTimeReal(record.GetTimeStamp())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal time stamp: %w", err)
	}
	copy(canvas[offset:offset+4], timeStampBytes)
	offset += 4

	// cardNumberAndGenDriverSlot (20 bytes)
	cardNumberDriverSlotBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetCardNumberDriverSlot())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card number driver slot: %w", err)
	}
	copy(canvas[offset:offset+19], cardNumberDriverSlotBytes)
	offset += 19

	// cardNumberAndGenCodriverSlot (19 bytes)
	cardNumberCodriverSlotBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetCardNumberCodriverSlot())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card number codriver slot: %w", err)
	}
	copy(canvas[offset:offset+19], cardNumberCodriverSlotBytes)
	offset += 19

	// gnssPlaceAuthRecord (12 bytes)
	gnssPlaceAuthRecordBytes, err := opts.MarshalGNSSPlaceAuthRecord(record.GetGnssPlaceAuthRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GNSS place auth record: %w", err)
	}
	copy(canvas[offset:offset+12], gnssPlaceAuthRecordBytes)
	offset += 12

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerBytes, err := opts.MarshalOdometer(record.GetVehicleOdometerKm())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle odometer value: %w", err)
	}
	copy(canvas[offset:offset+3], vehicleOdometerBytes)

	return canvas[:], nil
}
