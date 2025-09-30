# AGENTS.md

This is a Go SDK for parsing and creating tachograph data files.

## Tachograph data files

Tachograph data files are binary data dumps from tachograph vehicle units (VU), tachograph cards, and external GNSS systems. The format is specified in EU regulation for digital tachographs. Files usually have `.DDD`, `.V1B`, or `.C1B` extensions. The regulation can be vague, so careful implementation is critical.

## Goals

- Full alignment with the EU digital tachograph regulation.
- Full binary roundtrip parsing with no data loss.
- Easy-to-use and high-fidelity protobuf data model.
- Support for all types of tachograph files.

## Project Scope

To ensure a high degree of quality and alignment with the specification, this project is being developed with a phased scope.

**Phase 1 (Current Focus):**

- Full, high-fidelity support for **Driver Card** files.
- Full, high-fidelity support for **Vehicle Unit (VU)** files.

**Future Phases:**

- Support for Workshop Card, Control Card, and Company Card files is intentionally deferred. The protobuf schema has placeholders for these types, but their implementation will be addressed after the core support for Driver and VU files is complete and stable.

## Regulation

A full PDF copy of the regulation is in `docs/regulation`. Relevant chapters are OCR'd into Markdown in `docs/regulation/chapters`.

Key chapters:

- **[03-data-dictionary.md](docs/regulation/chapters/03-data-dictionary.md)**: Critical for data parsing. Contains ASN.1 specifications.
- **[05-tachograph-cards-file-structure.md](docs/regulation/chapters/05-tachograph-cards-file-structure.md)**: Essential for card data.
- **[11-response-message-content.md](docs/regulation/chapters/11-response-message-content.md)**: Essential for vehicle (VU) data.
- **[12-card-downloading.md](docs/regulation/chapters/12-card-downloading.md)**: Essential for card data.
- **[16-common-security-mechanisms.md](docs/regulation/chapters/16-common-security-mechanisms.md)**: Essential for certificates and signatures.

**IMPORTANT**: Always read [docs/asn-1.md](docs/asn-1.md) before working with ASN.1 data.

## API overview

The API is designed for simplicity and orthogonality. The top-level package contains only the main entry points, with all implementation details organized into internal packages.

### `tachograph.UnmarshalFile(data []byte) (*tachographv1.File, error)`

Main entry point for parsing a tachograph data file. It delegates to type-specific unmarshal functions in the internal packages (`internal/vu`, `internal/card`, `internal/dd`).

Usage:

```go
data, err := os.ReadFile("tachograph.DDD")
// ... handle error
file, err := tachograph.UnmarshalFile(data)
// ... handle error
fmt.Println(protojson.Format(file))
```

### `tachograph.MarshalFile(file *tachographv1.File) ([]byte, error)`

Main entry point for serializing a tachograph file. It delegates to type-specific marshal functions in the internal packages that use the `encoding.BinaryAppender` pattern.

Usage:

```go
file := &tachographv1.File{ /* ... */ }
data, err := tachograph.MarshalFile(file)
// ... handle error
os.WriteFile("tachograph.DDD", data, 0600)
```

## Project Structure and Guidelines

- **[Protobuf Schemas](./proto/AGENTS.md)**: Guidelines for developing the protobuf schemas.
- **[Development Tools](./tools/AGENTS.md)**: Guidance on build scripts and build targets.

### Package Organization

The codebase is organized into a clean, modular structure with clear separation of concerns:

#### Top-Level Package (`github.com/way-platform/tachograph-go`)

Contains only the main public API entry points:

- `UnmarshalFile(data []byte) (*tachographv1.File, error)`: Main parsing entry point
- `MarshalFile(file *tachographv1.File) ([]byte, error)`: Main serialization entry point

All implementation details are delegated to internal packages.

#### Internal Packages

**`internal/vu`**: Vehicle Unit (VU) file processing

- `UnmarshalVehicleUnitFile(data []byte) (*vuv1.File, error)`: VU-specific unmarshaling
- `MarshalVehicleUnitFile(file *vuv1.File) ([]byte, error)`: VU-specific marshaling
- Private functions for VU-specific data structures and logic

**`internal/card`**: Tachograph card file processing

- `UnmarshalDriverCardFile(data []byte) (*cardv1.DriverFile, error)`: Card-specific unmarshaling
- `MarshalDriverCardFile(file *cardv1.DriverFile) ([]byte, error)`: Card-specific marshaling
- Private functions for card-specific data structures and logic

**`internal/dd`**: Data Dictionary types and utilities

- Public `Unmarshal*` and `Append*` functions for data dictionary types
- These functions are used by both `internal/vu` and `internal/card` packages
- Contains shared parsing logic for common ASN.1 data structures

#### Visibility Rules

- **Top-level package**: Only `UnmarshalFile` and `MarshalFile` are public
- **`internal/vu`**: Only `UnmarshalVehicleUnitFile` and `MarshalVehicleUnitFile` are public
- **`internal/card`**: Only card file marshal/unmarshal functions are public
- **`internal/dd`**: Most `Unmarshal*` and `Append*` functions are public since they're shared between VU and card packages

### VU Data Modeling for Gen1 and Gen2

Vehicle Unit (VU) data structures differ between generations. Generation 1 (Gen1) often specifies a single data record or a `SET OF` records, while Generation 2 (Gen2) introduces `RecordArray` types that explicitly contain multiple records.

To create a unified and forward-compatible protobuf data model, we have adopted the following policy:

**Always use `repeated` fields for data structures that can contain multiple records in any generation.**

- For data that is a single record in Gen1 and a `RecordArray` in Gen2 (e.g., `VuCalibrationData`), we define a single `repeated` field in our protobuf message (e.g., `repeated CalibrationRecord`).
- When parsing Gen1 data, this `repeated` field will contain one element.
- When parsing Gen2 data, it will contain multiple elements.

This approach avoids the need for separate `_gen1` and `_gen2` fields, simplifying both the data model and the client-side logic required to interact with it.

### Handling Generational Differences in Marshalling

When using a unified "superset" protobuf message to represent data that differs between generations (e.g., using `FullCardNumberAndGeneration` for both Gen1 and Gen2 records), the following marshalling policy applies:

- **Populate for Consumer:** The unmarshalling process should always populate the full superset message for the consumer's benefit. For example, when parsing a Gen1 record, the `generation` field in `FullCardNumberAndGeneration` should be explicitly set to `GENERATION_1`. This makes the in-memory data model self-describing.
- **Marshal for Compliance:** The marshalling (serialization) logic **must** be version-aware. When writing a binary file, the marshaller must check the target generation and only write the fields that are compliant with that generation's specification. For example, when marshalling a `FullCardNumberAndGeneration` message for a Gen1 file, it must _only_ write the bytes for the nested `FullCardNumber` and **must not** write the `generation` field.

This ensures that we provide a rich, easy-to-use in-memory data model while maintaining perfect binary compliance and fidelity for the serialized output.

### Principle: Be Aware of Generation and Version in All Parsing and Marshalling

Tachograph specifications evolve, and data structures frequently differ between generations (Gen1 vs. Gen2) or even minor versions. It is a critical principle of this SDK to handle these differences correctly to ensure data fidelity.

**Never assume a single, fixed structure for a given data element.** Always consider the card/VU generation when parsing or marshalling data.

**Example: `EF_Places`**

A clear example of this principle is the `EF_Places` file on a driver card.

- **Gen1 `PlaceRecord`:** Has a fixed size of 12 bytes.
- **Gen2 `PlaceRecord`:** Extends the record to 22 bytes to include GNSS location data.

A parser that assumes a fixed 12-byte size will fail to correctly read Gen2 card data, leading to data corruption and loss. The correct implementation must:

1.  Determine the card generation _before_ parsing `EF_Places`.
2.  Use the generation to select the correct record size and parsing logic.
3.  The protobuf model must be a "superset" capable of holding data from both generations.

This principle applies to all data structures, not just `EF_Places`. Always verify the structure against the specification for all relevant generations.

### Marshalling and Unmarshalling

To ensure the codebase remains maintainable and easy to extend, we follow a specific structure for marshalling and unmarshalling logic organized by internal packages.

#### File Structure

The core principle is to organize files by type rather than by operation, with a direct correspondence to the protobuf schema definitions. Each internal package contains files that correspond to protobuf message types:

**`internal/vu/`**: For each VU-related protobuf file (e.g., `vu/v1/activities.proto`), there should be one corresponding file:

- `<typename>.go`: Handles both marshalling and unmarshalling for VU-specific protobuf message types

**`internal/card/`**: For each card-related protobuf file (e.g., `card/v1/activity.proto`), there should be one corresponding file:

- `<typename>.go`: Handles both marshalling and unmarshalling for card-specific protobuf message types

**`internal/dd/`**: For each data dictionary protobuf file (e.g., `dd/v1/time.proto`), there should be one corresponding file:

- `<typename>.go`: Handles both marshalling and unmarshalling for data dictionary types

This convention improves locality of context by keeping related marshalling and unmarshalling logic together, making it easier to spot inconsistencies and ensuring the operations stay in sync.

**Migration Note**: The existing codebase uses files in the top-level package with `(vu|card|dd)_` prefixes. As we refactor the codebase, we will move these into the appropriate internal packages and remove the prefixes. This migration should be done incrementally, one proto file at a time.

#### Marshalling Pattern

Marshalling is implemented using a multi-level approach to balance efficiency and simplicity:

1.  **Top-Level Function (`MarshalFile`)**: The main entry point, `tachograph.MarshalFile`, conforms to the standard `encoding.BinaryMarshaler` interface. It determines the file type and delegates to the appropriate internal package.

2.  **Package-Level Functions**: Each internal package provides its own marshal function (e.g., `vu.MarshalVehicleUnitFile`, `card.MarshalDriverCardFile`). These functions are responsible for allocating a sufficiently large `[]byte` buffer and orchestrating the serialization process for their specific file type.

3.  **Appending Functions (`Append*`)**: The detailed work of writing data is delegated to `Append*` functions (primarily in `internal/dd`) that follow the `BinaryAppender` pattern. They take a pre-allocated `[]byte` slice and append their binary representation to it, returning the updated slice. This approach avoids multiple small allocations and is more efficient.

#### Unmarshalling Pattern

Unmarshalling is implemented with a similar multi-level approach:

1.  **Top-Level Function (`UnmarshalFile`)**: The main entry point determines the file type and delegates to the appropriate internal package.

2.  **Package-Level Functions**: Each internal package provides its own unmarshal function (e.g., `vu.UnmarshalVehicleUnitFile`, `card.UnmarshalDriverCardFile`).

3.  **Unmarshalling Functions (`Unmarshal*`)**: These functions (primarily in `internal/dd`) are responsible for parsing a `[]byte` slice (or a sub-slice) into the target protobuf message. Each function takes the binary data as input and returns the populated struct and any error encountered.

### Helper Functions and Code Co-location

To maintain clarity and prevent the accumulation of disconnected utility functions, we avoid creating generic "helper" or "utility" files. Files with names like `helpers.go`, `utils.go`, or the existing `append_helpers.go`, `append_vu_helpers.go`, and `enum_helpers.go` are examples of a pattern we seek to avoid in new code.

The preferred approach is to co-locate helper functions with the code that uses them:

- **Package-specific helpers**: If a function is only used within a single file in an internal package, it should be a private function within that same file.
- **Cross-package helpers**: Functions needed across packages should generally be placed in `internal/dd` with public visibility, since this package serves as the shared foundation for data dictionary types.
- **Package-wide helpers**: If a helper function is needed across multiple files within the same internal package, it should be placed in a file with a clear, semantic name that describes its purpose (e.g., `time.go` for time-parsing helpers within `internal/dd`).

This approach maintains clear boundaries between packages while ensuring shared functionality is accessible where needed.

### Using Protobuf Reflection for Annotations

Our protobuf schemas are the single source of truth for data structures and their metadata. We use custom options to annotate fields and values with protocol-specific information, such as `protocol_enum_value` for enums or `code_page` for string encodings.

To avoid hardcoding these values in Go code and to prevent duplicative helper functions (like the ones in `enum_helpers.go`), we should use protobuf reflection to access these annotations at runtime. This approach ensures that our Go code automatically adapts to changes in the protobuf schemas.

#### Example: Accessing Custom Enum Value Options

Here is a generic example in Go for retrieving the `protocol_enum_value` annotation from an enum value. This pattern should be adapted for other annotations and types.

```go
import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// GetProtocolValueFromEnum retrieves the custom 'protocol_enum_value' annotation
// from a protoreflect.Enum value.
func GetProtocolValueFromEnum(enumValue protoreflect.Enum) (int32, bool) {
	// An enum value's Descriptor() method returns the EnumDescriptor for its type.
	// We then look up the descriptor for the specific value by its number.
	valueDesc := enumValue.Descriptor().Values().ByNumber(enumValue.Number())
	if valueDesc == nil {
		return 0, false
	}

	// Get the options for that value descriptor.
	opts := valueDesc.Options()
	if !proto.HasExtension(opts, ddv1.E_ProtocolEnumValue) {
		return 0, false
	}

	// Retrieve the value of the custom extension.
	protocolValue := proto.GetExtension(opts, ddv1.E_ProtocolEnumValue).(int32)
	return protocolValue, true
}

/*
// Example usage:
import (
    "fmt"
    ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

func main() {
    activity := ddv1.DriverActivityValue_DRIVING
    if val, ok := GetProtocolValueFromEnum(activity); ok {
        fmt.Printf("The protocol value for %s is %d\n", activity.String(), val)
        // Output: The protocol value for DRIVING is 3
    }
}
*/
```

This approach makes the code more robust and maintainable, as the logic is driven directly by the schema definitions.

### In-Code Documentation and Context

To make the marshalling and unmarshalling logic as robust and maintainable as possible, we bring critical context from the regulation specifications directly into the code. This practice reduces ambiguity and the need to cross-reference external documents.

#### ASN.1 Definitions in Comments

Every function that marshals or unmarshals a data structure defined in the ASN.1 specification should include the corresponding ASN.1 definition in its function-level comment block. This provides immediate context for the binary layout.

#### Constants for Binary Layout

Avoid using "magic numbers" for sizes and offsets. Instead, define a `const` block within the function to specify the byte layout (offsets, lengths) of the structure being processed. Use the `idx` prefix for offsets and `len` for lengths to make them easy to identify.

**Example (Fixed-Length):**

```go
// UnmarshalCardIccIdentification parses the CardIccIdentification structure.
//
// The data type `CardIccIdentification` is specified in the Data Dictionary, Section 2.23.
//
// ASN.1 Specification:
//
//     CardIccIdentification ::= SEQUENCE {
//         clockStop                   OCTET STRING (SIZE(1)),
//         cardExtendedSerialNumber    ExtendedSerialNumber,    -- 8 bytes
//         cardApprovalNumber          CardApprovalNumber,      -- 8 bytes
//         cardPersonaliserID          ManufacturerCode,        -- 1 byte
//         embedderIcAssemblerId       EmbedderIcAssemblerId,   -- 5 bytes
//         icIdentifier                OCTET STRING (SIZE(2))
//     }
func UnmarshalCardIccIdentification(data []byte) (*cardv1.Icc, error) {
    const (
        idxClockStop              = 0
        idxExtendedSerialNumber   = 1
        idxApprovalNumber         = 9
        idxPersonaliserID         = 17
        idxEmbedderID             = 18
        idxIcIdentifier           = 23
        lenCardIccIdentification  = 25
    )

    if len(data) < lenCardIccIdentification {
        return nil, fmt.Errorf("not enough data for CardIccIdentification: got %d, want %d", len(data), lenCardIccIdentification)
    }
    // ... parsing logic ...
}
```

#### Handling Variable-Length Data

Some data structures have a variable length, often specified as a range (e.g., `SIZE(8..16)`). This information should also be included in the ASN.1 comment. The parsing function must validate that the input data length falls within this expected range.

**Example (Variable-Length):**

```go
// UnmarshalVuApprovalNumber parses the VuApprovalNumber.
//
// The data type `VuApprovalNumber` is specified in the Data Dictionary, Section 2.172.
//
// ASN.1 Specification:
//
//     VuApprovalNumber ::= IA5String(SIZE(8 | 16))
func UnmarshalVuApprovalNumber(data []byte) (*ddv1.StringValue, error) {
    const (
        lenGen1 = 8
        lenGen2 = 16
    )

    if len(data) != lenGen1 && len(data) != lenGen2 {
        return nil, fmt.Errorf("invalid length for VuApprovalNumber: got %d, want %d or %d", len(data), lenGen1, lenGen2)
    }

    // ... parsing logic ...
}
```

## Testing and validation

The testing strategy uses real files to validate the parser.

### Testing Framework

All tests must use **only** the standard library `testing` package and `github.com/google/go-cmp/cmp` for comparisons. Do not use third-party testing frameworks like `testify`.

**Rationale:** This keeps dependencies minimal and ensures tests are portable and maintainable using only well-supported, stable libraries.

**Guidelines:**

- Use `t.Errorf()` for non-fatal errors and `t.Fatalf()` for fatal errors
- Use `cmp.Diff()` for comparing complex structures (slices, maps, structs)
- Use standard equality checks (`==`, `!=`) for simple types
- Check for nil explicitly before accessing pointers
- Always check errors before proceeding with test logic

**Example:**

```go
func TestParseData(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        want    *Data
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseData(tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Errorf("ParseData() expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Fatalf("ParseData() unexpected error: %v", err)
            }
            if got == nil {
                t.Fatal("ParseData() returned nil")
            }
            if diff := cmp.Diff(tt.want, got); diff != "" {
                t.Errorf("ParseData() mismatch (-want +got):\n%s", diff)
            }
        })
    }
}
```

### Golden file tests

Golden file tests for the parser are in [unmarshal_test.go](unmarshal_test.go). Example files are in the [testdata](testdata) directory. These files may contain personal data and are often in `.gitignore`.

Golden files can be updated by running tests with the `-update` flag. This should typically only be done by the user to confirm parser output changes.

### Roundtrip tests

Roundtrip tests exist but are currently in a failing state and should not be relied upon.

## Principles

Leave all version control and git to the user/developer. If you see a build error related to having a git diff, this is normal.

### Bufio Scanner Pattern for Record Parsing

Use `bufio.Scanner` with custom `SplitFunc` for all contiguous binary data parsing that advances forward through memory.

**Use for:** Fixed-size records, variable-length records, record arrays, complex structures
**Avoid for:** Backward iteration, linked lists with pointers, non-contiguous data, cyclic buffers

**Guidelines:**

- Co-locate `SplitFunc` in same file with descriptive name (e.g., `splitVuBorderCrossingRecord`)
- Never reuse `SplitFunc` across different record types
- Use `unmarshal<ProtoMessage>` naming for parsing functions
- Return errors for invalid data in `SplitFunc` (fail-fast)
- Include proper size validation

**Pattern:**

```go
func splitRecordType(data []byte, atEOF bool) (advance int, token []byte, err error) {
    const recordSize = 59
    if len(data) < recordSize {
        if atEOF { return 0, nil, nil }
        return 0, nil, nil
    }
    return recordSize, data[:recordSize], nil
}

func parseRecordArray(data []byte, offset int) ([]*Type, int, error) {
    scanner := bufio.NewScanner(bytes.NewReader(data[offset:]))
    scanner.Split(splitRecordType)
    var records []*Type
    for scanner.Scan() {
        record, err := unmarshalRecordType(scanner.Bytes())
        if err != nil { return nil, offset, err }
        records = append(records, record)
    }
    if err := scanner.Err(); err != nil {
        return nil, offset, err
    }
    return records, offset + len(data[offset:]), nil
}
```

**Benefits:** Cleaner code, better error handling, efficient memory usage, easier testing.

### Code Quality

- **No `//nolint` comments**: Never suppress linter warnings with `//nolint` comments. Instead, fix the underlying issues by removing unused code, implementing missing functionality, or restructuring the code properly.
- **Zero linter errors**: The codebase must have zero linter errors at all times. This ensures code quality and maintainability.
