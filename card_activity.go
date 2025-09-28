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
// The data type `CardDriverActivity` is specified in the Data Dictionary, Section 2.17.
//
// ASN.1 Definition:
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
func unmarshalDriverActivityData(data []byte) (*cardv1.DriverActivityData, error) {
	const (
		lenCardDriverActivityHeader = 4 // 2 bytes oldest + 2 bytes newest pointer
	)

	if len(data) < lenCardDriverActivityHeader {
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
	const (
		lenMinDailyRecord = 12 // Minimum size: 4 bytes header + 4 bytes date + 2 bytes counter + 2 bytes distance
	)

	if len(data) < lenMinDailyRecord {
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
		setEnumFromProtocolValue(ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED.Descriptor(), slot, func(en protoreflect.EnumNumber) {
			activityChange.SetSlot(ddv1.CardSlotNumber(en))
		}, nil)
		setEnumFromProtocolValue(ddv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED.Descriptor(), drivingStatus, func(en protoreflect.EnumNumber) {
			activityChange.SetDrivingStatus(ddv1.DrivingStatus(en))
		}, nil)
		activityChange.SetInserted(cardStatus != 0) // Convert to boolean (1 = inserted, 0 = not inserted)
		setEnumFromProtocolValue(ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED.Descriptor(), activity, func(en protoreflect.EnumNumber) {
			activityChange.SetActivity(ddv1.DriverActivityValue(en))
		}, nil)

		activityChange.SetTimeOfChangeMinutes(timeOfChange)

		activityChanges = append(activityChanges, activityChange)
	}

	record.SetActivityChangeInfo(activityChanges)
	return record, nil
}

// AppendDriverActivity appends the binary representation of DriverActivityData to dst.
//
// The data type `CardDriverActivity` is specified in the Data Dictionary, Section 2.17.
//
// ASN.1 Definition:
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
func appendDriverActivity(dst []byte, activity *cardv1.DriverActivityData) ([]byte, error) {
	if activity == nil {
		return dst, nil
	}

	// Append header (pointers)
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetOldestDayRecordIndex()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetNewestDayRecordIndex()))

	// Append the records
	var err error
	for _, rec := range activity.GetDailyRecords() {
		if rec.GetValid() {
			// It's a parsed record, so we need to serialize it.
			dst, err = appendParsedDailyRecord(dst, rec)
			if err != nil {
				return nil, err
			}
		} else {
			// It's a raw record, just append the bytes.
			dst = append(dst, rec.GetRaw()...)
		}
	}
	return dst, nil
}

// appendParsedDailyRecord appends a single parsed daily activity record.
func appendParsedDailyRecord(dst []byte, rec *cardv1.DriverActivityData_DailyRecord) ([]byte, error) {
	// Check if this is an empty record (for roundtrip compatibility)
	isEmpty := rec.GetActivityDayDistance() == 0 &&
		len(rec.GetActivityChangeInfo()) == 0 &&
		(rec.GetActivityRecordDate() == nil || rec.GetActivityRecordDate().GetSeconds() == 0)

	if isEmpty && rec.GetActivityRecordLength() > 0 {
		// For empty records, use the original record length and write zeros
		originalLength := int(rec.GetActivityRecordLength())
		dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetActivityPreviousRecordLength()))
		dst = binary.BigEndian.AppendUint16(dst, uint16(originalLength))

		// Write the content as zeros (originalLength - 4 for the header we already wrote)
		contentLength := originalLength - 4
		if contentLength > 0 {
			dst = append(dst, make([]byte, contentLength)...)
		}
		return dst, nil
	}

	// Normal record processing - serialize content to temporary buffer first
	contentBuf := make([]byte, 0, 2048)

	// Activity record date (4 bytes BCD)
	contentBuf = appendDatef(contentBuf, rec.GetActivityRecordDate())

	// Activity daily presence counter (2 bytes BCD)
	if bcdCounter := rec.GetActivityDailyPresenceCounter(); bcdCounter != nil {
		contentBuf = append(contentBuf, bcdCounter.GetEncoded()...)
	}

	// Activity day distance (2 bytes)
	contentBuf = binary.BigEndian.AppendUint16(contentBuf, uint16(rec.GetActivityDayDistance()))

	// Activity change info (2 bytes each)
	for _, ac := range rec.GetActivityChangeInfo() {
		var err error
		contentBuf, err = appendActivityChange(contentBuf, ac)
		if err != nil {
			return nil, err
		}
	}

	// Calculate total record length (content + 4-byte header)
	recordLength := len(contentBuf) + 4

	// Append header with lengths
	dst = binary.BigEndian.AppendUint16(dst, uint16(rec.GetActivityPreviousRecordLength()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(recordLength))

	// Append content
	dst = append(dst, contentBuf...)

	return dst, nil
}

// AppendActivityChange appends the binary representation of ActivityChangeInfo to dst.
//
// The data type `ActivityChangeInfo` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	ActivityChangeInfo ::= OCTET STRING (SIZE (2))
//
// Binary Layout (2 bytes):
//
//	Bit layout: 'scpaattttttttttt'B (16 bits)
//	- s: Slot (1 bit): '0'B: DRIVER, '1'B: CO-DRIVER
//	- c: Driving status (1 bit): '0'B: SINGLE, '1'B: CREW
//	- p: Card status (1 bit): '0'B: INSERTED, '1'B: NOT INSERTED
//	- aa: Activity (2 bits): '00'B: BREAK/REST, '01'B: AVAILABILITY, '10'B: WORK, '11'B: DRIVING
//	- ttttttttttt: Time of change (11 bits): Number of minutes since 00h00 on the given day
func appendActivityChange(dst []byte, ac *ddv1.ActivityChangeInfo) ([]byte, error) {
	var aci uint16

	// Reconstruct the bitfield from enum values
	slot := getCardSlotNumberProtocolValue(ac.GetSlot(), 0)
	drivingStatus := getDrivingStatusProtocolValue(ac.GetDrivingStatus(), 0)
	cardInserted := getCardInsertedFromBool(ac.GetInserted())
	activity := getDriverActivityValueProtocolValue(ac.GetActivity(), 0)

	aci |= (uint16(slot) & 0x1) << 15
	aci |= (uint16(drivingStatus) & 0x1) << 14
	aci |= (uint16(cardInserted) & 0x1) << 13
	aci |= (uint16(activity) & 0x3) << 11
	aci |= (uint16(ac.GetTimeOfChangeMinutes()) & 0x7FF)

	return binary.BigEndian.AppendUint16(dst, aci), nil
}
