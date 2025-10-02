package card

import (
	"errors"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (opts UnmarshalOptions) unmarshalIdentification(data []byte) (*cardv1.Identification, error) {
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
	if nation, err := dd.UnmarshalEnum[ddv1.NationNumeric](data[offset]); err == nil {
		cardId.SetCardIssuingMemberState(nation)
	} else {
		// Value not recognized - set UNRECOGNIZED (no unrecognized field for this type)
		cardId.SetCardIssuingMemberState(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}
	offset++

	// Handle the CardNumber CHOICE type (16 bytes total)
	// CardNumber ::= CHOICE {
	//     -- Driver Card: 14 bytes identification + 1 byte replacement + 1 byte renewal
	//     -- Other Cards: 13 bytes identification + 1 byte consecutive + 1 byte replacement + 1 byte renewal
	// }
	if offset+16 > len(data) {
		return nil, fmt.Errorf("insufficient data for card number")
	}

	cardNumberData := data[offset : offset+16]
	offset += 16

	// Determine card type based on the data structure
	// Driver cards have 14 bytes for identification, other cards have 13 bytes
	// We can detect this by checking if the 14th byte is a valid IA5String character
	// and if the 15th and 16th bytes are single digits (replacement/renewal indices)

	// Try to parse as driver card first (14 + 1 + 1 format)
	driverIdentification, err := opts.UnmarshalIA5StringValue(cardNumberData[0:14])
	if err == nil {
		// Check if bytes 14 and 15 are single digits (0-9)
		replacementByte := cardNumberData[14]
		renewalByte := cardNumberData[15]
		if replacementByte >= '0' && replacementByte <= '9' && renewalByte >= '0' && renewalByte <= '9' {
			// This looks like a driver card format
			driverID := &ddv1.DriverIdentification{}
			driverID.SetDriverIdentificationNumber(driverIdentification)

			// Parse replacement and renewal indices (1 byte each)
			replacementIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[14:15])
			if err == nil {
				driverID.SetCardReplacementIndex(replacementIndex)
			}
			renewalIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[15:16])
			if err == nil {
				driverID.SetCardRenewalIndex(renewalIndex)
			}

			cardId.SetDriverIdentification(driverID)
			identification.SetCardType(cardv1.CardType_DRIVER_CARD)
		} else {
			// Fall back to other card format
			ownerID := &ddv1.OwnerIdentification{}

			// Owner identification (13 bytes)
			ownerIdentification, err := opts.UnmarshalIA5StringValue(cardNumberData[0:13])
			if err != nil {
				return nil, fmt.Errorf("failed to read owner identification: %w", err)
			}
			ownerID.SetOwnerIdentification(ownerIdentification)

			// Consecutive index (1 byte)
			consecutiveIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[13:14])
			if err != nil {
				return nil, fmt.Errorf("failed to read consecutive index: %w", err)
			}
			ownerID.SetConsecutiveIndex(consecutiveIndex)

			// Replacement index (1 byte)
			replacementIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[14:15])
			if err != nil {
				return nil, fmt.Errorf("failed to read replacement index: %w", err)
			}
			ownerID.SetReplacementIndex(replacementIndex)

			// Renewal index (1 byte)
			renewalIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[15:16])
			if err != nil {
				return nil, fmt.Errorf("failed to read renewal index: %w", err)
			}
			ownerID.SetRenewalIndex(renewalIndex)

			cardId.SetOwnerIdentification(ownerID)
			identification.SetCardType(cardv1.CardType_WORKSHOP_CARD) // Default to workshop card
		}
	} else {
		// Try to parse as other card format (13 + 1 + 1 + 1 format)
		ownerID := &ddv1.OwnerIdentification{}

		// Owner identification (13 bytes)
		ownerIdentification, err := opts.UnmarshalIA5StringValue(cardNumberData[0:13])
		if err != nil {
			return nil, fmt.Errorf("failed to read owner identification: %w", err)
		}
		ownerID.SetOwnerIdentification(ownerIdentification)

		// Consecutive index (1 byte)
		consecutiveIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[13:14])
		if err != nil {
			return nil, fmt.Errorf("failed to read consecutive index: %w", err)
		}
		ownerID.SetConsecutiveIndex(consecutiveIndex)

		// Replacement index (1 byte)
		replacementIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[14:15])
		if err != nil {
			return nil, fmt.Errorf("failed to read replacement index: %w", err)
		}
		ownerID.SetReplacementIndex(replacementIndex)

		// Renewal index (1 byte)
		renewalIndex, err := opts.UnmarshalIA5StringValue(cardNumberData[15:16])
		if err != nil {
			return nil, fmt.Errorf("failed to read renewal index: %w", err)
		}
		ownerID.SetRenewalIndex(renewalIndex)

		cardId.SetOwnerIdentification(ownerID)
		identification.SetCardType(cardv1.CardType_WORKSHOP_CARD) // Default to workshop card
	}

	// Authority name (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for card issuing authority name")
	}
	authorityName, err := opts.UnmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read card issuing authority name: %w", err)
	}
	cardId.SetCardIssuingAuthorityName(authorityName)
	offset += 36

	// Card issue date (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card issue date")
	}
	cardIssueDate, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card issue date: %w", err)
	}
	cardId.SetCardIssueDate(cardIssueDate)
	offset += 4

	// Card validity begin (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card validity begin")
	}
	cardValidityBegin, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card validity begin: %w", err)
	}
	cardId.SetCardValidityBegin(cardValidityBegin)
	offset += 4

	// Card expiry date (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card expiry date")
	}
	cardExpiryDate, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card expiry date: %w", err)
	}
	cardId.SetCardExpiryDate(cardExpiryDate)
	offset += 4

	identification.SetCard(cardId)

	// Create and populate DriverCardHolderIdentification part (78 bytes)
	holderId := &cardv1.Identification_DriverCardHolder{}

	// Card holder surname (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder surname")
	}
	surname, err := opts.UnmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder surname: %w", err)
	}
	holderId.SetCardHolderSurname(surname)
	offset += 36

	// Card holder first names (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder first names")
	}
	firstNames, err := opts.UnmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder first names: %w", err)
	}
	holderId.SetCardHolderFirstNames(firstNames)
	offset += 36

	// Card holder birth date (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder birth date")
	}
	birthDate, err := opts.UnmarshalDate(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card holder birth date: %w", err)
	}
	holderId.SetCardHolderBirthDate(birthDate)
	offset += 4

	// Card holder preferred language (2 bytes) - Language ::= IA5String(SIZE(2))
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for card holder preferred language")
	}
	preferredLanguage, err := opts.UnmarshalIA5StringValue(data[offset : offset+2])
	if err != nil {
		return nil, fmt.Errorf("failed to read card holder preferred language: %w", err)
	}
	holderId.SetCardHolderPreferredLanguage(preferredLanguage)
	// offset += 2 // Not needed as this is the last field

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
	// Append cardIssuingMemberState (1 byte) - get protocol value from enum
	memberState := id.GetCardIssuingMemberState()
	var memberStateByte byte
	if memberState == ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED {
		// UNRECOGNIZED values should not occur during marshalling
		return nil, fmt.Errorf("cannot marshal UNRECOGNIZED member state (no unrecognized field)")
	} else {
		var err error
		memberStateByte, err = dd.MarshalEnum(memberState)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal member state: %w", err)
		}
	}
	dst = append(dst, memberStateByte)
	var err error
	// Handle the CardNumber CHOICE type
	// CardNumber ::= CHOICE {
	//     -- Driver Card
	//     SEQUENCE {
	//         driverIdentification    IA5String(SIZE(14)),
	//         cardReplacementIndex    CardReplacementIndex, -- 1 byte
	//         cardRenewalIndex        CardRenewalIndex,     -- 1 byte
	//     },
	//     -- Other Cards (Workshop, Control, Company)
	//     SEQUENCE {
	//         ownerIdentification     IA5String(SIZE(13)),
	//         cardConsecutiveIndex    CardConsecutiveIndex, -- 1 byte
	//         cardReplacementIndex    CardReplacementIndex, -- 1 byte
	//         cardRenewalIndex        CardRenewalIndex      -- 1 byte
	//     }
	// }
	// Total size is always 16 bytes
	cardNumberBytes := make([]byte, 16)
	if driverID := id.GetDriverIdentification(); driverID != nil {
		// Driver card: 14 bytes identification + 1 byte replacement + 1 byte renewal
		identification := driverID.GetDriverIdentificationNumber()
		if identification != nil {
			// Pad or truncate to exactly 14 bytes
			identStr := identification.GetValue()
			if len(identStr) > 14 {
				identStr = identStr[:14]
			}
			copy(cardNumberBytes[0:14], []byte(identStr))
			// Pad with spaces if needed
			for i := len(identStr); i < 14; i++ {
				cardNumberBytes[i] = ' '
			}
		}
		// Write replacement and renewal indices (1 byte each)
		replacementIndex := driverID.GetCardReplacementIndex()
		if replacementIndex != nil && len(replacementIndex.GetValue()) > 0 {
			cardNumberBytes[14] = replacementIndex.GetValue()[0]
		} else {
			cardNumberBytes[14] = '0' // Default replacement index
		}
		renewalIndex := driverID.GetCardRenewalIndex()
		if renewalIndex != nil && len(renewalIndex.GetValue()) > 0 {
			cardNumberBytes[15] = renewalIndex.GetValue()[0]
		} else {
			cardNumberBytes[15] = '0' // Default renewal index
		}
	} else if ownerID := id.GetOwnerIdentification(); ownerID != nil {
		// Other cards: 13 bytes identification + 1 byte consecutive + 1 byte replacement + 1 byte renewal
		identification := ownerID.GetOwnerIdentification()
		if identification != nil {
			// Pad or truncate to exactly 13 bytes
			identStr := identification.GetValue()
			if len(identStr) > 13 {
				identStr = identStr[:13]
			}
			copy(cardNumberBytes[0:13], []byte(identStr))
			// Pad with spaces if needed
			for i := len(identStr); i < 13; i++ {
				cardNumberBytes[i] = ' '
			}
		}
		consecutive := ownerID.GetConsecutiveIndex()
		if consecutive != nil {
			// Convert string to byte value
			consecutiveStr := consecutive.GetValue()
			if len(consecutiveStr) > 0 {
				cardNumberBytes[13] = consecutiveStr[0]
			}
		}
		replacement := ownerID.GetReplacementIndex()
		if replacement != nil {
			// Convert string to byte value
			replacementStr := replacement.GetValue()
			if len(replacementStr) > 0 {
				cardNumberBytes[14] = replacementStr[0]
			}
		}
		renewal := ownerID.GetRenewalIndex()
		if renewal != nil {
			// Convert string to byte value
			renewalStr := renewal.GetValue()
			if len(renewalStr) > 0 {
				cardNumberBytes[15] = renewalStr[0]
			}
		}
	}
	dst = append(dst, cardNumberBytes...)
	dst, err = dd.AppendStringValue(dst, id.GetCardIssuingAuthorityName())
	if err != nil {
		return nil, fmt.Errorf("failed to append card issuing authority name: %w", err)
	}
	dst, err = dd.AppendTimeReal(dst, id.GetCardIssueDate())
	if err != nil {
		return nil, fmt.Errorf("failed to append card issue date: %w", err)
	}
	dst, err = dd.AppendTimeReal(dst, id.GetCardValidityBegin())
	if err != nil {
		return nil, fmt.Errorf("failed to append card validity begin: %w", err)
	}
	dst, err = dd.AppendTimeReal(dst, id.GetCardExpiryDate())
	if err != nil {
		return nil, fmt.Errorf("failed to append card expiry date: %w", err)
	}
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
	nameBlock, err = dd.AppendStringValue(nameBlock, h.GetCardHolderSurname())
	if err != nil {
		return nil, fmt.Errorf("failed to append card holder surname: %w", err)
	}
	nameBlock, err = dd.AppendStringValue(nameBlock, h.GetCardHolderFirstNames())
	if err != nil {
		return nil, fmt.Errorf("failed to append card holder first names: %w", err)
	}
	dst = append(dst, nameBlock...)

	birthDate := h.GetCardHolderBirthDate()
	if birthDate != nil {
		dst, err = dd.AppendDate(dst, birthDate)
		if err != nil {
			return nil, fmt.Errorf("failed to append birth date: %w", err)
		}
	} else {
		// Append default date (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	dst, err = dd.AppendStringValue(dst, h.GetCardHolderPreferredLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to append preferred language: %w", err)
	}
	return dst, nil
}

// AnonymizeIdentification creates an anonymized copy of Identification, replacing all
// personally identifiable information with safe, deterministic test values while
// preserving the structure and validity for testing.
//
// Anonymization strategy:
// - Names: Replaced with generic test names
// - Card numbers: Replaced with test values
// - Addresses: Replaced with generic test addresses
// - Birth dates: Replaced with static test date (2000-01-01)
// - Card dates: Replaced with static test dates (issue/validity: 2020-01-01, expiry: 2024-12-31)
// - Countries: Preserved (structural info)
// - Signatures: Cleared (will be invalid after anonymization anyway)
func AnonymizeIdentification(id *cardv1.Identification) *cardv1.Identification {
	if id == nil {
		return nil
	}

	result := &cardv1.Identification{}
	result.SetCardType(id.GetCardType())

	// Anonymize card identification
	if card := id.GetCard(); card != nil {
		anonymizedCard := &cardv1.Identification_Card{}

		// Preserve country (structural info)
		anonymizedCard.SetCardIssuingMemberState(card.GetCardIssuingMemberState())

		// Anonymize driver identification
		if driverID := card.GetDriverIdentification(); driverID != nil {
			anonymizedCard.SetDriverIdentification(dd.AnonymizeDriverIdentification(driverID))
		}

		// Anonymize owner identification (for workshop/control/company cards)
		if ownerID := card.GetOwnerIdentification(); ownerID != nil {
			anonymizedOwner := &ddv1.OwnerIdentification{}
			// Use generic owner ID
			anonymizedOwner.SetOwnerIdentification(
				dd.AnonymizeStringValue(ownerID.GetOwnerIdentification(), "OWNER00000001"),
			)
			// Preserve indices (structural info)
			anonymizedOwner.SetConsecutiveIndex(ownerID.GetConsecutiveIndex())
			anonymizedOwner.SetReplacementIndex(ownerID.GetReplacementIndex())
			anonymizedOwner.SetRenewalIndex(ownerID.GetRenewalIndex())
			anonymizedCard.SetOwnerIdentification(anonymizedOwner)
		}

		// Anonymize issuing authority name
		anonymizedCard.SetCardIssuingAuthorityName(
			dd.AnonymizeStringValue(card.GetCardIssuingAuthorityName(), "TEST_AUTHORITY"),
		)

		// Replace card dates with static test dates (valid 5-year period)
		// Issue/validity: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
		// Expiry: 2024-12-31 23:59:59 UTC (epoch: 1735689599)
		anonymizedCard.SetCardIssueDate(&timestamppb.Timestamp{Seconds: 1577836800})
		anonymizedCard.SetCardValidityBegin(&timestamppb.Timestamp{Seconds: 1577836800})
		anonymizedCard.SetCardExpiryDate(&timestamppb.Timestamp{Seconds: 1735689599})

		result.SetCard(anonymizedCard)
	}

	// Anonymize holder identification based on card type
	switch id.GetCardType() {
	case cardv1.CardType_DRIVER_CARD:
		if holder := id.GetDriverCardHolder(); holder != nil {
			anonymizedHolder := &cardv1.Identification_DriverCardHolder{}

			// Replace names with test values
			anonymizedHolder.SetCardHolderSurname(
				dd.AnonymizeStringValue(holder.GetCardHolderSurname(), "TEST_SURNAME"),
			)
			anonymizedHolder.SetCardHolderFirstNames(
				dd.AnonymizeStringValue(holder.GetCardHolderFirstNames(), "TEST_FIRSTNAME"),
			)

			// Replace birth date with static test date (2000-01-01)
			birthDate := &ddv1.Date{}
			birthDate.SetYear(2000)
			birthDate.SetMonth(1)
			birthDate.SetDay(1)
			// Regenerate raw_data for binary fidelity
			if rawData, err := dd.AppendDate(nil, birthDate); err == nil {
				birthDate.SetRawData(rawData)
			}
			anonymizedHolder.SetCardHolderBirthDate(birthDate)

			// Preserve language (not sensitive)
			anonymizedHolder.SetCardHolderPreferredLanguage(holder.GetCardHolderPreferredLanguage())

			result.SetDriverCardHolder(anonymizedHolder)
		}

	case cardv1.CardType_WORKSHOP_CARD:
		if holder := id.GetWorkshopCardHolder(); holder != nil {
			anonymizedHolder := &cardv1.Identification_WorkshopCardHolder{}

			// Anonymize workshop details
			anonymizedHolder.SetWorkshopName(
				dd.AnonymizeStringValue(holder.GetWorkshopName(), "TEST_WORKSHOP"),
			)
			anonymizedHolder.SetWorkshopAddress(
				dd.AnonymizeStringValue(holder.GetWorkshopAddress(), "TEST_ADDRESS"),
			)

			// Anonymize holder names
			anonymizedHolder.SetCardHolderSurname(
				dd.AnonymizeStringValue(holder.GetCardHolderSurname(), "TEST_SURNAME"),
			)
			anonymizedHolder.SetCardHolderFirstNames(
				dd.AnonymizeStringValue(holder.GetCardHolderFirstNames(), "TEST_FIRSTNAME"),
			)

			// Preserve language
			anonymizedHolder.SetCardHolderPreferredLanguage(holder.GetCardHolderPreferredLanguage())

			result.SetWorkshopCardHolder(anonymizedHolder)
		}

	case cardv1.CardType_CONTROL_CARD:
		if holder := id.GetControlCardHolder(); holder != nil {
			anonymizedHolder := &cardv1.Identification_ControlCardHolder{}

			// Anonymize control body details
			anonymizedHolder.SetControlBodyName(
				dd.AnonymizeStringValue(holder.GetControlBodyName(), "TEST_CONTROL_BODY"),
			)
			anonymizedHolder.SetControlBodyAddress(
				dd.AnonymizeStringValue(holder.GetControlBodyAddress(), "TEST_ADDRESS"),
			)

			// Anonymize holder names
			anonymizedHolder.SetCardHolderSurname(
				dd.AnonymizeStringValue(holder.GetCardHolderSurname(), "TEST_SURNAME"),
			)
			anonymizedHolder.SetCardHolderFirstNames(
				dd.AnonymizeStringValue(holder.GetCardHolderFirstNames(), "TEST_FIRSTNAME"),
			)

			// Preserve language
			anonymizedHolder.SetCardHolderPreferredLanguage(holder.GetCardHolderPreferredLanguage())

			result.SetControlCardHolder(anonymizedHolder)
		}

	case cardv1.CardType_COMPANY_CARD:
		if holder := id.GetCompanyCardHolder(); holder != nil {
			anonymizedHolder := &cardv1.Identification_CompanyCardHolder{}

			// Anonymize company details
			anonymizedHolder.SetCompanyName(
				dd.AnonymizeStringValue(holder.GetCompanyName(), "TEST_COMPANY"),
			)
			anonymizedHolder.SetCompanyAddress(
				dd.AnonymizeStringValue(holder.GetCompanyAddress(), "TEST_ADDRESS"),
			)

			// Preserve language
			anonymizedHolder.SetCardHolderPreferredLanguage(holder.GetCardHolderPreferredLanguage())

			result.SetCompanyCardHolder(anonymizedHolder)
		}
	}

	// Don't preserve signature - it will be invalid after anonymization
	// Users can re-sign if needed for testing

	return result
}
