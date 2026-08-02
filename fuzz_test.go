package decimal

import (
	"strings"
	"testing"
)

// smallConfig returns a constructor with sane settings for fuzzing. Each
// invocation gets a fresh clone so fuzz iterations never share mutable state.
func smallConfig() *Constructor {
	return Default.Clone(&Config{
		Precision: I64(20),
		Rounding:  I64(4),
		MaxE:      I64(9e15),
		MinE:      I64(-9e15),
	})
}

// expectDecimalError absorbs the documented "[DecimalError]" panic (the Go
// equivalent of decimal.js throwing on invalid input) and turns anything else
// into a test failure.
func expectDecimalError(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		msg := ""
		switch e := r.(type) {
		case error:
			msg = e.Error()
		case string:
			msg = e
		}
		if !strings.HasPrefix(msg, decimalError) {
			t.Fatalf("unexpected panic: %v", r)
		}
	}
}

// FuzzNewString parses arbitrary strings: a [DecimalError] panic is the
// documented contract, anything else is a bug.
func FuzzNewString(f *testing.F) {
	for _, s := range []string{
		"0", "-0", "1", "0.5", "1e3", "1e-9", "-0.001", "NaN",
		"Infinity", "0x1f", "0b101", "0o17", "1e100000",
		"3.141592653589793238462643383279",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		defer expectDecimalError(t)
		c := smallConfig()
		x := c.New(s)
		// Parsing must be reproducible: parse(x.String()) == x.
		y := c.New(x.String())
		if x.Cmp(y) != 0 && !(x.IsNaN() && y.IsNaN()) {
			t.Fatalf("round-trip mismatch: %q -> %q -> %q", s, x, y)
		}
		_ = x.String()
	})
}

// maxQuotientWords bounds the number of base-1e7 words a fuzzed Div/Mod
// quotient may have. Division (and hence Mod, which divides first) internal-
// expands the integer quotient digit-by-digit up to x.e-y.e words; with a
// p/hex-pair input like "0b1p2738415519256" that is astronomically large
// (~8e11 words) and kills the fuzz worker with OOM — as does the reference
// decimal.js, which hangs identically on the same input. Pairs beyond this
// bound are skipped: they exercise nothing the port could ever compute.
// See docs/DECISIONS.md §7.
const maxQuotientWords = 1e5

// FuzzArith crosses the binary operators over arbitrary string pairs;
// sanitizers panic recovery only allows [DecimalError] failures.
func FuzzArith(f *testing.F) {
	for _, p := range [][2]string{
		{"0", "0"}, {"1", "1"}, {"2", "0.5"}, {"-1", "3"},
		{"1e-5", "1e5"}, {"9999999999", "0.0000000001"},
		{"1", "0.00000000000000000001"}, {"-0", "2"},
	} {
		f.Add(p[0], p[1])
	}
	f.Fuzz(func(t *testing.T, a, b string) {
		defer expectDecimalError(t)
		c := smallConfig()
		x, y := c.New(a), c.New(b)
		_ = x.Plus(y)
		_ = x.Minus(y)
		_ = x.Times(y)
		// Div and Mod expand a quotient with roughly x.e-y.e words;
		// skip pairs whose quotient cannot be materialized (see above).
		if y.d == nil || x.d == nil || y.s == 0 || y.d[0] == 0 || x.e-y.e > maxQuotientWords {
			return
		}
		_ = x.Div(y)
		_ = x.Mod(y)
	})
}

// FuzzPowNew exercises Pow and Sqrt.
func FuzzPowNew(f *testing.F) {
	for _, p := range [][2]string{
		{"0", "0"}, {"2", "10"}, {"-1", "0.5"}, {"2", "1e6"},
		{"1e10", "2"}, {"0.9", "-1e10"},
	} {
		f.Add(p[0], p[1])
	}
	f.Fuzz(func(t *testing.T, a, b string) {
		defer expectDecimalError(t)
		c := smallConfig()
		_ = c.New(a).Pow(c.New(b))
		_ = c.New(a).Sqrt()
	})
}

// FuzzTrig repairs large arguments into quadrant handling; a DecimalError
// panic (e.g. the legit 1024-dig distortion guard) is acceptable.
func FuzzTrig(f *testing.F) {
	for _, s := range []string{"0", "1", "-1", "3.14159", "1e-3", "0.5", "100", "-100"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		defer expectDecimalError(t)
		c := smallConfig()
		x := c.New(s)
		// An [DecimalError] from the trig guard must not leave the shared
		// config in an unusable state for the next clone.
		_ = x.Sin()
		// A second constructor still works.
		c2 := smallConfig()
		_ = c2.New(s).Cos()
	})
}
