package security

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Clock abstracts time.Now() so certificate validity checks are deterministic
// in tests. Production code uses RealClock; tests inject a fakeClock.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the real wall clock.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time { return time.Now() }

// CheckCertificateValidity verifies that the current time (from clock) falls
// within the certificate's validity window [effectiveDate, expirationDate].
//
// Either or both timestamps may be nil:
//   - nil effectiveDate: no "not yet valid" check
//   - nil expirationDate: no "expired" check (e.g. Gen1 certs with 0xFFFFFFFF EOV)
//
// This is the common implementation used by both RSA and ECC verification paths.
func CheckCertificateValidity(effectiveDate, expirationDate *timestamppb.Timestamp, clock Clock) error {
	now := clock.Now()

	if effectiveDate != nil {
		eff := effectiveDate.AsTime()
		if now.Before(eff) {
			return fmt.Errorf("certificate not yet valid: effective %s, now %s",
				eff.Format(time.RFC3339), now.Format(time.RFC3339))
		}
	}

	if expirationDate != nil {
		exp := expirationDate.AsTime()
		if now.After(exp) {
			return fmt.Errorf("certificate expired: expiration %s, now %s",
				exp.Format(time.RFC3339), now.Format(time.RFC3339))
		}
	}

	return nil
}
