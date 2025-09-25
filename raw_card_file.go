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

// DriverCardFileToRaw converts a DriverCardFile to RawCardFile
func DriverCardFileToRaw(card *cardv1.DriverCardFile) (*cardv1.RawCardFile, error) {
	return DriverCardFileToRawWithSignatures(card, nil)
}

// DriverCardFileToRawWithSignatures converts a DriverCardFile to RawCardFile,
// preserving original signatures from originalRawFile if provided
func DriverCardFileToRawWithSignatures(card *cardv1.DriverCardFile, originalRawFile *cardv1.RawCardFile) (*cardv1.RawCardFile, error) {
	rawFile := &cardv1.RawCardFile{}
	var records []*cardv1.RawCardFile_Record

	// Follow the observed file order from TLV analysis

	// 1. EF_ICC (0x0002) - no signature
	if icc := card.GetIcc(); icc != nil {
		data, err := AppendIcc(nil, icc)
		if err != nil {
			return nil, err
		}
		record := createRawRecord(0x000200, cardv1.ElementaryFileType_EF_ICC, cardv1.ContentType_DATA, data)
		records = append(records, record)
	}

	// 2. EF_IC (0x0005) - no signature
	if ic := card.GetIc(); ic != nil {
		data, err := AppendCardIc(nil, ic)
		if err != nil {
			return nil, err
		}
		record := createRawRecord(0x000500, cardv1.ElementaryFileType_EF_IC, cardv1.ContentType_DATA, data)
		records = append(records, record)
	}

	// 3. EF_APPLICATION_IDENTIFICATION (0x0501) - with signature
	if appId := card.GetApplicationIdentification(); appId != nil {
		data, err := AppendCardApplicationIdentification(nil, appId)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050100, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050100)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050101, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 4. EF_IDENTIFICATION (0x0520) - with signature (composite)
	if identification := card.GetIdentification(); identification != nil || card.GetHolderIdentification() != nil {
		var data []byte
		var err error

		if identification != nil {
			data, err = AppendCardIdentification(data, identification)
			if err != nil {
				return nil, err
			}
		}

		if holderIdentification := card.GetHolderIdentification(); holderIdentification != nil {
			data, err = AppendDriverCardHolderIdentification(data, holderIdentification)
			if err != nil {
				return nil, err
			}
		}

		// Data record
		record := createRawRecord(0x052000, cardv1.ElementaryFileType_EF_IDENTIFICATION, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x052000)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x052001, cardv1.ElementaryFileType_EF_IDENTIFICATION, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 5. EF_EVENTS_DATA (0x0502) - with signature
	if eventsData := card.GetEventsData(); eventsData != nil {
		data, err := appendEventsDataToBytes(nil, card, eventsData)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050200, cardv1.ElementaryFileType_EF_EVENTS_DATA, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050200)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050201, cardv1.ElementaryFileType_EF_EVENTS_DATA, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 6. EF_FAULTS_DATA (0x0503) - with signature
	if faultsData := card.GetFaultsData(); faultsData != nil {
		data, err := appendFaultsDataToBytes(nil, card, faultsData)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050300, cardv1.ElementaryFileType_EF_FAULTS_DATA, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050300)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050301, cardv1.ElementaryFileType_EF_FAULTS_DATA, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 7. EF_DRIVER_ACTIVITY_DATA (0x0504) - with signature
	if driverActivityData := card.GetDriverActivityData(); driverActivityData != nil {
		data, err := AppendDriverActivity(nil, driverActivityData)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050400, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050400)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050401, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 8. EF_VEHICLES_USED (0x0505) - with signature
	if vehiclesUsed := card.GetVehiclesUsed(); vehiclesUsed != nil {
		data, err := AppendVehiclesUsed(nil, vehiclesUsed)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050500, cardv1.ElementaryFileType_EF_VEHICLES_USED, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050500)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050501, cardv1.ElementaryFileType_EF_VEHICLES_USED, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 9. EF_PLACES (0x0506) - with signature
	if places := card.GetPlaces(); places != nil {
		data, err := AppendPlaces(nil, places)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050600, cardv1.ElementaryFileType_EF_PLACES, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050600)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050601, cardv1.ElementaryFileType_EF_PLACES, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 10. EF_CURRENT_USAGE (0x0508) - with signature (Note: FID 0x0508, not 0x0507!)
	if currentUsage := card.GetCurrentUsage(); currentUsage != nil {
		data, err := AppendCurrentUsage(nil, currentUsage)
		if err != nil {
			return nil, err
		}
		// Data record - using 0x0508 as observed in actual file, not 0x0507 from spec
		record := createRawRecord(0x050800, cardv1.ElementaryFileType_EF_CURRENT_USAGE, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050800)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050801, cardv1.ElementaryFileType_EF_CURRENT_USAGE, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 10. EF_CONTROL_ACTIVITY_DATA (0x0508) - with signature
	if controlActivityData := card.GetControlActivityData(); controlActivityData != nil {
		data, err := AppendCardControlActivityData(nil, controlActivityData)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x050800, cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x050800)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x050801, cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 11. EF_SPECIFIC_CONDITIONS (0x0522) - with signature
	if specificConditions := card.GetSpecificConditions(); specificConditions != nil {
		data, err := AppendCardSpecificConditions(nil, specificConditions)
		if err != nil {
			return nil, err
		}
		// Data record
		record := createRawRecord(0x052200, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, cardv1.ContentType_DATA, data)
		records = append(records, record)
		// Signature record (preserve original signature if available)
		originalSig := findOriginalSignature(originalRawFile, 0x052200)
		if originalSig == nil {
			originalSig = make([]byte, 128) // Fallback to zeros if no original
		}
		sigRecord := createRawRecord(0x052201, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, cardv1.ContentType_SIGNATURE, originalSig)
		records = append(records, sigRecord)
	}

	// 12. EF_CARD_CERTIFICATE (0xC100) - no signature
	if certificates := card.GetCertificates(); certificates != nil {
		if cardCert := certificates.GetCardCertificate(); len(cardCert) > 0 {
			record := createRawRecord(0xC10000, cardv1.ElementaryFileType_EF_CARD_CERTIFICATE, cardv1.ContentType_DATA, cardCert)
			records = append(records, record)
		}

		// 13. EF_CA_CERTIFICATE (0xC108) - no signature
		if caCert := certificates.GetCaCertificate(); len(caCert) > 0 {
			record := createRawRecord(0xC10800, cardv1.ElementaryFileType_EF_CA_CERTIFICATE, cardv1.ContentType_DATA, caCert)
			records = append(records, record)
		}
	}

	rawFile.SetRecords(records)
	return rawFile, nil
}

// findOriginalSignature finds the signature for a given data tag in the original RawCardFile
func findOriginalSignature(originalRawFile *cardv1.RawCardFile, dataTag int32) []byte {
	if originalRawFile == nil {
		return nil
	}

	// Signature tag is data tag with appendix 0x01 instead of 0x00
	signatureTag := dataTag + 1 // 0xXXXX00 -> 0xXXXX01

	for _, record := range originalRawFile.GetRecords() {
		if record.GetTag() == signatureTag && record.GetContentType() == cardv1.ContentType_SIGNATURE {
			return record.GetValue()
		}
	}

	return nil
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
