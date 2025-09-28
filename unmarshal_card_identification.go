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
	offset := 0

	// Create and populate CardIdentification part (65 bytes)
	cardId := &cardv1.Identification_Card{}

	// Read nation as byte and convert to NationNumeric
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for card issuing member state")
	}
	nation, err := unmarshalNationNumeric(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read card issuing member state: %w", err)
	}
	cardId.SetCardIssuingMemberState(nation)
	offset++

	// Handle the inlined CardNumber structure
	// For now, assume this is a driver card and create driver identification
	driverID := &cardv1.Identification_DriverIdentification{}

	// Card number (14 bytes)
	if offset+14 > len(data) {
		return nil, fmt.Errorf("insufficient data for card number")
	}
	cardNumber, err := unmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return nil, fmt.Errorf("failed to read card number: %w", err)
	}
	driverID.SetIdentification(cardNumber)
	offset += 14

	// Consecutive index (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for consecutive index")
	}
	consecutiveIndex, err := unmarshalIA5StringValue(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read consecutive index: %w", err)
	}
	driverID.SetConsecutiveIndex(consecutiveIndex)
	offset++

	// Replacement index (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for replacement index")
	}
	replacementIndex, err := unmarshalIA5StringValue(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read replacement index: %w", err)
	}
	driverID.SetReplacementIndex(replacementIndex)
	offset++

	// Renewal index (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for renewal index")
	}
	renewalIndex, err := unmarshalIA5StringValue(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read renewal index: %w", err)
	}
	driverID.SetRenewalIndex(renewalIndex)
	offset++

	cardId.SetDriverIdentification(driverID)

	// Authority name (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for card issuing authority name")
	}
	authorityName, err := unmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read card issuing authority name: %w", err)
	}
	cardId.SetCardIssuingAuthorityName(authorityName)
	offset += 36

	// Card issue date (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card issue date")
	}
	cardId.SetCardIssueDate(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Card validity begin (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card validity begin")
	}
	cardId.SetCardValidityBegin(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	// Card expiry date (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card expiry date")
	}
	cardId.SetCardExpiryDate(readTimeReal(bytes.NewReader(data[offset : offset+4])))
	offset += 4

	identification.SetCard(cardId)

	// Set card type to DRIVER_CARD for driver cards
	identification.SetCardType(cardv1.CardType_DRIVER_CARD)

	// Create and populate DriverCardHolderIdentification part (78 bytes)
	holderId := &cardv1.Identification_DriverCardHolder{}

	// Card holder surname (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder surname")
	}
	surname, err := unmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder surname: %w", err)
	}
	holderId.SetCardHolderSurname(surname)
	offset += 36

	// Card holder first names (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder first names")
	}
	firstNames, err := unmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder first names: %w", err)
	}
	holderId.SetCardHolderFirstNames(firstNames)
	offset += 36

	// Card holder birth date (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder birth date")
	}
	birthDate, err := readDatef(bytes.NewReader(data[offset : offset+4]))
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder birth date: %w", err)
	}
	holderId.SetCardHolderBirthDate(birthDate)
	offset += 4

	// Card holder preferred language (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder preferred language")
	}
	preferredLanguage, err := unmarshalIA5StringValue(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder preferred language: %w", err)
	}
	holderId.SetCardHolderPreferredLanguage(preferredLanguage)
	offset++

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
