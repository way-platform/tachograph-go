package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCurrentUsage appends the binary representation of CurrentUsage to dst.
//
// ASN.1 Specification (Data Dictionary 2.16):
//
//	CardCurrentUse ::= SEQUENCE {
//	    sessionOpenTime                   TimeReal,
//	    sessionOpenVehicle                VehicleRegistrationIdentification
//	}
//
// Binary Layout (19 bytes):
//
//	0-3:   sessionOpenTime (4 bytes, TimeReal)
//	4-18:  sessionOpenVehicle (15 bytes: 1 byte nation + 14 bytes number)
//

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
//
// ASN.1 Specification (Data Dictionary 2.15):
//
//	CardControlActivityDataRecord ::= SEQUENCE {
//	    controlType                        ControlType,
//	    controlTime                        TimeReal,
//	    controlCardNumber                  FullCardNumber,
//	    controlVehicleRegistration         VehicleRegistrationIdentification,
//	    controlDownloadPeriodBegin         TimeReal,
//	    controlDownloadPeriodEnd           TimeReal
//	}
//
// Binary Layout (46 bytes):
//
//	0-0:   controlType (1 byte)
//	1-4:   controlTime (4 bytes, TimeReal)
//	5-22:  controlCardNumber (18 bytes, FullCardNumber)
//	23-37: controlVehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	38-41: controlDownloadPeriodBegin (4 bytes, TimeReal)
//	42-45: controlDownloadPeriodEnd (4 bytes, TimeReal)
//

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
//
// ASN.1 Specification (Data Dictionary 2.153):
//
//	SpecificConditions ::= SEQUENCE {
//	    conditionPointerNewestRecord      INTEGER(0..NoOfSpecificConditionRecords-1),
//	    specificConditionRecords          SET SIZE(NoOfSpecificConditionRecords) OF SpecificConditionRecord
//	}
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime                          TimeReal,
//	    specificConditionType              SpecificConditionType
//	}
//
// Binary Layout (variable size):
//
//	Each record: 5 bytes (4 bytes time + 1 byte condition type)
//	Total: 56 records * 5 bytes = 280 bytes
//

func AppendSpecificConditions(dst []byte, sc *cardv1.SpecificConditions) ([]byte, error) {
	var err error
	for i := 0; i < specificConditionTotalRecords; i++ {
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

// AppendSpecificConditionRecord appends a single specific condition record to dst.
//
// ASN.1 Specification (Data Dictionary 2.152):
//
//	SpecificConditionRecord ::= SEQUENCE {
//	    entryTime                          TimeReal,
//	    specificConditionType              SpecificConditionType
//	}
//
// Binary Layout (5 bytes):
//
//	0-3:   entryTime (4 bytes, TimeReal)
//	4-4:   specificConditionType (1 byte)
func AppendSpecificConditionRecord(dst []byte, rec *cardv1.SpecificConditions_Record) ([]byte, error) {
	dst = appendTimeReal(dst, rec.GetEntryTime())
	dst = append(dst, byte(rec.GetSpecificConditionType()))
	return dst, nil
}

// AppendLastCardDownload appends the binary representation of LastCardDownload to dst.
//
// ASN.1 Specification (Data Dictionary 2.16):
//
//	CardDownloadDriver ::= TimeReal
//
// Binary Layout (4 bytes):
//
//	0-3:   timestamp (4 bytes, TimeReal)
//

func AppendLastCardDownload(dst []byte, lcd *cardv1.CardDownloadDriver) ([]byte, error) {
	if lcd == nil {
		return dst, nil
	}
	return appendTimeReal(dst, lcd.GetTimestamp()), nil
}
