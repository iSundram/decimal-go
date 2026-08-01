# Porting decimal.js test modules to Go

Each JS module in `/root/hackathon/decimal.js/test/modules/<name>.js` is ported
to `/root/hackathon/decimal-go/<name>_test.go` in package `decimal` (white-box).

Goal: a 1:1 behavioral port of every assertion, keeping the same order and
coverage. Do not skip tests.

## Harness (harness_test.go)

- `assert(t, cond bool)` — mirrors `T.assert(cond)`
- `assertEq(t, expected, actual any)` — mirrors `T.assertEqual(expected, actual)`
  (NaN == NaN for float64)
- `assertEqProps(t, digits []int, exponent int64, sign int8, n *Decimal)` —
  mirrors `T.assertEqualProps(digits, exponent, sign, new Decimal(n))`
- `assertEqDecimal(t, x, y *Decimal)` — mirrors `T.assertEqualDecimal`
- `assertException(t, f func(), msg string)` — mirrors `T.assertException(func, msg)`:
  f must panic with an error containing "DecimalError"
- `setCfg(precision, rounding, toExpNeg, toExpPos, maxE, minE, modulo int64, crypto bool)` —
  full `Decimal.config({...})`
- `i64(v) *int64`, `bptr(v) *bool` — for partial `Default.Config(&Config{...})`
- `nan() float64`, `inf(sign int) float64`, `math_NegZero() float64`

## Conversions

- `T('name', function () { ... })` → `func TestXxx(t *testing.T) { ... }`
- `Decimal` (the constructor) → the package default: `New(v)`,
  `Default.M(...)` for static methods, `Default.Precision` etc for props.
- `Decimal.config({ precision: 40, rounding: 4, ... })` (full) →
  `setCfg(40, 4, -9e15, 9e15, 9e15, -9e15, 1, false)`; partial →
  `Default.Config(&Config{Precision: i64(100)})`
- `Decimal.toExpNeg = Decimal.toExpPos = 0;` → `Default.ToExpNeg = 0; Default.ToExpPos = 0`
- `new Decimal(v)` → `New(v)`
- methods camelCase same name: `.plus` → `.Plus`, `.divToInt` → `.DivToInt`,
  `.toDP(2, 1)` → `.ToDP(2, 1)`, `valueOf()` → `.ValueOf()`, etc.
- `NaN` → `nan()`, `Infinity` or `1/0` → `inf(1)`, `-Infinity` or `-1/0` →
  `inf(-1)`, `-0` → `math_NegZero()`.
- **JS numeric literals are float64 in Go**: write them with `.0` or exponent
  so they are floats, NOT ints (e.g. `t('4', 4)` → `tf("4", 4.0)`). This is
  essential: Go would treat a bare `4` as int64 and parse it exactly, while
  JS parses the float. Integers ≤ 2^53 give identical results either way, but
  use float literals consistently. String literals stay strings.
  *Decimal values*: `new Decimal(123)` in JS → keep `New(123)` — int is fine
  there if the JS number was an exact integer ≤ 2^53; when in doubt use floats.
- `T.assertEqual(NaN, x.sd())` → `assertEq(t, nan(), x.Sd())` (Sd returns
  float64, NaN for non-finite).
- `T.assertException(function () { new Decimal(x); }, '...')` →
  `assertException(t, func() { New(x) }, "...")`
- `N = new Decimal(...)` local vars just become Go locals.
- Loops like `for (i = 0; i < n; i++)` port directly. Random-data tests stay.

## Rules

- Do NOT modify any non-`_test.go` file, `harness_test.go`, or this file.
- Keep test order and do not drop assertions.
- If a test fails because of a suspected library bug, double-check your
  conversion first (compare against running the same snippet in node with
  the original decimal.js: `cd /root/hackathon/decimal.js && node -e "..."`).
  If it's a real library bug, report it in your final message, do NOT fix it.
- Run `cd /root/hackathon/decimal-go && go test -run TestXxx` until green.
- Do NOT run git commands.
- Test function names: `Test<Module>` with first letter capitalised,
  e.g. `toFixed.js` → `TestToFixed`, `toDP.js` → `TestToDP`, `dpSd.js` →
  `TestDpSd`, `isFiniteEtc.js` → `TestIsFiniteEtc`, `minAndMax.js` →
  `TestMinAndMax`, `powSqrt.js` → `TestPowSqrt`, `intPow.js` → `TestIntPow`.
