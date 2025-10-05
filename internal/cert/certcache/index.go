package certcache

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"
)

//go:embed index.json
var index []byte

// LoadIndex loads the cached certificate index.
func LoadIndex() (*Index, error) {
	var result Index
	if err := json.Unmarshal(index, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index: %w", err)
	}
	return &result, nil
}

// Entry represents an entry in the cached certificate index.
type Entry struct {
	// CAR is the Certificate Authority Reference.
	CAR string `json:"car"`
	// CHR is the Certificate Holder Reference.
	CHR string `json:"chr"`
	// Country is the country of the certificate.
	Country string `json:"country"`
	// URL is the URL of the certificate.
	URL string `json:"url"`
	// ExpirationDate is the expiration date of the certificate.
	ExpirationDate string `json:"expirationDate,omitzero"`
}

// Index represents the cached certificate index.
type Index struct {
	// CreateTime is the time the index was created.
	CreateTime time.Time `json:"createTime"`
	// Root is the root certificate.
	Root Entry `json:"root"`
	// G1 is the Gen1 certificates.
	G1 []Entry `json:"g1"`
	// G2 is the Gen2 certificates.
	G2 []Entry `json:"g2"`
}
