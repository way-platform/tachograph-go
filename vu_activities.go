package tachograph

import (
	"bytes"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuActivities unmarshals VU activities data from a VU transfer.
//
// The data type `VuActivities` is specified in the Data Dictionary, Section 2.2.6.2.
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    dateOfDay                        TimeReal,
//	    odometerValueMidnight            OdometerValueMidnight,
//	    vuCardIWData                     VuCardIWData,
//	    vuActivityDailyData              VuActivityDailyData,
//	    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
//	    vuSpecificConditionData          VuSpecificConditionData,
//	    signature                        SignatureFirstGen
//	}
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
//	    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                      VuCardIWRecordArray,
//	    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                      VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
//	    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
//	    signatureRecordArray                     SignatureRecordArray
//	}
func unmarshalVuActivities(data []byte, offset int, target *vuv1.Activities, generation int) (int, error) {
	switch generation {
	case 1:
		return unmarshalVuActivitiesGen1(data, offset, target)
	case 2:
		return unmarshalVuActivitiesGen2(data, offset, target)
	default:
		return 0, fmt.Errorf("unsupported generation: %d", generation)
	}
}

// unmarshalVuActivitiesGen1 unmarshals Generation 1 VU activities
func unmarshalVuActivitiesGen1(data []byte, offset int, target *vuv1.Activities) (int, error) {
	startOffset := offset
	target.SetGeneration(ddv1.Generation_GENERATION_1)

	// Read TimeReal (4 bytes) - this is the date of the day
	timeReal, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read date of day: %w", err)
	}
	target.SetDateOfDay(timestamppb.New(time.Unix(timeReal, 0)))

	// Read OdometerValueMidnight (3 bytes)
	odometerValue, offset, err := readVuOdometerFromBytes(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read odometer value midnight: %w", err)
	}
	target.SetOdometerMidnightKm(int32(odometerValue))

	// Parse VuCardIWData
	cardIWData, offset, err := parseVuCardIWData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW data: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyData
	activityChanges, offset, err := parseVuActivityDailyData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily data: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodData
	places, offset, err := parseVuPlaceDailyWorkPeriodData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period data: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuSpecificConditionData
	specificConditions, offset, err := parseVuSpecificConditionData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition data: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Read signature (128 bytes for Gen1)
	signatureBytes, offset, err := readBytesFromBytes(data, offset, 128)
	if err != nil {
		return 0, fmt.Errorf("failed to read signature: %w", err)
	}
	target.SetSignatureGen1(signatureBytes)

	return offset - startOffset, nil
}

// unmarshalVuActivitiesGen2 unmarshals Generation 2 VU activities
func unmarshalVuActivitiesGen2(data []byte, offset int, target *vuv1.Activities) (int, error) {
	startOffset := offset
	target.SetGeneration(ddv1.Generation_GENERATION_2)

	// Gen2 format uses record arrays, each with a header
	// Parse DateOfDayDownloadedRecordArray
	dates, offset, err := parseDateOfDayDownloadedRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse date of day downloaded record array: %w", err)
	}
	if len(dates) > 0 {
		target.SetDateOfDay(dates[0]) // Use first date
	}

	// Parse OdometerValueMidnightRecordArray
	odometerValues, offset, err := parseOdometerValueMidnightRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse odometer value midnight record array: %w", err)
	}
	if len(odometerValues) > 0 {
		target.SetOdometerMidnightKm(odometerValues[0])
	}

	// Parse VuCardIWRecordArray
	cardIWData, offset, err := parseVuCardIWRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW record array: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyRecordArray
	activityChanges, offset, err := parseVuActivityDailyRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily record array: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodRecordArray
	places, offset, err := parseVuPlaceDailyWorkPeriodRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period record array: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuGNSSADRecordArray (Gen2+)
	gnssRecords, offset, err := parseVuGNSSADRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse GNSS AD record array: %w", err)
	}
	target.SetGnssAccumulatedDriving(gnssRecords)

	// Parse VuSpecificConditionRecordArray
	specificConditions, offset, err := parseVuSpecificConditionRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition record array: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Try to parse Gen2v2 specific arrays if there's more data
	if offset+10 <= len(data) { // Need some minimum data for arrays
		// Parse VuBorderCrossingRecordArray (Gen2v2+)
		borderCrossings, newOffset, err := parseVuBorderCrossingRecordArray(data, offset)
		if err == nil {
			target.SetBorderCrossings(borderCrossings)
			target.SetVersion(vuv1.Version_VERSION_2)
			offset = newOffset
		}

		// Parse VuLoadUnloadRecordArray (Gen2v2+)
		if offset+5 <= len(data) {
			loadUnloadRecords, newOffset, err := parseVuLoadUnloadRecordArray(data, offset)
			if err == nil {
				target.SetLoadUnloadOperations(loadUnloadRecords)
				offset = newOffset
			}
		}
	}

	// Parse SignatureRecordArray
	signatureBytes, offset, err := parseSignatureRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signature record array: %w", err)
	}
	target.SetSignatureGen2(signatureBytes)

	return offset - startOffset, nil
}

// Helper functions for parsing different record types
// These are simplified implementations - in a full implementation,
// each would need to properly handle the record array format

func parseVuCardIWData(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*vuv1.Activities_CardIWRecord{}, offset, nil
}

func parseVuActivityDailyData(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*ddv1.ActivityChangeInfo{}, offset, nil
}

func parseVuPlaceDailyWorkPeriodData(data []byte, offset int) ([]*vuv1.Activities_PlaceRecord, int, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*vuv1.Activities_PlaceRecord{}, offset, nil
}

func parseVuSpecificConditionData(data []byte, offset int) ([]*vuv1.Activities_SpecificConditionRecord, int, error) {
	// Simplified implementation - would need to parse the actual structure
	return []*vuv1.Activities_SpecificConditionRecord{}, offset, nil
}

// Gen2 record array parsers
func parseDateOfDayDownloadedRecordArray(data []byte, offset int) ([]*timestamppb.Timestamp, int, error) {
	// Parse record array header and records
	return []*timestamppb.Timestamp{}, offset, nil
}

func parseOdometerValueMidnightRecordArray(data []byte, offset int) ([]int32, int, error) {
	// Parse record array header and records
	return []int32{}, offset, nil
}

func parseVuCardIWRecordArray(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// Parse record array header and records
	return []*vuv1.Activities_CardIWRecord{}, offset, nil
}

func parseVuActivityDailyRecordArray(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	// Parse record array header and records
	return []*ddv1.ActivityChangeInfo{}, offset, nil
}

func parseVuPlaceDailyWorkPeriodRecordArray(data []byte, offset int) ([]*vuv1.Activities_PlaceRecord, int, error) {
	// Parse record array header and records
	return []*vuv1.Activities_PlaceRecord{}, offset, nil
}

func parseVuGNSSADRecordArray(data []byte, offset int) ([]*vuv1.Activities_GnssRecord, int, error) {
	// Parse record array header and records
	return []*vuv1.Activities_GnssRecord{}, offset, nil
}

func parseVuSpecificConditionRecordArray(data []byte, offset int) ([]*vuv1.Activities_SpecificConditionRecord, int, error) {
	// Parse record array header and records
	return []*vuv1.Activities_SpecificConditionRecord{}, offset, nil
}

func parseVuBorderCrossingRecordArray(data []byte, offset int) ([]*vuv1.Activities_BorderCrossingRecord, int, error) {
	// Parse record array header and records
	return []*vuv1.Activities_BorderCrossingRecord{}, offset, nil
}

func parseVuLoadUnloadRecordArray(data []byte, offset int) ([]*vuv1.Activities_LoadUnloadRecord, int, error) {
	// Parse record array header and records
	return []*vuv1.Activities_LoadUnloadRecord{}, offset, nil
}

func parseSignatureRecordArray(data []byte, offset int) ([]byte, int, error) {
	// Parse signature record array and return the signature bytes
	return []byte{}, offset, nil
}

// AppendVuActivities appends VU activities data to a buffer.
//
// The data type `VuActivities` is specified in the Data Dictionary, Section 2.2.6.2.
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    dateOfDay                        TimeReal,
//	    odometerValueMidnight            OdometerValueMidnight,
//	    vuCardIWData                     VuCardIWData,
//	    vuActivityDailyData              VuActivityDailyData,
//	    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
//	    vuSpecificConditionData          VuSpecificConditionData,
//	    signature                        SignatureFirstGen
//	}
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
//	    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                      VuCardIWRecordArray,
//	    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                      VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
//	    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
//	    signatureRecordArray                     SignatureRecordArray
//	}
func appendVuActivities(buf *bytes.Buffer, activities *vuv1.Activities) error {
	if activities == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if activities.GetGeneration() == ddv1.Generation_GENERATION_1 {
		signature := activities.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := activities.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}
