package dd

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/security"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
)

// UnmarshalEccCertificate parses an ECC certificate from binary data.
//
// This function delegates to internal/security and converts the result to ddv1 format
// for backwards compatibility with existing card file proto definitions.
//
// See Appendix 11, Section 9.3.2 for the certificate format specification.
func (opts UnmarshalOptions) UnmarshalEccCertificate(data []byte) (*ddv1.EccCertificate, error) {
	secCert, err := security.UnmarshalEccCertificate(data)
	if err != nil {
		return nil, err
	}
	return ConvertEccCertificateFromSecurity(secCert)
}

// AppendEccCertificate marshals an ECC certificate to binary format.
//
// This function delegates to internal/security for the actual marshalling.
//
// See Appendix 11, Section 9.3.2 for the certificate format specification.
func AppendEccCertificate(dst []byte, cert *ddv1.EccCertificate) ([]byte, error) {
	secCert, err := ConvertEccCertificateToSecurity(cert)
	if err != nil {
		return nil, err
	}
	return security.AppendEccCertificate(dst, secCert)
}

// ConvertEccCertificateToSecurity converts a ddv1.EccCertificate to securityv1.EccCertificate.
func ConvertEccCertificateToSecurity(cert *ddv1.EccCertificate) (*securityv1.EccCertificate, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	sec := &securityv1.EccCertificate{}

	sec.SetCertificateProfileIdentifier(cert.GetCertificateProfileIdentifier())

	// Convert uint64 references to string
	car := cert.GetCertificateAuthorityReference()
	carStr := fmt.Sprintf("%d", car)
	sec.SetCertificateAuthorityReference(carStr)

	sec.SetCertificateHolderAuthorisation(cert.GetCertificateHolderAuthorisation())

	// Convert public key
	if pk := cert.GetPublicKey(); pk != nil {
		secPK := &securityv1.EccCertificate_PublicKey{}
		secPK.SetDomainParametersOid(pk.GetDomainParametersOid())
		secPK.SetPublicPointX(pk.GetPublicPointX())
		secPK.SetPublicPointY(pk.GetPublicPointY())
		sec.SetPublicKey(secPK)
	}

	chr := cert.GetCertificateHolderReference()
	chrStr := fmt.Sprintf("%d", chr)
	sec.SetCertificateHolderReference(chrStr)

	sec.SetCertificateEffectiveDate(cert.GetCertificateEffectiveDate())
	sec.SetCertificateExpirationDate(cert.GetCertificateExpirationDate())

	// Convert signature
	if sig := cert.GetSignature(); sig != nil {
		secSig := &securityv1.EccCertificate_EccSignature{}
		secSig.SetR(sig.GetR())
		secSig.SetS(sig.GetS())
		sec.SetSignature(secSig)
	}

	sec.SetSignatureValid(cert.GetSignatureValid())
	sec.SetRawData(cert.GetRawData())

	return sec, nil
}

// ConvertEccCertificateFromSecurity converts a securityv1.EccCertificate to ddv1.EccCertificate.
func ConvertEccCertificateFromSecurity(cert *securityv1.EccCertificate) (*ddv1.EccCertificate, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	dd := &ddv1.EccCertificate{}

	dd.SetCertificateProfileIdentifier(cert.GetCertificateProfileIdentifier())

	// Convert string references to uint64
	var car, chr uint64
	if _, err := fmt.Sscanf(cert.GetCertificateAuthorityReference(), "%d", &car); err != nil {
		return nil, fmt.Errorf("invalid CAR format: %w", err)
	}
	dd.SetCertificateAuthorityReference(car)

	dd.SetCertificateHolderAuthorisation(cert.GetCertificateHolderAuthorisation())

	// Convert public key
	if pk := cert.GetPublicKey(); pk != nil {
		ddPK := &ddv1.EccCertificate_PublicKey{}
		ddPK.SetDomainParametersOid(pk.GetDomainParametersOid())
		ddPK.SetPublicPointX(pk.GetPublicPointX())
		ddPK.SetPublicPointY(pk.GetPublicPointY())
		dd.SetPublicKey(ddPK)
	}

	if chrStr := cert.GetCertificateHolderReference(); chrStr != "" {
		if _, err := fmt.Sscanf(chrStr, "%d", &chr); err != nil {
			return nil, fmt.Errorf("invalid CHR format: %w", err)
		}
		dd.SetCertificateHolderReference(chr)
	}

	dd.SetCertificateEffectiveDate(cert.GetCertificateEffectiveDate())
	dd.SetCertificateExpirationDate(cert.GetCertificateExpirationDate())

	// Convert signature
	if sig := cert.GetSignature(); sig != nil {
		ddSig := &ddv1.EccCertificate_EccSignature{}
		ddSig.SetR(sig.GetR())
		ddSig.SetS(sig.GetS())
		dd.SetSignature(ddSig)
	}

	dd.SetSignatureValid(cert.GetSignatureValid())
	dd.SetRawData(cert.GetRawData())

	return dd, nil
}
