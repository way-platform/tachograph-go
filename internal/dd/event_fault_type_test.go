package dd

import (
	"testing"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestParseEventFaultType checks that the protocol values of Appendix 1,
// EventFaultType, map to the right enum values. The protocol value is not the
// enum number: the two are related by the protocol_enum_value annotation.
func TestParseEventFaultType(t *testing.T) {
	var opts UnmarshalOptions
	for _, tt := range []struct {
		raw  byte
		want ddv1.EventFaultType
	}{
		{0x00, ddv1.EventFaultType_GENERAL_NO_FURTHER_DETAILS},
		{0x04, ddv1.EventFaultType_GENERAL_DRIVING_WITHOUT_APPROPRIATE_CARD},
		{0x05, ddv1.EventFaultType_GENERAL_CARD_INSERTION_WHILE_DRIVING},
		{0x07, ddv1.EventFaultType_GENERAL_OVER_SPEEDING},
		{0x10, ddv1.EventFaultType_VU_SEC_NO_FURTHER_DETAILS},
		{0x30, ddv1.EventFaultType_FAULT_REC_EQ_NO_FURTHER_DETAILS},
	} {
		got, unrecognized := opts.parseEventFaultType(tt.raw)
		if got != tt.want || unrecognized != 0 {
			t.Errorf("parseEventFaultType(%#02x) = %v, %d; want %v, 0", tt.raw, got, unrecognized, tt.want)
		}
		var marshalOpts MarshalOptions
		if back := marshalOpts.marshalEventFaultType(got, 0); back != tt.raw {
			t.Errorf("marshalEventFaultType(%v) = %#02x; want %#02x", got, back, tt.raw)
		}
	}
}

// TestParseEventFaultTypeUnrecognized checks that a value outside the data
// dictionary is preserved verbatim.
func TestParseEventFaultTypeUnrecognized(t *testing.T) {
	var opts UnmarshalOptions
	got, unrecognized := opts.parseEventFaultType(0x8A)
	if got != ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED || unrecognized != 0x8A {
		t.Errorf("parseEventFaultType(0x8A) = %v, %#02x; want UNSPECIFIED, 0x8a", got, unrecognized)
	}
	var marshalOpts MarshalOptions
	if back := marshalOpts.marshalEventFaultType(got, unrecognized); back != 0x8A {
		t.Errorf("marshalEventFaultType kept %#02x; want 0x8a", back)
	}
}

// TestParseEventFaultRecordPurpose checks the same mapping for the record
// purpose of Appendix 1, EventFaultRecordPurpose.
func TestParseEventFaultRecordPurpose(t *testing.T) {
	var opts UnmarshalOptions
	for _, tt := range []struct {
		raw  byte
		want ddv1.EventFaultRecordPurpose
	}{
		{0x00, ddv1.EventFaultRecordPurpose_TEN_MOST_RECENT},
		{0x02, ddv1.EventFaultRecordPurpose_FIVE_LONGEST_IN_LAST_365_DAYS},
		{0x04, ddv1.EventFaultRecordPurpose_MOST_SERIOUS_IN_LAST_10_DAYS},
		{0x06, ddv1.EventFaultRecordPurpose_FIRST_AFTER_LAST_CALIBRATION},
	} {
		got, unrecognized := opts.parseEventFaultRecordPurpose(tt.raw)
		if got != tt.want || unrecognized != 0 {
			t.Errorf("parseEventFaultRecordPurpose(%#02x) = %v, %d; want %v, 0", tt.raw, got, unrecognized, tt.want)
		}
		var marshalOpts MarshalOptions
		if back := marshalOpts.marshalEventFaultRecordPurpose(got, 0); back != tt.raw {
			t.Errorf("marshalEventFaultRecordPurpose(%v) = %#02x; want %#02x", got, back, tt.raw)
		}
	}
}
