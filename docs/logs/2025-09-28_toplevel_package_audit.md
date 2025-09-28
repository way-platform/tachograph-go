# 2025-09-28: Top-Level Package Audit Log

This log is a list of tasks to improve the alignment of the top-level Go package with the guidelines in `AGENTS.md`.

*   **[COMPLETE]** This log contains the full list of tasks from the initial audit.

---
## Helper File Refactoring

These tasks focus on eliminating generic helper files (`*_helpers.go`, `enum_helpers.go`) by refactoring their functionality into more appropriate, semantically-named files or co-locating them with their usage, as per `AGENTS.md`.

### Task: Eliminate `enum_helpers.go` via Protobuf Reflection

*   **File:** `enum_helpers.go`
*   **Problem:** This file contains dozens of hardcoded functions (`SetEventFaultType`, `GetEquipmentTypeProtocolValue`, etc.) to convert between protobuf enums and their raw protocol values. This is brittle and directly contradicts the `AGENTS.md` guideline to use protobuf reflection.
*   **Action:**
    1.  Create a new file, `protobuf_helpers.go`.
    2.  In this new file, implement two generic, reflection-based functions as proposed in `AGENTS.md`:
        *   `getProtocolValueFromEnum(enumValue protoreflect.Enum) (int32, bool)` for marshalling.
        *   `setEnumFromProtocolValue(enumDesc protoreflect.EnumDescriptor, rawValue int32) (protoreflect.EnumNumber, bool)` for unmarshalling.
    3.  Systematically replace all calls to the legacy functions in `enum_helpers.go` across the entire codebase with calls to these new, generic reflection helpers.
    4.  Delete the `enum_helpers.go` file.

### Task: Refactor Marshalling and Unmarshalling Helpers

*   **Files:** `append_helpers.go`, `append_vu_helpers.go`, `unmarshal_card_helpers.go`, `unmarshal_vu_helpers.go`
*   **Problem:** These files contain generic helper functions (e.g., `appendTimeReal`, `readTimeReal`, `bcdBytesToInt`, `readString`) that should be organized into semantically-named files.
*   **Action:**
    1.  Create new, semantically-named files for common data types. For example:
        *   `unmarshal_time.go` and `append_time.go` for time-related functions (`readTimeReal`, `readDatef`, `appendTimeReal`, `appendDatef`).
        *   `unmarshal_string.go` and `append_string.go` for string and BCD-related functions (`readString`, `bcdBytesToInt`, `appendString`, `appendBCD`).
    2.  Move the relevant helper functions from the generic `*_helpers.go` files into these new, more specific files.
    3.  Update all call sites throughout the codebase to use the functions from their new locations.
    4.  Delete the `append_helpers.go`, `append_vu_helpers.go`, `unmarshal_card_helpers.go`, and `unmarshal_vu_helpers.go` files.

---

## Card File Documentation and Layout (`unmarshal_card_*`)

The following tasks are to bring individual parsing functions into alignment with `AGENTS.md` by adding ASN.1 specification comments and using `const` blocks for binary layouts.

### Task: Document `unmarshalIcc` in `unmarshal_card_icc.go`

*   **File:** `unmarshal_card_icc.go`
*   **Function:** `unmarshalIcc`
*   **Problem:** The function lacks the ASN.1 definition in its comments and uses magic numbers for parsing.
*   **Action:**
    1.  Add the following ASN.1 documentation to the function comment.
    2.  Introduce a `const` block to define the offsets and lengths for the fields.
    3.  Refactor the parsing logic to use these constants.
*   **ASN.1 Specification (Data Dictionary 2.23):**
    ```asn.1
    // CardIccIdentification ::= SEQUENCE {
    //     clockStop                   OCTET STRING (SIZE(1)),
    //     cardExtendedSerialNumber    ExtendedSerialNumber,    -- 8 bytes
    //     cardApprovalNumber          CardApprovalNumber,      -- 8 bytes
    //     cardPersonaliserID          ManufacturerCode,        -- 1 byte
    //     embedderIcAssemblerId       EmbedderIcAssemblerId,   -- 5 bytes
    //     icIdentifier                OCTET STRING (SIZE(2))
    // }
    ```

### Task: Document `unmarshalCardIc` in `unmarshal_card_ic.go`

*   **File:** `unmarshal_card_ic.go`
*   **Function:** `unmarshalCardIc`
*   **Problem:** The function lacks the ASN.1 definition and layout constants.
*   **Action:**
    1.  Add the ASN.1 documentation to the function comment.
    2.  Introduce a `const` block for the layout.
    3.  Refactor the parsing logic to use these constants.
*   **ASN.1 Specification (Data Dictionary 2.13):**
    ```asn.1
    // CardChipIdentification ::= SEQUENCE {
    //     icSerialNumber              OCTET STRING (SIZE(4)),
    //     icManufacturingReferences   OCTET STRING (SIZE(4))
    // }
    ```

### Task: Document `unmarshalIdentification` in `unmarshal_card_identification.go`

*   **File:** `unmarshal_card_identification.go`
*   **Function:** `unmarshalIdentification`
*   **Problem:** The function parses a large, composite file structure (`EF_Identification`) without any documentation on the layout or use of constants.
*   **Action:**
    1.  Add a detailed comment explaining the 143-byte structure of the `EF_Identification` file, which contains both `CardIdentification` and `DriverCardHolderIdentification`.
    2.  Introduce a `const` block defining the offsets and lengths for all parsed fields (e.g., `idxCardIssuingMemberState`, `lenCardIssuingAuthorityName`, `idxCardHolderSurname`, etc.).
    3.  Refactor the parsing logic to use these constants instead of magic numbers.
*   **File Structure Context (Data Dictionary 2.24 & 2.62):**
    The `EF_Identification` file is a concatenation of `CardIdentification` (approx. 65-66 bytes) and `DriverCardHolderIdentification` (78 bytes). The documentation should clarify this structure.

### Task: Document `unmarshalDrivingLicenceInfo` in `unmarshal_card_driving_licence.go`

*   **File:** `unmarshal_card_driving_licence.go`
*   **Function:** `unmarshalDrivingLicenceInfo`
*   **Problem:** The function lacks the ASN.1 definition and layout constants.
*   **Action:**
    1.  Add the ASN.1 documentation to the function comment.
    2.  Introduce a `const` block for the layout.
    3.  Refactor the parsing logic to use these constants.
*   **ASN.1 Specification (Data Dictionary 2.18):**
    ```asn.1
    // CardDrivingLicenceInformation ::= SEQUENCE {
    //     drivingLicenceIssuingAuthority  Name,          -- 36 bytes
    //     drivingLicenceIssuingNation     NationNumeric, -- 1 byte
    //     drivingLicenceNumber            IA5String(SIZE(16))
    // }
    ```

### Task: Document `unmarshalCardApplicationIdentification`

*   **File:** `unmarshal_card_application_identification.go`
*   **Function:** `unmarshalCardApplicationIdentification`
*   **Problem:** Lacks ASN.1 definition and layout constants.
*   **Action:** Add ASN.1 documentation and a `const` block for the layout of `DriverCardApplicationIdentification`.
*   **ASN.1 Specification (Data Dictionary 2.61):**
    ```asn.1
    // DriverCardApplicationIdentification ::= SEQUENCE {
    //     typeOfTachographCardId    EquipmentType,          -- 1 byte
    //     cardStructureVersion      CardStructureVersion,   -- 2 bytes
    //     noOfEventsPerType         NoOfEventsPerType,      -- 1 byte
    //     noOfFaultsPerType         NoOfFaultsPerType,      -- 1 byte
    //     activityStructureLength   CardActivityLengthRange,-- 2 bytes
    //     noOfCardVehicleRecords    NoOfCardVehicleRecords, -- 1 byte
    //     noOfCardPlaceRecords      NoOfCardPlaceRecords    -- 1 byte
    //     -- Gen2 additions follow
    // }
    ```

### Task: Document `unmarshalCardControlActivityData`

*   **File:** `unmarshal_card_control_activity.go`
*   **Function:** `unmarshalCardControlActivityData`
*   **Problem:** Lacks ASN.1 definition and layout constants.
*   **Action:** Add ASN.1 documentation and a `const` block for `CardControlActivityDataRecord`.
*   **ASN.1 Specification (Data Dictionary 2.15):**
    ```asn.1
    // CardControlActivityDataRecord ::= SEQUENCE {
    //     controlType                 ControlType,                      -- 1 byte
    //     controlTime                 TimeReal,                         -- 4 bytes
    //     controlCardNumber           FullCardNumber,                   -- 18 bytes
    //     controlVehicleRegistration  VehicleRegistrationIdentification,-- 15 bytes
    //     controlDownloadPeriodBegin  TimeReal,                         -- 4 bytes
    //     controlDownloadPeriodEnd    TimeReal                          -- 4 bytes
    // }
    ```

### Task: Document `unmarshalCardCurrentUsage`

*   **File:** `unmarshal_card_current_usage.go`
*   **Function:** `unmarshalCardCurrentUsage`
*   **Problem:** Lacks ASN.1 definition and layout constants.
*   **Action:** Add ASN.1 documentation and a `const` block for `CardCurrentUse`.
*   **ASN.1 Specification (Data Dictionary 2.16):**
    ```asn.1
    // CardCurrentUse ::= SEQUENCE {
    //     sessionOpenTime     TimeReal,                         -- 4 bytes
    //     sessionOpenVehicle  VehicleRegistrationIdentification -- 15 bytes
    // }
    ```
### Task: Document `unmarshalEventRecord` in `unmarshal_card_events.go`

*   **File:** `unmarshal_card_events.go`
*   **Function:** `unmarshalEventRecord`
*   **Problem:** The function parses a 24-byte event record without the ASN.1 definition in its comments or a `const` block for the layout.
*   **Action:**
    1.  Add the ASN.1 documentation to the function comment.
    2.  Introduce a `const` block for the layout.
    3.  Refactor the parsing logic to use these constants.
*   **ASN.1 Specification (Data Dictionary 2.20):**
    ```asn.1
    // CardEventRecord ::= SEQUENCE {
    //     eventType                   EventFaultType,                     -- 1 byte
    //     eventBeginTime              TimeReal,                         -- 4 bytes
    //     eventEndTime                TimeReal,                         -- 4 bytes
    //     eventVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
    // }
    ```

### Task: Document `UnmarshalFaultRecord` in `unmarshal_card_faults.go`

*   **File:** `unmarshal_card_faults.go`
*   **Function:** `UnmarshalFaultRecord`
*   **Problem:** The function parses a 24-byte fault record without the ASN.1 definition in its comments or a `const` block for the layout.
*   **Action:**
    1.  Add the ASN.1 documentation to the function comment.
    2.  Introduce a `const` block for the layout.
    3.  Refactor the parsing logic to use these constants.
*   **ASN.1 Specification (Data Dictionary 2.22):**
    ```asn.1
    // CardFaultRecord ::= SEQUENCE {
    //     faultType                   EventFaultType,                     -- 1 byte
    //     faultBeginTime              TimeReal,                         -- 4 bytes
    //     faultEndTime                TimeReal,                         -- 4 bytes
    //     faultVehicleRegistration    VehicleRegistrationIdentification -- 15 bytes
    // }
    ```

### Task: Document `parsePlaceRecord` in `unmarshal_card_places.go`

*   **File:** `unmarshal_card_places.go`
*   **Function:** `parsePlaceRecord`
*   **Problem:** The function parses a 12-byte record that deviates slightly from the specification without clear documentation.
*   **Action:**
    1.  Add the ASN.1 documentation to the function comment.
    2.  The comment should explicitly note that the implementation reads a 12-byte record, including a 2-byte region and a 1-byte reserved field, which differs from the 10-byte structure implied by a strict reading of the `PlaceRecord` and `RegionNumeric` specs.
    3.  Introduce a `const` block for the 12-byte layout being parsed.
*   **ASN.1 Specification (Data Dictionary 2.117):**
    ```asn.1
    // PlaceRecord ::= SEQUENCE {
    //     entryTime                   TimeReal,                   -- 4 bytes
    //     entryTypeDailyWorkPeriod    EntryTypeDailyWorkPeriod,   -- 1 byte
    //     dailyWorkPeriodCountry      NationNumeric,              -- 1 byte
    //     dailyWorkPeriodRegion       RegionNumeric,              -- 1 byte
    //     vehicleOdometerValue        OdometerShort               -- 3 bytes
    // }
    ```

### Task: Document `parseSpecificConditionRecord` in `unmarshal_card_specific_conditions.go`

*   **File:** `unmarshal_card_specific_conditions.go`
*   **Function:** `parseSpecificConditionRecord`
*   **Problem:** The function lacks the ASN.1 definition and layout constants.
*   **Action:**
    1.  Add the ASN.1 documentation to the function comment.
    2.  Introduce a `const` block for the 5-byte layout.
*   **ASN.1 Specification (Data Dictionary 2.152):**
    ```asn.1
    // SpecificConditionRecord ::= SEQUENCE {
    //     entryTime                   TimeReal,               -- 4 bytes
    //     specificConditionType       SpecificConditionType   -- 1 byte
    // }
    ```

### Task: Document `parseVehicleRecord` in `unmarshal_card_vehicles.go`

*   **File:** `unmarshal_card_vehicles.go`
*   **Function:** `parseVehicleRecord`
*   **Problem:** The function handles both Gen1 and Gen2 record formats but lacks clear documentation and constants for their layouts.
*   **Action:**
    1.  Add the ASN.1 documentation for both Gen1 (31 bytes) and Gen2 (48 bytes) versions of `CardVehicleRecord`.
    2.  Introduce `const` blocks for both layouts.
    3.  Ensure the comments clarify the logic for determining the generation and parsing accordingly.
*   **ASN.1 Specification (Data Dictionary 2.37):**
    ```asn.1
    // CardVehicleRecord (Generation 1) ::= SEQUENCE {
    //     vehicleOdometerBegin        OdometerShort,                      -- 3 bytes
    //     vehicleOdometerEnd          OdometerShort,                      -- 3 bytes
    //     vehicleFirstUse             TimeReal,                         -- 4 bytes
    //     vehicleLastUse              TimeReal,                         -- 4 bytes
    //     vehicleRegistration         VehicleRegistrationIdentification,  -- 15 bytes
    //     vuDataBlockCounter          VuDataBlockCounter                  -- 2 bytes
    // }
    //
    // CardVehicleRecord (Generation 2) ::= SEQUENCE {
    //     ... (same as Gen1) ...
    //     vehicleIdentificationNumber VehicleIdentificationNumber         -- 17 bytes
    // }
    ```

### Task: Document `parseSingleActivityDailyRecord` in `unmarshal_card_activity.go`

*   **File:** `unmarshal_card_activity.go`
*   **Function:** `parseSingleActivityDailyRecord`
*   **Problem:** This function parses a complex, variable-length structure without documentation for its fixed-header part.
*   **Action:**
    1.  Add the ASN.1 documentation for `CardActivityDailyRecord`.
    2.  Introduce a `const` block for the fixed-size header of the record (12 bytes).
    3.  The function comment should clarify that the `activityChangeInfo` section is variable in length.
*   **ASN.1 Specification (Data Dictionary 2.9):**
    ```asn.1
    // CardActivityDailyRecord ::= SEQUENCE {
    //     activityPreviousRecordLength   INTEGER(0..CardActivityLengthRange), -- 2 bytes
    //     activityRecordLength           INTEGER(0..CardActivityLengthRange), -- 2 bytes
    //     activityRecordDate             TimeReal,                            -- 4 bytes
    //     activityDailyPresenceCounter   DailyPresenceCounter,                -- 2 bytes
    //     activityDayDistance            Distance,                            -- 2 bytes
    //     activityChangeInfo             SET SIZE (1..1440) OF ActivityChangeInfo -- 2 * n bytes
    // }
    ```

---

## Card File Documentation (`append_card_*`)

The following tasks are to bring the marshalling functions into alignment with `AGENTS.md` by adding ASN.1 specification comments and using `const` blocks for binary layouts.

### Task: Document `AppendIcc` in `append_card_icc.go`

*   **File:** `append_card_icc.go`
*   **Function:** `AppendIcc`
*   **Problem:** The function serializes a 25-byte structure without ASN.1 documentation or layout constants.
*   **Action:** Add the ASN.1 specification for `CardIccIdentification` (DD 2.23) to the function comment and add a `const` block for the layout.

### Task: Document `AppendCardIc` in `append_card_ic.go`

*   **File:** `append_card_ic.go`
*   **Function:** `AppendCardIc`
*   **Problem:** The function serializes an 8-byte structure without ASN.1 documentation or layout constants.
*   **Action:** Add the ASN.1 specification for `CardChipIdentification` (DD 2.13) to the function comment and add a `const` block for the layout.

### Task: Document `AppendCardIdentification` and `AppendDriverCardHolderIdentification`

*   **File:** `append_card_identification.go`
*   **Problem:** These functions serialize the `EF_Identification` file structure without documentation or layout constants.
*   **Action:** Add comments to both functions explaining the 143-byte file structure and introduce `const` blocks for the field layouts, mirroring the unmarshalling task.

### Task: Document `AppendDrivingLicenceInfo` in `append_card_driving_licence.go`

*   **File:** `append_card_driving_licence.go`
*   **Function:** `AppendDrivingLicenceInfo`
*   **Problem:** The function serializes a 53-byte structure without ASN.1 documentation or layout constants.
*   **Action:** Add the ASN.1 specification for `CardDrivingLicenceInformation` (DD 2.18) to the function comment and add a `const` block for the layout.

### Task: Document `AppendCardApplicationIdentification`

*   **File:** `append_card_application_identification.go`
*   **Function:** `AppendCardApplicationIdentification`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `DriverCardApplicationIdentification` (DD 2.61) to the function comment and add a `const` block for the layout.

### Task: Document `AppendCardControlActivityData` in `append_card_control_activity.go`

*   **File:** `append_card_control_activity.go`
*   **Function:** `AppendCardControlActivityData`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `CardControlActivityDataRecord` (DD 2.15) to the function comment and add a `const` block for the layout.

### Task: Document `AppendCurrentUsage` in `append_card_misc.go`

*   **File:** `append_card_misc.go`
*   **Function:** `AppendCurrentUsage`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `CardCurrentUse` (DD 2.16) to the function comment and add a `const` block for the layout.

### Task: Document `AppendEventRecord` in `append_card_events.go`

*   **File:** `append_card_events.go`
*   **Function:** `AppendEventRecord`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `CardEventRecord` (DD 2.20) to the function comment and add a `const` block for the layout.

### Task: Document `AppendFaultRecord` in `append_card_faults.go`

*   **File:** `append_card_faults.go`
*   **Function:** `AppendFaultRecord`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `CardFaultRecord` (DD 2.22) to the function comment and add a `const` block for the layout.

### Task: Document `AppendPlaceRecord` in `append_card_places.go`

*   **File:** `append_card_places.go`
*   **Function:** `AppendPlaceRecord`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `PlaceRecord` (DD 2.117) to the function comment, note the implementation deviation, and add a `const` block for the layout.

### Task: Document `AppendSpecificConditionRecord` in `append_card_misc.go`

*   **File:** `append_card_misc.go`
*   **Function:** `AppendSpecificConditionRecord`
*   **Problem:** The function lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `SpecificConditionRecord` (DD 2.152) to the function comment and add a `const` block for the layout.

### Task: Document `AppendVehicleRecord` in `append_card_vehicles.go`

*   **File:** `append_card_vehicles.go`
*   **Function:** `AppendVehicleRecord`
*   **Problem:** The function lacks ASN.1 documentation and layout constants for the Gen1/Gen2 formats.
*   **Action:** Add the ASN.1 specification for `CardVehicleRecord` (DD 2.37) to the function comment and add `const` blocks for the layouts.

### Task: Document `appendParsedDailyRecord` in `append_card_activity.go`

*   **File:** `append_card_activity.go`
*   **Function:** `appendParsedDailyRecord`
*   **Problem:** The function serializes a variable-length structure without documentation for its fixed-header part.
*   **Action:** Add the ASN.1 specification for `CardActivityDailyRecord` (DD 2.9) to the function comment and add a `const` block for the fixed-header portion.

---

## VU File Implementation and Documentation

The following tasks address gaps in the implementation and documentation of the Vehicle Unit (VU) data processing files (`unmarshal_vu_*` and `append_vu_*`). Many of these files are currently stubs and require full implementation to meet the project goal of full roundtrip parsing.

### Task: Document and Enhance `unmarshal_vu_overview.go`

*   **File:** `unmarshal_vu_overview.go`
*   **Function:** `unmarshalOverviewGen1`
*   **Problem:** The function has a partial implementation for parsing a Gen1 VU overview but lacks documentation, uses magic numbers, and skips complex parts like `VuCompanyLocksData`.
*   **Action:**
    1.  Add a function comment explaining the known layout of the Gen1 overview structure. Since this is a composite download structure, a direct ASN.1 mapping may not exist; document the field order, sizes, and offsets based on the specification for VU data downloads.
    2.  Create a `const` block for all known field lengths and offsets (e.g., `lenMemberStateCertificate = 194`, `lenVuCertificate = 194`, `lenVehicleIdentificationNumber = 17`).
    3.  Refactor the parsing logic to use these constants.
    4.  Create sub-tasks to fully implement the parsing for skipped sections like `VuCompanyLocksData` and `VuControlActivityData`.

### Task: Document and Implement `unmarshal_vu_activities.go`

*   **File:** `unmarshal_vu_activities.go`
*   **Function:** `unmarshalVuActivitiesGen1`, `unmarshalVuActivitiesGen2`
*   **Problem:** The file contains stubs for parsing VU activity data. The helper functions (`parseVuCardIWData`, `parseVuActivityDailyData`, etc.) are empty.
*   **Action:**
    1.  Fully implement the parsing logic for `VuActivityDailyData` (DD 2.170), `VuCardIWData` (DD 2.176), `VuPlaceDailyWorkPeriodData` (DD 2.218), and `VuSpecificConditionData` (DD 2.227) according to the regulation.
    2.  For each parsing function, add the corresponding ASN.1 documentation and layout constants as required by `AGENTS.md`.

### Task: Document `UnmarshalDownloadInterfaceVersion`

*   **File:** `unmarshal_vu_download_interface_version.go`
*   **Function:** `UnmarshalDownloadInterfaceVersion`
*   **Problem:** The function is implemented but lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `DownloadInterfaceVersion` (DD 2.60a) to the function comment and add a `const` block for its 2-byte layout.
*   **ASN.1 Specification (Data Dictionary 2.60a):**
    ```asn.1
    // DownloadInterfaceVersion ::= OCTET STRING (SIZE (2))
    // Value assignment: 'aabb'H
    // 'aa'H: Generation ('01'H for Gen1, '02'H for Gen2)
    // 'bb'H: Version
    ```

### Task: Implement Stubbed VU Unmarshalling Files

*   **Files:** `unmarshal_vu_detailed_speed.go`, `unmarshal_vu_events_faults.go`, `unmarshal_vu_technical_data.go`
*   **Problem:** These files are stubs that read the entire data block into a signature field instead of parsing it. This violates the "full binary roundtrip" goal.
*   **Action:** For each file, create a high-level implementation task:
    *   **`unmarshal_vu_detailed_speed.go`:** Implement the parser for `VuDetailedSpeedData` (DD 2.192).
    *   **`unmarshal_vu_events_faults.go`:** Implement the parser for `VuEventData` (DD 2.197) and `VuFaultData` (DD 2.200).
    *   **`unmarshal_vu_technical_data.go`:** Implement the parser for `VuIdentification` (DD 2.205), `VuCalibrationData` (DD 2.173), and other technical data blocks.
    *   Each implementation must include proper ASN.1 documentation and layout constants.

### Task: Implement Stubbed VU Marshalling Files

*   **Files:** `append_vu_activities.go`, `append_vu_detailed_speed.go`, `append_vu_events_faults.go`, `append_vu_technical_data.go`, `append_vu_overview.go`
*   **Problem:** These files are stubs that only write back a signature or are incomplete, preventing roundtrip serialization.
*   **Action:** For each file, create a high-level implementation task to correctly serialize the corresponding protobuf message (`Activities`, `DetailedSpeed`, etc.) into the binary format specified by the regulation. Each implementation must be documented with ASN.1 specifications and use layout constants.

### Task: Document `AppendDownloadInterfaceVersion`

*   **File:** `append_vu_download_interface_version.go`
*   **Function:** `AppendDownloadInterfaceVersion`
*   **Problem:** The function is implemented but lacks ASN.1 documentation and layout constants.
*   **Action:** Add the ASN.1 specification for `DownloadInterfaceVersion` (DD 2.60a) and a `const` block for its layout.

This concludes the initial audit. The log now contains a comprehensive set of tasks covering helper file refactoring, documentation alignment, and implementation of stubbed functionality for both card and VU files.
