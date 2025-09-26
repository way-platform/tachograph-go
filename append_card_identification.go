package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardIdentification appends the binary representation of CardIdentification to dst.
func AppendCardIdentification(dst []byte, id *cardv1.CardIdentification) ([]byte, error) {
	if id == nil {
		return dst, nil
	}
	var err error
	dst, err = appendBCDNation(dst, id.GetCardIssuingMemberState()) // NationNumeric is BCD-encoded
	if err != nil {
		return nil, err
	}
	dst, err = appendString(dst, id.GetCardNumber(), 16)
	if err != nil {
		return nil, err
	}
	dst, err = appendString(dst, id.GetCardIssuingAuthorityName(), 36)
	if err != nil {
		return nil, err
	}
	dst = appendTimeReal(dst, id.GetCardIssueDate())
	dst = appendTimeReal(dst, id.GetCardValidityBegin())
	dst = appendTimeReal(dst, id.GetCardExpiryDate())
	return dst, nil
}

// AppendDriverCardHolderIdentification appends the binary representation of DriverCardHolderIdentification to dst.
func AppendDriverCardHolderIdentification(dst []byte, h *cardv1.DriverCardHolderIdentification) ([]byte, error) {
	if h == nil {
		return dst, nil
	}
	var err error
	nameBlock := make([]byte, 0, 72)
	nameBlock, err = appendString(nameBlock, h.GetCardHolderSurname(), 36)
	if err != nil {
		return nil, err
	}
	nameBlock, err = appendString(nameBlock, h.GetCardHolderFirstNames(), 36)
	if err != nil {
		return nil, err
	}
	dst = append(dst, nameBlock...)

	dst = appendDatef(dst, h.GetCardHolderBirthDate())

	dst, err = appendString(dst, h.GetCardHolderPreferredLanguage(), 2)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
