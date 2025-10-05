package cert

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/tachograph-go/internal/cert/certcache"
	"github.com/way-platform/tachograph-go/internal/security"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// Client retrieves certificates from the Digital Tachograph Joint Research Centre.
type Client struct {
	httpClient *http.Client
}

var _ Resolver = &Client{}

// NewClient creates a new [Client].
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// GetRootCertificate retrieves the European Root CA certificate.
func (c *Client) GetRootCertificate(ctx context.Context) (*securityv1.RootCertificate, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetRsaCertificate retrieves an RSA certificate by its CHR.
func (c *Client) GetRsaCertificate(ctx context.Context, chr string) (*securityv1.RsaCertificate, error) {
	index, err := certcache.LoadIndex()
	if err != nil {
		return nil, err
	}
	var entry *certcache.Entry
	for _, e := range index.G1 {
		if e.CHR == chr {
			entry = &e
			break
		}
	}
	if entry == nil {
		return nil, fmt.Errorf("certificate not found in index: CHR %s", chr)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, entry.URL, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download certificate: %s", httpResponse.Status)
	}
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	return security.UnmarshalRsaCertificate(body)
}

// GetEccCertificate retrieves an ECC certificate by its CHR.
func (c *Client) GetEccCertificate(ctx context.Context, chr string) (*securityv1.EccCertificate, error) {
	index, err := certcache.LoadIndex()
	if err != nil {
		return nil, err
	}
	var entry *certcache.Entry
	for _, e := range index.G2 {
		if e.CHR == chr {
			entry = &e
			break
		}
	}
	if entry == nil {
		return nil, fmt.Errorf("certificate not found in index: CHR %s", chr)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, entry.URL, nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download certificate: %s", httpResponse.Status)
	}
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	return security.UnmarshalEccCertificate(body)
}
