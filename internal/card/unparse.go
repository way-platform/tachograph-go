package card

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/security"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnparseDriverCardFile converts a parsed DriverCardFile back into its raw TLV representation.
// This is the inverse of ParseRawDriverCardFile.
func UnparseDriverCardFile(file *cardv1.DriverCardFile) (*cardv1.RawCardFile, error) {
	if file == nil {
		return nil, fmt.Errorf("driver card file cannot be nil")
	}

	var records []*cardv1.RawCardFile_Record
	marshalOpts := MarshalOptions{}

	// Helper to append a TLV record (data + optional signature)
	appendRecord := func(
		fileType cardv1.ElementaryFileType,
		generation ddv1.Generation,
		dataBytes []byte,
		signature []byte,
	) error {
		if dataBytes == nil {
			return nil
		}

		// Get FID for this file type
		fid, ok := getFileId(fileType)
		if !ok {
			return fmt.Errorf("no FID found for file type %v", fileType)
		}

		// Calculate appendix byte from generation
		// Bit 0: 0 = DATA, 1 = SIGNATURE
		// Bit 1: 0 = Gen1, 1 = Gen2
		var dataAppendix byte
		if generation == ddv1.Generation_GENERATION_2 {
			dataAppendix = 0x02 // Gen2 data
		} else {
			dataAppendix = 0x00 // Gen1 data
		}

		// Create data record
		dataTag := (int32(fid) << 8) | int32(dataAppendix)
		dataRecord := &cardv1.RawCardFile_Record{}
		dataRecord.SetTag(dataTag)
		dataRecord.SetFile(fileType)
		dataRecord.SetGeneration(generation)
		dataRecord.SetContentType(cardv1.ContentType_DATA)
		dataRecord.SetLength(int32(len(dataBytes)))
		dataRecord.SetValue(dataBytes)
		records = append(records, dataRecord)

		// Create signature record if present
		if len(signature) > 0 {
			sigAppendix := dataAppendix + 1 // 0x01 for Gen1, 0x03 for Gen2
			sigTag := (int32(fid) << 8) | int32(sigAppendix)
			sigRecord := &cardv1.RawCardFile_Record{}
			sigRecord.SetTag(sigTag)
			sigRecord.SetFile(fileType)
			sigRecord.SetGeneration(generation)
			sigRecord.SetContentType(cardv1.ContentType_SIGNATURE)
			sigRecord.SetLength(int32(len(signature)))
			sigRecord.SetValue(signature)
			records = append(records, sigRecord)
		}

		return nil
	}

	// 1. EF_ICC - common file
	if icc := file.GetIcc(); icc != nil {
		dataBytes, err := marshalOpts.MarshalIcc(icc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal EF_ICC: %w", err)
		}
		if err := appendRecord(cardv1.ElementaryFileType_EF_ICC, ddv1.Generation_GENERATION_1, dataBytes, nil); err != nil {
			return nil, err
		}
	}

	// 2. EF_IC - common file
	if ic := file.GetIc(); ic != nil {
		dataBytes, err := marshalOpts.MarshalCardIc(ic)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal EF_IC: %w", err)
		}
		if err := appendRecord(cardv1.ElementaryFileType_EF_IC, ddv1.Generation_GENERATION_1, dataBytes, nil); err != nil {
			return nil, err
		}
	}

	// 3. Gen1 DF (Tachograph)
	if tachograph := file.GetTachograph(); tachograph != nil {
		// EF_APPLICATION_IDENTIFICATION
		if appId := tachograph.GetApplicationIdentification(); appId != nil {
			dataBytes, err := marshalOpts.MarshalCardApplicationIdentification(appId)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_APPLICATION_IDENTIFICATION (Gen1): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, ddv1.Generation_GENERATION_1, dataBytes, appId.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CARD_CERTIFICATE
		if cert := tachograph.GetCardCertificate(); cert != nil {
			dataBytes, err := security.MarshalRsaCertificate(cert.GetRsaCertificate())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CARD_CERTIFICATE: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CARD_CERTIFICATE, ddv1.Generation_GENERATION_1, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_CA_CERTIFICATE (Gen1)
		if cert := tachograph.GetCaCertificate(); cert != nil {
			dataBytes, err := security.MarshalRsaCertificate(cert.GetRsaCertificate())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CA_CERTIFICATE (Gen1): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CA_CERTIFICATE, ddv1.Generation_GENERATION_1, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_IDENTIFICATION
		if identification := tachograph.GetIdentification(); identification != nil {
			dataBytes, err := marshalOpts.MarshalDriverCardIdentification(identification)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal DriverCardIdentification: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_IDENTIFICATION, ddv1.Generation_GENERATION_1, dataBytes, identification.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CARD_DOWNLOAD
		if cardDownload := tachograph.GetCardDownload(); cardDownload != nil {
			dataBytes, err := marshalOpts.MarshalCardDownload(cardDownload)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CARD_DOWNLOAD: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, ddv1.Generation_GENERATION_1, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_DRIVING_LICENCE_INFO
		if dli := tachograph.GetDrivingLicenceInfo(); dli != nil {
			dataBytes, err := marshalOpts.MarshalDrivingLicenceInfo(dli)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_DRIVING_LICENCE_INFO: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO, ddv1.Generation_GENERATION_1, dataBytes, dli.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_EVENTS_DATA
		if events := tachograph.GetEventsData(); events != nil {
			dataBytes, err := marshalOpts.MarshalEventsData(events)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_EVENTS_DATA: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_EVENTS_DATA, ddv1.Generation_GENERATION_1, dataBytes, events.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_FAULTS_DATA
		if faults := tachograph.GetFaultsData(); faults != nil {
			dataBytes, err := marshalOpts.MarshalFaultsData(faults)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_FAULTS_DATA: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_FAULTS_DATA, ddv1.Generation_GENERATION_1, dataBytes, faults.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_DRIVER_ACTIVITY_DATA
		if activity := tachograph.GetDriverActivityData(); activity != nil {
			dataBytes, err := marshalOpts.MarshalDriverActivity(activity)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_DRIVER_ACTIVITY_DATA: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, ddv1.Generation_GENERATION_1, dataBytes, activity.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_VEHICLES_USED
		if vehicles := tachograph.GetVehiclesUsed(); vehicles != nil {
			dataBytes, err := marshalOpts.MarshalVehiclesUsed(vehicles)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_VEHICLES_USED: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_VEHICLES_USED, ddv1.Generation_GENERATION_1, dataBytes, vehicles.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_PLACES
		if places := tachograph.GetPlaces(); places != nil {
			dataBytes, err := marshalOpts.MarshalPlaces(places)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_PLACES: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_PLACES, ddv1.Generation_GENERATION_1, dataBytes, places.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CURRENT_USAGE
		if currentUsage := tachograph.GetCurrentUsage(); currentUsage != nil {
			dataBytes, err := marshalOpts.MarshalCurrentUsage(currentUsage)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CURRENT_USAGE: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CURRENT_USAGE, ddv1.Generation_GENERATION_1, dataBytes, currentUsage.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CONTROL_ACTIVITY_DATA
		if controlActivity := tachograph.GetControlActivityData(); controlActivity != nil {
			dataBytes, err := marshalOpts.MarshalCardControlActivityData(controlActivity)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CONTROL_ACTIVITY_DATA: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, ddv1.Generation_GENERATION_1, dataBytes, controlActivity.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_SPECIFIC_CONDITIONS
		if specificConditions := tachograph.GetSpecificConditions(); specificConditions != nil {
			dataBytes, err := marshalOpts.MarshalCardSpecificConditions(specificConditions)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_SPECIFIC_CONDITIONS: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, ddv1.Generation_GENERATION_1, dataBytes, specificConditions.GetSignature()); err != nil {
				return nil, err
			}
		}
	}

	// 4. Gen2 DF (Tachograph_G2)
	if tachographG2 := file.GetTachographG2(); tachographG2 != nil {
		// EF_APPLICATION_IDENTIFICATION (Gen2)
		if appIdG2 := tachographG2.GetApplicationIdentification(); appIdG2 != nil {
			dataBytes, err := marshalOpts.MarshalCardApplicationIdentificationG2(appIdG2)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_APPLICATION_IDENTIFICATION (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, ddv1.Generation_GENERATION_2, dataBytes, appIdG2.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CARD_MA_CERTIFICATE
		if cert := tachographG2.GetCardMaCertificate(); cert != nil {
			dataBytes, err := security.MarshalEccCertificate(cert.GetEccCertificate())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CARD_MA_CERTIFICATE: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CARD_MA_CERTIFICATE, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_CARD_SIGN_CERTIFICATE
		if cert := tachographG2.GetCardSignCertificate(); cert != nil {
			dataBytes, err := security.MarshalEccCertificate(cert.GetEccCertificate())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CARD_SIGN_CERTIFICATE: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CARD_SIGN_CERTIFICATE, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_CA_CERTIFICATE (Gen2)
		if cert := tachographG2.GetCaCertificate(); cert != nil {
			dataBytes, err := security.MarshalEccCertificate(cert.GetEccCertificate())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CA_CERTIFICATE (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CA_CERTIFICATE, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_LINK_CERTIFICATE
		if cert := tachographG2.GetLinkCertificate(); cert != nil {
			dataBytes, err := security.MarshalEccCertificate(cert.GetEccCertificate())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_LINK_CERTIFICATE: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_LINK_CERTIFICATE, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_IDENTIFICATION (Gen2)
		if identification := tachographG2.GetIdentification(); identification != nil {
			dataBytes, err := marshalOpts.MarshalDriverCardIdentification(identification)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal DriverCardIdentification (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_IDENTIFICATION, ddv1.Generation_GENERATION_2, dataBytes, identification.GetSignature()); err != nil {
				return nil, err
			}
		}

		// Continue with other Gen2 EFs...
		// EF_CARD_DOWNLOAD
		if cardDownload := tachographG2.GetCardDownload(); cardDownload != nil {
			dataBytes, err := marshalOpts.MarshalCardDownload(cardDownload)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CARD_DOWNLOAD (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_DRIVING_LICENCE_INFO (Gen2)
		if dli := tachographG2.GetDrivingLicenceInfo(); dli != nil {
			dataBytes, err := marshalOpts.MarshalDrivingLicenceInfo(dli)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_DRIVING_LICENCE_INFO (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO, ddv1.Generation_GENERATION_2, dataBytes, dli.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_EVENTS_DATA (Gen2)
		if events := tachographG2.GetEventsData(); events != nil {
			dataBytes, err := marshalOpts.MarshalEventsData(events)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_EVENTS_DATA (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_EVENTS_DATA, ddv1.Generation_GENERATION_2, dataBytes, events.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_FAULTS_DATA (Gen2)
		if faults := tachographG2.GetFaultsData(); faults != nil {
			dataBytes, err := marshalOpts.MarshalFaultsData(faults)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_FAULTS_DATA (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_FAULTS_DATA, ddv1.Generation_GENERATION_2, dataBytes, faults.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_DRIVER_ACTIVITY_DATA (Gen2)
		if activity := tachographG2.GetDriverActivityData(); activity != nil {
			dataBytes, err := marshalOpts.MarshalDriverActivity(activity)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_DRIVER_ACTIVITY_DATA (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, ddv1.Generation_GENERATION_2, dataBytes, activity.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_VEHICLES_USED (Gen2)
		if vehicles := tachographG2.GetVehiclesUsed(); vehicles != nil {
			dataBytes, err := marshalOpts.MarshalVehiclesUsedG2(vehicles)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_VEHICLES_USED (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_VEHICLES_USED, ddv1.Generation_GENERATION_2, dataBytes, vehicles.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_PLACES (Gen2)
		if places := tachographG2.GetPlaces(); places != nil {
			dataBytes, err := marshalOpts.MarshalPlacesG2(places)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_PLACES (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_PLACES, ddv1.Generation_GENERATION_2, dataBytes, places.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CURRENT_USAGE (Gen2)
		if currentUsage := tachographG2.GetCurrentUsage(); currentUsage != nil {
			dataBytes, err := marshalOpts.MarshalCurrentUsage(currentUsage)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CURRENT_USAGE (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CURRENT_USAGE, ddv1.Generation_GENERATION_2, dataBytes, currentUsage.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_CONTROL_ACTIVITY_DATA (Gen2)
		if controlActivity := tachographG2.GetControlActivityData(); controlActivity != nil {
			dataBytes, err := marshalOpts.MarshalCardControlActivityData(controlActivity)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_CONTROL_ACTIVITY_DATA (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, ddv1.Generation_GENERATION_2, dataBytes, controlActivity.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_SPECIFIC_CONDITIONS (Gen2)
		if specificConditions := tachographG2.GetSpecificConditions(); specificConditions != nil {
			dataBytes, err := marshalOpts.MarshalCardSpecificConditionsG2(specificConditions)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_SPECIFIC_CONDITIONS (Gen2): %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, ddv1.Generation_GENERATION_2, dataBytes, specificConditions.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_VEHICLE_UNITS_USED (Gen2 only)
		if vehicleUnitsUsed := tachographG2.GetVehicleUnitsUsed(); vehicleUnitsUsed != nil {
			dataBytes, err := marshalOpts.MarshalCardVehicleUnitsUsed(vehicleUnitsUsed)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_VEHICLE_UNITS_USED: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED, ddv1.Generation_GENERATION_2, dataBytes, vehicleUnitsUsed.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_GNSS_PLACES (Gen2 only)
		if gnssPlaces := tachographG2.GetGnssPlaces(); gnssPlaces != nil {
			dataBytes, err := marshalOpts.MarshalCardGnssPlaces(gnssPlaces)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_GNSS_PLACES: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_GNSS_PLACES, ddv1.Generation_GENERATION_2, dataBytes, gnssPlaces.GetSignature()); err != nil {
				return nil, err
			}
		}

		// EF_BORDER_CROSSINGS (Gen2 only)
		if borderCrossings := tachographG2.GetBorderCrossings(); borderCrossings != nil {
			dataBytes, err := marshalOpts.MarshalCardBorderCrossings(borderCrossings)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_BORDER_CROSSINGS: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_BORDER_CROSSINGS, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_PLACES_AUTHENTICATION (Gen2 only)
		if placesAuthentication := tachographG2.GetPlacesAuthentication(); placesAuthentication != nil {
			dataBytes, err := marshalOpts.MarshalPlacesAuthentication(placesAuthentication)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_PLACES_AUTHENTICATION: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_PLACES_AUTHENTICATION, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_GNSS_PLACES_AUTHENTICATION (Gen2 only)
		if gnssPlacesAuthentication := tachographG2.GetGnssPlacesAuthentication(); gnssPlacesAuthentication != nil {
			dataBytes, err := marshalOpts.MarshalGnssPlacesAuthentication(gnssPlacesAuthentication)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_GNSS_PLACES_AUTHENTICATION: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_GNSS_PLACES_AUTHENTICATION, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_LOAD_UNLOAD_OPERATIONS (Gen2 only)
		if loadUnloadOperations := tachographG2.GetLoadUnloadOperations(); loadUnloadOperations != nil {
			dataBytes, err := marshalOpts.MarshalLoadUnloadOperations(loadUnloadOperations)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_LOAD_UNLOAD_OPERATIONS: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_LOAD_UNLOAD_OPERATIONS, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_LOAD_TYPE_ENTRIES (Gen2 only)
		if loadTypeEntries := tachographG2.GetLoadTypeEntries(); loadTypeEntries != nil {
			dataBytes, err := marshalOpts.MarshalLoadTypeEntries(loadTypeEntries)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_LOAD_TYPE_ENTRIES: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_LOAD_TYPE_ENTRIES, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}

		// EF_COMPANY_ACTIVITY_DATA (Gen2 only)
		if companyActivityData := tachographG2.GetCompanyActivityData(); companyActivityData != nil {
			dataBytes, err := marshalOpts.MarshalCompanyActivityData(companyActivityData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EF_COMPANY_ACTIVITY_DATA: %w", err)
			}
			if err := appendRecord(cardv1.ElementaryFileType_EF_COMPANY_ACTIVITY_DATA, ddv1.Generation_GENERATION_2, dataBytes, nil); err != nil {
				return nil, err
			}
		}
	}

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords(records)
	return rawFile, nil
}
