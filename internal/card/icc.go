package card

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
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
func (opts UnmarshalOptions) unmarshalIcc(data []byte) (*cardv1.Icc, error) {
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
	// Convert clock stop byte to ClockStopMode enum using generic helper
	if clockStopMode, err := dd.UnmarshalEnum[ddv1.ClockStopMode](data[offset]); err == nil {
		icc.SetClockStop(clockStopMode)
	} else {
		return nil, fmt.Errorf("invalid clock stop mode: %w", err)
	}
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
			monthYear, err := opts.UnmarshalMonthYear(serialBytes[4:6])
			if err != nil {
				return nil, fmt.Errorf("failed to parse month/year: %w", err)
			}
			esn.SetMonthYear(monthYear)
		}

		// Next byte: equipment type (convert from protocol value using generic helper)
		if len(serialBytes) > 6 {
			if equipmentType, err := dd.UnmarshalEnum[ddv1.EquipmentType](serialBytes[6]); err == nil {
				esn.SetType(equipmentType)
			} else {
				return nil, fmt.Errorf("invalid equipment type in extended serial number: %w", err)
			}
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
	cardApprovalNumber, err := opts.UnmarshalIA5StringValue(data[offset : offset+lenCardApprovalNumber])
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
		// Country code (2 bytes, IA5String)
		countryCode, err := opts.UnmarshalIA5StringValue(embedder[0:2])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal country code: %w", err)
		}
		eia.SetCountryCode(countryCode)

		// Module embedder (2 bytes, IA5String)
		moduleEmbedder, err := opts.UnmarshalIA5StringValue(embedder[2:4])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module embedder: %w", err)
		}
		eia.SetModuleEmbedder(moduleEmbedder)

		// Manufacturer information (1 byte)
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
	clockStopByte, err := dd.MarshalEnum(icc.GetClockStop())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal clock stop: %w", err)
	}
	dst = append(dst, clockStopByte)

	// Append extended serial number (8 bytes)
	dst, err = dd.AppendExtendedSerialNumberAsString(dst, icc.GetCardExtendedSerialNumber(), lenCardExtendedSerialNumber)
	if err != nil {
		return nil, err
	}

	// Append card approval number (8 bytes)
	dst, err = dd.AppendStringValue(dst, icc.GetCardApprovalNumber())
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

// appendEmbedderIcAssemblerId appends an EmbedderIcAssemblerId structure (5 bytes total)
func appendEmbedderIcAssemblerId(dst []byte, eia *cardv1.Icc_EmbedderIcAssemblerId) ([]byte, error) {
	const lenEmbedderIcAssemblerId = 5
	const lenCountryCode = 2
	const lenModuleEmbedder = 2

	if eia == nil {
		// Append default values: 5 zero bytes
		return append(dst, make([]byte, lenEmbedderIcAssemblerId)...), nil
	}

	// Append country code (2 bytes, IA5String)
	// Note: IA5String format has no code page byte, just the raw string data
	countryCode := eia.GetCountryCode()
	if countryCode != nil && len(countryCode.GetRawData()) == lenCountryCode {
		dst = append(dst, countryCode.GetRawData()...)
	} else if countryCode != nil {
		// Pad or truncate value to 2 bytes
		value := countryCode.GetValue()
		if len(value) > lenCountryCode {
			dst = append(dst, []byte(value[:lenCountryCode])...)
		} else {
			dst = append(dst, []byte(value)...)
			for i := len(value); i < lenCountryCode; i++ {
				dst = append(dst, ' ')
			}
		}
	} else {
		dst = append(dst, 0x00, 0x00)
	}

	// Append module embedder (2 bytes, IA5String)
	moduleEmbedder := eia.GetModuleEmbedder()
	if moduleEmbedder != nil && len(moduleEmbedder.GetRawData()) == lenModuleEmbedder {
		dst = append(dst, moduleEmbedder.GetRawData()...)
	} else if moduleEmbedder != nil {
		// Pad or truncate value to 2 bytes
		value := moduleEmbedder.GetValue()
		if len(value) > lenModuleEmbedder {
			dst = append(dst, []byte(value[:lenModuleEmbedder])...)
		} else {
			dst = append(dst, []byte(value)...)
			for i := len(value); i < lenModuleEmbedder; i++ {
				dst = append(dst, ' ')
			}
		}
	} else {
		dst = append(dst, 0x00, 0x00)
	}

	// Append manufacturer information (1 byte)
	dst = append(dst, byte(eia.GetManufacturerInformation()))

	return dst, nil
}
