package card

import (
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalApplicationIdentificationV2 parses the binary data for an EF_ApplicationIdentificationV2 record.
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
func (opts UnmarshalOptions) unmarshalApplicationIdentificationV2(
	data []byte, cardType cardv1.CardType,
) (*cardv1.ApplicationIdentificationV2, error) {
	const lenEfApplicationIdentificationV2 = 4

	if len(data) < lenEfApplicationIdentificationV2 {
		return nil, fmt.Errorf("insufficient data for application identification V2: got %d bytes, need %d", len(data), lenEfApplicationIdentificationV2)
	}
	var target cardv1.ApplicationIdentificationV2
	switch cardType {
	case cardv1.CardType_DRIVER_CARD:
		driver := &cardv1.ApplicationIdentificationV2_Driver{}
		driver.SetBorderCrossingRecordsCount(int32(data[0]))
		driver.SetLoadUnloadRecordsCount(int32(data[1]))
		driver.SetLoadTypeEntryRecordsCount(int32(data[2]))
		driver.SetVuConfigurationLengthRange(int32(data[3]))
		target.SetDriver(driver)
	case cardv1.CardType_WORKSHOP_CARD:
		workshop := &cardv1.ApplicationIdentificationV2_Workshop{}
		workshop.SetBorderCrossingRecordsCount(int32(data[0]))
		workshop.SetLoadUnloadRecordsCount(int32(data[1]))
		workshop.SetLoadTypeEntryRecordsCount(int32(data[2]))
		workshop.SetVuConfigurationLengthRange(int32(data[3]))
		target.SetWorkshop(workshop)
	case cardv1.CardType_COMPANY_CARD:
		company := &cardv1.ApplicationIdentificationV2_Company{}
		company.SetVuConfigurationLengthRange(int32(data[3]))
		target.SetCompany(company)
	case cardv1.CardType_CONTROL_CARD:
		control := &cardv1.ApplicationIdentificationV2_Control{}
		control.SetVuConfigurationLengthRange(int32(data[3]))
		target.SetControl(control)
	}
	target.SetCardType(cardType)
	return &target, nil
}

// MarshalCardApplicationIdentificationV2 marshals application identification V2 data.
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
func (opts MarshalOptions) MarshalCardApplicationIdentificationV2(appIdV2 *cardv1.ApplicationIdentificationV2) ([]byte, error) {
	if appIdV2 == nil {
		return nil, nil
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

	var data []byte

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
