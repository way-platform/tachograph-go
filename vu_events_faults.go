package tachograph

import (
	"bytes"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

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
func unmarshalVuEventsAndFaults(data []byte, offset int, target *vuv1.EventsAndFaults, generation int) (int, error) {
	startOffset := offset

	// Set generation
	if generation == 1 {
		target.SetGeneration(ddv1.Generation_GENERATION_1)
	} else {
		target.SetGeneration(ddv1.Generation_GENERATION_2)
	}

	// For now, implement a simplified version that just reads the data
	// without fully parsing all the complex structures
	// This ensures the interface is complete while allowing for future enhancement

	// Read all remaining data as signature for now
	remainingData, offset, err := readBytesFromBytes(data, offset, len(data)-offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read events and faults data: %w", err)
	}

	// Set as signature based on generation
	if generation == 1 {
		target.SetSignatureGen1(remainingData)
	} else {
		target.SetSignatureGen2(remainingData)
	}

	return offset - startOffset, nil
}

// AppendVuEventsAndFaults appends VU events and faults data to a buffer.
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
func appendVuEventsAndFaults(buf *bytes.Buffer, eventsAndFaults *vuv1.EventsAndFaults) error {
	if eventsAndFaults == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if eventsAndFaults.GetGeneration() == ddv1.Generation_GENERATION_1 {
		signature := eventsAndFaults.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := eventsAndFaults.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}
