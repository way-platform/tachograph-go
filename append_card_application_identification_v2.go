package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardApplicationIdentificationV2 appends application identification V2 data to a byte slice.
func AppendCardApplicationIdentificationV2(data []byte, appIdV2 *cardv1.ApplicationIdentificationV2) ([]byte, error) {
	if appIdV2 == nil {
		return data, nil
	}

	// Border crossing records count (1 byte)
	if appIdV2.HasBorderCrossingRecordsCount() {
		data = append(data, byte(appIdV2.GetBorderCrossingRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// Load/unload records count (1 byte)
	if appIdV2.HasLoadUnloadRecordsCount() {
		data = append(data, byte(appIdV2.GetLoadUnloadRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// Load type entry records count (1 byte)
	if appIdV2.HasLoadTypeEntryRecordsCount() {
		data = append(data, byte(appIdV2.GetLoadTypeEntryRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	// VU configuration length range (1 byte)
	if appIdV2.HasVuConfigurationLengthRange() {
		data = append(data, byte(appIdV2.GetVuConfigurationLengthRange()))
	} else {
		data = append(data, 0x00)
	}

	return data, nil
}
