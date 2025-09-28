package tachograph

import (
	"encoding/binary"
	"fmt"
	"time"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalDownloadablePeriod parses downloadable period data.
//
// The data type `DownloadablePeriod` is specified in the Data Dictionary, Section 2.193.
//
// ASN.1 Definition:
//
//	VuDownloadablePeriod ::= SEQUENCE {
//	    minDownloadableTime TimeReal,
//	    maxDownloadableTime TimeReal
//	}
//
//	TimeReal ::= INTEGER (0..2^32-1)
//
// Binary Layout (8 bytes total):
//   - Min Downloadable Time (4 bytes): Unsigned integer (seconds since epoch)
//   - Max Downloadable Time (4 bytes): Unsigned integer (seconds since epoch)
func unmarshalDownloadablePeriod(data []byte) (*ddv1.DownloadablePeriod, error) {
	const (
		lenDownloadablePeriod = 8 // 4 bytes min + 4 bytes max
	)

	if len(data) < lenDownloadablePeriod {
		return nil, fmt.Errorf("insufficient data for DownloadablePeriod: got %d, want %d", len(data), lenDownloadablePeriod)
	}

	downloadablePeriod := &ddv1.DownloadablePeriod{}

	// Parse min downloadable time (4 bytes, unsigned big-endian)
	minTime := binary.BigEndian.Uint32(data[0:4])
	downloadablePeriod.SetMinTime(timestamppb.New(time.Unix(int64(minTime), 0)))

	// Parse max downloadable time (4 bytes, unsigned big-endian)
	maxTime := binary.BigEndian.Uint32(data[4:8])
	downloadablePeriod.SetMaxTime(timestamppb.New(time.Unix(int64(maxTime), 0)))

	return downloadablePeriod, nil
}

// appendDownloadablePeriod appends downloadable period data to dst.
//
// The data type `DownloadablePeriod` is specified in the Data Dictionary, Section 2.193.
//
// ASN.1 Definition:
//
//	VuDownloadablePeriod ::= SEQUENCE {
//	    minDownloadableTime TimeReal,
//	    maxDownloadableTime TimeReal
//	}
//
//	TimeReal ::= INTEGER (0..2^32-1)
//
// Binary Layout (8 bytes total):
//   - Min Downloadable Time (4 bytes): Unsigned integer (seconds since epoch)
//   - Max Downloadable Time (4 bytes): Unsigned integer (seconds since epoch)
func appendDownloadablePeriod(dst []byte, downloadablePeriod *ddv1.DownloadablePeriod) ([]byte, error) {
	if downloadablePeriod == nil {
		// Append default values (8 zero bytes)
		return append(dst, make([]byte, 8)...), nil
	}

	// Append min downloadable time (4 bytes, unsigned big-endian)
	minTime := downloadablePeriod.GetMinTime()
	if minTime != nil {
		minTimeUnix := minTime.GetSeconds()
		dst = binary.BigEndian.AppendUint32(dst, uint32(minTimeUnix))
	} else {
		dst = binary.BigEndian.AppendUint32(dst, 0)
	}

	// Append max downloadable time (4 bytes, unsigned big-endian)
	maxTime := downloadablePeriod.GetMaxTime()
	if maxTime != nil {
		maxTimeUnix := maxTime.GetSeconds()
		dst = binary.BigEndian.AppendUint32(dst, uint32(maxTimeUnix))
	} else {
		dst = binary.BigEndian.AppendUint32(dst, 0)
	}

	return dst, nil
}
