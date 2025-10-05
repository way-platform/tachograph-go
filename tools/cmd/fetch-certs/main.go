package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/asn1"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

// TimeReal represents a TimeReal timestamp from tachograph data.
// This is a simplified implementation for certificate parsing.
type TimeReal struct {
	Timedata [4]byte
}

// Decode converts TimeReal bytes to time.Time.
func (tr TimeReal) Decode() time.Time {
	timeVal := binary.BigEndian.Uint32(tr.Timedata[:])
	if timeVal == 0 {
		return time.Time{} // Zero time
	}
	return time.Unix(int64(timeVal), 0)
}

// CertificateIndex contains metadata about all available certificates.
type CertificateIndex struct {
	// CreateTime is when this index was created
	CreateTime time.Time `json:"createTime"`

	// Root is the European Root CA (ERCA) certificate that signs all member state certificates
	Root *CertificateEntry `json:"root"`

	// G1 contains Gen1 (RSA) member state certificates
	G1 []CertificateEntry `json:"g1"`

	// G2 contains Gen2 (ECC) member state certificates
	G2 []CertificateEntry `json:"g2"`
}

// MarshalJSON implements custom JSON marshaling for CertificateIndex
func (ci CertificateIndex) MarshalJSON() ([]byte, error) {
	type Alias CertificateIndex
	return json.Marshal(&struct {
		CreateTime string `json:"createTime"`
		Alias
	}{
		CreateTime: ci.CreateTime.Format(time.RFC3339),
		Alias:      (Alias)(ci),
	})
}

// CertificateEntry contains metadata about a single certificate.
type CertificateEntry struct {
	// CAR is the Certificate Authority Reference
	CAR string `json:"car"`

	// CHR is the Certificate Holder Reference (subject of this certificate)
	// For CA certificates, this is what other certificates reference via their CAR
	CHR string `json:"chr,omitempty"`

	// Country is the full country name for member state certificates (empty for root certificate)
	Country string `json:"country,omitempty"`

	// URL is the full URL to download the certificate
	URL string `json:"url"`

	// EffectiveDate is when the certificate becomes valid (ISO 8601 format)
	EffectiveDate string `json:"effectiveDate,omitempty"`

	// ExpirationDate is when the certificate expires (ISO 8601 format)
	ExpirationDate string `json:"expirationDate,omitempty"`

	// Filename is the suggested filename for embedded storage
	Filename string `json:"filename"`

	// Path is the relative path to the certificate file from the index location
	Path string `json:"path"`
}

func main() {
	var (
		outputDir   = flag.String("dir", "certs", "Output directory for certificates and index")
		baseURL     = flag.String("base-url", "https://dtc.jrc.ec.europa.eu", "Base URL for the DTC website")
		concurrency = flag.Int("concurrency", 10, "Maximum number of concurrent downloads")
	)
	flag.Parse()

	log.Printf("Building certificate index from %s (concurrency: %d)", *baseURL, *concurrency)
	log.Printf("Output directory: %s", *outputDir)

	ctx := context.Background()

	// Create directory structure
	if err := createDirectoryStructure(*outputDir); err != nil {
		log.Fatalf("Failed to create directory structure: %v", err)
	}

	index := CertificateIndex{
		CreateTime: time.Now().UTC(),
	}

	// Download and save root certificate
	log.Println("Downloading root certificate...")
	ercaCertData, err := downloadAndSaveRootCertificate(ctx, *outputDir)
	if err != nil {
		log.Fatalf("Failed to download root certificate: %v", err)
	}
	log.Printf("Root certificate saved (%d bytes)", len(ercaCertData))

	// Extract key identifier from root certificate (first 8 bytes)
	// The root certificate is in raw format: [8 bytes key ID][128 bytes modulus][8 bytes exponent]
	var ercaKeyID uint64
	if len(ercaCertData) >= 8 {
		ercaKeyID = binary.BigEndian.Uint64(ercaCertData[0:8])
	}

	// Set root certificate in index
	// For the root certificate, CAR and CHR are the same (self-signed)
	// Country is not set for the root certificate as it's the European-level CA
	index.Root = &CertificateEntry{
		CAR:      fmt.Sprintf("%d", ercaKeyID),
		CHR:      fmt.Sprintf("%d", ercaKeyID), // Self-signed
		URL:      "https://dtc.jrc.ec.europa.eu/erca_of_doc/EC_PK.zip",
		Filename: "EC_PK.bin",
		Path:     "root/EC_PK.bin",
	}
	log.Printf("Root certificate: %s (Key ID: %s)", index.Root.Filename, index.Root.CAR)

	// Extract public key components for signature recovery
	ercaModulus, ercaExponent, err := extractERCA(ercaCertData)
	if err != nil {
		log.Fatalf("Failed to extract ERCA public key: %v", err)
	}

	// Build Gen1 certificate index
	log.Println("Indexing and downloading Gen1 certificates...")
	index.G1, err = indexAndDownloadGen1Certificates(ctx, *baseURL, *outputDir, *concurrency, ercaModulus, ercaExponent)
	if err != nil {
		log.Fatalf("Failed to index Gen1 certificates: %v", err)
	}
	log.Printf("Downloaded %d Gen1 member state certificates", len(index.G1))

	// Build Gen2 certificate index
	log.Println("Indexing and downloading Gen2 certificates...")
	index.G2, err = indexAndDownloadGen2Certificates(ctx, *baseURL, *outputDir, *concurrency)
	if err != nil {
		log.Fatalf("Failed to index Gen2 certificates: %v", err)
	}
	log.Printf("Downloaded %d Gen2 member state certificates", len(index.G2))

	// Sort certificates by country (ascending) and expiration date (descending)
	log.Println("Sorting certificate entries...")
	sortCertificates(index.G1)
	sortCertificates(index.G2)

	// Write index to file
	indexPath := filepath.Join(*outputDir, "index.json")
	if err := writeIndexFile(&index, indexPath); err != nil {
		log.Fatalf("Failed to write index file: %v", err)
	}

	log.Printf("Certificate index written to %s", indexPath)
	log.Printf("Total certificates: 1 root + %d Gen1 + %d Gen2",
		len(index.G1), len(index.G2))
}

// createDirectoryStructure creates the required directory structure for certificates
// and removes any existing .bin files to ensure a clean state.
func createDirectoryStructure(baseDir string) error {
	dirs := []string{
		filepath.Join(baseDir, "root"),
		filepath.Join(baseDir, "g1"),
		filepath.Join(baseDir, "g2"),
	}

	// Create directories if they don't exist
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Remove existing .bin files to ensure clean state
		files, err := filepath.Glob(filepath.Join(dir, "*.bin"))
		if err != nil {
			return fmt.Errorf("failed to list files in %s: %w", dir, err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", file, err)
			}
		}
	}

	log.Printf("Created directory structure in %s", baseDir)
	return nil
}

// downloadAndSaveRootCertificate downloads the ERCA root certificate and saves it to disk.
func downloadAndSaveRootCertificate(ctx context.Context, baseDir string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	certData, err := downloadERCA(ctx, client)
	if err != nil {
		return nil, err
	}

	// Save to root/EC_PK.bin
	certPath := filepath.Join(baseDir, "root", "EC_PK.bin")
	if err := os.WriteFile(certPath, certData, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write root certificate: %w", err)
	}

	log.Printf("  Saved root certificate to %s", certPath)
	return certData, nil
}

// sortCertificates sorts certificate entries by country (ascending) and expiration date (descending).
func sortCertificates(certs []CertificateEntry) {
	sort.Slice(certs, func(i, j int) bool {
		// First, sort by country (ascending)
		if certs[i].Country != certs[j].Country {
			return certs[i].Country < certs[j].Country
		}

		// If same country, sort by expiration date (descending - newest first)
		// Parse the expiration dates
		timeI, errI := time.Parse(time.RFC3339, certs[i].ExpirationDate)
		timeJ, errJ := time.Parse(time.RFC3339, certs[j].ExpirationDate)

		// Handle parse errors by putting invalid dates at the end
		if errI != nil && errJ != nil {
			return false // Keep original order if both are invalid
		}
		if errI != nil {
			return false // Put invalid date at the end
		}
		if errJ != nil {
			return true // Put invalid date at the end
		}

		// Sort by expiration date descending (newer certificates first)
		return timeI.After(timeJ)
	})
}

func indexAndDownloadGen1Certificates(ctx context.Context, baseURL, outputDir string, concurrency int, ercaModulus, ercaExponent *big.Int) ([]CertificateEntry, error) {
	url := baseURL + "/dtc_public_key_certificates_dt.php.html"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Gen1 certificate list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Parse the HTML table to extract certificate metadata including dates
	certMetadata := parseCertificateTable(doc, baseURL)
	log.Printf("Found %d Gen1 certificate links", len(certMetadata))

	// Process certificates concurrently
	var mu sync.Mutex
	var certificates []CertificateEntry

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	_ = gctx // Context is embedded in errgroup

	for _, meta := range certMetadata {
		meta := meta // capture loop variable
		g.Go(func() error {
			certData, err := downloadCertificate(meta.downloadURL, 194)
			if err != nil {
				log.Printf("Warning: Failed to download %s: %v", meta.name, err)
				return nil // Don't fail entire batch, just skip this cert
			}

			car, err := extractGen1CAR(certData)
			if err != nil {
				log.Printf("Warning: Failed to extract CAR from %s: %v", meta.name, err)
				return nil
			}

			carStr := fmt.Sprintf("%d", car)

			// Perform signature recovery to extract CHR
			chrStr := ""
			eov := meta.expirationDate
			cPrime, err := recoverGen1Certificate(certData, ercaModulus, ercaExponent)
			if err == nil {
				chr, extractedEOV, extractErr := extractGen1CHRAndEOV(cPrime)
				if extractErr == nil {
					chrStr = fmt.Sprintf("%d", chr)
					if extractedEOV != "" {
						eov = extractedEOV
					}
				}
			}

			// Use CHR as filename if available, otherwise use the metadata name
			filename := meta.name + ".bin"
			if chrStr != "" {
				filename = chrStr + ".bin"
			}

			// Save certificate to g1/<CHR>.bin
			certPath := filepath.Join(outputDir, "g1", filename)
			if err := os.WriteFile(certPath, certData, 0o644); err != nil {
				log.Printf("Warning: Failed to save %s: %v", filename, err)
				return nil
			}

			entry := CertificateEntry{
				CAR:            carStr,
				CHR:            chrStr,
				Country:        meta.country,
				URL:            meta.downloadURL,
				ExpirationDate: eov,
				Filename:       filename,
				Path:           filepath.Join("g1", filename),
			}

			mu.Lock()
			certificates = append(certificates, entry)
			mu.Unlock()

			if eov != "" {
				log.Printf("  Downloaded: %s (%s, CAR: %s, CHR: %s, Expires: %s)", filename, meta.country, carStr, chrStr, eov)
			} else {
				log.Printf("  Downloaded: %s (%s, CAR: %s, CHR: %s)", filename, meta.country, carStr, chrStr)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return certificates, nil
}

func indexAndDownloadGen2Certificates(ctx context.Context, baseURL, outputDir string, concurrency int) ([]CertificateEntry, error) {
	url := baseURL + "/dtc_public_key_certificates_st.php.html"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Gen2 certificate list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Parse the HTML table to extract certificate metadata including dates
	certMetadata := parseCertificateTable(doc, baseURL)
	log.Printf("Found %d Gen2 certificate links", len(certMetadata))

	// Process certificates concurrently
	var mu sync.Mutex
	var certificates []CertificateEntry

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	_ = gctx // Context is embedded in errgroup

	for _, meta := range certMetadata {
		meta := meta // capture loop variable
		g.Go(func() error {
			certData, err := downloadCertificate(meta.downloadURL, 341)
			if err != nil {
				log.Printf("Warning: Failed to download %s: %v", meta.name, err)
				return nil // Don't fail entire batch
			}

			if len(certData) < 204 || len(certData) > 341 {
				log.Printf("Warning: Invalid certificate size for %s: %d bytes", meta.name, len(certData))
				return nil
			}

			car, err := extractGen2CAR(certData)
			if err != nil {
				log.Printf("Warning: Failed to extract CAR from %s: %v", meta.name, err)
				return nil
			}

			// Extract CHR from Gen2 certificate (it's in the ASN.1 structure)
			chr, err := extractGen2CHR(certData)
			chrStr := ""
			if err == nil {
				chrStr = fmt.Sprintf("%d", chr)
			}

			// Use CHR as filename if available, otherwise use the metadata name
			filename := meta.name + ".bin"
			if chrStr != "" {
				filename = chrStr + ".bin"
			}

			// Save certificate to g2/<CHR>.bin
			certPath := filepath.Join(outputDir, "g2", filename)
			if err := os.WriteFile(certPath, certData, 0o644); err != nil {
				log.Printf("Warning: Failed to save %s: %v", filename, err)
				return nil
			}

			carStr := fmt.Sprintf("%d", car)
			entry := CertificateEntry{
				CAR:            carStr,
				CHR:            chrStr,
				Country:        meta.country,
				URL:            meta.downloadURL,
				ExpirationDate: meta.expirationDate,
				Filename:       filename,
				Path:           filepath.Join("g2", filename),
			}

			mu.Lock()
			certificates = append(certificates, entry)
			mu.Unlock()

			if meta.expirationDate != "" {
				log.Printf("  Downloaded: %s (%s, CAR: %s, CHR: %s, Expires: %s)", filename, meta.country, carStr, chrStr, meta.expirationDate)
			} else {
				log.Printf("  Downloaded: %s (%s, CAR: %s, CHR: %s)", filename, meta.country, carStr, chrStr)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return certificates, nil
}

// extractGen1CAR extracts the Certificate Authority Reference from a Gen1 certificate.
// Gen1 certificates are 194 bytes with CAR at bytes 186-193 (8 bytes, big-endian uint64).
func extractGen1CAR(certData []byte) (uint64, error) {
	const (
		lenCert = 194
		idxCAR  = 186
		lenCAR  = 8
	)

	if len(certData) != lenCert {
		return 0, fmt.Errorf("invalid Gen1 certificate size: got %d, want %d", len(certData), lenCert)
	}

	car := binary.BigEndian.Uint64(certData[idxCAR : idxCAR+lenCAR])
	return car, nil
}

// downloadERCA downloads the ERCA root certificate from the EU DTC website.
// The ERCA certificate is distributed as a ZIP file containing EC_PK.bin.
func downloadERCA(ctx context.Context, client *http.Client) ([]byte, error) {
	const ercaZipURL = "https://dtc.jrc.ec.europa.eu/erca_of_doc/EC_PK.zip"

	log.Println("  Downloading ERCA root certificate ZIP...")
	req, err := http.NewRequestWithContext(ctx, "GET", ercaZipURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERCA request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download ERCA ZIP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ERCA download failed with status %d", resp.StatusCode)
	}

	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read ERCA ZIP: %w", err)
	}

	// Parse ZIP file
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ERCA ZIP: %w", err)
	}

	// Extract EC_PK.bin from ZIP
	for _, file := range zipReader.File {
		if file.Name == "EC_PK.bin" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open EC_PK.bin in ZIP: %w", err)
			}
			defer rc.Close()

			certData, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read EC_PK.bin: %w", err)
			}

			log.Printf("  Successfully downloaded ERCA certificate (%d bytes)", len(certData))
			return certData, nil
		}
	}

	return nil, fmt.Errorf("EC_PK.bin not found in ERCA ZIP")
}

// extractERCA extracts the RSA public key from the ERCA root certificate.
// The ERCA certificate is in a raw RSA public key format:
// - Bytes 0-7: Key identifier (8 bytes)
// - Bytes 8-135: RSA modulus (128 bytes)
// - Bytes 136-143: RSA exponent (8 bytes)
func extractERCA(certData []byte) (*big.Int, *big.Int, error) {
	const (
		lenCert     = 144
		idxModulus  = 8
		idxExponent = 136
		lenModulus  = 128
		lenExponent = 8
	)

	log.Printf("  ERCA certificate size: %d bytes", len(certData))

	if len(certData) != lenCert {
		return nil, nil, fmt.Errorf("invalid ERCA certificate size: got %d, want %d", len(certData), lenCert)
	}

	// Extract the RSA public key components from the raw format
	modulus := new(big.Int).SetBytes(certData[idxModulus : idxModulus+lenModulus])
	exponent := new(big.Int).SetBytes(certData[idxExponent : idxExponent+lenExponent])

	if modulus.Sign() == 0 || exponent.Sign() == 0 {
		return nil, nil, fmt.Errorf("ERCA certificate has invalid public key")
	}

	log.Printf("  Successfully extracted ERCA public key (raw format, modulus: %d bytes, exponent: %d bytes)",
		len(modulus.Bytes()), len(exponent.Bytes()))

	return modulus, exponent, nil
}

// recoverGen1Certificate performs signature recovery on a Gen1 certificate using the provided CA public key.
// Returns the recovered certificate content (Cr' || Cn') and any error encountered.
func recoverGen1Certificate(certData []byte, caModulus, caExponent *big.Int) ([]byte, error) {
	const (
		lenCert      = 194
		lenSignature = 128
		lenCnPrime   = 58
		lenCAR       = 8
	)

	if len(certData) != lenCert {
		return nil, fmt.Errorf("invalid certificate size: got %d, want %d", len(certData), lenCert)
	}

	// Extract components
	signature := certData[0:lenSignature]
	cnPrime := certData[lenSignature : lenSignature+lenCnPrime]
	_ = binary.BigEndian.Uint64(certData[lenSignature+lenCnPrime : lenSignature+lenCnPrime+lenCAR]) // CAR' (unused for now)

	// Perform RSA signature recovery: Sr' = signature^e mod n
	srPrimeBig := new(big.Int).SetBytes(signature)
	srPrimeBig.Exp(srPrimeBig, caExponent, caModulus)
	srPrime := srPrimeBig.Bytes()

	// Pad Sr' to 128 bytes if necessary
	if len(srPrime) < lenSignature {
		paddedSrPrime := make([]byte, lenSignature)
		copy(paddedSrPrime[lenSignature-len(srPrime):], srPrime)
		srPrime = paddedSrPrime
	} else if len(srPrime) > lenSignature {
		return nil, fmt.Errorf("recovered Sr' is too long: %d bytes", len(srPrime))
	}

	// Verify structure: 6A || Cr' || H' || BC
	const (
		headerByte    = 0x6A
		trailerByte   = 0xBC
		lenCrPrime    = 106
		lenHPrime     = 20
		expectedSrLen = 128
	)

	if len(srPrime) != expectedSrLen || srPrime[0] != headerByte || srPrime[expectedSrLen-1] != trailerByte {
		return nil, fmt.Errorf("invalid recovered Sr' structure")
	}

	crPrime := srPrime[1 : 1+lenCrPrime]
	hPrime := srPrime[1+lenCrPrime : 1+lenCrPrime+lenHPrime]

	// Reconstruct C' = Cr' || Cn'
	cPrime := append(crPrime, cnPrime...)
	if len(cPrime) != 164 {
		return nil, fmt.Errorf("invalid C' length: got %d, want 164", len(cPrime))
	}

	// Verify hash: SHA-1(C') == H'
	hash := sha1.Sum(cPrime)
	if !bytes.Equal(hPrime, hash[:]) {
		return nil, fmt.Errorf("certificate content hash mismatch")
	}

	return cPrime, nil
}

// extractGen1CHRAndEOV extracts the Certificate Holder Reference and End of Validity from recovered Gen1 certificate content.
func extractGen1CHRAndEOV(cPrime []byte) (chr uint64, eov string, err error) {
	const (
		idxCHR    = 20
		idxEOV    = 16
		lenCHR    = 8
		lenEOV    = 4
		lenCPrime = 164
	)

	if len(cPrime) != lenCPrime {
		return 0, "", fmt.Errorf("invalid C' length: got %d, want %d", len(cPrime), lenCPrime)
	}

	chr = binary.BigEndian.Uint64(cPrime[idxCHR : idxCHR+lenCHR])

	// Extract EOV (End of Validity) - TimeReal format
	eovTimeReal := TimeReal{}
	copy(eovTimeReal.Timedata[:], cPrime[idxEOV:idxEOV+lenEOV])
	eovTime := eovTimeReal.Decode()

	// Convert to UTC and format as RFC3339
	eov = eovTime.UTC().Format(time.RFC3339)

	return chr, eov, nil
}

// extractGen1PublicKey extracts the RSA public key (modulus and exponent) from a Gen1 certificate.
// This requires signature recovery using the issuer's public key.
func extractGen1PublicKey(certData []byte, issuerModulus, issuerExponent []byte) (modulus []byte, exponent []byte, err error) {
	const (
		lenCert = 194
		idxSig  = 0
		idxCn   = 128
		idxCAR  = 186
	)

	if len(certData) != lenCert {
		return nil, nil, fmt.Errorf("invalid certificate size: got %d, want %d", len(certData), lenCert)
	}

	// Extract signature and non-recoverable part
	signature := certData[idxSig:idxCn]
	cnPrime := certData[idxCn:idxCAR]

	// Perform RSA signature recovery
	n := new(big.Int).SetBytes(issuerModulus)
	e := new(big.Int).SetBytes(issuerExponent)
	sigBig := new(big.Int).SetBytes(signature)

	srPrimeBig := new(big.Int).Exp(sigBig, e, n)
	srPrime := srPrimeBig.Bytes()

	// Pad to 128 bytes if necessary
	if len(srPrime) < 128 {
		padded := make([]byte, 128)
		copy(padded[128-len(srPrime):], srPrime)
		srPrime = padded
	} else if len(srPrime) > 128 {
		return nil, nil, fmt.Errorf("recovered signature too long: %d bytes", len(srPrime))
	}

	// Verify structure
	const (
		srPrimeHeader  = 0x6A
		srPrimeTrailer = 0xBC
		lenCrPrime     = 106
		lenHPrime      = 20
	)

	if srPrime[0] != srPrimeHeader || srPrime[127] != srPrimeTrailer {
		return nil, nil, fmt.Errorf("invalid recovered signature format")
	}

	crPrime := srPrime[1 : 1+lenCrPrime]
	hPrime := srPrime[1+lenCrPrime : 1+lenCrPrime+lenHPrime]

	// Reconstruct and verify
	cPrime := append(crPrime, cnPrime...)
	hash := sha1.Sum(cPrime)
	for i := range hPrime {
		if hPrime[i] != hash[i] {
			return nil, nil, fmt.Errorf("hash mismatch")
		}
	}

	// Extract modulus (bytes 28-155) and exponent (bytes 156-163)
	const (
		idxModulus  = 28
		lenModulus  = 128
		idxExponent = 156
		lenExponent = 8
	)

	modulus = make([]byte, lenModulus)
	copy(modulus, cPrime[idxModulus:idxModulus+lenModulus])

	exponent = make([]byte, lenExponent)
	copy(exponent, cPrime[idxExponent:idxExponent+lenExponent])

	return modulus, exponent, nil
}

// extractGen2CAR extracts the Certificate Authority Reference from a Gen2 certificate.
// Gen2 certificates use ASN.1 encoding. We need to parse the outer SEQUENCE,
// then the certificate body SEQUENCE, then extract the CAR field.
func extractGen2CAR(certData []byte) (uint64, error) {
	// Parse outer SEQUENCE
	var outerSeq asn1.RawValue
	_, err := asn1.Unmarshal(certData, &outerSeq)
	if err != nil {
		return 0, fmt.Errorf("failed to parse outer SEQUENCE: %w", err)
	}

	// Parse certificate body SEQUENCE
	var bodySeq asn1.RawValue
	_, err = asn1.Unmarshal(outerSeq.Bytes, &bodySeq)
	if err != nil {
		return 0, fmt.Errorf("failed to parse body SEQUENCE: %w", err)
	}

	// Parse CPI (skip it)
	var cpiRaw asn1.RawValue
	restAfterCPI, err := asn1.Unmarshal(bodySeq.Bytes, &cpiRaw)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CPI: %w", err)
	}

	// Parse CAR
	var carRaw asn1.RawValue
	_, err = asn1.Unmarshal(restAfterCPI, &carRaw)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CAR: %w", err)
	}

	if len(carRaw.Bytes) != 8 {
		return 0, fmt.Errorf("invalid CAR length: got %d, want 8", len(carRaw.Bytes))
	}

	car := binary.BigEndian.Uint64(carRaw.Bytes)
	return car, nil
}

// extractGen2CHR extracts the Certificate Holder Reference from a Gen2 certificate.
// CHR is the fifth field in the certificate body (after CPI, CAR, CHA, PublicKey).
func extractGen2CHR(certData []byte) (uint64, error) {
	// Parse outer SEQUENCE
	var outerSeq asn1.RawValue
	_, err := asn1.Unmarshal(certData, &outerSeq)
	if err != nil {
		return 0, fmt.Errorf("failed to parse outer SEQUENCE: %w", err)
	}

	// Parse certificate body SEQUENCE
	var bodySeq asn1.RawValue
	_, err = asn1.Unmarshal(outerSeq.Bytes, &bodySeq)
	if err != nil {
		return 0, fmt.Errorf("failed to parse body SEQUENCE: %w", err)
	}

	// Parse and skip: CPI, CAR, CHA, PublicKey
	rest := bodySeq.Bytes
	for i := 0; i < 4; i++ {
		var skip asn1.RawValue
		rest, err = asn1.Unmarshal(rest, &skip)
		if err != nil {
			return 0, fmt.Errorf("failed to skip field %d: %w", i, err)
		}
	}

	// Parse CHR (5th field)
	var chrRaw asn1.RawValue
	_, err = asn1.Unmarshal(rest, &chrRaw)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CHR: %w", err)
	}

	if len(chrRaw.Bytes) != 8 {
		return 0, fmt.Errorf("invalid CHR length: got %d, want 8", len(chrRaw.Bytes))
	}

	chr := binary.BigEndian.Uint64(chrRaw.Bytes)
	return chr, nil
}

// certificateMetadata holds parsed certificate metadata from HTML table
type certificateMetadata struct {
	name           string
	downloadURL    string
	country        string
	expirationDate string
}

// parseCertificateTable extracts certificate metadata from the HTML table
func parseCertificateTable(n *html.Node, baseURL string) []certificateMetadata {
	var certs []certificateMetadata

	// Find all table rows
	var findRows func(*html.Node)
	findRows = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			// Parse the row
			if cert := parseTableRow(node, baseURL); cert != nil {
				certs = append(certs, *cert)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findRows(c)
		}
	}

	findRows(n)
	return certs
}

// parseTableRow extracts certificate data from a single table row
func parseTableRow(row *html.Node, baseURL string) *certificateMetadata {
	var cells []string
	var downloadLink string

	// Helper function to extract text from a cell
	extractCell := func(node *html.Node) string {
		var text string
		var collectText func(*html.Node)
		collectText = func(n *html.Node) {
			if n.Type == html.TextNode {
				text += n.Data
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				collectText(c)
			}
		}
		collectText(node)
		return text
	}

	// Extract all cell contents and look for download links
	for td := row.FirstChild; td != nil; td = td.NextSibling {
		if td.Type == html.ElementNode && (td.Data == "td" || td.Data == "th") {
			cellText := extractCell(td)
			cells = append(cells, strings.TrimSpace(cellText))

			// Look for download link in this cell
			if link := findDownloadLink(td); link != "" {
				downloadLink = link
			}
		}
	}

	// Skip header rows and rows without download links
	if downloadLink == "" || len(cells) < 3 {
		return nil
	}

	// Table structure typically has:
	// - Country name column
	// - Certificate name/ID column (often has the download link)
	// - Expiration date column
	// We extract country name and expiration date from the cells

	var country string
	var expirationDate string

	for _, cell := range cells {
		if isDateString(cell) {
			expirationDate = parseDateString(cell)
		} else if cell != "" && !strings.Contains(cell, "Download") && len(cell) > 2 {
			// Heuristic: country name is typically a non-empty string that's not a date,
			// not "Download" text, and has more than 2 characters
			// We take the first such value that looks like a country name
			if country == "" && !strings.Contains(strings.ToLower(cell), "certificate") {
				country = cell
			}
		}
	}

	name := extractCertificateName(downloadLink)
	fullURL := baseURL + "/" + downloadLink

	return &certificateMetadata{
		name:           name,
		downloadURL:    fullURL,
		country:        country,
		expirationDate: expirationDate,
	}
}

// findDownloadLink finds the certificate download link in a table cell
func findDownloadLink(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && (strings.Contains(attr.Val, "BIN_DT") || strings.Contains(attr.Val, "BIN_ST")) {
				return attr.Val
			}
			if attr.Key == "title" && strings.Contains(attr.Val, "Download certificate file") {
				for _, a := range node.Attr {
					if a.Key == "href" {
						return a.Val
					}
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if link := findDownloadLink(c); link != "" {
			return link
		}
	}

	return ""
}

// isDateString checks if a string looks like a date
func isDateString(s string) bool {
	s = strings.TrimSpace(s)
	// Match DD/MM/YYYY or YYYY-MM-DD patterns
	return len(s) == 10 && (strings.Count(s, "/") == 2 || strings.Count(s, "-") == 2)
}

// parseDateString converts various date formats to RFC3339 UTC format
func parseDateString(s string) string {
	s = strings.TrimSpace(s)

	var dateStr string

	// Try DD/MM/YYYY format
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) == 3 {
			// Convert DD/MM/YYYY to YYYY-MM-DD
			dateStr = fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])
		}
	} else if strings.Contains(s, "-") {
		// Try YYYY-MM-DD format
		parts := strings.Split(s, "-")
		if len(parts) == 3 && len(parts[0]) == 4 {
			dateStr = s
		}
	}

	if dateStr == "" {
		return ""
	}

	// Parse the date and convert to RFC3339 format in UTC (at midnight)
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ""
	}

	// Convert to UTC and format as RFC3339
	return t.UTC().Format(time.RFC3339)
}

func extractCertificateName(url string) string {
	parts := strings.Split(url, "=")
	if len(parts) < 2 {
		return filepath.Base(url)
	}

	name := parts[len(parts)-1]
	name = strings.TrimPrefix(name, "iot_doc/")
	name = strings.TrimPrefix(name, "erca_of_doc/")

	return name
}

func downloadCertificate(url string, maxSize int) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 && data[0] == '<' {
		return nil, fmt.Errorf("received HTML instead of certificate data")
	}

	if len(data) > maxSize {
		return nil, fmt.Errorf("certificate too large: %d bytes (max %d)", len(data), maxSize)
	}

	return data, nil
}

func writeIndexFile(index *CertificateIndex, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
