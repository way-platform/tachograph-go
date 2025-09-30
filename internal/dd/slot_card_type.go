package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalSlotCardType parses slot card type from raw data.
//
// The data type `SlotCardType` is specified in the Data Dictionary, Section 2.34.
//
// ASN.1 Definition:
//
//	CardSlotsStatus ::= OCTET STRING (SIZE(1))
//	-- Bitfield: 'ccccdddd'B where 'dddd' is driver slot, 'cccc' is co-driver slot
//	-- Each 4-bit nibble represents SlotCardType
//
// Binary Layout (4 bits):
//   - Slot Card Type (4 bits): Raw integer value (0-4)
//
//nolint:unused
func UnmarshalSlotCardType(data []byte) (ddv1.SlotCardType, error) {
	if len(data) != 1 {
		return ddv1.SlotCardType_SLOT_CARD_TYPE_UNSPECIFIED, fmt.Errorf("invalid data length for SlotCardType: got %d, want 1", len(data))
	}
	rawValue := int32(data[0])
	// Use the protocol enum value mapping
	if enumNumber, found := GetEnumForProtocolValue(ddv1.SlotCardType_SLOT_CARD_TYPE_UNSPECIFIED.Descriptor(), rawValue); found {
		return ddv1.SlotCardType(enumNumber), nil
	} else {
		return ddv1.SlotCardType_SLOT_CARD_TYPE_UNRECOGNIZED, nil
	}
}

// appendSlotCardType appends slot card type as a single byte.
//
// The data type `SlotCardType` is specified in the Data Dictionary, Section 2.34.
//
// ASN.1 Definition:
//
//	CardSlotsStatus ::= OCTET STRING (SIZE(1))
//	-- Bitfield: 'ccccdddd'B where 'dddd' is driver slot, 'cccc' is co-driver slot
//	-- Each 4-bit nibble represents SlotCardType
//
// Binary Layout (4 bits):
//   - Slot Card Type (4 bits): Raw integer value (0-4)
//
//nolint:unused
func AppendSlotCardType(dst []byte, slotCardType ddv1.SlotCardType) []byte {
	// Get the protocol value for the enum
	if protocolValue, ok := GetProtocolValueForEnum(slotCardType); ok {
		return append(dst, byte(protocolValue))
	}
	return append(dst, 0) // Default to 0 if not found
}

// mapSlotCardType maps a raw slot value to SlotCardType enum (legacy compatibility).
// This function maintains compatibility with existing code that uses direct mapping.
func MapSlotCardType(slotValue uint8) ddv1.SlotCardType {
	switch slotValue {
	case 0:
		return ddv1.SlotCardType_NO_CARD
	case 1:
		return ddv1.SlotCardType_DRIVER_CARD_INSERTED
	case 2:
		return ddv1.SlotCardType_WORKSHOP_CARD_INSERTED
	case 3:
		return ddv1.SlotCardType_CONTROL_CARD_INSERTED
	case 4:
		return ddv1.SlotCardType_COMPANY_CARD_INSERTED
	default:
		return ddv1.SlotCardType_SLOT_CARD_TYPE_UNSPECIFIED
	}
}
