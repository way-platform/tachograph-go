package dd

import (
	"encoding/binary"
	"fmt"
	"strings"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalCardVehicleRecordG2 unmarshals a Generation 2 CardVehicleRecord (48 bytes).
//
// The data type `CardVehicleRecord` (Gen2 variant) is specified in the Data Dictionary, Section 2.37.
//
// ASN.1 Definition (Gen2):
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBegin          OdometerShort,
//	    vehicleOdometerEnd            OdometerShort,
//	    vehicleFirstUse               TimeReal,
//	    vehicleLastUse                TimeReal,
//	    vehicleRegistration           VehicleRegistrationIdentification,
//	    vuDataBlockCounter            VuDataBlockCounter,
//	    vehicleIdentificationNumber   VehicleIdentificationNumber
//	}
func (opts UnmarshalOptions) UnmarshalCardVehicleRecordG2(data []byte) (*ddv1.CardVehicleRecordG2, error) {
	const (
		idxOdometerBegin       = 0
		idxOdometerEnd         = 3
		idxVehicleFirstUse     = 6
		idxVehicleLastUse      = 10
		idxVehicleRegistration = 14
		idxVuDataBlockCounter  = 29
		idxVIN                 = 31
		lenCardVehicleRecord   = 48 // Fixed size for Gen2
	)

	if len(data) != lenCardVehicleRecord {
		return nil, fmt.Errorf("invalid data length for Gen2 CardVehicleRecord: got %d, want %d", len(data), lenCardVehicleRecord)
	}

	record := &ddv1.CardVehicleRecordG2{}
	record.SetRawData(data)

	// Parse odometer begin (3 bytes)
	odometerBeginBytes := data[idxOdometerBegin : idxOdometerBegin+3]
	odometerBegin := int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBeginBytes...)))
	record.SetVehicleOdometerBeginKm(odometerBegin)

	// Parse odometer end (3 bytes)
	odometerEndBytes := data[idxOdometerEnd : idxOdometerEnd+3]
	odometerEnd := int32(binary.BigEndian.Uint32(append([]byte{0}, odometerEndBytes...)))
	record.SetVehicleOdometerEndKm(odometerEnd)

	// Parse vehicle first use (TimeReal - 4 bytes)
	vehicleFirstUse, err := opts.UnmarshalTimeReal(data[idxVehicleFirstUse : idxVehicleFirstUse+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle first use: %w", err)
	}
	record.SetVehicleFirstUse(vehicleFirstUse)

	// Parse vehicle last use (TimeReal - 4 bytes)
	vehicleLastUse, err := opts.UnmarshalTimeReal(data[idxVehicleLastUse : idxVehicleLastUse+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle last use: %w", err)
	}
	record.SetVehicleLastUse(vehicleLastUse)

	// Parse vehicle registration (15 bytes)
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[idxVehicleRegistration : idxVehicleRegistration+15])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle registration: %w", err)
	}
	record.SetVehicleRegistration(vehicleReg)

	// Parse VU data block counter (2 bytes as BCD)
	vuDataBlockCounter, err := opts.UnmarshalBcdString(data[idxVuDataBlockCounter : idxVuDataBlockCounter+2])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal VU data block counter: %w", err)
	}
	record.SetVuDataBlockCounter(vuDataBlockCounter)

	// Parse VIN (17 bytes IA5String)
	vinBytes := data[idxVIN : idxVIN+17]
	vin := strings.TrimRight(string(vinBytes), "\x00 ") // Trim null bytes and spaces
	record.SetVehicleIdentificationNumber(vin)

	return record, nil
}

// AppendCardVehicleRecordG2 appends a Generation 2 CardVehicleRecord (48 bytes).
func AppendCardVehicleRecordG2(dst []byte, record *ddv1.CardVehicleRecordG2) ([]byte, error) {
	const lenCardVehicleRecord = 48 // Fixed size for Gen2

	// Use raw data painting strategy if available
	var canvas [lenCardVehicleRecord]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenCardVehicleRecord {
			return nil, fmt.Errorf("invalid raw_data length for CardVehicleRecordG2: got %d, want %d", len(rawData), lenCardVehicleRecord)
		}
		copy(canvas[:], rawData)
	}

	// Paint semantic values over the canvas
	var err error

	// Odometer begin (3 bytes)
	odometerBegin := uint32(record.GetVehicleOdometerBeginKm())
	canvas[0] = byte(odometerBegin >> 16)
	canvas[1] = byte(odometerBegin >> 8)
	canvas[2] = byte(odometerBegin)

	// Odometer end (3 bytes)
	odometerEnd := uint32(record.GetVehicleOdometerEndKm())
	canvas[3] = byte(odometerEnd >> 16)
	canvas[4] = byte(odometerEnd >> 8)
	canvas[5] = byte(odometerEnd)

	// Vehicle first use (4 bytes)
	firstUseBytes, err := AppendTimeReal(nil, record.GetVehicleFirstUse())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle first use: %w", err)
	}
	copy(canvas[6:10], firstUseBytes)

	// Vehicle last use (4 bytes)
	lastUseBytes, err := AppendTimeReal(nil, record.GetVehicleLastUse())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle last use: %w", err)
	}
	copy(canvas[10:14], lastUseBytes)

	// Vehicle registration (15 bytes)
	vehicleRegBytes, err := AppendVehicleRegistration(nil, record.GetVehicleRegistration())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle registration: %w", err)
	}
	copy(canvas[14:29], vehicleRegBytes)

	// VU data block counter (2 bytes as BCD)
	vuDataBlockCounterBytes, err := AppendBcdString(nil, record.GetVuDataBlockCounter())
	if err != nil {
		return nil, fmt.Errorf("failed to append VU data block counter: %w", err)
	}
	copy(canvas[29:31], vuDataBlockCounterBytes)

	// VIN (17 bytes IA5String)
	vin := record.GetVehicleIdentificationNumber()
	vinBytes := make([]byte, 17)
	copy(vinBytes, []byte(vin))
	// Pad with spaces if shorter than 17 bytes
	for i := len(vin); i < 17; i++ {
		vinBytes[i] = ' '
	}
	copy(canvas[31:48], vinBytes)

	return append(dst, canvas[:]...), nil
}
