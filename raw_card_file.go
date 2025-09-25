package tachograph

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// MarshalRawCardFile converts a RawCardFile back to binary data
func MarshalRawCardFile(rawFile *cardv1.RawCardFile) ([]byte, error) {
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

// createRawRecord creates a RawCardFile_Record with the given parameters
func createRawRecord(tag int32, fileType cardv1.ElementaryFileType, contentType cardv1.ContentType, data []byte) *cardv1.RawCardFile_Record {
	record := &cardv1.RawCardFile_Record{}
	record.SetTag(tag)
	record.SetFile(fileType)
	record.SetGeneration(cardv1.ApplicationGeneration_GENERATION_1)
	record.SetContentType(contentType)
	record.SetLength(int32(len(data)))
	record.SetValue(data)
	return record
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

// appendEventsDataToBytes marshals events data using the tagged union approach
func appendEventsDataToBytes(dst []byte, card *cardv1.DriverCardFile, eventsData *cardv1.EventData) ([]byte, error) {
	eventsValBuf := make([]byte, 0, 1728) // Max size for Gen1

	// With the tagged union approach, we simply iterate through all records in order
	// Each record either contains valid semantic data or preserved raw bytes
	for _, record := range eventsData.GetRecords() {
		var err error
		eventsValBuf, err = AppendEventRecord(eventsValBuf, record)
		if err != nil {
			return nil, err
		}
	}

	return append(dst, eventsValBuf...), nil
}

// appendFaultsDataToBytes marshals faults data using the tagged union approach
func appendFaultsDataToBytes(dst []byte, card *cardv1.DriverCardFile, faultsData *cardv1.FaultData) ([]byte, error) {
	faultsValBuf := make([]byte, 0, 1152) // Max size for Gen1

	// With the tagged union approach, we simply iterate through all records in order
	// Each record either contains valid semantic data or preserved raw bytes
	for _, record := range faultsData.GetRecords() {
		var err error
		faultsValBuf, err = AppendFaultRecord(faultsValBuf, record)
		if err != nil {
			return nil, err
		}
	}

	return append(dst, faultsValBuf...), nil
}
