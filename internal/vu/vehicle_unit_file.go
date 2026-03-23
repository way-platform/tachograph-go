package vu

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MarshalVehicleUnitFile serializes a VehicleUnitFile into binary format.
//
// The VehicleUnitFile is marshaled in TV (Tag-Value) format as specified in
// Appendix 7, Section 2.2.6 of the regulation.
func (opts MarshalOptions) MarshalVehicleUnitFile(file *vuv1.VehicleUnitFile) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("vehicle unit file is nil")
	}

	var dst []byte

	switch file.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		gen1 := file.GetGen1()
		if gen1 == nil {
			return nil, fmt.Errorf("Gen1 data is nil")
		}

		// Marshal Overview (TREP 01)
		if overview := gen1.GetOverview(); overview != nil {
			transferData, err := opts.MarshalOverviewGen1(overview)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Overview Gen1: %w", err)
			}
			dst = appendTransfer(dst, vuv1.TransferType_OVERVIEW_GEN1, transferData)
		}

		// Marshal Activities (TREP 02) - multiple transfers
		for i, activities := range gen1.GetActivities() {
			transferData, err := opts.MarshalActivitiesGen1(activities)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Activities Gen1 [%d]: %w", i, err)
			}
			dst = appendTransfer(dst, vuv1.TransferType_ACTIVITIES_GEN1, transferData)
		}

		// Marshal Events and Faults (TREP 03) - multiple transfers
		for i, eventsAndFaults := range gen1.GetEventsAndFaults() {
			transferData, err := opts.MarshalEventsAndFaultsGen1(eventsAndFaults)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EventsAndFaults Gen1 [%d]: %w", i, err)
			}
			dst = appendTransfer(dst, vuv1.TransferType_EVENTS_AND_FAULTS_GEN1, transferData)
		}

		// Marshal Detailed Speed (TREP 04) - multiple transfers
		for i, detailedSpeed := range gen1.GetDetailedSpeed() {
			transferData, err := opts.MarshalDetailedSpeedGen1(detailedSpeed)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal DetailedSpeed Gen1 [%d]: %w", i, err)
			}
			dst = appendTransfer(dst, vuv1.TransferType_DETAILED_SPEED_GEN1, transferData)
		}

		// Marshal Technical Data (TREP 05) - multiple transfers
		for i, technicalData := range gen1.GetTechnicalData() {
			transferData, err := opts.MarshalTechnicalDataGen1(technicalData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal TechnicalData Gen1 [%d]: %w", i, err)
			}
			dst = appendTransfer(dst, vuv1.TransferType_TECHNICAL_DATA_GEN1, transferData)
		}

	case ddv1.Generation_GENERATION_2:
		if file.GetVersion() == ddv1.Version_VERSION_2 {
			// Handle Gen2 V2
			gen2v2 := file.GetGen2V2()
			if gen2v2 == nil {
				return nil, fmt.Errorf("Gen2V2 data is nil")
			}

			// Marshal Overview (TREP 31)
			if overview := gen2v2.GetOverview(); overview != nil {
				transferData, err := opts.MarshalOverviewGen2V2(overview)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Overview Gen2V2: %w", err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_OVERVIEW_GEN2_V2, transferData)
			}

			// Marshal Activities (TREP 32) - multiple transfers
			for i, activities := range gen2v2.GetActivities() {
				transferData, err := opts.MarshalActivitiesGen2V2(activities)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Activities Gen2V2 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_ACTIVITIES_GEN2_V2, transferData)
			}

			// Marshal Events and Faults (TREP 33) - multiple transfers
			for i, eventsAndFaults := range gen2v2.GetEventsAndFaults() {
				transferData, err := opts.MarshalEventsAndFaultsGen2V2(eventsAndFaults)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal EventsAndFaults Gen2V2 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2, transferData)
			}

			// Marshal Detailed Speed (TREP 34) - multiple transfers
			for i, detailedSpeed := range gen2v2.GetDetailedSpeed() {
				transferData, err := opts.MarshalDetailedSpeedGen2(detailedSpeed)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal DetailedSpeed Gen2V2 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_DETAILED_SPEED_GEN2, transferData)
			}

			// Marshal Technical Data (TREP 35) - multiple transfers
			for i, technicalData := range gen2v2.GetTechnicalData() {
				transferData, err := opts.MarshalTechnicalDataGen2V2(technicalData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal TechnicalData Gen2V2 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_TECHNICAL_DATA_GEN2_V2, transferData)
			}

		} else {
			// Handle Gen2 V1
			gen2v1 := file.GetGen2V1()
			if gen2v1 == nil {
				return nil, fmt.Errorf("Gen2V1 data is nil")
			}

			// Marshal Overview (TREP 11)
			if overview := gen2v1.GetOverview(); overview != nil {
				transferData, err := opts.MarshalOverviewGen2V1(overview)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Overview Gen2V1: %w", err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_OVERVIEW_GEN2_V1, transferData)
			}

			// Marshal Activities (TREP 12) - multiple transfers
			for i, activities := range gen2v1.GetActivities() {
				transferData, err := opts.MarshalActivitiesGen2V1(activities)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Activities Gen2V1 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_ACTIVITIES_GEN2_V1, transferData)
			}

			// Marshal Events and Faults (TREP 13) - multiple transfers
			for i, eventsAndFaults := range gen2v1.GetEventsAndFaults() {
				transferData, err := opts.MarshalEventsAndFaultsGen2V1(eventsAndFaults)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal EventsAndFaults Gen2V1 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1, transferData)
			}

			// Marshal Detailed Speed (TREP 14) - multiple transfers
			for i, detailedSpeed := range gen2v1.GetDetailedSpeed() {
				transferData, err := opts.MarshalDetailedSpeedGen2(detailedSpeed)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal DetailedSpeed Gen2V1 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_DETAILED_SPEED_GEN2, transferData)
			}

			// Marshal Technical Data (TREP 15) - multiple transfers
			for i, technicalData := range gen2v1.GetTechnicalData() {
				transferData, err := opts.MarshalTechnicalDataGen2V1(technicalData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal TechnicalData Gen2V1 [%d]: %w", i, err)
				}
				dst = appendTransfer(dst, vuv1.TransferType_TECHNICAL_DATA_GEN2_V1, transferData)
			}
		}

	default:
		return nil, fmt.Errorf("unsupported generation: %v", file.GetGeneration())
	}

	// Append trailing data for round-trip fidelity
	dst = append(dst, file.GetTrailingData()...)

	return dst, nil
}

// ParseRawVehicleUnitFile parses a RawVehicleUnitFile into a fully parsed VehicleUnitFile message.
// Authentication results from the raw file records are propagated to the parsed messages.
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
func (opts ParseOptions) ParseRawVehicleUnitFile(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFile, error) {
	// Determine generation/version
	if len(rawFile.GetRecords()) == 0 {
		return nil, fmt.Errorf("empty VU file")
	}

	firstRecord := rawFile.GetRecords()[0]

	// Dispatch to generation-specific unmarshaller
	output := &vuv1.VehicleUnitFile{}

	switch firstRecord.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		gen1File, err := opts.unmarshalVehicleUnitFileGen1(rawFile)
		if err != nil {
			return nil, err
		}
		output.SetGeneration(ddv1.Generation_GENERATION_1)
		output.SetGen1(gen1File)

	case ddv1.Generation_GENERATION_2:
		if hasGen2V2Transfers(rawFile) {
			gen2v2File, err := opts.unmarshalVehicleUnitFileGen2V2(rawFile)
			if err != nil {
				return nil, err
			}
			output.SetGeneration(ddv1.Generation_GENERATION_2)
			output.SetVersion(ddv1.Version_VERSION_2)
			output.SetGen2V2(gen2v2File)
		} else {
			gen2v1File, err := opts.unmarshalVehicleUnitFileGen2V1(rawFile)
			if err != nil {
				return nil, err
			}
			output.SetGeneration(ddv1.Generation_GENERATION_2)
			output.SetVersion(ddv1.Version_VERSION_1)
			output.SetGen2V1(gen2v1File)
		}

	default:
		return nil, fmt.Errorf("unknown generation: %v", firstRecord.GetGeneration())
	}

	// Propagate trailing data for round-trip fidelity
	if td := rawFile.GetTrailingData(); len(td) > 0 {
		output.SetTrailingData(td)
	}

	return output, nil
}

// hasGen2V2Transfers checks if the raw file contains Gen2 V2 transfers.
// Gen2 V2 is identified by the presence of TREP 00 (DownloadInterfaceVersion)
// or TREP 31-35 transfers.
func hasGen2V2Transfers(rawFile *vuv1.RawVehicleUnitFile) bool {
	for _, record := range rawFile.GetRecords() {
		switch record.GetType() {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION,
			vuv1.TransferType_OVERVIEW_GEN2_V2,
			vuv1.TransferType_ACTIVITIES_GEN2_V2,
			vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2,
			vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
			return true
		}
	}
	return false
}

// unmarshalVehicleUnitFileGen1 unmarshals a Gen1 VU file from raw records.
func (opts ParseOptions) unmarshalVehicleUnitFileGen1(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFileGen1, error) {
	var output vuv1.VehicleUnitFileGen1
	// unmarshalOpts := opts.unmarshal()  // Available for future use

	for _, record := range rawFile.GetRecords() {
		// Get complete transfer value (already combined)
		transferValue := record.GetValue()

		switch record.GetType() {
		case vuv1.TransferType_OVERVIEW_GEN1:
			overview, err := unmarshalOverviewGen1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Overview Gen1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				overview.SetAuthentication(auth)
			}
			output.SetOverview(overview)

		case vuv1.TransferType_ACTIVITIES_GEN1:
			activities, err := unmarshalActivitiesGen1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Activities Gen1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				activities.SetAuthentication(auth)
			}
			output.SetActivities(append(output.GetActivities(), activities))

		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
			eventsAndFaults, err := unmarshalEventsAndFaultsGen1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Events and Faults Gen1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				eventsAndFaults.SetAuthentication(auth)
			}
			output.SetEventsAndFaults(append(output.GetEventsAndFaults(), eventsAndFaults))

		case vuv1.TransferType_DETAILED_SPEED_GEN1:
			detailedSpeed, err := unmarshalDetailedSpeedGen1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Detailed Speed Gen1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				detailedSpeed.SetAuthentication(auth)
			}
			output.SetDetailedSpeed(append(output.GetDetailedSpeed(), detailedSpeed))

		case vuv1.TransferType_TECHNICAL_DATA_GEN1:
			technicalData, err := unmarshalTechnicalDataGen1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Technical Data Gen1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				technicalData.SetAuthentication(auth)
			}
			output.SetTechnicalData(append(output.GetTechnicalData(), technicalData))

		default:
			return nil, fmt.Errorf("unexpected transfer type %v in Gen1 file", record.GetType())
		}
	}

	return &output, nil
}

// unmarshalVehicleUnitFileGen2V1 unmarshals a Gen2 V1 VU file from raw records.
func (opts ParseOptions) unmarshalVehicleUnitFileGen2V1(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFileGen2V1, error) {
	var output vuv1.VehicleUnitFileGen2V1
	// unmarshalOpts := opts.unmarshal()  // Available for future use

	for _, record := range rawFile.GetRecords() {
		// Get complete transfer value (already combined)
		transferValue := record.GetValue()

		switch record.GetType() {
		case vuv1.TransferType_OVERVIEW_GEN2_V1:
			overview, err := unmarshalOverviewGen2V1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Overview Gen2 V1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				overview.SetAuthentication(auth)
			}
			output.SetOverview(overview)

		case vuv1.TransferType_ACTIVITIES_GEN2_V1:
			activities, err := unmarshalActivitiesGen2V1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Activities Gen2 V1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				activities.SetAuthentication(auth)
			}
			output.SetActivities(append(output.GetActivities(), activities))

		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			eventsAndFaults, err := unmarshalEventsAndFaultsGen2V1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Events and Faults Gen2 V1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				eventsAndFaults.SetAuthentication(auth)
			}
			output.SetEventsAndFaults(append(output.GetEventsAndFaults(), eventsAndFaults))

		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed, err := unmarshalDetailedSpeedGen2(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Detailed Speed Gen2: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				detailedSpeed.SetAuthentication(auth)
			}
			output.SetDetailedSpeed(append(output.GetDetailedSpeed(), detailedSpeed))

		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			technicalData, err := unmarshalTechnicalDataGen2V1(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Technical Data Gen2 V1: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				technicalData.SetAuthentication(auth)
			}
			output.SetTechnicalData(append(output.GetTechnicalData(), technicalData))

		default:
			return nil, fmt.Errorf("unexpected transfer type %v in Gen2 V1 file", record.GetType())
		}
	}

	return &output, nil
}

// unmarshalVehicleUnitFileGen2V2 unmarshals a Gen2 V2 VU file from raw records.
func (opts ParseOptions) unmarshalVehicleUnitFileGen2V2(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFileGen2V2, error) {
	var output vuv1.VehicleUnitFileGen2V2
	// unmarshalOpts := opts.unmarshal()  // Available for future use

	for _, record := range rawFile.GetRecords() {
		// Get complete transfer value (already combined)
		transferValue := record.GetValue()

		switch record.GetType() {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			// Download interface version can be parsed if needed
			// For now, skip as it's mainly used for version detection
			// output.SetDownloadInterfaceVersion(...)

		case vuv1.TransferType_OVERVIEW_GEN2_V2:
			overview, err := unmarshalOverviewGen2V2(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Overview Gen2 V2: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				overview.SetAuthentication(auth)
			}
			output.SetOverview(overview)

		case vuv1.TransferType_ACTIVITIES_GEN2_V2:
			activities, err := unmarshalActivitiesGen2V2(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Activities Gen2 V2: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				activities.SetAuthentication(auth)
			}
			output.SetActivities(append(output.GetActivities(), activities))

		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2:
			eventsAndFaults, err := unmarshalEventsAndFaultsGen2V2(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Events and Faults Gen2 V2: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				eventsAndFaults.SetAuthentication(auth)
			}
			output.SetEventsAndFaults(append(output.GetEventsAndFaults(), eventsAndFaults))

		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed, err := unmarshalDetailedSpeedGen2(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Detailed Speed Gen2: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				detailedSpeed.SetAuthentication(auth)
			}
			output.SetDetailedSpeed(append(output.GetDetailedSpeed(), detailedSpeed))

		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
			technicalData, err := unmarshalTechnicalDataGen2V2(transferValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshal Technical Data Gen2 V2: %w", err)
			}
			// Propagate authentication
			if auth := record.GetAuthentication(); auth != nil {
				technicalData.SetAuthentication(auth)
			}
			output.SetTechnicalData(append(output.GetTechnicalData(), technicalData))

		default:
			return nil, fmt.Errorf("unexpected transfer type %v in Gen2 V2 file", record.GetType())
		}
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

// appendTransfer appends a transfer in TV format: [Tag: 2 bytes][Value: N bytes]
func appendTransfer(dst []byte, transferType vuv1.TransferType, data []byte) []byte {
	tag := getTagForTransferType(transferType)
	dst = binary.BigEndian.AppendUint16(dst, tag)
	dst = append(dst, data...)
	return dst
}

// getTagForTransferType returns the TV format tag for a given transfer type
func getTagForTransferType(transferType vuv1.TransferType) uint16 {
	valueDesc := transferType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(transferType))
	if valueDesc == nil {
		return 0
	}

	opts := valueDesc.Options()
	if !proto.HasExtension(opts, vuv1.E_TrepValue) {
		return 0
	}

	trepValue := proto.GetExtension(opts, vuv1.E_TrepValue).(int32)
	// VU tags are constructed as 0x76XX where XX is the TREP value
	return uint16(0x7600 | (uint16(trepValue) & 0xFF))
}
