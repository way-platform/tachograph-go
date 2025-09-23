# Specification: Protobuf Enum Refactoring

## 1. Introduction

This document specifies a refactoring of the tachograph protobuf schemas to improve type safety, readability, and long-term maintainability. Currently, many fields that represent fixed sets of coded values (e.g., activity types, event types, purposes) are defined as raw integers (`int32`).

This refactoring will introduce strongly-typed `enum` definitions for these fields, based on the official specifications in the EU regulation's Data Dictionary.

## 2. Guiding Principles

- **Information Preservation**: The refactoring must be lossless. The ability to round-trip data (unmarshal then marshal) without any change to the binary output is critical.
- **Clarity and Type Safety**: The primary goal is to replace "magic numbers" with named, self-documenting enum values.
- **Forward Compatibility**: The system must gracefully handle unknown or manufacturer-specific values that are not defined in the standard enums.

## 3. Refactoring Strategy

To ensure information preservation and forward compatibility, the following pattern will be applied to each field being converted from an integer to an enum:

1.  **Create a comprehensive `enum`**: A new `enum` will be defined, typically in a shared `.proto` file (e.g., `proto/wayplatform/connect/tachograph/common/v1/enums.proto`) to promote reuse. The enum will contain all known values from the Data Dictionary.
2.  **Add an `_UNRECOGNIZED` value**: Each enum will include a value named `<ENUM_NAME>_UNRECOGNIZED`. This will serve as the fallback for any raw integer value encountered during unmarshalling that does not correspond to a known enum member.
3.  **Modify the message**: The original `int32` field in the message will be changed to the new `enum` type.
4.  **Add an `unrecognized_` field**: A new `int32` field, named `unrecognized_<original_field_name>`, will be added to the message.

### Unmarshalling Logic:

- When parsing the raw data, if the integer value corresponds to a known member of the enum, the enum field is set to that member. The `unrecognized_` field is left at its default value.
- If the integer value does not correspond to a known member, the enum field is set to `<ENUM_NAME>_UNRECOGNIZED`, and the `unrecognized_<original_field_name>` field is populated with the original, raw integer value.

### Marshalling Logic:

- If the enum field is set to any value other than `<ENUM_NAME>_UNRECOGNIZED`, the marshaller will convert it back to its corresponding raw integer value.
- If the enum field is set to `<ENUM_NAME>_UNRECOGNIZED`, the marshaller will use the value from the `unrecognized_<original_field_name>` field as the raw integer value.

This strategy ensures that even if we encounter data with new or proprietary codes, we can still parse, store, and re-serialize it without losing the original information.

## 4. Proposed Schema Changes

The following sections detail the specific fields to be refactored.

### 4.1. Common Enums

A new file, `proto/wayplatform/connect/tachograph/common/v1/enums.proto`, will be created to house enums shared across card and VU schemas.

---

### 4.2. `EventFaultType`

- **Data Dictionary**: Section 2.70
- **Description**: Qualifies dozens of different events and faults across all generations.
- **Affected Files**:
    - `proto/wayplatform/connect/tachograph/card/v1/event_data.proto`
    - `proto/wayplatform/connect/tachograph/card/v1/fault_data.proto`
    - `proto/wayplatform/connect/tachograph/vu/v1/events_and_faults.proto`

#### New Enum Definition (`common/v1/enums.proto`)

```protobuf
enum EventFaultType {
  EVENT_FAULT_TYPE_UNSPECIFIED = 0;
  EVENT_FAULT_TYPE_UNRECOGNIZED = 1;

  // General Events (0x00 - 0x0F)
  GENERAL_NO_FURTHER_DETAILS = 2; // 0x00
  GENERAL_INSERTION_OF_NON_VALID_CARD = 3; // 0x01
  GENERAL_CARD_CONFLICT = 4; // 0x02
  GENERAL_TIME_OVERLAP = 5; // 0x03
  GENERAL_DRIVING_WITHOUT_APPROPRIATE_CARD = 6; // 0x04
  GENERAL_CARD_INSERTION_WHILE_DRIVING = 7; // 0x05
  GENERAL_LAST_CARD_SESSION_NOT_CORRECTLY_CLOSED = 8; // 0x06
  GENERAL_OVER_SPEEDING = 9; // 0x07
  GENERAL_POWER_SUPPLY_INTERRUPTION = 10; // 0x08
  GENERAL_MOTION_DATA_ERROR = 11; // 0x09
  GENERAL_VEHICLE_MOTION_CONFLICT = 12; // 0x0A
  GENERAL_TIME_CONFLICT_GNSS_VS_VU = 13; // 0x0B (Gen2+)
  GENERAL_COMM_ERROR_REMOTE_COMM_FACILITY = 14; // 0x0C (Gen2+)
  GENERAL_ABSENCE_OF_POSITION_INFO_FROM_GNSS = 15; // 0x0D (Gen2+)
  GENERAL_COMM_ERROR_EXTERNAL_GNSS_FACILITY = 16; // 0x0E (Gen2+)
  GENERAL_GNSS_ANOMALY = 17; // 0x0F (Gen2v2+)

  // VU Security Breach Events (0x10 - 0x1F)
  VU_SEC_NO_FURTHER_DETAILS = 18; // 0x10
  VU_SEC_MOTION_SENSOR_AUTH_FAILURE = 19; // 0x11
  VU_SEC_TACHOGRAPH_CARD_AUTH_FAILURE = 20; // 0x12
  VU_SEC_UNAUTHORISED_CHANGE_OF_MOTION_SENSOR = 21; // 0x13
  VU_SEC_CARD_DATA_INPUT_INTEGRITY_ERROR = 22; // 0x14
  VU_SEC_STORED_USER_DATA_INTEGRITY_ERROR = 23; // 0x15
  VU_SEC_INTERNAL_DATA_TRANSFER_ERROR = 24; // 0x16
  VU_SEC_UNAUTHORISED_CASE_OPENING = 25; // 0x17
  VU_SEC_HARDWARE_SABOTAGE = 26; // 0x18
  VU_SEC_TAMPER_DETECTION_OF_GNSS = 27; // 0x19 (Gen2+)
  VU_SEC_EXTERNAL_GNSS_FACILITY_AUTH_FAILURE = 28; // 0x1A (Gen2+)
  VU_SEC_EXTERNAL_GNSS_FACILITY_CERT_EXPIRED = 29; // 0x1B (Gen2+)
  VU_SEC_INCONSISTENCY_MOTION_VS_ACTIVITY = 30; // 0x1C (Gen2v2+)

  // Sensor Security Breach Events (0x20 - 0x2F)
  SENSOR_SEC_NO_FURTHER_DETAILS = 31; // 0x20
  SENSOR_SEC_AUTHENTICATION_FAILURE = 32; // 0x21
  SENSOR_SEC_STORED_DATA_INTEGRITY_ERROR = 33; // 0x22
  SENSOR_SEC_INTERNAL_DATA_TRANSFER_ERROR = 34; // 0x23
  SENSOR_SEC_UNAUTHORISED_CASE_OPENING = 35; // 0x24
  SENSOR_SEC_HARDWARE_SABOTAGE = 36; // 0x25

  // Recording Equipment Faults (0x30 - 0x3F)
  FAULT_REC_EQ_NO_FURTHER_DETAILS = 37; // 0x30
  FAULT_REC_EQ_VU_INTERNAL_FAULT = 38; // 0x31
  FAULT_REC_EQ_PRINTER_FAULT = 39; // 0x32
  FAULT_REC_EQ_DISPLAY_FAULT = 40; // 0x33
  FAULT_REC_EQ_DOWNLOADING_FAULT = 41; // 0x34
  FAULT_REC_EQ_SENSOR_FAULT = 42; // 0x35
  FAULT_REC_EQ_INTERNAL_GNSS_RECEIVER = 43; // 0x36 (Gen2+)
  FAULT_REC_EQ_EXTERNAL_GNSS_FACILITY = 44; // 0x37 (Gen2+)
  FAULT_REC_EQ_REMOTE_COMM_FACILITY = 45; // 0x38 (Gen2+)
  FAULT_REC_EQ_ITS_INTERFACE = 46; // 0x39 (Gen2+)
  FAULT_REC_EQ_INTERNAL_SENSOR_FAULT = 47; // 0x3A (Gen2v2+)

  // Card Faults (0x40 - 0x4F)
  FAULT_CARD_NO_FURTHER_DETAILS = 48; // 0x40
}
```

#### Message Modification Example (`card/v1/event_data.proto`)

**Before:**
```protobuf
message Record {
  int32 event_type = 1;
  // ...
}
```

**After:**
```protobuf
// import "wayplatform/connect/tachograph/common/v1/enums.proto";

message Record {
  EventFaultType event_type = 1;
  // ...
  string vehicle_registration_number = 5;
  // Populated only when event_type is EVENT_FAULT_TYPE_UNRECOGNIZED.
  int32 unrecognized_event_type = 6;
}
```
*(This change would be applied to all messages using `EventFaultType`)*

---

### 4.3. `CalibrationPurpose`

- **Data Dictionary**: Section 2.8
- **Description**: Explains why a calibration was performed.
- **Affected Files**:
    - `proto/wayplatform/connect/tachograph/card/v1/calibrations.proto`
    - `proto/wayplatform/connect/tachograph/vu/v1/technical_data.proto`

#### New Enum Definition (`common/v1/enums.proto`)

```protobuf
enum CalibrationPurpose {
  CALIBRATION_PURPOSE_UNSPECIFIED = 0;
  CALIBRATION_PURPOSE_UNRECOGNIZED = 1;
  RESERVED = 2; // 0x00
  ACTIVATION = 3; // 0x01
  FIRST_INSTALLATION = 4; // 0x02
  INSTALLATION = 5; // 0x03
  PERIODIC_INSPECTION = 6; // 0x04
  VRN_ENTRY_BY_COMPANY = 7; // 0x05 (Gen2+)
  TIME_ADJUSTMENT = 8; // 0x06 (Gen2+)
}
```

#### Message Modification Example (`card/v1/calibrations.proto`)

**Before:**
```protobuf
message Record {
  int32 calibration_purpose = 1;
  // ...
}
```

**After:**
```protobuf
// import "wayplatform/connect/tachograph/common/v1/enums.proto";

message Record {
  CalibrationPurpose calibration_purpose = 1;
  // ...
  string sensor_serial_number = 17;
  // Populated only when calibration_purpose is CALIBRATION_PURPOSE_UNRECOGNIZED.
  int32 unrecognized_calibration_purpose = 18;
}
```

---

### 4.4. Driver Activity Enums

- **Data Dictionary**: Section 2.1
- **Description**: A set of values describing driver status.
- **Affected Files**:
    - `proto/wayplatform/connect/tachograph/card/v1/driver_activity.proto`
    - `proto/wayplatform/connect/tachograph/vu/v1/activities.proto`

#### New Enum Definitions (`common/v1/enums.proto`)

```protobuf
enum DriverActivityValue {
  DRIVER_ACTIVITY_UNSPECIFIED = 0;
  DRIVER_ACTIVITY_UNRECOGNIZED = 1;
  BREAK_REST = 2; // 0b00
  AVAILABILITY = 3; // 0b01
  WORK = 4; // 0b10
  DRIVING = 5; // 0b11
}

enum DrivingStatus {
  DRIVING_STATUS_UNSPECIFIED = 0;
  DRIVING_STATUS_UNRECOGNIZED = 1;
  SINGLE = 2; // 0
  CREW = 3; // 1
}

enum CardStatus {
  CARD_STATUS_UNSPECIFIED = 0;
  CARD_STATUS_UNRECOGNIZED = 1;
  INSERTED = 2; // 0
  NOT_INSERTED = 3; // 1
}

enum CardSlotNumber {
  CARD_SLOT_NUMBER_UNSPECIFIED = 0;
  CARD_SLOT_NUMBER_UNRECOGNIZED = 1;
  DRIVER_SLOT = 2; // 0
  CO_DRIVER_SLOT = 3; // 1
}
```

#### Message Modification Example (`card/v1/driver_activity.proto`)

**Before:**
```protobuf
message ActivityChange {
  int32 slot = 1;
  int32 driving_status = 2;
  int32 card_status = 3;
  int32 activity = 4;
  int32 time_of_change_minutes = 5;
}
```

**After:**
```protobuf
// import "wayplatform/connect/tachograph/common/v1/enums.proto";

message ActivityChange {
  CardSlotNumber slot = 1;
  DrivingStatus driving_status = 2;
  CardStatus card_status = 3;
  DriverActivityValue activity = 4;
  int32 time_of_change_minutes = 5;

  // Unrecognized value fields
  int32 unrecognized_slot = 6;
  int32 unrecognized_driving_status = 7;
  int32 unrecognized_card_status = 8;
  int32 unrecognized_activity = 9;
}
```

---

### 4.5. `EventFaultRecordPurpose`

- **Data Dictionary**: Section 2.69
- **Description**: Explains why an event or fault was recorded by the VU.
- **Affected Files**:
    - `proto/wayplatform/connect/tachograph/vu/v1/events_and_faults.proto`

#### New Enum Definition (`common/v1/enums.proto`)

```protobuf
enum EventFaultRecordPurpose {
  EVENT_FAULT_RECORD_PURPOSE_UNSPECIFIED = 0;
  EVENT_FAULT_RECORD_PURPOSE_UNRECOGNIZED = 1;
  TEN_MOST_RECENT = 2; // 0x00
  LONGEST_IN_LAST_10_DAYS = 3; // 0x01
  FIVE_LONGEST_IN_LAST_365_DAYS = 4; // 0x02
  LAST_IN_LAST_10_DAYS = 5; // 0x03
  MOST_SERIOUS_IN_LAST_10_DAYS = 6; // 0x04
  FIVE_MOST_SERIOUS_IN_LAST_365_DAYS = 7; // 0x05
  FIRST_AFTER_LAST_CALIBRATION = 8; // 0x06
  ACTIVE_OR_ONGOING = 9; // 0x07
}
```

#### Message Modification Example (`vu/v1/events_and_faults.proto`)

**Before:**
```protobuf
message FaultRecord {
  int32 fault_type = 1;
  int32 record_purpose = 2;
  // ...
}
```

**After:**
```protobuf
// import "wayplatform/connect/tachograph/common/v1/enums.proto";

message FaultRecord {
  EventFaultType fault_type = 1;
  EventFaultRecordPurpose record_purpose = 2;
  // ...
  Generation card_generation = 6;
  int32 unrecognized_fault_type = 7;
  int32 unrecognized_record_purpose = 8;
}
```
*(This change would apply to `FaultRecord` and `EventRecord` in the same file)*

---

### 4.6. Other Candidates

The following fields are also strong candidates for this refactoring and should be implemented using the same strategy.

| Field Name(s) | Data Dictionary | Affected Schemas |
|---|---|---|
| `entry_type` | 2.66 `EntryTypeDailyWorkPeriod` | `card/v1/places.proto`, `vu/v1/activities.proto` |
| `specific_condition_type` | 2.154 `SpecificConditionType` | `card/v1/specific_conditions.proto`, `vu/v1/activities.proto` |
| `type_of_tachograph_card_id` | 2.67 `EquipmentType` | `card/v1/*_application_identification.proto` |
| `company_activity_type` | 2.47 `CompanyActivityType` | `card/v1/company_activity_data.proto` |
| `by_default_load_type` | 2.90a `LoadType` | `card/v1/calibrations_add_data.proto` |
| `authentication_status` | 2.117a `PositionAuthenticationStatus` | `card/v1/gnss_*.proto`, `card/v1/places_*.proto`, `vu/v1/activities.proto` |
| `manual_input_flag` | 2.93 `ManualInputFlag` | `vu/v1/activities.proto` |
| `country_*`, `*nation` | 2.101 `NationNumeric` | Multiple schemas |

Due to the large number of nations, the `NationNumeric` enum will be extensive but provides significant value. It should be included.

## 5. Implementation Plan

1.  Create the new shared file: `proto/wayplatform/connect/tachograph/common/v1/enums.proto`.
2.  Populate this file with all the new `enum` definitions as specified above.
3.  Iterate through each of the "Affected Files" listed in this document.
4.  In each file, add an `import` statement for the new `enums.proto` file.
5.  Modify the relevant messages to replace the `int32` field with the new `enum` type and add the corresponding `unrecognized_<fieldname>` field.
6.  After schema changes are complete, update the unmarshalling and marshalling logic in the Go source code to correctly map between raw integer values and the new enum types, following the strategy outlined in Section 3.
