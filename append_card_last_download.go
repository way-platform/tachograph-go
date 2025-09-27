package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardLastDownload appends last card download data to a byte slice.
func AppendCardLastDownload(data []byte, lastDownload *cardv1.CardDownloadDriver) ([]byte, error) {
	if lastDownload == nil {
		return data, nil
	}

	// Timestamp (4 bytes)
	data = appendTimeReal(data, lastDownload.GetTimestamp())

	return data, nil
}
