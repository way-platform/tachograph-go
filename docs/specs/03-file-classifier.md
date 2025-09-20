# PRD: DDD File Classifier

## 1. Overview

This document specifies the requirements for a feature that classifies tachograph `.DDD` files as either driver card or vehicle unit files based on their binary content.

### 1.1. Problem Statement

Driver card and vehicle unit `.DDD` files have fundamentally different binary structures. Before the system can parse a file, it must first reliably identify its type. Without an automated classifier, the system cannot safely select the correct parsing logic, leading to errors and an inability to process files.

### 1.2. Proposed Solution

We will implement a public function, `InferFileType`, in the top-level `tacho` package. This function will analyze the first few bytes of a file's content to identify unique structural markers that distinguish card files from unit files. The logic will be strictly based on the data storage formats defined in **Appendix 7 of the EU Tachograph regulation**.

## 2. Goals and Objectives

-   **Reliability**: The classification must be 100% accurate and verifiable against the official regulation documents.
-   **Simplicity**: Provide a clean, simple, and easy-to-use API for developers.
-   **Performance**: Ensure the classification is fast and efficient, reading only the minimum number of bytes required.
-   **Clarity**: The implementation must be well-documented with clear comments referencing the specific sections of the regulation that justify the logic.

## 3. Scope

### 3.1. In Scope

-   A new `tacho.go` file in the project root to house the new public API.
-   A `FileType` enum with values for `CardFileType`, `UnitFileType`, and `UnknownFileType`.
-   A public function `InferFileType(data []byte) FileType`.
-   Implementation logic based on the file format definitions in **Appendix 7 (DATA DOWNLOADING PROTOCOLS)**.
-   Comprehensive unit tests in `tacho_test.go` to verify the classifier's correctness.

### 3.2. Out of Scope

-   Parsing the full content of the `.DDD` files. This feature is strictly for classification.
-   Support for any file formats other than the specified driver card and vehicle unit files.

## 4. Functional Requirements

-   The `InferFileType` function must accept a `[]byte` slice as input.
-   It must return `CardFileType` if the file begins with the 2-byte tag for the `EF_ICC` file (`0x0002`).
-   It must return `UnitFileType` if the file begins with a known 1-byte `TRTP` (Transfer Request Parameter) code.
-   It must return `UnknownFileType` for all other cases, including empty files or files too short to contain a valid marker.
-   The code must contain comments that trace the logic back to the specific requirements in the regulation (e.g., DDP_019, DDP_021, DDP_023 of Appendix 7).

## 5. Technical Architecture

The feature is implemented as a single function within the `tacho` package, making it a core, easily accessible part of the library. The logic is self-contained and relies on the standard `encoding/binary` package.

The classification is based on two key markers derived from **Appendix 7 (DATA DOWNLOADING PROTOCOLS)**:

-   **Card File Marker**: A downloaded card file is a concatenation of Elementary Files (EFs). As per requirements `DDP_021` and `DDP_023`, the first EF must be `EF_ICC`, which is identified by the 2-byte tag `0x0002`.
-   **Unit File Marker**: A downloaded vehicle unit file is a concatenation of data blocks from the download protocol. As per requirement `DDP_019` and the message structures in section 2.2.6, the file begins with a 1-byte `TRTP` code identifying the first block of data.

## 6. Success Metrics

-   The `InferFileType` function is successfully implemented and exported from the `tacho` package.
-   The unit tests in `tacho_test.go` pass, demonstrating correct classification for valid card files, valid unit files, and various invalid/unknown file inputs.
-   The implementation is confirmed to be compliant with the rules defined in Appendix 7 of the regulation.
