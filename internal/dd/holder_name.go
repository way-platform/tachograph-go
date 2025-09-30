package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalHolderName parses holder name data.
//
// The data type `HolderName` is specified in the Data Dictionary, Section 2.83.
//
// ASN.1 Definition:
//
//	HolderName ::= SEQUENCE {
//	    holderSurname Name,
//	    holderFirstNames Name
//	}
//
//	Name ::= SEQUENCE {
//	    codePage INTEGER,
//	    name OCTET STRING
//	}
//
// Binary Layout (variable length):
//   - Holder Surname (variable): Name structure
//   - Holder First Names (variable): Name structure
func UnmarshalHolderName(data []byte) (*ddv1.HolderName, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for HolderName: got %d, want at least 2", len(data))
	}

	holderName := &ddv1.HolderName{}

	// Parse holder surname (Name structure)
	// Name structure: 1 byte codePage + variable length name
	nameLength := int(data[1])

	if len(data) < 2+nameLength {
		return nil, fmt.Errorf("insufficient data for holder surname: got %d, want %d", len(data), 2+nameLength)
	}

	surnameData := data[0 : 2+nameLength] // Include codePage in the data
	surname, err := UnmarshalStringValue(surnameData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse holder surname: %w", err)
	}
	holderName.SetHolderSurname(surname)

	// Parse holder first names (Name structure)
	remainingData := data[2+nameLength:]
	if len(remainingData) < 2 {
		return nil, fmt.Errorf("insufficient data for holder first names: got %d, want at least 2", len(remainingData))
	}

	firstNamesData := remainingData // Include codePage in the data
	firstNames, err := UnmarshalStringValue(firstNamesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse holder first names: %w", err)
	}
	holderName.SetHolderFirstNames(firstNames)

	return holderName, nil
}

// appendHolderName appends holder name data to dst.
//
// The data type `HolderName` is specified in the Data Dictionary, Section 2.83.
//
// ASN.1 Definition:
//
//	HolderName ::= SEQUENCE {
//	    holderSurname Name,
//	    holderFirstNames Name
//	}
//
//	Name ::= SEQUENCE {
//	    codePage INTEGER,
//	    name OCTET STRING
//	}
//
// Binary Layout (variable length):
//   - Holder Surname (variable): Name structure
//   - Holder First Names (variable): Name structure
func AppendHolderName(dst []byte, holderName *ddv1.HolderName) ([]byte, error) {
	if holderName == nil {
		// Append default values: empty name structures
		dst = append(dst, 0, 0) // Empty surname
		dst = append(dst, 0, 0) // Empty first names
		return dst, nil
	}

	// Append holder surname (Name structure)
	surname := holderName.GetHolderSurname()
	if surname != nil {
		dst = append(dst, byte(surname.GetEncoding()))
		nameBytes := []byte(surname.GetDecoded())
		dst = append(dst, byte(len(nameBytes)))
		dst = append(dst, nameBytes...)
	} else {
		dst = append(dst, 0, 0) // Empty surname
	}

	// Append holder first names (Name structure)
	firstNames := holderName.GetHolderFirstNames()
	if firstNames != nil {
		dst = append(dst, byte(firstNames.GetEncoding()))
		nameBytes := []byte(firstNames.GetDecoded())
		dst = append(dst, byte(len(nameBytes)))
		dst = append(dst, nameBytes...)
	} else {
		dst = append(dst, 0, 0) // Empty first names
	}

	return dst, nil
}
