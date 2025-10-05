package cert

import (
	"context"
	"errors"
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// ChainStore is a store that chains multiple stores.
type ChainStore struct {
	stores []Store
}

var _ Store = &ChainStore{}

// NewChainStore creates a new [ChainStore].
func NewChainStore(stores ...Store) *ChainStore {
	return &ChainStore{
		stores: stores,
	}
}

// GetCertificateG1 implements [Store.GetCertificateG1].
func (s *ChainStore) GetCertificateG1(ctx context.Context, chr string) (*ddv1.RsaCertificate, error) {
	var errs []error
	for _, store := range s.stores {
		cert, err := store.GetCertificateG1(ctx, chr)
		if err == nil {
			return cert, nil
		}
		errs = append(errs, err)
	}
	return nil, fmt.Errorf("failed to get certificate: %w", errors.Join(errs...))
}

// GetCertificateG2 implements [Store.GetCertificateG2].
func (s *ChainStore) GetCertificateG2(ctx context.Context, chr string) (*ddv1.EccCertificate, error) {
	var errs []error
	for _, store := range s.stores {
		cert, err := store.GetCertificateG2(ctx, chr)
		if err == nil {
			return cert, nil
		}
		errs = append(errs, err)
	}
	return nil, fmt.Errorf("failed to get certificate: %w", errors.Join(errs...))
}
