package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalDriverCardFile parses driver card data into a protobuf DriverCardFile message.
func UnmarshalDriverCardFile(rawCard *cardv1.RawCardFile) (*cardv1.DriverCardFile, error) {
	return unmarshalDriverCardFile(rawCard)
}

// MarshalDriverCardFile serializes a DriverCardFile into binary format.
func MarshalDriverCardFile(file *cardv1.DriverCardFile) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("driver card file is nil")
	}

	// Allocate a buffer large enough for the card file
	buf := make([]byte, 0, 1024*1024) // 1MB initial capacity

	// Use the existing appendDriverCard function
	return appendDriverCard(buf, file)
}

// UnmarshalRawCardFile parses raw card data.
func UnmarshalRawCardFile(input []byte) (*cardv1.RawCardFile, error) {
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

// MarshalRawCardFile serializes a RawCardFile into binary format.
func MarshalRawCardFile(file *cardv1.RawCardFile) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("raw card file is nil")
	}

	var result []byte
	for _, record := range file.GetRecords() {
		// Write tag (FID + appendix)
		result = binary.BigEndian.AppendUint16(result, uint16(record.GetTag()>>8))
		result = append(result, byte(record.GetTag()&0xFF))

		// Write length
		result = binary.BigEndian.AppendUint16(result, uint16(record.GetLength()))

		// Write value
		result = append(result, record.GetValue()...)
	}

	return result, nil
}

// InferCardFileType determines the card type from raw card data.
func InferCardFileType(input *cardv1.RawCardFile) cardv1.CardType {
	// Create a copy of the records with File fields set for inference
	recordsWithFileTypes := make([]*cardv1.RawCardFile_Record, len(input.GetRecords()))
	for i, record := range input.GetRecords() {
		// Create a copy of the record
		recordCopy := &cardv1.RawCardFile_Record{}
		recordCopy.SetTag(record.GetTag())
		recordCopy.SetLength(record.GetLength())
		recordCopy.SetValue(record.GetValue())
		recordCopy.SetContentType(record.GetContentType())

		// Set the File field based on the FID
		fid := uint16(record.GetTag() >> 8)
		fileType := MapFidToElementaryFileType(fid)
		recordCopy.SetFile(fileType)

		recordsWithFileTypes[i] = recordCopy
	}

	enumDesc := cardv1.CardType_CARD_TYPE_UNSPECIFIED.Descriptor()
	for i := 0; i < enumDesc.Values().Len(); i++ {
		enumValue := enumDesc.Values().Get(i)
		fileStructure, ok := proto.GetExtension(enumValue.Options(), cardv1.E_FileStructure).(*cardv1.FileDescriptor)
		if !ok {
			continue
		}
		if hasAllElementaryFiles(fileStructure, recordsWithFileTypes) {
			return cardv1.CardType(enumValue.Number())
		}
	}
	return cardv1.CardType_CARD_TYPE_UNSPECIFIED
}

// scanCardFile is a [bufio.SplitFunc] that splits a card file into separate TLV records.
func scanCardFile(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Need at least 5 bytes for TLV header (3 bytes tag + 2 bytes length)
	if len(data) < 5 {
		if atEOF {
			if len(data) == 0 {
				// No more data - this is normal EOF
				return 0, nil, nil
			}
			// We have some data but not enough for a complete header
			return 0, nil, io.ErrUnexpectedEOF
		}
		// Request more data
		return 0, nil, nil
	}
	// Read the length field (bytes 3-4, big-endian)
	length := binary.BigEndian.Uint16(data[3:5])
	// Calculate total record size: 5 bytes header + length bytes value
	totalSize := 5 + int(length)
	// Check if we have enough data for the complete record
	if len(data) < totalSize {
		if atEOF {
			// We're at EOF but don't have enough data - this is an error condition
			return 0, nil, io.ErrUnexpectedEOF
		}
		// Request more data
		return 0, nil, nil
	}
	// We have a complete TLV record.
	return totalSize, data[:totalSize], nil
}

// unmarshalRawCardFileRecord unmarshals a single raw card file record
func unmarshalRawCardFileRecord(input []byte) (*cardv1.RawCardFile_Record, error) {
	var output cardv1.RawCardFile_Record
	// Parse tag: FID (2 bytes) + appendix (1 byte)
	fid := binary.BigEndian.Uint16(input[0:2])
	appendix := input[2]
	output.SetTag((int32(fid) << 8) | int32(appendix))

	// Parse length (2 bytes)
	length := binary.BigEndian.Uint16(input[3:5])
	output.SetLength(int32(length))

	// Parse value - make a copy since input slice may be reused by scanner
	value := make([]byte, length)
	copy(value, input[5:5+length])
	output.SetValue(value)

	// Determine content type based on appendix byte
	// Appendix 0x00 = DATA, 0x01 = SIGNATURE
	if appendix == 0x01 {
		output.SetContentType(cardv1.ContentType_SIGNATURE)
	} else {
		output.SetContentType(cardv1.ContentType_DATA)
	}

	// Don't set File field in raw card records to maintain compatibility
	// The File field will be set only when needed for card type inference

	return &output, nil
}

// MapFidToElementaryFileType maps a FID to its ElementaryFileType using protobuf annotations.
func MapFidToElementaryFileType(fid uint16) cardv1.ElementaryFileType {
	enumDesc := cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED.Descriptor()
	for i := 0; i < enumDesc.Values().Len(); i++ {
		enumValue := enumDesc.Values().Get(i)
		fileId, ok := proto.GetExtension(enumValue.Options(), cardv1.E_FileId).(int32)
		if !ok {
			continue
		}
		if uint16(fileId) == fid {
			return cardv1.ElementaryFileType(enumValue.Number())
		}
	}
	return cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED
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

// unmarshalDriverCardFile unmarshals a driver card file from raw card file data.
//
// The data type `DriverCardFile` represents a complete driver card file structure.
//
// ASN.1 Definition:
//
//	DriverCardFile ::= SEQUENCE {
//	    icc                              CardIccIdentification,
//	    ic                               CardChipIdentification,
//	    identification                   CardIdentification,
//	    applicationIdentification        ApplicationIdentification,
//	    drivingLicenceInfo               CardDrivingLicenceInformation,
//	    eventsData                       CardEventData,
//	    faultsData                       CardFaultData,
//	    driverActivityData               CardDriverActivity,
//	    vehiclesUsed                     CardVehicleRecord,
//	    places                           CardPlaceDailyWorkPeriod,
//	    currentUsage                     CardCurrentUse,
//	    controlActivityData              CardControlActivityDataRecord,
//	    specificConditions               CardSpecificConditionRecord,
//	    cardDownloadDriver               CardDownloadDriver,
//	    vehicleUnitsUsed                 CardVehicleUnitsUsed,
//	    gnssPlaces                       CardGNSSPlaceRecord,
//	    applicationIdentificationV2      ApplicationIdentificationV2,
//	    certificates                     Certificates
//	}
func unmarshalDriverCardFile(input *cardv1.RawCardFile) (*cardv1.DriverCardFile, error) {
	// Initialize with default generation/version (Gen1, Version1)
	// This will be updated when we encounter EF_Application_Identification
	var opts UnmarshalOptions
	var output cardv1.DriverCardFile

	for i := 0; i < len(input.GetRecords()); i++ {
		record := input.GetRecords()[i]
		if record.GetContentType() != cardv1.ContentType_DATA {
			return nil, fmt.Errorf("record %d has unexpected content type", i)
		}
		if !record.HasFile() {
			// Set the File field based on the FID if not already set
			fid := uint16(record.GetTag() >> 8)
			fileType := MapFidToElementaryFileType(fid)
			record.SetFile(fileType)
		}
		var signature []byte
		if i+1 < len(input.GetRecords()) {
			nextRecord := input.GetRecords()[i+1]
			// Set File field for next record too if not set, so we can compare
			if !nextRecord.HasFile() {
				nextFid := uint16(nextRecord.GetTag() >> 8)
				nextFileType := MapFidToElementaryFileType(nextFid)
				nextRecord.SetFile(nextFileType)
			}
			if nextRecord.GetFile() == record.GetFile() && nextRecord.GetContentType() == cardv1.ContentType_SIGNATURE {
				signature = nextRecord.GetValue()
				i++
			}
		}
		switch record.GetFile() {
		case cardv1.ElementaryFileType_EF_ICC:
			icc, err := opts.unmarshalIcc(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_ICC")
			}
			output.SetIcc(icc)

		case cardv1.ElementaryFileType_EF_IC:
			ic, err := opts.unmarshalIc(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_IC")
			}
			output.SetIc(ic)

		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification, err := opts.unmarshalIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				identification.SetSignature(signature)
			}
			output.SetIdentification(identification)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId, err := opts.unmarshalApplicationIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appId.SetSignature(signature)
			}
			output.SetApplicationIdentification(appId)

			// Update opts with the actual generation/version for subsequent EFs
			opts.SetFromCardStructureVersion(appId.GetCardStructureVersion())

		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo, err := opts.unmarshalDrivingLicenceInfo(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				drivingLicenceInfo.SetSignature(signature)
			}
			output.SetDrivingLicenceInfo(drivingLicenceInfo)

		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData, err := opts.unmarshalEventsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				eventsData.SetSignature(signature)
			}
			output.SetEventsData(eventsData)

		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData, err := opts.unmarshalFaultsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				faultsData.SetSignature(signature)
			}
			output.SetFaultsData(faultsData)

		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData, err := opts.unmarshalDriverActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				activityData.SetSignature(signature)
			}
			output.SetDriverActivityData(activityData)

		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed, err := opts.unmarshalVehiclesUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehiclesUsed.SetSignature(signature)
			}
			output.SetVehiclesUsed(vehiclesUsed)

		case cardv1.ElementaryFileType_EF_PLACES:
			places, err := opts.unmarshalPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			// Store the generation in the Places message itself (from opts now, not record)
			places.SetGeneration(opts.Generation)
			if signature != nil {
				places.SetSignature(signature)
			}
			output.SetPlaces(places)

		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage, err := opts.unmarshalCurrentUsage(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				currentUsage.SetSignature(signature)
			}
			output.SetCurrentUsage(currentUsage)

		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity, err := opts.unmarshalControlActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				controlActivity.SetSignature(signature)
			}
			output.SetControlActivityData(controlActivity)

		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions, err := opts.unmarshalSpecificConditions(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				specificConditions.SetSignature(signature)
			}
			output.SetSpecificConditions(specificConditions)

		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload, err := opts.unmarshalLastDownload(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				lastDownload.SetSignature(signature)
			}
			output.SetCardDownloadDriver(lastDownload)

		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits, err := opts.unmarshalVehicleUnitsUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehicleUnits.SetSignature(signature)
			}
			output.SetVehicleUnitsUsed(vehicleUnits)

		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces, err := opts.unmarshalGnssPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				gnssPlaces.SetSignature(signature)
			}
			output.SetGnssPlaces(gnssPlaces)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2, err := opts.unmarshalApplicationIdentificationV2(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appIdV2.SetSignature(signature)
			}
			output.SetApplicationIdentificationV2(appIdV2)

		case cardv1.ElementaryFileType_EF_CARD_CERTIFICATE:
			if output.GetCertificates() == nil {
				output.SetCertificates(&cardv1.Certificates{})
			}
			output.GetCertificates().SetCardCertificate(record.GetValue())
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			if output.GetCertificates() == nil {
				output.SetCertificates(&cardv1.Certificates{})
			}
			output.GetCertificates().SetCaCertificate(record.GetValue())
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CA_CERTIFICATE")
			}
		}
	}
	return &output, nil
}

// appendDriverCard orchestrates the writing of a driver card file.
// The order follows the actual file structure observed in real DDD files.
func appendDriverCard(dst []byte, card *cardv1.DriverCardFile) ([]byte, error) {
	var err error

	// 1. EF_ICC (0x0002) - no signature
	dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_ICC, card.GetIcc(), appendIcc)
	if err != nil {
		return nil, err
	}

	// 2. EF_IC (0x0005) - no signature
	dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_IC, card.GetIc(), appendCardIc)
	if err != nil {
		return nil, err
	}

	// 3. EF_APPLICATION_IDENTIFICATION (0x0501)
	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, card.GetApplicationIdentification(), appendCardApplicationIdentification)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO, card.GetDrivingLicenceInfo(), appendDrivingLicenceInfo)
	if err != nil {
		return nil, err
	}

	// 4. EF_IDENTIFICATION (0x0520) - composite file
	valBuf := make([]byte, 0, 143)
	valBuf, err = appendCardIdentification(valBuf, card.GetIdentification().GetCard())
	if err != nil {
		return nil, err
	}
	valBuf, err = appendDriverCardHolderIdentification(valBuf, card.GetIdentification().GetDriverCardHolder())
	if err != nil {
		return nil, err
	}
	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_IDENTIFICATION, &compositeMessage{data: valBuf}, func(dst []byte, msg *compositeMessage) ([]byte, error) {
		return append(dst, msg.data...), nil
	})
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_EVENTS_DATA, card.GetEventsData(), appendEventsData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_FAULTS_DATA, card.GetFaultsData(), appendFaultsData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, card.GetDriverActivityData(), appendDriverActivity)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLES_USED, card.GetVehiclesUsed(), appendVehiclesUsed)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_PLACES, card.GetPlaces(), func(dst []byte, places *cardv1.Places) ([]byte, error) {
		// Use the generation stored in the Places message itself
		generation := places.GetGeneration()
		if generation == ddv1.Generation_GENERATION_UNSPECIFIED {
			// Default to Generation 1 if not specified
			generation = ddv1.Generation_GENERATION_1
		}
		return appendPlaces(dst, places, generation)
	})
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CURRENT_USAGE, card.GetCurrentUsage(), appendCurrentUsage)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, card.GetControlActivityData(), appendCardControlActivityData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, card.GetSpecificConditions(), appendCardSpecificConditions)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, card.GetCardDownloadDriver(), appendCardLastDownload)
	if err != nil {
		return nil, err
	}

	// Remove duplicate EF_APPLICATION_IDENTIFICATION - already added above

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED, card.GetVehicleUnitsUsed(), appendCardVehicleUnitsUsed)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_GNSS_PLACES, card.GetGnssPlaces(), appendCardGnssPlaces)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2, card.GetApplicationIdentificationV2(), appendCardApplicationIdentificationV2)
	if err != nil {
		return nil, err
	}

	// Append certificate EFs
	if certificates := card.GetCertificates(); certificates != nil {
		dst, err = appendCertificateEF(dst, cardv1.ElementaryFileType_EF_CARD_CERTIFICATE, certificates.GetCardCertificate())
		if err != nil {
			return nil, err
		}

		dst, err = appendCertificateEF(dst, cardv1.ElementaryFileType_EF_CA_CERTIFICATE, certificates.GetCaCertificate())
		if err != nil {
			return nil, err
		}
	}

	// Note: Any remaining proprietary EFs would be handled here if needed

	return dst, nil
}

// compositeMessage is a helper type for marshalling composite TLV values
type compositeMessage struct {
	data []byte
}

// ProtoReflect implements proto.Message
func (m *compositeMessage) ProtoReflect() protoreflect.Message {
	return nil // Not needed for our use case
}

// appendCertificateEF appends a certificate EF (which are not signed)
func appendCertificateEF(dst []byte, fileType cardv1.ElementaryFileType, certData []byte) ([]byte, error) {
	if len(certData) == 0 {
		return dst, nil // Skip empty certificates
	}

	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	// Write data tag (FID + appendix 0x00) - certificates are NOT signed
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x00) // appendix for data
	dst = binary.BigEndian.AppendUint16(dst, uint16(len(certData)))
	dst = append(dst, certData...)

	// Note: Certificates do NOT have signature blocks
	return dst, nil
}

// appendTlv is a generic helper for writing TLV records with zero value-allocation.
func appendTlv[T proto.Message](
	dst []byte,
	fileType cardv1.ElementaryFileType,
	msg T,
	appenderFunc func([]byte, T) ([]byte, error),
) ([]byte, error) {
	// Use reflection to check if the message is nil
	msgValue := reflect.ValueOf(msg)
	if !msgValue.IsValid() || (msgValue.Kind() == reflect.Ptr && msgValue.IsNil()) {
		return dst, nil // Don't write anything if the message is nil
	}

	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	// Write data tag (FID + appendix 0x00) first
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x00) // appendix for data

	// Placeholder for length
	lenPos := len(dst)
	dst = binary.BigEndian.AppendUint16(dst, 0) // Will be updated later

	valPos := len(dst)

	var err error
	dst, err = appenderFunc(dst, msg)
	if err != nil {
		return nil, err
	}

	valLen := len(dst) - valPos

	// Update the length field
	binary.BigEndian.PutUint16(dst[lenPos:], uint16(valLen))

	// Add signature block (FID + appendix 0x01, 128 bytes of zeros for now)
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x01)                       // appendix for signature
	dst = binary.BigEndian.AppendUint16(dst, 128) // Signature length
	// Add 128 bytes of signature data (zeros for now)
	signature := make([]byte, 128)
	dst = append(dst, signature...)

	return dst, nil
}

// appendTlvUnsigned is like appendTlv but doesn't add a signature block
func appendTlvUnsigned[T proto.Message](
	dst []byte,
	fileType cardv1.ElementaryFileType,
	msg T,
	appenderFunc func([]byte, T) ([]byte, error),
) ([]byte, error) {
	// Use reflection to check if the message is nil
	msgValue := reflect.ValueOf(msg)
	if !msgValue.IsValid() || (msgValue.Kind() == reflect.Ptr && msgValue.IsNil()) {
		return dst, nil // Don't write anything if the message is nil
	}

	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	// Write data tag (FID + appendix 0x00) first
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x00) // appendix for data

	// Placeholder for length
	lenPos := len(dst)
	dst = binary.BigEndian.AppendUint16(dst, 0) // Will be updated later

	valPos := len(dst)

	var err error
	dst, err = appenderFunc(dst, msg)
	if err != nil {
		return nil, err
	}

	valLen := len(dst) - valPos

	// Update the length field
	binary.BigEndian.PutUint16(dst[lenPos:], uint16(valLen))

	// No signature block for unsigned EFs
	return dst, nil
}
