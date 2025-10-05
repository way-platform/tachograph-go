package cert

import (
	"github.com/way-platform/tachograph-go/internal/cert/certcache"
	"github.com/way-platform/tachograph-go/internal/security"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// Root returns the embedded ERCA root certificate.
//
// The ERCA root certificate is the trust anchor for the entire tachograph PKI
// hierarchy. It is trusted a priori and is used to verify Member State CA
// certificates.
//
// This function reads the embedded 144-byte root certificate file and parses
// it into a RootCertificate message.
//
// See Appendix 11, Section 2.1 "European Root Public Key".
func Root() (*securityv1.RootCertificate, error) {
	return security.UnmarshalRootCertificate(certcache.Root())
}
