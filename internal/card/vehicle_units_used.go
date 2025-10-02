package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalVehicleUnitsUsed unmarshals vehicle units used data from a card EF.
//
// The data type `CardVehicleUnitsUsed` is specified in the Data Dictionary, Section 2.40.
//
// ASN.1 Definition:
//
//	CardVehicleUnitsUsed ::= SEQUENCE {
//	    vehicleUnitPointerNewestRecord     INTEGER(0..NoOfCardVehicleUnitRecords-1),
//	    cardVehicleUnitRecords             SET SIZE(NoOfCardVehicleUnitRecords) OF CardVehicleUnitRecord
//	}
//
//	CardVehicleUnitRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,             -- 4 bytes
//	    manufacturerCode                   ManufacturerCode,     -- 1 byte
//	    deviceID                           OCTET STRING(SIZE(1)), -- 1 byte
//	    vuSoftwareVersion                  VuSoftwareVersion     -- 4 bytes
//	}
//
// Binary structure:
//   - 2 bytes: vehicleUnitPointerNewestRecord
//   - N * 10 bytes: fixed-size array of CardVehicleUnitRecord (N determined by data length)
//
// Typical sizes:
//   - Driver cards: 2002 bytes (200 records)
//   - Control cards: 82 bytes (8 records)
func (opts UnmarshalOptions) unmarshalVehicleUnitsUsed(data []byte) (*cardv1.VehicleUnitsUsed, error) {
	const (
		idxNewestRecordPointer         = 0
		lenNewestRecordPointer         = 2
		lenCardVehicleUnitRecord       = 10
		lenCardVehicleUnitsUsedMinimum = lenNewestRecordPointer
	)

	if len(data) < lenCardVehicleUnitsUsedMinimum {
		return nil, fmt.Errorf("invalid data length for CardVehicleUnitsUsed: got %d, want at least %d", len(data), lenCardVehicleUnitsUsedMinimum)
	}

	// Validate that the records section is a multiple of record size
	recordsDataLen := len(data) - lenNewestRecordPointer
	if recordsDataLen%lenCardVehicleUnitRecord != 0 {
		return nil, fmt.Errorf("invalid records data length for CardVehicleUnitsUsed: got %d bytes, not a multiple of %d", recordsDataLen, lenCardVehicleUnitRecord)
	}

	var target cardv1.VehicleUnitsUsed

	// Parse newest record pointer
	newestRecordPointer := binary.BigEndian.Uint16(data[idxNewestRecordPointer:])
	target.SetVehicleUnitPointerNewestRecord(int32(newestRecordPointer))

	// Parse records using bufio.Scanner pattern
	records, err := parseCardVehicleUnitRecords(data[lenNewestRecordPointer:], opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle unit records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// splitCardVehicleUnitRecord is a bufio.SplitFunc for parsing CardVehicleUnitRecord entries.
func splitCardVehicleUnitRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const lenCardVehicleUnitRecord = 10

	if len(data) < lenCardVehicleUnitRecord {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return lenCardVehicleUnitRecord, data[:lenCardVehicleUnitRecord], nil
}

// parseCardVehicleUnitRecords parses the fixed-size array of vehicle unit records.
func parseCardVehicleUnitRecords(data []byte, opts UnmarshalOptions) ([]*cardv1.VehicleUnitsUsed_Record, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(splitCardVehicleUnitRecord)

	var records []*cardv1.VehicleUnitsUsed_Record
	for scanner.Scan() {
		record, err := unmarshalCardVehicleUnitRecord(scanner.Bytes(), opts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal vehicle unit record: %w", err)
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return records, nil
}

// unmarshalCardVehicleUnitRecord unmarshals a single CardVehicleUnitRecord.
//
// Binary structure (10 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 1 byte: ManufacturerCode
//   - 1 byte: deviceID
//   - 4 bytes: VuSoftwareVersion (IA5String)
func unmarshalCardVehicleUnitRecord(data []byte, opts UnmarshalOptions) (*cardv1.VehicleUnitsUsed_Record, error) {
	const (
		idxTimeStamp         = 0
		idxManufacturerCode  = 4
		idxDeviceID          = 5
		idxVuSoftwareVersion = 6
		lenRecord            = 10
	)

	if len(data) != lenRecord {
		return nil, fmt.Errorf("invalid data length for CardVehicleUnitRecord: got %d, want %d", len(data), lenRecord)
	}

	var record cardv1.VehicleUnitsUsed_Record

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse manufacturer code (1 byte)
	record.SetManufacturerCode(int32(data[idxManufacturerCode]))

	// Parse device ID (1 byte)
	record.SetDeviceId(data[idxDeviceID : idxDeviceID+1])

	// Parse VU software version (4 bytes, IA5String)
	record.SetVuSoftwareVersion(data[idxVuSoftwareVersion : idxVuSoftwareVersion+4])

	return &record, nil
}

// appendCardVehicleUnitsUsed appends vehicle units used data to a byte slice.
//
// The data type `CardVehicleUnitsUsed` is specified in the Data Dictionary, Section 2.40.
//
// ASN.1 Definition:
//
//	CardVehicleUnitsUsed ::= SEQUENCE {
//	    vehicleUnitPointerNewestRecord     INTEGER(0..NoOfCardVehicleUnitRecords-1),
//	    cardVehicleUnitRecords             SET SIZE(NoOfCardVehicleUnitRecords) OF CardVehicleUnitRecord
//	}
//
//	CardVehicleUnitRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,             -- 4 bytes
//	    manufacturerCode                   ManufacturerCode,     -- 1 byte
//	    deviceID                           OCTET STRING(SIZE(1)), -- 1 byte
//	    vuSoftwareVersion                  VuSoftwareVersion     -- 4 bytes
//	}
//
// Binary structure:
//   - 2 bytes: vehicleUnitPointerNewestRecord
//   - N * 10 bytes: fixed-size array of CardVehicleUnitRecord
//
// The number of records (N) is determined from the original data, not explicitly stored.
func appendCardVehicleUnitsUsed(dst []byte, vehicleUnits *cardv1.VehicleUnitsUsed) ([]byte, error) {
	if vehicleUnits == nil {
		return dst, nil
	}

	// Append newest record pointer (2 bytes)
	newestRecordPointer := vehicleUnits.GetVehicleUnitPointerNewestRecord()
	dst = binary.BigEndian.AppendUint16(dst, uint16(newestRecordPointer))

	// Append all vehicle unit records
	// The binary format is a fixed-size array, so we write exactly what we have
	records := vehicleUnits.GetRecords()
	for _, record := range records {
		var err error
		dst, err = appendCardVehicleUnitRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append vehicle unit record: %w", err)
		}
	}

	return dst, nil
}

// appendCardVehicleUnitRecord appends a single vehicle unit record to dst.
//
// Binary structure (10 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 1 byte: ManufacturerCode
//   - 1 byte: deviceID
//   - 4 bytes: VuSoftwareVersion (IA5String)
func appendCardVehicleUnitRecord(dst []byte, record *cardv1.VehicleUnitsUsed_Record) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// Append timestamp (TimeReal - 4 bytes)
	var err error
	dst, err = dd.AppendTimeReal(dst, record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to append timestamp: %w", err)
	}

	// Append manufacturer code (1 byte)
	manufacturerCode := record.GetManufacturerCode()
	if manufacturerCode < 0 || manufacturerCode > 255 {
		return nil, fmt.Errorf("invalid manufacturer code: %d", manufacturerCode)
	}
	dst = append(dst, byte(manufacturerCode))

	// Append device ID (1 byte)
	deviceID := record.GetDeviceId()
	if len(deviceID) > 1 {
		return nil, fmt.Errorf("device ID too long: %d bytes", len(deviceID))
	}
	if len(deviceID) == 1 {
		dst = append(dst, deviceID[0])
	} else {
		dst = append(dst, 0x00)
	}

	// Append VU software version (4 bytes)
	vuSoftwareVersion := record.GetVuSoftwareVersion()
	if len(vuSoftwareVersion) > 4 {
		return nil, fmt.Errorf("VU software version too long: %d bytes", len(vuSoftwareVersion))
	}
	if len(vuSoftwareVersion) == 4 {
		dst = append(dst, vuSoftwareVersion...)
	} else {
		// Pad with zeros if shorter than 4 bytes
		padded := make([]byte, 4)
		copy(padded, vuSoftwareVersion)
		dst = append(dst, padded...)
	}

	return dst, nil
}
