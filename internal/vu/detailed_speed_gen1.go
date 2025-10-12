package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalDetailedSpeedGen1 parses Gen1 Detailed Speed data from the complete transfer value.
//
// Gen1 Detailed Speed structure (from Data Dictionary and Appendix 7, Section 2.2.6.6):
//
// ASN.1 Definition:
//
//	VuDetailedSpeedFirstGen ::= SEQUENCE {
//	    vuDetailedSpeedBlocks      VuDetailedSpeedBlocksFirstGen,
//	    signature                  SignatureFirstGen
//	}
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalDetailedSpeedGen1(value []byte) (*vuv1.DetailedSpeedGen1, error) {
	detailedSpeed := &vuv1.DetailedSpeedGen1{}
	detailedSpeed.SetRawData(value)

	// TODO: Implement full semantic parsing
	// For now, validate that we have enough data for the structure
	if len(value) < 128 { // At minimum, signature is 128 bytes
		return nil, fmt.Errorf("insufficient data for Detailed Speed Gen1")
	}

	// Store the signature (last 128 bytes)
	signatureStart := len(value) - 128
	detailedSpeed.SetSignature(value[signatureStart:])

	return detailedSpeed, nil
}

// appendDetailedSpeedGen1 marshals Gen1 Detailed Speed data using raw data painting.
func appendDetailedSpeedGen1(dst []byte, detailedSpeed *vuv1.DetailedSpeedGen1) ([]byte, error) {
	if detailedSpeed == nil {
		return nil, fmt.Errorf("detailedSpeed cannot be nil")
	}

	raw := detailedSpeed.GetRawData()
	if len(raw) > 0 {
		return append(dst, raw...), nil
	}

	return nil, fmt.Errorf("cannot marshal Detailed Speed Gen1 without raw_data (semantic marshalling not yet implemented)")
}
