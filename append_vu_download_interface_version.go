package tachograph

import (
	"bytes"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AppendDownloadInterfaceVersion marshals the download interface version to VU data
func AppendDownloadInterfaceVersion(buf *bytes.Buffer, version *vuv1.DownloadInterfaceVersion) {
	if version == nil {
		return
	}

	// DownloadInterfaceVersion structure (2 bytes: generation + version)
	// See Appendix 7, Section 2.2.6.1

	// Map generation enum to byte
	var generationByte uint8
	switch version.GetGeneration() {
	case vuv1.Generation_GENERATION_1:
		generationByte = 1
	case vuv1.Generation_GENERATION_2:
		generationByte = 2
	default:
		generationByte = 0
	}
	appendUint8(buf, generationByte)

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
	appendUint8(buf, versionByte)
}
