package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCurrentUsage appends the binary representation of CurrentUsage to dst.
func AppendCurrentUsage(dst []byte, cu *cardv1.CurrentUsage) ([]byte, error) {
	if cu == nil {
		return dst, nil
	}
	var err error
	dst = appendTimeReal(dst, cu.GetSessionOpenTime())
	dst, err = appendVehicleRegistration(dst, cu.GetSessionOpenVehicle())
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// AppendControlActivityData appends the binary representation of ControlActivityData to dst.
func AppendControlActivityData(dst []byte, cad *cardv1.ControlActivityData) ([]byte, error) {
	if cad == nil {
		return dst, nil
	}
	dst = appendControlType(dst, cad.GetControlType())
	dst = appendTimeReal(dst, cad.GetControlTime())
	// TODO: Append FullCardNumber and VehicleRegistrationIdentification correctly
	dst = append(dst, make([]byte, 18+15)...)
	dst = appendTimeReal(dst, cad.GetControlDownloadPeriodBegin())
	dst = appendTimeReal(dst, cad.GetControlDownloadPeriodEnd())
	return dst, nil
}

// AppendSpecificConditions appends the binary representation of SpecificConditions to dst.
func AppendSpecificConditions(dst []byte, sc *cardv1.SpecificConditions) ([]byte, error) {
	const totalRecords = 56
	var err error
	for i := 0; i < totalRecords; i++ {
		if i < len(sc.GetRecords()) {
			dst, err = AppendSpecificConditionRecord(dst, sc.GetRecords()[i])
			if err != nil {
				return nil, err
			}
		} else {
			dst = append(dst, make([]byte, 5)...) // Pad with empty records
		}
	}
	return dst, nil
}

// AppendSpecificConditionRecord appends a single 5-byte specific condition record.
func AppendSpecificConditionRecord(dst []byte, rec *cardv1.SpecificConditions_Record) ([]byte, error) {
	dst = appendTimeReal(dst, rec.GetEntryTime())
	dst = append(dst, byte(rec.GetSpecificConditionType()))
	return dst, nil
}

// AppendLastCardDownload appends the binary representation of LastCardDownload to dst.
func AppendLastCardDownload(dst []byte, lcd *cardv1.CardDownloadDriver) ([]byte, error) {
	if lcd == nil {
		return dst, nil
	}
	return appendTimeReal(dst, lcd.GetTimestamp()), nil
}
