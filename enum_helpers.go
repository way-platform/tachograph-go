package tachograph

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// Helper function to find enum value by protocol annotation
func findEnumValueByProtocol(enumDesc protoreflect.EnumDescriptor, rawValue int32) (protoreflect.EnumNumber, bool) {
	values := enumDesc.Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()

		// Check if this value has the protocol_enum_value annotation
		if proto.HasExtension(opts, datadictionaryv1.E_ProtocolEnumValue) {
			protocolValue := proto.GetExtension(opts, datadictionaryv1.E_ProtocolEnumValue).(int32)
			if protocolValue == rawValue {
				return valueDesc.Number(), true
			}
		}
	}
	return 0, false
}

// Helper function to get protocol value from enum
func getProtocolValueFromEnumNumber(enumDesc protoreflect.EnumDescriptor, enumNumber protoreflect.EnumNumber) (int32, bool) {
	valueDesc := enumDesc.Values().ByNumber(enumNumber)
	if valueDesc == nil {
		return 0, false
	}

	opts := valueDesc.Options()
	if !proto.HasExtension(opts, datadictionaryv1.E_ProtocolEnumValue) {
		return 0, false
	}

	protocolValue := proto.GetExtension(opts, datadictionaryv1.E_ProtocolEnumValue).(int32)
	return protocolValue, true
}

// Specific helper functions for common enum types

// SetEventFaultType converts a raw protocol value to EventFaultType enum
func SetEventFaultType(rawValue int32, setEnum func(datadictionaryv1.EventFaultType), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.EventFaultType_EVENT_FAULT_TYPE_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.EventFaultType(enumNumber))
	} else {
		setEnum(datadictionaryv1.EventFaultType_EVENT_FAULT_TYPE_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetCalibrationPurpose converts a raw protocol value to CalibrationPurpose enum
func SetCalibrationPurpose(rawValue int32, setEnum func(datadictionaryv1.CalibrationPurpose), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.CalibrationPurpose(enumNumber))
	} else {
		setEnum(datadictionaryv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetDriverActivityValue converts a raw protocol value to DriverActivityValue enum
func SetDriverActivityValue(rawValue int32, setEnum func(datadictionaryv1.DriverActivityValue), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.DriverActivityValue_DRIVER_ACTIVITY_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.DriverActivityValue(enumNumber))
	} else {
		setEnum(datadictionaryv1.DriverActivityValue_DRIVER_ACTIVITY_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetCardSlotNumber converts a raw protocol value to CardSlotNumber enum
func SetCardSlotNumber(rawValue int32, setEnum func(datadictionaryv1.CardSlotNumber), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.CardSlotNumber(enumNumber))
	} else {
		setEnum(datadictionaryv1.CardSlotNumber_CARD_SLOT_NUMBER_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetCardStatus converts a raw protocol value to CardStatus enum
func SetCardStatus(rawValue int32, setEnum func(datadictionaryv1.CardStatus), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.CardStatus_CARD_STATUS_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.CardStatus(enumNumber))
	} else {
		setEnum(datadictionaryv1.CardStatus_CARD_STATUS_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetDrivingStatus converts a raw protocol value to DrivingStatus enum
func SetDrivingStatus(rawValue int32, setEnum func(datadictionaryv1.DrivingStatus), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.DrivingStatus_DRIVING_STATUS_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.DrivingStatus(enumNumber))
	} else {
		setEnum(datadictionaryv1.DrivingStatus_DRIVING_STATUS_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetEquipmentType converts a raw protocol value to EquipmentType enum
func SetEquipmentType(rawValue int32, setEnum func(datadictionaryv1.EquipmentType), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.EquipmentType(enumNumber))
	} else {
		setEnum(datadictionaryv1.EquipmentType_EQUIPMENT_TYPE_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetEventFaultRecordPurpose converts a raw protocol value to EventFaultRecordPurpose enum
func SetEventFaultRecordPurpose(rawValue int32, setEnum func(datadictionaryv1.EventFaultRecordPurpose), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.EventFaultRecordPurpose(enumNumber))
	} else {
		setEnum(datadictionaryv1.EventFaultRecordPurpose_EVENT_FAULT_RECORD_PURPOSE_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// Marshalling helper functions

// GetEventFaultTypeProtocolValue returns the protocol value for marshalling
func GetEventFaultTypeProtocolValue(enumValue datadictionaryv1.EventFaultType, unrecognizedValue int32) int32 {
	// Check if this is the UNRECOGNIZED value
	if enumValue == datadictionaryv1.EventFaultType_EVENT_FAULT_TYPE_UNRECOGNIZED {
		return unrecognizedValue
	}

	enumDesc := enumValue.Descriptor()
	if protocolValue, ok := getProtocolValueFromEnumNumber(enumDesc, protoreflect.EnumNumber(enumValue)); ok {
		return protocolValue
	}

	// Fallback - this shouldn't happen in well-formed data
	return int32(enumValue)
}

// GetCalibrationPurposeProtocolValue returns the protocol value for marshalling
func GetCalibrationPurposeProtocolValue(enumValue datadictionaryv1.CalibrationPurpose, unrecognizedValue int32) int32 {
	if enumValue == datadictionaryv1.CalibrationPurpose_CALIBRATION_PURPOSE_UNRECOGNIZED {
		return unrecognizedValue
	}

	enumDesc := enumValue.Descriptor()
	if protocolValue, ok := getProtocolValueFromEnumNumber(enumDesc, protoreflect.EnumNumber(enumValue)); ok {
		return protocolValue
	}

	return int32(enumValue)
}

// GetDriverActivityValueProtocolValue returns the protocol value for marshalling
func GetDriverActivityValueProtocolValue(enumValue datadictionaryv1.DriverActivityValue, unrecognizedValue int32) int32 {
	if enumValue == datadictionaryv1.DriverActivityValue_DRIVER_ACTIVITY_UNRECOGNIZED {
		return unrecognizedValue
	}

	enumDesc := enumValue.Descriptor()
	if protocolValue, ok := getProtocolValueFromEnumNumber(enumDesc, protoreflect.EnumNumber(enumValue)); ok {
		return protocolValue
	}

	return int32(enumValue)
}

// SetEntryTypeDailyWorkPeriod converts a raw protocol value to EntryTypeDailyWorkPeriod enum
func SetEntryTypeDailyWorkPeriod(rawValue int32, setEnum func(datadictionaryv1.EntryTypeDailyWorkPeriod), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.EntryTypeDailyWorkPeriod(enumNumber))
	} else {
		setEnum(datadictionaryv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetNationNumeric converts a raw protocol value to NationNumeric enum
func SetNationNumeric(rawValue int32, setEnum func(datadictionaryv1.NationNumeric), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.NationNumeric(enumNumber))
	} else {
		setEnum(datadictionaryv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}

// SetSpecificConditionType converts a raw protocol value to SpecificConditionType enum
func SetSpecificConditionType(rawValue int32, setEnum func(datadictionaryv1.SpecificConditionType), setUnrecognized func(int32)) {
	enumDesc := datadictionaryv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNSPECIFIED.Descriptor()
	if enumNumber, found := findEnumValueByProtocol(enumDesc, rawValue); found {
		setEnum(datadictionaryv1.SpecificConditionType(enumNumber))
	} else {
		setEnum(datadictionaryv1.SpecificConditionType_SPECIFIC_CONDITION_TYPE_UNRECOGNIZED)
		setUnrecognized(rawValue)
	}
}
