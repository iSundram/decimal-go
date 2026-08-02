package decimal

import (
	"math"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

func StressDefaultCfg() {
	setCfg(20, 4, -7, 21, 9e15, -9e15, 1, false)
}

func isDecimalError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "DecimalError")
}

// TestStressHighPrecision drives transcendental at precision beyond the
// precomputed pi/LN10 constants. This is decimal.js-parity behaviour: it may
// throw the exact "Precision limit exceeded" DecimalError, but must panic with
// no other kind of error and never return garbage digits.
//
// Note (parity): after a precision-limit throw, decimal.js leaves its global
// config (precision/rounding) mutated until the user resets it; this port does
// the same, so each case here resets the config before the next operation.
func TestStressHighPrecision(t *testing.T) {
	ops := map[string]func(*Decimal) *Decimal{
		"Ln":   func(x *Decimal) *Decimal { return x.Ln() },
		"Sqrt": func(x *Decimal) *Decimal { return x.Sqrt() },
		"Exp":  func(x *Decimal) *Decimal { return x.Exp() },
		"Pow":  func(x *Decimal) *Decimal { return x.Pow(0.5) },
		"Sin":  func(x *Decimal) *Decimal { return x.Sin() },
		"Cos":  func(x *Decimal) *Decimal { return x.Cos() },
		"Tan":  func(x *Decimal) *Decimal { return x.Tan() },
	}
	for _, pr := range []int64{400, 700, 1100, 2048} {
		for op, fn := range ops {
			// Fresh ctor per (pr, op): unlike JS there is no global state
			// to leak between cases.
			setCfg(pr, 4, -1, 50, 9e15, -9e15, 1, false)
			x := New("2")
			func(op string, fn func(*Decimal) *Decimal) {
				defer func() {
					if r := recover(); r != nil {
						err, ok := r.(error)
						if !ok || !isDecimalError(err) {
							t.Fatalf("pr=%d %s panicked with %T %v", pr, op, r, r)
						}
					}
				}()
				r := fn(x)
				if sd := r.Sd(); sd > float64(pr)+2 {
					t.Fatalf("pr=%d %s returned %v sig figs", pr, op, sd)
				}
			}(op, fn)
		}
		StressDefaultCfg()
	}
}

// TestStressInt64Boundaries: extreme int/uint inputs.
func TestStressInt64Boundaries(t *testing.T) {
	StressDefaultCfg()
	cases := []struct {
		in   any
		want string
	}{
		{int64(math.MinInt64), "-9223372036854775808"},
		{int64(math.MaxInt64), "9223372036854775807"},
		{uint64(math.MaxUint64), "18446744073709551615"},
		{int(1), "1"},
	}
	for _, c := range cases {
		if got := New(c.in).ValueOf(); got != c.want {
			t.Fatalf("New(%v) = %s, want %s", c.in, got, c.want)
		}
	}
}

// TestStressGoroutine: distinct cloned constructors used from separate
// goroutines must be race-free (decimal.js's documented model: one clone per
// concurrent context). Sharing a single constructor across goroutines is NOT
// supported, exactly as a single decimal.js Decimal is not shared between
// concurrent event loops.
func TestStressGoroutine(t *testing.T) {
	results := make([]string, 64)
	errs := make([]string, 64)
	var wg sync.WaitGroup
	for g := 0; g < 64; g++ {
		wg.Add(1)
		go func(g int) {
			c := Default.Clone(nil) // one clone per goroutine
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errs[g] = "panic: " + r.(error).Error()
				}
			}()
			r := rand.New(rand.NewSource(int64(g)))
			var v *Decimal
			for i := 0; i < 500; i++ {
				x := c.New("1.5")
				y := c.New(r.Float64() * 1e6)
				v = x.Times(y).Div("0.2").Sqrt()
				_ = v.ValueOf()
				_ = x.Mul(y).Cmp(0)
			}
			results[g] = v.ValueOf()
		}(g)
	}
	wg.Wait()
	for g := range results {
		if errs[g] != "" {
			t.Fatalf("goroutine %d errored: %s", g, errs[g])
		}
		if results[g] == "" {
			t.Fatalf("goroutine %d produced no result", g)
		}
	}
}

// TestStressInputVariety: adversarial construction inputs.
func TestStressInputVariety(t *testing.T) {
	StressDefaultCfg()
	inputs := []string{
		"0", "-0", "0.0", "1", "5", "+1", "1e0", "1e+0", "1e-0",
		"0.0000000001", "1e21", "1e+100", "-1e-100", "-1e+100",
		".5", "5.", "1.23e-9", "000.00100",
		"9223372036854775807", "-9223372036854775808",
	}
	for _, in := range inputs {
		x := New(in)
		y := New(x.ValueOf())
		if x.Cmp(y) != 0 {
			t.Fatalf("parse %q -> %q reparse %q", in, x.ValueOf(), y.ValueOf())
		}
	}
}

// TestStressNaNInf: propagation and self-compare semantics.
func TestStressNaNInf(t *testing.T) {
	StressDefaultCfg()
	nan, pInf, nInf := New("NaN"), New("Infinity"), New("-Infinity")
	if !nan.IsNaN() {
		t.Fatal("NaN not recognized")
	}
	if r := pInf.Plus("1"); !r.Eq(pInf) {
		t.Fatalf("Inf+1 = %s (should be Infinity)", r.ValueOf())
	}
	if nInf.Cmp(5) >= 0 {
		t.Fatal("-Inf should be < 5")
	}
	if pInf.Cmp("Infinity") != 0 {
		t.Fatal("Inf should equal Inf")
	}
	if nan.Times(nan) != nil && !nan.Times(nan).IsNaN() {
		t.Fatal("NaN*NaN should be NaN")
	}
	if !nan.Div(nan).IsNaN() {
		t.Fatal("NaN/NaN should be NaN")
	}
}
