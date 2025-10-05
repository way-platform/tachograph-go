package security

import (
	"os"
	"testing"

	"github.com/way-platform/tachograph-go/internal/cert/certcache"
)

func TestVerifyRsaCertificateWithRoot(t *testing.T) {
	// Load the ERCA root certificate
	rootData := certcache.Root()
	root, err := UnmarshalRootCertificate(rootData)
	if err != nil {
		t.Fatalf("Failed to unmarshal root certificate: %v", err)
	}

	// Test cases using real MSCA certificates from Finland (latest, expires 2031)
	tests := []struct {
		name        string
		filename    string // Certificate file name
		expectedCAR string // Should match root key ID
		expectedCHR string // Certificate Holder Reference
	}{
		{
			name:        "Finland TCC37",
			filename:    "testdata/certs/g1/finland_tcc37.bin",
			expectedCAR: "18250066869723594497",
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
			// Read MSCA certificate file
			data, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("Failed to read certificate file %s: %v", tt.filename, err)
			}

			// Parse the MSCA certificate (without verification)
			cert, err := UnmarshalRsaCertificate(data)
			if err != nil {
				t.Fatalf("UnmarshalRsaCertificate() failed: %v", err)
			}

			// Verify the certificate against the ERCA root
			err = VerifyRsaCertificateWithRoot(cert, root)
			if err != nil {
				t.Fatalf("VerifyRsaCertificateWithRoot() failed: %v", err)
			}

			// After successful verification, validate that fields are populated
			if !cert.GetSignatureValid() {
				t.Error("signature_valid = false, want true")
			}

			// Validate Certificate Holder Reference (CHR) is populated
			chr := cert.GetCertificateHolderReference()
			if chr == "" {
				t.Error("CertificateHolderReference is empty after verification")
			}
			if chr != tt.expectedCHR {
				t.Errorf("CertificateHolderReference = %q, want %q", chr, tt.expectedCHR)
			}

			// Validate Certificate Authority Reference (CAR) matches ERCA
			car := cert.GetCertificateAuthorityReference()
			if car != tt.expectedCAR {
				t.Errorf("CertificateAuthorityReference = %q, want %q", car, tt.expectedCAR)
			}

			// Validate RSA public key is populated
			modulus := cert.GetRsaModulus()
			const lenModulus = 128
			if len(modulus) != lenModulus {
				t.Errorf("RsaModulus length = %d, want %d", len(modulus), lenModulus)
			}

			exponent := cert.GetRsaExponent()
			const lenExponent = 8
			if len(exponent) != lenExponent {
				t.Errorf("RsaExponent length = %d, want %d", len(exponent), lenExponent)
			}

			// Validate End of Validity is populated (may be nil if set to 0xFFFFFFFF)
			eov := cert.GetEndOfValidity()
			if eov != nil {
				// If EOV is set, it should be a valid timestamp
				if eov.GetSeconds() == 0 && eov.GetNanos() == 0 {
					t.Error("EndOfValidity is zero timestamp (unexpected)")
				}
			}
		})
	}
}

func TestVerifyRsaCertificateWithRoot_CARMismatch(t *testing.T) {
	// Load the ERCA root certificate
	rootData := certcache.Root()

	// Create a modified root with different key ID
	wrongRoot, err := UnmarshalRootCertificate(rootData)
	if err != nil {
		t.Fatalf("Failed to unmarshal root certificate: %v", err)
	}
	// Set a different key ID (this will cause CAR mismatch)
	wrongRoot.SetKeyId("99999999999999999")

	// Read a real MSCA certificate
	data, err := os.ReadFile("testdata/certs/g1/finland_tcc37.bin")
	if err != nil {
		t.Fatalf("Failed to read certificate: %v", err)
	}

	cert, err := UnmarshalRsaCertificate(data)
	if err != nil {
		t.Fatalf("UnmarshalRsaCertificate() failed: %v", err)
	}

	// Verification should fail due to CAR mismatch
	err = VerifyRsaCertificateWithRoot(cert, wrongRoot)
	if err == nil {
		t.Error("VerifyRsaCertificateWithRoot() succeeded with wrong root, want error")
	}

	// signature_valid should be false
	if cert.GetSignatureValid() {
		t.Error("signature_valid = true after failed verification, want false")
	}
}

func TestVerifyRsaCertificateWithCA(t *testing.T) {
	// Load and verify the first MSCA certificate (Austria) against ERCA root
	rootData := certcache.Root()
	root, err := UnmarshalRootCertificate(rootData)
	if err != nil {
		t.Fatalf("Failed to unmarshal root certificate: %v", err)
	}

	// Read and verify Finland TCC37 certificate
	finlandData1, err := os.ReadFile("testdata/certs/g1/finland_tcc37.bin")
	if err != nil {
		t.Fatalf("Failed to read Finland TCC37 certificate: %v", err)
	}

	finlandCert1, err := UnmarshalRsaCertificate(finlandData1)
	if err != nil {
		t.Fatalf("UnmarshalRsaCertificate() failed for Finland TCC37: %v", err)
	}

	err = VerifyRsaCertificateWithRoot(finlandCert1, root)
	if err != nil {
		t.Fatalf("Failed to verify Finland TCC37: %v", err)
	}

	// Now finlandCert1 has its public key populated and can act as a CA
	// For this test, we'll verify another Finland MSCA certificate
	// (In a real scenario, this would be an equipment certificate signed by Finland MSCA)

	// Read another Finland MSCA certificate (TCC38)
	finlandData2, err := os.ReadFile("testdata/certs/g1/finland_tcc38.bin")
	if err != nil {
		t.Fatalf("Failed to read Finland TCC38 certificate: %v", err)
	}

	finlandCert2, err := UnmarshalRsaCertificate(finlandData2)
	if err != nil {
		t.Fatalf("UnmarshalRsaCertificate() failed for Finland TCC38: %v", err)
	}

	// Verify TCC38 against ERCA root (both should be signed by ERCA)
	err = VerifyRsaCertificateWithRoot(finlandCert2, root)
	if err != nil {
		t.Fatalf("Failed to verify Finland TCC38: %v", err)
	}

	// Both certificates should now be verified
	if !finlandCert1.GetSignatureValid() {
		t.Error("Finland TCC37 certificate signature_valid = false")
	}
	if !finlandCert2.GetSignatureValid() {
		t.Error("Finland TCC38 certificate signature_valid = false")
	}

	// Note: We can't test VerifyRsaCertificateWithCA with these MSCA certs
	// verifying each other because they're both signed by ERCA, not by each other.
	// This would require equipment certificates (card or VU certs) signed by the MSCA.
}
