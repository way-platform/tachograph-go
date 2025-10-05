package security

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalEccCertificate(t *testing.T) {
	// Test cases using real MSCA Gen2 certificates from Finland (latest, expires 2031)
	tests := []struct {
		name        string
		filename    string // Certificate file name
		expectedCAR string // Certificate Authority Reference (Gen2 ERCA root)
		expectedCHR string // Certificate Holder Reference
		minKeySize  int    // Minimum expected key size in bytes
	}{
		{
			name:        "Finland MSCA Card42",
			filename:    "testdata/certs/g2/finland_msca_card42.bin",
			expectedCAR: "18250066869740371713", // Gen2 ERCA root
			expectedCHR: "1316820541130145537",
			minKeySize:  32, // P-256 or BrainpoolP256r1
		},
		{
			name:        "Finland MSCA Card43",
			filename:    "testdata/certs/g2/finland_msca_card43.bin",
			expectedCAR: "18250066869740371713",
			expectedCHR: "1316820541146922753",
			minKeySize:  32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read certificate file
			data, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Skipf("Certificate file not found %s: %v", tt.filename, err)
			}

			// Validate file size is within expected range
			const minLen = 204
			const maxLen = 341
			if len(data) < minLen || len(data) > maxLen {
				t.Fatalf("Certificate file size = %d, want %d-%d", len(data), minLen, maxLen)
			}

			// Parse the certificate
			cert, err := UnmarshalEccCertificate(data)
			if err != nil {
				t.Fatalf("UnmarshalEccCertificate() failed: %v", err)
			}

			// Validate the certificate
			if cert == nil {
				t.Fatal("UnmarshalEccCertificate() returned nil")
			}

			// Validate Certificate Profile Identifier
			cpi := cert.GetCertificateProfileIdentifier()
			const expectedCPI = 0 // Version 1 profile
			if cpi != expectedCPI {
				t.Errorf("CertificateProfileIdentifier = %d, want %d", cpi, expectedCPI)
			}

			// Validate Certificate Authority Reference
			car := cert.GetCertificateAuthorityReference()
			if car != tt.expectedCAR {
				t.Errorf("CertificateAuthorityReference = %q, want %q", car, tt.expectedCAR)
			}

			// Validate Certificate Holder Authorisation is present
			cha := cert.GetCertificateHolderAuthorisation()
			const lenCHA = 7
			if len(cha) != lenCHA {
				t.Errorf("CertificateHolderAuthorisation length = %d, want %d", len(cha), lenCHA)
			}

			// Validate Public Key structure
			pubKey := cert.GetPublicKey()
			if pubKey == nil {
				t.Fatal("PublicKey is nil")
			}

			// Validate Domain Parameters OID is present
			oid := pubKey.GetDomainParametersOid()
			if oid == "" {
				t.Error("DomainParametersOid is empty")
			}

			// Validate public point coordinates are present
			pointX := pubKey.GetPublicPointX()
			pointY := pubKey.GetPublicPointY()
			if len(pointX) < tt.minKeySize {
				t.Errorf("PublicPointX length = %d, want >= %d", len(pointX), tt.minKeySize)
			}
			if len(pointY) < tt.minKeySize {
				t.Errorf("PublicPointY length = %d, want >= %d", len(pointY), tt.minKeySize)
			}
			if len(pointX) != len(pointY) {
				t.Errorf("PublicPoint coordinates length mismatch: X=%d, Y=%d", len(pointX), len(pointY))
			}

			// Validate Certificate Holder Reference
			chr := cert.GetCertificateHolderReference()
			if chr != tt.expectedCHR {
				t.Errorf("CertificateHolderReference = %q, want %q", chr, tt.expectedCHR)
			}

			// Validate timestamps are present
			cefd := cert.GetCertificateEffectiveDate()
			if cefd == nil {
				t.Error("CertificateEffectiveDate is nil")
			} else if cefd.GetSeconds() == 0 {
				t.Error("CertificateEffectiveDate is zero")
			}

			cexd := cert.GetCertificateExpirationDate()
			if cexd == nil {
				t.Error("CertificateExpirationDate is nil")
			} else if cexd.GetSeconds() == 0 {
				t.Error("CertificateExpirationDate is zero")
			}

			// Validate effective date is before expiration date
			if cefd != nil && cexd != nil {
				if cefd.GetSeconds() >= cexd.GetSeconds() {
					t.Errorf("CertificateEffectiveDate >= CertificateExpirationDate")
				}
			}

			// Validate signature is present
			sig := cert.GetSignature()
			if sig == nil {
				t.Fatal("Signature is nil")
			}

			// Validate signature components R and S
			r := sig.GetR()
			s := sig.GetS()
			if len(r) == 0 {
				t.Error("Signature R component is empty")
			}
			if len(s) == 0 {
				t.Error("Signature S component is empty")
			}
			if len(r) != len(s) {
				t.Errorf("Signature components length mismatch: R=%d, S=%d", len(r), len(s))
			}

			// Validate raw data is preserved
			rawData := cert.GetRawData()
			if diff := cmp.Diff(data, rawData); diff != "" {
				t.Errorf("Raw data mismatch (-want +got):\n%s", diff)
			}

			// Note: signature_valid is NOT set during parsing
			// Signature verification is performed separately via VerifyEccCertificateWithCA
			if cert.GetSignatureValid() {
				t.Error("signature_valid = true after parsing (should be false until verified)")
			}
		})
	}
}

func TestUnmarshalEccCertificate_RoundTrip(t *testing.T) {
	// Use Finland MSCA Card42 certificate
	originalData, err := os.ReadFile("testdata/certs/g2/finland_msca_card42.bin")
	if err != nil {
		t.Skipf("Certificate file not found: %v", err)
	}

	// First unmarshal
	cert1, err := UnmarshalEccCertificate(originalData)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal back to binary
	marshaledData, err := AppendEccCertificate(nil, cert1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Assert binary equality
	if diff := cmp.Diff(originalData, marshaledData); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	cert2, err := UnmarshalEccCertificate(marshaledData)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Assert structural equality of key fields
	if diff := cmp.Diff(cert1.GetCertificateProfileIdentifier(), cert2.GetCertificateProfileIdentifier()); diff != "" {
		t.Errorf("CPI mismatch after round-trip (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cert1.GetCertificateAuthorityReference(), cert2.GetCertificateAuthorityReference()); diff != "" {
		t.Errorf("CAR mismatch after round-trip (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cert1.GetCertificateHolderReference(), cert2.GetCertificateHolderReference()); diff != "" {
		t.Errorf("CHR mismatch after round-trip (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cert1.GetPublicKey().GetDomainParametersOid(), cert2.GetPublicKey().GetDomainParametersOid()); diff != "" {
		t.Errorf("OID mismatch after round-trip (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cert1.GetPublicKey().GetPublicPointX(), cert2.GetPublicKey().GetPublicPointX()); diff != "" {
		t.Errorf("PublicPointX mismatch after round-trip (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cert1.GetPublicKey().GetPublicPointY(), cert2.GetPublicKey().GetPublicPointY()); diff != "" {
		t.Errorf("PublicPointY mismatch after round-trip (-want +got):\n%s", diff)
	}
}

func TestUnmarshalEccCertificate_InvalidLength(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "too short",
			data: make([]byte, 203),
		},
		{
			name: "too long",
			data: make([]byte, 342),
		},
		{
			name: "empty",
			data: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalEccCertificate(tt.data)
			if err == nil {
				t.Error("UnmarshalEccCertificate() succeeded, want error")
			}
		})
	}
}
