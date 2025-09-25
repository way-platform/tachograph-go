package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func unmarshalCardCurrentUsage(data []byte) (*cardv1.CurrentUsage, error) {
	if len(data) < 19 {
		return nil, fmt.Errorf("insufficient data for current usage")
	}
	var target cardv1.CurrentUsage
	r := bytes.NewReader(data)
	// Read session open time (4 bytes)
	target.SetSessionOpenTime(readTimeReal(r))
	// Read session open vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	var nation byte
	if err := binary.Read(r, binary.BigEndian, &nation); err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	target.SetVehicleRegistrationNation(fmt.Sprintf("%02X", nation))
	registrationBytes := make([]byte, 14)
	if _, err := r.Read(registrationBytes); err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	target.SetVehicleRegistrationNumber(readString(bytes.NewReader(registrationBytes), 14))
	return &target, nil
}
