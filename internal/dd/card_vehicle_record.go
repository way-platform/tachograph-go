package dd

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UnmarshalCardVehicleRecord unmarshals a Generation 1 CardVehicleRecord (31 bytes).
//
// The data type `CardVehicleRecord` is specified in the Data Dictionary, Section 2.37.
//
// ASN.1 Definition (Gen1):
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBegin       OdometerShort,
//	    vehicleOdometerEnd         OdometerShort,
//	    vehicleFirstUse            TimeReal,
//	    vehicleLastUse             TimeReal,
//	    vehicleRegistration        VehicleRegistrationIdentification,
//	    vuDataBlockCounter         VuDataBlockCounter
//	}
func (opts UnmarshalOptions) UnmarshalCardVehicleRecord(data []byte) (*ddv1.CardVehicleRecord, error) {
	const (
		idxOdometerBegin       = 0
		idxOdometerEnd         = 3
		idxVehicleFirstUse     = 6
		idxVehicleLastUse      = 10
		idxVehicleRegistration = 14
		idxVuDataBlockCounter  = 29
		lenCardVehicleRecord   = 31 // Fixed size for Gen1
	)

	if len(data) != lenCardVehicleRecord {
		return nil, fmt.Errorf("invalid data length for Gen1 CardVehicleRecord: got %d, want %d", len(data), lenCardVehicleRecord)
	}

	record := &ddv1.CardVehicleRecord{}
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

	return record, nil
}

// AppendCardVehicleRecord appends a Generation 1 CardVehicleRecord (31 bytes).
func AppendCardVehicleRecord(dst []byte, record *ddv1.CardVehicleRecord) ([]byte, error) {
	const lenCardVehicleRecord = 31 // Fixed size for Gen1

	// Use raw data painting strategy if available
	var canvas [lenCardVehicleRecord]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenCardVehicleRecord {
			return nil, fmt.Errorf("invalid raw_data length for CardVehicleRecord: got %d, want %d", len(rawData), lenCardVehicleRecord)
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

	return append(dst, canvas[:]...), nil
}

// AnonymizeCardVehicleRecord creates an anonymized copy, replacing sensitive data
// with static, deterministic test values while preserving structure.
//
// Parameters:
//   - record: The original record to anonymize
//   - index: Record index, used to create incrementing timestamps for determinism
//
// Anonymization strategy:
// - Vehicle registration: Replaced with "TEST-VRN"
// - Timestamps: Static base (2020-01-01 00:00:00) + 1 day per index
// - Odometer readings: Rounded to nearest 1000km
// - Country: Preserved (structural info)
// - VU counter: Preserved (structural info)
func AnonymizeCardVehicleRecord(record *ddv1.CardVehicleRecord, index int) *ddv1.CardVehicleRecord {
	if record == nil {
		return nil
	}

	anonymized := &ddv1.CardVehicleRecord{}

	// Round odometer readings to nearest 1000km
	odometerBegin := (record.GetVehicleOdometerBeginKm() / 1000) * 1000
	odometerEnd := (record.GetVehicleOdometerEndKm() / 1000) * 1000
	anonymized.SetVehicleOdometerBeginKm(odometerBegin)
	anonymized.SetVehicleOdometerEndKm(odometerEnd)

	// Static timestamps: 2020-01-01 00:00:00 UTC (epoch: 1577836800) + 1 day per record
	const secondsPerDay = 86400
	baseEpoch := int64(1577836800)
	firstUseEpoch := baseEpoch + int64(index)*secondsPerDay
	lastUseEpoch := firstUseEpoch + secondsPerDay - 1 // End of same day

	anonymized.SetVehicleFirstUse(&timestamppb.Timestamp{Seconds: firstUseEpoch})
	anonymized.SetVehicleLastUse(&timestamppb.Timestamp{Seconds: lastUseEpoch})

	// Anonymize vehicle registration
	if vreg := record.GetVehicleRegistration(); vreg != nil {
		anonymizedReg := &ddv1.VehicleRegistrationIdentification{}
		
		// Preserve country (structural info)
		anonymizedReg.SetNation(vreg.GetNation())
		
		// Replace VRN with test value
		// VehicleRegistrationNumber is: 1 byte code page + 13 bytes data
		testRegNum := &ddv1.StringValue{}
		testRegNum.SetValue("TEST-VRN")
		testRegNum.SetEncoding(ddv1.Encoding_ISO_8859_1) // Code page 1 (Latin-1)
		testRegNum.SetLength(13) // Length of data bytes (not including code page)
		anonymizedReg.SetNumber(testRegNum)
		
		anonymized.SetVehicleRegistration(anonymizedReg)
	}

	// Preserve VU counter (structural info)
	anonymized.SetVuDataBlockCounter(record.GetVuDataBlockCounter())

	// Regenerate raw_data for binary fidelity
	if rawData, err := AppendCardVehicleRecord(nil, anonymized); err == nil {
		anonymized.SetRawData(rawData)
	}

	return anonymized
}
