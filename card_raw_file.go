package tachograph

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalRawCardFile unmarshals a raw card file from binary data.
//
// The data type `RawCardFile` represents a raw card file structure with TLV records.
//
// ASN.1 Definition:
//
//	RawCardFile ::= SEQUENCE OF RawCardFileRecord
//
//	RawCardFileRecord ::= SEQUENCE {
//	    tag        INTEGER,  -- FID (2 bytes) + appendix (1 byte)
//	    length     INTEGER,  -- Length of value (2 bytes)
//	    value      OCTET STRING
//	}
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
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &output, nil
}

// unmarshalRawCardFileRecord unmarshals a single raw card file record.
func unmarshalRawCardFileRecord(input []byte) (*cardv1.RawCardFile_Record, error) {
	var output cardv1.RawCardFile_Record
	// Parse tag: FID (2 bytes) + appendix (1 byte)
	fid := binary.BigEndian.Uint16(input[0:2])
	appendix := input[2]
	output.SetTag((int32(fid) << 8) | int32(appendix))
	length := int32(binary.BigEndian.Uint16(input[3:5]))
	output.SetLength(length)
	// Make a copy of the value bytes to avoid slice sharing issues with the scanner buffer
	value := make([]byte, length)
	copy(value, input[5:5+length])
	output.SetValue(value)
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

// marshalRawCardFile converts a RawCardFile back to binary data.
//
// The data type `RawCardFile` represents a raw card file structure with TLV records.
//
// ASN.1 Definition:
//
//	RawCardFile ::= SEQUENCE OF RawCardFileRecord
//
//	RawCardFileRecord ::= SEQUENCE {
//	    tag        INTEGER,  -- FID (2 bytes) + appendix (1 byte)
//	    length     INTEGER,  -- Length of value (2 bytes)
//	    value      OCTET STRING
//	}
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
