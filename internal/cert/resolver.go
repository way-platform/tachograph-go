package cert

import (
	"context"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// Resolver is an interface for resolving tachograph certificates.
type Resolver interface {
	// GetRootCertificate retrieves the European Root CA certificate.
	GetRootCertificate(ctx context.Context) (*securityv1.RootCertificate, error)

	// GetRsaCertificate retrieves an RSA certificate (Generation 1) by its CHR.
	GetRsaCertificate(ctx context.Context, chr string) (*securityv1.RsaCertificate, error)

	// GetEccCertificate retrieves an ECC certificate (Generation 2) by its CHR.
	GetEccCertificate(ctx context.Context, chr string) (*securityv1.EccCertificate, error)
}
