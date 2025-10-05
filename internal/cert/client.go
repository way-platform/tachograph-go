package cert

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/way-platform/tachograph-go/internal/cert/certcache"
	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// Client is a client retreiving certificates from the Digital Tachograph Joint Research Centre.
type Client struct {
	httpClient *http.Client
}

var _ Store = &Client{}

// NewClient creates a new [Client].
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// GetCertificateG1 retrieves a Gen1 certificate by its CHR.
func (c *Client) GetCertificateG1(ctx context.Context, chr string) (*ddv1.RsaCertificate, error) {
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
	result, err := (dd.UnmarshalOptions{}).UnmarshalRsaCertificate(body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCertificateG2 retrieves a Gen2 certificate by its CHR.
func (c *Client) GetCertificateG2(ctx context.Context, chr string) (*ddv1.EccCertificate, error) {
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
	result, err := (dd.UnmarshalOptions{}).UnmarshalEccCertificate(body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
