package tachograph

import (
	"bufio"
	"bytes"
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

func unmarshalRawCardRecord(data []byte) (*cardv1.RawCardFile_Record, error) {
	var output cardv1.RawCardFile_Record
	output.SetTag(int32(binary.BigEndian.Uint16(data[:2])))
	output.SetFile(getElementaryFileTypeFromTag(int32(binary.BigEndian.Uint16(data[2:4]))))
	output.SetGeneration(cardv1.ApplicationGeneration_GENERATION_1)
	output.SetLength(int32(binary.BigEndian.Uint16(data[4:6])))
	output.SetValue(data[6:])
	return &output, nil
}

func unmarshalRawCardFile(input []byte) (*cardv1.RawCardFile, error) {
	var output cardv1.RawCardFile
	sc := bufio.NewScanner(bytes.NewReader(input))
	sc.Split(ScanCardFile)
	for sc.Scan() {
		record, err := unmarshalRawCardRecord(sc.Bytes())
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

func unmarshalCardFile(data []byte) (*tachographv1.File, error) {
	rawCardFile, err := unmarshalRawCardFile(data)
	if err != nil {
		return nil, err
	}
	// Pass 2: Process each data record from RawCardFile and parse into protobuf messages
	for _, record := range rawCardFile.GetRecords() {
		// Only process data records (signatures are preserved in RawCardFile for marshalling)
		if record.GetContentType() != cardv1.ContentType_DATA {
			continue
		}

		// Extract FID from tag (remove appendix) and find file type
		fid := record.GetTag() >> 8 // Remove appendix byte
		fileType := findFileTypeByTag(fid)

		// Debug: log all FIDs being processed
		// fmt.Printf("DEBUG: Processing tag 0x%06X, FID 0x%04X, fileType=%v\n",
		//	record.GetTag(), fid, fileType)

		if fileType == cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
			// Skip unknown EF types (true proprietary EFs would be handled here)
			continue
		}

		driverCard := file.GetDriverCard()
		switch fileType {
		case cardv1.ElementaryFileType_EF_ICC:
			icc := &cardv1.IccIdentification{}
			if err := UnmarshalIcc(record.GetValue(), icc); err != nil {
				return nil, err
			}
			driverCard.SetIcc(icc)
		case cardv1.ElementaryFileType_EF_IC:
			ic := &cardv1.ChipIdentification{}
			if err := UnmarshalCardIc(record.GetValue(), ic); err != nil {
				return nil, err
			}
			driverCard.SetIc(ic)
		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification := &cardv1.CardIdentification{}
			holderIdentification := &cardv1.DriverCardHolderIdentification{}
			if err := UnmarshalIdentification(record.GetValue(), identification, holderIdentification); err != nil {
				return nil, err
			}
			driverCard.SetIdentification(identification)
			driverCard.SetHolderIdentification(holderIdentification)
		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo := &cardv1.DrivingLicenceInfo{}
			if err := UnmarshalDrivingLicenceInfo(record.GetValue(), drivingLicenceInfo); err != nil {
				return nil, err
			}
			driverCard.SetDrivingLicenceInfo(drivingLicenceInfo)
		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData := &cardv1.EventData{}
			if err := UnmarshalEventsData(record.GetValue(), eventsData); err != nil {
				return nil, err
			}
			driverCard.SetEventsData(eventsData)
		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData := &cardv1.FaultData{}
			if err := UnmarshalFaultsData(record.GetValue(), faultsData); err != nil {
				return nil, err
			}
			driverCard.SetFaultsData(faultsData)
		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData := &cardv1.DriverActivity{}
			if err := UnmarshalCardActivityData(record.GetValue(), activityData); err != nil {
				return nil, err
			}
			driverCard.SetDriverActivityData(activityData)
		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed := &cardv1.VehiclesUsed{}
			if err := UnmarshalCardVehiclesUsed(record.GetValue(), vehiclesUsed); err != nil {
				return nil, err
			}
			driverCard.SetVehiclesUsed(vehiclesUsed)
		case cardv1.ElementaryFileType_EF_PLACES:
			places := &cardv1.Places{}
			if err := UnmarshalCardPlaces(record.GetValue(), places); err != nil {
				return nil, err
			}
			driverCard.SetPlaces(places)
		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage := &cardv1.CurrentUsage{}
			if err := UnmarshalCardCurrentUsage(record.GetValue(), currentUsage); err != nil {
				return nil, err
			}
			driverCard.SetCurrentUsage(currentUsage)
		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId := &cardv1.DriverCardApplicationIdentification{}
			if err := UnmarshalCardApplicationIdentification(record.GetValue(), appId); err != nil {
				return nil, err
			}
			driverCard.SetApplicationIdentification(appId)
		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity := &cardv1.ControlActivityData{}
			if err := UnmarshalCardControlActivityData(record.GetValue(), controlActivity); err != nil {
				return nil, err
			}
			driverCard.SetControlActivityData(controlActivity)
		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions := &cardv1.SpecificConditions{}
			if err := UnmarshalCardSpecificConditions(record.GetValue(), specificConditions); err != nil {
				return nil, err
			}
			driverCard.SetSpecificConditions(specificConditions)
		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload := &cardv1.LastCardDownload{}
			if err := UnmarshalCardLastDownload(record.GetValue(), lastDownload); err != nil {
				return nil, err
			}
			driverCard.SetLastCardDownload(lastDownload)
		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits := &cardv1.VehicleUnitsUsed{}
			if err := UnmarshalCardVehicleUnitsUsed(record.GetValue(), vehicleUnits); err != nil {
				return nil, err
			}
			driverCard.SetVehicleUnitsUsed(vehicleUnits)
		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces := &cardv1.GnssPlaces{}
			if err := UnmarshalCardGnssPlaces(record.GetValue(), gnssPlaces); err != nil {
				return nil, err
			}
			driverCard.SetGnssPlaces(gnssPlaces)
		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2 := &cardv1.ApplicationIdentificationV2{}
			if err := UnmarshalCardApplicationIdentificationV2(record.GetValue(), appIdV2); err != nil {
				return nil, err
			}
			driverCard.SetApplicationIdentificationV2(appIdV2)
		case cardv1.ElementaryFileType_EF_CARD_CERTIFICATE:
			// Initialize certificates if needed
			if driverCard.GetCertificates() == nil {
				driverCard.SetCertificates(&cardv1.Certificates{})
			}
			driverCard.GetCertificates().SetCardCertificate(record.GetValue())
		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			// Initialize certificates if needed
			if driverCard.GetCertificates() == nil {
				driverCard.SetCertificates(&cardv1.Certificates{})
			}
			driverCard.GetCertificates().SetCaCertificate(record.GetValue())
		}
	}

	// Note: Any proprietary EFs would be stored here if needed

	return file, nil
}
