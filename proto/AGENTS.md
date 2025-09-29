# Protobuf Schemas

This document provides guidance on how the protobuf schemas are developed for this project.

Our data model is based on the protobuf schemas defined in the [`proto`](.) directory. We use protobuf edition 2023, which makes all fields optional by default. Our generated code uses the opaque API, requiring accessors and setters.

## Design Principles

- Align high-level conventions to the AIP (https://aip.dev) design system.
- Prefer tagged unions (e.g., a `type` enum) over `oneof` for ergonomic reasons.
- Use `google.protobuf.Timestamp` for timestamp fields.
- Avoid unsigned integers due to limited support in some languages.
- **Use Superset Messages for Generational Differences:** When a Generation 2 data structure is a clear superset of its Generation 1 equivalent (e.g., `FullCardNumberAndGeneration` contains `FullCardNumber`), we will use the Gen2 superset message for all related fields. This simplifies the API by providing a single, unified field for consumers. When parsing Gen1 data, the additional Gen2 fields in the message will be left unset.
- **Use `StringValue` for special strings:** Many string-like types in the data dictionary are not simple UTF-8. This includes `IA5String` (which may have padding) and complex `SEQUENCE` types (like `Name` and `Address`) that contain a code page. To ensure lossless round-trips, these fields **must** use `datadictionary.v1.StringValue`. This message provides the original `encoded` bytes, the `Encoding` enum (which includes a value for `IA5`), and a `decoded` field for display.
- **Use `Date` for BCD dates:** The `Datef` data type (DD 2.57) is an `OCTET STRING (SIZE(4))` representing a BCD-encoded `yyyymmdd` date. Any field corresponding to this type **must** use the `datadictionary.v1.Date` message, which provides decoded `year`, `month`, and `day` fields.
- **Use `bytes` for `OCTET STRING`:** To maintain semantic fidelity with the ASN.1 specification, fields defined as `OCTET STRING` should be represented as `bytes` in Protobuf, even if they are single-byte values that could be losslessly stored in an `int32`. This makes it clear to consumers that the data is a raw byte string, not necessarily a number.
- **Combine EF Signature:** For Elementary Files (EFs) that are followed by a signature block on the card, the signature is included as a `signature` field within the corresponding protobuf message. This provides a complete, self-contained representation of the signed data block, even though the signature is technically separate from the EF content in the raw card file structure.

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