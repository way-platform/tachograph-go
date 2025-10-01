package card

import (
	"encoding/binary"
	"fmt"
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

// unmarshalDriverCardFile unmarshals a driver card file from raw card file data.
//
// The driver card file structure is organized into Dedicated Files (DFs):
// - Common EFs (ICC, IC) reside in the Master File (MF)
// - Tachograph DF contains Generation 1 application data
// - Tachograph_G2 DF contains Generation 2 application data
//
// The generation of each EF is determined by the TLV tag appendix byte:
// - '00'/'01' indicates Gen1 (Tachograph DF)
// - '02'/'03' indicates Gen2 (Tachograph_G2 DF)
func unmarshalDriverCardFile(input *cardv1.RawCardFile) (*cardv1.DriverCardFile, error) {
	// File-level version context (extracted from CardStructureVersion)
	// This represents the card's overall version capability
	var fileVersion ddv1.Version = ddv1.Version_VERSION_1
	var output cardv1.DriverCardFile

	// DF-level containers - we populate these as we encounter EFs
	var tachographDF *cardv1.DriverCardFile_Tachograph
	var tachographG2DF *cardv1.DriverCardFile_TachographG2

	for i := 0; i < len(input.GetRecords()); i++ {
		record := input.GetRecords()[i]
		if record.GetContentType() != cardv1.ContentType_DATA {
			return nil, fmt.Errorf("record %d has unexpected content type", i)
		}

		// Use generation already parsed from the TLV tag appendix
		// (set during unmarshalRawCardFileRecord)
		efGeneration := record.GetGeneration()

		// Create UnmarshalOptions with EF-specific generation and file-level version
		opts := UnmarshalOptions{}
		opts.Generation = efGeneration
		opts.Version = fileVersion

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

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetIdentification(identification)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetIdentification(identification)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_IDENTIFICATION: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId, err := opts.unmarshalApplicationIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appId.SetSignature(signature)
			}

			// Extract file-level version from CardStructureVersion for subsequent EFs
			// Generation comes from TLV tag appendix, but version is file-level
			if csv := appId.GetCardStructureVersion(); csv != nil {
				var versionOpts UnmarshalOptions
				versionOpts.SetFromCardStructureVersion(csv)
				fileVersion = versionOpts.Version
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetApplicationIdentification(appId)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetApplicationIdentification(appId)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_APPLICATION_IDENTIFICATION: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo, err := opts.unmarshalDrivingLicenceInfo(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				drivingLicenceInfo.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetDrivingLicenceInfo(drivingLicenceInfo)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetDrivingLicenceInfo(drivingLicenceInfo)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_DRIVING_LICENCE_INFO: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData, err := opts.unmarshalEventsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				eventsData.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetEventsData(eventsData)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetEventsData(eventsData)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_EVENTS_DATA: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData, err := opts.unmarshalFaultsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				faultsData.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetFaultsData(faultsData)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetFaultsData(faultsData)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_FAULTS_DATA: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData, err := opts.unmarshalDriverActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				activityData.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetDriverActivityData(activityData)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetDriverActivityData(activityData)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_DRIVER_ACTIVITY_DATA: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed, err := opts.unmarshalVehiclesUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehiclesUsed.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetVehiclesUsed(vehiclesUsed)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetVehiclesUsed(vehiclesUsed)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_VEHICLES_USED: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_PLACES:
			places, err := opts.unmarshalPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			// Store the EF-specific generation in the Places message
			places.SetGeneration(opts.Generation)
			if signature != nil {
				places.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetPlaces(places)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetPlaces(places)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_PLACES: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage, err := opts.unmarshalCurrentUsage(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				currentUsage.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetCurrentUsage(currentUsage)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetCurrentUsage(currentUsage)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CURRENT_USAGE: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity, err := opts.unmarshalControlActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				controlActivity.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetControlActivityData(controlActivity)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetControlActivityData(controlActivity)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CONTROL_ACTIVITY_DATA: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions, err := opts.unmarshalSpecificConditions(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				specificConditions.SetSignature(signature)
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetSpecificConditions(specificConditions)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetSpecificConditions(specificConditions)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_SPECIFIC_CONDITIONS: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload, err := opts.unmarshalLastDownload(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_DOWNLOAD_DRIVER")
			}

			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetCardDownload(lastDownload)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetCardDownload(lastDownload)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CARD_DOWNLOAD_DRIVER: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits, err := opts.unmarshalVehicleUnitsUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehicleUnits.SetSignature(signature)
			}

			// Only Gen2
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			tachographG2DF.SetVehicleUnitsUsed(vehicleUnits)

		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces, err := opts.unmarshalGnssPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				gnssPlaces.SetSignature(signature)
			}

			// Only Gen2
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			tachographG2DF.SetGnssPlaces(gnssPlaces)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2, err := opts.unmarshalApplicationIdentificationV2(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appIdV2.SetSignature(signature)
			}

			// Only Gen2
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			tachographG2DF.SetApplicationIdentificationV2(appIdV2)

		case cardv1.ElementaryFileType_EF_CARD_CERTIFICATE:
			// Certificates are shared between Gen1 and Gen2 (based on which DF they appear in)
			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				if tachographDF.GetCertificates() == nil {
					tachographDF.SetCertificates(&cardv1.Certificates{})
				}
				tachographDF.GetCertificates().SetCardCertificate(record.GetValue())
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				if tachographG2DF.GetCertificates() == nil {
					tachographG2DF.SetCertificates(&cardv1.Certificates{})
				}
				tachographG2DF.GetCertificates().SetCardCertificate(record.GetValue())
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CARD_CERTIFICATE: %v", efGeneration)
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				if tachographDF.GetCertificates() == nil {
					tachographDF.SetCertificates(&cardv1.Certificates{})
				}
				tachographDF.GetCertificates().SetCaCertificate(record.GetValue())
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				if tachographG2DF.GetCertificates() == nil {
					tachographG2DF.SetCertificates(&cardv1.Certificates{})
				}
				tachographG2DF.GetCertificates().SetCaCertificate(record.GetValue())
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CA_CERTIFICATE: %v", efGeneration)
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CA_CERTIFICATE")
			}
		}
	}

	// Set the DFs on the output if they have content
	if tachographDF != nil {
		output.SetTachograph(tachographDF)
	}
	if tachographG2DF != nil {
		output.SetTachographG2(tachographG2DF)
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
	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, card.GetTachograph().GetApplicationIdentification(), appendCardApplicationIdentification)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO, card.GetTachograph().GetDrivingLicenceInfo(), appendDrivingLicenceInfo)
	if err != nil {
		return nil, err
	}

	// 4. EF_IDENTIFICATION (0x0520) - composite file
	if identification := card.GetTachograph().GetIdentification(); identification != nil {
		valBuf := make([]byte, 0, 143)
		valBuf, err = appendCardIdentification(valBuf, identification.GetCard())
		if err != nil {
			return nil, err
		}
		valBuf, err = appendDriverCardHolderIdentification(valBuf, identification.GetDriverCardHolder())
		if err != nil {
			return nil, err
		}
		dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_IDENTIFICATION, &compositeMessage{data: valBuf}, func(dst []byte, msg *compositeMessage) ([]byte, error) {
			return append(dst, msg.data...), nil
		})
		if err != nil {
			return nil, err
		}
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_EVENTS_DATA, card.GetTachograph().GetEventsData(), appendEventsData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_FAULTS_DATA, card.GetTachograph().GetFaultsData(), appendFaultsData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, card.GetTachograph().GetDriverActivityData(), appendDriverActivity)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLES_USED, card.GetTachograph().GetVehiclesUsed(), appendVehiclesUsed)
	if err != nil {
		return nil, err
	}

	if places := card.GetTachograph().GetPlaces(); places != nil {
		dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_PLACES, places, func(dst []byte, places *cardv1.Places) ([]byte, error) {
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
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CURRENT_USAGE, card.GetTachograph().GetCurrentUsage(), appendCurrentUsage)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, card.GetTachograph().GetControlActivityData(), appendCardControlActivityData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, card.GetTachograph().GetSpecificConditions(), appendCardSpecificConditions)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, card.GetTachograph().GetCardDownload(), appendCardLastDownload)
	if err != nil {
		return nil, err
	}

	// Gen2 DFs - only append if present
	if tachographG2 := card.GetTachographG2(); tachographG2 != nil {
		dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED, tachographG2.GetVehicleUnitsUsed(), appendCardVehicleUnitsUsed)
		if err != nil {
			return nil, err
		}

		dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_GNSS_PLACES, tachographG2.GetGnssPlaces(), appendCardGnssPlaces)
		if err != nil {
			return nil, err
		}

		dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2, tachographG2.GetApplicationIdentificationV2(), appendCardApplicationIdentificationV2)
		if err != nil {
			return nil, err
		}
	}

	// Append certificate EFs from Gen1 DF
	if tachograph := card.GetTachograph(); tachograph != nil {
		if certificates := tachograph.GetCertificates(); certificates != nil {
			dst, err = appendCertificateEF(dst, cardv1.ElementaryFileType_EF_CARD_CERTIFICATE, certificates.GetCardCertificate())
			if err != nil {
				return nil, err
			}

			dst, err = appendCertificateEF(dst, cardv1.ElementaryFileType_EF_CA_CERTIFICATE, certificates.GetCaCertificate())
			if err != nil {
				return nil, err
			}
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
