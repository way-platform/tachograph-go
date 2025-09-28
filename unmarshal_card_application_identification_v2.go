package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardApplicationIdentificationV2 parses the binary data for an EF_ApplicationIdentificationV2 record.
//
// ASN.1 Specification (Data Dictionary 2.2):
//
//	ApplicationIdentificationV2 ::= SEQUENCE {
//	    noOfBorderCrossingRecords    INTEGER(0..255),
//	    noOfLoadUnloadRecords        INTEGER(0..255),
//	    noOfLoadTypeEntryRecords     INTEGER(0..255),
//	    vuConfigurationLengthRange   INTEGER(0..255)
//	}
//
// Binary Layout (4 bytes):
//
//	0-0: noOfBorderCrossingRecords (1 byte)
//	1-1: noOfLoadUnloadRecords (1 byte)
//	2-2: noOfLoadTypeEntryRecords (1 byte)
//	3-3: vuConfigurationLengthRange (1 byte)
func unmarshalCardApplicationIdentificationV2(data []byte) (*cardv1.ApplicationIdentificationV2, error) {
	const (
		// EF_ApplicationIdentificationV2 record size
		EF_APPLICATION_IDENTIFICATION_V2_SIZE = 4

		// Field offsets
		NO_OF_BORDER_CROSSING_RECORDS_OFFSET = 0
		NO_OF_LOAD_UNLOAD_RECORDS_OFFSET     = 1
		NO_OF_LOAD_TYPE_ENTRY_RECORDS_OFFSET = 2
		VU_CONFIGURATION_LENGTH_RANGE_OFFSET = 3
	)

	if len(data) < EF_APPLICATION_IDENTIFICATION_V2_SIZE {
		return nil, fmt.Errorf("insufficient data for application identification V2: got %d bytes, need %d", len(data), EF_APPLICATION_IDENTIFICATION_V2_SIZE)
	}
	var target cardv1.ApplicationIdentificationV2
	r := bytes.NewReader(data)

	// For now, assume this is a driver card and create the driver data
	driver := &cardv1.ApplicationIdentificationV2_Driver{}

	// Read border crossing records count (1 byte)
	var borderCrossingCount byte
	if err := binary.Read(r, binary.BigEndian, &borderCrossingCount); err != nil {
		return nil, fmt.Errorf("failed to read border crossing records count: %w", err)
	}
	driver.SetBorderCrossingRecordsCount(int32(borderCrossingCount))

	// Read load/unload records count (1 byte)
	var loadUnloadCount byte
	if err := binary.Read(r, binary.BigEndian, &loadUnloadCount); err != nil {
		return nil, fmt.Errorf("failed to read load/unload records count: %w", err)
	}
	driver.SetLoadUnloadRecordsCount(int32(loadUnloadCount))

	// Read load type entry records count (1 byte)
	var loadTypeCount byte
	if err := binary.Read(r, binary.BigEndian, &loadTypeCount); err != nil {
		return nil, fmt.Errorf("failed to read load type entry records count: %w", err)
	}
	driver.SetLoadTypeEntryRecordsCount(int32(loadTypeCount))

	// Read VU configuration length range (1 byte)
	var vuConfigRange byte
	if err := binary.Read(r, binary.BigEndian, &vuConfigRange); err != nil {
		return nil, fmt.Errorf("failed to read VU configuration length range: %w", err)
	}
	driver.SetVuConfigurationLengthRange(int32(vuConfigRange))

	// Set the driver data and card type
	target.SetDriver(driver)
	target.SetCardType(cardv1.CardType_DRIVER_CARD)

	return &target, nil
}
