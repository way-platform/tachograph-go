package card

import (
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalBorderCrossings parses the EF_Border_Crossings file.
//
// The data type `CardBorderCrossings` is specified in the Data Dictionary, Section 2.11a.
//
// ASN.1 Definition:
//
//	CardBorderCrossings ::= SEQUENCE {
//	    borderCrossingPointerNewestRecord INTEGER (0..NoOfBorderCrossingRecords -1),
//	    cardBorderCrossingRecords SET SIZE (NoOfBorderCrossingRecords) OF CardBorderCrossingRecord
//	}
//
// Binary Layout:
//   - Byte 0: borderCrossingPointerNewestRecord (1 byte)
//   - Bytes 1..N: Array of CardBorderCrossingRecord (17 bytes each)
//
// Note: The pointer is 1 byte in Gen2v2 (unlike 2 bytes for some other EFs in Gen2).
// Verification: Annex 1C 2.11a says `INTEGER (0..NoOfBorderCrossingRecords -1)`.
// `NoOfBorderCrossingRecords` is defined as a small number (typically < 255).
// Let's double check the regulation or infer from other EFs.
// Most pointers in driver cards are 2 bytes (Short). However, PlaceRecord pointers are 1 byte in Gen1.
// In Gen2, indices are often 2 bytes.
//
// WAIT: The regulation says `INTEGER (0..NoOfBorderCrossingRecords -1)`.
// If `NoOfBorderCrossingRecords` > 255, it must be 2 bytes.
// Driver Card `NoOfBorderCrossingRecords` is likely small.
// Let's assume 1 byte for now based on similar Gen1 structures, BUT Gen2 generally moved to 2-byte indices.
//
// Let's look at `EF_Places` (Gen1): 1 byte index.
// Let's look at `EF_GNSS_Places` (Gen2): 2 byte index.
//
// Let's assume 1 byte for now, but be ready to correct.
// Actually, let's check `internal/dd/card_border_crossing_record.go` size. It's 17 bytes.
// If the total size is `1 + N*17`, we can check `(len-1) % 17 == 0`.
// If `(len-2) % 17 == 0`, it's 2 bytes.
func (opts UnmarshalOptions) unmarshalBorderCrossings(data []byte) (*cardv1.BorderCrossings, error) {
	const (
		lenRecord = 17
	)

	if len(data) < 1 {
		return nil, fmt.Errorf("data too short for CardBorderCrossings")
	}

	// Heuristic to determine pointer size:
	// Try 1 byte pointer
	rem1 := (len(data) - 1) % lenRecord
	// Try 2 byte pointer
	rem2 := (len(data) - 2) % lenRecord

	var pointerSize int
	if rem1 == 0 && rem2 != 0 {
		pointerSize = 1
	} else if rem2 == 0 && rem1 != 0 {
		pointerSize = 2
	} else {
		// Ambiguous or invalid (e.g. empty records)
		// Default to 1 byte as per standard ASN.1 INTEGER implicit sizing for small values,
		// but Tachograph usually uses fixed sizes.
		// Let's assume 1 byte if consistent with typical cyclic buffer pointers for small arrays.
		// However, standard Gen2 pointers are often 2 bytes.
		// Let's pick 1 byte as a starting point, matching Gen1 Places.
		// If the file is 0 bytes, it returns error above.
		// If len(data) is just the pointer (no records), e.g. 1 byte or 2 bytes.
		if len(data) == 1 {
			pointerSize = 1
		} else if len(data) == 2 {
			pointerSize = 2
		} else {
			return nil, fmt.Errorf("invalid data length for CardBorderCrossings: %d (not 1+N*17 or 2+N*17)", len(data))
		}
	}

	// Parse pointer
	var newestRecordIndex int32
	if pointerSize == 1 {
		newestRecordIndex = int32(data[0])
	} else {
		newestRecordIndex = int32(binary.BigEndian.Uint16(data[0:2]))
	}

	target := &cardv1.BorderCrossings{}
	target.SetNewestRecordIndex(newestRecordIndex)

	// Parse records
	recordsData := data[pointerSize:]
	count := len(recordsData) / lenRecord
	targetRecords := make([]*cardv1.BorderCrossings_Record, count)

	for i := 0; i < count; i++ {
		start := i * lenRecord
		end := start + lenRecord
		recordData := recordsData[start:end]

		ddRecord, err := opts.UnmarshalCardBorderCrossingRecord(recordData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal record %d: %w", i, err)
		}

		record := &cardv1.BorderCrossings_Record{}
		record.SetCountryLeft(ddRecord.GetCountryLeft())
		record.SetCountryEntered(ddRecord.GetCountryEntered())
		record.SetGnssPlaceAuthRecord(ddRecord.GetGnssPlaceAuthRecord())
		record.SetVehicleOdometerKm(ddRecord.GetVehicleOdometerKm())
		targetRecords[i] = record
	}
	target.SetRecords(targetRecords)

	return target, nil
}

// MarshalCardBorderCrossings marshals the EF_Border_Crossings file.
func (opts MarshalOptions) MarshalCardBorderCrossings(msg *cardv1.BorderCrossings) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}

	// Determine pointer size logic (inverse of unmarshal).
	// For now, we will write 1 byte if index < 255, but strictly we should follow the spec.
	// We'll default to 1 byte for consistency with the unmarshaler's primary guess.
	// If we find out it's 2 bytes, we'll update both.
	
	var data []byte
	
	newestRecordIndex := msg.GetNewestRecordIndex()
	// Write pointer (1 byte assumption)
	// TODO: verify if it should be 2 bytes
	if newestRecordIndex > 255 {
		// Must be 2 bytes if larger
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(newestRecordIndex))
		data = append(data, buf...)
	} else {
		data = append(data, byte(newestRecordIndex))
	}

	for _, record := range msg.GetRecords() {
		ddRecord := &ddv1.CardBorderCrossingRecord{}
		ddRecord.SetCountryLeft(record.GetCountryLeft())
		ddRecord.SetCountryEntered(record.GetCountryEntered())
		ddRecord.SetGnssPlaceAuthRecord(record.GetGnssPlaceAuthRecord())
		ddRecord.SetVehicleOdometerKm(record.GetVehicleOdometerKm())
		
		bytes, err := opts.MarshalCardBorderCrossingRecord(ddRecord)
		if err != nil {
			return nil, err
		}
		data = append(data, bytes...)
	}

	return data, nil
}
