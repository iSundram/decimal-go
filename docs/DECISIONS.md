# Design decisions

This document records the decisions made while porting `decimal.js` v10.6.0 to
Go, and the reasoning behind them. It exists so that future contributors
understand *why* the code looks the way it does before changing it.

## 1. Preserve behavior over redesign

The goal is a **behavioral port**: where a Go idiom and the original behavior
conflict, the original behavior wins. Examples:

- `Cmp` returns a `float64` (`-1`/`0`/`1`, `NaN` when an operand is `NaN`)
  instead of the Go-conventional `int` — matching `decimal.js`'s `cmp`.
- `String()` omits the leading minus of `-0` while `ValueOf()` keeps it, just
  like `toString()` vs `valueOf()`.
- Transcendental functions panic with `[DecimalError] Precision limit exceeded`
  when they would need more than 1024 significant digits of π or ln10 — a
  faithful port of decimal.js's `precisionLimitExceeded`, including the
  original's quirk of leaving the constructor config mutated until reset.

Rationale: the value of the port is provable fidelity (see `xvalidate/`, which
diffs byte-for-byte against the real decimal.js). Introducing "improvements"
silently breaks parity and makes the cross-validation meaningless.

## 2. Representation: base-1e7 coefficient words + int64 exponent

`Decimal` is `sign * coefficient * 10^exponent`, where the coefficient is a
`[]int32` slice of base-1e7 words, exactly mirroring decimal.js's
`base = 1e7` digit arrays, with a nil digits slice for `NaN` (sign == 0) or
`±Infinity`.

Alternatives considered:

| Alternative | Why rejected |
|---|---|
| `math/big.Int` coefficient | Would change arithmetic cost, rounding behavior and exactness guarantees, breaking byte-for-byte parity (decimal.js's algorithms are written against base-1e7 arrays, e.g. long division with trial digits). |
| `math/big.Float` | Wrong behavior class: binary exponent, probabilistic rounding — not decimal semantics. |
| base-10^n with n > 7 | decimal.js normalizes on 1e7 (fits JS float64 mult. precision and Go `int32` ops); any change ripples into every division/log constant cache in the ported test corpus. |

The `int64` exponent replaces JS's IEEE-754 double exponent field (decimal.js
stores exponents in a JS number, so ±9e15 bounds are natural there) — Go's
`int64` is both safe and equals the JS MaxE/MinE headroom (±9e15).

## 3. Naming: idiomatic Go first, parity aliases second

Methods use short Go names (`Div`, `Sqrt`, `Ln`) as the primary API. All 42
decimal.js long-form names (`DividedBy`, `SquareRoot`, `NaturalLogarithm`,
…) exist as thin aliases in `parity.go`, verified by tests to be *equal* to
the primary methods.

Rationale: Go users get `x.Div(y)`; migrated JS code keeps compiling almost
verbatim. The cost is one extra public name per op and is documented in
`docs/migrating-from-decimal-js.md`.

The single deliberate divergence: decimal.js v10.6.0 has instance `max`/`min`
methods commented out in source (documented but not implemented). The port
matches the *real* runtime API: constructor statics only.

## 4. Errors: panic with `[DecimalError]` prefix

decimal.js throws on invalid input; Go would conventionally return `(T,
error)`. We chose panic with a `[DecimalError]`-prefixed message for these
reasons:

- **Parity**: every JS throw point maps to a Go panic point; tests compare
  message strings directly (`assertException`).
- **API ergonomics**: arithmetic ops returning `(Decimal, error)` would
  contradict the chaining style the library is designed for
  (`a.Times(b).Plus(c)`).
- **Precedent**: `math/big` parsers panic on invalid input in some paths too.

Garbage input never produces undefined behavior — only the documented panic —
which is exactly what the fuzz targets assert.

## 5. Constructor/config model: one clone per concurrent context

decimal.js holds counters on module-level globals (`external`, `inexact`,
`quadrant`); users MUST follow the doc advice that "each concurrent context use
its own clone." In Go that advice is a hard requirement: cloned constructors
each get their own copy of those flags (moved onto `Constructor`), which
`go test -race` verified (64 concurrent clone-instances).

Rationale: this keeps parity with the *meaning* of decimal.js's
config-and-clone design while removing the one real race hazard. A global
`sync.Mutex` around all arithmetic was rejected (kills parallelism and is
also the wrong abstraction for a formatting/orchestration state).

## 6. No rounding at construction time

`New("123.4500")` retains trailing zeros / full precision in the internal
value, matching decimal.js exactly. Rounding happens only in *operations*,
via `finalise(pr, rm)`. During work the code runs with `external=false` to
suppress any rounding of intermediates (e.g. `Mod` computes
`x - floor(x/y)*y` without intermediate rounding).

This is a subtle behavioral difference vs. packages like shopspring/decimal,
and users notice it first in `String()` output, so the migration guide calls
it out explicitly.

## 7. Fuzzing: contract assertions, not assertions on unbounded work

The four fuzz targets (`FuzzNewString`, `FuzzArith`, `FuzzPowNew`,
`FuzzTrig`) assert:

1. Anything that panics must be a `[DecimalError]` (the documented contract).
2. Parsing round-trips deterministically.

Inputs like a binary literal with a huge binary exponent (e.g.
`0b1p2738415519256`, which is 2^2.7e12 ≈ 10^824 billion) create numbers whose
decimal expansion is astronomically long; both the original decimal.js and
this port cannot compute e.g. `Mod` on them in finite time and finite memory
(this is a **property of arbitrary-precision decimal arithmetic**, not a port
bug: the real decimal.js hangs identically on the same input). The *library*
deliberately has no artificial cap — that would diverge from the reference.
Instead, the fuzz target itself skips pairs whose quotient would exceed
`maxQuotientWords` (`1e5` base-1e7 words ≈ 700k digits, ~3 MB), so the
worker can never OOM and coverage-loss is nil (even decimal.js cannot
compute those cases).

## 8. Benchmark honesty

`bench/README.md` documents the boundaries of every ratio:
op-only figures isolate the arithmetic kernel (Go wins ~2–6×);
`New(string)` is ~1.45× slower (Go allocation overhead vs. JIT-optimised JS
string parsing); round trips that re-parse inputs are also slower. We do not
publish "decimal-go is faster than decimal.js" without the qualification —
comparison notes live in `bench/README.md` § "Honest comparison notes".

## 9. Dependencies: none

The module imports only the standard library. This is a direct consequence of
the base-1e7 representation (§2): all arithmetic is hand-rolled on `int32`
slices, and decimal.js itself is dependency-free.

## 9a. Browser playground

`cmd/playground` compiles the library to WebAssembly for the docs site. It is
build-tagged `js && wasm` so `go test ./...` and ordinary builds never see it.
The playground is a demo surface, not an API: it exposes a small, string-based
subset of operations via `syscall/js` and is deliberately not a supported
interface. The binary is rebuilt in CI by the docs-sync workflow, so the
deployed playground always matches the pushed library.

## 10. Versioning

`v0.1.0` is the first tagged release. While the API is behavior-complete
relative to decimal.js v10.6.0, the `0.x` series signals that small
ergonomic refinements (e.g. helper names unique to Go) may still change.
Behavioral parity is *not* expected to change — any such change would be a
regression, caught by `xvalidate/` and the ported suite.
 
 
