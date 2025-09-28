package tachograph

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalIcc parses the binary data for an EF_ICC record.
func unmarshalIcc(data []byte) (*cardv1.Icc, error) {
	var icc cardv1.Icc
	if len(data) < 25 {
		return nil, errors.New("not enough data for IccIdentification")
	}
	r := bytes.NewReader(data)
	clockStop, _ := r.ReadByte()
	// Convert clock stop byte to ClockStopMode enum
	// This is a simplified mapping - the actual mapping should be based on the bit pattern
	clockStopMode := datadictionaryv1.ClockStopMode(clockStop)
	icc.SetClockStop(clockStopMode)
	// Create ExtendedSerialNumber structure
	esn := &datadictionaryv1.ExtendedSerialNumber{}
	// Read the 8-byte extended serial number
	serialBytes := make([]byte, 8)
	if _, err := r.Read(serialBytes); err == nil && len(serialBytes) >= 8 {
		// Parse the fields according to ExtendedSerialNumber structure
		// First 4 bytes: serial number (big-endian)
		serialNum := binary.BigEndian.Uint32(serialBytes[0:4])
		esn.SetSerialNumber(int64(serialNum))

		// Next 2 bytes: month/year BCD (MMYY format)
		if len(serialBytes) > 5 {
			// Decode BCD month/year as 4-digit number MMYY
			monthYearInt, err := bcdBytesToInt(serialBytes[4:6])
			if err == nil && monthYearInt > 0 {
				month := int32(monthYearInt / 100)
				year := int32(monthYearInt % 100)

				// Convert 2-digit year to 4-digit (assuming 20xx for years 00-99)
				if year >= 0 && year <= 99 {
					year += 2000
				}

				esn.SetMonth(month)
				esn.SetYear(year)
			}
		}

		// Next byte: equipment type (convert from protocol value)
		if len(serialBytes) > 6 {
			SetEquipmentType(int32(serialBytes[6]), esn.SetType, nil)
		}

		// Last byte: manufacturer code
		if len(serialBytes) > 7 {
			esn.SetManufacturerCode(int32(serialBytes[7]))
		}
	}
	icc.SetCardExtendedSerialNumber(esn)
	cardApprovalNumber, err := unmarshalIA5StringValueFromReader(r, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to read card approval number: %w", err)
	}
	icc.SetCardApprovalNumber(cardApprovalNumber)
	personaliser, _ := r.ReadByte()
	icc.SetCardPersonaliserId(int32(personaliser))
	// Create EmbedderIcAssemblerId structure
	embedder := make([]byte, 5)
	if _, err := r.Read(embedder); err != nil {
		return nil, fmt.Errorf("failed to read embedder IC assembler ID: %w", err)
	}
	eia := &cardv1.Icc_EmbedderIcAssemblerId{}
	if len(embedder) >= 5 {
		// Store as hex string to avoid UTF-8 validation issues with binary data
		countryCode := &datadictionaryv1.StringValue{}
		countryCode.SetEncoding(datadictionaryv1.Encoding_IA5)
		countryCode.SetEncoded([]byte(fmt.Sprintf("%02X%02X", embedder[0], embedder[1])))
		countryCode.SetDecoded(fmt.Sprintf("%02X%02X", embedder[0], embedder[1]))
		eia.SetCountryCode(countryCode)

		moduleEmbedder := &datadictionaryv1.StringValue{}
		moduleEmbedder.SetEncoding(datadictionaryv1.Encoding_IA5)
		moduleEmbedder.SetEncoded([]byte(fmt.Sprintf("%02X%02X", embedder[2], embedder[3])))
		moduleEmbedder.SetDecoded(fmt.Sprintf("%02X%02X", embedder[2], embedder[3]))
		eia.SetModuleEmbedder(moduleEmbedder)

		eia.SetManufacturerInformation(int32(embedder[4]))
	}
	icc.SetEmbedderIcAssemblerId(eia)
	icIdentifier := make([]byte, 2)
	r.Read(icIdentifier)
	icc.SetIcIdentifier(icIdentifier)
	return &icc, nil
}
