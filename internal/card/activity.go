package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (opts UnmarshalOptions) unmarshalDriverActivityData(data []byte) (*cardv1.DriverActivityData, error) {
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
	target.SetRawData(activityData)

	// Parse records using the iterator
	dailyRecords, err := opts.parseActivityRecordsWithIterator(activityData, int(newestDayRecordPointer), int(oldestDayRecordPointer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cyclic activity daily records: %w", err)
	}
	target.SetDailyRecords(dailyRecords)

	return target, nil
}

// parseActivityRecordsWithIterator parses activity records using the CyclicRecordIterator.
// This separates the complex traversal logic from the parsing logic, improving maintainability
// and enabling the buffer painting strategy for perfect round-trip fidelity.
func (opts UnmarshalOptions) parseActivityRecordsWithIterator(buffer []byte, startPos int, oldestPos int) ([]*cardv1.DriverActivityData_DailyRecord, error) {
	var records []*cardv1.DriverActivityData_DailyRecord

	iterator := NewCyclicRecordIterator(buffer, startPos, oldestPos)
	for iterator.Next() {
		recordBytes, _, _ := iterator.Record()

		// Try to parse the record semantically
		parsedRecord, err := opts.parseSingleActivityDailyRecord(recordBytes)
		dailyRecord := &cardv1.DriverActivityData_DailyRecord{}

		if err != nil {
			// Parsing failed, store as raw
			dailyRecord.SetValid(false)
			dailyRecord.SetRawData(recordBytes)
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
func (opts UnmarshalOptions) parseSingleActivityDailyRecord(data []byte) (*cardv1.DriverActivityData_DailyRecord, error) {
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
	date, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse activity record date: %w", err)
	}
	record.SetActivityRecordDate(date)
	offset += 4

	// Read activity daily presence counter (2 bytes BCD)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for presence counter")
	}
	bcdCounter, err := opts.UnmarshalBcdString(data[offset : offset+2])
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
		// Check for invalid entries before parsing (all zeros or all ones)
		changeData := binary.BigEndian.Uint16(data[offset : offset+2])
		if changeData == 0 || changeData == 0xFFFF {
			offset += 2
			continue
		}

		// Parse ActivityChangeInfo using centralized helper
		activityChange, err := opts.UnmarshalActivityChangeInfo(data[offset : offset+2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse activity change info at offset %d: %w", offset, err)
		}
		offset += 2

		activityChanges = append(activityChanges, activityChange)
	}

	record.SetActivityChangeInfo(activityChanges)

	// Store raw_data for round-trip fidelity (enables buffer painting strategy)
	record.SetRawData(data)
	record.SetValid(true)

	return record, nil
}

// MarshalDriverActivity marshals the binary representation of DriverActivityData.
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
func (opts MarshalOptions) MarshalDriverActivity(activity *cardv1.DriverActivityData) ([]byte, error) {
	if activity == nil {
		return nil, nil
	}

	var dst []byte

	// Append header (pointers)
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetOldestDayRecordIndex()))
	dst = binary.BigEndian.AppendUint16(dst, uint16(activity.GetNewestDayRecordIndex()))

	// For perfect round-trip fidelity, use raw buffer directly when available.
	// This preserves all padding, reserved bits, and linked-list structure.
	if rawBuffer := activity.GetRawData(); len(rawBuffer) > 0 {
		dst = append(dst, rawBuffer...)
	} else {
		// Fallback: Build cyclic buffer from scratch with proper linked-list structure
		buffer, err := buildCyclicBufferFromRecords(activity.GetDailyRecords(), int(activity.GetNewestDayRecordIndex()))
		if err != nil {
			return nil, fmt.Errorf("failed to build cyclic buffer: %w", err)
		}
		dst = append(dst, buffer...)
	}

	return dst, nil
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
	oldestPos   int
	recordCount int
	err         error
	seen        map[int]struct{}

	// Current record state
	recordStart  int
	recordLength int
	recordBytes  []byte
}

// NewCyclicRecordIterator creates a new iterator for traversing activity records
// in the cyclic buffer, starting from the newest record position.
// oldestPos is used as the termination condition for wrapped cyclic buffers.
func NewCyclicRecordIterator(buffer []byte, startPos int, oldestPos int) *cyclicRecordIterator {
	return &cyclicRecordIterator{
		buffer:     buffer,
		currentPos: startPos,
		oldestPos:  oldestPos,
		seen:       make(map[int]struct{}),
	}
}

// Next advances to the next record in the cyclic buffer.
// Returns true if a record was found, false if end of chain or error.
// The iterator traverses backwards from newest to oldest record.
func (it *cyclicRecordIterator) Next() bool {
	// Safety backstop: gen2 buffer is 55140 bytes; minimum 12-byte records = ~4595 max.
	const maxRecords = 5000
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
	// Cycle detection: if we've already visited this position, stop.
	if _, ok := it.seen[it.currentPos]; ok {
		return false
	}
	it.seen[it.currentPos] = struct{}{}
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
	// Move to previous record for next iteration.
	// Stop after processing the oldest record: if the current record is the oldest
	// (at oldestPos), do not follow its prevRecordLength further — the buffer may be
	// fully wrapped and the pointer would lead into overwritten/garbage data.
	if it.recordStart == it.oldestPos || prevRecordLength == 0 {
		it.currentPos = -1 // Mark as finished after this record
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

// buildCyclicBufferFromRecords constructs a cyclic buffer from scratch with proper
// linked-list structure. This is used when raw_data is not available (e.g., after anonymization).
//
// LIMITATION: This function does not perfectly reconstruct the original cyclic buffer because:
// - We don't know the original buffer's total size (only the records we parsed)
// - We don't know the original absolute positions of records (only relative prev/current lengths)
// - We create a sequential buffer sized to fit all records contiguously
//
// This means the reconstructed buffer may differ from the original in:
// - Total buffer size
// - Record positions (we place sequentially, original may have gaps/wrapping)
// - The order records appear when re-parsed (cyclic iterator may traverse differently)
//
// For perfect fidelity, callers should preserve and use the original raw_data buffer directly.
// This fallback is primarily for testing scenarios where we need to marshal modified records.
//
// TODO: To fix this limitation:
// - Store original buffer size during parsing
// - Store absolute positions of records (not just prev/next lengths)
// - Allocate buffer of original size and place records at original positions
//
// The cyclic buffer structure:
// - Records are stored sequentially in chronological order (oldest to newest)
// - Each record has a header: [prevRecordLength: 2 bytes][currentRecordLength: 2 bytes]
// - The prevRecordLength points backward to enable traversal from newest to oldest
// - The newest record is at position newestRecordPos
//
// The buffer is sized to accommodate all records sequentially starting from position 0.
func buildCyclicBufferFromRecords(records []*cardv1.DriverActivityData_DailyRecord, newestRecordPos int) ([]byte, error) {
	if len(records) == 0 {
		return nil, nil
	}

	// First pass: calculate the size of each record
	recordSizes := make([]int, len(records))
	totalRecordsSize := 0
	for i, rec := range records {
		if !rec.GetValid() {
			// For invalid records, use raw_data length if available
			if raw := rec.GetRawData(); len(raw) > 0 {
				recordSizes[i] = len(raw)
			} else {
				return nil, fmt.Errorf("invalid record %d has no raw data", i)
			}
		} else {
			// Calculate size for valid record
			size, err := calculateRecordSize(rec)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate size for record %d: %w", i, err)
			}
			recordSizes[i] = size
		}
		totalRecordsSize += recordSizes[i]
	}

	// Calculate buffer size: must be large enough for all records starting at position 0
	// In a real cyclic buffer, we'd place records at their original positions, but since
	// we don't know the original buffer size, we create a buffer that fits all records sequentially.
	bufferSize := totalRecordsSize

	// Allocate buffer (zero-filled by default)
	buffer := make([]byte, bufferSize)

	// Second pass: write records to buffer with proper linked-list pointers
	// Records are written in chronological order (oldest to newest)
	currentPos := 0
	for i, rec := range records {
		recordSize := recordSizes[i]

		// Calculate prevRecordLength (0 for first/oldest record, previous record's size for others)
		prevRecordLength := 0
		if i > 0 {
			prevRecordLength = recordSizes[i-1]
		}

		// Write the record
		if !rec.GetValid() {
			// For invalid records, copy raw_data as-is (it already has the correct header)
			copy(buffer[currentPos:], rec.GetRawData())
		} else {
			// For valid records, marshal with proper header
			recordWithHeader, err := marshalRecordWithHeader(MarshalOptions{}, rec, prevRecordLength, recordSize)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal record %d: %w", i, err)
			}
			copy(buffer[currentPos:], recordWithHeader)
		}

		currentPos += recordSize
	}

	return buffer, nil
}

// calculateRecordSize calculates the size of a marshalled activity record.
// For records with activity_record_length set, we use that to preserve padding.
// Otherwise, we calculate from content.
func calculateRecordSize(rec *cardv1.DriverActivityData_DailyRecord) (int, error) {
	// Use original record length if available (preserves padding)
	if recordLength := rec.GetActivityRecordLength(); recordLength > 0 {
		return int(recordLength), nil
	}

	// Fallback: calculate from content
	const (
		lenHeader               = 4 // prevRecordLength (2) + currentRecordLength (2)
		lenTimeReal             = 4 // activity record date
		lenDailyPresenceCounter = 2 // BCD counter
		lenDayDistance          = 2 // distance
		lenActivityChangeInfo   = 2 // each activity change
	)

	size := lenHeader + lenTimeReal + lenDailyPresenceCounter + lenDayDistance
	size += len(rec.GetActivityChangeInfo()) * lenActivityChangeInfo

	return size, nil
}

// marshalRecordWithHeader marshals a single activity record with the correct header values.
// This ensures the linked-list structure is properly maintained.
func marshalRecordWithHeader(opts MarshalOptions, rec *cardv1.DriverActivityData_DailyRecord, prevRecordLength, currentRecordLength int) ([]byte, error) {
	var buf []byte

	// Write header
	buf = binary.BigEndian.AppendUint16(buf, uint16(prevRecordLength))
	buf = binary.BigEndian.AppendUint16(buf, uint16(currentRecordLength))

	// Write fixed content

	dateBytes, err := opts.MarshalTimeReal(rec.GetActivityRecordDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activity record date: %w", err)
	}
	buf = append(buf, dateBytes...)

	counterBytes, err := opts.MarshalBcdString(rec.GetActivityDailyPresenceCounter())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activity daily presence counter: %w", err)
	}
	buf = append(buf, counterBytes...)

	buf = binary.BigEndian.AppendUint16(buf, uint16(rec.GetActivityDayDistance()))

	// Write activity change info
	for _, ac := range rec.GetActivityChangeInfo() {
		acBytes, err := opts.MarshalActivityChangeInfo(ac)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal activity change info: %w", err)
		}
		buf = append(buf, acBytes...)
	}

	// Add padding if the current length is less than the expected record length
	// This preserves any padding bytes that were in the original record
	if len(buf) < currentRecordLength {
		padding := make([]byte, currentRecordLength-len(buf))
		buf = append(buf, padding...)
	}

	return buf, nil
}

// anonymizeDriverActivityData creates an anonymized copy of DriverActivityData,
// replacing sensitive information with static, deterministic test values.
func (opts AnonymizeOptions) anonymizeDriverActivityData(activity *cardv1.DriverActivityData) *cardv1.DriverActivityData {
	if activity == nil {
		return nil
	}

	anonymized := &cardv1.DriverActivityData{}

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Note: We do NOT preserve raw_data or the cyclic buffer pointers here, as we're modifying
	// the semantic fields (dates and activity change times), which means the buffer must be
	// rebuilt from scratch. The marshaller will recalculate appropriate indices.
	// For simplicity, we'll use indices that correspond to a sequential buffer layout.

	// Base timestamp for anonymization: 2020-01-01 00:00:00 UTC
	baseEpoch := int64(1577836800)

	// Anonymize the parsed record dates and activity changes
	var anonymizedRecords []*cardv1.DriverActivityData_DailyRecord
	for i, record := range activity.GetDailyRecords() {
		anonymizedRecord := &cardv1.DriverActivityData_DailyRecord{}

		// For invalid records, preserve them as-is with their raw_data
		if !record.GetValid() {
			anonymizedRecord.SetValid(false)
			anonymizedRecord.SetRawData(record.GetRawData())
			anonymizedRecords = append(anonymizedRecords, anonymizedRecord)
			continue
		}

		// For valid records, anonymize semantic fields
		anonymizedRecord.SetValid(true)
		// Note: raw_data is NOT preserved for valid records since we're anonymizing
		// the activity change info, which would make raw_data inconsistent with
		// the semantic fields. The marshaller will regenerate the binary representation.

		// Anonymize the date field
		recordDate := &timestamppb.Timestamp{Seconds: baseEpoch + int64(i)*86400}
		anonymizedRecord.SetActivityRecordDate(recordDate)

		// Preserve data fields including record lengths (needed for consistent buffer layout)
		anonymizedRecord.SetActivityPreviousRecordLength(record.GetActivityPreviousRecordLength())
		anonymizedRecord.SetActivityRecordLength(record.GetActivityRecordLength())
		anonymizedRecord.SetActivityDailyPresenceCounter(record.GetActivityDailyPresenceCounter())
		anonymizedRecord.SetActivityDayDistance(record.GetActivityDayDistance())

		// Anonymize activity change info (time intervals)
		if changes := record.GetActivityChangeInfo(); changes != nil {
			var anonymizedChanges []*ddv1.ActivityChangeInfo
			for j, change := range changes {
				anonymizedChange := ddOpts.AnonymizeActivityChangeInfo(change, j)
				anonymizedChanges = append(anonymizedChanges, anonymizedChange)
			}
			anonymizedRecord.SetActivityChangeInfo(anonymizedChanges)
		}

		anonymizedRecords = append(anonymizedRecords, anonymizedRecord)
	}

	anonymized.SetDailyRecords(anonymizedRecords)

	// Calculate buffer indices for sequential layout
	// Oldest record is at position 0, newest is at sum of all record sizes except the last
	if len(anonymizedRecords) > 0 {
		anonymized.SetOldestDayRecordIndex(0)

		// Calculate position of newest (last) record
		newestPos := 0
		for i := 0; i < len(anonymizedRecords)-1; i++ {
			recordSize := opts.calculateAnonymizedRecordSize(anonymizedRecords[i])
			newestPos += recordSize
		}
		anonymized.SetNewestDayRecordIndex(int32(newestPos))
	}

	// Signature and raw_data fields left unset (nil) - TLV marshaller will omit these blocks

	return anonymized
}

// calculateAnonymizedRecordSize calculates the size of an anonymized record.
// For valid records, we use the original activity_record_length to preserve
// any padding bytes that were in the original record.
func (opts AnonymizeOptions) calculateAnonymizedRecordSize(rec *cardv1.DriverActivityData_DailyRecord) int {
	if !rec.GetValid() {
		return len(rec.GetRawData())
	}

	// For valid records, use the original record length if available
	// This preserves padding and ensures consistent buffer layout
	if recordLength := rec.GetActivityRecordLength(); recordLength > 0 {
		return int(recordLength)
	}

	// Fallback: calculate from content
	// For valid records: 4 byte header + 4 byte date + 2 byte counter + 2 byte distance + (2 bytes * num changes)
	const (
		lenHeader               = 4
		lenTimeReal             = 4
		lenDailyPresenceCounter = 2
		lenDayDistance          = 2
		lenActivityChangeInfo   = 2
	)
	return lenHeader + lenTimeReal + lenDailyPresenceCounter + lenDayDistance + (len(rec.GetActivityChangeInfo()) * lenActivityChangeInfo)
}
