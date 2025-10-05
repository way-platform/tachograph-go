// Package security implements the Common Security Mechanisms specified in
// Appendix 11 of the EU digital tachograph regulation.
//
// This package provides cryptographic operations for certificate verification,
// including:
//
//   - RSA signature recovery using ISO/IEC 9796-2 (Generation 1)
//   - ECDSA signature verification with Brainpool curves (Generation 2)
//   - Root certificate parsing and trust anchor management
//
// The security mechanisms form the foundation of the tachograph PKI hierarchy:
//
//	European Root CA (ERCA)
//	    ↓
//	Member State CA (MSCA)
//	    ↓
//	Equipment Certificates (Cards, Vehicle Units)
//
// Certificate verification follows the chain from equipment certificates up
// through Member State CAs to the trusted European Root CA.
//
// See Appendix 11 "Common Security Mechanisms" in the regulation for the
// complete cryptographic specifications.
package security
