package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalOwnerIdentification parses owner identification data.
//
// The data type `OwnerIdentification` is specified in the Data Dictionary, Section 2.26.
//
// ASN.1 Definition:
//
//	ownerIdentification SEQUENCE {
//	    ownerIdentificationNumber IA5String(SIZE(13)),
//	    cardConsecutiveIndex CardConsecutiveIndex,
//	    cardReplacementIndex CardReplacementIndex,
//	    cardRenewalIndex CardRenewalIndex
//	}
//
// Binary Layout (16 bytes):
//   - Owner Identification Number (13 bytes): IA5String
//   - Card Consecutive Index (1 byte): IA5String
//   - Card Replacement Index (1 byte): IA5String
//   - Card Renewal Index (1 byte): IA5String
func UnmarshalOwnerIdentification(data []byte) (*ddv1.OwnerIdentification, error) {
	const (
		lenOwnerIdentification = 16 // 13 + 1 + 1 + 1
	)

	if len(data) != lenOwnerIdentification {
		return nil, fmt.Errorf("invalid data length for OwnerIdentification: got %d, want %d", len(data), lenOwnerIdentification)
	}

	ownerID := &ddv1.OwnerIdentification{}

	// Parse owner identification number (13 bytes)
	identificationNumber, err := UnmarshalIA5StringValue(data[0:13])
	if err != nil {
		return nil, fmt.Errorf("failed to parse owner identification number: %w", err)
	}
	ownerID.SetOwnerIdentification(identificationNumber)

	// Parse card consecutive index (1 byte)
	consecutiveIndex, err := UnmarshalIA5StringValue(data[13:14])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card consecutive index: %w", err)
	}
	ownerID.SetConsecutiveIndex(consecutiveIndex)

	// Parse card replacement index (1 byte)
	replacementIndex, err := UnmarshalIA5StringValue(data[14:15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card replacement index: %w", err)
	}
	ownerID.SetReplacementIndex(replacementIndex)

	// Parse card renewal index (1 byte)
	renewalIndex, err := UnmarshalIA5StringValue(data[15:16])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card renewal index: %w", err)
	}
	ownerID.SetRenewalIndex(renewalIndex)

	return ownerID, nil
}

// appendOwnerIdentification appends owner identification data to dst.
//
// The data type `OwnerIdentification` is specified in the Data Dictionary, Section 2.26.
//
// ASN.1 Definition:
//
//	ownerIdentification SEQUENCE {
//	    ownerIdentificationNumber IA5String(SIZE(13)),
//	    cardConsecutiveIndex CardConsecutiveIndex,
//	    cardReplacementIndex CardReplacementIndex,
//	    cardRenewalIndex CardRenewalIndex
//	}
//
// Binary Layout (16 bytes):
//   - Owner Identification Number (13 bytes): IA5String
//   - Card Consecutive Index (1 byte): IA5String
//   - Card Replacement Index (1 byte): IA5String
//   - Card Renewal Index (1 byte): IA5String
func AppendOwnerIdentification(dst []byte, ownerID *ddv1.OwnerIdentification) ([]byte, error) {
	if ownerID == nil {
		return nil, fmt.Errorf("ownerID cannot be nil")
	}

	// Append owner identification number (13 bytes)
	var err error
	dst, err = AppendStringValue(dst, ownerID.GetOwnerIdentification())
	if err != nil {
		return nil, fmt.Errorf("failed to append owner identification number: %w", err)
	}

	// Append card consecutive index (1 byte)
	dst, err = AppendStringValue(dst, ownerID.GetConsecutiveIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to append card consecutive index: %w", err)
	}

	// Append card replacement index (1 byte)
	dst, err = AppendStringValue(dst, ownerID.GetReplacementIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to append card replacement index: %w", err)
	}

	// Append card renewal index (1 byte)
	dst, err = AppendStringValue(dst, ownerID.GetRenewalIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to append card renewal index: %w", err)
	}

	return dst, nil
}
