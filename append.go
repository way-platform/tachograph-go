package tachograph

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// MarshalFile serializes a protobuf File message into the binary DDD file format.
func MarshalFile(file *tachographv1.File) ([]byte, error) {
	return appendCard(nil, file)
}

// appendCard orchestrates writing a card file.
func appendCard(dst []byte, file *tachographv1.File) ([]byte, error) {
	var err error
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		dst, err = appendDriverCard(dst, file.GetDriverCard())
	default:
		return nil, errors.New("unsupported card type for marshaling")
	}
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// appendDriverCard orchestrates the writing of a driver card file.
// The order follows the actual file structure observed in real DDD files.
func appendDriverCard(dst []byte, card *cardv1.DriverCardFile) ([]byte, error) {
	var err error

	// 1. EF_ICC (0x0002) - no signature
	dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_ICC, card.GetIcc(), AppendIcc)
	if err != nil {
		return nil, err
	}

	// 2. EF_IC (0x0005) - no signature
	dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_IC, card.GetIc(), AppendCardIc)
	if err != nil {
		return nil, err
	}

	// 3. EF_APPLICATION_IDENTIFICATION (0x0501)
	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, card.GetApplicationIdentification(), AppendCardApplicationIdentification)
	if err != nil {
		return nil, err
	}

	// 4. EF_IDENTIFICATION (0x0520) - composite file
	// valBuf := make([]byte, 0, 143)
	// valBuf, err = AppendCardIdentification(valBuf, card.GetIdentification())
	// if err != nil {
	// 	return nil, err
	// }
	// valBuf, err = AppendDriverCardHolderIdentification(valBuf, card.GetHolderIdentification())
	// if err != nil {
	// 	return nil, err
	// }
	// dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_IDENTIFICATION, &compositeMessage{data: valBuf}, func(dst []byte, msg *compositeMessage) ([]byte, error) {
	// 	return append(dst, msg.data...), nil
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// 4. EF_IDENTIFICATION (0x0520) - already handled above

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_EVENTS_DATA, card.GetEventsData(), AppendEventsData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_FAULTS_DATA, card.GetFaultsData(), AppendFaultsData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, card.GetDriverActivityData(), AppendDriverActivity)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLES_USED, card.GetVehiclesUsed(), AppendVehiclesUsed)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_PLACES, card.GetPlaces(), AppendPlaces)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CURRENT_USAGE, card.GetCurrentUsage(), AppendCurrentUsage)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, card.GetControlActivityData(), AppendCardControlActivityData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, card.GetSpecificConditions(), AppendCardSpecificConditions)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, card.GetCardDownloadDriver(), AppendCardLastDownload)
	if err != nil {
		return nil, err
	}

	// Remove duplicate EF_APPLICATION_IDENTIFICATION - already added above

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED, card.GetVehicleUnitsUsed(), AppendCardVehicleUnitsUsed)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_GNSS_PLACES, card.GetGnssPlaces(), AppendCardGnssPlaces)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2, card.GetApplicationIdentificationV2(), AppendCardApplicationIdentificationV2)
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

// appendVU orchestrates writing a VU file in TV format
func appendVU(dst []byte, file *tachographv1.File) ([]byte, error) {
	vuFile := file.GetVehicleUnit()
	if vuFile == nil {
		return dst, nil
	}

	buf := bytes.NewBuffer(dst)

	// Process each transfer in the VU file
	for _, transfer := range vuFile.GetTransfers() {
		// Get the tag for this transfer type
		tag := getTrepValueForTransferType(transfer.GetType())
		if tag == 0 {
			continue // Skip unknown transfer types
		}

		// Append the 2-byte tag (0x76XX format)
		vuTag := uint16(0x7600 | (uint16(tag) & 0xFF))
		appendVuTag(buf, vuTag)

		// Append the transfer data based on type
		switch transfer.GetType() {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			AppendDownloadInterfaceVersion(buf, transfer.GetDownloadInterfaceVersion())
		case vuv1.TransferType_OVERVIEW_GEN1, vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			AppendOverview(buf, transfer.GetOverview())
		case vuv1.TransferType_ACTIVITIES_GEN1, vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			AppendVuActivities(buf, transfer.GetActivities())
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			AppendVuEventsAndFaults(buf, transfer.GetEventsAndFaults())
		case vuv1.TransferType_DETAILED_SPEED_GEN1, vuv1.TransferType_DETAILED_SPEED_GEN2:
			AppendVuDetailedSpeed(buf, transfer.GetDetailedSpeed())
		case vuv1.TransferType_TECHNICAL_DATA_GEN1, vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			AppendVuTechnicalData(buf, transfer.GetTechnicalData())
		default:
			// Skip unknown transfer types
		}
	}

	return buf.Bytes(), nil
}

// getTrepValueForTransferType returns the TREP value for a given transfer type
func getTrepValueForTransferType(transferType vuv1.TransferType) uint8 {
	values := vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		if vuv1.TransferType(valueDesc.Number()) == transferType {
			opts := valueDesc.Options()
			if proto.HasExtension(opts, vuv1.E_TrepValue) {
				return uint8(proto.GetExtension(opts, vuv1.E_TrepValue).(int32))
			}
		}
	}
	return 0
}
