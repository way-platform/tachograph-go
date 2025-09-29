package tachograph

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnmarshalVuActivities unmarshals VU activities data from a VU transfer.
//
// The data type `VuActivities` is specified in the Data Dictionary, Section 2.2.6.2.
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    dateOfDay                        TimeReal,
//	    odometerValueMidnight            OdometerValueMidnight,
//	    vuCardIWData                     VuCardIWData,
//	    vuActivityDailyData              VuActivityDailyData,
//	    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
//	    vuSpecificConditionData          VuSpecificConditionData,
//	    signature                        SignatureFirstGen
//	}
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
//	    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                      VuCardIWRecordArray,
//	    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                      VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
//	    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
//	    signatureRecordArray                     SignatureRecordArray
//	}
func unmarshalVuActivities(data []byte, offset int, target *vuv1.Activities, generation int) (int, error) {
	switch generation {
	case 1:
		return unmarshalVuActivitiesGen1(data, offset, target)
	case 2:
		return unmarshalVuActivitiesGen2(data, offset, target)
	default:
		return 0, fmt.Errorf("unsupported generation: %d", generation)
	}
}

// unmarshalVuActivitiesGen1 unmarshals Generation 1 VU activities
func unmarshalVuActivitiesGen1(data []byte, offset int, target *vuv1.Activities) (int, error) {
	startOffset := offset
	target.SetGeneration(ddv1.Generation_GENERATION_1)

	// Read TimeReal (4 bytes) - this is the date of the day
	timeReal, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read date of day: %w", err)
	}
	target.SetDateOfDay(timestamppb.New(time.Unix(timeReal, 0)))

	// Read OdometerValueMidnight (3 bytes)
	odometerValue, offset, err := readOdometerFromBytes(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to read odometer value midnight: %w", err)
	}
	target.SetOdometerMidnightKm(int32(odometerValue))

	// Parse VuCardIWData
	cardIWData, offset, err := parseVuCardIWData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW data: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyData
	activityChanges, offset, err := parseVuActivityDailyData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily data: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodData
	places, offset, err := parseVuPlaceDailyWorkPeriodData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period data: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuSpecificConditionData
	specificConditions, offset, err := parseVuSpecificConditionData(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition data: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Read signature (128 bytes for Gen1)
	signatureBytes, offset, err := readBytesFromBytes(data, offset, 128)
	if err != nil {
		return 0, fmt.Errorf("failed to read signature: %w", err)
	}
	target.SetSignatureGen1(signatureBytes)

	return offset - startOffset, nil
}

// unmarshalVuActivitiesGen2 unmarshals Generation 2 VU activities
func unmarshalVuActivitiesGen2(data []byte, offset int, target *vuv1.Activities) (int, error) {
	startOffset := offset
	target.SetGeneration(ddv1.Generation_GENERATION_2)

	// Gen2 format uses record arrays, each with a header
	// Parse DateOfDayDownloadedRecordArray
	dates, offset, err := parseDateOfDayDownloadedRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse date of day downloaded record array: %w", err)
	}
	if len(dates) > 0 {
		target.SetDateOfDay(dates[0]) // Use first date
	}

	// Parse OdometerValueMidnightRecordArray
	odometerValues, offset, err := parseOdometerValueMidnightRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse odometer value midnight record array: %w", err)
	}
	if len(odometerValues) > 0 {
		target.SetOdometerMidnightKm(odometerValues[0])
	}

	// Parse VuCardIWRecordArray
	cardIWData, offset, err := parseVuCardIWRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse card IW record array: %w", err)
	}
	target.SetCardIwData(cardIWData)

	// Parse VuActivityDailyRecordArray
	activityChanges, offset, err := parseVuActivityDailyRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse activity daily record array: %w", err)
	}
	target.SetActivityChanges(activityChanges)

	// Parse VuPlaceDailyWorkPeriodRecordArray
	places, offset, err := parseVuPlaceDailyWorkPeriodRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse place daily work period record array: %w", err)
	}
	target.SetPlaces(places)

	// Parse VuGNSSADRecordArray (Gen2+)
	gnssRecords, offset, err := parseVuGNSSADRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse GNSS AD record array: %w", err)
	}
	target.SetGnssAccumulatedDriving(gnssRecords)

	// Parse VuSpecificConditionRecordArray
	specificConditions, offset, err := parseVuSpecificConditionRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse specific condition record array: %w", err)
	}
	target.SetSpecificConditions(specificConditions)

	// Try to parse Gen2v2 specific arrays if there's more data
	if offset+10 <= len(data) { // Need some minimum data for arrays
		// Parse VuBorderCrossingRecordArray (Gen2v2+)
		borderCrossings, newOffset, err := parseVuBorderCrossingRecordArray(data, offset)
		if err == nil {
			target.SetBorderCrossings(borderCrossings)
			target.SetVersion(vuv1.Version_VERSION_2)
			offset = newOffset
		}

		// Parse VuLoadUnloadRecordArray (Gen2v2+)
		if offset+5 <= len(data) {
			loadUnloadRecords, newOffset, err := parseVuLoadUnloadRecordArray(data, offset)
			if err == nil {
				target.SetLoadUnloadOperations(loadUnloadRecords)
				offset = newOffset
			}
		}
	}

	// Parse SignatureRecordArray
	signatureBytes, offset, err := parseSignatureRecordArray(data, offset)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signature record array: %w", err)
	}
	target.SetSignatureGen2(signatureBytes)

	return offset - startOffset, nil
}

// Helper functions for parsing different record types
// These are simplified implementations - in a full implementation,
// each would need to properly handle the record array format

// splitVuCardIWRecord splits data into 126-byte VuCardIWRecord records
func splitVuCardIWRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const cardIWRecordSize = 126

	if len(data) < cardIWRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return cardIWRecordSize, data[:cardIWRecordSize], nil
}

func parseVuCardIWData(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWData ::= SEQUENCE {
	//     noOfIWRecords INTEGER(0..255),
	//     vuCardIWRecords SET SIZE(noOfIWRecords) OF VuCardIWRecord -- 126 bytes each
	// }

	// Read number of records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of IW records: %w", err)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuCardIWRecord)

	var records []*vuv1.Activities_CardIWRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuCardIWRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse IW record %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 126

	return records, offset, nil
}

// unmarshalVuCardIWRecord parses a single VuCardIWRecord from a byte slice
func unmarshalVuCardIWRecord(data []byte) (*vuv1.Activities_CardIWRecord, error) {
	// VuCardIWRecord ::= SEQUENCE {
	//     cardHolderName HolderName,                    -- 72 bytes
	//     fullCardNumber FullCardNumber,                -- 19 bytes
	//     cardExpiryDate Datef,                         -- 4 bytes
	//     cardInsertionTime TimeReal,                   -- 4 bytes
	//     vehicleOdometerValueAtInsertion OdometerShort, -- 3 bytes
	//     cardSlotNumber CardSlotNumber,                -- 1 byte
	//     cardWithdrawalTime TimeReal,                  -- 4 bytes
	//     vehicleOdometerValueAtWithdrawal OdometerShort, -- 3 bytes
	//     previousVehicleInfo PreviousVehicleInfo,      -- 15 bytes
	//     manualInputFlag ManualInputFlag               -- 1 byte
	// }

	if len(data) < 126 {
		return nil, fmt.Errorf("insufficient data for card IW record: got %d, need 126", len(data))
	}

	record := &vuv1.Activities_CardIWRecord{}

	// Parse cardHolderName (HolderName - 72 bytes)
	holderNameData := data[0:72]
	holderName, err := unmarshalHolderName(holderNameData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// Parse fullCardNumber (FullCardNumber - 19 bytes)
	fullCardNumberData := data[72:91]
	fullCardNumber, err := unmarshalFullCardNumber(fullCardNumberData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal full card number: %w", err)
	}
	// Create FullCardNumberAndGeneration wrapper
	fullCardNumberAndGeneration := &ddv1.FullCardNumberAndGeneration{}
	fullCardNumberAndGeneration.SetFullCardNumber(fullCardNumber)
	fullCardNumberAndGeneration.SetGeneration(ddv1.Generation_GENERATION_1) // Default to Gen1
	record.SetFullCardNumberAndGeneration(fullCardNumberAndGeneration)

	// Parse cardExpiryDate (Datef - 4 bytes)
	datefData := data[91:95]
	cardExpiryDate, err := readDatef(bytes.NewReader(datefData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse card expiry date: %w", err)
	}
	record.SetCardExpiryDate(cardExpiryDate)

	// Parse cardInsertionTime (TimeReal - 4 bytes)
	insertionTime, _, err := readVuTimeRealFromBytes(data, 95)
	if err != nil {
		return nil, fmt.Errorf("failed to read card insertion time: %w", err)
	}
	record.SetCardInsertionTime(timestamppb.New(time.Unix(insertionTime, 0)))

	// Parse vehicleOdometerValueAtInsertion (OdometerShort - 3 bytes)
	odometerAtInsertion, _, err := readOdometerFromBytes(data, 99)
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer at insertion: %w", err)
	}
	record.SetOdometerAtInsertionKm(int32(odometerAtInsertion))

	// Parse cardSlotNumber (CardSlotNumber - 1 byte)
	slotNumber := data[102]
	record.SetCardSlotNumber(ddv1.CardSlotNumber(slotNumber))

	// Parse cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime, _, err := readVuTimeRealFromBytes(data, 103)
	if err != nil {
		return nil, fmt.Errorf("failed to read card withdrawal time: %w", err)
	}
	record.SetCardWithdrawalTime(timestamppb.New(time.Unix(withdrawalTime, 0)))

	// Parse vehicleOdometerValueAtWithdrawal (OdometerShort - 3 bytes)
	odometerAtWithdrawal, _, err := readOdometerFromBytes(data, 107)
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer at withdrawal: %w", err)
	}
	record.SetOdometerAtWithdrawalKm(int32(odometerAtWithdrawal))

	// Parse previousVehicleInfo (PreviousVehicleInfo - 15 bytes)
	previousVehicleData := data[110:125]
	previousVehicleInfo, err := parsePreviousVehicleInfo(previousVehicleData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previous vehicle info: %w", err)
	}
	record.SetPreviousVehicleInfo(previousVehicleInfo)

	// Parse manualInputFlag (ManualInputFlag - 1 byte)
	manualInputFlag := data[125]
	record.SetManualInputFlag(manualInputFlag != 0)

	return record, nil
}

// parsePreviousVehicleInfo parses PreviousVehicleInfo structure
func parsePreviousVehicleInfo(data []byte) (*vuv1.Activities_CardIWRecord_PreviousVehicleInfo, error) {
	// PreviousVehicleInfo ::= SEQUENCE {
	//     vehicleRegistrationIdentification VehicleRegistrationIdentification -- 15 bytes
	// }

	if len(data) < 15 {
		return nil, fmt.Errorf("insufficient data for previous vehicle info: got %d, need 15", len(data))
	}

	// Parse VehicleRegistrationIdentification (15 bytes: 1 byte nation + 14 bytes number)
	nation := ddv1.NationNumeric(data[0])
	regNumber, err := unmarshalIA5StringValue(data[1:15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration number: %w", err)
	}

	vehicleReg := &ddv1.VehicleRegistrationIdentification{}
	vehicleReg.SetNation(nation)
	vehicleReg.SetNumber(regNumber)

	info := &vuv1.Activities_CardIWRecord_PreviousVehicleInfo{}
	info.SetVehicleRegistration(vehicleReg)

	return info, nil
}

func parseVuActivityDailyData(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	// VuActivityDailyData ::= SEQUENCE {
	//     noOfActivityChanges INTEGER(0..255),
	//     activityChanges SET SIZE(noOfActivityChanges) OF ActivityChangeInfo
	// }

	// Read number of activity changes (1 byte)
	noOfChanges, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of activity changes: %w", err)
	}

	var changes []*ddv1.ActivityChangeInfo

	// Parse each ActivityChangeInfo
	for i := 0; i < int(noOfChanges); i++ {
		change, newOffset, err := parseActivityChangeInfo(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse activity change %d: %w", i, err)
		}
		changes = append(changes, change)
		offset = newOffset
	}

	return changes, offset, nil
}

// unmarshalActivityChangeInfo parses a single ActivityChangeInfo record from a byte slice
func unmarshalActivityChangeInfo(data []byte) (*ddv1.ActivityChangeInfo, error) {
	// ActivityChangeInfo ::= OCTET STRING (SIZE (2))
	// Bit-packed format: 'scpaattttttttttt'B (16 bits)
	// s: Slot (0=DRIVER, 1=CO-DRIVER)
	// c: Driving status (0=SINGLE, 1=CREW)
	// p: Card status (0=INSERTED, 1=NOT_INSERTED)
	// aa: Activity (00=BREAK/REST, 01=AVAILABILITY, 10=WORK, 11=DRIVING)
	// ttttttttttt: Time (11 bits for time in minutes)

	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for activity change info: got %d, need 2", len(data))
	}

	// Read 2 bytes
	value := uint16(data[0])<<8 | uint16(data[1])

	// Extract bit fields
	slot := (value >> 15) & 0x1          // bit 15
	drivingStatus := (value >> 14) & 0x1 // bit 14
	cardStatus := (value >> 13) & 0x1    // bit 13
	activity := (value >> 11) & 0x3      // bits 12-11
	timeMinutes := value & 0x7FF         // bits 10-0

	// Create ActivityChangeInfo
	change := &ddv1.ActivityChangeInfo{}

	// Set slot
	if slot == 0 {
		change.SetSlot(ddv1.CardSlotNumber_DRIVER_SLOT)
	} else {
		change.SetSlot(ddv1.CardSlotNumber_CO_DRIVER_SLOT)
	}

	// Set driving status
	if drivingStatus == 0 {
		change.SetDrivingStatus(ddv1.DrivingStatus_SINGLE)
	} else {
		change.SetDrivingStatus(ddv1.DrivingStatus_CREW)
	}

	// Set card status (note: bit p=0 means INSERTED, p=1 means NOT_INSERTED)
	change.SetInserted(cardStatus == 0)

	// Set activity
	switch activity {
	case 0:
		change.SetActivity(ddv1.DriverActivityValue_BREAK_REST)
	case 1:
		change.SetActivity(ddv1.DriverActivityValue_AVAILABILITY)
	case 2:
		change.SetActivity(ddv1.DriverActivityValue_WORK)
	case 3:
		change.SetActivity(ddv1.DriverActivityValue_DRIVING)
	}

	// Set time (in minutes)
	change.SetTimeOfChangeMinutes(int32(timeMinutes))

	return change, nil
}

// parseActivityChangeInfo parses a single ActivityChangeInfo record (legacy function for Gen1)
func parseActivityChangeInfo(data []byte, offset int) (*ddv1.ActivityChangeInfo, int, error) {
	change, err := unmarshalActivityChangeInfo(data[offset : offset+2])
	if err != nil {
		return nil, offset, err
	}
	return change, offset + 2, nil
}

func parseVuPlaceDailyWorkPeriodData(data []byte, offset int) ([]*vuv1.Activities_PlaceRecord, int, error) {
	// VuPlaceDailyWorkPeriodData ::= SEQUENCE {
	//     noOfPlaceRecords INTEGER(0..255),
	//     placeRecords SET SIZE(noOfPlaceRecords) OF PlaceRecord
	// }

	// Read number of place records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of place records: %w", err)
	}

	var records []*vuv1.Activities_PlaceRecord

	// Parse each PlaceRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuPlaceRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse place record %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// unmarshalVuPlaceRecord parses a single VuPlaceRecord from a byte slice
func unmarshalVuPlaceRecord(data []byte) (*vuv1.Activities_PlaceRecord, error) {
	// PlaceRecord ::= SEQUENCE {
	//     entryTime TimeReal,                           -- 4 bytes
	//     entryTypeDailyWorkPeriod EntryTypeDailyWorkPeriod, -- 1 byte
	//     dailyWorkPeriodCountry NationNumeric,         -- 1 byte
	//     dailyWorkPeriodRegion RegionNumeric,          -- 1 byte
	//     vehicleOdometerValue OdometerShort            -- 3 bytes
	// }

	if len(data) < 10 {
		return nil, fmt.Errorf("insufficient data for place record: got %d, need 10", len(data))
	}

	record := &vuv1.Activities_PlaceRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, _, err := readVuTimeRealFromBytes(data, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read entry time: %w", err)
	}
	record.SetEntryTime(timestamppb.New(time.Unix(entryTime, 0)))

	// Parse entryTypeDailyWorkPeriod (1 byte)
	entryType := data[4]
	record.SetEntryType(ddv1.EntryTypeDailyWorkPeriod(entryType))

	// Parse dailyWorkPeriodCountry (NationNumeric - 1 byte)
	country := data[5]
	record.SetCountry(ddv1.NationNumeric(country))

	// Parse dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	region := data[6]
	record.SetRegion([]byte{region})

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, _, err := readOdometerFromBytes(data, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, nil
}

// parseVuPlaceRecord parses a single PlaceRecord (legacy function for Gen1)
func parseVuPlaceRecord(data []byte, offset int) (*vuv1.Activities_PlaceRecord, int, error) {
	record, err := unmarshalVuPlaceRecord(data[offset : offset+10])
	if err != nil {
		return nil, offset, err
	}
	return record, offset + 10, nil
}

func parseVuSpecificConditionData(data []byte, offset int) ([]*vuv1.Activities_SpecificConditionRecord, int, error) {
	// VuSpecificConditionData ::= SEQUENCE {
	//     noOfSpecificConditionRecords INTEGER(0..255),
	//     specificConditionRecords SET SIZE(noOfSpecificConditionRecords) OF SpecificConditionRecord
	// }

	// Read number of specific condition records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of specific condition records: %w", err)
	}

	var records []*vuv1.Activities_SpecificConditionRecord

	// Parse each SpecificConditionRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuSpecificConditionRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse specific condition record %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// unmarshalSpecificConditionRecord parses a single SpecificConditionRecord from a byte slice
func unmarshalSpecificConditionRecord(data []byte) (*vuv1.Activities_SpecificConditionRecord, error) {
	// SpecificConditionRecord ::= SEQUENCE {
	//     entryTime TimeReal,                    -- 4 bytes
	//     specificConditionType SpecificConditionType -- 1 byte
	// }

	if len(data) < 5 {
		return nil, fmt.Errorf("insufficient data for specific condition record: got %d, need 5", len(data))
	}

	record := &vuv1.Activities_SpecificConditionRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, _, err := readVuTimeRealFromBytes(data, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read entry time: %w", err)
	}
	record.SetEntryTime(timestamppb.New(time.Unix(entryTime, 0)))

	// Parse specificConditionType (1 byte)
	conditionType := data[4]
	record.SetSpecificConditionType(ddv1.SpecificConditionType(conditionType))

	return record, nil
}

// parseVuSpecificConditionRecord parses a single SpecificConditionRecord (legacy function for Gen1)
func parseVuSpecificConditionRecord(data []byte, offset int) (*vuv1.Activities_SpecificConditionRecord, int, error) {
	record, err := unmarshalSpecificConditionRecord(data[offset : offset+5])
	if err != nil {
		return nil, offset, err
	}
	return record, offset + 5, nil
}

// Gen2 record array parsers
func parseDateOfDayDownloadedRecordArray(data []byte, offset int) ([]*timestamppb.Timestamp, int, error) {
	// DateOfDayDownloadedRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF TimeReal -- 4 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x01 for TimeReal)
	if recordType != 0x01 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 1 (TimeReal)", recordType)
	}

	// Validate record size (should be 4 bytes for TimeReal)
	if recordSize != 4 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 4", recordSize)
	}

	var timestamps []*timestamppb.Timestamp

	// Parse each TimeReal record
	for i := 0; i < int(noOfRecords); i++ {
		if offset+4 > len(data) {
			return nil, offset, fmt.Errorf("insufficient data for TimeReal record %d: got %d, need 4", i, len(data)-offset)
		}

		timeValue, newOffset, err := readVuTimeRealFromBytes(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to read TimeReal record %d: %w", i, err)
		}

		timestamps = append(timestamps, timestamppb.New(time.Unix(timeValue, 0)))
		offset = newOffset
	}

	return timestamps, offset, nil
}

func parseOdometerValueMidnightRecordArray(data []byte, offset int) ([]int32, int, error) {
	// OdometerValueMidnightRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF OdometerShort -- 3 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x02 for OdometerShort)
	if recordType != 0x02 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 2 (OdometerShort)", recordType)
	}

	// Validate record size (should be 3 bytes for OdometerShort)
	if recordSize != 3 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 3", recordSize)
	}

	var odometerValues []int32

	// Parse each OdometerShort record
	for i := 0; i < int(noOfRecords); i++ {
		odometerValue, newOffset, err := readOdometerFromBytes(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to read OdometerShort record %d: %w", i, err)
		}

		odometerValues = append(odometerValues, int32(odometerValue))
		offset = newOffset
	}

	return odometerValues, offset, nil
}

func parseVuCardIWRecordArray(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuCardIWRecord -- Variable size each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes) - not used for variable-length records
	_ = binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x0D for VuCardIWRecord)
	if recordType != 0x0D {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 13 (VuCardIWRecord)", recordType)
	}

	var records []*vuv1.Activities_CardIWRecord

	// Parse each VuCardIWRecord
	for i := 0; i < int(noOfRecords); i++ {
		// For Gen2, we need to determine the record size dynamically
		// VuCardIWRecord has variable length due to HolderName and PreviousVehicleInfo
		// We'll use a conservative approach and try to parse with known minimum size
		record, newOffset, err := parseVuCardIWRecordGen2(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuCardIWRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuCardIWRecordGen2 parses a single VuCardIWRecord for Gen2
func parseVuCardIWRecordGen2(data []byte, offset int) (*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWRecord (Gen2) ::= SEQUENCE {
	//     cardHolderName HolderName,                    -- Variable length (72 bytes)
	//     fullCardNumberAndGeneration FullCardNumberAndGeneration, -- 20 bytes
	//     cardExpiryDate Datef,                         -- 4 bytes
	//     cardInsertionTime TimeReal,                   -- 4 bytes
	//     vehicleOdometerValueAtInsertion OdometerShort, -- 3 bytes
	//     cardSlotNumber CardSlotNumber,                -- 1 byte
	//     cardWithdrawalTime TimeReal,                  -- 4 bytes
	//     vehicleOdometerValueAtWithdrawal OdometerShort, -- 3 bytes
	//     previousVehicleInfo PreviousVehicleInfo,      -- Variable length (15 bytes)
	//     manualInputFlag ManualInputFlag               -- 1 byte
	// }

	record := &vuv1.Activities_CardIWRecord{}

	// Parse cardHolderName (HolderName - 72 bytes)
	holderNameData, offset, err := readBytesFromBytes(data, offset, 72)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card holder name: %w", err)
	}
	holderName, err := unmarshalHolderName(holderNameData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// Parse fullCardNumberAndGeneration (FullCardNumberAndGeneration - 20 bytes)
	fullCardNumberAndGenData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read full card number and generation: %w", err)
	}
	_, err = unmarshalFullCardNumberAndGeneration(fullCardNumberAndGenData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal full card number and generation: %w", err)
	}
	// Note: The protobuf might not have this field yet, so we'll set the regular fullCardNumber
	// This is a limitation of the current schema that should be addressed
	// Create placeholder FullCardNumberAndGeneration
	placeholder := &ddv1.FullCardNumberAndGeneration{}
	placeholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	placeholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetFullCardNumberAndGeneration(placeholder)

	// Parse cardExpiryDate (Datef - 4 bytes)
	datefData, offset, err := readBytesFromBytes(data, offset, 4)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card expiry date: %w", err)
	}
	cardExpiryDate, err := readDatef(bytes.NewReader(datefData))
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse card expiry date: %w", err)
	}
	record.SetCardExpiryDate(cardExpiryDate)

	// Parse cardInsertionTime (TimeReal - 4 bytes)
	insertionTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card insertion time: %w", err)
	}
	record.SetCardInsertionTime(timestamppb.New(time.Unix(insertionTime, 0)))

	// Parse vehicleOdometerValueAtInsertion (OdometerShort - 3 bytes)
	odometerAtInsertion, offset, err := readOdometerFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer at insertion: %w", err)
	}
	record.SetOdometerAtInsertionKm(int32(odometerAtInsertion))

	// Parse cardSlotNumber (CardSlotNumber - 1 byte)
	slotNumber, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card slot number: %w", err)
	}
	record.SetCardSlotNumber(ddv1.CardSlotNumber(slotNumber))

	// Parse cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card withdrawal time: %w", err)
	}
	record.SetCardWithdrawalTime(timestamppb.New(time.Unix(withdrawalTime, 0)))

	// Parse vehicleOdometerValueAtWithdrawal (OdometerShort - 3 bytes)
	odometerAtWithdrawal, offset, err := readOdometerFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer at withdrawal: %w", err)
	}
	record.SetOdometerAtWithdrawalKm(int32(odometerAtWithdrawal))

	// Parse previousVehicleInfo (PreviousVehicleInfo - 15 bytes)
	previousVehicleData, offset, err := readBytesFromBytes(data, offset, 15)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read previous vehicle info: %w", err)
	}
	previousVehicleInfo, err := parsePreviousVehicleInfo(previousVehicleData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse previous vehicle info: %w", err)
	}
	record.SetPreviousVehicleInfo(previousVehicleInfo)

	// Parse manualInputFlag (ManualInputFlag - 1 byte)
	manualInputFlag, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read manual input flag: %w", err)
	}
	record.SetManualInputFlag(manualInputFlag != 0)

	return record, offset, nil
}

// splitActivityChangeInfoRecord splits data into 2-byte ActivityChangeInfo records
func splitActivityChangeInfoRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const activityChangeInfoRecordSize = 2

	if len(data) < activityChangeInfoRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return activityChangeInfoRecordSize, data[:activityChangeInfoRecordSize], nil
}

func parseVuActivityDailyRecordArray(data []byte, offset int) ([]*ddv1.ActivityChangeInfo, int, error) {
	// VuActivityDailyRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF ActivityChangeInfo -- 2 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x03 for ActivityChangeInfo)
	if recordType != 0x03 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 3 (ActivityChangeInfo)", recordType)
	}

	// Validate record size (should be 2 bytes for ActivityChangeInfo)
	if recordSize != 2 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 2", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitActivityChangeInfoRecord)

	var changes []*ddv1.ActivityChangeInfo
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		change, err := unmarshalActivityChangeInfo(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse ActivityChangeInfo record %d: %w", recordCount, err)
		}
		changes = append(changes, change)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 2

	return changes, offset, nil
}

// splitVuPlaceRecord splits data into 10-byte VuPlaceRecord records
func splitVuPlaceRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const placeRecordSize = 10

	if len(data) < placeRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return placeRecordSize, data[:placeRecordSize], nil
}

func parseVuPlaceDailyWorkPeriodRecordArray(data []byte, offset int) ([]*vuv1.Activities_PlaceRecord, int, error) {
	// VuPlaceDailyWorkPeriodRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF PlaceRecord -- 10 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x04 for PlaceRecord)
	if recordType != 0x04 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 4 (PlaceRecord)", recordType)
	}

	// Validate record size (should be 10 bytes for PlaceRecord)
	if recordSize != 10 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 10", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuPlaceRecord)

	var records []*vuv1.Activities_PlaceRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuPlaceRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse PlaceRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 10

	return records, offset, nil
}

// splitVuGNSSADRecord splits data into 14-byte VuGNSSADRecord records
func splitVuGNSSADRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const gnssRecordSize = 14

	if len(data) < gnssRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return gnssRecordSize, data[:gnssRecordSize], nil
}

func parseVuGNSSADRecordArray(data []byte, offset int) ([]*vuv1.Activities_GnssRecord, int, error) {
	// VuGNSSADRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuGNSSADRecord -- 14 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for GNSS record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x16 for VuGNSSADRecord)
	if recordType != 0x16 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 22 (VuGNSSADRecord)", recordType)
	}

	// Validate record size (should be 14 bytes for VuGNSSADRecord)
	if recordSize != 14 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 14", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuGNSSADRecord)

	var records []*vuv1.Activities_GnssRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuGNSSADRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuGNSSADRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 14

	return records, offset, nil
}

// unmarshalVuGNSSADRecord parses a single VuGNSSADRecord from a byte slice
func unmarshalVuGNSSADRecord(data []byte) (*vuv1.Activities_GnssRecord, error) {
	// VuGNSSADRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                    -- 4 bytes
	//     gnssAccuracy GNSSAccuracy,            -- 1 byte
	//     geoCoordinates GeoCoordinates,        -- 8 bytes (latitude + longitude)
	//     positionAuthenticationStatus PositionAuthenticationStatus -- 1 byte
	// }

	if len(data) < 14 {
		return nil, fmt.Errorf("insufficient data for GNSS record: got %d, need 14", len(data))
	}

	record := &vuv1.Activities_GnssRecord{}

	// Parse timeStamp (TimeReal - 4 bytes)
	timestamp, _, err := readVuTimeRealFromBytes(data, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read timestamp: %w", err)
	}
	record.SetTimestamp(timestamppb.New(time.Unix(timestamp, 0)))

	// Parse gnssAccuracy (GNSSAccuracy - 1 byte)
	accuracy := data[4]
	record.SetGnssAccuracy(int32(accuracy))

	// Parse geoCoordinates (GeoCoordinates - 8 bytes: 4 bytes latitude + 4 bytes longitude)
	// Latitude (4 bytes, signed integer)
	latBytes := data[5:9]
	latitude := int32(binary.BigEndian.Uint32(latBytes))

	// Longitude (4 bytes, signed integer)
	lonBytes := data[9:13]
	longitude := int32(binary.BigEndian.Uint32(lonBytes))

	// Create GeoCoordinates
	geoCoords := &ddv1.GeoCoordinates{}
	geoCoords.SetLatitude(latitude)
	geoCoords.SetLongitude(longitude)
	record.SetGeoCoordinates(geoCoords)

	// Parse positionAuthenticationStatus (PositionAuthenticationStatus - 1 byte)
	authStatus := data[13]
	record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus(authStatus))

	return record, nil
}

// splitSpecificConditionRecord splits data into 5-byte SpecificConditionRecord records
func splitSpecificConditionRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const specificConditionRecordSize = 5

	if len(data) < specificConditionRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return specificConditionRecordSize, data[:specificConditionRecordSize], nil
}

func parseVuSpecificConditionRecordArray(data []byte, offset int) ([]*vuv1.Activities_SpecificConditionRecord, int, error) {
	// VuSpecificConditionRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF SpecificConditionRecord -- 5 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x09 for SpecificConditionRecord)
	if recordType != 0x09 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 9 (SpecificConditionRecord)", recordType)
	}

	// Validate record size (should be 5 bytes for SpecificConditionRecord)
	if recordSize != 5 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 5", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitSpecificConditionRecord)

	var records []*vuv1.Activities_SpecificConditionRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalSpecificConditionRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse SpecificConditionRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 5

	return records, offset, nil
}

// splitVuBorderCrossingRecord splits data into 59-byte VuBorderCrossingRecord records
func splitVuBorderCrossingRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const borderCrossingRecordSize = 59

	if len(data) < borderCrossingRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return borderCrossingRecordSize, data[:borderCrossingRecordSize], nil
}

func parseVuBorderCrossingRecordArray(data []byte, offset int) ([]*vuv1.Activities_BorderCrossingRecord, int, error) {
	// VuBorderCrossingRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuBorderCrossingRecord -- 59 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for border crossing record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x17 for VuBorderCrossingRecord)
	if recordType != 0x17 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 23 (VuBorderCrossingRecord)", recordType)
	}

	// Validate record size (should be 59 bytes for VuBorderCrossingRecord)
	if recordSize != 59 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 59", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuBorderCrossingRecord)

	var records []*vuv1.Activities_BorderCrossingRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuBorderCrossingRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuBorderCrossingRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 59

	return records, offset, nil
}

// unmarshalVuBorderCrossingRecord parses a single VuBorderCrossingRecord from a byte slice
func unmarshalVuBorderCrossingRecord(data []byte) (*vuv1.Activities_BorderCrossingRecord, error) {
	// VuBorderCrossingRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     countryLeft NationNumeric,                              -- 1 byte
	//     countryEntered NationNumeric,                           -- 1 byte
	//     placeRecord GNSSPlaceRecord,                            -- 14 bytes (timestamp + coords + auth)
	//     odometerValue OdometerShort                             -- 3 bytes
	// }

	if len(data) < 59 {
		return nil, fmt.Errorf("insufficient data for border crossing record: got %d, need 59", len(data))
	}

	record := &vuv1.Activities_BorderCrossingRecord{}

	// Parse cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	driverCardData := data[0:20]
	_, err := unmarshalFullCardNumberAndGeneration(driverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal driver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	driverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	driverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	driverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberDriverSlot(driverPlaceholder)

	// Parse cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	codriverCardData := data[20:40]
	_, err = unmarshalFullCardNumberAndGeneration(codriverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal codriver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	codriverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	codriverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	codriverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberCodriverSlot(codriverPlaceholder)

	// Parse countryLeft (NationNumeric - 1 byte)
	countryLeft := data[40]
	record.SetCountryLeft(ddv1.NationNumeric(countryLeft))

	// Parse countryEntered (NationNumeric - 1 byte)
	countryEntered := data[41]
	record.SetCountryEntered(ddv1.NationNumeric(countryEntered))

	// Parse placeRecord (GNSSPlaceRecord - 14 bytes)
	placeData := data[42:56]
	placeRecord, err := unmarshalGNSSPlaceRecord(placeData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal place record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	// Parse odometerValue (OdometerShort - 3 bytes)
	odometerValue, _, err := readOdometerFromBytes(data, 56)
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, nil
}

// unmarshalGNSSPlaceRecord parses a single GNSSPlaceRecord from a byte slice
func unmarshalGNSSPlaceRecord(data []byte) (*vuv1.Activities_GnssRecord, error) {
	// GNSSPlaceRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                    -- 4 bytes
	//     gnssAccuracy GNSSAccuracy,            -- 1 byte
	//     geoCoordinates GeoCoordinates,        -- 8 bytes
	//     positionAuthenticationStatus PositionAuthenticationStatus -- 1 byte
	// }

	if len(data) < 14 {
		return nil, fmt.Errorf("insufficient data for GNSS place record: got %d, need 14", len(data))
	}

	record := &vuv1.Activities_GnssRecord{}

	// Parse timeStamp (TimeReal - 4 bytes)
	timestamp, _, err := readVuTimeRealFromBytes(data, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read timestamp: %w", err)
	}
	record.SetTimestamp(timestamppb.New(time.Unix(timestamp, 0)))

	// Parse gnssAccuracy (GNSSAccuracy - 1 byte)
	accuracy := data[4]
	record.SetGnssAccuracy(int32(accuracy))

	// Parse geoCoordinates (GeoCoordinates - 8 bytes: 4 bytes latitude + 4 bytes longitude)
	// Latitude (4 bytes, signed integer)
	latBytes := data[5:9]
	latitude := int32(binary.BigEndian.Uint32(latBytes))

	// Longitude (4 bytes, signed integer)
	lonBytes := data[9:13]
	longitude := int32(binary.BigEndian.Uint32(lonBytes))

	// Create GeoCoordinates
	geoCoords := &ddv1.GeoCoordinates{}
	geoCoords.SetLatitude(latitude)
	geoCoords.SetLongitude(longitude)
	record.SetGeoCoordinates(geoCoords)

	// Parse positionAuthenticationStatus (PositionAuthenticationStatus - 1 byte)
	authStatus := data[13]
	record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus(authStatus))

	return record, nil
}

// splitVuLoadUnloadRecord splits data into 58-byte VuLoadUnloadRecord records
func splitVuLoadUnloadRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const loadUnloadRecordSize = 58

	if len(data) < loadUnloadRecordSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return loadUnloadRecordSize, data[:loadUnloadRecordSize], nil
}

func parseVuLoadUnloadRecordArray(data []byte, offset int) ([]*vuv1.Activities_LoadUnloadRecord, int, error) {
	// VuLoadUnloadRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuLoadUnloadRecord -- 58 bytes each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for load/unload record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes)
	noOfRecords := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x18 for VuLoadUnloadRecord)
	if recordType != 0x18 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 24 (VuLoadUnloadRecord)", recordType)
	}

	// Validate record size (should be 58 bytes for VuLoadUnloadRecord)
	if recordSize != 58 {
		return nil, offset, fmt.Errorf("unexpected record size: got %d, expected 58", recordSize)
	}

	// Use bufio.Scanner to parse the records
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitVuLoadUnloadRecord)

	var records []*vuv1.Activities_LoadUnloadRecord
	recordCount := 0

	for scanner.Scan() {
		if recordCount >= int(noOfRecords) {
			break
		}

		recordData := scanner.Bytes()
		record, err := unmarshalVuLoadUnloadRecord(recordData)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuLoadUnloadRecord %d: %w", recordCount, err)
		}
		records = append(records, record)
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, offset, fmt.Errorf("scanner error: %w", err)
	}

	// Update offset to reflect consumed data
	offset += recordCount * 58

	return records, offset, nil
}

// unmarshalVuLoadUnloadRecord parses a single VuLoadUnloadRecord from a byte slice
func unmarshalVuLoadUnloadRecord(data []byte) (*vuv1.Activities_LoadUnloadRecord, error) {
	// VuLoadUnloadRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     operationType OperationType,                            -- 1 byte
	//     placeRecord GNSSPlaceRecord,                            -- 14 bytes
	//     odometerValue OdometerShort                             -- 3 bytes
	// }

	if len(data) < 58 {
		return nil, fmt.Errorf("insufficient data for load/unload record: got %d, need 58", len(data))
	}

	record := &vuv1.Activities_LoadUnloadRecord{}

	// Parse cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	driverCardData := data[0:20]
	_, err := unmarshalFullCardNumberAndGeneration(driverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal driver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	driverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	driverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	driverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberDriverSlot(driverPlaceholder)

	// Parse cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	codriverCardData := data[20:40]
	_, err = unmarshalFullCardNumberAndGeneration(codriverCardData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal codriver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	codriverPlaceholder := &ddv1.FullCardNumberAndGeneration{}
	codriverPlaceholder.SetFullCardNumber(&ddv1.FullCardNumber{})
	codriverPlaceholder.SetGeneration(ddv1.Generation_GENERATION_1)
	record.SetCardNumberCodriverSlot(codriverPlaceholder)

	// Parse operationType (OperationType - 1 byte)
	operationType := data[40]
	record.SetOperationType(ddv1.OperationType(operationType))

	// Parse placeRecord (GNSSPlaceRecord - 14 bytes)
	placeData := data[41:55]
	placeRecord, err := unmarshalGNSSPlaceRecord(placeData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal place record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	// Parse odometerValue (OdometerShort - 3 bytes)
	odometerValue, _, err := readOdometerFromBytes(data, 55)
	if err != nil {
		return nil, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, nil
}

// splitSignatureRecord splits data into variable-length Signature records
func splitSignatureRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// We need at least 6 bytes for the record array header
	if len(data) < 6 {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	// Read recordSize from the header (bytes 2-3)
	recordSize := binary.BigEndian.Uint16(data[2:4])

	// Validate record size
	if recordSize == 0 {
		return 0, nil, fmt.Errorf("invalid record size: got 0")
	}

	// Total record size = header (6 bytes) + signature data (recordSize bytes)
	totalSize := 6 + int(recordSize)

	if len(data) < totalSize {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return totalSize, data[:totalSize], nil
}

func parseSignatureRecordArray(data []byte, offset int) ([]byte, int, error) {
	// SignatureRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF Signature -- Variable size each
	// }

	// Use bufio.Scanner to parse the signature record
	recordsData := data[offset:]
	scanner := bufio.NewScanner(bytes.NewReader(recordsData))
	scanner.Split(splitSignatureRecord)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, offset, fmt.Errorf("scanner error: %w", err)
		}
		return nil, offset, fmt.Errorf("no signature record found")
	}

	recordData := scanner.Bytes()

	// Validate record type (should be 0x08 for Signature)
	recordType := binary.BigEndian.Uint16(recordData[0:2])
	if recordType != 0x08 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 8 (Signature)", recordType)
	}

	// Extract signature data (skip 6-byte header)
	signature := make([]byte, len(recordData)-6)
	copy(signature, recordData[6:])

	// Update offset to reflect consumed data
	offset += len(recordData)

	return signature, offset, nil
}

// AppendVuActivities appends VU activities data to a buffer.
//
// The data type `VuActivities` is specified in the Data Dictionary, Section 2.2.6.2.
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    dateOfDay                        TimeReal,
//	    odometerValueMidnight            OdometerValueMidnight,
//	    vuCardIWData                     VuCardIWData,
//	    vuActivityDailyData              VuActivityDailyData,
//	    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
//	    vuSpecificConditionData          VuSpecificConditionData,
//	    signature                        SignatureFirstGen
//	}
//
//	VuActivitiesSecondGen ::= SEQUENCE {
//	    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
//	    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
//	    vuCardIWRecordArray                      VuCardIWRecordArray,
//	    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
//	    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
//	    vuGNSSADRecordArray                      VuGNSSADRecordArray,
//	    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
//	    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
//	    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
//	    signatureRecordArray                     SignatureRecordArray
//	}

// appendVuActivitiesBytes appends VU activities data to a byte slice
func appendVuActivitiesBytes(dst []byte, activities *vuv1.Activities) ([]byte, error) {
	if activities == nil {
		return dst, nil
	}

	if activities.GetGeneration() == ddv1.Generation_GENERATION_1 {
		return appendVuActivitiesGen1Bytes(dst, activities)
	} else {
		return appendVuActivitiesGen2Bytes(dst, activities)
	}
}

// appendVuActivitiesGen1Bytes appends Generation 1 VU activities data
func appendVuActivitiesGen1Bytes(dst []byte, activities *vuv1.Activities) ([]byte, error) {
	// DateOfDay (TimeReal - 4 bytes)
	dst = appendVuTimeReal(dst, activities.GetDateOfDay())

	// OdometerValueMidnight (3 bytes)
	dst = appendVuOdometer(dst, activities.GetOdometerMidnightKm())

	// VuCardIWData
	var err error
	dst, err = appendVuCardIWData(dst, activities.GetCardIwData())
	if err != nil {
		return nil, fmt.Errorf("failed to append card IW data: %w", err)
	}

	// VuActivityDailyData
	dst, err = appendVuActivityDailyData(dst, activities.GetActivityChanges())
	if err != nil {
		return nil, fmt.Errorf("failed to append activity daily data: %w", err)
	}

	// VuPlaceDailyWorkPeriodData
	dst, err = appendVuPlaceDailyWorkPeriodData(dst, activities.GetPlaces())
	if err != nil {
		return nil, fmt.Errorf("failed to append place daily work period data: %w", err)
	}

	// VuSpecificConditionData
	dst, err = appendVuSpecificConditionData(dst, activities.GetSpecificConditions())
	if err != nil {
		return nil, fmt.Errorf("failed to append specific condition data: %w", err)
	}

	// Signature (128 bytes for Gen1)
	signature := activities.GetSignatureGen1()
	if len(signature) > 0 {
		dst = append(dst, signature...)
	} else {
		// Pad with zeros if no signature
		dst = append(dst, make([]byte, 128)...)
	}

	return dst, nil
}

// appendVuActivitiesGen2Bytes appends Generation 2 VU activities data
func appendVuActivitiesGen2Bytes(dst []byte, activities *vuv1.Activities) ([]byte, error) {
	// Gen2 format uses record arrays, each with a header
	// Append DateOfDayDownloadedRecordArray
	dates := []*timestamppb.Timestamp{activities.GetDateOfDay()}
	dst, err := appendDateOfDayDownloadedRecordArray(dst, dates)
	if err != nil {
		return nil, fmt.Errorf("failed to append date of day downloaded record array: %w", err)
	}

	// Append OdometerValueMidnightRecordArray
	odometerValues := []int32{activities.GetOdometerMidnightKm()}
	dst, err = appendOdometerValueMidnightRecordArray(dst, odometerValues)
	if err != nil {
		return nil, fmt.Errorf("failed to append odometer value midnight record array: %w", err)
	}

	// Append VuCardIWRecordArray
	dst, err = appendVuCardIWRecordArray(dst, activities.GetCardIwData())
	if err != nil {
		return nil, fmt.Errorf("failed to append card IW record array: %w", err)
	}

	// Append VuActivityDailyRecordArray
	dst, err = appendVuActivityDailyRecordArray(dst, activities.GetActivityChanges())
	if err != nil {
		return nil, fmt.Errorf("failed to append activity daily record array: %w", err)
	}

	// Append VuPlaceDailyWorkPeriodRecordArray
	dst, err = appendVuPlaceDailyWorkPeriodRecordArray(dst, activities.GetPlaces())
	if err != nil {
		return nil, fmt.Errorf("failed to append place daily work period record array: %w", err)
	}

	// Append VuGNSSADRecordArray (Gen2+)
	dst, err = appendVuGNSSADRecordArray(dst, activities.GetGnssAccumulatedDriving())
	if err != nil {
		return nil, fmt.Errorf("failed to append GNSS AD record array: %w", err)
	}

	// Append VuSpecificConditionRecordArray
	dst, err = appendVuSpecificConditionRecordArray(dst, activities.GetSpecificConditions())
	if err != nil {
		return nil, fmt.Errorf("failed to append specific condition record array: %w", err)
	}

	// Append Gen2v2 specific arrays if present
	if activities.GetVersion() == vuv1.Version_VERSION_2 {
		// Append VuBorderCrossingRecordArray (Gen2v2+)
		dst, err = appendVuBorderCrossingRecordArray(dst, activities.GetBorderCrossings())
		if err != nil {
			return nil, fmt.Errorf("failed to append border crossing record array: %w", err)
		}

		// Append VuLoadUnloadRecordArray (Gen2v2+)
		dst, err = appendVuLoadUnloadRecordArray(dst, activities.GetLoadUnloadOperations())
		if err != nil {
			return nil, fmt.Errorf("failed to append load/unload record array: %w", err)
		}
	}

	// Append SignatureRecordArray
	dst, err = appendSignatureRecordArray(dst, activities.GetSignatureGen2())
	if err != nil {
		return nil, fmt.Errorf("failed to append signature record array: %w", err)
	}

	return dst, nil
}

// Helper functions for appending VU data types

// appendVuOdometer appends an odometer value (3 bytes) to dst
func appendVuOdometer(dst []byte, value int32) []byte {
	// Convert to 3-byte big-endian
	odometerBytes := make([]byte, 3)
	odometerBytes[0] = byte((value >> 16) & 0xFF)
	odometerBytes[1] = byte((value >> 8) & 0xFF)
	odometerBytes[2] = byte(value & 0xFF)
	return append(dst, odometerBytes...)
}

// appendVuCardIWData appends VuCardIWData to dst
func appendVuCardIWData(dst []byte, cardIWData []*vuv1.Activities_CardIWRecord) ([]byte, error) {
	if cardIWData == nil {
		// Write number of records (1 byte) as 0
		return append(dst, 0), nil
	}

	// Write number of records (1 byte)
	if len(cardIWData) > 255 {
		return nil, fmt.Errorf("too many card IW records: %d", len(cardIWData))
	}
	dst = append(dst, byte(len(cardIWData)))

	// Write each record
	for _, record := range cardIWData {
		var err error
		dst, err = appendVuCardIWRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append card IW record: %w", err)
		}
	}

	return dst, nil
}

// appendVuCardIWRecord appends a single VuCardIWRecord to dst
func appendVuCardIWRecord(dst []byte, record *vuv1.Activities_CardIWRecord) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// VuCardIWRecord ::= SEQUENCE {
	//     cardHolderName HolderName,                    -- 72 bytes
	//     fullCardNumber FullCardNumber,                -- 19 bytes
	//     cardExpiryDate Datef,                         -- 4 bytes
	//     cardInsertionTime TimeReal,                   -- 4 bytes
	//     vehicleOdometerValueAtInsertion OdometerShort, -- 3 bytes
	//     cardSlotNumber CardSlotNumber,                -- 1 byte
	//     cardWithdrawalTime TimeReal,                  -- 4 bytes
	//     vehicleOdometerValueAtWithdrawal OdometerShort, -- 3 bytes
	//     previousVehicleInfo PreviousVehicleInfo,      -- 15 bytes
	//     manualInputFlag ManualInputFlag               -- 1 byte
	// }

	var err error

	// Append cardHolderName (HolderName - 72 bytes)
	holderName := record.GetCardHolderName()
	if holderName != nil {
		dst, err = appendHolderName(dst, holderName)
		if err != nil {
			return nil, fmt.Errorf("failed to append holder name: %w", err)
		}
	} else {
		// Pad with spaces if no holder name
		dst = append(dst, make([]byte, 72)...)
	}

	// Append fullCardNumber (FullCardNumber - 19 bytes)
	fullCardNumberAndGeneration := record.GetFullCardNumberAndGeneration()
	var fullCardNumber *ddv1.FullCardNumber
	if fullCardNumberAndGeneration != nil {
		fullCardNumber = fullCardNumberAndGeneration.GetFullCardNumber()
	}
	if fullCardNumber != nil {
		dst, err = appendFullCardNumber(dst, fullCardNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to append full card number: %w", err)
		}
	} else {
		// Pad with zeros if no card number
		dst = append(dst, make([]byte, 19)...)
	}

	// Append cardExpiryDate (Datef - 4 bytes)
	cardExpiryDate := record.GetCardExpiryDate()
	if cardExpiryDate != nil {
		dst = appendDate(dst, cardExpiryDate)
	} else {
		// Append default date (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	// Append cardInsertionTime (TimeReal - 4 bytes)
	insertionTime := record.GetCardInsertionTime()
	if insertionTime != nil {
		dst = appendVuTimeReal(dst, insertionTime)
	} else {
		// Append default time (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	// Append vehicleOdometerValueAtInsertion (OdometerShort - 3 bytes)
	dst = appendVuOdometer(dst, record.GetOdometerAtInsertionKm())

	// Append cardSlotNumber (CardSlotNumber - 1 byte)
	dst = append(dst, byte(record.GetCardSlotNumber()))

	// Append cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime := record.GetCardWithdrawalTime()
	if withdrawalTime != nil {
		dst = appendVuTimeReal(dst, withdrawalTime)
	} else {
		// Append default time (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	// Append vehicleOdometerValueAtWithdrawal (OdometerShort - 3 bytes)
	dst = appendVuOdometer(dst, record.GetOdometerAtWithdrawalKm())

	// Append previousVehicleInfo (PreviousVehicleInfo - 15 bytes)
	previousVehicleInfo := record.GetPreviousVehicleInfo()
	if previousVehicleInfo != nil {
		vehicleReg := previousVehicleInfo.GetVehicleRegistration()
		if vehicleReg != nil {
			// Nation (1 byte)
			dst = append(dst, byte(vehicleReg.GetNation()))
			// Registration number (14 bytes)
			regNumber := vehicleReg.GetNumber()
			if regNumber != nil {
				regStr := regNumber.GetDecoded()
				if len(regStr) > 14 {
					regStr = regStr[:14]
				}
				dst = append(dst, []byte(regStr)...)
				// Pad with spaces if needed
				for i := len(regStr); i < 14; i++ {
					dst = append(dst, ' ')
				}
			} else {
				dst = append(dst, make([]byte, 14)...)
			}
		} else {
			dst = append(dst, make([]byte, 15)...)
		}
	} else {
		dst = append(dst, make([]byte, 15)...)
	}

	// Append manualInputFlag (ManualInputFlag - 1 byte)
	if record.GetManualInputFlag() {
		dst = append(dst, 1)
	} else {
		dst = append(dst, 0)
	}

	return dst, nil
}

// appendVuActivityDailyData appends VuActivityDailyData to dst
func appendVuActivityDailyData(dst []byte, activityChanges []*ddv1.ActivityChangeInfo) ([]byte, error) {
	if activityChanges == nil {
		// Write number of activity changes (1 byte) as 0
		return append(dst, 0), nil
	}

	// Write number of activity changes (1 byte)
	if len(activityChanges) > 255 {
		return nil, fmt.Errorf("too many activity changes: %d", len(activityChanges))
	}
	dst = append(dst, byte(len(activityChanges)))

	// Write each activity change (2 bytes each)
	for _, change := range activityChanges {
		var err error
		dst, err = appendVuActivityChangeInfo(dst, change)
		if err != nil {
			return nil, fmt.Errorf("failed to append activity change: %w", err)
		}
	}

	return dst, nil
}

// appendVuActivityChangeInfo appends an ActivityChangeInfo to dst
func appendVuActivityChangeInfo(dst []byte, change *ddv1.ActivityChangeInfo) ([]byte, error) {
	if change == nil {
		return dst, nil
	}

	// ActivityChangeInfo is 2 bytes packed as bitfield
	var aci uint16

	// Extract values and pack into bitfield
	slot := getCardSlotNumberProtocolValue(change.GetSlot(), 0)
	drivingStatus := getDrivingStatusProtocolValue(change.GetDrivingStatus(), 0)
	cardInserted := getCardInsertedFromBool(change.GetInserted())
	activity := getDriverActivityValueProtocolValue(change.GetActivity(), 0)

	aci |= (uint16(slot) & 0x1) << 15
	aci |= (uint16(drivingStatus) & 0x1) << 14
	aci |= (uint16(cardInserted) & 0x1) << 13
	aci |= (uint16(activity) & 0x3) << 11
	aci |= (uint16(change.GetTimeOfChangeMinutes()) & 0x7FF)

	return binary.BigEndian.AppendUint16(dst, aci), nil
}

// appendVuPlaceDailyWorkPeriodData appends VuPlaceDailyWorkPeriodData to dst
func appendVuPlaceDailyWorkPeriodData(dst []byte, places []*vuv1.Activities_PlaceRecord) ([]byte, error) {
	if places == nil {
		// Write number of place records (1 byte) as 0
		return append(dst, 0), nil
	}

	// Write number of place records (1 byte)
	if len(places) > 255 {
		return nil, fmt.Errorf("too many place records: %d", len(places))
	}
	dst = append(dst, byte(len(places)))

	// Write each place record
	for _, place := range places {
		var err error
		dst, err = appendVuPlaceRecord(dst, place)
		if err != nil {
			return nil, fmt.Errorf("failed to append place record: %w", err)
		}
	}

	return dst, nil
}

// appendVuPlaceRecord appends a single PlaceRecord to dst
func appendVuPlaceRecord(dst []byte, place *vuv1.Activities_PlaceRecord) ([]byte, error) {
	if place == nil {
		return dst, nil
	}

	// PlaceRecord ::= SEQUENCE {
	//     entryTime TimeReal,                           -- 4 bytes
	//     entryTypeDailyWorkPeriod EntryTypeDailyWorkPeriod, -- 1 byte
	//     dailyWorkPeriodCountry NationNumeric,         -- 1 byte
	//     dailyWorkPeriodRegion RegionNumeric,          -- 1 byte
	//     vehicleOdometerValue OdometerShort            -- 3 bytes
	// }

	// Append entryTime (TimeReal - 4 bytes)
	entryTime := place.GetEntryTime()
	if entryTime != nil {
		dst = appendVuTimeReal(dst, entryTime)
	} else {
		// Append default time (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	// Append entryTypeDailyWorkPeriod (1 byte)
	dst = append(dst, byte(place.GetEntryType()))

	// Append dailyWorkPeriodCountry (NationNumeric - 1 byte)
	dst = append(dst, byte(place.GetCountry()))

	// Append dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	region := place.GetRegion()
	if len(region) > 0 {
		dst = append(dst, region[0])
	} else {
		dst = append(dst, 0)
	}

	// Append vehicleOdometerValue (OdometerShort - 3 bytes)
	dst = appendVuOdometer(dst, place.GetOdometerKm())

	return dst, nil
}

// appendVuSpecificConditionData appends VuSpecificConditionData to dst
func appendVuSpecificConditionData(dst []byte, specificConditions []*vuv1.Activities_SpecificConditionRecord) ([]byte, error) {
	if specificConditions == nil {
		// Write number of specific condition records (1 byte) as 0
		return append(dst, 0), nil
	}

	// Write number of specific condition records (1 byte)
	if len(specificConditions) > 255 {
		return nil, fmt.Errorf("too many specific condition records: %d", len(specificConditions))
	}
	dst = append(dst, byte(len(specificConditions)))

	// Write each specific condition record
	for _, condition := range specificConditions {
		var err error
		dst, err = appendVuSpecificConditionRecord(dst, condition)
		if err != nil {
			return nil, fmt.Errorf("failed to append specific condition record: %w", err)
		}
	}

	return dst, nil
}

// appendVuSpecificConditionRecord appends a single SpecificConditionRecord to dst
func appendVuSpecificConditionRecord(dst []byte, condition *vuv1.Activities_SpecificConditionRecord) ([]byte, error) {
	if condition == nil {
		return dst, nil
	}

	// SpecificConditionRecord ::= SEQUENCE {
	//     entryTime TimeReal,                    -- 4 bytes
	//     specificConditionType SpecificConditionType -- 1 byte
	// }

	// Append entryTime (TimeReal - 4 bytes)
	entryTime := condition.GetEntryTime()
	if entryTime != nil {
		dst = appendVuTimeReal(dst, entryTime)
	} else {
		// Append default time (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	// Append specificConditionType (1 byte)
	dst = append(dst, byte(condition.GetSpecificConditionType()))

	return dst, nil
}

// Gen2 record array append functions

// appendDateOfDayDownloadedRecordArray appends DateOfDayDownloadedRecordArray to dst
func appendDateOfDayDownloadedRecordArray(dst []byte, dates []*timestamppb.Timestamp) ([]byte, error) {
	// DateOfDayDownloadedRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF TimeReal -- 4 bytes each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x01) // TimeReal
	recordSize := uint16(4)    // 4 bytes for TimeReal
	noOfRecords := uint16(len(dates))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each TimeReal record
	for _, date := range dates {
		if date != nil {
			dst = appendVuTimeReal(dst, date)
		} else {
			// Append default time (00000000)
			dst = append(dst, 0x00, 0x00, 0x00, 0x00)
		}
	}

	return dst, nil
}

// appendOdometerValueMidnightRecordArray appends OdometerValueMidnightRecordArray to dst
func appendOdometerValueMidnightRecordArray(dst []byte, odometerValues []int32) ([]byte, error) {
	// OdometerValueMidnightRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF OdometerShort -- 3 bytes each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x02) // OdometerShort
	recordSize := uint16(3)    // 3 bytes for OdometerShort
	noOfRecords := uint16(len(odometerValues))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each OdometerShort record
	for _, value := range odometerValues {
		dst = appendVuOdometer(dst, value)
	}

	return dst, nil
}

// appendVuCardIWRecordArray appends VuCardIWRecordArray to dst
func appendVuCardIWRecordArray(dst []byte, cardIWData []*vuv1.Activities_CardIWRecord) ([]byte, error) {
	// VuCardIWRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes (not used for variable-length)
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuCardIWRecord -- Variable size each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x0D) // VuCardIWRecord
	recordSize := uint16(0)    // Not used for variable-length records
	noOfRecords := uint16(len(cardIWData))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each VuCardIWRecord
	for _, record := range cardIWData {
		var err error
		dst, err = appendVuCardIWRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append card IW record: %w", err)
		}
	}

	return dst, nil
}

// appendVuActivityDailyRecordArray appends VuActivityDailyRecordArray to dst
func appendVuActivityDailyRecordArray(dst []byte, activityChanges []*ddv1.ActivityChangeInfo) ([]byte, error) {
	// VuActivityDailyRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF ActivityChangeInfo -- 2 bytes each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x03) // ActivityChangeInfo
	recordSize := uint16(2)    // 2 bytes for ActivityChangeInfo
	noOfRecords := uint16(len(activityChanges))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each ActivityChangeInfo record
	for _, change := range activityChanges {
		var err error
		dst, err = appendVuActivityChangeInfo(dst, change)
		if err != nil {
			return nil, fmt.Errorf("failed to append activity change: %w", err)
		}
	}

	return dst, nil
}

// appendVuPlaceDailyWorkPeriodRecordArray appends VuPlaceDailyWorkPeriodRecordArray to dst
func appendVuPlaceDailyWorkPeriodRecordArray(dst []byte, places []*vuv1.Activities_PlaceRecord) ([]byte, error) {
	// VuPlaceDailyWorkPeriodRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF PlaceRecord -- 10 bytes each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x04) // PlaceRecord
	recordSize := uint16(10)   // 10 bytes for PlaceRecord
	noOfRecords := uint16(len(places))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each PlaceRecord
	for _, place := range places {
		var err error
		dst, err = appendVuPlaceRecord(dst, place)
		if err != nil {
			return nil, fmt.Errorf("failed to append place record: %w", err)
		}
	}

	return dst, nil
}

// appendVuGNSSADRecordArray appends VuGNSSADRecordArray to dst
func appendVuGNSSADRecordArray(dst []byte, gnssRecords []*vuv1.Activities_GnssRecord) ([]byte, error) {
	// VuGNSSADRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuGNSSADRecord -- 14 bytes each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x16) // VuGNSSADRecord
	recordSize := uint16(14)   // 14 bytes for VuGNSSADRecord
	noOfRecords := uint16(len(gnssRecords))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each VuGNSSADRecord
	for _, record := range gnssRecords {
		var err error
		dst, err = appendVuGNSSADRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append GNSS record: %w", err)
		}
	}

	return dst, nil
}

// appendVuGNSSADRecord appends a single VuGNSSADRecord to dst
func appendVuGNSSADRecord(dst []byte, record *vuv1.Activities_GnssRecord) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// VuGNSSADRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                    -- 4 bytes
	//     gnssAccuracy GNSSAccuracy,            -- 1 byte
	//     geoCoordinates GeoCoordinates,        -- 8 bytes (latitude + longitude)
	//     positionAuthenticationStatus PositionAuthenticationStatus -- 1 byte
	// }

	// Append timeStamp (TimeReal - 4 bytes)
	timestamp := record.GetTimestamp()
	if timestamp != nil {
		dst = appendVuTimeReal(dst, timestamp)
	} else {
		// Append default time (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00)
	}

	// Append gnssAccuracy (GNSSAccuracy - 1 byte)
	dst = append(dst, byte(record.GetGnssAccuracy()))

	// Append geoCoordinates (GeoCoordinates - 8 bytes: 4 bytes latitude + 4 bytes longitude)
	geoCoords := record.GetGeoCoordinates()
	if geoCoords != nil {
		// Latitude (4 bytes, signed integer)
		latitude := geoCoords.GetLatitude()
		dst = binary.BigEndian.AppendUint32(dst, uint32(latitude))
		// Longitude (4 bytes, signed integer)
		longitude := geoCoords.GetLongitude()
		dst = binary.BigEndian.AppendUint32(dst, uint32(longitude))
	} else {
		// Append default coordinates (00000000)
		dst = append(dst, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	}

	// Append positionAuthenticationStatus (PositionAuthenticationStatus - 1 byte)
	dst = append(dst, byte(record.GetAuthenticationStatus()))

	return dst, nil
}

// appendVuSpecificConditionRecordArray appends VuSpecificConditionRecordArray to dst
func appendVuSpecificConditionRecordArray(dst []byte, specificConditions []*vuv1.Activities_SpecificConditionRecord) ([]byte, error) {
	// VuSpecificConditionRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF SpecificConditionRecord -- 5 bytes each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x09) // SpecificConditionRecord
	recordSize := uint16(5)    // 5 bytes for SpecificConditionRecord
	noOfRecords := uint16(len(specificConditions))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each SpecificConditionRecord
	for _, condition := range specificConditions {
		var err error
		dst, err = appendVuSpecificConditionRecord(dst, condition)
		if err != nil {
			return nil, fmt.Errorf("failed to append specific condition record: %w", err)
		}
	}

	return dst, nil
}

// appendVuBorderCrossingRecordArray appends VuBorderCrossingRecordArray to dst
func appendVuBorderCrossingRecordArray(dst []byte, borderCrossings []*vuv1.Activities_BorderCrossingRecord) ([]byte, error) {
	// VuBorderCrossingRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuBorderCrossingRecord -- Variable size each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x17) // VuBorderCrossingRecord
	recordSize := uint16(0)    // Not used for variable-length records
	noOfRecords := uint16(len(borderCrossings))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each VuBorderCrossingRecord
	for _, record := range borderCrossings {
		var err error
		dst, err = appendVuBorderCrossingRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append border crossing record: %w", err)
		}
	}

	return dst, nil
}

// appendVuBorderCrossingRecord appends a single VuBorderCrossingRecord to dst
func appendVuBorderCrossingRecord(dst []byte, record *vuv1.Activities_BorderCrossingRecord) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// VuBorderCrossingRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     countryLeft NationNumeric,                              -- 1 byte
	//     countryEntered NationNumeric,                           -- 1 byte
	//     placeRecord GNSSPlaceRecord,                            -- 14 bytes (timestamp + coords + auth)
	//     odometerValue OdometerShort                             -- 3 bytes
	// }

	// For now, implement simplified versions due to schema limitations
	// These would need to be completed when the protobuf schema is updated

	// Append cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	// Note: Schema limitation - using placeholder for now
	dst = append(dst, make([]byte, 20)...)

	// Append cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	// Note: Schema limitation - using placeholder for now
	dst = append(dst, make([]byte, 20)...)

	// Append countryLeft (NationNumeric - 1 byte)
	dst = append(dst, byte(record.GetCountryLeft()))

	// Append countryEntered (NationNumeric - 1 byte)
	dst = append(dst, byte(record.GetCountryEntered()))

	// Append placeRecord (GNSSPlaceRecord - 14 bytes)
	placeRecord := record.GetPlaceRecord()
	if placeRecord != nil {
		var err error
		dst, err = appendVuGNSSADRecord(dst, placeRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to append place record: %w", err)
		}
	} else {
		// Append default place record (14 bytes of zeros)
		dst = append(dst, make([]byte, 14)...)
	}

	// Append odometerValue (OdometerShort - 3 bytes)
	dst = appendVuOdometer(dst, record.GetOdometerKm())

	return dst, nil
}

// appendVuLoadUnloadRecordArray appends VuLoadUnloadRecordArray to dst
func appendVuLoadUnloadRecordArray(dst []byte, loadUnloadRecords []*vuv1.Activities_LoadUnloadRecord) ([]byte, error) {
	// VuLoadUnloadRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuLoadUnloadRecord -- Variable size each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x18) // VuLoadUnloadRecord
	recordSize := uint16(0)    // Not used for variable-length records
	noOfRecords := uint16(len(loadUnloadRecords))

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write each VuLoadUnloadRecord
	for _, record := range loadUnloadRecords {
		var err error
		dst, err = appendVuLoadUnloadRecord(dst, record)
		if err != nil {
			return nil, fmt.Errorf("failed to append load/unload record: %w", err)
		}
	}

	return dst, nil
}

// appendVuLoadUnloadRecord appends a single VuLoadUnloadRecord to dst
func appendVuLoadUnloadRecord(dst []byte, record *vuv1.Activities_LoadUnloadRecord) ([]byte, error) {
	if record == nil {
		return dst, nil
	}

	// VuLoadUnloadRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     operationType OperationType,                            -- 1 byte
	//     placeRecord GNSSPlaceRecord,                            -- 14 bytes
	//     odometerValue OdometerShort                             -- 3 bytes
	// }

	// For now, implement simplified versions due to schema limitations
	// These would need to be completed when the protobuf schema is updated

	// Append cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	// Note: Schema limitation - using placeholder for now
	dst = append(dst, make([]byte, 20)...)

	// Append cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	// Note: Schema limitation - using placeholder for now
	dst = append(dst, make([]byte, 20)...)

	// Append operationType (OperationType - 1 byte)
	dst = append(dst, byte(record.GetOperationType()))

	// Append placeRecord (GNSSPlaceRecord - 14 bytes)
	placeRecord := record.GetPlaceRecord()
	if placeRecord != nil {
		var err error
		dst, err = appendVuGNSSADRecord(dst, placeRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to append place record: %w", err)
		}
	} else {
		// Append default place record (14 bytes of zeros)
		dst = append(dst, make([]byte, 14)...)
	}

	// Append odometerValue (OdometerShort - 3 bytes)
	dst = appendVuOdometer(dst, record.GetOdometerKm())

	return dst, nil
}

// appendSignatureRecordArray appends SignatureRecordArray to dst
func appendSignatureRecordArray(dst []byte, signature []byte) ([]byte, error) {
	// SignatureRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes (not used for single signature)
	//     records SET SIZE(noOfRecords) OF Signature -- Variable size each
	// }

	// Write record array header (6 bytes total)
	recordType := uint16(0x08) // Signature
	recordSize := uint16(len(signature))
	noOfRecords := uint16(1) // Single signature

	dst = binary.BigEndian.AppendUint16(dst, recordType)
	dst = binary.BigEndian.AppendUint16(dst, recordSize)
	dst = binary.BigEndian.AppendUint16(dst, noOfRecords)

	// Write the signature data
	if len(signature) > 0 {
		dst = append(dst, signature...)
	}

	return dst, nil
}
