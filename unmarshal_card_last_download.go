package tachograph

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalCardLastDownload unmarshals last card download data from a card EF.
//
// ASN.1 Specification (Data Dictionary 2.16):
//
//	CardDownloadDriver ::= TimeReal
//
// Binary Layout (4 bytes):
//
//	0-3:   timestamp (4 bytes, TimeReal)
//
// Constants:
const (
	// CardDownloadDriver total size
	cardDownloadDriverSize = 4
)

func unmarshalCardLastDownload(data []byte) (*cardv1.CardDownloadDriver, error) {
	if len(data) < cardDownloadDriverSize {
		return nil, fmt.Errorf("insufficient data for last card download")
	}

	var target cardv1.CardDownloadDriver
	r := bytes.NewReader(data)

	// Read timestamp (4 bytes)
	target.SetTimestamp(readTimeReal(r))

	return &target, nil
}

// UnmarshalCardLastDownload unmarshals last card download data from a card EF (legacy function).
// Deprecated: Use unmarshalCardLastDownload instead.
func UnmarshalCardLastDownload(data []byte, target *cardv1.CardDownloadDriver) error {
	result, err := unmarshalCardLastDownload(data)
	if err != nil {
		return err
	}
	proto.Merge(target, result)
	return nil
}
