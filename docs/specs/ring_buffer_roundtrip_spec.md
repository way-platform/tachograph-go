# Specification: Ensuring Roundtrip Integrity for Ring Buffer Data

## 1. Introduction

This document specifies the strategy for correctly unmarshalling and marshalling ring buffer data structures within tachograph files, specifically focusing on the `EF_Driver_Activity` file. The primary goal is to ensure that a file can be unmarshalled and then marshalled back to a byte-for-byte identical binary, a process known as a "roundtrip".

The key challenge this specification addresses is the handling of empty, zeroed-out, or placeholder records within the ring buffer, which are currently dropped during the unmarshalling process, leading to failed roundtrips.

## 2. Problem Analysis

Certain data structures in tachograph cards, most notably `EF_Driver_Activity`, are implemented as a cyclic ring buffer. This buffer contains a series of daily records, where each record contains a header with pointers (`activityPreviousRecordLength`, `activityRecordLength`) that link it to the previous entry.

The current implementation exhibits the following issue:

1.  **Premature Termination**: The unmarshalling logic in `unmarshal_card_activity.go` iterates backwards through the linked list of daily records. It interprets a record with zeroed-out content or a zero-length header as the end of the valid data chain and stops parsing.
2.  **Incomplete Data Model**: This causes any "empty" records, which act as placeholders in the fixed-size ring buffer, to be discarded. The resulting in-memory protobuf message (e.g., `DriverActivity`) contains a list of daily records that is shorter than the actual number of slots (both used and empty) in the original file.
3.  **Roundtrip Failure**: When this incomplete data structure is marshalled back into binary, the empty records are missing. The resulting file is shorter than the original and not byte-identical, causing the roundtrip tests in `roundtrip_test.go` to fail.

## 3. Proposed Solution

To guarantee a perfect roundtrip, the unmarshalling and marshalling logic must be made symmetrical. The unmarshaller must recognize and preserve empty records, and the marshaller must be able to write them back to their original binary form.

### 3.1. Definition of an "Empty Record"

An "empty record" is not a record of zero length. It is a record that has a valid length in its header but whose content fields (e.g., date, distance, activity changes) are filled with zero or placeholder values.

The exact byte pattern that constitutes an "empty record" must be determined by empirical analysis of real-world `.DDD` files that contain such records. This is a prerequisite for implementation.

### 3.2. Unmarshalling Logic

The unmarshalling process will be enhanced to preserve the integrity of the buffer structure.

-   The parsing loop in `parseActivityDailyRecords` will be modified. When it encounters a byte pattern matching the defined "empty record", it will no longer terminate.
-   Instead, it will proceed to parse this record in `parseActivityDailyRecord`. This function will be adjusted to correctly handle the empty content without returning an error.
-   It will produce a valid `cardv1.DriverActivity_DailyRecord` protobuf message, but one where all fields have their default values (e.g., zero for integers, empty for lists).
-   This "empty" message will be appended to the list of daily records. The final list will therefore contain the exact same number of records as the original file's ring buffer.

### 3.3. Marshalling Logic

The marshalling process will be updated to symmetrically handle the "empty" records now present in the in-memory data structure.

-   The `AppendActivityDailyRecord` function will be modified to detect when it is given an "empty" `DailyRecord` message. This detection will be based on its fields having default values (e.g., `rec.GetActivityDayDistance() == 0 && len(rec.GetActivityChangeInfo()) == 0`).
-   When an empty record is detected, the function will write the exact, predefined byte pattern for an empty record (as determined in section 3.1) to the output buffer.
-   If the record is not empty, the existing logic for serializing a data-filled record will be used.

This symmetry ensures that the presence and position of empty records are preserved through the unmarshal-marshal cycle.

## 4. Implementation Guidance

The following steps will be taken to implement this specification:

1.  **Analysis**: Obtain a sample `.DDD` file that fails the roundtrip test due to this issue. Using a hex analysis tool, inspect the `EF_Driver_Activity` data block (FID `0x0504`) and document the precise byte pattern of an "empty" daily record, including its header.
2.  **Unmarshaller Modification**: Update the logic in `unmarshal_card_activity.go`.
    -   Modify `parseActivityDailyRecords` and `parseActivityDailyRecord` to correctly identify the "empty record" pattern from Step 1 and generate an empty `DriverActivity_DailyRecord` message in its place in the records slice.
3.  **Marshaller Modification**: Update the logic in `append_card_activity.go`.
    -   Modify `AppendActivityDailyRecord` to detect an empty `DriverActivity_DailyRecord` message and write the exact byte pattern identified in Step 1.
4.  **Verification**: Add the failing `.DDD` file to the `testdata/card/` directory. Run the test suite and confirm that `TestRoundtripCard` now passes for this file, verifying that the output is byte-for-byte identical to the input.
