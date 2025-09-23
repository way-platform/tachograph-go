# Specification: DDD File Signature and Certificate Validation

## 1. Introduction

This document specifies the plan for implementing digital signature and certificate chain validation for tachograph data files (`.DDD`) within this library. The goal is to ensure that all parsed data can be cryptographically verified for authenticity and integrity, adhering to the EU regulations for both Gen1 (digital) and Gen2 (smart) tachographs.

A key requirement is a hybrid key management strategy that combines the reliability of bundled certificates with the flexibility of on-demand downloading.

## 2. Background: The Tachograph PKI

The security of tachograph data relies on a Public Key Infrastructure (PKI) with a clear chain of trust:

1.  **European Root Certification Authority (ERCA)**: The top-level trust anchor. It issues certificates for Member State CAs. Its public key must be implicitly trusted.
2.  **Member State Certification Authority (MSCA)**: A national body that issues certificates for equipment manufacturers and personalisers. Its certificate is signed by the ERCA.
3.  **Equipment (VU or Card)**: Each Vehicle Unit or Tachograph Card has its own unique key pair and a certificate for its public key, signed by an MSCA.

To trust a signature from a piece of equipment, its certificate must be validated by checking the signature from the MSCA, whose certificate must in turn be validated against the ERCA root.

- **Gen1 (Digital)**: Uses RSA-1024 keys and SHA-1 hashes.
- **Gen2 (Smart)**: Uses Elliptic Curve Cryptography (ECC) and SHA-2/3 hashes.

## 3. Analysis of the `tachoparser` Benchmark

The `tachoparser` project provides a complete, working example of the validation process.

### 3.1. Key Downloading and Storage

`tachoparser` adopts a "download-then-bundle" strategy.

**1. Downloading:** Python scripts are used to download all public certificates from the official JRC source.

*Code Example (`benchmark/tachoparser/scripts/pks1/dl_all_pks1.py`):*
```python
# ... imports
WWW_PK_URL = "https://dtc.jrc.ec.europa.eu/dtc_public_key_certificates_dt.php.html"
PK_BASE_URL = "https://dtc.jrc.ec.europa.eu/"
TARGET = "../../internal/pkg/certificates/pks1/"

# ...
r = requests.get(WWW_PK_URL)
tree = html.fromstring(r.content)
# Find all links with the title "Download certificate file"
pkas = tree.xpath('//a[@title="Download certificate file"]')
for pka in pkas:
    key_identifier = pka.xpath('text()')[0]
    if not exists(TARGET + key_identifier + ".bin"):
        link = pka.xpath('@href')[0]
        # Download the certificate file
        r_cert = requests.get(PK_BASE_URL + link)
        # ... save the file
```

**2. Storage:** The downloaded `.bin` files are stored in the project tree and embedded directly into the Go binary using the `go:embed` directive.

*Code Example (`benchmark/tachoparser/internal/pkg/certificates/certificates.go`):*
```go
import (
	"embed"
	"log"
	// ... other imports
)

//go:embed pks1/*.bin
var pks1 embed.FS

//go:embed pks2/*.bin
var pks2 embed.FS

func init() {
    // ... logic to read from the embedded filesystem `pks1` and `pks2`
    // and load the keys into a global map.
    f, err := getPks1Fs().Open("EC_PK.bin")
    // ...
}
```
This approach is robust and self-contained but lacks runtime flexibility.

### 3.2. Chain of Trust Validation Logic

The core validation logic in `tachoparser` correctly implements the certificate "unwrapping" process defined in the regulations.

*Code Example (Simplified from `benchmark/tachoparser/pkg/decoder/definitions.go` - `CertificateFirstGen.Decode()`):*
```go
// Simplified logic for Gen1 certificate validation
func (c *CertificateFirstGen) Decode() error {
	cert := new(DecodedCertificateFirstGen)
	data := c.Certificate // The raw 194-byte certificate

	// 1. Get the Certificate Authority Reference (CAR) from the certificate data
	var CARPrime uint64
	buf := bytes.NewBuffer(data[186:194])
	binary.Read(buf, binary.BigEndian, &CARPrime)

	// 2. Look up the CA's public key from the global store of known keys
	ca, ok := PKsFirstGen[CARPrime]
	if !ok {
		return errors.New("could not find CA public key")
	}

	// 3. Use the CA's public key to decrypt the signature part of the certificate
	// This is the "unwrapping" step.
	SrPrime := ca.Perform(data[0:128]) // Perform does the RSA operation

	// 4. Check for correct padding bytes
	if SrPrime[0] != 0x6a || SrPrime[127] != 0xbc {
		return errors.New("invalid signature padding")
	}

	// 5. Reconstruct the original certificate content and its hash
	CrPrime := SrPrime[1 : 1+106]
	HPrime := SrPrime[1+106 : 1+106+20]
	CnPrime := data[128:186]
	CPrime := append(CrPrime, CnPrime...) // This is the full certificate body
	hash := sha1.Sum(CPrime)              // Hash the reconstructed body

	// 6. Compare the computed hash with the hash from the unwrapped signature
	if !reflect.DeepEqual(HPrime, hash[:]) {
		return errors.New("certificate content hash mismatch")
	}

	// 7. If hashes match, the certificate is genuine. Parse its content.
	// ... logic to parse CAR, CHR, public key, etc., from CPrime ...
	c.DecodedCertificate = cert
	return nil
}
```

## 4. Proposed Implementation for `tachograph-go`

We will implement a more flexible `CertificateStore` that supports a hybrid approach.

### 4.1. The Hybrid `CertificateStore`

A new `crypto` package will be created. It will define a `CertificateStore` interface and a default `HybridStore` implementation.

*Proposed `CertificateStore` Interface:*
```go
package crypto

// KeyIdentifier is a unique ID for a certificate, typically the
// Certificate Holder Reference (CHR).
type KeyIdentifier [8]byte

// Certificate represents a parsed public key certificate.
type Certificate interface {
    // Verify checks the certificate's signature against its issuer.
    Verify(store CertificateStore) error
    // ... other methods
}

// CertificateStore is responsible for retrieving certificates.
type CertificateStore interface {
    // GetCertificate retrieves a certificate by its unique identifier.
    GetCertificate(id KeyIdentifier) (Certificate, error)
}
```

*Proposed `HybridStore` Implementation Logic:*
```go
import "embed"

//go:embed certs/root/*.crt
var embeddedRootCerts embed.FS

//go:embed certs/msca/*.crt
var embeddedMscaCerts embed.FS

type HybridStore struct {
    memoryCache       map[KeyIdentifier]Certificate
    fileSystemPath    string // User-configurable path for caching/sideloading
    disableDownloads  bool   // Option to prevent network access
}

func (s *HybridStore) GetCertificate(id KeyIdentifier) (Certificate, error) {
    // 1. Check in-memory cache
    if cert, found := s.memoryCache[id]; found {
        return cert, nil
    }

    // 2. Check embedded root certificates
    if cert, err := s.loadFromFS(embeddedRootCerts, id); err == nil {
        s.memoryCache[id] = cert // Cache it
        return cert, nil
    }

    // 3. Check embedded MSCA certificates
    if cert, err := s.loadFromFS(embeddedMscaCerts, id); err == nil {
        s.memoryCache[id] = cert // Cache it
        return cert, nil
    }

    // 4. Check user-provided filesystem path
    if cert, err := s.loadFromFileSystem(id); err == nil {
        s.memoryCache[id] = cert // Cache it
        return cert, nil
    }

    // 5. If allowed, download from the JRC website
    if !s.disableDownloads {
        if cert, err := s.downloadAndCache(id); err == nil {
            s.memoryCache[id] = cert // Cache it
            return cert, nil
        }
    }

    return nil, errors.New("certificate not found")
}

// ... other helper methods: loadFromFS, loadFromFileSystem, downloadAndCache
```

### 4.2. Verification Workflow

The unmarshalling process will be enhanced to trigger validation.

1.  The unmarshaller parses a signed data block (e.g., `CardEventData`) and its corresponding signature from the DDD file.
2.  It also parses the equipment's certificate (`EF_Card_Certificate`).
3.  It calls `cardCertificate.Verify(certificateStore)`.
4.  The `Verify` method recursively uses the `CertificateStore` to fetch and verify the issuer's certificate (MSCA), and then its issuer (ERCA), establishing a chain of trust.
5.  If the chain is valid, the equipment's public key is trusted.
6.  The unmarshaller then calls a `VerifySignature` function, passing the data block, the signature, and the trusted public key.
7.  The result (`verified: true/false`) is stored in the final protobuf message.

## 5. Implementation Plan

1.  **Create `crypto` Package**: Set up the directory and the `CertificateStore` interface.
2.  **Implement `HybridStore`**: Build the store with the layered lookup logic (memory, embed, filesystem, download).
3.  **Add Certificate Downloader Script**: Create a Go or Python script in `tools/` to download all official certificates into a `dist/certs` directory, which will be used by `go:embed`.
4.  **Define Crypto Structs**: Add Go structs for `CertificateFirstGen`, `CertificateSecondGen`, etc., in the `crypto` package.
5.  **Implement Verification Logic**: Write the `Verify` methods for certificates and signatures, adapting the battle-tested logic from `tachoparser`.
6.  **Integrate into Unmarshallers**: Modify the `unmarshal_*.go` files to call the verification logic and populate a `verified` field in the target protobuf messages (this may require schema changes as per `card_security_data_representation_spec.md`).
