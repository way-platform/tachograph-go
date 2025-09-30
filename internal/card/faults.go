package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalFaultsData parses the binary data for an EF_Faults_Data record.
//
// The data type `CardFaultData` is specified in the Data Dictionary, Section 2.22.
//
// ASN.1 Definition:
//
//	CardFaultData ::= SEQUENCE OF CardFaultRecord
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                   EventFaultType,                     -- 1 byte
//	    faultBeginTime              TimeReal,                         -- 4 bytes
//	    faultEndTime                TimeReal,                         -- 4 bytes
//	    faultVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
//	}
const (
	// CardFaultRecord size (24 bytes total)
	cardFaultRecordSize = 24
)

// splitCardFaultRecord returns a SplitFunc that splits data into 24-byte fault records
func splitCardFaultRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < cardFaultRecordSize {
		if atEOF {
			return 0, nil, nil // No more complete records, but not an error
		}
		return 0, nil, nil // Need more data
	}

	return cardFaultRecordSize, data[:cardFaultRecordSize], nil
}

func unmarshalFaultsData(data []byte) (*cardv1.FaultsData, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(splitCardFaultRecord)

	var records []*cardv1.FaultsData_Record
	for scanner.Scan() {
		recordData := scanner.Bytes()
		// Check if this is a valid record by examining the fault begin time (first 4 bytes after fault type)
		// Fault type is 1 byte, so fault begin time starts at byte 1
		faultBeginTime := binary.BigEndian.Uint32(recordData[1:5])

		rec := &cardv1.FaultsData_Record{}

		if faultBeginTime == 0 {
			// Non-valid record: preserve original bytes
			rec.SetValid(false)
			rec.SetRawData(recordData)
		} else {
			// Valid record: parse semantic data
			rec.SetValid(true)
			if err := unmarshalFaultRecord(recordData, rec); err != nil {
				return nil, err
			}
		}

		records = append(records, rec)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Use simplified schema with single faults array in chronological order
	var fd cardv1.FaultsData
	fd.SetFaults(records)
	return &fd, nil
}

// UnmarshalFaultRecord parses a single fault record.
//
// The data type `CardFaultRecord` is specified in the Data Dictionary, Section 2.22.
//
// ASN.1 Definition:
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                   EventFaultType,                     -- 1 byte
//	    faultBeginTime              TimeReal,                         -- 4 bytes
//	    faultEndTime                TimeReal,                         -- 4 bytes
//	    faultVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
//	}
func unmarshalFaultRecord(data []byte, rec *cardv1.FaultsData_Record) error {
	const (
		lenFaultType                = 1
		lenFaultBeginTime           = 4
		lenFaultEndTime             = 4
		lenFaultVehicleRegistration = 15
		lenCardFaultRecord          = lenFaultType + lenFaultBeginTime + lenFaultEndTime + lenFaultVehicleRegistration
	)

	if len(data) < lenCardFaultRecord {
		return fmt.Errorf("insufficient data for fault record: got %d bytes, need %d", len(data), lenCardFaultRecord)
	}

	offset := 0

	// Read fault type (1 byte) and convert using generic enum helper
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for fault type")
	}
	faultType := data[offset]
	enumDesc := ddv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED.Descriptor()
	if enumNum, found := dd.GetEnumForProtocolValue(enumDesc, int32(faultType)); found {
		rec.SetFaultType(ddv1.EventFaultType(enumNum))
	}
	offset++

	// Read fault begin time (4 bytes)
	if offset+4 > len(data) {
		return fmt.Errorf("insufficient data for fault begin time")
	}
	faultBeginTime, err := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return fmt.Errorf("failed to parse fault begin time: %w", err)
	}
	rec.SetFaultBeginTime(faultBeginTime)
	offset += 4

	// Read fault end time (4 bytes)
	if offset+4 > len(data) {
		return fmt.Errorf("insufficient data for fault end time")
	}
	faultEndTime, err := dd.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return fmt.Errorf("failed to parse fault end time: %w", err)
	}
	rec.SetFaultEndTime(faultEndTime)
	offset += 4

	// Read vehicle registration nation (1 byte)
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for vehicle registration nation")
	}
	nationByte := data[offset]
	offset++

	// Create VehicleRegistrationIdentification structure
	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	if enumNum, found := dd.GetEnumForProtocolValue(ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor(), int32(nationByte)); found {
		vehicleReg.SetNation(ddv1.NationNumeric(enumNum))
	} else {
		vehicleReg.SetNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}

	// Read vehicle registration number (14 bytes)
	if offset+14 > len(data) {
		return fmt.Errorf("insufficient data for vehicle registration number")
	}
	regNumber, err := dd.UnmarshalIA5StringValue(data[offset : offset+14])
	if err != nil {
		return fmt.Errorf("failed to read vehicle registration number: %w", err)
	}
	// offset += 14 // Not needed as this is the last field
	vehicleReg.SetNumber(regNumber)
	rec.SetFaultVehicleRegistration(vehicleReg)
	return nil
}

// AppendFaultsData appends the binary representation of FaultData to dst.
//
// The data type `CardFaultData` is specified in the Data Dictionary, Section 2.22.
//
// ASN.1 Definition:
//
//	CardFaultData ::= SEQUENCE OF CardFaultRecord
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                         EventFaultType,                     -- 1 byte
//	    faultBeginTime                    TimeReal,                         -- 4 bytes
//	    faultEndTime                      TimeReal,                         -- 4 bytes
//	    faultVehicleRegistration          VehicleRegistrationIdentification -- 15 bytes
//	}
func appendFaultsData(dst []byte, data *cardv1.FaultsData) ([]byte, error) {
	var err error

	// Process faults in their chronological order
	for _, r := range data.GetFaults() {
		dst, err = appendFaultRecord(dst, r)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// AppendFaultRecord appends a single fault record to dst.
//
// The data type `CardFaultRecord` is specified in the Data Dictionary, Section 2.22.
//
// ASN.1 Definition:
//
//	CardFaultRecord ::= SEQUENCE {
//	    faultType                         EventFaultType,                     -- 1 byte
//	    faultBeginTime                    TimeReal,                         -- 4 bytes
//	    faultEndTime                      TimeReal,                         -- 4 bytes
//	    faultVehicleRegistration          VehicleRegistrationIdentification -- 15 bytes
//	}
func appendFaultRecord(dst []byte, record *cardv1.FaultsData_Record) ([]byte, error) {
	if !record.GetValid() {
		return append(dst, record.GetRawData()...), nil
	}

	protocolValue := dd.GetEventFaultTypeProtocolValue(record.GetFaultType(), 0)
	dst = append(dst, byte(protocolValue))
	dst = dd.AppendTimeReal(dst, record.GetFaultBeginTime())
	dst = dd.AppendTimeReal(dst, record.GetFaultEndTime())
	dst, err := dd.AppendVehicleRegistration(dst, record.GetFaultVehicleRegistration())
	if err != nil {
		return nil, err
	}

	return dst, nil
}
