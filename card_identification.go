package tachograph

import (
	"bytes"
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalIdentification parses the binary data for an EF_Identification record.
//
// The data type `CardIdentification` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	CardIdentification ::= SEQUENCE {
//	    cardIssuingMemberState    NationNumeric,
//	    cardNumber                CardNumber,
//	    cardIssuingAuthorityName  Name,
//	    cardIssueDate            TimeReal,
//	    cardValidityBegin        TimeReal,
//	    cardExpiryDate           TimeReal
//	}
//
//	DriverCardHolderIdentification ::= SEQUENCE {
//	    cardHolderSurname            Name,
//	    cardHolderFirstNames         Name,
//	    cardHolderBirthDate          Datef,
//	    cardHolderPreferredLanguage  Language
//	}
func unmarshalIdentification(data []byte) (*cardv1.Identification, error) {
	const (
		lenMinIdentification = 143 // Minimum size for EF_Identification
	)

	if len(data) < lenMinIdentification {
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
	driverID := &ddv1.DriverIdentification{}

	// Card number (14 bytes)
	if offset+14 > len(data) {
		return nil, fmt.Errorf("insufficient data for card number")
	}
	cardNumber, err := unmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return nil, fmt.Errorf("failed to read card number: %w", err)
	}
	driverID.SetIdentificationNumber(cardNumber)
	offset += 14

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
	// offset++ // Not needed as this is the last field

	identification.SetDriverCardHolder(holderId)

	return &identification, nil
}

// AppendCardIdentification appends the binary representation of CardIdentification to dst.
//
// The data type `CardIdentification` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	CardIdentification ::= SEQUENCE {
//	    cardIssuingMemberState    NationNumeric,
//	    cardNumber                CardNumber,
//	    cardIssuingAuthorityName  Name,
//	    cardIssueDate            TimeReal,
//	    cardValidityBegin        TimeReal,
//	    cardExpiryDate           TimeReal
//	}
func appendCardIdentification(dst []byte, id *cardv1.Identification_Card) ([]byte, error) {
	if id == nil {
		return dst, nil
	}
	var err error
	dst, err = appendBCDNation(dst, fmt.Sprintf("%d", int32(id.GetCardIssuingMemberState()))) // NationNumeric is BCD-encoded
	if err != nil {
		return nil, err
	}
	// Handle the CardNumber CHOICE - for now, we'll use a placeholder
	// This needs to be implemented based on the actual card type
	cardNumberStr := ""
	if driverID := id.GetDriverIdentification(); driverID != nil {
		// Build driver identification string
		if identification := driverID.GetIdentificationNumber(); identification != nil {
			cardNumberStr += identification.GetDecoded()
		}
	} else if ownerID := id.GetOwnerIdentification(); ownerID != nil {
		if identification := ownerID.GetIdentificationNumber(); identification != nil {
			cardNumberStr = identification.GetDecoded()
		}
		if consecutive := ownerID.GetConsecutiveIndex(); consecutive != nil {
			cardNumberStr += consecutive.GetDecoded()
		}
		if replacement := ownerID.GetReplacementIndex(); replacement != nil {
			cardNumberStr += replacement.GetDecoded()
		}
		if renewal := ownerID.GetRenewalIndex(); renewal != nil {
			cardNumberStr += renewal.GetDecoded()
		}
	}
	dst, err = appendString(dst, cardNumberStr, 16)
	if err != nil {
		return nil, err
	}
	authorityName := id.GetCardIssuingAuthorityName()
	if authorityName != nil {
		dst, err = appendString(dst, authorityName.GetDecoded(), 36)
	} else {
		dst, err = appendString(dst, "", 36)
	}
	if err != nil {
		return nil, err
	}
	dst = appendTimeReal(dst, id.GetCardIssueDate())
	dst = appendTimeReal(dst, id.GetCardValidityBegin())
	dst = appendTimeReal(dst, id.GetCardExpiryDate())
	return dst, nil
}

// AppendDriverCardHolderIdentification appends the binary representation of DriverCardHolderIdentification to dst.
//
// The data type `DriverCardHolderIdentification` is specified in the Data Dictionary, Section 2.1.
//
// ASN.1 Definition:
//
//	DriverCardHolderIdentification ::= SEQUENCE {
//	    cardHolderSurname            Name,
//	    cardHolderFirstNames         Name,
//	    cardHolderBirthDate          Datef,
//	    cardHolderPreferredLanguage  Language
//	}
func appendDriverCardHolderIdentification(dst []byte, h *cardv1.Identification_DriverCardHolder) ([]byte, error) {
	if h == nil {
		return dst, nil
	}
	var err error
	nameBlock := make([]byte, 0, 72)
	surname := h.GetCardHolderSurname()
	if surname != nil {
		nameBlock, err = appendString(nameBlock, surname.GetDecoded(), 36)
	} else {
		nameBlock, err = appendString(nameBlock, "", 36)
	}
	if err != nil {
		return nil, err
	}
	firstNames := h.GetCardHolderFirstNames()
	if firstNames != nil {
		nameBlock, err = appendString(nameBlock, firstNames.GetDecoded(), 36)
	} else {
		nameBlock, err = appendString(nameBlock, "", 36)
	}
	if err != nil {
		return nil, err
	}
	dst = append(dst, nameBlock...)

	birthDate := h.GetCardHolderBirthDate()
	if birthDate != nil {
		dst = appendDate(dst, birthDate)
	} else {
		// Append default date (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	preferredLanguage := h.GetCardHolderPreferredLanguage()
	if preferredLanguage != nil {
		dst, err = appendString(dst, preferredLanguage.GetDecoded(), 2)
	} else {
		dst, err = appendString(dst, "", 2)
	}
	if err != nil {
		return nil, err
	}
	return dst, nil
}
