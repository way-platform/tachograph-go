package tachograph

import (
	"bytes"
	"errors"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalDrivingLicenceInfo parses the binary data for an EF_Driving_Licence_Info record.
func UnmarshalDrivingLicenceInfo(data []byte, dli *cardv1.DrivingLicenceInfo) error {
	if len(data) < 53 {
		return errors.New("not enough data for DrivingLicenceInfo")
	}
	r := bytes.NewReader(data)

	dli.SetDrivingLicenceIssuingAuthority(readString(r, 36))
	nation, _ := r.ReadByte()
	dli.SetDrivingLicenceIssuingNation(int32(nation))
	dli.SetDrivingLicenceNumber(readString(r, 16))

	return nil
}
