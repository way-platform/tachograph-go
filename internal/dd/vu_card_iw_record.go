package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuCardIWRecord parses a Generation 1 VuCardIWRecord (129 bytes).
//
// The data type `VuCardIWRecord` is specified in the Data Dictionary, Section 2.177.
//
// ASN.1 Definition (Gen1):
//
//	VuCardIWRecord ::= SEQUENCE {
//	    cardHolderName                     HolderName,
//	    fullCardNumber                     FullCardNumber,
//	    cardExpiryDate                     Datef,
//	    cardInsertionTime                  TimeReal,
//	    vehicleOdometerValueAtInsertion    OdometerShort,
//	    cardSlotNumber                     CardSlotNumber,
//	    cardWithdrawalTime                 TimeReal,
//	    vehicleOdometerValueAtWithdrawal   OdometerShort,
//	    previousVehicleInfo                PreviousVehicleInfo,
//	    manualInputFlag                    ManualInputFlag
//	}
//
// Binary Layout (fixed length, 129 bytes):
//   - Bytes 0-71: cardHolderName (HolderName)
//   - Bytes 72-89: fullCardNumber (FullCardNumber)
//   - Bytes 90-93: cardExpiryDate (Datef)
//   - Bytes 94-97: cardInsertionTime (TimeReal)
//   - Bytes 98-100: vehicleOdometerValueAtInsertion (OdometerShort)
//   - Byte 101: cardSlotNumber (CardSlotNumber)
//   - Bytes 102-105: cardWithdrawalTime (TimeReal)
//   - Bytes 106-108: vehicleOdometerValueAtWithdrawal (OdometerShort)
//   - Bytes 109-127: previousVehicleInfo (PreviousVehicleInfo)
//   - Byte 128: manualInputFlag (ManualInputFlag)
func (opts UnmarshalOptions) UnmarshalVuCardIWRecord(data []byte) (*ddv1.VuCardIWRecord, error) {
	const (
		idxCardHolderName       = 0
		idxFullCardNumber       = 72
		idxCardExpiryDate       = 90
		idxCardInsertionTime    = 94
		idxOdometerAtInsertion  = 98
		idxCardSlotNumber       = 101
		idxCardWithdrawalTime   = 102
		idxOdometerAtWithdrawal = 106
		idxPreviousVehicleInfo  = 109
		idxManualInputFlag      = 128
		lenVuCardIWRecord       = 129

		lenHolderName          = 72
		lenFullCardNumber      = 18
		lenDatef               = 4
		lenTimeReal            = 4
		lenOdometerShort       = 3
		lenCardSlotNumber      = 1
		lenPreviousVehicleInfo = 19
		lenManualInputFlag     = 1
	)

	if len(data) != lenVuCardIWRecord {
		return nil, fmt.Errorf("invalid data length for VuCardIWRecord: got %d, want %d", len(data), lenVuCardIWRecord)
	}

	record := &ddv1.VuCardIWRecord{}
	record.SetRawData(data)

	// cardHolderName (72 bytes)
	holderName, err := opts.UnmarshalHolderName(data[idxCardHolderName : idxCardHolderName+lenHolderName])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// fullCardNumber (18 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumber(data[idxFullCardNumber : idxFullCardNumber+lenFullCardNumber])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number: %w", err)
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

	// previousVehicleInfo (19 bytes)
	previousVehicleInfo, err := opts.UnmarshalPreviousVehicleInfo(data[idxPreviousVehicleInfo : idxPreviousVehicleInfo+lenPreviousVehicleInfo])
	if err != nil {
		return nil, fmt.Errorf("unmarshal previous vehicle info: %w", err)
	}
	record.SetPreviousVehicleInfo(previousVehicleInfo)

	// manualInputFlag (1 byte)
	manualInputFlag := data[idxManualInputFlag] != 0
	record.SetManualInputFlag(manualInputFlag)

	return record, nil
}

// AppendVuCardIWRecord appends a Generation 1 VuCardIWRecord (129 bytes).
func AppendVuCardIWRecord(dst []byte, record *ddv1.VuCardIWRecord) ([]byte, error) {
	const lenVuCardIWRecord = 129

	// Use raw data painting strategy if available
	var canvas [lenVuCardIWRecord]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenVuCardIWRecord {
			return nil, fmt.Errorf("invalid raw_data length for VuCardIWRecord: got %d, want %d", len(rawData), lenVuCardIWRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// cardHolderName (72 bytes)
	holderNameBytes, err := AppendHolderName(nil, record.GetCardHolderName())
	if err != nil {
		return nil, fmt.Errorf("failed to append card holder name: %w", err)
	}
	copy(canvas[offset:offset+72], holderNameBytes)
	offset += 72

	// fullCardNumber (18 bytes)
	fullCardNumberBytes, err := AppendFullCardNumber(nil, record.GetFullCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to append full card number: %w", err)
	}
	copy(canvas[offset:offset+18], fullCardNumberBytes)
	offset += 18

	// cardExpiryDate (4 bytes)
	expiryDateBytes, err := AppendDate(nil, record.GetCardExpiryDate())
	if err != nil {
		return nil, fmt.Errorf("failed to append card expiry date: %w", err)
	}
	copy(canvas[offset:offset+4], expiryDateBytes)
	offset += 4

	// cardInsertionTime (4 bytes)
	insertionTimeBytes, err := AppendTimeReal(nil, record.GetCardInsertionTime())
	if err != nil {
		return nil, fmt.Errorf("failed to append card insertion time: %w", err)
	}
	copy(canvas[offset:offset+4], insertionTimeBytes)
	offset += 4

	// vehicleOdometerValueAtInsertion (3 bytes)
	odometerAtInsertionBytes := AppendOdometer(nil, uint32(record.GetOdometerAtInsertionKm()))
	copy(canvas[offset:offset+3], odometerAtInsertionBytes)
	offset += 3

	// cardSlotNumber (1 byte)
	cardSlotNumberByte, _ := MarshalEnum(record.GetCardSlotNumber())
	canvas[offset] = cardSlotNumberByte
	offset += 1

	// cardWithdrawalTime (4 bytes)
	withdrawalTimeBytes, err := AppendTimeReal(nil, record.GetCardWithdrawalTime())
	if err != nil {
		return nil, fmt.Errorf("failed to append card withdrawal time: %w", err)
	}
	copy(canvas[offset:offset+4], withdrawalTimeBytes)
	offset += 4

	// vehicleOdometerValueAtWithdrawal (3 bytes)
	odometerAtWithdrawalBytes := AppendOdometer(nil, uint32(record.GetOdometerAtWithdrawalKm()))
	copy(canvas[offset:offset+3], odometerAtWithdrawalBytes)
	offset += 3

	// previousVehicleInfo (19 bytes)
	previousVehicleInfoBytes, err := AppendPreviousVehicleInfo(nil, record.GetPreviousVehicleInfo())
	if err != nil {
		return nil, fmt.Errorf("failed to append previous vehicle info: %w", err)
	}
	copy(canvas[offset:offset+19], previousVehicleInfoBytes)
	offset += 19

	// manualInputFlag (1 byte)
	if record.GetManualInputFlag() {
		canvas[offset] = 1
	} else {
		canvas[offset] = 0
	}

	return append(dst, canvas[:]...), nil
}
