# Work Log: Refactoring Parsing with bufio.Scanner

**Date:** 2025-09-29

## Objective

Refactor the parsing of repeated records to use `bufio.Scanner` with custom `SplitFunc` implementations. This will decouple data scanning from data parsing, improving code clarity, robustness, and adherence to idiomatic Go practices.

## Analysis

An analysis of the current codebase identified several areas where repeated records are parsed. These can be categorized into three groups:

1.  **Fixed-Size Contiguous Records:** Ideal candidates for `bufio.Scanner`.
    *   `card_events.go`: `unmarshalEventsData` parses a sequence of 24-byte `CardEventRecord`s.
    *   `card_faults.go`: `unmarshalFaultsData` parses a sequence of 24-byte `CardFaultRecord`s.
    *   `card_places.go`: `unmarshalCardPlaces` parses a sequence of 12-byte `PlaceRecord`s after an initial pointer.

2.  **Length-Prefixed Record Arrays:** Good candidates where a scanner can be used on a sub-slice of the data.
    *   `vu_activities.go` (Gen1): Functions like `parseVuCardIWData` read a count and then loop.
    *   `vu_activities.go` (Gen2): "Record Array" structures explicitly provide `noOfRecords` and `recordSize`.

3.  **Non-Sequential Records:** Not suitable for `bufio.Scanner`.
    *   `card_activity.go`: `parseCyclicActivityDailyRecords` parses a cyclic buffer with backward pointers, which is incompatible with the forward-scanning nature of `bufio.Scanner`.

## Implementation Plan

The refactoring will be implemented in the following order:

1.  **`card_events.go`**:
    *   Create a `splitCardEventRecord` function that returns 24-byte chunks.
    *   Refactor `unmarshalEventsData` to use `bufio.NewScanner` with the new split function.

2.  **`card_faults.go`**:
    *   Create a `splitCardFaultRecord` function (similar to the one for events).
    *   Refactor `unmarshalFaultsData` to use the new scanner.

3.  **`card_places.go`**:
    *   Create a `splitPlaceRecord` function that returns 12-byte chunks.
    *   Refactor `unmarshalCardPlaces` to read the initial pointer, then use a scanner on the remaining byte slice.

4.  **`vu_activities.go`**:
    *   For Gen2 "Record Array" parsing functions (e.g., `parseVuActivityDailyRecordArray`):
        *   After parsing the array header (`noOfRecords`, `recordSize`), create a `bufio.Scanner` for the data portion of the array.
        *   Implement a generic fixed-size split function or specific ones as needed.
        *   Refactor the parsing loops to use the scanner.
    *   Apply a similar pattern to the Gen1 length-prefixed parsing functions.

This phased approach will allow for incremental improvements and testing at each stage. The `card_activity.go` file will be left as is, since its parsing logic is not a good fit for this pattern.

## Deeper Dive: `card_activity.go` and Custom Iterators

A deeper analysis was conducted on `card_activity.go`, which uses a cyclic buffer with backward pointers, making it unsuitable for a standard `bufio.Scanner`. The goal was to find an alternative way to decouple traversal from parsing while ensuring 100% round-trip fidelity.

### Chosen Approach: Custom Iterator

The most idiomatic and robust Go approach is to create a custom, stateful iterator struct (e.g., `CyclicRecordIterator`).

*   **Structure:** The iterator will have methods like `Next()`, `Record()`, and `Err()`.
*   **Responsibilities:**
    *   The `Next()` method will contain all the complex traversal logic: following the `prevRecordLength` pointers, handling buffer wrap-around, and detecting the end of the chain.
    *   The `Record()` method will handle the slicing of the current record's bytes, also managing buffer wrap-around. The iterator must also expose the position and length of the yielded record to enable the "buffer painting" strategy.
    *   The `Err()` method will report any errors encountered *during* traversal or slicing (e.g., corrupted pointers).
*   **Benefit:** This design perfectly separates the concern of navigating the complex buffer from the concern of parsing the content of a record. The calling code becomes a clean, idiomatic `for iterator.Next() { ... }` loop.

### Round-Trip Fidelity: The "Buffer Painting" Strategy

To ensure 100% byte-for-byte fidelity on round-trip tests for the complex cyclic buffer, a "buffer painting" strategy will be used. This avoids the pitfalls of a purely semantic comparison, which could mask parsing omissions.

1.  **Store Raw Buffer:** Upon unmarshalling, the entire `activityDailyRecords` octet string (the cyclic buffer) will be stored in memory alongside the parsed protobuf object.

2.  **Parse Semantically:** The custom iterator will traverse the buffer and parse records into their semantic protobuf fields as planned.

3.  **Marshal by Overwriting:** The marshalling process will be enhanced:
    a. Start with a copy of the original raw buffer.
    b. For each parsed semantic record, marshal it back into its binary form.
    c. "Paint" these new bytes back onto the buffer copy, overwriting the exact byte range of the original record.

4.  **Guaranteed Fidelity:** The final output is this modified buffer. This process ensures that any data in the "holes" of the buffer is preserved. Crucially, if the parser omits any fields, the re-serialized record will differ from the original, causing a `bytes.Equal` comparison against the original file to **correctly fail**. This provides a strong and byte-accurate round-trip guarantee.

### Contiguity of the Cyclic Buffer

The analysis also confirmed that the linked list of records in the cyclic buffer is **not** guaranteed to be physically contiguous. Records are linked by offsets and can be located anywhere in the buffer, reinforcing the need for a pointer-following iterator and the "buffer painting" strategy to handle interspersed unused data.
