package security

import (
	"encoding/binary"
	"fmt"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// UnmarshalRsaCertificate parses an RSA certificate from Generation 1 tachograph cards.
//
// The certificate format is specified in Appendix 11, Section 3.3 (PART A).
// It uses ISO/IEC 9796-2 digital signature scheme with partial message recovery.
//
// Certificate Structure (194 bytes fixed):
//
//	Bytes 0-127:   Sr (Digital signature with partial message recovery)
//	Bytes 128-185: Cn' (non-recoverable part of certificate content)
//	Bytes 186-193: CAR' (Certificate Authority Reference)
//
// The recovered message Sr contains:
//
//	Byte 0:        Header (0x6A)
//	Bytes 1-106:   Cr' (recoverable part of certificate content)
//	Bytes 107-126: H' (SHA-1 hash of C' = Cr' || Cn')
//	Byte 127:      Trailer (0xBC)
//
// The complete certificate content C' = Cr' || Cn' (164 bytes) contains:
//
//	Byte 0:        CPI (Certificate Profile Identifier, 0x01)
//	Bytes 1-8:     CAR (Certification Authority Reference)
//	Bytes 9-15:    CHA (Certificate Holder Authorisation)
//	Bytes 16-19:   EOV (End Of Validity, TimeReal, or 0xFFFFFFFF)
//	Bytes 20-27:   CHR (Certificate Holder Reference)
//	Bytes 28-155:  n (RSA modulus, 128 bytes)
//	Bytes 156-163: e (RSA exponent, 8 bytes)
//
// Note: Full extraction of CHR, EOV, modulus, and exponent requires signature
// recovery using the CA's RSA public key, which is not performed here.
// Only the CAR can be extracted without signature verification.
//
// See Appendix 11, Section 3.3 for the complete certificate format specification.
func UnmarshalRsaCertificate(data []byte) (*securityv1.RsaCertificate, error) {
	const (
		lenRsaCertificate = 194
		idxCAR            = 186
	)

	if len(data) != lenRsaCertificate {
		return nil, fmt.Errorf("invalid data length for RsaCertificate: got %d, want %d", len(data), lenRsaCertificate)
	}

	// Extract CAR' (Certificate Authority Reference) from bytes 186-193
	// This can be extracted without signature verification
	car := binary.BigEndian.Uint64(data[idxCAR : idxCAR+8])
	carStr := fmt.Sprintf("%d", car)

	// Note: CHR, EOV, RSA modulus, and RSA exponent can only be extracted
	// after signature recovery, which requires the CA's public key.
	// Signature verification is intentionally not performed during parsing.

	cert := &securityv1.RsaCertificate{}
	cert.SetCertificateAuthorityReference(carStr)
	cert.SetRawData(data)

	return cert, nil
}

// AppendRsaCertificate marshals an RSA certificate to binary format.
//
// This function uses the raw data painting strategy: if raw_data is available,
// it is used as-is. Otherwise, the certificate would need to be reconstructed
// from semantic fields (CHR, CAR, EOV, modulus, exponent) and signed, which
// requires private key access.
//
// See Appendix 11, Section 3.3 for the certificate format specification.
func AppendRsaCertificate(dst []byte, cert *securityv1.RsaCertificate) ([]byte, error) {
	const lenRsaCertificate = 194

	if cert == nil {
		return nil, fmt.Errorf("RsaCertificate cannot be nil")
	}

	// Use raw_data if available (raw data painting strategy)
	if rawData := cert.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenRsaCertificate {
			return nil, fmt.Errorf("invalid raw_data length for RsaCertificate: got %d, want %d", len(rawData), lenRsaCertificate)
		}
		return append(dst, rawData...), nil
	}

	// If no raw_data, we would need to construct the certificate from semantic fields
	// and sign it, which requires CA private key access. This is not typically needed
	// for parsing/marshalling existing card data.
	return nil, fmt.Errorf("cannot marshal RsaCertificate without raw_data (certificate signing requires CA private key)")
}
