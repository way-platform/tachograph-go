package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendDrivingLicenceInfo appends the binary representation of DrivingLicenceInfo to dst.
func AppendDrivingLicenceInfo(dst []byte, dli *cardv1.DrivingLicenceInfo) ([]byte, error) {
	if dli == nil {
		return dst, nil
	}
	dst = appendString(dst, dli.GetDrivingLicenceIssuingAuthority(), 36)
	dst = append(dst, byte(dli.GetDrivingLicenceIssuingNation()))
	dst = appendString(dst, dli.GetDrivingLicenceNumber(), 16)
	return dst, nil
}
