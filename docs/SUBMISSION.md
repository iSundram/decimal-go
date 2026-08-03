# Port Mortem 2026 — submission: decimal-go

## The project

**decimal-go** is a Go port of the [`decimal.js`](https://github.com/MikeMcl/decimal.js)
library (v10.6.0, by Michael Mclaughlin) — an arbitrary-precision decimal
arithmetic library for Go, with **zero dependencies** and full behavioral
parity with the original, verified byte-for-byte.

Repository: `github.com/iSundram/decimal-go` — install with:

```sh
go get github.com/iSundram/decimal-go
```

```go
price := decimal.New("19.99")
total := price.Times(decimal.New("3"))
fmt.Println(total) // 59.97 — no float64 surprises
```

## Links

- **Repository**: <https://github.com/iSundram/decimal-go>
- **Documentation site**: <https://decimal-go.github.io/>
- **Browser playground** (the library compiled to WebAssembly — try `sqrt(2)`,
  `plus(0.1, 0.2)` or `toFixed(pi, 30)` with no install):
  <https://decimal-go.github.io/playground/>
- **Package docs**: <https://pkg.go.dev/github.com/iSundram/decimal-go>

## Why this port

`float64` cannot represent money, rates, or high-precision measurements
exactly (`0.1 + 0.2 != 0.3`). decimal.js is one of the most widely used
exact-decimal libraries in the JavaScript world; Go had no drop-in
behavioral equivalent (shopspring/decimal et al. differ in construction,
rounding, and formatting semantics). This port brings the *reference
behavior* — not just the idea — to Go.

## Architecture in one breath

- A `Decimal` is `sign * coefficient * 10^exponent`; the coefficient is a
  `[]int32` of base-1e7 words (the same base decimal.js uses), the exponent is
  an `int64` (`decimal.go`).
- Arithmetic, the long-division kernel, comparisons, and rounding live in
  `ops.go`; transcendentals (exp/ln/log/sqrt/cbrt/trig) in `trans.go` and
  `trig.go`; string↔value in `convert.go`.
- Settings (precision, rounding, exponent bounds) live on a `Constructor`;
  `Default` mirrors the decimal.js defaults, `Clone` creates an isolated one —
  the exact analogue of `Decimal.clone()` (which decimal.js documents: one
  clone per concurrent context).
- The package has **no imports beyond the standard library**.

## The four hardest parts

### 1. Proving API parity, not just claiming it

Every decimal.js long-form method name (`DividedBy`, `NaturalLogarithm`,
`SquareRoot`, `ToPower`, … — 39 as alias methods here, with `Plus`, `Minus`,
`Times`, `ValueOf` as the primary Go names) exists, and `parity.go`'s tests
check the aliases are *equal to* the primary Go methods. The Go API is
idiomatic (`x.Div(y)`), while copied JS code keeps compiling via the aliases.
One deliberate, documented divergence: decimal.js v10.6.0 ships *commented
out* instance `min`/`max` (docs-only); the port matches the real runtime API
(constructor statics only).

### 2. Cross-validation against the live library

`xvalidate/` runs a shared corpus through both the Go port (`go run`) and the
real decimal.js (`node`): **1518 result lines, byte-for-byte identical**
(`bash xvalidate/compare.sh`). On top of that, all 61 decimal.js test modules
were ported as white-box tests (every assertion kept, same order), which is
where the three real porting bugs surfaced:

- `pow10` int64 overflow/hang,
- `Pow` with a negative base hitting an out-of-range index,
- `1^±Infinity` incorrectly returning 1.

All three are locked in by `regression_test.go`.

### 3. Constructor concurrency

decimal.js keeps per-operation flags (`external`, `inexact`, `quadrant`) in
module-level globals — "safe" only because JS is single-threaded. Under
`go test -race` the same design is a real data race the moment methods are
called concurrently. The fix moves those flags onto the `Constructor` (the
per-clone state decimal.js's own docs already prescribe), and
`stress_test.go` runs 64 goroutines, each with its own clone, race-clean.

### 4. Faithfulness vs. "fixing" the original

Two of decimal.js's quirks were tempting to smooth over, and both were kept
on purpose (documented in `docs/DECISIONS.md`):

- Transcendentals that would need more than 1024 significant digits of π/ln10
  panic with `[DecimalError] Precision limit exceeded`, including the quirk
  that the *trig* path may leave the constructor mutated until you reset it
  (the ln10 path restores config before panicking).
- The constructor does *not* round at construction time; rounding happens at
  operation time, so `New("123.4500").String()` is `123.45`.

## Performance — honestly

Ratios are only comparable when the same work is measured on both sides
(full documentation in `bench/README.md`, raw data in `bench/results.txt`).

**Op-only benchmarks** (pre-created operands, same precision, same machine,
Node v22 vs Go 1.25): **23 of 25 directly comparable pairs are faster in Go**
(Mul 0.32×, Div 0.48×, Div@1000 0.39×, Cbrt 0.43×, Atan 0.45×, … — ratio < 1
means Go is faster).

The two slower paths, disclosed rather than hidden:

- `New(string)` parsing is ~1.45× slower — V8's JIT is extremely good at this
  string-parsing workload while Go pays per-parse allocation overhead.
- Full parse → op → format round trips therefore land at ~1.7× slower.

Go's advantage *widens* as precision grows (Div@1000: 0.39×).

## Test & validation inventory

| Layer | Where | What it proves |
|---|---|---|
| Ported test suite | `*_test.go` (61 modules) | Every decimal.js assertion, same order, green |
| Cross-validation | `xvalidate/` | Byte-for-byte identical output vs. live decimal.js, 1518 result lines (66 inputs × 23 ops) |
| Property tests | `property_test.go` | Round-trip, commutativity, add/sub & mul/div inverses, sqrt/cbrt/exp-ln inverses, cmp antisymmetry, modulo range |
| Stress + race | `stress_test.go` | Precision 400–2048, int64/uint64 edge values, NaN/∞, 64-way race-clean |
| Regression tests | `regression_test.go` | The 3 real porting bugs locked in |
| Input matrix | `input_test.go` | Every decimal.js input type incl. `*big.Int`, hex/bin/oct, ±0, NaN, ∞ |
| Fuzzing | `fuzz_test.go` | 4 targets; only `[DecimalError]` panics are permitted |
| Encoding | `encoding.go` | JSON, text, and `database/sql` integration |
| CI | `.github/workflows/ci.yml` | gofmt, vet, staticcheck, unit + race tests, coverage, per-target fuzz |

## Live playground (WASM)

The library is compiled to WebAssembly and runs in the browser at
`decimal-go.github.io/playground/` — no server, no install. `cmd/playground`
exposes ~45 operations through `syscall/js` (build-tagged `js && wasm`, so it
is excluded from host builds), and a small expression editor lets visitors
type things like `plus(0.1, 0.2)`, `sqrt(2)`, `div(1, 3)` or
`toFixed(pi, 30)` and see the real Go implementation's answers. The WASM
binary is rebuilt by the docs-sync workflow on every push to `main`, keeping
the deployed binary in lockstep with the library.

## Status

Coverage: **96%** statements. Race detector: clean. staticcheck: clean.
Released as **v0.1.0** (installable via `go get`). The only known limitation,
shared with the original: astronomical operands (e.g. `0b1p2738415519256`,
≈10^824 billion digits) make `Mod`/`divide` run for non-finite time — a
property of arbitrary-precision decimals, documented in `docs/DECISIONS.md` §7.
