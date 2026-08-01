package decimal

import (
	"math"
	"strings"
	"testing"
)

// Test harness mirroring test/setup.js of decimal.js: T.assert,
// T.assertEqual, T.assertEqualDecimal, T.assertEqualProps, T.assertException.

func i64(v int64) *int64 { return &v }
func bptr(v bool) *bool  { return &v }

// setCfg applies a full configuration to the Default constructor.
func setCfg(precision, rounding, toExpNeg, toExpPos, maxE, minE, modulo int64, crypto bool) {
	Default.Config(&Config{
		Precision: i64(precision),
		Rounding:  i64(rounding),
		ToExpNeg:  i64(toExpNeg),
		ToExpPos:  i64(toExpPos),
		MaxE:      i64(maxE),
		MinE:      i64(minE),
		Modulo:    i64(modulo),
		Crypto:    bptr(crypto),
	})
}

func assert(t *testing.T, actual bool) {
	t.Helper()
	if actual != true {
		t.Errorf("assert: expected true, got %v", actual)
	}
}

// assertEq mirrors T.assertEqual: NaN equals NaN.
func assertEq(t *testing.T, expected, actual any) {
	t.Helper()
	if jsEqual(expected, actual) {
		return
	}
	t.Errorf("assertEqual:\n  expected: %v (%T)\n  actual:   %v (%T)", expected, expected, actual, actual)
}

func jsEqual(a, b any) bool {
	if af, ok := a.(float64); ok && math.IsNaN(af) {
		bf, ok := b.(float64)
		return ok && math.IsNaN(bf)
	}
	return sameKind(a, b)
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func sameKind(a, b any) bool {
	an, aok := a.(string)
	bn, bok := b.(string)
	if aok || bok {
		return aok && bok && an == bn
	}
	fa, aok := toFloat(a)
	fb, bok := toFloat(b)
	return aok && bok && fa == fb
}

// assertEqProps mirrors T.assertEqualProps: checks d, e and s.
func assertEqProps(t *testing.T, digits []int, exponent int64, sign int8, n *Decimal) {
	t.Helper()
	if len(digits) != len(n.d) {
		t.Errorf("assertEqualProps: expected digits %v, got %v", digits, n.d)
		return
	}
	for i, w := range digits {
		if int32(w) != n.d[i] {
			t.Errorf("assertEqualProps: expected digits %v, got %v", digits, n.d)
			return
		}
	}
	if n.e != exponent || n.s != sign {
		t.Errorf("assertEqualProps: expected e=%v s=%v, got e=%v s=%v", exponent, sign, n.e, n.s)
	}
}

// assertEqDecimal mirrors T.assertEqualDecimal.
func assertEqDecimal(t *testing.T, x, y *Decimal) {
	t.Helper()
	if x.Eq(y) || x.IsNaN() && y.IsNaN() {
		return
	}
	t.Errorf("assertEqualDecimal:\n  x: %s\n  y: %s", x.ValueOf(), y.ValueOf())
}

// assertException mirrors T.assertException: the function must panic with a
// "[DecimalError]" error.
func assertException(t *testing.T, f func(), msg string) {
	t.Helper()
	ok := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err, isErr := r.(error); isErr && strings.Contains(err.Error(), "DecimalError") {
					ok = true
				}
			}
		}()
		f()
	}()
	if !ok {
		t.Errorf("assertException: expected %q to raise a DecimalError", msg)
	}
}

// nan is math.NaN() shorthand for ported tests.
func nan() float64 { return math.NaN() }

// inf is math.Inf(sign).
func inf(sign int) float64 { return math.Inf(sign) }

// math_NegZero returns negative zero.
func math_NegZero() float64 { return math.Copysign(0, -1) }
