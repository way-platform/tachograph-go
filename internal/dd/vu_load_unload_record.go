package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuLoadUnloadRecord parses a VuLoadUnloadRecord (60 bytes).
//
// The data type `VuLoadUnloadRecord` is specified in the Data Dictionary, Section 2.208a.
//
// ASN.1 Definition (Gen2 V2):
//
//	VuLoadUnloadRecord ::= SEQUENCE {
//	    timeStamp                       TimeReal,
//	    operationType                   OperationType,
//	    cardNumberAndGenDriverSlot      FullCardNumberAndGeneration,
//	    cardNumberAndGenCodriverSlot    FullCardNumberAndGeneration,
//	    gnssPlaceAuthRecord             GNSSPlaceAuthRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 58 bytes):
//   - Bytes 0-3: timeStamp (TimeReal)
//   - Byte 4: operationType (OperationType)
//   - Bytes 5-23: cardNumberAndGenDriverSlot (FullCardNumberAndGeneration, 19 bytes)
//   - Bytes 24-42: cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration, 19 bytes)
//   - Bytes 43-54: gnssPlaceAuthRecord (GNSSPlaceAuthRecord)
//   - Bytes 55-57: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalVuLoadUnloadRecord(data []byte) (*ddv1.VuLoadUnloadRecord, error) {
	const (
		idxTimeStamp              = 0
		idxOperationType          = 4
		idxCardNumberDriverSlot   = 5
		idxCardNumberCodriverSlot = 24
		idxGnssPlaceAuthRecord    = 43
		idxVehicleOdometerValue   = 55
		lenVuLoadUnloadRecord     = 58

		lenTimeReal                    = 4
		lenOperationType               = 1
		lenFullCardNumberAndGeneration = 19
		lenGNSSPlaceAuthRecord         = 12
		lenOdometerShort               = 3
	)

	if len(data) != lenVuLoadUnloadRecord {
		return nil, fmt.Errorf("invalid data length for VuLoadUnloadRecord: got %d, want %d", len(data), lenVuLoadUnloadRecord)
	}

	record := &ddv1.VuLoadUnloadRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// timeStamp (4 bytes)
	timeStamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal time stamp: %w", err)
	}
	record.SetTimeStamp(timeStamp)

	// operationType (1 byte)
	operationType, err := UnmarshalEnum[ddv1.OperationType](data[idxOperationType])
	if err != nil {
		return nil, fmt.Errorf("unmarshal operation type: %w", err)
	}
	record.SetOperationType(operationType)

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

// MarshalVuLoadUnloadRecord marshals a VuLoadUnloadRecord (60 bytes) to bytes.
func (opts MarshalOptions) MarshalVuLoadUnloadRecord(record *ddv1.VuLoadUnloadRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuLoadUnloadRecord = 58

	// Use raw data painting strategy if available
	var canvas [lenVuLoadUnloadRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenVuLoadUnloadRecord {
			return nil, fmt.Errorf("invalid raw_data length for VuLoadUnloadRecord: got %d, want %d", len(rawData), lenVuLoadUnloadRecord)
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

	// operationType (1 byte)
	operationTypeByte, _ := MarshalEnum(record.GetOperationType())
	canvas[offset] = operationTypeByte
	offset += 1

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
