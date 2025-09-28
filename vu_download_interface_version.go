package tachograph

import (
	"bytes"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalDownloadInterfaceVersion parses the download interface version from VU data
//
// The data type `DownloadInterfaceVersion` is specified in the Data Dictionary, Section 2.2.6.1.
//
// ASN.1 Definition:
//
//	DownloadInterfaceVersion ::= SEQUENCE {
//	    generation    Generation,
//	    version       Version
//	}
func unmarshalDownloadInterfaceVersion(data []byte, offset int, version *vuv1.DownloadInterfaceVersion) (int, error) {
	startOffset := offset

	// DownloadInterfaceVersion structure (2 bytes: generation + version)
	// See Appendix 7, Section 2.2.6.1
	generationByte, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return 0, err
	}

	versionByte, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return 0, err
	}

	// Map generation byte to enum
	switch generationByte {
	case 1:
		version.SetGeneration(ddv1.Generation_GENERATION_1)
	case 2:
		version.SetGeneration(ddv1.Generation_GENERATION_2)
	default:
		version.SetGeneration(ddv1.Generation_GENERATION_UNSPECIFIED)
	}

	// Map version byte to enum
	switch versionByte {
	case 1:
		version.SetVersion(vuv1.Version_VERSION_1)
	case 2:
		version.SetVersion(vuv1.Version_VERSION_2)
	default:
		version.SetVersion(vuv1.Version_VERSION_UNSPECIFIED)
	}

	bytesRead := offset - startOffset
	return bytesRead, nil
}

// AppendDownloadInterfaceVersion marshals the download interface version to VU data
//
// The data type `DownloadInterfaceVersion` is specified in the Data Dictionary, Section 2.2.6.1.
//
// ASN.1 Definition:
//
//	DownloadInterfaceVersion ::= SEQUENCE {
//	    generation    Generation,
//	    version       Version
//	}
func appendDownloadInterfaceVersion(buf *bytes.Buffer, version *vuv1.DownloadInterfaceVersion) {
	if version == nil {
		return
	}

	// DownloadInterfaceVersion structure (2 bytes: generation + version)
	// See Appendix 7, Section 2.2.6.1

	// Map generation enum to byte
	var generationByte uint8
	switch version.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		generationByte = 1
	case ddv1.Generation_GENERATION_2:
		generationByte = 2
	default:
		generationByte = 0
	}
	buf.WriteByte(generationByte)

	// Map version enum to byte
	var versionByte uint8
	switch version.GetVersion() {
	case vuv1.Version_VERSION_1:
		versionByte = 1
	case vuv1.Version_VERSION_2:
		versionByte = 2
	default:
		versionByte = 0
	}
	buf.WriteByte(versionByte)
}
