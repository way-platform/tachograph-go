package tachograph

import (
	"encoding/binary"
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalIcc parses the binary data for an EF_ICC record.
//
// The data type `CardIccIdentification` is specified in the Data Dictionary, Section 2.23.
//
// ASN.1 Definition:
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
		lenClockStop                = 1
		lenCardExtendedSerialNumber = 8
		lenCardApprovalNumber       = 8
		lenCardPersonaliserId       = 1
		lenEmbedderIcAssemblerId    = 5
		lenIcIdentifier             = 2
		lenCardIccIdentification    = lenClockStop + lenCardExtendedSerialNumber + lenCardApprovalNumber + lenCardPersonaliserId + lenEmbedderIcAssemblerId + lenIcIdentifier
	)

	var icc cardv1.Icc
	if len(data) < lenCardIccIdentification {
		return nil, errors.New("not enough data for IccIdentification")
	}
	offset := 0

	// Read clock stop (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for clock stop")
	}
	clockStop := data[offset]
	// Convert clock stop byte to ClockStopMode enum using generic helper
	enumDesc := ddv1.ClockStopMode_CLOCK_STOP_MODE_UNSPECIFIED.Descriptor()
	setEnumFromProtocolValue(enumDesc, int32(clockStop),
		func(enumNum protoreflect.EnumNumber) {
			icc.SetClockStop(ddv1.ClockStopMode(enumNum))
		}, nil)
	offset++

	// Create ExtendedSerialNumber structure
	esn := &ddv1.ExtendedSerialNumber{}
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
			enumDesc := ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED.Descriptor()
			setEnumFromProtocolValue(enumDesc, int32(serialBytes[6]),
				func(enumNum protoreflect.EnumNumber) {
					esn.SetType(ddv1.EquipmentType(enumNum))
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
		countryCode := &ddv1.StringValue{}
		countryCode.SetEncoding(ddv1.Encoding_IA5)
		countryCode.SetEncoded([]byte(fmt.Sprintf("%02X%02X", embedder[0], embedder[1])))
		countryCode.SetDecoded(fmt.Sprintf("%02X%02X", embedder[0], embedder[1]))
		eia.SetCountryCode(countryCode)

		moduleEmbedder := &ddv1.StringValue{}
		moduleEmbedder.SetEncoding(ddv1.Encoding_IA5)
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

// AppendIcc appends the binary representation of an EF_ICC message to dst.
//
// The data type `CardIccIdentification` is specified in the Data Dictionary, Section 2.23.
//
// ASN.1 Definition:
//
//	CardIccIdentification ::= SEQUENCE {
//	    clockStop                   OCTET STRING (SIZE(1)),
//	    cardExtendedSerialNumber    ExtendedSerialNumber,    -- 8 bytes
//	    cardApprovalNumber          CardApprovalNumber,      -- 8 bytes
//	    cardPersonaliserID          ManufacturerCode,        -- 1 byte
//	    embedderIcAssemblerId       EmbedderIcAssemblerId,   -- 5 bytes
//	    icIdentifier                OCTET STRING (SIZE(2))
//	}
func appendIcc(dst []byte, icc *cardv1.Icc) ([]byte, error) {
	const (
		lenClockStop                = 1
		lenCardExtendedSerialNumber = 8
		lenCardApprovalNumber       = 8
		lenCardPersonaliserId       = 1
		lenEmbedderIcAssemblerId    = 5
		lenIcIdentifier             = 2
	)

	var err error
	// Append clock stop (1 byte)
	dst = append(dst, byte(icc.GetClockStop()))

	// Append extended serial number (8 bytes)
	dst, err = appendExtendedSerialNumberAsString(dst, icc.GetCardExtendedSerialNumber(), lenCardExtendedSerialNumber)
	if err != nil {
		return nil, err
	}

	// Append card approval number (8 bytes)
	dst, err = appendStringValue(dst, icc.GetCardApprovalNumber(), lenCardApprovalNumber)
	if err != nil {
		return nil, err
	}

	// Append card personaliser ID (1 byte)
	dst = append(dst, byte(icc.GetCardPersonaliserId()))

	// Append embedder IC assembler ID (5 bytes)
	dst, err = appendEmbedderIcAssemblerId(dst, icc.GetEmbedderIcAssemblerId())
	if err != nil {
		return nil, err
	}

	// Append IC identifier (2 bytes)
	dst = append(dst, icc.GetIcIdentifier()...)
	return dst, nil
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
