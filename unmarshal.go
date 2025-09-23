package tachograph

import (
	"bytes"
	"encoding/binary"
	"errors"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalFile parses a .DDD file's byte data into a protobuf File message.
func UnmarshalFile(data []byte) (*tachographv1.File, error) {
	fileType := InferFileType(data)
	switch fileType {
	case CardFileType:
		return unmarshalCard(data)
	case UnitFileType:
		return unmarshalVU(data)
	}
	return nil, errors.New("unknown or unsupported file type")
}

func unmarshalCard(data []byte) (*tachographv1.File, error) {
	file := &tachographv1.File{}
	file.SetType(tachographv1.File_DRIVER_CARD) // Assume Driver card for now
	file.SetDriverCard(&cardv1.DriverCardFile{})

	// Track any truly proprietary EFs (none expected for standard cards)

	r := bytes.NewReader(data)

	for r.Len() > 0 {
		// Read Tag - 3 bytes (FID + appendix) according to DDD format spec
		tagBytes := make([]byte, 3)
		if _, err := r.Read(tagBytes); err != nil {
			return nil, err
		}
		// Extract FID (first 2 bytes) and appendix (last byte)
		fid := binary.BigEndian.Uint16(tagBytes[0:2])
		appendix := tagBytes[2]

		// Read Length - 2 bytes
		var length uint16
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		value := make([]byte, length)
		if _, err := r.Read(value); err != nil {
			return nil, err
		}

		// Skip signatures (appendix '01' and '03') - we only process data (appendix '00' and '02')
		if appendix == 0x01 || appendix == 0x03 {
			continue // Skip signature TLV objects
		}

		// Find file type from FID and dispatch
		fileType := findFileTypeByTag(int32(fid))

		// Debug: log all FIDs being processed
		// if fid >= 0xC000 {
		//	fmt.Printf("DEBUG: Processing FID 0x%04X, fileType=%v, appendix=0x%02X\n", fid, fileType, appendix)
		// }

		if fileType == cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
			// Skip unknown EF types (true proprietary EFs would be handled here)
			continue
		}

		driverCard := file.GetDriverCard()
		switch fileType {
		case cardv1.ElementaryFileType_EF_ICC:
			icc := &cardv1.IccIdentification{}
			if err := UnmarshalIcc(value, icc); err != nil {
				return nil, err
			}
			driverCard.SetIcc(icc)
		case cardv1.ElementaryFileType_EF_IC:
			ic := &cardv1.ChipIdentification{}
			if err := UnmarshalCardIc(value, ic); err != nil {
				return nil, err
			}
			driverCard.SetIc(ic)
		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification := &cardv1.CardIdentification{}
			holderIdentification := &cardv1.DriverCardHolderIdentification{}
			if err := UnmarshalIdentification(value, identification, holderIdentification); err != nil {
				return nil, err
			}
			driverCard.SetIdentification(identification)
			driverCard.SetHolderIdentification(holderIdentification)
		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo := &cardv1.DrivingLicenceInfo{}
			if err := UnmarshalDrivingLicenceInfo(value, drivingLicenceInfo); err != nil {
				return nil, err
			}
			driverCard.SetDrivingLicenceInfo(drivingLicenceInfo)
		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData := &cardv1.EventData{}
			if err := UnmarshalEventsData(value, eventsData); err != nil {
				return nil, err
			}
			driverCard.SetEventsData(eventsData)
		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData := &cardv1.FaultData{}
			if err := UnmarshalFaultsData(value, faultsData); err != nil {
				return nil, err
			}
			driverCard.SetFaultsData(faultsData)
		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData := &cardv1.DriverActivity{}
			if err := UnmarshalCardActivityData(value, activityData); err != nil {
				return nil, err
			}
			driverCard.SetDriverActivityData(activityData)
		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed := &cardv1.VehiclesUsed{}
			if err := UnmarshalCardVehiclesUsed(value, vehiclesUsed); err != nil {
				return nil, err
			}
			driverCard.SetVehiclesUsed(vehiclesUsed)
		case cardv1.ElementaryFileType_EF_PLACES:
			places := &cardv1.Places{}
			if err := UnmarshalCardPlaces(value, places); err != nil {
				return nil, err
			}
			driverCard.SetPlaces(places)
		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage := &cardv1.CurrentUsage{}
			if err := UnmarshalCardCurrentUsage(value, currentUsage); err != nil {
				return nil, err
			}
			driverCard.SetCurrentUsage(currentUsage)
		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId := &cardv1.DriverCardApplicationIdentification{}
			if err := UnmarshalCardApplicationIdentification(value, appId); err != nil {
				return nil, err
			}
			driverCard.SetApplicationIdentification(appId)
		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity := &cardv1.ControlActivityData{}
			if err := UnmarshalCardControlActivityData(value, controlActivity); err != nil {
				return nil, err
			}
			driverCard.SetControlActivityData(controlActivity)
		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions := &cardv1.SpecificConditions{}
			if err := UnmarshalCardSpecificConditions(value, specificConditions); err != nil {
				return nil, err
			}
			driverCard.SetSpecificConditions(specificConditions)
		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload := &cardv1.LastCardDownload{}
			if err := UnmarshalCardLastDownload(value, lastDownload); err != nil {
				return nil, err
			}
			driverCard.SetLastCardDownload(lastDownload)
		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits := &cardv1.VehicleUnitsUsed{}
			if err := UnmarshalCardVehicleUnitsUsed(value, vehicleUnits); err != nil {
				return nil, err
			}
			driverCard.SetVehicleUnitsUsed(vehicleUnits)
		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces := &cardv1.GnssPlaces{}
			if err := UnmarshalCardGnssPlaces(value, gnssPlaces); err != nil {
				return nil, err
			}
			driverCard.SetGnssPlaces(gnssPlaces)
		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2 := &cardv1.ApplicationIdentificationV2{}
			if err := UnmarshalCardApplicationIdentificationV2(value, appIdV2); err != nil {
				return nil, err
			}
			driverCard.SetApplicationIdentificationV2(appIdV2)
		case cardv1.ElementaryFileType_EF_CARD_CERTIFICATE:
			// Initialize certificates if needed
			if driverCard.GetCertificates() == nil {
				driverCard.SetCertificates(&cardv1.Certificates{})
			}
			driverCard.GetCertificates().SetCardCertificate(value)
		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			// Initialize certificates if needed
			if driverCard.GetCertificates() == nil {
				driverCard.SetCertificates(&cardv1.Certificates{})
			}
			driverCard.GetCertificates().SetCaCertificate(value)
		}
	}

	// Note: Any proprietary EFs would be stored here if needed

	return file, nil
}

func unmarshalVU(data []byte) (*tachographv1.File, error) {
	file := &tachographv1.File{}
	file.SetType(tachographv1.File_VEHICLE_UNIT)
	file.SetVehicleUnit(&vuv1.VehicleUnitFile{})

	r := bytes.NewReader(data)
	vuFile := file.GetVehicleUnit()

	for r.Len() > 1 { // Need at least 2 bytes for tag
		// Read Tag - 2 bytes for VU files (TV format)
		var tag uint16
		if err := binary.Read(r, binary.BigEndian, &tag); err != nil {
			return nil, err
		}

		// Determine transfer type from tag
		transferType := findTransferTypeByTag(tag)
		if transferType == vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED {
			// Skip unknown tags - we need to determine how much data to skip
			// For now, we'll break out of the loop on unknown tags
			break
		}

		// Parse the transfer data based on type
		transfer := &vuv1.VehicleUnitFile_Transfer{}
		transfer.SetType(transferType)

		// Parse the specific data type - this will determine how much data to consume
		switch transferType {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			version := &vuv1.DownloadInterfaceVersion{}
			_, err := UnmarshalDownloadInterfaceVersion(r, version)
			if err != nil {
				return nil, err
			}
			transfer.SetDownloadInterfaceVersion(version)
			// Move reader position
			// Reader position is already advanced by the unmarshal function
		case vuv1.TransferType_OVERVIEW_GEN1:
			overview := &vuv1.Overview{}
			_, err := UnmarshalOverview(r, overview, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetOverview(overview)
			// Reader position is already advanced by the unmarshal function
		case vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			overview := &vuv1.Overview{}
			generation := 2
			_, err := UnmarshalOverview(r, overview, generation)
			if err != nil {
				return nil, err
			}
			transfer.SetOverview(overview)
			// Reader position is already advanced by the unmarshal function
		case vuv1.TransferType_ACTIVITIES_GEN1:
			activities := &vuv1.Activities{}
			_, err := UnmarshalVuActivities(r, activities, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			activities := &vuv1.Activities{}
			generation := 2
			_, err := UnmarshalVuActivities(r, activities, generation)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
			eventsAndFaults := &vuv1.EventsAndFaults{}
			_, err := UnmarshalVuEventsAndFaults(r, eventsAndFaults, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetEventsAndFaults(eventsAndFaults)
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			eventsAndFaults := &vuv1.EventsAndFaults{}
			_, err := UnmarshalVuEventsAndFaults(r, eventsAndFaults, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetEventsAndFaults(eventsAndFaults)
		case vuv1.TransferType_DETAILED_SPEED_GEN1:
			detailedSpeed := &vuv1.DetailedSpeed{}
			_, err := UnmarshalVuDetailedSpeed(r, detailedSpeed, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetDetailedSpeed(detailedSpeed)
		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed := &vuv1.DetailedSpeed{}
			_, err := UnmarshalVuDetailedSpeed(r, detailedSpeed, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetDetailedSpeed(detailedSpeed)
		case vuv1.TransferType_TECHNICAL_DATA_GEN1:
			technicalData := &vuv1.TechnicalData{}
			_, err := UnmarshalVuTechnicalData(r, technicalData, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetTechnicalData(technicalData)
		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			technicalData := &vuv1.TechnicalData{}
			_, err := UnmarshalVuTechnicalData(r, technicalData, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetTechnicalData(technicalData)
		default:
			// For now, skip unknown transfer types
			break
		}

		vuFile.SetTransfers(append(vuFile.GetTransfers(), transfer))
	}

	return file, nil
}

func findFileTypeByTag(tag int32) cardv1.ElementaryFileType {
	values := cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, cardv1.E_FileId) {
			if proto.GetExtension(opts, cardv1.E_FileId).(int32) == tag {
				return cardv1.ElementaryFileType(valueDesc.Number())
			}
		}
	}
	return cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED
}

func findTransferTypeByTag(tag uint16) vuv1.TransferType {
	values := vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, vuv1.E_TrepValue) {
			trepValue := proto.GetExtension(opts, vuv1.E_TrepValue).(int32)
			// VU tags are constructed as 0x76XX where XX is the TREP value
			expectedTag := uint16(0x7600 | (uint16(trepValue) & 0xFF))
			if expectedTag == tag {
				return vuv1.TransferType(valueDesc.Number())
			}
		}
	}
	return vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED
}
