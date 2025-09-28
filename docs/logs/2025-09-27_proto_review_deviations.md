# Proto Review Deviations - 2025-09-27

This document lists the deviations found in the protobuf schemas in `proto/wayplatform/connect/tachograph/card/v1` and `proto/wayplatform/connect/tachograph/vu/v1` from the principles outlined in `AGENTS.md` and the data dictionary.

## **Shared Data Types (`datadictionary/v1`)**

### `extended_serial_number.proto`

- **Deviation:** The `serial_number` field is a `uint32`.
- **Problem:** The `AGENTS.md` file states to "Avoid unsigned integers, since they are not well supported in some languages."
- **Specification:** Data Dictionary, Section 2.72, `ExtendedSerialNumber`.
- **ASN.1 Definition:** `serialNumber INTEGER(0..2^32-1)`.
- **Intended Action:** Change the type of `serial_number` to `int64` to avoid using an unsigned integer.

### `full_card_number.proto`

- **Deviation 1:** The `card_issuing_member_state` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Intended Action:** Change the type of `card_issuing_member_state` to use the existing `NationNumeric` enum from `datadictionary/v1`.

- **Deviation 2:** The `card_number` field is a `string`.
- **Problem:** The `CardNumber` data type (Data Dictionary, Section 2.26) is a `CHOICE` of two complex sequences. It should be a message.
- **Specification:** Data Dictionary, Section 2.26, `CardNumber`.
- **ASN.1 Definition:** `CardNumber ::= CHOICE { driverIdentification SEQUENCE { ... }, ownerIdentification SEQUENCE { ... } }`.
- **Intended Action:** Create a new message `CardNumber` in `datadictionary/v1` to represent this `CHOICE` and use it in the `FullCardNumber` message.

### `vehicle_registration_identification.proto`

- **Deviation 1:** The `nation` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Intended Action:** Change the type of `nation` to use the existing `NationNumeric` enum from `datadictionary/v1`.

- **Deviation 2:** The `number` field is a `string`.
- **Problem:** The `VehicleRegistrationNumber` data type (Data Dictionary, Section 2.167) is a `CHOICE` of a `codePage` and a `vehicleRegNumber`. It should be a message.
- **Specification:** Data Dictionary, Section 2.167, `VehicleRegistrationNumber`.
- **ASN.1 Definition:** `VehicleRegistrationNumber ::= CHOICE { codePage INTEGER(0..255), vehicleRegNumber OCTET STRING(SIZE(13)) }`.
- **Intended Action:** Create a new message `VehicleRegistrationNumber` in `datadictionary/v1` to represent this `CHOICE` and use it in the `VehicleRegistrationIdentification` message.

## **Card Specific Data (`card/v1`)**

### `ic.proto`

- **Deviation:** The `ic_serial_number` and `ic_manufacturing_references` fields are `string`.
- **Problem:** The ASN.1 type is `OCTET STRING`, which should be represented as `bytes` in protobuf, not `string`. Using `string` can lead to incorrect parsing if the data is not valid UTF-8.
- **Specification:** Data Dictionary, Section 2.13, `CardChipIdentification`.
- **ASN.1 Definition:** `CardChipIdentification ::= SEQUENCE { icSerialNumber OCTET STRING (SIZE(4)), icManufacturingReferences OCTET STRING (SIZE(4)) }`.
- **Intended Action:** Change the type of `ic_serial_number` and `ic_manufacturing_references` to `bytes`.

### `icc.proto`

- **Deviation:** The `clock_stop` field is an `int32`.
- **Problem:** The `clockStop` data type (Data Dictionary, Section 2.23) is an `OCTET STRING (SIZE(1))` that represents a bitmask with specific meanings for different bit combinations. `int32` is not a good representation for this.
- **Specification:** Data Dictionary, Section 2.23, `clockStop` and Appendix 2 of the regulation (found in `04-tachograph-cards-specification.md`).
- **ASN.1 Definition:** `OCTET STRING (SIZE(1))`.
- **Interpretation:** The byte is a bitmask that defines the clock stop mode. The meaning of the bits is defined in `04-tachograph-cards-specification.md` as follows:
    | Bit 3 | Bit 2 | Bit 1 | Meaning                                 |
    |-------|-------|-------|-----------------------------------------|
    | 0     | 0     | 1     | Clockstop allowed, no preferred level   |
    | 0     | 1     | 1     | Clockstop allowed, high level preferred |
    | 1     | 0     | 1     | Clockstop allowed, low level preferred  |
    | 0     | 0     | 0     | Clockstop not allowed                   |
    | 0     | 1     | 0     | Clockstop only allowed on high level    |
    | 1     | 0     | 0     | Clockstop only allowed on low level     |
- **Intended Action:**
    1. Create a new enum `ClockStopMode` in `icc.proto` with values for each mode.
    2. Change the type of `clock_stop` to the new `ClockStopMode` enum.

### `load_type_entries.proto`

- **Deviation:** The `load_type_entered` field is an `int32`.
- **Problem:** The `LoadType` data type (Data Dictionary, Section 2.90a) is an enum.
- **Specification:** Data Dictionary, Section 2.90a, `LoadType`.
- **ASN.1 Definition:** `LoadType ::= INTEGER { not-defined(0), passengers(1), goods(2) } (0..255)`.
- **Intended Action:** Update `load_type_entries.proto` to import and use the existing `LoadType` enum from `datadictionary/v1`.

### `driver_activity_data.proto`

- **Deviation:** The `activity_daily_presence_counter` field in the `DailyRecord` message is an `int32`.
- **Problem:** The `DailyPresenceCounter` data type (Data Dictionary, Section 2.56) is a `BCDString(SIZE(2))`. Storing it as a plain integer loses the BCD encoding information, which might be important for certain applications or for round-trip serialization.
- **Specification:** Data Dictionary, Section 2.56, `DailyPresenceCounter`.
- **ASN.1 Definition:** `DailyPresenceCounter ::= BCDString(SIZE(2))`.
- **Intended Action:** Change the type of `activity_daily_presence_counter` to `bytes` to store the raw BCD value. Alternatively, create a new message `BcdString` in `datadictionary/v1` that can represent BCD-encoded strings and use it here. The comment should also be updated to clarify the BCD encoding.

### `identification.proto`

- **Deviation:** The fields for `driverIdentification` and `ownerIdentification` are swapped in the `Card` message. The `DriverIdentification` message incorrectly contains `consecutive_index`, `replacement_index`, and `renewal_index`, while the `OwnerIdentification` message is missing them.
- **Problem:** The protobuf message does not correctly represent the ASN.1 specification for `CardNumber`. This will lead to incorrect parsing and serialization of card numbers for drivers and owners.
- **Specification:** Data Dictionary, Section 2.26, `CardNumber`.
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
- **Intended Action:** Change the type of `driving_licence_issuing_nation` to use the existing `NationNumeric` enum from `datadictionary/v1`.

- **Deviation 2:** The `driving_licence_number` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.18, `drivingLicenceNumber`.
- **ASN.1 Definition:** `IA5String(SIZE(16))`.
- **Intended Action:** Change the type of `driving_licence_number` to `wayplatform.connect.tachograph.datadictionary.v1.StringValue`.

### `calibration_add_data.proto`

- **Deviation 1:** The `vehicle_identification_number` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.164, `VehicleIdentificationNumber`.
- **ASN.1 Definition:** `VehicleIdentificationNumber ::= IA5String(SIZE(17))`.
- **Intended Action:** Change the type of `vehicle_identification_number` to `wayplatform.connect.tachograph.datadictionary.v1.StringValue`.

- **Deviation 2:** The `calibration_country` field is an `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum. Using `int32` loses the semantic meaning of the values.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Intended Action:** Change the type of `calibration_country` to use the existing `NationNumeric` enum from `datadictionary/v1`.

### `vehicles_used.proto`

- **Deviation:** The `vehicle_identification_number` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.164, `VehicleIdentificationNumber`.
- **ASN.1 Definition:** `VehicleIdentificationNumber ::= IA5String(SIZE(17))`.
- **Intended Action:** Change the type of `vehicle_identification_number` to `wayplatform.connect.tachograph.datadictionary.v1.StringValue`.

### `vehicle_units_used.proto`

- **Deviation:** The `device_id` field is an `int32`.
- **Problem:** The `deviceID` data type (Data Dictionary, Section 2.39) is an `OCTET STRING(SIZE(1))`. Using `int32` is incorrect for an octet string. It should be `bytes`.
- **Specification:** Data Dictionary, Section 2.39, `deviceID`.
- **ASN.1 Definition:** `OCTET STRING(SIZE(1))`.
- **Intended Action:** Change the type of `device_id` to `bytes`.

### `calibration.proto`

- **Deviation:** The `vehicle_identification_number`, `tyre_size`, and `vu_part_number` fields are `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Sections 2.164, 2.163, 2.217.
- **ASN.1 Definition:** `IA5String`.
- **Intended Action:** Change the type of these fields to `wayplatform.connect.tachograph.datadictionary.v1.StringValue`.

### `border_crossings.proto`

- **Deviation:** The `country_left` and `country_entered` fields are `int32`.
- **Problem:** The `NationNumeric` data type (Data Dictionary, Section 2.101) is an enum. Using `int32` loses the semantic meaning of the values.
- **Specification:** Data Dictionary, Section 2.101, `NationNumeric`.
- **ASN.1 Definition:** `NationNumeric ::= INTEGER(0..255)`.
- **Intended Action:** Change the type of `country_left` and `country_entered` to use the existing `NationNumeric` enum from `datadictionary/v1`.

## **Vehicle Unit Specific Data (`vu/v1`)**

### `technical_data.proto`

- **Deviation 1:** `software_version` in `VuSoftwareIdentification` is `string`.
- **Problem:** `VuSoftwareVersion` (2.226) is an `IA5String(SIZE(4))`. `IA5String` is a subset of ASCII. While `string` is acceptable, `bytes` is a more faithful representation of a fixed-size string, especially if the content is not guaranteed to be valid UTF-8.
- **Specification:** Data Dictionary, Section 2.226, `VuSoftwareVersion`.
- **ASN.1 Definition:** `VuSoftwareVersion ::= IA5String(SIZE(4))`.
- **Intended Action:** Change the type of `software_version` to `bytes`.

- **Deviation 2:** `card_structure_version` in `CardRecord` is `bytes`.
- **Problem:** `CardStructureVersion` (2.36) is an `OCTET STRING (SIZE (2))` with a specific structure (`'aabb'H`). It should be a nested message with `major` and `minor` versions.
- **Specification:** Data Dictionary, Section 2.36, `CardStructureVersion`.
- **ASN.1 Definition:** `CardStructureVersion ::= OCTET STRING (SIZE (2))`.
- **Intended Action:** Change the type of `card_structure_version` to a nested `CardStructureVersion` message with `major` and `minor` fields.

- **Deviation 3:** `consent_status` in `ItsConsentRecord` is `int32`.
- **Problem:** The ASN.1 type is `BOOLEAN`. `int32` is not the correct type.
- **Specification:** Data Dictionary, Section 2.207, `VuITSConsentRecord`.
- **ASN.1 Definition:** `VuITSConsentRecord ::= SEQUENCE { ..., consent BOOLEAN }`.
- **Intended Action:** Change the type of `consent_status` to `bool`.

### `overview.proto`

- **Deviation:** The `vehicle_registration_number_only` field is a `string`.
- **Problem:** The project uses the `StringValue` message to represent complex string types like `IA5String`. Using a raw `string` is inconsistent with this pattern.
- **Specification:** Data Dictionary, Section 2.167, `VehicleRegistrationNumber`.
- **ASN.1 Definition:** `IA5String(SIZE(13))`.
- **Intended Action:** Change the type of `vehicle_registration_number_only` to `wayplatform.connect.tachograph.datadictionary.v1.StringValue`.

### `activities.proto`

- **Deviation:** The `card_holder_name` field in `CardIWRecord` is a `string`.
- **Problem:** The `HolderName` data type (Data Dictionary, Section 2.83) is a `SEQUENCE` of two `Name`s (surname and first names). Storing it as a single string loses the structure and the encoding information of the names.
- **Specification:** Data Dictionary, Section 2.83, `HolderName`.
- **ASN.1 Definition:** `HolderName ::= SEQUENCE { holderSurname Name, holderFirstNames Name }`.
- **Intended Action:** Replace the `card_holder_name` field with two fields `card_holder_surname` and `card_holder_first_names` of type `wayplatform.connect.tachograph.datadictionary.v1.StringValue`, similar to what is done in `identification.proto`.

## **Cross-cutting Deviations**

### `Name` and `Address` types

- **Affected Files:** `card/v1/driving_licence_info.proto`, `vu/v1/overview.proto`, `vu/v1/events_and_faults.proto`
- **Deviation:** Fields representing `Name` and `Address` are `string`.
- **Problem:** `Name` (2.99) and `Address` (2.2) are sequences containing a `codePage`. Storing them as `string` loses this information, which is crucial for correct text interpretation.
- **Specification:** Data Dictionary, Sections 2.99 (`Name`) and 2.2 (`Address`).
- **ASN.1 Definition:** `Name ::= SEQUENCE { codePage INTEGER, name OCTET STRING }`, `Address ::= SEQUENCE { codePage INTEGER, address OCTET STRING }`.
- **Code Page:** A code page is a table of values that describes a character set. In this context, it's an integer that specifies which character encoding to use to interpret the `name` or `address` byte string. For example, code page 1 corresponds to `ISO/IEC 8859-1` (Latin-1), a common encoding for Western European languages. Without the code page, non-ASCII characters could be misinterpreted. The mapping from code page values to character sets is defined in Chapter 4 of the Data Dictionary.
- **Intended Action:**
    1.  Create `name.proto` and `address.proto` in `datadictionary/v1`.
    2.  Define `Name` and `Address` messages with `code_page` and `name`/`address` fields.
    3.  Update all fields of type `Name` and `Address` in the affected proto files to use these messages.

### `Datef` type

- **Affected Files:** `vu/v1/activities.proto`, `vu/v1/technical_data.proto`
- **Deviation:** Fields representing `Datef` are `google.protobuf.Timestamp`.
- **Problem:** The `Datef` data type (Data Dictionary, Section 2.57) is an `OCTET STRING (SIZE(4))` representing a BCD encoded date `yyyymmdd`. `Timestamp` is not the correct type.
- **Specification:** Data Dictionary, Section 2.57, `Datef`.
- **ASN.1 Definition:** `Datef ::= OCTET STRING (SIZE(4))`.
- **Intended Action:** Change the type to a message with `year`, `month`, `day` fields, or a `string` to hold the raw BCD. A message is preferred for better semantics.

### `GeoCoordinates` type

- **Affected Files:** `card/v1/gnss_place_auth_record.proto`, `vu/v1/activities.proto`, `vu/v1/gnss_places.proto`
- **Deviation:** `longitude` and `latitude` are separate fields.
- **Problem:** They should be in a `GeoCoordinates` message, as `GeoCoordinates` is a reusable sequence.
- **Specification:** Data Dictionary, Section 2.76, `GeoCoordinates`.
- **ASN.1 Definition:** `GeoCoordinates ::= SEQUENCE { latitude INTEGER, longitude INTEGER }`.
- **Intended Action:**
    1. Create a new file `proto/wayplatform/connect/tachograph/datadictionary/v1/geo_coordinates.proto`.
    2. Define a new message `GeoCoordinates` in this file with `latitude` and `longitude` fields.
    3. Update the affected proto files to import and use the new `GeoCoordinates` message.

### `GNSSAccuracy` type

- **Affected Files:** `card/v1/gnss_place_auth_record.proto`, `vu/v1/activities.proto`, `vu/v1/gnss_places.proto`
- **Deviation:** The comment for the `gnss_accuracy` field is incorrect.
- **Problem:** The comment states that `GNSSAccuracy` is an `OCTET STRING(SIZE(1))`, but the data dictionary (Section 2.77) defines it as `INTEGER (1..100)`.
- **Specification:** Data Dictionary, Section 2.77, `GNSSAccuracy`.
- **ASN.1 Definition:** `GNSSAccuracy ::= INTEGER (1..100)`.
- **Intended Action:** Correct the comment for the `gnss_accuracy` field.

### `PositionAuthenticationStatus` type

- **Affected Files:** `card/v1/gnss_place_auth_record.proto`, `card/v1/gnss_places_authentication.proto`, `card/v1/places_authentication.proto`
- **Deviation:** The `authentication_status` field is an `int32`.
- **Problem:** The `PositionAuthenticationStatus` data type (Data Dictionary, Section 2.117a) is an enum.
- **Specification:** Data Dictionary, Section 2.117a, `PositionAuthenticationStatus`.
- **ASN.1 Definition:** `PositionAuthenticationStatus ::= INTEGER { notAvailable(0), authenticated(1), notAuthenticated(2), authenticationCorrupted(3) } (0..255)`.
- **Intended Action:** Update the affected proto files to import and use the existing `PositionAuthenticationStatus` enum from `datadictionary/v1`.

### `ControlType` type

- **Affected Files:** `card/v1/control_activity_data.proto`, `card/v1/controller_activity_data.proto`, `vu/v1/overview.proto`
- **Deviation:** The `control_type` field is `bytes`.
- **Problem:** `ControlType` (2.53) is a bitmask.
- **Specification:** Data Dictionary, Section 2.53, `ControlType`.
- **ASN.1 Definition:** `ControlType ::= OCTET STRING (SIZE(1))`.
- **Intended Action:** Create a new message `ControlType` in `datadictionary/v1` with boolean fields for each flag and use it in the affected files.

### `OperationType` type

- **Affected Files:** `card/v1/load_unload_operations.proto`, `vu/v1/activities.proto`
- **Deviation:** The `OperationType` enum is defined inside the `Record` message in `load_unload_operations.proto`.
- **Problem:** It should be a shared type in `datadictionary/v1`.
- **Specification:** Data Dictionary, Section 2.114a, `OperationType`.
- **ASN.1 Definition:** `OperationType ::= INTEGER { load(1), unload(2), simultaneous(3) } (0..255)`.
- **Intended Action:** Move the `OperationType` enum to `datadictionary/v1`.

### `RegionNumeric` type

- **Affected Files:** `card/v1/places.proto`, `vu/v1/activities.proto`
- **Deviation:** The `daily_work_period_region` and `region` fields are typed as `int32`.
- **Problem:** The `RegionNumeric` data type (Data Dictionary, Section 2.122) is an `OCTET STRING (SIZE (1))`. Using `int32` is incorrect for an octet string. It should be `bytes`.
- **Specification:** Data Dictionary, Section 2.122, `RegionNumeric`.
- **ASN.1 Definition:** `RegionNumeric ::= OCTET STRING (SIZE (1))`.
- **Intended Action:** Change the type of `daily_work_period_region` and `region` fields to `bytes` (or `bytes` of size 1).

### `full_card_number.proto`

- **Deviation:** The fields for `driverIdentification` and `ownerIdentification` are swapped in the `FullCardNumber` message. The `DriverIdentification` message incorrectly contains `consecutive_index`, `replacement_index`, and `renewal_index`, while the `OwnerIdentification` message is missing them.
- **Problem:** The protobuf message does not correctly represent the ASN.1 specification for `CardNumber`. This will lead to incorrect parsing and serialization of card numbers for drivers and owners.
- **Specification:** Data Dictionary, Section 2.26, `CardNumber`.
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
    1.  In `full_card_number.proto`, move the `consecutive_index`, `replacement_index`, and `renewal_index` fields from the `DriverIdentification` message to the `OwnerIdentification` message.
    2.  Update the comments in both messages to accurately reflect the ASN.1 structure.
    3.  The `driver_identification` field in `CardNumber` should only contain the `identification` field.

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
