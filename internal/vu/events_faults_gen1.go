package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalEventsAndFaultsGen1 parses Gen1 Events and Faults data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Events and Faults structure (from Data Dictionary and Appendix 7, Section 2.2.6.4 and 2.2.6.5):
//
// ASN.1 Definition:
//
//	VuEventsAndFaultsFirstGen ::= SEQUENCE {
//	    vuFaultData                   VuFaultDataFirstGen,
//	    vuEventData                   VuEventDataFirstGen,
//	    vuOverSpeedingControlData     VuOverSpeedingControlData,
//	    vuOverSpeedingEventData       VuOverSpeedingEventDataFirstGen,
//	    vuTimeAdjustmentData          VuTimeAdjustmentDataFirstGen,
//	    signature                     SignatureFirstGen
//	}
func unmarshalEventsAndFaultsGen1(value []byte) (*vuv1.EventsAndFaultsGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	eventsAndFaults := &vuv1.EventsAndFaultsGen1{}
	eventsAndFaults.SetRawData(value) // Store complete transfer value for painting
	offset := 0
	opts := dd.UnmarshalOptions{PreserveRawData: true}

	// Parse VuFaultData (1 byte count + fault records)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfVuFaults")
	}
	noOfVuFaults := data[offset]
	offset += 1

	faultRecords := make([]*ddv1.VuFaultRecord, noOfVuFaults)
	for i := 0; i < int(noOfVuFaults); i++ {
		const faultRecordSize = 82
		if offset+faultRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for VuFaultRecord %d", i)
		}
		faultRecord, err := opts.UnmarshalVuFaultRecord(data[offset : offset+faultRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuFaultRecord %d: %w", i, err)
		}
		faultRecords[i] = faultRecord
		offset += faultRecordSize
	}
	eventsAndFaults.SetFaults(faultRecords)

	// Parse VuEventData (1 byte count + event records)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfVuEvents")
	}
	noOfVuEvents := data[offset]
	offset += 1

	eventRecords := make([]*ddv1.VuEventRecord, noOfVuEvents)
	for i := 0; i < int(noOfVuEvents); i++ {
		const eventRecordSize = 83
		if offset+eventRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for VuEventRecord %d", i)
		}
		eventRecord, err := opts.UnmarshalVuEventRecord(data[offset : offset+eventRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuEventRecord %d: %w", i, err)
		}
		eventRecords[i] = eventRecord
		offset += eventRecordSize
	}
	eventsAndFaults.SetEvents(eventRecords)

	// Parse VuOverSpeedingControlData (9 bytes, no count byte)
	const overspeedControlSize = 9
	if offset+overspeedControlSize > len(data) {
		return nil, fmt.Errorf("insufficient data for VuOverSpeedingControlData")
	}
	overspeedControl, err := opts.UnmarshalVuOverspeedControlData(data[offset : offset+overspeedControlSize])
	if err != nil {
		return nil, fmt.Errorf("unmarshal VuOverSpeedingControlData: %w", err)
	}
	eventsAndFaults.SetOverspeedingControl(overspeedControl)
	offset += overspeedControlSize

	// Parse VuOverSpeedingEventData (1 byte count + overspeed event records)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfVuOverSpeedingEvents")
	}
	noOfVuOverSpeedingEvents := data[offset]
	offset += 1

	overspeedEventRecords := make([]*ddv1.VuOverspeedEventRecord, noOfVuOverSpeedingEvents)
	for i := 0; i < int(noOfVuOverSpeedingEvents); i++ {
		const overspeedEventRecordSize = 31
		if offset+overspeedEventRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for VuOverSpeedingEventRecord %d", i)
		}
		overspeedEventRecord, err := opts.UnmarshalVuOverspeedEventRecord(data[offset : offset+overspeedEventRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuOverSpeedingEventRecord %d: %w", i, err)
		}
		overspeedEventRecords[i] = overspeedEventRecord
		offset += overspeedEventRecordSize
	}
	eventsAndFaults.SetOverspeedingEvents(overspeedEventRecords)

	// Parse VuTimeAdjustmentData (1 byte count + time adjustment records)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfVuTimeAdjustments")
	}
	noOfVuTimeAdjustments := data[offset]
	offset += 1

	timeAdjustmentRecords := make([]*ddv1.VuTimeAdjustmentRecord, noOfVuTimeAdjustments)
	for i := 0; i < int(noOfVuTimeAdjustments); i++ {
		const timeAdjustmentRecordSize = 98
		if offset+timeAdjustmentRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for VuTimeAdjustmentRecord %d", i)
		}
		timeAdjustmentRecord, err := opts.UnmarshalVuTimeAdjustmentRecord(data[offset : offset+timeAdjustmentRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuTimeAdjustmentRecord %d: %w", i, err)
		}
		timeAdjustmentRecords[i] = timeAdjustmentRecord
		offset += timeAdjustmentRecordSize
	}
	eventsAndFaults.SetTimeAdjustments(timeAdjustmentRecords)

	// Verify we consumed all data
	if offset != len(data) {
		return nil, fmt.Errorf("unexpected extra data: consumed %d bytes, total %d bytes", offset, len(data))
	}

	// Store signature
	eventsAndFaults.SetSignature(signature)

	return eventsAndFaults, nil
}

// MarshalEventsAndFaultsGen1 marshals Gen1 Events and Faults data using raw data painting.
func (opts MarshalOptions) MarshalEventsAndFaultsGen1(eventsAndFaults *vuv1.EventsAndFaultsGen1) ([]byte, error) {
	if eventsAndFaults == nil {
		return nil, fmt.Errorf("eventsAndFaults cannot be nil")
	}

	// Calculate data size
	noOfFaults := len(eventsAndFaults.GetFaults())
	noOfEvents := len(eventsAndFaults.GetEvents())
	noOfOverspeedEvents := len(eventsAndFaults.GetOverspeedingEvents())
	noOfTimeAdjustments := len(eventsAndFaults.GetTimeAdjustments())

	dataSize := 1 + (noOfFaults * 82) + // VuFaultData: 1 byte count + records
		1 + (noOfEvents * 83) + // VuEventData: 1 byte count + records
		9 + // VuOverSpeedingControlData: fixed 9 bytes
		1 + (noOfOverspeedEvents * 31) + // VuOverSpeedingEventData: 1 byte count + records
		1 + (noOfTimeAdjustments * 98) // VuTimeAdjustmentData: 1 byte count + records

	// Use raw data as canvas if available
	var canvas []byte
	raw := eventsAndFaults.GetRawData()
	if len(raw) == dataSize+128 {
		canvas = make([]byte, dataSize)
		copy(canvas, raw[:dataSize])
	} else if len(raw) == dataSize {
		canvas = make([]byte, dataSize)
		copy(canvas, raw)
	} else {
		canvas = make([]byte, dataSize)
	}

	offset := 0
	marshalOpts := dd.MarshalOptions{}

	// Marshal VuFaultData
	canvas[offset] = byte(noOfFaults)
	offset += 1
	for i, faultRecord := range eventsAndFaults.GetFaults() {
		faultBytes, err := marshalOpts.MarshalVuFaultRecord(faultRecord)
		if err != nil {
			return nil, fmt.Errorf("marshal VuFaultRecord %d: %w", i, err)
		}
		if len(faultBytes) != 82 {
			return nil, fmt.Errorf("VuFaultRecord %d has invalid length: got %d, want 82", i, len(faultBytes))
		}
		copy(canvas[offset:offset+82], faultBytes)
		offset += 82
	}

	// Marshal VuEventData
	canvas[offset] = byte(noOfEvents)
	offset += 1
	for i, eventRecord := range eventsAndFaults.GetEvents() {
		eventBytes, err := marshalOpts.MarshalVuEventRecord(eventRecord)
		if err != nil {
			return nil, fmt.Errorf("marshal VuEventRecord %d: %w", i, err)
		}
		if len(eventBytes) != 83 {
			return nil, fmt.Errorf("VuEventRecord %d has invalid length: got %d, want 83", i, len(eventBytes))
		}
		copy(canvas[offset:offset+83], eventBytes)
		offset += 83
	}

	// Marshal VuOverSpeedingControlData
	overspeedControlBytes, err := marshalOpts.MarshalVuOverspeedControlData(eventsAndFaults.GetOverspeedingControl())
	if err != nil {
		return nil, fmt.Errorf("marshal VuOverSpeedingControlData: %w", err)
	}
	if len(overspeedControlBytes) != 9 {
		return nil, fmt.Errorf("VuOverSpeedingControlData has invalid length: got %d, want 9", len(overspeedControlBytes))
	}
	copy(canvas[offset:offset+9], overspeedControlBytes)
	offset += 9

	// Marshal VuOverSpeedingEventData
	canvas[offset] = byte(noOfOverspeedEvents)
	offset += 1
	for i, overspeedEvent := range eventsAndFaults.GetOverspeedingEvents() {
		overspeedEventBytes, err := marshalOpts.MarshalVuOverspeedEventRecord(overspeedEvent)
		if err != nil {
			return nil, fmt.Errorf("marshal VuOverSpeedingEventRecord %d: %w", i, err)
		}
		if len(overspeedEventBytes) != 31 {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d has invalid length: got %d, want 31", i, len(overspeedEventBytes))
		}
		copy(canvas[offset:offset+31], overspeedEventBytes)
		offset += 31
	}

	// Marshal VuTimeAdjustmentData
	canvas[offset] = byte(noOfTimeAdjustments)
	offset += 1
	for i, timeAdjustment := range eventsAndFaults.GetTimeAdjustments() {
		timeAdjustmentBytes, err := marshalOpts.MarshalVuTimeAdjustmentRecord(timeAdjustment)
		if err != nil {
			return nil, fmt.Errorf("marshal VuTimeAdjustmentRecord %d: %w", i, err)
		}
		if len(timeAdjustmentBytes) != 98 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d has invalid length: got %d, want 98", i, len(timeAdjustmentBytes))
		}
		copy(canvas[offset:offset+98], timeAdjustmentBytes)
		offset += 98
	}

	// Verify we wrote all data
	if offset != dataSize {
		return nil, fmt.Errorf("Events and Faults Gen1 marshalling mismatch: wrote %d bytes, expected %d", offset, dataSize)
	}

	// Append signature
	signature := eventsAndFaults.GetSignature()
	if len(signature) == 0 {
		signature = make([]byte, 128)
	}
	if len(signature) != 128 {
		return nil, fmt.Errorf("invalid signature length: got %d, want 128", len(signature))
	}
	transferValue := append(canvas, signature...)

	return transferValue, nil
}

// anonymizeEventsAndFaultsGen1 anonymizes Gen1 Events and Faults data.
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen1(ef *vuv1.EventsAndFaultsGen1) *vuv1.EventsAndFaultsGen1 {
	if ef == nil {
		return nil
	}
	result := proto.Clone(ef).(*vuv1.EventsAndFaultsGen1)

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize fault records
	var anonymizedFaults []*ddv1.VuFaultRecord
	for _, fault := range result.GetFaults() {
		anonymizedFaults = append(anonymizedFaults, ddOpts.AnonymizeVuFaultRecord(fault))
	}
	result.SetFaults(anonymizedFaults)

	// Anonymize event records
	var anonymizedEvents []*ddv1.VuEventRecord
	for _, event := range result.GetEvents() {
		anonymizedEvents = append(anonymizedEvents, ddOpts.AnonymizeVuEventRecord(event))
	}
	result.SetEvents(anonymizedEvents)

	// Anonymize overspeed event records
	var anonymizedOverspeedEvents []*ddv1.VuOverspeedEventRecord
	for _, overspeedEvent := range result.GetOverspeedingEvents() {
		anonymizedOverspeedEvents = append(anonymizedOverspeedEvents, ddOpts.AnonymizeVuOverspeedEventRecord(overspeedEvent))
	}
	result.SetOverspeedingEvents(anonymizedOverspeedEvents)

	// Anonymize time adjustment records
	var anonymizedTimeAdjustments []*ddv1.VuTimeAdjustmentRecord
	for _, timeAdj := range result.GetTimeAdjustments() {
		anonymizedTimeAdjustments = append(anonymizedTimeAdjustments, ddOpts.AnonymizeVuTimeAdjustmentRecord(timeAdj))
	}
	result.SetTimeAdjustments(anonymizedTimeAdjustments)

	// Overspeed control data has no PII - just keep as-is
	// (It only contains last overspeed control time, max speed, average speed, etc.)

	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Clear raw_data to force semantic marshalling
	result.ClearRawData()

	return result
}
