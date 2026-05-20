package security

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"math/big"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// VerifyRsaDataSignature verifies an RSA signature on data using PKCS#1 v1.5 with SHA-1.
//
// This function implements the signature verification for Generation 1 VU data downloads
// as specified in Appendix 11, Section 6 (CSM_034):
//
//	Signature = EQT.SK['00' || '01' || PS || '00' || DER(SHA-1(Data))]
//
// Where:
//   - PS = Padding string of octets with value 'FF' such that total length is 128 bytes
//   - DER(SHA-1(M)) is the ASN.1 DigestInfo encoding
//
// The certificate must have been verified (signature_valid = true) and contain
// the RSA public key components (modulus and exponent) needed for verification.
//
// See Appendix 11, Section 6.1 for the complete specification.
func VerifyRsaDataSignature(data []byte, signature []byte, cert *securityv1.RsaCertificate) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	// Ensure the certificate has been verified and contains public key
	if !cert.GetSignatureValid() {
		return fmt.Errorf("certificate has not been verified (signature_valid = false)")
	}

	modulus := cert.GetRsaModulus()
	exponent := cert.GetRsaExponent()

	if len(modulus) == 0 || len(exponent) == 0 {
		return fmt.Errorf("certificate does not contain RSA public key")
	}

	// Construct the RSA public key
	modulusInt := new(big.Int).SetBytes(modulus)
	exponentInt := new(big.Int).SetBytes(exponent)

	// Standard RSA public key format expects exponent as int
	// For tachograph certs, exponent is typically 65537 (0x00010001)
	if exponentInt.BitLen() > 31 {
		return fmt.Errorf("RSA exponent too large: %d bits", exponentInt.BitLen())
	}
	expInt := int(exponentInt.Int64())

	pubKey := &rsa.PublicKey{
		N: modulusInt,
		E: expInt,
	}

	// SHA-1 is mandated by the EU tachograph standard (EC 2016/799 Annex 1C) for Gen1 devices.
	hash := sha1.Sum(data) //nolint:gosec // nosemgrep: go.lang.security.audit.crypto.use_of_weak_crypto.use-of-sha1

	// Verify the signature using PKCS#1 v1.5
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hash[:], signature); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
