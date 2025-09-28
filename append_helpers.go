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

// appendDate appends a 4-byte BCD-encoded date from the new Date type.
func appendDate(dst []byte, date *datadictionaryv1.Date) []byte {
	if date == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := int(date.GetYear())
	month := int(date.GetMonth())
	day := int(date.GetDay())

	dst = append(dst, byte((year/1000)%10<<4|(year/100)%10))
	dst = append(dst, byte((year/10)%10<<4|year%10))
	dst = append(dst, byte((month/10)%10<<4|month%10))
	dst = append(dst, byte((day/10)%10<<4|day%10))
	return dst
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

// appendStringValue appends a fixed-length string from a StringValue, padding with spaces.
func appendStringValue(dst []byte, sv *datadictionaryv1.StringValue, length int) ([]byte, error) {
	if sv == nil {
		return appendString(dst, "", length)
	}
	// Use the decoded string if available, otherwise use the encoded bytes as string
	s := sv.GetDecoded()
	if s == "" {
		s = string(sv.GetEncoded())
	}
	return appendString(dst, s, length)
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

	// Handle the CardNumber CHOICE based on card type
	switch cardNumber.GetCardType() {
	case datadictionaryv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			// Concatenate the driver identification components
			identification := driverID.GetIdentification()
			consecutive := driverID.GetConsecutiveIndex()
			replacement := driverID.GetReplacementIndex()
			renewal := driverID.GetRenewalIndex()

			// Build the full driver identification string
			driverStr := ""
			if identification != nil {
				driverStr += identification.GetDecoded()
			}
			if consecutive != nil {
				driverStr += consecutive.GetDecoded()
			}
			if replacement != nil {
				driverStr += replacement.GetDecoded()
			}
			if renewal != nil {
				driverStr += renewal.GetDecoded()
			}
			return appendString(dst, driverStr, maxLen)
		}
	case datadictionaryv1.EquipmentType_WORKSHOP_CARD, datadictionaryv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			identification := ownerID.GetIdentification()
			if identification != nil {
				return appendString(dst, identification.GetDecoded(), maxLen)
			}
		}
	}

	// Fallback to empty string
	return appendString(dst, "", maxLen)
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
	number := vehicleReg.GetNumber()
	if number != nil {
		return appendString(dst, number.GetDecoded(), 14)
	}
	return appendString(dst, "", 14)
}

// appendExtendedSerialNumber appends an ExtendedSerialNumber structure as string (legacy compatibility)
func appendExtendedSerialNumber(dst []byte, esn *datadictionaryv1.ExtendedSerialNumber, maxLen int) ([]byte, error) {
	if esn == nil {
		return append(dst, make([]byte, maxLen)...), nil
	}

	// Create an 8-byte extended serial number based on the structured fields
	serialBytes := make([]byte, 8)

	// First 4 bytes: serial number (big-endian)
	binary.BigEndian.PutUint32(serialBytes[0:4], uint32(esn.GetSerialNumber()))

	// Next 2 bytes: month/year BCD (MMYY format)
	month := esn.GetMonth()
	year := esn.GetYear()
	if month > 0 && year > 0 {
		// Convert 4-digit year to 2-digit year
		year2digit := year % 100
		// Create MMYY as 4-digit integer, then encode as BCD in 2 bytes
		monthYear := int(month*100 + year2digit)
		serialBytes[4] = byte((monthYear/1000)%10<<4 | (monthYear/100)%10)
		serialBytes[5] = byte((monthYear/10)%10<<4 | monthYear%10)
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
	if countryCode != nil {
		countryCodeStr := countryCode.GetDecoded()
		if len(countryCodeStr) == 4 { // hex string format
			if b, err := hex.DecodeString(countryCodeStr); err == nil && len(b) == 2 {
				dst = append(dst, b[0], b[1])
			} else {
				dst = append(dst, 0, 0)
			}
		} else {
			dst = append(dst, 0, 0)
		}
	} else {
		dst = append(dst, 0, 0)
	}

	// Module embedder (2 bytes) - parse hex string back to bytes
	moduleEmbedder := eia.GetModuleEmbedder()
	if moduleEmbedder != nil {
		moduleEmbedderStr := moduleEmbedder.GetDecoded()
		if len(moduleEmbedderStr) == 4 { // hex string format
			if b, err := hex.DecodeString(moduleEmbedderStr); err == nil && len(b) == 2 {
				dst = append(dst, b[0], b[1])
			} else {
				dst = append(dst, 0, 0)
			}
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
