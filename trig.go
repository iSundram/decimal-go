package decimal

import (
	"math"
	"strconv"
)

// Exponent e must be positive and non-zero.
func tinyPow(b, e int64) float64 {
	n := float64(b)
	for e--; e > 0; e-- {
		n *= float64(b)
	}
	return n
}

// cos(x) = 1 - x^2/2! + x^4/4! - ...
// |x| < pi/2
func cosine(c *Constructor, x *Decimal) *Decimal {
	if x.IsZero() {
		return x
	}

	// Argument reduction: cos(4x) = 8*(cos^4(x) - cos^2(x)) + 1
	// i.e. cos(x) = 8*(cos^4(x/4) - cos^2(x/4)) + 1

	// Estimate the optimum number of times to use the argument reduction.
	var k int64
	var y string
	ln := len(x.d)
	if ln < 32 {
		k = int64(math.Ceil(float64(ln) / 3))
		y = strconv.FormatFloat(1/tinyPow(4, k), 'g', -1, 64)
	} else {
		k = 16
		y = "2.3283064365386962890625e-10"
	}

	c.Precision += k

	x = taylorSeries(c, 1, x.Times(c.New(y)), c.New(1), false)

	// Reverse argument reduction
	for i := k; i > 0; i-- {
		cos2x := x.Times(x)
		x = cos2x.Times(cos2x).Minus(cos2x).Times(8).Plus(1)
	}

	c.Precision -= k

	return x
}

// sin(x) = x - x^3/3! + x^5/5! - ...
// |x| < pi/2
func sine(c *Constructor, x *Decimal) *Decimal {
	ln := len(x.d)

	if ln < 3 {
		if x.IsZero() {
			return x
		}
		return taylorSeries(c, 2, x, x, false)
	}

	// Argument reduction: sin(5x) = 16*sin^5(x) - 20*sin^3(x) + 5*sin(x)
	// i.e. sin(x) = sin(x/5)(5 + sin^2(x/5)(16sin^2(x/5) - 20))

	// Estimate the optimum number of times to use the argument reduction.
	k := int64(1.4 * math.Sqrt(float64(ln)))
	if k > 16 {
		k = 16
	}

	x = x.Times(1 / tinyPow(5, k))
	x = taylorSeries(c, 2, x, x, false)

	// Reverse argument reduction
	d5 := c.New(5)
	d16 := c.New(16)
	d20 := c.New(20)
	for ; k > 0; k-- {
		sin2x := x.Times(x)
		x = x.Times(d5.Plus(sin2x.Times(d16.Times(sin2x).Minus(d20))))
	}

	return x
}

// Calculate Taylor series for cos, cosh, sin and sinh.
func taylorSeries(c *Constructor, n int64, x, y *Decimal, isHyperbolic bool) *Decimal {
	pr := c.Precision
	k := int(math.Ceil(float64(pr) / float64(logBase)))

	c.external = false
	x2 := x.Times(x)
	u := c.New(y)

	var t *Decimal
	var j int
	for {
		factor := n * (n + 1)
		n += 2
		t = divide(u.Times(x2), c.New(factor), pr, 1, false, 0)
		if isHyperbolic {
			u = y.Plus(t)
		} else {
			u = y.Minus(t)
		}
		factor = n * (n + 1)
		n += 2
		y = divide(t.Times(x2), c.New(factor), pr, 1, false, 0)
		t = u.Plus(y)

		if k < len(t.d) {
			j = k
			for {
				if j >= len(u.d) || t.d[j] != u.d[j] {
					break
				}
				if j == 0 {
					j = -1
					break
				}
				j--
			}
			if j == -1 {
				break
			}
		}

		u, y, t = y, t, u //lint:ignore SA4006 parity port rotation
	}

	c.external = true
	if len(t.d) > k+1 {
		t.d = t.d[:k+1]
	}

	return t
}

// Return the absolute value of x reduced to less than or equal to half pi.
func toLessThanHalfPi(c *Constructor, x *Decimal) *Decimal {
	isNeg := x.s < 0
	pi := getPi(c, c.Precision, 1)
	halfPi := pi.Times(0.5)

	x = x.Abs()

	if x.Lte(halfPi) {
		if isNeg {
			c.quadrant = 4
		} else {
			c.quadrant = 1
		}
		return x
	}

	t := x.DivToInt(pi)

	if t.IsZero() {
		if isNeg {
			c.quadrant = 3
		} else {
			c.quadrant = 2
		}
	} else {
		x = x.Minus(t.Times(pi))

		// 0 <= x < pi
		if x.Lte(halfPi) {
			q := 1
			if isOdd(t) {
				if isNeg {
					q = 2
				} else {
					q = 3
				}
			} else {
				if isNeg {
					q = 4
				} else {
					q = 1
				}
			}
			c.quadrant = q
			return x
		}

		if isOdd(t) {
			if isNeg {
				c.quadrant = 1
			} else {
				c.quadrant = 4
			}
		} else {
			if isNeg {
				c.quadrant = 3
			} else {
				c.quadrant = 2
			}
		}
	}

	return x.Minus(pi).Abs()
}

// Cos returns the cosine of x (in radians), rounded to the constructor's
// precision.
func (x *Decimal) Cos() *Decimal {
	c := x.c

	if x.d == nil {
		return c.newNaN()
	}

	// cos(0) = cos(-0) = 1
	if x.d[0] == 0 {
		return c.New(1)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + max64(x.e, x.sdInt()) + logBase
	c.Rounding = 1

	x = cosine(c, toLessThanHalfPi(c, x))

	c.Precision = pr
	c.Rounding = rm

	if x.c.quadrant == 2 || x.c.quadrant == 3 {
		x = x.Neg()
	}
	return finalise(x, pr, rm, true)
}

// Sin returns the sine of x (in radians), rounded to the constructor's
// precision.
func (x *Decimal) Sin() *Decimal {
	c := x.c

	if !x.IsFinite() {
		return c.newNaN()
	}
	if x.IsZero() {
		return c.New(x)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + max64(x.e, x.sdInt()) + logBase
	c.Rounding = 1

	x = sine(c, toLessThanHalfPi(c, x))

	c.Precision = pr
	c.Rounding = rm

	if x.c.quadrant > 2 {
		x = x.Neg()
	}
	return finalise(x, pr, rm, true)
}

// Tan returns the tangent of x (in radians), rounded to the constructor's
// precision.
func (x *Decimal) Tan() *Decimal {
	c := x.c

	if !x.IsFinite() {
		return c.newNaN()
	}
	if x.IsZero() {
		return c.New(x)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + 10
	c.Rounding = 1

	xx := x.Sin()
	xx.s = 1
	xx = divide(xx, c.New(1).Minus(xx.Times(xx)).Sqrt(), pr+10, 0, false, 0)

	c.Precision = pr
	c.Rounding = rm

	if x.c.quadrant == 2 || x.c.quadrant == 4 {
		xx = xx.Neg()
	}
	return finalise(xx, pr, rm, true)
}

// Cosh returns the hyperbolic cosine of x, rounded to the constructor's
// precision.
func (x *Decimal) Cosh() *Decimal {
	c := x.c
	one := c.New(1)

	if !x.IsFinite() {
		if x.s != 0 {
			return c.newInf(1)
		}
		return c.newNaN()
	}
	if x.IsZero() {
		return one
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + max64(x.e, x.sdInt()) + 4
	c.Rounding = 1
	ln := len(x.d)

	// Argument reduction: cos(4x) = 1 - 8cos^2(x) + 8cos^4(x) + 1
	var k int64
	var n string
	if ln < 32 {
		k = int64(math.Ceil(float64(ln) / 3))
		n = strconv.FormatFloat(1/tinyPow(4, k), 'g', -1, 64)
	} else {
		k = 16
		n = "2.3283064365386962890625e-10"
	}

	x = taylorSeries(c, 1, x.Times(c.New(n)), c.New(1), true)

	// Reverse argument reduction
	d8 := c.New(8)
	for i := k; i > 0; i-- {
		cosh2x := x.Times(x)
		x = one.Minus(cosh2x.Times(d8.Minus(cosh2x.Times(d8))))
	}

	c.Precision = pr
	c.Rounding = rm

	return finalise(x, pr, rm, true)
}

// Sinh returns the hyperbolic sine of x, rounded to the constructor's
// precision.
func (x *Decimal) Sinh() *Decimal {
	c := x.c

	if !x.IsFinite() || x.IsZero() {
		return c.New(x)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + max64(x.e, x.sdInt()) + 4
	c.Rounding = 1
	ln := len(x.d)

	if ln < 3 {
		x = taylorSeries(c, 2, x, x, true)
	} else {
		// Argument reduction: sinh(5x) = sinh(x)(5 + sinh^2(x)(20 + 16sinh^2(x)))
		// i.e. sinh(x) = sinh(x/5)(5 + sinh^2(x/5)(20 + 16sinh^2(x/5)))

		// Estimate the optimum number of times to use the argument reduction.
		k := int64(1.4 * math.Sqrt(float64(ln)))
		if k > 16 {
			k = 16
		}

		x = x.Times(1 / tinyPow(5, k))
		x = taylorSeries(c, 2, x, x, true)

		// Reverse argument reduction
		d5 := c.New(5)
		d16 := c.New(16)
		d20 := c.New(20)
		for ; k > 0; k-- {
			sinh2x := x.Times(x)
			x = x.Times(d5.Plus(sinh2x.Times(d16.Times(sinh2x).Plus(d20))))
		}
	}

	c.Precision = pr
	c.Rounding = rm

	return finalise(x, pr, rm, true)
}

// Tanh returns the hyperbolic tangent of x, rounded to the constructor's
// precision.
func (x *Decimal) Tanh() *Decimal {
	c := x.c

	if !x.IsFinite() {
		if x.s == 0 {
			return c.newNaN()
		}
		return c.New(int(x.s))
	}
	if x.IsZero() {
		return c.New(x)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + 7
	c.Rounding = 1

	x2 := x.Sinh().Div(x.Cosh())

	c.Precision = pr
	c.Rounding = rm

	return divide(x2, c.New(1), pr, rm, false, 0)
}

// Acos returns the arccosine in radians of x, rounded to the constructor's
// precision. Domain: [-1, 1]; Range: [0, pi].
func (x *Decimal) Acos() *Decimal {
	c := x.c
	k := x.Abs().Cmp(1)
	pr := c.Precision
	rm := c.Rounding

	if k != -1 {
		if k == 0 {
			// |x| is 1
			if x.IsNeg() {
				return getPi(c, pr, rm)
			}
			return c.New(0)
		}
		// |x| > 1 or x is NaN
		return c.newNaN()
	}

	if x.IsZero() {
		return getPi(c, pr+4, rm).Times(0.5)
	}

	c.Precision = pr + 6
	c.Rounding = 1

	x = c.New(1).Minus(x).Div(x.Plus(1)).Sqrt().Atan()

	c.Precision = pr
	c.Rounding = rm

	return x.Times(2)
}

// Acosh returns the inverse hyperbolic cosine of x, rounded to the
// constructor's precision.
func (x *Decimal) Acosh() *Decimal {
	c := x.c

	if x.Lte(1) {
		if x.Eq(1) {
			return c.New(0)
		}
		return c.newNaN()
	}
	if !x.IsFinite() {
		return c.New(x)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + max64(int64(math.Abs(float64(x.e))), x.sdInt()) + 4
	c.Rounding = 1
	x.c.external = false

	xx := x.Times(x).Minus(1).Sqrt().Plus(x)

	x.c.external = true
	c.Precision = pr
	c.Rounding = rm

	return xx.Ln()
}

// Asinh returns the inverse hyperbolic sine of x, rounded to the
// constructor's precision.
func (x *Decimal) Asinh() *Decimal {
	c := x.c

	if !x.IsFinite() || x.IsZero() {
		return c.New(x)
	}

	pr := c.Precision
	rm := c.Rounding
	c.Precision = pr + 2*max64(int64(math.Abs(float64(x.e))), x.sdInt()) + 6
	c.Rounding = 1
	x.c.external = false

	xx := x.Times(x).Plus(1).Sqrt().Plus(x)

	x.c.external = true
	c.Precision = pr
	c.Rounding = rm

	return xx.Ln()
}

// Atanh returns the inverse hyperbolic tangent of x, rounded to the
// constructor's precision.
func (x *Decimal) Atanh() *Decimal {
	c := x.c

	if !x.IsFinite() {
		return c.newNaN()
	}
	if x.e >= 0 {
		if x.Abs().Eq(1) {
			return c.newInf(x.s)
		}
		if x.IsZero() {
			return c.New(x)
		}
		return c.newNaN()
	}

	pr := c.Precision
	rm := c.Rounding
	xsd := x.sdInt()

	if max64(xsd, pr) < 2*-x.e-1 {
		return finalise(c.New(x), pr, rm, true)
	}

	wpr := xsd - x.e
	c.Precision = wpr

	xx := divide(x.Plus(1), c.New(1).Minus(x), wpr+pr, 1, false, 0)

	c.Precision = pr + 4
	c.Rounding = 1

	xx = xx.Ln()

	c.Precision = pr
	c.Rounding = rm

	return xx.Times(0.5)
}

// Asin returns the arcsine in radians of x, rounded to the constructor's
// precision. Domain: [-1, 1]; Range: [-pi/2, pi/2].
func (x *Decimal) Asin() *Decimal {
	c := x.c

	if x.IsZero() {
		return c.New(x)
	}

	k := x.Abs().Cmp(1)
	pr := c.Precision
	rm := c.Rounding

	if k != -1 {
		// |x| is 1
		if k == 0 {
			halfPi := getPi(c, pr+4, rm).Times(0.5)
			halfPi.s = x.s
			return halfPi
		}

		// |x| > 1 or x is NaN
		return c.newNaN()
	}

	c.Precision = pr + 6
	c.Rounding = 1

	x = x.Div(c.New(1).Minus(x).Times(c.New(1).Plus(x)).Sqrt().Plus(1)).Atan()

	c.Precision = pr
	c.Rounding = rm

	return x.Times(2)
}

// Atan returns the arctangent in radians of x, rounded to the constructor's
// precision. Range: [-pi/2, pi/2].
func (x *Decimal) Atan() *Decimal {
	c := x.c
	pr := c.Precision
	rm := c.Rounding

	if !x.IsFinite() {
		if x.s == 0 {
			return c.newNaN()
		}
		if pr+4 <= piPrecision {
			r := getPi(c, pr+4, rm).Times(0.5)
			r.s = x.s
			return r
		}
	} else if x.IsZero() {
		return c.New(x)
	} else if x.Abs().Eq(1) && pr+4 <= piPrecision {
		r := getPi(c, pr+4, rm).Times(0.25)
		r.s = x.s
		return r
	}

	wpr := pr + 10
	c.Precision = wpr
	c.Rounding = 1

	// Argument reduction: ensure |x| < 0.42
	// atan(x) = 2 * atan(x / (1 + sqrt(1 + x^2)))

	k := int64(float64(wpr)/float64(logBase) + 2)
	if 28 < k {
		k = 28
	}

	for i := k; i > 0; i-- {
		x = x.Div(x.Times(x).Plus(1).Sqrt().Plus(1))
	}

	x.c.external = false

	j := int(math.Ceil(float64(wpr) / float64(logBase)))
	n := int64(1)
	x2 := x.Times(x)
	r := c.New(x)
	px := x

	// atan(x) = x - x^3/3 + x^5/5 - x^7/7 + ...
	i := int64(0)
	for i != -1 {
		px = px.Times(x2)
		n += 2
		t := r.Minus(px.Div(n))

		px = px.Times(x2)
		n += 2
		r = t.Plus(px.Div(n))

		if j < len(r.d) {
			i = int64(j)
			for {
				if int(i) >= len(t.d) || r.d[i] != t.d[i] {
					break
				}
				if i == 0 {
					i = -1
					break
				}
				i--
			}
		}
	}

	if k != 0 {
		r = r.Times(int64(2 << (k - 1)))
	}

	x.c.external = true

	c.Precision = pr
	c.Rounding = rm

	return finalise(r, pr, rm, true)
}

// Atan2 returns the arctangent in radians of y/x in the range -pi to pi,
// rounded to the constructor's precision.
func (c *Constructor) Atan2(y, x any) *Decimal {
	yv := c.New(y)
	xv := c.New(x)
	pr := c.Precision
	rm := c.Rounding
	wpr := pr + 4

	var r *Decimal

	switch {
	// Either NaN.
	case yv.s == 0 || xv.s == 0:
		r = c.newNaN()

	// Both ±Infinity.
	case yv.d == nil && xv.d == nil:
		if xv.s > 0 {
			r = getPi(c, wpr, 1).Times(0.25)
		} else {
			r = getPi(c, wpr, 1).Times(0.75)
		}
		r.s = yv.s

	// x is ±Infinity or y is ±0.
	case xv.d == nil || yv.IsZero():
		if xv.s < 0 {
			r = getPi(c, pr, rm)
		} else {
			r = c.newZero(1)
		}
		r.s = yv.s

	// y is ±Infinity or x is ±0.
	case yv.d == nil || xv.IsZero():
		r = getPi(c, wpr, 1).Times(0.5)
		r.s = yv.s

	// Both non-zero and finite, x negative.
	case xv.s < 0:
		c.Precision = wpr
		c.Rounding = 1
		r = divide(yv, xv, wpr, 1, false, 0).Atan()
		pi := getPi(c, wpr, 1)
		c.Precision = pr
		c.Rounding = rm
		if yv.s < 0 {
			r = r.Minus(pi)
		} else {
			r = r.Plus(pi)
		}

	default:
		r = divide(yv, xv, wpr, 1, false, 0).Atan()
	}

	return r
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
