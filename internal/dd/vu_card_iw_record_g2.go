package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuCardIWRecordG2 parses a Generation 2 VuCardIWRecord (131 bytes).
//
// The data type `VuCardIWRecord` is specified in the Data Dictionary, Section 2.177.
//
// ASN.1 Definition (Gen2):
//
//	VuCardIWRecord ::= SEQUENCE {
//	    cardHolderName                     HolderName,
//	    fullCardNumberAndGeneration        FullCardNumberAndGeneration,
//	    cardExpiryDate                     Datef,
//	    cardInsertionTime                  TimeReal,
//	    vehicleOdometerValueAtInsertion    OdometerShort,
//	    cardSlotNumber                     CardSlotNumber,
//	    cardWithdrawalTime                 TimeReal,
//	    vehicleOdometerValueAtWithdrawal   OdometerShort,
//	    previousVehicleInfo                PreviousVehicleInfoGen2,
//	    manualInputFlag                    ManualInputFlag
//	}
//
// Binary Layout (fixed length, 131 bytes):
//   - Bytes 0-71: cardHolderName (HolderName)
//   - Bytes 72-90: fullCardNumberAndGeneration (FullCardNumberAndGeneration, 19 bytes)
//   - Bytes 91-94: cardExpiryDate (Datef)
//   - Bytes 95-98: cardInsertionTime (TimeReal)
//   - Bytes 99-101: vehicleOdometerValueAtInsertion (OdometerShort)
//   - Byte 102: cardSlotNumber (CardSlotNumber)
//   - Bytes 103-106: cardWithdrawalTime (TimeReal)
//   - Bytes 107-109: vehicleOdometerValueAtWithdrawal (OdometerShort)
//   - Bytes 110-129: previousVehicleInfo (PreviousVehicleInfoGen2)
//   - Byte 130: manualInputFlag (ManualInputFlag)
func (opts UnmarshalOptions) UnmarshalVuCardIWRecordG2(data []byte) (*ddv1.VuCardIWRecordG2, error) {
	const (
		idxCardHolderName       = 0
		idxFullCardNumber       = 72
		idxCardExpiryDate       = 91
		idxCardInsertionTime    = 95
		idxOdometerAtInsertion  = 99
		idxCardSlotNumber       = 102
		idxCardWithdrawalTime   = 103
		idxOdometerAtWithdrawal = 107
		idxPreviousVehicleInfo  = 110
		idxManualInputFlag      = 130
		lenVuCardIWRecordG2     = 131

		lenHolderName                  = 72
		lenFullCardNumberAndGeneration = 19
		lenDatef                       = 4
		lenTimeReal                    = 4
		lenOdometerShort               = 3
		lenCardSlotNumber              = 1
		lenPreviousVehicleInfoG2       = 20
		lenManualInputFlag             = 1
	)

	if len(data) != lenVuCardIWRecordG2 {
		return nil, fmt.Errorf("invalid data length for VuCardIWRecordG2: got %d, want %d", len(data), lenVuCardIWRecordG2)
	}

	record := &ddv1.VuCardIWRecordG2{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// cardHolderName (72 bytes)
	holderName, err := opts.UnmarshalHolderName(data[idxCardHolderName : idxCardHolderName+lenHolderName])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// fullCardNumberAndGeneration (20 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxFullCardNumber : idxFullCardNumber+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number and generation: %w", err)
	}
	record.SetFullCardNumber(fullCardNumber)

	// cardExpiryDate (4 bytes)
	expiryDate, err := opts.UnmarshalDate(data[idxCardExpiryDate : idxCardExpiryDate+lenDatef])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card expiry date: %w", err)
	}
	record.SetCardExpiryDate(expiryDate)

	// cardInsertionTime (4 bytes)
	insertionTime, err := opts.UnmarshalTimeReal(data[idxCardInsertionTime : idxCardInsertionTime+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card insertion time: %w", err)
	}
	record.SetCardInsertionTime(insertionTime)

	// vehicleOdometerValueAtInsertion (3 bytes)
	odometerAtInsertion, err := opts.UnmarshalOdometer(data[idxOdometerAtInsertion : idxOdometerAtInsertion+lenOdometerShort])
	if err != nil {
		return nil, fmt.Errorf("unmarshal odometer at insertion: %w", err)
	}
	record.SetOdometerAtInsertionKm(int32(odometerAtInsertion))

	// cardSlotNumber (1 byte)
	cardSlotNumber, err := UnmarshalEnum[ddv1.CardSlotNumber](data[idxCardSlotNumber])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card slot number: %w", err)
	}
	record.SetCardSlotNumber(cardSlotNumber)

	// cardWithdrawalTime (4 bytes)
	withdrawalTime, err := opts.UnmarshalTimeReal(data[idxCardWithdrawalTime : idxCardWithdrawalTime+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card withdrawal time: %w", err)
	}
	record.SetCardWithdrawalTime(withdrawalTime)

	// vehicleOdometerValueAtWithdrawal (3 bytes)
	odometerAtWithdrawal, err := opts.UnmarshalOdometer(data[idxOdometerAtWithdrawal : idxOdometerAtWithdrawal+lenOdometerShort])
	if err != nil {
		return nil, fmt.Errorf("unmarshal odometer at withdrawal: %w", err)
	}
	record.SetOdometerAtWithdrawalKm(int32(odometerAtWithdrawal))

	// previousVehicleInfo (20 bytes)
	previousVehicleInfo, err := opts.UnmarshalPreviousVehicleInfoG2(data[idxPreviousVehicleInfo : idxPreviousVehicleInfo+lenPreviousVehicleInfoG2])
	if err != nil {
		return nil, fmt.Errorf("unmarshal previous vehicle info: %w", err)
	}
	record.SetPreviousVehicleInfo(previousVehicleInfo)

	// manualInputFlag (1 byte)
	manualInputFlag := data[idxManualInputFlag] != 0
	record.SetManualInputFlag(manualInputFlag)

	return record, nil
}

// MarshalVuCardIWRecordG2 marshals a Generation 2 VuCardIWRecord (132 bytes) to bytes.
func (opts MarshalOptions) MarshalVuCardIWRecordG2(record *ddv1.VuCardIWRecordG2) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuCardIWRecordG2 = 131

	// Use raw data painting strategy if available
	var canvas [lenVuCardIWRecordG2]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenVuCardIWRecordG2 {
			return nil, fmt.Errorf("invalid raw_data length for VuCardIWRecordG2: got %d, want %d", len(rawData), lenVuCardIWRecordG2)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// cardHolderName (72 bytes)
	holderNameBytes, err := opts.MarshalHolderName(record.GetCardHolderName())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card holder name: %w", err)
	}
	copy(canvas[offset:offset+72], holderNameBytes)
	offset += 72

	// fullCardNumberAndGeneration (19 bytes)
	fullCardNumberBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetFullCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal full card number and generation: %w", err)
	}
	copy(canvas[offset:offset+19], fullCardNumberBytes)
	offset += 19

	// cardExpiryDate (4 bytes)
	expiryDateBytes, err := opts.MarshalDate(record.GetCardExpiryDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card expiry date: %w", err)
	}
	copy(canvas[offset:offset+4], expiryDateBytes)
	offset += 4

	// cardInsertionTime (4 bytes)
	insertionTimeBytes, err := opts.MarshalTimeReal(record.GetCardInsertionTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card insertion time: %w", err)
	}
	copy(canvas[offset:offset+4], insertionTimeBytes)
	offset += 4

	// vehicleOdometerValueAtInsertion (3 bytes)
	odometerAtInsertionBytes, err := opts.MarshalOdometer(record.GetOdometerAtInsertionKm())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal odometer at insertion: %w", err)
	}
	copy(canvas[offset:offset+3], odometerAtInsertionBytes)
	offset += 3

	// cardSlotNumber (1 byte)
	cardSlotNumberByte, _ := MarshalEnum(record.GetCardSlotNumber())
	canvas[offset] = cardSlotNumberByte
	offset += 1

	// cardWithdrawalTime (4 bytes)
	withdrawalTimeBytes, err := opts.MarshalTimeReal(record.GetCardWithdrawalTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card withdrawal time: %w", err)
	}
	copy(canvas[offset:offset+4], withdrawalTimeBytes)
	offset += 4

	// vehicleOdometerValueAtWithdrawal (3 bytes)
	odometerAtWithdrawalBytes, err := opts.MarshalOdometer(record.GetOdometerAtWithdrawalKm())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal odometer at withdrawal: %w", err)
	}
	copy(canvas[offset:offset+3], odometerAtWithdrawalBytes)
	offset += 3

	// previousVehicleInfo (20 bytes)
	previousVehicleInfoBytes, err := opts.MarshalPreviousVehicleInfoG2(record.GetPreviousVehicleInfo())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal previous vehicle info: %w", err)
	}
	copy(canvas[offset:offset+20], previousVehicleInfoBytes)
	offset += 20

	// manualInputFlag (1 byte)
	if record.GetManualInputFlag() {
		canvas[offset] = 1
	} else {
		canvas[offset] = 0
	}

	return canvas[:], nil
}
