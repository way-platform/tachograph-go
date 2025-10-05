package cert

import (
	"context"
	"errors"
	"fmt"

	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// ChainResolver chains multiple certificate resolvers, trying each in sequence
// until one succeeds.
type ChainResolver struct {
	resolvers []Resolver
}

var _ Resolver = &ChainResolver{}

// NewChainResolver creates a new [ChainResolver].
func NewChainResolver(resolvers ...Resolver) *ChainResolver {
	return &ChainResolver{
		resolvers: resolvers,
	}
}

// GetRootCertificate implements [Resolver.GetRootCertificate].
func (r *ChainResolver) GetRootCertificate(ctx context.Context) (*securityv1.RootCertificate, error) {
	var errs []error
	for _, resolver := range r.resolvers {
		cert, err := resolver.GetRootCertificate(ctx)
		if err == nil {
			return cert, nil
		}
		errs = append(errs, err)
	}
	return nil, fmt.Errorf("failed to get root certificate: %w", errors.Join(errs...))
}

// GetRsaCertificate implements [Resolver.GetRsaCertificate].
func (r *ChainResolver) GetRsaCertificate(ctx context.Context, chr string) (*securityv1.RsaCertificate, error) {
	var errs []error
	for _, resolver := range r.resolvers {
		cert, err := resolver.GetRsaCertificate(ctx, chr)
		if err == nil {
			return cert, nil
		}
		errs = append(errs, err)
	}
	return nil, fmt.Errorf("failed to get RSA certificate: %w", errors.Join(errs...))
}

// GetEccCertificate implements [Resolver.GetEccCertificate].
func (r *ChainResolver) GetEccCertificate(ctx context.Context, chr string) (*securityv1.EccCertificate, error) {
	var errs []error
	for _, resolver := range r.resolvers {
		cert, err := resolver.GetEccCertificate(ctx, chr)
		if err == nil {
			return cert, nil
		}
		errs = append(errs, err)
	}
	return nil, fmt.Errorf("failed to get ECC certificate: %w", errors.Join(errs...))
}
