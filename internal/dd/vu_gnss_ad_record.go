package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuGNSSADRecord parses a VuGNSSADRecord (Generation 2, version 1 - 56 bytes).
//
// The data type `VuGNSSADRecord` is specified in the Data Dictionary, Section 2.203.
//
// ASN.1 Definition (Gen2 V1):
//
//	VuGNSSADRecord ::= SEQUENCE {
//	    timeStamp                       TimeReal,
//	    cardNumberAndGenDriverSlot      FullCardNumberAndGeneration,
//	    cardNumberAndGenCodriverSlot    FullCardNumberAndGeneration,
//	    gnssPlaceRecord                 GNSSPlaceRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 56 bytes):
//   - Bytes 0-3: timeStamp (TimeReal)
//   - Bytes 4-22: cardNumberAndGenDriverSlot (FullCardNumberAndGeneration)
//   - Bytes 23-41: cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration)
//   - Bytes 42-52: gnssPlaceRecord (GNSSPlaceRecord)
//   - Bytes 53-55: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalVuGNSSADRecord(data []byte) (*ddv1.VuGNSSADRecord, error) {
	const (
		idxTimeStamp              = 0
		idxCardNumberDriverSlot   = 4
		idxCardNumberCodriverSlot = 23
		idxGnssPlaceRecord        = 42
		idxVehicleOdometerValue   = 53
		lenVuGNSSADRecord         = 56

		lenTimeReal                    = 4
		lenFullCardNumberAndGeneration = 19
		lenGNSSPlaceRecord             = 11
		lenOdometerShort               = 3
	)

	if len(data) != lenVuGNSSADRecord {
		return nil, fmt.Errorf("invalid data length for VuGNSSADRecord: got %d, want %d", len(data), lenVuGNSSADRecord)
	}

	record := &ddv1.VuGNSSADRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// timeStamp (4 bytes)
	timeStamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal time stamp: %w", err)
	}
	record.SetTimeStamp(timeStamp)

	// cardNumberAndGenDriverSlot (19 bytes)
	cardNumberDriverSlot, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxCardNumberDriverSlot : idxCardNumberDriverSlot+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card number driver slot: %w", err)
	}
	record.SetCardNumberDriverSlot(cardNumberDriverSlot)

	// cardNumberAndGenCodriverSlot (19 bytes)
	cardNumberCodriverSlot, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxCardNumberCodriverSlot : idxCardNumberCodriverSlot+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card number codriver slot: %w", err)
	}
	record.SetCardNumberCodriverSlot(cardNumberCodriverSlot)

	// gnssPlaceRecord (11 bytes)
	gnssPlaceRecord, err := opts.UnmarshalGNSSPlaceRecord(data[idxGnssPlaceRecord : idxGnssPlaceRecord+lenGNSSPlaceRecord])
	if err != nil {
		return nil, fmt.Errorf("unmarshal GNSS place record: %w", err)
	}
	record.SetGnssPlaceRecord(gnssPlaceRecord)

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerValue, err := opts.UnmarshalOdometer(data[idxVehicleOdometerValue : idxVehicleOdometerValue+lenOdometerShort])
	if err != nil {
		return nil, fmt.Errorf("unmarshal vehicle odometer value: %w", err)
	}
	record.SetVehicleOdometerKm(int32(vehicleOdometerValue))

	return record, nil
}

// MarshalVuGNSSADRecord marshals a VuGNSSADRecord (56 bytes) to bytes.
func (opts MarshalOptions) MarshalVuGNSSADRecord(record *ddv1.VuGNSSADRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuGNSSADRecord = 56

	// Use raw data painting strategy if available
	var canvas [lenVuGNSSADRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenVuGNSSADRecord {
			return nil, fmt.Errorf("invalid raw_data length for VuGNSSADRecord: got %d, want %d", len(rawData), lenVuGNSSADRecord)
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

	// cardNumberAndGenDriverSlot (19 bytes)
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

	// gnssPlaceRecord (11 bytes)
	gnssPlaceRecordBytes, err := opts.MarshalGNSSPlaceRecord(record.GetGnssPlaceRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GNSS place record: %w", err)
	}
	copy(canvas[offset:offset+11], gnssPlaceRecordBytes)
	offset += 11

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerBytes, err := opts.MarshalOdometer(record.GetVehicleOdometerKm())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle odometer value: %w", err)
	}
	copy(canvas[offset:offset+3], vehicleOdometerBytes)

	return canvas[:], nil
}
