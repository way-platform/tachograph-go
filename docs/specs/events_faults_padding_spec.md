# Specification: EF_EVENTS_DATA and EF_FAULTS_DATA Padding Mismatch

## Problem Statement

During roundtrip testing, EF_EVENTS_DATA and EF_FAULTS_DATA exhibit complex padding inconsistencies that prevent byte-perfect marshalling. The issue manifests as the first binary difference in roundtrip tests.

## Current Behavior

### Original File Pattern (from binary analysis)

```
Bytes 0-10:   00 00 00 00 00 00 00 00 00 00 00    (null bytes - likely header/real data)
Bytes 11-23:  20 20 20 20 20 20 20 20 20 20 20 20 20  (space-padded empty record)
Bytes 24-34:  00 00 00 00 00 00 00 00 00 00 00    (null bytes - likely header/real data)
Bytes 35-47:  20 20 20 20 20 20 20 20 20 20 20 20 20  (space-padded empty record)
... (pattern continues)
```

### Marshalled Output Pattern

```
All padding records: 00 00 00 00 00 00 00 00 00 00 00 00 00  (null-padded)
```

### Test Results

- **EF_EVENTS_DATA**: Original=1728 bytes, Marshalled=1728 bytes, **First difference at byte 11: original=0x20, marshalled=0x00**
- **EF_FAULTS_DATA**: Original=1152 bytes, Marshalled=1152 bytes, **First difference at byte 11: original=0x20, marshalled=0x00**

## Root Cause Analysis

### 1. Unmarshalling Logic Issue

The current empty record detection in `unmarshal_card_events.go` and `unmarshal_card_faults.go`:

```go
// Current logic - recognizes both null and space padding as empty
isEmpty := true
for _, b := range recordData {
    if b != 0 && b != 0x20 { // Allow both null and space padding
        isEmpty = false
        break
    }
}
if isEmpty {
    continue // Skip empty records
}
```

**Problem**: This logic treats space-padded records as empty and skips them during unmarshalling, losing the original padding information.

### 2. Marshalling Logic Issue

The current marshalling in `raw_card_file.go`:

```go
} else {
    // Pad with an empty 24-byte record
    eventsValBuf = append(eventsValBuf, make([]byte, 24)...) // Creates null padding
}
```

**Problem**: Always uses null padding for empty records, regardless of original format.

### 3. Data Structure Complexity

The original file appears to have:

- **Real record data**: May contain null bytes as legitimate data
- **Empty record padding**: Uses space bytes (0x20) for unused record slots
- **Mixed patterns**: Some sections have null bytes, others have space padding

## Attempted Solutions and Results

### Attempt 1: Change marshalling to use space padding

**Result**: All records became space-padded, causing opposite mismatch (original=0x00, marshalled=0x20)

### Attempt 2: Enhanced empty record detection

**Result**: Improved detection but didn't solve the fundamental issue of preserving original padding style

## Technical Challenges

1. **Information Loss**: Current unmarshalling strips padding information by skipping empty records
2. **Pattern Complexity**: Original file has mixed null/space patterns that aren't easily predictable
3. **Record Structure**: 24-byte records with complex internal structure make pattern detection difficult
4. **Roundtrip Requirement**: Must preserve exact original bytes for perfect roundtrip

## Impact

- **Binary Roundtrip**: First difference at byte 477-488 in test files
- **Semantic Roundtrip**: EF_EVENTS_DATA and EF_FAULTS_DATA show value mismatches
- **File Size**: Correct (1728/1152 bytes) but content differs

## Proposed Solution Approaches

### Option 1: Preserve Original Padding During Unmarshalling

- Store original padding style in protobuf messages
- Use preserved padding during marshalling
- **Pros**: Most accurate, handles complex patterns
- **Cons**: Requires protobuf schema changes, architectural complexity

### Option 2: Raw Byte Preservation for Empty Records

- Keep original raw bytes for detected empty records
- Use raw bytes during marshalling instead of generating padding
- **Pros**: No schema changes, preserves exact original
- **Cons**: Need to identify empty vs real records accurately

### Option 3: Pattern Analysis and Recreation

- Analyze original file structure to understand padding patterns
- Implement logic to recreate original patterns during marshalling
- **Pros**: No schema changes, algorithmic approach
- **Cons**: May be brittle, hard to generalize across different files

### Option 4: Hybrid Approach

- Use raw byte preservation for padding sections
- Keep current logic for real data sections
- **Pros**: Balanced complexity and accuracy
- **Cons**: Need clear boundary detection between padding and data

## Files Affected

- `unmarshal_card_events.go` - Event record parsing
- `unmarshal_card_faults.go` - Fault record parsing
- `raw_card_file.go` - Event/fault marshalling logic
- `append_card_events.go` - Event record serialization
- `append_card_faults.go` - Fault record serialization

## Test Cases

- Binary roundtrip test: First difference point tracking
- Semantic roundtrip test: Value comparison for EF_EVENTS_DATA/EF_FAULTS_DATA
- Multiple test files to verify pattern consistency

## Success Criteria

1. **Binary Perfect**: No differences in EF_EVENTS_DATA and EF_FAULTS_DATA sections
2. **Pattern Preservation**: Original null/space padding patterns maintained
3. **Semantic Correctness**: Parsed records remain unchanged
4. **Generalization**: Solution works across different card files

## Priority

**Medium-High**: This is the first binary difference in roundtrip tests, but the signature preservation was higher priority and is now complete. Should be addressed after simpler content issues (EF_VEHICLES_USED, EF_PLACES) are resolved.
