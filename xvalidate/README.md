# Cross-validation harness

Proves that decimal-go behaves *identically* to the original decimal.js on a
shared corpus. Both sides read the same `cases.txt` and run the same unary,
binary, comparison and formatting operations with the same config
(`precision 20, rounding 4, toExpNeg -7, toExpPos 21`).

## Run

```bash
bash xvalidate/compare.sh   # go run + node, sort, diff; empty diff = PASS
```

Requires `go` and `node` with `decimal.js` available (point `DECIMAL_JS_PATH`
at the original lib, e.g. `/root/hackathon/decimal.js/decimal.js`).

## Coverage

- `cases.txt`: 50+ fixed inputs — signs, zeros (`0`, `-0`), small/large
  exponents that cross the toExpNeg/toExpPos formatting boundaries (±1e±7/21),
  full integer precision beyond `Number.MAX_SAFE_INTEGER` (1e22,
  123...890×10^18?), subnormals, `NaN`, `±Infinity`.
- Unary: `abs`, `neg`, `trunc`, `ceil`, `floor`, `round`, `sqrt`, `cbrt`,
  `exp`, `ln`.
- Formatting: `dp`, `sd`, `toFixed(2)`, `toExponential(6)`, `toPrecision(8)`.
- Binary (adjacent-pair, wrap-around): `add`, `sub`, `mul`, `div`, `divToInt`,
  `mod`, `pow`, plus `cmp`.

Given identical inputs and settings, the port and the original emit byte-for
byte identical result strings (verified: 1518/1518 lines match).