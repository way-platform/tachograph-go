package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendDrivingLicenceInfo appends the binary representation of DrivingLicenceInfo to dst.
func AppendDrivingLicenceInfo(dst []byte, dli *cardv1.DrivingLicenceInfo) ([]byte, error) {
	if dli == nil {
		return dst, nil
	}
	var err error
	dst, err = appendStringValue(dst, dli.GetDrivingLicenceIssuingAuthority(), 36)
	if err != nil {
		return nil, err
	}
	dst = append(dst, byte(dli.GetDrivingLicenceIssuingNation()))
	if licenceNumber := dli.GetDrivingLicenceNumber(); licenceNumber != nil {
		dst, err = appendString(dst, licenceNumber.GetDecoded(), 16)
	}
	if err != nil {
		return nil, err
	}
	return dst, nil
}
