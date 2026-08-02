# Changelog

All notable changes to this project are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added (this milestone)

- **Cross-validation harness** (`xvalidate/`): 1518-line corpus of unary,
  binary, formatting and comparison cases run through both the Go port and the
  live `decimal.js`; outputs diff byte-for-byte identical
  (`bash xvalidate/compare.sh`).
- **decimal.js API parity**: all 42 long-form method aliases (`DividedBy`,
  `SquareRoot`, `NaturalLogarithm`, ...), plus `ToJSON` and `MarshalJSON`.
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
  semantics and edge cases), install/quick-start in the README,
  `docs/migrating-from-decimal-js.md`, package overview in `doc.go`.
- **CI**: gofmt, go vet, staticcheck, unit + race tests, coverage summary,
  and per-target fuzzing on every push.