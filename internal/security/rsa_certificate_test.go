package security

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestUnmarshalRsaCertificate tests parsing of RSA certificates from Generation 1.
//
// RSA certificates are 194 bytes with the following structure:
//
//	Bytes 0-127:   Signature (with message recovery, ISO/IEC 9796-2)
//	Bytes 128-185: Cn' (non-recoverable certificate content)
//	Bytes 186-193: CAR' (Certificate Authority Reference)
//
// Without signature verification, we can only extract the CAR field.
//
// See Appendix 11, Part A, Section 3.3.2 (CSM_018) "Certificates issued".
func TestUnmarshalRsaCertificate(t *testing.T) {
	// Test cases using real MSCA certificates from Finland (latest, expires 2031)
	tests := []struct {
		name        string
		filename    string // Certificate file name
		expectedCAR string // Certificate Authority Reference (ERCA root)
		expectedCHR string // Certificate Holder Reference
	}{
		{
			name:        "Finland TCC37",
			filename:    "testdata/certs/g1/finland_tcc37.bin",
			expectedCAR: "18250066869723594497", // ERCA root key ID
			expectedCHR: "1316820541096591105",
		},
		{
			name:        "Finland TCC38",
			filename:    "testdata/certs/g1/finland_tcc38.bin",
			expectedCAR: "18250066869723594497",
			expectedCHR: "1316820541113368321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read certificate file
			data, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("Failed to read certificate file %s: %v", tt.filename, err)
			}

			// Validate file size
			const lenRsaCertificate = 194
			if len(data) != lenRsaCertificate {
				t.Fatalf("Certificate file size = %d, want %d", len(data), lenRsaCertificate)
			}

			// Parse the certificate
			cert, err := UnmarshalRsaCertificate(data)
			if err != nil {
				t.Fatalf("UnmarshalRsaCertificate() failed: %v", err)
			}

			// Validate the certificate
			if cert == nil {
				t.Fatal("UnmarshalRsaCertificate() returned nil")
			}

			// Validate CAR (bytes 186-193)
			car := cert.GetCertificateAuthorityReference()
			if car != tt.expectedCAR {
				t.Errorf("CertificateAuthorityReference = %q, want %q", car, tt.expectedCAR)
			}

			// Note: We can't validate CHR here because it requires signature verification
			// CHR will be validated in TestVerifyRsaCertificateWithRoot

			// Validate raw data is preserved
			rawData := cert.GetRawData()
			if diff := cmp.Diff(data, rawData); diff != "" {
				t.Errorf("Raw data mismatch (-want +got):\n%s", diff)
			}

			// Note: CHR, EOV, modulus, and exponent cannot be validated here
			// because they require signature verification with the CA's public key.
			// These fields will be populated during VerifyRsaCertificateWithRoot test.
		})
	}
}

// TestUnmarshalRsaCertificate_RoundTrip tests that marshalling and
// unmarshalling an RSA certificate produces identical results.
func TestUnmarshalRsaCertificate_RoundTrip(t *testing.T) {
	// Use Finland TCC37 certificate
	originalData, err := os.ReadFile("testdata/certs/g1/finland_tcc37.bin")
	if err != nil {
		t.Fatalf("Failed to read certificate: %v", err)
	}

	// First unmarshal
	cert1, err := UnmarshalRsaCertificate(originalData)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal back to binary
	marshaledData, err := AppendRsaCertificate(nil, cert1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Assert binary equality
	if diff := cmp.Diff(originalData, marshaledData); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	cert2, err := UnmarshalRsaCertificate(marshaledData)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Assert structural equality
	if diff := cmp.Diff(cert1.GetCertificateAuthorityReference(), cert2.GetCertificateAuthorityReference()); diff != "" {
		t.Errorf("CAR mismatch after round-trip (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cert1.GetRawData(), cert2.GetRawData()); diff != "" {
		t.Errorf("RawData mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestUnmarshalRsaCertificate_InvalidLength tests that parsing fails
// when the input data has an incorrect length.
func TestUnmarshalRsaCertificate_InvalidLength(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "too short",
			data: make([]byte, 193),
		},
		{
			name: "too long",
			data: make([]byte, 195),
		},
		{
			name: "empty",
			data: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalRsaCertificate(tt.data)
			if err == nil {
				t.Error("UnmarshalRsaCertificate() succeeded, want error")
			}
		})
	}
}
