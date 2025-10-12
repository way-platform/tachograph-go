package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalEventsAndFaultsGen2V1 parses Gen2 V1 Events and Faults data from the complete transfer value.
//
// Gen2 V1 Events and Faults structure uses RecordArray format.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalEventsAndFaultsGen2V1(value []byte) (*vuv1.EventsAndFaultsGen2V1, error) {
	eventsAndFaults := &vuv1.EventsAndFaultsGen2V1{}
	eventsAndFaults.SetRawData(value)

	// Validate structure by skipping through all record arrays
	offset := 0
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(value, offset)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		offset += size
		return nil
	}

	// Skip all record arrays
	// VuFaultRecordArray
	if err := skipRecordArray("VuFault"); err != nil {
		return nil, err
	}
	// VuEventRecordArray
	if err := skipRecordArray("VuEvent"); err != nil {
		return nil, err
	}
	// VuOverSpeedingControlRecordArray
	if err := skipRecordArray("VuOverSpeedingControl"); err != nil {
		return nil, err
	}
	// VuTimeAdjustmentRecordArray
	if err := skipRecordArray("VuTimeAdjustment"); err != nil {
		return nil, err
	}
	// SignatureRecordArray
	if err := skipRecordArray("Signature"); err != nil {
		return nil, err
	}

	if offset != len(value) {
		return nil, fmt.Errorf("Events and Faults Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	return eventsAndFaults, nil
}

// appendEventsAndFaultsGen2V1 marshals Gen2 V1 Events and Faults data using raw data painting.
func appendEventsAndFaultsGen2V1(dst []byte, eventsAndFaults *vuv1.EventsAndFaultsGen2V1) ([]byte, error) {
	if eventsAndFaults == nil {
		return nil, fmt.Errorf("eventsAndFaults cannot be nil")
	}

	raw := eventsAndFaults.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	return nil, fmt.Errorf("cannot marshal Events and Faults Gen2 V1 without raw_data (semantic marshalling not yet implemented)")
}
