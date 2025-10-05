package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyEccCertificateWithCA(t *testing.T) {
	// First, we need to find and verify the Gen2 ERCA root certificate
	// The Gen2 ERCA root is self-signed, so we need to identify it
	// Based on index.json, the Gen2 ERCA has CAR = "18250066869740371713"

	// For this test, we'll need to:
	// 1. Find a Gen2 MSCA certificate
	// 2. Find the Gen2 ERCA root certificate (self-signed, CAR == CHR)
	// 3. Verify the MSCA cert against the root

	// Test cases using real Finland MSCA certificates
	tests := []struct {
		name     string
		mscaFile string // MSCA certificate file
		mscaCHR  string // MSCA Certificate Holder Reference
		mscaCAR  string // MSCA's CAR (should reference ERCA root)
		ercaCHR  string // ERCA root CHR (same as CAR for self-signed)
	}{
		{
			name:     "Finland MSCA Card42",
			mscaFile: "testdata/certs/g2/finland_msca_card42.bin",
			mscaCHR:  "1316820541130145537",
			mscaCAR:  "18250066869740371713", // Gen2 ERCA root
			ercaCHR:  "18250066869740371713", // Same for self-signed root
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Find and parse the Gen2 ERCA root certificate
			// The root is self-signed, so CAR == CHR
			ercaPath := filepath.Join("..", "cert", "certcache", "g2", tt.ercaCHR+".bin")
			ercaData, err := os.ReadFile(ercaPath)
			if err != nil {
				t.Skipf("ERCA root certificate not found at %s: %v", ercaPath, err)
			}

			ercaCert, err := UnmarshalEccCertificate(ercaData)
			if err != nil {
				t.Fatalf("Failed to unmarshal ERCA certificate: %v", err)
			}

			// Verify it's self-signed (CAR == CHR)
			if ercaCert.GetCertificateAuthorityReference() != ercaCert.GetCertificateHolderReference() {
				t.Fatalf("ERCA certificate is not self-signed: CAR=%s, CHR=%s",
					ercaCert.GetCertificateAuthorityReference(),
					ercaCert.GetCertificateHolderReference())
			}

			// Verify the ERCA root certificate against itself (self-signed)
			err = VerifyEccCertificateWithCA(ercaCert, ercaCert)
			if err != nil {
				t.Fatalf("Failed to verify self-signed ERCA certificate: %v", err)
			}

			if !ercaCert.GetSignatureValid() {
				t.Fatal("ERCA certificate signature_valid = false after self-verification")
			}

			// Step 2: Parse the MSCA certificate
			mscaData, err := os.ReadFile(tt.mscaFile)
			if err != nil {
				t.Skipf("MSCA certificate not found at %s: %v", tt.mscaFile, err)
			}

			mscaCert, err := UnmarshalEccCertificate(mscaData)
			if err != nil {
				t.Fatalf("Failed to unmarshal MSCA certificate: %v", err)
			}

			// Verify the MSCA's CAR matches the ERCA's CHR
			if mscaCert.GetCertificateAuthorityReference() != tt.mscaCAR {
				t.Errorf("MSCA CAR = %s, want %s",
					mscaCert.GetCertificateAuthorityReference(), tt.mscaCAR)
			}

			// Step 3: Verify the MSCA certificate against the ERCA root
			err = VerifyEccCertificateWithCA(mscaCert, ercaCert)
			if err != nil {
				t.Fatalf("Failed to verify MSCA certificate against ERCA: %v", err)
			}

			// Validate signature_valid is set to true
			if !mscaCert.GetSignatureValid() {
				t.Error("MSCA certificate signature_valid = false after verification, want true")
			}

			// Validate the MSCA certificate fields are still intact
			chr := mscaCert.GetCertificateHolderReference()
			if chr != tt.mscaCHR {
				t.Errorf("MSCA CHR = %s, want %s", chr, tt.mscaCHR)
			}

			// Validate public key is present
			pubKey := mscaCert.GetPublicKey()
			if pubKey == nil {
				t.Fatal("MSCA PublicKey is nil after verification")
			}

			pointX := pubKey.GetPublicPointX()
			pointY := pubKey.GetPublicPointY()
			if len(pointX) == 0 || len(pointY) == 0 {
				t.Error("MSCA public key coordinates are empty after verification")
			}
		})
	}
}

func TestVerifyEccCertificateWithCA_MultipleCountries(t *testing.T) {
	// Load the Gen2 ERCA root certificate
	const ercaCHR = "18250066869740371713"
	ercaPath := filepath.Join("..", "cert", "certcache", "g2", ercaCHR+".bin")
	ercaData, err := os.ReadFile(ercaPath)
	if err != nil {
		t.Skipf("ERCA root certificate not found: %v", err)
	}

	ercaCert, err := UnmarshalEccCertificate(ercaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal ERCA certificate: %v", err)
	}

	// Verify ERCA is self-signed
	err = VerifyEccCertificateWithCA(ercaCert, ercaCert)
	if err != nil {
		t.Fatalf("Failed to verify ERCA: %v", err)
	}

	// List of Finland MSCA certificates to test
	mscaCerts := []struct {
		filename string
		chr      string
	}{
		{"testdata/certs/g2/finland_msca_card42.bin", "1316820541130145537"},
		{"testdata/certs/g2/finland_msca_card43.bin", "1316820541146922753"},
	}

	successCount := 0
	for _, msca := range mscaCerts {
		t.Run(msca.chr, func(t *testing.T) {
			mscaData, err := os.ReadFile(msca.filename)
			if err != nil {
				t.Skipf("Certificate not found: %v", err)
			}

			mscaCert, err := UnmarshalEccCertificate(mscaData)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			// Verify against ERCA root
			err = VerifyEccCertificateWithCA(mscaCert, ercaCert)
			if err != nil {
				t.Errorf("Verification failed: %v", err)
				return
			}

			if !mscaCert.GetSignatureValid() {
				t.Error("signature_valid = false after verification")
				return
			}

			successCount++
		})
	}

	if successCount == 0 {
		t.Error("No MSCA certificates were successfully verified")
	}
}

func TestVerifyEccCertificateWithCA_InvalidSignature(t *testing.T) {
	// Load ERCA root
	const ercaCHR = "18250066869740371713"
	ercaPath := filepath.Join("..", "cert", "certcache", "g2", ercaCHR+".bin")
	ercaData, err := os.ReadFile(ercaPath)
	if err != nil {
		t.Skipf("ERCA root certificate not found: %v", err)
	}

	ercaCert, err := UnmarshalEccCertificate(ercaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal ERCA: %v", err)
	}

	err = VerifyEccCertificateWithCA(ercaCert, ercaCert)
	if err != nil {
		t.Fatalf("Failed to verify ERCA: %v", err)
	}

	// Load an MSCA certificate
	mscaData, err := os.ReadFile("testdata/certs/g2/finland_msca_card42.bin")
	if err != nil {
		t.Skipf("MSCA certificate not found: %v", err)
	}

	mscaCert, err := UnmarshalEccCertificate(mscaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal MSCA: %v", err)
	}

	// Corrupt the signature by flipping a bit
	sig := mscaCert.GetSignature()
	if sig != nil && len(sig.GetR()) > 0 {
		corruptedR := make([]byte, len(sig.GetR()))
		copy(corruptedR, sig.GetR())
		corruptedR[0] ^= 0x01 // Flip a bit
		sig.SetR(corruptedR)
	}

	// Verification should fail
	err = VerifyEccCertificateWithCA(mscaCert, ercaCert)
	if err == nil {
		t.Error("VerifyEccCertificateWithCA() succeeded with corrupted signature, want error")
	}

	// signature_valid should be false
	if mscaCert.GetSignatureValid() {
		t.Error("signature_valid = true after failed verification, want false")
	}
}

func TestVerifyEccCertificateWithCA_CARMismatch(t *testing.T) {
	// Load two Finland certificates
	cert1Data, err := os.ReadFile("testdata/certs/g2/finland_msca_card42.bin")
	if err != nil {
		t.Skipf("Certificate not found: %v", err)
	}

	cert1, err := UnmarshalEccCertificate(cert1Data)
	if err != nil {
		t.Fatalf("Failed to unmarshal cert1: %v", err)
	}

	cert2Data, err := os.ReadFile("testdata/certs/g2/finland_msca_card43.bin")
	if err != nil {
		t.Skipf("Certificate not found: %v", err)
	}

	cert2, err := UnmarshalEccCertificate(cert2Data)
	if err != nil {
		t.Fatalf("Failed to unmarshal cert2: %v", err)
	}

	// Try to verify cert1 using cert2 as CA
	// This should fail because cert1's CAR doesn't reference cert2
	err = VerifyEccCertificateWithCA(cert1, cert2)
	if err == nil {
		t.Error("VerifyEccCertificateWithCA() succeeded with mismatched CAR/CHR, want error")
	}
}
