# PRD: Data Dictionary Parser & Code Generation

## 1. Overview

This document extends the regulation parsing toolchain to include the Data Dictionary defined in Appendix 1 of the EU Tachograph regulation.

### 1.1. Problem Statement

While the existing toolchain handles identifiers, it does not cover the complex data structures that comprise the actual content of tachograph files (DDD files). Manually defining Go `structs` for the hundreds of data types in the Data Dictionary is error-prone, time-consuming, and difficult to maintain against regulatory updates.

### 1.2. Proposed Solution

We will create a new, dedicated two-stage pipeline to process the Data Dictionary:

1.  **Parser (`tachomodel-parser`)**: A new tool that reads the cleaned HTML of Appendix 1, parses the data type definitions, and outputs a structured Intermediate Representation (IR) as `tachomodel.json`.
2.  **Generator (`tachomodel-generator`)**: A new tool that consumes `tachomodel.json` and generates idiomatic Go `struct` types for each data definition.

This will establish Appendix 1 as the single source of truth for all data models used in the project.

## 2. Goals and Objectives

-   **Accuracy**: Ensure Go `structs` perfectly match the data structures specified in the regulation.
-   **Maintainability**: Automate the update process for all data models when the regulation changes.
-   **Completeness**: Cover all data types defined in the Data Dictionary to ensure full DDD file parsing capabilities.
-   **Code Quality**: Generate clean, idiomatic, and well-documented Go code.

## 3. Scope

### 3.1. In Scope

-   **`tachomodel-parser` Tool**:
    -   Parses `appendix-1-data-dictionary.html`.
    -   Outputs a `tachomodel.json` file to `/internal/gen`.
    -   Extracts each data type's name, description, and detailed structure (field name, ASN.1 type, description, size in bytes) from the definition tables.
-   **`tachomodel-generator` Tool**:
    -   Consumes `tachomodel.json`.
    -   Generates a `tacho/datadictionary.go` file containing Go `struct` definitions for all parsed data types.
-   The parser must handle nested data structures and correctly represent the relationships in the JSON IR.
-   The generator must map ASN.1 types to appropriate Go types (e.g., `OCTET STRING`, `BCDString`, `TimeReal`) and include field descriptions as code comments.

### 3.2. Out of Scope

-   Modifying the existing `tachocard` or `tachounit` parsers and generators.
-   Parsing any part of Appendix 1 other than the data type definition sections (Chapter 2 onwards).

## 4. Functional Requirements

### 4.1. `tachomodel-parser`

-   **Input**: Path to the cleaned `appendix-1-data-dictionary.html` file.
-   **Output**: A `tachomodel.json` file.
-   **JSON IR Structure**: The JSON should be an array of objects, where each object represents a data type and contains:
    -   `name`: The name of the data type (e.g., "CardActivityDailyRecord").
    -   `description`: The textual description of the data type.
    -   `fields`: An array of field objects, each containing:
        -   `name`: The field name.
        -   `type`: The ASN.1 data type.
        -   `description`: The field's description.
        -   `size`: The size in bytes.

### 4.2. `tachomodel-generator`

-   **Input**: Path to the `tachomodel.json` IR file.
-   **Output**: A single Go source file at `tacho/datadictionary.go`.
-   **Logic**:
    -   Reads and unmarshals the JSON IR.
    -   Iterates through each data type definition.
    -   Generates a Go `struct` for each, mapping field names and types.
    -   Adds the description of the data type and each field as Go comments.
    -   Uses the `@internal/codegen` utility for consistent file generation.

## 5. Technical Architecture

A new pipeline will be created alongside the existing ones:

**Tachomodel Pipeline:**
`[Cleaned Appendix 1 HTML]` --> **`tachomodel-parser`** --> `[tachomodel.json]` --> **`tachomodel-generator`** --> `[tacho/datadictionary.go]`

## 6. Success Metrics

-   The `tachomodel-parser` and `tachomodel-generator` tools are created and function as specified.
-   The generated `tacho/models/models.go` file contains accurate Go `structs` for all data types in Appendix 1.
-   The generated code is well-formatted, commented, and passes linting checks.
-   The process is repeatable and can be integrated into the project's build scripts.
