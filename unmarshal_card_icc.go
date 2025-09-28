package tachograph

import (
	"encoding/binary"
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalIcc parses the binary data for an EF_ICC record.
//
// ASN.1 Specification (Data Dictionary 2.23):
//
//	CardIccIdentification ::= SEQUENCE {
//	    clockStop                   OCTET STRING (SIZE(1)),
//	    cardExtendedSerialNumber    ExtendedSerialNumber,    -- 8 bytes
//	    cardApprovalNumber          CardApprovalNumber,      -- 8 bytes
//	    cardPersonaliserID          ManufacturerCode,        -- 1 byte
//	    embedderIcAssemblerId       EmbedderIcAssemblerId,   -- 5 bytes
//	    icIdentifier                OCTET STRING (SIZE(2))
//	}
func unmarshalIcc(data []byte) (*cardv1.Icc, error) {
	const (
		// CardIccIdentification layout constants
		lenClockStop                = 1
		lenCardExtendedSerialNumber = 8
		lenCardApprovalNumber       = 8
		lenCardPersonaliserId       = 1
		lenEmbedderIcAssemblerId    = 5
		lenIcIdentifier             = 2
		totalLength                 = lenClockStop + lenCardExtendedSerialNumber + lenCardApprovalNumber + lenCardPersonaliserId + lenEmbedderIcAssemblerId + lenIcIdentifier
	)

	var icc cardv1.Icc
	if len(data) < totalLength {
		return nil, errors.New("not enough data for IccIdentification")
	}
	offset := 0

	// Read clock stop (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for clock stop")
	}
	clockStop := data[offset]
	// Convert clock stop byte to ClockStopMode enum using generic helper
	enumDesc := datadictionaryv1.ClockStopMode_CLOCK_STOP_MODE_UNSPECIFIED.Descriptor()
	SetEnumFromProtocolValue(enumDesc, int32(clockStop),
		func(enumNum protoreflect.EnumNumber) {
			icc.SetClockStop(datadictionaryv1.ClockStopMode(enumNum))
		}, nil)
	offset++

	// Create ExtendedSerialNumber structure
	esn := &datadictionaryv1.ExtendedSerialNumber{}
	// Read the 8-byte extended serial number
	if offset+lenCardExtendedSerialNumber > len(data) {
		return nil, fmt.Errorf("insufficient data for card extended serial number")
	}
	serialBytes := data[offset : offset+lenCardExtendedSerialNumber]
	offset += lenCardExtendedSerialNumber
	if len(serialBytes) >= lenCardExtendedSerialNumber {
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

		// Next byte: equipment type (convert from protocol value using generic helper)
		if len(serialBytes) > 6 {
			enumDesc := datadictionaryv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED.Descriptor()
			SetEnumFromProtocolValue(enumDesc, int32(serialBytes[6]),
				func(enumNum protoreflect.EnumNumber) {
					esn.SetType(datadictionaryv1.EquipmentType(enumNum))
				}, nil)
		}

		// Last byte: manufacturer code
		if len(serialBytes) > 7 {
			esn.SetManufacturerCode(int32(serialBytes[7]))
		}
	}
	icc.SetCardExtendedSerialNumber(esn)

	// Read card approval number (8 bytes)
	if offset+lenCardApprovalNumber > len(data) {
		return nil, fmt.Errorf("insufficient data for card approval number")
	}
	cardApprovalNumber, err := unmarshalIA5StringValue(data[offset : offset+lenCardApprovalNumber])
	if err != nil {
		return nil, fmt.Errorf("failed to read card approval number: %w", err)
	}
	icc.SetCardApprovalNumber(cardApprovalNumber)
	offset += lenCardApprovalNumber

	// Read card personaliser ID (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for card personaliser ID")
	}
	personaliser := data[offset]
	icc.SetCardPersonaliserId(int32(personaliser))
	offset++

	// Create EmbedderIcAssemblerId structure (5 bytes)
	if offset+lenEmbedderIcAssemblerId > len(data) {
		return nil, fmt.Errorf("insufficient data for embedder IC assembler ID")
	}
	embedder := data[offset : offset+lenEmbedderIcAssemblerId]
	offset += lenEmbedderIcAssemblerId
	eia := &cardv1.Icc_EmbedderIcAssemblerId{}
	if len(embedder) >= lenEmbedderIcAssemblerId {
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

	// Read IC identifier (2 bytes)
	if offset+lenIcIdentifier > len(data) {
		return nil, fmt.Errorf("insufficient data for IC identifier")
	}
	icIdentifier := data[offset : offset+lenIcIdentifier]
	// offset += lenIcIdentifier // Not needed as this is the last field
	icc.SetIcIdentifier(icIdentifier)
	return &icc, nil
}
