package decimal

// Tests immutability of operand[s] for all applicable methods.
// Also tests each Decimal method against its equivalent Constructor method
// where applicable.

import (
	"math"
	"math/rand"
	"testing"
)

func TestImmutability(t *testing.T) {
	setCfg(20, 4, -7, 21, 9e15, -9e15, 1, false)

	// Integer [0, 9e15), with each possible number of digits, [1, 16], equally likely.
	randInt := func() float64 {
		return math.Floor(rand.Float64() * 9e15 / math.Pow(10, float64(int64(rand.Float64()*16))))
	}

	teq := func(x, y *Decimal) {
		t.Helper()
		assertEqDecimal(t, x, y)
	}

	v := make([]any, 14)
	v[0] = 0.0
	v[1] = nan()
	v[2] = inf(1)
	v[3] = inf(-1)
	v[4] = 0.5
	v[5] = -0.5
	v[6] = 1.0
	v[7] = -1.0
	{
		x := Default.Random()
		v[8] = x
		v[9] = x.Neg()
		xi := randInt()
		v[10] = xi
		v[11] = -xi
		x2 := Default.Random().Plus(randInt())
		v[12] = x2
		v[13] = x2.Neg()
	}

	for i := 0; i < len(v); i++ {
		a := New(v[i])
		aa := New(v[i])
		k := float64(int64(rand.Float64()*10)) / 10
		var b *Decimal
		if k == 0.5 {
			b = New(a)
		} else if k < 0.5 {
			b = a.Plus(Default.Random().Plus(randInt()))
		} else {
			b = a.Minus(Default.Random().Plus(randInt()))
		}
		bb := New(b)
		n := int64(rand.Float64()*20 + 1)

		var dx, dy, dz *Decimal
		var fx, fy float64
		var bx, by bool
		var sx, sy string

		dx = a.Abs()
		teq(a, aa)
		dy = a.Abs()
		teq(a, aa)
		dz = Default.Abs(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Ceil()
		teq(a, aa)
		dy = Default.Ceil(a)
		teq(a, aa)

		teq(dx, dy)

		fx = a.Cmp(b)
		teq(a, aa)
		teq(b, bb)
		fy = a.Cmp(b)
		teq(a, aa)
		teq(b, bb)

		assertEq(t, fx, fy)

		dx = a.Cos()
		teq(a, aa)
		dy = a.Cos()
		teq(a, aa)
		dz = Default.Cos(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Cbrt()
		teq(a, aa)
		dy = a.Cbrt()
		teq(a, aa)
		dz = Default.Cbrt(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		fx = a.Dp()
		teq(a, aa)
		fy = a.Dp()
		teq(a, aa)

		assertEq(t, fx, fy)

		dx = a.Div(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.Div(b)
		teq(a, aa)
		teq(b, bb)
		dz = Default.Div(a, b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.DivToInt(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.DivToInt(b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)

		bx = a.Eq(b)
		teq(a, aa)
		teq(b, bb)
		by = a.Eq(b)
		teq(a, aa)
		teq(b, bb)

		assert(t, bx == by)

		dx = a.Floor()
		teq(a, aa)
		dy = Default.Floor(a)
		teq(a, aa)

		teq(dx, dy)

		bx = a.Gt(b)
		teq(a, aa)
		teq(b, bb)
		by = a.Gt(b)
		teq(a, aa)
		teq(b, bb)

		assert(t, bx == by)

		bx = a.Gte(b)
		teq(a, aa)
		teq(b, bb)
		by = a.Gte(b)
		teq(a, aa)
		teq(b, bb)

		assert(t, bx == by)

		// Omit hyperbolic methods if a is large, as they are too time-consuming.
		if a.Abs().Lt(1000) {
			dx = a.Cosh()
			teq(a, aa)
			dy = a.Cosh()
			teq(a, aa)
			dz = Default.Cosh(a)
			teq(a, aa)

			teq(dx, dy)
			teq(dy, dz)

			dx = a.Sinh()
			teq(a, aa)
			dy = a.Sinh()
			teq(a, aa)
			dz = Default.Sinh(a)
			teq(a, aa)

			teq(dx, dy)
			teq(dy, dz)

			dx = a.Tanh()
			teq(a, aa)
			dy = a.Tanh()
			teq(a, aa)
			dz = Default.Tanh(a)
			teq(a, aa)

			teq(dx, dy)
			teq(dy, dz)
		}

		// a [-1, 1]
		dx = a.Acos()
		teq(a, aa)
		dy = a.Acos()
		teq(a, aa)
		dz = Default.Acos(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		// a [1, Infinity]
		dx = a.Acosh()
		teq(a, aa)
		dy = a.Acosh()
		teq(a, aa)
		dz = Default.Acosh(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Asinh()
		teq(a, aa)
		dy = a.Asinh()
		teq(a, aa)
		dz = Default.Asinh(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		// a [-1, 1]
		dx = a.Atanh()
		teq(a, aa)
		dy = a.Atanh()
		teq(a, aa)
		dz = Default.Atanh(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		// a [-1, 1]
		dx = a.Asin()
		teq(a, aa)
		dy = a.Asin()
		teq(a, aa)
		dz = Default.Asin(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Atan()
		teq(a, aa)
		dy = a.Atan()
		teq(a, aa)
		dz = Default.Atan(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		a.IsFinite()
		teq(a, aa)

		bx = a.IsInt()
		teq(a, aa)
		by = a.IsInt()
		teq(a, aa)

		assert(t, bx == by)

		a.IsNaN()
		teq(a, aa)

		bx = a.IsNeg()
		teq(a, aa)
		by = a.IsNeg()
		teq(a, aa)

		assert(t, bx == by)

		bx = a.IsPos()
		teq(a, aa)
		by = a.IsPos()
		teq(a, aa)

		assert(t, bx == by)

		a.IsZero()
		teq(a, aa)

		bx = a.Lt(b)
		teq(a, aa)
		teq(b, bb)
		by = a.Lt(b)
		teq(a, aa)
		teq(b, bb)

		assert(t, bx == by)

		bx = a.Lte(b)
		teq(a, aa)
		teq(b, bb)
		by = a.Lte(b)
		teq(a, aa)
		teq(b, bb)

		assert(t, bx == by)

		dx = a.Log()
		teq(a, aa)
		dy = a.Log()
		teq(a, aa)
		dz = Default.Log(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Minus(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.Sub(b)
		teq(a, aa)
		teq(b, bb)
		dz = Default.Sub(a, b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Mod(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.Mod(b)
		teq(a, aa)
		teq(b, bb)
		dz = Default.Mod(a, b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Exp()
		teq(a, aa)
		dy = a.Exp()
		teq(a, aa)
		dz = Default.Exp(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Ln()
		teq(a, aa)
		dy = a.Ln()
		teq(a, aa)
		dz = Default.Ln(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Neg()
		teq(a, aa)
		dy = a.Neg()
		teq(a, aa)

		teq(dx, dy)

		dx = a.Plus(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.Add(b)
		teq(a, aa)
		teq(b, bb)
		dz = Default.Add(a, b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)
		teq(dy, dz)

		fx = a.Sd()
		teq(a, aa)
		fy = a.Sd()
		teq(a, aa)

		assertEq(t, fx, fy)

		dx = a.Round()
		teq(a, aa)
		dy = Default.Round(a)
		teq(a, aa)

		teq(dx, dy)

		dx = a.Sin()
		teq(a, aa)
		dy = a.Sin()
		teq(a, aa)
		dz = Default.Sin(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Sqrt()
		teq(a, aa)
		dy = a.Sqrt()
		teq(a, aa)
		dz = Default.Sqrt(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Tan()
		teq(a, aa)
		dy = a.Tan()
		teq(a, aa)
		dz = Default.Tan(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		dx = a.Times(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.Mul(b)
		teq(a, aa)
		teq(b, bb)
		dz = Default.Mul(a, b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)
		teq(dy, dz)

		a.ToBinary()
		teq(a, aa)

		dx = a.ToDP(n)
		teq(a, aa)
		dy = a.ToDP(n)
		teq(a, aa)

		teq(dx, dy)

		a.ToExponential(n)
		teq(a, aa)

		a.ToFixed(n)
		teq(a, aa)

		a.ToFraction()
		teq(a, aa)

		sx = a.ToHex()
		teq(a, aa)
		sy = a.ToHex()
		teq(a, aa)

		assertEq(t, sx, sy)

		a.ValueOf() // toJSON
		teq(a, aa)

		a.ToNearest(b)
		teq(a, aa)
		teq(b, bb)

		a.Float64()
		teq(a, aa)

		a.ToOctal()
		teq(a, aa)

		dx = a.Pow(b)
		teq(a, aa)
		teq(b, bb)
		dy = a.Pow(b)
		teq(a, aa)
		teq(b, bb)
		dz = Default.Pow(a, b)
		teq(a, aa)
		teq(b, bb)

		teq(dx, dy)
		teq(dy, dz)

		a.ToPrecision(n)
		teq(a, aa)

		dx = a.ToSD(n)
		teq(a, aa)
		dy = a.ToSD(n)
		teq(a, aa)

		teq(dx, dy)

		_ = a.String()
		teq(a, aa)

		dx = a.Trunc()
		teq(a, aa)
		dy = a.Trunc()
		teq(a, aa)
		dz = Default.Trunc(a)
		teq(a, aa)

		teq(dx, dy)
		teq(dy, dz)

		a.ValueOf()
		teq(a, aa)

		Default.Atan2(a, b)
		teq(a, aa)
		teq(b, bb)

		Default.Hypot(a, b)
		teq(a, aa)
		teq(b, bb)

		dx = Default.Log(a, 10)
		teq(a, aa)
		dy = Default.Log10(a)
		teq(a, aa)

		teq(dx, dy)

		dx = Default.Log(a, 2)
		teq(a, aa)
		dy = Default.Log2(a)
		teq(a, aa)

		teq(dx, dy)

		Default.Max(a, b)
		teq(a, aa)
		teq(b, bb)

		Default.Min(a, b)
		teq(a, aa)
		teq(b, bb)

		Default.Sign(a)
		teq(a, aa)
	}
}
