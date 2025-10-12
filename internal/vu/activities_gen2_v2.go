package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalActivitiesGen2V2 parses Gen2 V2 Activities data from the complete transfer value.
//
// Gen2 V2 Activities structure is identical to Gen2 V1 (from Data Dictionary):
//
// ASN.1 Definition:
//
//	VuActivitiesSecondGenV2 ::= SEQUENCE {
//	    timeRealRecordArray                   TimeRealRecordArray,
//	    odometerValueMidnightRecordArray      OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                   VuCardIWRecordArray,
//	    vuActivityDailyRecordArray            VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray     VuPlaceDailyWorkPeriodRecordArray,
//	    vuSpecificConditionRecordArray        VuSpecificConditionRecordArray,
//	    signatureRecordArray                  SignatureRecordArray
//	}
//
// Each RecordArray has a 5-byte header:
//
//	recordType (1 byte) + recordSize (2 bytes, big-endian) + noOfRecords (2 bytes, big-endian)
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
// Full semantic parsing of all RecordArrays is TODO.
func unmarshalActivitiesGen2V2(value []byte) (*vuv1.ActivitiesGen2V2, error) {
	activities := &vuv1.ActivitiesGen2V2{}
	activities.SetRawData(value)

	// For now, store the raw data and validate structure by skipping through all record arrays
	offset := 0

	// Helper to skip a RecordArray
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(value, offset)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		offset += size
		return nil
	}

	// TimeRealRecordArray
	if err := skipRecordArray("TimeReal"); err != nil {
		return nil, err
	}

	// OdometerValueMidnightRecordArray
	if err := skipRecordArray("OdometerValueMidnight"); err != nil {
		return nil, err
	}

	// VuCardIWRecordArray
	if err := skipRecordArray("VuCardIW"); err != nil {
		return nil, err
	}

	// VuActivityDailyRecordArray
	if err := skipRecordArray("VuActivityDaily"); err != nil {
		return nil, err
	}

	// VuPlaceDailyWorkPeriodRecordArray
	if err := skipRecordArray("VuPlaceDailyWorkPeriod"); err != nil {
		return nil, err
	}

	// VuSpecificConditionRecordArray
	if err := skipRecordArray("VuSpecificCondition"); err != nil {
		return nil, err
	}

	// SignatureRecordArray (last)
	if err := skipRecordArray("Signature"); err != nil {
		return nil, err
	}

	// Verify we consumed exactly the right amount of data
	if offset != len(value) {
		return nil, fmt.Errorf("Activities Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	// TODO: Implement full semantic parsing of all record arrays
	// For now, raw_data contains all the information needed for round-trip testing

	return activities, nil
}

// appendActivitiesGen2V2 marshals Gen2 V2 Activities data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and matches the structure, it uses it as the output.
func appendActivitiesGen2V2(dst []byte, activities *vuv1.ActivitiesGen2V2) ([]byte, error) {
	if activities == nil {
		return nil, fmt.Errorf("activities cannot be nil")
	}

	// For Gen2 structures with RecordArrays, raw data painting is straightforward
	raw := activities.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	// TODO: Implement marshalling from semantic fields
	// This would require constructing all RecordArrays from semantic data
	return nil, fmt.Errorf("cannot marshal Activities Gen2 V2 without raw_data (semantic marshalling not yet implemented)")
}
