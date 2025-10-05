package security

import (
	"encoding/binary"
	"fmt"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// UnmarshalRootCertificate parses the European Root Certificate Authority (ERCA)
// public key from its 144-byte binary format.
//
// The ERCA root certificate is distributed in a special format that contains
// the RSA public key components directly, without a signature. This root
// certificate is trusted a priori and serves as the trust anchor for all Member
// State certificates in the tachograph PKI hierarchy.
//
// See Appendix 11, Section 2.1 "European Root Public Key".
//
// Binary Structure (144 bytes):
//
//	Bytes 0-7:     Key Identifier (serves as both CAR and CHR)
//	Bytes 8-135:   RSA Modulus (128 bytes)
//	Bytes 136-143: RSA Public Exponent (8 bytes)
func UnmarshalRootCertificate(data []byte) (*securityv1.RootCertificate, error) {
	const (
		lenRootCert = 144
		idxKeyID    = 0
		idxModulus  = 8
		idxExponent = 136
		lenKeyID    = 8
		lenModulus  = 128
		lenExponent = 8
	)

	if len(data) != lenRootCert {
		return nil, fmt.Errorf("invalid ERCA root certificate size: got %d, want %d", len(data), lenRootCert)
	}

	// Extract key ID (used as both CAR and CHR) and convert to decimal string
	keyID := binary.BigEndian.Uint64(data[idxKeyID : idxKeyID+lenKeyID])
	keyIDStr := fmt.Sprintf("%d", keyID)

	// Extract RSA public key components
	modulus := data[idxModulus : idxModulus+lenModulus]
	exponent := data[idxExponent : idxExponent+lenExponent]

	cert := &securityv1.RootCertificate{}
	cert.SetKeyId(keyIDStr)
	cert.SetRsaModulus(modulus)
	cert.SetRsaExponent(exponent)

	return cert, nil
}

// AppendRootCertificate marshals a RootCertificate to binary format.
//
// This function is provided for API completeness, but root certificates are not
// typically included in tachograph files. Root certificates are distributed
// separately and loaded from the embedded certificate cache.
//
// See Appendix 11, Section 2.1 for the root certificate format specification.
func AppendRootCertificate(dst []byte, root *securityv1.RootCertificate) ([]byte, error) {
	const (
		lenRootCert = 144
		lenKeyID    = 8
		lenModulus  = 128
		lenExponent = 8
	)

	if root == nil {
		return nil, fmt.Errorf("RootCertificate cannot be nil")
	}

	// Parse key ID from decimal string
	var keyID uint64
	if _, err := fmt.Sscanf(root.GetKeyId(), "%d", &keyID); err != nil {
		return nil, fmt.Errorf("invalid key_id format: %w", err)
	}

	modulus := root.GetRsaModulus()
	exponent := root.GetRsaExponent()

	// Validate lengths
	if len(modulus) != lenModulus {
		return nil, fmt.Errorf("invalid rsa_modulus length: got %d, want %d", len(modulus), lenModulus)
	}
	if len(exponent) != lenExponent {
		return nil, fmt.Errorf("invalid rsa_exponent length: got %d, want %d", len(exponent), lenExponent)
	}

	// Allocate buffer
	var buf [lenRootCert]byte

	// Encode key ID (8 bytes, big-endian)
	binary.BigEndian.PutUint64(buf[0:8], keyID)

	// Copy RSA modulus (128 bytes)
	copy(buf[8:136], modulus)

	// Copy RSA public exponent (8 bytes)
	copy(buf[136:144], exponent)

	return append(dst, buf[:]...), nil
}
