package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardApplicationIdentificationV2 appends application identification V2 data to a byte slice.
func AppendCardApplicationIdentificationV2(data []byte, appIdV2 *cardv1.ApplicationIdentificationV2) ([]byte, error) {
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
