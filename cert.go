package tachograph

import (
	"context"
	"net/http"

	"github.com/way-platform/tachograph-go/internal/cert"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// CertificateStore provides access to Certificate Authority (CA) certificates
// needed for signature verification of tachograph files.
type CertificateStore interface {
	// GetCertificateG1 retrieves a Generation 1 RSA CA certificate
	// by its Certificate Authority Reference (CAR).
	GetCertificateG1(ctx context.Context, chr string) (*ddv1.RsaCertificate, error)

	// GetCertificateG2 retrieves a Generation 2 ECC CA certificate
	// by its Certificate Authority Reference (CAR).
	GetCertificateG2(ctx context.Context, chr string) (*ddv1.EccCertificate, error)
}

// DefaultCertificateStore returns a certificate store with embedded
// common CA certificates from EU member states.
func DefaultCertificateStore() CertificateStore {
	return cert.NewChainStore(
		cert.NewEmbeddedStore(),
		cert.NewClient(http.DefaultClient),
	)
}
