package cert

import (
	"github.com/way-platform/tachograph-go/internal/cert/certcache"
	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// Root returns the ERCA root certificate.
func Root() (*ddv1.RsaCertificate, error) {
	result, err := (dd.UnmarshalOptions{}).UnmarshalRsaCertificate(certcache.Root())
	if err != nil {
		return nil, err
	}
	return result, nil
}
