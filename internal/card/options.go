package card

import (
	"github.com/way-platform/tachograph-go/internal/dd"
)

// UnmarshalOptions provides context for parsing binary card data.
//
// This struct embeds dd.UnmarshalOptions to inherit all data dictionary unmarshal methods,
// and extends it with card-specific unmarshal methods.
//
// For card files, generation and version information is determined from two sources:
// 1. The CardStructureVersion field in EF_Application_Identification (file-level)
// 2. The TLV tag appendix byte for each Elementary File (EF-level)
//
// The zero value (UnmarshalOptions{}) is valid and represents Generation 1,
// Version 1, which is the default for tachograph card data parsing.
//
// Functions check for Generation == GENERATION_2; all other values (including
// GENERATION_UNSPECIFIED) are treated as Generation 1. Similarly for Version.
type UnmarshalOptions struct {
	// Embed dd.UnmarshalOptions to inherit all data dictionary unmarshal methods.
	// This allows card.UnmarshalOptions to be used wherever dd.UnmarshalOptions is needed.
	dd.UnmarshalOptions
}
