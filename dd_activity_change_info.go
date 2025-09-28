package tachograph

import (
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalActivityChangeInfo parses a 2-byte ActivityChangeInfo bitfield.
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
func unmarshalActivityChangeInfo(data []byte) (*ddv1.ActivityChangeInfo, error) {
	const (
		lenActivityChangeInfo = 2
	)

	if len(data) < lenActivityChangeInfo {
		return nil, fmt.Errorf("insufficient data for ActivityChangeInfo: got %d, want %d", len(data), lenActivityChangeInfo)
	}

	// Parse 2-byte bitfield according to spec
	changeData := binary.BigEndian.Uint16(data[0:2])

	// Skip invalid entries (all zeros or all ones)
	if changeData == 0 || changeData == 0xFFFF {
		return nil, fmt.Errorf("invalid ActivityChangeInfo: all zeros or all ones")
	}

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

	return activityChange, nil
}

// appendActivityChangeInfo appends the binary representation of ActivityChangeInfo to dst.
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
func appendActivityChangeInfo(dst []byte, ac *ddv1.ActivityChangeInfo) ([]byte, error) {
	if ac == nil {
		return dst, nil
	}

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
