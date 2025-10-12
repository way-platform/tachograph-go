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
//	    codePage INTEGER (0..255),    -- 1 byte
//	    name OCTET STRING (SIZE (35))  -- 35 bytes (fixed, no length prefix with PER)
//	}
//
// Binary Layout (fixed length, 72 bytes):
//   - Bytes 0-35: Holder Surname (1 byte codePage + 35 bytes name)
//   - Bytes 36-71: Holder First Names (1 byte codePage + 35 bytes name)
func (opts UnmarshalOptions) UnmarshalHolderName(data []byte) (*ddv1.HolderName, error) {
	const (
		lenName       = 36 // 1 byte codePage + 35 bytes name
		lenHolderName = 72 // 2 * lenName
	)

	if len(data) != lenHolderName {
		return nil, fmt.Errorf("invalid data length for HolderName: got %d, want %d", len(data), lenHolderName)
	}

	holderName := &ddv1.HolderName{}

	// Parse holder surname (first 36 bytes)
	surnameData := data[0:lenName]
	surname, err := opts.UnmarshalStringValue(surnameData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse holder surname: %w", err)
	}
	holderName.SetHolderSurname(surname)

	// Parse holder first names (second 36 bytes)
	firstNamesData := data[lenName:lenHolderName]
	firstNames, err := opts.UnmarshalStringValue(firstNamesData)
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
//	    codePage INTEGER (0..255),    -- 1 byte
//	    name OCTET STRING (SIZE (35))  -- 35 bytes (fixed, no length prefix with PER)
//	}
//
// Binary Layout (fixed length, 72 bytes):
//   - Bytes 0-35: Holder Surname (1 byte codePage + 35 bytes name)
//   - Bytes 36-71: Holder First Names (1 byte codePage + 35 bytes name)
func AppendHolderName(dst []byte, holderName *ddv1.HolderName) ([]byte, error) {
	if holderName == nil {
		return nil, fmt.Errorf("holderName cannot be nil")
	}

	var err error

	// Append holder surname (36 bytes)
	dst, err = AppendStringValue(dst, holderName.GetHolderSurname())
	if err != nil {
		return nil, fmt.Errorf("failed to append holder surname: %w", err)
	}

	// Append holder first names (36 bytes)
	dst, err = AppendStringValue(dst, holderName.GetHolderFirstNames())
	if err != nil {
		return nil, fmt.Errorf("failed to append holder first names: %w", err)
	}

	return dst, nil
}
