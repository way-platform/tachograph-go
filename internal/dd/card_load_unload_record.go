package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalCardLoadUnloadRecord parses a CardLoadUnloadRecord (20 bytes).
//
// The data type `CardLoadUnloadRecord` is specified in the Data Dictionary, Section 2.24d.
//
// ASN.1 Definition (Gen2 V2):
//
//	CardLoadUnloadRecord ::= SEQUENCE {
//	    timeStamp                       TimeReal,
//	    operationType                   OperationType,
//	    gnssPlaceAuthRecord             GNSSPlaceAuthRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 20 bytes):
//   - Bytes 0-3: timeStamp (TimeReal)
//   - Byte 4: operationType (OperationType)
//   - Bytes 5-16: gnssPlaceAuthRecord (GNSSPlaceAuthRecord)
//   - Bytes 17-19: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalCardLoadUnloadRecord(data []byte) (*ddv1.CardLoadUnloadRecord, error) {
	const (
		idxTimeStamp            = 0
		idxOperationType        = 4
		idxGnssPlaceAuthRecord  = 5
		idxVehicleOdometerValue = 17
		lenCardLoadUnloadRecord = 20

		lenTimeReal            = 4
		lenOperationType       = 1
		lenGNSSPlaceAuthRecord = 12
		lenOdometerShort       = 3
	)

	if len(data) != lenCardLoadUnloadRecord {
		return nil, fmt.Errorf("invalid data length for CardLoadUnloadRecord: got %d, want %d", len(data), lenCardLoadUnloadRecord)
	}

	record := &ddv1.CardLoadUnloadRecord{}
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
	if vehicleOdometerValue != nil {
		record.SetVehicleOdometerKm(*vehicleOdometerValue)
	}

	return record, nil
}

// MarshalCardLoadUnloadRecord marshals a CardLoadUnloadRecord (20 bytes) to bytes.
func (opts MarshalOptions) MarshalCardLoadUnloadRecord(record *ddv1.CardLoadUnloadRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenCardLoadUnloadRecord = 20

	// Use raw data painting strategy if available
	var canvas [lenCardLoadUnloadRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenCardLoadUnloadRecord {
			return nil, fmt.Errorf("invalid raw_data length for CardLoadUnloadRecord: got %d, want %d", len(rawData), lenCardLoadUnloadRecord)
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
