package security

import (
	"os"
	"testing"
	"time"

	"github.com/way-platform/tachograph-go/internal/cert/certcache"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ---------------------------------------------------------------------------
// RSA (Generation 1) — EndOfValidity only, no EffectiveDate
// ---------------------------------------------------------------------------

func TestRsaCertificate_Expired(t *testing.T) {
	// Load ERCA root certificate
	rootData := certcache.Root()
	root, err := UnmarshalRootCertificate(rootData)
	if err != nil {
		t.Fatalf("Failed to unmarshal root certificate: %v", err)
	}

	// Load a real MSCA certificate
	data, err := os.ReadFile("testdata/certs/g1/finland_tcc37.bin")
	if err != nil {
		t.Fatalf("Failed to read certificate: %v", err)
	}

	cert, err := UnmarshalRsaCertificate(data)
	if err != nil {
		t.Fatalf("UnmarshalRsaCertificate() failed: %v", err)
	}

	// Use a clock set far in the future (year 2099) — well past any certificate's EOV.
	futureClock := newFakeClock(time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC))

	// Verify with the future clock.  The certificate's EndOfValidity (EOV)
	// is in the past relative to the clock → should return an expiration error.
	err = VerifyRsaCertificateWithRootAndClock(cert, root, futureClock)
	if err == nil {
		t.Fatal("VerifyRsaCertificateWithRootAndClock() succeeded for expired certificate, want error")
	}

	// The error message should mention expiration.
	if !containsAny(err.Error(), "expired", "expir") {
		t.Errorf("expected expiration error, got: %v", err)
	}
}

func TestRsaCertificate_NotYetValid(t *testing.T) {
	// Generation 1 RSA certificates have only EndOfValidity (EOV), not an
	// EffectiveDate.  Therefore there is no "not yet valid" state that can be
	// detected from the certificate content alone.
	//
	// This test verifies that when the clock is set BEFORE the EOV
	// (i.e. the certificate is still valid), verification succeeds.

	rootData := certcache.Root()
	root, err := UnmarshalRootCertificate(rootData)
	if err != nil {
		t.Fatalf("Failed to unmarshal root certificate: %v", err)
	}

	data, err := os.ReadFile("testdata/certs/g1/finland_tcc37.bin")
	if err != nil {
		t.Fatalf("Failed to read certificate: %v", err)
	}

	cert, err := UnmarshalRsaCertificate(data)
	if err != nil {
		t.Fatalf("UnmarshalRsaCertificate() failed: %v", err)
	}

	// Use a clock set to 2025 — the Finland TCC37 cert should still be valid.
	validClock := newFakeClock(time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC))

	err = VerifyRsaCertificateWithRootAndClock(cert, root, validClock)
	if err != nil {
		t.Fatalf("VerifyRsaCertificateWithRootAndClock() failed for valid certificate: %v", err)
	}
}

// ---------------------------------------------------------------------------
// ECC (Generation 2) — CertificateEffectiveDate + CertificateExpirationDate
// ---------------------------------------------------------------------------

func TestEccCertificate_Expired(t *testing.T) {
	// Load Gen2 ERCA root certificate
	const ercaCHR = "18250066869740371713"
	ercaData, err := os.ReadFile("../cert/certcache/g2/" + ercaCHR + ".bin")
	if err != nil {
		t.Skipf("ERCA root certificate not found: %v", err)
	}

	ercaCert, err := UnmarshalEccCertificate(ercaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal ERCA certificate: %v", err)
	}

	// Load an MSCA certificate
	mscaData, err := os.ReadFile("testdata/certs/g2/finland_msca_card42.bin")
	if err != nil {
		t.Skipf("MSCA certificate not found: %v", err)
	}

	mscaCert, err := UnmarshalEccCertificate(mscaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal MSCA certificate: %v", err)
	}

	// Use a clock set far in the future — well past the certificate's ExpirationDate.
	futureClock := newFakeClock(time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC))

	err = VerifyEccCertificateWithCAClock(mscaCert, ercaCert, futureClock)
	if err == nil {
		t.Fatal("VerifyEccCertificateWithCAClock() succeeded for expired certificate, want error")
	}

	if !containsAny(err.Error(), "expired", "expir") {
		t.Errorf("expected expiration error, got: %v", err)
	}
}

func TestEccCertificate_NotYetValid(t *testing.T) {
	// Load Gen2 ERCA root certificate
	const ercaCHR = "18250066869740371713"
	ercaData, err := os.ReadFile("../cert/certcache/g2/" + ercaCHR + ".bin")
	if err != nil {
		t.Skipf("ERCA root certificate not found: %v", err)
	}

	ercaCert, err := UnmarshalEccCertificate(ercaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal ERCA certificate: %v", err)
	}

	// Load an MSCA certificate
	mscaData, err := os.ReadFile("testdata/certs/g2/finland_msca_card42.bin")
	if err != nil {
		t.Skipf("MSCA certificate not found: %v", err)
	}

	mscaCert, err := UnmarshalEccCertificate(mscaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal MSCA certificate: %v", err)
	}

	// Use a clock set to epoch (1970) — before the certificate's EffectiveDate.
	pastClock := newFakeClock(time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC))

	err = VerifyEccCertificateWithCAClock(mscaCert, ercaCert, pastClock)
	if err == nil {
		t.Fatal("VerifyEccCertificateWithCAClock() succeeded for not-yet-valid certificate, want error")
	}

	if !containsAny(err.Error(), "not yet valid") {
		t.Errorf("expected 'not yet valid' error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// ECC: valid certificate with correct clock
// ---------------------------------------------------------------------------

func TestEccCertificate_ValidWithClock(t *testing.T) {
	// Load Gen2 ERCA root certificate
	const ercaCHR = "18250066869740371713"
	ercaData, err := os.ReadFile("../cert/certcache/g2/" + ercaCHR + ".bin")
	if err != nil {
		t.Skipf("ERCA root certificate not found: %v", err)
	}

	ercaCert, err := UnmarshalEccCertificate(ercaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal ERCA certificate: %v", err)
	}

	mscaData, err := os.ReadFile("testdata/certs/g2/finland_msca_card42.bin")
	if err != nil {
		t.Skipf("MSCA certificate not found: %v", err)
	}

	mscaCert, err := UnmarshalEccCertificate(mscaData)
	if err != nil {
		t.Fatalf("Failed to unmarshal MSCA certificate: %v", err)
	}

	// Use a clock set to 2025 — the certificate should be within its validity window.
	validClock := newFakeClock(time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC))

	err = VerifyEccCertificateWithCAClock(mscaCert, ercaCert, validClock)
	if err != nil {
		t.Fatalf("VerifyEccCertificateWithCAClock() failed for valid certificate: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Clock interface tests
// ---------------------------------------------------------------------------

func TestRealClock_ReturnsCurrentTime(t *testing.T) {
	c := RealClock{}
	before := time.Now()
	got := c.Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Errorf("RealClock.Now() = %v, want between %v and %v", got, before, after)
	}
}

func TestFakeClock_ReturnsFixedTime(t *testing.T) {
	fixed := time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)
	c := newFakeClock(fixed)
	got := c.Now()

	if !got.Equal(fixed) {
		t.Errorf("fakeClock.Now() = %v, want %v", got, fixed)
	}
}

// ---------------------------------------------------------------------------
// CheckCertificateValidity unit tests (timestamp-level, no real certs)
// ---------------------------------------------------------------------------

func TestCheckCertificateValidity_Expired(t *testing.T) {
	// Certificate valid 2020-01-01 to 2022-01-01
	effective := timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	expiration := timestamppb.New(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))

	// Clock at 2023 — past expiration
	clock := newFakeClock(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))

	err := CheckCertificateValidity(effective, expiration, clock)
	if err == nil {
		t.Fatal("CheckCertificateValidity() succeeded for expired cert, want error")
	}
	if !containsAny(err.Error(), "expired") {
		t.Errorf("expected expiration error, got: %v", err)
	}
}

func TestCheckCertificateValidity_NotYetValid(t *testing.T) {
	// Certificate valid 2025-01-01 to 2030-01-01
	effective := timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	expiration := timestamppb.New(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))

	// Clock at 2024 — before effective date
	clock := newFakeClock(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))

	err := CheckCertificateValidity(effective, expiration, clock)
	if err == nil {
		t.Fatal("CheckCertificateValidity() succeeded for not-yet-valid cert, want error")
	}
	if !containsAny(err.Error(), "not yet valid") {
		t.Errorf("expected 'not yet valid' error, got: %v", err)
	}
}

func TestCheckCertificateValidity_Valid(t *testing.T) {
	// Certificate valid 2020-01-01 to 2030-01-01
	effective := timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	expiration := timestamppb.New(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))

	// Clock at 2025 — within window
	clock := newFakeClock(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC))

	err := CheckCertificateValidity(effective, expiration, clock)
	if err != nil {
		t.Fatalf("CheckCertificateValidity() failed for valid cert: %v", err)
	}
}

func TestCheckCertificateValidity_NilDates(t *testing.T) {
	// When both dates are nil (e.g. 0xFFFFFFFF in Gen1), no error.
	clock := newFakeClock(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC))

	err := CheckCertificateValidity(nil, nil, clock)
	if err != nil {
		t.Fatalf("CheckCertificateValidity() failed with nil dates: %v", err)
	}
}

func TestCheckCertificateValidity_NilExpiration(t *testing.T) {
	// Only effective date set, no expiration — should still pass if now >= effective.
	effective := timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	clock := newFakeClock(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC))

	err := CheckCertificateValidity(effective, nil, clock)
	if err != nil {
		t.Fatalf("CheckCertificateValidity() failed with nil expiration: %v", err)
	}
}

// ---------------------------------------------------------------------------
// test doubles
// ---------------------------------------------------------------------------

// fakeClock is a test double that returns a fixed time.
type fakeClock struct{ t time.Time }

func newFakeClock(t time.Time) Clock { return fakeClock{t: t} }

func (c fakeClock) Now() time.Time { return c.t }

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// containsAny returns true if s contains any of the substrings.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}
