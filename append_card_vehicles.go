package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendVehiclesUsed appends the binary representation of VehiclesUsed to dst.
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

// AppendVehicleRecord appends a single 31-byte vehicle record (Gen1).
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
