package tachograph

import (
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
