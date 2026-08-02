package decimal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	prNull = int64(math.MinInt64)     // sentinel: use constructor precision
	sdNull = int64(math.MinInt64 + 1) // sentinel: no rounding requested
)

func pow10(n int64) int64 {
	r := int64(1)
	for ; n > 0; n-- {
		r *= 10
	}
	return r
}

// divPow10 returns w / 10^exp, mirroring JS's `w / Math.pow(10, exp)`.
// 10^exp overflows int64 for exp > 18 (and Math.pow gives Infinity for
// exp > 308); in both cases the mathematically-correct truncated result is
// 0, since w is a base-1e7 word (< 10^7).
func divPow10(w, exp int64) int64 {
	if exp > 18 {
		return 0
	}
	return w / pow10(exp)
}

func ceilDiv(a, b int64) int64 {
	if a <= 0 {
		return a / b
	}
	return (a + b - 1) / b
}

func itoa(w int32) string { return strconv.Itoa(int(w)) }

func digitsToString(d []int32) string {
	indexOfLastWord := len(d) - 1
	str := ""
	w := int64(d[0])

	if indexOfLastWord > 0 {
		str += strconv.FormatInt(w, 10)
		for i := 1; i < indexOfLastWord; i++ {
			ws := itoa(d[i])
			if k := int(logBase) - len(ws); k != 0 {
				str += strings.Repeat("0", k)
			}
			str += ws
		}
		w = int64(d[indexOfLastWord])
		ws := strconv.FormatInt(w, 10)
		if k := int(logBase) - len(ws); k != 0 {
			str += strings.Repeat("0", k)
		}
	} else if w == 0 {
		return "0"
	}

	// Remove trailing zeros of last word.
	for w%10 == 0 {
		w /= 10
	}

	return str + strconv.FormatInt(w, 10)
}

// Check 5 rounding digits if repeating is null (-1), 4 otherwise.
func checkRoundingDigits(d []int32, i, rm, repeating int64) bool {
	// Get the length of the first word of the array d.
	k := int64(d[0])
	for k >= 10 {
		k /= 10
		i--
	}

	var di int64
	i--
	if i < 0 {
		i += logBase
		di = 0
	} else {
		di = (i + 1 + logBase - 1) / logBase
		i %= logBase
	}

	// i is the index (0 - 6) of the rounding digit within the word.
	k = pow10(logBase - i)
	rd := int64(word(d, int(di))) % k

	var r bool
	if repeating < 0 {
		if i < 3 {
			if i == 0 {
				rd /= 100
			} else if i == 1 {
				rd /= 10
			}
			r = rm < 4 && rd == 99999 || rm > 3 && rd == 49999 || rd == 50000 || rd == 0
		} else {
			r = (rm < 4 && rd+1 == k || rm > 3 && rd+1 == k/2) &&
				(int64(word(d, int(di)+1))/k/100) == pow10(i-2)-1 ||
				(rd == k/2 || rd == 0) && (int64(word(d, int(di)+1))/k/100) == 0
		}
	} else {
		if i < 4 {
			if i == 0 {
				rd /= 1000
			} else if i == 1 {
				rd /= 100
			} else if i == 2 {
				rd /= 10
			}
			r = (repeating != 0 || rm < 4) && rd == 9999 || repeating == 0 && rm > 3 && rd == 4999
		} else {
			r = ((repeating != 0 || rm < 4) && rd+1 == k ||
				(repeating == 0 && rm > 3) && rd+1 == k/2) &&
				(int64(word(d, int(di)+1))/k/1000) == pow10(i-3)-1
		}
	}

	return r
}

// Convert string of baseIn to an array of numbers of baseOut.
func convertBase(str string, baseIn, baseOut int64) []int32 {
	arr := []int64{0}

	for i := 0; i < len(str); i++ {
		for j := range arr {
			arr[j] *= baseIn
		}
		arr[0] += int64(strings.IndexByte(numerals, str[i]))
		for j := 0; j < len(arr); j++ {
			if arr[j] > baseOut-1 {
				if j+1 == len(arr) {
					arr = append(arr, 0)
				}
				arr[j+1] += arr[j] / baseOut
				arr[j] %= baseOut
			}
		}
	}

	// Reverse.
	out := make([]int32, len(arr))
	for i, v := range arr {
		out[len(arr)-1-i] = int32(v)
	}
	return out
}

// Round x to sd significant digits using rounding mode rm.
// Check for over/under-flow. If sd == sdNull, no rounding is performed.
func finalise(x *Decimal, sd, rm int64, isTruncated bool) *Decimal {
	c := x.c
	xd := x.d

	if sd != sdNull {
		// Infinity/NaN.
		if xd == nil {
			return x
		}

		var j, w, xdi, rd, kPow, digits int64

		// Get the length of the first word of the digits array xd.
		digits = 1
		k := int64(xd[0])
		for k >= 10 {
			k /= 10
			digits++
		}
		i := sd - digits

		// Is the rounding digit in the first word of xd?
		if i < 0 {
			i += logBase
			j = sd
			xdi = 0
			w = int64(xd[0])

			// Get the rounding digit at index j of w.
			rd = divPow10(w, digits-j-1) % 10
		} else {
			xdi = (i + 1 + logBase - 1) / logBase
			k = int64(len(xd))
			if xdi >= k {
				if isTruncated {
					// Needed by naturalExponential, naturalLogarithm and squareRoot.
					for ; k <= xdi; k++ {
						xd = append(xd, 0)
					}
					x.d = xd
					w, rd, digits = 0, 0, 1
					i %= logBase
					j = i - logBase + 1
				} else {
					goto out
				}
			} else {
				k = int64(xd[xdi])
				w = k
				// Get the number of digits of w.
				digits = 1
				for k >= 10 {
					k /= 10
					digits++
				}
				// Get the index of rd within w.
				i %= logBase
				// Get the index of rd within w, adjusted for leading zeros.
				j = i - logBase + digits
				if j < 0 {
					rd = 0
				} else {
					rd = divPow10(w, digits-j-1) % 10
				}
			}
		}

		// Are there any non-zero digits after the rounding digit?
		isTruncated = isTruncated || sd < 0 ||
			int64(len(xd)) > xdi+1 ||
			(j < 0 && w != 0) || (j >= 0 && w%pow10(digits-j-1) != 0)

		var roundUp bool
		if rm < 4 {
			roundUp = (rd != 0 || isTruncated) && (rm == 0 || rm == ifelse64(x.s < 0, 3, 2))
		} else {
			var leftDigit int64
			if i > 0 {
				if j > 0 {
					leftDigit = w / pow10(digits-j)
				} else {
					leftDigit = 0
				}
			} else {
				leftDigit = int64(word(xd, int(xdi)-1))
			}
			roundUp = rd > 5 || rd == 5 && (rm == 4 || isTruncated ||
				rm == 6 && (leftDigit%10)&1 == 1 ||
				rm == ifelse64(x.s < 0, 8, 7))
		}

		if sd < 1 || xd[0] == 0 {
			xd = xd[:0]
			if roundUp {
				// Convert sd to decimal places.
				sd -= x.e + 1
				// 1, 0.1, 0.01, 0.001, 0.0001 etc.
				xd = append(xd, int32(pow10((logBase-sd%logBase)%logBase)))
				if sd == 0 {
					x.e = 0
				} else {
					x.e = -sd
				}
			} else {
				// Zero.
				xd = append(xd, 0)
				x.e = 0
			}
			x.d = xd
			return x
		}

		// Remove excess digits.
		if i == 0 {
			xd = xd[:xdi]
			kPow = 1
			xdi--
		} else {
			xd = xd[:xdi+1]
			kPow = pow10(logBase - i)
			// j > 0 means i > number of leading zeros of w.
			if j > 0 {
				xd[xdi] = int32((w / pow10(digits-j) % pow10(j)) * kPow)
			} else {
				xd[xdi] = 0
			}
		}

		if roundUp {
			for {
				// Is the digit to be rounded up in the first word of xd?
				if xdi == 0 {
					// i2 will be the length of xd[0] before kPow is added.
					i2 := int64(1)
					for j2 := int64(xd[0]); j2 >= 10; j2 /= 10 {
						i2++
					}
					nw := int64(xd[0]) + kPow
					xd[0] = int32(nw)
					k3 := int64(1)
					for ; nw >= 10; nw /= 10 {
						k3++
					}
					// If i2 != k3 the length has increased.
					if i2 != k3 {
						x.e++
						if xd[0] == int32(base) {
							xd[0] = 1
						}
					}
					break
				}
				xd[xdi] += int32(kPow)
				if xd[xdi] != int32(base) {
					break
				}
				xd[xdi] = 0
				xdi--
				kPow = 1
			}
		}

		// Remove trailing zeros.
		for len(xd) > 0 && xd[len(xd)-1] == 0 {
			xd = xd[:len(xd)-1]
		}
		x.d = xd
	}

out:
	if external {
		// Overflow?
		if x.e > c.MaxE {
			// Infinity.
			x.d = nil
			x.e = 0
		} else if x.e < c.MinE {
			// Zero.
			x.e = 0
			x.d = []int32{0}
		}
	}

	return x
}

func ifelse64(cond bool, a, b int64) int64 {
	if cond {
		return a
	}
	return b
}

// Calculate the base 10 exponent from the base 1e7 exponent.
func getBase10Exponent(digits []int32, e int64) int64 {
	w := int32(digits[0])
	e *= logBase
	for ; w >= 10; w /= 10 {
		e++
	}
	return e
}

func getPrecision(digits []int32) int64 {
	w := len(digits) - 1
	ln := int64(w)*logBase + 1
	wv := int64(digits[w])
	if wv != 0 {
		// Subtract the number of trailing zeros of the last word.
		for wv%10 == 0 {
			wv /= 10
			ln--
		}
		// Add the number of digits of the first word.
		for wv = int64(digits[0]); wv >= 10; wv /= 10 {
			ln++
		}
	}
	return ln
}

func zeroString(k int64) string {
	if k <= 0 {
		return ""
	}
	return strings.Repeat("0", int(k))
}

// intPow returns x^n for non-negative integer n, using exponentiation by
// squaring. Called by Pow and parseOther.
func intPow(c *Constructor, x *Decimal, n, pr int64) *Decimal {
	isTruncated := false
	r := c.New(1)

	// Maximum digits array length; leaves [28, 34] guard digits.
	k := int(math.Ceil(float64(pr)/float64(logBase) + 4))

	external = false

	for {
		if n%2 != 0 {
			r = r.Times(x)
			if len(r.d) > k {
				r.d = r.d[:k]
				isTruncated = true
			}
		}

		n /= 2
		if n == 0 {
			// To ensure correct rounding when r.d is truncated, increment the
			// last word if it is zero.
			nn := len(r.d) - 1
			if isTruncated && r.d[nn] == 0 {
				r.d[nn]++
			}
			break
		}

		x = x.Times(x)
		if len(x.d) > k {
			x.d = x.d[:k]
		}
	}

	external = true

	return r
}

func isOdd(x *Decimal) bool { return x.d[len(x.d)-1]&1 == 1 }

// --- Division ---

// Assumes non-zero x and k, and hence non-zero result.
func multiplyInteger(x []int32, k, baseI int64) []int32 {
	out := make([]int32, len(x))
	copy(out, x)
	carry := int64(0)
	for i := len(out) - 1; i >= 0; i-- {
		temp := int64(out[i])*k + carry
		out[i] = int32(temp % baseI)
		carry = temp / baseI
	}
	if carry != 0 {
		out = append([]int32{int32(carry)}, out...)
	}
	return out
}

func compareDigits(a, b []int32, aL, bL int) int {
	if aL != bL {
		if aL > bL {
			return 1
		}
		return -1
	}
	for i := 0; i < aL; i++ {
		if a[i] != b[i] {
			if a[i] > b[i] {
				return 1
			}
			return -1
		}
	}
	return 0
}

func subtractDigits(a, b []int32, aL int, baseI int64) []int32 {
	i := int64(0)
	// Subtract b from a.
	for ; aL > 0; aL-- {
		av := int64(a[aL-1]) - i
		bv := int64(word(b, aL-1))
		if av < bv {
			i = 1
		} else {
			i = 0
		}
		a[aL-1] = int32(i*baseI + av - bv)
	}
	// Remove leading zeros.
	for len(a) > 1 && a[0] == 0 {
		a = a[1:]
	}
	return a
}

// divide performs division in base 1e7 (baseI == 0) or in the specified
// base (used for base conversion, where logBase is 1).
// pr == prNull means use the constructor's precision and rounding.
// dp selects integer-division precision handling (used by DivToInt and Mod).
func divide(x, y *Decimal, pr, rm int64, dp bool, baseI int64) *Decimal {
	c := x.c
	var sign int8 = 1
	if x.s != y.s {
		sign = -1
	}
	xd, yd := x.d, y.d

	// Either NaN, Infinity or 0?
	if xd == nil || xd[0] == 0 || yd == nil || yd[0] == 0 {
		nan := x.s == 0 || y.s == 0
		if !nan {
			if xd != nil {
				nan = yd != nil && xd[0] == yd[0]
			} else {
				nan = yd == nil
			}
		}
		if nan {
			// Return NaN if either NaN, or both Infinity or 0.
			return c.newNaN()
		}
		// Return ±0 if x is 0 or y is ±Infinity, or return ±Infinity as y is 0.
		if (xd != nil && xd[0] == 0) || yd == nil {
			return c.newZero(sign)
		}
		return c.newInf(sign)
	}

	var logBaseI, e int64
	if baseI != 0 {
		logBaseI = 1
		e = x.e - y.e
	} else {
		baseI = base
		logBaseI = logBase
		e = floorDiv(x.e, logBase) - floorDiv(y.e, logBase)
	}

	yL, xL := len(yd), len(xd)
	q := c.newSign(sign)
	qd := q.d

	// Result exponent may be one less than e.
	i := 0
	for {
		ydef := i < yL
		if !(ydef && yd[i] == word(xd, i)) {
			break
		}
		i++
	}
	if i < yL && yd[i] > word(xd, i) {
		e--
	}

	var sd int64
	if pr == prNull {
		sd = c.Precision
		pr = sd
		rm = c.Rounding
	} else if dp {
		sd = pr + (x.e - y.e) + 1
	} else {
		sd = pr
	}

	var more bool
	remValid := false

	if sd < 0 {
		qd = append(qd, 1)
		more = true
	} else {
		sd0 := sd/logBaseI + 2
		i = 0

		// divisor < 1e7
		if yL == 1 {
			k := int64(0)
			yd0 := int64(yd[0])
			sd0++

			// k is the carry.
			for i < xL || k != 0 {
				if sd0 == 0 {
					break
				}
				sd0--
				t := k*baseI + int64(word(xd, i))
				qd = append(qd, int32(t/yd0))
				k = t % yd0
				i++
			}
			more = k != 0 || i < xL

			// divisor >= 1e7
		} else {
			// Normalise xd and yd so highest order digit of yd is >= base/2.
			k := baseI / (int64(yd[0]) + 1)

			if k > 1 {
				yd = multiplyInteger(yd, k, baseI)
				xd = multiplyInteger(xd, k, baseI)
				yL = len(yd)
				xL = len(xd)
			}

			xi := yL
			remLen := yL
			if xL < yL {
				remLen = xL
			}
			rem := append([]int32(nil), xd[:remLen]...)
			for len(rem) < yL {
				rem = append(rem, 0)
			}
			remL := len(rem)

			yz := append([]int32{0}, yd...)
			yd0 := int64(yd[0])
			if yd[1] >= int32(baseI/2) {
				yd0++
			}

			cmp := 0
			for {
				k = 0
				cmp = compareDigits(yd, rem, yL, remL)

				// If divisor < remainder.
				if cmp < 0 {
					// Calculate trial digit, k.
					rem0 := int64(rem[0])
					if yL != remL {
						rem0 = rem0*baseI + int64(word(rem, 1))
					}
					// k will be how many times the divisor goes into the
					// current remainder.
					k = rem0 / yd0

					var prod []int32
					var prodL int
					if k > 1 {
						if k >= baseI {
							k = baseI - 1
						}
						// product = divisor * trial digit.
						prod = multiplyInteger(yd, k, baseI)
						prodL = len(prod)
						remL = len(rem)
						// Compare product and remainder.
						cmp = compareDigits(prod, rem, prodL, remL)
						// product > remainder.
						if cmp == 1 {
							k--
							// Subtract divisor from product.
							if yL < prodL {
								prod = subtractDigits(prod, yz, prodL, baseI)
							} else {
								prod = subtractDigits(prod, yd, prodL, baseI)
							}
						}
					} else {
						// cmp is -1. If k is 0, change cmp to 1 to avoid
						// comparing yd and rem again below.
						if k == 0 {
							cmp = 1
							k = 1
						}
						prod = append([]int32(nil), yd...)
					}

					prodL = len(prod)
					if prodL < remL {
						prod = append([]int32{0}, prod...)
					}

					// Subtract product from remainder.
					rem = subtractDigits(rem, prod, remL, baseI)

					// If product was < previous remainder.
					if cmp == -1 {
						remL = len(rem)
						cmp = compareDigits(yd, rem, yL, remL)
						if cmp < 1 {
							k++
							if yL < remL {
								rem = subtractDigits(rem, yz, remL, baseI)
							} else {
								rem = subtractDigits(rem, yd, remL, baseI)
							}
						}
					}

					remL = len(rem)
				} else if cmp == 0 {
					k++
					rem = []int32{0}
				} // if cmp == 1, k will be 0

				// Add the next digit, k, to the result array.
				qd = append(qd, int32(k))
				i++

				// Update the remainder.
				remL = len(rem)
				if cmp != 0 && rem[0] != 0 {
					rem = append(rem, word(xd, xi))
					remL = len(rem)
					remValid = true
				} else {
					if xi < xL {
						rem = []int32{xd[xi]}
						remValid = true
					} else {
						rem = []int32{0}
						remValid = false
					}
					remL = 1
				}

				cond := xi < xL || remValid
				xi++
				if !cond {
					break
				}
				if sd0 == 0 {
					break
				}
				sd0--
			}
			more = remValid
		}

		// Leading zero?
		if len(qd) > 0 && qd[0] == 0 {
			qd = qd[1:]
		}
	}

	q.d = qd

	// logBaseI is 1 when divide is being used for base conversion.
	if logBaseI == 1 {
		q.e = e
		inexact = more
	} else {
		// To calculate q.e, first get the number of digits of qd[0].
		i2 := int64(1)
		for k2 := int64(qd[0]); k2 >= 10; k2 /= 10 {
			i2++
		}
		q.e = i2 + e*logBaseI - 1

		if dp {
			finalise(q, pr+q.e+1, rm, more)
		} else {
			finalise(q, pr, rm, more)
		}
	}

	return q
}

// --- Instance methods: comparisons and basic arithmetic ---

// Cmp returns
//
//	1    if the value of x is greater than the value of y,
//	-1   if the value of x is less than the value of y,
//	0    if they have the same value,
//	NaN  if the value of either is NaN.
func (x *Decimal) Cmp(y any) float64 {
	yv := x.c.New(y)
	xd, yd := x.d, yv.d
	xs, ys := x.s, yv.s

	// Either NaN or ±Infinity?
	if xd == nil || yd == nil {
		if xs == 0 || ys == 0 {
			return math.NaN()
		}
		if xs != ys {
			return float64(xs)
		}
		if xd == nil && yd == nil {
			return 0
		}
		if (xd == nil) != (xs < 0) {
			return 1
		}
		return -1
	}

	// Either zero?
	if xd[0] == 0 || yd[0] == 0 {
		if xd[0] != 0 {
			return float64(xs)
		}
		if yd[0] != 0 {
			return float64(-ys)
		}
		return 0
	}

	// Signs differ?
	if xs != ys {
		return float64(xs)
	}

	// Compare exponents.
	if x.e != yv.e {
		if (x.e > yv.e) != (xs < 0) {
			return 1
		}
		return -1
	}

	// Compare digit by digit.
	xdL, ydL := len(xd), len(yd)
	n := xdL
	if ydL < n {
		n = ydL
	}
	for i := 0; i < n; i++ {
		if xd[i] != yd[i] {
			if (xd[i] > yd[i]) != (xs < 0) {
				return 1
			}
			return -1
		}
	}

	// Compare lengths.
	if xdL == ydL {
		return 0
	}
	if (xdL > ydL) != (xs < 0) {
		return 1
	}
	return -1
}

// Eq returns true if x == y.
func (x *Decimal) Eq(y any) bool { return x.Cmp(y) == 0 }

// Gt returns true if x > y.
func (x *Decimal) Gt(y any) bool { return x.Cmp(y) > 0 }

// Gte returns true if x >= y.
func (x *Decimal) Gte(y any) bool {
	k := x.Cmp(y)
	return k == 1 || k == 0
}

// Lt returns true if x < y.
func (x *Decimal) Lt(y any) bool { return x.Cmp(y) < 0 }

// Lte returns true if x <= y.
func (x *Decimal) Lte(y any) bool { return x.Cmp(y) < 1 }

// IsFinite returns true if x is a finite number.
func (x *Decimal) IsFinite() bool { return x.d != nil }

// IsInt returns true if x is a finite integer.
func (x *Decimal) IsInt() bool {
	return x.d != nil && floorDiv(x.e, logBase) > int64(len(x.d))-2
}

// IsNaN returns true if x is NaN.
func (x *Decimal) IsNaN() bool { return x.s == 0 }

// IsNeg returns true if x is negative.
func (x *Decimal) IsNeg() bool { return x.s < 0 }

// IsPos returns true if x is positive.
func (x *Decimal) IsPos() bool { return x.s > 0 }

// IsZero returns true if x is 0 or -0.
func (x *Decimal) IsZero() bool { return x.d != nil && x.d[0] == 0 }

// Dp returns the number of decimal places of x, or NaN if x is not finite.
func (x *Decimal) Dp() float64 {
	d := x.d
	if d == nil {
		return math.NaN()
	}
	w := len(d) - 1
	n := (int64(w) - floorDiv(x.e, logBase)) * logBase
	wv := int64(d[w])
	if wv != 0 {
		for wv%10 == 0 {
			wv /= 10
			n--
		}
	}
	if n < 0 {
		n = 0
	}
	return float64(n)
}

// Sd returns the number of significant digits of x, or NaN if x is not
// finite. If z is true, integer-part trailing zeros are counted.
func (x *Decimal) Sd(z ...bool) float64 {
	if x.d == nil {
		return math.NaN()
	}
	k := x.sdInt()
	if len(z) > 0 && z[0] && x.e+1 > k {
		k = x.e + 1
	}
	return float64(k)
}

func (x *Decimal) sdInt() int64 {
	if x.d == nil {
		return math.MinInt64
	}
	return getPrecision(x.d)
}

// Abs returns a new Decimal whose value is |x|.
func (x *Decimal) Abs() *Decimal {
	xx := x.c.New(x)
	if xx.s < 0 {
		xx.s = 1
	}
	return finalise(xx, sdNull, 0, false)
}

// Ceil returns a new Decimal whose value is x rounded to a whole number in
// the direction of positive Infinity.
func (x *Decimal) Ceil() *Decimal {
	xx := x.c.New(x)
	return finalise(xx, xx.e+1, RoundCeil, false)
}

// Floor returns a new Decimal whose value is x rounded to a whole number in
// the direction of negative Infinity.
func (x *Decimal) Floor() *Decimal {
	xx := x.c.New(x)
	return finalise(xx, xx.e+1, RoundFloor, false)
}

// Round returns a new Decimal whose value is x rounded to a whole number
// using the constructor's rounding mode.
func (x *Decimal) Round() *Decimal {
	xx := x.c.New(x)
	return finalise(xx, xx.e+1, x.c.Rounding, false)
}

// Trunc returns a new Decimal whose value is x truncated to a whole number.
func (x *Decimal) Trunc() *Decimal {
	xx := x.c.New(x)
	return finalise(xx, xx.e+1, RoundDown, false)
}

// Neg returns a new Decimal whose value is -x.
func (x *Decimal) Neg() *Decimal {
	xx := x.c.New(x)
	xx.s = -xx.s
	return finalise(xx, sdNull, 0, false)
}

// Plus returns x + y rounded to the constructor's precision.
func (x *Decimal) Plus(y any) *Decimal {
	c := x.c
	yv := c.New(y)

	// If either is not finite...
	if x.d == nil || yv.d == nil {
		// Return NaN if either is NaN.
		if x.s == 0 || yv.s == 0 {
			return c.newNaN()
		}

		// Return x if y is finite and x is ±Infinity.
		// Return x if both are ±Infinity with the same sign.
		// Return NaN if both are ±Infinity with different signs.
		// Return y if x is finite and y is ±Infinity.
		if x.d == nil {
			if yv.d != nil || x.s == yv.s {
				return c.New(x)
			}
			return c.newNaN()
		}

		return yv
	}

	// If signs differ...
	if x.s != yv.s {
		yv.s = -yv.s
		return x.Minus(yv)
	}

	xd, yd := x.d, yv.d
	pr, rm := c.Precision, c.Rounding

	// If either is zero...
	if xd[0] == 0 || yd[0] == 0 {
		// Return x if y is zero.
		if yd[0] == 0 {
			yv = c.New(x)
		}
		if external {
			return finalise(yv, pr, rm, false)
		}
		return yv
	}

	// x and y are finite, non-zero numbers with the same sign.

	// Calculate base 1e7 exponents.
	k := floorDiv(x.e, logBase)
	e := floorDiv(yv.e, logBase)

	xdc := append([]int32(nil), xd...)
	i := k - e

	// If base 1e7 exponents differ...
	if i != 0 {
		var d []int32
		var ln int64
		wasXd := false
		if i < 0 {
			d = xdc
			wasXd = true
			i = -i
			ln = int64(len(yd))
		} else {
			d = yd
			e = k
			ln = int64(len(xdc))
		}

		// Limit number of zeros prepended to max(ceil(pr / LOG_BASE), len) + 1.
		if kk := ceilDiv(pr, logBase); kk > ln {
			ln = kk + 1
		} else {
			ln++
		}

		if i > ln {
			i = ln
			d = d[:1]
		}

		// Prepend zeros to equalise exponents.
		nd := make([]int32, 0, int(i)+len(d))
		for ; i > 0; i-- {
			nd = append(nd, 0)
		}
		d = append(nd, d...)
		if wasXd {
			xdc = d
		} else {
			yd = d
		}
	}

	ln := len(xdc)
	i = int64(len(yd))

	// If yd is longer than xd, swap so xd points to the longer array.
	if int64(ln)-i < 0 {
		i = int64(ln)
		xdc, yd = yd, xdc
	}

	// Only start adding at len(yd)-1 as the further digits of xd can be
	// left as they are.
	var carry int64
	for ; i > 0; i-- {
		t := int64(xdc[i-1]) + int64(yd[i-1]) + carry
		xdc[i-1] = int32(t % base)
		carry = t / base
	}

	if carry != 0 {
		xdc = append([]int32{int32(carry)}, xdc...)
		e++
	}

	// Remove trailing zeros.
	for len(xdc) > 0 && xdc[len(xdc)-1] == 0 {
		xdc = xdc[:len(xdc)-1]
	}

	yv.d = xdc
	yv.e = getBase10Exponent(xdc, e)

	if external {
		return finalise(yv, c.Precision, c.Rounding, false)
	}
	return yv
}

// Add is an alias of Plus.
func (x *Decimal) Add(y any) *Decimal { return x.Plus(y) }

// Minus returns x - y rounded to the constructor's precision.
func (x *Decimal) Minus(y any) *Decimal {
	c := x.c
	yv := c.New(y)

	// If either is not finite...
	if x.d == nil || yv.d == nil {
		// Return NaN if either is NaN.
		if x.s == 0 || yv.s == 0 {
			return c.newNaN()
		}
		if x.d != nil {
			// Return y negated if x is finite and y is ±Infinity.
			yv.s = -yv.s
			return yv
		}
		// Return x if y is finite and x is ±Infinity.
		// Return x if both are ±Infinity with different signs.
		// Return NaN if both are ±Infinity with the same sign.
		if yv.d != nil || x.s != yv.s {
			return c.New(x)
		}
		return c.newNaN()
	}

	// If signs differ...
	if x.s != yv.s {
		yv.s = -yv.s
		return x.Plus(yv)
	}

	xd, yd := x.d, yv.d
	pr, rm := c.Precision, c.Rounding

	// If either is zero...
	if xd[0] == 0 || yd[0] == 0 {
		if yd[0] != 0 {
			// Return y negated if x is zero and y is non-zero.
			yv.s = -yv.s
		} else if xd[0] != 0 {
			// Return x if y is zero and x is non-zero.
			yv = c.New(x)
		} else {
			// Return zero if both are zero.
			// From IEEE 754 (2008) 6.3: 0 - 0 = -0 - -0 = -0 when rounding
			// to -Infinity.
			return c.newZero(ifelse8(rm == RoundFloor, -1, 1))
		}
		if external {
			return finalise(yv, pr, rm, false)
		}
		return yv
	}

	// x and y are finite, non-zero numbers with the same sign.

	// Calculate base 1e7 exponents.
	e := floorDiv(yv.e, logBase)
	xe := floorDiv(x.e, logBase)

	xdc := append([]int32(nil), xd...)
	k := xe - e

	var xLTy bool
	if k != 0 {
		var d []int32
		var ln int64
		xLTy = k < 0
		wasXd := false
		if xLTy {
			d = xdc
			wasXd = true
			k = -k
			ln = int64(len(yd))
		} else {
			d = yd
			e = xe
			ln = int64(len(xdc))
		}

		// Limit the number of zeros prepended.
		i := ceilDiv(pr, logBase)
		if i < ln {
			i = ln
		}
		i += 2
		if k > i {
			k = i
			d = d[:1]
		}

		// Prepend zeros to equalise exponents.
		nd := make([]int32, 0, int(k)+len(d))
		for jj := int64(0); jj < k; jj++ {
			nd = append(nd, 0)
		}
		d = append(nd, d...)
		if wasXd {
			xdc = d
		} else {
			yd = d
		}
	} else {
		// Check digits to determine which is the bigger number.
		xLTy = len(xdc) < len(yd)
		n := len(xdc)
		if xLTy {
			n = len(xdc)
		} else if len(yd) < n {
			n = len(yd)
		}
		for i := 0; i < n; i++ {
			if xdc[i] != yd[i] {
				xLTy = xdc[i] < yd[i]
				break
			}
		}
		k = 0
	}

	if xLTy {
		xdc, yd = yd, xdc
		yv.s = -yv.s
	}

	// Append zeros to xd if shorter. Don't add zeros to yd if shorter as
	// subtraction only needs to start at yd length.
	for len(xdc) < len(yd) {
		xdc = append(xdc, 0)
	}

	// Subtract yd from xd.
	for i := len(yd) - 1; i >= int(k); i-- {
		if xdc[i] < yd[i] {
			j := i
			for j > 0 {
				j--
				if xdc[j] != 0 {
					break
				}
				xdc[j] = int32(base - 1)
			}
			xdc[j]--
			xdc[i] += int32(base)
		}
		xdc[i] -= yd[i]
	}

	// Remove trailing zeros.
	for len(xdc) > 0 && xdc[len(xdc)-1] == 0 {
		xdc = xdc[:len(xdc)-1]
	}
	// Remove leading zeros and adjust exponent accordingly.
	for len(xdc) > 0 && xdc[0] == 0 {
		xdc = xdc[1:]
		e--
	}

	// Zero?
	if len(xdc) == 0 || xdc[0] == 0 {
		return c.newZero(ifelse8(rm == RoundFloor, -1, 1))
	}

	yv.d = xdc
	yv.e = getBase10Exponent(xdc, e)

	if external {
		return finalise(yv, pr, rm, false)
	}
	return yv
}

// Sub is an alias of Minus.
func (x *Decimal) Sub(y any) *Decimal { return x.Minus(y) }

func ifelse8(cond bool, a, b int8) int8 {
	if cond {
		return a
	}
	return b
}

// Times returns x * y rounded to the constructor's precision.
func (x *Decimal) Times(y any) *Decimal {
	c := x.c
	yv := c.New(y)
	xd, yd := x.d, yv.d

	yv.s *= x.s

	// If either is NaN, ±Infinity or ±0...
	if xd == nil || xd[0] == 0 || yd == nil || yd[0] == 0 {
		if yv.s == 0 || (xd != nil && xd[0] == 0 && yd == nil) || (yd != nil && yd[0] == 0 && xd == nil) {
			// Return NaN if either is NaN, or if one is ±0 and the other
			// is ±Infinity.
			return c.newNaN()
		}
		if xd == nil || yd == nil {
			// Return ±Infinity if either is ±Infinity.
			return c.newInf(yv.s)
		}
		// Return ±0 if either is ±0.
		return c.newZero(yv.s)
	}

	e := floorDiv(x.e, logBase) + floorDiv(yv.e, logBase)
	xdL, ydL := len(xd), len(yd)

	// Ensure xd points to the longer array.
	if xdL < ydL {
		xd, yd = yd, xd
		xdL, ydL = ydL, xdL
	}

	// Initialise the result array with zeros.
	r := make([]int32, xdL+ydL)

	// Multiply!
	var carry int64
	for i := ydL - 1; i >= 0; i-- {
		carry = 0
		k := xdL + i
		for ; k > i; k-- {
			t := int64(r[k]) + int64(yd[i])*int64(xd[k-i-1]) + carry
			r[k] = int32(t % base)
			carry = t / base
		}
		r[k] = int32((int64(r[k]) + carry) % base)
	}

	// Remove trailing zeros.
	for len(r) > 0 && r[len(r)-1] == 0 {
		r = r[:len(r)-1]
	}

	if carry != 0 {
		e++
	} else {
		r = r[1:]
	}

	yv.d = r
	yv.e = getBase10Exponent(r, e)

	if external {
		return finalise(yv, c.Precision, c.Rounding, false)
	}
	return yv
}

// Mul is an alias of Times.
func (x *Decimal) Mul(y any) *Decimal { return x.Times(y) }

// Div returns x / y rounded to the constructor's precision.
func (x *Decimal) Div(y any) *Decimal {
	return divide(x, x.c.New(y), prNull, 0, false, 0)
}

// DivToInt returns a new Decimal whose value is the integer part of x / y,
// rounded to the constructor's precision.
func (x *Decimal) DivToInt(y any) *Decimal {
	c := x.c
	return finalise(divide(x, c.New(y), 0, 1, true, 0), c.Precision, c.Rounding, false)
}

// Mod returns x modulo y, rounded to the constructor's precision. The
// result depends on the constructor's modulo mode.
func (x *Decimal) Mod(y any) *Decimal {
	c := x.c
	yv := c.New(y)

	// Return NaN if x is ±Infinity or NaN, or y is NaN or ±0.
	if x.d == nil || yv.s == 0 || (yv.d != nil && yv.d[0] == 0) {
		return c.newNaN()
	}

	// Return x if y is ±Infinity or x is ±0.
	if yv.d == nil || (x.d != nil && x.d[0] == 0) {
		return finalise(c.New(x), c.Precision, c.Rounding, false)
	}

	// Prevent rounding of intermediate calculations.
	external = false

	var q *Decimal
	if c.Modulo == Euclid {
		// Euclidian division: q = sign(y) * floor(x / abs(y))
		// result = x - q * y    where  0 <= result < abs(y)
		q = divide(x, yv.Abs(), 0, 3, true, 0)
		q.s *= yv.s
	} else {
		q = divide(x, yv, 0, c.Modulo, true, 0)
	}

	q = q.Times(yv)

	external = true

	return x.Minus(q)
}

// Clamp returns x clamped to the range delineated by min and max.
func (x *Decimal) Clamp(min, max any) *Decimal {
	c := x.c
	minv := c.New(min)
	maxv := c.New(max)
	if minv.s == 0 || maxv.s == 0 {
		return c.newNaN()
	}
	if minv.Gt(maxv) {
		panic(errors.New(invalidArgument + maxv.String()))
	}
	if x.Cmp(minv) < 0 {
		return minv
	}
	if x.Cmp(maxv) > 0 {
		return maxv
	}
	return c.New(x)
}
