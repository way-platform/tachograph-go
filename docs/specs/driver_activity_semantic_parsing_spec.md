# Specification: Semantic Parsing of EF_Driver_Activity_Data

## 1. Introduction

This document specifies a comprehensive refactoring of the parsing logic for the `EF_Driver_Activity_Data` (FID `0x0504`) elementary file. The goal is to replace the current placeholder implementation with a fully semantic, robust parsing strategy that aligns with the project's principles of data completeness and roundtrip integrity.

The current implementation in `unmarshal_card_activity.go` is a stub that prioritizes raw data preservation over semantic understanding. It marks the entire data block as invalid and stores it as a raw byte slice. This specification outlines the necessary changes to both the Protobuf schema and the Go parsing logic to correctly interpret the daily activity records, handle data corruption gracefully, and ensure byte-perfect roundtrips.

## 2. Problem Analysis

The current implementation falls short of the desired state in several key areas:

1.  **Incomplete Parsing Logic**: The primary function `unmarshalCardActivityData` does not perform a semantic parse. It reads the two header pointers (`activityPointerOldestDayRecord` and `activityPointerNewestRecord`) and then stores the entire remaining ring buffer as a single `raw_data` blob. The existing but unused `parseActivityDailyRecords` function contains logic for parsing, but it is not integrated.

2.  **Insufficient Protobuf Schema**: The `DriverActivity` message in `driver_activity.proto` uses a top-level `bool valid` flag and `bytes raw_data` field. This is an all-or-nothing approach. If a single byte in one daily record is corrupted, the entire dataset is treated as a raw blob, losing all semantic information from the valid records within the buffer.

3.  **Principle Violation**: This approach violates the core principle of isolating invalidity. The goal is to parse as much data as possible semantically and only store raw bytes for the specific, individual records that are malformed.

4.  **Naming Inconsistency**: The top-level message is named `DriverActivity`, while the Elementary File is `EF_Driver_Activity_Data`. A more consistent name would be `DriverActivityData`.

## 3. Proposed Solution

The solution involves a two-pronged approach: refactoring the Protobuf schema to support per-record validity and rewriting the Go unmarshalling code to implement full semantic parsing with robust error handling.

### 3.1. Protobuf Schema Refactoring (`driver_activity.proto`)

The schema will be updated to precisely model the desired state: a collection of daily records, each of which can be either valid or raw.

1.  **Rename Top-Level Message**: `DriverActivity` will be renamed to `DriverActivityData` to align with the EF name.
2.  **Remove Global Raw Fields**: The top-level `bool valid` and `bytes raw_data` fields will be removed. The new design will always attempt a full semantic parse.
3.  **Introduce Per-Record Validity**: A `oneof` field will be introduced in the `DailyRecord` message. This allows a record to be represented as either:
    *   A `ParsedDailyRecord` message containing all the structured fields (`activityRecordDate`, `activityDayDistance`, `activityChangeInfo`, etc.).
    *   A `bytes raw_record` field containing the raw byte slice for a single, unparseable daily record.

This change is the key to isolating failures and achieving both semantic parsing and perfect roundtrip capability.

### 3.2. Go Unmarshalling Logic (`unmarshal_card_activity.go`)

The Go code will be updated to implement the new schema and the correct parsing strategy for a cyclic, backward-linked buffer.

1.  **Activate Parsing**: `unmarshalDriverActivityData` will be modified to call a new `parseCyclicActivityDailyRecords` function, passing it the ring buffer data slice and the `newestDayRecordPointer`.
2.  **Implement Cyclic Parsing**: `parseCyclicActivityDailyRecords` will be implemented to parse the data as a cyclic, backward-linked list.
    *   The parsing loop will start at the offset given by `newestDayRecordPointer`.
    *   In each iteration, it will read the current record's data (handling buffer wrap-around if the record spans the end of the buffer). It will then use the `activityPreviousRecordLength` from that record's header to calculate the starting position of the *next* record to parse (which is the previous one chronologically).
    *   The loop will terminate when a `previousRecordLength` of zero is found, or after a safeguard number of iterations (e.g., 366) to prevent infinite loops on corrupted data.
3.  **Robust Record Handling**: Inside the loop, it will pass the byte slice for a single record to a `parseSingleActivityDailyRecord` function.
    *   **On Success**: If parsing succeeds, the returned `DailyRecord` message will be used to populate the `parsed` field of the `ActivityDailyRecord` container.
    *   **On Failure**: If parsing fails, the original byte slice for that record will be stored in the `raw` field of the `ActivityDailyRecord` container.
4.  **Reverse Order**: Since the records are parsed from newest to oldest, the resulting slice of records must be reversed at the end to restore chronological order.
5.  **Complete `parseSingleActivityDailyRecord`**: This function will be finalized to correctly parse all fields within a single daily record, including the variable-length list of 2-byte `ActivityChangeInfo` bitfields.

## 4. Detailed Schema Changes

The file `driver_activity_data.proto` will be modified as follows.
(Note: The schema below reflects the final version after several iterations of feedback, using a nested `DailyRecord` with a `valid` flag and `raw` field, which is the simplest and most direct representation.)

**Before:**
```protobuf
// Represents driver activity data from a tachograph card.
//
// Corresponds to the `CardDriverActivity` data type.
// See Data Dictionary, Section 2.17.
message DriverActivity {
  // If true, the fields below are populated with parsed, semantic data.
  // If false, the 'raw_data' field contains the original, unprocessed ring buffer
  // bytes for perfect roundtrip accuracy.
  bool valid = 1;

  // --- Fields for valid data (when valid = true) ---

  // Represents a record of driver activity for a single day.
  //
  // Corresponds to the `CardActivityDailyRecord` data type.
  // See Data Dictionary, Section 2.9.
  message DailyRecord {
    // ... fields
  }

  // See Data Dictionary, Section 2.17, `activityPointerOldestDayRecord`.
  int32 oldest_day_record_index = 2;

  // See Data Dictionary, Section 2.17, `activityPointerNewestRecord`.
  int32 newest_day_record_index = 3;

  // See Data Dictionary, Section 2.17, `activityDailyRecords`.
  repeated DailyRecord daily_records = 4;

  // --- Field for raw data preservation (when valid = false) ---
  // Holds the raw ring buffer bytes (after the 4-byte header) for perfect roundtrip.
  bytes raw_data = 5;

  // Digital signature for the EF_Driver_Activity_Data file content.
  bytes signature = 6;
}
```

**After:**
```protobuf
edition = "2023";

package wayplatform.connect.tachograph.card.v1;

import "google/protobuf/timestamp.proto";
import "wayplatform/connect/tachograph/datadictionary/v1/card_slot_number.proto";
import "wayplatform/connect/tachograph/datadictionary/v1/card_status.proto";
import "wayplatform/connect/tachograph/datadictionary/v1/driver_activity_value.proto";
import "wayplatform/connect/tachograph/datadictionary/v1/driving_status.proto";

// Represents the EF_Driver_Activity_Data file from a tachograph card.
//
// Corresponds to the `CardDriverActivity` data type.
// See Data Dictionary, Section 2.17.
message DriverActivityData {
  // Represents a single daily activity record, which can either be fully parsed
  // or stored as raw bytes, indicated by the `valid` flag.
  // Corresponds to the `CardActivityDailyRecord` data type.
  message DailyRecord {
    // Represents a change in driver activity, driving status, or card status.
    //
    // Corresponds to the `ActivityChangeInfo` data type.
    // See Data Dictionary, Section 2.1.
    message ActivityChange {
      // Slot of the driver/co-driver.
      datadictionary.v1.CardSlotNumber slot = 1;
      // Populated only when slot is CARD_SLOT_NUMBER_UNRECOGNIZED.
      int32 unrecognized_slot = 2;

      // Driving status (single or crew).
      datadictionary.v1.DrivingStatus driving_status = 3;
      // Populated only when driving_status is DRIVING_STATUS_UNRECOGNIZED.
      int32 unrecognized_driving_status = 4;

      // Card status (inserted or not inserted).
      datadictionary.v1.CardStatus card_status = 5;
      // Populated only when card_status is CARD_STATUS_UNRECOGNIZED.
      int32 unrecognized_card_status = 6;

      // Driver's activity (break/rest, availability, work, driving).
      datadictionary.v1.DriverActivityValue activity = 7;
      // Populated only when activity is DRIVER_ACTIVITY_UNRECOGNIZED.
      int32 unrecognized_activity = 8;

      // Time of the change in minutes since 00:00.
      int32 time_of_change_minutes = 9;
    }

    // If true, the fields below are populated with semantic data.
    // If false, the `raw` field contains the original, unprocessed record bytes.
    bool valid = 1;

    // --- Fields for valid data (when valid = true) ---

    // See Data Dictionary, Section 2.9, `activityPreviousRecordLength`.
    int32 activity_previous_record_length = 2;

    // See Data Dictionary, Section 2.9, `activityRecordLength`.
    int32 activity_record_length = 3;

    // See Data Dictionary, Section 2.9, `activityRecordDate`.
    google.protobuf.Timestamp activity_record_date = 4;

    // See Data Dictionary, Section 2.9, `activityDailyPresenceCounter`.
    int32 activity_daily_presence_counter = 5;

    // See Data Dictionary, Section 2.9, `activityDayDistance`.
    int32 activity_day_distance = 6;

    // See Data Dictionary, Section 2.9, `activityChangeInfo`.
    repeated ActivityChange activity_change_info = 7;

    // --- Field for raw data (when valid = false) ---

    // The raw bytes of a daily record that could not be parsed.
    bytes raw = 8;
  }

  // See Data Dictionary, Section 2.17, `activityPointerOldestDayRecord`.
  int32 oldest_day_record_index = 1;

  // See Data Dictionary, Section 2.17, `activityPointerNewestRecord`.
  int32 newest_day_record_index = 2;

  // A collection of daily activity records, which may be parsed or raw.
  // See Data Dictionary, Section 2.17, `activityDailyRecords`.
  repeated DailyRecord daily_records = 3;

  // Digital signature for the EF_Driver_Activity_Data file content.
  bytes signature = 4;
}
```

## 5. Implementation Guidance

1.  **Modify the Protobuf Schema**: Apply the changes to `driver_activity_data.proto` as detailed above and regenerate the Go code from it.
2.  **Update `unmarshal_card_activity.go`**:
    *   Rename `unmarshalCardActivityData` to `unmarshalDriverActivityData` and change its signature to return `*cardv1.DriverActivityData`.
    *   In this function, read the two 2-byte pointers and store them.
    *   Pass the remaining data slice and the `newestDayRecordPointer` to a new `parseCyclicActivityDailyRecords` function.
    *   The `parseCyclicActivityDailyRecords` function should loop backwards from the `newestDayRecordPointer`, using `prevRecordLength` to find the next record and handling buffer wrap-around.
    *   In each iteration, it reads the full data for the current record and passes it to `parseSingleActivityDailyRecord`.
    *   If `parseSingleActivityDailyRecord` succeeds, it populates a `DailyRecord` with `valid = true` and the parsed fields. If it fails, it populates a `DailyRecord` with `valid = false` and the `raw` field set to the record's byte slice.
    *   The loop terminates when `prevRecordLength` is 0 or after a safeguard limit.
    *   Finally, reverse the collected slice of records.
3.  **Update `unmarshal.go`**: Update the `case` statement for `ElementaryFileType_EF_DRIVER_ACTIVITY` to call the new `unmarshalDriverActivityData` function and handle the new `DriverActivityData` message.
