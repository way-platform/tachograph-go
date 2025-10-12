package vu

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfTechnicalData dispatches to generation-specific size calculation.
func sizeOfTechnicalData(data []byte, transferType vuv1.TransferType) (int, error) {
	switch transferType {
	case vuv1.TransferType_TECHNICAL_DATA_GEN1:
		return sizeOfTechnicalDataGen1(data)
	case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
		return sizeOfTechnicalDataGen2V1(data)
	case vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
		return sizeOfTechnicalDataGen2V2(data)
	default:
		return 0, fmt.Errorf("unsupported transfer type for TechnicalData: %v", transferType)
	}
}

// sizeOfTechnicalDataGen1 calculates total size for Gen1 Technical Data including signature.
//
// Technical Data Gen1 structure (from Appendix 7, Section 2.2.6.6):
// - VuIdentification (Data Dictionary 2.205): 116 bytes
//   - vuManufacturerName: Name = 36 bytes (1 codepage + 35 bytes)
//   - vuManufacturerAddress: Address = 36 bytes (1 codepage + 35 bytes)
//   - vuPartNumber: 16 bytes
//   - vuSerialNumber: ExtendedSerialNumber = 8 bytes (4+2+1+1)
//   - vuSoftwareIdentification: 8 bytes (VuSoftwareVersion 4 + VuSoftInstallationDate 4)
//   - vuManufacturingDate: TimeReal = 4 bytes
//   - vuApprovalNumber: IA5String(SIZE(8)) = 8 bytes
//
// - SensorPaired: 20 bytes
//   - sensorSerialNumber: ExtendedSerialNumber = 8 bytes
//   - sensorApprovalNumber: 8 bytes (Gen1 SIZE(8))
//   - sensorPairingDateFirst: TimeReal = 4 bytes
//
// - VuCalibrationData (Data Dictionary 2.173): 1 byte + (noOfVuCalibrationRecords * 167 bytes)
//   - noOfVuCalibrationRecords: 1 byte
//   - vuCalibrationRecords: SET OF VuCalibrationRecordFirstGen
//   - VuCalibrationRecordFirstGen (Data Dictionary 2.174): 167 bytes
//   - calibrationPurpose: 1 byte
//   - workshopName: Name = 36 bytes
//   - workshopAddress: Address = 36 bytes
//   - workshopCardNumber: FullCardNumber = 18 bytes
//   - workshopCardExpiryDate: TimeReal = 4 bytes
//   - vehicleIdentificationNumber: 17 bytes
//   - vehicleRegistrationIdentification: 15 bytes (1+1+13)
//   - wVehicleCharacteristicConstant: 2 bytes
//   - kConstantOfRecordingEquipment: 2 bytes
//   - lTyreCircumference: 2 bytes
//   - tyreSize: 15 bytes
//   - authorisedSpeed: 1 byte
//   - oldOdometerValue: 3 bytes
//   - newOdometerValue: 3 bytes
//   - oldTimeValue: TimeReal = 4 bytes
//   - newTimeValue: TimeReal = 4 bytes
//   - nextCalibrationDate: TimeReal = 4 bytes
//
// - Signature: 128 bytes (RSA)
func sizeOfTechnicalDataGen1(data []byte) (int, error) {
	offset := 0

	// VuIdentification: 116 bytes (fixed structure, per Data Dictionary 2.205)
	offset += 116

	// SensorPaired: 20 bytes (fixed structure)
	offset += 20

	// VuCalibrationData: 1 byte count + variable calibration records
	if len(data[offset:]) < 1 {
		return 0, fmt.Errorf("insufficient data for noOfVuCalibrationRecords")
	}
	noOfVuCalibrationRecords := data[offset]
	offset += 1

	// Each VuCalibrationRecordFirstGen: 167 bytes (per Data Dictionary 2.174)
	const vuCalibrationRecordSize = 167
	offset += int(noOfVuCalibrationRecords) * vuCalibrationRecordSize

	// Signature: 128 bytes for Gen1 RSA
	offset += 128

	return offset, nil
}

// sizeOfTechnicalDataGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
func sizeOfTechnicalDataGen2V1(data []byte) (int, error) {
	offset := 0

	// VuIdentificationRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuIdentificationRecordArray: %w", err)
	}
	offset += size

	// SensorPairedRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SensorPairedRecordArray: %w", err)
	}
	offset += size

	// VuCalibrationRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuCalibrationRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// sizeOfTechnicalDataGen2V2 calculates size by parsing all Gen2 V2 RecordArrays.
func sizeOfTechnicalDataGen2V2(data []byte) (int, error) {
	offset := 0

	// VuIdentificationRecordArray
	size, err := sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuIdentificationRecordArray: %w", err)
	}
	offset += size

	// SensorPairedRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SensorPairedRecordArray: %w", err)
	}
	offset += size

	// SensorExternalGNSSCoupledRecordArray (Gen2 V2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SensorExternalGNSSCoupledRecordArray: %w", err)
	}
	offset += size

	// VuCalibrationRecordArray
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuCalibrationRecordArray: %w", err)
	}
	offset += size

	// VuITSConsentRecordArray (Gen2 V2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuITSConsentRecordArray: %w", err)
	}
	offset += size

	// VuPowerSupplyInterruptionRecordArray (Gen2 V2+)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("VuPowerSupplyInterruptionRecordArray: %w", err)
	}
	offset += size

	// SignatureRecordArray (last)
	size, err = sizeOfRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("SignatureRecordArray: %w", err)
	}
	offset += size

	return offset, nil
}

// ===== Unmarshal Functions =====

// UnmarshalVuTechnicalData unmarshals VU technical data from a VU transfer.
//
// The data type `VuTechnicalData` is specified in the Data Dictionary, Section 2.2.6.5.
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuIdentification                  VuIdentification,
//	    vuCalibrationData                 VuCalibrationData,
//	    vuCardData                        VuCardData,
//	    signature                         SignatureFirstGen
//	}
//
//	VuTechnicalDataSecondGen ::= SEQUENCE {
//	    vuIdentificationRecordArray       VuIdentificationRecordArray,
//	    vuCalibrationRecordArray          VuCalibrationRecordArray,
//	    vuCardRecordArray                 VuCardRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}
func UnmarshalVuTechnicalData(data []byte, offset int, target *vuv1.TechnicalData, generation int) (int, error) {
	startOffset := offset

	// Set generation
	if generation == 1 {
		target.SetGeneration(ddv1.Generation_GENERATION_1)
	} else {
		target.SetGeneration(ddv1.Generation_GENERATION_2)
	}

	// For now, implement a simplified version that just reads the data
	// This ensures the interface is complete while allowing for future enhancement

	// Read all remaining data
	remainingData, offset, err := readBytesFromBytes(data, offset, len(data)-offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read technical data: %w", err)
	}

	// Set as signature based on generation
	if generation == 1 {
		target.SetSignatureGen1(remainingData)
	} else {
		target.SetSignatureGen2(remainingData)
	}

	return offset - startOffset, nil
}

// AppendVuTechnicalData appends VU technical data to a buffer.
//
// The data type `VuTechnicalData` is specified in the Data Dictionary, Section 2.2.6.5.
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuIdentification                  VuIdentification,
//	    vuCalibrationData                 VuCalibrationData,
//	    vuCardData                        VuCardData,
//	    signature                         SignatureFirstGen
//	}
//
//	VuTechnicalDataSecondGen ::= SEQUENCE {
//	    vuIdentificationRecordArray       VuIdentificationRecordArray,
//	    vuCalibrationRecordArray          VuCalibrationRecordArray,
//	    vuCardRecordArray                 VuCardRecordArray,
//	    signatureRecordArray              SignatureRecordArray
//	}

// appendVuTechnicalDataBytes appends VU technical data to a byte slice
