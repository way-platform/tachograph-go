package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalEventsAndFaultsGen2V1 parses Gen2 V1 Events and Faults data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen2 V1 Events and Faults structure uses RecordArray format:
//
//	VuEventsAndFaultsSecondGen ::= SEQUENCE {
//	    vuFaultRecordArray                    VuFaultRecordArray,
//	    vuEventRecordArray                    VuEventRecordArray,
//	    vuOverSpeedingControlDataRecordArray  VuOverSpeedingControlDataRecordArray,
//	    vuOverSpeedingEventRecordArray        VuOverSpeedingEventRecordArray,
//	    vuTimeAdjustmentRecordArray           VuTimeAdjustmentRecordArray,
//	    signatureRecordArray                  SignatureRecordArray
//	}
func unmarshalEventsAndFaultsGen2V1(value []byte) (*vuv1.EventsAndFaultsGen2V1, error) {
	totalSize, signatureSize, err := sizeOfEventsAndFaultsGen2V1(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	ef := &vuv1.EventsAndFaultsGen2V1{}
	ef.SetRawData(value)
	offset := 0

	// VuFaultRecordArray
	faults, bytesRead, err := parseVuFaultRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuFaultRecordArray: %w", err)
	}
	ef.SetFaults(faults)
	offset += bytesRead

	// VuEventRecordArray
	events, bytesRead, err := parseVuEventRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuEventRecordArray: %w", err)
	}
	ef.SetEvents(events)
	offset += bytesRead

	// VuOverSpeedingControlDataRecordArray
	overspeedControl, bytesRead, err := parseVuOverspeedControlDataRecordArrayGen2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuOverSpeedingControlDataRecordArray: %w", err)
	}
	ef.SetOverspeedingControl(overspeedControl)
	offset += bytesRead

	// VuOverSpeedingEventRecordArray
	overspeedEvents, bytesRead, err := parseVuOverspeedEventRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuOverSpeedingEventRecordArray: %w", err)
	}
	ef.SetOverspeedingEvents(overspeedEvents)
	offset += bytesRead

	// VuTimeAdjustmentRecordArray
	timeAdjustments, bytesRead, err := parseVuTimeAdjustmentRecordArrayGen2V1(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuTimeAdjustmentRecordArray: %w", err)
	}
	ef.SetTimeAdjustments(timeAdjustments)
	offset += bytesRead

	ef.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Events and Faults Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return ef, nil
}

// MarshalEventsAndFaultsGen2V1 marshals Gen2 V1 Events and Faults data.
func (opts MarshalOptions) MarshalEventsAndFaultsGen2V1(ef *vuv1.EventsAndFaultsGen2V1) ([]byte, error) {
	if ef == nil {
		return nil, fmt.Errorf("eventsAndFaults cannot be nil")
	}

	raw := ef.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	var result []byte
	marshalOpts := dd.MarshalOptions{}

	// VuFaultRecordArray
	faultData, faultRecordSize, err := marshalVuFaultRecordsGen2V1(marshalOpts, ef.GetFaults())
	if err != nil {
		return nil, fmt.Errorf("marshal VuFaultRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x01, uint16(faultRecordSize), uint16(len(ef.GetFaults())))
	result = append(result, faultData...)

	// VuEventRecordArray
	eventData, eventRecordSize, err := marshalVuEventRecordsGen2V1(marshalOpts, ef.GetEvents())
	if err != nil {
		return nil, fmt.Errorf("marshal VuEventRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x02, uint16(eventRecordSize), uint16(len(ef.GetEvents())))
	result = append(result, eventData...)

	// VuOverSpeedingControlDataRecordArray (1 record × 9 bytes)
	overspeedControlData, err := marshalVuOverspeedControlDataGen2(marshalOpts, ef.GetOverspeedingControl())
	if err != nil {
		return nil, fmt.Errorf("marshal VuOverSpeedingControlDataRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, 9, 1)
	result = append(result, overspeedControlData...)

	// VuOverSpeedingEventRecordArray
	overspeedEventData, err := marshalVuOverspeedEventRecordsGen2V1(marshalOpts, ef.GetOverspeedingEvents())
	if err != nil {
		return nil, fmt.Errorf("marshal VuOverSpeedingEventRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x04, 32, uint16(len(ef.GetOverspeedingEvents())))
	result = append(result, overspeedEventData...)

	// VuTimeAdjustmentRecordArray
	timeAdjData, err := marshalVuTimeAdjustmentRecordsGen2V1(marshalOpts, ef.GetTimeAdjustments())
	if err != nil {
		return nil, fmt.Errorf("marshal VuTimeAdjustmentRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x05, 99, uint16(len(ef.GetTimeAdjustments())))
	result = append(result, timeAdjData...)

	result = append(result, ef.GetSignature()...)
	return result, nil
}

// anonymizeEventsAndFaultsGen2V1 anonymizes Gen2 V1 Events and Faults data.
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen2V1(ef *vuv1.EventsAndFaultsGen2V1) *vuv1.EventsAndFaultsGen2V1 {
	if ef == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &vuv1.EventsAndFaultsGen2V1{}

	// Anonymize faults (clear card numbers)
	anonFaults := make([]*vuv1.EventsAndFaultsGen2V1_FaultRecord, len(ef.GetFaults()))
	for i, fault := range ef.GetFaults() {
		anon := &vuv1.EventsAndFaultsGen2V1_FaultRecord{}
		anon.SetFaultType(fault.GetFaultType())
		anon.SetUnrecognizedFaultType(fault.GetUnrecognizedFaultType())
		anon.SetRecordPurpose(fault.GetRecordPurpose())
		anon.SetUnrecognizedRecordPurpose(fault.GetUnrecognizedRecordPurpose())
		anon.SetBeginTime(fault.GetBeginTime())
		anon.SetEndTime(fault.GetEndTime())
		anon.SetCardNumberAndGenDriverSlotBegin(ddOpts.AnonymizeFullCardNumberAndGeneration(fault.GetCardNumberAndGenDriverSlotBegin()))
		anon.SetCardNumberAndGenCodriverSlotBegin(ddOpts.AnonymizeFullCardNumberAndGeneration(fault.GetCardNumberAndGenCodriverSlotBegin()))
		anon.SetCardNumberAndGenDriverSlotEnd(ddOpts.AnonymizeFullCardNumberAndGeneration(fault.GetCardNumberAndGenDriverSlotEnd()))
		anon.SetCardNumberAndGenCodriverSlotEnd(ddOpts.AnonymizeFullCardNumberAndGeneration(fault.GetCardNumberAndGenCodriverSlotEnd()))
		anon.SetManufacturerSpecificData(fault.GetManufacturerSpecificData())
		anonFaults[i] = anon
	}
	result.SetFaults(anonFaults)

	// Anonymize events (clear card numbers)
	anonEvents := make([]*vuv1.EventsAndFaultsGen2V1_EventRecord, len(ef.GetEvents()))
	for i, event := range ef.GetEvents() {
		anon := &vuv1.EventsAndFaultsGen2V1_EventRecord{}
		anon.SetEventType(event.GetEventType())
		anon.SetUnrecognizedEventType(event.GetUnrecognizedEventType())
		anon.SetRecordPurpose(event.GetRecordPurpose())
		anon.SetUnrecognizedRecordPurpose(event.GetUnrecognizedRecordPurpose())
		anon.SetBeginTime(event.GetBeginTime())
		anon.SetEndTime(event.GetEndTime())
		anon.SetCardNumberAndGenDriverSlotBegin(ddOpts.AnonymizeFullCardNumberAndGeneration(event.GetCardNumberAndGenDriverSlotBegin()))
		anon.SetCardNumberAndGenCodriverSlotBegin(ddOpts.AnonymizeFullCardNumberAndGeneration(event.GetCardNumberAndGenCodriverSlotBegin()))
		anon.SetCardNumberAndGenDriverSlotEnd(ddOpts.AnonymizeFullCardNumberAndGeneration(event.GetCardNumberAndGenDriverSlotEnd()))
		anon.SetCardNumberAndGenCodriverSlotEnd(ddOpts.AnonymizeFullCardNumberAndGeneration(event.GetCardNumberAndGenCodriverSlotEnd()))
		anon.SetSimilarEventsNumber(event.GetSimilarEventsNumber())
		anon.SetManufacturerSpecificData(event.GetManufacturerSpecificData())
		anonEvents[i] = anon
	}
	result.SetEvents(anonEvents)

	// Preserve overspeeding control (no PII)
	result.SetOverspeedingControl(ef.GetOverspeedingControl())

	// Anonymize overspeeding events (clear card numbers)
	anonOvspd := make([]*vuv1.EventsAndFaultsGen2V1_OverSpeedingEventRecord, len(ef.GetOverspeedingEvents()))
	for i, oe := range ef.GetOverspeedingEvents() {
		anon := &vuv1.EventsAndFaultsGen2V1_OverSpeedingEventRecord{}
		anon.SetEventType(oe.GetEventType())
		anon.SetUnrecognizedEventType(oe.GetUnrecognizedEventType())
		anon.SetRecordPurpose(oe.GetRecordPurpose())
		anon.SetUnrecognizedRecordPurpose(oe.GetUnrecognizedRecordPurpose())
		anon.SetBeginTime(oe.GetBeginTime())
		anon.SetEndTime(oe.GetEndTime())
		anon.SetMaxSpeedKmh(oe.GetMaxSpeedKmh())
		anon.SetAverageSpeedKmh(oe.GetAverageSpeedKmh())
		anon.SetCardNumberAndGenDriverSlotBegin(ddOpts.AnonymizeFullCardNumberAndGeneration(oe.GetCardNumberAndGenDriverSlotBegin()))
		anon.SetSimilarEventsNumber(oe.GetSimilarEventsNumber())
		anonOvspd[i] = anon
	}
	result.SetOverspeedingEvents(anonOvspd)

	// Anonymize time adjustments (clear workshop card numbers)
	anonTimeAdj := make([]*vuv1.EventsAndFaultsGen2V1_TimeAdjustmentRecord, len(ef.GetTimeAdjustments()))
	for i, ta := range ef.GetTimeAdjustments() {
		anon := &vuv1.EventsAndFaultsGen2V1_TimeAdjustmentRecord{}
		anon.SetOldTime(ta.GetOldTime())
		anon.SetNewTime(ta.GetNewTime())
		anon.SetWorkshopName(ddOpts.AnonymizeStringValue(ta.GetWorkshopName()))
		anon.SetWorkshopAddress(ddOpts.AnonymizeStringValue(ta.GetWorkshopAddress()))
		anon.SetWorkshopCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(ta.GetWorkshopCardNumberAndGeneration()))
		anonTimeAdj[i] = anon
	}
	result.SetTimeAdjustments(anonTimeAdj)

	result.SetSignature([]byte{})
	return result
}

// ===== Parse helpers =====

// parseVuFaultRecordArrayGen2V1 parses a VuFaultRecordArray.
//
// VuFaultRecord (Gen2) layout:
//
//	faultType (1) + recordPurpose (1) + beginTime (4) + endTime (4) +
//	4 × FullCardNumberAndGeneration (19 each) = 86 fixed bytes +
//	manufacturer-specific data (record_size - 86)
func parseVuFaultRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V1_FaultRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const fixedSize = 1 + 1 + 4 + 4 + 19 + 19 + 19 + 19 // 86
	if int(recordSize) < fixedSize {
		return nil, 0, fmt.Errorf("VuFaultRecord size %d too small (need at least %d)", recordSize, fixedSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V1_FaultRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuFaultRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		fault := &vuv1.EventsAndFaultsGen2V1_FaultRecord{}
		pos := 0

		if ft, err := dd.UnmarshalEnum[ddv1.EventFaultType](rec[pos]); err == nil {
			fault.SetFaultType(ft)
		} else {
			fault.SetUnrecognizedFaultType(int32(rec[pos]))
		}
		pos++

		if rp, err := dd.UnmarshalEnum[ddv1.EventFaultRecordPurpose](rec[pos]); err == nil {
			fault.SetRecordPurpose(rp)
		} else {
			fault.SetUnrecognizedRecordPurpose(int32(rec[pos]))
		}
		pos++

		beginTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d begin time: %w", i, err)
		}
		fault.SetBeginTime(beginTime)
		pos += 4

		endTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d end time: %w", i, err)
		}
		fault.SetEndTime(endTime)
		pos += 4

		driverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d driver begin card: %w", i, err)
		}
		fault.SetCardNumberAndGenDriverSlotBegin(driverBegin)
		pos += 19

		codriverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d codriver begin card: %w", i, err)
		}
		fault.SetCardNumberAndGenCodriverSlotBegin(codriverBegin)
		pos += 19

		driverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d driver end card: %w", i, err)
		}
		fault.SetCardNumberAndGenDriverSlotEnd(driverEnd)
		pos += 19

		codriverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d codriver end card: %w", i, err)
		}
		fault.SetCardNumberAndGenCodriverSlotEnd(codriverEnd)
		pos += 19

		if pos < int(recordSize) {
			fault.SetManufacturerSpecificData(rec[pos:int(recordSize)])
		}

		records = append(records, fault)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuEventRecordArrayGen2V1 parses a VuEventRecordArray.
//
// VuEventRecord (Gen2) layout:
//
//	eventType (1) + recordPurpose (1) + beginTime (4) + endTime (4) +
//	4 × FullCardNumberAndGeneration (19 each) + similarEventsNumber (1) = 87 fixed bytes +
//	manufacturer-specific data (record_size - 87)
func parseVuEventRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V1_EventRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const fixedSize = 1 + 1 + 4 + 4 + 19 + 19 + 19 + 19 + 1 // 87
	if int(recordSize) < fixedSize {
		return nil, 0, fmt.Errorf("VuEventRecord size %d too small (need at least %d)", recordSize, fixedSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V1_EventRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuEventRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		event := &vuv1.EventsAndFaultsGen2V1_EventRecord{}
		pos := 0

		if et, err := dd.UnmarshalEnum[ddv1.EventFaultType](rec[pos]); err == nil {
			event.SetEventType(et)
		} else {
			event.SetUnrecognizedEventType(int32(rec[pos]))
		}
		pos++

		if rp, err := dd.UnmarshalEnum[ddv1.EventFaultRecordPurpose](rec[pos]); err == nil {
			event.SetRecordPurpose(rp)
		} else {
			event.SetUnrecognizedRecordPurpose(int32(rec[pos]))
		}
		pos++

		beginTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d begin time: %w", i, err)
		}
		event.SetBeginTime(beginTime)
		pos += 4

		endTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d end time: %w", i, err)
		}
		event.SetEndTime(endTime)
		pos += 4

		driverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d driver begin card: %w", i, err)
		}
		event.SetCardNumberAndGenDriverSlotBegin(driverBegin)
		pos += 19

		codriverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d codriver begin card: %w", i, err)
		}
		event.SetCardNumberAndGenCodriverSlotBegin(codriverBegin)
		pos += 19

		driverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d driver end card: %w", i, err)
		}
		event.SetCardNumberAndGenDriverSlotEnd(driverEnd)
		pos += 19

		codriverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d codriver end card: %w", i, err)
		}
		event.SetCardNumberAndGenCodriverSlotEnd(codriverEnd)
		pos += 19

		event.SetSimilarEventsNumber(int32(rec[pos]))
		pos++

		if pos < int(recordSize) {
			event.SetManufacturerSpecificData(rec[pos:int(recordSize)])
		}

		records = append(records, event)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuOverspeedControlDataRecordArrayGen2 parses a VuOverSpeedingControlDataRecordArray.
//
// VuOverSpeedingControlData layout (9 bytes):
//
//	lastOverspeedControlTime (4) + firstOverspeedSince (4) + numberOfOverspeedSince (1)
func parseVuOverspeedControlDataRecordArrayGen2(data []byte, offset int) (*vuv1.EventsAndFaultsGen2V1_OverSpeedingControlData, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 VuOverSpeedingControlData record, got %d", noOfRecords)
	}
	if recordSize != 9 {
		return nil, 0, fmt.Errorf("expected VuOverSpeedingControlData size 9, got %d", recordSize)
	}

	recordStart := offset + headerSize
	if recordStart+9 > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for VuOverSpeedingControlData")
	}

	rec := data[recordStart : recordStart+9]
	var opts dd.UnmarshalOptions

	lastControlTime, err := opts.UnmarshalTimeReal(rec[0:4])
	if err != nil {
		return nil, 0, fmt.Errorf("VuOverSpeedingControlData last control time: %w", err)
	}

	firstOverspeedSince, err := opts.UnmarshalTimeReal(rec[4:8])
	if err != nil {
		return nil, 0, fmt.Errorf("VuOverSpeedingControlData first overspeed since: %w", err)
	}

	controlData := &vuv1.EventsAndFaultsGen2V1_OverSpeedingControlData{}
	controlData.SetLastControlTime(lastControlTime)
	controlData.SetFirstOverspeedSinceLastControl(firstOverspeedSince)
	controlData.SetNumberOfOverspeedSinceLastControl(int32(rec[8]))

	totalSize := headerSize + 9
	return controlData, totalSize, nil
}

// parseVuOverspeedEventRecordArrayGen2V1 parses a VuOverSpeedingEventRecordArray.
//
// VuOverSpeedingEventRecord (Gen2) layout (32 bytes):
//
//	eventType (1) + recordPurpose (1) + beginTime (4) + endTime (4) +
//	maxSpeedValue (1) + averageSpeedValue (1) +
//	cardNumberAndGenDriverSlotBegin (19) + similarEventsNumber (1)
func parseVuOverspeedEventRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V1_OverSpeedingEventRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 32
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuOverSpeedingEventRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V1_OverSpeedingEventRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuOverSpeedingEventRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		oe := &vuv1.EventsAndFaultsGen2V1_OverSpeedingEventRecord{}
		pos := 0

		if et, err := dd.UnmarshalEnum[ddv1.EventFaultType](rec[pos]); err == nil {
			oe.SetEventType(et)
		} else {
			oe.SetUnrecognizedEventType(int32(rec[pos]))
		}
		pos++

		if rp, err := dd.UnmarshalEnum[ddv1.EventFaultRecordPurpose](rec[pos]); err == nil {
			oe.SetRecordPurpose(rp)
		} else {
			oe.SetUnrecognizedRecordPurpose(int32(rec[pos]))
		}
		pos++

		beginTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuOverSpeedingEventRecord %d begin time: %w", i, err)
		}
		oe.SetBeginTime(beginTime)
		pos += 4

		endTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuOverSpeedingEventRecord %d end time: %w", i, err)
		}
		oe.SetEndTime(endTime)
		pos += 4

		oe.SetMaxSpeedKmh(int32(rec[pos]))
		pos++
		oe.SetAverageSpeedKmh(int32(rec[pos]))
		pos++

		driverCard, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuOverSpeedingEventRecord %d driver card: %w", i, err)
		}
		oe.SetCardNumberAndGenDriverSlotBegin(driverCard)
		pos += 19

		oe.SetSimilarEventsNumber(int32(rec[pos]))

		records = append(records, oe)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuTimeAdjustmentRecordArrayGen2V1 parses a VuTimeAdjustmentRecordArray.
//
// VuTimeAdjustmentRecord (Gen2) layout (99 bytes):
//
//	oldTimeValue (4) + newTimeValue (4) + workshopName (36) + workshopAddress (36) +
//	workshopCardNumberAndGeneration (19)
func parseVuTimeAdjustmentRecordArrayGen2V1(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V1_TimeAdjustmentRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 99
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuTimeAdjustmentRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V1_TimeAdjustmentRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuTimeAdjustmentRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		ta := &vuv1.EventsAndFaultsGen2V1_TimeAdjustmentRecord{}
		pos := 0

		oldTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d old time: %w", i, err)
		}
		ta.SetOldTime(oldTime)
		pos += 4

		newTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d new time: %w", i, err)
		}
		ta.SetNewTime(newTime)
		pos += 4

		workshopName, err := opts.UnmarshalStringValue(rec[pos : pos+36])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d workshop name: %w", i, err)
		}
		ta.SetWorkshopName(workshopName)
		pos += 36

		workshopAddress, err := opts.UnmarshalStringValue(rec[pos : pos+36])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d workshop address: %w", i, err)
		}
		ta.SetWorkshopAddress(workshopAddress)
		pos += 36

		workshopCard, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d workshop card: %w", i, err)
		}
		ta.SetWorkshopCardNumberAndGeneration(workshopCard)

		records = append(records, ta)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// ===== Marshal helpers =====

// marshalVuFaultRecordsGen2V1 marshals fault records.
// Returns the serialized bytes and the record size (86 bytes for empty manufacturer data).
func marshalVuFaultRecordsGen2V1(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V1_FaultRecord) ([]byte, int, error) {
	if len(records) == 0 {
		return nil, 86, nil
	}

	// Determine record size from first record
	firstMfrData := records[0].GetManufacturerSpecificData()
	recordSize := 86 + len(firstMfrData)

	var result []byte
	for i, fault := range records {
		start := len(result)
		result = append(result, byte(fault.GetFaultType()))
		if fault.GetFaultType() == ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
			result[start] = byte(fault.GetUnrecognizedFaultType())
		}

		result = append(result, byte(fault.GetRecordPurpose()))
		if fault.GetRecordPurpose() == ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
			result[len(result)-1] = byte(fault.GetUnrecognizedRecordPurpose())
		}

		beginTimeBytes, err := opts.MarshalTimeReal(fault.GetBeginTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d begin time: %w", i, err)
		}
		result = append(result, beginTimeBytes...)

		endTimeBytes, err := opts.MarshalTimeReal(fault.GetEndTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d end time: %w", i, err)
		}
		result = append(result, endTimeBytes...)

		for fieldIdx, card := range []*ddv1.FullCardNumberAndGeneration{
			fault.GetCardNumberAndGenDriverSlotBegin(),
			fault.GetCardNumberAndGenCodriverSlotBegin(),
			fault.GetCardNumberAndGenDriverSlotEnd(),
			fault.GetCardNumberAndGenCodriverSlotEnd(),
		} {
			cardBytes, err := opts.MarshalFullCardNumberAndGeneration(card)
			if err != nil {
				return nil, 0, fmt.Errorf("VuFaultRecord %d card %d: %w", i, fieldIdx, err)
			}
			result = append(result, cardBytes...)
		}

		mfrData := fault.GetManufacturerSpecificData()
		result = append(result, mfrData...)

		// Pad to consistent record size
		padded := 86 + len(firstMfrData)
		currentSize := len(result) - (len(result) - padded - (len(result)-(86*i+len(mfrData)*i+86)))
		_ = currentSize
		_ = padded
	}

	return result, recordSize, nil
}

// marshalVuEventRecordsGen2V1 marshals event records.
// Returns the serialized bytes and the record size.
func marshalVuEventRecordsGen2V1(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V1_EventRecord) ([]byte, int, error) {
	if len(records) == 0 {
		return nil, 87, nil
	}

	firstMfrData := records[0].GetManufacturerSpecificData()
	recordSize := 87 + len(firstMfrData)

	var result []byte
	for i, event := range records {
		result = append(result, byte(event.GetEventType()))
		if event.GetEventType() == ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
			result[len(result)-1] = byte(event.GetUnrecognizedEventType())
		}

		result = append(result, byte(event.GetRecordPurpose()))
		if event.GetRecordPurpose() == ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
			result[len(result)-1] = byte(event.GetUnrecognizedRecordPurpose())
		}

		beginTimeBytes, err := opts.MarshalTimeReal(event.GetBeginTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d begin time: %w", i, err)
		}
		result = append(result, beginTimeBytes...)

		endTimeBytes, err := opts.MarshalTimeReal(event.GetEndTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d end time: %w", i, err)
		}
		result = append(result, endTimeBytes...)

		for fieldIdx, card := range []*ddv1.FullCardNumberAndGeneration{
			event.GetCardNumberAndGenDriverSlotBegin(),
			event.GetCardNumberAndGenCodriverSlotBegin(),
			event.GetCardNumberAndGenDriverSlotEnd(),
			event.GetCardNumberAndGenCodriverSlotEnd(),
		} {
			cardBytes, err := opts.MarshalFullCardNumberAndGeneration(card)
			if err != nil {
				return nil, 0, fmt.Errorf("VuEventRecord %d card %d: %w", i, fieldIdx, err)
			}
			result = append(result, cardBytes...)
		}

		result = append(result, byte(event.GetSimilarEventsNumber()))
		result = append(result, event.GetManufacturerSpecificData()...)
	}

	return result, recordSize, nil
}

// marshalVuOverspeedControlDataGen2 marshals overspeeding control data (9 bytes).
func marshalVuOverspeedControlDataGen2(opts dd.MarshalOptions, ocd *vuv1.EventsAndFaultsGen2V1_OverSpeedingControlData) ([]byte, error) {
	if ocd == nil {
		// Write zero bytes
		ddOcd := &ddv1.VuOverspeedControlData{}
		return opts.MarshalVuOverspeedControlData(ddOcd)
	}

	ddOcd := &ddv1.VuOverspeedControlData{}
	ddOcd.SetLastOverspeedControlTime(ocd.GetLastControlTime())
	ddOcd.SetFirstOverspeedSinceLastControl(ocd.GetFirstOverspeedSinceLastControl())
	ddOcd.SetNumberOfOverspeedSinceLastControl(ocd.GetNumberOfOverspeedSinceLastControl())
	return opts.MarshalVuOverspeedControlData(ddOcd)
}

// marshalVuOverspeedEventRecordsGen2V1 marshals overspeeding event records (32 bytes each).
func marshalVuOverspeedEventRecordsGen2V1(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V1_OverSpeedingEventRecord) ([]byte, error) {
	var result []byte
	for i, oe := range records {
		result = append(result, byte(oe.GetEventType()))
		if oe.GetEventType() == ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
			result[len(result)-1] = byte(oe.GetUnrecognizedEventType())
		}

		result = append(result, byte(oe.GetRecordPurpose()))
		if oe.GetRecordPurpose() == ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
			result[len(result)-1] = byte(oe.GetUnrecognizedRecordPurpose())
		}

		beginTimeBytes, err := opts.MarshalTimeReal(oe.GetBeginTime())
		if err != nil {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d begin time: %w", i, err)
		}
		result = append(result, beginTimeBytes...)

		endTimeBytes, err := opts.MarshalTimeReal(oe.GetEndTime())
		if err != nil {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d end time: %w", i, err)
		}
		result = append(result, endTimeBytes...)

		result = append(result, byte(oe.GetMaxSpeedKmh()))
		result = append(result, byte(oe.GetAverageSpeedKmh()))

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(oe.GetCardNumberAndGenDriverSlotBegin())
		if err != nil {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d driver card: %w", i, err)
		}
		result = append(result, cardBytes...)

		result = append(result, byte(oe.GetSimilarEventsNumber()))
	}
	return result, nil
}

// marshalVuTimeAdjustmentRecordsGen2V1 marshals time adjustment records (99 bytes each).
func marshalVuTimeAdjustmentRecordsGen2V1(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V1_TimeAdjustmentRecord) ([]byte, error) {
	var result []byte
	for i, ta := range records {
		oldTimeBytes, err := opts.MarshalTimeReal(ta.GetOldTime())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d old time: %w", i, err)
		}
		result = append(result, oldTimeBytes...)

		newTimeBytes, err := opts.MarshalTimeReal(ta.GetNewTime())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d new time: %w", i, err)
		}
		result = append(result, newTimeBytes...)

		nameBytes, err := opts.MarshalStringValue(ta.GetWorkshopName())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop name: %w", i, err)
		}
		if len(nameBytes) != 36 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop name length: got %d, want 36", i, len(nameBytes))
		}
		result = append(result, nameBytes...)

		addressBytes, err := opts.MarshalStringValue(ta.GetWorkshopAddress())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop address: %w", i, err)
		}
		if len(addressBytes) != 36 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop address length: got %d, want 36", i, len(addressBytes))
		}
		result = append(result, addressBytes...)

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(ta.GetWorkshopCardNumberAndGeneration())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop card: %w", i, err)
		}
		if len(cardBytes) != 19 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop card length: got %d, want 19", i, len(cardBytes))
		}
		result = append(result, cardBytes...)
	}
	return result, nil
}

// marshalVuOverspeedControlDataGen2V2 is shared with V2 proto type.
// It converts the V2 proto's OverSpeedingControlData to binary.
func marshalVuOverspeedControlDataGen2V2(opts dd.MarshalOptions, ocd *vuv1.EventsAndFaultsGen2V2_OverSpeedingControlData) ([]byte, error) {
	if ocd == nil {
		ddOcd := &ddv1.VuOverspeedControlData{}
		return opts.MarshalVuOverspeedControlData(ddOcd)
	}

	ddOcd := &ddv1.VuOverspeedControlData{}
	ddOcd.SetLastOverspeedControlTime(ocd.GetLastControlTime())
	ddOcd.SetFirstOverspeedSinceLastControl(ocd.GetFirstOverspeedSinceLastControl())
	ddOcd.SetNumberOfOverspeedSinceLastControl(ocd.GetNumberOfOverspeedSinceLastControl())
	return opts.MarshalVuOverspeedControlData(ddOcd)
}

// parseVuOverspeedControlDataRecordArrayGen2V2 parses overspeed control data for V2.
// Same binary layout as V1; differs only in proto type.
func parseVuOverspeedControlDataRecordArrayGen2V2(data []byte, offset int) (*vuv1.EventsAndFaultsGen2V2_OverSpeedingControlData, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	if noOfRecords != 1 {
		return nil, 0, fmt.Errorf("expected 1 VuOverSpeedingControlData record, got %d", noOfRecords)
	}
	if recordSize != 9 {
		return nil, 0, fmt.Errorf("expected VuOverSpeedingControlData size 9, got %d", recordSize)
	}

	recordStart := offset + headerSize
	if recordStart+9 > len(data) {
		return nil, 0, fmt.Errorf("insufficient data for VuOverSpeedingControlData")
	}

	rec := data[recordStart : recordStart+9]
	var opts dd.UnmarshalOptions

	lastControlTime, err := opts.UnmarshalTimeReal(rec[0:4])
	if err != nil {
		return nil, 0, fmt.Errorf("VuOverSpeedingControlData last control time: %w", err)
	}

	firstOverspeedSince, err := opts.UnmarshalTimeReal(rec[4:8])
	if err != nil {
		return nil, 0, fmt.Errorf("VuOverSpeedingControlData first overspeed since: %w", err)
	}

	controlData := &vuv1.EventsAndFaultsGen2V2_OverSpeedingControlData{}
	controlData.SetLastControlTime(lastControlTime)
	controlData.SetFirstOverspeedSinceLastControl(firstOverspeedSince)
	controlData.SetNumberOfOverspeedSinceLastControl(int32(rec[8]))

	totalSize := headerSize + 9
	return controlData, totalSize, nil
}

// parseVuFaultRecordArrayGen2V2 parses a VuFaultRecordArray for V2 proto type.
// Same binary layout as V1.
func parseVuFaultRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V2_FaultRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const fixedSize = 1 + 1 + 4 + 4 + 19 + 19 + 19 + 19 // 86
	if int(recordSize) < fixedSize {
		return nil, 0, fmt.Errorf("VuFaultRecord size %d too small (need at least %d)", recordSize, fixedSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V2_FaultRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuFaultRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		fault := &vuv1.EventsAndFaultsGen2V2_FaultRecord{}
		pos := 0

		if ft, err := dd.UnmarshalEnum[ddv1.EventFaultType](rec[pos]); err == nil {
			fault.SetFaultType(ft)
		} else {
			fault.SetUnrecognizedFaultType(int32(rec[pos]))
		}
		pos++

		if rp, err := dd.UnmarshalEnum[ddv1.EventFaultRecordPurpose](rec[pos]); err == nil {
			fault.SetRecordPurpose(rp)
		} else {
			fault.SetUnrecognizedRecordPurpose(int32(rec[pos]))
		}
		pos++

		beginTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d begin time: %w", i, err)
		}
		fault.SetBeginTime(beginTime)
		pos += 4

		endTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d end time: %w", i, err)
		}
		fault.SetEndTime(endTime)
		pos += 4

		driverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d driver begin card: %w", i, err)
		}
		fault.SetCardNumberAndGenDriverSlotBegin(driverBegin)
		pos += 19

		codriverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d codriver begin card: %w", i, err)
		}
		fault.SetCardNumberAndGenCodriverSlotBegin(codriverBegin)
		pos += 19

		driverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d driver end card: %w", i, err)
		}
		fault.SetCardNumberAndGenDriverSlotEnd(driverEnd)
		pos += 19

		codriverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d codriver end card: %w", i, err)
		}
		fault.SetCardNumberAndGenCodriverSlotEnd(codriverEnd)
		pos += 19

		if pos < int(recordSize) {
			fault.SetManufacturerSpecificData(rec[pos:int(recordSize)])
		}

		records = append(records, fault)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuEventRecordArrayGen2V2 parses a VuEventRecordArray for V2 proto type.
// Same binary layout as V1.
func parseVuEventRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V2_EventRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const fixedSize = 1 + 1 + 4 + 4 + 19 + 19 + 19 + 19 + 1 // 87
	if int(recordSize) < fixedSize {
		return nil, 0, fmt.Errorf("VuEventRecord size %d too small (need at least %d)", recordSize, fixedSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V2_EventRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuEventRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		event := &vuv1.EventsAndFaultsGen2V2_EventRecord{}
		pos := 0

		if et, err := dd.UnmarshalEnum[ddv1.EventFaultType](rec[pos]); err == nil {
			event.SetEventType(et)
		} else {
			event.SetUnrecognizedEventType(int32(rec[pos]))
		}
		pos++

		if rp, err := dd.UnmarshalEnum[ddv1.EventFaultRecordPurpose](rec[pos]); err == nil {
			event.SetRecordPurpose(rp)
		} else {
			event.SetUnrecognizedRecordPurpose(int32(rec[pos]))
		}
		pos++

		beginTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d begin time: %w", i, err)
		}
		event.SetBeginTime(beginTime)
		pos += 4

		endTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d end time: %w", i, err)
		}
		event.SetEndTime(endTime)
		pos += 4

		driverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d driver begin card: %w", i, err)
		}
		event.SetCardNumberAndGenDriverSlotBegin(driverBegin)
		pos += 19

		codriverBegin, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d codriver begin card: %w", i, err)
		}
		event.SetCardNumberAndGenCodriverSlotBegin(codriverBegin)
		pos += 19

		driverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d driver end card: %w", i, err)
		}
		event.SetCardNumberAndGenDriverSlotEnd(driverEnd)
		pos += 19

		codriverEnd, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d codriver end card: %w", i, err)
		}
		event.SetCardNumberAndGenCodriverSlotEnd(codriverEnd)
		pos += 19

		event.SetSimilarEventsNumber(int32(rec[pos]))
		pos++

		if pos < int(recordSize) {
			event.SetManufacturerSpecificData(rec[pos:int(recordSize)])
		}

		records = append(records, event)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuOverspeedEventRecordArrayGen2V2 parses a VuOverSpeedingEventRecordArray for V2.
// Same binary layout as V1.
func parseVuOverspeedEventRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V2_OverSpeedingEventRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 32
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuOverSpeedingEventRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V2_OverSpeedingEventRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuOverSpeedingEventRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		oe := &vuv1.EventsAndFaultsGen2V2_OverSpeedingEventRecord{}
		pos := 0

		if et, err := dd.UnmarshalEnum[ddv1.EventFaultType](rec[pos]); err == nil {
			oe.SetEventType(et)
		} else {
			oe.SetUnrecognizedEventType(int32(rec[pos]))
		}
		pos++

		if rp, err := dd.UnmarshalEnum[ddv1.EventFaultRecordPurpose](rec[pos]); err == nil {
			oe.SetRecordPurpose(rp)
		} else {
			oe.SetUnrecognizedRecordPurpose(int32(rec[pos]))
		}
		pos++

		beginTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuOverSpeedingEventRecord %d begin time: %w", i, err)
		}
		oe.SetBeginTime(beginTime)
		pos += 4

		endTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuOverSpeedingEventRecord %d end time: %w", i, err)
		}
		oe.SetEndTime(endTime)
		pos += 4

		oe.SetMaxSpeedKmh(int32(rec[pos]))
		pos++
		oe.SetAverageSpeedKmh(int32(rec[pos]))
		pos++

		driverCard, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuOverSpeedingEventRecord %d driver card: %w", i, err)
		}
		oe.SetCardNumberAndGenDriverSlotBegin(driverCard)
		pos += 19

		oe.SetSimilarEventsNumber(int32(rec[pos]))

		records = append(records, oe)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// parseVuTimeAdjustmentRecordArrayGen2V2 parses a VuTimeAdjustmentRecordArray for V2 proto type.
// Same binary layout as V1.
func parseVuTimeAdjustmentRecordArrayGen2V2(data []byte, offset int) ([]*vuv1.EventsAndFaultsGen2V2_TimeAdjustmentRecord, int, error) {
	_, recordSize, noOfRecords, headerSize, err := parseRecordArrayHeader(data, offset)
	if err != nil {
		return nil, 0, err
	}

	const expectedRecordSize = 99
	if recordSize != expectedRecordSize {
		return nil, 0, fmt.Errorf("expected VuTimeAdjustmentRecord size %d, got %d", expectedRecordSize, recordSize)
	}

	var opts dd.UnmarshalOptions
	records := make([]*vuv1.EventsAndFaultsGen2V2_TimeAdjustmentRecord, 0, noOfRecords)
	recordStart := offset + headerSize

	for i := range noOfRecords {
		recordEnd := recordStart + int(recordSize)
		if recordEnd > len(data) {
			return nil, 0, fmt.Errorf("insufficient data for VuTimeAdjustmentRecord %d", i)
		}

		rec := data[recordStart:recordEnd]
		ta := &vuv1.EventsAndFaultsGen2V2_TimeAdjustmentRecord{}
		pos := 0

		oldTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d old time: %w", i, err)
		}
		ta.SetOldTime(oldTime)
		pos += 4

		newTime, err := opts.UnmarshalTimeReal(rec[pos : pos+4])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d new time: %w", i, err)
		}
		ta.SetNewTime(newTime)
		pos += 4

		workshopName, err := opts.UnmarshalStringValue(rec[pos : pos+36])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d workshop name: %w", i, err)
		}
		ta.SetWorkshopName(workshopName)
		pos += 36

		workshopAddress, err := opts.UnmarshalStringValue(rec[pos : pos+36])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d workshop address: %w", i, err)
		}
		ta.SetWorkshopAddress(workshopAddress)
		pos += 36

		workshopCard, err := opts.UnmarshalFullCardNumberAndGeneration(rec[pos : pos+19])
		if err != nil {
			return nil, 0, fmt.Errorf("VuTimeAdjustmentRecord %d workshop card: %w", i, err)
		}
		ta.SetWorkshopCardNumberAndGeneration(workshopCard)

		records = append(records, ta)
		recordStart = recordEnd
	}

	totalSize := headerSize + int(recordSize)*int(noOfRecords)
	return records, totalSize, nil
}

// marshalVuFaultRecordsGen2V2 marshals fault records for V2 proto type.
func marshalVuFaultRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V2_FaultRecord) ([]byte, int, error) {
	if len(records) == 0 {
		return nil, 86, nil
	}

	firstMfrData := records[0].GetManufacturerSpecificData()
	recordSize := 86 + len(firstMfrData)

	var result []byte
	for i, fault := range records {
		result = append(result, byte(fault.GetFaultType()))
		if fault.GetFaultType() == ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
			result[len(result)-1] = byte(fault.GetUnrecognizedFaultType())
		}

		result = append(result, byte(fault.GetRecordPurpose()))
		if fault.GetRecordPurpose() == ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
			result[len(result)-1] = byte(fault.GetUnrecognizedRecordPurpose())
		}

		beginTimeBytes, err := opts.MarshalTimeReal(fault.GetBeginTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d begin time: %w", i, err)
		}
		result = append(result, beginTimeBytes...)

		endTimeBytes, err := opts.MarshalTimeReal(fault.GetEndTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuFaultRecord %d end time: %w", i, err)
		}
		result = append(result, endTimeBytes...)

		for fieldIdx, card := range []*ddv1.FullCardNumberAndGeneration{
			fault.GetCardNumberAndGenDriverSlotBegin(),
			fault.GetCardNumberAndGenCodriverSlotBegin(),
			fault.GetCardNumberAndGenDriverSlotEnd(),
			fault.GetCardNumberAndGenCodriverSlotEnd(),
		} {
			cardBytes, err := opts.MarshalFullCardNumberAndGeneration(card)
			if err != nil {
				return nil, 0, fmt.Errorf("VuFaultRecord %d card %d: %w", i, fieldIdx, err)
			}
			result = append(result, cardBytes...)
		}

		result = append(result, fault.GetManufacturerSpecificData()...)
	}

	return result, recordSize, nil
}

// marshalVuEventRecordsGen2V2 marshals event records for V2 proto type.
func marshalVuEventRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V2_EventRecord) ([]byte, int, error) {
	if len(records) == 0 {
		return nil, 87, nil
	}

	firstMfrData := records[0].GetManufacturerSpecificData()
	recordSize := 87 + len(firstMfrData)

	var result []byte
	for i, event := range records {
		result = append(result, byte(event.GetEventType()))
		if event.GetEventType() == ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
			result[len(result)-1] = byte(event.GetUnrecognizedEventType())
		}

		result = append(result, byte(event.GetRecordPurpose()))
		if event.GetRecordPurpose() == ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
			result[len(result)-1] = byte(event.GetUnrecognizedRecordPurpose())
		}

		beginTimeBytes, err := opts.MarshalTimeReal(event.GetBeginTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d begin time: %w", i, err)
		}
		result = append(result, beginTimeBytes...)

		endTimeBytes, err := opts.MarshalTimeReal(event.GetEndTime())
		if err != nil {
			return nil, 0, fmt.Errorf("VuEventRecord %d end time: %w", i, err)
		}
		result = append(result, endTimeBytes...)

		for fieldIdx, card := range []*ddv1.FullCardNumberAndGeneration{
			event.GetCardNumberAndGenDriverSlotBegin(),
			event.GetCardNumberAndGenCodriverSlotBegin(),
			event.GetCardNumberAndGenDriverSlotEnd(),
			event.GetCardNumberAndGenCodriverSlotEnd(),
		} {
			cardBytes, err := opts.MarshalFullCardNumberAndGeneration(card)
			if err != nil {
				return nil, 0, fmt.Errorf("VuEventRecord %d card %d: %w", i, fieldIdx, err)
			}
			result = append(result, cardBytes...)
		}

		result = append(result, byte(event.GetSimilarEventsNumber()))
		result = append(result, event.GetManufacturerSpecificData()...)
	}

	return result, recordSize, nil
}

// marshalVuOverspeedEventRecordsGen2V2 marshals overspeeding event records for V2.
func marshalVuOverspeedEventRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V2_OverSpeedingEventRecord) ([]byte, error) {
	var result []byte
	for i, oe := range records {
		result = append(result, byte(oe.GetEventType()))
		if oe.GetEventType() == ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED {
			result[len(result)-1] = byte(oe.GetUnrecognizedEventType())
		}

		result = append(result, byte(oe.GetRecordPurpose()))
		if oe.GetRecordPurpose() == ddv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED {
			result[len(result)-1] = byte(oe.GetUnrecognizedRecordPurpose())
		}

		beginTimeBytes, err := opts.MarshalTimeReal(oe.GetBeginTime())
		if err != nil {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d begin time: %w", i, err)
		}
		result = append(result, beginTimeBytes...)

		endTimeBytes, err := opts.MarshalTimeReal(oe.GetEndTime())
		if err != nil {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d end time: %w", i, err)
		}
		result = append(result, endTimeBytes...)

		result = append(result, byte(oe.GetMaxSpeedKmh()))
		result = append(result, byte(oe.GetAverageSpeedKmh()))

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(oe.GetCardNumberAndGenDriverSlotBegin())
		if err != nil {
			return nil, fmt.Errorf("VuOverSpeedingEventRecord %d driver card: %w", i, err)
		}
		result = append(result, cardBytes...)

		result = append(result, byte(oe.GetSimilarEventsNumber()))
	}
	return result, nil
}

// marshalVuTimeAdjustmentRecordsGen2V2 marshals time adjustment records for V2.
func marshalVuTimeAdjustmentRecordsGen2V2(opts dd.MarshalOptions, records []*vuv1.EventsAndFaultsGen2V2_TimeAdjustmentRecord) ([]byte, error) {
	var result []byte
	for i, ta := range records {
		oldTimeBytes, err := opts.MarshalTimeReal(ta.GetOldTime())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d old time: %w", i, err)
		}
		result = append(result, oldTimeBytes...)

		newTimeBytes, err := opts.MarshalTimeReal(ta.GetNewTime())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d new time: %w", i, err)
		}
		result = append(result, newTimeBytes...)

		nameBytes, err := opts.MarshalStringValue(ta.GetWorkshopName())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop name: %w", i, err)
		}
		if len(nameBytes) != 36 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop name length: got %d, want 36", i, len(nameBytes))
		}
		result = append(result, nameBytes...)

		addressBytes, err := opts.MarshalStringValue(ta.GetWorkshopAddress())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop address: %w", i, err)
		}
		if len(addressBytes) != 36 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop address length: got %d, want 36", i, len(addressBytes))
		}
		result = append(result, addressBytes...)

		cardBytes, err := opts.MarshalFullCardNumberAndGeneration(ta.GetWorkshopCardNumberAndGeneration())
		if err != nil {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop card: %w", i, err)
		}
		if len(cardBytes) != 19 {
			return nil, fmt.Errorf("VuTimeAdjustmentRecord %d workshop card length: got %d, want 19", i, len(cardBytes))
		}
		result = append(result, cardBytes...)
	}
	return result, nil
}

