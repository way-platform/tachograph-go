package card

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

// InferFileType determines the card type from raw card data.
func InferFileType(input *cardv1.RawCardFile) cardv1.CardType {
	// The File field is already set during raw parsing, so we can use the records directly
	enumDesc := cardv1.CardType_CARD_TYPE_UNSPECIFIED.Descriptor()
	for i := 0; i < enumDesc.Values().Len(); i++ {
		enumValue := enumDesc.Values().Get(i)
		fileStructure, ok := proto.GetExtension(enumValue.Options(), cardv1.E_FileStructure).(*cardv1.FileDescriptor)
		if !ok {
			continue
		}
		if hasAllElementaryFiles(fileStructure, input.GetRecords()) {
			return cardv1.CardType(enumValue.Number())
		}
	}
	return cardv1.CardType_CARD_TYPE_UNSPECIFIED
}

// mapFidToElementaryFileType maps a FID to its ElementaryFileType using protobuf annotations.
// Returns the file type and true if found, or ELEMENTARY_FILE_UNSPECIFIED and false if not found.
func mapFidToElementaryFileType(fid uint16) (cardv1.ElementaryFileType, bool) {
	enumDesc := cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED.Descriptor()
	for i := 0; i < enumDesc.Values().Len(); i++ {
		enumValue := enumDesc.Values().Get(i)
		fileId, ok := proto.GetExtension(enumValue.Options(), cardv1.E_FileId).(int32)
		if !ok {
			continue
		}
		if uint16(fileId) == fid {
			return cardv1.ElementaryFileType(enumValue.Number()), true
		}
	}
	return cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED, false
}

// hasAllElementaryFiles checks if all required elementary files are present
func hasAllElementaryFiles(fileStructure *cardv1.FileDescriptor, records []*cardv1.RawCardFile_Record) bool {
	// Get all elementary files that should be present for this card type
	expectedFiles := getAllElementaryFiles(fileStructure)

	// Check if all present files are expected for this card type
	for _, record := range records {
		if record.GetContentType() == cardv1.ContentType_DATA {
			found := false
			for _, expectedFile := range expectedFiles {
				if record.GetFile() == expectedFile {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

// getAllElementaryFiles extracts all elementary files from a file structure
func getAllElementaryFiles(desc *cardv1.FileDescriptor) []cardv1.ElementaryFileType {
	var files []cardv1.ElementaryFileType
	if desc.GetType() == cardv1.FileType_EF {
		files = append(files, desc.GetEf())
	}
	for _, child := range desc.GetFiles() {
		files = append(files, getAllElementaryFiles(child)...)
	}
	return files
}
