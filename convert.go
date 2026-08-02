package decimal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// ±Infinity, NaN (without sign).
func nonFiniteToString(x *Decimal) string {
	if x.s == 0 {
		return "NaN"
	}
	return "Infinity"
}

func finiteToString(x *Decimal, isExp bool, sd int64) string {
	if !x.IsFinite() {
		return nonFiniteToString(x)
	}
	e := x.e
	str := digitsToString(x.d)
	ln := int64(len(str))
	var k int64

	if isExp {
		if sd != 0 && ln != 0 && sd > ln {
			k = sd - ln
			str = str[:1] + "." + str[1:] + zeroString(k)
		} else if ln > 1 {
			str = str[:1] + "." + str[1:]
		}
		if x.e < 0 {
			str += "e" + strconv.FormatInt(x.e, 10)
		} else {
			str += "e+" + strconv.FormatInt(x.e, 10)
		}
	} else if e < 0 {
		str = "0." + zeroString(-e-1) + str
		if sd != 0 && ln != 0 && sd > ln {
			str += zeroString(sd - ln)
		}
	} else if e >= ln {
		str += zeroString(e + 1 - ln)
		if sd != 0 && ln != 0 && sd > e+1 {
			str = str + "." + zeroString(sd-e-1)
		}
	} else {
		if e+1 < ln {
			str = str[:e+1] + "." + str[e+1:]
		}
		if sd != 0 && ln != 0 && sd > ln {
			if e+1 == ln {
				str += "."
			}
			str += zeroString(sd - ln)
		}
	}

	return str
}

// finiteToString0 is finiteToString without the sd argument (sd omitted).
func finiteToString0(x *Decimal, isExp bool) string {
	if !x.IsFinite() {
		return nonFiniteToString(x)
	}
	e := x.e
	str := digitsToString(x.d)
	ln := int64(len(str))

	if isExp {
		if ln > 1 {
			str = str[:1] + "." + str[1:]
		}
		if x.e < 0 {
			str += "e" + strconv.FormatInt(x.e, 10)
		} else {
			str += "e+" + strconv.FormatInt(x.e, 10)
		}
	} else if e < 0 {
		str = "0." + zeroString(-e-1) + str
	} else if e >= ln {
		str += zeroString(e + 1 - ln)
	} else if e+1 < ln {
		str = str[:e+1] + "." + str[e+1:]
	}

	return str
}

// String returns a string representing the value of x, using exponential
// notation if the exponent is >= ToExpPos or <= ToExpNeg.
func (x *Decimal) String() string {
	c := x.c
	str := finiteToString0(x, x.e <= c.ToExpNeg || x.e >= c.ToExpPos)
	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ValueOf is like String, but negative zero includes the minus sign.
func (x *Decimal) ValueOf() string {
	c := x.c
	str := finiteToString0(x, x.e <= c.ToExpNeg || x.e >= c.ToExpPos)
	if x.IsNeg() {
		return "-" + str
	}
	return str
}

// Float64 returns the value of x converted to a float64.
// Zero keeps its sign. NaN converts to NaN, ±Infinity to ±Inf.
func (x *Decimal) Float64() float64 {
	if x.d == nil {
		if x.s == 0 {
			return math.NaN()
		}
		return math.Inf(int(x.s))
	}
	if x.d[0] == 0 {
		if x.s < 0 {
			return math.Copysign(0, -1)
		}
		return 0
	}
	f, _ := strconv.ParseFloat(x.ValueOf(), 64)
	return f
}

// parseOptArgs parses optional (dp|sd, rm) argument pairs.
func parseOptArgs(c *Constructor, args []int64) (a, rm int64, hasA bool) {
	if len(args) > 0 {
		a = args[0]
		hasA = true
	}
	if len(args) > 1 {
		rm = args[1]
	} else {
		rm = c.Rounding
	}
	return
}

// ToExponential returns a string representing x in exponential notation
// rounded to dp fixed decimal places using rounding mode rm (or the
// constructor's rounding mode if omitted).
func (x *Decimal) ToExponential(args ...int64) string {
	c := x.c
	var str string
	var xx *Decimal
	dp, rm, hasDp := parseOptArgs(c, args)

	if !hasDp {
		xx = x
		str = finiteToString0(x, true)
	} else {
		checkInt32(dp, 0, maxDigits)
		checkInt32(rm, 0, 8)
		xx = finalise(c.New(x), dp+1, rm, false)
		str = finiteToString(xx, true, dp+1)
	}

	if xx.IsNeg() && !xx.IsZero() {
		return "-" + str
	}
	return str
}

// ToFixed returns a string representing x in normal (fixed-point) notation
// to dp fixed decimal places, rounded using rm (or the constructor's
// rounding mode if omitted).
func (x *Decimal) ToFixed(args ...int64) string {
	c := x.c
	var str string
	dp, rm, hasDp := parseOptArgs(c, args)

	if !hasDp {
		str = finiteToString0(x, false)
	} else {
		checkInt32(dp, 0, maxDigits)
		checkInt32(rm, 0, 8)
		y := finalise(c.New(x), dp+x.e+1, rm, false)
		str = finiteToString(y, false, dp+y.e+1)
	}

	// To determine whether to add the minus sign look at the value before
	// it was rounded, i.e. look at x rather than y.
	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToPrecision returns a string representing x rounded to sd significant
// digits. Exponential notation is used if necessary.
func (x *Decimal) ToPrecision(args ...int64) string {
	c := x.c
	var str string
	var xx *Decimal
	sd, rm, hasSd := parseOptArgs(c, args)

	if !hasSd {
		xx = x
		str = finiteToString0(x, x.e <= c.ToExpNeg || x.e >= c.ToExpPos)
	} else {
		checkInt32(sd, 1, maxDigits)
		checkInt32(rm, 0, 8)
		xx = finalise(c.New(x), sd, rm, false)
		str = finiteToString(xx, sd <= xx.e || xx.e <= c.ToExpNeg, sd)
	}

	if xx.IsNeg() && !xx.IsZero() {
		return "-" + str
	}
	return str
}

// ToDP returns a new Decimal whose value is x rounded to a maximum of dp
// decimal places using rounding mode rm (or the constructor's rounding mode
// if omitted).
func (x *Decimal) ToDP(args ...int64) *Decimal {
	c := x.c
	xx := c.New(x)
	if len(args) == 0 {
		return xx
	}
	dp, rm, _ := parseOptArgs(c, args)
	checkInt32(dp, 0, maxDigits)
	checkInt32(rm, 0, 8)
	return finalise(xx, dp+xx.e+1, rm, false)
}

// ToSD returns a new Decimal whose value is x rounded to a maximum of sd
// significant digits using rounding mode rm (or the constructor's
// precision and rounding mode if omitted).
func (x *Decimal) ToSD(args ...int64) *Decimal {
	c := x.c
	sd, rm, hasSd := parseOptArgs(c, args)
	if !hasSd {
		sd = c.Precision
	} else {
		checkInt32(sd, 1, maxDigits)
	}
	checkInt32(rm, 0, 8)
	return finalise(c.New(x), sd, rm, false)
}

// ToBinary returns a string representing x in base 2, rounded to sd
// significant digits using rm. If sd is present the result uses binary
// exponential notation, otherwise fixed-point.
func (x *Decimal) ToBinary(sd ...int64) string { return toStringBinary(x, 2, sd) }

// ToHex returns a string representing x in base 16. See ToBinary.
func (x *Decimal) ToHex(sd ...int64) string { return toStringBinary(x, 16, sd) }

// ToOctal returns a string representing x in base 8. See ToBinary.
func (x *Decimal) ToOctal(sd ...int64) string { return toStringBinary(x, 8, sd) }

// Return the value of Decimal x as a string in base baseOut.
func toStringBinary(x *Decimal, baseOut int64, opts []int64) string {
	c := x.c
	isExp := len(opts) > 0
	var sd, rm int64

	if isExp {
		sd = opts[0]
		checkInt32(sd, 1, maxDigits)
		if len(opts) > 1 {
			rm = opts[1]
			checkInt32(rm, 0, 8)
		} else {
			rm = c.Rounding
		}
	} else {
		sd = c.Precision
		rm = c.Rounding
	}

	var str string
	var roundUp bool
	var e int64

	if !x.IsFinite() {
		str = nonFiniteToString(x)
	} else {
		str = finiteToString0(x, false)
		i := int64(strings.IndexByte(str, '.'))

		var base int64
		if isExp {
			base = 2
			if baseOut == 16 {
				sd = sd*4 - 3
			} else if baseOut == 8 {
				sd = sd*3 - 2
			}
		} else {
			base = baseOut
		}

		var y *Decimal
		if i >= 0 {
			// Non-integer.
			str = strings.Replace(str, ".", "", 1)
			y = c.New(1)
			y.e = int64(len(str)) - i
			y.d = convertBase(finiteToString0(y, false), 10, base)
			y.e = int64(len(y.d))
		}

		xd := convertBase(str, 10, base)
		e = int64(len(xd))
		ln := e

		// Remove trailing zeros.
		for ; ln > 0 && xd[ln-1] == 0; ln-- {
			xd = xd[:ln-1]
		}

		if len(xd) == 0 || xd[0] == 0 {
			if isExp {
				str = "0p+0"
			} else {
				str = "0"
			}
		} else {
			if i < 0 {
				e--
			} else {
				xx := c.New(x)
				xx.d = xd
				xx.e = e
				xx = divide(xx, y, sd, rm, false, base)
				xd = xx.d
				e = xx.e
				roundUp = x.c.inexact
			}

			// The rounding digit.
			var rDigit int64
			rDigitDefined := false
			if int(sd) < len(xd) {
				rDigit = int64(xd[sd])
				rDigitDefined = true
			}
			k := base / 2
			roundUp = roundUp || int(sd)+1 < len(xd)

			if rm < 4 {
				roundUp = (rDigitDefined || roundUp) && (rm == 0 || rm == ifelse64(x.s < 0, 3, 2))
			} else {
				roundUp = rDigitDefined && (rDigit > k || rDigit == k && (rm == 4 || roundUp ||
					rm == 6 && int64(xd[sd-1])&1 == 1 || rm == ifelse64(x.s < 0, 8, 7)))
			}

			// Truncate to sd significant digits (JS would extend the array
			// with holes which read as undefined; equivalent to leaving the
			// slice as-is since the trailing-zero scan skips them).
			if int(sd) < len(xd) {
				xd = xd[:sd]
			}

			if roundUp {
				// Rounding up may mean the previous digit has to be rounded
				// up and so on.
				for sd > 0 {
					sd--
					xd[sd]++
					if xd[sd] <= int32(base-1) {
						break
					}
					xd[sd] = 0
					if sd == 0 {
						e++
						xd = append([]int32{1}, xd...)
						break
					}
				}
			}

			// Determine trailing zeros.
			ln = int64(len(xd))
			for xd[ln-1] == 0 {
				ln--
			}

			// E.g. [4, 11, 15] becomes 4bf.
			var b strings.Builder
			for i := int64(0); i < ln; i++ {
				b.WriteByte(numerals[xd[i]])
			}
			str = b.String()

			// Add binary exponent suffix?
			if isExp {
				if ln > 1 {
					if baseOut == 16 || baseOut == 8 {
						var ii int64 = 4
						if baseOut == 8 {
							ii = 3
						}
						for ln--; ln%ii != 0; ln++ {
							str += "0"
						}
						xd = convertBase(str, base, baseOut)
						ln = int64(len(xd))
						for xd[ln-1] == 0 {
							ln--
						}
						b.Reset()
						b.WriteString("1.")
						for i := int64(1); i < ln; i++ {
							b.WriteByte(numerals[xd[i]])
						}
						str = b.String()
					} else {
						str = str[:1] + "." + str[1:]
					}
				}
				if e < 0 {
					str += "p" + strconv.FormatInt(e, 10)
				} else {
					str += "p+" + strconv.FormatInt(e, 10)
				}
			} else if e < 0 {
				for ; e < -1; e++ {
					str = "0" + str
				}
				str = "0." + str
			} else {
				e++
				if e > ln {
					for e -= ln; e > 0; e-- {
						str += "0"
					}
				} else if e < ln {
					str = str[:e] + "." + str[e:]
				}
			}
		}

		switch baseOut {
		case 16:
			str = "0x" + str
		case 2:
			str = "0b" + str
		case 8:
			str = "0o" + str
		}
	}

	if x.s < 0 {
		return "-" + str
	}
	return str
}

// ToFraction returns x as a simple fraction with integer numerator and
// denominator, each a new Decimal. The denominator will be positive and
// at most maxD (if omitted, the lowest denominator representing x exactly).
func (x *Decimal) ToFraction(maxD ...any) []*Decimal {
	c := x.c
	xd := x.d
	if xd == nil {
		return []*Decimal{c.New(x)}
	}

	n1 := c.New(1)
	d0 := c.New(1)
	d1 := c.New(0)
	n0 := c.New(0)

	d := c.New(d1)
	e := getPrecision(xd) - x.e - 1
	d.e = e
	k := e % logBase
	if k < 0 {
		k += logBase
	}
	d.d[0] = int32(pow10(k))

	var maxDv *Decimal
	if len(maxD) == 0 || maxD[0] == nil {
		// d is 10**e, the minimum max-denominator needed.
		if e > 0 {
			maxDv = d
		} else {
			maxDv = n1
		}
	} else {
		n := c.New(maxD[0])
		if !n.IsInt() || n.Lt(n1) {
			panic(errors.New(invalidArgument + n.String()))
		}
		if n.Gt(d) {
			if e > 0 {
				maxDv = d
			} else {
				maxDv = n1
			}
		} else {
			maxDv = n
		}
	}

	x.c.external = false
	n := c.New(digitsToString(xd))
	pr := c.Precision
	c.Precision = xdLen2(xd)
	e = c.Precision

	var q, d2 *Decimal
	for {
		q = divide(n, d, 0, 1, true, 0)
		d2 = d0.Plus(q.Times(d1))
		if d2.Cmp(maxDv) == 1 {
			break
		}
		d0 = d1
		d1 = d2
		d2 = n1
		n1 = n0.Plus(q.Times(d2))
		n0 = d2
		d2 = d
		d = n.Minus(q.Times(d2))
		n = d2
	}

	d2 = divide(maxDv.Minus(d0), d1, 0, 1, true, 0)
	n0 = n0.Plus(d2.Times(n1))
	d0 = d0.Plus(d2.Times(d1))
	n0.s = x.s
	n1.s = x.s

	// Determine which fraction is closer to x, n0/d0 or n1/d1?
	var r []*Decimal
	if divide(n1, d1, e, 1, false, 0).Minus(x).Abs().Cmp(divide(n0, d0, e, 1, false, 0).Minus(x).Abs()) < 1 {
		r = []*Decimal{n1, d1}
	} else {
		r = []*Decimal{n0, d0}
	}

	c.Precision = pr
	x.c.external = true

	return r
}

func xdLen2(xd []int32) int64 { return int64(len(xd)) * logBase * 2 }

// ToNearest returns a new Decimal whose value is the nearest multiple of y
// in the direction of rounding mode rm (or the constructor's rounding mode
// if omitted).
func (x *Decimal) ToNearest(y any, rm ...int64) *Decimal {
	c := x.c
	xx := c.New(x)

	var yv *Decimal
	var r int64
	if y == nil {
		// If x is not finite, return x.
		if xx.d == nil {
			return xx
		}
		yv = c.New(1)
		r = c.Rounding
	} else {
		yv = c.New(y)
		if len(rm) > 0 {
			r = rm[0]
			checkInt32(r, 0, 8)
		} else {
			r = c.Rounding
		}
		// If x is not finite, return x if y is not NaN, else NaN.
		if xx.d == nil {
			if yv.s != 0 {
				return xx
			}
			return yv
		}
		// If y is not finite, return Infinity with the sign of x if y is
		// Infinity, else NaN.
		if yv.d == nil {
			if yv.s != 0 {
				yv.s = xx.s
			}
			return yv
		}
	}

	// If y is not zero, calculate the nearest multiple of y to x.
	if yv.d[0] != 0 {
		x.c.external = false
		xx = divide(xx, yv, 0, r, true, 0).Times(yv)
		x.c.external = true
		finalise(xx, sdNull, 0, false)
	} else {
		// If y is zero, return zero with the sign of x.
		yv.s = xx.s
		xx = yv
	}

	return xx
}
