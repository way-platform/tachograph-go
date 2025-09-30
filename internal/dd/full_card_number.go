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
// Binary Layout (variable length):
//   - Card Type (1 byte): EquipmentType
//   - Issuing Member State (1 byte): NationNumeric
//   - Card Number (variable): CardNumber CHOICE based on card type
func UnmarshalFullCardNumber(data []byte) (*ddv1.FullCardNumber, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for FullCardNumber: got %d, want at least 2", len(data))
	}

	cardNumber := &ddv1.FullCardNumber{}

	// Parse card type (1 byte)
	cardType := data[0]
	cardNumber.SetCardType(ddv1.EquipmentType(cardType))

	// Parse issuing member state (1 byte)
	issuingState := data[1]
	cardNumber.SetCardIssuingMemberState(ddv1.NationNumeric(issuingState))

	// Parse card number based on card type
	remainingData := data[2:]
	switch ddv1.EquipmentType(cardType) {
	case ddv1.EquipmentType_DRIVER_CARD:
		driverID, err := UnmarshalDriverIdentification(remainingData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse driver identification: %w", err)
		}
		cardNumber.SetDriverIdentification(driverID)
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		ownerID, err := UnmarshalOwnerIdentification(remainingData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse owner identification: %w", err)
		}
		cardNumber.SetOwnerIdentification(ownerID)
	default:
		return nil, fmt.Errorf("unsupported card type: %d", cardType)
	}

	return cardNumber, nil
}

// appendFullCardNumber appends full card number data to dst.
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
// Binary Layout (variable length):
//   - Card Type (1 byte): EquipmentType
//   - Issuing Member State (1 byte): NationNumeric
//   - Card Number (variable): CardNumber CHOICE based on card type
func AppendFullCardNumber(dst []byte, cardNumber *ddv1.FullCardNumber) ([]byte, error) {
	if cardNumber == nil {
		return dst, nil
	}

	// Append card type (1 byte)
	dst = append(dst, byte(cardNumber.GetCardType()))

	// Append issuing member state (1 byte)
	dst = append(dst, byte(cardNumber.GetCardIssuingMemberState()))

	// Append card number based on card type
	switch cardNumber.GetCardType() {
	case ddv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			var err error
			dst, err = AppendDriverIdentification(dst, driverID)
			if err != nil {
				return nil, fmt.Errorf("failed to append driver identification: %w", err)
			}
		}
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			var err error
			dst, err = AppendOwnerIdentification(dst, ownerID)
			if err != nil {
				return nil, fmt.Errorf("failed to append owner identification: %w", err)
			}
		}
	}

	return dst, nil
}

// appendFullCardNumberAsString appends a FullCardNumber structure as a string representation.
// This is used for display purposes and has a maximum length constraint.
func AppendFullCardNumberAsString(dst []byte, cardNumber *ddv1.FullCardNumber, maxLen int) ([]byte, error) {
	if cardNumber == nil {
		return AppendString(dst, "", maxLen)
	}

	// Handle the CardNumber CHOICE based on card type
	switch cardNumber.GetCardType() {
	case ddv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			// Concatenate the driver identification components
			identification := driverID.GetDriverIdentificationNumber()

			// Build the full driver identification string
			driverStr := ""
			if identification != nil {
				driverStr += identification.GetValue()
			}
			return AppendString(dst, driverStr, maxLen)
		}
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			identification := ownerID.GetOwnerIdentification()
			if identification != nil {
				return AppendString(dst, identification.GetValue(), maxLen)
			}
		}
	}

	return AppendString(dst, "", maxLen)
}
