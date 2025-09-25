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
		return unmarshalCardFile(data)
	case UnitFileType:
		return unmarshalVU(data)
	}
	return nil, errors.New("unknown or unsupported file type")
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
