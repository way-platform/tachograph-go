# Crono patches to `tachograph-go`

This branch (`crono-patches`) tracks the patches that the `ddd-parser`
project applies on top of upstream [`way-platform/tachograph-go`](https://github.com/way-platform/tachograph-go)
to parse real-world `.ddd` files produced by firmware variants that
upstream does not yet handle.

- **Upstream base:** `v0.17.3`
- **Fork location:** [github.com/Trackysat-SRL/tachograph-go](https://github.com/Trackysat-SRL/tachograph-go)
- **Fork tag consumed by ddd-parser:** `v0.17.3-crono.1`
- **Consumer:** `ddd-parser` via `replace` directive in its `go.mod`:
  ```
  replace github.com/way-platform/tachograph-go => github.com/Trackysat-SRL/tachograph-go v0.17.3-crono.1
  ```

## Patches

| # | Commit | File | Summary | Upstream status |
|---|--------|------|---------|-----------------|
| 1 | `34ceff2` | `internal/dd/full_card_number.go` | Return partial `FullCardNumber` for card types without a structured parser (CONTROL_CARD type 5, MANUFACTURING_CARD) instead of aborting the whole file. | PR upstream planned |
| 2 | `9deab28` | `internal/dd/vehicle_registration_identification.go` | Accept 14-byte compact VRI records (no codepage byte) alongside the canonical 15-byte layout. | PR upstream planned |
| 3 | `f16a7ed` | `internal/vu/technical_data_gen2_v1.go` | Skip `CalibrationRecord` arrays whose record size is not 168 bytes (observed 28-byte variants). Workaround, not a structural fix. | Crono-only — no upstream PR |

### Patch 1 — FullCardNumber: unknown card types

**Problem:** upstream fails the whole file with `unsupported card type: N`
when it encounters a card type that the enum recognises but for which
there is no structured parser branch. In practice this hit us with
`CONTROL_CARD` (type 5) in Gen1 VU files inside control-activity records.

**Fix:** add a `default` case to the switch that returns the
`FullCardNumber` already populated with `CardType` and `raw_data`,
leaving downstream callers free to decide what to do with the variant.

**Real-world trigger:** Gen1 `M_*.ddd` files containing a control-card
record.

### Patch 2 — VehicleRegistrationIdentification: 14-byte compact layout

**Problem:** the canonical VRI record per EU 165/2014 Annex I(C) is
15 bytes (1 nation + 1 codepage + 13 number). Some Gen2 V1 firmwares
emit 14-byte records instead (1 nation + 13 number, no explicit
codepage), and the strict length check rejects them outright.

**Fix:** accept both lengths. For 14-byte records synthesise a
`StringValue` by prepending codepage `0` (default) before delegating
to `UnmarshalStringValue`, so the rest of the parser treats both
variants identically.

**Real-world trigger:** `M_*.ddd` Gen2 V1 files from certain vehicle
registrations (e.g. targa `GC866DP`) produced by firmware that predates
the explicit codepage field.

### Patch 3 — CalibrationRecord: non-standard record size

**Problem:** upstream hard-requires `recordSize == 168` for Gen2 V1
`CalibrationRecord` arrays. In the wild we have observed 28-byte
records, likely a GNSS-coupled layout or a national deviation, which
causes `parseCalibrationRecordArrayGen2V1` to return an error and fail
the whole file.

**Fix (workaround):** when `recordSize != 168`, advance the byte
cursor by `headerSize + recordSize * noOfRecords` and return an empty
slice. The array framing stays honoured, subsequent sections parse
normally, and the downstream mapper does not currently extract
calibration data so no user-visible field is lost.

**Why this is crono-only:** the workaround drops data. A proper fix
requires decoding the 28-byte layout and emitting best-effort records,
which needs sample files and spec clarification. Kept locally until
then.

**Real-world trigger:** `M_*.ddd` Gen2 V1 files with GNSS-enabled
firmware.

## Upstream contribution plan

- Patches 1 and 2 will be submitted to
  `way-platform/tachograph-go` as independent PRs, each with the
  relevant test coverage added under `internal/dd/` (see
  existing `unmarshal_test.go` / `golden_test.go` for patterns).
- Patch 3 stays on the fork until either (a) we decode the alternate
  layout, or (b) upstream adds explicit support.

When a patch gets merged upstream, drop the corresponding commit from
this branch during the next rebase (see runbook below) and cut a new
`vX.Y.Z-crono.<n>` tag.

## Update runbook (rebase on new upstream tag)

```bash
# Inside third_party/tachograph-go
git fetch origin --tags              # upstream way-platform
git fetch fork                       # our public fork
git checkout crono-patches
git rebase vX.Y.Z                    # replay the 3 commits on the new tag
# resolve conflicts, if any
go test ./...                        # run library tests
git push --force-with-lease fork crono-patches
git tag vX.Y.Z-crono.1
git push fork vX.Y.Z-crono.1

# Inside ddd-parser
cd ../..
go get github.com/way-platform/tachograph-go@vX.Y.Z
# Update go.mod:
#   require github.com/way-platform/tachograph-go vX.Y.Z
#   replace github.com/way-platform/tachograph-go => github.com/Trackysat-SRL/tachograph-go vX.Y.Z-crono.1
go build ./...
./ddd-parser.exe -input ./input -output ./output   # smoke test on real files
```

## Contact

Maintainer: `ddd-parser` / Crono team.
