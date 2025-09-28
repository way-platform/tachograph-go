package tachograph

import (
	"encoding/binary"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

func unmarshalVehicleUnitFile(input []byte) (*vuv1.VehicleUnitFile, error) {
	var output vuv1.VehicleUnitFile
	offset := 0

	for offset+2 <= len(input) { // Need at least 2 bytes for tag
		// Read Tag - 2 bytes for VU files (TV format)
		tag := binary.BigEndian.Uint16(input[offset:])
		offset += 2

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
		var bytesRead int
		var err error

		switch transferType {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			version := &vuv1.DownloadInterfaceVersion{}
			bytesRead, err = UnmarshalDownloadInterfaceVersion(input, offset, version)
			if err != nil {
				return nil, err
			}
			transfer.SetDownloadInterfaceVersion(version)
		case vuv1.TransferType_OVERVIEW_GEN1:
			overview := &vuv1.Overview{}
			bytesRead, err = UnmarshalOverview(input, offset, overview, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetOverview(overview)
		case vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			overview := &vuv1.Overview{}
			generation := 2
			bytesRead, err = UnmarshalOverview(input, offset, overview, generation)
			if err != nil {
				return nil, err
			}
			transfer.SetOverview(overview)
		case vuv1.TransferType_ACTIVITIES_GEN1:
			activities := &vuv1.Activities{}
			bytesRead, err = UnmarshalVuActivities(input, offset, activities, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			activities := &vuv1.Activities{}
			generation := 2
			bytesRead, err = UnmarshalVuActivities(input, offset, activities, generation)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
			eventsAndFaults := &vuv1.EventsAndFaults{}
			bytesRead, err = UnmarshalVuEventsAndFaults(input, offset, eventsAndFaults, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetEventsAndFaults(eventsAndFaults)
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			eventsAndFaults := &vuv1.EventsAndFaults{}
			bytesRead, err = UnmarshalVuEventsAndFaults(input, offset, eventsAndFaults, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetEventsAndFaults(eventsAndFaults)
		case vuv1.TransferType_DETAILED_SPEED_GEN1:
			detailedSpeed := &vuv1.DetailedSpeed{}
			bytesRead, err = UnmarshalVuDetailedSpeed(input, offset, detailedSpeed, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetDetailedSpeed(detailedSpeed)
		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed := &vuv1.DetailedSpeed{}
			bytesRead, err = UnmarshalVuDetailedSpeed(input, offset, detailedSpeed, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetDetailedSpeed(detailedSpeed)
		case vuv1.TransferType_TECHNICAL_DATA_GEN1:
			technicalData := &vuv1.TechnicalData{}
			bytesRead, err = UnmarshalVuTechnicalData(input, offset, technicalData, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetTechnicalData(technicalData)
		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			technicalData := &vuv1.TechnicalData{}
			bytesRead, err = UnmarshalVuTechnicalData(input, offset, technicalData, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetTechnicalData(technicalData)
		default:
			// For now, skip unknown transfer types
			break
		}

		// Advance offset by the number of bytes read
		offset += bytesRead
		output.SetTransfers(append(output.GetTransfers(), transfer))
	}
	return &output, nil
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
