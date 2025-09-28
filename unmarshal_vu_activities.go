package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuActivities unmarshals VU activities data from a VU transfer.
func UnmarshalVuActivities(r *bytes.Reader, target *vuv1.Activities, generation int) (int, error) {
	switch generation {
	case 1:
		return unmarshalVuActivitiesGen1(r, target)
	case 2:
		return unmarshalVuActivitiesGen2(r, target)
	default:
		return 0, fmt.Errorf("unsupported generation: %d", generation)
	}
}

// unmarshalVuActivitiesGen1 unmarshals Generation 1 VU activities
func unmarshalVuActivitiesGen1(r *bytes.Reader, target *vuv1.Activities) (int, error) {
	initialLen := r.Len()
	target.SetGeneration(datadictionaryv1.Generation_GENERATION_1)

	// Read TimeReal (4 bytes) - this is the date of the day
	target.SetDateOfDay(readTimeReal(r))

	// Read OdometerValueMidnight (3 bytes)
	odometerBytes := make([]byte, 3)
	if _, err := r.Read(odometerBytes); err != nil {
		return 0, fmt.Errorf("failed to read odometer value midnight: %w", err)
	}
	target.SetOdometerMidnightKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBytes...))))

	// Parse VuCardIWData
	cardIWData, err := parseVuCardIWData(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW data: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyData
	activityChanges, err := parseVuActivityDailyData(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily data: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodData
	places, err := parseVuPlaceDailyWorkPeriodData(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period data: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuSpecificConditionData
	specificConditions, err := parseVuSpecificConditionData(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition data: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Read signature (128 bytes for Gen1)
	signatureBytes := make([]byte, 128)
	if _, err := r.Read(signatureBytes); err != nil {
		return 0, fmt.Errorf("failed to read signature: %w", err)
	}
	target.SetSignatureGen1(signatureBytes)

	return initialLen - r.Len(), nil
}

// unmarshalVuActivitiesGen2 unmarshals Generation 2 VU activities
func unmarshalVuActivitiesGen2(r *bytes.Reader, target *vuv1.Activities) (int, error) {
	initialLen := r.Len()
	target.SetGeneration(datadictionaryv1.Generation_GENERATION_2)

	// Gen2 format uses record arrays, each with a header
	// Parse DateOfDayDownloadedRecordArray
	dates, err := parseDateOfDayDownloadedRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse date of day downloaded record array: %w", err)
	}
	if len(dates) > 0 {
		target.SetDateOfDay(dates[0]) // Use first date
	}

	// Parse OdometerValueMidnightRecordArray
	odometerValues, err := parseOdometerValueMidnightRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse odometer value midnight record array: %w", err)
	}
	if len(odometerValues) > 0 {
		target.SetOdometerMidnightKm(odometerValues[0])
	}

	// Parse VuCardIWRecordArray
	cardIWData, err := parseVuCardIWRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW record array: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyRecordArray
	activityChanges, err := parseVuActivityDailyRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily record array: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodRecordArray
	places, err := parseVuPlaceDailyWorkPeriodRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period record array: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuGNSSADRecordArray (Gen2+)
	gnssRecords, err := parseVuGNSSADRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse GNSS AD record array: %w", err)
	}
	target.SetGnssAccumulatedDriving(gnssRecords)

	// Parse VuSpecificConditionRecordArray
	specificConditions, err := parseVuSpecificConditionRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition record array: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Try to parse Gen2v2 specific arrays if there's more data
	if r.Len() > 10 { // Need some minimum data for arrays
		// Parse VuBorderCrossingRecordArray (Gen2v2+)
		borderCrossings, err := parseVuBorderCrossingRecordArray(r)
		if err == nil {
			target.SetBorderCrossings(borderCrossings)
			target.SetVersion(vuv1.Version_VERSION_2)
		}

		// Parse VuLoadUnloadRecordArray (Gen2v2+)
		if r.Len() > 5 {
			loadUnloadRecords, err := parseVuLoadUnloadRecordArray(r)
			if err == nil {
				target.SetLoadUnloadOperations(loadUnloadRecords)
			}
		}
	}

	// Parse SignatureRecordArray
	signatureBytes, err := parseSignatureRecordArray(r)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signature record array: %w", err)
	}
	target.SetSignatureGen2(signatureBytes)

	return initialLen - r.Len(), nil
}

// Helper functions for parsing different record types
// These are simplified implementations - in a full implementation,
// each would need to properly handle the record array format

func parseVuCardIWData(r *bytes.Reader) ([]*vuv1.Activities_CardIWRecord, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*vuv1.Activities_CardIWRecord{}, nil
}

func parseVuActivityDailyData(r *bytes.Reader) ([]*datadictionaryv1.ActivityChangeInfo, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*datadictionaryv1.ActivityChangeInfo{}, nil
}

func parseVuPlaceDailyWorkPeriodData(r *bytes.Reader) ([]*vuv1.Activities_PlaceRecord, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*vuv1.Activities_PlaceRecord{}, nil
}

func parseVuSpecificConditionData(r *bytes.Reader) ([]*vuv1.Activities_SpecificConditionRecord, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*vuv1.Activities_SpecificConditionRecord{}, nil
}

// Gen2 record array parsers
func parseDateOfDayDownloadedRecordArray(r *bytes.Reader) ([]*timestamppb.Timestamp, error) {
	// Parse record array header and records
	return []*timestamppb.Timestamp{}, nil
}

func parseOdometerValueMidnightRecordArray(r *bytes.Reader) ([]int32, error) {
	// Parse record array header and records
	return []int32{}, nil
}

func parseVuCardIWRecordArray(r *bytes.Reader) ([]*vuv1.Activities_CardIWRecord, error) {
	// Parse record array header and records
	return []*vuv1.Activities_CardIWRecord{}, nil
}

func parseVuActivityDailyRecordArray(r *bytes.Reader) ([]*datadictionaryv1.ActivityChangeInfo, error) {
	// Parse record array header and records
	return []*datadictionaryv1.ActivityChangeInfo{}, nil
}

func parseVuPlaceDailyWorkPeriodRecordArray(r *bytes.Reader) ([]*vuv1.Activities_PlaceRecord, error) {
	// Parse record array header and records
	return []*vuv1.Activities_PlaceRecord{}, nil
}

func parseVuGNSSADRecordArray(r *bytes.Reader) ([]*vuv1.Activities_GnssRecord, error) {
	// Parse record array header and records
	return []*vuv1.Activities_GnssRecord{}, nil
}

func parseVuSpecificConditionRecordArray(r *bytes.Reader) ([]*vuv1.Activities_SpecificConditionRecord, error) {
	// Parse record array header and records
	return []*vuv1.Activities_SpecificConditionRecord{}, nil
}

func parseVuBorderCrossingRecordArray(r *bytes.Reader) ([]*vuv1.Activities_BorderCrossingRecord, error) {
	// Parse record array header and records
	return []*vuv1.Activities_BorderCrossingRecord{}, nil
}

func parseVuLoadUnloadRecordArray(r *bytes.Reader) ([]*vuv1.Activities_LoadUnloadRecord, error) {
	// Parse record array header and records
	return []*vuv1.Activities_LoadUnloadRecord{}, nil
}

func parseSignatureRecordArray(r *bytes.Reader) ([]byte, error) {
	// Parse signature record array and return the signature bytes
	return []byte{}, nil
}
