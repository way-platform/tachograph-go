package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func unmarshalCardApplicationIdentificationV2(data []byte) (*cardv1.ApplicationIdentificationV2, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for application identification V2")
	}
	var target cardv1.ApplicationIdentificationV2
	r := bytes.NewReader(data)

	// For now, assume this is a driver card and create the driver data
	driver := &cardv1.ApplicationIdentificationV2_Driver{}

	// Read border crossing records count (1 byte)
	var borderCrossingCount byte
	if err := binary.Read(r, binary.BigEndian, &borderCrossingCount); err != nil {
		return nil, fmt.Errorf("failed to read border crossing records count: %w", err)
	}
	driver.SetBorderCrossingRecordsCount(int32(borderCrossingCount))

	// Read load/unload records count (1 byte)
	var loadUnloadCount byte
	if err := binary.Read(r, binary.BigEndian, &loadUnloadCount); err != nil {
		return nil, fmt.Errorf("failed to read load/unload records count: %w", err)
	}
	driver.SetLoadUnloadRecordsCount(int32(loadUnloadCount))

	// Read load type entry records count (1 byte)
	var loadTypeCount byte
	if err := binary.Read(r, binary.BigEndian, &loadTypeCount); err != nil {
		return nil, fmt.Errorf("failed to read load type entry records count: %w", err)
	}
	driver.SetLoadTypeEntryRecordsCount(int32(loadTypeCount))

	// Read VU configuration length range (1 byte)
	var vuConfigRange byte
	if err := binary.Read(r, binary.BigEndian, &vuConfigRange); err != nil {
		return nil, fmt.Errorf("failed to read VU configuration length range: %w", err)
	}
	driver.SetVuConfigurationLengthRange(int32(vuConfigRange))

	// Set the driver data and card type
	target.SetDriver(driver)
	target.SetCardType(cardv1.CardType_DRIVER_CARD)

	return &target, nil
}
