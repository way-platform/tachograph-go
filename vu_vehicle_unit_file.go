package tachograph

import (
	"encoding/binary"
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

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
			bytesRead, err = unmarshalVuActivities(input, offset, activities, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetActivities(activities)
		case vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			activities := &vuv1.Activities{}
			generation := 2
			bytesRead, err = unmarshalVuActivities(input, offset, activities, generation)
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
			bytesRead, err = unmarshalVuDetailedSpeed(input, offset, detailedSpeed, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetDetailedSpeed(detailedSpeed)
		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed := &vuv1.DetailedSpeed{}
			bytesRead, err = unmarshalVuDetailedSpeed(input, offset, detailedSpeed, 2)
			if err != nil {
				return nil, err
			}
			transfer.SetDetailedSpeed(detailedSpeed)
		case vuv1.TransferType_TECHNICAL_DATA_GEN1:
			technicalData := &vuv1.TechnicalData{}
			bytesRead, err = unmarshalVuTechnicalData(input, offset, technicalData, 1)
			if err != nil {
				return nil, err
			}
			transfer.SetTechnicalData(technicalData)
		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			technicalData := &vuv1.TechnicalData{}
			bytesRead, err = unmarshalVuTechnicalData(input, offset, technicalData, 2)
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
				dst, err = appendDownloadInterfaceVersionBytes(dst, version)
				if err != nil {
					return nil, fmt.Errorf("failed to append download interface version: %w", err)
				}
			}
		case vuv1.TransferType_OVERVIEW_GEN1:
			if overview := transfer.GetOverview(); overview != nil {
				dst, err = appendOverviewBytes(dst, overview, 1)
				if err != nil {
					return nil, fmt.Errorf("failed to append overview gen1: %w", err)
				}
			}
		case vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			if overview := transfer.GetOverview(); overview != nil {
				dst, err = appendOverviewBytes(dst, overview, 2)
				if err != nil {
					return nil, fmt.Errorf("failed to append overview gen2: %w", err)
				}
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
				dst, err = appendCardDownloadBytes(dst, cardDownload)
				if err != nil {
					return nil, fmt.Errorf("failed to append card download: %w", err)
				}
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

// []byte-based append functions for VU data types

func appendDownloadInterfaceVersionBytes(dst []byte, version *vuv1.DownloadInterfaceVersion) ([]byte, error) {
	if version == nil {
		return dst, nil
	}

	// DownloadInterfaceVersion structure (2 bytes: generation + version)
	// Generation (1 byte)
	generation := version.GetGeneration()
	generationValue, ok := getProtocolValueFromEnumInternal(generation)
	if !ok {
		return nil, fmt.Errorf("invalid generation value")
	}
	dst = append(dst, byte(generationValue))

	// Version (1 byte)
	versionValue := version.GetVersion()
	versionByte, ok := getProtocolValueFromEnumInternal(versionValue)
	if !ok {
		return nil, fmt.Errorf("invalid version value")
	}
	dst = append(dst, byte(versionByte))

	return dst, nil
}

func appendOverviewBytes(dst []byte, overview *vuv1.Overview, generation int) ([]byte, error) {
	if overview == nil {
		return dst, nil
	}

	// For now, implement a simplified version that just writes signature data
	// This ensures the interface is complete while allowing for future enhancement
	if generation == 1 {
		signature := overview.GetSignatureGen1()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	} else {
		signature := overview.GetSignatureGen2()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	}

	return dst, nil
}

func appendActivitiesBytes(dst []byte, activities *vuv1.Activities, generation int) ([]byte, error) {
	if activities == nil {
		return dst, nil
	}

	// For now, implement a simplified version that just writes signature data
	// This ensures the interface is complete while allowing for future enhancement
	if generation == 1 {
		signature := activities.GetSignatureGen1()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	} else {
		signature := activities.GetSignatureGen2()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	}

	return dst, nil
}

func appendEventsAndFaultsBytes(dst []byte, eventsAndFaults *vuv1.EventsAndFaults, generation int) ([]byte, error) {
	if eventsAndFaults == nil {
		return dst, nil
	}

	// For now, implement a simplified version that just writes signature data
	// This ensures the interface is complete while allowing for future enhancement
	if generation == 1 {
		signature := eventsAndFaults.GetSignatureGen1()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	} else {
		signature := eventsAndFaults.GetSignatureGen2()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	}

	return dst, nil
}

func appendDetailedSpeedBytes(dst []byte, detailedSpeed *vuv1.DetailedSpeed, generation int) ([]byte, error) {
	if detailedSpeed == nil {
		return dst, nil
	}

	// For now, implement a simplified version that just writes signature data
	// This ensures the interface is complete while allowing for future enhancement
	if generation == 1 {
		signature := detailedSpeed.GetSignatureGen1()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	} else {
		signature := detailedSpeed.GetSignatureGen2()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	}

	return dst, nil
}

func appendTechnicalDataBytes(dst []byte, technicalData *vuv1.TechnicalData, generation int) ([]byte, error) {
	if technicalData == nil {
		return dst, nil
	}

	// For now, implement a simplified version that just writes signature data
	// This ensures the interface is complete while allowing for future enhancement
	if generation == 1 {
		signature := technicalData.GetSignatureGen1()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	} else {
		signature := technicalData.GetSignatureGen2()
		if len(signature) > 0 {
			dst = append(dst, signature...)
		}
	}

	return dst, nil
}

func appendCardDownloadBytes(dst []byte, cardDownload *vuv1.CardDownload) ([]byte, error) {
	// TODO: Implement card download append function
	// For now, return the destination unchanged
	return dst, nil
}
