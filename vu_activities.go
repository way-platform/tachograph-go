package tachograph

import (
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
	odometerValue, offset, err := readVuOdometerFromBytes(data, offset)
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

func parseVuCardIWData(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWData ::= SEQUENCE {
	//     noOfIWRecords INTEGER(0..255),
	//     vuCardIWRecords SET SIZE(noOfIWRecords) OF VuCardIWRecord
	// }

	// Read number of records (1 byte)
	noOfRecords, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read number of IW records: %w", err)
	}

	var records []*vuv1.Activities_CardIWRecord

	// Parse each VuCardIWRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuCardIWRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse IW record %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuCardIWRecord parses a single VuCardIWRecord
func parseVuCardIWRecord(data []byte, offset int) (*vuv1.Activities_CardIWRecord, int, error) {
	// VuCardIWRecord ::= SEQUENCE {
	//     cardHolderName HolderName,                    -- Variable length
	//     fullCardNumber FullCardNumber,                -- 19 bytes
	//     cardExpiryDate Datef,                         -- 4 bytes
	//     cardInsertionTime TimeReal,                   -- 4 bytes
	//     vehicleOdometerValueAtInsertion OdometerShort, -- 3 bytes
	//     cardSlotNumber CardSlotNumber,                -- 1 byte
	//     cardWithdrawalTime TimeReal,                  -- 4 bytes
	//     vehicleOdometerValueAtWithdrawal OdometerShort, -- 3 bytes
	//     previousVehicleInfo PreviousVehicleInfo,      -- Variable length
	//     manualInputFlag ManualInputFlag               -- 1 byte
	// }

	record := &vuv1.Activities_CardIWRecord{}

	// Parse cardHolderName (HolderName - variable length)
	// HolderName is typically 72 bytes (36 for surname + 36 for first names)
	holderNameData, offset, err := readBytesFromBytes(data, offset, 72)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read card holder name: %w", err)
	}
	holderName, err := unmarshalHolderName(holderNameData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal holder name: %w", err)
	}
	record.SetCardHolderName(holderName)

	// Parse fullCardNumber (FullCardNumber - 19 bytes)
	fullCardNumberData, offset, err := readBytesFromBytes(data, offset, 19)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read full card number: %w", err)
	}
	fullCardNumber, err := unmarshalFullCardNumber(fullCardNumberData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal full card number: %w", err)
	}
	record.SetFullCardNumber(fullCardNumber)

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
	odometerAtInsertion, offset, err := readVuOdometerFromBytes(data, offset)
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
	odometerAtWithdrawal, offset, err := readVuOdometerFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer at withdrawal: %w", err)
	}
	record.SetOdometerAtWithdrawalKm(int32(odometerAtWithdrawal))

	// Parse previousVehicleInfo (PreviousVehicleInfo - variable length)
	// This is typically 15 bytes (VehicleRegistrationIdentification)
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

// parseActivityChangeInfo parses a single ActivityChangeInfo record
func parseActivityChangeInfo(data []byte, offset int) (*ddv1.ActivityChangeInfo, int, error) {
	// ActivityChangeInfo ::= OCTET STRING (SIZE (2))
	// Bit-packed format: 'scpaattttttttttt'B (16 bits)
	// s: Slot (0=DRIVER, 1=CO-DRIVER)
	// c: Driving status (0=SINGLE, 1=CREW)
	// p: Card status (0=INSERTED, 1=NOT_INSERTED)
	// aa: Activity (00=BREAK/REST, 01=AVAILABILITY, 10=WORK, 11=DRIVING)
	// ttttttttttt: Time (11 bits for time in minutes)

	if offset+2 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for activity change info: got %d, need 2", len(data)-offset)
	}

	// Read 2 bytes
	value := uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2

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

	return change, offset, nil
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

// parseVuPlaceRecord parses a single PlaceRecord
func parseVuPlaceRecord(data []byte, offset int) (*vuv1.Activities_PlaceRecord, int, error) {
	// PlaceRecord ::= SEQUENCE {
	//     entryTime TimeReal,                           -- 4 bytes
	//     entryTypeDailyWorkPeriod EntryTypeDailyWorkPeriod, -- 1 byte
	//     dailyWorkPeriodCountry NationNumeric,         -- 1 byte
	//     dailyWorkPeriodRegion RegionNumeric,          -- 1 byte
	//     vehicleOdometerValue OdometerShort            -- 3 bytes
	// }

	record := &vuv1.Activities_PlaceRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read entry time: %w", err)
	}
	record.SetEntryTime(timestamppb.New(time.Unix(entryTime, 0)))

	// Parse entryTypeDailyWorkPeriod (1 byte)
	entryType, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read entry type: %w", err)
	}
	record.SetEntryType(ddv1.EntryTypeDailyWorkPeriod(entryType))

	// Parse dailyWorkPeriodCountry (NationNumeric - 1 byte)
	country, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read country: %w", err)
	}
	record.SetCountry(ddv1.NationNumeric(country))

	// Parse dailyWorkPeriodRegion (RegionNumeric - 1 byte)
	region, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read region: %w", err)
	}
	record.SetRegion([]byte{region})

	// Parse vehicleOdometerValue (OdometerShort - 3 bytes)
	odometerValue, offset, err := readVuOdometerFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, offset, nil
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

// parseVuSpecificConditionRecord parses a single SpecificConditionRecord
func parseVuSpecificConditionRecord(data []byte, offset int) (*vuv1.Activities_SpecificConditionRecord, int, error) {
	// SpecificConditionRecord ::= SEQUENCE {
	//     entryTime TimeReal,                    -- 4 bytes
	//     specificConditionType SpecificConditionType -- 1 byte
	// }

	record := &vuv1.Activities_SpecificConditionRecord{}

	// Parse entryTime (TimeReal - 4 bytes)
	entryTime, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read entry time: %w", err)
	}
	record.SetEntryTime(timestamppb.New(time.Unix(entryTime, 0)))

	// Parse specificConditionType (1 byte)
	conditionType, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read specific condition type: %w", err)
	}
	record.SetSpecificConditionType(ddv1.SpecificConditionType(conditionType))

	return record, offset, nil
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
		odometerValue, newOffset, err := readVuOdometerFromBytes(data, offset)
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
	record.SetFullCardNumber(&ddv1.FullCardNumber{}) // Placeholder

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
	odometerAtInsertion, offset, err := readVuOdometerFromBytes(data, offset)
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
	odometerAtWithdrawal, offset, err := readVuOdometerFromBytes(data, offset)
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

	var changes []*ddv1.ActivityChangeInfo

	// Parse each ActivityChangeInfo record
	for i := 0; i < int(noOfRecords); i++ {
		change, newOffset, err := parseActivityChangeInfo(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse ActivityChangeInfo record %d: %w", i, err)
		}
		changes = append(changes, change)
		offset = newOffset
	}

	return changes, offset, nil
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

	var records []*vuv1.Activities_PlaceRecord

	// Parse each PlaceRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuPlaceRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse PlaceRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

func parseVuGNSSADRecordArray(data []byte, offset int) ([]*vuv1.Activities_GnssRecord, int, error) {
	// VuGNSSADRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuGNSSADRecord -- Variable size each
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

	var records []*vuv1.Activities_GnssRecord

	// Parse each VuGNSSADRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuGNSSADRecord(data, offset, int(recordSize))
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuGNSSADRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuGNSSADRecord parses a single VuGNSSADRecord
func parseVuGNSSADRecord(data []byte, offset int, recordSize int) (*vuv1.Activities_GnssRecord, int, error) {
	// VuGNSSADRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                    -- 4 bytes
	//     gnssAccuracy GNSSAccuracy,            -- 1 byte
	//     geoCoordinates GeoCoordinates,        -- 8 bytes (latitude + longitude)
	//     positionAuthenticationStatus PositionAuthenticationStatus -- 1 byte
	// }

	if offset+recordSize > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for GNSS record: got %d, need %d", len(data)-offset, recordSize)
	}

	record := &vuv1.Activities_GnssRecord{}

	// Parse timeStamp (TimeReal - 4 bytes)
	timestamp, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read timestamp: %w", err)
	}
	record.SetTimestamp(timestamppb.New(time.Unix(timestamp, 0)))

	// Parse gnssAccuracy (GNSSAccuracy - 1 byte)
	accuracy, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read GNSS accuracy: %w", err)
	}
	record.SetGnssAccuracy(int32(accuracy))

	// Parse geoCoordinates (GeoCoordinates - 8 bytes: 4 bytes latitude + 4 bytes longitude)
	// Latitude (4 bytes, signed integer)
	latBytes := data[offset : offset+4]
	latitude := int32(binary.BigEndian.Uint32(latBytes))
	offset += 4

	// Longitude (4 bytes, signed integer)
	lonBytes := data[offset : offset+4]
	longitude := int32(binary.BigEndian.Uint32(lonBytes))
	offset += 4

	// Create GeoCoordinates
	geoCoords := &ddv1.GeoCoordinates{}
	geoCoords.SetLatitude(latitude)
	geoCoords.SetLongitude(longitude)
	record.SetGeoCoordinates(geoCoords)

	// Parse positionAuthenticationStatus (PositionAuthenticationStatus - 1 byte)
	authStatus, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read authentication status: %w", err)
	}
	record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus(authStatus))

	return record, offset, nil
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

	var records []*vuv1.Activities_SpecificConditionRecord

	// Parse each SpecificConditionRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuSpecificConditionRecord(data, offset)
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse SpecificConditionRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

func parseVuBorderCrossingRecordArray(data []byte, offset int) ([]*vuv1.Activities_BorderCrossingRecord, int, error) {
	// VuBorderCrossingRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuBorderCrossingRecord -- Variable size each
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

	var records []*vuv1.Activities_BorderCrossingRecord

	// Parse each VuBorderCrossingRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuBorderCrossingRecord(data, offset, int(recordSize))
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuBorderCrossingRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuBorderCrossingRecord parses a single VuBorderCrossingRecord
func parseVuBorderCrossingRecord(data []byte, offset int, recordSize int) (*vuv1.Activities_BorderCrossingRecord, int, error) {
	// VuBorderCrossingRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     countryLeft NationNumeric,                              -- 1 byte
	//     countryEntered NationNumeric,                           -- 1 byte
	//     placeRecord GNSSPlaceRecord,                            -- 14 bytes (timestamp + coords + auth)
	//     odometerValue OdometerShort                             -- 3 bytes
	// }

	if offset+recordSize > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for border crossing record: got %d, need %d", len(data)-offset, recordSize)
	}

	record := &vuv1.Activities_BorderCrossingRecord{}

	// Parse cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	driverCardData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read driver card data: %w", err)
	}
	_, err = unmarshalFullCardNumberAndGeneration(driverCardData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal driver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	record.SetFullCardNumber(&ddv1.FullCardNumber{})

	// Parse cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	codriverCardData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read codriver card data: %w", err)
	}
	_, err = unmarshalFullCardNumberAndGeneration(codriverCardData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal codriver card data: %w", err)
	}

	// Parse countryLeft (NationNumeric - 1 byte)
	countryLeft, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read country left: %w", err)
	}
	record.SetCountryLeft(ddv1.NationNumeric(countryLeft))

	// Parse countryEntered (NationNumeric - 1 byte)
	countryEntered, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read country entered: %w", err)
	}
	record.SetCountryEntered(ddv1.NationNumeric(countryEntered))

	// Parse placeRecord (GNSSPlaceRecord - 14 bytes: 4 timestamp + 1 accuracy + 8 coords + 1 auth)
	placeRecord, offset, err := parseGNSSPlaceRecord(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse place record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	// Parse odometerValue (OdometerShort - 3 bytes)
	odometerValue, offset, err := readVuOdometerFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, offset, nil
}

// parseGNSSPlaceRecord parses a GNSSPlaceRecord (simplified version)
func parseGNSSPlaceRecord(data []byte, offset int) (*vuv1.Activities_GnssRecord, int, error) {
	// GNSSPlaceRecord ::= SEQUENCE {
	//     timeStamp TimeReal,                    -- 4 bytes
	//     gnssAccuracy GNSSAccuracy,            -- 1 byte
	//     geoCoordinates GeoCoordinates,        -- 8 bytes
	//     positionAuthenticationStatus PositionAuthenticationStatus -- 1 byte
	// }

	record := &vuv1.Activities_GnssRecord{}

	// Parse timeStamp (TimeReal - 4 bytes)
	timestamp, offset, err := readVuTimeRealFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read timestamp: %w", err)
	}
	record.SetTimestamp(timestamppb.New(time.Unix(timestamp, 0)))

	// Parse gnssAccuracy (GNSSAccuracy - 1 byte)
	accuracy, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read GNSS accuracy: %w", err)
	}
	record.SetGnssAccuracy(int32(accuracy))

	// Parse geoCoordinates (GeoCoordinates - 8 bytes: 4 bytes latitude + 4 bytes longitude)
	// Latitude (4 bytes, signed integer)
	latBytes := data[offset : offset+4]
	latitude := int32(binary.BigEndian.Uint32(latBytes))
	offset += 4

	// Longitude (4 bytes, signed integer)
	lonBytes := data[offset : offset+4]
	longitude := int32(binary.BigEndian.Uint32(lonBytes))
	offset += 4

	// Create GeoCoordinates
	geoCoords := &ddv1.GeoCoordinates{}
	geoCoords.SetLatitude(latitude)
	geoCoords.SetLongitude(longitude)
	record.SetGeoCoordinates(geoCoords)

	// Parse positionAuthenticationStatus (PositionAuthenticationStatus - 1 byte)
	authStatus, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read authentication status: %w", err)
	}
	record.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus(authStatus))

	return record, offset, nil
}

func parseVuLoadUnloadRecordArray(data []byte, offset int) ([]*vuv1.Activities_LoadUnloadRecord, int, error) {
	// VuLoadUnloadRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF VuLoadUnloadRecord -- Variable size each
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

	var records []*vuv1.Activities_LoadUnloadRecord

	// Parse each VuLoadUnloadRecord
	for i := 0; i < int(noOfRecords); i++ {
		record, newOffset, err := parseVuLoadUnloadRecord(data, offset, int(recordSize))
		if err != nil {
			return nil, offset, fmt.Errorf("failed to parse VuLoadUnloadRecord %d: %w", i, err)
		}
		records = append(records, record)
		offset = newOffset
	}

	return records, offset, nil
}

// parseVuLoadUnloadRecord parses a single VuLoadUnloadRecord
func parseVuLoadUnloadRecord(data []byte, offset int, recordSize int) (*vuv1.Activities_LoadUnloadRecord, int, error) {
	// VuLoadUnloadRecord ::= SEQUENCE {
	//     cardNumberAndGenDriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     cardNumberAndGenCodriverSlot FullCardNumberAndGeneration, -- 20 bytes
	//     operationType OperationType,                            -- 1 byte
	//     placeRecord GNSSPlaceRecord,                            -- 14 bytes
	//     odometerValue OdometerShort                             -- 3 bytes
	// }

	if offset+recordSize > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for load/unload record: got %d, need %d", len(data)-offset, recordSize)
	}

	record := &vuv1.Activities_LoadUnloadRecord{}

	// Parse cardNumberAndGenDriverSlot (FullCardNumberAndGeneration - 20 bytes)
	driverCardData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read driver card data: %w", err)
	}
	_, err = unmarshalFullCardNumberAndGeneration(driverCardData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal driver card data: %w", err)
	}
	// Note: Schema limitation - using placeholder for now
	record.SetFullCardNumber(&ddv1.FullCardNumber{})

	// Parse cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration - 20 bytes)
	codriverCardData, offset, err := readBytesFromBytes(data, offset, 20)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read codriver card data: %w", err)
	}
	_, err = unmarshalFullCardNumberAndGeneration(codriverCardData)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to unmarshal codriver card data: %w", err)
	}

	// Parse operationType (OperationType - 1 byte)
	operationType, offset, err := readUint8FromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read operation type: %w", err)
	}
	record.SetOperationType(ddv1.OperationType(operationType))

	// Parse placeRecord (GNSSPlaceRecord - 14 bytes)
	placeRecord, offset, err := parseGNSSPlaceRecord(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to parse place record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	// Parse odometerValue (OdometerShort - 3 bytes)
	odometerValue, offset, err := readVuOdometerFromBytes(data, offset)
	if err != nil {
		return nil, offset, fmt.Errorf("failed to read odometer value: %w", err)
	}
	record.SetOdometerKm(int32(odometerValue))

	return record, offset, nil
}

func parseSignatureRecordArray(data []byte, offset int) ([]byte, int, error) {
	// SignatureRecordArray ::= SEQUENCE {
	//     recordType INTEGER(1..65535),           -- 2 bytes
	//     recordSize INTEGER(0..65535),           -- 2 bytes
	//     noOfRecords INTEGER(0..65535),          -- 2 bytes
	//     records SET SIZE(noOfRecords) OF Signature -- Variable size each
	// }

	// Read record array header (6 bytes total)
	if offset+6 > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for signature record array header: got %d, need 6", len(data)-offset)
	}

	// Read recordType (2 bytes)
	recordType := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read recordSize (2 bytes)
	recordSize := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Read noOfRecords (2 bytes) - not used for single signature
	_ = binary.BigEndian.Uint16(data[offset:])
	offset += 2

	// Validate record type (should be 0x08 for Signature)
	if recordType != 0x08 {
		return nil, offset, fmt.Errorf("unexpected record type: got %d, expected 8 (Signature)", recordType)
	}

	// For Gen2, signatures are typically ECC (64 bytes) or RSA (128 bytes)
	// We'll use the recordSize to determine the actual signature size
	if recordSize == 0 {
		return nil, offset, fmt.Errorf("invalid record size: got 0")
	}

	// Read the signature data
	if offset+int(recordSize) > len(data) {
		return nil, offset, fmt.Errorf("insufficient data for signature: got %d, need %d", len(data)-offset, recordSize)
	}

	signature := make([]byte, recordSize)
	copy(signature, data[offset:offset+int(recordSize)])
	offset += int(recordSize)

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
func appendVuActivities(buf *bytes.Buffer, activities *vuv1.Activities) error {
	if activities == nil {
		return nil
	}

	// For now, implement a simplified version that writes the signature data
	// This ensures the interface is complete while allowing for future enhancement

	if activities.GetGeneration() == ddv1.Generation_GENERATION_1 {
		signature := activities.GetSignatureGen1()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	} else {
		signature := activities.GetSignatureGen2()
		if len(signature) > 0 {
			buf.Write(signature)
		}
	}

	return nil
}
