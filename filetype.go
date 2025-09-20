package tacho

import (
	"encoding/binary"

	"github.com/way-platform/tacho-go/tachocard"
	"github.com/way-platform/tacho-go/tachounit"
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
	firstTag := binary.BigEndian.Uint16(data[0:2])
	if firstTag == uint16(tachocard.EF_ICC) {
		return CardFileType
	}

	// Check for vehicle unit file marker.
	// Vehicle unit files use Tag-Value (TV) encoding where each record starts
	// with a 2-byte tag. According to Appendix 7, sections 2.2.6.1-2.2.6.6,
	// VU tags are formed by combining SID 76 Hex + TREP values.
	if tachounit.VuTag(firstTag).IsValid() {
		return UnitFileType
	}

	// Fallback: check for TRTP-based format (less common in practice)
	// According to Appendix 7, Section 2.3, some VU files may start with
	// a TRTP (Transfer Request Parameter) byte.
	if tachounit.DownloadType(data[0]).IsValid() {
		return UnitFileType
	}

	return UnknownFileType
}
