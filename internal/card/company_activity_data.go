package card

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalCompanyActivityData parses the EF_Company_Activity_Data file.
//
// The data type `CompanyActivityData` is specified in the Data Dictionary, Section 2.46.
//
// ASN.1 Definition:
//
//	CompanyActivityData ::= SEQUENCE {
//	    companyPointerNewestRecord INTEGER(0..NoOfCompanyActivityRecords-1),
//	    companyActivityRecords SET SIZE(NoOfCompanyActivityRecords) OF
//	    companyActivityRecord SEQUENCE {
//	        companyActivityType   CompanyActivityType,
//	        companyActivityTime   TimeReal,
//	        cardNumberInformation FullCardNumberAndGeneration,
//	        vehicleRegistrationInformation VehicleRegistrationIdentification,
//	        downloadPeriodBegin   TimeReal,
//	        downloadPeriodEnd     TimeReal
//	    }
//	}
//
// Binary Layout:
//   - Bytes 0-1: companyPointerNewestRecord (2 bytes, big-endian)
//   - Bytes 2..N: Array of companyActivityRecord (47 bytes each)
//
// Record layout (47 bytes):
//   - Byte  0:     companyActivityType (1 byte)
//   - Bytes 1-4:   companyActivityTime (4 bytes, TimeReal)
//   - Bytes 5-23:  cardNumberInformation (19 bytes, FullCardNumberAndGeneration: 18B FullCardNumber + 1B Generation)
//   - Bytes 24-38: vehicleRegistrationInformation (15 bytes)
//   - Bytes 39-42: downloadPeriodBegin (4 bytes, TimeReal)
//   - Bytes 43-46: downloadPeriodEnd (4 bytes, TimeReal)
func (opts UnmarshalOptions) unmarshalCompanyActivityData(data []byte) (*cardv1.CompanyActivityData, error) {
	const (
		lenPointer = 2
		lenRecord  = 47

		idxActivityType    = 0
		idxActivityTime    = 1
		lenActivityTime    = 4
		idxCardNumber      = 5
		lenCardNumber      = 19
		idxVehicleReg      = 24
		lenVehicleReg      = 15
		idxDownloadBegin   = 39
		lenDownloadBegin   = 4
		idxDownloadEnd     = 43
		lenDownloadEnd     = 4
	)

	if len(data) < lenPointer {
		return nil, fmt.Errorf("data too short for CompanyActivityData: got %d, want at least %d", len(data), lenPointer)
	}

	recordsData := data[lenPointer:]
	if len(recordsData)%lenRecord != 0 {
		return nil, fmt.Errorf("invalid data length for CompanyActivityData records: %d is not a multiple of %d", len(recordsData), lenRecord)
	}

	newestRecordIndex := int32(binary.BigEndian.Uint16(data[0:lenPointer]))

	target := &cardv1.CompanyActivityData{}
	target.SetNewestRecordIndex(newestRecordIndex)

	count := len(recordsData) / lenRecord
	records := make([]*cardv1.CompanyActivityData_Record, count)

	for i := 0; i < count; i++ {
		start := i * lenRecord
		rec := recordsData[start : start+lenRecord]

		r := &cardv1.CompanyActivityData_Record{}

		// companyActivityType (1 byte)
		activityType, err := dd.UnmarshalEnum[ddv1.CompanyActivityType](rec[idxActivityType])
		if err != nil {
			return nil, fmt.Errorf("record %d: failed to unmarshal company activity type: %w", i, err)
		}
		r.SetCompanyActivityType(activityType)

		// companyActivityTime (4 bytes, TimeReal)
		activityTime, err := opts.UnmarshalTimeReal(rec[idxActivityTime : idxActivityTime+lenActivityTime])
		if err != nil {
			return nil, fmt.Errorf("record %d: failed to unmarshal company activity time: %w", i, err)
		}
		r.SetCompanyActivityTime(activityTime)

		// cardNumberInformation (19 bytes, FullCardNumberAndGeneration)
		cardNumber, err := opts.UnmarshalFullCardNumberAndGeneration(rec[idxCardNumber : idxCardNumber+lenCardNumber])
		if err != nil {
			return nil, fmt.Errorf("record %d: failed to unmarshal card number information: %w", i, err)
		}
		r.SetCardNumberInformation(cardNumber)

		// vehicleRegistrationInformation (15 bytes)
		vehicleReg, err := opts.UnmarshalVehicleRegistrationIdentification(rec[idxVehicleReg : idxVehicleReg+lenVehicleReg])
		if err != nil {
			return nil, fmt.Errorf("record %d: failed to unmarshal vehicle registration information: %w", i, err)
		}
		r.SetVehicleRegistrationInformation(vehicleReg)

		// downloadPeriodBegin (4 bytes, TimeReal)
		downloadBegin, err := opts.UnmarshalTimeReal(rec[idxDownloadBegin : idxDownloadBegin+lenDownloadBegin])
		if err != nil {
			return nil, fmt.Errorf("record %d: failed to unmarshal download period begin: %w", i, err)
		}
		r.SetDownloadPeriodBegin(downloadBegin)

		// downloadPeriodEnd (4 bytes, TimeReal)
		downloadEnd, err := opts.UnmarshalTimeReal(rec[idxDownloadEnd : idxDownloadEnd+lenDownloadEnd])
		if err != nil {
			return nil, fmt.Errorf("record %d: failed to unmarshal download period end: %w", i, err)
		}
		r.SetDownloadPeriodEnd(downloadEnd)

		records[i] = r
	}
	target.SetRecords(records)

	return target, nil
}
