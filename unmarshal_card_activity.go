package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardActivityData unmarshals driver activity data from a card EF.
func UnmarshalCardActivityData(data []byte, target *cardv1.DriverActivity) error {
	r := bytes.NewReader(data)

	// Read pointers (2 bytes each)
	var oldestDayRecordPointer uint16
	var newestDayRecordPointer uint16
	if err := binary.Read(r, binary.BigEndian, &oldestDayRecordPointer); err != nil {
		return fmt.Errorf("failed to read oldest day record pointer: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &newestDayRecordPointer); err != nil {
		return fmt.Errorf("failed to read newest day record pointer: %w", err)
	}

	target.SetOldestDayRecordIndex(int32(oldestDayRecordPointer))
	target.SetNewestDayRecordIndex(int32(newestDayRecordPointer))

	// Read the remaining data as activity daily records
	remainingData := make([]byte, r.Len())
	if _, err := r.Read(remainingData); err != nil {
		return fmt.Errorf("failed to read activity daily records: %w", err)
	}

	// Parse daily records in a cyclic manner
	dailyRecords, err := parseActivityDailyRecords(remainingData, int(newestDayRecordPointer)) // Test: pointer might already be relative to remaining data
	if err != nil {
		return fmt.Errorf("failed to parse activity daily records: %w", err)
	}

	target.SetDailyRecords(dailyRecords)
	return nil
}

// parseActivityDailyRecords parses the cyclic activity daily records structure
func parseActivityDailyRecords(data []byte, newestRecordPos int) ([]*cardv1.DriverActivity_DailyRecord, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for activity daily records")
	}

	var records []*cardv1.DriverActivity_DailyRecord
	currentPos := newestRecordPos

	// We'll parse up to 28 days worth of records (maximum)
	for i := 0; i < 28; i++ {
		if currentPos < 0 || currentPos >= len(data) {
			break
		}

		// Read the record header (first 4 bytes)
		headerBytes := make([]byte, 4)
		for j := 0; j < 4; j++ {
			headerBytes[j] = data[(currentPos+j)%len(data)]
		}

		// Parse header
		prevRecordLength := binary.BigEndian.Uint16(headerBytes[0:2])
		currentRecordLength := binary.BigEndian.Uint16(headerBytes[2:4])

		// Invalid records: truly malformed data that should stop parsing
		if currentRecordLength == 0 || currentRecordLength > uint16(len(data)) || currentRecordLength < 4 {
			break // Invalid record length (must be at least 4 bytes for header)
		}

		// Read the full record
		recordData := make([]byte, currentRecordLength)
		for j := 0; j < int(currentRecordLength); j++ {
			recordData[j] = data[(currentPos+j)%len(data)]
		}

		// Parse the daily record
		dailyRecord, err := parseActivityDailyRecord(recordData)
		if err != nil {
			// Create an empty record for roundtrip compatibility instead of stopping
			dailyRecord = &cardv1.DriverActivity_DailyRecord{}
			// Set the record lengths from the header we already parsed
			dailyRecord.SetActivityPreviousRecordLength(int32(prevRecordLength))
			dailyRecord.SetActivityRecordLength(int32(currentRecordLength))
		}

		records = append(records, dailyRecord)

		// Move to previous record
		// Continue parsing even if prevRecordLength is 0 (empty record)
		// Only break if we've parsed all 28 expected records or hit an invalid position
		if prevRecordLength == 0 {
			// For empty records, assume a minimum record size to continue
			prevRecordLength = 4 // Minimum record size (header only)
		}
		currentPos -= int(prevRecordLength)
		if currentPos < 0 {
			currentPos += len(data) // Wrap around
		}
	}

	// Reverse the order since we parsed backwards
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return records, nil
}

// parseActivityDailyRecord parses a single daily record
func parseActivityDailyRecord(data []byte) (*cardv1.DriverActivity_DailyRecord, error) {
	if len(data) < 16 { // Minimum size for a daily record
		return nil, fmt.Errorf("insufficient data for daily record")
	}

	r := bytes.NewReader(data)

	record := &cardv1.DriverActivity_DailyRecord{}

	// Skip the header we already parsed (4 bytes)
	r.Seek(4, 0)

	// Read activity record date (4 bytes BCD)
	record.SetActivityRecordDate(readDatef(r))

	// Read activity daily presence counter (2 bytes)
	var presenceCounter uint16
	if err := binary.Read(r, binary.BigEndian, &presenceCounter); err != nil {
		return nil, fmt.Errorf("failed to read presence counter: %w", err)
	}
	record.SetActivityDailyPresenceCounter(int32(presenceCounter))

	// Read activity day distance (2 bytes)
	var dayDistance uint16
	if err := binary.Read(r, binary.BigEndian, &dayDistance); err != nil {
		return nil, fmt.Errorf("failed to read day distance: %w", err)
	}
	record.SetActivityDayDistance(int32(dayDistance))

	// Parse activity change info
	var activityChanges []*cardv1.DriverActivity_DailyRecord_ActivityChange

	// The rest of the data contains activity changes (4 bytes each)
	for r.Len() >= 4 {
		var changeData uint32
		if err := binary.Read(r, binary.BigEndian, &changeData); err != nil {
			break
		}

		// Parse the 32-bit activity change info
		// Bits 31-30: Slot
		// Bits 29-28: Driving status
		// Bits 27-26: Card status
		// Bits 25-23: Activity
		// Bits 22-11: Time of change (in minutes)
		// Bits 10-0: Reserved

		slot := int32((changeData >> 30) & 0x3)
		drivingStatus := int32((changeData >> 28) & 0x3)
		cardStatus := int32((changeData >> 26) & 0x3)
		activity := int32((changeData >> 23) & 0x7)
		timeOfChange := int32((changeData >> 11) & 0xFFF)

		// Skip invalid entries (all zeros or all ones)
		if changeData == 0 || changeData == 0xFFFFFFFF {
			continue
		}

		activityChange := &cardv1.DriverActivity_DailyRecord_ActivityChange{}

		// Convert raw values to enums using protocol annotations
		SetCardSlotNumber(slot, activityChange.SetSlot, activityChange.SetUnrecognizedSlot)
		SetDrivingStatus(drivingStatus, activityChange.SetDrivingStatus, activityChange.SetUnrecognizedDrivingStatus)
		SetCardStatus(cardStatus, activityChange.SetCardStatus, activityChange.SetUnrecognizedCardStatus)
		SetDriverActivityValue(activity, activityChange.SetActivity, activityChange.SetUnrecognizedActivity)

		activityChange.SetTimeOfChangeMinutes(timeOfChange)

		activityChanges = append(activityChanges, activityChange)
	}

	record.SetActivityChangeInfo(activityChanges)
	return record, nil
}
