package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalEventsAndFaultsGen2V2 parses Gen2 V2 Events and Faults data from the complete transfer value.
//
// Gen2 V2 adds VuTimeAdjustmentGNSSRecordArray after the regular VuTimeAdjustmentRecordArray.
// Both arrays are parsed into the same time_adjustments field.
//
// Structure:
//
//	VuEventsAndFaultsSecondGenV2 ::= SEQUENCE {
//	    vuFaultRecordArray                    VuFaultRecordArray,
//	    vuEventRecordArray                    VuEventRecordArray,
//	    vuOverSpeedingControlDataRecordArray  VuOverSpeedingControlDataRecordArray,
//	    vuOverSpeedingEventRecordArray        VuOverSpeedingEventRecordArray,
//	    vuTimeAdjustmentRecordArray           VuTimeAdjustmentRecordArray,
//	    vuTimeAdjustmentGNSSRecordArray       VuTimeAdjustmentGNSSRecordArray,
//	    signatureRecordArray                  SignatureRecordArray
//	}
func unmarshalEventsAndFaultsGen2V2(value []byte) (*vuv1.EventsAndFaultsGen2V2, error) {
	totalSize, signatureSize, err := sizeOfEventsAndFaultsGen2V2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	ef := &vuv1.EventsAndFaultsGen2V2{}
	ef.SetRawData(value)
	offset := 0

	// VuFaultRecordArray
	faults, bytesRead, err := parseVuFaultRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuFaultRecordArray: %w", err)
	}
	ef.SetFaults(faults)
	offset += bytesRead

	// VuEventRecordArray
	events, bytesRead, err := parseVuEventRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuEventRecordArray: %w", err)
	}
	ef.SetEvents(events)
	offset += bytesRead

	// VuOverSpeedingControlDataRecordArray
	overspeedControl, bytesRead, err := parseVuOverspeedControlDataRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuOverSpeedingControlDataRecordArray: %w", err)
	}
	ef.SetOverspeedingControl(overspeedControl)
	offset += bytesRead

	// VuOverSpeedingEventRecordArray
	overspeedEvents, bytesRead, err := parseVuOverspeedEventRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuOverSpeedingEventRecordArray: %w", err)
	}
	ef.SetOverspeedingEvents(overspeedEvents)
	offset += bytesRead

	// VuTimeAdjustmentRecordArray
	timeAdjs, bytesRead, err := parseVuTimeAdjustmentRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuTimeAdjustmentRecordArray: %w", err)
	}
	offset += bytesRead

	// VuTimeAdjustmentGNSSRecordArray — appended to same field
	gnssAdjs, bytesRead, err := parseVuTimeAdjustmentRecordArrayGen2V2(data, offset)
	if err != nil {
		return nil, fmt.Errorf("parse VuTimeAdjustmentGNSSRecordArray: %w", err)
	}
	ef.SetTimeAdjustments(append(timeAdjs, gnssAdjs...))
	offset += bytesRead

	ef.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Events and Faults Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return ef, nil
}

// MarshalEventsAndFaultsGen2V2 marshals Gen2 V2 Events and Faults data.
func (opts MarshalOptions) MarshalEventsAndFaultsGen2V2(ef *vuv1.EventsAndFaultsGen2V2) ([]byte, error) {
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
	faultData, faultRecordSize, err := marshalVuFaultRecordsGen2V2(marshalOpts, ef.GetFaults())
	if err != nil {
		return nil, fmt.Errorf("marshal VuFaultRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x01, uint16(faultRecordSize), uint16(len(ef.GetFaults())))
	result = append(result, faultData...)

	// VuEventRecordArray
	eventData, eventRecordSize, err := marshalVuEventRecordsGen2V2(marshalOpts, ef.GetEvents())
	if err != nil {
		return nil, fmt.Errorf("marshal VuEventRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x02, uint16(eventRecordSize), uint16(len(ef.GetEvents())))
	result = append(result, eventData...)

	// VuOverSpeedingControlDataRecordArray (1 record × 9 bytes)
	overspeedControlData, err := marshalVuOverspeedControlDataGen2V2(marshalOpts, ef.GetOverspeedingControl())
	if err != nil {
		return nil, fmt.Errorf("marshal VuOverSpeedingControlDataRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x03, 9, 1)
	result = append(result, overspeedControlData...)

	// VuOverSpeedingEventRecordArray
	overspeedEventData, err := marshalVuOverspeedEventRecordsGen2V2(marshalOpts, ef.GetOverspeedingEvents())
	if err != nil {
		return nil, fmt.Errorf("marshal VuOverSpeedingEventRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x04, 32, uint16(len(ef.GetOverspeedingEvents())))
	result = append(result, overspeedEventData...)

	// VuTimeAdjustmentRecordArray (all adjustments; GNSS distinction lost without raw_data)
	timeAdjData, err := marshalVuTimeAdjustmentRecordsGen2V2(marshalOpts, ef.GetTimeAdjustments())
	if err != nil {
		return nil, fmt.Errorf("marshal VuTimeAdjustmentRecordArray: %w", err)
	}
	result = appendRecordArrayHeader(result, 0x05, 99, uint16(len(ef.GetTimeAdjustments())))
	result = append(result, timeAdjData...)

	// VuTimeAdjustmentGNSSRecordArray (empty — GNSS/regular distinction not preserved in proto)
	result = appendRecordArrayHeader(result, 0x06, 99, 0)

	result = append(result, ef.GetSignature()...)
	return result, nil
}

// anonymizeEventsAndFaultsGen2V2 anonymizes Gen2 V2 Events and Faults data.
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen2V2(ef *vuv1.EventsAndFaultsGen2V2) *vuv1.EventsAndFaultsGen2V2 {
	if ef == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &vuv1.EventsAndFaultsGen2V2{}

	// Anonymize faults (clear card numbers)
	anonFaults := make([]*vuv1.EventsAndFaultsGen2V2_FaultRecord, len(ef.GetFaults()))
	for i, fault := range ef.GetFaults() {
		anon := &vuv1.EventsAndFaultsGen2V2_FaultRecord{}
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
	anonEvents := make([]*vuv1.EventsAndFaultsGen2V2_EventRecord, len(ef.GetEvents()))
	for i, event := range ef.GetEvents() {
		anon := &vuv1.EventsAndFaultsGen2V2_EventRecord{}
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
	anonOvspd := make([]*vuv1.EventsAndFaultsGen2V2_OverSpeedingEventRecord, len(ef.GetOverspeedingEvents()))
	for i, oe := range ef.GetOverspeedingEvents() {
		anon := &vuv1.EventsAndFaultsGen2V2_OverSpeedingEventRecord{}
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
	anonTimeAdj := make([]*vuv1.EventsAndFaultsGen2V2_TimeAdjustmentRecord, len(ef.GetTimeAdjustments()))
	for i, ta := range ef.GetTimeAdjustments() {
		anon := &vuv1.EventsAndFaultsGen2V2_TimeAdjustmentRecord{}
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
