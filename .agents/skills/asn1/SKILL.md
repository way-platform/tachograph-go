---
name: asn1
description: Reference for ASN.1 notation, BER (Basic Encoding Rules), and DER (Distinguished Encoding Rules). Use when decoding/encoding binary data, debugging ASN.1 structures, or understanding PKCS standards.
---

# ASN.1, BER, and DER

This skill provides reference material for Abstract Syntax Notation One (ASN.1) and its encoding rules (BER/DER).

## Quick Reference

### Universal Tags

| Tag (Hex) | Type |
| :--- | :--- |
| `01` | BOOLEAN |
| `02` | INTEGER |
| `03` | BIT STRING |
| `04` | OCTET STRING |
| `05` | NULL |
| `06` | OBJECT IDENTIFIER |
| `10` | SEQUENCE / SEQUENCE OF |
| `11` | SET / SET OF |
| `13` | PrintableString |
| `14` | T61String |
| `16` | IA5String |
| `17` | UTCTime |

### Identifier Octet (First Byte)

Bits 8-7: Class
- `00`: Universal
- `01`: Application
- `10`: Context-specific
- `11`: Private

Bit 6: PC (Primitive/Constructed)
- `0`: Primitive (Value is directly in contents)
- `1`: Constructed (Contents are nested TLV elements)

Bits 5-1: Tag Number (if < 31)

## Detailed Guide

For a comprehensive guide on notation, encoding rules, and specific type handling, read the bundled reference:

`references/guide.md`

### Contents of the Guide:
1.  **Introduction**: Overview of OSI, ASN.1, BER, DER.
2.  **ASN.1 Notation**: Simple types, structured types, tagging.
3.  **BER**: Primitive vs. Constructed, Definite vs. Indefinite length.
4.  **DER**: subset of BER for unique encoding (Canonical).
5.  **Type-Specific Encodings**: Detailed breakdown of how to encode INTEGER, BIT STRING, SEQUENCE, etc.
6.  **Example**: Walkthrough of an X.500 Name encoding.

## Common Operations

### Deciphering a Tag
If you encounter byte `0x30`:
- Binary: `0011 0000`
- Class: `00` (Universal)
- P/C: `1` (Constructed)
- Tag: `10000` (16 -> SEQUENCE)

If you encounter byte `0xA0`:
- Binary: `1010 0000`
- Class: `10` (Context-specific)
- P/C: `1` (Constructed)
- Tag: `00000` (0 -> [0])

### Variable Lengths
- **Short Form**: `0x00` - `0x7F` (0-127 bytes).
- **Long Form**: Bit 8 is 1. Bits 7-1 = number of length bytes following.
  - Example: `0x81 0x80` -> Length is in next 1 byte (`0x81`), value is `0x80` (128).
  - Example: `0x82 0x01 0x00` -> Length is in next 2 bytes (`0x82`), value is `0x0100` (256).