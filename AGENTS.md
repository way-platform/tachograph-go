# Agent Instructions

Go library for parsing, authenticating, anonymizing, and marshalling EU Digital
Tachograph binary files (.DDD/.V1B) into protobuf messages.

## Build

```bash
go mod tidy                # dependencies
./tools/mage test          # tests
./tools/mage lint          # golangci-lint (zero tolerance)
./tools/mage generate      # proto codegen + go generate
./tools/mage build         # full CI: download → generate → lint → test → tidy → cli → diff
```

## Skills

- **Tachograph** — `.agents/skills/tachograph/SKILL.md`
- **ASN.1** — `.agents/skills/asn1/SKILL.md`
- **Protobuf** — `.agents/skills/protobuf/SKILL.md`
- **Mage** — `.agents/skills/way-magefile/SKILL.md`

## Processing Pipeline

Five-phase pipeline with two intermediate representations:

```
Binary (.DDD)
  │  Unmarshal (binary → proto)
  ▼
RawFile (TLV/TV records with raw_data)
  │  Authenticate (signature verification at raw layer)
  │  Parse (raw records → semantic structures)
  ▼
File (semantic proto messages)
  │  Anonymize (PII removal)
  │  Unparse (semantic → raw records)
  ▼
RawFile
  │  Marshal (proto → binary)
  ▼
Binary (.DDD)
```

- **Unmarshal/Marshal** — binary ↔ `RawFile` (preserves raw bytes, no semantic
  interpretation)
- **Parse/Unparse** — `RawFile` ↔ `File` (semantic interpretation, generation
  dispatch)
- **Authenticate** — operates on `RawFile` (before parse), populates
  `Authentication` fields
- **Anonymize** — operates on `File` (after parse), replaces PII with
  deterministic test values

## Package Layout

```
tachograph.go          # public API (Unmarshal, Parse, Authenticate, Marshal, etc.)
internal/
  dd/                  # shared Data Dictionary types (TimeReal, HolderName, etc.)
  card/                # driver card TLV parsing/marshalling
  vu/                  # vehicle unit TV parsing/marshalling
  security/            # RSA/ECC certificate and signature verification
  cert/                # certificate resolution (embedded + remote)
  brainpool/           # Brainpool elliptic curves (Gen2)
  hexdump/             # hexdump ↔ binary conversion (test fixtures)
proto/                 # protobuf schemas + generated Go code
cmd/tachograph/        # CLI (cobra + charmbracelet/fang)
tools/                 # magefile, golangci-lint, buf, fetch-certs
```

## Architecture Patterns

### Options-on-Struct

All operations are methods on options structs (mirrors `protojson.MarshalOptions`):

```go
type ParseOptions struct {
    PreserveRawData bool
}
func (o ParseOptions) Parse(rawFile *tachographv1.RawFile) (*tachographv1.File, error)
```

Top-level convenience functions delegate with defaults:

```go
func Parse(rawFile *tachographv1.RawFile) (*tachographv1.File, error) {
    return ParseOptions{PreserveRawData: true}.Parse(rawFile)
}
```

### Options Embedding for Capability Inheritance

Internal options embed parent-layer options:

```go
// card.UnmarshalOptions embeds dd.UnmarshalOptions → inherits all dd unmarshal methods
type UnmarshalOptions struct {
    dd.UnmarshalOptions
    Strict bool
}
```

### Binary Parsing: Two Strategies

- **Card files (TLV)** — `bufio.Scanner` with custom split function. 5-byte
  header: 3-byte file ID + 2-byte big-endian length.
- **VU files (TV)** — direct offset slicing. 2-byte tag, then type-specific
  size calculation via `sizeOfTransferValue()`.

### Raw Data Painting (Marshalling)

When `raw_data` is present, marshal uses it as a canvas:

1. Copy `raw_data` into a fresh byte slice
2. Paint semantic fields at their known offsets
3. Unknown/reserved/padding bits preserved from original

This ensures byte-for-byte round-trip fidelity even when the binary format
contains bits not captured in the proto schema.

### Fixed-Size Layout Constants

```go
const (
    idxIssuingMemberState   = 0
    lenIssuingMemberState   = 1
    idxDriverIdentification = 1
    lenDriverIdentification = 14
)
```

Always validate: `len(data) == lenExpectedType`.

## Naming Conventions

### Functions

| Prefix       | Direction                     | Example                   |
| ------------ | ----------------------------- | ------------------------- |
| `Unmarshal*` | binary → proto (no semantics) | `UnmarshalRawCardFile`    |
| `Parse*`     | raw proto → semantic proto    | `ParseRawDriverCardFile`  |
| `Unparse*`   | semantic → raw proto          | `UnparseDriverCardFile`   |
| `Marshal*`   | proto → binary                | `MarshalDriverCardFile`   |
| `Anonymize*` | semantic → anonymized copy    | `AnonymizeDriverCardFile` |

### Constants

- `len*` — byte size (`lenHolderName = 72`)
- `idx*` — byte offset (`idxIssuingMemberState = 0`)

### Types

- `Raw*File` — intermediate TLV/TV records with `raw_data`
- `*File` — semantic structures
- `*Options` — configuration structs

## Proto Schema Conventions

See `proto/AGENTS.md` for full protobuf design guidelines. Key points:

- **Opaque API** — proto uses `API_OPAQUE` (getters/setters, not public fields)
- **AEP alignment** — follows https://aep.dev
- **Custom options** — `protocol_enum_value` maps binary protocol values to
  proto enum values for lossless round-trips
- **No unsigned integers** — use `int32`/`int64` for cross-language compat
- **`bytes` for OCTET STRING** — even single-byte values
- **ASN.1 in comments** — all DD messages include ASN.1 definition from
  regulation
- **Generation tracking** — EF messages include `dd.v1.Generation generation`
  field

## Testing

### Three Tiers

| Tier                  | Source                | Location                                              | Committed?    |
| --------------------- | --------------------- | ----------------------------------------------------- | ------------- |
| Unit (dd/)            | Inline test cases     | `internal/dd/*_test.go`                               | N/A           |
| Record fixtures       | Anonymized hexdumps   | `internal/{card,vu}/testdata/records/NNN-anonymized/` | ✅            |
| Full-file integration | Proprietary DDD files | `testdata/{card,vu}/*.DDD` (local only)               | ❌ gitignored |

### Integration Test Coverage

VU and card have different integration test depth:

| Test                                      | VU            | Card               |
| ----------------------------------------- | ------------- | ------------------ |
| Binary round-trip (`Unmarshal → Marshal`) | ✅            | ✅                 |
| Raw proto golden (per transfer/EF)        | ✅ gitignored | ❌ not implemented |
| Semantic golden (per transfer/EF)         | ✅ gitignored | ❌ not implemented |

Card full-file golden tests are a known gap. Record-level fixture tests provide
good coverage in practice.

### Gitignore Design

- `/testdata/` (root-anchored) — gitignores only the root `testdata/` directory;
  `internal/*/testdata/` is **committed** and version-controlled
- `internal/card/testdata/raw/` — gitignored (full-file JSONs, may contain PII)
- `internal/vu/testdata/*.golden.json` — gitignored (full-file goldens from
  local DDD files)
- `*.DDD`, `*.ddd`, `*.V1B`, `*.v1b` — proprietary binaries, never commit

### Record Fixture Layout

```
internal/vu/testdata/records/
  NNN-anonymized/               ← COMMITTED: one per source DDD file
    NNN-<TRANSFER_TYPE>.hexdump
    NNN-<TRANSFER_TYPE>.golden.json

internal/card/testdata/records/
  NNN-anonymized/               ← COMMITTED
    NNN-<EF>-<GEN>-<CONTENT>.hexdump
    NNN-<EF>-<GEN>-<CONTENT>.golden.json
```

### Adding New Proprietary DDD Files

1. Place DDD in `testdata/vu/` or `testdata/card/driver/` (gitignored)
2. Determine N = count of existing `NNN-anonymized/` dirs in target records dir
3. Run extract tool — creates TWO dirs per file:
   ```bash
   # VU
   go run ./internal/vu/cmd/extract-testdata-records/ -start N testdata/vu/new.DDD
   # Card (use -i flag for directory walk)
   go run ./internal/card/cmd/extract-testdata-records/ -start N -i testdata/card/driver/
   ```

   - `N-<filename>/` — original hexdumps with PII
   - `N-anonymized/` — anonymized hexdumps — commit these
4. **Delete PII dirs from disk** (not just gitignored — remove the actual files):
   ```bash
   rm -rf internal/vu/testdata/records/[0-9]*-[^a]*/
   rm -rf internal/card/testdata/records/[0-9]*-[^a]*/
   ```
5. Regenerate record-level golden JSONs (no DDD files required):
   ```bash
   cd internal/vu && go test -update . && cd ../card && go test -update .
   ```
6. Verify + commit `N-anonymized/` directories

### Golden Files

Two kinds — important to distinguish:

| Kind         | Location                                 | Committed?    | When generated                    |
| ------------ | ---------------------------------------- | ------------- | --------------------------------- |
| Record-level | `NNN-anonymized/*.golden.json`           | ✅            | Always (from hexdumps)            |
| Full-file    | `internal/vu/testdata/*.raw.golden.json` | ❌ gitignored | Only when local DDD files present |

Update record-level goldens (run from package dir):

```bash
cd internal/vu && go test -update .
cd internal/card && go test -update .
```

Full-file integration goldens are regenerated automatically when running tests
with local DDD files in `testdata/vu/` — ephemeral, never committed.

### Hexdump Format

`hexdump -C` style, handled by `internal/hexdump`. Unmarshaler is forgiving
(ignores offsets and ASCII column).

### Assertions

- `github.com/google/go-cmp/cmp` — `cmp.Diff` for all comparisons
- Standard `testing` only — no testify, no gomock
- Table-driven for parametric types; round-trip tests for binary codecs
- Round-trip tests use `UnmarshalOptions{PreserveRawData: true}` to ensure
  lossless binary fidelity (null vs space padding preserved via `raw_data`)
- `buf.build/go/protovalidate` for proto constraint validation

### Test Data Policy

- All committed fixtures must be deterministically anonymized
- Hexdump content is verified by golden JSON — never contains raw PII
- For unimplemented EFs without real data: manual byte-slice tests from spec

## Known Limitations

- `MarshalOptions.UseRawData = false` path (marshal from semantic fields only)
  is not yet implemented
- Workshop/company card types not implemented (driver card only)
- VU unmarshal cannot skip unknown transfer types (size not determinable)

## Dev Logs

`docs/logs/YYYY-MM-DDTHH-MM-description.md` — see `docs/logs/AGENTS.md` for
format.
