package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGeneration parses generation from raw data.
//
// The data type `Generation` is specified in the Data Dictionary.
//
// ASN.1 Definition:
//
//	Generation ::= INTEGER (1..2)
//
// Binary Layout (1 byte):
//   - Generation (1 byte): Raw integer value (1-2)
func UnmarshalGeneration(data []byte) (ddv1.Generation, error) {
	if len(data) != 1 {
		return ddv1.Generation_GENERATION_UNSPECIFIED, fmt.Errorf("invalid data length for Generation: got %d, want 1", len(data))
	}

	generationByte := data[0]

	// Map generation byte to enum
	switch generationByte {
	case 1:
		return ddv1.Generation_GENERATION_1, nil
	case 2:
		return ddv1.Generation_GENERATION_2, nil
	default:
		return ddv1.Generation_GENERATION_UNSPECIFIED, fmt.Errorf("invalid generation value: %d", generationByte)
	}
}

// AppendGeneration appends generation as a single byte.
//
// The data type `Generation` is specified in the Data Dictionary.
//
// ASN.1 Definition:
//
//	Generation ::= INTEGER (1..2)
//
// Binary Layout (1 byte):
//   - Generation (1 byte): Raw integer value (1-2)
func AppendGeneration(dst []byte, generation ddv1.Generation) []byte {
	// Map enum to generation byte
	switch generation {
	case ddv1.Generation_GENERATION_1:
		return append(dst, 1)
	case ddv1.Generation_GENERATION_2:
		return append(dst, 2)
	default:
		return append(dst, 0) // Default to 0 for unspecified
	}
}
