# File Structure Migration - Unresolved Issues

**Date**: 2025-01-27
**Status**: ✅ **MIGRATION COMPLETED** - All file structure changes implemented successfully

## 🎉 **Migration Summary**

**File Structure Migration**: ✅ **COMPLETED** (100%)

- All files migrated from separate `unmarshal_*` and `append_*` files to consolidated `(vu|card|dd)_<typename>.go` structure
- Function visibility fixed (only `UnmarshalFile` and `MarshalFile` public)
- Data dictionary consolidation completed (34/34 types)

**Test Status**: ✅ All tests passing - no regressions introduced

## 🚨 **Unresolved Issues Requiring Attention**

### 1. ASN.1 Documentation Inconsistencies

**Priority**: Medium
**Files Affected**: Multiple card and VU files

**Issue**: ASN.1 definitions in comments don't match between unmarshal and append functions:

- `card_events.go`: Unmarshal shows 24-byte records, append shows 35-byte records
- `card_faults.go`: Same pattern as events - 24 vs 35 byte documentation mismatch
- `card_identification.go`: Unmarshal uses 14-byte card number, append shows 16-byte
- `card_vehicles.go`: Unmarshal shows 3-byte odometer, append shows 4-byte

**Resolution Needed**: Verify correct ASN.1 definitions and standardize documentation across all files.

### 2. Incomplete Append Implementations

**Priority**: High
**Files Affected**: Multiple VU and card files

**Issue**: Several append functions are severely incomplete or simplified:

- `vu_activities.go`: Append only writes signature data, missing all other fields
- `vu_events_faults.go`: Simplified implementation, only handles signature
- `vu_technical_data.go`: Simplified implementation, only handles signature
- `vu_detailed_speed.go`: Simplified implementation, only handles signature
- `card_gnss_places.go`: Basic implementation, skips complex record structures
- `card_vehicle_units_used.go`: Basic implementation, skips complex record structures
- `vu_vehicle_unit_file.go`: `appendVU` function is placeholder that returns unchanged data

**Resolution Needed**: Complete all append implementations to match the complexity of their corresponding unmarshal functions.

### 3. Incomplete Gen2 Implementations

**Priority**: Medium
**Files Affected**: VU files

**Issue**: Gen2 support is incomplete or missing:

- `vu_overview.go`: Gen2 functions are placeholders that do nothing
- `vu_activities.go`: Gen2 support exists but helper functions are simplified

**Resolution Needed**: Implement complete Gen2 support for all VU file types.

### 4. Simplified Helper Functions

**Priority**: Medium
**Files Affected**: VU activities

**Issue**: Many helper functions in `vu_activities.go` are marked as "simplified implementation" and return empty data:

- `parseVuCardIWData`, `parseVuActivityDailyData`, `parseVuPlaceDailyWorkPeriodData`, etc.
- These functions have TODO-style comments indicating they need proper implementation

**Resolution Needed**: Complete all helper functions to provide full VU activities parsing.

### 5. Constants Naming Convention

**Priority**: Low
**Files Affected**: Multiple files

**Issue**: Inconsistent constant naming conventions across files:

- Some files use `cardDriverActivityHeaderSize` vs `lenCardDriverActivityHeader`
- AGENTS.md guidelines suggest using `idx` prefix for offsets and `len` for lengths

**Resolution Needed**: Standardize constant naming to follow the `idx`/`len` prefix convention consistently.

### 6. Incomplete Card Number Handling

**Priority**: Medium
**Files Affected**: `card_identification.go`

**Issue**: The append function has placeholder logic for handling CardNumber CHOICE:

- Comment says "This needs to be implemented based on the actual card type"
- Logic tries to build a string from different identification types but may not handle all cases properly

**Resolution Needed**: Implement proper CardNumber handling logic to match the unmarshal function's behavior.

## 📋 **Recommended Resolution Order**

1. **High Priority**: Complete all incomplete append implementations
2. **Medium Priority**: Fix ASN.1 documentation inconsistencies
3. **Medium Priority**: Implement missing Gen2 support
4. **Medium Priority**: Complete simplified helper functions
5. **Low Priority**: Standardize constants naming convention

## 📊 **Impact Assessment**

- **Functionality**: Core parsing works, but marshalling is incomplete for many types
- **Roundtrip Support**: Limited due to incomplete append implementations
- **Gen2 Support**: Partial implementation may cause issues with newer tachograph files
- **Maintenance**: Documentation inconsistencies make code harder to maintain

## 🎯 **Next Steps**

1. Prioritize completing append implementations for critical file types
2. Verify ASN.1 definitions against the regulation specification
3. Implement missing Gen2 support for VU files
4. Complete simplified helper functions
5. Standardize documentation and naming conventions
