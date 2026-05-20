package vu

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/cert"
	"github.com/way-platform/tachograph-go/internal/security"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AuthenticateOptions configures the VU authentication process.
type AuthenticateOptions struct {
	// CertificateResolver is used to resolve CA certificates by their Certificate Authority Reference (CAR).
	CertificateResolver cert.Resolver
}

// AuthenticateRawVehicleUnitFile performs cryptographic authentication on all records
// in a raw vehicle unit file. For each record, it verifies the signature against the
// data and populates the authentication field.
//
// This function mutates the rawFile parameter by setting the authentication field
// on each record.
//
// Authentication is performed by:
// 1. Extracting the equipment certificate from the appropriate transfer type
// 2. Verifying the certificate chain against the European Root CA
// 3. Verifying the data signature using the equipment certificate
//
// If any record fails authentication, its authentication.status will be set to the
// appropriate error status, and this function will return an error after processing
// all records.
func (opts AuthenticateOptions) AuthenticateRawVehicleUnitFile(ctx context.Context, rawFile *vuv1.RawVehicleUnitFile) error {
	if rawFile == nil {
		return fmt.Errorf("rawFile cannot be nil")
	}
	if opts.CertificateResolver == nil {
		opts.CertificateResolver = cert.DefaultResolver()
	}

	var errs []error
	records := rawFile.GetRecords()

	for _, record := range records {
		if err := opts.authenticateRecord(ctx, record, records); err != nil {
			errs = append(errs, err)
			// Continue processing other records even if one fails
		}
	}

	return errors.Join(errs...)
}

// authenticateRecord authenticates a single VU record by verifying its signature
// against the data. This function mutates the record by setting its authentication field.
func (opts AuthenticateOptions) authenticateRecord(ctx context.Context, record *vuv1.RawVehicleUnitFile_Record, allRecords []*vuv1.RawVehicleUnitFile_Record) error {
	generation := record.GetGeneration()
	transferType := record.GetType()

	// Initialize authentication result
	auth := &securityv1.Authentication{}
	record.SetAuthentication(auth)

	// Select authentication strategy based on generation
	switch generation {
	case ddv1.Generation_GENERATION_1:
		return opts.authenticateGen1Record(ctx, record, allRecords, auth)
	case ddv1.Generation_GENERATION_2:
		return opts.authenticateGen2Record(ctx, record, allRecords, auth)
	default:
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("unsupported generation for %v: %v", transferType, generation)
	}
}

// authenticateGen1Record authenticates a Generation 1 VU record using RSA signature recovery.
func (opts AuthenticateOptions) authenticateGen1Record(ctx context.Context, record *vuv1.RawVehicleUnitFile_Record, allRecords []*vuv1.RawVehicleUnitFile_Record, auth *securityv1.Authentication) error {
	// Step 1: Find the Overview record to get the VU and MSCA certificates
	overviewRecord := opts.findOverviewRecord(allRecords)
	if overviewRecord == nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("overview record not found for authentication")
	}

	// Step 2: Parse the Overview to extract certificates
	vuCert, mscaCert, err := opts.extractGen1Certificates(overviewRecord)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("failed to extract Gen1 certificates: %w", err)
	}

	// Step 3: Verify certificate chain
	if err := opts.verifyGen1CertificateChain(ctx, vuCert, mscaCert, auth); err != nil {
		return err
	}

	// Step 4: Verify the data signature
	if err := opts.verifyGen1DataSignature(record, vuCert, auth); err != nil {
		return err
	}

	// Authentication succeeded
	auth.SetStatus(securityv1.Authentication_VERIFIED)
	auth.SetSignatureAlgorithm(securityv1.SignatureAlgorithm_SHA1_WITH_RSA_ENCRYPTION)

	return nil
}

// authenticateGen2Record authenticates a Generation 2 VU record using ECDSA signature verification.
func (opts AuthenticateOptions) authenticateGen2Record(ctx context.Context, record *vuv1.RawVehicleUnitFile_Record, allRecords []*vuv1.RawVehicleUnitFile_Record, auth *securityv1.Authentication) error {
	// Step 1: Find the Overview record to get the VU and MSCA certificates
	overviewRecord := opts.findGen2OverviewRecord(allRecords)
	if overviewRecord == nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("gen2 overview record not found for authentication")
	}

	// Step 2: Parse the Overview to extract certificates
	vuCert, mscaCert, err := opts.extractGen2Certificates(overviewRecord)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("failed to extract Gen2 certificates: %w", err)
	}

	// Step 3: Verify certificate chain
	if err := opts.verifyGen2CertificateChain(ctx, vuCert, mscaCert, auth); err != nil {
		return err
	}

	// Step 4: Verify the data signature
	if err := opts.verifyGen2DataSignature(record, vuCert, auth); err != nil {
		return err
	}

	// Authentication succeeded
	auth.SetStatus(securityv1.Authentication_VERIFIED)

	// Infer signature algorithm from VU certificate's curve
	if pubKey := vuCert.GetPublicKey(); pubKey != nil {
		oid := pubKey.GetDomainParametersOid()
		switch oid {
		case "1.3.36.3.3.2.8.1.1.7", "1.2.840.10045.3.1.7": // brainpoolP256r1, NIST P-256
			auth.SetSignatureAlgorithm(securityv1.SignatureAlgorithm_ECDSA_WITH_SHA256)
		case "1.3.36.3.3.2.8.1.1.11", "1.3.132.0.34": // brainpoolP384r1, NIST P-384
			auth.SetSignatureAlgorithm(securityv1.SignatureAlgorithm_ECDSA_WITH_SHA384)
		case "1.3.36.3.3.2.8.1.1.13", "1.3.132.0.35": // brainpoolP512r1, NIST P-521
			auth.SetSignatureAlgorithm(securityv1.SignatureAlgorithm_ECDSA_WITH_SHA512)
		default:
			auth.SetSignatureAlgorithm(securityv1.SignatureAlgorithm_SIGNATURE_ALGORITHM_UNSPECIFIED)
		}
	}

	return nil
}

// findGen2OverviewRecord finds the first Gen2 Overview record in the file.
// This record contains the VU and MSCA certificates needed for authentication.
func (opts AuthenticateOptions) findGen2OverviewRecord(records []*vuv1.RawVehicleUnitFile_Record) *vuv1.RawVehicleUnitFile_Record {
	for _, rec := range records {
		transferType := rec.GetType()
		if transferType == vuv1.TransferType_OVERVIEW_GEN2_V1 || transferType == vuv1.TransferType_OVERVIEW_GEN2_V2 {
			return rec
		}
	}
	return nil
}

// extractGen2Certificates extracts the VU and MSCA ECC certificates from a Gen2 Overview record.
func (opts AuthenticateOptions) extractGen2Certificates(overviewRecord *vuv1.RawVehicleUnitFile_Record) (*securityv1.EccCertificate, *securityv1.EccCertificate, error) {
	data, _, err := splitTransferValue(overviewRecord)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to split Overview transfer value: %w", err)
	}

	// Gen2 Overview structure (from regulation):
	// Both V1 and V2 start with RecordArray structures for certificates:
	// - MemberStateCertificateRecordArray: 5-byte header + certificate data
	// - VuCertificateRecordArray: 5-byte header + certificate data
	//
	// RecordArray header:
	// - Byte 0: recordType
	// - Bytes 1-2: recordSize (big-endian, uint16)
	// - Bytes 3-4: noOfRecords (big-endian, uint16)

	// Parse MSCA certificate RecordArray
	if len(data) < 5 {
		return nil, nil, fmt.Errorf("insufficient data for MSCA RecordArray header")
	}

	mscaRecordSize := int(binary.BigEndian.Uint16(data[1:3]))
	mscaNoOfRecords := int(binary.BigEndian.Uint16(data[3:5]))

	if mscaNoOfRecords != 1 {
		return nil, nil, fmt.Errorf("expected exactly 1 MSCA certificate, got %d", mscaNoOfRecords)
	}

	mscaArraySize := 5 + mscaRecordSize // header + data
	if len(data) < mscaArraySize {
		return nil, nil, fmt.Errorf("insufficient data for MSCA certificate: need %d, have %d", mscaArraySize, len(data))
	}

	mscaCertData := data[5:mscaArraySize]

	// Parse VU certificate RecordArray
	vuArrayStart := mscaArraySize
	if len(data) < vuArrayStart+5 {
		return nil, nil, fmt.Errorf("insufficient data for VU RecordArray header")
	}

	vuRecordSize := int(binary.BigEndian.Uint16(data[vuArrayStart+1 : vuArrayStart+3]))
	vuNoOfRecords := int(binary.BigEndian.Uint16(data[vuArrayStart+3 : vuArrayStart+5]))

	if vuNoOfRecords != 1 {
		return nil, nil, fmt.Errorf("expected exactly 1 VU certificate, got %d", vuNoOfRecords)
	}

	vuArraySize := 5 + vuRecordSize
	if len(data) < vuArrayStart+vuArraySize {
		return nil, nil, fmt.Errorf("insufficient data for VU certificate: need %d, have %d", vuArrayStart+vuArraySize, len(data))
	}

	vuCertData := data[vuArrayStart+5 : vuArrayStart+vuArraySize]

	// Unmarshal ECC certificates
	mscaCert, err := security.UnmarshalEccCertificate(mscaCertData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal MSCA certificate: %w", err)
	}

	vuCert, err := security.UnmarshalEccCertificate(vuCertData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal VU certificate: %w", err)
	}

	return vuCert, mscaCert, nil
}

// verifyGen2CertificateChain verifies the Gen2 certificate chain: EUR Root (ECC) -> MSCA (ECC) -> VU (ECC)
func (opts AuthenticateOptions) verifyGen2CertificateChain(ctx context.Context, vuCert *securityv1.EccCertificate, mscaCert *securityv1.EccCertificate, auth *securityv1.Authentication) error {
	// Get Gen2 ECC root certificate
	rootCert, err := opts.CertificateResolver.GetEccRootCertificate(ctx)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("failed to get Gen2 root certificate: %w", err)
	}

	// Verify MSCA certificate against EUR Gen2 root (both ECC)
	if err := security.VerifyEccCertificateWithEccRoot(mscaCert, rootCert); err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("MSCA certificate verification failed: %w", err)
	}

	// Verify VU certificate against MSCA
	if err := security.VerifyEccCertificateWithCA(vuCert, mscaCert); err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("VU certificate verification failed: %w", err)
	}

	return nil
}

// verifyGen2DataSignature verifies the ECDSA signature on the data portion of a Gen2 record.
func (opts AuthenticateOptions) verifyGen2DataSignature(record *vuv1.RawVehicleUnitFile_Record, vuCert *securityv1.EccCertificate, auth *securityv1.Authentication) error {
	data, signature, err := splitTransferValue(record)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("failed to split transfer value: %w", err)
	}

	if len(signature) == 0 {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("no signature present in Gen2 record")
	}

	// For Gen2, the signature is over all the data in the transfer
	// The signature format is plain ECDSA (R || S)
	if err := security.VerifyEccDataSignature(data, signature, vuCert); err != nil {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("data signature verification failed: %w", err)
	}

	return nil
}

// findOverviewRecord finds the first Gen1 Overview record in the file.
// This record contains the VU and MSCA certificates needed for authentication.
func (opts AuthenticateOptions) findOverviewRecord(records []*vuv1.RawVehicleUnitFile_Record) *vuv1.RawVehicleUnitFile_Record {
	for _, rec := range records {
		if rec.GetType() == vuv1.TransferType_OVERVIEW_GEN1 {
			return rec
		}
	}
	return nil
}

// extractGen1Certificates extracts the VU and MSCA certificates from a Gen1 Overview record.
func (opts AuthenticateOptions) extractGen1Certificates(overviewRecord *vuv1.RawVehicleUnitFile_Record) (*securityv1.RsaCertificate, *securityv1.RsaCertificate, error) {
	data, _, err := splitTransferValue(overviewRecord)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to split Overview transfer value: %w", err)
	}

	// Gen1 Overview structure (from regulation):
	// - MemberStateCertificate: 194 bytes
	// - VuCertificate: 194 bytes
	// - VehicleIdentificationNumber: variable
	// - ...rest of data

	const lenRsaCertificate = 194
	const idxMscaCert = 0
	const idxVuCert = 194
	const minDataLen = 388 // Two certificates

	if len(data) < minDataLen {
		return nil, nil, fmt.Errorf("insufficient data for certificates: got %d, need at least %d", len(data), minDataLen)
	}

	// Extract MSCA certificate
	mscaCertData := data[idxMscaCert : idxMscaCert+lenRsaCertificate]
	mscaCert := &securityv1.RsaCertificate{}
	mscaCert.SetRawData(mscaCertData)

	// Extract VU certificate
	vuCertData := data[idxVuCert : idxVuCert+lenRsaCertificate]
	vuCert := &securityv1.RsaCertificate{}
	vuCert.SetRawData(vuCertData)

	return vuCert, mscaCert, nil
}

// verifyGen1CertificateChain verifies the certificate chain for Gen1:
// EUR Root -> MSCA -> VU
func (opts AuthenticateOptions) verifyGen1CertificateChain(ctx context.Context, vuCert *securityv1.RsaCertificate, mscaCert *securityv1.RsaCertificate, auth *securityv1.Authentication) error {
	// Step 1: Get EUR root certificate
	rootCert, err := opts.CertificateResolver.GetRootCertificate(ctx)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("failed to get root certificate: %w", err)
	}

	// Step 2: Verify MSCA certificate against EUR root
	if err := security.VerifyRsaCertificateWithRoot(mscaCert, rootCert); err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("MSCA certificate verification failed: %w", err)
	}

	// Step 3: Verify VU certificate against MSCA
	if err := security.VerifyRsaCertificateWithCA(vuCert, mscaCert); err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("VU certificate verification failed: %w", err)
	}

	// Step 4: Populate certificate info in auth
	// TODO: Extract certificate info (CHR, nation, validity dates) and populate
	// auth.signer_certificate and auth.root_certificate

	return nil
}

// verifyGen1DataSignature verifies the RSA signature on the data portion of a Gen1 record.
func (opts AuthenticateOptions) verifyGen1DataSignature(record *vuv1.RawVehicleUnitFile_Record, vuCert *securityv1.RsaCertificate, auth *securityv1.Authentication) error {
	allData, signature, err := splitTransferValue(record)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("failed to split transfer value: %w", err)
	}

	// Gen1 uses RSA-1024 with 128-byte signatures
	const lenRsaSignature = 128
	if len(signature) != lenRsaSignature {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("invalid signature length for Gen1: got %d, want %d", len(signature), lenRsaSignature)
	}

	// For Overview (TREP 01), Technical Data (TREP 05), and other Gen1 records,
	// the signature is over the data AFTER the certificates.
	// Each RSA certificate is 194 bytes: MSCA cert + VU cert = 388 bytes total.
	var signedData []byte
	transferType := record.GetType()

	switch transferType {
	case vuv1.TransferType_OVERVIEW_GEN1:
		// Overview: Signature covers data from VehicleIdentificationNumber onwards (after 2 certificates)
		const lenCertificates = 388 // 194 bytes MSCA + 194 bytes VU
		if len(allData) <= lenCertificates {
			auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
			return fmt.Errorf("insufficient data for Overview signature verification: got %d, need > %d", len(allData), lenCertificates)
		}
		signedData = allData[lenCertificates:]

	case vuv1.TransferType_TECHNICAL_DATA_GEN1:
		// Technical Data: Signature covers all data (no certificates in this record)
		signedData = allData

	case vuv1.TransferType_ACTIVITIES_GEN1:
		// Activities: Signature covers all data (no certificates in this record)
		signedData = allData

	case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
		// Events/Faults: Signature covers all data (no certificates in this record)
		signedData = allData

	case vuv1.TransferType_DETAILED_SPEED_GEN1:
		// Detailed Speed: Signature covers all data (no certificates in this record)
		signedData = allData

	default:
		// For unknown transfer types, assume signature covers all data
		signedData = allData
	}

	// Verify the signature using PKCS#1 v1.5
	if err := security.VerifyRsaDataSignature(signedData, signature, vuCert); err != nil {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("data signature verification failed: %w", err)
	}

	return nil
}
