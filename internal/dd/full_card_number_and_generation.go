package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalFullCardNumberAndGeneration parses full card number and generation data.
//
// The data type `FullCardNumberAndGeneration` is specified in the Data Dictionary, Section 2.74.
//
// ASN.1 Definition:
//
//	FullCardNumberAndGeneration ::= SEQUENCE {
//	    fullcardNumber FullCardNumber,
//	    generation Generation
//	}
//
// Binary Layout (variable length):
//   - Full Card Number (variable): FullCardNumber structure
//   - Generation (1 byte): Generation enum value
func UnmarshalFullCardNumberAndGeneration(data []byte) (*ddv1.FullCardNumberAndGeneration, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for FullCardNumberAndGeneration: got %d, want at least 1", len(data))
	}

	fullCardNumberAndGen := &ddv1.FullCardNumberAndGeneration{}

	// Parse full card number (variable length)
	// We need to determine the length of the FullCardNumber first
	// For now, we'll assume it's the last 1 byte is the generation
	// and everything before that is the FullCardNumber
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for FullCardNumberAndGeneration")
	}

	// Parse generation (last byte)
	generation, err := UnmarshalGeneration(data[len(data)-1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse generation: %w", err)
	}
	fullCardNumberAndGen.SetGeneration(generation)

	// Parse full card number (everything except the last byte)
	fullCardNumberData := data[:len(data)-1]
	fullCardNumber, err := UnmarshalFullCardNumber(fullCardNumberData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse full card number: %w", err)
	}
	fullCardNumberAndGen.SetFullCardNumber(fullCardNumber)

	return fullCardNumberAndGen, nil
}

// appendFullCardNumberAndGeneration appends full card number and generation data to dst.
//
// The data type `FullCardNumberAndGeneration` is specified in the Data Dictionary, Section 2.74.
//
// ASN.1 Definition:
//
//	FullCardNumberAndGeneration ::= SEQUENCE {
//	    fullcardNumber FullCardNumber,
//	    generation Generation
//	}
//
// Binary Layout (variable length):
//   - Full Card Number (variable): FullCardNumber structure
//   - Generation (1 byte): Generation enum value
//
//nolint:unused
func AppendFullCardNumberAndGeneration(dst []byte, fullCardNumberAndGen *ddv1.FullCardNumberAndGeneration) ([]byte, error) {
	if fullCardNumberAndGen == nil {
		// Append default values (empty FullCardNumber + default generation)
		dst = append(dst, 0) // Default generation
		return dst, nil
	}

	// Append full card number (variable length)
	fullCardNumber := fullCardNumberAndGen.GetFullCardNumber()
	if fullCardNumber != nil {
		var err error
		dst, err = AppendFullCardNumber(dst, fullCardNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to append full card number: %w", err)
		}
	}

	// Append generation (1 byte)
	generation := fullCardNumberAndGen.GetGeneration()
	dst = append(dst, byte(generation))

	return dst, nil
}
