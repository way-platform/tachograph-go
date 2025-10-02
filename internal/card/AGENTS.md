# Card Package Guidelines

This document provides comprehensive guidance for implementing and maintaining card file parsing in the `internal/card` package.

## Overview

The `internal/card` package handles parsing and marshalling of tachograph card files (`.DDD` files). These are binary data dumps from tachograph cards using a TLV (Tag-Length-Value) format organized into a hierarchical Dedicated File (DF) and Elementary File (EF) structure.

**Goals:**

- Full alignment with the EU digital tachograph regulation
- Full binary roundtrip parsing with no data loss
- Easy-to-use and high-fidelity protobuf data model
- Support for Generation 1 and Generation 2 card formats

## Regulation References

Key regulation chapters for card file implementation:

- **[03-data-dictionary.md](../../docs/regulation/chapters/03-data-dictionary.md)**: Critical for data parsing. Contains ASN.1 specifications.
- **[05-tachograph-cards-file-structure.md](../../docs/regulation/chapters/05-tachograph-cards-file-structure.md)**: Essential for understanding the DF/EF hierarchy and file structure.
- **[12-card-downloading.md](../../docs/regulation/chapters/12-card-downloading.md)**: Essential for card data format and TLV encoding.
- **[16-common-security-mechanisms.md](../../docs/regulation/chapters/16-common-security-mechanisms.md)**: Essential for certificates and signatures.

**IMPORTANT**: Always read [../../docs/asn-1.md](../../docs/asn-1.md) before working with ASN.1 data.

## Card File Structure Overview

Tachograph card files use a **TLV (Tag-Length-Value)** format where data is organized hierarchically into **Dedicated Files (DFs)** containing **Elementary Files (EFs)**.

### Physical vs. Logical Structure

**Physical Structure (TLV Format):**

```
[Tag: 3 bytes][Length: 2 bytes][Value: N bytes]
[Tag: 3 bytes][Length: 2 bytes][Value: M bytes]
...
```

**Logical Structure (DF/EF Hierarchy):**

```
Card File
├─ Master File (MF)
│  ├─ EF_ICC (common to all card types)
│  └─ EF_IC (common to all card types)
│
├─ DF Tachograph (Generation 1 application)
│  ├─ EF_Application_Identification
│  ├─ EF_Card_Certificate
│  ├─ EF_CA_Certificate
│  ├─ EF_Identification
│  ├─ EF_Card_Download
│  ├─ EF_Driving_Licence_Info
│  ├─ EF_Events_Data
│  ├─ EF_Faults_Data
│  ├─ EF_Driver_Activity_Data
│  ├─ EF_Vehicles_Used
│  ├─ EF_Places
│  ├─ EF_Current_Usage
│  ├─ EF_Control_Activity_Data
│  └─ EF_Specific_Conditions
│
└─ DF Tachograph_G2 (Generation 2 application)
   ├─ EF_Application_Identification
   ├─ EF_CardMA_Certificate
   ├─ EF_CardSignCertificate
   ├─ EF_CA_Certificate
   ├─ EF_Link_Certificate
   ├─ EF_Identification
   ├─ EF_Card_Download
   ├─ EF_Driving_Licence_Info
   ├─ EF_Events_Data
   ├─ EF_Faults_Data
   ├─ EF_Driver_Activity_Data
   ├─ EF_Vehicles_Used
   ├─ EF_Places
   ├─ EF_Current_Usage
   ├─ EF_Control_Activity_Data
   ├─ EF_Specific_Conditions
   ├─ EF_VehicleUnits_Used (Gen2 only)
   ├─ EF_GNSS_Places (Gen2 only)
   ├─ EF_Application_Identification_V2 (Gen2v2 only)
   ├─ EF_Places_Authentication (Gen2v2 only)
   ├─ EF_GNSS_Places_Authentication (Gen2v2 only)
   ├─ EF_Border_Crossings (Gen2v2 only)
   ├─ EF_Load_Unload_Operations (Gen2v2 only)
   ├─ EF_Load_Type_Entries (Gen2v2 only)
   ├─ EF_Company_Activity_Data (Gen2)
   └─ EF_VU_Configuration (Gen2)
```

### TLV Tag Structure

The 3-byte TLV tag encodes both the **File ID (FID)** and the **DF context**:

```
[Byte 0-1: File ID][Byte 2: Appendix/Generation]
```

**Tag Appendix Byte (Byte 2):**

- `0x00`: Data from Gen1 DF (Tachograph)
- `0x01`: Signature for Gen1 DF data
- `0x02`: Data from Gen2 DF (Tachograph_G2)
- `0x03`: Signature for Gen2 DF data

**Examples:**

- `0x050100`: EF_Application_Identification data, Gen1 DF
- `0x050101`: EF_Application_Identification signature, Gen1 DF
- `0x050102`: EF_Application_Identification data, Gen2 DF
- `0x050103`: EF_Application_Identification signature, Gen2 DF
- `0x000200`: EF_ICC data (always `00` - common file)
- `0x000500`: EF_IC data (always `00` - common file)

## Protobuf Data Model

### DF-Level Organization

The `DriverCardFile` protobuf message mirrors the logical DF/EF hierarchy:

```protobuf
message DriverCardFile {
  // Common files from Master File (MF)
  Icc icc = 1;
  Ic ic = 2;

  // Gen1 application data
  Tachograph tachograph = 3;

  // Gen2 application data
  TachographG2 tachograph_g2 = 4;

  // Nested DF messages
  message Tachograph {
    ApplicationIdentification application_identification = 1;
    Identification identification = 2;
    // ... other Gen1 EFs ...
  }

  message TachographG2 {
    ApplicationIdentification application_identification = 1;
    Identification identification = 2;
    // ... other Gen2 EFs (including Gen2-only EFs) ...
  }
}
```

### Key Design Principles

1. **DF Separation Prevents Data Loss**: Gen2 cards contain both Gen1 and Gen2 versions of many EFs for backward compatibility. Separate DF messages ensure both versions are preserved.

2. **Generation from Tag Appendix**: The generation is determined by the TLV tag's appendix byte during parsing, not from file content.

3. **Type Splitting for Structural Differences**: When a Data Dictionary type has different binary layouts between generations, create separate proto types (e.g., `PlaceRecord` vs `PlaceRecordG2`).

4. **One File Per Proto**: Each `.proto` file in `card/v1/` should have a corresponding `.go` file in `internal/card/` with the same base name.

## Parsing Flow (Unmarshal)

### Two-Pass Parsing

Card files are parsed in two passes:

**Pass 1: TLV Parsing → `RawCardFile`**

- Implemented in `rawcardfile.go`
- Splits the binary data into individual TLV records
- Extracts: tag, file type, generation (from appendix), content type (data vs signature), value bytes
- Output: `RawCardFile` with `repeated TlvRecord`

**Pass 2: Semantic Parsing → `DriverCardFile`**

- Implemented in `driver_card_file.go`
- Routes each TLV record to the appropriate DF based on generation
- Calls EF-specific unmarshal functions
- Attaches signatures to their corresponding EF data
- Output: Fully structured `DriverCardFile`

### Routing Logic

```go
func unmarshalDriverCardFile(input *RawCardFile) (*DriverCardFile, error) {
    var tachographDF *DriverCardFile_Tachograph
    var tachographG2DF *DriverCardFile_TachographG2

    for _, record := range input.GetRecords() {
        efGeneration := record.GetGeneration() // From tag appendix

        switch record.GetFile() {
        case EF_IDENTIFICATION:
            identification, err := opts.unmarshalIdentification(record.GetValue())
            // ...

            // Route to appropriate DF
            switch efGeneration {
            case Generation_GENERATION_1:
                if tachographDF == nil {
                    tachographDF = &DriverCardFile_Tachograph{}
                }
                tachographDF.SetIdentification(identification)

            case Generation_GENERATION_2:
                if tachographG2DF == nil {
                    tachographG2DF = &DriverCardFile_TachographG2{}
                }
                tachographG2DF.SetIdentification(identification)
            }
        }
    }

    output.SetTachograph(tachographDF)
    output.SetTachographG2(tachographG2DF)
    return output, nil
}
```

## Marshalling Flow

### DF-Aware Serialization

When marshalling, respect the DF hierarchy:

```go
func appendDriverCard(dst []byte, card *DriverCardFile) ([]byte, error) {
    // 1. Common files (MF) - always appendix 0x00
    dst = appendTlvUnsigned(dst, EF_ICC, card.GetIcc(), appendIcc)
    dst = appendTlvUnsigned(dst, EF_IC, card.GetIc(), appendCardIc)

    // 2. Gen1 DF - use appendix 0x00 for data, 0x01 for signatures
    if tachograph := card.GetTachograph(); tachograph != nil {
        dst = appendTlv(dst, EF_APPLICATION_IDENTIFICATION,
            tachograph.GetApplicationIdentification(),
            appendCardApplicationIdentification)
        // ... other Gen1 EFs ...
    }

    // 3. Gen2 DF - use appendix 0x02 for data, 0x03 for signatures
    if tachographG2 := card.GetTachographG2(); tachographG2 != nil {
        dst = appendTlvG2(dst, EF_APPLICATION_IDENTIFICATION,
            tachographG2.GetApplicationIdentification(),
            appendCardApplicationIdentification)
        // ... other Gen2 EFs ...
    }

    return dst, nil
}
```

### TLV Helpers

**For signed EFs (most EFs):**

```go
func appendTlv(dst []byte, fileType ElementaryFileType, msg Message,
    appenderFunc func([]byte, Message) ([]byte, error)) ([]byte, error) {

    // Write data block: [FID][0x00][Length][Value]
    // Write signature block: [FID][0x01][128][Signature]
    // ...
}
```

**For unsigned EFs (ICC, IC, certificates, Card_Download):**

```go
func appendTlvUnsigned(dst []byte, fileType ElementaryFileType, msg Message,
    appenderFunc func([]byte, Message) ([]byte, error)) ([]byte, error) {

    // Write data block only: [FID][0x00][Length][Value]
    // No signature block
    // ...
}
```

## Signature Handling

### Which EFs Are Signed?

**From Regulation Section 3.3 (DDP_035, DDP_037):**

**Unsigned EFs:**

- EF_ICC, EF_IC (common card information)
- All certificate EFs (Card_Certificate, CA_Certificate, CardSignCertificate, Link_Certificate)
- EF_Card_Download

**Signed EFs (all others):**

- EF_Application_Identification
- EF_Identification
- EF_Driving_Licence_Info
- EF_Events_Data
- EF_Faults_Data
- EF_Driver_Activity_Data
- EF_Vehicles_Used
- EF_Places
- EF_Current_Usage
- EF_Control_Activity_Data
- EF_Specific_Conditions
- EF_VehicleUnits_Used
- EF_GNSS_Places
- All other Gen2/Gen2v2 EFs

### Signature Attachment During Parsing

Signatures follow their data blocks in the TLV stream:

```go
for i := 0; i < len(records); i++ {
    record := records[i]

    // Look ahead for signature
    var signature []byte
    if i+1 < len(records) {
        nextRecord := records[i+1]
        if nextRecord.GetFile() == record.GetFile() &&
           nextRecord.GetContentType() == ContentType_SIGNATURE {
            signature = nextRecord.GetValue()
            i++ // Skip the signature record
        }
    }

    // Parse EF and attach signature
    ef, err := unmarshalEF(record.GetValue())
    if signature != nil {
        ef.SetSignature(signature)
    }
}
```

## Generation-Specific Type Splitting

### Principle: Split Types by Generation for Structural Differences

When a data structure has **different binary layouts or sizes** between generations, create separate protobuf types for each generation rather than using a superset message with conditional logic.

**Benefits of Type Splitting:**

1. **Fixed Sizes**: Each type has a deterministic, fixed size with no conditionals
2. **Type Safety**: The type system prevents mixing Gen1 and Gen2 data
3. **Simpler Code**: Parse/marshal functions are straightforward with no generation checks
4. **Better Testing**: Test each generation independently with clear expectations
5. **Clearer Schema**: The protobuf explicitly shows what exists in each generation

### When to Split Types

Split Data Dictionary types into separate Gen1 and Gen2 versions when:

1. **Different sizes**: Gen1 is X bytes, Gen2 is Y bytes
2. **Different layouts**: Fields at different offsets
3. **Structural changes**: Not just additive fields

**When to Use Superset (Don't Split):**

- **Pure addition**: Gen2 is Gen1 + extra byte(s) at the end with no layout changes (e.g., `FullCardNumberAndGeneration`)
- **Identical structures**: No differences across generations (e.g., `TimeReal`, `Date`)

### Example: PlaceRecord

**Gen1: 10 bytes (no GNSS)**

```go
// proto/wayplatform/connect/tachograph/dd/v1/place_record.proto
message PlaceRecord {
  google.protobuf.Timestamp entry_time = 1;  // 4 bytes
  EntryTypeDailyWorkPeriod entry_type = 2;   // 1 byte
  NationNumeric country = 4;                  // 1 byte
  bytes region = 6;                           // 1 byte
  int32 vehicle_odometer_km = 7;              // 3 bytes
  bytes raw_data = 8;                         // 10 bytes
  bool valid = 9;
}

// internal/dd/place_record.go
func (opts UnmarshalOptions) UnmarshalPlaceRecord(data []byte) (*PlaceRecord, error) {
    const lenPlaceRecord = 10  // Fixed!
    if len(data) != lenPlaceRecord {
        return nil, fmt.Errorf("invalid length: got %d, want %d", len(data), lenPlaceRecord)
    }
    // ... parse exactly 10 bytes ...
}
```

**Gen2: 21 bytes (includes GNSS)**

```go
// proto/wayplatform/connect/tachograph/dd/v1/place_record_g2.proto
message PlaceRecordG2 {
  google.protobuf.Timestamp entry_time = 1;     // 4 bytes
  EntryTypeDailyWorkPeriod entry_type = 2;      // 1 byte
  NationNumeric country = 4;                     // 1 byte
  bytes region = 6;                              // 1 byte
  int32 vehicle_odometer_km = 7;                 // 3 bytes
  GNSSPlaceRecord entry_gnss_place_record = 8;  // 11 bytes (NEW!)
  bytes raw_data = 9;                            // 21 bytes
  bool valid = 10;
}

// internal/dd/place_record_g2.go
func (opts UnmarshalOptions) UnmarshalPlaceRecordG2(data []byte) (*PlaceRecordG2, error) {
    const lenPlaceRecord = 21  // Fixed!
    if len(data) != lenPlaceRecord {
        return nil, fmt.Errorf("invalid length: got %d, want %d", len(data), lenPlaceRecord)
    }
    // ... parse exactly 21 bytes ...
}
```

**Usage in EF-level code:**

```go
// internal/card/places.go
func (opts UnmarshalOptions) unmarshalPlaces(data []byte) (*Places, error) {
    // ... parse header ...

    if opts.Generation == Generation_GENERATION_2 {
        records := parseCircularPlaceRecordsG2(data, opts)  // Uses PlaceRecordG2
        places.SetRecordsG2(records)
    } else {
        records := parseCircularPlaceRecordsGen1(data, opts)  // Uses PlaceRecord
        places.SetRecords(records)
    }

    return places, nil
}
```

### Benefits Achieved

✅ **Fixed sizes** - No conditionals in parsing
✅ **Type safety** - Can't mix Gen1/Gen2 records
✅ **Simpler code** - Each function has one job
✅ **Better testing** - Test generations independently
✅ **Clear schema** - Protobuf explicitly shows differences

## File Organization

### Proto Files (`proto/wayplatform/connect/tachograph/`)

```
dd/v1/
  ├─ place_record.proto           (Gen1 type)
  ├─ place_record_g2.proto        (Gen2 type)
  ├─ previous_vehicle_info.proto  (Gen1 type)
  ├─ previous_vehicle_info_g2.proto (Gen2 type)
  └─ ...

card/v1/
  ├─ driver_card_file.proto       (Top-level, contains DF messages)
  ├─ places.proto                 (EF with both Gen1/Gen2 record arrays)
  ├─ identification.proto         (EF)
  └─ ...
```

### Go Implementation Files (`internal/`)

```
dd/
  ├─ place_record.go              (Unmarshal/Append for Gen1)
  ├─ place_record_g2.go           (Unmarshal/Append for Gen2)
  ├─ previous_vehicle_info.go     (Unmarshal/Append for Gen1)
  ├─ previous_vehicle_info_g2.go  (Unmarshal/Append for Gen2)
  └─ ...

card/
  ├─ driver_card_file.go          (DF-level routing and TLV helpers)
  ├─ places.go                    (EF-level: uses dd/place_record*.go)
  ├─ identification.go            (EF-level)
  └─ ...
```

**One file per proto**: Each `.proto` file has a corresponding `.go` file with the same base name.

## Marshalling and Unmarshalling Implementation

### File Structure

The core principle is to organize files by type rather than by operation, with a direct correspondence to the protobuf schema definitions:

- **`internal/card/`**: For each card-related protobuf file (e.g., `card/v1/activity.proto`), there should be one corresponding file:

  - `<typename>.go`: Handles both marshalling and unmarshalling for card-specific protobuf message types

- **`internal/dd/`**: For each data dictionary protobuf file (e.g., `dd/v1/time.proto`), there should be one corresponding file:
  - `<typename>.go`: Handles both marshalling and unmarshalling for data dictionary types

This convention improves locality of context by keeping related marshalling and unmarshalling logic together, making it easier to spot inconsistencies and ensuring the operations stay in sync.

### Marshalling Pattern

Marshalling is implemented using a multi-level approach:

1. **Top-Level Function (`UnmarshalDriverCardFile`)**: Entry point in `internal/card` that orchestrates the full card file unmarshaling
2. **EF-Level Functions (`unmarshalPlaces`, `unmarshalIdentification`, etc.)**: Parse individual Elementary Files
3. **Appending Functions (`Append*`)**: Functions in `internal/dd` that follow the `BinaryAppender` pattern, taking a pre-allocated `[]byte` slice and appending their binary representation

### Unmarshalling Pattern

Unmarshalling follows a similar structure:

1. **Top-Level Function (`UnmarshalDriverCardFile`)**: Entry point that handles TLV parsing and DF routing
2. **EF-Level Functions (`unmarshalPlaces`, `unmarshalIdentification`, etc.)**: Parse individual Elementary Files
3. **Unmarshalling Functions (`Unmarshal*`)**: Functions in `internal/dd` responsible for parsing `[]byte` slices into protobuf messages

## Coding Principles

### Bufio Scanner Pattern for Record Parsing

Use `bufio.Scanner` with custom `SplitFunc` for all contiguous binary data parsing that advances forward through memory.

**Use for:** Fixed-size records, variable-length records, record arrays, complex structures
**Avoid for:** Backward iteration, linked lists with pointers, non-contiguous data, cyclic buffers

**Guidelines:**

- Co-locate `SplitFunc` in same file with descriptive name (e.g., `splitPlaceRecord`)
- Never reuse `SplitFunc` across different record types
- Use `unmarshal<ProtoMessage>` naming for parsing functions
- Return errors for invalid data in `SplitFunc` (fail-fast)
- Include proper size validation

**Pattern:**

```go
func splitPlaceRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
    const recordSize = 10
    if len(data) < recordSize {
        if atEOF { return 0, nil, nil }
        return 0, nil, nil
    }
    return recordSize, data[:recordSize], nil
}

func parseCircularPlaceRecords(data []byte, offset int) ([]*PlaceRecord, error) {
    scanner := bufio.NewScanner(bytes.NewReader(data[offset:]))
    scanner.Split(splitPlaceRecord)
    var records []*PlaceRecord
    for scanner.Scan() {
        record, err := unmarshalPlaceRecord(scanner.Bytes())
        if err != nil { return nil, err }
        records = append(records, record)
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return records, nil
}
```

### Nil Handling Policy

The binary tachograph protocol has no concept of `nil` or null values. Every field in the protocol is either present with valid data or absent (which is represented by specific zero/empty patterns in the binary format).

**Policy for `Append*` functions in `internal/dd`:**

- `Append*` functions **must error** when receiving a `nil` protobuf message parameter **if** the function needs to call nested `Append*` functions or access complex fields
- For functions that only read primitive fields (integers, bytes) where zero is a valid protocol value, **skip the nil check** and rely on protobuf's zero-value behavior
- Exception: `AppendStringValue` accepts `nil` and encodes it as an empty string (code page 255), as this is a valid protocol state
- When appending optional data, the caller should pass a properly initialized message with empty/zero values, not `nil`

**Rationale:** This policy catches bugs early by failing fast when data is missing, rather than silently writing incorrect/default values to the binary output.

### Exact Length Validation Policy

When parsing fixed-size binary structures, we must validate that the input data length exactly matches the expected size. The protocol is strictly defined - if we expect N bytes and receive a different amount, something has already gone wrong upstream and we should fail early.

**Policy for `Unmarshal*` functions in `internal/dd`:**

- For fixed-size structures, validate with `len(data) == expectedSize`, not `len(data) >= expectedSize`
- For variable-size structures with known minimums, validate with `len(data) < expectedSize` only when consuming from a stream
- When unmarshalling a complete structure from a byte slice, the slice should contain exactly the expected bytes
- Extra bytes indicate a parsing error upstream (incorrect offset calculation, wrong structure interpretation, etc.)

**Rationale:** Strict validation catches bugs early. If a 4-byte timestamp gets 5 bytes, that's an error that should be caught immediately, not silently ignored.

**Example:**

```go
// GOOD: Requires exact length
func UnmarshalTimeReal(data []byte) (*timestamppb.Timestamp, error) {
    const lenTimeReal = 4
    if len(data) != lenTimeReal {  // Correct! Exact match required
        return nil, fmt.Errorf("invalid data length for TimeReal: got %d, want %d", len(data), lenTimeReal)
    }
    // ...
}
```

### Raw Data Painting Policy

When marshalling data structures that have both semantic fields and a `raw_data` field preserving the original binary representation, use the "raw data painting" strategy to achieve optimal round-trip fidelity while ensuring semantic field correctness.

**Policy for `Append*` functions:**

- **Always prefer raw_data as a canvas**: If `raw_data` is available and has the correct length, make a copy of it and use it as a canvas for marshalling
- **Paint semantic values over the canvas**: Serialize semantic fields on top of the canvas at their designated byte offsets, overwriting those specific bytes. **Critical**: Do NOT just return raw_data as-is - you must encode the semantic values and write them over the canvas
- **Preserve unknown bits**: Any padding bytes, reserved bits, or unknown data in the original `raw_data` are automatically preserved in areas not overwritten by semantic fields
- **Fall back to zero canvas**: If `raw_data` is unavailable or has incorrect length, create a zero-filled buffer of the correct size and serialize semantic fields into it

**Rationale:** This approach provides three critical benefits:

1. **Round-trip fidelity**: Reserved bits, padding, and vendor-specific data are preserved exactly as they appeared in the original binary
2. **Semantic field validation**: When round-trip tests pass, it proves the semantic fields were correctly parsed and serialized, not just that raw bytes were copied
3. **Maximum trust**: The serialized output is guaranteed to match the original binary format because it literally uses the original as a template, while also validating that our semantic understanding is correct

**Example:**

```go
// GOOD: Raw data painting strategy with stack-allocated canvas
func AppendDate(dst []byte, date *ddv1.Date) ([]byte, error) {
    const lenDatef = 4

    // Use stack-allocated array for the canvas (fixed size, avoids heap allocation)
    var canvas [lenDatef]byte

    // Start with raw_data as canvas if available (raw data painting approach)
    if rawData := date.GetRawData(); len(rawData) > 0 {
        if len(rawData) != lenDatef {
            return nil, fmt.Errorf("invalid raw_data length for Date: got %d, want %d", len(rawData), lenDatef)
        }
        copy(canvas[:], rawData)
    }
    // Otherwise canvas is zero-initialized (Go default)

    // Paint semantic values over the canvas
    year := int(date.GetYear())
    month := int(date.GetMonth())
    day := int(date.GetDay())
    canvas[0] = byte((year/1000)%10<<4 | (year/100)%10)
    canvas[1] = byte((year/10)%10<<4 | year%10)
    canvas[2] = byte((month/10)%10<<4 | month%10)
    canvas[3] = byte((day/10)%10<<4 | day%10)

    return append(dst, canvas[:]...), nil
}
```

### In-Code Documentation and Context

To make the marshalling and unmarshalling logic as robust and maintainable as possible, we bring critical context from the regulation specifications directly into the code.

#### ASN.1 Definitions in Comments

Every function that marshals or unmarshals a data structure defined in the ASN.1 specification should include the corresponding ASN.1 definition in its function-level comment block. This provides immediate context for the binary layout.

#### Constants for Binary Layout

Avoid using "magic numbers" for sizes and offsets. Instead, define a `const` block within the function to specify the byte layout (offsets, lengths) of the structure being processed. Use the `idx` prefix for offsets and `len` for lengths to make them easy to identify.

**Example:**

```go
// unmarshalIdentification parses the CardIdentification structure.
//
// The data type `CardIdentification` is specified in the Data Dictionary, Section 2.24.
//
// ASN.1 Specification:
//
//     CardIdentification ::= SEQUENCE {
//         cardIssuingMemberState          NationNumeric,         -- 1 byte
//         cardNumber                      CardNumber,            -- 16 bytes
//         cardIssuingAuthorityName        Name,                  -- 36 bytes
//         cardIssueDate                   TimeReal,              -- 4 bytes
//         cardValidityBegin               TimeReal,              -- 4 bytes
//         cardExpiryDate                  TimeReal               -- 4 bytes
//     }
func (opts UnmarshalOptions) unmarshalIdentification(data []byte) (*cardv1.Identification, error) {
    const (
        idxIssuingMemberState = 0
        idxCardNumber         = 1
        idxAuthorityName      = 17
        idxIssueDate          = 53
        idxValidityBegin      = 57
        idxExpiryDate         = 61
        lenIdentification     = 65
    )

    if len(data) != lenIdentification {
        return nil, fmt.Errorf("invalid data length: got %d, want %d", len(data), lenIdentification)
    }
    // ... parsing logic ...
}
```

## Testing Strategy

### Testing Framework

All tests must use **only** the standard library `testing` package and `github.com/google/go-cmp/cmp` for comparisons. Do not use third-party testing frameworks like `testify`.

**Guidelines:**

- Use `t.Errorf()` for non-fatal errors and `t.Fatalf()` for fatal errors
- Use `cmp.Diff()` for comparing complex structures (slices, maps, structs)
- Use standard equality checks (`==`, `!=`) for simple types
- Check for nil explicitly before accessing pointers
- Always check errors before proceeding with test logic

### EF-Level Tests

Each Elementary File implementation should have comprehensive tests:

**Round-trip tests**: Verify that `unmarshal → marshal → unmarshal` produces identical results:

```go
func TestPlacesRoundTrip(t *testing.T) {
    // Read testdata/places.b64
    b64Data, err := os.ReadFile("testdata/places.b64")
    if err != nil {
        t.Fatalf("Failed to read test data: %v", err)
    }
    data, err := base64.StdEncoding.DecodeString(string(b64Data))
    if err != nil {
        t.Fatalf("Failed to decode base64: %v", err)
    }

    // First unmarshal
    opts := UnmarshalOptions{}
    places1, err := opts.unmarshalPlaces(data)
    if err != nil {
        t.Fatalf("First unmarshal failed: %v", err)
    }

    // Marshal
    marshaled, err := appendPlaces(nil, places1)
    if err != nil {
        t.Fatalf("Marshal failed: %v", err)
    }

    // Assert: binary equality
    if diff := cmp.Diff(data, marshaled); diff != "" {
        t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
    }

    // Second unmarshal
    places2, err := opts.unmarshalPlaces(marshaled)
    if err != nil {
        t.Fatalf("Second unmarshal failed: %v", err)
    }

    // Assert: structural equality
    if diff := cmp.Diff(places1, places2, protocmp.Transform()); diff != "" {
        t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
    }
}
```

**Anonymization tests**: Use the golden file pattern with the `-update` flag. See [testdata/AGENTS.md](testdata/AGENTS.md) for comprehensive guidance on creating anonymized test data.

### Data Dictionary Tests

Test generation-specific parsers independently:

```go
// internal/dd/place_record_test.go
func TestUnmarshalPlaceRecord(t *testing.T) {
    const lenPlaceRecord = 10  // Always 10!

    tests := []struct {
        name string
        data []byte
        want *PlaceRecord
    }{
        // ... test cases with exactly 10-byte inputs ...
    }
    // ...
}

// internal/dd/place_record_g2_test.go
func TestUnmarshalPlaceRecordG2(t *testing.T) {
    const lenPlaceRecord = 21  // Always 21!

    tests := []struct {
        name string
        data []byte
        want *PlaceRecordG2
    }{
        // ... test cases with exactly 21-byte inputs ...
    }
    // ...
}
```

### Golden File Tests with Anonymization

For comprehensive guidance on creating anonymized test data using the golden file pattern, see **[testdata/AGENTS.md](testdata/AGENTS.md)**.

Key points:

- All test data must be deterministically anonymized
- Use the `-update` flag to regenerate golden files when logic changes
- Round-trip tests validate binary fidelity
- Anonymization tests validate that anonymization is deterministic

## Code Quality

- **No `//nolint` comments**: Never suppress linter warnings with `//nolint` comments. Instead, fix the underlying issues by removing unused code, implementing missing functionality, or restructuring the code properly.
- **Zero linter errors**: The codebase must have zero linter errors at all times. This ensures code quality and maintainability.

## Summary

The card package implements a **two-pass parsing strategy** that respects the **DF/EF hierarchy** of tachograph card files:

1. **TLV layer** (`RawCardFile`): Splits binary stream into tagged records
2. **DF routing layer** (`DriverCardFile`): Routes EFs to Gen1/Gen2 DFs based on tag appendix
3. **EF parsing layer**: Calls generation-specific unmarshal functions from `internal/dd`
4. **DD type layer**: Fixed-size, generation-specific types with no conditionals

This architecture ensures **perfect binary fidelity**, **type safety**, and **maintainable code** that aligns precisely with the EU tachograph regulation.
