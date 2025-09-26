package tachograph

import (
	"strconv"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardControlActivityData appends control activity data to a byte slice.
func AppendCardControlActivityData(data []byte, controlData *cardv1.ControlActivityData) ([]byte, error) {
	if controlData == nil {
		return data, nil
	}

	if !controlData.GetValid() {
		// Non-valid record: use preserved raw data
		rawData := controlData.GetRawData()
		if len(rawData) != 46 {
			// Fallback to zeros if raw data is invalid
			return append(data, make([]byte, 46)...), nil
		}
		return append(data, rawData...), nil
	}

	// Valid record: serialize semantic data
	// Control type (1 byte)
	controlType := controlData.GetControlType()
	if len(controlType) > 0 {
		data = append(data, controlType[0])
	} else {
		data = append(data, 0x00)
	}

	// Control time (4 bytes)
	data = appendTimeReal(data, controlData.GetControlTime())

	var err error
	// Control card number (18 bytes)
	cardNumber := controlData.GetControlCardNumber()
	data, err = appendString(data, cardNumber, 18)
	if err != nil {
		return nil, err
	}

	// Vehicle registration nation (1 byte)
	nationStr := controlData.GetVehicleRegistrationNation()
	var nationByte byte = 0x00
	if len(nationStr) >= 2 {
		if val, err := strconv.ParseUint(nationStr, 16, 8); err == nil {
			nationByte = byte(val)
		}
	}
	data = append(data, nationByte)

	// Vehicle registration number (14 bytes)
	regNumber := controlData.GetVehicleRegistrationNumber()
	data, err = appendString(data, regNumber, 14)
	if err != nil {
		return nil, err
	}

	// Control download period begin (4 bytes)
	data = appendTimeReal(data, controlData.GetControlDownloadPeriodBegin())

	// Control download period end (4 bytes)
	data = appendTimeReal(data, controlData.GetControlDownloadPeriodEnd())

	return data, nil
}
