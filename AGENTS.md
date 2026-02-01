# Agent Instructions

## Package Manager
Use **Mage** for all build tasks.
- `go mod tidy`
- `./tools/mage test`
- `./tools/mage lint`
- `./tools/mage generate` (Run after modifying .proto files)

## Local Skills
Reference these skills for deep procedural guides:

- **Tachograph**: For parsing patterns, DF/EF structure, regulation mapping, and implementation mandates. See `.gemini/skills/tachograph/SKILL.md`.
- **ASN.1**: For understanding BER/DER encoding and notation. See `.gemini/skills/asn1/SKILL.md`.
- **Mage**: For build script internals. See `.gemini/skills/way-magefile/SKILL.md`.
- **Protobuf**: For advanced proto operations. See `.gemini/skills/protobuf/SKILL.md`.

## Key Conventions

### Coding & Testing
- **No linter errors**: Zero tolerance.
- **Testing**: Use standard `testing` and `cmp` packages only.
- **Golden Files**:
  - Location: `testdata/`
  - Update: Run tests with `-update` flag (if supported).
  - Policy: All test data must be deterministically anonymized.

### Documentation
- **Log Files**: `docs/logs/YYYY-MM-DDTHH-MM-description.md`.

## Implementation Patterns

### Binary Parsing
- **Bufio Scanner**: Use for contiguous binary data (records, arrays).
- **Direct Slicing**: Use for fixed-size structures with known offsets. Enforce `len(data) == expected`.
- **Raw Data Painting**: For marshalling, copy original `raw_data` to a canvas, then paint semantic fields over it to preserve padding/reserved bits.

### Project Architecture
- **Package Layout**: `internal/vu` (Vehicle Unit), `internal/card` (Driver Card), `internal/dd` (Shared Data Dictionary).
- **Separation**: Public API in root, implementation in `internal/`.
