# 2025-09-28: Missing VU Proto Schemas Analysis

This log documents the missing protobuf schemas that are preventing complete VU (Vehicle Unit) data parsing and marshalling implementation.

## Overview

During the code cleanup phase, we discovered that our VU implementation is incomplete due to missing protobuf data structures. While we have basic VU structures like `Overview`, we are missing several critical VU-specific data types that are required by the EU digital tachograph regulation.

## Current VU Proto Files

We currently have these VU proto files:

- `activities.proto` - VU activities data
- `card_download.proto` - Card download functionality
- `detailed_speed.proto` - Detailed speed data
- `download_interface_version.proto` - Download interface version
- `events_and_faults.proto` - Events and faults data
- `overview.proto` - VU overview data
- `technical_data.proto` - Technical data
- `transfer_type.proto` - Transfer type definitions
- `trep.proto` - TREP (Transfer Record Element Protocol) definitions
- `vehicle_unit_file.proto` - Main VU file structure
- `versioning.proto` - Version information

## Missing VU Data Structures

### 1. **VuDownloadActivityData** (Critical Missing)

**Regulation Reference:** Data Dictionary Section 2.195

**Purpose:** Information stored in a vehicle unit related to its last download (Annex 1B requirement 105 and Annex 1C requirement 129).

**ASN.1 Definition:**

```
VuDownloadActivityData ::= SEQUENCE {
    downloadingTime TimeReal,
    fullCardNumber FullCardNumber,
    companyOrWorkshopName Name
}
```

**Generation Differences:**

- **Gen1:** Uses `fullCardNumber` directly
- **Gen2:** Uses `fullCardNumberAndGeneration` instead

**Impact:** This structure is referenced in our VU overview parsing code but doesn't exist in our proto schemas, causing compilation errors.

### 2. **VuDownloadActivityDataRecordArray** (Critical Missing)

**Regulation Reference:** Data Dictionary Section 2.196

**Purpose:** Generation 2 version of download activity data with record array format.

**ASN.1 Definition:**

```
VuDownloadActivityDataRecordArray ::= SEQUENCE {
    recordType RecordType,
    recordSize INTEGER(0..2^16-1),
    noOfRecords INTEGER(0..2^16-1),
    records SET SIZE(noOfRecords) OF VuDownloadActivityData
}
```

**Impact:** Required for Gen2 VU data parsing.

### 3. **VuCompanyLocksData** (Partially Missing)

**Regulation Reference:** Data Dictionary Section 2.184

**Current State:** We have `CompanyLock` message in `overview.proto` but missing the main `VuCompanyLocksData` structure.

**Missing ASN.1 Definition:**

```
VuCompanyLocksData ::= SEQUENCE {
    companyLocks SET SIZE(NoOfCompanyLocks) OF VuCompanyLocksRecord
}
```

**Impact:** VU overview parsing references this structure but it's not properly modeled.

### 4. **VuControlActivityData** (Partially Missing)

**Regulation Reference:** Data Dictionary Section 2.185

**Current State:** We have `ControlActivity` message in `overview.proto` but missing the main `VuControlActivityData` structure.

**Missing ASN.1 Definition:**

```
VuControlActivityData ::= SEQUENCE {
    controlActivities SET SIZE(NoOfControlActivityRecords) OF VuControlActivityRecord
}
```

**Impact:** VU overview parsing references this structure but it's not properly modeled.

### 5. **VuActivitiesData** (Missing Detailed Structures)

**Regulation Reference:** Data Dictionary Section 2.2.6.2

**Current State:** We have basic `Activities` message but missing detailed sub-structures.

**Missing Structures:**

- `VuCardIWData` - Card insertion/withdrawal data
- `VuActivityDailyData` - Daily activity data
- `VuPlaceDailyWorkPeriodData` - Place and work period data
- `VuSpecificConditionData` - Specific condition data
- `VuGNSSADRecordArray` - GNSS accumulated driving records (Gen2+)
- `VuBorderCrossingRecordArray` - Border crossing records (Gen2v2+)
- `VuLoadUnloadRecordArray` - Load/unload operation records (Gen2v2+)

### 6. **VuEventsAndFaultsData** (Missing Detailed Structures)

**Regulation Reference:** Data Dictionary Section 2.2.6.3

**Current State:** We have basic `EventsAndFaults` message but missing detailed sub-structures.

**Missing Structures:**

- `VuEventData` - Event data
- `VuFaultData` - Fault data
- `VuOverSpeedingEventData` - Over-speeding event data
- `VuTimeAdjustmentData` - Time adjustment data
- `VuEventRecordArray` - Event record arrays (Gen2)
- `VuFaultRecordArray` - Fault record arrays (Gen2)

### 7. **VuDetailedSpeedData** (Missing Detailed Structures)

**Regulation Reference:** Data Dictionary Section 2.2.6.4

**Current State:** We have basic `DetailedSpeed` message but missing detailed sub-structures.

**Missing Structures:**

- `VuDetailedSpeedBlock` - Detailed speed blocks
- `VuDetailedSpeedBlockRecordArray` - Speed block record arrays (Gen2)

### 8. **VuTechnicalData** (Missing Detailed Structures)

**Regulation Reference:** Data Dictionary Section 2.2.6.5

**Current State:** We have basic `TechnicalData` message but missing detailed sub-structures.

**Missing Structures:**

- `VuIdentification` - VU identification data
- `VuCalibrationData` - Calibration data
- `VuCardData` - Card data
- `VuIdentificationRecordArray` - Identification record arrays (Gen2)
- `VuCalibrationRecordArray` - Calibration record arrays (Gen2)
- `VuCardRecordArray` - Card record arrays (Gen2)

## Compilation Errors Due to Missing Schemas

The missing schemas are causing several compilation errors:

1. **`undefined: vuv1.DownloadActivityData`** - Missing download activity data structure
2. **Missing helper functions** - Functions that depend on missing proto types
3. **Incomplete VU parsing** - Can't parse VU files completely without these structures

## Impact on VU Implementation

### **Current Limitations:**

1. **Incomplete VU Parsing** - We can only parse basic VU overview data, not the detailed VU-specific structures
2. **No Roundtrip Testing** - Can't marshal VU data back to binary without complete data structures
3. **Missing VU Functionality** - Can't implement full VU data processing as required by the regulation

### **Regulation Compliance:**

The missing structures are **required by the EU digital tachograph regulation** for:

- Complete VU data parsing (Annex 1B requirement 105)
- VU download functionality (Annex 1C requirement 129)
- Full compliance with the data dictionary specifications

## Recommended Action Plan

### **Phase 1: Critical Missing Structures (High Priority)**

1. **Add `VuDownloadActivityData`** to `overview.proto` or create new `download_activity.proto`
2. **Add `VuDownloadActivityDataRecordArray`** for Gen2 support
3. **Complete `VuCompanyLocksData`** structure
4. **Complete `VuControlActivityData`** structure

### **Phase 2: Detailed VU Structures (Medium Priority)**

1. **Expand `activities.proto`** with missing activity-related structures
2. **Expand `events_and_faults.proto`** with missing event/fault structures
3. **Expand `detailed_speed.proto`** with missing speed structures
4. **Expand `technical_data.proto`** with missing technical structures

### **Phase 3: Generation-Specific Support (Lower Priority)**

1. **Add Gen2-specific record arrays** for all VU data types
2. **Add Gen2v2-specific structures** (border crossing, load/unload)
3. **Add proper versioning support** for different VU generations

## Implementation Notes

### **ASN.1 to Protobuf Mapping:**

When implementing these structures, follow these principles:

- Use `google.protobuf.Timestamp` for `TimeReal` fields
- Use `repeated` for `SET SIZE(n) OF` constructs
- Use `optional` for `OPTIONAL` fields
- Use `oneof` for `CHOICE` constructs
- Add proper field comments with ASN.1 references

### **File Organization:**

Consider creating separate proto files for:

- `download_activity.proto` - Download activity data
- `company_locks.proto` - Company locks data
- `control_activity.proto` - Control activity data
- `gnss.proto` - GNSS-related data (Gen2+)

### **Testing Strategy:**

1. **Add test data** for each new VU structure
2. **Implement roundtrip testing** to ensure data integrity
3. **Validate against regulation** to ensure compliance
4. **Use benchmark implementation** as reference for complex structures

## Conclusion

The missing VU proto schemas represent a **critical gap** in our implementation that prevents us from achieving full VU data processing capabilities. These structures are not optional - they are **required by the EU regulation** for complete tachograph data processing.

**Priority:** This should be addressed **immediately** after fixing the current compilation errors, as it's blocking the core VU functionality that was identified as a major goal in the original audit.

**Estimated Impact:** Completing these schemas will enable:

- Full VU file parsing and marshalling
- Complete roundtrip testing
- Full compliance with EU digital tachograph regulation
- Achievement of the "full binary roundtrip parsing" goal

---

## Files Referenced

- `docs/regulation/chapters/03-data-dictionary.md` - ASN.1 specifications
- `proto/wayplatform/connect/tachograph/vu/v1/` - Current VU proto files
- `benchmark/tachoparser/` - Reference implementation
- `unmarshal_vu_*.go` - VU parsing code that references missing structures
- `append_vu_*.go` - VU marshalling code that references missing structures

---

## 2025-09-28: Resolution Summary

This issue was analyzed and resolved. The key findings and actions are summarized below:

1.  **Confirmation of Analysis:** The analysis that `VuDownloadActivityData` and its corresponding `RecordArray` for Gen2 were missing was confirmed to be correct. The analysis concerning other detailed structures (`Activities`, `EventsAndFaults`, etc.) was found to be outdated, as those structures were already correctly implemented using `repeated` fields to handle both Gen1 and Gen2 data.

2.  **Design Policy Documentation:** To prevent future misunderstandings, the project's design policy of using `repeated` fields to model both single (Gen1) and multiple (Gen2) records was explicitly documented in the main `AGENTS.md` file.

3.  **Schema Implementation:**
    *   A new `message DownloadActivity` was added to `overview.proto` to represent the `VuDownloadActivityData` structure.
    *   The `Overview` message was updated to replace the temporary `last_download_time` field with `repeated DownloadActivity download_activities`. This aligns the implementation with the newly documented design policy and provides a complete data model for download activities.

4.  **Verification:** After the changes, the protobuf Go files were successfully regenerated, confirming the schema is healthy and free of errors. The primary gap identified in this analysis is now closed.
