package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

func inferCardFileType(input *cardv1.RawCardFile) cardv1.CardType {
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

func hasAllElementaryFiles(fileStructure *cardv1.FileDescriptor, records []*cardv1.RawCardFile_Record) bool {
	for _, record := range records {
		if record.GetContentType() == cardv1.ContentType_DATA {
			if !containsEF(fileStructure, record.GetFile()) {
				return false
			}
		}
	}
	return true
}

func containsEF(desc *cardv1.FileDescriptor, targetEF cardv1.ElementaryFileType) bool {
	if desc.GetType() == cardv1.FileType_EF && desc.GetEf() == targetEF {
		return true
	}
	for _, child := range desc.GetFiles() {
		if containsEF(child, targetEF) {
			return true
		}
	}
	return false
}
