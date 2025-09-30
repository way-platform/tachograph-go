package dd

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalExtendedSerialNumber parses extended serial number data.
//
// The data type `ExtendedSerialNumber` is specified in the Data Dictionary, Section 2.72.
//
// ASN.1 Definition:
//
//	ExtendedSerialNumber ::= SEQUENCE {
//	    serialNumber INTEGER(0..2^32-1),
//	    monthYear BCDString(SIZE(2)),
//	    type EquipmentType,
//	    manufacturerCode ManufacturerCode
//	}
//
// Binary Layout (8 bytes total):
//   - Serial Number (4 bytes): Big-endian uint32
//   - Month/Year (2 bytes): BCD-encoded MMYY format
//   - Equipment Type (1 byte): EquipmentType
//   - Manufacturer Code (1 byte): ManufacturerCode
//
//nolint:unused
func (opts UnmarshalOptions) UnmarshalExtendedSerialNumber(data []byte) (*ddv1.ExtendedSerialNumber, error) {
	const (
		lenExtendedSerialNumber = 8
	)

	if len(data) != lenExtendedSerialNumber {
		return nil, fmt.Errorf("invalid data length for ExtendedSerialNumber: got %d, want %d", len(data), lenExtendedSerialNumber)
	}

	esn := &ddv1.ExtendedSerialNumber{}

	// Parse serial number (4 bytes, big-endian)
	serialNum := binary.BigEndian.Uint32(data[0:4])
	esn.SetSerialNumber(int64(serialNum))

	// Parse month/year BCD (2 bytes, MMYY format)
	monthYear, err := opts.UnmarshalMonthYear(data[4:6])
	if err != nil {
		return nil, fmt.Errorf("failed to parse month/year: %w", err)
	}
	esn.SetMonthYear(monthYear)

	// Parse equipment type (1 byte)
	if equipmentType, err := UnmarshalEnum[ddv1.EquipmentType](data[6]); err == nil {
		esn.SetType(equipmentType)
	} else {
		// Return UNRECOGNIZED for unknown values
		esn.SetType(ddv1.EquipmentType_EQUIPMENT_TYPE_UNRECOGNIZED)
	}

	// Parse manufacturer code (1 byte)
	esn.SetManufacturerCode(int32(data[7]))

	return esn, nil
}

// appendExtendedSerialNumber appends extended serial number data to dst.
//
// The data type `ExtendedSerialNumber` is specified in the Data Dictionary, Section 2.72.
//
// ASN.1 Definition:
//
//	ExtendedSerialNumber ::= SEQUENCE {
//	    serialNumber INTEGER(0..2^32-1),
//	    monthYear BCDString(SIZE(2)),
//	    type EquipmentType,
//	    manufacturerCode ManufacturerCode
//	}
//
// Binary Layout (8 bytes total):
//   - Serial Number (4 bytes): Big-endian uint32
//   - Month/Year (2 bytes): BCD-encoded MMYY format
//   - Equipment Type (1 byte): EquipmentType
//   - Manufacturer Code (1 byte): ManufacturerCode
//
//nolint:unused
func AppendExtendedSerialNumber(dst []byte, esn *ddv1.ExtendedSerialNumber) ([]byte, error) {
	if esn == nil {
		return nil, fmt.Errorf("extendedSerialNumber cannot be nil")
	}

	// Append serial number (4 bytes, big-endian)
	serialNumber := esn.GetSerialNumber()
	dst = binary.BigEndian.AppendUint32(dst, uint32(serialNumber))

	// Append month/year BCD (2 bytes, MMYY format)
	var err error
	dst, err = AppendMonthYear(dst, esn.GetMonthYear())
	if err != nil {
		return nil, fmt.Errorf("failed to append month/year: %w", err)
	}

	// Append equipment type (1 byte)
	equipmentTypeByte, err := MarshalEnum(esn.GetType())
	if err != nil {
		return nil, fmt.Errorf("failed to append equipment type: %w", err)
	}
	dst = append(dst, equipmentTypeByte)

	// Append manufacturer code (1 byte)
	dst = append(dst, byte(esn.GetManufacturerCode()))

	return dst, nil
}

// appendExtendedSerialNumberAsString appends an ExtendedSerialNumber structure as string (legacy compatibility).
// This function maintains compatibility with existing code that expects a string representation.
func AppendExtendedSerialNumberAsString(dst []byte, esn *ddv1.ExtendedSerialNumber, maxLen int) ([]byte, error) {
	if esn == nil {
		return nil, fmt.Errorf("extendedSerialNumber cannot be nil")
	}

	// Create a byte slice for the extended serial number
	serialBytes := make([]byte, 8)

	// First 4 bytes: serial number (big-endian)
	serialNumber := esn.GetSerialNumber()
	if serialNumber != 0 {
		binary.BigEndian.PutUint32(serialBytes[0:4], uint32(serialNumber))
	}

	// Next 2 bytes: month/year (BCD)
	// Append month/year to a temporary buffer, then copy to serialBytes
	tempBytes := make([]byte, 0, 2)
	tempBytes, err := AppendMonthYear(tempBytes, esn.GetMonthYear())
	if err != nil {
		return nil, err
	}
	if len(tempBytes) >= 2 {
		copy(serialBytes[4:6], tempBytes[:2])
	}

	// Next byte: equipment type (converted to protocol value using generic helper)
	if esn.GetType() != ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED {
		protocolValue, err := MarshalEnum(esn.GetType())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal equipment type: %w", err)
		}
		serialBytes[6] = protocolValue
	}

	// Last byte: manufacturer code
	manufacturerCode := esn.GetManufacturerCode()
	if manufacturerCode != 0 {
		serialBytes[7] = byte(manufacturerCode)
	}

	// Truncate or pad to maxLen
	if len(serialBytes) > maxLen {
		serialBytes = serialBytes[:maxLen]
	} else if len(serialBytes) < maxLen {
		padding := make([]byte, maxLen-len(serialBytes))
		serialBytes = append(serialBytes, padding...)
	}

	return append(dst, serialBytes...), nil
}
