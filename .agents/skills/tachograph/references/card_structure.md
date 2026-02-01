# Card Parsing Structure

Comprehensive guidance for Tachograph Card Files (Driver, Workshop, Control, Company).

## Card File Structure (TLV)

Tachograph card files use a **TLV (Tag-Length-Value)** format organized into a **Dedicated File (DF)** and **Elementary File (EF)** hierarchy.

### TLV Tag Structure
The 3-byte TLV tag encodes the File ID (FID) and the DF context.
`[Byte 0-1: File ID][Byte 2: Appendix/Generation]`

**Appendix Byte (Byte 2):**
- `0x00`: Gen1 Data
- `0x01`: Gen1 Signature
- `0x02`: Gen2 Data
- `0x03`: Gen2 Signature

## Parsing Flow (Two-Pass)

1.  **Pass 1: TLV Parsing**
    - Splits binary into `TlvRecord`s.
    - Extracts tag, file type, generation (from appendix), and value.

2.  **Pass 2: Semantic Parsing**
    - Routes `TlvRecord` to the appropriate DF message (Gen1 or Gen2).
    - Parses EF-specific data (e.g., Identification, Activity).
    - Attaches signatures to their EFs (if present).

## Signature Handling

- **Unsigned EFs**: EF_ICC, EF_IC, Certificates, EF_Card_Download.
    - **Parsing**: Silently ignore signatures if found (be liberal).
    - **Marshalling**: NEVER write signatures for these files (be strict).
- **Signed EFs**: Most others (e.g., Identification, Activity).
    - **Structure**: `[Data TLV] [Signature TLV]`.
    - **Model**: The data model for the EF should include a signature field.

## Generation-Specific Types

**Principle**: Split data types by generation if the binary layout or size differs.

- **Example**: `PlaceRecord` (10 bytes) vs `PlaceRecordG2` (21 bytes).
- **Don't Split**: If the structure is identical or purely additive at the end without layout changes (superset).

## Marshalling Flow

Respect the DF hierarchy:
1.  **Common Files (MF)**: EF_ICC, EF_IC (Appendix 0x00).
2.  **Gen1 DF**: Data (0x00) and Signatures (0x01).
3.  **Gen2 DF**: Data (0x02) and Signatures (0x03).
