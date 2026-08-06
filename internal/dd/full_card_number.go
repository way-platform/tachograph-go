package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalFullCardNumber parses full card number data.
//
// The data type `FullCardNumber` is specified in the Data Dictionary, Section 2.73.
//
// ASN.1 Definition:
//
//	FullCardNumber ::= SEQUENCE {
//	    cardType EquipmentType,
//	    cardIssuingMemberState NationNumeric,
//	    cardNumber CardNumber
//	}
//
//	CardNumber ::= CHOICE {
//	    driverIdentification   SEQUENCE { ... },
//	    ownerIdentification    SEQUENCE { ... }
//	}
//
// Binary Layout (fixed length, 18 bytes):
//   - Card Type (1 byte): EquipmentType
//   - Issuing Member State (1 byte): NationNumeric
//   - Card Number (16 bytes): CardNumber CHOICE based on card type (padded to 16 bytes)
func (opts UnmarshalOptions) UnmarshalFullCardNumber(data []byte) (*ddv1.FullCardNumber, error) {
	const lenFullCardNumber = 18

	if len(data) != lenFullCardNumber {
		return nil, fmt.Errorf("invalid data length for FullCardNumber: got %d, want %d", len(data), lenFullCardNumber)
	}

	cardNumber := &ddv1.FullCardNumber{}

	// Preserve raw data for round-trip fidelity
	if opts.PreserveRawData {
		cardNumber.SetRawData(data)
	}

	// Special case: Various byte patterns indicate "no card inserted" or invalid card data.
	// Common patterns: 0xFF (255), 0x00 (0), or any unrecognized equipment type
	// In these cases, we return an empty FullCardNumber with only raw_data preserved.
	cardTypeByte := data[0]
	if cardTypeByte == 0xFF || cardTypeByte == 0x00 {
		// Return empty card number with only raw_data preserved
		return cardNumber, nil
	}

	// Parse card type (1 byte)
	// If the card type is unrecognized, treat it as an empty/invalid card
	cardType, err := UnmarshalEnum[ddv1.EquipmentType](cardTypeByte)
	if err != nil {
		// Unrecognized card type - treat as empty card
		return cardNumber, nil
	}
	cardNumber.SetCardType(cardType)

	// Parse issuing member state (1 byte)
	issuingState := data[1]
	cardNumber.SetCardIssuingMemberState(ddv1.NationNumeric(issuingState))

	// Parse card number based on card type (16 bytes)
	cardNumberData := data[2:18]
	switch cardType {
	case ddv1.EquipmentType_DRIVER_CARD:
		// DriverIdentification is 16 bytes
		driverID, err := opts.UnmarshalDriverIdentification(cardNumberData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse driver identification: %w", err)
		}
		cardNumber.SetDriverIdentification(driverID)
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		// OwnerIdentification is 16 bytes
		ownerID, err := opts.UnmarshalOwnerIdentification(cardNumberData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse owner identification: %w", err)
		}
		cardNumber.SetOwnerIdentification(ownerID)
	default:
		// Card type is recognized by the enum (and already set on cardNumber
		// above) but no structured parser exists for this variant — e.g.
		// CONTROL_CARD (type 5) or MANUFACTURING_CARD. Return the partial
		// FullCardNumber with raw_data preserved instead of failing the
		// whole file parse; downstream callers can decide how to handle it.
		return cardNumber, nil
	}

	return cardNumber, nil
}

// MarshalFullCardNumber marshals full card number data to bytes.
//
// The data type `FullCardNumber` is specified in the Data Dictionary, Section 2.73.
//
// ASN.1 Definition:
//
//	FullCardNumber ::= SEQUENCE {
//	    cardType EquipmentType,
//	    cardIssuingMemberState NationNumeric,
//	    cardNumber CardNumber
//	}
//
//	CardNumber ::= CHOICE {
//	    driverIdentification   SEQUENCE { ... },
//	    ownerIdentification    SEQUENCE { ... }
//	}
//
// Binary Layout (fixed length, 18 bytes):
//   - Card Type (1 byte): EquipmentType
//   - Issuing Member State (1 byte): NationNumeric
//   - Card Number (16 bytes): CardNumber CHOICE based on card type
func (opts MarshalOptions) MarshalFullCardNumber(cardNumber *ddv1.FullCardNumber) ([]byte, error) {
	if cardNumber == nil {
		return nil, fmt.Errorf("cardNumber cannot be nil")
	}

	const lenFullCardNumber = 18

	// Use raw data painting strategy if available
	var canvas [lenFullCardNumber]byte
	if cardNumber.HasRawData() {
		rawData := cardNumber.GetRawData()
		if len(rawData) != lenFullCardNumber {
			return nil, fmt.Errorf("invalid raw_data length for FullCardNumber: got %d, want %d", len(rawData), lenFullCardNumber)
		}
		copy(canvas[:], rawData)

		// Special case: If raw_data starts with 0xFF (no card), return it as-is
		if rawData[0] == 0xFF {
			return canvas[:], nil
		}
	}

	// Check if this is an empty card number (no card inserted)
	// This happens when CardType is UNSPECIFIED and there's no driver/owner identification
	if cardNumber.GetCardType() == ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED &&
		cardNumber.GetDriverIdentification() == nil &&
		cardNumber.GetOwnerIdentification() == nil {
		// Fill with 0xFF to indicate "no card"
		for i := range canvas {
			canvas[i] = 0xFF
		}
		return canvas[:], nil
	}

	// Paint card type (1 byte)
	cardTypeByte, err := MarshalEnum(cardNumber.GetCardType())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card type: %w", err)
	}
	canvas[0] = cardTypeByte

	// Paint issuing member state (1 byte)
	canvas[1] = byte(cardNumber.GetCardIssuingMemberState())

	// Paint card number based on card type (bytes 2-17, 16 bytes total)
	switch cardNumber.GetCardType() {
	case ddv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			// DriverIdentification is 16 bytes
			driverBytes, err := opts.MarshalDriverIdentification(driverID)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal driver identification: %w", err)
			}
			copy(canvas[2:18], driverBytes)
		}
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			// OwnerIdentification is 16 bytes
			ownerBytes, err := opts.MarshalOwnerIdentification(ownerID)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal owner identification: %w", err)
			}
			copy(canvas[2:18], ownerBytes)
		}
	}

	return canvas[:], nil
}

// MarshalFullCardNumberAsString marshals a FullCardNumber structure as a string representation.
// This is used for display purposes and has a maximum length constraint.
func (opts MarshalOptions) MarshalFullCardNumberAsString(cardNumber *ddv1.FullCardNumber, maxLen int) ([]byte, error) {
	if cardNumber == nil {
		return nil, fmt.Errorf("cardNumber cannot be nil")
	}

	// Handle the CardNumber CHOICE based on card type
	switch cardNumber.GetCardType() {
	case ddv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			return opts.MarshalIa5StringValue(driverID.GetDriverIdentificationNumber())
		}
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			return opts.MarshalIa5StringValue(ownerID.GetOwnerIdentification())
		}
	}

	return opts.MarshalStringValue(nil)
}

// AnonymizeFullCardNumber replaces a card number with test values while preserving structure.
func (opts AnonymizeOptions) AnonymizeFullCardNumber(fc *ddv1.FullCardNumber) *ddv1.FullCardNumber {
	if fc == nil {
		return nil
	}

	result := &ddv1.FullCardNumber{}
	// Preserve the card type from the original
	result.SetCardType(fc.GetCardType())
	// Set issuing member state to UNSPECIFIED
	result.SetCardIssuingMemberState(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED)

	// Anonymize driver identification if present
	if driverID := fc.GetDriverIdentification(); driverID != nil {
		anonDriverID := &ddv1.DriverIdentification{}
		anonDriverID.SetDriverIdentificationNumber(opts.AnonymizeIa5StringValue(driverID.GetDriverIdentificationNumber()))
		anonDriverID.SetCardReplacementIndex(opts.AnonymizeIa5StringValue(driverID.GetCardReplacementIndex()))
		anonDriverID.SetCardRenewalIndex(opts.AnonymizeIa5StringValue(driverID.GetCardRenewalIndex()))
		result.SetDriverIdentification(anonDriverID)
	} else if ownerID := fc.GetOwnerIdentification(); ownerID != nil {
		// Anonymize owner identification if present (company cards)
		anonOwnerID := &ddv1.OwnerIdentification{}
		anonOwnerID.SetOwnerIdentification(opts.AnonymizeIa5StringValue(ownerID.GetOwnerIdentification()))
		anonOwnerID.SetConsecutiveIndex(opts.AnonymizeIa5StringValue(ownerID.GetConsecutiveIndex()))
		anonOwnerID.SetReplacementIndex(opts.AnonymizeIa5StringValue(ownerID.GetReplacementIndex()))
		anonOwnerID.SetRenewalIndex(opts.AnonymizeIa5StringValue(ownerID.GetRenewalIndex()))
		result.SetOwnerIdentification(anonOwnerID)
	}

	return result
}
