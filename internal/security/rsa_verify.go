package security

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/big"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// VerifyRsaCertificateWithCA performs signature recovery and verification on an RSA certificate
// using another RSA certificate as the Certificate Authority.
//
// The CA certificate must have its public key components (modulus and exponent) populated,
// typically from a previous verification against a higher-level CA or root certificate.
//
// See Appendix 11, Section 3.3 for the complete specification.
func VerifyRsaCertificateWithCA(cert *securityv1.RsaCertificate, caCert *securityv1.RsaCertificate) error {
	if caCert == nil {
		return fmt.Errorf("CA certificate cannot be nil")
	}

	caModulus := caCert.GetRsaModulus()
	caExponent := caCert.GetRsaExponent()
	caCHR := caCert.GetCertificateHolderReference()

	return verifyRsaCertificate(cert, caModulus, caExponent, caCHR)
}

// VerifyRsaCertificateWithRoot performs signature recovery and verification on an RSA certificate
// using the ERCA root certificate as the Certificate Authority.
//
// The root certificate is trusted a priori and contains the public key components directly.
// This is typically used to verify Member State CA certificates against the European Root CA.
//
// See Appendix 11, Section 2.1 for the root certificate format and Section 3.3 for
// RSA certificate verification.
func VerifyRsaCertificateWithRoot(cert *securityv1.RsaCertificate, root *securityv1.RootCertificate) error {
	if root == nil {
		return fmt.Errorf("root certificate cannot be nil")
	}

	rootModulus := root.GetRsaModulus()
	rootExponent := root.GetRsaExponent()
	rootKeyID := root.GetKeyId()

	return verifyRsaCertificate(cert, rootModulus, rootExponent, rootKeyID)
}

// verifyRsaCertificate is the internal implementation that performs signature recovery
// and verification on an RSA certificate using the provided RSA public key components.
//
// The function implements ISO/IEC 9796-2 digital signature scheme with partial message recovery.
//
// Certificate Structure (194 bytes):
//
//	Bytes 0-127:   Sr (Signature with recoverable message)
//	Bytes 128-185: Cn' (Non-recoverable certificate content)
//	Bytes 186-193: CAR' (Certificate Authority Reference)
//
// Recovered Structure Sr' (128 bytes after RSA operation):
//
//	Byte 0:        Header (0x6A)
//	Bytes 1-106:   Cr' (Recoverable certificate content)
//	Bytes 107-126: H' (SHA-1 hash of C' = Cr' || Cn')
//	Byte 127:      Trailer (0xBC)
//
// Complete Certificate Content C' = Cr' || Cn' (164 bytes):
//
//	Byte 0:        CPI (Certificate Profile Identifier, 0x01)
//	Bytes 1-8:     CAR (Certification Authority Reference)
//	Bytes 9-15:    CHA (Certificate Holder Authorisation)
//	Bytes 16-19:   EOV (End Of Validity, TimeReal)
//	Bytes 20-27:   CHR (Certificate Holder Reference)
//	Bytes 28-155:  n (RSA modulus, 128 bytes)
//	Bytes 156-163: e (RSA public exponent, 8 bytes)
//
// This function mutates the certificate by:
//   - Setting signature_valid to true if verification succeeds
//   - Populating certificate_holder_reference
//   - Populating end_of_validity
//   - Populating rsa_modulus
//   - Populating rsa_exponent
//
// If verification fails, signature_valid is set to false and other fields remain unchanged.
//
// See Appendix 11, Section 3.3 for the complete specification.
func verifyRsaCertificate(cert *securityv1.RsaCertificate, caModulus, caExponent []byte, caCHR string) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	rawData := cert.GetRawData()
	if len(rawData) != 194 {
		return fmt.Errorf("invalid certificate length: got %d, want 194", len(rawData))
	}

	// CA public key must be provided
	if len(caModulus) == 0 || len(caExponent) == 0 {
		return fmt.Errorf("CA public key missing (modulus or exponent empty)")
	}

	// Extract components from certificate
	const (
		idxSignature = 0
		lenSignature = 128
		idxCnPrime   = 128
		lenCnPrime   = 58
		idxCARPrime  = 186
		lenCARPrime  = 8
	)

	signature := rawData[idxSignature : idxSignature+lenSignature]
	cnPrime := rawData[idxCnPrime : idxCnPrime+lenCnPrime]
	carPrime := binary.BigEndian.Uint64(rawData[idxCARPrime : idxCARPrime+lenCARPrime])
	carPrimeStr := fmt.Sprintf("%d", carPrime)

	// Verify that CAR' matches the CA's CHR (or root key ID)
	// (The CAR in the certificate should reference the CA that signed it)
	if carPrimeStr != caCHR {
		cert.SetSignatureValid(false)
		return fmt.Errorf("CAR mismatch: certificate references CAR %s, but CA has CHR %s",
			carPrimeStr, caCHR)
	}

	// Perform RSA signature recovery: signature^e mod n
	sigInt := new(big.Int).SetBytes(signature)
	modulus := new(big.Int).SetBytes(caModulus)
	exponent := new(big.Int).SetBytes(caExponent)

	recovered := new(big.Int).Exp(sigInt, exponent, modulus)
	srPrime := recovered.Bytes()

	// Pad to 128 bytes if necessary (leading zeros may be dropped)
	if len(srPrime) < 128 {
		padded := make([]byte, 128)
		copy(padded[128-len(srPrime):], srPrime)
		srPrime = padded
	}

	// Verify recovered message structure: 0x6A || Cr' || H' || 0xBC
	if len(srPrime) != 128 {
		cert.SetSignatureValid(false)
		return fmt.Errorf("invalid recovered message length: got %d, want 128", len(srPrime))
	}

	if srPrime[0] != 0x6A {
		cert.SetSignatureValid(false)
		return fmt.Errorf("invalid recovered message header: got 0x%02X, want 0x6A", srPrime[0])
	}

	if srPrime[127] != 0xBC {
		cert.SetSignatureValid(false)
		return fmt.Errorf("invalid recovered message trailer: got 0x%02X, want 0xBC", srPrime[127])
	}

	// Extract Cr' (recoverable part) and H' (hash)
	const (
		idxCrPrime = 1
		lenCrPrime = 106
		idxHPrime  = 107
		lenHPrime  = 20
	)

	crPrime := srPrime[idxCrPrime : idxCrPrime+lenCrPrime]
	hPrime := srPrime[idxHPrime : idxHPrime+lenHPrime]

	// Reconstruct complete certificate content: C' = Cr' || Cn'
	cPrime := make([]byte, 0, lenCrPrime+lenCnPrime)
	cPrime = append(cPrime, crPrime...)
	cPrime = append(cPrime, cnPrime...)

	if len(cPrime) != 164 {
		cert.SetSignatureValid(false)
		return fmt.Errorf("invalid certificate content length: got %d, want 164", len(cPrime))
	}

	// Verify hash: SHA-1(C') should equal H'
	hash := sha1.Sum(cPrime)
	if !bytes.Equal(hPrime, hash[:]) {
		cert.SetSignatureValid(false)
		return fmt.Errorf("certificate content hash mismatch")
	}

	// Hash verified! Now extract semantic fields from C'
	const (
		idxCPI      = 0
		idxCAR      = 1
		idxCHA      = 9
		idxEOV      = 16
		idxCHR      = 20
		idxModulus  = 28
		idxExponent = 156
		lenCAR      = 8
		lenEOV      = 4
		lenCHR      = 8
		lenModulus  = 128
		lenExponent = 8
	)

	// Extract CAR
	// Note: CAR inside C' should match CAR' at the end of the certificate,
	// but we don't enforce this check as it doesn't affect signature validity
	car := binary.BigEndian.Uint64(cPrime[idxCAR : idxCAR+lenCAR])
	carStr := fmt.Sprintf("%d", car)

	// Extract CHR
	chr := binary.BigEndian.Uint64(cPrime[idxCHR : idxCHR+lenCHR])
	chrStr := fmt.Sprintf("%d", chr)

	// Extract EOV (End Of Validity)
	// If parsing fails (e.g., 0xFFFFFFFF indicating no expiry), we continue anyway
	eovBytes := cPrime[idxEOV : idxEOV+lenEOV]
	eov, _ := unmarshalTimeReal(eovBytes)

	// Extract RSA public key
	certModulus := cPrime[idxModulus : idxModulus+lenModulus]
	certExponent := cPrime[idxExponent : idxExponent+lenExponent]

	// Signature verification successful! Populate the certificate
	cert.SetSignatureValid(true)
	cert.SetCertificateHolderReference(chrStr)
	cert.SetCertificateAuthorityReference(carStr)
	if eov != nil {
		cert.SetEndOfValidity(eov)
	}
	cert.SetRsaModulus(certModulus)
	cert.SetRsaExponent(certExponent)

	return nil
}
