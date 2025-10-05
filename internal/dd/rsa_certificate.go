package dd

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/security"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// UnmarshalRsaCertificate parses an RSA certificate from binary data.
//
// This function delegates to internal/security and converts the result to ddv1 format
// for backwards compatibility with existing card file proto definitions.
//
// See Appendix 11, Section 3.3 for the certificate format specification.
func (opts UnmarshalOptions) UnmarshalRsaCertificate(data []byte) (*ddv1.RsaCertificate, error) {
	secCert, err := security.UnmarshalRsaCertificate(data)
	if err != nil {
		return nil, err
	}
	return ConvertRsaCertificateFromSecurity(secCert)
}

// AppendRsaCertificate marshals an RSA certificate to binary format.
//
// This function delegates to internal/security for the actual marshalling.
//
// See Appendix 11, Section 3.3 for the certificate format specification.
func AppendRsaCertificate(dst []byte, cert *ddv1.RsaCertificate) ([]byte, error) {
	secCert, err := ConvertRsaCertificateToSecurity(cert)
	if err != nil {
		return nil, err
	}
	return security.AppendRsaCertificate(dst, secCert)
}

// ConvertRsaCertificateToSecurity converts a ddv1.RsaCertificate to securityv1.RsaCertificate.
func ConvertRsaCertificateToSecurity(cert *ddv1.RsaCertificate) (*securityv1.RsaCertificate, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	sec := &securityv1.RsaCertificate{}

	// Convert uint64 references to string
	car := cert.GetCertificateAuthorityReference()
	carStr := fmt.Sprintf("%d", car)
	sec.SetCertificateAuthorityReference(carStr)

	chr := cert.GetCertificateHolderReference()
	chrStr := fmt.Sprintf("%d", chr)
	sec.SetCertificateHolderReference(chrStr)

	sec.SetEndOfValidity(cert.GetEndOfValidity())
	sec.SetRsaModulus(cert.GetRsaModulus())
	sec.SetRsaExponent(cert.GetRsaExponent())
	sec.SetRawData(cert.GetRawData())
	sec.SetSignatureValid(cert.GetSignatureValid())

	return sec, nil
}

// ConvertRsaCertificateFromSecurity converts a securityv1.RsaCertificate to ddv1.RsaCertificate.
func ConvertRsaCertificateFromSecurity(cert *securityv1.RsaCertificate) (*ddv1.RsaCertificate, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	dd := &ddv1.RsaCertificate{}

	// Convert string references to uint64
	var car, chr uint64
	if _, err := fmt.Sscanf(cert.GetCertificateAuthorityReference(), "%d", &car); err != nil {
		return nil, fmt.Errorf("invalid CAR format: %w", err)
	}
	dd.SetCertificateAuthorityReference(car)

	if chrStr := cert.GetCertificateHolderReference(); chrStr != "" {
		if _, err := fmt.Sscanf(chrStr, "%d", &chr); err != nil {
			return nil, fmt.Errorf("invalid CHR format: %w", err)
		}
		dd.SetCertificateHolderReference(chr)
	}

	dd.SetEndOfValidity(cert.GetEndOfValidity())
	dd.SetRsaModulus(cert.GetRsaModulus())
	dd.SetRsaExponent(cert.GetRsaExponent())
	dd.SetRawData(cert.GetRawData())
	dd.SetSignatureValid(cert.GetSignatureValid())

	return dd, nil
}
