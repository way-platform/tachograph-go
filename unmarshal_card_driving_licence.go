package tachograph

import (
	"bytes"
	"errors"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalDrivingLicenceInfo parses the binary data for an EF_Driving_Licence_Info record.
func unmarshalDrivingLicenceInfo(data []byte) (*cardv1.DrivingLicenceInfo, error) {
	if len(data) < 53 {
		return nil, errors.New("not enough data for DrivingLicenceInfo")
	}
	var dli cardv1.DrivingLicenceInfo
	r := bytes.NewReader(data)
	dli.SetDrivingLicenceIssuingAuthority(readString(r, 36))
	nation, _ := r.ReadByte()
	dli.SetDrivingLicenceIssuingNation(int32(nation))
	dli.SetDrivingLicenceNumber(readString(r, 16))
	return &dli, nil
}
