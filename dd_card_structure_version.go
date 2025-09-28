package tachograph

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// bcdByteToInt converts a single BCD byte to an integer
func bcdByteToInt(b byte) (int, error) {
	high := int((b >> 4) & 0x0F)
	low := int(b & 0x0F)

	// Validate that both nibbles are valid BCD (0-9)
	if high > 9 || low > 9 {
		return 0, fmt.Errorf("invalid BCD byte: 0x%02X", b)
	}

	return high*10 + low, nil
}

// intToBcd converts an integer (0-99) to a BCD byte
func intToBcd(i int) byte {
	if i < 0 || i > 99 {
		panic(fmt.Sprintf("intToBcd: value %d out of range [0, 99]", i))
	}

	high := byte(i / 10)
	low := byte(i % 10)
	return (high << 4) | low
}

// unmarshalCardStructureVersion parses card structure version data.
//
// The data type `CardStructureVersion` is specified in the Data Dictionary, Section 2.36.
//
// ASN.1 Definition:
//
//	CardStructureVersion ::= OCTET STRING (SIZE (2))
//
// Binary Layout (2 bytes):
//   - Major Version (1 byte): BCD-encoded major version
//   - Minor Version (1 byte): BCD-encoded minor version
func unmarshalCardStructureVersion(data []byte) (*ddv1.CardStructureVersion, error) {
	const (
		lenCardStructureVersion = 2 // 1 byte major + 1 byte minor
	)

	if len(data) < lenCardStructureVersion {
		return nil, fmt.Errorf("insufficient data for CardStructureVersion: got %d, want %d", len(data), lenCardStructureVersion)
	}

	cardStructureVersion := &ddv1.CardStructureVersion{}

	// Parse major version (1 byte BCD)
	majorVersion, err := bcdByteToInt(data[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse major version: %w", err)
	}
	cardStructureVersion.SetMajor(int32(majorVersion))

	// Parse minor version (1 byte BCD)
	minorVersion, err := bcdByteToInt(data[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse minor version: %w", err)
	}
	cardStructureVersion.SetMinor(int32(minorVersion))

	return cardStructureVersion, nil
}

// appendCardStructureVersion appends card structure version data to dst.
//
// The data type `CardStructureVersion` is specified in the Data Dictionary, Section 2.36.
//
// ASN.1 Definition:
//
//	CardStructureVersion ::= OCTET STRING (SIZE (2))
//
// Binary Layout (2 bytes):
//   - Major Version (1 byte): BCD-encoded major version
//   - Minor Version (1 byte): BCD-encoded minor version
func appendCardStructureVersion(dst []byte, cardStructureVersion *ddv1.CardStructureVersion) ([]byte, error) {
	if cardStructureVersion == nil {
		// Append default values (2 zero bytes)
		return append(dst, make([]byte, 2)...), nil
	}

	// Append major version (1 byte BCD)
	majorVersion := int(cardStructureVersion.GetMajor())
	if majorVersion < 0 || majorVersion > 99 {
		return nil, fmt.Errorf("major version out of range: %d", majorVersion)
	}
	dst = append(dst, byte(intToBcd(majorVersion)))

	// Append minor version (1 byte BCD)
	minorVersion := int(cardStructureVersion.GetMinor())
	if minorVersion < 0 || minorVersion > 99 {
		return nil, fmt.Errorf("minor version out of range: %d", minorVersion)
	}
	dst = append(dst, byte(intToBcd(minorVersion)))

	return dst, nil
}
