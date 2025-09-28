package tachograph

import (
	"encoding/binary"
	"fmt"
	"time"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalSoftwareIdentification parses software identification data.
//
// The data type `SoftwareIdentification` is specified in the Data Dictionary, Section 2.225.
//
// ASN.1 Definition:
//
//	VuSoftwareIdentification ::= SEQUENCE {
//	    vuSoftwareVersion VuSoftwareVersion,
//	    vuSoftInstallationDate VuSoftInstallationDate
//	}
//
//	VuSoftwareVersion ::= IA5String(SIZE(4))
//	VuSoftInstallationDate ::= TimeReal ::= INTEGER (0..2^32-1)
//
// Binary Layout (8 bytes total):
//   - Software Version (4 bytes): IA5String (ASCII)
//   - Software Installation Date (4 bytes): Unsigned integer (seconds since epoch)
func unmarshalSoftwareIdentification(data []byte) (*ddv1.SoftwareIdentification, error) {
	const (
		lenSoftwareIdentification = 8 // 4 bytes version + 4 bytes date
	)

	if len(data) < lenSoftwareIdentification {
		return nil, fmt.Errorf("insufficient data for SoftwareIdentification: got %d, want %d", len(data), lenSoftwareIdentification)
	}

	softwareID := &ddv1.SoftwareIdentification{}

	// Parse software version (4 bytes IA5String)
	versionData := data[0:4]
	version, err := unmarshalIA5StringValue(versionData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse software version: %w", err)
	}
	softwareID.SetSoftwareVersion(version)

	// Parse software installation date (4 bytes, unsigned big-endian)
	installationDate := binary.BigEndian.Uint32(data[4:8])
	softwareID.SetSoftwareInstallationDate(timestamppb.New(time.Unix(int64(installationDate), 0)))

	return softwareID, nil
}

// appendSoftwareIdentification appends software identification data to dst.
//
// The data type `SoftwareIdentification` is specified in the Data Dictionary, Section 2.225.
//
// ASN.1 Definition:
//
//	VuSoftwareIdentification ::= SEQUENCE {
//	    vuSoftwareVersion VuSoftwareVersion,
//	    vuSoftInstallationDate VuSoftInstallationDate
//	}
//
//	VuSoftwareVersion ::= IA5String(SIZE(4))
//	VuSoftInstallationDate ::= TimeReal ::= INTEGER (0..2^32-1)
//
// Binary Layout (8 bytes total):
//   - Software Version (4 bytes): IA5String (ASCII)
//   - Software Installation Date (4 bytes): Unsigned integer (seconds since epoch)
func appendSoftwareIdentification(dst []byte, softwareID *ddv1.SoftwareIdentification) ([]byte, error) {
	if softwareID == nil {
		// Append default values (8 zero bytes)
		return append(dst, make([]byte, 8)...), nil
	}

	// Append software version (4 bytes IA5String)
	version := softwareID.GetSoftwareVersion()
	if version != nil {
		var err error
		dst, err = appendStringValue(dst, version, 4)
		if err != nil {
			return nil, fmt.Errorf("failed to append software version: %w", err)
		}
	} else {
		// Append 4 zero bytes for empty version
		dst = append(dst, make([]byte, 4)...)
	}

	// Append software installation date (4 bytes, unsigned big-endian)
	installationDate := softwareID.GetSoftwareInstallationDate()
	if installationDate != nil {
		installationDateUnix := installationDate.GetSeconds()
		dst = binary.BigEndian.AppendUint32(dst, uint32(installationDateUnix))
	} else {
		dst = binary.BigEndian.AppendUint32(dst, 0)
	}

	return dst, nil
}
