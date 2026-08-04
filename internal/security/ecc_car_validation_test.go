package security

import (
	"strings"
	"testing"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

func TestVerifyEccCertificateWithEccRoot_CARMismatch(t *testing.T) {
	// Build a cert whose CAR does not match the root's CHR.
	// The CAR validation must fail BEFORE reaching signature verification,
	// so we don't need valid keys or signatures — just the CAR/CHR fields.

	cert := &securityv1.EccCertificate{}
	cert.SetCertificateAuthorityReference("FAKE_CA_123")

	root := &securityv1.EccCertificate{}
	root.SetCertificateHolderReference("REAL_ROOT_456")
	// Root needs a public key to pass nil checks that come before CAR validation
	// would come after nil checks in the current code, but the fix should place
	// CAR validation early enough that even without a full key it fires.
	// We set a minimal public key so we don't fail on the "no public key" check.
	rootPK := &securityv1.EccCertificate_PublicKey{}
	rootPK.SetDomainParametersOid("1.3.36.3.3.2.8.1.1.7") // brainpoolP256r1
	rootPK.SetPublicPointX([]byte{0x01})
	rootPK.SetPublicPointY([]byte{0x01})
	root.SetPublicKey(rootPK)

	err := VerifyEccCertificateWithEccRoot(cert, root)
	if err == nil {
		t.Fatal("VerifyEccCertificateWithEccRoot() returned nil error for CAR/CHR mismatch, want error")
	}
	if !strings.Contains(err.Error(), "CAR mismatch") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "CAR mismatch")
	}
	if !strings.Contains(err.Error(), "FAKE_CA_123") {
		t.Errorf("error = %q, want it to mention the cert's CAR %q", err.Error(), "FAKE_CA_123")
	}
	if !strings.Contains(err.Error(), "REAL_ROOT_456") {
		t.Errorf("error = %q, want it to mention the root's CHR %q", err.Error(), "REAL_ROOT_456")
	}
}

func TestVerifyEccCertificateWithEccRoot_CARValidationSkippedWhenEmpty(t *testing.T) {
	// When either CAR or CHR is empty (not yet parsed), we must NOT reject
	// the certificate on CAR grounds — the check should be skipped and
	// verification continues to the signature step (which will fail for
	// other reasons with our stub data, but NOT with "CAR mismatch").

	tests := []struct {
		name    string
		certCAR string
		rootCHR string
	}{
		{
			name:    "empty cert CAR",
			certCAR: "",
			rootCHR: "REAL_ROOT_456",
		},
		{
			name:    "empty root CHR",
			certCAR: "FAKE_CA_123",
			rootCHR: "",
		},
		{
			name:    "both empty",
			certCAR: "",
			rootCHR: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert := &securityv1.EccCertificate{}
			if tt.certCAR != "" {
				cert.SetCertificateAuthorityReference(tt.certCAR)
			}
			// Need raw data to pass the "certificate has no raw data" check
			// that comes after CAR validation (to prove we got past CAR).
			// A zero-length raw data will trigger that error, which is fine —
			// the point is it must NOT be "CAR mismatch".

			root := &securityv1.EccCertificate{}
			if tt.rootCHR != "" {
				root.SetCertificateHolderReference(tt.rootCHR)
			}
			rootPK := &securityv1.EccCertificate_PublicKey{}
			rootPK.SetDomainParametersOid("1.3.36.3.3.2.8.1.1.7") // brainpoolP256r1
			rootPK.SetPublicPointX([]byte{0x01})
			rootPK.SetPublicPointY([]byte{0x01})
			root.SetPublicKey(rootPK)

			err := VerifyEccCertificateWithEccRoot(cert, root)
			// We expect an error (no raw data, bad signature, etc.) but NOT
			// a CAR mismatch error.
			if err != nil && strings.Contains(err.Error(), "CAR mismatch") {
				t.Errorf("got CAR mismatch error when CAR/CHR are empty: %v", err)
			}
		})
	}
}

func TestVerifyEccCertificateWithCA_CARMismatchUnit(t *testing.T) {
	// VerifyEccCertificateWithCA delegates to VerifyEccCertificateWithEccRoot,
	// so CAR validation must also work through that entry point.

	cert := &securityv1.EccCertificate{}
	cert.SetCertificateAuthorityReference("CA_AAA")

	ca := &securityv1.EccCertificate{}
	ca.SetCertificateHolderReference("CA_BBB")
	caPK := &securityv1.EccCertificate_PublicKey{}
	caPK.SetDomainParametersOid("1.3.36.3.3.2.8.1.1.7")
	caPK.SetPublicPointX([]byte{0x01})
	caPK.SetPublicPointY([]byte{0x01})
	ca.SetPublicKey(caPK)

	err := VerifyEccCertificateWithCA(cert, ca)
	if err == nil {
		t.Fatal("VerifyEccCertificateWithCA() returned nil error for CAR/CHR mismatch, want error")
	}
	if !strings.Contains(err.Error(), "CAR mismatch") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "CAR mismatch")
	}
}
