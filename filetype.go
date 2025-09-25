package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	vupb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// FileType represents the type of a tachograph file.
type FileType string

const (
	// UnknownFileType is the default file type.
	UnknownFileType FileType = "unknown"
	// CardFileType indicates a driver card file.
	CardFileType FileType = "card"
	// UnitFileType indicates a vehicle unit file.
	UnitFileType FileType = "unit"
)

// InferFileType determines the type of a tachograph file based on its content.
func InferFileType(data []byte) FileType {
	if len(data) < 2 {
		return UnknownFileType
	}
	// Check for the card file marker.
	// According to Appendix 7, Section 3.4.2 of the tachograph regulation,
	// a downloaded card file is a concatenation of Elementary Files (EFs).
	// Each EF is preceded by a 2-byte tag and a 2-byte length.
	// Section 3.3.2 mandates that the first file downloaded is always EF_ICC,
	// which has the File Identifier (tag) 0x0002.
	opts := cardv1.ElementaryFileType_EF_ICC.Descriptor().Values().ByNumber(1).Options()
	efIccTag := proto.GetExtension(opts, cardv1.E_FileId).(int32)
	firstTag := binary.BigEndian.Uint16(data[0:2])
	if firstTag == uint16(efIccTag) {
		return CardFileType
	}
	// Check for VU file markers
	// VU files use TV format with 2-byte tags starting with 0x76xx
	// Check if this looks like a VU tag by examining the first byte
	if data[0] == 0x76 {
		// Check if the second byte corresponds to a valid TREP value
		secondByte := data[1]
		// Check against known VU transfer types
		values := vupb.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
		for i := 0; i < values.Len(); i++ {
			valueDesc := values.Get(i)
			opts := valueDesc.Options()
			if proto.HasExtension(opts, vupb.E_TrepValue) {
				trepValue := proto.GetExtension(opts, vupb.E_TrepValue).(int32)
				if uint8(trepValue) == secondByte {
					return UnitFileType
				}
			}
		}
	}
	return UnknownFileType
}
