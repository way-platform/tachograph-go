package tachograph

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// TLVRecord represents a parsed TLV record with optional signature
type TLVRecord struct {
	FID       uint16
	Appendix  uint8
	Value     []byte
	Signature []byte // Optional signature if this is a signed EF
}

// ScanTLV is a split function for bufio.Scanner that splits TLV records
func ScanTLV(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Need at least 5 bytes for TLV header (3 bytes tag + 2 bytes length)
	if len(data) < 5 {
		if atEOF {
			return 0, nil, errors.New("incomplete TLV header")
		}
		return 0, nil, nil // Need more data
	}

	// Parse TLV header
	length := binary.BigEndian.Uint16(data[3:5])
	totalLength := 5 + int(length) // 5 bytes header + value length

	if len(data) < totalLength {
		if atEOF {
			return 0, nil, errors.New("incomplete TLV record")
		}
		return 0, nil, nil // Need more data
	}

	// Return the complete TLV record
	return totalLength, data[:totalLength], nil
}

// ParseTLVRecords parses binary data into TLV records, automatically pairing data records with their signatures
func ParseTLVRecords(data []byte) ([]*TLVRecord, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(ScanTLV)

	var records []*TLVRecord
	var pendingDataRecord *TLVRecord

	for scanner.Scan() {
		tlvData := scanner.Bytes()
		if len(tlvData) < 5 {
			continue
		}

		// Parse TLV header
		fid := binary.BigEndian.Uint16(tlvData[0:2])
		appendix := tlvData[2]
		length := binary.BigEndian.Uint16(tlvData[3:5])
		value := tlvData[5 : 5+length]

		record := &TLVRecord{
			FID:      fid,
			Appendix: appendix,
			Value:    make([]byte, len(value)),
		}
		copy(record.Value, value)

		// Handle signature pairing
		if appendix == 0x01 || appendix == 0x03 {
			// This is a signature record
			if pendingDataRecord != nil && pendingDataRecord.FID == fid {
				// Pair with the pending data record
				pendingDataRecord.Signature = record.Value
				records = append(records, pendingDataRecord)
				pendingDataRecord = nil
			}
			// If no matching data record, skip this signature
		} else {
			// This is a data record (appendix 0x00 or 0x02)
			if pendingDataRecord != nil {
				// Add the previous unpaired data record without signature
				records = append(records, pendingDataRecord)
			}
			pendingDataRecord = record
		}
	}

	// Add any remaining unpaired data record
	if pendingDataRecord != nil {
		records = append(records, pendingDataRecord)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

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

	// Parse TLV records with automatic signature pairing
	tlvRecords, err := ParseTLVRecords(data)
	if err != nil {
		return nil, err
	}

	// Process each TLV record
	for _, record := range tlvRecords {
		// Find file type from FID and dispatch
		fileType := findFileTypeByTag(int32(record.FID))

		// Debug: log all FIDs being processed
		// if record.FID >= 0xC000 {
		//	fmt.Printf("DEBUG: Processing FID 0x%04X, fileType=%v, appendix=0x%02X, hasSignature=%v\n",
		//		record.FID, fileType, record.Appendix, len(record.Signature) > 0)
		// }

		if fileType == cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
			// Skip unknown EF types (true proprietary EFs would be handled here)
			continue
		}

		driverCard := file.GetDriverCard()
		switch fileType {
		case cardv1.ElementaryFileType_EF_ICC:
			icc := &cardv1.IccIdentification{}
			if err := UnmarshalIcc(record.Value, icc); err != nil {
				return nil, err
			}
			driverCard.SetIcc(icc)
		case cardv1.ElementaryFileType_EF_IC:
			ic := &cardv1.ChipIdentification{}
			if err := UnmarshalCardIc(record.Value, ic); err != nil {
				return nil, err
			}
			driverCard.SetIc(ic)
		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification := &cardv1.CardIdentification{}
			holderIdentification := &cardv1.DriverCardHolderIdentification{}
			if err := UnmarshalIdentification(record.Value, identification, holderIdentification); err != nil {
				return nil, err
			}
			driverCard.SetIdentification(identification)
			driverCard.SetHolderIdentification(holderIdentification)
		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo := &cardv1.DrivingLicenceInfo{}
			if err := UnmarshalDrivingLicenceInfo(record.Value, drivingLicenceInfo); err != nil {
				return nil, err
			}
			driverCard.SetDrivingLicenceInfo(drivingLicenceInfo)
		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData := &cardv1.EventData{}
			if err := UnmarshalEventsData(record.Value, eventsData); err != nil {
				return nil, err
			}
			driverCard.SetEventsData(eventsData)
		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData := &cardv1.FaultData{}
			if err := UnmarshalFaultsData(record.Value, faultsData); err != nil {
				return nil, err
			}
			driverCard.SetFaultsData(faultsData)
		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData := &cardv1.DriverActivity{}
			if err := UnmarshalCardActivityData(record.Value, activityData); err != nil {
				return nil, err
			}
			driverCard.SetDriverActivityData(activityData)
		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed := &cardv1.VehiclesUsed{}
			if err := UnmarshalCardVehiclesUsed(record.Value, vehiclesUsed); err != nil {
				return nil, err
			}
			driverCard.SetVehiclesUsed(vehiclesUsed)
		case cardv1.ElementaryFileType_EF_PLACES:
			places := &cardv1.Places{}
			if err := UnmarshalCardPlaces(record.Value, places); err != nil {
				return nil, err
			}
			driverCard.SetPlaces(places)
		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage := &cardv1.CurrentUsage{}
			if err := UnmarshalCardCurrentUsage(record.Value, currentUsage); err != nil {
				return nil, err
			}
			driverCard.SetCurrentUsage(currentUsage)
		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId := &cardv1.DriverCardApplicationIdentification{}
			if err := UnmarshalCardApplicationIdentification(record.Value, appId); err != nil {
				return nil, err
			}
			driverCard.SetApplicationIdentification(appId)
		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity := &cardv1.ControlActivityData{}
			if err := UnmarshalCardControlActivityData(record.Value, controlActivity); err != nil {
				return nil, err
			}
			driverCard.SetControlActivityData(controlActivity)
		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions := &cardv1.SpecificConditions{}
			if err := UnmarshalCardSpecificConditions(record.Value, specificConditions); err != nil {
				return nil, err
			}
			driverCard.SetSpecificConditions(specificConditions)
		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload := &cardv1.LastCardDownload{}
			if err := UnmarshalCardLastDownload(record.Value, lastDownload); err != nil {
				return nil, err
			}
			driverCard.SetLastCardDownload(lastDownload)
		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits := &cardv1.VehicleUnitsUsed{}
			if err := UnmarshalCardVehicleUnitsUsed(record.Value, vehicleUnits); err != nil {
				return nil, err
			}
			driverCard.SetVehicleUnitsUsed(vehicleUnits)
		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces := &cardv1.GnssPlaces{}
			if err := UnmarshalCardGnssPlaces(record.Value, gnssPlaces); err != nil {
				return nil, err
			}
			driverCard.SetGnssPlaces(gnssPlaces)
		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2 := &cardv1.ApplicationIdentificationV2{}
			if err := UnmarshalCardApplicationIdentificationV2(record.Value, appIdV2); err != nil {
				return nil, err
			}
			driverCard.SetApplicationIdentificationV2(appIdV2)
		case cardv1.ElementaryFileType_EF_CARD_CERTIFICATE:
			// Initialize certificates if needed
			if driverCard.GetCertificates() == nil {
				driverCard.SetCertificates(&cardv1.Certificates{})
			}
			driverCard.GetCertificates().SetCardCertificate(record.Value)
		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			// Initialize certificates if needed
			if driverCard.GetCertificates() == nil {
				driverCard.SetCertificates(&cardv1.Certificates{})
			}
			driverCard.GetCertificates().SetCaCertificate(record.Value)
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
