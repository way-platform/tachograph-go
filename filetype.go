package tachograph

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"
	cardpb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
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
	opts := cardpb.ElementaryFileType_EF_ICC.Descriptor().Values().ByNumber(1).Options()
	efIccTag := proto.GetExtension(opts, cardpb.E_FileId).(int32)

	firstTag := binary.BigEndian.Uint16(data[0:2])
	if firstTag == uint16(efIccTag) {
		return CardFileType
	}

	// TODO: Add check for vehicle unit file marker using protobuf enums

	return UnknownFileType
}
