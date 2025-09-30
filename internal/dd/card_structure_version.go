package dd

import (
	"bytes"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalCardStructureVersion parses a BCD-encoded CardStructureVersion from a 2-byte slice.
//
// The data type `CardStructureVersion` is specified in the Data Dictionary, Section 2.36.
//
// ASN.1 Specification:
//
//	CardStructureVersion ::= OCTET STRING (SIZE (2))
//
// Binary Layout (2 bytes):
//
//	The two bytes are BCD-encoded, representing major and minor versions.
//	For example, version '01.02' is coded as '0102'H.
//	- Byte 0: Major version in BCD (e.g., 0x01 = 01)
//	- Byte 1: Minor version in BCD (e.g., 0x02 = 02)
func UnmarshalCardStructureVersion(data []byte) (*ddv1.CardStructureVersion, error) {
	const lenCardStructureVersion = 2
	if len(data) != lenCardStructureVersion {
		return nil, fmt.Errorf("invalid data length for CardStructureVersion: got %d, want %d", len(data), lenCardStructureVersion)
	}
	var output ddv1.CardStructureVersion
	output.SetRawData(bytes.Clone(data))
	output.SetMajor(int32((data[0]&0xF0)>>4)*10 + int32(data[0]&0x0F))
	output.SetMinor(int32((data[1]&0xF0)>>4)*10 + int32(data[1]&0x0F))
	return &output, nil
}

// AppendCardStructureVersion appends the binary representation of CardStructureVersion to dst.
//
// The data type `CardStructureVersion` is specified in the Data Dictionary, Section 2.36.
//
// ASN.1 Specification:
//
//	CardStructureVersion ::= OCTET STRING (SIZE (2))
//
// Binary Layout (2 bytes):
//
//	The two bytes are BCD-encoded, representing major and minor versions.
//	For example, version '01.02' is coded as '0102'H.
//	- Byte 0: Major version in BCD (e.g., 0x01 = 01)
//	- Byte 1: Minor version in BCD (e.g., 0x02 = 02)
func AppendCardStructureVersion(dst []byte, csv *ddv1.CardStructureVersion) ([]byte, error) {
	const lenCardStructureVersion = 2
	var canvas [lenCardStructureVersion]byte
	if csv.HasRawData() {
		if len(csv.GetRawData()) != lenCardStructureVersion {
			return nil, fmt.Errorf(
				"invalid raw_data length for CardStructureVersion: got %d, want %d",
				len(csv.GetRawData()), lenCardStructureVersion,
			)
		}
		copy(canvas[:], csv.GetRawData())
	}
	major := csv.GetMajor()
	majorBCD := byte((major/10)<<4) | byte(major%10)
	canvas[0] = majorBCD
	minor := csv.GetMinor()
	minorBCD := byte((minor/10)<<4) | byte(minor%10)
	canvas[1] = minorBCD
	return append(dst, canvas[:]...), nil
}
