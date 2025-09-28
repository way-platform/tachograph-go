# Log: Protobuf Schema Deviation Fixes

**Date:** 2025-09-28

## 1. Summary

This log entry documents the systematic review and resolution of all issues outlined in the `2025-09-27_proto_review_deviations.md` document. The effort focused on correcting data types, improving data modeling to better reflect the ASN.1 specification, and enhancing inline documentation to provide better context for future development.

## 2. Resolved Deviations

The following fixes were implemented across the protobuf schemas.

### 2.1. `ExtendedSerialNumber` Type

-   **Problem**: The `serial_number` field was `uint32`, violating the project guideline to avoid unsigned integers.
-   **Fix**: The field type was changed to `int64`, which can safely represent the entire `INTEGER(0..2^32-1)` range while adhering to our standards.

### 2.2. `CardNumber` Modeling

-   **Problem**: The `CardNumber` type was incorrectly modeled as a primitive `string`, losing the complex `CHOICE { SEQUENCE }` structure defined in the data dictionary.
-   **Fix**: After discussing the trade-offs between DRY and context colocation, we adopted an inlining strategy. The full `CardNumber` structure, including nested `DriverIdentification` and `OwnerIdentification` messages, was inlined directly into the `FullCardNumber` and `Identification.Card` messages. This provides developers with the full context of the data structure in one place and includes detailed comments explaining the ASN.1 hierarchy and the use of an external discriminant (`card_type`).

### 2.3. `VehicleRegistrationNumber` Modeling

-   **Problem**: The `VehicleRegistrationNumber` was incorrectly modeled as a `string`.
-   **Fix**: We clarified that the underlying ASN.1 type is a `SEQUENCE` of a `codePage` and the number string, not a `CHOICE`. Accordingly, the field's type was changed to `StringValue`, which is the correct, existing message for representing a string with its associated encoding information. The comments were updated with the correct `SEQUENCE` definition.

### 2.4. `technical_data.proto` Refinements

-   **`card_structure_version`**: Corrected from `bytes` to a nested `CardStructureVersion` message with `major` and `minor` fields, accurately representing the 2-byte BCD structure.
-   **`consent_status`**: Corrected from `int32` to `bool` to match the ASN.1 `BOOLEAN` type.
-   **`software_version`**: Confirmed that `StringValue` is the correct type, as the `Encoding` enum is designed to handle `IA5String` types. The initial deviation was based on a misunderstanding of the `StringValue` message's purpose.

### 2.5. `Datef` Type Implementation

-   **Problem**: Fields representing the BCD-encoded `Datef` type were incorrectly using `google.protobuf.Timestamp`.
-   **Fix**: A new, reusable `Date` message was created in the data dictionary with `year`, `month`, and `day` fields and detailed comments on BCD decoding. All incorrect `Timestamp` usages across `identification.proto`, `activities.proto`, and `technical_data.proto` were replaced with this new, correct type. The `proto/AGENTS.md` file was also updated with guidance on this new type.

### 2.6. `GeoCoordinates` Type Implementation

-   **Problem**: `latitude` and `longitude` were modeled as separate primitive fields.
-   **Fix**: A new, reusable `GeoCoordinates` message was created in the data dictionary. All affected files (`gnss_place_auth_record.proto`, `activities.proto`, `gnss_places.proto`) were refactored to use this new, more cohesive message.

### 2.7. Comment-Only Fixes

-   **`GNSSAccuracy`**: The ASN.1 definition in the comments for this field was corrected from `OCTET STRING(SIZE(1))` to `INTEGER (1..100)` in all affected files.
-   **`RegionNumeric`**: The missing explanatory comment regarding the country-specific nature of the codes and the link to the JRC website was added to the relevant fields in `places.proto` and `activities.proto`.

## 3. Conclusion

All deviations identified in the `2025-09-27` review have been addressed. The protobuf schemas are now significantly more accurate, type-safe, and better documented, providing a solid foundation for future development.