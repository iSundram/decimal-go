package decimal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// jsNumZero reports whether +s === 0 for a string of digits (or empty).
func jsNumZero(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= '1' && s[i] <= '9' {
			return false
		}
	}
	return true
}

func allNines(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '9' {
			return false
		}
	}
	return len(s) > 0
}

// sliceTo mimics JS str.slice(0, n) with tolerant bounds.
func sliceTo(s string, n int64) string {
	if n <= 0 {
		return ""
	}
	if int(n) >= len(s) {
		return s
	}
	return s[:n]
}

// sliceRange mimics JS str.slice(from, to) for from >= 0.
func sliceRange(s string, from, to int64) string {
	if from < 0 {
		from = 0
	}
	if int(from) >= len(s) {
		return ""
	}
	if int(to) > len(s) {
		to = int64(len(s))
	}
	if to < from {
		return ""
	}
	return s[from:to]
}

// powNumber is math.Pow with JS Math.pow semantics for |a| == 1, b = ±Inf.
func powNumber(a, b float64) float64 {
	if math.IsInf(b, 0) && math.Abs(a) == 1 && math.Signbit(a) {
		return math.NaN()
	}
	return math.Pow(a, b)
}

// formatSciToExp formats f like JS Number.prototype.toExponential() with no
// argument (shortest round-trip digits), then replaces the exponent with e.
func formatSciWithExp(f float64, e int64) string {
	ns := strconv.FormatFloat(f, 'e', -1, 64)
	idx := strings.IndexByte(ns, 'e')
	return ns[:idx+1] + strconv.FormatInt(e, 10)
}

func getLn10(c *Constructor, sd, pr int64) *Decimal {
	if sd > ln10Precision {
		// Reset global state in case the panic is caught.
		external = true
		if pr > 0 {
			c.Precision = pr
		}
		panic(errors.New(precisionLimitExceeded))
	}
	return finalise(c.New(ln10Num), sd, 1, true)
}

func getPi(c *Constructor, sd, rm int64) *Decimal {
	if sd > piPrecision {
		panic(errors.New(precisionLimitExceeded))
	}
	return finalise(c.New(piNum), sd, rm, true)
}

// Sqrt returns a new Decimal whose value is the square root of x, rounded
// to the constructor's precision.
func (x *Decimal) Sqrt() *Decimal {
	var m bool
	c := x.c
	d := x.d
	e := x.e
	s := x.s

	// Negative/NaN/Infinity/zero?
	if s != 1 || d == nil || d[0] == 0 {
		if s == 0 || (s < 0 && (d == nil || d[0] != 0)) {
			return c.newNaN()
		}
		if d != nil {
			return c.New(x)
		}
		return c.newInf(1)
	}

	external = false

	// Initial estimate.
	sf := math.Sqrt(x.Float64())

	var r *Decimal
	var n string

	// Sqrt underflow/overflow?
	// Pass x to Sqrt as integer, then adjust the exponent of the result.
	if sf == 0 || math.IsInf(sf, 0) {
		n = digitsToString(d)
		if (int64(len(n))+e)%2 == 0 {
			n += "0"
		}
		f, _ := strconv.ParseFloat(n, 64)
		sf = math.Sqrt(f)

		adj := floorDiv(e+1, 2)
		if e < 0 || e%2 != 0 {
			adj--
		}
		e = adj

		if math.IsInf(sf, 0) {
			n = "5e" + strconv.FormatInt(e, 10)
		} else {
			n = formatSciWithExp(sf, e)
		}

		r = c.New(n)
	} else {
		r = c.New(strconv.FormatFloat(sf, 'g', -1, 64))
	}

	sd := c.Precision
	e = sd
	sd += 3
	rep := false

	// Newton-Raphson iteration.
	for {
		t := r
		r = t.Plus(divide(x, t, sd+2, 1, false, 0)).Times(0.5)

		tsd := sliceTo(digitsToString(t.d), sd)
		ns := digitsToString(r.d)
		if tsd == sliceTo(ns, sd) {
			nn := sliceRange(ns, sd-3, sd+1)

			// The 4th rounding digit may be in error by -1 so if the 4
			// rounding digits are 9999 or 4999, i.e. approaching a rounding
			// boundary, continue the iteration.
			if nn == "9999" || !rep && nn == "4999" {
				// On the first iteration only, check to see if rounding up
				// gives the exact result as the nines may infinitely repeat.
				if !rep {
					finalise(t, e+1, 0, false)
					if t.Times(t).Eq(x) {
						r = t
						break
					}
				}

				sd += 4
				rep = true
			} else {
				// If the rounding digits are null, 0{0,4} or 50{0,3},
				// check for an exact result. If not, then there are further
				// digits and m will be true.
				if jsNumZero(nn) || (len(nn) > 0 && jsNumZero(nn[1:]) && nn[0] == '5') {
					// Truncate to the first rounding digit.
					finalise(r, e+1, 1, false)
					m = !r.Times(r).Eq(x)
				}

				break
			}
		}
	}

	external = true

	return finalise(r, e, c.Rounding, m)
}

// Cbrt returns a new Decimal whose value is the cube root of x, rounded to
// the constructor's precision.
func (x *Decimal) Cbrt() *Decimal {
	var m bool
	c := x.c

	if !x.IsFinite() || x.IsZero() {
		return c.New(x)
	}
	external = false

	// Initial estimate.
	s := float64(x.s) * math.Pow(float64(x.s)*x.Float64(), 1.0/3)

	var r *Decimal
	var n string
	var e int64

	// Math.cbrt underflow/overflow?
	// Pass x to math.Pow as integer, then adjust the exponent of the result.
	if s == 0 || math.IsNaN(s) || math.IsInf(s, 0) || math.Abs(s) == math.Inf(1) {
		n = digitsToString(x.d)
		e = x.e

		// Adjust n exponent so it is a multiple of 3 away from x exponent.
		if adj := (e - int64(len(n)) + 1) % 3; adj != 0 {
			if adj == 1 || adj == -2 {
				n += "0"
			} else {
				n += "00"
			}
		}
		f, _ := strconv.ParseFloat(n, 64)
		s = math.Pow(f, 1.0/3)

		// Rarely, e may be one less than the result exponent value.
		sub := int64(0)
		rem := e % 3
		exp := -1
		if e >= 0 {
			exp = 2
		}
		if rem == int64(exp) {
			sub = 1
		}
		e = floorDiv(e+1, 3) - sub

		if math.IsInf(s, 0) {
			n = "5e" + strconv.FormatInt(e, 10)
		} else {
			n = formatSciWithExp(s, e)
		}

		r = c.New(n)
		r.s = x.s
	} else {
		r = c.New(strconv.FormatFloat(s, 'g', -1, 64))
	}

	sd := c.Precision
	e = sd
	sd += 3
	rep := false

	// Halley's method.
	for {
		t := r
		t3 := t.Times(t).Times(t)
		t3plusx := t3.Plus(x)
		r = divide(t3plusx.Plus(x).Times(t), t3plusx.Plus(t3), sd+2, 1, false, 0)

		if sliceTo(digitsToString(t.d), sd) == sliceTo(digitsToString(r.d), sd) {
			nn := sliceRange(digitsToString(r.d), sd-3, sd+1)

			if nn == "9999" || !rep && nn == "4999" {
				if !rep {
					finalise(t, e+1, 0, false)
					if t.Times(t).Times(t).Eq(x) {
						r = t
						break
					}
				}

				sd += 4
				rep = true
			} else {
				if jsNumZero(nn) || (len(nn) > 0 && jsNumZero(nn[1:]) && nn[0] == '5') {
					// Truncate to the first rounding digit.
					finalise(r, e+1, 1, false)
					m = !r.Times(r).Times(r).Eq(x)
				}

				break
			}
		}
	}

	external = true

	return finalise(r, e, c.Rounding, m)
}

// Exp returns e^x rounded to the constructor's precision.
func (x *Decimal) Exp() *Decimal { return naturalExponential(x, prNull) }

// Ln returns the natural logarithm of x rounded to the constructor's
// precision.
func (x *Decimal) Ln() *Decimal { return naturalLogarithm(x, prNull) }

// Log returns the logarithm of x to the given base (default base 10),
// rounded to the constructor's precision.
func (x *Decimal) Log(bases ...any) *Decimal {
	arg := x
	c := x.c
	pr := c.Precision
	rm := c.Rounding
	guard := int64(5)

	var baseD *Decimal
	var isBase10 bool

	// Default base is 10.
	if len(bases) == 0 || bases[0] == nil {
		baseD = c.New(10)
		isBase10 = true
	} else {
		baseD = c.New(bases[0])
		d := baseD.d
		// Return NaN if base is negative, or non-finite, or is 0 or 1.
		if baseD.s < 0 || d == nil || d[0] == 0 || baseD.Eq(1) {
			return c.newNaN()
		}
		isBase10 = baseD.Eq(10)
	}

	d := arg.d

	// Is arg negative, non-finite, 0 or 1?
	if arg.s < 0 || d == nil || d[0] == 0 || arg.Eq(1) {
		return c.newLogSpecial(arg, d)
	}

	// The result will have a non-terminating decimal expansion if base is
	// 10 and arg is not an integer power of 10.
	inf := false
	if isBase10 {
		if len(d) > 1 {
			inf = true
		} else {
			k := d[0]
			for k%10 == 0 {
				k /= 10
			}
			inf = k != 1
		}
	}

	external = false
	sd := pr + guard
	num := naturalLogarithm(arg, sd)
	var denominator *Decimal
	if isBase10 {
		denominator = getLn10(c, sd+10, 0)
	} else {
		denominator = naturalLogarithm(baseD, sd)
	}

	// The result will have 5 rounding digits.
	r := divide(num, denominator, sd, 1, false, 0)

	// If at a rounding boundary, i.e. the result's rounding digits are
	// [49]9999 or [50]0000, calculate 10 further digits.
	if checkRoundingDigits(r.d, pr, rm, -1) {
		k := pr
		for {
			sd += 10
			num = naturalLogarithm(arg, sd)
			if isBase10 {
				denominator = getLn10(c, sd+10, 0)
			} else {
				denominator = naturalLogarithm(baseD, sd)
			}
			r = divide(num, denominator, sd, 1, false, 0)

			if !inf {
				// Check for 14 nines from the 2nd rounding digit, as the
				// first may be 4.
				if nn := sliceRange(digitsToString(r.d), k+1, k+15); len(nn) == 14 && allNines(nn) {
					r = finalise(r, pr+1, 0, false)
				}
				break
			}
			k += 10
			if !checkRoundingDigits(r.d, k, rm, -1) {
				break
			}
		}
	}

	external = true

	return finalise(r, pr, rm, false)
}

func (c *Constructor) newLogSpecial(arg *Decimal, d []int32) *Decimal {
	if d != nil && d[0] == 0 {
		return c.newInf(-1)
	}
	if arg.s != 1 {
		return c.newNaN()
	}
	if d != nil {
		return c.newZero(1)
	}
	return c.newInf(1)
}

// Return a new Decimal whose value is x^y.
func naturalExponential(x *Decimal, sd int64) *Decimal {
	c := x.c
	rm := c.Rounding
	pr := c.Precision
	rep := int64(0)
	i := 0
	k := 0

	// 0/NaN/Infinity?
	if x.d == nil || x.d[0] == 0 || x.e > 17 {
		switch {
		case x.d != nil:
			switch {
			case x.d[0] == 0:
				return c.New(1)
			case x.s < 0:
				return c.newZero(1)
			default:
				return c.newInf(1)
			}
		case x.s != 0:
			if x.s < 0 {
				return c.newZero(1)
			}
			return c.New(x)
		default:
			return c.newNaN()
		}
	}

	var wpr int64
	if sd == prNull {
		external = false
		wpr = pr
	} else {
		wpr = sd
	}

	t := c.New(0.03125)

	// while abs(x) >= 0.1
	for x.e > -2 {
		x = x.Times(t)
		k += 5
	}

	// Estimate the increase in precision necessary.
	guard := int64(math.Log(math.Pow(2, float64(k)))/math.Ln10*2 + 5)
	wpr += guard
	denominator := c.New(1)
	pow := denominator
	sum := c.New(1)
	c.Precision = wpr

	for {
		pow = finalise(pow.Times(x), wpr, 1, false)
		i++
		denominator = denominator.Times(i)
		t = sum.Plus(divide(pow, denominator, wpr, 1, false, 0))

		if sliceTo(digitsToString(t.d), wpr) == sliceTo(digitsToString(sum.d), wpr) {
			for j := k; j > 0; j-- {
				sum = finalise(sum.Times(sum), wpr, 1, false)
			}

			// Check to see if the first 4 rounding digits are [49]999.
			if sd == prNull {
				if rep < 3 && checkRoundingDigits(sum.d, wpr-guard, rm, rep) {
					wpr += 10
					c.Precision = wpr
					t = c.New(1)
					denominator = c.New(1)
					pow = c.New(1)
					i = 0
					rep++
				} else {
					c.Precision = pr
					external = true
					return finalise(sum, pr, rm, true)
				}
			} else {
				c.Precision = pr
				return sum
			}
		}

		sum = t
	}
}

// Return the natural logarithm of y rounded to sd significant digits.
func naturalLogarithm(y *Decimal, sd int64) *Decimal {
	n := int64(1)
	guard := int64(10)
	x := y
	xd := x.d
	c := x.c
	rm := c.Rounding
	pr := c.Precision
	rep := int64(0)

	// Is x negative or Infinity, NaN, 0 or 1?
	if x.s < 0 || xd == nil || xd[0] == 0 || (x.e == 0 && xd[0] == 1 && len(xd) == 1) {
		if xd != nil && xd[0] == 0 {
			return c.newInf(-1)
		}
		if x.s != 1 {
			return c.newNaN()
		}
		if xd != nil {
			return c.newZero(1)
		}
		return c.New(x)
	}

	var wpr int64
	if sd == prNull {
		external = false
		wpr = pr
	} else {
		wpr = sd
	}

	wpr += guard
	c.Precision = wpr
	cc := digitsToString(xd)
	c0 := cc[0]
	var e int64
	var t *Decimal

	if e = x.e; math.Abs(float64(e)) < 1.5e15 {
		// Argument reduction.
		for (c0 < '7' && c0 != '1') || (c0 == '1' && len(cc) > 1 && cc[1] > '3') {
			x = x.Times(y)
			cc = digitsToString(x.d)
			c0 = cc[0]
			n++
		}

		e = x.e

		if c0 > '1' {
			x = c.New("0." + cc)
			e++
		} else {
			x = c.New(string(c0) + "." + cc[1:])
		}
	} else {
		// The argument reduction method above may result in overflow if the
		// argument y is a massive number with exponent >= 1500000000000000,
		// so instead recall this function using ln(x*10^e) = ln(x) + e*ln(10).
		t = getLn10(c, wpr+2, pr).Times(strconv.FormatInt(e, 10))
		x = naturalLogarithm(c.New(string(c0)+"."+cc[1:]), wpr-guard).Plus(t)
		c.Precision = pr

		if sd == prNull {
			external = true
			return finalise(x, pr, rm, true)
		}
		return x
	}

	// x1 is x reduced to a value near 1.
	x1 := x

	// Taylor series.
	// ln(y) = ln((1 + x)/(1 - x)) = 2(x + x^3/3 + x^5/5 + x^7/7 + ...)
	// where x = (y - 1)/(y + 1)    (|x| < 1)
	sum := divide(x.Minus(1), x.Plus(1), wpr, 1, false, 0)
	numerator := sum
	x = sum
	x2 := finalise(x.Times(x), wpr, 1, false)
	denominator := int64(3)

	for {
		numerator = finalise(numerator.Times(x2), wpr, 1, false)
		t = sum.Plus(divide(numerator, c.New(denominator), wpr, 1, false, 0))

		if sliceTo(digitsToString(t.d), wpr) == sliceTo(digitsToString(sum.d), wpr) {
			sum = sum.Times(2)

			// Reverse the argument reduction.
			if e != 0 {
				sum = sum.Plus(getLn10(c, wpr+2, pr).Times(strconv.FormatInt(e, 10)))
			}
			sum = divide(sum, c.New(n), wpr, 1, false, 0)

			// If rm > 3 and the first 4 rounding digits 4999, or rm < 4
			// (or the summation has been repeated previously) and the first
			// 4 rounding digits 9999, restart with a higher precision.
			if sd == prNull {
				if checkRoundingDigits(sum.d, wpr-guard, rm, rep) {
					wpr += guard
					c.Precision = wpr
					x = divide(x1.Minus(1), x1.Plus(1), wpr, 1, false, 0)
					numerator = x
					t = x
					x2 = finalise(x.Times(x), wpr, 1, false)
					denominator = 1
					rep = 1
				} else {
					c.Precision = pr
					external = true
					return finalise(sum, pr, rm, true)
				}
			} else {
				c.Precision = pr
				return sum
			}
		}

		sum = t
		denominator += 2
	}
}

// Pow returns a new Decimal whose value is x raised to the power y, rounded
// to the constructor's precision.
func (x *Decimal) Pow(y any) *Decimal {
	c := x.c
	yv := c.New(y)
	yn := yv.Float64()

	// Either ±Infinity, NaN or ±0?
	if x.d == nil || yv.d == nil || x.d[0] == 0 || yv.d[0] == 0 {
		return c.New(powNumber(x.Float64(), yn))
	}

	xx := c.New(x)

	if xx.Eq(1) {
		return xx
	}

	pr := c.Precision
	rm := c.Rounding

	if yv.Eq(1) {
		return finalise(xx, pr, rm, false)
	}

	// y exponent
	e := floorDiv(yv.e, logBase)

	// If y is a small integer use the 'exponentiation by squaring' algorithm.
	if e >= int64(len(yv.d))-1 {
		k := yn
		if k < 0 {
			k = -k
		}
		if !math.IsInf(k, 0) && !math.IsNaN(k) && k <= maxInteger {
			r := intPow(c, xx, int64(k), pr)
			if yv.s < 0 {
				return c.New(1).Div(r)
			}
			return finalise(r, pr, rm, false)
		}
	}

	s := xx.s

	// If x is negative.
	if s < 0 {
		// If y is not an integer.
		if e < int64(len(yv.d))-1 {
			return c.newNaN()
		}

		// Result is positive if x is negative and the last digit of integer
		// y is even.
		if yv.d[e]&1 == 0 {
			s = 1
		}

		// If x = -1.
		if xx.e == 0 && xx.d[0] == 1 && len(xx.d) == 1 {
			xx.s = s
			return xx
		}
	}

	// Estimate result exponent.
	// x^y = 10^e,  where e = y * log10(x)
	kF := powNumber(xx.Float64(), yn)
	var ee int64
	if kF == 0 || math.IsInf(kF, 0) || math.IsNaN(kF) {
		f, _ := strconv.ParseFloat("0."+digitsToString(xx.d), 64)
		ee = int64(math.Floor(yn * (math.Log(f)/math.Ln10 + float64(xx.e) + 1)))
	} else {
		ee = c.New(strconv.FormatFloat(kF, 'g', -1, 64)).e
	}

	// Overflow/underflow?
	if ee > c.MaxE+1 || ee < c.MinE-1 {
		if ee > 0 {
			return c.newInf(s)
		}
		return c.newZero(1)
	}

	external = false
	c.Rounding = 1
	xx.s = 1

	// Estimate the extra guard digits needed to ensure five correct
	// rounding digits from naturalLogarithm(x).
	kk := int64(len(strconv.FormatInt(ee, 10)))
	if kk > 12 {
		kk = 12
	}

	// r = x^y = exp(y*ln(x))
	r := naturalExponential(yv.Times(naturalLogarithm(xx, pr+kk)), pr)

	// r may be Infinity, e.g. (0.9999999999999999).pow(-1e+40)
	if r.d != nil {
		// Truncate to the required precision plus five rounding digits.
		r = finalise(r, pr+5, 1, false)

		// If the rounding digits are [49]9999 or [50]0000 increase the
		// precision by 10 and recalculate the result.
		if checkRoundingDigits(r.d, pr, rm, -1) {
			ee = pr + 10

			// Truncate to the increased precision plus five rounding digits.
			r = finalise(naturalExponential(yv.Times(naturalLogarithm(xx, ee+kk)), ee), ee+5, 1, false)

			// Check for 14 nines from the 2nd rounding digit (the first
			// rounding digit may be 4 or 9).
			if nn := sliceRange(digitsToString(r.d), pr+1, pr+15); len(nn) == 14 && allNines(nn) {
				r = finalise(r, pr+1, 0, false)
			}
		}
	}

	r.s = s
	external = true
	c.Rounding = rm

	return finalise(r, pr, rm, false)
}
