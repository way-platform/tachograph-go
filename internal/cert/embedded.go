package cert

import (
	"context"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/cert/certcache"
	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// EmbeddedStore is a certificate store that uses embedded certificates.
type EmbeddedStore struct{}

var _ Store = &EmbeddedStore{}

// NewEmbeddedStore creates a new [EmbeddedStore].
func NewEmbeddedStore() *EmbeddedStore {
	return &EmbeddedStore{}
}

// GetCertificateG1 retrieves a Gen1 certificate by its CHR.
func (s *EmbeddedStore) GetCertificateG1(ctx context.Context, chr string) (*ddv1.RsaCertificate, error) {
	data, ok := certcache.ReadG1(chr)
	if !ok {
		return nil, fmt.Errorf("certificate not found in embedded store: CHR %s", chr)
	}
	return (dd.UnmarshalOptions{}).UnmarshalRsaCertificate(data)
}

// GetCertificateG2 retrieves a Gen2 certificate by its CHR.
func (s *EmbeddedStore) GetCertificateG2(ctx context.Context, chr string) (*ddv1.EccCertificate, error) {
	data, ok := certcache.ReadG2(chr)
	if !ok {
		return nil, fmt.Errorf("certificate not found in embedded store: CHR %s", chr)
	}
	return (dd.UnmarshalOptions{}).UnmarshalEccCertificate(data)
}
