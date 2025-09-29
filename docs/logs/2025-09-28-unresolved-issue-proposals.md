# Proposed Solutions for Unresolved Issues (2025-09-28)

**Objective**: This document provides analysis and actionable solutions for the unresolved issues identified in the `2025-01-27-file-structure-migration.md` log. Each proposal includes context from the relevant regulations and a clear path for implementation, adhering to the principles outlined in `AGENTS.md`.

---

## 1. ASN.1 Documentation Inconsistencies

**Issue**: ASN.1 definitions in code comments are inconsistent with the regulation specifications, leading to confusion and potential maintenance issues. This proposal provides the ground truth from the regulations to ensure future conformity.

### 1.1 `card_events.go` & `card_faults.go`

**Problem**: The documentation for `appendCardEventRecord` and `appendCardFaultRecord` incorrectly states a record size of 35 bytes, while the unmarshal functions use the correct 24-byte size.

**Regulatory Proof**:
The structure for `CardEventRecord` and `CardFaultRecord` is defined in `docs/regulation/chapters/05-tachograph-cards-file-structure.md` (Table TCS_150 for Gen1 driver cards).

**Verbatim Table Data (`TCS_150`):**

```
| File / Data element        | No of<br>Records | Size (bytes)<br>Min | Size (bytes)<br>Max | Default<br>Values |
|----------------------------|------------------|---------------------|---------------------|-------------------|
| EF Events Data             |                  | 864                 | 1728                |                   |
| └CardEventData             |                  | 864                 | 1728                |                   |
| └cardEventRecords          | 6                | 144                 | 288                 |                   |
| └ CardEventRecord          | n1               | 24                  | 24                  |                   |
| └eventType                 |                  | 1                   | 1                   | {00}              |
| └eventBeginTime            |                  | 4                   | 4                   | {00..00}          |
| └eventEndTime              |                  | 4                   | 4                   | {00..00}          |
| └eventVehicleRegistration  |                  |                     |                     |                   |
| └vehicleRegistrationNation |                  | 1                   | 1                   | {00}              |
| └vehicleRegistrationNumber |                  | 14                  | 14                  | {00, 20..2        |
```

The sum of the fields (1 + 4 + 4 + 1 + 14) confirms the total size is **24 bytes**. The same structure and size apply to `CardFaultRecord`.

**Proposed Solution**:
Update the comments in `card_events.go` and `card_faults.go` to reflect the correct 24-byte record size, referencing the regulation.

### 1.2 `card_identification.go`

**Problem**: The `unmarshal` function comment refers to a 14-byte card number, but the regulation specifies 16 bytes.

**Regulatory Proof**:
The size of the `cardNumber` field within the `CardIdentification` structure is explicitly defined in `docs/regulation/chapters/05-tachograph-cards-file-structure.md` (Table TCS_150).

**Verbatim Table Data (`TCS_150`):**

```
| File / Data element        | No of<br>Records | Size (bytes)<br>Min | Size (bytes)<br>Max | Default<br>Values |
|----------------------------|------------------|---------------------|---------------------|-------------------|
| EF Identification          |                  | 143                 | 143                 |                   |
| └CardIdentification        |                  | 65                  | 65                  |                   |
| └cardIssuingMemberState    |                  | 1                   | 1                   | {00}              |
| └cardNumber                |                  | 16                  | 16                  | {20..20}          |
| └cardIssuingAuthorityName  |                  | 36                  | 36                  | {00, 20..2        |
...
```

This table confirms the `cardNumber` field has a fixed size of **16 bytes**.

**Proposed Solution**:
Update the comment in `card_identification.go` to state the correct `cardNumber` size of **16 bytes**.

### 1.3 `card_vehicles.go`

**Problem**: The `append` function comment incorrectly documents the odometer fields as 4 bytes, while the regulation specifies 3 bytes for Gen1 cards.

**Regulatory Proof**:
The `CardVehicleRecord` for Gen1 cards uses `OdometerShort`, which is a 3-byte value. This is defined in `docs/regulation/chapters/03-data-dictionary.md` (Section 2.37) and sized in `docs/regulation/chapters/05-tachograph-cards-file-structure.md` (Table TCS_150).

**Verbatim ASN.1 Definition (`DD 2.37`):**

```
Generation 1:
CardVehicleRecord ::= SEQUENCE {
    vehicleOdometerBegin OdometerShort,
    vehicleOdometerEnd OdometerShort,
    vehicleFirstUse TimeReal,
    vehicleLastUse TimeReal,
    vehicleRegistration VehicleRegistrationIdentification,
    vuDataBlockCounter VuDataBlockCounter
}
```

**Verbatim Table Data (`TCS_150`):**

```
| File / Data element   |      | Size (bytes)<br>Min | Size (bytes)<br>Max | Default<br>Values |
|-----------------------|------|---------------------|---------------------|-------------------|
| └CardVehicleRecord    | n3   | 31                  | 31                  |                   |
| -vehicleOdometerBegin |      | 3                   | 3                   | {00..00}          |
| -vehicleOdometerEnd   |      | 3                   | 3                   | {00..00}          |
...
```

The table explicitly sizes `vehicleOdometerBegin` and `vehicleOdometerEnd` at **3 bytes** each.

**Proposed Solution**:
Update the comment in `card_vehicles.go` to reflect the correct 3-byte size for the odometer fields in Gen1 records.

---

## 2. Incomplete Append and Helper Implementations

**Issue**: Multiple `append` and `parse` helper functions across the codebase are placeholders. They lack full implementation, which prevents roundtrip data integrity (writing back the data that was read) and complete data parsing. This is a high-priority issue as it affects core functionality.

**Analysis**:
Files like `vu_activities.go`, `vu_events_faults.go`, and `card_gnss_places.go` contain `append` functions that only write a signature or return unchanged data. Additionally, helper functions responsible for parsing specific data blocks (e.g., `parseVuActivityDailyData` in `vu_activities.go`) are stubbed out, returning empty slices instead of parsed data.

**Example from `vu_activities.go`:**

```go
// Simplified implementation - would need to parse the actual structure
func parseVuCardIWData(data []byte, offset int) ([]*vuv1.Activities_CardIWRecord, int, error) {
	return []*vuv1.Activities_CardIWRecord{}, offset, nil
}

func appendVuActivities(buf *bytes.Buffer, activities *vuv1.Activities) error {
	// ... simplified version that only writes the signature data
	return nil
}
```

**Proposed Solution**:

Implement the full logic for all placeholder functions according to the EU regulations. This is a significant undertaking that should be broken down by file.

**Action Plan for `vu_activities.go`**:

1.  **Implement `appendVuActivities` for Gen1 and Gen2**:

    - The function should iterate through the fields of the `vuv1.Activities` protobuf message.
    - For each field, call a corresponding `append` helper (e.g., `appendVuCardIWData`, `appendVuActivityDailyData`) to write the data to the buffer in the correct ASN.1 format.
    - Ensure all data, including date, odometer, and all record arrays, is written—not just the signature.

2.  **Implement `parse...` Helper Functions**:

    - Complete the implementation for all `parse...` functions (`parseVuCardIWData`, `parseVuActivityDailyData`, etc.).
    - These functions must read the raw byte data, parse the specific record structures according to the regulation, and return populated protobuf message slices.
    - This is critical for the `unmarshal` functions to provide complete data to the user.

3.  **Follow `AGENTS.md` Principles**:
    - Each new `append` helper should be co-located within `vu_activities.go`.
    - The logic should be clearly documented with references to the specific ASN.1 structures being implemented.

This systematic approach should be applied to all files listed under the "Incomplete Append Implementations" and "Simplified Helper Functions" sections of the migration log. Addressing this will enable full roundtrip testing and ensure complete data processing capabilities.

---

## 3. Detailed Implementation Guide for `vu_activities.go`

**Objective**: Provide developers with the necessary regulatory context to fully implement the parsing and appending logic for `vu_activities.go`.

**Developer Instruction**: The following sections contain the verbatim ASN.1 definitions for the data structures within the VU Activities file. You must cross-reference these definitions meticulously during implementation to ensure correctness and conformity with the regulation. Do not rely on the existing simplified code.

### 3.1 Gen1 `VuActivities` Structure

**ASN.1 Definition (`DD 2.170`, `DD 2.176`, etc.)**

```
VuActivitiesFirstGen ::= SEQUENCE {
    dateOfDay                        TimeReal,
    odometerValueMidnight            OdometerValueMidnight,
    vuCardIWData                     VuCardIWData,
    vuActivityDailyData              VuActivityDailyData,
    vuPlaceDailyWorkPeriodData       VuPlaceDailyWorkPeriodData,
    vuSpecificConditionData          VuSpecificConditionData,
    signature                        SignatureFirstGen
}
```

**Implementation Details**:

- **`dateOfDay`**: Implement `readVuTimeRealFromBytes` and `appendVuTimeReal` for the 4-byte `TimeReal` type (DD 2.162).
- **`odometerValueMidnight`**: Implement `readVuOdometerFromBytes` and `appendVuOdometer` for the 3-byte `OdometerValueMidnight` type (DD 2.114).
- **`vuCardIWData`**: Implement `parseVuCardIWData` and `appendVuCardIWData`. This involves parsing a `VuCardIWRecord` (DD 2.177) which contains `cardSlotsStatus` and `vuCardRecord`.
- **`vuActivityDailyData`**: Implement `parseVuActivityDailyData` and `appendVuActivityDailyData`. This involves parsing a sequence of `ActivityChangeInfo` records (DD 2.1).
- **`vuPlaceDailyWorkPeriodData`**: Implement `parseVuPlaceDailyWorkPeriodData` and `appendVuPlaceDailyWorkPeriodData`. This involves parsing `VuPlaceDailyWorkPeriodRecord` structures (DD 2.219).
- **`vuSpecificConditionData`**: Implement `parseVuSpecificConditionData` and `appendVuSpecificConditionData`. This involves parsing a sequence of `SpecificConditionRecord`s (DD 2.152).
- **`signature`**: Ensure the 128-byte signature is correctly read and written.

### 3.2 Gen2 `VuActivities` Structure

**ASN.1 Definition (`DD 2.171`, `DD 2.178`, etc.)**

```
VuActivitiesSecondGen ::= SEQUENCE {
    dateOfDayDownloadedRecordArray           DateOfDayDownloadedRecordArray,
    odometerValueMidnightRecordArray         OdometerValueMidnightRecordArray,
    vuCardIWRecordArray                      VuCardIWRecordArray,
    vuActivityDailyRecordArray               VuActivityDailyRecordArray,
    vuPlaceDailyWorkPeriodRecordArray        VuPlaceDailyWorkPeriodRecordArray,
    vuGNSSADRecordArray                      VuGNSSADRecordArray,
    vuSpecificConditionRecordArray           VuSpecificConditionRecordArray,
    vuBorderCrossingRecordArray              VuBorderCrossingRecordArray OPTIONAL,
    vuLoadUnloadRecordArray                  VuLoadUnloadRecordArray OPTIONAL,
    signatureRecordArray                     SignatureRecordArray
}
```

**Implementation Details**:
Each of these fields is a "RecordArray" type, which typically consists of a header indicating the number and size of records, followed by the records themselves. The `parse...` and `append...` helpers for each array type must correctly handle this header + data structure.

- **`DateOfDayDownloadedRecordArray` (DD 2.59)**: Parse/append `TimeReal` records.
- **`OdometerValueMidnightRecordArray` (DD 2.115)**: Parse/append `OdometerValueMidnight` records.
- **`VuCardIWRecordArray` (DD 2.178)**: Parse/append `VuCardIWRecord` records.
- **`VuActivityDailyRecordArray` (DD 2.171)**: Parse/append `VuActivityDailyRecord` records.
- **`VuPlaceDailyWorkPeriodRecordArray` (DD 2.220)**: Parse/append `VuPlaceDailyWorkPeriodRecord` records.
- **`VuGNSSADRecordArray` (DD 2.204)**: Parse/append `VuGNSSADRecord` records (DD 2.203).
- **`VuSpecificConditionRecordArray` (DD 2.228)**: Parse/append `VuSpecificConditionRecord` records.
- **`VuBorderCrossingRecordArray` (DD 2.203b, Gen2v2)**: Parse/append `VuBorderCrossingRecord` records (DD 2.203a).
- **`VuLoadUnloadRecordArray` (DD 2.208b, Gen2v2)**: Parse/append `VuLoadUnloadRecord` records (DD 2.208a).
- **`SignatureRecordArray` (DD 2.150)**: Parse/append the signature.

By providing this detailed context, developers will have a clear and authoritative reference to complete the required implementations correctly.

---

## 4. Incomplete Gen2 Implementation in `vu_overview.go`

**Issue**: The functions `unmarshalOverviewGen2` and `appendOverviewGen2` in `vu_overview.go` are placeholders and do not implement the logic for handling Generation 2 VU overview data.

**Analysis**: The current implementation returns empty data for Gen2, preventing the parsing and creation of Gen2 overview files. The Gen2 structure is significantly different from Gen1, using a series of record arrays instead of fixed-size fields.

**Developer Instruction**: The developer must implement the parsing and appending logic based on the following ASN.1 definition for the `VuOverviewSecondGen` data type. Pay close attention to the `RecordArray` structures, which require parsing a header to determine the number and size of subsequent records.

**Verbatim ASN.1 Definition**:

```
VuOverviewSecondGen ::= SEQUENCE {
    memberStateCertificateRecordArray    MemberStateCertificateRecordArray,
    vuCertificateRecordArray             VuCertificateRecordArray,
    vehicleIdentificationNumberRecordArray VehicleIdentificationNumberRecordArray,
    vehicleRegistrationIdentificationRecordArray VehicleRegistrationIdentificationRecordArray,
    currentDateTimeRecordArray           CurrentDateTimeRecordArray,
    vuDownloadablePeriodRecordArray      VuDownloadablePeriodRecordArray,
    cardSlotsStatusRecordArray           CardSlotsStatusRecordArray,
    vuDownloadActivityDataRecordArray    VuDownloadActivityDataRecordArray,
    vuCompanyLocksRecordArray            VuCompanyLocksRecordArray,
    vuControlActivityRecordArray         VuControlActivityRecordArray,
    signatureRecordArray                 SignatureRecordArray
}
```

**Action Plan for `vu_overview.go`**:

1.  **Implement `unmarshalOverviewGen2`**:

    - Sequentially parse each `RecordArray` field from the byte data.
    - For each `...RecordArray` (e.g., `MemberStateCertificateRecordArray`), create a helper function (e.g., `parseMemberStateCertificateRecordArray`).
    - This helper must first parse the `RecordArray` header to get the count and size of the records.
    - Then, it must loop that many times, parsing each individual record and appending it to a slice in the target `vuv1.Overview` protobuf message.

2.  **Implement `appendOverviewGen2`**:

    - Sequentially append each field from the `vuv1.Overview` protobuf message to the buffer.
    - For each `...RecordArray`, create an `append...RecordArray` helper.
    - This helper must first write the `RecordArray` header (calculating the total size and record count from the protobuf data).
    - Then, it must loop through the slice in the protobuf message, appending each record to the buffer.

3.  **Cross-Reference Documentation**: The definitions for each `RecordArray` and the underlying record types (e.g., `VuCompanyLocksRecord`, `VuControlActivityRecord`) are specified in the **Data Dictionary (`docs/regulation/chapters/03-data-dictionary.md`)**. The developer must consult this document for the precise structure and size of each field within these records.

---

## 5. Final Cleanup: Naming Conventions and CardNumber Logic

**Objective**: Address the remaining low and medium-priority issues to ensure the codebase is fully aligned with project standards.

### 5.1 Constants Naming Convention

**Issue**: Inconsistent constant naming (e.g., `cardDriverActivityHeaderSize` vs. `lenCardDriverActivityHeader`).

**Proposed Solution**:

- **Standardize All Constants**: Globally search for and refactor all size and length-related constants to adhere to the `AGENTS.md` guidelines.
- **Use `len` prefix for byte lengths** (e.g., `lenCardDriverActivityHeader`).
- **Use `idx` prefix for byte offsets** (e.g., `idxCardHolderName`).
- This is a low-priority, but high-impact task for code readability and consistency. It can be done as a single, sweeping change across the codebase.

### 5.2 Incomplete Card Number Handling in `card_identification.go`

**Issue**: The `appendCardIdentification` function has placeholder logic for serializing the `CardNumber` CHOICE type.

**Analysis**: The current code attempts to concatenate fields but does not correctly handle the two different structures within the `CardNumber` CHOICE type (one for drivers, one for other card types).

**Developer Instruction**: The developer must implement the logic to correctly serialize the `CardNumber` based on which field in the `oneof` is populated in the protobuf message.

**Verbatim ASN.1 Definition (`DD 2.26`):**

```
CardNumber ::= CHOICE {
    -- Driver Card
    SEQUENCE {
        driverIdentification    IA5String(SIZE(14)),
        cardReplacementIndex    CardReplacementIndex, -- 1 byte
        cardRenewalIndex        CardRenewalIndex,     -- 1 byte
    },
    -- Other Cards (Workshop, Control, Company)
    SEQUENCE {
        ownerIdentification     IA5String(SIZE(13)),
        cardConsecutiveIndex    CardConsecutiveIndex, -- 1 byte
        cardReplacementIndex    CardReplacementIndex, -- 1 byte
        cardRenewalIndex        CardRenewalIndex      -- 1 byte
    }
}
```

_Note: The total size is always 16 bytes._

**Action Plan**:

1.  **Refactor `appendCardIdentification`**:
    - Check which of `DriverIdentification` or `OwnerIdentification` is set in the `cardId` protobuf.
    - **If `DriverIdentification` is set**: Append the 14-byte `identificationNumber`, 1-byte `replacementIndex`, and 1-byte `renewalIndex` to a 16-byte buffer.
    - **If `OwnerIdentification` is set**: Append the 13-byte `identificationNumber`, 1-byte `consecutiveIndex`, 1-byte `replacementIndex`, and 1-byte `renewalIndex` to a 16-byte buffer.
    - Ensure proper padding and field order as specified in the regulation.
2.  **Verify `unmarshalIdentification`**: Review the unmarshalling logic to ensure it correctly parses the 16-byte `CardNumber` field based on the card type (which may need to be inferred or passed in), correctly populating either the `DriverIdentification` or `OwnerIdentification` protobuf message.

---

## 6. Protobuf Schema Limitations Discovered During Implementation

**Issue**: During the implementation of the CardNumber handling logic, we discovered that the current protobuf schema is incomplete and doesn't fully represent the ASN.1 specification.

### 6.1 CardNumber CHOICE Type Incomplete Implementation

**Problem**: The `CardNumber` CHOICE type in the current protobuf schema doesn't include all the fields specified in the ASN.1 definition from the regulation.

**ASN.1 Definition (DD 2.26)**:

```
CardNumber ::= CHOICE {
    -- Driver Card
    SEQUENCE {
        driverIdentification    IA5String(SIZE(14)),
        cardReplacementIndex    CardReplacementIndex, -- 1 byte
        cardRenewalIndex        CardRenewalIndex,     -- 1 byte
    },
    -- Other Cards (Workshop, Control, Company)
    SEQUENCE {
        ownerIdentification     IA5String(SIZE(13)),
        cardConsecutiveIndex    CardConsecutiveIndex, -- 1 byte
        cardReplacementIndex    CardReplacementIndex, -- 1 byte
        cardRenewalIndex        CardRenewalIndex      -- 1 byte
    }
}
```

**Current Protobuf Schema Issues**:

1. **DriverIdentification Missing Fields**: The `DriverIdentification` message only contains the `identificationNumber` field but is missing:

   - `cardReplacementIndex` (1 byte)
   - `cardRenewalIndex` (1 byte)

2. **OwnerIdentification Complete**: The `OwnerIdentification` message correctly includes all required fields:
   - `identificationNumber` (13 bytes)
   - `consecutiveIndex` (1 byte)
   - `replacementIndex` (1 byte)
   - `renewalIndex` (1 byte)

**Impact**: This limitation prevents proper roundtrip parsing and serialization of driver card numbers, as the replacement and renewal indices are lost during the unmarshal process and cannot be restored during marshaling.

**Temporary Workaround**: The current implementation assumes driver card format and only parses the 14-byte identification number, with a comment noting the schema limitation.

**Required Schema Updates**:

1. Add `cardReplacementIndex` and `cardRenewalIndex` fields to the `DriverIdentification` message
2. Ensure the field types match the ASN.1 specification (1-byte values)
3. Update the generated Go code accordingly

**Files Affected**:

- `proto/wayplatform/connect/tachograph/dd/v1/driver_identification.proto`
- Generated Go files in `proto/gen/go/wayplatform/connect/tachograph/dd/v1/`
- `card_identification.go` (unmarshal and append functions)

---

## 7. Protobuf API Improvements for Better Go Idioms

**Issue**: During the implementation of CardNumber handling logic, we discovered that the current protobuf-generated Go code could be more idiomatic and easier to work with.

### 7.1 Add GetValue Method to StringValue

**Problem**: The `StringValue` type currently only provides `GetDecoded()` and `GetEncoded()` methods, but many use cases require direct access to the underlying string value without encoding concerns.

**Current Usage**:

```go
// Current approach - verbose and not idiomatic
consecutiveStr := consecutive.GetDecoded()
if len(consecutiveStr) > 0 {
    cardNumberBytes[13] = consecutiveStr[0]
}

// Or for numeric values stored as strings
replacementStr := replacement.GetDecoded()
if len(replacementStr) > 0 {
    cardNumberBytes[14] = replacementStr[0]
}
```

**Proposed Solution**:
Add a `GetValue()` method to `StringValue` that returns the decoded string directly, making the code more idiomatic:

```go
// Proposed approach - more idiomatic
if consecutive.GetValue() != "" {
    cardNumberBytes[13] = consecutive.GetValue()[0]
}

// Or for numeric values
if replacement.GetValue() != "" {
    cardNumberBytes[14] = replacement.GetValue()[0]
}
```

**Benefits**:

1. **Idiomatic Go**: Follows Go conventions where `GetValue()` is a common pattern for accessing the primary value of a type
2. **Reduced Verbosity**: Eliminates the need to call `GetDecoded()` and check length in many cases
3. **Consistency**: Aligns with other protobuf-generated types that provide `GetValue()` methods
4. **Better Error Handling**: Could potentially return an error for invalid states

**Implementation Details**:

- Add `GetValue() string` method to `StringValue` struct
- The method should return the decoded string value directly
- Consider adding `GetValueOrEmpty() string` for cases where nil safety is important
- Update the protobuf generator configuration if needed

**Files Affected**:

- `proto/wayplatform/connect/tachograph/dd/v1/string_value.proto`
- Generated Go files in `proto/gen/go/wayplatform/connect/tachograph/dd/v1/`
- All files that use `StringValue` for simple value access

**Priority**: Medium - This is a quality-of-life improvement that would make the codebase more maintainable and idiomatic.

---
