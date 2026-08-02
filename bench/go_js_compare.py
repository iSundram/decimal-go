#!/usr/bin/env python3
"""Join Go sub-benchmark output with JS benchmark output into a ratio table.

Go lines look like:
    BenchmarkDiv/precision=1000-2  42  20414 ns/op ...
JS lines look like:
    Div,1000,60936.7,0
"""
import re
import sys


def parse_go(path):
    out = {}
    for line in open(path):
        m = re.match(r"^Benchmark(\w+)(?:/(?:precision=)(\d+))?-\d+\s+\d+\s+([\d.]+) ns/op", line)
        if not m:
            m = re.match(r"^Benchmark(\w+)(?:/(\w+)=\d+)?-\d+\s+\d+\s+([\d.]+) ns/op", line)
        if m:
            name, pr, ns = m.group(1), m.group(2) or "20", m.group(3)
            out[f"{name}@{pr}"] = float(ns)
    return out


def parse_js(path):
    out = {}
    for line in open(path):
        parts = line.strip().split(",")
        if len(parts) >= 3 and parts[-3].replace(".", "", 1).isdigit():
            name, pr, ns = parts[0], parts[1], parts[2]
            out[f"{name}@{pr}"] = float(ns)
    return out


def main():
    go = parse_go(sys.argv[1])
    js = parse_js(sys.argv[2])
    print(f"{'op@precision':<18} {'Go ns/op':>10} {'JS ns/op':>10} {'Go/JS':>9}")
    for key in sorted(set(go) | set(js)):
        g = go.get(key)
        j = js.get(key)
        if g is not None and j is not None:
            ratio = g / j
            flag = " <-- Go slower" if ratio > 1.5 else ""
            print(f"{key:<18} {g:>10.1f} {j:>10.1f} {ratio:>9.2f}{flag}")
        elif g is not None:
            print(f"{key:<18} {g:>10.1f} {'(no JS)':>10} ")
        else:
            print(f"{key:<18} {'(no Go)':>10} {j:>10.1f} ")


if __name__ == "__main__":
    main()