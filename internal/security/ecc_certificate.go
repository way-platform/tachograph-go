package security

import (
	"encoding/asn1"
	"encoding/binary"
	"fmt"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// UnmarshalEccCertificate parses an ECC certificate from Generation 2 tachograph cards.
//
// The certificate format is specified in Appendix 11, Section 9.3.2 (PART B), Table 4.
// It uses ASN.1 DER encoding with the following structure:
//
// Certificate Structure (204-341 bytes variable):
//
//	SEQUENCE {
//	  SEQUENCE (Certificate Body) {
//	    [APPLICATION 41] CertificateProfileIdentifier (1 byte)
//	    [APPLICATION 42] CertificationAuthorityReference (8 bytes)
//	    [APPLICATION 76] CertificateHolderAuthorisation (7 bytes)
//	    [APPLICATION 73] PublicKey {
//	      OBJECT IDENTIFIER (Domain Parameters OID)
//	      OCTET STRING (Public Point - uncompressed EC point: 04 || X || Y)
//	    }
//	    [APPLICATION 32] CertificateHolderReference (8 bytes)
//	    [APPLICATION 37] CertificateEffectiveDate (4 bytes, TimeReal)
//	    [APPLICATION 36] CertificateExpirationDate (4 bytes, TimeReal)
//	  }
//	  [APPLICATION 55] ECCCertificateSignature (variable, R || S in plain format)
//	}
//
// See Appendix 11, Section 9.3.2 for the complete specification.
func UnmarshalEccCertificate(data []byte) (*securityv1.EccCertificate, error) {
	const (
		minLenEccCertificate = 204
		maxLenEccCertificate = 341
	)

	if len(data) < minLenEccCertificate || len(data) > maxLenEccCertificate {
		return nil, fmt.Errorf("invalid data length for EccCertificate: got %d, want %d-%d", len(data), minLenEccCertificate, maxLenEccCertificate)
	}

	cert := &securityv1.EccCertificate{}
	cert.SetRawData(data)

	// Parse outer SEQUENCE (the certificate wrapper)
	var outerSeq asn1.RawValue
	_, err := asn1.Unmarshal(data, &outerSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate outer SEQUENCE: %w", err)
	}

	// Parse certificate body SEQUENCE
	var bodySeq asn1.RawValue
	restAfterBody, err := asn1.Unmarshal(outerSeq.Bytes, &bodySeq)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate body SEQUENCE: %w", err)
	}

	// Parse Certificate Profile Identifier (CPI)
	var cpiRaw asn1.RawValue
	restAfterCPI, err := asn1.Unmarshal(bodySeq.Bytes, &cpiRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPI: %w", err)
	}
	if len(cpiRaw.Bytes) != 1 {
		return nil, fmt.Errorf("invalid CPI length: got %d, want 1", len(cpiRaw.Bytes))
	}
	cert.SetCertificateProfileIdentifier(int32(cpiRaw.Bytes[0]))

	// Parse Certificate Authority Reference (CAR)
	var carRaw asn1.RawValue
	restAfterCAR, err := asn1.Unmarshal(restAfterCPI, &carRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CAR: %w", err)
	}
	if len(carRaw.Bytes) != 8 {
		return nil, fmt.Errorf("invalid CAR length: got %d, want 8", len(carRaw.Bytes))
	}
	car := binary.BigEndian.Uint64(carRaw.Bytes)
	carStr := fmt.Sprintf("%d", car)
	cert.SetCertificateAuthorityReference(carStr)

	// Parse Certificate Holder Authorisation (CHA)
	var chaRaw asn1.RawValue
	restAfterCHA, err := asn1.Unmarshal(restAfterCAR, &chaRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CHA: %w", err)
	}
	if len(chaRaw.Bytes) != 7 {
		return nil, fmt.Errorf("invalid CHA length: got %d, want 7", len(chaRaw.Bytes))
	}
	cert.SetCertificateHolderAuthorisation(chaRaw.Bytes)

	// Parse Public Key SEQUENCE
	var pkSeq asn1.RawValue
	restAfterPK, err := asn1.Unmarshal(restAfterCHA, &pkSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key SEQUENCE: %w", err)
	}

	// Parse Domain Parameters (OID)
	var dpOID asn1.ObjectIdentifier
	restAfterDP, err := asn1.Unmarshal(pkSeq.Bytes, &dpOID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse domain parameters OID: %w", err)
	}

	// Parse Public Point (uncompressed EC point: 04 || X || Y)
	var ppRaw asn1.RawValue
	_, err = asn1.Unmarshal(restAfterDP, &ppRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public point: %w", err)
	}

	// Parse the uncompressed point (first byte should be 0x04)
	if len(ppRaw.Bytes) < 1 || ppRaw.Bytes[0] != 0x04 {
		return nil, fmt.Errorf("invalid public point format: expected uncompressed point (0x04)")
	}

	// The remaining bytes are X || Y, each of equal length
	coordLen := (len(ppRaw.Bytes) - 1) / 2
	if len(ppRaw.Bytes) != 1+2*coordLen {
		return nil, fmt.Errorf("invalid public point length: got %d bytes", len(ppRaw.Bytes))
	}

	publicKey := &securityv1.EccCertificate_PublicKey{}
	publicKey.SetDomainParametersOid(dpOID.String())
	publicKey.SetPublicPointX(ppRaw.Bytes[1 : 1+coordLen])
	publicKey.SetPublicPointY(ppRaw.Bytes[1+coordLen:])
	cert.SetPublicKey(publicKey)

	// Parse Certificate Holder Reference (CHR)
	var chrRaw asn1.RawValue
	restAfterCHR, err := asn1.Unmarshal(restAfterPK, &chrRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CHR: %w", err)
	}
	if len(chrRaw.Bytes) != 8 {
		return nil, fmt.Errorf("invalid CHR length: got %d, want 8", len(chrRaw.Bytes))
	}
	chr := binary.BigEndian.Uint64(chrRaw.Bytes)
	chrStr := fmt.Sprintf("%d", chr)
	cert.SetCertificateHolderReference(chrStr)

	// Parse Certificate Effective Date (CEfD)
	var cefdRaw asn1.RawValue
	restAfterCEfD, err := asn1.Unmarshal(restAfterCHR, &cefdRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CEfD: %w", err)
	}
	if len(cefdRaw.Bytes) != 4 {
		return nil, fmt.Errorf("invalid CEfD length: got %d, want 4", len(cefdRaw.Bytes))
	}
	cefd, err := unmarshalTimeReal(cefdRaw.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CEfD timestamp: %w", err)
	}
	cert.SetCertificateEffectiveDate(cefd)

	// Parse Certificate Expiration Date (CExD)
	var cexdRaw asn1.RawValue
	_, err = asn1.Unmarshal(restAfterCEfD, &cexdRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CExD: %w", err)
	}
	if len(cexdRaw.Bytes) != 4 {
		return nil, fmt.Errorf("invalid CExD length: got %d, want 4", len(cexdRaw.Bytes))
	}
	cexd, err := unmarshalTimeReal(cexdRaw.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CExD timestamp: %w", err)
	}
	cert.SetCertificateExpirationDate(cexd)

	// Parse signature from rest after body
	var sigRaw asn1.RawValue
	_, err = asn1.Unmarshal(restAfterBody, &sigRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	// Signature is in plain format: R || S, each of equal length
	if len(sigRaw.Bytes)%2 != 0 {
		return nil, fmt.Errorf("invalid signature length: got %d bytes (must be even)", len(sigRaw.Bytes))
	}
	sigLen := len(sigRaw.Bytes) / 2

	signature := &securityv1.EccCertificate_EccSignature{}
	signature.SetR(sigRaw.Bytes[:sigLen])
	signature.SetS(sigRaw.Bytes[sigLen:])
	cert.SetSignature(signature)

	// Note: Signature verification is not performed during parsing.
	// The signature_valid field is left unset (false).
	// Signature verification should be performed separately when needed,
	// which requires CA certificate lookup and ECDSA verification.

	return cert, nil
}

// AppendEccCertificate marshals an ECC certificate to binary format.
//
// This function uses the raw data painting strategy: if raw_data is available,
// it is used as-is since the certificate is already in ASN.1 DER format.
//
// Reconstructing a certificate from semantic fields would require re-encoding
// the ASN.1 structure and re-signing with the CA's private key, which is not
// typically needed for parsing/marshalling existing card data.
//
// See Appendix 11, Section 9.3.2 for the certificate format specification.
func AppendEccCertificate(dst []byte, cert *securityv1.EccCertificate) ([]byte, error) {
	const (
		minLenEccCertificate = 204
		maxLenEccCertificate = 341
	)

	if cert == nil {
		return nil, fmt.Errorf("EccCertificate cannot be nil")
	}

	// Use raw_data if available (raw data painting strategy)
	if rawData := cert.GetRawData(); len(rawData) > 0 {
		if len(rawData) < minLenEccCertificate || len(rawData) > maxLenEccCertificate {
			return nil, fmt.Errorf("invalid raw_data length for EccCertificate: got %d, want %d-%d", len(rawData), minLenEccCertificate, maxLenEccCertificate)
		}
		return append(dst, rawData...), nil
	}

	// If no raw_data, we would need to construct the certificate from semantic fields
	// and sign it, which requires CA private key access. This is not typically needed
	// for parsing/marshalling existing card data.
	return nil, fmt.Errorf("cannot marshal EccCertificate without raw_data (certificate signing requires CA private key)")
}
