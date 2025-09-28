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

## Testing and validation

The testing strategy uses real files to validate the parser.

### Golden file tests

Golden file tests for the parser are in [unmarshal_test.go](unmarshal_test.go). Example files are in the [testdata](testdata) directory. These files may contain personal data and are often in `.gitignore`.

Golden files can be updated by running tests with the `-update` flag. This should typically only be done by the user to confirm parser output changes.

### Roundtrip tests

Roundtrip tests exist but are currently in a failing state and should not be relied upon.

## Principles

Leave all version control and git to the user/developer. If you see a build error related to having a git diff, this is normal.