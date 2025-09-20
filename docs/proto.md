# Protobuf Schema Design for Tachograph Data

## 1. Overview

This document outlines the architectural design for the Protobuf (Protocol Buffers) schemas used to represent parsed tachograph `.DDD` files. The primary goal is to establish Protobuf as the **single source of truth** for all data models within the `tacho-go` project, from which Go data structures will be generated.

This approach ensures accuracy, type safety, and long-term maintainability, moving away from manual struct definitions and the risks of model divergence.

## 2. Core Principles

The schema design is guided by the following core principles:

-   **Single Source of Truth**: The `.proto` files are the canonical definition of all data structures. All Go types will be generated from them, eliminating the need to maintain parallel models in Go and Protobuf.
-   **Preservation of Data Fidelity**: The model must be able to represent the original data with 100% accuracy, including the precise order of records as they appear in the binary files.
-   **Clarity and Explicitness**: The structure should be easy to understand and navigate. Versioning and data types should be explicit.
-   **Decoupling and Scalability**: Versioned schemas should be decoupled from one another to allow for future extensions (e.g., a new generation of tachograph) without breaking or cluttering existing models.
-   **Robustness**: The model must be resilient enough to handle known variations in the specification, such as the reuse of data structures across different generations.

## 3. Key Design Decisions

Several key architectural decisions were made to satisfy the core principles.

### 3.1. Top-Level Structure: Hierarchical Typed Payload

A `.DDD` file can be one of two types: a Card file or a Vehicle Unit (VU) file. The schema must represent this primary distinction cleanly.

-   **Decision**: A root `DDDFile` message will use a `oneof` field to contain either a `CardFile` or a `VuFile` message.
-   **Rationale**: This provides immediate, high-level classification of the file type and enforces a clear separation of concerns at the top of the data model. It is superior to a flat "bag of all possible records" approach, which would be difficult for a consumer to navigate.

### 3.2. Record Organization: Polymorphic Sequence

Card and VU files are composed of a sequence of data records. The original order of these records can be important for auditing and debugging.

-   **Decision**: Instead of grouping records by type (e.g., a list of all "activity" records, a list of all "event" records), the model will use a single `repeated Record` list. The `Record` message itself will be a `oneof` containing every possible record type.
-   **Rationale**: This "polymorphic list" approach perfectly preserves the original sequence of records as found in the binary file, ensuring complete data fidelity. While it makes accessing all records of a single type slightly less direct (requiring an iteration and type check), this is a worthwhile trade-off for maintaining the integrity of the original data sequence.

### 3.3. Versioning Strategy: Separate Packages

The tachograph specification has multiple generations (Gen1, Gen2v1, Gen2v2), each with different data structures.

-   **Decision**: Each generation's schemas will be organized into its own package and directory (e.g., `proto/vu/v1`, `proto/vu/v2`).
-   **Rationale**: This is a standard software engineering pattern for versioning. It provides strong namespacing, makes the project structure self-documenting (all Gen1 VU types are in the `v1` folder), and scales cleanly for future generations without cluttering existing packages. It is superior to using name suffixes (e.g., `VuOverview_V1`), which leads to a cluttered and less maintainable central package.

### 3.4. Code Reuse Strategy: Central `common` Package

Some primitive data types (`TimeReal`, `VehicleIdentificationNumber`, etc.) and even some entire record structures are identical across different generations.

-   **Decision**: A central `common` package (`proto/common.proto`) will be used to define all primitive and shared record types. Versioned packages will import from `common` as needed.
-   **Rationale**: This decision was made to avoid two anti-patterns:
    1.  **Duplication**: Copy-pasting definitions into each versioned package would violate the DRY (Don't Repeat Yourself) principle and create a maintenance nightmare.
    2.  **Brittle Dependency Chains**: Having `v2` import from `v1` would create a fragile dependency chain where a change in `v1` could break `v2`.
-   The `common` package establishes a clean "hub-and-spoke" dependency model (`v1 -> common`, `v2 -> common`), ensuring that shared types are defined once and that versioned schemas remain decoupled from each other.

### 3.5. Lossless Round-Tripping and Unknown Data

A critical requirement for a robust parser is the ability to perform "round-tripping": parsing a binary file into a model and then serializing that model back into an identical binary file.

-   **Decision**: The polymorphic `Record` messages (`CardRecord`, `VuRecord`) will include a final "catch-all" `UnknownRecord` type.
-   **Rationale**: This ensures that if the parser encounters a proprietary or future tag that it does not recognize, it will not discard the data. Instead, it will store the raw tag and value in an `UnknownRecord`. This allows for a perfect, lossless reconstruction of the original file, which is essential for data fidelity and enables powerful property-based testing of the entire serialization/deserialization pipeline.

## 4. Final Schema and Package Structure

The following structure is the result of the design decisions above.

### 4.1. Directory Structure

```
proto/
├── tacho.proto                   # Root DDDFile message
├── card.proto                    # CardFile and CardRecord messages
├── vu.proto                      # VuFile and VuRecord messages
├── datadictionary/
│   └── v1/
│       └── data.proto            # Shared primitive types from the Data Dictionary
├── card/
│   ├── v1/data.proto             # Gen1 card-specific records
│   └── v2/data.proto             # Gen2 card-specific records
└── vu/
    ├── v1/data.proto             # Gen1 VU-specific records
    └── v2/data.proto             # Gen2 VU-specific records
```

### 4.2. Schema Examples

#### Top-Level (`tacho.proto`)

```protobuf
syntax = "proto3";
package tacho;

import "card.proto";
import "vu.proto";
import "validate/validate.proto";

// DDDFile is the root message representing a complete .DDD file.
message DDDFile {
  option (validate.message) = {
    oneof: { required: true, fields: ["card_file", "vu_file"] }
  };

  enum FileType {
    FILE_TYPE_UNSPECIFIED = 0;
    CARD = 1;
    VEHICLE_UNIT = 2;
  }

  // The determined type of the file, which dictates which of the
  // optional fields below is populated.
  FileType type = 1;

  optional CardFile card_file = 2;
  optional VuFile vu_file = 3;
}
```

#### File-Level (`vu.proto`)

```protobuf
syntax = "proto3";
package tacho;

import "validate/validate.proto";
import "vu/v1/data.proto";
import "vu/v2/data.proto";

// VuFile contains a sequence of all records found in a Vehicle Unit file.
message VuFile {
  repeated VuRecord records = 1;
}

// VuRecord is a polymorphic container for any possible VU record type.
message VuRecord {
  option (validate.message) = {
    oneof: { required: true, fields: ["overview_g1", "activities_g1", "overview_g2"] }
  };

  enum RecordType {
    RECORD_TYPE_UNSPECIFIED = 0;
    VU_OVERVIEW_GEN1 = 1;
    VU_ACTIVITIES_GEN1 = 2;
    // ... other Gen1 record types

    VU_OVERVIEW_GEN2 = 101;
    // ... other Gen2 record types
  }

  // The specific type of this record.
  RecordType type = 1;

  // Gen 1 Records
  optional vu.v1.VuOverview overview_g1 = 2;
  optional vu.v1.VuActivities activities_g1 = 3;
  // ... other Gen1 record fields

  // Gen 2 Records
  optional vu.v2.VuOverview overview_g2 = 102;
  // ... other Gen2 record fields
}
```

#### Shared Primitives (`datadictionary/v1/data.proto`)

This package directly corresponds to the regulation's **Appendix 1: DATA DICTIONARY**.

```protobuf
syntax = "proto3";
package tacho.datadictionary.v1;

// TimeReal represents a number of seconds elapsed since 00:00:00 on 1-1-1970 UTC.
message TimeReal {
  uint32 seconds_since_epoch = 1;
}

// A vehicle's unique identification number.
message VehicleIdentificationNumber {
  string vin = 1;
}

// A shared record structure that is identical across generations.
message CardChipIdentification {
  bytes ic_serial_number = 1;
  bytes ic_manufacturing_reference = 2;
}
```

#### Versioned Records with Reuse (`card/v1/data.proto`)

This file defines messages specific to Gen1 cards and uses wrappers for shared types.

```protobuf
syntax = "proto3";
package tacho.card.v1;

import "datadictionary/v1/data.proto";

// VuOverview contains the data for a Gen1 Card ICC block (FID 0x0002).
message CardIccIdentification {
  // ... fields specific to the Gen1 version of this record
}

// Wrapper for the common CardChipIdentification type for use in Gen1.
message CardChipIdentification {
  tacho.datadictionary.v1.CardChipIdentification data = 1;
}
```

## 5. Parser Implementation Implications

This schema design directly informs the implementation of the binary parser:

1.  The parser reads a raw `.DDD` file as a byte slice.
2.  It determines if it's a Card or VU file (e.g., by checking for known initial tags). It then instantiates the correct top-level message (`CardFile` or `VuFile`).
3.  It iterates through the file, processing one binary record (TLV or TV) at a time.
4.  For each binary record, it identifies the tag and thus the corresponding Protobuf message type (e.g., tag `0x7601` maps to `vu.v1.VuOverview`).
5.  It decodes the binary data into an instance of the appropriate **generated Go struct**.
6.  It creates a new `VuRecord` (or `CardRecord`) message, sets the correct `oneof` field to the newly populated struct, and appends this `VuRecord` to the `VuFile.records` list.
7.  The final result is a single `DDDFile` message that contains a complete, ordered, and strongly-typed representation of the original file.
