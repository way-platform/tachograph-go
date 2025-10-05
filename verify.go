package tachograph

import (
	"context"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// VerifyOptions configures the signature verification process.
type VerifyOptions struct {
	// CertificateResolver is used to resolve CA certificates by their Certificate Authority Reference (CAR).
	// If nil, this defaults to using [DefaultCertificateResolver].
	CertificateResolver CertificateResolver
}

// VerifyFile verifies the certificates in a tachograph file.
//
// See [VerifyOptions] if you need more control over the verification process.
func VerifyFile(ctx context.Context, file *tachographv1.File) error {
	return VerifyOptions{}.VerifyFile(ctx, file)
}

// VerifyFile verifies the certificates in a tachograph file.
//
// For driver card files, this function verifies:
//   - Generation 1: Card certificate using the CA certificate
//   - Generation 2: Card sign certificate using the CA certificate
//
// The verification process uses a certificate resolver to fetch CA certificates
// by their Certificate Authority Reference (CAR). If no resolver is configured,
// it defaults to using [DefaultCertificateResolver], which includes embedded
// certificates from EU member states.
//
// For vehicle unit files, certificate verification is not currently implemented
// as VU certificates are stored as raw bytes and require additional parsing.
//
// This function mutates the certificate structures by setting their signature_valid
// fields to true or false based on the verification result.
//
// Returns an error if verification fails for any certificate.
func (o VerifyOptions) VerifyFile(ctx context.Context, file *tachographv1.File) error {
	if file == nil {
		return fmt.Errorf("file cannot be nil")
	}
	if o.CertificateResolver == nil {
		o.CertificateResolver = DefaultCertificateResolver()
	}
	// Create card-level options with the certificate resolver
	cardOpts := card.VerifyOptions{
		CertificateResolver: o.CertificateResolver,
	}
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		return cardOpts.VerifyDriverCardFile(ctx, file.GetDriverCard())
	case tachographv1.File_VEHICLE_UNIT:
		// VU certificate verification not yet implemented
		return nil
	case tachographv1.File_RAW_CARD:
		// Raw card files don't have parsed certificate structures
		return nil
	default:
		return fmt.Errorf("unsupported file type: %v", file.GetType())
	}
}
