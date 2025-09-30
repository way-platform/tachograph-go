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
func UnmarshalExtendedSerialNumber(data []byte) (*ddv1.ExtendedSerialNumber, error) {
	const (
		lenExtendedSerialNumber = 8
	)

	if len(data) < lenExtendedSerialNumber {
		return nil, fmt.Errorf("insufficient data for ExtendedSerialNumber: got %d, want %d", len(data), lenExtendedSerialNumber)
	}

	esn := &ddv1.ExtendedSerialNumber{}

	// Parse serial number (4 bytes, big-endian)
	serialNum := binary.BigEndian.Uint32(data[0:4])
	esn.SetSerialNumber(int64(serialNum))

	// Parse month/year BCD (2 bytes, MMYY format)
	monthYear := &ddv1.MonthYear{}
	monthYear.SetEncoded(data[4:6])

	monthYearInt, err := BcdBytesToInt(data[4:6])
	if err == nil && monthYearInt > 0 {
		month := int32(monthYearInt / 100)
		year := int32(monthYearInt % 100)

		// Convert 2-digit year to 4-digit (assuming 20xx for years 00-99)
		if year >= 0 && year <= 99 {
			year += 2000
		}

		monthYear.SetMonth(month)
		monthYear.SetYear(year)
	}

	esn.SetMonthYear(monthYear)

	// Parse equipment type (1 byte)
	equipmentType, err := UnmarshalEquipmentType(data[6:7])
	if err != nil {
		return nil, fmt.Errorf("failed to parse equipment type: %w", err)
	}
	esn.SetType(equipmentType)

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
		// Append default values (8 zero bytes)
		return append(dst, make([]byte, 8)...), nil
	}

	// Append serial number (4 bytes, big-endian)
	serialNumber := esn.GetSerialNumber()
	dst = binary.BigEndian.AppendUint32(dst, uint32(serialNumber))

	// Append month/year BCD (2 bytes, MMYY format)
	monthYear := esn.GetMonthYear()
	if monthYear != nil && len(monthYear.GetEncoded()) >= 2 {
		// Use the original encoded bytes for perfect round-trip fidelity
		dst = append(dst, monthYear.GetEncoded()[:2]...)
	} else if monthYear != nil {
		// Fall back to encoding from decoded values
		month := monthYear.GetMonth()
		year := monthYear.GetYear()

		// Convert 4-digit year to 2-digit for BCD encoding
		year2Digit := year % 100

		// Encode month as BCD
		monthBCD := byte(((month / 10) << 4) | (month % 10))
		dst = append(dst, monthBCD)

		// Encode year as BCD
		yearBCD := byte(((year2Digit / 10) << 4) | (year2Digit % 10))
		dst = append(dst, yearBCD)
	} else {
		// No month/year data, append zeros
		dst = append(dst, 0, 0)
	}

	// Append equipment type (1 byte)
	dst = AppendEquipmentType(dst, esn.GetType())

	// Append manufacturer code (1 byte)
	dst = append(dst, byte(esn.GetManufacturerCode()))

	return dst, nil
}

// appendExtendedSerialNumberAsString appends an ExtendedSerialNumber structure as string (legacy compatibility).
// This function maintains compatibility with existing code that expects a string representation.
func AppendExtendedSerialNumberAsString(dst []byte, esn *ddv1.ExtendedSerialNumber, maxLen int) ([]byte, error) {
	if esn == nil {
		return append(dst, make([]byte, maxLen)...), nil
	}

	// Create a byte slice for the extended serial number
	serialBytes := make([]byte, 8)

	// First 4 bytes: serial number (big-endian)
	serialNumber := esn.GetSerialNumber()
	if serialNumber != 0 {
		binary.BigEndian.PutUint32(serialBytes[0:4], uint32(serialNumber))
	}

	// Next 2 bytes: month/year (BCD)
	monthYear := esn.GetMonthYear()
	if monthYear != nil && len(monthYear.GetEncoded()) >= 2 {
		// Use the original encoded bytes for perfect round-trip fidelity
		copy(serialBytes[4:6], monthYear.GetEncoded()[:2])
	} else if monthYear != nil {
		// Fall back to encoding from decoded values
		month := monthYear.GetMonth()
		year := monthYear.GetYear()

		if month != 0 {
			serialBytes[4] = byte(((month / 10) << 4) | (month % 10))
		}

		// Convert 4-digit year to 2-digit for BCD encoding
		year2Digit := year % 100
		if year != 0 {
			serialBytes[5] = byte(((year2Digit / 10) << 4) | (year2Digit % 10))
		}
	}

	// Next byte: equipment type (converted to protocol value using generic helper)
	if esn.GetType() != ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED {
		if protocolValue, ok := GetProtocolValueFromEnum(esn.GetType()); ok {
			serialBytes[6] = byte(protocolValue)
		}
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
