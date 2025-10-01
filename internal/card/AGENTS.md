# Card Package Guidelines

This document provides detailed guidance for implementing and maintaining card file parsing in the `internal/card` package.

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

### When to Split Types

Split Data Dictionary types into separate Gen1 and Gen2 versions when:

1. **Different sizes**: Gen1 is X bytes, Gen2 is Y bytes
2. **Different layouts**: Fields at different offsets
3. **Structural changes**: Not just additive fields

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

## Testing Strategy

### Golden File Tests

Test full card file parsing with real-world examples:

```go
// golden_test.go
func TestDriverCardGolden(t *testing.T) {
    data, err := os.ReadFile("testdata/card.DDD")
    // ...

    file, err := UnmarshalDriverCardFile(rawCard)
    // ...

    // Compare with golden JSON
    // ...
}
```

### Unit Tests for DD Types

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

## Summary

The card package implements a **two-pass parsing strategy** that respects the **DF/EF hierarchy** of tachograph card files:

1. **TLV layer** (`RawCardFile`): Splits binary stream into tagged records
2. **DF routing layer** (`DriverCardFile`): Routes EFs to Gen1/Gen2 DFs based on tag appendix
3. **EF parsing layer**: Calls generation-specific unmarshal functions from `internal/dd`
4. **DD type layer**: Fixed-size, generation-specific types with no conditionals

This architecture ensures **perfect binary fidelity**, **type safety**, and **maintainable code** that aligns precisely with the EU tachograph regulation.
