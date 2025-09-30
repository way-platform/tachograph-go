package card

import (
	"errors"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDrivingLicenceInfo parses the binary data for an EF_Driving_Licence_Info record.
//
// The data type `CardDrivingLicenceInformation` is specified in the Data Dictionary, Section 2.18.
//
// ASN.1 Definition:
//
//	CardDrivingLicenceInformation ::= SEQUENCE {
//	    drivingLicenceIssuingAuthority     Name,
//	    drivingLicenceIssuingNation        NationNumeric,
//	    drivingLicenceNumber               Name
//	}
func (opts UnmarshalOptions) unmarshalDrivingLicenceInfo(data []byte) (*cardv1.DrivingLicenceInfo, error) {
	const (
		lenCardDrivingLicenceInformation = 53 // CardDrivingLicenceInformation total size
	)

	if len(data) < lenCardDrivingLicenceInformation {
		return nil, errors.New("not enough data for DrivingLicenceInfo")
	}
	var dli cardv1.DrivingLicenceInfo
	offset := 0

	// Read driving licence issuing authority (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence issuing authority")
	}
	authority, err := opts.UnmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence issuing authority: %w", err)
	}
	dli.SetDrivingLicenceIssuingAuthority(authority)
	offset += 36

	// Read driving licence issuing nation (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence issuing nation")
	}
	if nation, err := dd.UnmarshalEnum[ddv1.NationNumeric](data[offset]); err == nil {
		dli.SetDrivingLicenceIssuingNation(nation)
	} else {
		// Value not recognized - set UNRECOGNIZED (no unrecognized field for this type)
		dli.SetDrivingLicenceIssuingNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}
	offset++

	// Read driving licence number (16 bytes)
	if offset+16 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence number")
	}
	licenceNumber, err := opts.UnmarshalIA5StringValue(data[offset : offset+16])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence number: %w", err)
	}
	dli.SetDrivingLicenceNumber(licenceNumber)
	// offset += 16 // Not needed as this is the last field

	return &dli, nil
}

// AppendDrivingLicenceInfo appends the binary representation of DrivingLicenceInfo to dst.
//
// The data type `CardDrivingLicenceInformation` is specified in the Data Dictionary, Section 2.18.
//
// ASN.1 Definition:
//
//	CardDrivingLicenceInformation ::= SEQUENCE {
//	    drivingLicenceIssuingAuthority     Name,
//	    drivingLicenceIssuingNation        NationNumeric,
//	    drivingLicenceNumber               Name
//	}
func appendDrivingLicenceInfo(dst []byte, dli *cardv1.DrivingLicenceInfo) ([]byte, error) {
	if dli == nil {
		return dst, nil
	}
	var err error
	dst, err = dd.AppendStringValue(dst, dli.GetDrivingLicenceIssuingAuthority())
	if err != nil {
		return nil, err
	}
	dst = append(dst, byte(dli.GetDrivingLicenceIssuingNation()))
	dst, err = dd.AppendStringValue(dst, dli.GetDrivingLicenceNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to append driving licence number: %w", err)
	}
	return dst, nil
}
