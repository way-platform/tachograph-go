package tachograph

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
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
		vehicleUnitFile, err := vu.UnmarshalVehicleUnitFile(data)
		if err != nil {
			return nil, err
		}
		output.SetType(tachographv1.File_VEHICLE_UNIT)
		output.SetVehicleUnit(vehicleUnitFile)
		return &output, nil

	// Card file (starts with EF_ICC prefix).
	case binary.BigEndian.Uint16(data[0:2]) == 0x0002:
		rawCardFile, err := card.UnmarshalRawCardFile(data)
		if err != nil {
			return nil, err
		}

		// Infer the card type
		cardType := card.InferCardFileType(rawCardFile)

		// Parse structured card data based on type
		switch cardType {
		case cardv1.CardType_DRIVER_CARD:
			driverCard, err := card.UnmarshalDriverCardFile(rawCardFile)
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver card: %w", err)
			}
			output.SetType(tachographv1.File_DRIVER_CARD)
			output.SetDriverCard(driverCard)
			return &output, nil
		default:
			// For unsupported card types, return raw card data
			output.SetType(tachographv1.File_RAW_CARD)
			output.SetRawCard(rawCardFile)
			return &output, nil
		}

	default:
		return nil, errors.New("unknown or unsupported file type")
	}
}

// MarshalFile serializes a protobuf File message into the binary DDD file format.
func MarshalFile(file *tachographv1.File) ([]byte, error) {
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		return card.MarshalDriverCardFile(file.GetDriverCard())
	case tachographv1.File_VEHICLE_UNIT:
		return vu.MarshalVehicleUnitFile(file.GetVehicleUnit())
	case tachographv1.File_RAW_CARD:
		return card.MarshalRawCardFile(file.GetRawCard())
	default:
		return nil, fmt.Errorf("unsupported file type for marshaling: %v", file.GetType())
	}
}
