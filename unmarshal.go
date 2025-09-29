package tachograph

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// UnmarshalFile parses a .DDD file's byte data into a protobuf File message.
func UnmarshalFile(data []byte) (*tachographv1.File, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for tachograph file: %w", io.ErrUnexpectedEOF)
	}
	var output tachographv1.File
	switch {

	// Vehicle unit file (starts with TREP prefix).
	case data[0] == 0x76:
		vehicleUnitFile, err := unmarshalVehicleUnitFile(data)
		if err != nil {
			return nil, err
		}
		output.SetType(tachographv1.File_VEHICLE_UNIT)
		output.SetVehicleUnit(vehicleUnitFile)
		return &output, nil

	// Card file (starts with EF_ICC prefix).
	case binary.BigEndian.Uint16(data[0:2]) == 0x0002:
		rawCardFile, err := unmarshalRawCardFile(data)
		if err != nil {
			return nil, err
		}
		switch fileType := inferCardFileType(rawCardFile); fileType {
		case cardv1.CardType_DRIVER_CARD:
			driverCardFile, err := unmarshalDriverCardFile(rawCardFile)
			if err != nil {
				return nil, err
			}
			output.SetType(tachographv1.File_DRIVER_CARD)
			output.SetDriverCard(driverCardFile)
		case cardv1.CardType_WORKSHOP_CARD:
			// TODO: Implement workshop card.
			fallthrough
		case cardv1.CardType_CONTROL_CARD:
			// TODO: Implement control card.
			fallthrough
		case cardv1.CardType_COMPANY_CARD:
			// TODO: Implement company card.
			fallthrough
		default:
			output.SetType(tachographv1.File_RAW_CARD)
			output.SetRawCard(rawCardFile)
		}
		return &output, nil
	default:
		return nil, errors.New("unknown or unsupported file type")
	}
}

// MarshalFile serializes a protobuf File message into the binary DDD file format.
func MarshalFile(file *tachographv1.File) ([]byte, error) {
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		return appendCard(nil, file)
	case tachographv1.File_VEHICLE_UNIT:
		return AppendVU(nil, file.GetVehicleUnit())
	case tachographv1.File_RAW_CARD:
		return appendCard(nil, file) // Raw cards use the same format as driver cards
	default:
		return nil, fmt.Errorf("unsupported file type for marshaling: %v", file.GetType())
	}
}
