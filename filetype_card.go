package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

func inferCardFileType(input *cardv1.RawCardFile) cardv1.CardType {
	if input == nil || len(input.GetRecords()) == 0 {
		return cardv1.CardType_CARD_TYPE_UNSPECIFIED
	}
	cardTypeDescriptor := cardv1.CardType(0).Descriptor()
	for i := 0; i < cardTypeDescriptor.Values().Len(); i++ {
		cardType := cardTypeDescriptor.Values().Get(i)
		fileStructure, ok := proto.GetExtension(cardType.Options(), cardv1.E_FileStructure).(*cardv1.FileDescriptor)
		if !ok {
			continue
		}
		if matchesStructure(input.GetRecords(), flattenFileStructure(fileStructure)) {
			return cardv1.CardType(cardType.Number())
		}
	}
	return cardv1.CardType_CARD_TYPE_UNSPECIFIED
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

func flattenFileStructure(desc *cardv1.FileDescriptor) []*cardv1.FileDescriptor {
	var flatList []*cardv1.FileDescriptor
	if desc.GetType() == cardv1.FileType_EF {
		flatList = append(flatList, desc)
	}
	for _, child := range desc.GetFiles() {
		flatList = append(flatList, flattenFileStructure(child)...)
	}
	return flatList
}
