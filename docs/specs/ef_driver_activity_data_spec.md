# Specification: EF_DRIVER_ACTIVITY_DATA Ring Buffer Implementation

## Problem Statement

After achieving **95.5% success rate (21/22 records)** with byte-perfect roundtrip marshalling, the final remaining issue is **EF_DRIVER_ACTIVITY_DATA**. This document provides comprehensive analysis of the ring buffer implementation challenge that prevents 100% roundtrip accuracy.

## Current Achievement Context

### ✅ Successfully Resolved (21/22 records)

- **All signature records**: Perfect preservation using multi-pass TLV architecture
- **EF_EVENTS_DATA & EF_FAULTS_DATA**: Tagged union approach for padding preservation
- **EF_CONTROL_ACTIVITY_DATA**: Tagged union approach for string field padding
- **EF_SPECIFIC_CONDITIONS**: Enum protocol value conversion
- **EF_VEHICLES_USED**: Vehicle registration nation field conversion
- **EF_PLACES**: Added reserved_byte and trailing_bytes fields
- **EF_CARD_CERTIFICATE & EF_CA_CERTIFICATE**: Fixed ordering and content
- **All other standard EFs**: Perfect matches

### 🔧 Proven Techniques Available

1. **Tagged Union Pattern**: Raw byte preservation for complex padding
2. **Enum Protocol Value Conversion**: Proper enum field marshalling
3. **Multi-pass TLV Architecture**: Signature preservation
4. **Field-by-field Analysis**: Systematic byte difference resolution

## EF_DRIVER_ACTIVITY_DATA: The Final Challenge

### Severity Assessment

- **Impact**: **CRITICAL** - Prevents 100% roundtrip accuracy
- **Complexity**: **HIGH** - Complex ring buffer with variable-length records
- **Data Volume**: **MAJOR** - 13,352 byte difference (52% of total file size)
- **Regulatory Importance**: **ESSENTIAL** - Activity data is core compliance requirement

### Current Behavior Analysis

#### Binary Roundtrip Results

```
Test File: proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD
Original file size:    26,145 bytes
Marshalled file size:  12,793 bytes
Size difference:       13,352 bytes (51% data loss)

First binary difference: byte 3631
Context: CD34621EEA05040035D4 vs CD34621EEA05040001AC
         └─ Tag ─┘└─App─┘└─Length─┘    └─ Tag ─┘└─App─┘└─Length─┘
         0x0504   0x00   13,780 bytes   0x0504   0x00   428 bytes
```

#### Semantic Roundtrip Results

```
Record 10 (EF_DRIVER_ACTIVITY_DATA):
Original:   Tag=0x050400, Length=13780 bytes
Marshalled: Tag=0x050400, Length=428 bytes
Status: ❌ Issues: [Length Value]
First difference at byte 4: original=0x67, marshalled=0x00
```

#### Hex Data Analysis

```
Original (first 32 bytes):   2A2C298A67D0CE800627005C600000F218F700F8192A112C1930013F114E194F
Marshalled (first 32 bytes): 2A2C298A000E000E00000000000000000000000E000E00000000000000000000
                              ────────└─ Divergence starts here (byte 4)
```

### Root Cause Analysis

#### Current Implementation Limitations

Looking at `unmarshal_card_activity.go` and `append_card_activity.go`:

1. **Incomplete Ring Buffer Parsing**:

   ```go
   // Current parseActivityDailyRecord creates minimal empty records
   func parseActivityDailyRecord(data []byte, startIndex int) (*cardv1.DriverActivity_DailyRecord, int, error) {
       // Only reads header, creates empty record with basic metadata
       // MISSING: Actual activity change parsing
       // MISSING: Variable-length record content
   ```

2. **Simplified Marshalling**:

   ```go
   // Current AppendActivityDailyRecord handles empty records
   isEmpty := rec.GetActivityDayDistance() == 0 &&
       len(rec.GetActivityChangeInfo()) == 0 &&
       (rec.GetActivityRecordDate() == nil || rec.GetActivityRecordDate().GetSeconds() == 0)
   ```

3. **Data Structure Mismatch**:
   - **Expected**: Complex ring buffer with variable-length activity records
   - **Current**: Simple array of mostly empty records
   - **Result**: 13,352 bytes of activity data completely lost

### Technical Deep Dive

#### Ring Buffer Structure (from regulation)

```
EF_Driver_Activity_Data structure:
- Newest record index (2 bytes)
- Oldest record index (2 bytes)
- Activity daily records (variable length, up to ~13KB)

Each Daily Record:
- Previous record length (2 bytes)
- Current record length (2 bytes)
- Activity record date (4 bytes)
- Daily presence counter (2 bytes BCD)
- Activity day distance (2 bytes)
- Activity change info (variable length sequence)
```

#### Activity Change Info Structure

```
Each Activity Change:
- Slot (1 byte): Driver/Co-driver slot
- Status (1 byte): Available/Break/Work/Drive/Unknown
- Card inserted (1 byte): Yes/No
- Driving status (1 byte): Crew/Single/Not known
- Vehicle speed (1 byte): 0-250 km/h or special values
- Time (2 bytes): Minutes since 00:00 UTC
```

#### Ring Buffer Navigation Challenge

```
Original binary data analysis:
Byte 0-1:   2A2C = Newest record index (10796)
Byte 2-3:   298A = Oldest record index (10634)
Byte 4-5:   67D0 = First record previous length (26576)
Byte 6-7:   CE80 = First record current length (52864)
...

Current marshalled output:
Byte 0-1:   2A2C = Newest record index (preserved)
Byte 2-3:   298A = Oldest record index (preserved)
Byte 4-5:   000E = First record previous length (14) - WRONG
Byte 6-7:   000E = First record current length (14) - WRONG
```

### Proposed Solution Architecture

#### Phase 1: Enhanced Ring Buffer Parsing

1. **Implement proper ring buffer navigation**:

   - Calculate actual buffer positions using indices
   - Handle wraparound cases correctly
   - Parse records in chronological order

2. **Parse complete activity change sequences**:

   - Read variable number of 6-byte activity changes
   - Convert to protobuf ActivityChange messages
   - Preserve exact timing and status data

3. **Calculate accurate record lengths**:
   - Base length calculation on actual parsed content
   - Match original variable-length records
   - Ensure marshalled lengths match original

#### Phase 2: Byte-Perfect Marshalling

1. **Implement reverse ring buffer construction**:

   - Convert protobuf records back to ring buffer format
   - Maintain original cyclic index relationships
   - Preserve exact byte sequences

2. **Activity change serialization**:

   - Convert ActivityChange messages to 6-byte sequences
   - Maintain original activity timing
   - Preserve enum protocol values

3. **Length field accuracy**:
   - Calculate previous/current record lengths correctly
   - Ensure total buffer size matches original
   - Handle variable-length record marshalling

#### Phase 3: Tagged Union Fallback (if needed)

If full ring buffer implementation proves too complex:

1. **Apply tagged union pattern**:
   - Add `valid` boolean and `raw_data` bytes field
   - Preserve original 13,780 bytes as raw data
   - Ensure perfect roundtrip without semantic parsing

### Success Criteria

#### Functional Requirements

1. **Length Match**: Marshalled length = 13,780 bytes (exact match)
2. **Content Match**: First difference beyond byte 3631 (current limit)
3. **Index Preservation**: Newest/oldest indices match original
4. **Activity Data**: All activity changes correctly parsed and marshalled

### Conclusion

EF_DRIVER_ACTIVITY_DATA represents the final and most complex challenge in achieving 100% roundtrip accuracy. The ring buffer structure with variable-length records and cyclic indexing requires careful implementation, but the systematic approach that achieved 95.5% success provides a proven foundation.

The tagged union fallback ensures that even if full semantic parsing proves challenging, perfect roundtrip can still be achieved through raw byte preservation. This guarantees that 100% accuracy is achievable regardless of implementation complexity.

Success in this final challenge will complete an exceptional engineering achievement: byte-perfect roundtrip marshalling for complex regulatory tachograph data with full semantic understanding.

## Final Implementation Plan: Full Semantic Parsing

This section supersedes the previous "Raw Data Preservation" plan. The new objective is a complete and correct semantic implementation of the EF_DRIVER_ACTIVITY_DATA ring buffer to achieve a byte-perfect roundtrip while also providing meaningful data.

### 1. Analysis of `ActivityChangeInfo` Structure

A key challenge is the ambiguity in the `ActivityChangeInfo` structure. The existing codebase contains conflicting implementations:
- `unmarshal_card_activity.go` attempts to parse a **4-byte** structure.
- `append_card_activity.go` serializes a **2-byte** bitfield.

After careful review, the **2-byte bitfield structure is determined to be correct**. It is a compact and common pattern in embedded systems and tachograph data. The 4-byte and 6-byte interpretations are considered incorrect. The implementation will proceed with the following 2-byte structure for `ActivityChangeInfo`:

- **Bit 15**: `Slot` (0 = Driver, 1 = Co-driver)
- **Bit 14**: `Driving Status` (0 = Single, 1 = Crew)
- **Bit 13**: `Card Status` (0 = Not inserted, 1 = Inserted)
- **Bits 11-12**: `Activity` (0 = Rest/Break, 1 = Available, 2 = Work, 3 = Driving)
- **Bits 0-10**: `Time of Change` (Minutes since 00:00 UTC)

### 2. Detailed Implementation Steps

#### Step 1: `unmarshal_card_activity.go` — Full Parsing Logic

The entire file will be rewritten to correctly parse the activity data records linearly.

1.  **`UnmarshalCardActivityData` function:**
    *   Reads the 2-byte `oldestDayRecordPointer` and 2-byte `newestDayRecordPointer` from the start of the data.
    *   Passes the remaining byte slice (`activityData`) to `parseActivityDailyRecords`.

2.  **`parseActivityDailyRecords` function:**
    *   This function will iterate through the `activityData` slice as a simple sequence of variable-length records, not a cyclic buffer.
    *   In a loop, it reads `currentRecordLength` from the header of each record to determine its boundary.
    *   It slices the data for the current record and passes it to `parseActivityDailyRecord`.
    *   The loop continues, advancing the offset by `currentRecordLength` after each record, until all data is consumed.

3.  **`parseActivityDailyRecord` function:**
    *   Takes a byte slice representing a single daily record.
    *   Parses the fixed-size content: `activityRecordDate` (4 bytes, BCD), `activityDailyPresenceCounter` (2 bytes, BCD), and `activityDayDistance` (2 bytes).
    *   Loops through the remainder of the slice in 2-byte chunks, parsing each chunk as an `ActivityChangeInfo` bitfield according to the structure defined above.
    *   Populates and returns a complete `DailyRecord` protobuf message.

#### Step 2: `append_card_activity.go` — Reconstructive Marshalling Logic

This file will be rewritten to perfectly serialize the semantic data back into the correct binary format.

1.  **`AppendDriverActivity` function:**
    *   Appends the 2-byte `oldest_day_record_index` and 2-byte `newest_day_record_index`.
    *   Loops through the `DailyRecord` messages, calling `AppendActivityDailyRecord` for each.

2.  **`AppendActivityDailyRecord` function:**
    *   This function is critical for correctness. It will first serialize the *content* of the record (date, distance, all activity changes) into a temporary buffer.
    *   The size of this temporary buffer plus the 4-byte header gives the `currentRecordLength`.
    *   It then appends the `activityPreviousRecordLength` (from the proto field) and the newly calculated `currentRecordLength` to the destination buffer.
    *   Finally, it appends the content from the temporary buffer.

3.  **`AppendActivityChange` function:**
    *   This helper function will take an `ActivityChange` message and construct the precise 2-byte bitfield, which is then appended to the temporary buffer in `AppendActivityDailyRecord`.

