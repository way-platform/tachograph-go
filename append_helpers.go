package tachograph

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// appendTimeReal appends a 4-byte TimeReal value.
func appendTimeReal(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	return binary.BigEndian.AppendUint32(dst, uint32(ts.GetSeconds()))
}

// appendDatef appends a 4-byte BCD-encoded date.
func appendDatef(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := ts.AsTime().Year()
	month := int(ts.AsTime().Month())
	day := ts.AsTime().Day()

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
}

func appendOdometer(dst []byte, odometer uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, odometer)
	return append(dst, b[1:]...)
}

// appendString appends a fixed-length string, padding with spaces.
func appendString(dst []byte, s string, length int) ([]byte, error) {
	if len(s) > length {
		return nil, fmt.Errorf("string '%s' is longer than the allowed length %d", s, length)
	}
	result := make([]byte, length)
	copy(result, []byte(s))
	for i := len(s); i < length; i++ {
		result[i] = ' '
	}
	return append(dst, result...), nil
}

// appendBCDNation appends a BCD-encoded nation number.
func appendBCDNation(dst []byte, nation string) ([]byte, error) {
	// This is a placeholder. A real implementation would convert the nation string
	// to its numeric code and then to BCD.
	return append(dst, 0), nil // Append a single zero byte for now
}

// AppendVehicleRegistration appends a VehicleRegistrationIdentification structure.
func AppendVehicleRegistration(dst []byte, nation string, number string) ([]byte, error) {
	// This is also a placeholder.
	dst = append(dst, 0) // Nation
	dst = append(dst, []byte(strings.Repeat(" ", 14))...)
	copy(dst[1:], []byte(number))
	return dst, nil
}

// Helper functions for structured data types

// appendFullCardNumber appends a FullCardNumber structure as a string
func appendFullCardNumber(dst []byte, cardNumber *datadictionaryv1.FullCardNumber, maxLen int) ([]byte, error) {
	if cardNumber == nil {
		return appendString(dst, "", maxLen)
	}
	// For now, just use the card_number field - this may need refinement based on actual format requirements
	return appendString(dst, cardNumber.GetCardNumber(), maxLen)
}

// appendVehicleRegistration appends vehicle registration from VehicleRegistrationIdentification
func appendVehicleRegistration(dst []byte, vehicleReg *datadictionaryv1.VehicleRegistrationIdentification) ([]byte, error) {
	if vehicleReg == nil {
		// Append default values: 1 byte nation (0xFF) + 14 bytes registration number (spaces)
		dst = append(dst, 0xFF)
		return appendString(dst, "", 14)
	}

	// Append nation (1 byte)
	dst = append(dst, byte(vehicleReg.GetNation()))

	// Append registration number (14 bytes, padded with spaces)
	return appendString(dst, vehicleReg.GetNumber(), 14)
}

// appendExtendedSerialNumber appends an ExtendedSerialNumber structure as string (legacy compatibility)
func appendExtendedSerialNumber(dst []byte, esn *datadictionaryv1.ExtendedSerialNumber, maxLen int) ([]byte, error) {
	if esn == nil {
		return append(dst, make([]byte, maxLen)...), nil
	}

	// Create an 8-byte extended serial number based on the structured fields
	serialBytes := make([]byte, 8)

	// First 4 bytes: serial number (big-endian)
	binary.BigEndian.PutUint32(serialBytes[0:4], esn.GetSerialNumber())

	// Next 2 bytes: month/year BCD
	monthYear := esn.GetMonthYear()
	if len(monthYear) >= 2 {
		serialBytes[4] = monthYear[0]
		serialBytes[5] = monthYear[1]
	}

	// Next byte: equipment type (converted to protocol value)
	if esn.GetType() != datadictionaryv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED {
		serialBytes[6] = byte(GetEquipmentTypeProtocolValue(esn.GetType(), 0))
	}

	// Last byte: manufacturer code
	serialBytes[7] = byte(esn.GetManufacturerCode())

	// Append the appropriate number of bytes
	if len(serialBytes) <= maxLen {
		dst = append(dst, serialBytes...)
		// Pad with zeros if needed
		if len(serialBytes) < maxLen {
			dst = append(dst, make([]byte, maxLen-len(serialBytes))...)
		}
	} else {
		dst = append(dst, serialBytes[:maxLen]...)
	}
	return dst, nil
}

// appendEmbedderIcAssemblerId appends an EmbedderIcAssemblerId structure (5 bytes)
func appendEmbedderIcAssemblerId(dst []byte, eia *cardv1.Icc_EmbedderIcAssemblerId) ([]byte, error) {
	if eia == nil {
		// Append 5 bytes of zeros
		return append(dst, 0, 0, 0, 0, 0), nil
	}

	// Country code (2 bytes) - parse hex string back to bytes
	countryCode := eia.GetCountryCode()
	if len(countryCode) == 4 { // hex string format
		if b, err := hex.DecodeString(countryCode); err == nil && len(b) == 2 {
			dst = append(dst, b[0], b[1])
		} else {
			dst = append(dst, 0, 0)
		}
	} else {
		dst = append(dst, 0, 0)
	}

	// Module embedder (2 bytes) - parse hex string back to bytes
	moduleEmbedder := eia.GetModuleEmbedder()
	if len(moduleEmbedder) == 4 { // hex string format
		if b, err := hex.DecodeString(moduleEmbedder); err == nil && len(b) == 2 {
			dst = append(dst, b[0], b[1])
		} else {
			dst = append(dst, 0, 0)
		}
	} else {
		dst = append(dst, 0, 0)
	}

	// Manufacturer information (1 byte)
	dst = append(dst, byte(eia.GetManufacturerInformation()))

	return dst, nil
}
