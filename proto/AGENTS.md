# Protobuf Schemas

This document provides guidance on how the protobuf schemas are developed for this project.

Our data model is based on the protobuf schemas defined in the [`proto`](.) directory. We use protobuf edition 2023, which makes all fields optional by default. Our generated code uses the opaque API, requiring accessors and setters.

## Design Principles

- Align high-level conventions to the AIP (https://aip.dev) design system.
- Prefer tagged unions (e.g., a `type` enum) over `oneof` for ergonomic reasons.
- Use `google.protobuf.Timestamp` for timestamp fields.
- Avoid unsigned integers due to limited support in some languages.
- **Use boolean flags for single-bit fields within bitfields:** When a field represents a single bit within a larger bitfield structure (e.g., bits within `ActivityChangeInfo`), use a `bool` field instead of a two-value enum. This provides better ergonomics for API consumers and aligns with AIP style guidelines. The field name should describe the `true` state without an `is_` prefix (e.g., `crew` for crew mode, `inserted` for card insertion status). However, keep an enum when: (1) the field occupies an entire byte with reserved/RFU values for future extension (e.g., `PositionAuthenticationStatus` with values 0, 1, and 2-255 RFU), or (2) the field represents a categorical choice rather than a boolean state (e.g., `CardSlotNumber` distinguishing driver vs co-driver slot). Always document the protocol bit/byte values in field comments.
- **Use Superset Messages for Generational Differences:** When a Generation 2 data structure is a clear superset of its Generation 1 equivalent (e.g., `FullCardNumberAndGeneration` contains `FullCardNumber`), we will use the Gen2 superset message for all related fields. This simplifies the API by providing a single, unified field for consumers. When parsing Gen1 data, the additional Gen2 fields in the message will be left unset.
- **Use `StringValue` for special strings:** Many string-like types in the data dictionary are not simple UTF-8. This includes `IA5String` (which may have padding) and complex `SEQUENCE` types (like `Name` and `Address`) that contain a code page. To ensure lossless round-trips, these fields **must** use `datadictionary.v1.StringValue`. This message provides the original `encoded` bytes, the `Encoding` enum (which includes a value for `IA5`), and a `decoded` field for display.
- **Use `Date` for BCD dates:** The `Datef` data type (DD 2.57) is an `OCTET STRING (SIZE(4))` representing a BCD-encoded `yyyymmdd` date. Any field corresponding to this type **must** use the `datadictionary.v1.Date` message, which provides decoded `year`, `month`, and `day` fields.
- **Use `bytes` for `OCTET STRING`:** To maintain semantic fidelity with the ASN.1 specification, fields defined as `OCTET STRING` should be represented as `bytes` in Protobuf, even if they are single-byte values that could be losslessly stored in an `int32`. This makes it clear to consumers that the data is a raw byte string, not necessarily a number.
- **Combine EF Signature:** The ASN.1 definitions in the Data Dictionary describe the content of an Elementary File (EF) itself and will not include a signature. However, the physical card file structure may include signature data blocks for certain EFs. For usability, our policy is to model this by including a `bytes signature` field within the protobuf message that represents the EF. This provides a complete, self-contained representation of the signed data block when present. **Permissive Signature Policy:** Unless the source material explicitly states that a specific EF will not have a signature, we include a `signature` field in the protobuf message. The signature field should be documented to indicate that it contains signature data from the following file block, if tagged as a signature for this EF according to the card file format specification (Appendix 2). This approach ensures compatibility with real-world card data while maintaining clear documentation of the signature's source.
- **Flatten Type-Grouped Structures for API Usability:** When ASN.1 specifications define data as "a sequence of sets grouped by type" (e.g., `SEQUENCE OF SET` structures), prefer flattening these into a single chronological array in the protobuf API. The type information is preserved in each record's type field, allowing consumers to filter by type while maintaining temporal order. This approach prioritizes API usability over strict ASN.1 structural fidelity. Examples include `EventsData.events` and `FaultsData.faults`. The protobuf message comment should explain this design choice and reference the original ASN.1 structure.
- **Preserve Unrecognized Enum Values:** When parsing binary data that contains enum values not defined in our protobuf schemas (due to incomplete specifications, vendor extensions, or future protocol versions), preserve the original raw values to maintain perfect data fidelity. For any message containing enum fields that might encounter unrecognized values, include corresponding `unrecognized_<field_name>` fields of type `int32` to capture the raw protocol values. This ensures lossless round-trip operations and allows consumers to inspect and handle unknown enum values. The enum field itself should be set to the appropriate `UNRECOGNIZED` value when an unknown protocol value is encountered.
- **Track Generation at EF Level:** Generation information is fundamentally a property of individual Elementary Files (EFs), not entire card files. Each EF message that has generational differences **must** include a `dd.v1.Generation generation` field to track which generation the parsed data represents. This generation information comes from bit 1 of the TLV tag's appendix byte during parsing and must be preserved for correct marshalling. The generation field enables generation-aware parsing and marshalling of EF content, ensuring that the correct binary format is used when writing data back to card format. For EFs that are identical across generations, the generation field may be omitted.

## ASN.1 Documentation

All messages, fields, and enums corresponding to a data dictionary type **must** be documented with a comment containing the original ASN.1 definition.

**Source Material Only:** All comments and documentation within `.proto` files must be self-contained and based on first principles from the source regulations. **Do not** reference internal project documents like `AGENTS.md` or internal policies. The rationale for a design choice should be evident from the regulatory context provided in the comment itself.

The comment **must** follow this structure:

1.  A brief summary of the element's purpose.
2.  A blank line.
3.  A reference to the Data Dictionary section (e.g., `See Data Dictionary, Section 2.53.`).
4.  A blank line.
5.  The heading `ASN.1 Definition:`. Use generational headings (e.g., `ASN.1 Definition (Gen1):`) if the definition varies.
6.  A blank line.
7.  The full, indented ASN.1 definition.

**Example:**

```protobuf
// Represents the activities carried out during a control.
//
// See Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//     ControlType ::= OCTET STRING (SIZE(1))
message ControlType {
  // ...
}
```

## File Structure Documentation

To provide clear, at-a-glance context for whether a message should contain a `signature` field, any message that represents a signed Elementary File (EF) **must** include a file structure diagram in its message-level comment.

This diagram should be a small, focused snippet from the tables in Appendix 2 (e.g., TCS_158), illustrating the EF's content and the explicit `Signature` block that follows it. This makes the "Combine EF Signature" design principle verifiable directly within the schema.

**Example (for a signed EF):**

```protobuf
// Represents data from EF_Identification for a workshop card.
//
// The file structure, including the signature, is defined in Appendix 2,
// table TCS_158.
//
// File Structure:
//
//     EF Identification
//     ├─CardIdentification
//     └─WorkshopCardHolderIdentification
//     Signature
//
message Identification {
  // ... fields ...
  bytes signature = 99; // This field is present because of the structure above.
}
```

**Example (for an unsigned EF):**

```protobuf
// Represents data from EF_Driving_Licence_Info.
//
// The file structure is defined in Appendix 2, table TCS_154. Note the
// absence of a `Signature` block.
//
// File Structure:
//
//     EF Driving_Licence_Info
//     └─CardDrivingLicenceInformation
//
message DrivingLicenceInfo {
  // ... fields ...
  // NO signature field is present.
}
```

## Package Structure

### [`wayplatform.connect.tachograph.v1`](./wayplatform/connect/tachograph/v1)

Top-level package for all tachograph data.

- `wayplatform.connect.tachograph.v1.File`: Represents any type of tachograph file.

### [`wayplatform.connect.tachograph.vu.v1`](./wayplatform/connect/tachograph/vu/v1)

Package for vehicle unit (VU) data.

- `wayplatform.connect.tachograph.vu.v1.VehicleUnitFile`: Represents a VU file.

### [`wayplatform.connect.tachograph.card.v1`](./wayplatform/connect/tachograph/card/v1)

Package for tachograph card data.

- `wayplatform.connect.tachograph.card.v1.DriverCardFile`: Represents a driver card file.
- `wayplatform.connect.tachograph.card.v1.RawCardFile`: Represents a generic card file (TLV records).
- Each EF (elementary file) has a corresponding top-level message named after it.

### [`wayplatform.connect.tachograph.datadictionary.v1`](./wayplatform/connect/tachograph/datadictionary/v1)

Package for shared types from the data dictionary ([03-data-dictionary.md](../../docs/regulation/chapters/03-data-dictionary.md)).

- Contains types used across multiple card EFs or VU data transfers.
- Types used in only a single context should be defined inline within that message.
