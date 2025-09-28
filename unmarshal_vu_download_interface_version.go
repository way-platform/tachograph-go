package tachograph

import (
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalDownloadInterfaceVersion parses the download interface version from VU data
//
// ASN.1 Specification (Data Dictionary 2.2.6.1):
//
//	DownloadInterfaceVersion ::= SEQUENCE {
//	    generation    Generation,
//	    version       Version
//	}
//
// Binary Layout (2 bytes):
//
//	0-0:   generation (1 byte)
//	1-1:   version (1 byte)
//
// Constants:
const (
// DownloadInterfaceVersion total size
)

func UnmarshalDownloadInterfaceVersion(data []byte, offset int, version *vuv1.DownloadInterfaceVersion) (int, error) {
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
		version.SetGeneration(datadictionaryv1.Generation_GENERATION_1)
	case 2:
		version.SetGeneration(datadictionaryv1.Generation_GENERATION_2)
	default:
		version.SetGeneration(datadictionaryv1.Generation_GENERATION_UNSPECIFIED)
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
