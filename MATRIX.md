# Operations matrix

A mapping of every decimal.js v10.6.0 operation to its decimal-go counterpart,
its rounding semantics, and the edge cases covered by the ported test suite
plus the proof tests added in this worker. All Go behavior matches the
original JS implementation, verified by `xvalidate/` (byte-for-byte output
parity on the shared corpus) and the worker commits.

## Core arithmetic

| Op | decimal.js | Go function | Rounding semantics | Edge cases |
|----|-----------|-------------|--------------------|------------|
| Add | `plus` | `Plus` (`ops.go`) | exact; only `finalise`'d when operands have max fractional digits | `-0 + +0 = +0`; huge+small exponent; Infinity/NaN propagation |
| Subtract | `minus` | `Minus` | same as add | `-0 - +0 = +0`, unequal signs; `Inf - Inf = NaN` |
| Multiply | `times` | `Times` | exact product then `finalise(pr, rm)`; guard for 9999/0000 | mixed exponent overflow; `0 × Inf = NaN` |
| Divide | `div` | `Div` | result rounded to `precision` via long division (`divide`, `ops.go:493`) | `0/0=NaN`, `Inf/Inf=NaN`, `x/0=±Inf`; `maxE` overflow; non-terminating guarded by `finalise` |
| Divide-to-int | `divToInt` | `DivToInt` | truncates to integer (`rm=1`) regardless of rounding mode | division by zero throws; negative truncation toward zero |
| Modulo | `mod` | `Mod` | sign of divisor; `modulo` 0..8 + `EUCLID` (9) | remainder magnitude < |divisor|; modulo rounding-mode matrix tested |
| Power | `pow` | `Pow` | integer exponent -> `intPow`; else `powNumber`+`exp/ln` | `0^0=1`; negative base, non-integer exponent -> NaN; `(-1)^±Inf=NaN`; `pow10` int64 overflow regression |

## Unary & rounding

| Operation | decimal.js | Go function | Edge cases |
|---|---|---|---|
| Negation | `neg` | `Neg` | preserves `-0` |
| Absolute | `abs` | `Abs` | `abs(-0)=+0` |
| Truncate | `trunc` | `Trunc` | truncation toward zero |
| Ceil | `ceil` | `Ceil` | rounding mode independent |
| Floor | `floor` | `Floor` | `-0` sign handling |
| Round | `round` | `Round` | ties to rounding mode 0 (half-up default) |
| Square root | `sqrt` | `Sqrt` (`trans.go:95`) | exact squares no guard; 4999/9999 guard in `trans.go:370`; negative -> NaN |
| Cube root | `cbrt` | `Cbrt` (`trans.go:201`) | sign-preserving; same guard mechanism |
| Natural exp | `exp` | `Exp` (`trans.go:417`) | `getLn10` fallback precision; Infinity guard |
| Natural log | `ln` | `Ln` (`trans.go:508`) | `ln(0)=-Inf`, `ln(<0)=NaN`; `[49]9999` rounding re-run guard |
| Log base | `log(x[,base])` | `Log` | `log(0)`, `log(1)=0`, negative base/arg -> NaN; computed as `ln(x)/ln(base)` with truncation guards |
| Log2 / log10 | `log2`, `log10` | `constructor.Log2`/`Log10` | base 2/10 special-cased |
| Inverse trig | `asin acos atan atan2` | `Asin Acos Atan` + `Constructor.Atan2` | domain errors NaN; `atan2(0,0)=0` |
| Hyperbolic | `sinh cosh tanh` | `Sinh Cosh Tanh` | `sinh(±Inf)=±Inf`; `tanh(±Inf)=±1`; large arg precision-limit like origin |
| Inverse hyperbolic | `asinh acosh atanh` | `Asinh Acosh Atanh` | `acosh(<1)=NaN`; `atanh(±1)=∓Inf` |

## Trig implementation notes

- `toLessThanHalfPi` (`trig.go:144`) performs quadrant reduction; the `quadrant`
  state is per-constructor (concurrency-safe), matching decimal.js's
  `external` flag model.
- `taylorSeries` (`trig.go:91`) terminates on convergence digits; the
  `precisionLimitExceeded` panic (sd > 1024 bits) is a **faithful** port of
  decimal.js: `getPi` estimates required precision and throws when the result
  needs more than the configured `maxDigits` (1024), leaving config mutated
  until the user resets it, exactly like the original.
- `getPi` / `getLn10` (`trans.go:74/86`) cache constants per constructor.

## Rounding & digit helpers

| Operation | decimal.js | Go function | Edge cases |
|---|---|---|---|
| Decimal places | `dp` | `Dp` | returns float64 (JS number), NaN for NaN |
| Significant digits | `sd` | `Sd` | `sd(true)` counts trailing zeros; `-sd` for non-integer |
| To DP | `toDecimalPlaces` | `ToDP` | rounding mode arg (4 default) |
| To SD / precision | `toSignificantDigits` | `ToSD`, `ToPrecision` | dp from param optional; ties handled |
| Clamp | `clamp` | `Clamp` | `min > max` no-op; `-0` preserved |

## String / base conversion

| Function | decimal.js | Go function | Edge cases |
|---|---|---|---|
| `toString` | `toString` | `String` | respects `toExpPos`/`toExpNeg`; `-0` omitted (JS parity) |
| `valueOf` / `toJSON` | `valueOf` | `ValueOf` / `ToJSON` | includes `-0` and exponent; JSON-safe |
| `toFixed` | `toFixed` | `ToFixed` | dp beyond 1e9 throws; ties per rounding |
| `toExponential` | `toExponential` | `ToExponential` | dp required; `toExpNeg` ignored |
| `toPrecision` | `toPrecision` | `ToPrecision` | sd from 0..prec; exponent form when needed |
| Binary | `toBinary` | `ToBinary` | base 2 with fractional part |
| Octal | `toOctal` | `ToOctal` | base 8 |
| Hexadecimal | `toHexadecimal` | `ToHex` | base 16 uppercase; `base256` unused |
| Fraction | `toFraction` | `ToFraction` | `max` input; IEEE-style |
| `toNearest` | `toNearest` | `ToNearest` | nearest integer / `y` |
| `toNumber` / JS Number | `toNumber` | `ToNumber` / `Float64` | NaN on overflow/underflow |

## Comparison & predicates

`Cmp` returns `float64 -1/0/1` (or NaN for NaN operand) for JS parity; the
boolean comparators `Eq Gt Gte Lt Lte`, predicates `IsNaN IsFinite IsInt
IsNeg IsPos IsZero` all present, with `IsInteger`/`IsNegative`/ants aliases.

## Statics (constructor)

All decimal.js statics present: `abs acos acosh add asin asinh atan atan2
atanh cbrt ceil clamp cos cosh div exp floor hypot isDecimal ln log log10
log2 max min mod mul pow random round sign sin sinh sqrt sub sum tan tanh
trunc`, plus `clone`/`config`/`Default`/`New`/`IsDecimal`.

## Edge-case coverage summary

- **Bug regressions** (`regression_test.go`): pow10 int64 overflow/hang,
  Pow negative-base index-out-of-range, `1^±Inf`. All pass.
- **Property tests** (`property_test.go`): round-trip, add/mul commutativity,
  inverse add/sub & mul/div, sqrt², cbrt³, exp∘ln, cmp antisymmetry, exact
  sqrt, modulo range — all green.
- **Stress tests** (`stress_test.go`): precision 400–2048, int64/uint64
  boundary, adversarial inputs, NaN/Inf, 64-goroutine clones under `-race`.
- **Input matrix** (`input_test.go`): number/string/bigInt/Decimal, -0, NaN,
  ±Infinity, hex/bin/oct strings, maxE/minE deep-underflow/overflow.
- **Cross-validation** (`xvalidate/`): 1518 case lines identical to real
  decimal.js.
- **Per-op suites**: each op port inherits the original decimal.js test table
  (TestPlus, TestMinus, TestTimes, TestDiv, TestDivToInt, TestIntPow, TestPow,
  TestSqrt, TestCbrt, TestExp, TestLn, TestLog, TestTrig/Hyperbolics, TestTo*
  conversions, TestConfig/TestRandom/TestClone, TestRound etc.) — all green.

## Honest comparison notes

- Go `new(string)` parsing benchmarks ~1.3× slower than decimal.js; the full
  add-from-string round trip ~1.5× slower; op-only arithmetic (mul/div/sqrt/
  ln/trig) is 2–6× faster. See `bench/README.md` for the full ratio table.
- deconv: the port preserves decimal.js's error and configuration semantics,
  including the 1024-digit `precisionLimitExceeded` for transcendental
  constants and the shared-constructor model: one clone per concurrent
  context (module-level mutable flags were moved onto the Constructor to make
  per-clone use race-free).