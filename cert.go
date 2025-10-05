package tachograph

import (
	"context"
	"net/http"

	"github.com/way-platform/tachograph-go/internal/cert"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// CertificateResolver provides access to tachograph certificates needed for
// signature verification.
//
// Implementations of this interface are responsible for fetching certificates
// by their Certificate Holder Reference (CHR). The default implementation uses
// embedded certificates from EU member states and falls back to fetching from
// remote sources.
type CertificateResolver interface {
	// GetRootCertificate retrieves the European Root CA certificate.
	GetRootCertificate(ctx context.Context) (*securityv1.RootCertificate, error)

	// GetRsaCertificate retrieves an RSA certificate (Generation 1)
	// by its Certificate Holder Reference (CHR).
	GetRsaCertificate(ctx context.Context, chr string) (*securityv1.RsaCertificate, error)

	// GetEccCertificate retrieves an ECC certificate (Generation 2)
	// by its Certificate Holder Reference (CHR).
	GetEccCertificate(ctx context.Context, chr string) (*securityv1.EccCertificate, error)
}

// DefaultCertificateResolver returns the default certificate resolver.
//
// The default resolver uses a chain of certificate sources:
//  1. Embedded certificates from EU member states (fast, offline)
//  2. Remote certificate fetching via HTTP (fallback, requires network)
//
// This resolver is suitable for most use cases and provides good performance
// while ensuring compatibility with certificates from all EU member states.
func DefaultCertificateResolver() CertificateResolver {
	return cert.NewChainResolver(
		cert.NewEmbeddedResolver(),
		cert.NewClient(http.DefaultClient),
	)
}
