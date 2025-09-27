package tachograph

import (
	"bytes"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalDownloadInterfaceVersion parses the download interface version from VU data
func UnmarshalDownloadInterfaceVersion(r *bytes.Reader, version *vuv1.DownloadInterfaceVersion) (int, error) {
	startPos := int64(r.Len())

	// DownloadInterfaceVersion structure (2 bytes: generation + version)
	// See Appendix 7, Section 2.2.6.1
	generationByte, err := readUint8(r)
	if err != nil {
		return 0, err
	}

	versionByte, err := readUint8(r)
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

	bytesRead := int(startPos - int64(r.Len()))
	return bytesRead, nil
}
