package card

import (
	"github.com/way-platform/tachograph-go/internal/dd"
	"bytes"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardLastDownload unmarshals last card download data from a card EF.
//
// The data type `CardDownloadDriver` is specified in the Data Dictionary, Section 2.16.
//
// ASN.1 Definition:
//
//	CardDownloadDriver ::= TimeReal
func unmarshalCardLastDownload(data []byte) (*cardv1.CardDownloadDriver, error) {
	const (
		lenCardDownloadDriver = 4 // CardDownloadDriver total size
	)

	if len(data) < lenCardDownloadDriver {
		return nil, fmt.Errorf("insufficient data for last card download")
	}

	var target cardv1.CardDownloadDriver
	r := bytes.NewReader(data)

	// Read timestamp (4 bytes)
	target.SetTimestamp(dd.ReadTimeReal(r))

	return &target, nil
}

// AppendCardLastDownload appends last card download data to a byte slice.
//
// The data type `CardDownloadDriver` is specified in the Data Dictionary, Section 2.16.
//
// ASN.1 Definition:
//
//	CardDownloadDriver ::= TimeReal
func appendCardLastDownload(data []byte, lastDownload *cardv1.CardDownloadDriver) ([]byte, error) {
	if lastDownload == nil {
		return data, nil
	}

	// Timestamp (4 bytes)
	data = dd.AppendTimeReal(data, lastDownload.GetTimestamp())

	return data, nil
}
