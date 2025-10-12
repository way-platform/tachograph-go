package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalOverviewGen2V1 parses Gen2 V1 Overview data from the complete transfer value.
//
// Gen2 V1 Overview structure uses RecordArray format (from Data Dictionary):
//
// ASN.1 Definition:
//
//	VuOverviewSecondGen ::= SEQUENCE {
//	    memberStateCertificateRecordArray                MemberStateCertificateRecordArray,
//	    vuCertificateRecordArray                         VuCertificateRecordArray,
//	    vehicleIdentificationNumberRecordArray           VehicleIdentificationNumberRecordArray,
//	    vehicleRegistrationIdentificationRecordArray     VehicleRegistrationIdentificationRecordArray,
//	    currentDateTimeRecordArray                       CurrentDateTimeRecordArray,
//	    vuDownloadablePeriodRecordArray                  VuDownloadablePeriodRecordArray,
//	    cardSlotsStatusRecordArray                       CardSlotsStatusRecordArray,
//	    vuDownloadActivityDataRecordArray                VuDownloadActivityDataRecordArray,
//	    vuCompanyLocksRecordArray                        VuCompanyLocksRecordArray,
//	    vuControlActivityRecordArray                     VuControlActivityRecordArray,
//	    signatureRecordArray                             SignatureRecordArray
//	}
//
// Each RecordArray has a 5-byte header:
//
//	recordType (1 byte) + recordSize (2 bytes, big-endian) + noOfRecords (2 bytes, big-endian)
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
// Full semantic parsing of all RecordArrays is TODO.
func unmarshalOverviewGen2V1(value []byte) (*vuv1.OverviewGen2V1, error) {
	overview := &vuv1.OverviewGen2V1{}
	overview.SetRawData(value)

	// For now, store the raw data and validate structure by skipping through all record arrays
	offset := 0

	// Helper to skip a RecordArray
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(value, offset)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		offset += size
		return nil
	}

	// MemberStateCertificateRecordArray
	if err := skipRecordArray("MemberStateCertificate"); err != nil {
		return nil, err
	}

	// VUCertificateRecordArray
	if err := skipRecordArray("VUCertificate"); err != nil {
		return nil, err
	}

	// VehicleIdentificationNumberRecordArray
	if err := skipRecordArray("VehicleIdentificationNumber"); err != nil {
		return nil, err
	}

	// VehicleRegistrationIdentificationRecordArray
	if err := skipRecordArray("VehicleRegistrationIdentification"); err != nil {
		return nil, err
	}

	// CurrentDateTimeRecordArray
	if err := skipRecordArray("CurrentDateTime"); err != nil {
		return nil, err
	}

	// VuDownloadablePeriodRecordArray
	if err := skipRecordArray("VuDownloadablePeriod"); err != nil {
		return nil, err
	}

	// CardSlotsStatusRecordArray
	if err := skipRecordArray("CardSlotsStatus"); err != nil {
		return nil, err
	}

	// VuDownloadActivityDataRecordArray
	if err := skipRecordArray("VuDownloadActivityData"); err != nil {
		return nil, err
	}

	// VuCompanyLocksRecordArray
	if err := skipRecordArray("VuCompanyLocks"); err != nil {
		return nil, err
	}

	// VuControlActivityRecordArray
	if err := skipRecordArray("VuControlActivity"); err != nil {
		return nil, err
	}

	// SignatureRecordArray (last)
	if err := skipRecordArray("Signature"); err != nil {
		return nil, err
	}

	// Verify we consumed exactly the right amount of data
	if offset != len(value) {
		return nil, fmt.Errorf("Overview Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	// TODO: Implement full semantic parsing of all record arrays
	// For now, raw_data contains all the information needed for round-trip testing

	return overview, nil
}

// appendOverviewGen2V1 marshals Gen2 V1 Overview data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and matches the structure, it uses it as the output. Otherwise, it would need to
// construct from semantic fields (currently not implemented).
func appendOverviewGen2V1(dst []byte, overview *vuv1.OverviewGen2V1) ([]byte, error) {
	if overview == nil {
		return nil, fmt.Errorf("overview cannot be nil")
	}

	// For Gen2 structures with RecordArrays, raw data painting is straightforward:
	// we use the raw_data if available
	raw := overview.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	// TODO: Implement marshalling from semantic fields
	// This would require constructing all RecordArrays from semantic data
	return nil, fmt.Errorf("cannot marshal Overview Gen2 V1 without raw_data (semantic marshalling not yet implemented)")
}
