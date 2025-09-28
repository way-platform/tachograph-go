package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendVehiclesUsed appends the binary representation of VehiclesUsed to dst.
//
// ASN.1 Specification (Data Dictionary 2.6):
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBeginKm       OdometerShort,
//	    vehicleOdometerEndKm         OdometerShort,
//	    vehicleFirstUse              TimeReal,
//	    vehicleLastUse               TimeReal,
//	    vehicleRegistration          VehicleRegistrationIdentification,
//	    vehicleIdentificationNumber  VehicleIdentificationNumber OPTIONAL,
//	    vehicleRegistrationNation    NationNumeric OPTIONAL,
//	    vehicleRegistrationNumber    RegistrationNumber OPTIONAL
//	}
//
// Binary Layout (variable size):
//
//	0-1:   newestRecordIndex (2 bytes, big-endian)
//	2+:    vehicle records (31 bytes Gen1, 48 bytes Gen2 each)
//	  - 0-3:   vehicleOdometerBeginKm (4 bytes, big-endian)
//	  - 4-7:   vehicleOdometerEndKm (4 bytes, big-endian)
//	  - 8-11:  vehicleFirstUse (4 bytes, TimeReal)
//	  - 12-15: vehicleLastUse (4 bytes, TimeReal)
//	  - 16-30: vehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	  - 31+:   vehicleIdentificationNumber (17 bytes, Gen2 only)
//	  - 48+:   vehicleRegistrationNation (1 byte, Gen2 only)
//	  - 49+:   vehicleRegistrationNumber (14 bytes, Gen2 only)
//

func AppendVehiclesUsed(dst []byte, vu *cardv1.VehiclesUsed) ([]byte, error) {
	if vu == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(vu.GetNewestRecordIndex()))

	var err error
	for _, rec := range vu.GetRecords() {
		dst, err = AppendVehicleRecord(dst, rec)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendVehicleRecord appends a single vehicle record to dst.
//
// ASN.1 Specification (Data Dictionary 2.6):
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBeginKm       OdometerShort,
//	    vehicleOdometerEndKm         OdometerShort,
//	    vehicleFirstUse              TimeReal,
//	    vehicleLastUse               TimeReal,
//	    vehicleRegistration          VehicleRegistrationIdentification,
//	    vehicleIdentificationNumber  VehicleIdentificationNumber OPTIONAL,
//	    vehicleRegistrationNation    NationNumeric OPTIONAL,
//	    vehicleRegistrationNumber    RegistrationNumber OPTIONAL
//	}
//
// Binary Layout (31 bytes Gen1, 48 bytes Gen2):
//
//	0-3:   vehicleOdometerBeginKm (4 bytes, big-endian)
//	4-7:   vehicleOdometerEndKm (4 bytes, big-endian)
//	8-11:  vehicleFirstUse (4 bytes, TimeReal)
//	12-15: vehicleLastUse (4 bytes, TimeReal)
//	16-30: vehicleRegistration (15 bytes: 1 byte nation + 14 bytes number)
//	31+:   vehicleIdentificationNumber (17 bytes, Gen2 only)
//	48+:   vehicleRegistrationNation (1 byte, Gen2 only)
//	49+:   vehicleRegistrationNumber (14 bytes, Gen2 only)
func AppendVehicleRecord(dst []byte, rec *cardv1.VehiclesUsed_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, 31)...), nil
	}
	dst = appendOdometer(dst, uint32(rec.GetVehicleOdometerBeginKm()))
	dst = appendOdometer(dst, uint32(rec.GetVehicleOdometerEndKm()))
	dst = appendTimeReal(dst, rec.GetVehicleFirstUse())
	dst = appendTimeReal(dst, rec.GetVehicleLastUse())
	// Vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	var err error
	dst, err = appendVehicleRegistration(dst, rec.GetVehicleRegistration())
	if err != nil {
		return nil, err
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetVuDataBlockCounter()))
	return dst, nil
}
