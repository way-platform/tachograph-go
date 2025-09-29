package tachograph

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// compositeMessage is a helper type for marshalling composite TLV values
type compositeMessage struct {
	data []byte
}

// ProtoReflect implements proto.Message
func (m *compositeMessage) ProtoReflect() protoreflect.Message {
	return nil // Not needed for our use case
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
	var output cardv1.DriverCardFile
	for i := 0; i < len(input.GetRecords()); i++ {
		record := input.GetRecords()[i]
		if record.GetContentType() != cardv1.ContentType_DATA {
			return nil, fmt.Errorf("record %d has unexpected content type", i)
		}
		if !record.HasFile() {
			return nil, fmt.Errorf("record %d has no file type", i)
		}
		var signature []byte
		if i+1 < len(input.GetRecords()) {
			nextRecord := input.GetRecords()[i+1]
			if nextRecord.GetFile() == record.GetFile() && nextRecord.GetContentType() == cardv1.ContentType_SIGNATURE {
				signature = nextRecord.GetValue()
				i++
			}
		}
		switch record.GetFile() {
		case cardv1.ElementaryFileType_EF_ICC:
			icc, err := unmarshalIcc(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_ICC")
			}
			output.SetIcc(icc)

		case cardv1.ElementaryFileType_EF_IC:
			ic, err := unmarshalCardIc(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_IC")
			}
			output.SetIc(ic)

		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification, err := unmarshalIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				identification.SetSignature(signature)
			}
			output.SetIdentification(identification)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId, err := unmarshalCardApplicationIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appId.SetSignature(signature)
			}
			output.SetApplicationIdentification(appId)

		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo, err := unmarshalDrivingLicenceInfo(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				drivingLicenceInfo.SetSignature(signature)
			}
			output.SetDrivingLicenceInfo(drivingLicenceInfo)

		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData, err := unmarshalEventsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				eventsData.SetSignature(signature)
			}
			output.SetEventsData(eventsData)

		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData, err := unmarshalFaultsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				faultsData.SetSignature(signature)
			}
			output.SetFaultsData(faultsData)

		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData, err := unmarshalDriverActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				activityData.SetSignature(signature)
			}
			output.SetDriverActivityData(activityData)

		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed, err := unmarshalCardVehiclesUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehiclesUsed.SetSignature(signature)
			}
			output.SetVehiclesUsed(vehiclesUsed)

		case cardv1.ElementaryFileType_EF_PLACES:
			places, err := unmarshalCardPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				places.SetSignature(signature)
			}
			output.SetPlaces(places)

		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage, err := unmarshalCardCurrentUsage(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				currentUsage.SetSignature(signature)
			}
			output.SetCurrentUsage(currentUsage)

		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity, err := unmarshalCardControlActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				controlActivity.SetSignature(signature)
			}
			output.SetControlActivityData(controlActivity)

		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions, err := unmarshalCardSpecificConditions(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				specificConditions.SetSignature(signature)
			}
			output.SetSpecificConditions(specificConditions)

		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload, err := unmarshalCardLastDownload(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				lastDownload.SetSignature(signature)
			}
			output.SetCardDownloadDriver(lastDownload)

		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits, err := unmarshalCardVehicleUnitsUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehicleUnits.SetSignature(signature)
			}
			output.SetVehicleUnitsUsed(vehicleUnits)

		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces, err := unmarshalCardGnssPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				gnssPlaces.SetSignature(signature)
			}
			output.SetGnssPlaces(gnssPlaces)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2, err := unmarshalCardApplicationIdentificationV2(record.GetValue())
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

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_PLACES, card.GetPlaces(), appendPlaces)
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
