package decimal

import "testing"

// fullVals covers arithmetic/comparison inc. extreme exponents, which all
// basic ops handle without issue.
var fullVals = []string{
	"-123.456", "-2", "-0", "0", "0.5", "2", "123.456", "1e-2000", "1e2000",
}

// transVals keeps magnitudes small; transcendental ops (and inverse trig
// domain) need well-scaled inputs. Extremes like 1e2000 legitimately throw
// or hang in both decimal.js and this port.
var transVals = []string{
	"-123.456", "-2", "-0.5", "-0", "0", "0.5", "2", "3.14", "123.456", "1e-10", "-1e-10",
}

// parityTest verifies that every JS long-name alias returns exactly the same
// result as its idiomatic Go counterpart, across representative values.
func TestParityAliases(t *testing.T) {
	setCfg(40, 4, -9e15, 9e15, 9e15, -9e15, 1, false)

	// Unary aliases on the transcendental-safe set.
	type unary struct {
		alias, primary func(*Decimal) *Decimal
	}
	unaries := []unary{
		{func(x *Decimal) *Decimal { return x.AbsoluteValue() }, func(x *Decimal) *Decimal { return x.Abs() }},
		{func(x *Decimal) *Decimal { return x.Negated() }, func(x *Decimal) *Decimal { return x.Neg() }},
		{func(x *Decimal) *Decimal { return x.Truncated() }, func(x *Decimal) *Decimal { return x.Trunc() }},
		{func(x *Decimal) *Decimal { return x.SquareRoot() }, func(x *Decimal) *Decimal { return x.Sqrt() }},
		{func(x *Decimal) *Decimal { return x.CubeRoot() }, func(x *Decimal) *Decimal { return x.Cbrt() }},
		{func(x *Decimal) *Decimal { return x.NaturalExponential() }, func(x *Decimal) *Decimal { return x.Exp() }},
		{func(x *Decimal) *Decimal { return x.NaturalLogarithm() }, func(x *Decimal) *Decimal { return x.Ln() }},
		{func(x *Decimal) *Decimal { return x.Sine() }, func(x *Decimal) *Decimal { return x.Sin() }},
		{func(x *Decimal) *Decimal { return x.Cosine() }, func(x *Decimal) *Decimal { return x.Cos() }},
		{func(x *Decimal) *Decimal { return x.Tangent() }, func(x *Decimal) *Decimal { return x.Tan() }},
		{func(x *Decimal) *Decimal { return x.HyperbolicSine() }, func(x *Decimal) *Decimal { return x.Sinh() }},
		{func(x *Decimal) *Decimal { return x.HyperbolicCosine() }, func(x *Decimal) *Decimal { return x.Cosh() }},
		{func(x *Decimal) *Decimal { return x.HyperbolicTangent() }, func(x *Decimal) *Decimal { return x.Tanh() }},
		{func(x *Decimal) *Decimal { return x.InverseSine() }, func(x *Decimal) *Decimal { return x.Asin() }},
		{func(x *Decimal) *Decimal { return x.InverseCosine() }, func(x *Decimal) *Decimal { return x.Acos() }},
		{func(x *Decimal) *Decimal { return x.InverseTangent() }, func(x *Decimal) *Decimal { return x.Atan() }},
		{func(x *Decimal) *Decimal { return x.InverseHyperbolicSine() }, func(x *Decimal) *Decimal { return x.Asinh() }},
		{func(x *Decimal) *Decimal { return x.InverseHyperbolicCosine() }, func(x *Decimal) *Decimal { return x.Acosh() }},
		{func(x *Decimal) *Decimal { return x.InverseHyperbolicTangent() }, func(x *Decimal) *Decimal { return x.Atanh() }},
	}
	for _, v := range transVals {
		x := New(v)
		for _, u := range unaries {
			al, pr := u.alias(x), u.primary(x)
			if al.ValueOf() != pr.ValueOf() || al.IsNaN() != pr.IsNaN() {
				t.Fatalf("%s: alias=%s primary=%s", v, al.ValueOf(), pr.ValueOf())
			}
		}
	}

	// Binary aliases on the full range (arithmetic only).
	for _, a := range fullVals {
		for _, b := range fullVals {
			x, y := New(a), New(b)
			if x.DividedBy(y).ValueOf() != x.Div(y).ValueOf() {
				t.Fatalf("DividedBy %s/%s", a, b)
			}
			if x.DividedToIntegerBy(y).ValueOf() != x.DivToInt(y).ValueOf() {
				t.Fatalf("DividedToIntegerBy %s/%s", a, b)
			}
			if x.Modulo(y).ValueOf() != x.Mod(y).ValueOf() {
				t.Fatalf("Modulo %s %s", a, b)
			}
			if x.ToPower(y).ValueOf() != x.Pow(y).ValueOf() {
				t.Fatalf("ToPower %s^%s", a, b)
			}
			if x.ComparedTo(y) != x.Cmp(y) {
				t.Fatalf("ComparedTo %s %s", a, b)
			}
			if x.Equals(y) != x.Eq(y) || x.GreaterThan(y) != x.Gt(y) ||
				x.GreaterThanOrEqualTo(y) != x.Gte(y) || x.LessThan(y) != x.Lt(y) ||
				x.LessThanOrEqualTo(y) != x.Lte(y) {
				t.Fatalf("comparisons %s %s", a, b)
			}
		}
	}

	// Clamp/ClampedTo with ordered bounds (min<=max).
	for _, a := range transVals {
		for _, lo := range transVals {
			for _, hi := range transVals {
				if New(lo).Gt(New(hi)) {
					continue
				}
				x := New(a)
				c1 := x.ClampedTo(lo, hi)
				c2 := x.Clamp(lo, hi)
				if c1.ValueOf() != c2.ValueOf() || c1.IsNaN() != c2.IsNaN() {
					t.Fatalf("ClampedTo clamp(%s,%s,%s) alias=%s primary=%s", a, lo, hi, c1.ValueOf(), c2.ValueOf())
				}
			}
		}
	}

	// Precision/DecimalPlaces/Sd/Dp.
	for _, v := range fullVals {
		x := New(v)
		if x.Precision() != x.Sd() {
			t.Fatalf("Precision %s: %v vs %v", v, x.Precision(), x.Sd())
		}
		if x.Precision(true) != x.Sd(true) {
			t.Fatalf("Precision(true) %s", v)
		}
		if x.DecimalPlaces() != x.Dp() {
			t.Fatalf("DecimalPlaces %s", v)
		}
		if x.IsInteger() != x.IsInt() || x.IsNegative() != x.IsNeg() || x.IsPositive() != x.IsPos() {
			t.Fatalf("is* %s", v)
		}
	}

	// Log/Logarithm (default and explicit bases).
	for _, base := range []string{"2", "10", "0.5", "3.14"} {
		for _, v := range transVals {
			x := New(v)
			l1 := x.Logarithm(base)
			l2 := x.Log(base)
			if l1.ValueOf() != l2.ValueOf() {
				t.Fatalf("Logarithm(%s) %s", base, v)
			}
		}
	}
	if x := New("2"); x.Logarithm().ValueOf() != x.Log().ValueOf() {
		t.Fatal("Logarithm default")
	}
}
