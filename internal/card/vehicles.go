package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardVehiclesUsed unmarshals vehicles used data from a card EF.
//
// The data type `CardVehicleRecord` is specified in the Data Dictionary, Section 2.6.
//
// ASN.1 Definition:
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBeginKm       OdometerShort,
//	    vehicleOdometerEndKm         OdometerShort,
//	    vehicleFirstUse              TimeReal,
//	    vehicleLastUse               TimeReal,
//	    vehicleRegistration          VehicleRegistrationIdentification,
//	    vehicleIdentificationNumber  VehicleIdentificationNumber OPTIONAL,
//	    vehicleRegistrationNation    NationNumeric OPTIONAL,
//	    vehicleRegistrationNumber    RegistrationNumber OPTIONAL
//	}
func unmarshalCardVehiclesUsed(data []byte) (*cardv1.VehiclesUsed, error) {
	const (
		lenMinEfVehiclesUsed = 2 // Minimum EF_VehiclesUsed record size
	)

	if len(data) < lenMinEfVehiclesUsed {
		return nil, fmt.Errorf("insufficient data for vehicles used: got %d bytes, need at least %d", len(data), lenMinEfVehiclesUsed)
	}

	var target cardv1.VehiclesUsed
	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse vehicle records
	var records []*cardv1.VehiclesUsed_Record
	recordSize := determineVehicleRecordSize(data)

	for r.Len() >= recordSize {
		record, err := parseVehicleRecord(r, recordSize)
		if err != nil {
			break // Stop parsing on error, but return what we have
		}
		records = append(records, record)
	}

	target.SetRecords(records)
	return &target, nil
}

// determineVehicleRecordSize determines the size of vehicle records based on data length
// Gen1: 31 bytes per record, Gen2: 48 bytes per record
func determineVehicleRecordSize(data []byte) int {
	remainingData := len(data) - 2 // Subtract pointer size

	// Try to determine generation based on typical record counts and sizes
	if remainingData%48 == 0 {
		return 48 // Gen2 record size
	}
	return 31 // Gen1 record size (default)
}

// parseVehicleRecord parses a single vehicle record
func parseVehicleRecord(r *bytes.Reader, recordSize int) (*cardv1.VehiclesUsed_Record, error) {
	if r.Len() < recordSize {
		return nil, fmt.Errorf("insufficient data for vehicle record")
	}

	record := &cardv1.VehiclesUsed_Record{}

	// Read odometer begin (3 bytes)
	odometerBeginBytes := make([]byte, 3)
	if _, err := r.Read(odometerBeginBytes); err != nil {
		return nil, fmt.Errorf("failed to read odometer begin: %w", err)
	}
	record.SetVehicleOdometerBeginKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerBeginBytes...))))

	// Read odometer end (3 bytes)
	odometerEndBytes := make([]byte, 3)
	if _, err := r.Read(odometerEndBytes); err != nil {
		return nil, fmt.Errorf("failed to read odometer end: %w", err)
	}
	record.SetVehicleOdometerEndKm(int32(binary.BigEndian.Uint32(append([]byte{0}, odometerEndBytes...))))

	// Read vehicle first use (4 bytes)
	firstUseBytes := make([]byte, 4)
	if _, err := r.Read(firstUseBytes); err != nil {
		return nil, fmt.Errorf("failed to read vehicle first use: %w", err)
	}
	vehicleFirstUse, err := dd.UnmarshalTimeReal(firstUseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle first use: %w", err)
	}
	record.SetVehicleFirstUse(vehicleFirstUse)

	// Read vehicle last use (4 bytes)
	lastUseBytes := make([]byte, 4)
	if _, err := r.Read(lastUseBytes); err != nil {
		return nil, fmt.Errorf("failed to read vehicle last use: %w", err)
	}
	vehicleLastUse, err := dd.UnmarshalTimeReal(lastUseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle last use: %w", err)
	}
	record.SetVehicleLastUse(vehicleLastUse)

	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes registration number)
	vehicleRegBytes := make([]byte, 15)
	if _, err := r.Read(vehicleRegBytes); err != nil {
		return nil, fmt.Errorf("failed to read vehicle registration: %w", err)
	}
	vehicleReg, err := dd.UnmarshalVehicleRegistration(vehicleRegBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration: %w", err)
	}
	record.SetVehicleRegistration(vehicleReg)

	// Read VU data block counter (2 bytes)
	var vuDataBlockCounter uint16
	if err := binary.Read(r, binary.BigEndian, &vuDataBlockCounter); err != nil {
		return nil, fmt.Errorf("failed to read VU data block counter: %w", err)
	}

	// Convert to BCD bytes (2 bytes)
	bcdBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bcdBytes, vuDataBlockCounter)
	bcdCounter, err := dd.UnmarshalBcdString(bcdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create BCD string for VU data block counter: %w", err)
	}
	record.SetVuDataBlockCounter(bcdCounter)

	// For Gen2 records, read VIN (17 bytes IA5String)
	if recordSize == 48 {
		vinBytes := make([]byte, 17)
		if _, err := r.Read(vinBytes); err != nil {
			return nil, fmt.Errorf("failed to read vehicle identification number: %w", err)
		}
		vin, err := dd.UnmarshalIA5StringValue(vinBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vehicle identification number: %w", err)
		}
		record.SetVehicleIdentificationNumber(vin)
	}

	return record, nil
}

// AppendVehiclesUsed appends the binary representation of VehiclesUsed to dst.
//
// The data type `CardVehicleRecord` is specified in the Data Dictionary, Section 2.6.
//
// ASN.1 Definition:
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBeginKm       OdometerShort,
//	    vehicleOdometerEndKm         OdometerShort,
//	    vehicleFirstUse              TimeReal,
//	    vehicleLastUse               TimeReal,
//	    vehicleRegistration          VehicleRegistrationIdentification,
//	    vehicleIdentificationNumber  VehicleIdentificationNumber OPTIONAL,
//	    vehicleRegistrationNation    NationNumeric OPTIONAL,
//	    vehicleRegistrationNumber    RegistrationNumber OPTIONAL
//	}
func appendVehiclesUsed(dst []byte, vu *cardv1.VehiclesUsed) ([]byte, error) {
	if vu == nil {
		return dst, nil
	}
	dst = binary.BigEndian.AppendUint16(dst, uint16(vu.GetNewestRecordIndex()))

	var err error
	for _, rec := range vu.GetRecords() {
		dst, err = appendVehicleRecord(dst, rec)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendVehicleRecord appends a single vehicle record to dst.
//
// The data type `CardVehicleRecord` is specified in the Data Dictionary, Section 2.6.
//
// ASN.1 Definition:
//
//	CardVehicleRecord ::= SEQUENCE {
//	    vehicleOdometerBeginKm       OdometerShort,                    -- 3 bytes
//	    vehicleOdometerEndKm         OdometerShort,                    -- 3 bytes
//	    vehicleFirstUse              TimeReal,                         -- 4 bytes
//	    vehicleLastUse               TimeReal,                         -- 4 bytes
//	    vehicleRegistration          VehicleRegistrationIdentification, -- 15 bytes
//	    vehicleIdentificationNumber  VehicleIdentificationNumber OPTIONAL,
//	    vehicleRegistrationNation    NationNumeric OPTIONAL,
//	    vehicleRegistrationNumber    RegistrationNumber OPTIONAL
//	}
func appendVehicleRecord(dst []byte, rec *cardv1.VehiclesUsed_Record) ([]byte, error) {
	if rec == nil {
		return append(dst, make([]byte, 31)...), nil
	}
	dst = dd.AppendOdometer(dst, uint32(rec.GetVehicleOdometerBeginKm()))
	dst = dd.AppendOdometer(dst, uint32(rec.GetVehicleOdometerEndKm()))

	var err error
	dst, err = dd.AppendTimeReal(dst, rec.GetVehicleFirstUse())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle first use: %w", err)
	}
	dst, err = dd.AppendTimeReal(dst, rec.GetVehicleLastUse())
	if err != nil {
		return nil, fmt.Errorf("failed to append vehicle last use: %w", err)
	}

	// Vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	// Note: Odometer fields are 3 bytes each (OdometerShort) for Gen1 cards
	dst, err = dd.AppendVehicleRegistration(dst, rec.GetVehicleRegistration())
	if err != nil {
		return nil, err
	}
	dst, err = dd.AppendBcdString(dst, rec.GetVuDataBlockCounter())
	if err != nil {
		return nil, fmt.Errorf("failed to append vu data block counter: %w", err)
	}
	return dst, nil
}
