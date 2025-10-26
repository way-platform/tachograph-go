package card

import (
	"context"
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/way-platform/tachograph-go/internal/security"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// MarshalDriverCardFile serializes a DriverCardFile into binary format.
func (opts MarshalOptions) MarshalDriverCardFile(file *cardv1.DriverCardFile) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("driver card file is nil")
	}

	// Allocate a buffer large enough for the card file
	buf := make([]byte, 0, 1024*1024) // 1MB initial capacity

	// Use the existing appendDriverCard function
	// TODO: Pass opts.UseRawData through to append functions
	return appendDriverCard(buf, file)
}

// ParseRawDriverCardFile parses driver card data into a protobuf DriverCardFile message.
//
// The driver card file structure is organized into Dedicated Files (DFs):
// - Common EFs (ICC, IC) reside in the Master File (MF)
// - Tachograph DF contains Generation 1 application data
// - Tachograph_G2 DF contains Generation 2 application data
//
// The generation of each EF is determined by the TLV tag appendix byte:
// - '00'/'01' indicates Gen1 (Tachograph DF)
// - '02'/'03' indicates Gen2 (Tachograph_G2 DF)
func (opts ParseOptions) ParseRawDriverCardFile(input *cardv1.RawCardFile) (*cardv1.DriverCardFile, error) {
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

		// Create UnmarshalOptions with PreserveRawData from ParseOptions
		unmarshalOpts := opts.unmarshal()

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
			icc, err := unmarshalOpts.unmarshalIcc(record.GetValue())
			if err != nil {
				return nil, err
			}
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				icc.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				icc.SetAuthentication(auth)
			}
			output.SetIcc(icc)

		case cardv1.ElementaryFileType_EF_IC:
			ic, err := unmarshalOpts.unmarshalIc(record.GetValue())
			if err != nil {
				return nil, err
			}
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				ic.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				ic.SetAuthentication(auth)
			}
			output.SetIc(ic)

		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification, err := unmarshalOpts.unmarshalDriverCardIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				identification.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				identification.SetAuthentication(auth)
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
				appId, err := unmarshalOpts.unmarshalApplicationIdentification(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					appId.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					appId.SetAuthentication(auth)
				}

				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetApplicationIdentification(appId)

			case ddv1.Generation_GENERATION_2:
				appIdG2, err := unmarshalOpts.unmarshalApplicationIdentificationG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					appIdG2.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					appIdG2.SetAuthentication(auth)
				}

				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetApplicationIdentification(appIdG2)

			default:
				return nil, fmt.Errorf("unexpected generation for EF_APPLICATION_IDENTIFICATION: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo, err := unmarshalOpts.unmarshalDrivingLicenceInfo(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				drivingLicenceInfo.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				drivingLicenceInfo.SetAuthentication(auth)
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
			eventsData, err := unmarshalOpts.unmarshalEventsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				eventsData.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				eventsData.SetAuthentication(auth)
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
			faultsData, err := unmarshalOpts.unmarshalFaultsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				faultsData.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				faultsData.SetAuthentication(auth)
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
			activityData, err := unmarshalOpts.unmarshalDriverActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				activityData.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				activityData.SetAuthentication(auth)
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
				vehiclesUsed, err := unmarshalOpts.unmarshalVehiclesUsed(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					vehiclesUsed.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					vehiclesUsed.SetAuthentication(auth)
				}
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetVehiclesUsed(vehiclesUsed)

			case ddv1.Generation_GENERATION_2:
				vehiclesUsedG2, err := unmarshalOpts.unmarshalVehiclesUsedG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					vehiclesUsedG2.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					vehiclesUsedG2.SetAuthentication(auth)
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
				places, err := unmarshalOpts.unmarshalPlaces(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					places.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					places.SetAuthentication(auth)
				}
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetPlaces(places)

			case ddv1.Generation_GENERATION_2:
				placesG2, err := unmarshalOpts.unmarshalPlacesG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					placesG2.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					placesG2.SetAuthentication(auth)
				}
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetPlaces(placesG2)

			default:
				return nil, fmt.Errorf("unexpected generation for EF_PLACES: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage, err := unmarshalOpts.unmarshalCurrentUsage(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				currentUsage.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				currentUsage.SetAuthentication(auth)
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
			controlActivity, err := unmarshalOpts.unmarshalControlActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				controlActivity.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				controlActivity.SetAuthentication(auth)
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
				specificConditions, err := unmarshalOpts.unmarshalSpecificConditions(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					specificConditions.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					specificConditions.SetAuthentication(auth)
				}

				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				tachographDF.SetSpecificConditions(specificConditions)

			case ddv1.Generation_GENERATION_2:
				specificConditionsG2, err := unmarshalOpts.unmarshalSpecificConditionsG2(record.GetValue())
				if err != nil {
					return nil, err
				}
				if signature != nil {
					specificConditionsG2.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					specificConditionsG2.SetAuthentication(auth)
				}

				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				tachographG2DF.SetSpecificConditions(specificConditionsG2)

			default:
				return nil, fmt.Errorf("unexpected generation for EF_SPECIFIC_CONDITIONS: %v", efGeneration)
			}

		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			cardDownload, err := unmarshalOpts.unmarshalCardDownload(record.GetValue())
			if err != nil {
				return nil, err
			}
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				cardDownload.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				cardDownload.SetAuthentication(auth)
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
			vehicleUnits, err := unmarshalOpts.unmarshalVehicleUnitsUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehicleUnits.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				vehicleUnits.SetAuthentication(auth)
			}

			// Only Gen2
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			tachographG2DF.SetVehicleUnitsUsed(vehicleUnits)

		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces, err := unmarshalOpts.unmarshalGnssPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				gnssPlaces.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				gnssPlaces.SetAuthentication(auth)
			}

			// Only Gen2
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			tachographG2DF.SetGnssPlaces(gnssPlaces)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2, err := unmarshalOpts.unmarshalApplicationIdentificationV2(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appIdV2.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				appIdV2.SetAuthentication(auth)
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
			rsaCert, err := security.UnmarshalRsaCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_CARD_CERTIFICATE: %w", err)
			}
			cert := &cardv1.CardCertificate{}
			cert.SetRsaCertificate(rsaCert)
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				cert.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				cert.SetAuthentication(auth)
			}
			tachographDF.SetCardCertificate(cert)

		case cardv1.ElementaryFileType_EF_CARD_MA_CERTIFICATE:
			// Gen2: Card mutual authentication certificate (replaces Gen1 Card_Certificate)
			// Only appears in Gen2 DF (Tachograph_G2)
			if efGeneration != ddv1.Generation_GENERATION_2 {
				return nil, fmt.Errorf("EF_CARD_MA_CERTIFICATE should only appear in Gen2 DF, got generation: %v", efGeneration)
			}
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			eccCert, err := security.UnmarshalEccCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_CARD_MA_CERTIFICATE: %w", err)
			}
			cert := &cardv1.CardMaCertificate{}
			cert.SetEccCertificate(eccCert)
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				cert.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				cert.SetAuthentication(auth)
			}
			tachographG2DF.SetCardMaCertificate(cert)

		case cardv1.ElementaryFileType_EF_CARD_SIGN_CERTIFICATE:
			// Gen2: Card signature certificate
			// Only appears in Gen2 DF (Tachograph_G2) on driver and workshop cards
			if efGeneration != ddv1.Generation_GENERATION_2 {
				return nil, fmt.Errorf("EF_CARD_SIGN_CERTIFICATE should only appear in Gen2 DF, got generation: %v", efGeneration)
			}
			if tachographG2DF == nil {
				tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
			}
			eccCert, err := security.UnmarshalEccCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_CARD_SIGN_CERTIFICATE: %w", err)
			}
			cert := &cardv1.CardSignCertificate{}
			cert.SetEccCertificate(eccCert)
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				cert.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				cert.SetAuthentication(auth)
			}
			tachographG2DF.SetCardSignCertificate(cert)

		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			// CA certificate - present in both Gen1 and Gen2
			// Route to appropriate DF based on generation
			switch efGeneration {
			case ddv1.Generation_GENERATION_1:
				if tachographDF == nil {
					tachographDF = &cardv1.DriverCardFile_Tachograph{}
				}
				rsaCert, err := security.UnmarshalRsaCertificate(record.GetValue())
				if err != nil {
					return nil, fmt.Errorf("failed to parse EF_CA_CERTIFICATE (Gen1): %w", err)
				}
				cert := &cardv1.CaCertificate{}
				cert.SetRsaCertificate(rsaCert)
				// Per regulation (Chapter 12, Section 3.3), certificate EFs should NOT have signatures.
				// However, some real-world cards may incorrectly include one. We capture it for
				// data fidelity while noting it's non-compliant. It will not be written during marshalling.
				if signature != nil {
					cert.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					cert.SetAuthentication(auth)
				}
				tachographDF.SetCaCertificate(cert)
			case ddv1.Generation_GENERATION_2:
				if tachographG2DF == nil {
					tachographG2DF = &cardv1.DriverCardFile_TachographG2{}
				}
				eccCert, err := security.UnmarshalEccCertificate(record.GetValue())
				if err != nil {
					return nil, fmt.Errorf("failed to parse EF_CA_CERTIFICATE (Gen2): %w", err)
				}
				cert := &cardv1.CaCertificateG2{}
				cert.SetEccCertificate(eccCert)
				// Per regulation (Chapter 12, Section 3.3), certificate EFs should NOT have signatures.
				// However, some real-world cards may incorrectly include one. We capture it for
				// data fidelity while noting it's non-compliant. It will not be written during marshalling.
				if signature != nil {
					cert.SetSignature(signature)
				}
				// Propagate authentication
				if auth := record.GetAuthentication(); auth != nil {
					cert.SetAuthentication(auth)
				}
				tachographG2DF.SetCaCertificate(cert)
			default:
				return nil, fmt.Errorf("unexpected generation for EF_CA_CERTIFICATE: %v", efGeneration)
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
			eccCert, err := security.UnmarshalEccCertificate(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("failed to parse EF_LINK_CERTIFICATE: %w", err)
			}
			cert := &cardv1.LinkCertificate{}
			cert.SetEccCertificate(eccCert)
			// Signature is non-compliant per regulation but captured for data fidelity.
			if signature != nil {
				cert.SetSignature(signature)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				cert.SetAuthentication(auth)
			}
			tachographG2DF.SetLinkCertificate(cert)
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

	// Create default MarshalOptions for internal calls
	opts := MarshalOptions{}

	// 1. EF_ICC (0x0002) - no signature
	if icc := card.GetIcc(); icc != nil {
		dataBytes, err := opts.MarshalIcc(icc)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_ICC,
			dataBytes,
			nil, // no signature
			0x00)
		if err != nil {
			return nil, err
		}
	}

	// 2. EF_IC (0x0005) - no signature
	if ic := card.GetIc(); ic != nil {
		dataBytes, err := opts.MarshalCardIc(ic)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_IC,
			dataBytes,
			nil, // no signature
			0x00)
		if err != nil {
			return nil, err
		}
	}

	// 3. EF_APPLICATION_IDENTIFICATION (0x0501)
	if appId := card.GetTachograph().GetApplicationIdentification(); appId != nil {
		dataBytes, err := opts.MarshalCardApplicationIdentification(appId)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION,
			dataBytes,
			appId.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if drivingLicence := card.GetTachograph().GetDrivingLicenceInfo(); drivingLicence != nil {
		dataBytes, err := opts.MarshalDrivingLicenceInfo(drivingLicence)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO,
			dataBytes,
			drivingLicence.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	// 4. EF_IDENTIFICATION (0x0520)
	if identification := card.GetTachograph().GetIdentification(); identification != nil {
		dataBytes, err := opts.MarshalDriverCardIdentification(identification)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_IDENTIFICATION,
			dataBytes,
			identification.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if eventsData := card.GetTachograph().GetEventsData(); eventsData != nil {
		dataBytes, err := opts.MarshalEventsData(eventsData)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_EVENTS_DATA,
			dataBytes,
			eventsData.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if faultsData := card.GetTachograph().GetFaultsData(); faultsData != nil {
		dataBytes, err := opts.MarshalFaultsData(faultsData)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_FAULTS_DATA,
			dataBytes,
			faultsData.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if driverActivity := card.GetTachograph().GetDriverActivityData(); driverActivity != nil {
		dataBytes, err := opts.MarshalDriverActivity(driverActivity)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA,
			dataBytes,
			driverActivity.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if vehiclesUsed := card.GetTachograph().GetVehiclesUsed(); vehiclesUsed != nil {
		dataBytes, err := opts.MarshalVehiclesUsed(vehiclesUsed)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_VEHICLES_USED,
			dataBytes,
			vehiclesUsed.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if places := card.GetTachograph().GetPlaces(); places != nil {
		dataBytes, err := opts.MarshalPlaces(places)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_PLACES,
			dataBytes,
			places.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if currentUsage := card.GetTachograph().GetCurrentUsage(); currentUsage != nil {
		dataBytes, err := opts.MarshalCurrentUsage(currentUsage)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_CURRENT_USAGE,
			dataBytes,
			currentUsage.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if controlActivity := card.GetTachograph().GetControlActivityData(); controlActivity != nil {
		dataBytes, err := opts.MarshalCardControlActivityData(controlActivity)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA,
			dataBytes,
			controlActivity.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if specificConditions := card.GetTachograph().GetSpecificConditions(); specificConditions != nil {
		dataBytes, err := opts.MarshalCardSpecificConditions(specificConditions)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS,
			dataBytes,
			specificConditions.GetSignature(),
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	if cardDownload := card.GetTachograph().GetCardDownload(); cardDownload != nil {
		dataBytes, err := opts.MarshalCardDownload(cardDownload)
		if err != nil {
			return nil, err
		}
		dst, err = appendTlvBlock(dst,
			cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER,
			dataBytes,
			nil,  // no signature
			0x00) // Gen1
		if err != nil {
			return nil, err
		}
	}

	// Gen2 DF - marshal all Gen2 EFs with appendix 0x02/0x03
	if tachographG2 := card.GetTachographG2(); tachographG2 != nil {
		// Marshal Gen2 versions of shared EFs

		// ApplicationIdentification (Gen2)
		if appId := tachographG2.GetApplicationIdentification(); appId != nil {
			dataBytes, err := opts.MarshalCardApplicationIdentificationG2(appId)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION,
				dataBytes,
				appId.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
		}

		if vehiclesUsed := tachographG2.GetVehiclesUsed(); vehiclesUsed != nil {
			dataBytes, err := opts.MarshalVehiclesUsedG2(vehiclesUsed)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_VEHICLES_USED,
				dataBytes,
				vehiclesUsed.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
		}

		if places := tachographG2.GetPlaces(); places != nil {
			dataBytes, err := opts.MarshalPlacesG2(places)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_PLACES,
				dataBytes,
				places.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
		}

		// SpecificConditions (Gen2)
		if specificConditions := tachographG2.GetSpecificConditions(); specificConditions != nil {
			dataBytes, err := opts.MarshalCardSpecificConditionsG2(specificConditions)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS,
				dataBytes,
				specificConditions.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
		}

		// Marshal Gen2-exclusive EFs
		if vehicleUnitsUsed := tachographG2.GetVehicleUnitsUsed(); vehicleUnitsUsed != nil {
			dataBytes, err := opts.MarshalCardVehicleUnitsUsed(vehicleUnitsUsed)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED,
				dataBytes,
				vehicleUnitsUsed.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
		}

		if gnssPlaces := tachographG2.GetGnssPlaces(); gnssPlaces != nil {
			dataBytes, err := opts.MarshalCardGnssPlaces(gnssPlaces)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_GNSS_PLACES,
				dataBytes,
				gnssPlaces.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
		}

		if appIdV2 := tachographG2.GetApplicationIdentificationV2(); appIdV2 != nil {
			dataBytes, err := opts.MarshalCardApplicationIdentificationV2(appIdV2)
			if err != nil {
				return nil, err
			}
			dst, err = appendTlvBlock(dst,
				cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2,
				dataBytes,
				appIdV2.GetSignature(),
				0x02) // Gen2
			if err != nil {
				return nil, err
			}
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

// appendTlvBlock writes a TLV block (and optional signature block) to dst.
//
// Parameters:
//   - fileType: The Elementary File type (used to look up the FID tag)
//   - dataBytes: Pre-marshalled data bytes (nil = skip writing entirely)
//   - signature: Signature bytes (nil/empty = skip signature block)
//   - appendix: Tag appendix byte (0x00 for Gen1 data, 0x02 for Gen2 data)
//
// The function writes:
//  1. Data block: [FID:2][appendix:1][length:2][data:N]
//  2. Signature block (if signature present): [FID:2][appendix+1:1][length:2][signature:N]
func appendTlvBlock(
	dst []byte,
	fileType cardv1.ElementaryFileType,
	dataBytes []byte,
	signature []byte,
	appendix byte,
) ([]byte, error) {
	// Skip if no data to write
	if dataBytes == nil {
		return dst, nil
	}

	// Get FID tag from protobuf enum options
	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	// Write data block: [FID][appendix][length][value]
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, appendix)
	dst = binary.BigEndian.AppendUint16(dst, uint16(len(dataBytes)))
	dst = append(dst, dataBytes...)

	// Write signature block if present: [FID][appendix+1][length][signature]
	if len(signature) > 0 {
		sigAppendix := appendix + 1 // 0x01 for Gen1, 0x03 for Gen2
		dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
		dst = append(dst, sigAppendix)
		dst = binary.BigEndian.AppendUint16(dst, uint16(len(signature)))
		dst = append(dst, signature...)
	}

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

	var caCert *securityv1.RsaCertificate
	var err error

	if o.CertificateResolver != nil {
		// Use certificate resolver to fetch CA certificate
		car := cardCert.GetCertificateAuthorityReference()
		caCert, err = o.CertificateResolver.GetRsaCertificate(ctx, car)
		if err != nil {
			return fmt.Errorf("failed to fetch CA certificate from resolver: %w", err)
		}

		// For RSA certificates, the public key is extracted during signature recovery.
		// If the CA certificate doesn't have its public key yet, we need to verify it
		// against the root CA first to populate it.
		if len(caCert.GetRsaModulus()) == 0 || len(caCert.GetRsaExponent()) == 0 {
			// Fetch the root CA certificate
			rootCert, err := o.CertificateResolver.GetRootCertificate(ctx)
			if err != nil {
				return fmt.Errorf("failed to get root CA certificate: %w", err)
			}

			// Verify the CA certificate against the root CA to populate its public key
			if err := security.VerifyRsaCertificateWithRoot(caCert, rootCert); err != nil {
				return fmt.Errorf("CA certificate verification failed: %w", err)
			}
		}
	} else {
		// Fall back to embedded CA certificate from card file
		caCert = tachograph.GetCaCertificate().GetRsaCertificate()
		if caCert == nil {
			return fmt.Errorf("CA certificate is missing from card file")
		}
	}

	// Verify the card certificate using the CA certificate
	if err := security.VerifyRsaCertificateWithCA(cardCert, caCert); err != nil {
		return fmt.Errorf("card certificate verification failed: %w", err)
	}

	return nil
}

// verifyGen2Certificates verifies Generation 2 ECC certificates.
// If a certificate resolver is configured, it fetches CA certificates from the resolver.
// Otherwise, it uses the embedded CA certificate from the card file.
func (o VerifyOptions) verifyGen2Certificates(ctx context.Context, tachographG2 *cardv1.DriverCardFile_TachographG2) error {
	cardSignCert := tachographG2.GetCardSignCertificate().GetEccCertificate()

	if cardSignCert == nil {
		return fmt.Errorf("card sign certificate is missing")
	}

	var caCert *securityv1.EccCertificate
	var err error

	if o.CertificateResolver != nil {
		// Use certificate resolver to fetch CA certificate
		car := cardSignCert.GetCertificateAuthorityReference()
		caCert, err = o.CertificateResolver.GetEccCertificate(ctx, car)
		if err != nil {
			return fmt.Errorf("failed to fetch CA certificate from resolver: %w", err)
		}
	} else {
		// Fall back to embedded CA certificate from card file
		caCert = tachographG2.GetCaCertificate().GetEccCertificate()
		if caCert == nil {
			return fmt.Errorf("CA certificate is missing from card file")
		}
	}

	// Verify the card sign certificate using the CA certificate
	if err := security.VerifyEccCertificateWithCA(cardSignCert, caCert); err != nil {
		return fmt.Errorf("card sign certificate verification failed: %w", err)
	}

	return nil
}
