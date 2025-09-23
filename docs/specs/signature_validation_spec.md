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

The core validation logic involves "unwrapping" a certificate to verify its authenticity and extract the public key of the holder. This process is critically dependent on finding the **Certificate Authority Reference (CAR)** within the certificate data, which acts as a pointer to the issuer's certificate. The method for finding the CAR differs significantly between Gen1 and Gen2.

#### 3.2.1. Locating the Certificate Authority Reference (CAR)

**Gen1 (RSA) Certificates:**

Gen1 certificates have a fixed 194-byte structure. The CAR is not part of the signed data but is appended at the end of the certificate block for easy lookup.

-   **Structure**: `[128-byte Signature]` + `[58-byte Content]` + `[8-byte CAR]`
-   **Location**: The CAR is always located at bytes **186-193** of the 194-byte certificate block.

*Code Example (from `benchmark/tachoparser/pkg/decoder/definitions.go`):*
```go
// Simplified logic for Gen1 certificate validation
func (c *CertificateFirstGen) Decode() error {
	// ...
	data := c.Certificate // The raw 194-byte certificate

	// 1. Get the CAR from the fixed offset at the end of the certificate
	var CARPrime uint64
	buf := bytes.NewBuffer(data[186:194]) // Read the last 8 bytes
	binary.Read(buf, binary.BigEndian, &CARPrime)

	// 2. Look up the CA's public key from the global store using the CAR
	ca, ok := PKsFirstGen[CARPrime]
	if !ok {
		return errors.New("could not find CA public key")
	}

	// 3. Use the CA's public key to verify the signature and unwrap the content
    // ... (rest of the verification logic)
}
```

**Gen2 (ECC) Certificates:**

Gen2 certificates use a flexible ASN.1 DER (Tag-Length-Value) format. There are no fixed offsets; the certificate must be parsed field by field.

-   **Structure**: A nested series of TLV-encoded data objects.
-   **Location**: The CAR is the value of the data object with the **Tag `'42'`**.

The parser must iterate through the ASN.1 structure of the certificate body, identify the correct tag, and then read its 8-byte value.

*Code Example (Simplified from `benchmark/tachoparser/pkg/decoder/definitions.go`):*
```go
// Simplified logic for Gen2 certificate validation
func (c *CertificateSecondGen) Decode() error {
    // ...
    // asn1Body.Bytes contains the raw bytes of the certificate body

    // 1. Parse the CPI (first field) to get the remaining bytes
    var asn1CPI asn1.RawValue
    restCPI, err := asn1.Unmarshal(asn1Body.Bytes, &asn1CPI)
    // ...

    // 2. Parse the CAR (second field) from the remaining bytes
    var asn1CAR asn1.RawValue
    _, err = asn1.Unmarshal(restCPI, &asn1CAR)
    if err != nil || asn1CAR.Tag != 0x42 { // Tag '42' identifies the CAR
        return errors.New("could not parse CAR or tag mismatch")
    }

    // 3. Convert the 8-byte value to a uint64
    var car uint64
    buf := bytes.NewBuffer(asn1CAR.Bytes)
    binary.Read(buf, binary.BigEndian, &car)

    // 4. Look up the CA's public key from the global store using the CAR
    if caPK, ok := PKsSecondGen[car]; ok {
        // ... proceed with verification
    }
    // ...
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
