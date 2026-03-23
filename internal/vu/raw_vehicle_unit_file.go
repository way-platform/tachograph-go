package vu

import (
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// splitTransferValue splits a record's value into data and signature portions.
// The signature size was computed once during unmarshal and stored in the record,
// making this operation a simple slice operation with no complex size calculations.
//
// Returns:
//   - data: The data portion (everything except the signature)
//   - signature: The signature portion (last N bytes)
//   - error: If the value is too short for the stored signature size
func splitTransferValue(record *vuv1.RawVehicleUnitFile_Record) (data, signature []byte, err error) {
	value := record.GetValue()
	sigSize := int(record.GetSignatureSize())

	if len(value) < sigSize {
		return nil, nil, fmt.Errorf("value too short for signature: got %d bytes, need at least %d", len(value), sigSize)
	}

	dataSize := len(value) - sigSize
	return value[:dataSize], value[dataSize:], nil
}

// UnmarshalRawVehicleUnitFile performs the first parsing pass, identifying TV record
// boundaries and extracting complete values including embedded signatures.
//
// This function does NOT parse the semantic content of records - it only slices the
// binary data into individual records for later processing.
//
// The TV (Tag-Value) format is defined in Appendix 7, Section 2.2.6.
// Each record consists of:
//   - Tag: 2 bytes (0x76XX where XX is the TREP value)
//   - Value: Variable length (determined by parsing the structure)
//
// The challenge of the TV format is that the length is not explicitly encoded - it must
// be calculated by understanding the structure of each transfer type.
func (opts UnmarshalOptions) UnmarshalRawVehicleUnitFile(data []byte) (*vuv1.RawVehicleUnitFile, error) {
	var rawFile vuv1.RawVehicleUnitFile
	offset := 0

	for offset < len(data) {
		// Read tag (2 bytes)
		if offset+2 > len(data) {
			return nil, fmt.Errorf("insufficient data for tag at offset %d: need 2 bytes, have %d", offset, len(data)-offset)
		}
		tag := binary.BigEndian.Uint16(data[offset:])

		// Tags must follow 0x76XX pattern (SID 0x76 + TREP byte).
		// Non-matching bytes indicate end of transfer data section
		// (e.g., trailing file-level signature in Gen2v2).
		if tag>>8 != 0x76 {
			rawFile.SetTrailingData(data[offset:])
			break
		}
		offset += 2

		// Determine transfer type from tag
		transferType := findTransferTypeByTag(tag)
		if transferType == vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED {
			if opts.Strict {
				return nil, fmt.Errorf("unknown tag: 0x%04X at offset %d", tag, offset-2)
			}
			// In non-strict mode, we can't know the structure without
			// knowing the transfer type, so we have to stop here.
			rawFile.SetTrailingData(data[offset-2:])
			break
		}

		// Calculate size of value (including embedded signature)
		totalSize, sigSize, err := sizeOfTransferValue(data[offset:], transferType)
		if err != nil {
			return nil, fmt.Errorf("sizeOf failed for %v at offset %d: %w", transferType, offset, err)
		}

		// Extract complete value (includes signature)
		if offset+totalSize > len(data) {
			return nil, fmt.Errorf("insufficient data for %v value: need %d bytes, have %d", transferType, totalSize, len(data)-offset)
		}
		value := data[offset : offset+totalSize]
		offset += totalSize

		// Create record with complete transfer value
		record := &vuv1.RawVehicleUnitFile_Record{}
		record.SetTag(uint32(tag))
		record.SetType(transferType)
		record.SetGeneration(generationFromTransferType(transferType))
		record.SetValue(value)                  // Store complete value
		record.SetSignatureSize(int32(sigSize)) // Store signature size for efficient splitting

		rawFile.SetRecords(append(rawFile.GetRecords(), record))
	}

	return &rawFile, nil
}

// sizeOfTransferValue dispatches to type-specific sizeOf functions.
// Returns both the total byte size (including signature) and the signature size.
// The data portion size can be calculated as: totalSize - signatureSize.
func sizeOfTransferValue(data []byte, transferType vuv1.TransferType) (totalSize, signatureSize int, err error) {
	switch transferType {
	case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
		return sizeOfDownloadInterfaceVersion(data, transferType)
	case vuv1.TransferType_OVERVIEW_GEN1, vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
		return sizeOfOverview(data, transferType)
	case vuv1.TransferType_ACTIVITIES_GEN1, vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
		return sizeOfActivities(data, transferType)
	case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2:
		return sizeOfEventsAndFaults(data, transferType)
	case vuv1.TransferType_DETAILED_SPEED_GEN1, vuv1.TransferType_DETAILED_SPEED_GEN2:
		return sizeOfDetailedSpeed(data, transferType)
	case vuv1.TransferType_TECHNICAL_DATA_GEN1, vuv1.TransferType_TECHNICAL_DATA_GEN2_V1, vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
		return sizeOfTechnicalData(data, transferType)
	case vuv1.TransferType_CARD_DOWNLOAD:
		return sizeOfCardDownload(data, transferType)
	default:
		return 0, 0, fmt.Errorf("unsupported transfer type: %v", transferType)
	}
}

// sizeOfRecordArray parses a Gen2 RecordArray header and returns the total size.
// Gen2 uses RecordArray structures with a 5-byte header containing size information.
//
// RecordArray header format (See Appendix 1, Section 1.1.3):
//   - recordType: 1 byte (identifies the type of records in the array)
//   - recordSize: 2 bytes (big-endian, size in bytes of each record)
//   - noOfRecords: 2 bytes (big-endian, number of records in the array)
//
// Total size = 5 + (recordSize * noOfRecords)
func sizeOfRecordArray(data []byte, offset int) (int, error) {
	const headerSize = 5
	if len(data[offset:]) < headerSize {
		return 0, fmt.Errorf("insufficient data for RecordArray header: need %d, have %d", headerSize, len(data[offset:]))
	}

	recordSize := binary.BigEndian.Uint16(data[offset+1:])
	noOfRecords := binary.BigEndian.Uint16(data[offset+3:])

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return totalSize, nil
}

// generationFromTransferType extracts generation from transfer type using protobuf reflection.
func generationFromTransferType(transferType vuv1.TransferType) ddv1.Generation {
	// Use protobuf reflection to get generation from enum options
	valueDesc := transferType.Descriptor().Values().ByNumber(transferType.Number())
	if valueDesc == nil {
		return ddv1.Generation_GENERATION_UNSPECIFIED
	}
	opts := valueDesc.Options()
	if proto.HasExtension(opts, ddv1.E_Generation) {
		return proto.GetExtension(opts, ddv1.E_Generation).(ddv1.Generation)
	}
	return ddv1.Generation_GENERATION_UNSPECIFIED
}
