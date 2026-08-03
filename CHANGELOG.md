# Changelog

All notable changes to this project are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

No changes yet beyond v0.1.0.

## [v0.1.0] - 2026-08-02

### Added (this milestone)

- **Cross-validation harness** (`xvalidate/`): 66 shared inputs run through
  both the Go port and the live `decimal.js`; all 1518 result lines (66 × 23
  ops) diff byte-for-byte identical (`bash xvalidate/compare.sh`).
- **decimal.js API parity**: every decimal.js long-form method name exists
  (39 alias methods, with `Plus`/`Minus`/`Times`/`ValueOf` as primary Go
  names), plus `ToJSON` and `MarshalJSON`.
- **Input parity**: `*big.Int` construction (JS `bigint`), with an input-matrix
  test covering number, string, decimal, `-0`, `NaN`, `±Infinity`, hex/bin/oct.
- **Concurrency fix**: decimal.js's module-level mutable flags
  (`external`, `inexact`, `quadrant`) were moved onto `Constructor`, removing a
  data race; `go test -race` is green (one clone per concurrent context).
- **Property tests**: round-trip, commutativity, inverse operations,
  sqrt/cbrt/exp-log inverses, comparison antisymmetry, exact sqrt, modulo
  range.
- **Stress tests**: precision 400–2048, int64/uint64 boundaries, NaN/Infinity,
  adversarial inputs, 64-goroutine clone-per-goroutine race.
- **Regression tests** locking in the three bugs fixed during porting
  (pow10 overflow, Pow negative exponent index, `1^±Inf`).
- **Encoding integrations**: `MarshalText`/`UnmarshalText`
  (`encoding.TextMarshaler`), `MarshalJSON`/`Unmarshal` (JSON), and
  `Scan`/`Value` (`database/sql`).
- **Fuzzing**: `FuzzNewString`, `FuzzArith`, `FuzzPowNew`, `FuzzTrig`
  (assert the documented `[DecimalError]` panic contract, catch everything
  else).
- **Benchmarks**: Go vs decimal.js op benchmarks with an honest comparison
  (`bench/README.md`, `bench/results.txt`).
- **Documentation**: `MATRIX.md` (operation-by-operation matrix with rounding
  semantics and edge cases), `docs/DECISIONS.md` (design-decision record),
  install/quick-start in the README,
  `docs/migrating-from-decimal-js.md`, package overview in `doc.go`.
- **Documentation site**: `decimal-go.github.io` — a dark-themed static site
  (same palette and Raleway typeface as the logo) with install, usage, parity,
  benchmark and docs sections; auto-synced from `docs/` on every push.
- **Browser playground**: `decimal-go.github.io/playground/` — the real
  library compiled to WebAssembly (`cmd/playground`, `syscall/js`), exposing
  ~45 operations in an expression editor; the WASM binary is built by the
  sync workflow in CI.
- **CI**: gofmt, go vet, staticcheck, unit + race tests, coverage summary,
  and per-target fuzzing on every push; a second workflow syncs `docs/`
  (including the freshly built playground binary) to the org Pages repo.

[unreleased]: https://github.com/iSundram/decimal-go/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/iSundram/decimal-go/releases/tag/v0.1.0