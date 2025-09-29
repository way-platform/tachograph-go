# Work Log: Fix Generation-Aware Parsing for EF_Places

**Date:** 2025-09-29

## Objective

Correct a critical bug in `EF_Places` parsing where the code did not correctly handle differences between Gen1 and Gen2 card data. The existing implementation was also non-compliant with the specification for Gen1 records, leading to incorrect data parsing.

## Analysis

An investigation into `card_places.go` revealed that the parser for `EF_Places` incorrectly assumed a fixed record size of 12 bytes.

This was wrong for two reasons:

1.  **Gen2 Incompatibility**: Gen2 cards use a larger 22-byte record to include GNSS data, which the parser could not handle.
2.  **Gen1 Non-Compliance**: The official specification for Gen1 defines a 10-byte record. The existing code was incorrectly parsing a 12-byte structure (by assuming a 2-byte region and an extra reserved byte), which did not align with the source material.

This bug resulted in incorrect data for all Gen1 cards and a complete failure to parse place records from Gen2 cards.

## Resolution

To address this, a full refactoring of the `EF_Places` handling was performed, guided by the principle of adhering strictly to the source specification.

1.  **Principle Update (`AGENTS.md`)**: A new principle was added to `AGENTS.md` to emphasize the importance of generation-aware parsing, using this issue as a key example of why it is critical.

2.  **Protobuf Schema (`places.proto`)**:
    *   The comments in `places.proto` were updated to accurately document the spec-compliant binary structures: **10 bytes for a Gen1 `PlaceRecord`** and **22 bytes for a Gen2 `PlaceRecord_G2`**.
    *   A new `GnssPlaceRecord` message was added to encapsulate the GNSS data.
    *   An `optional GnssPlaceRecord gnss_place_record` field was added to the `Record` message to support Gen2 data.

3.  **Go Implementation (`card_places.go`)**:
    *   The parsing and marshalling logic in `card_places.go` was completely refactored.
    *   The incorrect hardcoded 12-byte record size was removed.
    *   The `unmarshalCardPlaces` and `appendPlaces` functions were updated to require a `generation` parameter, making them generation-aware.
    *   The core `unmarshalPlaceRecord` and `appendPlaceRecord` functions now contain distinct logic paths to correctly handle the 10-byte Gen1 records and 22-byte Gen2 records, ensuring full compliance with the specification.

### Files Modified

-   `AGENTS.md`
-   `proto/wayplatform/connect/tachograph/card/v1/places.proto`
-   `card_places.go`
