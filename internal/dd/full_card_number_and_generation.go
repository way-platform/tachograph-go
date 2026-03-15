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
func (opts UnmarshalOptions) UnmarshalFullCardNumberAndGeneration(data []byte) (*ddv1.FullCardNumberAndGeneration, error) {
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

	// Parse generation (last byte).
	// 0xFF means no card present — treat as unspecified rather than failing.
	if generation, err := UnmarshalEnum[ddv1.Generation](data[len(data)-1]); err == nil {
		fullCardNumberAndGen.SetGeneration(generation)
	} else {
		fullCardNumberAndGen.SetGeneration(ddv1.Generation_GENERATION_UNSPECIFIED)
	}

	// Parse full card number (everything except the last byte)
	fullCardNumberData := data[:len(data)-1]
	fullCardNumber, err := opts.UnmarshalFullCardNumber(fullCardNumberData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse full card number: %w", err)
	}
	fullCardNumberAndGen.SetFullCardNumber(fullCardNumber)

	return fullCardNumberAndGen, nil
}

// MarshalFullCardNumberAndGeneration marshals full card number and generation data to bytes.
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
func (opts MarshalOptions) MarshalFullCardNumberAndGeneration(fullCardNumberAndGen *ddv1.FullCardNumberAndGeneration) ([]byte, error) {
	if fullCardNumberAndGen == nil {
		return nil, fmt.Errorf("fullCardNumberAndGeneration cannot be nil")
	}

	var dst []byte

	// Marshal full card number (variable length)
	fullCardNumber := fullCardNumberAndGen.GetFullCardNumber()
	if fullCardNumber != nil {
		fullCardNumberBytes, err := opts.MarshalFullCardNumber(fullCardNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal full card number: %w", err)
		}
		dst = append(dst, fullCardNumberBytes...)
	}

	// Marshal generation (1 byte).
	// GENERATION_UNSPECIFIED means "no card present" — encode as 0xFF.
	gen := fullCardNumberAndGen.GetGeneration()
	if gen == ddv1.Generation_GENERATION_UNSPECIFIED {
		dst = append(dst, 0xFF)
	} else {
		generationByte, err := MarshalEnum(gen)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal generation: %w", err)
		}
		dst = append(dst, generationByte)
	}

	return dst, nil
}

// AnonymizeFullCardNumberAndGeneration anonymizes a full card number with generation.
func (opts AnonymizeOptions) AnonymizeFullCardNumberAndGeneration(fc *ddv1.FullCardNumberAndGeneration) *ddv1.FullCardNumberAndGeneration {
	if fc == nil {
		return nil
	}

	result := &ddv1.FullCardNumberAndGeneration{}

	// Anonymize the nested FullCardNumber
	result.SetFullCardNumber(opts.AnonymizeFullCardNumber(fc.GetFullCardNumber()))

	// Preserve the generation
	result.SetGeneration(fc.GetGeneration())

	return result
}
