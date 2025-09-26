package tachograph

import (
	"encoding/binary"
	"encoding/hex"

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
	// Convert hex string nation back to byte
	nationByte := byte(0) // Default fallback
	if nationStr := rec.GetVehicleRegistrationNation(); len(nationStr) > 0 {
		if b, err := hex.DecodeString(nationStr); err == nil && len(b) > 0 {
			nationByte = b[0]
		}
	}
	dst = append(dst, nationByte)

	var err error
	dst, err = appendString(dst, rec.GetVehicleRegistrationNumber(), 14)
	if err != nil {
		return nil, err
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetVuDataBlockCounter()))
	return dst, nil
}
