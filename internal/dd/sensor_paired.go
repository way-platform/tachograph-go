package dd

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalSensorPaired parses the SensorPaired structure.
//
// See Data Dictionary, Section 2.144, `SensorPaired`.
//
// ASN.1 Specification (Generation 1):
//
//	SensorPaired ::= SEQUENCE {
//	    sensorSerialNumber SensorSerialNumber,    -- 8 bytes
//	    sensorApprovalNumber SensorApprovalNumber,-- 8 bytes (Gen1), 16 bytes (Gen2)
//	    sensorPairingDate SensorPairingDate       -- 4 bytes
//	}
//
// Binary Layout:
//   - Generation 1: 20 bytes total
//   - Generation 2: 28 bytes total (approval number is 16 bytes)
func (opts UnmarshalOptions) UnmarshalSensorPaired(data []byte) (*ddv1.SensorPaired, error) {
	// Determine generation by data length
	const (
		lenGen1 = 20
		lenGen2 = 28
	)

	if len(data) != lenGen1 && len(data) != lenGen2 {
		return nil, fmt.Errorf(
			"invalid data length for SensorPaired: got %d, want %d (Gen1) or %d (Gen2)",
			len(data), lenGen1, lenGen2,
		)
	}

	sensorPaired := &ddv1.SensorPaired{}
	if opts.PreserveRawData {
		sensorPaired.SetRawData(data)
	}

	const (
		idxSerialNumber   = 0
		lenSerialNumber   = 8
		idxApprovalNumber = 8
	)

	// Parse sensor serial number (8 bytes)
	serialNumber, err := opts.UnmarshalExtendedSerialNumber(data[idxSerialNumber : idxSerialNumber+lenSerialNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to parse sensor serial number: %w", err)
	}
	sensorPaired.SetSerialNumber(serialNumber)

	// Determine approval number length based on total data length
	var lenApprovalNumber int
	if len(data) == lenGen1 {
		lenApprovalNumber = 8 // Gen1
	} else {
		lenApprovalNumber = 16 // Gen2
	}

	// Parse sensor approval number (8 or 16 bytes depending on generation)
	approvalNumber, err := opts.UnmarshalIa5StringValue(
		data[idxApprovalNumber : idxApprovalNumber+lenApprovalNumber],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sensor approval number: %w", err)
	}
	sensorPaired.SetApprovalNumber(approvalNumber)

	// Parse sensor pairing date (4 bytes)
	idxPairingDate := idxApprovalNumber + lenApprovalNumber
	const lenPairingDate = 4
	pairingDate, err := opts.UnmarshalTimeReal(data[idxPairingDate : idxPairingDate+lenPairingDate])
	if err != nil {
		return nil, fmt.Errorf("failed to parse sensor pairing date: %w", err)
	}
	sensorPaired.SetPairingDate(pairingDate)

	return sensorPaired, nil
}

// MarshalSensorPaired marshals the SensorPaired structure using raw data painting.
//
// See Data Dictionary, Section 2.144, `SensorPaired`.
func (opts MarshalOptions) MarshalSensorPaired(sensorPaired *ddv1.SensorPaired) ([]byte, error) {
	if sensorPaired == nil {
		return nil, fmt.Errorf("sensorPaired cannot be nil")
	}

	// Determine size based on approval number length
	approvalNumberLen := int(sensorPaired.GetApprovalNumber().GetLength())
	var size int
	switch approvalNumberLen {
	case 8:
		size = 20 // Gen1
	case 16:
		size = 28 // Gen2
	default:
		return nil, fmt.Errorf(
			"invalid approval number length: got %d, want 8 (Gen1) or 16 (Gen2)",
			approvalNumberLen,
		)
	}

	// Use raw data painting strategy
	canvas := make([]byte, size)
	if raw := sensorPaired.GetRawData(); len(raw) > 0 {
		if len(raw) != size {
			return nil, fmt.Errorf(
				"invalid raw_data length for SensorPaired: got %d, want %d",
				len(raw), size,
			)
		}
		copy(canvas, raw)
	}

	offset := 0

	// Marshal sensor serial number (8 bytes)
	serialNumberBytes, err := opts.MarshalExtendedSerialNumber(sensorPaired.GetSerialNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sensor serial number: %w", err)
	}
	if len(serialNumberBytes) != 8 {
		return nil, fmt.Errorf(
			"invalid serial number length: got %d, want 8",
			len(serialNumberBytes),
		)
	}
	copy(canvas[offset:offset+8], serialNumberBytes)
	offset += 8

	// Marshal sensor approval number (8 or 16 bytes)
	approvalNumberBytes, err := opts.MarshalIa5StringValue(sensorPaired.GetApprovalNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sensor approval number: %w", err)
	}
	if len(approvalNumberBytes) != approvalNumberLen {
		return nil, fmt.Errorf(
			"invalid approval number length: got %d, want %d",
			len(approvalNumberBytes), approvalNumberLen,
		)
	}
	copy(canvas[offset:offset+approvalNumberLen], approvalNumberBytes)
	offset += approvalNumberLen

	// Marshal sensor pairing date (4 bytes)
	pairingDateBytes, err := opts.MarshalTimeReal(sensorPaired.GetPairingDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sensor pairing date: %w", err)
	}
	if len(pairingDateBytes) != 4 {
		return nil, fmt.Errorf(
			"invalid pairing date length: got %d, want 4",
			len(pairingDateBytes),
		)
	}
	copy(canvas[offset:offset+4], pairingDateBytes)
	offset += 4

	if offset != size {
		return nil, fmt.Errorf(
			"SensorPaired marshalling size mismatch: wrote %d bytes, expected %d",
			offset, size,
		)
	}

	return canvas, nil
}

// AnonymizeSensorPaired anonymizes sensor paired data.
func (opts AnonymizeOptions) AnonymizeSensorPaired(sensor *ddv1.SensorPaired) *ddv1.SensorPaired {
	if sensor == nil {
		return nil
	}

	result := proto.Clone(sensor).(*ddv1.SensorPaired)

	// Anonymize sensor serial number (ExtendedSerialNumber)
	// Preserve equipment type and manufacturer code (structural info) but anonymize the serial number
	serialNum := &ddv1.ExtendedSerialNumber{}
	if origSerial := sensor.GetSerialNumber(); origSerial != nil {
		serialNum.SetType(origSerial.GetType())
		serialNum.SetManufacturerCode(origSerial.GetManufacturerCode())
	}
	serialNum.SetSerialNumber(0)
	result.SetSerialNumber(serialNum)

	// Anonymize approval number (8 bytes IA5String for Gen1)
	result.SetApprovalNumber(NewIa5StringValue(8, "SENSOR01"))

	// Keep pairing date as-is (could be anonymized if needed)

	// Clear raw_data
	result.ClearRawData()

	return result
}
