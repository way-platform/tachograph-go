package dd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/keybase/go-crypto/brainpool"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// VerifyEccCertificate performs ECDSA signature verification on an ECC certificate.
// It uses the CA certificate's public key to verify the digital signature.
//
// The certificate uses ASN.1 DER encoding as defined in Appendix 11, Section 9.3.2 (PART B), Table 4.
//
// The signature is computed over the certificate body (including ASN.1 tag and length)
// using ECDSA with SHA-256, SHA-384, or SHA-512 depending on the curve parameters.
//
// Supported elliptic curves:
//   - brainpoolP256r1 (OID: 1.3.36.3.3.2.8.1.1.7) - uses SHA-256
//   - brainpoolP384r1 (OID: 1.3.36.3.3.2.8.1.1.11) - uses SHA-384
//   - brainpoolP512r1 (OID: 1.3.36.3.3.2.8.1.1.13) - uses SHA-512
//   - NIST P-256 (OID: 1.2.840.10045.3.1.7) - uses SHA-256
//   - NIST P-384 (OID: 1.3.132.0.34) - uses SHA-384
//   - NIST P-521 (OID: 1.3.132.0.35) - uses SHA-512
//
// This function mutates the certificate by setting signature_valid to true
// if verification succeeds, or false if it fails.
//
// See Appendix 11, Section 9.3.2 for the complete specification.
func VerifyEccCertificate(cert *ddv1.EccCertificate, caCert *ddv1.EccCertificate) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	if caCert == nil {
		return fmt.Errorf("CA certificate cannot be nil")
	}

	rawData := cert.GetRawData()
	if len(rawData) < 204 || len(rawData) > 341 {
		return fmt.Errorf("invalid certificate length: got %d, want 204-341", len(rawData))
	}

	// Get CA's public key
	caPubKey := caCert.GetPublicKey()
	if caPubKey == nil {
		return fmt.Errorf("CA certificate has no public key")
	}

	caPointX := caPubKey.GetPublicPointX()
	caPointY := caPubKey.GetPublicPointY()
	if len(caPointX) == 0 || len(caPointY) == 0 {
		return fmt.Errorf("CA certificate public key is incomplete")
	}

	// Parse CA's curve parameters
	hashBits, curve, err := parseCurveOID(caPubKey.GetDomainParametersOid())
	if err != nil {
		return fmt.Errorf("failed to parse CA curve: %w", err)
	}

	// Re-parse the certificate to extract the certificate body with full ASN.1 bytes
	// The signature is computed over the body INCLUDING the ASN.1 tag and length
	var outerSeq asn1.RawValue
	_, err = asn1.Unmarshal(rawData, &outerSeq)
	if err != nil {
		return fmt.Errorf("failed to parse certificate outer SEQUENCE: %w", err)
	}

	var bodySeq asn1.RawValue
	_, err = asn1.Unmarshal(outerSeq.Bytes, &bodySeq)
	if err != nil {
		return fmt.Errorf("failed to parse certificate body SEQUENCE: %w", err)
	}

	// bodySeq.FullBytes contains the complete body including tag, length, and content
	// This is what the signature is computed over
	hashData := bodySeq.FullBytes

	// Compute hash based on curve size
	var hash []byte
	switch hashBits {
	case 256:
		h := sha256.Sum256(hashData)
		hash = h[:]
	case 384:
		h := sha512.Sum384(hashData)
		hash = h[:]
	case 512:
		h := sha512.Sum512(hashData)
		hash = h[:]
	default:
		return fmt.Errorf("unsupported hash size: %d bits", hashBits)
	}

	// Get signature components
	sig := cert.GetSignature()
	if sig == nil {
		cert.SetSignatureValid(false)
		return fmt.Errorf("certificate has no signature")
	}

	r := new(big.Int).SetBytes(sig.GetR())
	s := new(big.Int).SetBytes(sig.GetS())

	// Construct CA's public key
	caX := new(big.Int).SetBytes(caPointX)
	caY := new(big.Int).SetBytes(caPointY)

	caPub := &ecdsa.PublicKey{
		Curve: curve,
		X:     caX,
		Y:     caY,
	}

	// Verify ECDSA signature
	valid := ecdsa.Verify(caPub, hash, r, s)

	cert.SetSignatureValid(valid)

	if !valid {
		return fmt.Errorf("ECDSA signature verification failed")
	}

	return nil
}

// parseCurveOID parses an elliptic curve OID and returns the hash size (in bits)
// and the elliptic.Curve interface.
//
// Supported curves:
//   - brainpoolP256r1, brainpoolP384r1, brainpoolP512r1
//   - NIST P-256, P-384, P-521
func parseCurveOID(oidStr string) (hashBits int, curve elliptic.Curve, err error) {
	// Parse OID string to asn1.ObjectIdentifier
	var oid asn1.ObjectIdentifier
	_, err = asn1.UnmarshalWithParams([]byte(oidStr), &oid, "tag:6") // tag 6 is OBJECT IDENTIFIER
	if err != nil {
		// Try parsing as dot-separated string
		switch oidStr {
		case "1.3.36.3.3.2.8.1.1.7":
			return 256, brainpool.P256r1(), nil
		case "1.3.36.3.3.2.8.1.1.11":
			return 384, brainpool.P384r1(), nil
		case "1.3.36.3.3.2.8.1.1.13":
			return 512, brainpool.P512r1(), nil
		case "1.2.840.10045.3.1.7":
			return 256, elliptic.P256(), nil
		case "1.3.132.0.34":
			return 384, elliptic.P384(), nil
		case "1.3.132.0.35":
			return 512, elliptic.P521(), nil
		default:
			return 0, nil, fmt.Errorf("unknown elliptic curve OID: %s", oidStr)
		}
	}

	// Map OID to curve
	switch {
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 7}):
		return 256, brainpool.P256r1(), nil
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 11}):
		return 384, brainpool.P384r1(), nil
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 13}):
		return 512, brainpool.P512r1(), nil
	case oid.Equal(asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}):
		return 256, elliptic.P256(), nil
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 132, 0, 34}):
		return 384, elliptic.P384(), nil
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 132, 0, 35}):
		return 512, elliptic.P521(), nil
	default:
		return 0, nil, fmt.Errorf("unknown elliptic curve OID: %v", oid)
	}
}
