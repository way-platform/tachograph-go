package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfEventsAndFaults dispatches to generation-specific size calculation.
func sizeOfEventsAndFaults(data []byte, transferType vuv1.TransferType) (int, error) {
	switch transferType {
	case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
		return sizeOfEventsAndFaultsGen1(data)
	case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
		return sizeOfEventsAndFaultsGen2V1(data)
	case vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2:
		return sizeOfEventsAndFaultsGen2V2(data)
	default:
		return 0, fmt.Errorf("unsupported transfer type for EventsAndFaults: %v", transferType)
	}
}

// sizeOfEventsAndFaultsGen1 calculates total size for Gen1 Events and Faults including signature.
//
// Events and Faults Gen1 structure (from Appendix 7, Section 2.2.6.4):
// - VuFaultData: 1 byte (noOfVuFaults) + (noOfVuFaults * 82 bytes per record)
//   - VuFaultRecord (Data Dictionary 2.201): 82 bytes total
//   - faultType: EventFaultType = 1 byte (OCTET STRING SIZE(1))
//   - faultRecordPurpose: EventFaultRecordPurpose = 1 byte (OCTET STRING SIZE(1))
//   - faultBeginTime: TimeReal = 4 bytes
//   - faultEndTime: TimeReal = 4 bytes
//   - cardNumberDriverSlotBegin: FullCardNumber = 18 bytes (1+1+16)
//   - cardNumberCodriverSlotBegin: FullCardNumber = 18 bytes
//   - cardNumberDriverSlotEnd: FullCardNumber = 18 bytes
//   - cardNumberCodriverSlotEnd: FullCardNumber = 18 bytes
//
// - VuEventData: 1 byte (noOfVuEvents) + (noOfVuEvents * 83 bytes per record)
//   - VuEventRecord (Data Dictionary 2.198): 83 bytes total
//   - eventType: EventFaultType = 1 byte
//   - eventRecordPurpose: EventFaultRecordPurpose = 1 byte
//   - eventBeginTime: TimeReal = 4 bytes
//   - eventEndTime: TimeReal = 4 bytes
//   - cardNumberDriverSlotBegin: FullCardNumber = 18 bytes
//   - cardNumberCodriverSlotBegin: FullCardNumber = 18 bytes
//   - cardNumberDriverSlotEnd: FullCardNumber = 18 bytes
//   - cardNumberCodriverSlotEnd: FullCardNumber = 18 bytes
//   - similarEventsNumber: SimilarEventsNumber = 1 byte
//
// - VuOverSpeedingControlData: 9 bytes (fixed structure, no count)
//   - lastOverspeedControlTime: TimeReal = 4 bytes
//   - firstOverspeedSince: TimeReal = 4 bytes
//   - numberOfOverspeedSince: OverspeedNumber = 1 byte
//
// - VuOverSpeedingEventData: 1 byte (noOfVuOverSpeedingEvents) + (noOfVuOverSpeedingEvents * 31 bytes per record)
//   - VuOverSpeedingEventRecord: 31 bytes total
//   - eventType: EventFaultType = 1 byte
//   - eventRecordPurpose: EventFaultRecordPurpose = 1 byte
//   - eventBeginTime: TimeReal = 4 bytes
//   - eventEndTime: TimeReal = 4 bytes
//   - maxSpeedValue: SpeedMax = 1 byte
//   - averageSpeedValue: SpeedAverage = 1 byte
//   - cardNumberDriverSlotBegin: FullCardNumber = 18 bytes
//   - similarEventsNumber: SimilarEventsNumber = 1 byte
//
// - VuTimeAdjustmentData: 1 byte (noOfVuTimeAdjRecords) + (noOfVuTimeAdjRecords * 98 bytes per record)
//   - VuTimeAdjustmentRecord: 98 bytes total
//   - oldTimeValue: TimeReal = 4 bytes
//   - newTimeValue: TimeReal = 4 bytes
//   - workshopName: Name = 36 bytes (1 codepage + 35 bytes)
//   - workshopAddress: Address = 36 bytes (1 codepage + 35 bytes)
//   - workshopCardNumber: FullCardNumber = 18 bytes
//
// - Signature: 128 bytes (RSA)
func sizeOfEventsAndFaultsGen1(data []byte) (int, error) {
	offset := 0

	// VuFaultData: 1 byte count + variable fault records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfVuFaults")
	}
	noOfVuFaults := data[offset]
	offset += 1

	// Each VuFaultRecord: 82 bytes (per Data Dictionary 2.201)
	const vuFaultRecordSize = 82
	offset += int(noOfVuFaults) * vuFaultRecordSize

	// VuEventData: 1 byte count + variable event records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfVuEvents")
	}
	noOfVuEvents := data[offset]
	offset += 1

	// Each VuEventRecord: 83 bytes (per Data Dictionary 2.198)
	const vuEventRecordSize = 83
	offset += int(noOfVuEvents) * vuEventRecordSize

	// VuOverSpeedingControlData: 9 bytes (fixed structure)
	offset += 9

	// VuOverSpeedingEventData: 1 byte count + variable overspeed records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfVuOverSpeedingEvents")
	}
	noOfVuOverSpeedingEvents := data[offset]
	offset += 1

	// Each VuOverSpeedingEventRecord: 31 bytes
	const vuOverSpeedingEventRecordSize = 31
	offset += int(noOfVuOverSpeedingEvents) * vuOverSpeedingEventRecordSize

	// VuTimeAdjustmentData: 1 byte count + variable time adjustment records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfVuTimeAdjRecords")
	}
	noOfVuTimeAdjRecords := data[offset]
	offset += 1

	// Each VuTimeAdjustmentRecord: 98 bytes
	const vuTimeAdjustmentRecordSize = 98
	offset += int(noOfVuTimeAdjRecords) * vuTimeAdjustmentRecordSize

	// Signature: 128 bytes for Gen1 RSA
	offset += 128

	return offset, nil
}

// sizeOfEventsAndFaultsGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
func sizeOfEventsAndFaultsGen2V1(data []byte) (int, error) {
	offset := 0

	// VuFaultRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuFaultRecordArray: %w", err)
	}
	offset += size

	// VuEventRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuEventRecordArray: %w", err)
	}
	offset += size

	// VuOverSpeedingControlDataRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuOverSpeedingControlDataRecordArray: %w", err)
	}
	offset += size

	// VuOverSpeedingEventRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuOverSpeedingEventRecordArray: %w", err)
	}
	offset += size

	// VuTimeAdjustmentRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuTimeAdjustmentRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// sizeOfEventsAndFaultsGen2V2 calculates size by parsing all Gen2 V2 RecordArrays.
// Gen2 V2 has an additional VuTimeAdjustmentGNSSRecordArray.
func sizeOfEventsAndFaultsGen2V2(data []byte) (int, error) {
	offset := 0

	// VuFaultRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuFaultRecordArray: %w", err)
	}
	offset += size

	// VuEventRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuEventRecordArray: %w", err)
	}
	offset += size

	// VuOverSpeedingControlDataRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuOverSpeedingControlDataRecordArray: %w", err)
	}
	offset += size

	// VuOverSpeedingEventRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuOverSpeedingEventRecordArray: %w", err)
	}
	offset += size

	// VuTimeAdjustmentRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuTimeAdjustmentRecordArray: %w", err)
	}
	offset += size

	// VuTimeAdjustmentGNSSRecordArray (Gen2 V2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuTimeAdjustmentGNSSRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// ===== Unmarshal Functions =====

// UnmarshalVuEventsAndFaults unmarshals VU events and faults data from a VU transfer.
//
// The data type `VuEventsAndFaults` is specified in the Data Dictionary, Section 2.2.6.3.
//
// ASN.1 Definition:
//
//	VuEventsAndFaultsFirstGen ::= SEQUENCE {
//	    vuEventData                       VuEventData,
//	    vuFaultData                       VuFaultData,
//	    vuOverSpeedingEventData           VuOverSpeedingEventData,
//	    vuTimeAdjustmentData              VuTimeAdjustmentData,
//	    signature                         SignatureFirstGen
//	}
//
//	VuEventsAndFaultsSecondGen ::= SEQUENCE {
//	    vuEventRecordArray                VuEventRecordArray,
//	    vuFaultRecordArray                VuFaultRecordArray,
//	    vuOverSpeedingEventRecordArray    VuOverSpeedingEventRecordArray,
//	    vuTimeAdjustmentRecordArray       VuTimeAdjustmentRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}
