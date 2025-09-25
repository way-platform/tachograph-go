package tachograph

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardLastDownload unmarshals last card download data from a card EF.
func unmarshalCardLastDownload(data []byte) (*cardv1.LastCardDownload, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for last card download")
	}

	var target cardv1.LastCardDownload
	r := bytes.NewReader(data)

	// Read timestamp (4 bytes)
	target.SetTimestamp(readTimeReal(r))

	return &target, nil
}

// UnmarshalCardLastDownload unmarshals last card download data from a card EF (legacy function).
// Deprecated: Use unmarshalCardLastDownload instead.
func UnmarshalCardLastDownload(data []byte, target *cardv1.LastCardDownload) error {
	result, err := unmarshalCardLastDownload(data)
	if err != nil {
		return err
	}
	*target = *result
	return nil
}
