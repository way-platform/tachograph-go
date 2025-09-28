package tachograph

import (
	"bytes"
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalIdentification parses the binary data for an EF_Identification record.
func unmarshalIdentification(data []byte) (*cardv1.Identification, error) {
	if len(data) < 143 {
		return nil, errors.New("not enough data for EF_Identification")
	}

	var identification cardv1.Identification
	r := bytes.NewReader(data)

	// Create and populate CardIdentification part (65 bytes)
	cardId := &cardv1.Identification_Card{}
	// Read nation as byte and convert to NationNumeric
	nation, err := unmarshalNationNumericFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read card issuing member state: %w", err)
	}
	cardId.SetCardIssuingMemberState(nation)
	// Handle the inlined CardNumber structure
	// For now, assume this is a driver card and create driver identification
	driverID := &cardv1.Identification_DriverIdentification{}

	cardNumber, err := unmarshalIA5StringValueFromReader(r, 14)
	if err != nil {
		return nil, fmt.Errorf("failed to read card number: %w", err)
	}
	driverID.SetIdentification(cardNumber)

	consecutiveIndex, err := unmarshalIA5StringValueFromReader(r, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to read consecutive index: %w", err)
	}
	driverID.SetConsecutiveIndex(consecutiveIndex)

	replacementIndex, err := unmarshalIA5StringValueFromReader(r, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to read replacement index: %w", err)
	}
	driverID.SetReplacementIndex(replacementIndex)

	renewalIndex, err := unmarshalIA5StringValueFromReader(r, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to read renewal index: %w", err)
	}
	driverID.SetRenewalIndex(renewalIndex)

	cardId.SetDriverIdentification(driverID)

	authorityName, err := unmarshalStringValueFromReader(r, 36)
	if err != nil {
		return nil, fmt.Errorf("failed to read card issuing authority name: %w", err)
	}
	cardId.SetCardIssuingAuthorityName(authorityName)
	cardId.SetCardIssueDate(readTimeReal(r))
	cardId.SetCardValidityBegin(readTimeReal(r))
	cardId.SetCardExpiryDate(readTimeReal(r))
	identification.SetCard(cardId)

	// Set card type to DRIVER_CARD for driver cards
	identification.SetCardType(cardv1.CardType_DRIVER_CARD)

	// Create and populate DriverCardHolderIdentification part (78 bytes)
	holderId := &cardv1.Identification_DriverCardHolder{}

	surname, err := unmarshalStringValueFromReader(r, 36)
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder surname: %w", err)
	}
	holderId.SetCardHolderSurname(surname)

	firstNames, err := unmarshalStringValueFromReader(r, 36)
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder first names: %w", err)
	}
	holderId.SetCardHolderFirstNames(firstNames)

	birthDate, err := readDatef(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder birth date: %w", err)
	}
	holderId.SetCardHolderBirthDate(birthDate)

	preferredLanguage, err := unmarshalIA5StringValueFromReader(r, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder preferred language: %w", err)
	}
	holderId.SetCardHolderPreferredLanguage(preferredLanguage)
	identification.SetDriverCardHolder(holderId)

	return &identification, nil
}

// UnmarshalIdentification parses the binary data for an EF_Identification record (legacy function).
// Deprecated: Use unmarshalIdentification instead.
func UnmarshalIdentification(data []byte, id *cardv1.Identification_Card, hid *cardv1.Identification_DriverCardHolder) error {
	identification, err := unmarshalIdentification(data)
	if err != nil {
		return err
	}

	// Extract the separate components for backward compatibility
	if identification.GetCard() != nil {
		*id = *identification.GetCard()
	}
	if identification.GetDriverCardHolder() != nil {
		*hid = *identification.GetDriverCardHolder()
	}
	return nil
}
