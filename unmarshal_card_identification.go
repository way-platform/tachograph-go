package tachograph

import (
	"bytes"
	"errors"

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
	cardId.SetCardIssuingMemberState(readString(r, 1))
	cardId.SetCardNumber(readString(r, 16))
	cardId.SetCardIssuingAuthorityName(readString(r, 36))
	cardId.SetCardIssueDate(readTimeReal(r))
	cardId.SetCardValidityBegin(readTimeReal(r))
	cardId.SetCardExpiryDate(readTimeReal(r))
	identification.SetCard(cardId)

	// Set card type to DRIVER_CARD for driver cards
	identification.SetCardType(cardv1.CardType_DRIVER_CARD)

	// Create and populate DriverCardHolderIdentification part (78 bytes)
	holderId := &cardv1.Identification_DriverCardHolder{}
	holderId.SetCardHolderSurname(readString(r, 36))
	holderId.SetCardHolderFirstNames(readString(r, 36))
	birthDate, err := readDatef(r)
	if err != nil {
		return nil, err
	}
	holderId.SetCardHolderBirthDate(birthDate)
	holderId.SetCardHolderPreferredLanguage(readString(r, 2))
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
