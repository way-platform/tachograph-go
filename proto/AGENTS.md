# Protobuf Schema Design Guidelines

These guidelines are specific to the `proto/` directory and MUST be followed when modifying or creating protobuf definitions.

## Principles

- **AEP Alignment**: Follow https://aep.dev.
- **Tagged Unions**: Prefer `type` enum over `oneof`.
- **Timestamps**: Use `google.protobuf.Timestamp`.
- **No Unsigned Integers**: Avoid unsigned integers due to limited support in some languages.
- **Bitfields**: Use `bool` flags for single bits (e.g., `crew` not `is_crew`).
- **Superset Messages**: Use Gen2 structure if it's a strict superset of Gen1 (e.g., `FullCardNumberAndGeneration`).
- **Special Types**:
    - `StringValue`: For code-paged strings (includes encoding).
    - `Date`: For BCD dates (year, month, day).
    - `bytes`: For `OCTET STRING`, even if it's a single byte.
- **Flattening**: Flatten `SEQUENCE OF SET` into a single chronological `repeated` array (e.g., `EventsData.events`).
- **Unrecognized Enums**: Include `unrecognized_<field_name>` (`int32`) fields to capture raw protocol values for unrecognized enum members, ensuring lossless round-trips.

## Signature Handling in Proto
**Combine EF Signature**: Model signed EFs by including a `bytes signature` field within the EF's protobuf message.
- **Permissive Policy**: Include a `signature` field unless the regulation explicitly states the EF will NOT have one.

## Generation Tracking
- **EF Level**: Each EF message with generational differences MUST include `dd.v1.Generation generation`.
- **Source**: Derived from TLV tag appendix (Bit 1).

## Documentation Requirements

### Source Material Only
All documentation in `.proto` files must be self-contained and based **only** on regulations. **NEVER** reference `AGENTS.md` or internal project policies.

### ASN.1 definitions
All messages/fields mirroring Data Dictionary types MUST include the ASN.1 definition.
**Format:**
```protobuf
// Summary of purpose.
//
// See Data Dictionary, Section X.Y.
//
// ASN.1 Definition:
//
//     MyType ::= SEQUENCE { ... }
```

### File Structure Diagrams
For signed EFs, include a diagram in the message comment showing the `Signature` block from Appendix 2.
```protobuf
// File Structure:
//
//     EF Identification
//     └─CardIdentification
//     Signature
```

## Package Structure
- `wayplatform.connect.tachograph.v1`: Top-level (File).
- `wayplatform.connect.tachograph.vu.v1`: VU-specific.
- `wayplatform.connect.tachograph.card.v1`: Card-specific.
- `wayplatform.connect.tachograph.datadictionary.v1`: Shared types.
