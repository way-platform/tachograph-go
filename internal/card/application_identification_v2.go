package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardApplicationIdentificationV2 parses the binary data for an EF_ApplicationIdentificationV2 record.
//
// The data type `ApplicationIdentificationV2` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition:
//
//	ApplicationIdentificationV2 ::= SEQUENCE {
//	    noOfBorderCrossingRecords    INTEGER(0..255),
//	    noOfLoadUnloadRecords        INTEGER(0..255),
//	    noOfLoadTypeEntryRecords     INTEGER(0..255),
//	    vuConfigurationLengthRange   INTEGER(0..255)
//	}
func unmarshalCardApplicationIdentificationV2(data []byte) (*cardv1.ApplicationIdentificationV2, error) {
	const (
		lenEfApplicationIdentificationV2 = 4 // EF_ApplicationIdentificationV2 record size
	)

	if len(data) < lenEfApplicationIdentificationV2 {
		return nil, fmt.Errorf("insufficient data for application identification V2: got %d bytes, need %d", len(data), lenEfApplicationIdentificationV2)
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

// AppendCardApplicationIdentificationV2 appends application identification V2 data to a byte slice.
//
// The data type `ApplicationIdentificationV2` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition:
//
//	ApplicationIdentificationV2 ::= SEQUENCE {
//	    noOfBorderCrossingRecords    INTEGER(0..255),
//	    noOfLoadUnloadRecords        INTEGER(0..255),
//	    noOfLoadTypeEntryRecords     INTEGER(0..255),
//	    vuConfigurationLengthRange   INTEGER(0..255)
//	}
func appendCardApplicationIdentificationV2(data []byte, appIdV2 *cardv1.ApplicationIdentificationV2) ([]byte, error) {
	if appIdV2 == nil {
		return data, nil
	}

	// Get the appropriate nested message based on card type
	var borderCrossingRecords, loadUnloadRecords, loadTypeEntryRecords, vuConfigLength int32

	switch appIdV2.GetCardType() {
	case cardv1.CardType_DRIVER_CARD:
		if driver := appIdV2.GetDriver(); driver != nil {
			borderCrossingRecords = driver.GetBorderCrossingRecordsCount()
			loadUnloadRecords = driver.GetLoadUnloadRecordsCount()
			loadTypeEntryRecords = driver.GetLoadTypeEntryRecordsCount()
			vuConfigLength = driver.GetVuConfigurationLengthRange()
		}
	case cardv1.CardType_WORKSHOP_CARD:
		if workshop := appIdV2.GetWorkshop(); workshop != nil {
			borderCrossingRecords = workshop.GetBorderCrossingRecordsCount()
			loadUnloadRecords = workshop.GetLoadUnloadRecordsCount()
			loadTypeEntryRecords = workshop.GetLoadTypeEntryRecordsCount()
			vuConfigLength = workshop.GetVuConfigurationLengthRange()
		}
	case cardv1.CardType_COMPANY_CARD:
		if company := appIdV2.GetCompany(); company != nil {
			vuConfigLength = company.GetVuConfigurationLengthRange()
		}
	case cardv1.CardType_CONTROL_CARD:
		if control := appIdV2.GetControl(); control != nil {
			vuConfigLength = control.GetVuConfigurationLengthRange()
		}
	}

	// Border crossing records count (1 byte)
	data = append(data, byte(borderCrossingRecords))

	// Load/unload records count (1 byte)
	data = append(data, byte(loadUnloadRecords))

	// Load type entry records count (1 byte)
	data = append(data, byte(loadTypeEntryRecords))

	// VU configuration length range (1 byte)
	data = append(data, byte(vuConfigLength))

	return data, nil
}
