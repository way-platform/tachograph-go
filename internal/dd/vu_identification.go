package dd

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuIdentification parses the VuIdentification structure.
//
// See Data Dictionary, Section 2.205, `VuIdentification`.
//
// ASN.1 Specification (Generation 1):
//
//	VuIdentification ::= SEQUENCE {
//	    vuManufacturerName VuManufacturerName,           -- 36 bytes (StringValue: 1 code page + 35 data)
//	    vuManufacturerAddress VuManufacturerAddress,     -- 36 bytes (StringValue: 1 code page + 35 data)
//	    vuPartNumber VuPartNumber,                       -- 16 bytes (IA5String, no code page)
//	    vuSerialNumber VuSerialNumber,                   -- 8 bytes (ExtendedSerialNumber)
//	    vuSoftwareIdentification VuSoftwareIdentification,-- 8 bytes (4 IA5String + 4 TimeReal)
//	    vuManufacturingDate VuManufacturingDate,         -- 4 bytes (TimeReal)
//	    vuApprovalNumber VuApprovalNumber                -- 8 bytes (Gen1 IA5String), 16 bytes (Gen2)
//	}
//
// Binary Layout:
//   - Generation 1: 116 bytes total (36+36+16+8+8+4+8)
//   - Generation 2: varies (124+ bytes with 16-byte approval number, plus additional fields)
func (opts UnmarshalOptions) UnmarshalVuIdentification(data []byte) (*ddv1.VuIdentification, error) {
	// Minimum size check (Gen1)
	const minLen = 116
	if len(data) < minLen {
		return nil, fmt.Errorf(
			"insufficient data for VuIdentification: got %d bytes, need at least %d",
			len(data), minLen,
		)
	}

	vuIdent := &ddv1.VuIdentification{}
	if opts.PreserveRawData {
		vuIdent.SetRawData(data)
	}

	const (
		idxManufacturerName    = 0
		lenManufacturerName    = 36
		idxManufacturerAddress = 36
		lenManufacturerAddress = 36
		idxPartNumber          = 72
		lenPartNumber          = 16 // IA5String, not StringValue!
		idxSerialNumber        = 88 // 72 + 16
		lenSerialNumber        = 8
		idxSoftwareIdent       = 96  // 88 + 8
		lenSoftwareIdent       = 8   // 4 + 4, not 16!
		idxManufacturingDate   = 104 // 96 + 8
		lenManufacturingDate   = 4
		idxApprovalNumber      = 108 // 104 + 4
	)

	// Parse VU manufacturer name (36 bytes)
	manufacturerName, err := opts.UnmarshalStringValue(
		data[idxManufacturerName : idxManufacturerName+lenManufacturerName],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manufacturer name: %w", err)
	}
	vuIdent.SetManufacturerName(manufacturerName)

	// Parse VU manufacturer address (36 bytes)
	manufacturerAddress, err := opts.UnmarshalStringValue(
		data[idxManufacturerAddress : idxManufacturerAddress+lenManufacturerAddress],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manufacturer address: %w", err)
	}
	vuIdent.SetManufacturerAddress(manufacturerAddress)

	// Parse VU part number (16 bytes, IA5String)
	partNumber, err := opts.UnmarshalIa5StringValue(
		data[idxPartNumber : idxPartNumber+lenPartNumber],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse part number: %w", err)
	}
	vuIdent.SetPartNumber(partNumber)

	// Parse VU serial number (8 bytes)
	serialNumber, err := opts.UnmarshalExtendedSerialNumber(
		data[idxSerialNumber : idxSerialNumber+lenSerialNumber],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse serial number: %w", err)
	}
	vuIdent.SetSerialNumber(serialNumber)

	// Parse VU software identification (16 bytes)
	softwareIdent, err := opts.UnmarshalSoftwareIdentification(
		data[idxSoftwareIdent : idxSoftwareIdent+lenSoftwareIdent],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse software identification: %w", err)
	}
	vuIdent.SetSoftwareIdentification(softwareIdent)

	// Parse VU manufacturing date (4 bytes)
	manufacturingDate, err := opts.UnmarshalTimeReal(
		data[idxManufacturingDate : idxManufacturingDate+lenManufacturingDate],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manufacturing date: %w", err)
	}
	vuIdent.SetManufacturingDate(manufacturingDate)

	// Determine approval number length based on total data length
	// Gen1: 8 bytes (total 147), Gen2: 16 bytes (total 155+)
	remainingBytes := len(data) - idxApprovalNumber
	var lenApprovalNumber int
	if remainingBytes >= 16 {
		lenApprovalNumber = 16 // Gen2
	} else if remainingBytes >= 8 {
		lenApprovalNumber = 8 // Gen1
	} else {
		return nil, fmt.Errorf(
			"insufficient data for approval number: got %d bytes remaining, need at least 8",
			remainingBytes,
		)
	}

	// Parse VU approval number (8 or 16 bytes depending on generation)
	approvalNumber, err := opts.UnmarshalIa5StringValue(
		data[idxApprovalNumber : idxApprovalNumber+lenApprovalNumber],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse approval number: %w", err)
	}
	vuIdent.SetApprovalNumber(approvalNumber)

	return vuIdent, nil
}

// MarshalVuIdentification marshals the VuIdentification structure using raw data painting.
//
// See Data Dictionary, Section 2.205, `VuIdentification`.
func (opts MarshalOptions) MarshalVuIdentification(vuIdent *ddv1.VuIdentification) ([]byte, error) {
	if vuIdent == nil {
		return nil, fmt.Errorf("vuIdent cannot be nil")
	}

	// Determine size based on approval number length
	approvalNumberLen := int(vuIdent.GetApprovalNumber().GetLength())
	var size int
	switch approvalNumberLen {
	case 8:
		size = 116 // Gen1: 36+36+16+8+8+4+8
	case 16:
		size = 124 // Gen2 (without additional Gen2-specific fields): 36+36+16+8+8+4+16
	default:
		return nil, fmt.Errorf(
			"invalid approval number length: got %d, want 8 (Gen1) or 16 (Gen2)",
			approvalNumberLen,
		)
	}

	// Use raw data painting strategy
	canvas := make([]byte, size)
	if raw := vuIdent.GetRawData(); len(raw) > 0 {
		if len(raw) < size {
			return nil, fmt.Errorf(
				"invalid raw_data length for VuIdentification: got %d, want at least %d",
				len(raw), size,
			)
		}
		// Use only the portion we need
		copy(canvas, raw[:size])
	}

	offset := 0

	// Marshal VU manufacturer name (36 bytes)
	manufacturerNameBytes, err := opts.MarshalStringValue(vuIdent.GetManufacturerName())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manufacturer name: %w", err)
	}
	if len(manufacturerNameBytes) != 36 {
		return nil, fmt.Errorf(
			"invalid manufacturer name length: got %d, want 36",
			len(manufacturerNameBytes),
		)
	}
	copy(canvas[offset:offset+36], manufacturerNameBytes)
	offset += 36

	// Marshal VU manufacturer address (36 bytes)
	manufacturerAddressBytes, err := opts.MarshalStringValue(vuIdent.GetManufacturerAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manufacturer address: %w", err)
	}
	if len(manufacturerAddressBytes) != 36 {
		return nil, fmt.Errorf(
			"invalid manufacturer address length: got %d, want 36",
			len(manufacturerAddressBytes),
		)
	}
	copy(canvas[offset:offset+36], manufacturerAddressBytes)
	offset += 36

	// Marshal VU part number (16 bytes, IA5String)
	partNumberBytes, err := opts.MarshalIa5StringValue(vuIdent.GetPartNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal part number: %w", err)
	}
	if len(partNumberBytes) != 16 {
		return nil, fmt.Errorf(
			"invalid part number length: got %d, want 16",
			len(partNumberBytes),
		)
	}
	copy(canvas[offset:offset+16], partNumberBytes)
	offset += 16

	// Marshal VU serial number (8 bytes)
	serialNumberBytes, err := opts.MarshalExtendedSerialNumber(vuIdent.GetSerialNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal serial number: %w", err)
	}
	if len(serialNumberBytes) != 8 {
		return nil, fmt.Errorf(
			"invalid serial number length: got %d, want 8",
			len(serialNumberBytes),
		)
	}
	copy(canvas[offset:offset+8], serialNumberBytes)
	offset += 8

	// Marshal VU software identification (8 bytes)
	softwareIdentBytes, err := opts.MarshalSoftwareIdentification(vuIdent.GetSoftwareIdentification())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal software identification: %w", err)
	}
	if len(softwareIdentBytes) != 8 {
		return nil, fmt.Errorf(
			"invalid software identification length: got %d, want 8",
			len(softwareIdentBytes),
		)
	}
	copy(canvas[offset:offset+8], softwareIdentBytes)
	offset += 8

	// Marshal VU manufacturing date (4 bytes)
	manufacturingDateBytes, err := opts.MarshalTimeReal(vuIdent.GetManufacturingDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manufacturing date: %w", err)
	}
	if len(manufacturingDateBytes) != 4 {
		return nil, fmt.Errorf(
			"invalid manufacturing date length: got %d, want 4",
			len(manufacturingDateBytes),
		)
	}
	copy(canvas[offset:offset+4], manufacturingDateBytes)
	offset += 4

	// Marshal VU approval number (8 or 16 bytes)
	approvalNumberBytes, err := opts.MarshalIa5StringValue(vuIdent.GetApprovalNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal approval number: %w", err)
	}
	if len(approvalNumberBytes) != approvalNumberLen {
		return nil, fmt.Errorf(
			"invalid approval number length: got %d, want %d",
			len(approvalNumberBytes), approvalNumberLen,
		)
	}
	copy(canvas[offset:offset+approvalNumberLen], approvalNumberBytes)
	offset += approvalNumberLen

	if offset != size {
		return nil, fmt.Errorf(
			"VuIdentification marshalling size mismatch: wrote %d bytes, expected %d",
			offset, size,
		)
	}

	return canvas, nil
}

// AnonymizeVuIdentification anonymizes VU identification data.
func (opts AnonymizeOptions) AnonymizeVuIdentification(ident *ddv1.VuIdentification) *ddv1.VuIdentification {
	if ident == nil {
		return nil
	}

	result := proto.Clone(ident).(*ddv1.VuIdentification)

	// Anonymize manufacturer name
	result.SetManufacturerName(opts.AnonymizeStringValue(ident.GetManufacturerName()))

	// Anonymize manufacturer address
	result.SetManufacturerAddress(opts.AnonymizeStringValue(ident.GetManufacturerAddress()))

	// Anonymize part number
	result.SetPartNumber(opts.AnonymizeIa5StringValue(ident.GetPartNumber()))

	// Anonymize serial number (ExtendedSerialNumber)
	// Preserve equipment type and manufacturer code (structural info) but anonymize the serial number
	serialNum := &ddv1.ExtendedSerialNumber{}
	if origSerial := ident.GetSerialNumber(); origSerial != nil {
		serialNum.SetType(origSerial.GetType())
		serialNum.SetManufacturerCode(origSerial.GetManufacturerCode())
	}
	serialNum.SetSerialNumber(0)
	result.SetSerialNumber(serialNum)

	// Anonymize software identification
	if softwareIdent := ident.GetSoftwareIdentification(); softwareIdent != nil {
		anonymizedSoftware := &ddv1.SoftwareIdentification{}
		anonymizedSoftware.SetSoftwareVersion(opts.AnonymizeIa5StringValue(softwareIdent.GetSoftwareVersion()))
		// Keep installation date as-is or anonymize if needed
		if softwareIdent.GetSoftwareInstallationDate() != nil {
			anonymizedSoftware.SetSoftwareInstallationDate(softwareIdent.GetSoftwareInstallationDate())
		}
		result.SetSoftwareIdentification(anonymizedSoftware)
	}

	// Anonymize approval number (8 bytes IA5String for Gen1)
	result.SetApprovalNumber(NewIa5StringValue(8, "TEST0001"))

	// Clear raw_data
	result.ClearRawData()

	return result
}
