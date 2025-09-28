package tachograph

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// bcdBytesToInt converts BCD-encoded bytes to an integer
func bcdBytesToInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	s := hex.EncodeToString(b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid BCD value: %s", s)
	}
	return int(i), nil
}

// createBcdString creates a BcdString message from BCD-encoded bytes
func createBcdString(bcdBytes []byte) (*datadictionaryv1.BcdString, error) {
	decoded, err := bcdBytesToInt(bcdBytes)
	if err != nil {
		return nil, err
	}

	bcdString := &datadictionaryv1.BcdString{}
	bcdString.SetEncoded(bcdBytes)
	bcdString.SetDecoded(int32(decoded))
	return bcdString, nil
}

// createStringValue creates a StringValue message from a string
func createStringValue(s string) *datadictionaryv1.StringValue {
	stringValue := &datadictionaryv1.StringValue{}
	stringValue.SetDecoded(s)
	return stringValue
}

// appendOdometer appends a 3-byte odometer value
func appendOdometer(dst []byte, odometer uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, odometer)
	return append(dst, b[1:]...)
}

// appendControlType appends a ControlType as a single byte bitmask.
func appendControlType(dst []byte, ct *datadictionaryv1.ControlType) []byte {
	if ct == nil {
		return append(dst, 0)
	}

	var b byte
	if ct.GetCardDownloading() {
		b |= 0x80 // bit 'c'
	}
	if ct.GetVuDownloading() {
		b |= 0x40 // bit 'v'
	}
	if ct.GetPrinting() {
		b |= 0x20 // bit 'p'
	}
	if ct.GetDisplay() {
		b |= 0x10 // bit 'd'
	}
	if ct.GetCalibrationChecking() {
		b |= 0x08 // bit 'e'
	}
	// bits 0-2 are RFU (Reserved for Future Use)

	return append(dst, b)
}

// appendExtendedSerialNumber appends an ExtendedSerialNumber structure as string (legacy compatibility)
func appendExtendedSerialNumber(dst []byte, esn *datadictionaryv1.ExtendedSerialNumber, maxLen int) ([]byte, error) {
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

	// Next byte: month (BCD)
	month := esn.GetMonth()
	if month != 0 {
		serialBytes[4] = byte(((month / 10) << 4) | (month % 10))
	}

	// Next byte: year (BCD)
	year := esn.GetYear()
	if year != 0 {
		serialBytes[5] = byte(((year / 10) << 4) | (year % 10))
	}

	// Next byte: equipment type (converted to protocol value using generic helper)
	if esn.GetType() != datadictionaryv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED {
		serialBytes[6] = byte(GetProtocolValueFromEnum(esn.GetType(), 0))
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

// appendEmbedderIcAssemblerId appends an EmbedderIcAssemblerId structure
func appendEmbedderIcAssemblerId(dst []byte, eia *cardv1.Icc_EmbedderIcAssemblerId) ([]byte, error) {
	if eia == nil {
		// Append default values: 1 byte manufacturer code + 1 byte year
		return append(dst, 0x00, 0x00), nil
	}

	// Append manufacturer code (1 byte)
	dst = append(dst, byte(eia.GetManufacturerInformation()))

	// Append year (1 byte) - this is a placeholder since the structure doesn't have a year field
	dst = append(dst, 0x00)

	return dst, nil
}

// appendUint8 appends a single byte to dst
func appendUint8(dst []byte, value uint8) []byte {
	return append(dst, value)
}

// appendUint32 appends a 32-bit unsigned integer to dst
func appendUint32(dst []byte, value uint32) []byte {
	return binary.BigEndian.AppendUint32(dst, value)
}

// appendVuTag appends a 2-byte VU tag to dst

// appendVuBytes appends a byte slice to dst
func appendVuBytes(dst []byte, data []byte) []byte {
	return append(dst, data...)
}

// appendVuString appends a string to dst with a fixed length, padding with null bytes
func appendVuString(dst []byte, s string, length int) []byte {
	result := make([]byte, length)
	copy(result, []byte(s))
	// Pad with null bytes
	for i := len(s); i < length; i++ {
		result[i] = 0
	}
	return append(dst, result...)
}

// appendVuTimeReal appends a TimeReal value (4 bytes) to dst
func appendVuTimeReal(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	return binary.BigEndian.AppendUint32(dst, uint32(ts.GetSeconds()))
}

// appendVuFullCardNumber appends a FullCardNumber to dst with a fixed length
func appendVuFullCardNumber(dst []byte, cardNumber *datadictionaryv1.FullCardNumber, length int) []byte {
	if cardNumber == nil {
		return append(dst, make([]byte, length)...)
	}
	// TODO: Implement proper FullCardNumber serialization
	return append(dst, make([]byte, length)...)
}
