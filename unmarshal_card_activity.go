package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardActivityData unmarshals driver activity data from a card EF.
func unmarshalCardActivityData(data []byte) (*cardv1.DriverActivity, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for activity data header")
	}

	var target cardv1.DriverActivity
	r := bytes.NewReader(data)

	// Read pointers (2 bytes each)
	var oldestDayRecordPointer uint16
	var newestDayRecordPointer uint16
	if err := binary.Read(r, binary.BigEndian, &oldestDayRecordPointer); err != nil {
		return nil, fmt.Errorf("failed to read oldest day record pointer: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &newestDayRecordPointer); err != nil {
		return nil, fmt.Errorf("failed to read newest day record pointer: %w", err)
	}

	target.SetOldestDayRecordIndex(int32(oldestDayRecordPointer))
	target.SetNewestDayRecordIndex(int32(newestDayRecordPointer))

	// Read the remaining data as ring buffer
	remainingData := make([]byte, r.Len())
	if _, err := r.Read(remainingData); err != nil {
		return nil, fmt.Errorf("failed to read activity daily records: %w", err)
	}

	// Use tagged union approach for perfect roundtrip
	// For now, preserve raw data until full semantic parsing is implemented
	target.SetValid(false)
	target.SetRawData(remainingData)

	return &target, nil
}

// parseActivityDailyRecords parses the cyclic activity daily records structure
func parseActivityDailyRecords(data []byte, newestRecordPos int) ([]*cardv1.DriverActivity_DailyRecord, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for activity daily records")
	}

	var records []*cardv1.DriverActivity_DailyRecord
	currentPos := newestRecordPos

	// We'll parse up to 365 days worth of records (maximum for full ring buffer)
	for i := 0; i < 365; i++ {
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

		// Parse the daily record (now with 2-byte activity changes)
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
	if len(data) < 12 { // Minimum size: 4 bytes header + 4 bytes date + 2 bytes counter + 2 bytes distance
		return nil, fmt.Errorf("insufficient data for daily record")
	}

	record := &cardv1.DriverActivity_DailyRecord{}

	// Parse header (4 bytes)
	prevRecordLength := binary.BigEndian.Uint16(data[0:2])
	currentRecordLength := binary.BigEndian.Uint16(data[2:4])
	record.SetActivityPreviousRecordLength(int32(prevRecordLength))
	record.SetActivityRecordLength(int32(currentRecordLength))

	// Parse fixed-size content starting at offset 4
	offset := 4

	// Read activity record date (4 bytes BCD)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for activity record date")
	}
	record.SetActivityRecordDate(readDatef(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Read activity daily presence counter (2 bytes BCD)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for presence counter")
	}
	presenceCounter := binary.BigEndian.Uint16(data[offset : offset+2])
	record.SetActivityDailyPresenceCounter(int32(presenceCounter))
	offset += 2

	// Read activity day distance (2 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for day distance")
	}
	dayDistance := binary.BigEndian.Uint16(data[offset : offset+2])
	record.SetActivityDayDistance(int32(dayDistance))
	offset += 2

	// Parse activity change info - loop through remainder in 2-byte chunks
	var activityChanges []*cardv1.DriverActivity_DailyRecord_ActivityChange

	for offset+2 <= len(data) {
		// Parse 2-byte ActivityChangeInfo bitfield
		changeData := binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		// Skip invalid entries (all zeros or all ones)
		if changeData == 0 || changeData == 0xFFFF {
			continue
		}

		// Parse 2-byte bitfield according to spec:
		// Bit 15: Slot (0 = Driver, 1 = Co-driver)
		// Bit 14: Driving Status (0 = Single, 1 = Crew)
		// Bit 13: Card Status (0 = Not inserted, 1 = Inserted)
		// Bits 11-12: Activity (0 = Rest/Break, 1 = Available, 2 = Work, 3 = Driving)
		// Bits 0-10: Time of Change (Minutes since 00:00 UTC)

		slot := int32((changeData >> 15) & 0x1)
		drivingStatus := int32((changeData >> 14) & 0x1)
		cardStatus := int32((changeData >> 13) & 0x1)
		activity := int32((changeData >> 11) & 0x3)
		timeOfChange := int32(changeData & 0x7FF)

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
