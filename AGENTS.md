# AGENTS.md

This is a Go SDK for parsing and creating tachograph data files.

## Tachograph data files

Tachograph data files are binary data dumps from tachograph vehicle units (VU), tachograph cards, and external GNSS systems.

The format of the files is determined by the binary protocols used to download the data. These protocols are specified in the EU regulation for digital tachographs.

The data files are usually provided with a `.DDD`, `.V1B`, or `.C1B` extension, but this is a de-facto standard and not specified in the official regulation.

The regulation is vague in some areas, so extreme care must be exercised when reading the specification, to ensure that this SDK is fully compliant with the regulation.

## Goals

- Full alignment with the EU digital tachograph regulation
- Full binary roundtrip parsing, no data loss
- Easy to use and high-fidelity protobuf data model
- Support for all types of tachograph files

## Regulation

- A full PDF copy of the regulation is available in the `docs/regulation` directory.
- Chapters relevant to this project have been OCR'd into Markdown for easy reading.

The chapters that are relevant to this project are:

### [02-requirements.md](docs/regulation/chapters/02-requirements.md)

This chapter contains the requirements for the tachograph system.

### [03-data-dictionary.md](docs/regulation/chapters/03-data-dictionary.md)

This chapter contains the data dictionary for the tachograph system.
Read this chapter to understand the data types and their relationships.

This section contains ASN.1 specifications for data types, which is absolutely
critical information to this project.

**IMPORTANT**: Always read this chapter before working with data parsing.

### [04-tachograph-cards-specification.md](docs/regulation/chapters/04-tachograph-cards-specification.md)

The specification for tachograph cards.

### [05-tachograph-cards-file-structure.md](docs/regulation/chapters/05-tachograph-cards-file-structure.md)

The file structure for tachograph cards.

**IMPORTANT**: Always read this chapter before working with card data.

### [10-data-downloading-protocols.md](docs/regulation/chapters/10-data-downloading-protocols.md)

The data downloading protocols for the tachograph system.

### [11-response-message-content.md](docs/regulation/chapters/11-response-message-content.md)

The structure of downloaded vehicle data from tachograph files.

**IMPORTANT**: Always read this chapter before working with vehicle (VU) data.

### [12-card-downloading.md](docs/regulation/chapters/12-card-downloading.md)

The card downloading for the tachograph system.

**IMPORTANT**: Always read this chapter before working with card data.

### [15-security-requirements.md](docs/regulation/chapters/15-security-requirements.md)

The security requirements for the tachograph system.

### [16-common-security-mechanisms.md](docs/regulation/chapters/16-common-security-mechanisms.md)

The common security mechanisms for the tachograph system.

**IMPORTANT**: Always read this chapter before working with certificates and signatures.

## Protobuf schemas

Our data model is based on the protobuf schemas defined in the [`proto`](proto) directory.

We use protobuf edition 2023 for our schemas, since this enables all fields to be optional by default.

Our generated code uses the opaque API, meaning fields are hidden and accessor and setter methods must be used.

When designing our protobuf schemas, we apply the following principles:

- Align high-level conventions to the AEP (https://aep.dev) design system.
- Prefer tagged unions (e.g. a `type` enum) over `oneof` constructs, for ergonomic reasons.
- Use `google.protobuf.Timestamp` for timestamp fields.
- Avoid unsigned integers, since they are not well supported in some languages.
- Include ASN.1 definitions from [03-data-dictionary.md](docs/regulation/chapters/03-data-dictionary.md) for all messages and fields.
- Field comments should always start with a single sentence or paragraph summarizing the purpose and/or function of the field.
- **Use `StringValue` for complex strings:** Many string-like types in the data dictionary are not simple UTF-8 strings. They can be fixed-size `IA5String` types with padding, or complex `SEQUENCE` types (like `Name` and `Address`) that include a code page byte to define their character set. To handle these cases correctly and ensure lossless round-trips, any such field **must** be represented using the `datadictionary.v1.StringValue` message. This type provides the original `encoded` bytes, the `Encoding` enum to specify how to interpret them, and a convenient `decoded` field for display purposes.

### [`wayplatform.connect.tachograph.v1`](proto/wayplatform/connect/tachograph/v1)

Top-level package for all tachograph data.

Key entities are:

- `wayplatform.connect.tachograph.v1.File`, representing any type of tachograph file.

### [`wayplatform.connect.tachograph.vu.v1`](proto/wayplatform/connect/tachograph/vu/v1)

Package for vehicle unit (VU) data.

Key entities are:

- `wayplatform.connect.tachograph.vu.v1.VehicleUnitFile`, representing a vehicle unit (VU) file.

### [`wayplatform.connect.tachograph.card.v1`](proto/wayplatform/connect/tachograph/card/v1)

Package for tachograph card data.

Key entities are:

- `wayplatform.connect.tachograph.card.v1.DriverCardFile`, representing a driver card file.
- `wayplatform.connect.tachograph.card.v1.RawCardFile`, representing a generic card file (TLV records).
- One top-level message type for every EF (elementary file) that a card file can contain.
- The top-level messages should be named after the EF they represent.

### [`wayplatform.connect.tachograph.datadictionary.v1`](proto/wayplatform/connect/tachograph/datadictionary/v1)

Package for the data dictionary, with types defined in [03-data-dictionary.md](docs/regulation/chapters/03-data-dictionary.md).

This package contains shared types used across many card Elementary Files, and/or Vehicle Unit data transfers.

When a type in the data dictionary is only used in a single EF or VU transfer, it should be defined as an inline type of that specific message, instead of being included in this package.

## API overview

The API is meant to be simple, easy to use, and orthogonal. less is moere.

Any type, function or method that is not explicitly needed in the public API should be made private.

**IMPORTANT**: Always read [docs/asn-1.md](docs/asn-1.md) before working with ASN.1 data.

### `tachograph.UnmarshalFile(data []byte) (*tachographv1.File, error)`

The main entry point for parsing a tachograph data file.

Usage:

```go
data, err := os.ReadFile("tachograph.DDD")
if err != nil {
    panic(err)
}
file, err := tachograph.UnmarshalFile(data)
if err != nil {
    panic(err)
}
fmt.Println(protojson.Format(file))
```

This top-level function is implemented by delegating to more specific, private unmarshal functions for every specific type in the data model.

The typical structure of these functions is:

```go
func unmarshalTYPE(data []byte) (*TYPE, error) {
    // ...
}
```

### `tachograph.MarshalFile(file *tachographv1.File) ([]byte, error)`

Usage:

```go
file := &tachographv1.File{
    Type: tachographv1.File_DRIVER_CARD,
    DriverCard: &cardv1.DriverCardFile{
        // ...
    },
}
data, err := tachograph.MarshalFile(file)
if err != nil {
    panic(err)
}
if err := os.WriteFile("tachograph.DDD", data, 0600); err != nil {
    panic(err)
}
```

This top-level function is implemented by delegating to more specific, private functions for every specific type in the data model.

The marshal functions use the encoding.BinaryAppender pattern of appending binary data to a byte slice.

The typical structure of these functions is:

```go
func appendTYPE(dst []byte, msg *TYPE) ([]byte, error) {
    // ...
}
```

## Structure

## Developing

- The project uses a [tools](./tools) directory with a separate Go module containing tools for building, linting and generating code.
- The project uses Mage with build tasks declared in [magefile.go](./tools/magefile.go).
- Run tests with `./tools/mage test`
- Lint with `./tools/mage lint`
- Re-generate code with `./tools/mage generate`

Leave all version control and git to the user/developer. If you see a build error related to having a git diff, this is normal.

## Testing and validation

The testing strategy is built around using real files to validate the parser.

### Golden file tests

Golden file tests for the parser are defined in [unmarshal_test.go](unmarshal_test.go).

Example files are available in the [testdata](testdata) directory.

Since these files may contain personal data, they are typically included in the `.gitignore` file.

The golden files can be updated by running the tests with the `-update` flag. This should typically only be done by the user, to confirm that any changes to the parser output are accepted.

### Roundtrip tests

There are roundtrip tests currently, but they are in a failing state and should not be relied upon.

## Principles
