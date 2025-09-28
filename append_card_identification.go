package tachograph

import (
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardIdentification appends the binary representation of CardIdentification to dst.
//
// ASN.1 Specification (Data Dictionary 2.1):
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
// Binary Layout (65 bytes):
//
//	0-0:   cardIssuingMemberState (1 byte, BCD)
//	1-16:  cardNumber (16 bytes, padded)
//	17-52: cardIssuingAuthorityName (36 bytes, padded)
//	53-60: cardIssueDate (8 bytes)
//	61-68: cardValidityBegin (8 bytes)
//	69-76: cardExpiryDate (8 bytes)
func AppendCardIdentification(dst []byte, id *cardv1.Identification_Card) ([]byte, error) {
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
// ASN.1 Specification (Data Dictionary 2.1):
//
//	DriverCardHolderIdentification ::= SEQUENCE {
//	    cardHolderSurname            Name,
//	    cardHolderFirstNames         Name,
//	    cardHolderBirthDate          Datef,
//	    cardHolderPreferredLanguage  Language
//	}
//
// Binary Layout (78 bytes):
//
//	0-35:  cardHolderSurname (36 bytes, padded)
//	36-71: cardHolderFirstNames (36 bytes, padded)
//	72-75: cardHolderBirthDate (4 bytes)
//	76-77: cardHolderPreferredLanguage (2 bytes, padded)
func AppendDriverCardHolderIdentification(dst []byte, h *cardv1.Identification_DriverCardHolder) ([]byte, error) {
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
