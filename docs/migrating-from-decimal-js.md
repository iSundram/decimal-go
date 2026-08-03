# Migrating from decimal.js

`decimal-go` is a behavioral port of `decimal.js` v10.6.0. If you are moving
an existing decimal.js codebase (or a ported algorithm) to Go, this guide maps
the API and calls out the differences that matter.

## Construction

| decimal.js | decimal-go |
|---|---|
| `new Decimal("123.45")` | `decimal.New("123.45")` |
| `Decimal.clone({...})` | `decimal.Default.Clone(&decimal.Config{...})` |
| `Decimal.config({...})` | `ctor.Config(&decimal.Config{...})` / `ctor.Set(...)` |

Both accept numbers, strings, bigints (`*big.Int`), other decimals, `-0`,
`NaN`, `±Infinity`, and `0x`/`0b`/`0o`-prefixed strings. Go additionally
accepts all integer widths and `float64`. Invalid input panics with a
`[DecimalError]`-prefixed message (decimal.js throws).

## Method names

The idiomatic Go names are short. Every decimal.js long-form name also
exists as an alias (defined in `parity.go`) so copied JS code compiles
practically verbatim:

| decimal.js | Go (idiomatic) | Go (alias) |
|---|---|---|
| `div` / `dividedBy` | `Div` | `DividedBy` |
| `divToInt` / `dividedToIntegerBy` | `DivToInt` | `DividedToIntegerBy` |
| `mod` | `Mod` | `Modulo` |
| `neg` | `Neg` | `Negated` |
| `abs` | `Abs` | `AbsoluteValue` |
| `trunc` | `Trunc` | `Truncated` |
| `clamp` | `Clamp` | `ClampedTo` |
| `cmp` | `Cmp` | `ComparedTo` |
| `equals/greaterThan/greaterThanOrEqualTo/lessThan/lessThanOrEqualTo` | `Eq Gt Gte Lt Lte` | `Equals GreaterThan GreaterThanOrEqualTo LessThan LessThanOrEqualTo` |
| `isInteger` | `IsInt` | `IsInteger` |
| `isNegative` | `IsNeg` | `IsNegative` |
| `isPositive` | `IsPos` | `IsPositive` |
| `decimalPlaces` | `Dp` | `DecimalPlaces` |
| `precision` | `Sd` | `Precision` |
| `exp` / `ln` | `Exp` / `Ln` | `NaturalExponential` / `NaturalLogarithm` |
| `log` | `Log` | `Logarithm` |
| `sqrt` / `cbrt` | `Sqrt` / `Cbrt` | `SquareRoot` / `CubeRoot` |
| `pow` | `Pow` | `ToPower` |
| `toFixed` / `toExponential` / `toPrecision` | `ToFixed` / `ToExponential` / `ToPrecision` | — |
| `toDP` / `toSD` | `ToDP` / `ToSD` | `ToDecimalPlaces` / `ToSignificantDigits` |
| `toHexadecimal` | `ToHex` | `ToHexadecimal` |
| `toJSON` | `ValueOf` | `ToJSON` |
| `toString` | `String` | `ToString` |
| `sin cos tan asin acos atan asinh acosh atanh...` | same | `Sine Cosine Tangent InverseSine ...` |

Only `Max`/`Min` differ: in decimal.js v10.6.0 the **instance** methods
`max` and `min` are commented out in the source (docs-only), so they are not
part of the real runtime API. `decimal-go` matches the real API:
`Constructor.Max` and `Constructor.Min` exist as statics; instance `Max`/`Min`
do not.

## Behavior you should know about

- **Construction keeps the full coefficient.** `new Decimal("123.4500")` has
  value 123.45 — `String()`/`ValueOf()` print `123.45`. Rounding to the
  configured precision applies on *operations*, not construction.
- **`String` vs `ValueOf` (and `-0`).** Like decimal.js, `String()` omits the
  sign of `-0` (`"0"`) while `ValueOf()`/`ToJSON()` keep it (`"-0"`).
- **`Cmp` returns a `float64`** (`-1`/`0`/`1`, or `NaN` when an operand is
  `NaN`), exactly as decimal.js's `cmp`.
- **Rounding modes** are positional constants: `RoundUp=0`, `RoundDown=1`,
  `RoundCeil=2`, `RoundFloor=3`, `RoundHalfUp=4`, `RoundHalfDown=5`,
  `RoundHalfEven=6`, `RoundHalfCeil=7`, `RoundHalfFloor=8`,
  `Euclid=9` (modulo), matching `ROUND_*` in JS.
- **1024-digit limit.** Transcendental functions that would need more
  significant digits of π/ln10 than the embedded constants provide (1024)
  throw a
  `[DecimalError] Precision limit exceeded` panic — a faithful port of
  decimal.js's `precisionLimitExceeded`. After it throws, the constructor
  setting may be left mutated until you reset it (trig path), exactly like
  the JS module.

## Concurrency

decimal.js keeps module-level mutable state (the `external`, `inexact` and
`quadrant` flags), so the recommended pattern is one cloned constructor per
concurrent context. The same applies to this port: Constructor values carry
mutable state (precision,
rounding, intermediate flags) and are **not** safe to share across
goroutines. Create one cloned constructor per goroutine:

```go
go func() {
    c := decimal.Default.Clone(nil) // isolated copy
    _ = c.New("1.1").Plus(c.New("2.2"))
}()
```

A single `*Decimal` value, once created, is immutable and safe for concurrent
reads. The original decimal.js module globals that made this unsafe
(`external`, `inexact`, `quadrant`) were moved onto the `Constructor`; the
test suite runs race-clean (`go test -race`).

## Formatting and JSON

- `MarshalJSON` produces a JSON **string** in the `valueOf()` form (this
  matches decimal.js: `JSON.stringify(D(x))` calls `toJSON` → `valueOf`).
- `MarshalText`/`UnmarshalText` give `encoding.TextMarshaler` round-trips.
- `Scan`/`Value` implement `database/sql`: values are stored as strings so the
  full coefficient and exponent survive; `NULL` scans to `NaN`.