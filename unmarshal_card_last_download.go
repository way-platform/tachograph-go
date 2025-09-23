package tachograph

import (
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalCardLastDownload unmarshals last card download data from a card EF.
func UnmarshalCardLastDownload(data []byte, target *cardv1.LastCardDownload) error {
	if len(data) < 4 {
		return fmt.Errorf("insufficient data for last card download")
	}

	r := bytes.NewReader(data)

	// Read timestamp (4 bytes)
	target.SetTimestamp(readTimeReal(r))

	return nil
}
