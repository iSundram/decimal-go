#!/usr/bin/env bash
# Compare decimal-go vs decimal.js performance.
#
# Run from the repo root:   bash bench/compare.sh
# Requires: go, node, and the original decimal.js checkout (path via
# DECIMAL_JS_PATH, default /root/hackathon/decimal.js/decimal.js).
#
# Go sub-benchmarks are named e.g. "BenchmarkDiv/precision=1000"; JS ops are
# keyed "Div,1000". This script joins them and prints a ratio table
# (ratio < 1 = the Go port is faster).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DECIMAL_JS_PATH="${DECIMAL_JS_PATH:-/root/hackathon/decimal.js/decimal.js}"

echo "== Go benchmarks =="
(cd "$ROOT" && go test -bench . -run '^$' -count=1 . 2>/dev/null) | grep -E '^Benchmark' > /tmp/go_ops.txt
cat /tmp/go_ops.txt

echo
echo "== JS benchmarks =="
DECIMAL_JS_PATH="$DECIMAL_JS_PATH" node "$ROOT/bench/js/bench.js" | tee /tmp/js_ops.txt

echo
echo "== Comparison (ns/op; ratio>1 means Go slower) =="
printf "%-18s %-14s %-14s %-10s\n" "op@precision" "Go" "JS" "Go/JS"

python3 "$ROOT/bench/go_js_compare.py" /tmp/go_ops.txt /tmp/js_ops.txt