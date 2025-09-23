package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardApplicationIdentificationV2 unmarshals application identification V2 data from a card EF.
func UnmarshalCardApplicationIdentificationV2(data []byte, target *cardv1.ApplicationIdentificationV2) error {
	if len(data) < 4 { // Minimum size: 1 + 1 + 1 + 1 = 4 bytes
		return fmt.Errorf("insufficient data for application identification V2")
	}

	r := bytes.NewReader(data)

	// Read border crossing records count (1 byte)
	var borderCrossingCount byte
	if err := binary.Read(r, binary.BigEndian, &borderCrossingCount); err != nil {
		return fmt.Errorf("failed to read border crossing records count: %w", err)
	}
	target.SetBorderCrossingRecordsCount(int32(borderCrossingCount))

	// Read load/unload records count (1 byte)
	var loadUnloadCount byte
	if err := binary.Read(r, binary.BigEndian, &loadUnloadCount); err != nil {
		return fmt.Errorf("failed to read load/unload records count: %w", err)
	}
	target.SetLoadUnloadRecordsCount(int32(loadUnloadCount))

	// Read load type entry records count (1 byte)
	var loadTypeCount byte
	if err := binary.Read(r, binary.BigEndian, &loadTypeCount); err != nil {
		return fmt.Errorf("failed to read load type entry records count: %w", err)
	}
	target.SetLoadTypeEntryRecordsCount(int32(loadTypeCount))

	// Read VU configuration length range (1 byte)
	var vuConfigRange byte
	if err := binary.Read(r, binary.BigEndian, &vuConfigRange); err != nil {
		return fmt.Errorf("failed to read VU configuration length range: %w", err)
	}
	target.SetVuConfigurationLengthRange(int32(vuConfigRange))

	return nil
}
