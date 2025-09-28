package tachograph

import (
	"bytes"
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalDrivingLicenceInfo parses the binary data for an EF_Driving_Licence_Info record.
func unmarshalDrivingLicenceInfo(data []byte) (*cardv1.DrivingLicenceInfo, error) {
	if len(data) < 53 {
		return nil, errors.New("not enough data for DrivingLicenceInfo")
	}
	var dli cardv1.DrivingLicenceInfo
	r := bytes.NewReader(data)
	authority, err := unmarshalStringValueFromReader(r, 36)
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence issuing authority: %w", err)
	}
	dli.SetDrivingLicenceIssuingAuthority(authority)

	nation, err := unmarshalNationNumericFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence issuing nation: %w", err)
	}
	dli.SetDrivingLicenceIssuingNation(int32(nation))

	licenceNumber, err := unmarshalIA5StringValueFromReader(r, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence number: %w", err)
	}
	dli.SetDrivingLicenceNumber(licenceNumber.GetDecoded())
	return &dli, nil
}
