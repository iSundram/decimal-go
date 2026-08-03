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
  original's quirk of leaving the constructor config mutated until reset on
  the π/trig path (the ln10 path restores config before panicking).

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
op-only figures isolate the arithmetic kernel (Go wins, up to ~3× on Mul,
1.5–2× on most others); `New(string)` is ~1.45× slower (Go allocation
overhead vs. JIT-optimised JS string parsing); round trips that re-parse
inputs are ~1.7× slower. We do not
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

## 11. Upstream bug found during differential fuzzing: `log(0, base)`

Differential fuzzing of decimal.js v10.6.0 against Python `mpmath`
(200+ digits of precision) found a genuine upstream bug, faithfully
inherited by this port. We deliberately preserve it: behavioral parity is
the product (§1, §10), and "fixing" it unilaterally would make decimal-go
diverge from the reference it claims to match.

### The bug

`Decimal(0).log(base)` (and `Decimal(-0).log(base)`) always returns
`-Infinity`, regardless of the base. The correct value depends on the base:

```
log_b(0) = ln(0) / ln(b) = -∞ / ln(b)
```

- `b > 1`:   `ln(b) > 0`  ⇒ `-∞`  (decimal.js — and decimal-go — are correct)
- `0 < b < 1`: `ln(b) < 0` ⇒ `+∞`  (**both return `-∞` — wrong**)
- `b == 1`:  `NaN` (handled correctly upstream by the `base.eq(1)` check)

Confirmed live:

```
Node decimal.js:   new Decimal(0).log(0.5).valueOf()  → "-Infinity"  (want "+Infinity")
decimal-go:        Default.Log("0", "0.5").String()   → "-Infinity"  (want "+Infinity")
mpmath @200dps:    mp.log(0, 0.5)                     → +inf
```

### Root cause

In decimal.js's `P.log` the zero-argument check short-circuits before the
base is ever consulted:

```js
if (arg.s < 0 || !d || !d[0] || arg.eq(1)) {
  return new Ctor(d && !d[0] ? -1 / 0 : arg.s != 1 ? NaN : d ? 0 : 1 / 0);
}
```

decimal-go mirrors this exactly: `newLogSpecial` (trans.go) returns
`newInf(-1)` for a zero argument unconditionally.

### Why we keep it (and when we won't)

- **Parity is the contract.** The 1518-case cross-validation and the 61-module
  ported suite exist to make decimal-go *indistinguishable* from decimal.js.
  Fixing this in Go alone would make `Log` diverge from the reference for
  exactly the inputs the docs promise compatibility on.
- **When we'd fix it:** if upstream ships a correction, decimal-go follows
  (with a changelog entry and regression tests). The suggested upstream fix:

  ```js
  if (d && !d[0]) {
    return new Ctor(base.gt(1) ? -1 / 0 : base.lt(1) ? 1 / 0 : NaN);
  }
  ```

- **Impact on real code:** only `log(x, base)` with `x → 0` (or `x = 0`)
  and a fractional base `0 < b < 1` — rare in practice, but real in
  information-theory / entropy-style computations where probabilities are
  used as the base. Downstream comparisons like `x.log(b).gt(0)` can flip.

### Test coverage gap

The ported `test/modules/log.js` (log_test.go) exercises `log(1, 0)`,
`log(10, 0)`, `log(10, 1)`, `log(10, Infinity)`, etc., but never
`log(0, base)` with `base ∈ (0, 1)` — upstream never tested it either.
A regression test was added to `regression_test.go` pinning the *current*
parity behavior (`-Infinity`), so if upstream fixes decimal.js and we
follow, the change is caught deliberately rather than silently.

Full hunt report (methodology, harness locations, everything else that was
fuzzed and passed): `bugs.md` in the decimal.js checkout this port was
derived from.

## 12. Upstream bug #2: `toFraction()` infinite loop under `rounding: 3`

A second bug found during the differential hunt, and it is more severe than
#1: with `rounding: 3` (ROUND_FLOOR) set, `toFraction()` **never returns**
for essentially any finite non-zero value — an infinite loop, not a
crash or exception. Confirmed in:

- decimal.js v10.6.0 (CJS and ESM), and the published npm package
  `decimal.js@10.6.0` (fresh install) — the JS process hangs (verified by
  timeout; `npm test` passes 22658/22658 because `test/modules/toFraction.js`
  hard-codes `rounding: 4`).
- this port: `c.New("0.5").ToFraction()` with `Rounding: 3` hangs
  identically (verified by timeout).

### Why it happens

The continued-fraction loop's only exit is `d2.cmp(maxD) == 1` (decimal.js
~line 2085; this port convert.go `ToFraction`). Under ROUND_FLOOR the exact
remainder `n.minus(q.times(d2))` is `-0` (IEEE-correct: x−x → −0 under
round-toward-negative, decimal.js lines 1310/1405). Next iteration
`divide(n, d, ...)` with `d = -0` yields `q = -Infinity`, so
`d2 = d0.plus(q.times(d1)) = -Infinity`, and `-Infinity < maxD` — the loop
does not break, values become `NaN`, and `NaN.cmp(maxD) == 1` is false
**forever**. A `+0` remainder (any other rounding mode) produces
`+Infinity > maxD` and the loop exits.

### Our stance

Same as §11: decimal-go inherits the reference behavior and we do not
deviate unilaterally. Unlike bug #1 (a silently wrong *value*), this one is
a denial-of-service in both languages — so it is worth stating the
boundaries explicitly:

- Only `rounding: 3` triggers it; modes 0–2, 4–8 terminate.
- Only `toFraction()` (continued-fraction expansion) is affected; the loop
  is the algorithm, not an implementation accident.
- A `maxD` small enough to terminate the CF before the `-0` remainder
  (e.g. `toFraction(1)` on `0.5`) dodges the hang in both implementations.
- If upstream fixes it, we follow. The upstream-proposed guards: break when
  the remainder is exactly zero (`!d.d[0]`) or when `d2` becomes NaN.

No unit test can pin a hang; the boundary cases above are what our ported
`toFraction_test.go` covers. The regression pin for the *parity* behavior is
documented rather than tested, by design.

Full hunt report: `bugs.md` in the decimal.js checkout (same as §11).
 
 
