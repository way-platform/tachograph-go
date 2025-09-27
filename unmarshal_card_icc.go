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
	icc.SetClockStop(int32(clockStop))
	// Create ExtendedSerialNumber structure
	esn := &datadictionaryv1.ExtendedSerialNumber{}
	// Read the 8-byte extended serial number
	serialBytes := make([]byte, 8)
	if _, err := r.Read(serialBytes); err == nil && len(serialBytes) >= 8 {
		// Parse the fields according to ExtendedSerialNumber structure
		// First 4 bytes: serial number (big-endian)
		serialNum := binary.BigEndian.Uint32(serialBytes[0:4])
		esn.SetSerialNumber(serialNum)

		// Next 2 bytes: month/year BCD
		esn.SetMonthYear(serialBytes[4:6])

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
	icc.SetCardApprovalNumber(readString(r, 8))
	personaliser, _ := r.ReadByte()
	icc.SetCardPersonaliserId(int32(personaliser))
	// Create EmbedderIcAssemblerId structure
	embedder := make([]byte, 5)
	r.Read(embedder)
	eia := &cardv1.Icc_EmbedderIcAssemblerId{}
	if len(embedder) >= 5 {
		// Store as hex string to avoid UTF-8 validation issues with binary data
		eia.SetCountryCode(fmt.Sprintf("%02X%02X", embedder[0], embedder[1]))
		eia.SetModuleEmbedder(fmt.Sprintf("%02X%02X", embedder[2], embedder[3]))
		eia.SetManufacturerInformation(int32(embedder[4]))
	}
	icc.SetEmbedderIcAssemblerId(eia)
	icIdentifier := make([]byte, 2)
	r.Read(icIdentifier)
	icc.SetIcIdentifier(icIdentifier)
	return &icc, nil
}
