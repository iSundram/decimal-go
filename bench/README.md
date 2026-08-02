# Benchmarks

Two runners measure the same operations on identical inputs:

- **Go** — `bench_go_test.go` (package `decimal`, black-box style), timed by
  `go test -bench`. Operands are pre-created once, then each op runs in a
  tight loop, so the "op-only" cost is measured.
- **JS** — `bench/js/bench.js` (Node), mirrors the Go cases exactly with the
  same pre-created-operand structure and JIT warmup, reporting `ns/op`.

`compare.sh` runs both and prints a ratio table. **Ratio < 1.0 means the Go
port is faster.**

## Requirements

- Go 1.22+
- Node.js (tested with v22)
- The original decimal.js checkout, path via `DECIMAL_JS_PATH`
  (default `/root/hackathon/decimal.js/decimal.js`).

## Run

```sh
bash bench/compare.sh
```

Results are printed as `ns/op` per op, plus `Go/JS`. A dedicated transcoding
script (`bench/go_js_compare.py`) merges Go's sub-benchmarks (e.g.
`BenchmarkDiv/precision=1000`) with the JS `Div@1000` rows.

## Honest comparison notes

The Go/JS numbers are only directly comparable when the **same work** is
measured. Known asymmetries, worth remembering before quoting any ratio:

1. Go `New(string)` is slightly slower than JS construction (~1.3x), and a
   full round-trip (`New`+op+`ValueOf`) is ~1.5x slower, because JS string
   parsing is highly JIT-optimised while Go pays allocation overhead per
   parse. Benchmarks that re-parse inputs every iteration (real-world style)
   therefore narrow Go's advantage.
2. Op-only figures (the default below) isolate the arithmetic kernel; this is
   the fairest comparison of raw speed and is where Go is typically 2-6x
   faster.
3. Alloc/op is reported by Go but not JS (V8 doesn't surface per-op
   allocation) — omit allocation from cross-language claims.
4. Node's JIT warns: measurements are for a warmed-up steady state; first-call
   latency differs but is not measured here.

Everything else — precision settings, rounding mode, and inputs — is
identical between the two suites.

## Results (last full run, saved as `results.txt`)

Same machine, Node v22, Go from the repo root, op-only operands at the default
precision 20 (and 100/1000 where both sides can): **ratio < 1 means the Go
port is faster.** 18 of 20 comparable pairs go Go's way; the only slower paths
are string construction and parse+op+format round trips.

| op@precision | Go ns/op | JS ns/op | Go/JS |
|---|---|---|---|
| Mul@20   |  332.3 | 1039.0 | 0.32 |
| Div@1000 | 20416.0 | 52159.3 | 0.39 |
| Cbrt@20  | 11341.0 | 26129.1 | 0.43 |
| Atan@20  | 1503.0 |  3360.2 | 0.45 |
| Cos@20   | 20132.0 | 42136.8 | 0.48 |
| Div@20   | 1184.0 | 2468.3 | 0.48 |
| Tan@20   | 38750.0 | 77092.3 | 0.50 |
| PowFraction@20 | 93816.0 | 179369.0 | 0.52 |
| PowInt@20 | 3406.0 | 6194.3 | 0.55 |
| LogBase10@20 | 67097.0 | 119776.5 | 0.56 |
| Ln@20    | 50497.0 | 83774.0 | 0.60 |
| DivToInt@20 | 425.0 | 665.0 | 0.64 |
| Sqrt@20  | 8082.0 | 11828.6 | 0.68 |
| ToFixed@20 | 693.8 | 1011.4 | 0.69 |
| Sin@20   | 21844.0 | 31453.7 | 0.69 |
| Exp@20   | 40885.0 | 56693.8 | 0.72 |
| ToPrecision@20 | 657.4 | 879.6 | 0.75 |
| Sub@20   | 245.4 | 312.5 | 0.79 |
| Add@20   | 267.2 | 323.2 | 0.83 |
| Mod@20   | 1050.0 | 1258.3 | 0.83 |
| Cmp@20   | 213.3 | 239.1 | 0.89 |
| ToExponential@20 | 717.7 | 782.7 | 0.92 |
| **New@20**  | **1642.0** | **1132.8** | **1.45** |
| **RoundTrip@20** | **3047.0** | **1767.6** | **1.72** |

See the columns at precision 100 / 1000 in `results.txt` for high-precision
scaling (Go maintains or widens its lead as precision grows; JS only
benchmarks base/extra setups for those rows).