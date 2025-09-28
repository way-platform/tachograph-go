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

## **Cross-cutting Deviations**

### `Name` and `Address` types

- **Affected Files:** `card/v1/driving_licence_info.proto`, `card/v1/identification.proto`, `vu/v1/overview.proto`, `vu/v1/technical_data.proto`
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

- **Affected Files:** `card/v1/identification.proto`, `vu/v1/activities.proto`, `vu/v1/technical_data.proto`
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
- **Deviation:** The `daily_work_period_region` and `region` fields lack context on how to be interpreted.
- **Problem:** The `RegionNumeric` data type (Data Dictionary, Section 2.122) is an `OCTET STRING (SIZE (1))` that represents a numeric code for a region. The meaning of this code is country-specific and the lists of codes are maintained externally.
- **Specification:** Data Dictionary, Section 2.122, `RegionNumeric`.
- **ASN.1 Definition:** `RegionNumeric ::= OCTET STRING (SIZE (1))`.
- **Interpretation:** This is a code that identifies a region within a country. The meaning of the code depends on the country. For Generation 1, the data dictionary provides a list of values for Spain. For Generation 2, the codes are maintained on the website of the laboratory appointed to carry out interoperability testing: [dtlab.jrc.ec.europa.eu](https://dtlab.jrc.ec.europa.eu/).
- **Intended Action:** Add a comment to the `daily_work_period_region` and `region` fields with the interpretation context and the link to the external resource.