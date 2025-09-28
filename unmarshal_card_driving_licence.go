package tachograph

import (
	"errors"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalDrivingLicenceInfo parses the binary data for an EF_Driving_Licence_Info record.
//
// ASN.1 Specification (Data Dictionary 2.18):
//
//	CardDrivingLicenceInformation ::= SEQUENCE {
//	    drivingLicenceIssuingAuthority     Name,
//	    drivingLicenceIssuingNation        NationNumeric,
//	    drivingLicenceNumber               Name
//	}
//
// Binary Layout (53 bytes):
//
//	0-35:  drivingLicenceIssuingAuthority (36 bytes, Name)
//	36-36: drivingLicenceIssuingNation (1 byte, NationNumeric)
//	37-52: drivingLicenceNumber (16 bytes, Name)
//
// Constants:
const (
	// CardDrivingLicenceInformation total size
	cardDrivingLicenceInformationSize = 53
)

func unmarshalDrivingLicenceInfo(data []byte) (*cardv1.DrivingLicenceInfo, error) {
	if len(data) < cardDrivingLicenceInformationSize {
		return nil, errors.New("not enough data for DrivingLicenceInfo")
	}
	var dli cardv1.DrivingLicenceInfo
	offset := 0

	// Read driving licence issuing authority (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence issuing authority")
	}
	authority, err := unmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence issuing authority: %w", err)
	}
	dli.SetDrivingLicenceIssuingAuthority(authority)
	offset += 36

	// Read driving licence issuing nation (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence issuing nation")
	}
	nation, err := unmarshalNationNumeric(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence issuing nation: %w", err)
	}
	dli.SetDrivingLicenceIssuingNation(int32(nation))
	offset++

	// Read driving licence number (16 bytes)
	if offset+16 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence number")
	}
	licenceNumber, err := unmarshalIA5StringValue(data[offset : offset+16])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence number: %w", err)
	}
	dli.SetDrivingLicenceNumber(licenceNumber.GetDecoded())
	offset += 16

	return &dli, nil
}
