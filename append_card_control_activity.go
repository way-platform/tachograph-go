package tachograph

import (
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
	var controlTypeByte byte
	if controlType != nil {
		// Build bitmask from boolean fields
		// Structure: 'cvpdexxx'B
		// - 'c': card downloading
		// - 'v': VU downloading
		// - 'p': printing
		// - 'd': display
		// - 'e': calibration checking
		if controlType.GetCardDownloading() {
			controlTypeByte |= 0x80 // bit 7
		}
		if controlType.GetVuDownloading() {
			controlTypeByte |= 0x40 // bit 6
		}
		if controlType.GetPrinting() {
			controlTypeByte |= 0x20 // bit 5
		}
		if controlType.GetDisplay() {
			controlTypeByte |= 0x10 // bit 4
		}
		if controlType.GetCalibrationChecking() {
			controlTypeByte |= 0x08 // bit 3
		}
	}
	data = append(data, controlTypeByte)

	// Control time (4 bytes)
	data = appendTimeReal(data, controlData.GetControlTime())

	var err error
	// Control card number (18 bytes)
	data, err = appendFullCardNumber(data, controlData.GetControlCardNumber(), 18)
	if err != nil {
		return nil, err
	}

	// Vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	data, err = appendVehicleRegistration(data, controlData.GetControlVehicleRegistration())
	if err != nil {
		return nil, err
	}

	// Control download period begin (4 bytes)
	data = appendTimeReal(data, controlData.GetControlDownloadPeriodBegin())

	// Control download period end (4 bytes)
	data = appendTimeReal(data, controlData.GetControlDownloadPeriodEnd())

	return data, nil
}
