package card

import (
	"context"
	"encoding/binary"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/way-platform/tachograph-go/internal/dd"
	"github.com/way-platform/tachograph-go/internal/security"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
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
			// Parse and route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
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

				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetApplicationIdentification(appId)

			case ddv1.Generation_GENERATION_2:
				appIdG2, err := opts.unmarshalApplicationIdentificationG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					appIdG2.SetSignature(signature)
				}

				// Extract file-level version from CardStructureVersion for subsequent EFs
				if csv := appIdG2.GetCardStructureVersion(); csv != nil {
					var versionOpts UnmarshalOptions
					versionOpts.SetFromCardStructureVersion(csv)
					fileVersion = versionOpts.Version
				}

				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetApplicationIdentification(appIdG2)

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
			// Parse and route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				vehiclesUsed, err := opts.unmarshalVehiclesUsed(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					vehiclesUsed.SetSignature(signature)
				}
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetVehiclesUsed(vehiclesUsed)

			case ddv1.Generation_GENERATION_2:
				vehiclesUsedG2, err := opts.unmarshalVehiclesUsedG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					vehiclesUsedG2.SetSignature(signature)
				}
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetVehiclesUsed(vehiclesUsedG2)

			default:
				return nil, fmt.Errorf("unexpected generation for EF_VEHICLES_USED: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_PLACES:
			// Parse and route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				places, err := opts.unmarshalPlaces(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					places.SetSignature(signature)
				}
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetPlaces(places)

			case ddv1.Generation_GENERATION_2:
				placesG2, err := opts.unmarshalPlacesG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					placesG2.SetSignature(signature)
				}
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetPlaces(placesG2)

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
			// Parse and route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				specificConditions, err := opts.unmarshalSpecificConditions(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					specificConditions.SetSignature(signature)
				}

				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetSpecificConditions(specificConditions)

			case ddv1.Generation_GENERATION_2:
				specificConditionsG2, err := opts.unmarshalSpecificConditionsG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					specificConditionsG2.SetSignature(signature)
				}

				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetSpecificConditions(specificConditionsG2)

			default:
				return nil, fmt.Errorf("unexpected generation for EF_SPECIFIC_CONDITIONS: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			cardDownload, err := opts.unmarshalCardDownload(record.GetValue())
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
				tachographDF.SetCardDownload(cardDownload)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetCardDownload(cardDownload)
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
			// Gen1: Card authentication certificate
			// Only appears in Gen1 DF (Tachograph)
			if efGeneration != ddv1.Generation_GENERATION_1 {
				return nil, fmt.Errorf("EF_CARD_CERTIFICATE should only appear in Gen1 DF, got generation: %v", efGeneration)
			}
			if tachographDF == nil {
				tachographDF = &cardv1.DriverCardFile_Tachograph{}
			}
			rsaCert, err := dd.UnmarshalOptions{}.UnmarshalRsaCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_CARD_CERTIFICATE: %w", err)
			}
			cert := &cardv1.CardCertificate{}
			cert.SetRsaCertificate(rsaCert)
			tachographDF.SetCardCertificate(cert)
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_CARD_MA_CERTIFICATE:
			// Gen2: Card mutual authentication certificate (replaces Gen1 Card_Certificate)
			// Only appears in Gen2 DF (Tachograph_G2)
			if efGeneration != ddv1.Generation_GENERATION_2 {
				return nil, fmt.Errorf("EF_CARD_MA_CERTIFICATE should only appear in Gen2 DF, got generation: %v", efGeneration)
			}
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			eccCert, err := dd.UnmarshalOptions{}.UnmarshalEccCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_CARD_MA_CERTIFICATE: %w", err)
			}
			cert := &cardv1.CardMaCertificate{}
			cert.SetEccCertificate(eccCert)
			tachographG2DF.SetCardMaCertificate(cert)
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_MA_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_CARD_SIGN_CERTIFICATE:
			// Gen2: Card signature certificate
			// Only appears in Gen2 DF (Tachograph_G2) on driver and workshop cards
			if efGeneration != ddv1.Generation_GENERATION_2 {
				return nil, fmt.Errorf("EF_CARD_SIGN_CERTIFICATE should only appear in Gen2 DF, got generation: %v", efGeneration)
			}
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			eccCert, err := dd.UnmarshalOptions{}.UnmarshalEccCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_CARD_SIGN_CERTIFICATE: %w", err)
			}
			cert := &cardv1.CardSignCertificate{}
			cert.SetEccCertificate(eccCert)
			tachographG2DF.SetCardSignCertificate(cert)
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_SIGN_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			// CA certificate - present in both Gen1 and Gen2
			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				rsaCert, err := dd.UnmarshalOptions{}.UnmarshalRsaCertificate(record.GetValue())
				if err != nil {
					return nil, fmt.Errorf("failed to parse EF_CA_CERTIFICATE (Gen1): %w", err)
				}
				cert := &cardv1.CaCertificate{}
				cert.SetRsaCertificate(rsaCert)
				tachographDF.SetCaCertificate(cert)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				eccCert, err := dd.UnmarshalOptions{}.UnmarshalEccCertificate(record.GetValue())
				if err != nil {
					return nil, fmt.Errorf("failed to parse EF_CA_CERTIFICATE (Gen2): %w", err)
				}
				cert := &cardv1.CaCertificateG2{}
				cert.SetEccCertificate(eccCert)
				tachographG2DF.SetCaCertificate(cert)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CA_CERTIFICATE: %v", efGeneration)
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CA_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_LINK_CERTIFICATE:
			// Gen2: Link certificate for CA chaining
			// Only appears in Gen2 DF (Tachograph_G2)
			if efGeneration != ddv1.Generation_GENERATION_2 {
				return nil, fmt.Errorf("EF_LINK_CERTIFICATE should only appear in Gen2 DF, got generation: %v", efGeneration)
			}
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			eccCert, err := dd.UnmarshalOptions{}.UnmarshalEccCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_LINK_CERTIFICATE: %w", err)
			}
			cert := &cardv1.LinkCertificate{}
			cert.SetEccCertificate(eccCert)
			tachographG2DF.SetLinkCertificate(cert)
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_LINK_CERTIFICATE")
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

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_PLACES, card.GetTachograph().GetPlaces(), appendPlaces)
	if err != nil {
		return nil, err
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

	dst, err = appendTlvUnsigned(dst, cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, card.GetTachograph().GetCardDownload(), appendCardDownload)
	if err != nil {
		return nil, err
	}

	// Gen2 DF - marshal all Gen2 EFs with appendix 0x02/0x03
	if tachographG2 := card.GetTachographG2(); tachographG2 != nil {
		// Marshal Gen2 versions of shared EFs

		// ApplicationIdentification (Gen2)
		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, tachographG2.GetApplicationIdentification(), appendCardApplicationIdentificationG2)
		if err != nil {
			return nil, err
		}

		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_VEHICLES_USED, tachographG2.GetVehiclesUsed(), appendVehiclesUsedG2)
		if err != nil {
			return nil, err
		}

		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_PLACES, tachographG2.GetPlaces(), appendPlacesG2)
		if err != nil {
			return nil, err
		}

		// SpecificConditions (Gen2)
		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, tachographG2.GetSpecificConditions(), appendCardSpecificConditionsG2)
		if err != nil {
			return nil, err
		}

		// Marshal Gen2-exclusive EFs
		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED, tachographG2.GetVehicleUnitsUsed(), appendCardVehicleUnitsUsed)
		if err != nil {
			return nil, err
		}

		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_GNSS_PLACES, tachographG2.GetGnssPlaces(), appendCardGnssPlaces)
		if err != nil {
			return nil, err
		}

		dst, err = appendTlvG2(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2, tachographG2.GetApplicationIdentificationV2(), appendCardApplicationIdentificationV2)
		if err != nil {
			return nil, err
		}
	}

	// Append certificate EFs from Gen1 DF (in regulation order: SFID 2, 4)
	if tachograph := card.GetTachograph(); tachograph != nil {
		// Card authentication certificate (FID C100h)
		if cert := tachograph.GetCardCertificate(); cert != nil {
			if rsaCert := cert.GetRsaCertificate(); rsaCert != nil {
				dst, err = appendCertificateEF(dst, cardv1.ElementaryFileType_EF_CARD_CERTIFICATE, rsaCert.GetRawData())
				if err != nil {
					return nil, err
				}
			}
		}

		// CA certificate (FID C108h)
		if cert := tachograph.GetCaCertificate(); cert != nil {
			if rsaCert := cert.GetRsaCertificate(); rsaCert != nil {
				dst, err = appendCertificateEF(dst, cardv1.ElementaryFileType_EF_CA_CERTIFICATE, rsaCert.GetRawData())
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Append certificate EFs from Gen2 DF (in regulation order: SFID 2, 3, 4, 5)
	if tachographG2 := card.GetTachographG2(); tachographG2 != nil {
		// Card mutual authentication certificate (FID C100h)
		if cert := tachographG2.GetCardMaCertificate(); cert != nil {
			if eccCert := cert.GetEccCertificate(); eccCert != nil {
				dst, err = appendCertificateEFG2(dst, cardv1.ElementaryFileType_EF_CARD_MA_CERTIFICATE, eccCert.GetRawData())
				if err != nil {
					return nil, err
				}
			}
		}

		// Card signature certificate (FID C101h)
		if cert := tachographG2.GetCardSignCertificate(); cert != nil {
			if eccCert := cert.GetEccCertificate(); eccCert != nil {
				dst, err = appendCertificateEFG2(dst, cardv1.ElementaryFileType_EF_CARD_SIGN_CERTIFICATE, eccCert.GetRawData())
				if err != nil {
					return nil, err
				}
			}
		}

		// CA certificate (FID C108h)
		if cert := tachographG2.GetCaCertificate(); cert != nil {
			if eccCert := cert.GetEccCertificate(); eccCert != nil {
				dst, err = appendCertificateEFG2(dst, cardv1.ElementaryFileType_EF_CA_CERTIFICATE, eccCert.GetRawData())
				if err != nil {
					return nil, err
				}
			}
		}

		// Link certificate (FID C109h)
		if cert := tachographG2.GetLinkCertificate(); cert != nil {
			if eccCert := cert.GetEccCertificate(); eccCert != nil {
				dst, err = appendCertificateEFG2(dst, cardv1.ElementaryFileType_EF_LINK_CERTIFICATE, eccCert.GetRawData())
				if err != nil {
					return nil, err
				}
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

// appendCertificateEF appends a Gen1 certificate EF (which are not signed)
// Uses appendix 0x00 for Gen1 DF (Tachograph)
func appendCertificateEF(dst []byte, fileType cardv1.ElementaryFileType, certData []byte) ([]byte, error) {
	if len(certData) == 0 {
		return dst, nil // Skip empty certificates
	}

	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	// Write data tag (FID + appendix 0x00) - Gen1 DF certificates are NOT signed
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x00) // appendix for Gen1 data
	dst = binary.BigEndian.AppendUint16(dst, uint16(len(certData)))
	dst = append(dst, certData...)

	// Note: Certificates do NOT have signature blocks
	return dst, nil
}

// appendCertificateEFG2 appends a Gen2 certificate EF (which are not signed)
// Uses appendix 0x02 for Gen2 DF (Tachograph_G2)
func appendCertificateEFG2(dst []byte, fileType cardv1.ElementaryFileType, certData []byte) ([]byte, error) {
	if len(certData) == 0 {
		return dst, nil // Skip empty certificates
	}

	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	// Write data tag (FID + appendix 0x02) - Gen2 DF certificates are NOT signed
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x02) // appendix for Gen2 data
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

// appendTlvG2 is like appendTlv but uses Gen2 DF appendix (0x02/0x03 instead of 0x00/0x01)
func appendTlvG2[T proto.Message](
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

	// Write data tag (FID + appendix 0x02) first - Gen2 DF
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x02) // appendix for Gen2 data

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

	// Add signature block (FID + appendix 0x03, 128 bytes of zeros for now) - Gen2 DF
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x03)                       // appendix for Gen2 signature
	dst = binary.BigEndian.AppendUint16(dst, 128) // Signature length
	// Add 128 bytes of signature data (zeros for now)
	signature := make([]byte, 128)
	dst = append(dst, signature...)

	return dst, nil
}

// CertificateResolver provides access to tachograph certificates
// needed for signature verification.
type CertificateResolver interface {
	// GetRootCertificate retrieves the European Root CA certificate.
	GetRootCertificate(ctx context.Context) (*securityv1.RootCertificate, error)

	// GetRsaCertificate retrieves an RSA certificate (Generation 1)
	// by its Certificate Holder Reference (CHR).
	GetRsaCertificate(ctx context.Context, chr string) (*securityv1.RsaCertificate, error)

	// GetEccCertificate retrieves an ECC certificate (Generation 2)
	// by its Certificate Holder Reference (CHR).
	GetEccCertificate(ctx context.Context, chr string) (*securityv1.EccCertificate, error)
}

// VerifyOptions configures the signature verification process for driver card files.
type VerifyOptions struct {
	// CertificateResolver is used to resolve CA certificates by their Certificate Authority Reference (CAR).
	// If provided, it will be used to fetch CA certificates for verification.
	// If nil, verification will use the embedded CA certificates from the card file itself.
	CertificateResolver CertificateResolver
}

// VerifyDriverCardFile verifies the certificates in a driver card file.
//
// This function verifies:
//   - Generation 1: Card certificate using the CA certificate
//   - Generation 2: Card sign certificate using the CA certificate
//
// The verification process uses a certificate resolver to fetch CA certificates
// by their Certificate Authority Reference (CAR). If no resolver is configured,
// it falls back to using the embedded CA certificates from the card file itself,
// which contain the public keys needed to verify the card's certificates.
//
// This function mutates the certificate structures by setting their signature_valid
// fields to true or false based on the verification result.
//
// Returns an error if verification fails for any certificate.
func (o VerifyOptions) VerifyDriverCardFile(ctx context.Context, file *cardv1.DriverCardFile) error {
	if file == nil {
		return fmt.Errorf("driver card file cannot be nil")
	}

	// Verify Generation 1 certificates (RSA)
	if tachograph := file.GetTachograph(); tachograph != nil {
		if err := o.verifyGen1Certificates(ctx, tachograph); err != nil {
			return fmt.Errorf("Gen1 certificate verification failed: %w", err)
		}
	}

	// Verify Generation 2 certificates (ECC)
	if tachographG2 := file.GetTachographG2(); tachographG2 != nil {
		if err := o.verifyGen2Certificates(ctx, tachographG2); err != nil {
			return fmt.Errorf("Gen2 certificate verification failed: %w", err)
		}
	}

	return nil
}

// verifyGen1Certificates verifies Generation 1 RSA certificates.
// If a certificate resolver is configured, it fetches CA certificates from the resolver.
// Otherwise, it uses the embedded CA certificate from the card file.
func (o VerifyOptions) verifyGen1Certificates(ctx context.Context, tachograph *cardv1.DriverCardFile_Tachograph) error {
	cardCert := tachograph.GetCardCertificate().GetRsaCertificate()

	if cardCert == nil {
		return fmt.Errorf("card certificate is missing")
	}

	var caCert *ddv1.RsaCertificate
	var err error

	// Convert certificates to security types for verification
	cardCertSec, err := dd.ConvertRsaCertificateToSecurity(cardCert)
	if err != nil {
		return fmt.Errorf("failed to convert card certificate: %w", err)
	}

	var caCertSec *securityv1.RsaCertificate

	if o.CertificateResolver != nil {
		// Use certificate resolver to fetch CA certificate
		car := fmt.Sprintf("%d", cardCert.GetCertificateAuthorityReference())
		caCertSec, err = o.CertificateResolver.GetRsaCertificate(ctx, car)
		if err != nil {
			return fmt.Errorf("failed to fetch CA certificate from resolver: %w", err)
		}

		// For RSA certificates, the public key is extracted during signature recovery.
		// If the CA certificate doesn't have its public key yet, we need to verify it
		// against the root CA first to populate it.
		if len(caCertSec.GetRsaModulus()) == 0 || len(caCertSec.GetRsaExponent()) == 0 {
			// Fetch the root CA certificate
			rootCert, err := o.CertificateResolver.GetRootCertificate(ctx)
			if err != nil {
				return fmt.Errorf("failed to get root CA certificate: %w", err)
			}

			// Verify the CA certificate against the root CA to populate its public key
			if err := security.VerifyRsaCertificateWithRoot(caCertSec, rootCert); err != nil {
				return fmt.Errorf("CA certificate verification failed: %w", err)
			}
		}
	} else {
		// Fall back to embedded CA certificate from card file
		caCert = tachograph.GetCaCertificate().GetRsaCertificate()
		if caCert == nil {
			return fmt.Errorf("CA certificate is missing from card file")
		}
		caCertSec, err = dd.ConvertRsaCertificateToSecurity(caCert)
		if err != nil {
			return fmt.Errorf("failed to convert CA certificate: %w", err)
		}
	}

	// Verify the card certificate using the CA certificate
	if err := security.VerifyRsaCertificateWithCA(cardCertSec, caCertSec); err != nil {
		return fmt.Errorf("card certificate verification failed: %w", err)
	}

	// Copy the verification results back to the original ddv1 certificate
	cardCert.SetSignatureValid(cardCertSec.GetSignatureValid())
	cardCert.SetCertificateHolderReference(parseUint64(cardCertSec.GetCertificateHolderReference()))
	cardCert.SetEndOfValidity(cardCertSec.GetEndOfValidity())
	cardCert.SetRsaModulus(cardCertSec.GetRsaModulus())
	cardCert.SetRsaExponent(cardCertSec.GetRsaExponent())

	return nil
}

// parseUint64 converts a decimal string to uint64, returning 0 on error.
func parseUint64(s string) uint64 {
	var v uint64
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v
}

// verifyGen2Certificates verifies Generation 2 ECC certificates.
// If a certificate resolver is configured, it fetches CA certificates from the resolver.
// Otherwise, it uses the embedded CA certificate from the card file.
func (o VerifyOptions) verifyGen2Certificates(ctx context.Context, tachographG2 *cardv1.DriverCardFile_TachographG2) error {
	cardSignCert := tachographG2.GetCardSignCertificate().GetEccCertificate()

	if cardSignCert == nil {
		return fmt.Errorf("card sign certificate is missing")
	}

	// Convert card sign certificate to security type for verification
	cardSignCertSec, err := dd.ConvertEccCertificateToSecurity(cardSignCert)
	if err != nil {
		return fmt.Errorf("failed to convert card sign certificate: %w", err)
	}

	var caCertSec *securityv1.EccCertificate

	if o.CertificateResolver != nil {
		// Use certificate resolver to fetch CA certificate
		car := fmt.Sprintf("%d", cardSignCert.GetCertificateAuthorityReference())
		caCertSec, err = o.CertificateResolver.GetEccCertificate(ctx, car)
		if err != nil {
			return fmt.Errorf("failed to fetch CA certificate from resolver: %w", err)
		}
	} else {
		// Fall back to embedded CA certificate from card file
		caCert := tachographG2.GetCaCertificate().GetEccCertificate()
		if caCert == nil {
			return fmt.Errorf("CA certificate is missing from card file")
		}
		caCertSec, err = dd.ConvertEccCertificateToSecurity(caCert)
		if err != nil {
			return fmt.Errorf("failed to convert CA certificate: %w", err)
		}
	}

	// Verify the card sign certificate using the CA certificate
	if err := security.VerifyEccCertificateWithCA(cardSignCertSec, caCertSec); err != nil {
		return fmt.Errorf("card sign certificate verification failed: %w", err)
	}

	// Copy the verification results back to the original ddv1 certificate
	cardSignCert.SetSignatureValid(cardSignCertSec.GetSignatureValid())
	cardSignCert.SetCertificateHolderReference(parseUint64(cardSignCertSec.GetCertificateHolderReference()))

	return nil
}
