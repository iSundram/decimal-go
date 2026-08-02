#!/usr/bin/env bash
# Cross-validate decimal-go against the original decimal.js on a shared corpus.
#
#   bash xvalidate/compare.sh
#
# Reads xvalidate/cases.txt, runs the identical unary/format/binary cases
# through both decimal-go (go run) and decimal.js (node), sorts the two line
# sets and diffs them. Zero diff output means full behavioral parity.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DECIMAL_JS_PATH="${DECIMAL_JS_PATH:-/root/hackathon/decimal.js/decimal.js}"

echo "== Go results =="
(cd "$ROOT" && go run ./xvalidate/xvalid.go) | sort > /tmp/xv_go.txt
echo "  $ROOT lines: $(wc -l < /tmp/xv_go.txt)"

echo "== JS results =="
(cd "$ROOT" && DECIMAL_JS_PATH="$DECIMAL_JS_PATH" node xvalidate/x.js) | sort > /tmp/xv_js.txt
echo "  $(wc -l < /tmp/xv_js.txt) lines"

echo "== diff (empty = parity) =="
if diff -u /tmp/xv_js.txt /tmp/xv_go.txt > /tmp/xv_diff.txt; then
  echo "PASS: Go and decimal.js produce identical results."
else
  echo "DIFF FOUND:"
  cat /tmp/xv_diff.txt
  exit 1
fi