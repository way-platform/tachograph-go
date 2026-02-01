# Tachograph Generations and Versions

This document provides a comprehensive overview of tachograph generations and versions, their key differences, and capabilities.

## Overview

The tachograph system has evolved through multiple generations:
- **Generation 1**: Digital tachograph (legacy system)
- **Generation 2 Version 1**: Smart tachograph (enhanced security)
- **Generation 2 Version 2**: Smart tachograph with GNSS/OSNMA support

## Protocol and Generation Matrix

| Generation          | Card Protocol (TLV)      | VU Protocol (TV)          | TREP Range      | Card Tags              | Security                 |
| ------------------- | ------------------------ | ------------------------- | --------------- | ---------------------- | ------------------------ |
| **Generation 1**    | TLV with Gen1 structures | TV with simple structures | 0x01-0x05       | 0x0002-0x051E          | Part A (SHA-1, RSA-1024) |
| **Generation 2 V1** | TLV with record arrays   | TV with record arrays     | 0x21-0x25       | 0x0002-0x051E + Gen2   | Part B (SHA-256, ECC)    |
| **Generation 2 V2** | TLV with GNSS extensions | TV with GNSS + 0x7600     | 0x00, 0x31-0x35 | 0x0002-0x051E + Gen2V2 | Part B + OSNMA           |

## Protocol-Specific Characteristics

### Card Protocol (TLV)

- **Generation 1**:
  - Tags: `[FID][00]`
  - Structure: Simple, fixed-size structures.
- **Generation 2 V1**:
  - Tags: `[FID][01]`
  - Structure: Record arrays with headers (`RecordType` + `RecordSize` + `NumberOfRecords`).
- **Generation 2 V2**:
  - Tags: `[FID][02]`
  - Structure: Enhanced record arrays, often including GNSS capabilities.

### Vehicle Unit Protocol (TV)

- **Generation 1**:
  - Tags: `0x76[01-05]`
  - Structure: Direct field access.
- **Generation 2 V1**:
  - Tags: `0x76[21-25]`
  - Structure: All fields wrapped in record arrays.
- **Generation 2 V2**:
  - Tags: `0x76[00]`, `0x76[31-35]`
  - Structure: Extended record arrays with GNSS integration. Tag `0x7600` is unique for version identification.

## Generation Detection Logic

### Card File (TLV)
Detect generation by examining the **3rd byte** of the TLV tag:
- `0x00` = Generation 1
- `0x01` = Generation 2 V1
- `0x02` = Generation 2 V2

### Vehicle Unit File (TV)
Detect generation by examining the **TREP** (Transfer Entry Parameter) byte (low byte of the tag):
- `0x01-0x05` = Generation 1
- `0x21-0x25` = Generation 2 V1
- `0x00`, `0x31-0x35` = Generation 2 V2

## Interoperability

| Component               | Gen 1 VU         | Gen 2 V1 VU      | Gen 2 V2 VU      |
| ----------------------- | ---------------- | ---------------- | ---------------- |
| **Gen 1 Cards**         | ✅ Native        | ✅ Compatible\*  | ✅ Compatible\*  |
| **Gen 2 Cards**         | ❌ Not supported | ✅ Native        | ✅ Native        |
| **Gen 1 Motion Sensor** | ✅ Native        | ❌ Not supported | ❌ Not supported |
| **Gen 2 Motion Sensor** | ❌ Not supported | ✅ Native        | ✅ Native        |

_\* Can be disabled by workshop settings._

## Security Mechanisms

- **Part A (Gen 1)**: SHA-1, RSA-1024. Basic certificate validation.
- **Part B (Gen 2)**: SHA-256, ECC-256/384, RSA-2048. Enhanced validation, advanced tamper detection.
- **OSNMA (Gen 2 V2)**: Open Service Navigation Message Authentication for GNSS position verification.