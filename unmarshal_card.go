package tachograph

import (
	"bufio"
	"bytes"
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func unmarshalRawCardFile(input []byte) (*cardv1.RawCardFile, error) {
	var output cardv1.RawCardFile
	sc := bufio.NewScanner(bytes.NewReader(input))
	sc.Split(scanCardFile)
	for sc.Scan() {
		record, err := unmarshalRawCardFileRecord(sc.Bytes())
		if err != nil {
			return nil, err
		}
		output.SetRecords(append(output.GetRecords(), record))
	}
	return &output, nil
}

func unmarshalRawCardFileRecord(input []byte) (*cardv1.RawCardFile_Record, error) {
	var output cardv1.RawCardFile_Record
	// Parse tag: FID (2 bytes) + appendix (1 byte)
	fid := binary.BigEndian.Uint16(input[0:2])
	appendix := input[2]
	output.SetTag((int32(fid) << 8) | int32(appendix))
	output.SetLength(int32(binary.BigEndian.Uint16(input[3:5])))
	output.SetValue(input[5:])
	// Extract generation from bit 1 of appendix byte
	if appendix&0x02 != 0 { // bit 1 = 1
		output.SetGeneration(cardv1.ApplicationGeneration_GENERATION_2)
	} else { // bit 1 = 0
		output.SetGeneration(cardv1.ApplicationGeneration_GENERATION_1)
	}
	// Extract content type from bit 0 of appendix byte
	if appendix&0x01 != 0 { // bit 0 = 1
		output.SetContentType(cardv1.ContentType_SIGNATURE)
	} else { // bit 0 = 0
		output.SetContentType(cardv1.ContentType_DATA)
	}
	// Get file type from FID using existing helper function
	if fileType := getElementaryFileTypeFromTag(int32(fid)); fileType != cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
		output.SetFile(fileType)
	}
	return &output, nil
}
