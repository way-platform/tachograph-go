package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalCardBorderCrossingRecord parses a CardBorderCrossingRecord (17 bytes).
//
// The data type `CardBorderCrossingRecord` is specified in the Data Dictionary, Section 2.11b.
//
// ASN.1 Definition (Gen2 V2):
//
//	CardBorderCrossingRecord ::= SEQUENCE {
//	    countryLeft                     NationNumeric,
//	    countryEntered                  NationNumeric,
//	    gnssPlaceAuthRecord             GNSSPlaceAuthRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 17 bytes):
//   - Byte 0: countryLeft (NationNumeric)
//   - Byte 1: countryEntered (NationNumeric)
//   - Bytes 2-13: gnssPlaceAuthRecord (GNSSPlaceAuthRecord)
//   - Bytes 14-16: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalCardBorderCrossingRecord(data []byte) (*ddv1.CardBorderCrossingRecord, error) {
	const (
		idxCountryLeft              = 0
		idxCountryEntered           = 1
		idxGnssPlaceAuthRecord      = 2
		idxVehicleOdometerValue     = 14
		lenCardBorderCrossingRecord = 17

		lenNationNumeric       = 1
		lenGNSSPlaceAuthRecord = 12
		lenOdometerShort       = 3
	)

	if len(data) != lenCardBorderCrossingRecord {
		return nil, fmt.Errorf("invalid data length for CardBorderCrossingRecord: got %d, want %d", len(data), lenCardBorderCrossingRecord)
	}

	record := &ddv1.CardBorderCrossingRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// countryLeft (1 byte)
	countryLeft, err := UnmarshalEnum[ddv1.NationNumeric](data[idxCountryLeft])
	if err != nil {
		return nil, fmt.Errorf("unmarshal country left: %w", err)
	}
	record.SetCountryLeft(countryLeft)

	// countryEntered (1 byte)
	countryEntered, err := UnmarshalEnum[ddv1.NationNumeric](data[idxCountryEntered])
	if err != nil {
		return nil, fmt.Errorf("unmarshal country entered: %w", err)
	}
	record.SetCountryEntered(countryEntered)

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

// MarshalCardBorderCrossingRecord marshals a CardBorderCrossingRecord (17 bytes) to bytes.
func (opts MarshalOptions) MarshalCardBorderCrossingRecord(record *ddv1.CardBorderCrossingRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenCardBorderCrossingRecord = 17

	// Use raw data painting strategy if available
	var canvas [lenCardBorderCrossingRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenCardBorderCrossingRecord {
			return nil, fmt.Errorf("invalid raw_data length for CardBorderCrossingRecord: got %d, want %d", len(rawData), lenCardBorderCrossingRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// countryLeft (1 byte)
	countryLeftByte, _ := MarshalEnum(record.GetCountryLeft())
	canvas[offset] = countryLeftByte
	offset += 1

	// countryEntered (1 byte)
	countryEnteredByte, _ := MarshalEnum(record.GetCountryEntered())
	canvas[offset] = countryEnteredByte
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
