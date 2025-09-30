package dd

import (
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalOptions provides context for parsing binary tachograph data.
//
// This struct follows the pattern used in protojson.UnmarshalOptions and other
// Go standard library packages, where unmarshal functions are methods on the
// options struct.
//
// The zero value (UnmarshalOptions{}) is valid and represents Generation 1,
// Version 1, which is the most common case for tachograph data parsing.
//
// Functions check for Generation == GENERATION_2; all other values (including
// GENERATION_UNSPECIFIED) are treated as Generation 1. Similarly for Version.
type UnmarshalOptions struct {
	// Generation of the data being parsed (Gen1, Gen2, etc.)
	//
	// This is required for parsing generation-dependent data structures where
	// the binary format differs between generations (e.g., PlaceRecord is 10
	// bytes in Gen1, 22 bytes in Gen2).
	//
	// Unknown or unspecified generation is treated as Gen1.
	Generation ddv1.Generation

	// Version indicates the minor version within a generation (e.g., v2, v3).
	//
	// Functions check for specific version numbers (e.g., Version == VERSION_2);
	// all other values (including VERSION_UNSPECIFIED) are treated as version 1.
	//
	// This field is reserved for future use as new versions are introduced.
	Version ddv1.Version
}

// SetFromCardStructureVersion updates the generation and version fields based on
// a CardStructureVersion value.
//
// Per Data Dictionary Section 2.36, card structure versions are encoded as major.minor:
// - Major byte: 00H for Gen1, 01H for Gen2
// - Minor byte: 00H for Gen1, 00H for Gen2v1, 01H for Gen2v2
//
// Examples from the regulation:
// - {00 00} = Generation 1
// - {01 00} = Generation 2, Version 1
// - {01 01} = Generation 2, Version 2
//
// This method can be used by card and VU packages to extract generation/version
// context from CardStructureVersion fields.
func (opts *UnmarshalOptions) SetFromCardStructureVersion(csv *ddv1.CardStructureVersion) {
	switch csv.GetMajor() {
	case 0:
		opts.Generation = ddv1.Generation_GENERATION_1
		opts.Version = ddv1.Version_VERSION_1
	case 1:
		opts.Generation = ddv1.Generation_GENERATION_2
		switch csv.GetMinor() {
		case 0:
			opts.Version = ddv1.Version_VERSION_1
		case 1:
			opts.Version = ddv1.Version_VERSION_2
		default:
			opts.Version = ddv1.Version_VERSION_UNSPECIFIED
		}
	default:
		opts.Generation = ddv1.Generation_GENERATION_UNSPECIFIED
		opts.Version = ddv1.Version_VERSION_UNSPECIFIED
	}
}
