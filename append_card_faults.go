package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendFaultsData appends the binary representation of FaultData to dst.
//
// ASN.1 Specification (Data Dictionary 2.22):
//
//	CardFaultData ::= SEQUENCE OF CardFaultRecord
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                         EventFaultType,
//	    faultBeginTime                    TimeReal,
//	    faultEndTime                      TimeReal,
//	    faultVehicleRegistration          VehicleRegistrationIdentification,
//	    faultTypeSpecificData             OCTET STRING (SIZE (2))
//	}
//
// Binary Layout (variable size):
//
//	Each record: 35 bytes
//	- 0-0:   faultType (1 byte)
//	- 1-4:   faultBeginTime (4 bytes, TimeReal)
//	- 5-8:   faultEndTime (4 bytes, TimeReal)
//	- 9-23:  faultVehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	- 24-25: faultTypeSpecificData (2 bytes)
//

func AppendFaultsData(dst []byte, data *cardv1.FaultsData) ([]byte, error) {
	var err error
	for _, r := range data.GetRecords() {
		dst, err = AppendFaultRecord(dst, r)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendFaultRecord appends a single fault record to dst.
//
// ASN.1 Specification (Data Dictionary 2.22):
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                         EventFaultType,
//	    faultBeginTime                    TimeReal,
//	    faultEndTime                      TimeReal,
//	    faultVehicleRegistration          VehicleRegistrationIdentification,
//	    faultTypeSpecificData             OCTET STRING (SIZE (2))
//	}
//
// Binary Layout (35 bytes):
//
//	0-0:   faultType (1 byte)
//	1-4:   faultBeginTime (4 bytes, TimeReal)
//	5-8:   faultEndTime (4 bytes, TimeReal)
//	9-23:  faultVehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	24-25: faultTypeSpecificData (2 bytes)
func AppendFaultRecord(dst []byte, record *cardv1.FaultsData_Record) ([]byte, error) {
	if !record.GetValid() {
		return append(dst, record.GetRawData()...), nil
	}

	protocolValue := GetEventFaultTypeProtocolValue(record.GetFaultType(), 0)
	dst = append(dst, byte(protocolValue))
	dst = appendTimeReal(dst, record.GetFaultBeginTime())
	dst = appendTimeReal(dst, record.GetFaultEndTime())
	dst, err := appendVehicleRegistration(dst, record.GetFaultVehicleRegistration())
	if err != nil {
		return nil, err
	}

	return dst, nil
}
