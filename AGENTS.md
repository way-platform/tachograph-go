# AGENTS.md

This is a Go SDK for parsing and creating tachograph data files.

## Tachograph data files

Tachograph data files are binary data dumps from tachograph vehicle units (VU), tachograph cards, and external GNSS systems. The format is specified in EU regulation for digital tachographs. Files usually have `.DDD`, `.V1B`, or `.C1B` extensions. The regulation can be vague, so careful implementation is critical.

## Goals

- Full alignment with the EU digital tachograph regulation.
- Full binary roundtrip parsing with no data loss.
- Easy-to-use and high-fidelity protobuf data model.
- Support for all types of tachograph files.

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

The API is designed for simplicity and orthogonality. Non-essential types, functions, or methods should be private.

### `tachograph.UnmarshalFile(data []byte) (*tachographv1.File, error)`

Main entry point for parsing a tachograph data file. It delegates to private, type-specific unmarshal functions.

Usage:
```go
data, err := os.ReadFile("tachograph.DDD")
// ... handle error
file, err := tachograph.UnmarshalFile(data)
// ... handle error
fmt.Println(protojson.Format(file))
```

### `tachograph.MarshalFile(file *tachographv1.File) ([]byte, error)`

Main entry point for serializing a tachograph file. It delegates to private, type-specific append functions that use the `encoding.BinaryAppender` pattern.

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

### VU Data Modeling for Gen1 and Gen2

Vehicle Unit (VU) data structures differ between generations. Generation 1 (Gen1) often specifies a single data record or a `SET OF` records, while Generation 2 (Gen2) introduces `RecordArray` types that explicitly contain multiple records.

To create a unified and forward-compatible protobuf data model, we have adopted the following policy:

**Always use `repeated` fields for data structures that can contain multiple records in any generation.**

-   For data that is a single record in Gen1 and a `RecordArray` in Gen2 (e.g., `VuCalibrationData`), we define a single `repeated` field in our protobuf message (e.g., `repeated CalibrationRecord`).
-   When parsing Gen1 data, this `repeated` field will contain one element.
-   When parsing Gen2 data, it will contain multiple elements.

This approach avoids the need for separate `_gen1` and `_gen2` fields, simplifying both the data model and the client-side logic required to interact with it.

### Marshalling and Unmarshalling

To ensure the codebase remains maintainable and easy to extend, we follow a specific structure for marshalling and unmarshalling logic within this package.

#### File Structure

The core principle is to decompose the logic into separate files for marshalling and unmarshalling, with a direct correspondence to the protobuf schema definitions. For each protobuf file that defines a major entity (e.g., `card/v1/activity.proto`), there should be two corresponding files in the top-level package:

-   `unmarshal_(vu|card)_<typename>.go`: Handles parsing binary data into the corresponding protobuf message.
-   `append_(vu|card)_<typename>.go`: Handles serializing a protobuf message into binary data. The `append_` prefix is used to signify the `BinaryAppender` pattern.

This convention keeps the context for each operation small and self-contained, making the code easier to read, debug, and maintain.

#### Marshalling Pattern

Marshalling is implemented using a two-level approach to balance efficiency and simplicity:

1.  **Top-Level Functions (`MarshalFile`)**: The main entry point, `tachograph.MarshalFile`, conforms to the standard `encoding.BinaryMarshaler` interface. It is responsible for allocating a sufficiently large `[]byte` buffer and orchestrating the overall serialization process.

2.  **Appending Functions (`append_*`)**: The detailed work of writing data is delegated to `append_*` functions that follow the `BinaryAppender` pattern. They take a pre-allocated `[]byte` slice and append their binary representation to it, returning the updated slice. This approach avoids multiple small allocations and is more efficient.

#### Unmarshalling Pattern

Unmarshalling functions (`unmarshal_*`) are responsible for parsing a `[]byte` slice (or a sub-slice) into the target protobuf message. Each function takes the binary data as input and returns the populated struct and any error encountered.

### Helper Functions and Code Co-location

To maintain clarity and prevent the accumulation of disconnected utility functions, we avoid creating generic "helper" or "utility" files. Files with names like `helpers.go`, `utils.go`, or the existing `append_helpers.go`, `append_vu_helpers.go`, and `enum_helpers.go` are examples of a pattern we seek to avoid in new code.

The preferred approach is to co-locate helper functions with the code that uses them. If a function is only used within a single `append_*.go` or `unmarshal_*.go` file, it should be a private function within that same file.

If a helper function is needed across multiple, related files (e.g., for handling a specific data type that appears in different structures), it should be placed in a file with a clear, semantic name that describes its purpose (e.g., `unmarshal_time.go` for time-parsing helpers). This makes the codebase easier to navigate and understand.

### Using Protobuf Reflection for Annotations

Our protobuf schemas are the single source of truth for data structures and their metadata. We use custom options to annotate fields and values with protocol-specific information, such as `protocol_enum_value` for enums or `code_page` for string encodings.

To avoid hardcoding these values in Go code and to prevent duplicative helper functions (like the ones in `enum_helpers.go`), we should use protobuf reflection to access these annotations at runtime. This approach ensures that our Go code automatically adapts to changes in the protobuf schemas.

#### Example: Accessing Custom Enum Value Options

Here is a generic example in Go for retrieving the `protocol_enum_value` annotation from an enum value. This pattern should be adapted for other annotations and types.

```go
import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
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
	if !proto.HasExtension(opts, datadictionaryv1.E_ProtocolEnumValue) {
		return 0, false
	}

	// Retrieve the value of the custom extension.
	protocolValue := proto.GetExtension(opts, datadictionaryv1.E_ProtocolEnumValue).(int32)
	return protocolValue, true
}

/*
// Example usage:
import (
    "fmt"
    datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

func main() {
    activity := datadictionaryv1.DriverActivityValue_DRIVING
    if val, ok := GetProtocolValueFromEnum(activity); ok {
        fmt.Printf("The protocol value for %s is %d
", activity.String(), val)
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
// unmarshalCardIccIdentification parses the CardIccIdentification structure.
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
func unmarshalCardIccIdentification(data []byte) (*cardv1.Icc, error) {
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
// unmarshalVuApprovalNumber parses the VuApprovalNumber.
//
// The data type `VuApprovalNumber` is specified in the Data Dictionary, Section 2.172.
//
// ASN.1 Specification:
//
//     VuApprovalNumber ::= IA5String(SIZE(8 | 16))
func unmarshalVuApprovalNumber(data []byte) (*datadictionaryv1.StringValue, error) {
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

### Golden file tests

Golden file tests for the parser are in [unmarshal_test.go](unmarshal_test.go). Example files are in the [testdata](testdata) directory. These files may contain personal data and are often in `.gitignore`.

Golden files can be updated by running tests with the `-update` flag. This should typically only be done by the user to confirm parser output changes.

### Roundtrip tests

Roundtrip tests exist but are currently in a failing state and should not be relied upon.

## Principles

Leave all version control and git to the user/developer. If you see a build error related to having a git diff, this is normal.