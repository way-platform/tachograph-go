package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
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

	// Store the raw cyclic buffer for round-trip fidelity
	target.SetRawCyclicBuffer(activityData)

	// Parse records using the iterator
	dailyRecords, err := parseActivityRecordsWithIterator(activityData, int(newestDayRecordPointer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cyclic activity daily records: %w", err)
	}
	target.SetDailyRecords(dailyRecords)

	return target, nil
}

// parseActivityRecordsWithIterator parses activity records using the CyclicRecordIterator.
// This separates the complex traversal logic from the parsing logic, improving maintainability
// and enabling the buffer painting strategy for perfect round-trip fidelity.
func parseActivityRecordsWithIterator(buffer []byte, startPos int) ([]*cardv1.DriverActivityData_DailyRecord, error) {
	var records []*cardv1.DriverActivityData_DailyRecord

	iterator := NewCyclicRecordIterator(buffer, startPos)
	for iterator.Next() {
		recordBytes, _, _ := iterator.Record()

		// Try to parse the record semantically
		parsedRecord, err := parseSingleActivityDailyRecord(recordBytes)
		dailyRecord := &cardv1.DriverActivityData_DailyRecord{}

		if err != nil {
			// Parsing failed, store as raw
			dailyRecord.SetValid(false)
			dailyRecord.SetRaw(recordBytes)
		} else {
			// Parsing succeeded, store semantic data
			dailyRecord.SetValid(true)
			dailyRecord.SetActivityPreviousRecordLength(parsedRecord.GetActivityPreviousRecordLength())
			dailyRecord.SetActivityRecordLength(parsedRecord.GetActivityRecordLength())
			dailyRecord.SetActivityRecordDate(parsedRecord.GetActivityRecordDate())
			dailyRecord.SetActivityDailyPresenceCounter(parsedRecord.GetActivityDailyPresenceCounter())
			dailyRecord.SetActivityDayDistance(parsedRecord.GetActivityDayDistance())
			dailyRecord.SetActivityChangeInfo(parsedRecord.GetActivityChangeInfo())
		}

		// Position information is inferred during marshalling by walking the linked list

		records = append(records, dailyRecord)
	}

	if err := iterator.Err(); err != nil {
		return nil, err
	}

	// Reverse to get chronological order (oldest to newest)
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
	date, err := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse activity record date: %w", err)
	}
	record.SetActivityRecordDate(date)
	offset += 4

	// Read activity daily presence counter (2 bytes BCD)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for presence counter")
	}
	bcdCounter, err := dd.UnmarshalBcdString(data[offset : offset+2])
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
		if enumNum, found := dd.GetEnumForProtocolValue(ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED.Descriptor(), slot); found {
			activityChange.SetSlot(ddv1.CardSlotNumber(enumNum))
		}
		if enumNum, found := dd.GetEnumForProtocolValue(ddv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED.Descriptor(), drivingStatus); found {
			activityChange.SetDrivingStatus(ddv1.DrivingStatus(enumNum))
		}
		activityChange.SetInserted(cardStatus != 0) // Convert to boolean (1 = inserted, 0 = not inserted)
		if enumNum, found := dd.GetEnumForProtocolValue(ddv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED.Descriptor(), activity); found {
			activityChange.SetActivity(ddv1.DriverActivityValue(enumNum))
		}

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

	// Use buffer painting strategy for perfect fidelity when raw buffer is available
	if rawBuffer := activity.GetRawCyclicBuffer(); len(rawBuffer) > 0 {
		bufferCopy := make([]byte, len(rawBuffer))
		copy(bufferCopy, rawBuffer)

		// Paint updated records back onto the buffer by walking the linked list
		// to infer their original positions
		err := paintRecordsOntoBuffer(bufferCopy, activity.GetDailyRecords(), int(activity.GetNewestDayRecordIndex()))
		if err != nil {
			return nil, fmt.Errorf("failed to paint records onto buffer: %w", err)
		}

		dst = append(dst, bufferCopy...)
	} else {
		// Fallback to current logic if no raw buffer available (backward compatibility)
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
	}

	return dst, nil
}

// paintRecordsOntoBuffer paints updated records back onto the cyclic buffer
// by walking the linked list structure to infer their original positions.
// This implements the buffer painting strategy without storing explicit positions.
func paintRecordsOntoBuffer(buffer []byte, records []*cardv1.DriverActivityData_DailyRecord, startPos int) error {
	if len(records) == 0 {
		return nil
	}

	// Create a reverse mapping from chronological order to buffer positions
	// by walking the linked list structure like the iterator does
	iterator := NewCyclicRecordIterator(buffer, startPos)
	recordIndex := len(records) - 1 // Start from the newest (last in chronological order)

	for iterator.Next() && recordIndex >= 0 {
		_, recordStart, recordLength := iterator.Record()
		record := records[recordIndex]

		// Only paint valid (parsed) records that have been modified
		if record.GetValid() {
			// Re-serialize the semantic record
			recordBytes, err := marshalSingleActivityRecord(record)
			if err != nil {
				return fmt.Errorf("failed to marshal activity record at index %d: %w", recordIndex, err)
			}

			// Validate that we don't exceed the original buffer bounds
			if len(recordBytes) > recordLength {
				return fmt.Errorf("serialized record length (%d) exceeds original length (%d) at position %d",
					len(recordBytes), recordLength, recordStart)
			}

			// Paint the record back onto the buffer at its original position
			for i, b := range recordBytes {
				if i < recordLength {
					buffer[(recordStart+i)%len(buffer)] = b
				}
			}

			// If the new record is shorter than the original, fill remaining bytes with zeros
			for i := len(recordBytes); i < recordLength; i++ {
				buffer[(recordStart+i)%len(buffer)] = 0
			}
		}
		// Raw records are already in the buffer at their correct positions, no painting needed

		recordIndex--
	}

	if err := iterator.Err(); err != nil {
		return fmt.Errorf("error walking linked list during buffer painting: %w", err)
	}

	if recordIndex >= 0 {
		return fmt.Errorf("mismatch between records (%d) and linked list positions", len(records))
	}

	return nil
}

// marshalSingleActivityRecord marshals a single activity record to bytes.
// This function is used by the buffer painting strategy to re-serialize
// semantic records back to their binary form.
func marshalSingleActivityRecord(record *cardv1.DriverActivityData_DailyRecord) ([]byte, error) {
	if !record.GetValid() {
		// For invalid records, return the raw bytes if available
		if raw := record.GetRaw(); len(raw) > 0 {
			return raw, nil
		}
		return nil, fmt.Errorf("invalid record has no raw data")
	}

	// Start with empty buffer and use existing append logic
	var buf []byte
	var err error
	buf, err = appendParsedDailyRecord(buf, record)
	if err != nil {
		return nil, err
	}

	return buf, nil
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
	contentBuf = dd.AppendDatef(contentBuf, rec.GetActivityRecordDate())

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
	slot := dd.GetCardSlotNumberProtocolValue(ac.GetSlot(), 0)
	drivingStatus := dd.GetDrivingStatusProtocolValue(ac.GetDrivingStatus(), 0)
	cardInserted := dd.GetCardInsertedFromBool(ac.GetInserted())
	activity := dd.GetDriverActivityValueProtocolValue(ac.GetActivity(), 0)

	aci |= (uint16(slot) & 0x1) << 15
	aci |= (uint16(drivingStatus) & 0x1) << 14
	aci |= (uint16(cardInserted) & 0x1) << 13
	aci |= (uint16(activity) & 0x3) << 11
	aci |= (uint16(ac.GetTimeOfChangeMinutes()) & 0x7FF)

	return binary.BigEndian.AppendUint16(dst, aci), nil
}

// cyclicRecordIterator provides a clean interface for traversing the cyclic buffer
// of daily activity records, separating the complex pointer-following logic
// from the parsing of individual records.
//
// The iterator follows the linked list structure where each record contains a
// pointer to the previous record's length, allowing backward traversal through
// the cyclic buffer while handling wrap-around conditions.
type cyclicRecordIterator struct {
	buffer      []byte
	currentPos  int
	recordCount int
	err         error

	// Current record state
	recordStart  int
	recordLength int
	recordBytes  []byte
}

// NewCyclicRecordIterator creates a new iterator for traversing activity records
// in the cyclic buffer, starting from the newest record position.
func NewCyclicRecordIterator(buffer []byte, startPos int) *cyclicRecordIterator {
	return &cyclicRecordIterator{
		buffer:     buffer,
		currentPos: startPos,
	}
}

// Next advances to the next record in the cyclic buffer.
// Returns true if a record was found, false if end of chain or error.
// The iterator traverses backwards from newest to oldest record.
func (it *cyclicRecordIterator) Next() bool {
	const maxRecords = 366 // Safety limit to prevent infinite loops (max days per year + 1)
	if it.err != nil {
		return false
	}
	if it.recordCount >= maxRecords {
		it.err = fmt.Errorf("exceeded maximum record count (%d), possible infinite loop", maxRecords)
		return false
	}
	if len(it.buffer) == 0 {
		return false // No data to parse
	}
	// Validate current position for reading header
	if it.currentPos < 0 || it.currentPos+4 > len(it.buffer) {
		return false // Invalid position for header
	}
	// Read record header (4 bytes: prevRecordLength + currentRecordLength)
	prevRecordLength := int(binary.BigEndian.Uint16(it.buffer[it.currentPos : it.currentPos+2]))
	currentRecordLength := int(binary.BigEndian.Uint16(it.buffer[it.currentPos+2 : it.currentPos+4]))
	if currentRecordLength == 0 {
		return false // Zero-length record signifies end of chain
	}
	// Validate record length
	if currentRecordLength < 4 {
		it.err = fmt.Errorf("invalid record length %d at position %d", currentRecordLength, it.currentPos)
		return false
	}
	// Store current record information
	it.recordStart = it.currentPos
	it.recordLength = currentRecordLength
	// Extract record bytes, handling buffer wrap-around
	it.recordBytes = make([]byte, currentRecordLength)
	for i := 0; i < currentRecordLength; i++ {
		it.recordBytes[i] = it.buffer[(it.currentPos+i)%len(it.buffer)]
	}
	it.recordCount++
	// Move to previous record for next iteration
	if prevRecordLength == 0 {
		// End of chain marker - no more records
		it.currentPos = -1 // Mark as finished
	} else {
		// Move backwards by prevRecordLength, handling wrap-around
		it.currentPos -= prevRecordLength
		if it.currentPos < 0 {
			it.currentPos += len(it.buffer)
		}
	}
	return true
}

// Record returns the bytes of the current record along with its position and length
// in the original buffer. This information is needed for the buffer painting strategy.
func (it *cyclicRecordIterator) Record() (recordBytes []byte, position int, length int) {
	return it.recordBytes, it.recordStart, it.recordLength
}

// Err returns any error encountered during traversal.
func (it *cyclicRecordIterator) Err() error {
	return it.err
}
