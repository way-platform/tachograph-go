package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	vupb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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

// GetCardFileStructure retrieves the file structure annotation from a CardType enum value.
// It uses proto.GetExtension to access the custom file_structure option we defined.
func GetCardFileStructure(cardType cardv1.CardType) *cardv1.FileDescriptor {
	// Get the enum value descriptor for the specific card type
	enumValue := cardType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(cardType))
	if enumValue == nil {
		return nil
	}
	opts := enumValue.Options()
	if !proto.HasExtension(opts, cardv1.E_FileStructure) {
		return nil
	}
	ext := proto.GetExtension(opts, cardv1.E_FileStructure)
	fileDesc, ok := ext.(*cardv1.FileDescriptor)
	if !ok {
		return nil
	}
	return fileDesc
}

// InferRawCardFileType iterates through known card types, retrieves their
// file structure via protobuf annotations, and checks for a match.
func InferRawCardFileType(input *cardv1.RawCardFile) cardv1.CardType {
	if input == nil || len(input.GetRecords()) == 0 {
		return cardv1.CardType_CARD_TYPE_UNSPECIFIED
	}
	cardTypes := []cardv1.CardType{
		cardv1.CardType_DRIVER_CARD,
		cardv1.CardType_WORKSHOP_CARD,
		cardv1.CardType_CONTROL_CARD,
		cardv1.CardType_COMPANY_CARD,
	}
	for _, cardType := range cardTypes {
		fileStructure := GetCardFileStructure(cardType)
		if fileStructure == nil {
			// Log an error: card type is missing its structure annotation
			continue
		}
		// Flatten the file structure to get an ordered list of expected EFs.
		expectedEFs := FlattenFileStructure(fileStructure)
		if matchesStructure(input.GetRecords(), expectedEFs) {
			return cardType
		}
	}
	return cardv1.CardType_CARD_TYPE_UNSPECIFIED
}

// FlattenFileStructure recursively traverses the FileDescriptor tree
// and returns a flat, ordered list of file descriptors for all EFs.
func FlattenFileStructure(desc *cardv1.FileDescriptor) []*cardv1.FileDescriptor {
	var flatList []*cardv1.FileDescriptor
	if desc.GetType() == cardv1.FileType_EF {
		flatList = append(flatList, desc)
	}
	for _, child := range desc.GetFiles() {
		flatList = append(flatList, FlattenFileStructure(child)...)
	}
	return flatList
}

// matchesStructure uses a dual-cursor approach to match raw records against the expected EF structure.
func matchesStructure(records []*cardv1.RawCardFile_Record, expectedEFs []*cardv1.FileDescriptor) bool {
	recordCursor := 0
	efCursor := 0
	for recordCursor < len(records) && efCursor < len(expectedEFs) {
		record := records[recordCursor]
		expectedEF := expectedEFs[efCursor]
		if record.GetFile() == expectedEF.GetEf() {
			// Match found, advance both cursors.
			recordCursor++
			efCursor++
		} else if expectedEF.GetConditional() {
			// Expected EF is conditional and not present in the input file.
			// Advance the EF cursor to check the next expected file.
			efCursor++
		} else {
			// A required EF is missing or out of order. This card type does not match.
			return false
		}
	}
	// After consuming all records, ensure any remaining expected EFs are all conditional.
	// This handles cases where trailing files are conditional and not present.
	for efCursor < len(expectedEFs) {
		if !expectedEFs[efCursor].GetConditional() {
			return false // A required EF was not found in the input file.
		}
		efCursor++
	}
	return true
}
