package cert

import (
	"context"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// Store is an interface for a certificate store.
type Store interface {
	// GetCertificateG1 retrieves a Gen1 certificate by its CHR.
	GetCertificateG1(ctx context.Context, chr string) (*ddv1.RsaCertificate, error)
	// GetCertificateG2 retrieves a Gen2 certificate by its CHR.
	GetCertificateG2(ctx context.Context, chr string) (*ddv1.EccCertificate, error)
}
