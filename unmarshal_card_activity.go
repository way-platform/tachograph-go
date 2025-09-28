package tachograph

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// unmarshalDriverActivityData unmarshals driver activity data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.17):
//
//	CardDriverActivity ::= SEQUENCE {
//	    activityPointerOldestDayRecord    INTEGER(0..CardActivityLengthRange),
//	    activityPointerNewestRecord       INTEGER(0..CardActivityLengthRange),
//	    activityDailyRecords              OCTET STRING (SIZE (CardActivityLengthRange))
//	}
//
//	CardActivityDailyRecord ::= SEQUENCE {
//	    activityPreviousRecordLength      INTEGER(0..CardActivityLengthRange),
//	    activityRecordLength              INTEGER(0..CardActivityLengthRange),
//	    activityRecordDate                TimeReal,
//	    activityDailyPresenceCounter      DailyPresenceCounter,
//	    activityDayDistance               Distance,
//	    activityChangeInfo                SET SIZE (1..1440) OF ActivityChangeInfo
//	}
//
//	ActivityChangeInfo ::= OCTET STRING (SIZE (2))
//
// Binary Layout (variable size):
//
//	0-1:   activityPointerOldestDayRecord (2 bytes, big-endian)
//	2-3:   activityPointerNewestRecord (2 bytes, big-endian)
//	4+:    activityDailyRecords (cyclic buffer of CardActivityDailyRecord)
//	  - 0-1:   activityPreviousRecordLength (2 bytes, big-endian)
//	  - 2-3:   activityRecordLength (2 bytes, big-endian)
//	  - 4-7:   activityRecordDate (4 bytes, TimeReal)
//	  - 8-8:   activityDailyPresenceCounter (1 byte)
//	  - 9-12:  activityDayDistance (4 bytes, big-endian)
//	  - 13+:   activityChangeInfo (2 bytes each, up to 1440 records)
//
// Constants:
const (
	// CardDriverActivity header size
	cardDriverActivityHeaderSize = 4 // 2 bytes oldest + 2 bytes newest pointer

)

func unmarshalDriverActivityData(data []byte) (*cardv1.DriverActivityData, error) {
	if len(data) < cardDriverActivityHeaderSize {
		return nil, fmt.Errorf("insufficient data for activity data header")
	}

	target := &cardv1.DriverActivityData{}
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

	// The rest of the data is the cyclic buffer of daily records.
	activityData := make([]byte, r.Len())
	if _, err := r.Read(activityData); err != nil {
		return nil, fmt.Errorf("failed to read activity daily records: %w", err)
	}

	dailyRecords, err := parseCyclicActivityDailyRecords(activityData, int(newestDayRecordPointer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cyclic activity daily records: %w", err)
	}
	target.SetDailyRecords(dailyRecords)

	return target, nil
}

// parseCyclicActivityDailyRecords parses the cyclic activity daily records structure.
func parseCyclicActivityDailyRecords(data []byte, newestRecordPos int) ([]*cardv1.DriverActivityData_DailyRecord, error) {
	if len(data) == 0 {
		return nil, nil // No data to parse
	}

	var records []*cardv1.DriverActivityData_DailyRecord
	currentPos := newestRecordPos

	// We'll parse up to 366 days worth of records as a safeguard against infinite loops.
	for i := 0; i < 366; i++ {
		// The current position must be valid to read a header.
		if currentPos < 0 || currentPos+4 > len(data) {
			// Invalid starting position for a header, stop.
			break
		}

		// Read the record header
		prevRecordLength := int(binary.BigEndian.Uint16(data[currentPos : currentPos+2]))
		currentRecordLength := int(binary.BigEndian.Uint16(data[currentPos+2 : currentPos+4]))

		if currentRecordLength == 0 {
			// A zero-length record signifies the end of the chain.
			break
		}
		// Read the full record data, handling buffer wrap-around.
		recordData := make([]byte, currentRecordLength)
		for j := 0; j < currentRecordLength; j++ {
			recordData[j] = data[(currentPos+j)%len(data)]
		}

		// Attempt to parse the record semantically.
		parsedRecord, err := parseSingleActivityDailyRecord(recordData)
		dailyRecord := &cardv1.DriverActivityData_DailyRecord{}

		if err != nil {
			// Parsing failed, store as raw.
			dailyRecord.SetValid(false)
			dailyRecord.SetRaw(recordData)
		} else {
			// Parsing succeeded.
			dailyRecord.SetValid(true)
			dailyRecord.SetActivityPreviousRecordLength(parsedRecord.GetActivityPreviousRecordLength())
			dailyRecord.SetActivityRecordLength(parsedRecord.GetActivityRecordLength())
			dailyRecord.SetActivityRecordDate(parsedRecord.GetActivityRecordDate())
			dailyRecord.SetActivityDailyPresenceCounter(parsedRecord.GetActivityDailyPresenceCounter())
			dailyRecord.SetActivityDayDistance(parsedRecord.GetActivityDayDistance())
			dailyRecord.SetActivityChangeInfo(parsedRecord.GetActivityChangeInfo())
		}
		records = append(records, dailyRecord)

		if prevRecordLength == 0 {
			// End of the chain.
			break
		}

		// Move to the previous record, handling wrap-around.
		currentPos -= prevRecordLength
		if currentPos < 0 {
			currentPos += len(data)
		}
	}

	// Reverse the order since we parsed backwards from newest to oldest.
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return records, nil
}

// parseSingleActivityDailyRecord parses a single daily record byte slice.
func parseSingleActivityDailyRecord(data []byte) (*cardv1.DriverActivityData_DailyRecord, error) {
	if len(data) < 12 { // Minimum size: 4 bytes header + 4 bytes date + 2 bytes counter + 2 bytes distance
		return nil, fmt.Errorf("insufficient data for daily record, got %d bytes", len(data))
	}

	record := &cardv1.DriverActivityData_DailyRecord{}

	// Parse header (4 bytes)
	prevRecordLength := binary.BigEndian.Uint16(data[0:2])
	currentRecordLength := binary.BigEndian.Uint16(data[2:4])
	record.SetActivityPreviousRecordLength(int32(prevRecordLength))
	record.SetActivityRecordLength(int32(currentRecordLength))

	// Parse fixed-size content starting at offset 4
	offset := 4

	// Read activity record date (4 bytes TimeReal)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for activity record date")
	}
	date := readTimeReal(bytes.NewReader(data[offset : offset+4]))
	record.SetActivityRecordDate(date)
	offset += 4

	// Read activity daily presence counter (2 bytes BCD)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for presence counter")
	}
	bcdCounter, err := createBcdString(data[offset : offset+2])
	if err != nil {
		return nil, fmt.Errorf("failed to create BCD string for presence counter: %w", err)
	}
	record.SetActivityDailyPresenceCounter(bcdCounter)
	offset += 2

	// Read activity day distance (2 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for day distance")
	}
	dayDistance := binary.BigEndian.Uint16(data[offset : offset+2])
	record.SetActivityDayDistance(int32(dayDistance))
	offset += 2

	// Parse activity change info - loop through remainder in 2-byte chunks
	var activityChanges []*ddv1.ActivityChangeInfo

	for offset+2 <= len(data) {
		// Parse 2-byte ActivityChangeInfo bitfield
		changeData := binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		// Skip invalid entries (all zeros or all ones)
		if changeData == 0 || changeData == 0xFFFF {
			continue
		}

		// Parse 2-byte bitfield according to spec
		slot := int32((changeData >> 15) & 0x1)
		drivingStatus := int32((changeData >> 14) & 0x1)
		cardStatus := int32((changeData >> 13) & 0x1)
		activity := int32((changeData >> 11) & 0x3)
		timeOfChange := int32(changeData & 0x7FF)

		activityChange := &ddv1.ActivityChangeInfo{}

		// Convert raw values to enums using protocol annotations
		SetCardSlotNumber(ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED.Descriptor(), slot, func(en protoreflect.EnumNumber) {
			activityChange.SetSlot(ddv1.CardSlotNumber(en))
		}, nil)
		SetDrivingStatus(ddv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED.Descriptor(), drivingStatus, func(en protoreflect.EnumNumber) {
			activityChange.SetDrivingStatus(ddv1.DrivingStatus(en))
		}, nil)
		activityChange.SetInserted(cardStatus != 0) // Convert to boolean (1 = inserted, 0 = not inserted)
		SetDriverActivityValue(ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED.Descriptor(), activity, func(en protoreflect.EnumNumber) {
			activityChange.SetActivity(ddv1.DriverActivityValue(en))
		}, nil)

		activityChange.SetTimeOfChangeMinutes(timeOfChange)

		activityChanges = append(activityChanges, activityChange)
	}

	record.SetActivityChangeInfo(activityChanges)
	return record, nil
}
