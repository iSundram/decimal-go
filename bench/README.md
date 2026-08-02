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