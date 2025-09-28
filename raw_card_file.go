package tachograph

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// marshalRawCardFile converts a RawCardFile back to binary data
func marshalRawCardFile(rawFile *cardv1.RawCardFile) ([]byte, error) {
	var dst []byte

	for _, record := range rawFile.GetRecords() {
		tag := record.GetTag()
		fid := uint16((tag >> 8) & 0xFFFF)
		appendix := uint8(tag & 0xFF)
		length := record.GetLength()
		value := record.GetValue()

		// Write FID (2 bytes) + appendix (1 byte) + length (2 bytes) + value
		dst = binary.BigEndian.AppendUint16(dst, fid)
		dst = append(dst, appendix)
		dst = binary.BigEndian.AppendUint16(dst, uint16(length))
		dst = append(dst, value...)
	}

	return dst, nil
}

// getElementaryFileTypeFromTag maps FID to ElementaryFileType
func getElementaryFileTypeFromTag(fid int32) cardv1.ElementaryFileType {
	// Check all ElementaryFileType values to find matching file_id
	enumDesc := cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED.Descriptor()
	values := enumDesc.Values()

	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, cardv1.E_FileId) {
			fileId := proto.GetExtension(opts, cardv1.E_FileId).(int32)
			if fileId == fid {
				return cardv1.ElementaryFileType(valueDesc.Number())
			}
		}
	}

	return cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED
}
