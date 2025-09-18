# PRD: Regulation Parsing & Code Generation Toolchain

## 1. Overview

This document outlines the requirements for a toolchain designed to automate the creation of Go constants and types from the official EU Tachograph regulation HTML documents.

### 1.1. Problem Statement

The tachograph specification is defined across numerous appendices containing hundreds of specific identifiers, tags, and codes. Manually transcribing these values into Go code is:
-   **Error-Prone**: High risk of typos in names or values.
-   **Time-Consuming**: Requires significant developer time for initial creation and subsequent updates.
-   **Hard to Maintain**: When regulations are updated, changes must be found and applied manually, which is inefficient and unreliable.

### 1.2. Proposed Solution

We will build a dedicated, two-stage toolchain for each subdomain (`tachocard` and `tachounit`):

1.  **Parsers (`tachocard-tag-parser`, `tachounit-tag-parser`)**: A dedicated parser for each subdomain that reads the relevant cleaned HTML regulation documents, extracts key data (like tag identifiers), and outputs a clean, structured Intermediate Representation (IR) in a domain-specific JSON file (e.g., `tachocard-tags.json`).
2.  **Generators (`tachocard-tag-generator`, `tachounit-tag-generator`)**: A dedicated generator for each subdomain that reads the corresponding IR JSON file and generates well-formatted, idiomatic Go source code.

This establishes the regulation documents as the single source of truth while maintaining a strict separation of concerns between the subdomains.

## 2. Goals and Objectives

-   **Accuracy**: Eliminate human error by automating the extraction of values from the source documents.
-   **Maintainability**: Drastically reduce the effort required to update constants when the regulation changes. The process should be as simple as replacing the source documents and re-running the tools.
-   **Single Source of Truth**: Ensure that the constants in the Go codebase are verifiably derived from the regulation documents.
-   **Domain Separation**: Keep the tooling and generated artifacts for `tachocard` and `tachounit` completely separate.

## 3. Scope

### 3.1. In Scope

-   **`tachocard-tag-parser` Tool**:
    -   Parses `appendix-2-tachograph-cards-specification.html`.
    -   Outputs a `tachocard-tags.json` file to `/internal/gen`.
    -   Extracts Tachograph Card File IDs.
-   **`tachocard-tag-generator` Tool**:
    -   Consumes `tachocard-tags.json`.
    -   Generates `tacho/tachocard/tag.go` using the `@internal/codegen` utility.
-   **`tachounit-tag-parser` Tool**:
    -   Parses `appendix-7-data-downloading-protocols.html` and `appendix-8-calibration-protocol.html`.
    -   Outputs a `tachounit-tags.json` file to `/internal/gen`.
    -   Extracts Vehicle Unit Download Types (TRTPs) and Data Identifiers (RDIs).
-   **`tachounit-tag-generator` Tool**:
    -   Consumes `tachounit-tags.json`.
    -   Generates `tacho/tachounit/tag.go` using the `@internal/codegen` utility.
-   All parsers must convert hexadecimal strings from the source to numeric values in the JSON output.

### 3.2. Out of Scope

-   A single, monolithic parser or generator tool.
-   Modifying the existing `@tools/cmd/regulation-parser` tool, which acts as a pre-processor.
-   Parsing any data beyond the specific tag/identifier tables required for the initial scope.

## 4. Functional Requirements

### 4.1. Parsers

-   Each parser will be a standalone Go command-line application.
-   **Input**: Path to the directory containing cleaned regulation HTML files.
-   **Output**: A domain-specific JSON file (e.g., `tachocard-tags.json`).
-   **Logic**: Must locate and parse specific HTML tables, extract relevant data, and convert values to numbers.

### 4.2. Generators

-   Each generator will be a standalone Go command-line application.
-   **Input**: Path to a domain-specific JSON IR file.
-   **Output**: A single Go source file in the appropriate target package.
-   **Logic**: Must read the JSON, unmarshal it, and use the `@internal/codegen` utility to produce idiomatic Go `const` blocks.

## 5. Technical Architecture

The toolchain will consist of two parallel, two-stage pipelines:

**Tachocard Pipeline:**
`[Cleaned Appendix 2 HTML]` --> **`tachocard-tag-parser`** --> `[tachocard-tags.json]` --> **`tachocard-tag-generator`** --> `[tacho/tachocard/tag.go]`

**Tachounit Pipeline:**
`[Cleaned App. 7 & 8 HTML]` --> **`tachounit-tag-parser`** --> `[tachounit-tags.json]` --> **`tachounit-tag-generator`** --> `[tacho/tachounit/tag.go]`

## 6. Success Metrics

-   The four tools are created and function as specified.
-   The generated `tag.go` files in `/tacho/tachocard` and `/tacho/tachounit` contain the correct constants, numeric values, and comments, as verified against the source documents.
-   The entire process is repeatable and can be orchestrated via `go generate` or a build script.
