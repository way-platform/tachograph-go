package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardControlActivityData unmarshals control activity data from a card EF.
func UnmarshalCardControlActivityData(data []byte, target *cardv1.ControlActivityData) error {
	if len(data) < 46 { // Minimum size: 1 + 4 + 18 + 15 + 4 + 4 = 46 bytes
		return fmt.Errorf("insufficient data for control activity data")
	}

	r := bytes.NewReader(data)

	// Read control type (1 byte)
	var controlType byte
	if err := binary.Read(r, binary.BigEndian, &controlType); err != nil {
		return fmt.Errorf("failed to read control type: %w", err)
	}
	target.SetControlType([]byte{controlType})

	// Read control time (4 bytes)
	target.SetControlTime(readTimeReal(r))

	// Read control card number (18 bytes)
	controlCardBytes := make([]byte, 18)
	if _, err := r.Read(controlCardBytes); err != nil {
		return fmt.Errorf("failed to read control card number: %w", err)
	}
	target.SetControlCardNumber(readString(bytes.NewReader(controlCardBytes), 18))

	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	var nation byte
	if err := binary.Read(r, binary.BigEndian, &nation); err != nil {
		return fmt.Errorf("failed to read vehicle registration nation: %w", err)
	}
	target.SetVehicleRegistrationNation(fmt.Sprintf("%02X", nation))

	registrationBytes := make([]byte, 14)
	if _, err := r.Read(registrationBytes); err != nil {
		return fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	target.SetVehicleRegistrationNumber(readString(bytes.NewReader(registrationBytes), 14))

	// Read control download period begin (4 bytes)
	target.SetControlDownloadPeriodBegin(readTimeReal(r))

	// Read control download period end (4 bytes)
	target.SetControlDownloadPeriodEnd(readTimeReal(r))

	return nil
}
