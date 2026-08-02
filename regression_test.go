package decimal

import "testing"

// Regression tests for the three bugs found and fixed during porting
// (commit "Port 48 test modules; fix pow10 overflow in finalise, Pow index
// OOB, powNumber +1^Inf"). Each would crash or produce wrong results before
// the fix.

// Regression 1: pow10 int64 overflow in finalise. Rounding a value whose
// rounding digit lies far inside a base-1e7 word used to call
// w / pow10(big) which overflowed int64 and produced garbage digits or a
// precision-limit panic. divPow10 clamps, so deep rounding must succeed.
func TestRegressionPow10Overflow(t *testing.T) {
	StressDefaultCfg()
	for _, v := range []string{
		"1.0000000999999994",
		"999999999999999.00000005",
		"0.0000009999999999999",
		"12345678901234567.00000005",
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("%s ToDP panicked: %v", v, r)
				}
			}()
			if got := New(v).ToDP(6).ValueOf(); got == "" {
				t.Fatalf("empty ToDP result for %s", v)
			}
			if got := New(v).ToSD(6).ValueOf(); got == "" {
				t.Fatalf("empty ToSD result for %s", v)
			}
		}()
	}
}

// Regression 2: Pow index out of range. Raising a negative base to an integer
// exponent used to index y.d[e] without bounds checking (JS reads undefined
// and `undefined & 1 === 0`); Go slid out of the slice. word() returns 0,
// so negative bases to odd/even exponents stay correct.
func TestRegressionPowNegIndex(t *testing.T) {
	StressDefaultCfg()
	for _, a := range []string{"-2", "-5", "-123.456"} {
		for _, n := range []int64{2, 3, 4, 5, 100} {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("(%s)^%d panicked: %v", a, n, r)
					}
				}()
				_ = New(a).Pow(n)
			}()
		}
	}
	if r := New("-2").Pow(5); r.ValueOf() != "-32" {
		t.Fatalf("(-2)^5 = %s, want -32", r.ValueOf())
	}
	if r := New("-5").Pow(3); r.ValueOf() != "-125" {
		t.Fatalf("(-5)^3 = %s, want -125", r.ValueOf())
	}
}

// Regression 3: powNumber with |base| == 1 and exponent +-Infinity must be
// NaN (ES Lex 15.8.2.13). Go's math.Pow returns 1; the port wraps it.
func TestRegressionPowOneInf(t *testing.T) {
	StressDefaultCfg()
	for _, b := range []string{"1", "-1"} {
		for _, e := range []string{"Infinity", "-Infinity"} {
			r := New(b).Pow(New(e))
			if !r.IsNaN() {
				t.Fatalf("%s ^ %s = %s, want NaN", b, e, r.ValueOf())
			}
		}
	}
}

func TestRegressionAll(t *testing.T) {
	t.Run("Pow10Overflow", TestRegressionPow10Overflow)
	t.Run("PowNegIndex", TestRegressionPowNegIndex)
	t.Run("PowOneInf", TestRegressionPowOneInf)
}
