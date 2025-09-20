# Tachograph Data Parsing Guide

This document explains how to parse tachograph data structures for both TLV (card files) and TV (vehicle unit files) formats across all generations.

## Overview

After identifying file types using tags (see [tags.md](tags.md)), the next step is parsing the actual data structures. Each format has distinct parsing requirements:

- **TLV Format**: Used by card files with explicit length fields
- **TV Format**: Used by vehicle unit files with implicit length based on data type
- **Generation Differences**: Gen1, Gen2, and Gen2V2 have different data layouts

## TLV Format Parsing (Card Files)

### Binary Structure

Card files use **TLV (Tag-Length-Value)** encoding as specified in **Appendix 2, Section 4**:

```
┌─────────────┬─────────────┬─────────────┐
│ Tag (2B)    │ Length (2B) │ Value (N)   │
├─────────────┼─────────────┼─────────────┤
│ Big Endian  │ Big Endian  │ Raw Data    │
└─────────────┴─────────────┴─────────────┘
```

### Extended TLV Structure

The benchmark implementation shows that card files actually use a **5-byte header**:

```
Card TLV Structure (from benchmark/tachoparser/pkg/decoder/definitions.go):
┌─────────────┬─────────────┬─────────────┐
│ Tag (3B)    │ Length (2B) │ Value (N)   │
├─────────────┼─────────────┼─────────────┤
│ FID + Gen   │ Big Endian  │ Raw Data    │
└─────────────┴─────────────┴─────────────┘

Where Tag = [FID (2B)][Generation (1B)]
- Generation: 00=Data, 01=Signature, 02=Gen2Data, 03=Gen2Signature
```

### TLV Parsing Algorithm

```go
func parseTLV(data []byte) ([]TLVRecord, error) {
    var records []TLVRecord
    offset := 0

    for offset < len(data) {
        // 1. Check minimum header size
        if len(data)-offset < 5 {
            break
        }

        // 2. Extract tag (3 bytes) and length (2 bytes)
        tag := binary.BigEndian.Uint32(data[offset:offset+3])
        length := binary.BigEndian.Uint16(data[offset+3:offset+5])

        // 3. Check if complete record is available
        totalLength := 5 + int(length)
        if len(data)-offset < totalLength {
            return nil, fmt.Errorf("incomplete TLV record")
        }

        // 4. Extract value
        value := data[offset+5:offset+totalLength]

        records = append(records, TLVRecord{
            Tag:    tag,
            Length: length,
            Value:  value,
        })

        offset += totalLength
    }

    return records, nil
}
```

### Card Data Structure Selection

Based on **Appendix 2** and the benchmark implementation:

```go
func selectCardStructure(tag uint32) interface{} {
    // Extract FID (first 2 bytes) and generation (last byte)
    fid := uint16(tag >> 8)
    generation := uint8(tag & 0xFF)

    switch fid {
    case 0x0002: // EF_ICC
        switch generation {
        case 0x00: return &CardIccIdentificationFirstGen{}
        case 0x02: return &CardIccIdentificationSecondGen{}
        }
    case 0x0501: // EF_Application_Identification
        switch generation {
        case 0x00: return &DriverCardApplicationIdentificationFirstGen{}
        case 0x02: return &DriverCardApplicationIdentificationSecondGen{}
        }
    case 0x0507: // EF_Events_Data
        switch generation {
        case 0x00: return &CardEventDataFirstGen{}
        case 0x02: return &CardEventDataSecondGen{}
        }
    // ... more cases based on tachocard/tags.go
    }
    return nil
}
```

## TV Format Parsing (Vehicle Unit Files)

### Binary Structure

Vehicle unit files use **TV (Tag-Value)** encoding as specified in **Appendix 7, Section 2.2.6**:

```
┌─────────────┬─────────────┐
│ Tag (2B)    │ Value (N)   │
├─────────────┼─────────────┤
│ Big Endian  │ Raw Data    │
└─────────────┴─────────────┘
```

**Critical**: No explicit length field - value size is determined by the data structure type.

### TV Parsing Algorithm

```go
func parseTV(data []byte) ([]TVRecord, error) {
    var records []TVRecord
    offset := 0

    for offset < len(data) {
        // 1. Check minimum tag size
        if len(data)-offset < 2 {
            break
        }

        // 2. Extract tag
        tag := binary.BigEndian.Uint16(data[offset:offset+2])

        // 3. Determine structure size based on tag
        structSize, err := getVuStructureSize(tag)
        if err != nil {
            return nil, fmt.Errorf("unknown VU tag: 0x%04X", tag)
        }

        // 4. Check if complete record is available
        totalLength := 2 + structSize
        if len(data)-offset < totalLength {
            return nil, fmt.Errorf("incomplete TV record")
        }

        // 5. Extract value
        value := data[offset+2:offset+totalLength]

        records = append(records, TVRecord{
            Tag:   tag,
            Value: value,
        })

        offset += totalLength
    }

    return records, nil
}
```

### VU Data Structure Selection

Based on **Appendix 1 (Data Dictionary)** and **Appendix 7, Section 2.2.6**:

```go
func selectVuStructure(tag uint16) (interface{}, int) {
    switch tag {
    // Generation 1 structures
    case 0x7601: return &VuOverviewFirstGen{}, sizeOfVuOverviewFirstGen
    case 0x7602: return &VuActivitiesFirstGen{}, sizeOfVuActivitiesFirstGen
    case 0x7603: return &VuEventsAndFaultsFirstGen{}, sizeOfVuEventsAndFaultsFirstGen
    case 0x7604: return &VuDetailedSpeedFirstGen{}, sizeOfVuDetailedSpeedFirstGen
    case 0x7605: return &VuTechnicalDataFirstGen{}, sizeOfVuTechnicalDataFirstGen

    // Generation 2 Version 1 structures
    case 0x7621: return &VuOverviewSecondGen{}, sizeOfVuOverviewSecondGen
    case 0x7622: return &VuActivitiesSecondGen{}, sizeOfVuActivitiesSecondGen
    case 0x7623: return &VuEventsAndFaultsSecondGen{}, sizeOfVuEventsAndFaultsSecondGen
    case 0x7624: return &VuDetailedSpeedSecondGen{}, sizeOfVuDetailedSpeedSecondGen
    case 0x7625: return &VuTechnicalDataSecondGen{}, sizeOfVuTechnicalDataSecondGen

    // Generation 2 Version 2 structures
    case 0x7600: return &DownloadInterfaceVersion{}, sizeOfDownloadInterfaceVersion
    case 0x7631: return &VuOverviewSecondGenV2{}, sizeOfVuOverviewSecondGenV2
    case 0x7632: return &VuActivitiesSecondGenV2{}, sizeOfVuActivitiesSecondGenV2
    case 0x7633: return &VuEventsAndFaultsSecondGenV2{}, sizeOfVuEventsAndFaultsSecondGenV2
    case 0x7635: return &VuTechnicalDataSecondGenV2{}, sizeOfVuTechnicalDataSecondGenV2

    default: return nil, 0
    }
}
```

## Generation-Specific Data Structures

### Generation 1 (TREP 0x01-0x05)

**Regulation Source**: Original tachograph specification
**Security**: Appendix 11, Part A (legacy cryptography)
**Characteristics**:

- Simpler data layouts
- Fixed-size structures
- Basic certificate formats
- SHA-1 based signatures

**Example Structure** (from regulation analysis):

```go
type VuOverviewFirstGen struct {
    MemberStateCertificate [194]byte  // Certificate data
    VuCertificate         [194]byte   // VU certificate
    VehicleIdentificationNumber string // VIN
    VehicleRegistrationIdentification VehicleRegistrationIdentification
    CurrentDateTime       uint32      // Unix timestamp
    VuDownloadablePeriod  VuDownloadablePeriod
    CardSlotsStatus       uint8       // Card slot status
    VuDownloadActivityData VuDownloadActivityDataFirstGen
    VuCompanyLocksData    VuCompanyLocksDataFirstGen
    VuControlActivityData VuControlActivityDataFirstGen
    Signature             SignatureFirstGen // SHA-1 signature
}
```

### Generation 2 Version 1 (TREP 0x21-0x25)

**Regulation Source**: **Appendix 7, Section 2.2.6** + **Appendix 1**
**Security**: Appendix 11, Part B (enhanced cryptography)
**Characteristics**:

- Record array format with headers
- Enhanced security (SHA-256, ECC)
- Backward compatibility with Gen1

**Key Difference**: **Record Arrays** as specified in **Appendix 7**:

> "For generation 2 downloads, each top-level data element is represented by a record array, even if it contains only one record. A record array starts with a header; this header contains the record type, the record size and the number of records."

**Record Array Structure**:

```go
type RecordArrayHeader struct {
    RecordType   uint8   // Type identifier
    RecordSize   uint16  // Size of each record
    NumberOfRecords uint16 // Count of records
}

type VuOverviewSecondGen struct {
    MemberStateCertificateRecordArray MemberStateCertificateRecordArray
    VuCertificateRecordArray         VuCertificateRecordArray
    VehicleIdentificationNumberRecordArray VehicleIdentificationNumberRecordArray
    // ... other record arrays
    SignatureRecordArray             SignatureRecordArray
}
```

### Generation 2 Version 2 (TREP 0x00, 0x31-0x35)

**Regulation Source**: **Appendix 7, Section 2.2.6** (latest amendments)
**Security**: Appendix 11, Part B + OSNMA support
**Characteristics**:

- Extended record arrays
- GNSS/OSNMA integration
- Additional data fields
- Download Interface Version support

**Special Features**:

- **0x7600**: Download Interface Version (unique to Gen2V2)
- Enhanced GNSS data structures
- Extended certificate formats

## Data Type Definitions

### Primitive Types

Based on **Appendix 1, Section 1.1** and regulation analysis:

| Type            | Size     | Description                 | Regulation Reference |
| --------------- | -------- | --------------------------- | -------------------- |
| `uint8`         | 1 byte   | Unsigned integer            | Data Dictionary      |
| `uint16`        | 2 bytes  | Big-endian unsigned integer | Data Dictionary      |
| `uint32`        | 4 bytes  | Big-endian unsigned integer | Data Dictionary      |
| `TimeReal`      | 4 bytes  | Unix timestamp              | Appendix 1, 2.xxx    |
| `BCDString`     | Variable | BCD encoded string          | Appendix 1           |
| `OdometerShort` | 3 bytes  | Odometer reading            | Appendix 1           |

### Complex Types

**ActivityChangeInfo** (Appendix 1):

```go
type ActivityChangeInfo struct {
    Slot                uint8    // Driver slot (1 or 2)
    DrivingStatus       uint8    // Activity status
    CardStatus          uint8    // Card insertion status
    ActivityTime        TimeReal // Time of activity change
}
```

**VehicleRegistrationIdentification** (Appendix 1):

```go
type VehicleRegistrationIdentification struct {
    VehicleRegistrationNation    uint8     // Nation code
    VehicleRegistrationNumber    [14]byte  // Registration number
}
```

## Parsing Challenges and Solutions

### 1. Generation Detection

**Challenge**: Same data type, different structure sizes across generations.

**Solution**: Always check tag first to determine generation:

```go
trep := byte(tag & 0xFF)
switch {
case trep >= 0x01 && trep <= 0x05:
    // Use Generation 1 structures
case trep >= 0x21 && trep <= 0x25:
    // Use Generation 2 V1 structures (with record arrays)
case trep >= 0x31 && trep <= 0x35:
    // Use Generation 2 V2 structures (extended)
}
```

### 2. Record Array Parsing

**Challenge**: Gen2 uses record arrays with headers.

**Solution**: Parse header first, then iterate through records:

```go
func parseRecordArray(data []byte) (interface{}, error) {
    // Parse header
    header := RecordArrayHeader{
        RecordType:      data[0],
        RecordSize:      binary.BigEndian.Uint16(data[1:3]),
        NumberOfRecords: binary.BigEndian.Uint16(data[3:5]),
    }

    // Parse records
    offset := 5
    records := make([]interface{}, header.NumberOfRecords)

    for i := 0; i < int(header.NumberOfRecords); i++ {
        record := parseRecord(data[offset:offset+int(header.RecordSize)])
        records[i] = record
        offset += int(header.RecordSize)
    }

    return records, nil
}
```

### 3. Variable Length Fields

**Challenge**: Some fields have variable lengths (strings, arrays).

**Solution**: Use length prefixes or fixed-size buffers as specified in **Appendix 1**:

```go
// BCD String with length prefix
func parseBCDString(data []byte) (string, int) {
    length := data[0]
    bcdData := data[1:1+length]
    return decodeBCD(bcdData), 1 + int(length)
}
```

### 4. Mixed Generation Files

**Challenge**: Single file can contain multiple generations.

**Solution**: Parse each block independently:

```go
func parseVuFile(data []byte) (*VuData, error) {
    result := &VuData{}

    for offset := 0; offset < len(data); {
        tag := binary.BigEndian.Uint16(data[offset:offset+2])

        // Determine generation and structure
        structure, size := selectVuStructure(tag)
        if structure == nil {
            return nil, fmt.Errorf("unknown tag: 0x%04X", tag)
        }

        // Parse this block
        blockData := data[offset+2:offset+2+size]
        parseStructure(structure, blockData)

        // Add to appropriate generation collection
        addToResult(result, tag, structure)

        offset += 2 + size
    }

    return result, nil
}
```

## Implementation Examples

### Complete TLV Parser

```go
func NewTLVScanner(r io.Reader) *bufio.Scanner {
    scanner := bufio.NewScanner(r)
    scanner.Split(tachocard.SplitFunc) // 5-byte TLV header
    return scanner
}

func parseTLVFile(data []byte) (*CardData, error) {
    scanner := NewTLVScanner(bytes.NewReader(data))
    result := &CardData{}

    for scanner.Scan() {
        record := scanner.Bytes()

        // Extract tag and determine structure
        tag := binary.BigEndian.Uint32(record[0:3])
        length := binary.BigEndian.Uint16(record[3:5])
        value := record[5:5+length]

        structure := selectCardStructure(tag)
        if structure != nil {
            parseStructure(structure, value)
            addToCardData(result, tag, structure)
        }
    }

    return result, scanner.Err()
}
```

### Complete TV Parser

```go
func parseTVFile(data []byte) (*VuData, error) {
    result := &VuData{}
    offset := 0

    for offset < len(data) {
        if len(data)-offset < 2 {
            break
        }

        tag := binary.BigEndian.Uint16(data[offset:offset+2])
        structure, size := selectVuStructure(tag)

        if structure == nil {
            return nil, fmt.Errorf("unknown VU tag: 0x%04X", tag)
        }

        if len(data)-offset < 2+size {
            return nil, fmt.Errorf("incomplete TV record")
        }

        value := data[offset+2:offset+2+size]
        parseStructure(structure, value)
        addToVuData(result, tag, structure)

        offset += 2 + size
    }

    return result, nil
}
```

## Regulation References

| Component           | Regulation Section        | Description                       |
| ------------------- | ------------------------- | --------------------------------- |
| **TLV Format**      | Appendix 2, Section 4     | Card file TLV structure           |
| **TV Format**       | Appendix 7, Section 2.2.6 | VU file TV structure              |
| **Data Dictionary** | Appendix 1                | Complete type definitions         |
| **Record Arrays**   | Appendix 7, Section 2.2.6 | Gen2 record array format          |
| **Generation 1**    | Original regulation       | Legacy data structures            |
| **Generation 2**    | Amended regulation        | Enhanced data structures          |
| **Security**        | Appendix 11               | Signature and certificate formats |

## Tools and Implementation

- **`tachocard/scanner.go`**: TLV parsing utilities
- **`tachounit/vu_tags.go`**: TV tag definitions
- **`benchmark/tachoparser/`**: Complete reference implementation
- **`docs/tags.md`**: Tag identification guide

## Best Practices

1. **Always validate tags** before parsing data structures
2. **Handle mixed generations** in the same file
3. **Verify data sizes** against expected structure sizes
4. **Implement proper error handling** for malformed data
5. **Use regulation-compliant** data type definitions
6. **Test with real data** from all generations

This guide provides the foundation for implementing robust tachograph data parsing that complies with EU regulations across all generations.
