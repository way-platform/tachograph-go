package tachograph

import (
	"bytes"
	"errors"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalIdentification parses the binary data for an EF_Identification record.
func UnmarshalIdentification(data []byte, id *cardv1.CardIdentification, hid *cardv1.DriverCardHolderIdentification) error {
	if len(data) < 143 {
		return errors.New("not enough data for EF_Identification")
	}
	r := bytes.NewReader(data)

	// Parse CardIdentification part (65 bytes)
	id.SetCardIssuingMemberState(readString(r, 1))
	id.SetCardNumber(readString(r, 16))
	id.SetCardIssuingAuthorityName(readString(r, 36))
	id.SetCardIssueDate(readTimeReal(r))
	id.SetCardValidityBegin(readTimeReal(r))
	id.SetCardExpiryDate(readTimeReal(r))

	// Parse DriverCardHolderIdentification part (78 bytes)
	hid.SetCardHolderSurname(readString(r, 36))
	hid.SetCardHolderFirstNames(readString(r, 36))
	hid.SetCardHolderBirthDate(readDatef(r))
	hid.SetCardHolderPreferredLanguage(readString(r, 2))

	return nil
}
