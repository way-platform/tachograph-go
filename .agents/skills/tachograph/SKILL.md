---
name: tachograph
description: specialized knowledge for parsing/marshalling EU tachograph files. Use when needing to understand the .DDD/.V1B file format, TLV/TV protocols, ASN.1 mappings, or EU regulation compliance.
---

# EU Tachograph Domain Knowledge

This skill provides the domain expertise required to work with EU Tachograph files (Driver Card & Vehicle Unit).

## Core Concepts

- **Regulation First**: All implementation must derive directly from EU regulations (Data Dictionary).
- **Binary Fidelity**: Support perfect round-trip parsing/marshalling (byte-for-byte exactness).
- **Type Safety**: Use specific types for specific generations/formats.

## Workflows

### 1. Card File Structure (TLV)
- **Architecture**: See [card_structure.md](references/card_structure.md) for TLV, DF/EF structure, and signature handling.
- **Pattern**: Use a "Two-Pass" parsing strategy (Raw TLV -> Semantic Structure).

### 2. Regulation & ASN.1
- **Regulation**: Consult [regulation.md](references/regulation.md) to locate specific chapters (Data Dictionary is critical).
- **Definitions**: Use ASN.1 definitions from the regulation as the source of truth for all data types.

### 3. Generations & Versions
- **Details**: See [versions.md](references/versions.md) for a comprehensive overview of:
    - **Generation 1**: Legacy (SHA-1, RSA-1024).
    - **Generation 2 V1**: Smart Tachograph (SHA-256, ECC).
    - **Generation 2 V2**: Smart Tachograph with GNSS/OSNMA support.
    - **Protocol Matrix**: TLV (Card) vs TV (VU) tag ranges and structures.
- **Detection**: Use tag-based detection logic for both Card and VU files.

## Best Practices for Implementation

### Data Parsing
1.  **ASN.1 Compliance**: Every parser function should reference the ASN.1 definition from the spec.
2.  **No Magic Numbers**: Define constants for all binary offsets and lengths.
3.  **Exact Lengths**: Validate fixed-structure lengths exactly (`==`), not loosely (`>=`).
4.  **Nil Handling**: The protocol has no concept of "null", only "default values" (e.g., spaces for strings, 0 for numbers). Parsing functions should handle missing or zero-length input safely.

### Data Modeling
1.  **Integers**: Use signed integers (`int32`/`int64`) for general numbers; the protocol rarely uses unsigned types in a way that maps cleanly to `uint` in all languages.
2.  **Bitflags**: Prefer boolean fields over integer bitmasks for clarity.
3.  **Signatures**: Model signed EFs by including a signature field (unless Regulation forbids).