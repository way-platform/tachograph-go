package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuBorderCrossingRecord parses a VuBorderCrossingRecord (57 bytes).
//
// The data type `VuBorderCrossingRecord` is specified in the Data Dictionary, Section 2.203a.
//
// ASN.1 Definition (Gen2 V2):
//
//	VuBorderCrossingRecord ::= SEQUENCE {
//	    cardNumberAndGenDriverSlot      FullCardNumberAndGeneration,
//	    cardNumberAndGenCodriverSlot    FullCardNumberAndGeneration,
//	    countryLeft                     NationNumeric,
//	    countryEntered                  NationNumeric,
//	    gnssPlaceAuthRecord             GNSSPlaceAuthRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 55 bytes):
//   - Bytes 0-18: cardNumberAndGenDriverSlot (FullCardNumberAndGeneration, 19 bytes)
//   - Bytes 19-37: cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration, 19 bytes)
//   - Byte 38: countryLeft (NationNumeric)
//   - Byte 39: countryEntered (NationNumeric)
//   - Bytes 40-51: gnssPlaceAuthRecord (GNSSPlaceAuthRecord)
//   - Bytes 52-54: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalVuBorderCrossingRecord(data []byte) (*ddv1.VuBorderCrossingRecord, error) {
	const (
		idxCardNumberDriverSlot   = 0
		idxCardNumberCodriverSlot = 19
		idxCountryLeft            = 38
		idxCountryEntered         = 39
		idxGnssPlaceAuthRecord    = 40
		idxVehicleOdometerValue   = 52
		lenVuBorderCrossingRecord = 55

		lenFullCardNumberAndGeneration = 19
		lenNationNumeric               = 1
		lenGNSSPlaceAuthRecord         = 12
		lenOdometerShort               = 3
	)

	if len(data) != lenVuBorderCrossingRecord {
		return nil, fmt.Errorf("invalid data length for VuBorderCrossingRecord: got %d, want %d", len(data), lenVuBorderCrossingRecord)
	}

	record := &ddv1.VuBorderCrossingRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

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
	record.SetVehicleOdometerKm(int32(vehicleOdometerValue))

	return record, nil
}

// MarshalVuBorderCrossingRecord marshals a VuBorderCrossingRecord (57 bytes) to bytes.
func (opts MarshalOptions) MarshalVuBorderCrossingRecord(record *ddv1.VuBorderCrossingRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuBorderCrossingRecord = 55

	// Use raw data painting strategy if available
	var canvas [lenVuBorderCrossingRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenVuBorderCrossingRecord {
			return nil, fmt.Errorf("invalid raw_data length for VuBorderCrossingRecord: got %d, want %d", len(rawData), lenVuBorderCrossingRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

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
