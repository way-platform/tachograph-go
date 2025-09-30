package vu

import (
	"encoding/binary"
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// UnmarshalVehicleUnitFile parses VU file data into a protobuf VehicleUnitFile message.
func UnmarshalVehicleUnitFile(data []byte) (*vuv1.VehicleUnitFile, error) {
	return unmarshalVehicleUnitFile(data)
}

// MarshalVehicleUnitFile serializes a VehicleUnitFile into binary format.
func MarshalVehicleUnitFile(file *vuv1.VehicleUnitFile) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("vehicle unit file is nil")
	}

	// Allocate a buffer large enough for the VU file
	buf := make([]byte, 0, 1024*1024) // 1MB initial capacity

	// Use the existing AppendVU function
	return appendVU(buf, file)
}

// unmarshalVehicleUnitFile unmarshals a vehicle unit file from binary data.
//
// The data type `VehicleUnitFile` represents a complete vehicle unit file structure.
//
// ASN.1 Definition:
//
//	VehicleUnitFile ::= SEQUENCE OF Transfer
//
//	Transfer ::= SEQUENCE {
//	    type    TransferType,
//	    data    CHOICE {
//	        downloadInterfaceVersion    DownloadInterfaceVersion,
//	        overview                   Overview,
//	        activities                 Activities,
//	        eventsAndFaults           EventsAndFaults,
//	        detailedSpeed             DetailedSpeed,
//	        technicalData             TechnicalData
//	    }
//	}
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
			bytesRead, err = unmarshalDownloadInterfaceVersion(input, offset, version)
			if err != nil {
				return nil, err
			}
			transfer.SetDownloadInterfaceVersion(version)
		case vuv1.TransferType_OVERVIEW_GEN1:
			overview := &vuv1.Overview{}
			bytesRead, err = unmarshalOverview(input, offset, overview, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetOverview(overview)
		case vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			overview := &vuv1.Overview{}
			generation := 2
			bytesRead, err = unmarshalOverview(input, offset, overview, generation)
			if err != nil {
				return nil, err
			}
			transfer.SetOverview(overview)
		case vuv1.TransferType_ACTIVITIES_GEN1:
			activities := &vuv1.Activities{}
			bytesRead, err = unmarshalVuActivitiesGen1(input, offset, activities)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			activities := &vuv1.Activities{}
			bytesRead, err = unmarshalVuActivitiesGen2(input, offset, activities)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
			eventsAndFaults := &vuv1.EventsAndFaults{}
			bytesRead, err = unmarshalVuEventsAndFaults(input, offset, eventsAndFaults, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetEventsAndFaults(eventsAndFaults)
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			eventsAndFaults := &vuv1.EventsAndFaults{}
			bytesRead, err = unmarshalVuEventsAndFaults(input, offset, eventsAndFaults, 2)
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
		}

		// Advance offset by the number of bytes read
		offset += bytesRead
		output.SetTransfers(append(output.GetTransfers(), transfer))
	}
	return &output, nil
}

// findTransferTypeByTag maps VU transfer tags to TransferType enum values
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

// appendVU orchestrates writing a VU file in TV format
//
// The data type `VehicleUnitFile` represents a complete vehicle unit file structure.
//
// ASN.1 Definition:
//
//	VehicleUnitFile ::= SEQUENCE OF Transfer
//
//	Transfer ::= SEQUENCE {
//	    type    TransferType,
//	    data    CHOICE {
//	        downloadInterfaceVersion    DownloadInterfaceVersion,
//	        overview                   Overview,
//	        activities                 Activities,
//	        eventsAndFaults           EventsAndFaults,
//	        detailedSpeed             DetailedSpeed,
//	        technicalData             TechnicalData
//	    }
//	}
func appendVU(dst []byte, vuFile *vuv1.VehicleUnitFile) ([]byte, error) {
	if vuFile == nil {
		return dst, nil
	}

	for _, transfer := range vuFile.GetTransfers() {
		// Write TV format tag (2 bytes) for this transfer
		tag, err := getTagForTransferType(transfer.GetType())
		if err != nil {
			return nil, fmt.Errorf("failed to get tag for transfer type %v: %w", transfer.GetType(), err)
		}
		dst = binary.BigEndian.AppendUint16(dst, tag)

		// Write the transfer data based on type
		switch transfer.GetType() {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			if version := transfer.GetDownloadInterfaceVersion(); version != nil {
				// Placeholder implementation - just append some dummy data
				dst = append(dst, 0x01, 0x01) // 2 bytes for generation and version
			}
		case vuv1.TransferType_OVERVIEW_GEN1:
			if overview := transfer.GetOverview(); overview != nil {
				// Placeholder implementation - just append some dummy data
				dst = append(dst, make([]byte, 128)...) // 128 bytes of zeros
			}
		case vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			if overview := transfer.GetOverview(); overview != nil {
				// Placeholder implementation - just append some dummy data
				dst = append(dst, make([]byte, 128)...) // 128 bytes of zeros
			}
		case vuv1.TransferType_ACTIVITIES_GEN1:
			if activities := transfer.GetActivities(); activities != nil {
				dst, err = appendVuActivitiesBytes(dst, activities)
				if err != nil {
					return nil, fmt.Errorf("failed to append activities gen1: %w", err)
				}
			}
		case vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			if activities := transfer.GetActivities(); activities != nil {
				dst, err = appendVuActivitiesBytes(dst, activities)
				if err != nil {
					return nil, fmt.Errorf("failed to append activities gen2: %w", err)
				}
			}
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
			if eventsAndFaults := transfer.GetEventsAndFaults(); eventsAndFaults != nil {
				dst, err = appendVuEventsAndFaultsBytes(dst, eventsAndFaults)
				if err != nil {
					return nil, fmt.Errorf("failed to append events and faults gen1: %w", err)
				}
			}
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2:
			if eventsAndFaults := transfer.GetEventsAndFaults(); eventsAndFaults != nil {
				dst, err = appendVuEventsAndFaultsBytes(dst, eventsAndFaults)
				if err != nil {
					return nil, fmt.Errorf("failed to append events and faults gen2: %w", err)
				}
			}
		case vuv1.TransferType_DETAILED_SPEED_GEN1:
			if detailedSpeed := transfer.GetDetailedSpeed(); detailedSpeed != nil {
				dst, err = appendVuDetailedSpeedBytes(dst, detailedSpeed)
				if err != nil {
					return nil, fmt.Errorf("failed to append detailed speed gen1: %w", err)
				}
			}
		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			if detailedSpeed := transfer.GetDetailedSpeed(); detailedSpeed != nil {
				dst, err = appendVuDetailedSpeedBytes(dst, detailedSpeed)
				if err != nil {
					return nil, fmt.Errorf("failed to append detailed speed gen2: %w", err)
				}
			}
		case vuv1.TransferType_TECHNICAL_DATA_GEN1:
			if technicalData := transfer.GetTechnicalData(); technicalData != nil {
				dst, err = appendVuTechnicalDataBytes(dst, technicalData)
				if err != nil {
					return nil, fmt.Errorf("failed to append technical data gen1: %w", err)
				}
			}
		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1, vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
			if technicalData := transfer.GetTechnicalData(); technicalData != nil {
				dst, err = appendVuTechnicalDataBytes(dst, technicalData)
				if err != nil {
					return nil, fmt.Errorf("failed to append technical data gen2: %w", err)
				}
			}
		case vuv1.TransferType_CARD_DOWNLOAD:
			if cardDownload := transfer.GetCardDownload(); cardDownload != nil {
				// Placeholder implementation - just append some dummy data
				dst = append(dst, make([]byte, 64)...) // 64 bytes of zeros
			}
		default:
			return nil, fmt.Errorf("unsupported transfer type: %v", transfer.GetType())
		}
	}

	return dst, nil
}

// getTagForTransferType returns the TV format tag for a given transfer type
func getTagForTransferType(transferType vuv1.TransferType) (uint16, error) {
	values := vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		if vuv1.TransferType(valueDesc.Number()) == transferType {
			opts := valueDesc.Options()
			if proto.HasExtension(opts, vuv1.E_TrepValue) {
				trepValue := proto.GetExtension(opts, vuv1.E_TrepValue).(int32)
				// VU tags are constructed as 0x76XX where XX is the TREP value
				return uint16(0x7600 | (uint16(trepValue) & 0xFF)), nil
			}
		}
	}
	return 0, fmt.Errorf("no TREP value found for transfer type: %v", transferType)
}
