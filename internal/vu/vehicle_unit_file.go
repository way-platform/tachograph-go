package vu

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
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
	// Pass 1: Slice into RawVehicleUnitFile
	rawFile, err := unmarshalRawVehicleUnitFile(input)
	if err != nil {
		return nil, fmt.Errorf("raw unmarshal failed: %w", err)
	}

	// Determine generation/version
	if len(rawFile.GetRecords()) == 0 {
		return nil, fmt.Errorf("empty VU file")
	}

	firstRecord := rawFile.GetRecords()[0]

	// Dispatch to generation-specific unmarshaller
	output := &vuv1.VehicleUnitFile{}

	switch firstRecord.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		gen1File, err := unmarshalVehicleUnitFileGen1(rawFile)
		if err != nil {
			return nil, err
		}
		output.SetGeneration(ddv1.Generation_GENERATION_1)
		output.SetGen1(gen1File)

	case ddv1.Generation_GENERATION_2:
		if hasGen2V2Transfers(rawFile) {
			gen2v2File, err := unmarshalVehicleUnitFileGen2V2(rawFile)
			if err != nil {
				return nil, err
			}
			output.SetGeneration(ddv1.Generation_GENERATION_2)
			output.SetVersion(ddv1.Version_VERSION_2)
			output.SetGen2V2(gen2v2File)
		} else {
			gen2v1File, err := unmarshalVehicleUnitFileGen2V1(rawFile)
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
func unmarshalVehicleUnitFileGen1(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFileGen1, error) {
	var output vuv1.VehicleUnitFileGen1

	for _, record := range rawFile.GetRecords() {
		switch record.GetType() {
		case vuv1.TransferType_OVERVIEW_GEN1:
			overview, err := unmarshalOverviewGen1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Overview Gen1: %w", err)
			}
			output.SetOverview(overview)

		case vuv1.TransferType_ACTIVITIES_GEN1:
			activities, err := unmarshalActivitiesGen1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Activities Gen1: %w", err)
			}
			output.SetActivities(append(output.GetActivities(), activities))

		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
			eventsAndFaults, err := unmarshalEventsAndFaultsGen1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Events and Faults Gen1: %w", err)
			}
			output.SetEventsAndFaults(append(output.GetEventsAndFaults(), eventsAndFaults))

		case vuv1.TransferType_DETAILED_SPEED_GEN1:
			detailedSpeed, err := unmarshalDetailedSpeedGen1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Detailed Speed Gen1: %w", err)
			}
			output.SetDetailedSpeed(append(output.GetDetailedSpeed(), detailedSpeed))

		case vuv1.TransferType_TECHNICAL_DATA_GEN1:
			technicalData, err := unmarshalTechnicalDataGen1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Technical Data Gen1: %w", err)
			}
			output.SetTechnicalData(append(output.GetTechnicalData(), technicalData))

		default:
			return nil, fmt.Errorf("unexpected transfer type %v in Gen1 file", record.GetType())
		}
	}

	return &output, nil
}

// unmarshalVehicleUnitFileGen2V1 unmarshals a Gen2 V1 VU file from raw records.
func unmarshalVehicleUnitFileGen2V1(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFileGen2V1, error) {
	var output vuv1.VehicleUnitFileGen2V1

	for _, record := range rawFile.GetRecords() {
		switch record.GetType() {
		case vuv1.TransferType_OVERVIEW_GEN2_V1:
			overview, err := unmarshalOverviewGen2V1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Overview Gen2 V1: %w", err)
			}
			output.SetOverview(overview)

		case vuv1.TransferType_ACTIVITIES_GEN2_V1:
			activities, err := unmarshalActivitiesGen2V1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Activities Gen2 V1: %w", err)
			}
			output.SetActivities(append(output.GetActivities(), activities))

		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			eventsAndFaults, err := unmarshalEventsAndFaultsGen2V1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Events and Faults Gen2 V1: %w", err)
			}
			output.SetEventsAndFaults(append(output.GetEventsAndFaults(), eventsAndFaults))

		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed, err := unmarshalDetailedSpeedGen2(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Detailed Speed Gen2: %w", err)
			}
			output.SetDetailedSpeed(append(output.GetDetailedSpeed(), detailedSpeed))

		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			technicalData, err := unmarshalTechnicalDataGen2V1(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Technical Data Gen2 V1: %w", err)
			}
			output.SetTechnicalData(append(output.GetTechnicalData(), technicalData))

		default:
			return nil, fmt.Errorf("unexpected transfer type %v in Gen2 V1 file", record.GetType())
		}
	}

	return &output, nil
}

// unmarshalVehicleUnitFileGen2V2 unmarshals a Gen2 V2 VU file from raw records.
func unmarshalVehicleUnitFileGen2V2(rawFile *vuv1.RawVehicleUnitFile) (*vuv1.VehicleUnitFileGen2V2, error) {
	var output vuv1.VehicleUnitFileGen2V2

	for _, record := range rawFile.GetRecords() {
		switch record.GetType() {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			// Download interface version can be parsed if needed
			// For now, skip as it's mainly used for version detection
			// output.SetDownloadInterfaceVersion(...)

		case vuv1.TransferType_OVERVIEW_GEN2_V2:
			overview, err := unmarshalOverviewGen2V2(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Overview Gen2 V2: %w", err)
			}
			output.SetOverview(overview)

		case vuv1.TransferType_ACTIVITIES_GEN2_V2:
			activities, err := unmarshalActivitiesGen2V2(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Activities Gen2 V2: %w", err)
			}
			output.SetActivities(append(output.GetActivities(), activities))

		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2:
			eventsAndFaults, err := unmarshalEventsAndFaultsGen2V2(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Events and Faults Gen2 V2: %w", err)
			}
			output.SetEventsAndFaults(append(output.GetEventsAndFaults(), eventsAndFaults))

		case vuv1.TransferType_DETAILED_SPEED_GEN2:
			detailedSpeed, err := unmarshalDetailedSpeedGen2(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Detailed Speed Gen2: %w", err)
			}
			output.SetDetailedSpeed(append(output.GetDetailedSpeed(), detailedSpeed))

		case vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
			technicalData, err := unmarshalTechnicalDataGen2V2(record.GetValue())
			if err != nil {
				return nil, fmt.Errorf("unmarshal Technical Data Gen2 V2: %w", err)
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

	// Dispatch to generation-specific marshaller
	// Note: The plan specifies that VehicleUnitFile marshalling is NOT implemented,
	// only RawVehicleUnitFile marshalling is used for round-tripping.
	// Individual transfer marshalling is implemented for testing purposes.
	return nil, fmt.Errorf("VehicleUnitFile marshalling is not implemented; use RawVehicleUnitFile for binary round-tripping")
}

// getTagForTransferType returns the TV format tag for a given transfer type
