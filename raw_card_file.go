package tachograph

import (
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalRawCardFile parses binary card data into a RawCardFile with raw TLV records
func UnmarshalRawCardFile(data []byte) (*cardv1.RawCardFile, error) {
	rawFile := &cardv1.RawCardFile{}
	var records []*cardv1.RawCardFile_Record

	for offset := 0; offset < len(data); {
		if len(data[offset:]) < 5 { // Need at least 3 bytes tag + 2 bytes length
			break
		}

		// Read tag (FID + appendix)
		fid := binary.BigEndian.Uint16(data[offset:])
		appendix := data[offset+2]
		tag := (int32(fid) << 8) | int32(appendix)

		// Read length
		length := binary.BigEndian.Uint16(data[offset+3:])

		if len(data[offset+5:]) < int(length) {
			return nil, fmt.Errorf("insufficient data for record at offset %d: need %d bytes, have %d", offset, length, len(data[offset+5:]))
		}

		// Read value
		value := make([]byte, length)
		copy(value, data[offset+5:offset+5+int(length)])

		// Determine file type and content type
		fileType := getElementaryFileTypeFromTag(int32(fid))
		contentType := cardv1.ContentType_DATA
		if appendix == 0x01 {
			contentType = cardv1.ContentType_SIGNATURE
		}

		record := &cardv1.RawCardFile_Record{}
		record.SetTag(tag)
		record.SetFile(fileType)
		record.SetGeneration(cardv1.ApplicationGeneration_GENERATION_1) // Default to Gen1
		record.SetContentType(contentType)
		record.SetLength(int32(length))
		record.SetValue(value)

		records = append(records, record)
		offset += 5 + int(length)
	}

	rawFile.SetRecords(records)
	return rawFile, nil
}

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
		// Signature record (128 bytes of zeros for now)
		sigRecord := createRawRecord(0x050101, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x052001, cardv1.ElementaryFileType_EF_IDENTIFICATION, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x050201, cardv1.ElementaryFileType_EF_EVENTS_DATA, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x050301, cardv1.ElementaryFileType_EF_FAULTS_DATA, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x050401, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x050501, cardv1.ElementaryFileType_EF_VEHICLES_USED, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x050601, cardv1.ElementaryFileType_EF_PLACES, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x050801, cardv1.ElementaryFileType_EF_CURRENT_USAGE, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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
		// Signature record
		sigRecord := createRawRecord(0x052201, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, cardv1.ContentType_SIGNATURE, make([]byte, 128))
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

// appendEventsDataToBytes marshals events data using the complex logic from the original appendDriverCard
func appendEventsDataToBytes(dst []byte, card *cardv1.DriverCardFile, eventsData *cardv1.EventData) ([]byte, error) {
	eventsValBuf := make([]byte, 0, 1728) // Max size for Gen1
	eventsPerType := int(card.GetApplicationIdentification().GetEventsPerTypeCount())
	allEvents := eventsData.GetRecords()

	eventsByType := make(map[int32][]*cardv1.EventData_Record)
	for _, e := range allEvents {
		protocolValue := GetEventFaultTypeProtocolValue(e.GetEventType(), e.GetUnrecognizedEventType())
		eventsByType[protocolValue] = append(eventsByType[protocolValue], e)
	}

	// The 6 event groups in a Gen1 card file structure, ordered by type code.
	eventGroupTypeCodes := []int32{0x01, 0x02, 0x03, 0x04, 0x05, 0x07} // Example codes

	for _, eventTypeCode := range eventGroupTypeCodes {
		groupEvents := eventsByType[eventTypeCode]
		for j := 0; j < eventsPerType; j++ {
			if j < len(groupEvents) {
				var err error
				eventsValBuf, err = AppendEventRecord(eventsValBuf, groupEvents[j])
				if err != nil {
					return nil, err
				}
			} else {
				// Pad with an empty 24-byte record
				eventsValBuf = append(eventsValBuf, make([]byte, 24)...)
			}
		}
	}

	return append(dst, eventsValBuf...), nil
}

// appendFaultsDataToBytes marshals faults data using the complex logic from the original appendDriverCard
func appendFaultsDataToBytes(dst []byte, card *cardv1.DriverCardFile, faultsData *cardv1.FaultData) ([]byte, error) {
	faultsValBuf := make([]byte, 0, 1152) // Max size for Gen1
	faultsPerType := int(card.GetApplicationIdentification().GetFaultsPerTypeCount())
	allFaults := faultsData.GetRecords()

	faultsByType := make(map[bool][]*cardv1.FaultData_Record)
	for _, f := range allFaults {
		isEquipmentFault := (f.GetFaultType() >= 0x30 && f.GetFaultType() <= 0x3F)
		faultsByType[isEquipmentFault] = append(faultsByType[isEquipmentFault], f)
	}

	// Order: Equipment faults (true), then Card faults (false)
	faultGroupOrder := []bool{true, false}

	for _, isEquipmentFault := range faultGroupOrder {
		groupFaults := faultsByType[isEquipmentFault]
		for j := 0; j < faultsPerType; j++ {
			if j < len(groupFaults) {
				var err error
				faultsValBuf, err = AppendFaultRecord(faultsValBuf, groupFaults[j])
				if err != nil {
					return nil, err
				}
			} else {
				// Pad with an empty 24-byte record
				faultsValBuf = append(faultsValBuf, make([]byte, 24)...)
			}
		}
	}

	return append(dst, faultsValBuf...), nil
}
