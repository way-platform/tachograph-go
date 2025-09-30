package dd

import (
	"bytes"
	"encoding/binary"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalActivityChangeInfo parses a single ActivityChangeInfo record from a 2-byte slice.
//
// The data type `ActivityChangeInfo` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Specification:
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
func (opts UnmarshalOptions) UnmarshalActivityChangeInfo(input []byte) (*ddv1.ActivityChangeInfo, error) {
	const lenActivityChangeInfo = 2
	if len(input) != lenActivityChangeInfo {
		return nil, fmt.Errorf("invalid data length for ActivityChangeInfo: got %d, want %d", len(input), lenActivityChangeInfo)
	}
	var output ddv1.ActivityChangeInfo
	output.SetRawData(bytes.Clone(input))
	value := binary.BigEndian.Uint16(input)
	slot := int32((value >> 15) & 0x1)          // bit 15
	drivingStatus := int32((value >> 14) & 0x1) // bit 14
	cardStatus := (value >> 13) & 0x1           // bit 13
	activity := int32((value >> 11) & 0x3)      // bits 12-11
	timeMinutes := int32(value & 0x7FF)         // bits 10-0

	if slotEnum, err := UnmarshalEnum[ddv1.CardSlotNumber](byte(slot)); err == nil {
		output.SetSlot(slotEnum)
	} else {
		return nil, fmt.Errorf("invalid slot value %d: %w", slot, err)
	}

	if drivingStatusEnum, err := UnmarshalEnum[ddv1.DrivingStatus](byte(drivingStatus)); err == nil {
		output.SetDrivingStatus(drivingStatusEnum)
	} else {
		return nil, fmt.Errorf("invalid driving status value %d: %w", drivingStatus, err)
	}

	output.SetInserted(cardStatus == 0)

	if activityEnum, err := UnmarshalEnum[ddv1.DriverActivityValue](byte(activity)); err == nil {
		output.SetActivity(activityEnum)
	} else {
		return nil, fmt.Errorf("invalid activity value %d: %w", activity, err)
	}
	output.SetTimeOfChangeMinutes(timeMinutes)
	return &output, nil
}

// AppendActivityChangeInfo appends the binary representation of ActivityChangeInfo to dst.
//
// The data type `ActivityChangeInfo` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Specification:
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
func AppendActivityChangeInfo(dst []byte, ac *ddv1.ActivityChangeInfo) ([]byte, error) {
	const lenActivityChangeInfo = 2
	var canvas [lenActivityChangeInfo]byte
	if ac.HasRawData() {
		if len(ac.GetRawData()) != lenActivityChangeInfo {
			return nil, fmt.Errorf(
				"invalid raw_data length for ActivityChangeInfo: got %d, want %d",
				len(ac.GetRawData()), lenActivityChangeInfo,
			)
		}
		copy(canvas[:], ac.GetRawData())
	}
	slot, err := MarshalEnum(ac.GetSlot())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot: %w", err)
	}
	drivingStatus, err := MarshalEnum(ac.GetDrivingStatus())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal driving status: %w", err)
	}
	// Convert boolean to protocol value: 0 = inserted, 1 = not inserted
	cardInserted := int32(1)
	if ac.GetInserted() {
		cardInserted = 0
	}
	activity, err := MarshalEnum(ac.GetActivity())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activity: %w", err)
	}
	var aci uint16
	aci |= (uint16(slot) & 0x1) << 15
	aci |= (uint16(drivingStatus) & 0x1) << 14
	aci |= (uint16(cardInserted) & 0x1) << 13
	aci |= (uint16(activity) & 0x3) << 11
	aci |= (uint16(ac.GetTimeOfChangeMinutes()) & 0x7FF)
	binary.BigEndian.PutUint16(canvas[:], aci)
	return append(dst, canvas[:]...), nil
}
