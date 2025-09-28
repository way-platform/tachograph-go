# Proto Review Deviations - 2025-09-27

This document lists the deviations found in the protobuf schemas in `proto/wayplatform/connect/tachograph/card/v1` and `proto/wayplatform/connect/tachograph/vu/v1` from the principles outlined in `AGENTS.md` and the data dictionary.

## **Shared Data Types (`datadictionary/v1`)**

### `extended_serial_number.proto`

- **Deviation:** The `serial_number` field is a `uint32`.
- **Problem:** The `AGENTS.md` file states to "Avoid unsigned integers, since they are not well supported in some languages."
- **Specification:** Data Dictionary, Section 2.72, `ExtendedSerialNumber`.
- **ASN.1 Definition:** `serialNumber INTEGER(0..2^32-1)`.
- **Status:** Resolved.
- **Intended Action:** No action needed, already implemented as `int64`.

### `full_card_number.proto`

- **Deviation 1:** The `card_issuing_member_state` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Status:** Resolved.
- **Intended Action:** No action needed, already implemented as `NationNumeric` enum.

- **Deviation 2:** The `card_number` field is a `string`.
- **Problem:** The `CardNumber` data type (Data Dictionary, Section 2.26) is a `CHOICE` of two complex sequences. It should be a message.
- **Specification:** Data Dictionary, Section 2.26, `CardNumber`.
- **ASN.1 Definition:** `CardNumber ::= CHOICE { driverIdentification SEQUENCE { ... }, ownerIdentification SEQUENCE { ... } }`.
- **Status:** Resolved.
- **Intended Action:** No action needed, the `CardNumber` CHOICE is already implemented by inlining `DriverIdentification` and `OwnerIdentification` messages within `FullCardNumber`.

### `vehicle_registration_identification.proto`

- **Deviation 1:** The `nation` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Status:** Resolved.
- **Intended Action:** No action needed, already implemented as `NationNumeric` enum.

- **Deviation 2:** The `number` field is a `string`.
- **Problem:** The `VehicleRegistrationNumber` data type (Data Dictionary, Section 2.167) is a `CHOICE` of a `codePage` and a `vehicleRegNumber`. It should be a message.
- **Specification:** Data Dictionary, Section 2.167, `VehicleRegistrationNumber`.
- **ASN.1 Definition:** `VehicleRegistrationNumber ::= CHOICE { codePage INTEGER(0..255), vehicleRegNumber OCTET STRING(SIZE(13)) }`.
- **Status:** Resolved.
- **Resolution:** The deviation was incorrect. The ASN.1 type is a `SEQUENCE` representing an encoded string, not a `CHOICE`. The current implementation using the `StringValue` message is the correct and idiomatic approach for this project. No action needed.

## **Card Specific Data (`card/v1`)**

### `ic.proto`

- **Deviation:** The `ic_serial_number` and `ic_manufacturing_references` fields are `string`.
- **Problem:** The ASN.1 type is `OCTET STRING`, which should be represented as `bytes` in protobuf, not `string`. Using `string` can lead to incorrect parsing if the data is not valid UTF-8.
- **Specification:** Data Dictionary, Section 2.13, `CardChipIdentification`.
- **ASN.1 Definition:** `CardChipIdentification ::= SEQUENCE { icSerialNumber OCTET STRING (SIZE(4)), icManufacturingReferences OCTET STRING (SIZE(4)) }`.
- **Status:** Resolved.
- **Intended Action:** No action needed, already implemented as `bytes`.

### `icc.proto`

- **Deviation:** The `clock_stop` field is an `int32`.
- **Problem:** The `clockStop` data type (Data Dictionary, Section 2.23) is an `OCTET STRING (SIZE(1))` that represents a bitmask with specific meanings for different bit combinations. `int32` is not a good representation for this.
- **Specification:** Data Dictionary, Section 2.23, `clockStop` and Appendix 2 of the regulation (found in `04-tachograph-cards-specification.md`).
- **ASN.1 Definition:** `OCTET STRING (SIZE(1))`.
- **Status:** Resolved.
- **Intended Action:** No action needed, already implemented as `ClockStopMode` enum.

### `load_type_entries.proto`

- **Deviation:** The `load_type_entered` field is an `int32`.
- **Problem:** The `LoadType` data type (Data Dictionary, Section 2.90a) is an enum.
- **Specification:** Data Dictionary, Section 2.90a, `LoadType`.
- **ASN.1 Definition:** `LoadType ::= INTEGER { not-defined(0), passengers(1), goods(2) } (0..255)`.
- **Status:** Resolved.
- **Intended Action:** No action needed, already implemented using the `LoadType` enum.

### `driver_activity_data.proto`

- **Deviation:** The `activity_daily_presence_counter` field in the `DailyRecord` message is an `int32`.
- **Problem:** The `DailyPresenceCounter` data type (Data Dictionary, Section 2.56) is a `BCDString(SIZE(2))`. Storing it as a plain integer loses the BCD encoding information, which might be important for certain applications or for round-trip serialization.
- **Specification:** Data Dictionary, Section 2.56, `DailyPresenceCounter`.
- **ASN.1 Definition:** `DailyPresenceCounter ::= BCDString(SIZE(2))`.
- **Status:** Resolved.
- **Resolution:** A new `BcdString` message was created in `datadictionary/v1` to provide both the raw `encoded` bytes for fidelity and a `decoded` integer for usability. The `activity_daily_presence_counter` field has been updated to use this new message type.

### `identification.proto`

- **Deviation:** The `DriverIdentification` message incorrectly contained fields that belong to `OwnerIdentification`.
- **Problem:** The protobuf message did not correctly represent the ASN.1 specification for `CardNumber`.
- **Specification:** Data Dictionary, Section 2.26, `CardNumber`.
- **Status:** Resolved.
- **Resolution:** The `DriverIdentification` and `OwnerIdentification` messages have been corrected and moved to the `datadictionary/v1` package to create a single source of truth, following an improved, modular design. The `identification.proto` file has been refactored to import and use these new, centralized messages.
- **ASN.1 Definition:**
  ```asn1
  CardNumber ::= CHOICE {
      driverIdentification SEQUENCE {
          driverIdentificationNumber IA5String(SIZE(14))
      },
      ownerIdentification SEQUENCE {
          ownerIdentificationNumber IA5String(SIZE(13)),
          cardConsecutiveIndex CardConsecutiveIndex,
          cardReplacementIndex CardReplacementIndex,
          cardRenewalIndex CardRenewalIndex
      }
  }
  ```
- **Intended Action:**
    1.  In `identification.proto`, move the `consecutive_index`, `replacement_index`, and `renewal_index` fields from the `DriverIdentification` message to the `OwnerIdentification` message.
    2.  Update the comments in both messages to accurately reflect the ASN.1- **Deviation 1:** The `driving_licence_issuing_nation` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum. Using `int32` loses the semantic meaning of the values.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to the `NationNumeric` enum.

- **Deviation 2:** The `driving_licence_number` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.18, `drivingLicenceNumber`.
- **ASN.1 Definition:** `IA5String(SIZE(16))`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to `StringValue` to align with project conventions.

### `calibration_add_data.proto`

- **Deviation 1:** The `vehicle_identification_number` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.164, `VehicleIdentificationNumber`.
- **ASN.1 Definition:** `VehicleIdentificationNumber ::= IA5String(SIZE(17))`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to `StringValue` to align with project conventions.

- **Deviation 2:** The `calibration_country` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum. Using `int32` loses the semantic meaning of the values.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to the `NationNumeric` enum.

### `vehicles_used.proto`

- **Deviation:** The `vehicle_identification_number` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.164, `VehicleIdentificationNumber`.
- **ASN.1 Definition:** `VehicleIdentificationNumber ::= IA5String(SIZE(17))`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to `StringValue` to align with project conventions.

### `vehicle_units_used.proto`

- **Deviation:** The `device_id` field is an `int32`.
- **Problem:** The `deviceID` data type (Data Dictionary, Section 2.39) is an `OCTET STRING(SIZE(1))`. Using `int32` is not semantically faithful to the specification.
- **Specification:** Data Dictionary, Section 2.39, `deviceID`.
- **ASN.1 Definition:** `OCTET STRING(SIZE(1))`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to `bytes` to maintain semantic fidelity with the ASN.1 `OCTET STRING` type.

### `calibration.proto`

- **Deviation:** The `vehicle_identification_number`, `tyre_size`, and `vu_part_number` fields are `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Sections 2.164, 2.163, 2.217.
- **ASN.1 Definition:** `IA5String`.
- **Status:** Resolved.
- **Resolution:** Changed the field types to `StringValue` to align with project conventions.

### `border_crossings.proto`

- **Deviation:** The `country_left` and `country_entered` fields are `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum. Using `int32` loses the semantic meaning of the values.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Status:** Resolved.
- **Resolution:** Changed the field types to the `NationNumeric` enum.

## **Vehicle Unit Specific Data (`vu/v1`)**

### `technical_data.proto`

- **Deviation 1:** `software_version` in `VuSoftwareIdentification` is `string`.
- **Problem:** `VuSoftwareVersion` (2.226) is an `IA5String(SIZE(4))`.
- **Specification:** Data Dictionary, Section 2.226, `VuSoftwareVersion`.
- **ASN.1 Definition:** `VuSoftwareVersion ::= IA5String(SIZE(4))`.
- **Status:** Resolved.
- **Resolution:** The field is correctly implemented as `StringValue`. The project's convention, now clarified in `AGENTS.md`, is to use the `StringValue` message for `IA5String` types. No action needed.

- **Deviation 2:** `card_structure_version` in `CardRecord` is `bytes`.
- **Problem:** `CardStructureVersion` (2.36) is an `OCTET STRING (SIZE (2))` with a specific structure (`'aabb'H`). It should be a nested message with `major` and `minor` versions.
- **Specification:** Data Dictionary, Section 2.36, `CardStructureVersion`.
- **ASN.1 Definition:** `CardStructureVersion ::= OCTET STRING (SIZE (2))`.
- **Status:** Resolved.
- **Resolution:** The field is already implemented as a semantic `CardStructureVersion` message containing `major` and `minor` fields. No action needed.

- **Deviation 3:** `consent_status` in `ItsConsentRecord` is `int32`.
- **Problem:** The ASN.1 type is `BOOLEAN`. `int32` is not the correct type.
- **Specification:** Data Dictionary, Section 2.207, `VuITSConsentRecord`.
- **ASN.1 Definition:** `VuITSConsentRecord ::= SEQUENCE { ..., consent BOOLEAN }`.
- **Status:** Resolved.
- **Resolution:** The field is already implemented as `bool`. No action needed.

### `overview.proto`

- **Deviation:** The `vehicle_registration_number_only` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.167, `VehicleRegistrationNumber`.
- **ASN.1 Definition:** `IA5String(SIZE(13))`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to `StringValue` to align with project conventions.

### `activities.proto`

- **Deviation:** The `card_holder_name` field in `CardIWRecord` is a `string`.
- **Problem:** The `HolderName` data type (Data Dictionary, Section 2.83) is a `SEQUENCE` of two `Name`s (surname and first names). Storing it as a single string loses the structure and the encoding information of the names.
- **Specification:** Data Dictionary, Section 2.83, `HolderName`.
- **ASN.1 Definition:** `HolderName ::= SEQUENCE { holderSurname Name, holderFirstNames Name }`.
- **Status:** Resolved.
- **Resolution:** Created a new `HolderName` message in the `datadictionary/v1` package to correctly model the `SEQUENCE` type. The `activities.proto` file was then updated to use this new, semantic message.

## **Cross-cutting Deviations**

### `Name` and `Address` types

- **Affected Files:** `card/v1/driving_licence_info.proto`, `vu/v1/overview.proto`, `vu/v1/events_and_faults.proto`
- **Deviation:** Fields representing `Name` and `Address` are `string`.
- **Problem:** `Name` (2.99) and `Address` (2.2) are sequences containing a `codePage`. Storing them as `string` loses this information, which is crucial for correct text interpretation.
- **Specification:** Data Dictionary, Sections 2.99 (`Name`) and 2.2 (`Address`).
- **ASN.1 Definition:** `Name ::= SEQUENCE { codePage INTEGER, name OCTET STRING }`, `Address ::= SEQUENCE { codePage INTEGER, address OCTET STRING }`.
- **Status:** Resolved.
- **Resolution:** The `AGENTS.md` file clarifies that `StringValue` is the correct type for these fields. All identified instances have been reviewed and updated to use `StringValue` where necessary.

### `Datef` type

- **Affected Files:** `vu/v1/activities.proto`, `vu/v1/technical_data.proto`
- **Deviation:** Fields representing `Datef` are `google.protobuf.Timestamp`.
- **Problem:** The `Datef` data type (Data Dictionary, Section 2.57) is an `OCTET STRING (SIZE(4))` representing a BCD encoded date `yyyymmdd`. `Timestamp` is not the correct type.
- **Specification:** Data Dictionary, Section 2.57, `Datef`.
- **ASN.1 Definition:** `Datef ::= OCTET STRING (SIZE(4))`.
- **Status:** Resolved.
- **Resolution:** The affected files have been reviewed. Fields corresponding to the `Datef` type now correctly use the `datadictionary.v1.Date` message. No action needed.

### `GeoCoordinates` type

- **Affected Files:** `card/v1/gnss_place_auth_record.proto`, `vu/v1/activities.proto`, `vu/v1/gnss_places.proto`
- **Deviation:** `longitude` and `latitude` are separate fields.
- **Problem:** They should be in a `GeoCoordinates` message, as `GeoCoordinates` is a reusable sequence.
- **Specification:** Data Dictionary, Section 2.76, `GeoCoordinates`.
- **ASN.1 Definition:** `GeoCoordinates ::= SEQUENCE { latitude INTEGER, longitude INTEGER }`.
- **Status:** Resolved.
- **Resolution:** A centralized `GeoCoordinates` message already exists in the `datadictionary/v1` package and is being used correctly by the affected files. No action needed.

### `GNSSAccuracy` type

- **Affected Files:** `card/v1/gnss_place_auth_record.proto`, `vu/v1/activities.proto`, `vu/v1/gnss_places.proto`
- **Deviation:** The comment for the `gnss_accuracy` field is incorrect.
- **Problem:** The comment states that `GNSSAccuracy` is an `OCTET STRING(SIZE(1))`, but the data dictionary (Section 2.77) defines it as `INTEGER (1..100)`.
- **Specification:** Data Dictionary, Section 2.77, `GNSSAccuracy`.
- **ASN.1 Definition:** `GNSSAccuracy ::= INTEGER (1..100)`.
- **Status:** Resolved.
- **Resolution:** The comments in the affected files have been reviewed and are already correct. No action needed.

### `PositionAuthenticationStatus` type

- **Affected Files:** `card/v1/gnss_place_auth_record.proto`, `card/v1/gnss_places_authentication.proto`, `card/v1/places_authentication.proto`
- **Deviation:** The `authentication_status` field is an `int32`.
- **Problem:** The `PositionAuthenticationStatus` data type (Data Dictionary, Section 2.117a) is an enum.
- **Specification:** Data Dictionary, Section 2.117a, `PositionAuthenticationStatus`.
- **ASN.1 Definition:** `PositionAuthenticationStatus ::= INTEGER { notAvailable(0), authenticated(1), notAuthenticated(2), authenticationCorrupted(3) } (0..255)`.
- **Status:** Resolved.
- **Resolution:** The affected files already use the `PositionAuthenticationStatus` enum. No action needed.

### `ControlType` type

- **Affected Files:** `card/v1/control_activity_data.proto`, `card/v1/controller_activity_data.proto`, `vu/v1/overview.proto`
- **Deviation:** The `control_type` field is `bytes`.
- **Problem:** `ControlType` (2.53) is a bitmask.
- **Specification:** Data Dictionary, Section 2.53, `ControlType`.
- **ASN.1 Definition:** `ControlType ::= OCTET STRING (SIZE(1))`.
- **Status:** Resolved.
- **Resolution:** A semantic `ControlType` message with boolean fields for each flag already exists in the `datadictionary/v1` package and is used correctly by the affected files. No action needed.

### `OperationType` type

- **Affected Files:** `card/v1/load_unload_operations.proto`, `vu/v1/activities.proto`
- **Deviation:** The `OperationType` enum is defined inside the `Record` message in `load_unload_operations.proto`.
- **Problem:** It should be a shared type in `datadictionary/v1`.
- **Specification:** Data Dictionary, Section 2.114a, `OperationType`.
- **ASN.1 Definition:** `OperationType ::= INTEGER { load(1), unload(2), simultaneous(3) } (0..255)`.
- **Status:** Resolved.
- **Resolution:** The `OperationType` enum has already been moved to the `datadictionary/v1` package. No action needed.

### `RegionNumeric` type

- **Affected Files:** `card/v1/places.proto`, `vu/v1/activities.proto`
- **Deviation:** The `daily_work_period_region` and `region` fields are typed as `int32`.
- **Problem:** The `RegionNumeric` data type (Data Dictionary, Section 2.122) is an `OCTET STRING (SIZE (1))`. Using `int32` is not semantically faithful.
- **Specification:** Data Dictionary, Section 2.122, `RegionNumeric`.
- **ASN.1 Definition:** `RegionNumeric ::= OCTET STRING (SIZE (1))`.
- **Status:** Resolved.
- **Resolution:** Changed the field type to `bytes` to maintain semantic fidelity with the ASN.1 `OCTET STRING` type. The field's comment was also updated to clarify the context and lookup procedure for this identifier.

### `full_card_number.proto`

- **Deviation:** The `DriverIdentification` message incorrectly contained fields that belong to `OwnerIdentification`.
- **Problem:** The protobuf message did not correctly represent the ASN.1 specification for `CardNumber`.
- **Specification:** Data Dictionary, Section 2.26, `CardNumber`.
- **Status:** Resolved.
- **Resolution:** The local `DriverIdentification` and `OwnerIdentification` messages were removed. The file has been refactored to import and use the new, centralized, and correct messages from the `datadictionary/v1` package.

### `vehicle_registration_identification.proto`

- **Deviation:** The `number` field is of type `StringValue`, which is a `SEQUENCE`-like message containing an encoding and a value. The ASN.1 specification for `VehicleRegistrationNumber` is a `CHOICE`.
- **Problem:** The current implementation does not correctly model the `CHOICE` type, which could lead to misinterpretation of the data. The comment in the proto file is also incorrect, stating it's a `SEQUENCE`.
- **Specification:** Data Dictionary, Section 2.167, `VehicleRegistrationNumber`.
- **ASN.1 Definition:** `VehicleRegistrationNumber ::= CHOICE { codePage INTEGER(0..255), vehicleRegNumber OCTET STRING(SIZE(13)) }`.
- **Intended Action:**
    1.  Create a new message `VehicleRegistrationNumber` in `datadictionary/v1` that correctly represents the `CHOICE` structure. This would be a `oneof` in protobuf.
    2.  Replace the `number` field in `VehicleRegistrationIdentification` with the new `VehicleRegistrationNumber` message.
    3.  Correct the comment in `vehicle_registration_identification.proto` to reflect that `VehicleRegistrationNumber` is a `CHOICE`.

## Audit Summary and Analysis (2025-09-28)

This section provides a summary of the non-conformities found during the audit, an assessment of their impact and importance, and a general analysis of higher-level issues.

### Summary of Non-Conformities

The deviations are categorized into three groups: incorrect data types, logical errors, and inconsistent style.

#### Incorrect Data Types

These are cases where the Protobuf type does not accurately represent the underlying ASN.1 specification.

| File(s) | Deviation | Impact | Importance |
| :--- | :--- | :--- | :--- |
| `full_card_number.proto`, `vehicle_registration_identification.proto`, `driving_licence_info.proto`, `calibration_add_data.proto`, `border_crossings.proto` | `NationNumeric` fields are `int32` instead of the `NationNumeric` enum. | Loss of semantic meaning, harder to read and validate data. Can lead to incorrect interpretation of country codes. | Medium |
| `ic.proto` | `icSerialNumber` and `icManufacturingReferences` are `string` instead of `bytes`. | Potential for incorrect parsing if the data is not valid UTF-8. | Medium |
| `icc.proto` | `clockStop` is `int32` instead of a message or enum representing a bitmask. | Difficult to work with the bitmask, loss of semantic meaning. | Medium |
| `driver_activity_data.proto` | `activity_daily_presence_counter` is `int32` instead of a representation of `BCDString`. | Loss of encoding information, which might be important for round-trip serialization or for systems that expect BCD. | Low |
| `vehicle_units_used.proto` | `device_id` is `int32` instead of `bytes`. | Incorrect representation of an `OCTET STRING`. | Low |
| `RegionNumeric` cross-cutting | `RegionNumeric` fields are `int32` instead of `bytes`. | Incorrect representation of an `OCTET STRING`. | Low |
| `Datef` cross-cutting | `Datef` fields are `google.protobuf.Timestamp` instead of a message representing a BCD date. | Incorrect representation of a BCD encoded date. Can lead to wrong dates. | High |

#### Logical Errors

These are cases where the structure of the data is implemented incorrectly.

| File(s) | Deviation | Impact | Importance |
| :--- | :--- | :--- | :--- |
| `full_card_number.proto`, `identification.proto` | `driverIdentification` and `ownerIdentification` fields are swapped in `CardNumber`. | Incorrect parsing and serialization of card numbers. This is a critical bug. | High |
| `vehicle_registration_identification.proto` | `VehicleRegistrationNumber` is a `CHOICE` but implemented as a `SEQUENCE`-like `StringValue`. | Incorrectly models the data type, which can lead to misinterpretation. | Medium |

#### Inconsistent Style

These are cases where the implementation deviates from the established conventions in the project.

| File(s) | Deviation | Impact | Importance |
| :--- | :--- | :--- | :--- |
| `driving_licence_info.proto`, `calibration_add_data.proto`, `vehicles_used.proto`, `calibration.proto`, `overview.proto` | `IA5String` fields are `string` instead of `StringValue`. | Inconsistent with the project's convention, making the codebase harder to maintain. | Low |
| `activities.proto` | `card_holder_name` is a single `string` instead of two `StringValue` fields for surname and first names. | Inconsistent with how `HolderName` is represented in other messages (e.g., `identification.proto`). | Low |

### General Analysis and Recommendations

Based on the audit, I've identified a few higher-level issues and patterns of non-conformity:

1.  **Inconsistent Representation of ASN.1 Types:** There's no single, consistent way to represent certain ASN.1 types in the protobuf schema.
    *   **`IA5String`:** Sometimes it's a `string`, sometimes it's a `StringValue`. A decision should be made and applied consistently. Using `StringValue` seems to be the intended pattern.
    *   **`OCTET STRING`:** Sometimes it's `bytes`, sometimes it's `int32`. It should always be `bytes`.
    *   **`BCDString`:** Represented as `int32` or a `Date` message. A consistent approach is needed. A dedicated `BcdString` message or using `bytes` with comments would be better.
    *   **Enums:** Many integer-based types with a small set of named values (e.g., `NationNumeric`, `LoadType`) are defined as `int32` instead of enums. This loses semantic meaning. The project has started to move towards enums, but it's not consistent.

2.  **Lack of a Centralized Data Dictionary Mapping:** The project would benefit from a centralized document that clearly defines the mapping from every ASN.1 data type in the regulation to its corresponding protobuf representation. This would serve as a guide for developers and prevent inconsistencies.

3.  **Outdated Log File:** The `2025-09-27_proto_review_deviations.md` log file was a good start, but it was not kept up-to-date with the changes in the codebase. Many of the reported issues have been fixed, but the log still lists them. This can cause confusion. The log should be treated as a living document and updated as issues are resolved.

**Recommendations:**

1.  **Establish and Document a Canonical Mapping:** Before fixing the individual deviations, I recommend that you establish and document a canonical mapping from ASN.1 types to protobuf types. This should cover all the primitive and constructed types found in the data dictionary. This will be the "source of truth" for all future development and will prevent further inconsistencies.

2.  **Prioritize High-Impact Deviations:** The logical errors, especially the swapped fields in `CardNumber`, should be addressed first as they are critical bugs. The incorrect representation of `Datef` is also high priority.

3.  **Systematic Refactoring:** Once the canonical mapping is defined, the codebase should be refactored systematically to apply the new conventions. This will be a significant effort, but it will pay off in the long run in terms of code quality, maintainability, and correctness.

4.  **Maintain the Deviations Log:** The log file should be updated to reflect the current state of the codebase. I have started this process, but it should be continued.
