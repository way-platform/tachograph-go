package security

import (
	"testing"

	"github.com/way-platform/tachograph-go/internal/cert/certcache"
)

func TestUnmarshalRootCertificate(t *testing.T) {
	// Read the embedded ERCA root certificate
	rootData := certcache.Root()

	// Validate the data length before parsing
	const lenRootCert = 144
	if len(rootData) != lenRootCert {
		t.Fatalf("Root certificate data length = %d, want %d", len(rootData), lenRootCert)
	}

	// Parse the root certificate
	root, err := UnmarshalRootCertificate(rootData)
	if err != nil {
		t.Fatalf("UnmarshalRootCertificate() failed: %v", err)
	}

	// Validate the parsed certificate
	if root == nil {
		t.Fatal("UnmarshalRootCertificate() returned nil")
	}

	// Validate key identifier is present
	keyID := root.GetKeyId()
	if keyID == "" {
		t.Error("KeyId is empty")
	}

	// The key ID should be "18250066869723594497" for the current ERCA root
	const expectedKeyID = "18250066869723594497"
	if keyID != expectedKeyID {
		t.Errorf("KeyId = %q, want %q", keyID, expectedKeyID)
	}

	// Validate RSA modulus is present and correct length
	modulus := root.GetRsaModulus()
	const lenModulus = 128
	if len(modulus) != lenModulus {
		t.Errorf("RsaModulus length = %d, want %d", len(modulus), lenModulus)
	}

	// Validate RSA exponent is present and correct length
	exponent := root.GetRsaExponent()
	const lenExponent = 8
	if len(exponent) != lenExponent {
		t.Errorf("RsaExponent length = %d, want %d", len(exponent), lenExponent)
	}
}

// TestUnmarshalRootCertificate_InvalidLength tests that parsing fails
// when the input data has an incorrect length.
func TestUnmarshalRootCertificate_InvalidLength(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "too short",
			data: make([]byte, 143),
		},
		{
			name: "too long",
			data: make([]byte, 145),
		},
		{
			name: "empty",
			data: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalRootCertificate(tt.data)
			if err == nil {
				t.Error("UnmarshalRootCertificate() succeeded, want error")
			}
		})
	}
}
