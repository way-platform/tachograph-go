# Tachograph File Formats and Tag Systems

This document explains the binary formats used in tachograph files and how to derive tag definitions from the EU regulation.

## Overview

Tachograph files use two distinct encoding formats:

- **TLV (Tag-Length-Value)** for card files
- **TV (Tag-Value)** for vehicle unit files

Understanding these formats is crucial for proper file type identification and data parsing.

## Card Files: TLV Format

### Binary Structure

Card files use **TLV (Tag-Length-Value)** encoding as specified in **Appendix 7, Section 3.4.2**:

```
┌─────────────┬─────────────┬─────────────┐
│ Tag (2B)    │ Length (2B) │ Value (N)   │
├─────────────┼─────────────┼─────────────┤
│ Big Endian  │ Big Endian  │ Raw Data    │
└─────────────┴─────────────┴─────────────┘
```

### Actual TLV Structure (3-Byte Tag)

Real-world implementation analysis reveals that card files actually use a **5-byte header** with a 3-byte tag:

```
┌─────────────┬─────────────┬─────────────┐
│ Tag (3B)    │ Length (2B) │ Value (N)   │
├─────────────┼─────────────┼─────────────┤
│ FID + Gen   │ Big Endian  │ Raw Data    │
└─────────────┴─────────────┴─────────────┘

Where Tag = [FID (2B)][Generation (1B)]
- Generation: 00=Data, 01=Signature, 02=Gen2Data, 03=Gen2Signature
```

This structure allows for:

- **File Identification**: 2-byte FID identifies the Elementary File type
- **Generation Detection**: 1-byte generation flag distinguishes data/signature and Gen1/Gen2
- **Mixed Generation Support**: Single file can contain both Gen1 and Gen2 records

### File Identification

According to **Appendix 7, Section 3.3.2**, card files always start with **EF_ICC** (Elementary File - IC Card):

```go
// First 2 bytes = 0x0002 (EF_ICC tag)
if binary.BigEndian.Uint16(data[0:2]) == 0x0002 {
    return CardFileType
}
```

### Tag Sources

Card file tags are defined in **Appendix 2** and can be extracted from:

- HTML tables containing File Identifiers (FID)
- 4-character hex values (e.g., "0002", "0501", "0502")
- Associated descriptions for each Elementary File

**Implementation**: See `tachocard/tags.go` and `tools/cmd/tachocard-tag-parser/`

### Extended Card File Identifiers

Beyond the basic Elementary Files defined in **Appendix 2**, real-world tachograph cards contain additional FIDs discovered through analysis of actual card files:

#### Generation 2 Extended FIDs (0x052x Series)

These FIDs represent extended functionality introduced with Generation 2 tachographs:

| FID      | Name                                                             | Purpose                        | Source                   |
| -------- | ---------------------------------------------------------------- | ------------------------------ | ------------------------ |
| `0x0520` | **EF_Card_Identification_And_Driver_Card_Holder_Identification** | Enhanced driver identification | Benchmark implementation |
| `0x0521` | **EF_Card_Driving_Licence_Information**                          | Driving license details        | Benchmark implementation |
| `0x0522` | **EF_Specific_Conditions_Extended**                              | Extended specific conditions   | Benchmark implementation |
| `0x0523` | **EF_Card_Vehicle_Units_Used**                                   | Vehicle units usage history    | Benchmark implementation |
| `0x0524` | **EF_GNSS_Accumulated_Driving**                                  | GNSS-based driving data        | Benchmark implementation |
| `0x0525` | **EF_Driver_Card_Application_Identification_V2**                 | Enhanced application ID        | Benchmark implementation |
| `0x0526` | **EF_Card_Place_Auth_Daily_Work_Period**                         | Authenticated place data       | Benchmark implementation |
| `0x0527` | **EF_GNSS_Auth_Accumulated_Driving**                             | Authenticated GNSS driving     | Benchmark implementation |
| `0x0528` | **EF_Card_Border_Crossings**                                     | Border crossing records        | Benchmark implementation |
| `0x0529` | **EF_Card_Load_Unload_Operations**                               | Load/unload operations         | Benchmark implementation |

#### GNSS-Related FIDs (0xC10x Series)

These FIDs are defined in **Appendix 12 (GNSS)** for GNSS facility certificates:

| FID      | Name                         | Purpose               | Source               |
| -------- | ---------------------------- | --------------------- | -------------------- |
| `0xC100` | **EF_EGF_MACertificate**     | GNSS MA Certificate   | Appendix 12, Table 1 |
| `0xC108` | **EF_CA_Certificate_GNSS**   | GNSS CA Certificate   | Appendix 12, Table 1 |
| `0xC109` | **EF_Link_Certificate_GNSS** | GNSS Link Certificate | Appendix 12, Table 1 |

**Note**: The `0xC101` FID appears in real files but is not explicitly documented in the regulation - likely an extended or manufacturer-specific variant.

## Vehicle Unit Files: TV Format

### Binary Structure

Vehicle unit files use **TV (Tag-Value)** encoding as specified in **Appendix 7, Section 2.2.6**:

```
┌─────────────┬─────────────┐
│ Tag (2B)    │ Value (N)   │
├─────────────┼─────────────┤
│ Big Endian  │ Raw Data    │
└─────────────┴─────────────┘
```

**Key difference**: No explicit length field - the value size is determined by the data structure type.

### Tag Composition

VU tags are formed by combining two protocol components:

```
VU Tag = SID + TREP
┌─────────────┬─────────────┐
│ SID (0x76)  │ TREP (0xXX) │
├─────────────┼─────────────┤
│ Service ID  │ Transfer    │
│             │ Response    │
│             │ Parameter   │
└─────────────┴─────────────┘
```

### File Identification

```go
// Check if first 2 bytes form a valid VU TV tag
if tachounit.VuTag(firstTag).IsValid() {
    return UnitFileType
}
```

### Tag Derivation from Regulation

VU TV tags are **not explicitly listed in tables** but must be derived from protocol descriptions in **Appendix 7, Section 2.2.6**:

#### Step 1: Locate DDP Sections

Find the Data Download Protocol sections:

- **DDP_028a**: Download Interface Version
- **DDP_029**: Overview data
- **DDP_030**: Activities data
- **DDP_031**: Events and Faults data
- **DDP_032**: Detailed Speed data
- **DDP_033**: Technical Data

#### Step 2: Extract SID + TREP Patterns

Look for text patterns like:

```
"SID 76 Hex, the TREP 01, 21 or 31 Hex"
"SID 76 Hex, the TREP 02, 22 or 32 Hex"
"SID 76 Hex, the TREP 04 or 24 Hex"
```

#### Step 3: Generate Tag Values

Combine SID (0x76) with each TREP value:

```
SID 0x76 + TREP 0x01 = 0x7601 (VU_OverviewFirstGen)
SID 0x76 + TREP 0x21 = 0x7621 (VU_OverviewSecondGen)
SID 0x76 + TREP 0x31 = 0x7631 (VU_OverviewSecondGenV2)
```

#### Step 4: Map to Generations

TREP ranges indicate tachograph generations:

- **0x00-0x0F**: Generation 1 or special (e.g., 0x00 = Download Interface Version)
- **0x20-0x2F**: Generation 2 Version 1
- **0x30-0x3F**: Generation 2 Version 2

#### Generation Identification Logic

```go
func identifyGeneration(trep byte) string {
    switch {
    case trep == 0x00:
        return "DownloadInterfaceVersion" // Special case
    case trep >= 0x01 && trep <= 0x0F:
        return "FirstGen"                 // Generation 1
    case trep >= 0x20 && trep <= 0x2F:
        return "SecondGen"               // Generation 2 Version 1
    case trep >= 0x30 && trep <= 0x3F:
        return "SecondGenV2"             // Generation 2 Version 2
    default:
        return "Unknown"
    }
}
```

### Complete VU Tag Set

Based on **Appendix 7, Sections 2.2.6.1-2.2.6.6**:

| Data Type          | Gen 1    | Gen 2 V1 | Gen 2 V2 | Source   |
| ------------------ | -------- | -------- | -------- | -------- |
| Download Interface | -        | -        | `0x7600` | DDP_028a |
| Overview           | `0x7601` | `0x7621` | `0x7631` | DDP_029  |
| Activities         | `0x7602` | `0x7622` | `0x7632` | DDP_030  |
| Events & Faults    | `0x7603` | `0x7623` | `0x7633` | DDP_031  |
| Detailed Speed     | `0x7604` | `0x7624` | -        | DDP_032  |
| Technical Data     | `0x7605` | `0x7625` | `0x7635` | DDP_033  |

**Notes**:

- `0x7634` (Detailed Speed Gen2 V2) does not exist per regulation specifications
- `0x7600` is special - only exists for Gen2 V2 as Download Interface Version

### Generation Pattern Examples

```
TREP Pattern Analysis:
┌──────────┬──────────┬─────────────────┬─────────────────┐
│   TREP   │   Tag    │   Generation    │   Data Type     │
├──────────┼──────────┼─────────────────┼─────────────────┤
│   0x00   │  0x7600  │   Special       │ Download Iface  │
│   0x01   │  0x7601  │   Gen 1         │ Overview        │
│   0x02   │  0x7602  │   Gen 1         │ Activities      │
│   0x21   │  0x7621  │   Gen 2 V1      │ Overview        │
│   0x22   │  0x7622  │   Gen 2 V1      │ Activities      │
│   0x31   │  0x7631  │   Gen 2 V2      │ Overview        │
│   0x32   │  0x7632  │   Gen 2 V2      │ Activities      │
└──────────┴──────────┴─────────────────┴─────────────────┘
```

**Implementation**: See `tachounit/vu_tags.go`

## Data Structure Parsing

### Card Data Structures

Card data structures are defined in **Appendix 1 (Data Dictionary)** and **Appendix 2**:

- Each Elementary File has a specific structure
- ASN.1 encoding rules apply
- Certificate formats follow X.509 standards

### Vehicle Unit Data Structures

VU data structures are defined in **Appendix 1 (Data Dictionary)**:

- Structure names follow pattern: `Vu{DataType}{Generation}`
- Examples: `VuOverviewFirstGen`, `VuActivitiesSecondGen`
- Field definitions include data types and byte layouts

#### Generation-Specific Structure Selection

When parsing VU files, you must select the correct data structure based on the tag:

```go
func selectVuStructure(tag uint16) interface{} {
    switch tag {
    // Generation 1 structures
    case 0x7601: return &VuOverviewFirstGen{}
    case 0x7602: return &VuActivitiesFirstGen{}
    case 0x7603: return &VuEventsAndFaultsFirstGen{}
    case 0x7604: return &VuDetailedSpeedFirstGen{}
    case 0x7605: return &VuTechnicalDataFirstGen{}

    // Generation 2 Version 1 structures
    case 0x7621: return &VuOverviewSecondGen{}
    case 0x7622: return &VuActivitiesSecondGen{}
    case 0x7623: return &VuEventsAndFaultsSecondGen{}
    case 0x7624: return &VuDetailedSpeedSecondGen{}
    case 0x7625: return &VuTechnicalDataSecondGen{}

    // Generation 2 Version 2 structures
    case 0x7600: return &DownloadInterfaceVersion{}
    case 0x7631: return &VuOverviewSecondGenV2{}
    case 0x7632: return &VuActivitiesSecondGenV2{}
    case 0x7633: return &VuEventsAndFaultsSecondGenV2{}
    case 0x7635: return &VuTechnicalDataSecondGenV2{}

    default: return nil
    }
}
```

#### Key Differences Between Generations

| Aspect               | Generation 1     | Generation 2 V1   | Generation 2 V2            |
| -------------------- | ---------------- | ----------------- | -------------------------- |
| **TREP Range**       | 0x01-0x05        | 0x21-0x25         | 0x00, 0x31-0x35            |
| **Security**         | Part A (legacy)  | Part B (enhanced) | Part B (enhanced)          |
| **Data Structures**  | Simpler layouts  | Record arrays     | Extended record arrays     |
| **Certificates**     | First generation | Second generation | Second generation          |
| **Special Features** | -                | -                 | Download Interface Version |

**Critical**: Using the wrong generation structure will result in parsing errors or incorrect data interpretation.

## Parsing Challenges and Solutions

### 1. Inconsistent HTML Structure

**Problem**: Regulation HTML has varying layouts for tag definitions.

**Solution**: Focus on text patterns rather than HTML structure:

```regex
SID 76 Hex.*?TREP ((?:\d+(?:, \d+)*)|(?:\d+ or \d+(?:, \d+)*)) Hex
```

### 2. Missing Explicit Tables

**Problem**: VU tags are embedded in prose, not tabulated.

**Solution**: Manually extract patterns once, create definitive mapping with regulation traceability.

### 3. Complex Generation Mapping

**Problem**: Generation must be inferred from TREP ranges.

**Solution**: Use TREP value ranges to automatically determine generation and create appropriate Go identifiers.

## Implementation Guidelines

### File Type Detection

1. **Check minimum length** (at least 2 bytes for tag)
2. **Try card format first** (EF_ICC = 0x0002)
3. **Try VU TV format** (valid VuTag)
4. **Fallback to TRTP format** (legacy/alternative VU format)

### Generation Detection in Practice

```go
func detectVuGeneration(data []byte) (string, error) {
    if len(data) < 2 {
        return "", errors.New("insufficient data")
    }

    tag := binary.BigEndian.Uint16(data[0:2])
    trep := byte(tag & 0xFF) // Extract TREP (second byte)

    switch {
    case trep == 0x00:
        return "Gen2V2-DownloadInterface", nil
    case trep >= 0x01 && trep <= 0x05:
        return "Generation1", nil
    case trep >= 0x21 && trep <= 0x25:
        return "Generation2V1", nil
    case trep >= 0x31 && trep <= 0x35:
        return "Generation2V2", nil
    default:
        return "", fmt.Errorf("unknown TREP: 0x%02X", trep)
    }
}
```

### Mixed Generation Files

**Important**: A single VU file can contain multiple generations of data blocks. When parsing:

1. **Check each block individually** - don't assume the entire file is one generation
2. **Use appropriate structures** for each block based on its tag
3. **Handle security differences** - Gen1 uses different signature verification than Gen2

### Tag Parsing

1. **Use existing enums** (`tachocard.Tag`, `tachounit.VuTag`)
2. **Validate tags** with `.IsValid()` methods
3. **Trace back to regulation** sections for compliance

### Data Structure Parsing

1. **Identify file type** first
2. **Use appropriate scanner** (`tachocard.SplitFunc` for TLV, TV scanner for VU)
3. **Map tags to structures** using regulation-defined layouts
4. **Handle generation differences** (Gen1 vs Gen2 vs Gen2V2)

## Regulation Sections Reference

| Component            | Regulation Section               | Description                            |
| -------------------- | -------------------------------- | -------------------------------------- |
| **File Types**       | Appendix 7, Sections 2.3 & 3.4.2 | TLV vs TV format distinction           |
| **Card Tags**        | Appendix 2                       | Elementary File identifiers            |
| **VU Tags**          | Appendix 7, Section 2.2.6        | SID 76 + TREP combinations             |
| **Data Structures**  | Appendix 1                       | Complete data dictionary               |
| **Protocol Details** | Appendix 7                       | Download protocols and message formats |
| **Security**         | Appendix 11                      | Cryptographic specifications           |

## Tools and Generators

- **`tools/cmd/tachocard-tag-parser/`**: Extracts card tags from Appendix 2
- **`tools/cmd/tachounit-tag-parser/`**: Extracts TRTP/RDI from Appendix 7/8
- **`tachounit/vu_tags.go`**: Manually created VU TV tags with regulation traceability
- **`filetype.go`**: Unified file type detection using both tag systems

## Real-World Analysis Findings

### Test File Analysis Results

Analysis of actual tachograph card files reveals important patterns:

#### File Structure Variations

- **Simple Gen1 Cards**: 22 TLV records, ~26KB
- **Mixed Generation Cards**: 54 TLV records, ~67KB (both Gen1 and Gen2 data)
- **Signature Patterns**: Gen1 uses 128-byte RSA signatures, Gen2 uses 64-byte ECC signatures

#### Generation Migration Pattern

Real cards demonstrate the **"Migration" scenario** described in **Appendix 15**:

1. **First Section**: Complete Gen1 data set (all standard Elementary Files)
2. **Second Section**: Complete Gen2 data set (enhanced versions of same files)
3. **Coexistence**: Both generations present in single file for backward compatibility

#### Unknown FID Discovery Process

1. **Identify Unknown Tags**: Use CLI tools to scan real card files
2. **Cross-Reference**: Check benchmark implementations for tag definitions
3. **Regulation Mapping**: Trace back to specific appendices (especially Appendix 12 for GNSS)
4. **Validation**: Verify tag patterns match expected generation/signature structure

### Implementation Validation

The CLI implementation successfully handles:

- ✅ **Mixed Generation Files**: Correctly parses both Gen1 and Gen2 records
- ✅ **Extended FIDs**: Recognizes 0x052x and 0xC10x series tags
- ✅ **Generation Detection**: Properly identifies data vs signature records
- ✅ **GNSS Integration**: Handles GNSS facility certificates (0xC10x)
- ✅ **Regulation Compliance**: 5-byte TLV headers with 3-byte tags

## Future Considerations

1. **Regulation Updates**: Monitor for new TREP values or data types
2. **Generation 3**: Watch for new tag ranges beyond 0x30-0x3F
3. **Parser Improvements**: Could automate VU tag extraction with better HTML parsing
4. **Validation**: Cross-reference with official test vectors when available
5. **Extended FID Discovery**: Continue analysis of real-world files to identify additional proprietary FIDs
6. **GNSS Evolution**: Monitor for additional GNSS-related FIDs as OSNMA deployment expands

This documentation provides the foundation for understanding and maintaining the tachograph file parsing implementation while ensuring compliance with EU regulations and handling real-world file variations.
