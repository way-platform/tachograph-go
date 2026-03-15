package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfTechnicalData dispatches to generation-specific size calculation.
func sizeOfTechnicalData(data []byte, transferType vuv1.TransferType) (totalSize, signatureSize int, err error) {
	switch transferType {
	case vuv1.TransferType_TECHNICAL_DATA_GEN1:
		return sizeOfTechnicalDataGen1(data)
	case vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
		return sizeOfTechnicalDataGen2V1(data)
	case vuv1.TransferType_TECHNICAL_DATA_GEN2_V2:
		return sizeOfTechnicalDataGen2V2(data)
	default:
		return 0, 0, fmt.Errorf("unsupported transfer type for TechnicalData: %v", transferType)
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
func sizeOfTechnicalDataGen1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// VuIdentification: 116 bytes (fixed structure, per Data Dictionary 2.205)
	offset += 116

	// SensorPaired: 20 bytes (fixed structure)
	offset += 20

	// VuCalibrationData: 1 byte count + variable calibration records
	if len(data[offset:]) < 1 {
		return 0, 0, fmt.Errorf("insufficient data for noOfVuCalibrationRecords")
	}
	noOfVuCalibrationRecords := data[offset]
	offset += 1

	// Each VuCalibrationRecordFirstGen: 167 bytes (per Data Dictionary 2.174)
	const vuCalibrationRecordSize = 167
	offset += int(noOfVuCalibrationRecords) * vuCalibrationRecordSize

	// Signature: 128 bytes for Gen1 RSA
	const gen1SignatureSize = 128
	offset += gen1SignatureSize

	return offset, gen1SignatureSize, nil
}

// sizeOfTechnicalDataGen2V1 calculates size by parsing all Gen2 V1 RecordArrays.
func sizeOfTechnicalDataGen2V1(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0

	// VuIdentificationRecordArray
	size, sizeErr := sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuIdentificationRecordArray: %w", sizeErr)
	}
	offset += size

	// SensorPairedRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("SensorPairedRecordArray: %w", sizeErr)
	}
	offset += size

	// VuCalibrationRecordArray
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("VuCalibrationRecordArray: %w", sizeErr)
	}
	offset += size

	// SignatureRecordArray (last)
	size, sizeErr = sizeOfRecordArray(data, offset)
	if sizeErr != nil {
		return 0, 0, fmt.Errorf("SignatureRecordArray: %w", sizeErr)
	}
	signatureSizeGen2 := size
	offset += size

	return offset, signatureSizeGen2, nil
}

// sizeOfTechnicalDataGen2V2 calculates size by scanning all Gen2 V2 RecordArrays generically.
//
// Real devices may include unexpected arrays (e.g. VuCard) between the standard ones,
// so we scan until the end of data rather than expecting a fixed number of arrays.
func sizeOfTechnicalDataGen2V2(data []byte) (totalSize, signatureSize int, err error) {
	offset := 0
	lastSize := 0
	for offset < len(data) {
		size, sizeErr := sizeOfRecordArray(data, offset)
		if sizeErr != nil {
			return 0, 0, fmt.Errorf("RecordArray at offset %d: %w", offset, sizeErr)
		}
		lastSize = size
		offset += size
	}
	if lastSize == 0 {
		return 0, 0, fmt.Errorf("no RecordArrays found in TechnicalDataGen2V2")
	}
	return offset, lastSize, nil
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
